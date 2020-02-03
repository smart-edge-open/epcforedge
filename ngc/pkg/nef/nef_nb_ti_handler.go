/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	//"strconv"

	"github.com/gorilla/mux"
)

func createNewSub(nefCtx *nefContext, afID string,
	ti TrafficInfluSub) (loc string, rsp nefSBRspData, err error) {

	var af *afData
	nef := &nefCtx.nef

	af, err = nef.nefGetAf(afID)

	if err != nil {
		log.Err("NO AF PRESENT CREATE AF")
		af, err = nef.nefAddAf(nefCtx, afID)
		if err != nil {
			return loc, rsp, err
		}
	} else {
		log.Infoln("AF PRESENT")
	}

	loc, rsp, err = af.afAddSubscription(nefCtx, ti)

	if err != nil {
		return loc, rsp, err
	}

	return loc, rsp, nil
}

// ReadAllTrafficInfluenceSubscription : API to read all the subscritions
func ReadAllTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	var subslist []TrafficInfluSub
	var rsp nefSBRspData
	var err error

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID : %s", vars["afId"])

	af, err := nef.nefGetAf(vars["afId"])

	if err != nil {
		/* Failure in getting AF with afId received. In this case no
		 * subscription data will be returned to AF */
		log.Infoln(err)
	} else {
		rsp, subslist, err = af.afGetSubscriptionList(nefCtx)
		if err != nil {
			log.Err(err)
			sendErrorResponseToAF(w, rsp)
			return
		}
	}

	mdata, err2 := json.Marshal(subslist)

	if err2 != nil {
		sendCustomeErrorRspToAF(w, 400, "Failed to MARSHAL Subscription data ")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//Send Success response to Network
	_, err = w.Write(mdata)
	if err != nil {
		log.Errf("Write Failed: %v", err)
		return
	}

	log.Infof("HTTP Response sent: %d", http.StatusOK)
}

// CreateTrafficInfluenceSubscription : Handles the traffic influence requested
// by AF
func CreateTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])

	b, err := ioutil.ReadAll(r.Body)
	defer closeReqBody(r)

	if err != nil {
		sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP POST Body")
		return
	}

	//Traffic Influence data
	trInBody := TrafficInfluSub{}

	//Convert the json Traffic Influence data into struct
	err1 := json.Unmarshal(b, &trInBody)

	if err1 != nil {
		log.Err(err1)
		sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal POST data")
		return
	}

	//validate the mandatory parameters
	resRsp, status := validateAFTrafficInfluenceData(trInBody)
	if !status {
		log.Err(resRsp.pd.Title)
		sendErrorResponseToAF(w, resRsp)
		return
	}

	loc, rsp, err3 := createNewSub(nefCtx, vars["afId"], trInBody)

	if err3 != nil {
		log.Err(err3)
		// we return bad request here since we have reached the max
		rsp.errorCode = 400
		sendErrorResponseToAF(w, rsp)
		return
	}
	log.Infoln(loc)

	trInBody.Self = Link(loc)

	//Martshal data and send into the body
	mdata, err2 := json.Marshal(trInBody)

	if err2 != nil {
		log.Err(err2)
		sendCustomeErrorRspToAF(w, 400, "Failed to Marshal GET response data")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Location", loc)

	// Response should be 201 Created as per 3GPP 29.522
	w.WriteHeader(http.StatusCreated)
	log.Infof("CreateTrafficInfluenceSubscription responses => %d",
		http.StatusCreated)
	_, err = w.Write(mdata)
	if err != nil {
		log.Errf("Write Failed: %v", err)
		return
	}
	nef := &nefCtx.nef
	logNef(nef)

}

// ReadTrafficInfluenceSubscription : Read a particular subscription details
func ReadTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, ok := nef.nefGetAf(vars["afId"])

	if ok != nil {
		sendCustomeErrorRspToAF(w, 404, "Failed to find AF records")
		return
	}

	rsp, sub, err := af.afGetSubscription(nefCtx, vars["subscriptionId"])

	if err != nil {
		log.Err(err)
		sendErrorResponseToAF(w, rsp)
		return
	}

	mdata, err2 := json.Marshal(sub)
	if err2 != nil {
		log.Err(err2)
		sendCustomeErrorRspToAF(w, 400, "Failed to Marshal GET response data")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(mdata)
	if err != nil {
		log.Errf("Write Failed: %v", err)
		return
	}

	log.Infof("HTTP Response sent: %d", http.StatusOK)
}

// UpdatePutTrafficInfluenceSubscription : Updates a traffic influence created
// earlier (PUT Req)
func UpdatePutTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, ok := nef.nefGetAf(vars["afId"])
	if ok == nil {

		b, err := ioutil.ReadAll(r.Body)
		defer closeReqBody(r)

		if err != nil {
			log.Err(err)
			sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP PUT Body")
			return
		}

		//Traffic Influence data
		trInBody := TrafficInfluSub{}

		//Convert the json Traffic Influence data into struct
		err1 := json.Unmarshal(b, &trInBody)

		if err1 != nil {
			log.Err(err1)
			sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal PUT data")
			return
		}

		rsp, newTI, err := af.afUpdateSubscription(nefCtx,
			vars["subscriptionId"], trInBody)

		if err != nil {
			sendErrorResponseToAF(w, rsp)
			return
		}

		mdata, err2 := json.Marshal(newTI)

		if err2 != nil {
			log.Err(err2)
			sendCustomeErrorRspToAF(w, 400, "Failed to Marshal PUT"+
				"response data")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		_, err = w.Write(mdata)
		if err != nil {
			log.Errf("Write Failed: %v", err)
		}
		return

	}
	log.Infoln(ok)
	sendCustomeErrorRspToAF(w, 404, "Failed to find AF records")

}

