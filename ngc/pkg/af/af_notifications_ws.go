// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
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

// Upgrade upgrades the connection to webSocket, sends a success or failure
// response to client
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Infoln(err)
		return ws, err
	}
	return ws, nil
}

/* This goroutine reads and checks for errors(close by consumer)
If any read error then it closes the connection */
func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			_ = c.Close()
			break
		}
	}
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

	}

	// Upgrade the new connection to websocket and sends a success/failure
	// response code
	ws, err = Upgrade(w, r)
	if err != nil {
		return 0, err
	}

	// Store the connection for consumer
	afCtx.data.consumerConns[origin] = ws

	// Start the goroutine to check for read errors
	go readLoop(ws)
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
	// if ConsumerID is present for any other appSession, don't remove
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
func sendNotificationOnWs(consumerID string, afEvent interface{},
	afCtx *Context) error {

	if len(consumerID) == 0 {
		return errors.New("ConsumerID nil")
	}

	// Fetching the connection for the consumer
	conn := afCtx.data.consumerConns[consumerID]
	if conn == nil {
		return errors.New("Consumer Connection Not found")
	}

	err := conn.WriteJSON(afEvent)
	if err != nil {
		if err.Error() == "tls: use of closed connection" {
			log.Infoln("Removing connection as it is closed")
			delete(afCtx.data.consumerConns, consumerID)
		}
	}
	return err

}

/* This function builds the AF websocketURI to be shared with consumer */
func getWSNotificationURI(afCtx *Context) string {

	return ("https//" + afCtx.cfg.SrvCfg.Hostname +
		afCtx.cfg.SrvCfg.NotifWebsocketPort + afCtx.cfg.PrefixNotifications)

}

/* This function updates the AppSessionContextRespData with AF websocketURI
in resp*/
func updateAppSessRspForWS(appSess *AppSessionContext,
	afCtx *Context) {

	var (
		ascRespDataWS AppSessionContextRespData
		wsURI         string
	)

	// To fetch the AF websocketURI
	wsURI = getWSNotificationURI(afCtx)
	ascRespData := appSess.AscRespData

	/* If RspData is present, update that else create and send*/
	if ascRespData == nil {
		ascRespDataWS.WebsocketURI = wsURI
		appSess.AscRespData = &ascRespDataWS
	} else {
		ascRespData.WebsocketURI = wsURI
	}
}
