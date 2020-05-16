// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"errors"
	"net/url"
	"strings"
)

// NotificationURI to store the notification URI of consumer
type NotificationURI string

// EventInfo stores information for all Events of an appSessionID
type EventInfo struct {
	// This is true when webSocket delivery is used for all events
	wsReq bool
	// Unique identification of consumer
	consumerID string
	// UP_PATH_CHANGE event consumer information, key is notifCorrelID
	upPathEv map[string]NotificationURI
}

// AppSessEv This stores the Event Information for appSessions
type AppSessEv map[string]*EventInfo

func initNotify(afCtx *Context) (err error) {
	afCtx.appSessionsEv = make(AppSessEv)
	afCtx.data.notifyAPIClient, err =
		NewAFNotifyAPIClient(afCtx)
	if err != nil {
		log.Errf("Unable to create notification api client")
		return err
	}
	return nil

}

func validateNotifyURI(notifyURI string) error {
	/* Check the url type - if its https or http */
	u, err := url.Parse(notifyURI)
	if err != nil {
		return err
	}

	// https is not supported
	if u.Scheme != "http" {
		return errors.New("Unsupported url scheme")
	}
	return nil
}

func getAppSessFromCorrID(corrID string, afCtx *Context) (evInfo *EventInfo,
	err error) {

	for _, value := range afCtx.appSessionsEv {

		for key := range value.upPathEv {

			if key == corrID {
				return value, nil
			}
		}
	}
	return evInfo, errors.New("AppSession Event Info Not Found")
}

func getAppSessFromURL(url string) string {
	res := strings.Split(url, "app-sessions")
	aID := strings.Split(res[1], "/")
	return aID[1]
}

/* Updating the notificationURI in afRouteReq in response*/
func updateRouteReqInResp(afRouteReq *RoutingRequirement,
	afCtx *Context) (err error) {
	if afRouteReq != nil && afRouteReq.UpPathChgSub != nil {
		var evInfo *EventInfo
		corrID := afRouteReq.UpPathChgSub.NotifCorreID
		evInfo, err = getAppSessFromCorrID(corrID, afCtx)
		if err == nil {
			afRouteReq.UpPathChgSub.NotificationURI =
				string(evInfo.upPathEv[corrID])

		}
	}
	return err
}

/* This function updates the consumer notificationURI in response
which was replaced by AF and also sends back websocketURI
if websocket delivery is requested*/
func updateAppSessInResp(appSess *AppSessionContext,
	appSessionID string, afCtx *Context) (err error) {
	ascReqData := appSess.AscReqData
	if ascReqData == nil {
		log.Infoln("Nil Application Session Context Request Data")
		return nil
	}

	evInfo := afCtx.appSessionsEv[appSessionID]
	if evInfo == nil {
		log.Infoln("Event Information not present, Nothing to update")
		return nil
	}
	if evInfo.wsReq {
		//TODO update the websocketURI in AscRespData
		return nil
	}
	afRouteReq := ascReqData.AfRoutReq

	err = updateRouteReqInResp(afRouteReq, afCtx)
	if err != nil {
		return err
	}

	medCompList := ascReqData.MedComponents
	for _, medcomp := range medCompList {
		afRouteReq = medcomp.AfRoutReq
		err = updateRouteReqInResp(afRouteReq, afCtx)
		if err != nil {
			return err
		}
	}

	return err
}

/* To check if websocket delivery is requested in ascReqData. Both websocket and
notificationURI is not allowed in one appSession */
func chkAppSessCreateForWs(appSess AppSessionContext,
	evInfo *EventInfo) (err error) {

	ascReqData := appSess.AscReqData
	if ascReqData == nil {
		err = errors.New("Nil AppSessionContextReqData")
		return err
	}

	if ascReqData.AfwebsockNotifConfig != nil {
		if ascReqData.AfwebsockNotifConfig.RequestWebsocketURI {
			evInfo.wsReq = true
			evInfo.consumerID = ascReqData.AfwebsockNotifConfig.ConsumerID
			if len(evInfo.consumerID) == 0 {
				err = errors.New("Nil ConsumerID")
				return err
			}

		}
	}
	return nil
}

/* To check if websocket delivery is requested in ascUpdateData. Both websocket
and notificationURI is not allowed in one appSession */
func chkAppSessUpdateForWs(ascUpdateData AppSessionContextUpdateData,
	evInfo *EventInfo) (err error) {

	if ascUpdateData.AfwebsockNotifConfig != nil {
		if ascUpdateData.AfwebsockNotifConfig.RequestWebsocketURI {
			if evInfo.wsReq {
				err = errors.New("Websocket Already Established with consumer")
				return err
			}
			evInfo.wsReq = true
			evInfo.consumerID = ascUpdateData.AfwebsockNotifConfig.ConsumerID
		}
	}
	return nil

}