// UpdatePatchTrafficInfluenceSubscription : Updates a traffic influence created
//  earlier (PATCH Req)
func UpdatePatchTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, ok := nef.nefGetAf(vars["afId"])
	if ok == nil {

		b, err := ioutil.ReadAll(r.Body)

		defer closeReqBody(r)

		if err != nil {
			log.Err(err)
			sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP PATCH Body")
			return
		}

		//Traffic Influence Sub Patch data
		TrInSPBody := TrafficInfluSubPatch{}

		//Convert the json Traffic Influence data into struct
		err1 := json.Unmarshal(b, &TrInSPBody)

		if err1 != nil {
			log.Err(err1)
			sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal PATCH data")
			return
		}

		rsp, ti, err := af.afPartialUpdateSubscription(nefCtx,
			vars["subscriptionId"], TrInSPBody)

		if err != nil {
			sendErrorResponseToAF(w, rsp)
			return
		}

		mdata, err2 := json.Marshal(ti)

		if err2 != nil {
			log.Err(err2)
			sendCustomeErrorRspToAF(w, 400,
				"Failed to Marshal PATCH response data")
			return

		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		_, err = w.Write(mdata)
		if err != nil {
			log.Errf("Write Failed: %v", err)
		}
		return
	}

	log.Infoln(ok)
	sendCustomeErrorRspToAF(w, 404, "Failed to find AF records")
}

// DeleteTrafficInfluenceSubscription : Deletes a traffic influence created by
//  AF
func DeleteTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, err := nef.nefGetAf(vars["afId"])

	if err != nil {
		log.Err(err)
		sendCustomeErrorRspToAF(w, 404, "Failed to find AF entry")
		return
	}
	rsp, err := af.afDeleteSubscription(nefCtx, vars["subscriptionId"])

	if err != nil {
		log.Err(err)
		sendErrorResponseToAF(w, rsp)
		return
	}

	// Response should be 204 as per 3GPP 29.522
	w.WriteHeader(http.StatusNoContent)

	log.Infof("HTTP Response sent: %d", http.StatusNoContent)

	if af.afGetSubCount() == 0 {

		_ = nef.nefDeleteAf(vars["afId"])
	}

	logNef(nef)
}

