// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ConsumerConnection Websocket Connection details of a consumer
type ConsumerConnection struct {
	connection *websocket.Conn
	// gorilla websocket doesn't allow concurrent writes
	// Mutex per conn for avoiding concurrent writes
	control sync.Mutex
}

// ConsumerConns websocket connections where key is consumerID
type ConsumerConns map[string]*ConsumerConnection

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
If any read error then it closes the connection. Connection clearing from
connections map happens as part of Write or when a new connection establishement
 occurs */
func readLoop(c *websocket.Conn, consumerID string) {
	log.Infoln("Started the read loop on websocket for consumer ", consumerID)
	for {
		if _, _, err := c.NextReader(); err != nil {
			log.Errln("Closed the websocket conn in readLoop for consumer ",
				consumerID)
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
		err  error
		ws   *websocket.Conn
		conn ConsumerConnection
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

	conn.connection = ws
	// Store the connection for consumer
	afCtx.data.consumerConns[origin] = &conn

	// Start the goroutine to check for read errors
	go readLoop(conn.connection, origin)
	log.Infoln("Added consumer conn for consumerID ", origin)
	return 0, nil
}

/* This function checks whether the ws connection of a consumer can be closed
It is invoked when Delete/Update Policy Auth is recceived. If consumerID is
present for any other appSession, then websocket is not closed.*/
func chkRemoveWSConn(evInfo *EventInfo, appSessionID string,
	afCtx *Context) error {

	if evInfo == nil || !evInfo.wsReq {
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
sending a websocket Close message, mutex use before sending
Close message*/
func removeWsConn(consumerID string, afCtx *Context) error {

	foundConn, connFound := afCtx.data.consumerConns[consumerID]
	if connFound {

		msgType := websocket.CloseMessage
		closeMessage := websocket.FormatCloseMessage(
			websocket.CloseServiceRestart,
			"Closing this connection")
		foundConn.control.Lock()
		err := foundConn.connection.WriteMessage(msgType, closeMessage)
		if err != nil {
			log.Info("Failed to send close message to connection")
			return err
		}
		foundConn.control.Unlock()
		log.Infoln("Clearing the connection for consumer", consumerID)
		err = foundConn.connection.Close()
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

/* This function sends the Notification to consumer on websocket
Use of mutex before writing to the connection*/
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

	conn.control.Lock()
	err := conn.connection.WriteJSON(afEvent)
	if err != nil {
		if err.Error() == "tls: use of closed connection" {
			log.Infoln("Deleting connection as it is already closed")
			delete(afCtx.data.consumerConns, consumerID)
		}
	}
	conn.control.Unlock()
	return err

}

/* This function builds the AF websocketURI to be shared with consumer */
func getWSNotificationURI(afCtx *Context) string {

	return ("wss//" + afCtx.cfg.SrvCfg.Hostname +
		afCtx.cfg.SrvCfg.NotifWebsocketPort + afCtx.cfg.NotifWebsocketURI)

}

/* This function updates the AppSessionContextRespData with AF websocketURI
in Resp*/
func updateAppSessRspForWS(appSess *AppSessionContext,
	afCtx *Context) {

	var (
		ascRespDataWS AppSessionContextRespData
		wsURI         string
	)

	// To fetch the AF websocketURI
	wsURI = getWSNotificationURI(afCtx)
	ascRespData := appSess.AscRespData

	/* If RspData from PCF is present(It is an optional parameter),
	update that else create an AppSessionContextRespData and send*/
	if ascRespData == nil {
		ascRespDataWS.WebsocketURI = wsURI
		appSess.AscRespData = &ascRespDataWS
	} else {
		ascRespData.WebsocketURI = wsURI
	}
}