/* setAppSessNotifParams updates the notificationURI/NotifURI in ascReqData with
AF generated one and stores the notification related params*/
func setAppSessNotifParams(appSess *AppSessionContext,
	evInfo *EventInfo, afCtx *Context) (err error) {

	ascReqData := appSess.AscReqData

	err = chkAppSessCreateForWs(*appSess, evInfo)
	if err != nil {
		return err
	}

	// NotifURI to send terminate requests
	ascReqData.NotifURI = pcfPANotifURI

	if ascReqData.EvSubsc != nil {
		ascReqData.EvSubsc.NotifURI = pcfPANotifURI
	}

	// To update the notification URI in afRoueReq and store it for
	// sending notifications
	err = updateRouteReqParamsCreate(ascReqData.AfRoutReq, evInfo, afCtx)
	if err != nil {
		return err
	}

	medCompList := ascReqData.MedComponents
	for _, medcomp := range medCompList {
		afRouteReq := medcomp.AfRoutReq
		err = updateRouteReqParamsCreate(afRouteReq, evInfo, afCtx)
		if err != nil {
			return err
		}
	}
	return err
}

/* modifyAppSessNotifParams in ascUpdateData updates the
notificationURI/NotifURI with AF generated one and stores the
notification related params*/
func modifyAppSessNotifParams(ascUpdateData *AppSessionContextUpdateData,
	appSessionID string, afCtx *Context) (err error) {

	if ascUpdateData == nil {
		err = errors.New("Nil AppSessionContextUpdateData")
		return err
	}

	evInfo := afCtx.appSessionsEv[appSessionID]

	if evInfo == nil {
		// No event information was present prior for this appSessionID
		evInfo = new(EventInfo)
	}

	err = chkAppSessUpdateForWs(*ascUpdateData, evInfo)
	if err != nil {
		return err
	}

	if ascUpdateData.EvSubsc != nil {
		ascUpdateData.EvSubsc.NotifURI = pcfPANotifURI
	}

	err = updateRouteReqParamsUpdate(ascUpdateData.AfRoutReq, evInfo, afCtx)
	if err != nil {
		return err
	}
	medCompList := ascUpdateData.MedComponents
	for _, medcomp := range medCompList {
		afRouteReq := medcomp.AfRoutReq
		err = updateRouteReqParamsUpdate(afRouteReq, evInfo, afCtx)
		if err != nil {
			return err
		}
	}

	afCtx.appSessionsEv[appSessionID] = evInfo
	log.Infoln("Updated the Event Info for appSessionID ", appSessionID)
	return err
}

/* This function checks if correlID is unique*/
func chkCorrelIDExists(corrID string, evInfo *EventInfo,
	afCtx *Context) bool {

	// Check within the current request
	for k := range evInfo.upPathEv {
		if k == corrID {
			/*Already present*/
			return true
		}
	}
	// Check in all appSessions stored
	for _, value := range afCtx.appSessionsEv {

		for key := range value.upPathEv {

			if key == corrID {
				/*Already present*/
				return true
			}
		}
	}
	return false
}

func updateRouteReqParamsCreate(afRouteReq *RoutingRequirement,
	evInfo *EventInfo, afCtx *Context) (err error) {

	if afRouteReq != nil && afRouteReq.UpPathChgSub != nil {
		if evInfo.wsReq &&
			len(afRouteReq.UpPathChgSub.NotificationURI) != 0 {
			err = errors.New("Both WS and NotificationUri not allowed")
			return err
		}

		if len(afRouteReq.UpPathChgSub.NotifCorreID) == 0 {
			err = errors.New("NotifCorrelID missing")
			return err
		}

		notifyID := afRouteReq.UpPathChgSub.NotifCorreID

		if chkCorrelIDExists(notifyID, evInfo, afCtx) {
			err = errors.New("Notif CorrelID already exists")
			return err
		}

		if evInfo.upPathEv == nil {
			evInfo.upPathEv = make(map[string]NotificationURI)
		}

		if len(afRouteReq.UpPathChgSub.NotificationURI) != 0 {
			err = validateNotifyURI(afRouteReq.UpPathChgSub.NotificationURI)
			if err != nil {
				return err
			}
			// store the consumer notificationURI
			evInfo.upPathEv[notifyID] = NotificationURI(
				afRouteReq.UpPathChgSub.NotificationURI)
			// Update with AF generated URI
			afRouteReq.UpPathChgSub.NotificationURI = smfPANotifURI

		} else if !evInfo.wsReq {
			err = errors.New("Neither WS nor notificationURI present")
			return err

		} else {
			// This is the case of websocket, only notifyID mapping with
			// appSessionID is needed.
			evInfo.upPathEv[notifyID] = ""
		}

	} else {

		log.Infoln("AfRouteReq/UpPathChgEvent not set in the request")

	}

	return nil
}