// NotifySmfUPFEvent : Handles the SMF notification for UPF event
func NotifySmfUPFEvent(w http.ResponseWriter,
	r *http.Request) {

	var (
		smfEv    NsmfEventExposureNotification
		ev       EventNotification
		afURL    URI
		nsmEvNo  NsmEventNotification
		i        int
		upfFound bool
	)

	if r.Body == nil {
		log.Errf("NotifySmfUPFEvent Empty Body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Retrieve the event notification information from the request
	if err := json.NewDecoder(r.Body).Decode(&smfEv); err != nil {
		log.Errf("NotifySmfUPFEvent body parse: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate the content of the NsmfEventExposureNotification
	// Check if notification id is present
	if smfEv.NotifID == "" {
		log.Errf("NotifySmfUPFEvent missing notif id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if notification events with UP_PATH_CH is present
	if len(smfEv.EventNotifs) == 0 {
		log.Errf("NotifySmfUPFEvent missing event notifications")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, nsmEvNo = range smfEv.EventNotifs {
		if nsmEvNo.Event == "UP_PATH_CH" {
			log.Infof("NotifySmfUPFEvent found an entry for UP_PATH_CH"+
				"at index: %d", i)
			upfFound = true
			break
		}

	}

	if !upfFound {
		log.Errf("NotifySmfUPFEvent missing event with UP_PATH_CH")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Map the content of NsmfEventExposureNotification to EventNotificaiton
	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	afSubs, err1 := getSubFromCorrID(nefCtx, smfEv.NotifID)
	if err1 != nil {
		log.Errf("NotifySmfUPFEvent getSubFromCorrId [%s]: %s",
			smfEv.NotifID, err1.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Infof("NotifySmfUPFEvent [NotifID, TransId, URL] => [%s,%s,%s",
		smfEv.NotifID, afSubs.ti.AfTransID,
		afSubs.ti.NotificationDestination)

	ev.AfTransID = afSubs.ti.AfTransID
	afURL = URI(afSubs.ti.NotificationDestination)
	ev.Gpsi = nsmEvNo.Gpsi
	ev.DnaiChgType = nsmEvNo.DnaiChgType
	ev.SrcUeIpv4Addr = nsmEvNo.SourceUeIpv4Addr
	ev.SrcUeIpv6Prefix = nsmEvNo.SourceUeIpv6Prefix
	ev.TgtUeIpv4Addr = nsmEvNo.TargetUeIpv4Addr
	ev.TgtUeIpv6Prefix = nsmEvNo.TargetUeIpv6Prefix
	ev.UeMac = nsmEvNo.UeMac
	ev.SourceTrafficRoute = nsmEvNo.SourceTraRouting
	ev.SubscribedEvent = SubscribedEvent("UP_PATH_CHANGE")
	ev.TargetTrafficRoute = nsmEvNo.TargetTraRouting

	w.WriteHeader(http.StatusOK)

	// Send the request towards AF
	var afClient AfNotification = NewAfClient(&nefCtx.cfg)
	err := afClient.AfNotificationUpfEvent(r.Context(), afURL, ev)
	if err != nil {
		log.Errf("NotifySmfUPFEvent sending to AF failed : %s",
			err.Error())
	}
}

func getSubFromCorrID(nefCtx *nefContext, corrID string) (sub *afSubscription,
	err error) {

	nef := &nefCtx.nef

	/*Search across all the AF registered */
	for _, value := range nef.afs {

		/*Search across all the Subscription*/
		for _, vs := range value.subs {

			if vs.NotifCorreID == corrID {
				/*Match found return sub handle*/
				return vs, nil
			}
		}
	}
	return sub, errors.New("Subscription Not Found")
}

//validateAFTrafficInfluenceData: Function to validate mandatory parameters of
//TrafficInfluence received from AF
func validateAFTrafficInfluenceData(ti TrafficInfluSub) (rsp nefSBRspData,
	status bool) {

	if len(ti.AfTransID) == 0 {
		rsp.errorCode = 400
		rsp.pd.Title = "Missing AfTransID atttribute"
		return rsp, false
	}

	//In case AfServiceID  is not present then DNN has to be included in TI
	if len(ti.AfServiceID) == 0 && len(ti.Dnn) == 0 {

		rsp.errorCode = 400
		rsp.pd.Title = "Missing afServiceId atttribute"
		return rsp, false
	}

	if len(ti.AfAppID) == 0 && ti.TrafficFilters == nil &&
		ti.EthTrafficFilters == nil {
		rsp.errorCode = 400
		rsp.pd.Title = "missing one of afAppId, trafficFilters," +
			"ethTrafficFilters"
		return rsp, false
	}
	return rsp, true
}

//Creates a new subscription
func (af *afData) afAddSubscription(nefCtx *nefContext,
	ti TrafficInfluSub) (loc string, rsp nefSBRspData, err error) {

	/*Check if max subscription reached */
	if len(af.subs) >= nefCtx.cfg.MaxSubSupport {

		rsp.errorCode = 400
		rsp.pd.Title = "MAX Subscription Reached"
		return "", rsp, errors.New("MAX SUBS Created")
	}

	//Generate a unique subscription ID string
	subIDStr := strconv.Itoa(af.subIDnum)
	af.subIDnum++

	//Create Subscription data
	afsub := afSubscription{subid: subIDStr, ti: ti, appSessionID: "",
		NotifCorreID: "", iid: ""}

	if len(ti.Gpsi) > 0 || len(ti.Ipv4Addr) > 0 || len(ti.Ipv6Addr) > 0 {

		//Applicable to single UE, PCF case

		rsp, err = nefSBPCFPost(&afsub, nefCtx, ti)

		if err != nil {

			//Return error failed to create subscription
			return "", rsp, err
		}
		//Store Notification Destination URI
		afsub.afNotificationDestination = ti.NotificationDestination
		afsub.NEFSBGet = nefSBPCFGet
		afsub.NEFSBPut = nefSBPCFPut
		afsub.NEFSBPatch = nefSBPCFPatch
		afsub.NEFSBDelete = nefSBPCFDelete

	} else if len(ti.ExternalGroupID) > 0 || ti.AnyUeInd {

		//Applicable to Any UE, UDR case

		rsp, err = nefSBUDRPost(&afsub, nefCtx, ti)

		if err != nil {

			//Return error
			return "", rsp, err
		}
		//Store Notification Destination URI
		afsub.afNotificationDestination = ti.NotificationDestination

		afsub.NEFSBGet = nefSBUDRGet
		afsub.NEFSBPut = nefSBUDRPut
		afsub.NEFSBPatch = nefSBUDRPatch
		afsub.NEFSBDelete = nefSBUDRDelete

	} else {
		//Invalid case. Return Error
		rsp.errorCode = 400
		rsp.pd.Title = "Invalid Request"
		return "", rsp, errors.New("Invalid AF Request")
	}

	//Link the subscription with the AF
	af.subs[subIDStr] = &afsub

	//Create Location URI
	loc = nefCtx.nef.locationURLPrefix + af.afID + "/subscriptions/" +
		subIDStr

	afsub.ti.Self = Link(loc)

	log.Infoln(" NEW AF Subscription added " + subIDStr)

	return loc, rsp, nil
}

func (af *afData) afUpdateSubscription(nefCtx *nefContext, subID string,
	ti TrafficInfluSub) (rsp nefSBRspData, updtTI TrafficInfluSub, err error) {

	sub, ok := af.subs[subID]

	if !ok {
		rsp.errorCode = 400
		rsp.pd.Title = subNotFound

		return rsp, updtTI, errors.New(subNotFound)
	}

	rsp, err = sub.NEFSBPut(sub, nefCtx, ti)

	if err != nil {
		log.Err("Failed to Update Subscription")
		return rsp, updtTI, err
	}

	updtTI = ti
	updtTI.Self = sub.ti.Self
	sub.ti = updtTI

	log.Infoln("Update Subscription Successful")
	return rsp, updtTI, err
}

func updateTiFromTisp(ti *TrafficInfluSub, tisp TrafficInfluSubPatch) {

	if tisp.AppReloInd != ti.AppReloInd {
		log.Infoln("Updating AppReloInd...")
		ti.AppReloInd = tisp.AppReloInd
	}

	if tisp.TrafficFilters != nil {
		log.Infoln("Updating TrafficFilters...")
		ti.TrafficFilters = tisp.TrafficFilters
	}

	if tisp.EthTrafficFilters != nil {
		log.Infoln("Updating EthTrafficFilters")
		ti.EthTrafficFilters = tisp.EthTrafficFilters
	}

	if tisp.TrafficRoutes != nil {
		log.Infoln("Updating TrafficRoutes")
		ti.TrafficRoutes = tisp.TrafficRoutes
	}
	if tisp.TempValidities != nil {
		log.Infoln("Updating TempValidities")
		ti.TempValidities = tisp.TempValidities
	}
	if tisp.ValidGeoZoneIDs != nil {
		log.Infoln("Updating ValidGeoZoneIDs")
		ti.ValidGeoZoneIDs = tisp.ValidGeoZoneIDs
	}

}

func (af *afData) afPartialUpdateSubscription(nefCtx *nefContext, subID string,
	tisp TrafficInfluSubPatch) (rsp nefSBRspData, ti TrafficInfluSub,
	err error) {

	sub, ok := af.subs[subID]

	if !ok {
		rsp.errorCode = 400
		rsp.pd.Title = subNotFound
		return rsp, ti, errors.New(subNotFound)
	}

	rsp, err = sub.NEFSBPatch(sub, nefCtx, tisp)

	if err != nil {
		log.Err("Failed to Patch Subscription")
		return rsp, ti, err
	}
	updateTiFromTisp(&sub.ti, tisp)

	return rsp, sub.ti, err

}

func (af *afData) afGetSubscription(nefCtx *nefContext,
	subID string) (rsp nefSBRspData, ti TrafficInfluSub, err error) {

	sub, ok := af.subs[subID]

	if !ok {
		rsp.errorCode = 404
		rsp.pd.Title = subNotFound
		return rsp, ti, errors.New(subNotFound)
	}

	//ti, rsp, err = sub.NEFSBGet(sub, nefCtx)

	/*
		if err != nil {
			log.Infoln("Failed to Get Subscription")
			return rsp, ti, err
		}

		return rsp, ti, err
	*/

	//Return locally
	return rsp, sub.ti, err
}

func (af *afData) afGetSubscriptionList(nefCtx *nefContext) (rsp nefSBRspData,
	subsList []TrafficInfluSub, err error) {

	var ti TrafficInfluSub

	if len(af.subs) > 0 {

		for key := range af.subs {

			rsp, ti, err = af.afGetSubscription(nefCtx, key)

			if err != nil {
				return rsp, subsList, err
			}
			subsList = append(subsList, ti)
		}
	}
	return rsp, subsList, err
}

func (af *afData) afDeleteSubscription(nefCtx *nefContext,
	subID string) (rsp nefSBRspData, err error) {

	//Check if AF is already present
	sub, ok := af.subs[subID]

	if !ok {
		rsp.errorCode = 404
		rsp.pd.Title = subNotFound
		return rsp, errors.New(subNotFound)
	}

	rsp, err = sub.NEFSBDelete(sub, nefCtx)

	if err != nil {
		log.Err("Failed to Delete Subscription")
		return rsp, err
	}

	//Delete local entry in map
	delete(af.subs, subID)
	//af.subIDnum--

	return rsp, err
}

func (af *afData) afGetSubCount() (afCount int) {

	return len(af.subs)
}

/* unused function
func (af *afData) afDestroy(afid string) error {

	//Todo delete all subscriptions, needed in go ??
	//Needed for gracefully disconnecting
	return errors.New("AF data cleaned")
}
*/

// Generate the notification uri to be provided to PCF/UDR
func getNefNotificationURI(cfg *Config) URI {
	var uri string
	// If http2 port is configured use it else http port
	if cfg.HTTP2Config.Endpoint != "" {
		uri = "https://" + cfg.NefAPIRoot +
			cfg.HTTP2Config.Endpoint
	} else {
		uri = "http://" + cfg.NefAPIRoot +
			cfg.HTTPConfig.Endpoint
	}
	uri += cfg.UpfNotificationResURIPath
	return URI(uri)
}

func getNefLocationURLPrefix(cfg *Config) string {

	var uri string
	// If http2 port is configured use it else http port
	if cfg.HTTP2Config.Endpoint != "" {
		uri = "https://" + cfg.NefAPIRoot +
			cfg.HTTP2Config.Endpoint
	} else {
		uri = "http://" + cfg.NefAPIRoot +
			cfg.HTTPConfig.Endpoint
	}
	uri += cfg.LocationPrefix
	return uri

}

// nefSBPCFPost : This function sends HTTP POST Request to PCF to create Policy
//             Authorization.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
//   - ti: This is Traffic Influence Subscription Data.
// Output Args:
//    - rsp: This is Policy Authorization Create Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBPCFPost(pcfSub *afSubscription, nefCtx *nefContext,
	ti TrafficInfluSub) (rsp nefSBRspData, err error) {

	var appSessID AppSessionID
	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	pcfSub.NotifCorreID = strconv.Itoa(int(nef.corrID))
	nef.corrID++

	appSessCtx := AppSessionContext{}
	pcfPolicyResp := PcfPolicyResponse{}

	//Populating App Session Context Data Req
	appSessCtx.AscReqData.AfAppID = AfAppID(ti.AfAppID)
	appSessCtx.AscReqData.AfRoutReq.AppReloc = ti.AppReloInd

	//Populating UP Path Chnage Subbscription Data in App Session Context
	appSessCtx.AscReqData.AfRoutReq.UpPathChgSub.DnaiChgType = ti.DnaiChgType

	// If http2 port is configured use it else http port
	if nefCtx.cfg.HTTP2Config.Endpoint != "" {
		appSessCtx.AscReqData.AfRoutReq.UpPathChgSub.NotificationURI =
			URI("https://" + nefCtx.cfg.NefAPIRoot +
				nefCtx.cfg.HTTP2Config.Endpoint)
	} else {
		appSessCtx.AscReqData.AfRoutReq.UpPathChgSub.NotificationURI =
			URI("http://" + nefCtx.cfg.NefAPIRoot +
				nefCtx.cfg.HTTPConfig.Endpoint)
	}

	/*+ nefCtx.cfg.UpfNotificationResUriPath*/
	appSessCtx.AscReqData.AfRoutReq.UpPathChgSub.NotifCorreID =
		pcfSub.NotifCorreID

	//Populating Traffic Routes in App Session Context
	appSessCtx.AscReqData.AfRoutReq.RouteToLocs = make([]RouteToLocation,
		len(ti.TrafficRoutes))
	_ = copy(appSessCtx.AscReqData.AfRoutReq.RouteToLocs, ti.TrafficRoutes)

	//Populating Temporal Validity in App Session Context
	appSessCtx.AscReqData.AfRoutReq.TempVals = make([]TemporalValidity,
		len(ti.TempValidities))
	_ = copy(appSessCtx.AscReqData.AfRoutReq.TempVals, ti.TempValidities)

	//Populating Spatial Validity in App Session Context
	_ = getSpatialValidityData(cliCtx, nefCtx,
		&appSessCtx.AscReqData.AfRoutReq.SpVal)

	//Populating IP and Mac Addresses in App Session Context
	appSessCtx.AscReqData.UeIpv4 = ti.Ipv4Addr
	appSessCtx.AscReqData.UeIpv6 = ti.Ipv6Addr
	appSessCtx.AscReqData.UeMac = ti.MacAddr

	//Populating DNN and NW Slice Info and SUPI in App Session Context
	for _, afServIdcounter := range nefCtx.cfg.AfServiceIDs {
		afServiceID := afServIdcounter.(map[string]interface{})
		if 0 == strings.Compare(ti.AfServiceID, afServiceID["id"].(string)) {
			appSessCtx.AscReqData.Dnn = Dnn(afServiceID["dnn"].(string))
			appSessCtx.AscReqData.SliceInfo.Sd = afServiceID["snssai"].(string)
			appSessCtx.AscReqData.SliceInfo.Sst =
				uint8(len(appSessCtx.AscReqData.SliceInfo.Sd))
		}
	}

	//Populating SUPI in App Session Context
	_ = getSupiData(cliCtx, nefCtx, &appSessCtx.AscReqData.Supi)

	appSessID, pcfPolicyResp, err =
		nef.pcfClient.PolicyAuthorizationCreate(cliCtx, appSessCtx)

	if err != nil {
		rsp.errorCode = int(pcfPolicyResp.ResponseCode)
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Create Failure. Response Code: %d",
			rsp.errorCode)
		return rsp, err
	}

	rsp.errorCode = int(pcfPolicyResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Create Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		pcfSub.appSessionID = appSessID
		log.Infof("PCF Policy Authorization Create Success. Response Code: %d",
			rsp.errorCode)
	}
	return rsp, err
}

// nefSBPCFGet : This function sends HTTP GET Request to PCF to fetch Policy
//             Authorization using App Session Context Key.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
// Output Args:
//    - sub: This is Traffic Influence Subscription Data.
//    - rsp: This is Policy Authorization Get Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBPCFGet(pcfSub *afSubscription, nefCtx *nefContext) (
	sub TrafficInfluSub, rsp nefSBRspData, err error) {

	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	pcfPolicyResp, err :=
		nef.pcfClient.PolicyAuthorizationGet(cliCtx, pcfSub.appSessionID)
	if err != nil {
		rsp.errorCode = int(pcfPolicyResp.ResponseCode)
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Get Failure. Response Code: %d",
			rsp.errorCode)
		return sub, rsp, err
	}

	rsp.errorCode = int(pcfPolicyResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Get Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		sub = pcfSub.ti
		log.Infof("PCF Policy Authorization Get Success. Response Code: %d",
			rsp.errorCode)
	}

	return sub, rsp, err
}

// nefSBPCFPut : This function returns error as HTTP PUT Request to PCF is not
//            supported.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
//   - ti: This is Traffic Influence Subscription Data.
// Output Args:
//    - rsp: This is Policy Authorization Put Response Data
//    - error: retruns error .
func nefSBPCFPut(pcfSub *afSubscription, nefCtx *nefContext,
	ti TrafficInfluSub) (rsp nefSBRspData, err error) {
	err = errors.New("PUT Method Not Supported")
	log.Errf("PCF Policy Authorization Put Not Supported")
	return rsp, err
}

// nefSBPCFPatch : This function sends HTTP PATCH Request to PCF to update
//             Policy Authorization using App Session Context Key.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
//   - tisp: This is Traffic Influence Subscription Patch Data.
// Output Args:
//    - rsp: This is Policy Authorization Patch Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBPCFPatch(pcfSub *afSubscription, nefCtx *nefContext,
	tisp TrafficInfluSubPatch) (rsp nefSBRspData, err error) {

	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	appSessCtxUpdtData := AppSessionContextUpdateData{}

	//Populating App Session Context Data Req
	//TODO: appSessCtxUpdtData.AfAppID = ti.AfAppID
	appSessCtxUpdtData.AfRoutReq.AppReloc = tisp.AppReloInd

	//Populating UP Path Chnage Subbscription Data in App Session Context
	appSessCtxUpdtData.AfRoutReq.UpPathChgSub.NotificationURI =
		nefCtx.nef.upfNotificationURL

	appSessCtxUpdtData.AfRoutReq.UpPathChgSub.NotifCorreID =
		pcfSub.NotifCorreID

	//Populating Traffic Routes in App Session Context
	appSessCtxUpdtData.AfRoutReq.RouteToLocs = make([]RouteToLocation,
		len(tisp.TrafficRoutes))
	_ = copy(appSessCtxUpdtData.AfRoutReq.RouteToLocs, tisp.TrafficRoutes)

	//Populating Temporal Validity in App Session Context
	appSessCtxUpdtData.AfRoutReq.TempVals = make([]TemporalValidity,
		len(tisp.TempValidities))
	_ = copy(appSessCtxUpdtData.AfRoutReq.TempVals, tisp.TempValidities)

	//Populating Spatial Validity in App Session Context
	_ = getSpatialValidityData(cliCtx, nefCtx,
		&appSessCtxUpdtData.AfRoutReq.SpVal)

	pcfPolicyResp, err := nef.pcfClient.PolicyAuthorizationUpdate(cliCtx,
		appSessCtxUpdtData, pcfSub.appSessionID)
	if err != nil {
		rsp.errorCode = int(pcfPolicyResp.ResponseCode)
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Update Failure. Response Code: %d",
			rsp.errorCode)
		return rsp, err
	}

	rsp.errorCode = int(pcfPolicyResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Update Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		log.Infof("PCF Policy Authorization Update Success. Response Code: %d",
			rsp.errorCode)
	}

	return rsp, err
}

// nefSBPCFDelete : This function sends HTTP DELETE Request to PCF to delete
// 				    Policy Authorization using App Session Context Key.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
// Output Args:
//    - rsp: This is Policy Authorization Delete Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBPCFDelete(pcfSub *afSubscription, nefCtx *nefContext) (
	rsp nefSBRspData, err error) {

	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	pcfPolicyResp, err :=
		nef.pcfClient.PolicyAuthorizationDelete(cliCtx, pcfSub.appSessionID)
	if err != nil {
		rsp.errorCode = int(pcfPolicyResp.ResponseCode)
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Delete Failure. Response Code: %d",
			rsp.errorCode)
		return rsp, err
	}

	rsp.errorCode = int(pcfPolicyResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if pcfPolicyResp.Pd != nil {
			rsp.pd = *pcfPolicyResp.Pd
		}
		log.Errf("PCF Policy Authorization Delete Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		log.Infof("PCF Policy Authorization Delete Success. Response Code: %d",
			rsp.errorCode)
	}

	return rsp, err
}

// nefSBUDRPost :  HTTP POST Request to UDR is to trigger Put
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
//   - ti: This is Traffic Influence Subscription Data.
// Output Args:
//    - rsp: This is Traffic Influence Data Put Response Data
//    - error: retruns error .
func nefSBUDRPost(udrSub *afSubscription, nefCtx *nefContext,
	ti TrafficInfluSub) (rsp nefSBRspData, err error) {

	rsp, err = nefSBUDRPut(udrSub, nefCtx, ti)
	return rsp, err

}

// nefSBUDRGet : This function sends HTTP GET Request to UDR to fetch
//               Traffic Influence Data.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
// Output Args:
//    - sub: This is Traffic Influence Subscription Data.
//    - rsp: This is Traffic Influence Data Delete Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBUDRGet(udrSub *afSubscription, nefCtx *nefContext) (
	sub TrafficInfluSub, rsp nefSBRspData, err error) {

	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	udrInfluenceResp, err := nef.udrClient.UdrInfluenceDataGet(cliCtx)
	if err != nil {
		rsp.errorCode = int(udrInfluenceResp.ResponseCode)
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Get Failure. Response Code: %d",
			rsp.errorCode)
		return sub, rsp, err
	}

	rsp.errorCode = int(udrInfluenceResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Get Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		sub = udrSub.ti
		log.Infof("UDR Traffic Influence Data Get Success. Response Code: %d",
			rsp.errorCode)
	}

	return sub, rsp, err
}

