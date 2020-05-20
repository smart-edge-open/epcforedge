// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

// ConsumerConns websocket connections where key is consumerID
type ConsumerConns map[string]*websocket.Conn

/* This function validates the consumeID of Origin Header in
websocket establishment request, checks whether known at AF */
func chkConsumerID(r *http.Request) bool {

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	consumerID := r.Header.Get("Origin")

	if len(consumerID) == 0 {
		log.Errf("Authentication Failure")
		return false
	}

	for _, v := range afCtx.appSessionsEv {
		if v.wsReq {
			if v.consumerID == consumerID {
				log.Errf("Authentication Success for consumer %s", consumerID)
				return true
			}
		}
	}
	log.Errf("Authentication Failed for consumer %s", consumerID)
	return false
}

/* The Upgrader with CheckOrigin defined for validating the Origin*/
var upgrader = websocket.Upgrader{
	CheckOrigin: chkConsumerID,
}

// Upgrade upgrades the connection to webSocket
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Infoln(err)
		return ws, err
	}
	return ws, nil
}

// createWsConn creates a websocket connection for a consumer
// to send Notifications
func createWsConn(w http.ResponseWriter, r *http.Request,
	afCtx *Context) (int, error) {

	var (
		err error
		ws  *websocket.Conn
	)

	//Check for Origin Header
	// if yes, proceed, else return 403
	origin := r.Header.Get("Origin")
	if len(origin) == 0 {
		return http.StatusForbidden, errors.New("Nil Origin")
	}

	// If previous connection exists for this consumer, close the connection
	err = removeWsConn(origin, afCtx)
	if err != nil {
		log.Errf("Error in closing connection %s", err.Error())
		return http.StatusInternalServerError,
			errors.New("Unable to close previous connection")
	}

	// Upgrade the new connection to websocket
	ws, err = Upgrade(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Store the connection for consumer
	afCtx.data.consumerConns[origin] = ws
	log.Infoln("Added consumer conn for consumerID ", origin)

	return 0, nil
}

/* This function checks whether the ws connection of a consumer can be closed
It is invoked when Delete Policy Auth is recceived. If consumerID is present
for any other appSession, then webscoket is not closed.*/
func chkRemoveWSConn(evInfo *EventInfo, appSessionID string,
	afCtx *Context) error {

	if !evInfo.wsReq {
		return nil
	}
	consumerID := evInfo.consumerID
	// if ConsumerID is present for any other appSession, don't delete
	for key, v := range afCtx.appSessionsEv {
		if v.consumerID == consumerID {
			if key == appSessionID {
				continue
			}
			return nil
		}
	}
	err := removeWsConn(consumerID, afCtx)
	return err

}

/* This function closes the websocket connection of a consumer by
sending a websocket Close message */
func removeWsConn(consumerID string, afCtx *Context) error {

	foundConn, connFound := afCtx.data.consumerConns[consumerID]
	if connFound {

		msgType := websocket.CloseMessage
		closeMessage := websocket.FormatCloseMessage(
			websocket.CloseServiceRestart,
			"Closing this connection")
		err := foundConn.WriteMessage(msgType, closeMessage)
		if err != nil {
			log.Info("Failed to send close message to connection")
			return err
		}
		err = foundConn.Close()
		if err != nil {
			return err
		}
		delete(afCtx.data.consumerConns, consumerID)
	}
	return nil

}

//GetNotifications is invoked when a consumer connects to the
//AF websocketURI to make a websocket connection */
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		errRspHeader(&w, "GET NOTIFICATIONS",
			"af-ctx retrieved from request is nil",
			http.StatusInternalServerError)
		return
	}

	statCode, err := createWsConn(w, r, afCtx)
	if err != nil {
		log.Errf("Error in WebSocket Connection Creation: %#v ", err)
		if statCode != 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(statCode)
		}
		return
	}

}

/* This function sends the Notification to consumer on websocket */
func sendUpPathOnWs(evInfo *EventInfo, ev NotificationUpPathChg,
	afCtx *Context) error {

	var (
		afEvent    Afnotification
		consumerID string
	)

	consumerID = evInfo.consumerID

	if len(consumerID) == 0 {
		return errors.New("ConsumerID nil")
	}

	conn := afCtx.data.consumerConns[consumerID]
	if conn == nil {
		return errors.New("Consumer Connection Not found")
	}

	afEvent.Event = UPPathChangeEvent
	payload, err := json.Marshal(ev)
	if err != nil {
		log.Err(err)
		return err
	}

	afEvent.Payload = payload
	err = conn.WriteJSON(afEvent)
	return err

}

/* This function builds the AF websocketURI to be shared with consumer */
func getWSNotificationURI(afCtx *Context) (wsURI string) {

	wsURI = "https//" + afCtx.cfg.SrvCfg.Hostname +
		afCtx.cfg.SrvCfg.NotifWebsocketPort + afCtx.cfg.PrefixNotifications
	return wsURI
}

/* This function updates the AppSessionContextRespData with AF websocketURI
in resp*/
func updateAppSessRspForWS(appSess *AppSessionContext,
	afCtx *Context) {

	var (
		ascRespDataWS AppSessionContextRespData
		wsURI         string
	)

	wsURI = getWSNotificationURI(afCtx)
	ascRespData := appSess.AscRespData

	if ascRespData == nil {
		ascRespDataWS.WebsocketURI = wsURI
		appSess.AscRespData = &ascRespDataWS
	} else {
		ascRespData.WebsocketURI = wsURI
	}
}