func updateRouteReqParamsUpdate(afRouteReq *RoutingRequirement, evInfo *EventInfo,
	afCtx *Context) (err error) {

	if afRouteReq != nil && afRouteReq.UpPathChgSub != nil {
		if evInfo.wsReq &&
			len(afRouteReq.UpPathChgSub.NotificationURI) != 0 {
			err = errors.New("Both WS and NotificationUri not Allowed")
			return err
		}

		if len(afRouteReq.UpPathChgSub.NotifCorreID) == 0 {
			err = errors.New("NotifCorrelID missing")
			return err
		}

		notifyID := afRouteReq.UpPathChgSub.NotifCorreID

		if evInfo.upPathEv == nil {
			evInfo.upPathEv = make(map[string]NotificationURI)
		}

		if len(afRouteReq.UpPathChgSub.NotificationURI) != 0 {
			err = validateNotifyURI(afRouteReq.UpPathChgSub.NotificationURI)
			if err != nil {
				return err
			}
			// store the consumer notificationUri
			evInfo.upPathEv[notifyID] = NotificationURI(
				afRouteReq.UpPathChgSub.NotificationURI)
			// update with AF generated uri
			afRouteReq.UpPathChgSub.NotificationURI = smfPANotifURI

		} else if !evInfo.wsReq {
			err = errors.New("Neither WS nor notificationURI present")
			return err

		} else {
			// This is the case of websocket, only notifyID mapping with
			// appSessionID is needed.
			evInfo.upPathEv[notifyID] = ""
		}

	} else {

		log.Infoln("AfRouteReq/UpPathChgEvent not set in the request")

	}

	return nil
}

/* This function is called when SMF UP_PATH_CH notification is received.
It maps the notification into NotificationUpPathChg and sends to consumer*/
func sendUpPathEventNotification(corrID string, afCtx *Context,
	nsmEvNo NsmfEventNotification) {

	var (
		ev NotificationUpPathChg
	)
	// Map the content of NsmfEventExposureNotification to NotificationUpPathChg
	evInfo, err := getAppSessFromCorrID(corrID, afCtx)
	if err != nil {

		log.Errf("PolicyAuthSMFNotify getAppSessFromCorrID [%s]: [%s]",
			corrID, err.Error())
		return
	}

	ev.NotifyID = corrID
	ev.GPSI = string(nsmEvNo.Gpsi)
	ev.DNAIChgType = nsmEvNo.DnaiChgType
	ev.SrcUEIPv4Addr = nsmEvNo.SourceUeIpv4Addr
	ev.SrcUEIPv6Prefix = nsmEvNo.SourceUeIpv6Prefix
	ev.TgtUEIP4Addr = nsmEvNo.TargetUeIpv4Addr
	ev.TgtUEIPv6Prefix = nsmEvNo.TargetUeIpv6Prefix
	ev.UEMac = nsmEvNo.UeMac
	ev.SourceTrafficRoute = nsmEvNo.SourceTraRouting
	ev.SubscribedEvent = SubscribedEvent("UP_PATH_CHANGE")
	ev.TargetTrafficRoute = nsmEvNo.TargetTraRouting

	if evInfo.wsReq {
		// TODO: Send over websocket

	} else {

		notificationURI := evInfo.upPathEv[corrID]
		log.Infof("PolicyAuthSMFNotify [NotifID, URL] => [%s,%s]",
			corrID,
			notificationURI)

		apiClient := afCtx.data.notifyAPIClient
		if apiClient == nil {
			log.Errln("PolicyAuthSMFNotify nil notifyAPIClient")
			return
		}
		// Send the request towards Consumer
		err := apiClient.NotificationUpPathEvent(notificationURI, ev)
		if err != nil {
			log.Errf("UP_PATH_CHANGE Sending to consumer failed : %s",
				err.Error())
			return
		}
	}
}