// nefSBUDRPut : This function sends HTTP PUT Request to UDR to create Traffic
//            Influence Data.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
//   - ti: This is Traffic Influence Subscription Data.
// Output Args:
//    - rsp: This is Traffic Influence Data Create Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBUDRPut(udrSub *afSubscription, nefCtx *nefContext,
	ti TrafficInfluSub) (rsp nefSBRspData, err error) {

	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	trafficInfluData := TrafficInfluData{}

	//Populating Traffic Influence Data
	trafficInfluData.AfAppID = ti.AfAppID

	//Populating DNN and NW Slice Info in Traffic Influence Data
	for _, afServIdcounter := range nefCtx.cfg.AfServiceIDs {
		afServiceID := afServIdcounter.(map[string]interface{})
		if 0 == strings.Compare(ti.AfServiceID, afServiceID["id"].(string)) {
			trafficInfluData.Dnn = Dnn(afServiceID["dnn"].(string))
			trafficInfluData.Snssai.Sd = afServiceID["snssai"].(string)
			trafficInfluData.Snssai.Sst = uint8(len(trafficInfluData.Snssai.Sd))
		}
	}

	trafficInfluData.AppReloInd = ti.AppReloInd
	trafficInfluData.InterGroupID = string(ti.ExternalGroupID)

	//Populating UP Path Chnage Subbscription Data in Traffic Influence Data
	trafficInfluData.UpPathChgNotifURI = nefCtx.nef.upfNotificationURL

	if len(ti.SubscribedEvents) > 0 &&
		0 == strings.Compare(string(ti.SubscribedEvents[0]), "UP_PATH_CHANGE") {
		udrSub.NotifCorreID = strconv.Itoa(int(nef.corrID))
		nef.corrID++
		trafficInfluData.UpPathChgNotifCorreID = udrSub.NotifCorreID
	}

	//Populating Traffic Filters in Traffic Influence Data
	trafficInfluData.TrafficFilters = make([]FlowInfo,
		len(ti.TrafficFilters))
	_ = copy(trafficInfluData.TrafficFilters, ti.TrafficFilters)

	//Populating Eth Traffic Filters in Traffic Influence Data
	trafficInfluData.EthTrafficFilters = make([]EthFlowDescription,
		len(ti.EthTrafficFilters))
	_ = copy(trafficInfluData.EthTrafficFilters, ti.EthTrafficFilters)

	//Populating Traffic Routes in Traffic Influence Data
	trafficInfluData.TrafficRoutes = make([]RouteToLocation,
		len(ti.TrafficRoutes))
	_ = copy(trafficInfluData.TrafficRoutes, ti.TrafficRoutes)

	//Populating Temporal Validity in Traffic Influence Data
	if 0 < len(ti.TempValidities) {
		trafficInfluData.ValidStartTime =
			DateTime(ti.TempValidities[0].StartTime)
		trafficInfluData.ValidEndTime =
			DateTime(ti.TempValidities[0].StopTime)
	}

	//Populating Spatial Validity in Traffic Influence Data
	_ = getNetworkAreaInfo(cliCtx, nefCtx, &trafficInfluData.NwAreaInfo)

	udrInfluenceResp, err := nef.udrClient.UdrInfluenceDataCreate(
		cliCtx, trafficInfluData, udrSub.iid)
	if err != nil {
		rsp.errorCode = int(udrInfluenceResp.ResponseCode)
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Put Failure. Response Code: %d",
			rsp.errorCode)
		return rsp, err
	}

	rsp.errorCode = int(udrInfluenceResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Put Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		log.Infof("UDR Traffic Influence Data Put Success. Response Code: %d",
			rsp.errorCode)
	}
	return rsp, err
}

// nefSBUDRPatch : This function sends HTTP PATCH Request to UDR to update
//            Traffic Influence Data.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
//   - tisp: This is Traffic Influence Subscription Patch Data.
// Output Args:
//    - rsp: This is Traffic Influence Data Update Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBUDRPatch(udrSub *afSubscription, nefCtx *nefContext,
	tisp TrafficInfluSubPatch) (rsp nefSBRspData, err error) {

	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	trafficInfluDataPatch := TrafficInfluDataPatch{}

	trafficInfluDataPatch.AppReloInd = tisp.AppReloInd

	//Populating Traffic Filters in Traffic Influence Data
	trafficInfluDataPatch.TrafficFilters = make([]FlowInfo,
		len(tisp.TrafficFilters))
	_ = copy(trafficInfluDataPatch.TrafficFilters, tisp.TrafficFilters)

	//Populating Eth Traffic Filters in Traffic Influence Data
	trafficInfluDataPatch.EthTrafficFilters = make([]EthFlowDescription,
		len(tisp.EthTrafficFilters))
	_ = copy(trafficInfluDataPatch.EthTrafficFilters, tisp.EthTrafficFilters)

	//Populating Traffic Routes in Traffic Influence Data
	trafficInfluDataPatch.TrafficRoutes = make([]RouteToLocation,
		len(tisp.TrafficRoutes))
	_ = copy(trafficInfluDataPatch.TrafficRoutes, tisp.TrafficRoutes)

	//Populating Temporal Validity in Traffic Influence Data
	if 0 < len(tisp.TempValidities) {
		trafficInfluDataPatch.ValidStartTime =
			DateTime(tisp.TempValidities[0].StartTime)
		trafficInfluDataPatch.ValidEndTime =
			DateTime(tisp.TempValidities[0].StopTime)
	}

	udrInfluenceResp, err := nef.udrClient.UdrInfluenceDataUpdate(
		cliCtx, trafficInfluDataPatch, udrSub.iid)
	if err != nil {
		rsp.errorCode = int(udrInfluenceResp.ResponseCode)
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Update Failure. Response Code: %d",
			rsp.errorCode)
		return rsp, err
	}

	rsp.errorCode = int(udrInfluenceResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Update Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		log.Infof("UDR Traffic Influence Data Update Success.Response Code: %d",
			rsp.errorCode)
	}
	return rsp, err
}

// nefSBUDRDelete : This function sends HTTP DELETE Request to UDR to delete
//               Traffic Influence Data.
// Input Args:
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
// Output Args:
//    - rsp: This is Traffic Influence Data Delete Response Data
//    - error: retruns error in case there is failure happened in sending the
//             request or any failure response is received.
func nefSBUDRDelete(udrSub *afSubscription, nefCtx *nefContext) (
	rsp nefSBRspData, err error) {

	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	udrInfluenceResp, err := nef.udrClient.UdrInfluenceDataDelete(cliCtx,
		udrSub.iid)
	if err != nil {
		rsp.errorCode = int(udrInfluenceResp.ResponseCode)
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Delete Failure. Response Code: %d",
			rsp.errorCode)
		return rsp, err
	}

	rsp.errorCode = int(udrInfluenceResp.ResponseCode)
	if rsp.errorCode >= 300 && rsp.errorCode < 700 {
		if udrInfluenceResp.Pd != nil {
			rsp.pd = *udrInfluenceResp.Pd
		}
		log.Errf("UDR Traffic Influence Data Delete Failure. Response Code: %d",
			rsp.errorCode)
	} else {
		log.Infof("UDR Traffic Influence Data Delete Success."+
			"Response Code: %d", rsp.errorCode)
	}

	return rsp, err
}

func getSpatialValidityData(cliCtx context.Context, nefCtx *nefContext,
	spVal *SpatialValidity) error {
	_ = cliCtx
	_ = nefCtx

	spVal.PresenceInfoList.PraID = "PRA_01"
	spVal.PresenceInfoList.PresenceState = "IN_AREA"
	spVal.PresenceInfoList.EcgiList = make([]Ecgi, 1)
	spVal.PresenceInfoList.EcgiList[0].EutraCellID = "EUTRACELL_01"
	spVal.PresenceInfoList.EcgiList[0].PlmnID.Mcc = "634"
	spVal.PresenceInfoList.EcgiList[0].PlmnID.Mnc = "635"
	spVal.PresenceInfoList.NcgiList = make([]Ncgi, 1)
	spVal.PresenceInfoList.NcgiList[0].NrCellID = "NRCELL_01"
	spVal.PresenceInfoList.NcgiList[0].PlmnID.Mcc = "834"
	spVal.PresenceInfoList.NcgiList[0].PlmnID.Mnc = "835"
	spVal.PresenceInfoList.GlobalRanNodeIDList = make([]GlobalRanNodeID, 1)
	spVal.PresenceInfoList.GlobalRanNodeIDList[0].PlmnID.Mcc = "934"
	spVal.PresenceInfoList.GlobalRanNodeIDList[0].PlmnID.Mnc = "935"
	spVal.PresenceInfoList.GlobalRanNodeIDList[0].N3IwfID = "IWF_01"
	spVal.PresenceInfoList.GlobalRanNodeIDList[0].GNbID.BitLength = 48
	spVal.PresenceInfoList.GlobalRanNodeIDList[0].GNbID.GNBValue = "GNB_01"
	spVal.PresenceInfoList.GlobalRanNodeIDList[0].NgeNbID = "NB_01"

	return nil
}

func getSupiData(cliCtx context.Context, nefCtx *nefContext,
	supi *Supi) error {
	_ = cliCtx
	_ = nefCtx
	*supi = "imsi-8"
	return nil
}

func getNetworkAreaInfo(cliCtx context.Context, nefCtx *nefContext,
	nwAreaInfo *NetworkAreaInfo) error {
	_ = cliCtx
	_ = nefCtx

	nwAreaInfo.Ecgis = make([]Ecgi, 1)
	nwAreaInfo.Ecgis[0].EutraCellID = "EUTRACELL_01"
	nwAreaInfo.Ecgis[0].PlmnID.Mcc = "634"
	nwAreaInfo.Ecgis[0].PlmnID.Mnc = "635"
	nwAreaInfo.Ncgis = make([]Ncgi, 1)
	nwAreaInfo.Ncgis[0].NrCellID = "NRCELL_01"
	nwAreaInfo.Ncgis[0].PlmnID.Mcc = "834"
	nwAreaInfo.Ncgis[0].PlmnID.Mnc = "835"
	nwAreaInfo.GRanNodeIds = make([]GlobalRanNodeID, 1)
	nwAreaInfo.GRanNodeIds[0].PlmnID.Mcc = "934"
	nwAreaInfo.GRanNodeIds[0].PlmnID.Mnc = "935"
	nwAreaInfo.GRanNodeIds[0].N3IwfID = "IWF_01"
	nwAreaInfo.GRanNodeIds[0].GNbID.BitLength = 48
	nwAreaInfo.GRanNodeIds[0].GNbID.GNBValue = "GNB_01"
	nwAreaInfo.GRanNodeIds[0].NgeNbID = "NB_01"
	nwAreaInfo.Tais = make([]Tai, 1)
	nwAreaInfo.Tais[0].PlmnID.Mcc = "734"
	nwAreaInfo.Tais[0].PlmnID.Mnc = "735"
	nwAreaInfo.Tais[0].Tac = "TAC_01"

	return nil
}
