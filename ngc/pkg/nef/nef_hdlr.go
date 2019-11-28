// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

const correlationIDOffset = 20
const subNotFound string = "Subscription Not Found"

//NEF context data
type nefData struct {
	afcount   int
	apiRoot   string
	pcfClient PcfPolicyAuthorization
	udrClient UdrInfluenceData
	corrID    uint
	afs       map[string]*afData
}

//NEFSBGetFn is the callback for SB API
type NEFSBGetFn func(subData *afSubscription, nefCtx *nefContext) (
	sub TrafficInfluSub, rsp nefSBRspData, err error)

//NEFSBPutFn is the callback for SB API
type NEFSBPutFn func(subData *afSubscription, nefCtx *nefContext,
	ti TrafficInfluSub) (rsp nefSBRspData, err error)

//NEFSBPatchFn is the callback for SB API
type NEFSBPatchFn func(subData *afSubscription, nefCtx *nefContext,
	tisp TrafficInfluSubPatch) (rsp nefSBRspData, err error)

//NEFSBDeleteFn is the callback for SB API
type NEFSBDeleteFn func(subData *afSubscription, nefCtx *nefContext) (
	rsp nefSBRspData, err error)

//PCF Subscription data
type afSubscription struct {
	subid string
	ti    TrafficInfluSub

	//Applicable in case of single UE case only
	appSessionID AppSessionID

	//Applicable in case of Multiple UE only
	iid                       InfluenceID
	NotifCorreID              string
	afNotificationDestination Link
	NEFSBGet                  NEFSBGetFn
	NEFSBPut                  NEFSBPutFn
	NEFSBPatch                NEFSBPatchFn
	NEFSBDelete               NEFSBDeleteFn
}

//AF data
type afData struct {
	afid       string
	subIdnum   int
	maxSubSupp int
	subs       map[string]*afSubscription
}

type nefSBRspData struct {
	errorCode int
	pd        ProblemDetails
}

//Creates a AF instance
func (af *afData) afCreate(nefCtx *nefContext, afid string) error {

	//Validate afid ??

	af.afid = afid
	af.subIdnum = nefCtx.cfg.SubStartID //Start Number
	af.maxSubSupp = nefCtx.cfg.MaxSubSupport
	af.subs = make(map[string]*afSubscription)
	return nil
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
	subIDStr := strconv.Itoa(af.subIdnum)
	af.subIdnum++

	//Create Subscription data
	afsub := afSubscription{subid: subIDStr, ti: ti, appSessionID: "",
		NotifCorreID: "", iid: ""}

	if !ti.AnyUeInd {

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

	} else {
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

	}

	//Link the subscription with the AF
	af.subs[subIDStr] = &afsub

	//Create Location URI
	loc = nefCtx.nef.apiRoot + nefCtx.cfg.LocationPrefix + af.afid +
		"/subscriptions/" + subIDStr

	log.Infoln(" NEW AF Subscription added " + subIDStr)

	return loc, rsp, nil
}

func (af *afData) afUpdateSubscription(nefCtx *nefContext, subID string,
	ti TrafficInfluSub) (rsp nefSBRspData, err error) {

	sub, ok := af.subs[subID]

	if !ok {
		rsp.errorCode = 400
		rsp.pd.Title = subNotFound

		return rsp, errors.New(subNotFound)
	}

	rsp, err = sub.NEFSBPut(sub, nefCtx, ti)

	if err != nil {
		log.Infoln("Failed to Update Subscription")
		return rsp, err
	}
	sub.ti = ti

	log.Infoln("Update Subscription Successful")
	return rsp, err
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
		log.Infoln("Failed to Patch Subscription")
		return rsp, ti, err
	}
	updateTiFromTisp(&sub.ti, tisp)

	return rsp, sub.ti, err

}

func (af *afData) afGetSubscription(nefCtx *nefContext,
	subID string) (rsp nefSBRspData, ti TrafficInfluSub, err error) {

	sub, ok := af.subs[subID]

	if !ok {
		rsp.errorCode = 400
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
	subslist []TrafficInfluSub, err error) {

	var ti TrafficInfluSub

	if len(af.subs) > 0 {

		for key := range af.subs {

			rsp, ti, err = af.afGetSubscription(nefCtx, key)

			if err != nil {
				return rsp, subslist, err
			}
			subslist = append(subslist, ti)
		}
	}
	return rsp, subslist, err
}

func (af *afData) afDeleteSubscription(nefCtx *nefContext,
	subID string) (rsp nefSBRspData, err error) {

	//Check if AF is already present
	sub, ok := af.subs[subID]

	if !ok {
		rsp.errorCode = 400
		rsp.pd.Title = subNotFound
		return rsp, errors.New(subNotFound)
	}

	rsp, err = sub.NEFSBDelete(sub, nefCtx)

	if err != nil {
		log.Infoln("Failed to Delete Subscription")
		return rsp, err
	}

	//Delete local entry in map
	delete(af.subs, subID)
	af.subIdnum--

	return rsp, err
}

/* unused function
func (af *afData) afDestroy(afid string) error {

	//Todo delete all subscriptions, needed in go ??
	//Needed for gracefully disconnecting
	return errors.New("AF data cleaned")
}
*/

//Initialize the NEF component
func (nef *nefData) nefCreate(ctx context.Context, cfg Config) error {

	_ = ctx
	nef.afcount = 0
	nef.pcfClient = NewPCFClient(nil)
	nef.udrClient = NewUDRClient(nil)
	nef.afs = make(map[string]*afData)
	nef.corrID = uint(cfg.SubStartID + correlationIDOffset)
	return nil
}

/*
func NEFInit() error {

	return nef.nefCreate()
}*/

func (nef *nefData) nefAddAf(nefCtx *nefContext, afID string) (af *afData,
	err error) {

	var afe afData

	//Check if AF is already present
	_, ok := nef.afs[afID]

	if !ok {
		return nef.afs[afID], errors.New("AF already present")
	}

	//Create a new entry of AF

	_ = afe.afCreate(nefCtx, afID)
	nef.afs[afID] = &afe
	nef.afcount++

	return &afe, nil
}

func (nef *nefData) nefGetAf(afID string) (af *afData, err error) {

	//Check if AF is already present
	afe, ok := nef.afs[afID]

	if !ok {
		return afe, nil
	}
	err = errors.New("AF entry not present")
	return afe, err
}

/* unused function
func (nef *nefData) nefDeleteAf(afID string) (err error) {

	//Check if AF is already present
	_, ok := nef.afs[afID]

	if !ok {
		delete(nef.afs, afID)
		nef.afcount--
		return nil
	}

	err = errors.New("AF entry not present")
	return err
}
*/

func (nef *nefData) nefDestroy() {

	// Todo
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

	cliCtx, cancel := context.WithCancel(context.Background())
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
			URI("https://localhost" + nefCtx.cfg.HTTP2Config.Endpoint)
	} else {
		appSessCtx.AscReqData.AfRoutReq.UpPathChgSub.NotificationURI =
			URI("http://localhost" + nefCtx.cfg.HTTPConfig.Endpoint)
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
		nef.pcfClient.PcfPolicyAuthorizationCreate(cliCtx, appSessCtx)

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

	cliCtx, cancel := context.WithCancel(context.Background())
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

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appSessCtxUpdtData := AppSessionContextUpdateData{}

	//Populating App Session Context Data Req
	//TODO: appSessCtxUpdtData.AfAppID = ti.AfAppID
	appSessCtxUpdtData.AfRoutReq.AppReloc = tisp.AppReloInd

	//Populating UP Path Chnage Subbscription Data in App Session Context
	//TODO: appSessCtxUpdtData.AfRoutReq.UpPathChgSub.DnaiChgType =
	//         ti.DnaiChgType
	// If HTTP2 port is configured use it else HTTP port
	if nefCtx.cfg.HTTP2Config.Endpoint != "" {
		appSessCtxUpdtData.AfRoutReq.UpPathChgSub.NotificationURI =
			URI("https://localhost" + nefCtx.cfg.HTTP2Config.Endpoint)
	} else {
		appSessCtxUpdtData.AfRoutReq.UpPathChgSub.NotificationURI =
			URI("http://localhost" + nefCtx.cfg.HTTPConfig.Endpoint)
	}
	/*+ nefCtx.cfg.UpfNotificationResUriPath*/
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

	cliCtx, cancel := context.WithCancel(context.Background())
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

// nefSBUDRPost : This function returns error as HTTP POST Request to UDR is not
//            supported.
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

	cliCtx, cancel := context.WithCancel(context.Background())
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

	cliCtx, cancel := context.WithCancel(context.Background())
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
	// If HTTP2 port is configured use it else HTTP port
	if nefCtx.cfg.HTTP2Config.Endpoint != "" {
		trafficInfluData.UpPathChgNotifURI =
			URI("https://localhost" + nefCtx.cfg.HTTP2Config.Endpoint)
	} else {
		trafficInfluData.UpPathChgNotifURI =
			URI("http://localhost" + nefCtx.cfg.HTTPConfig.Endpoint)
	}
	/*+ nefCtx.cfg.UpfNotificationResUriPath*/
	if len(string(ti.SubscribedEvents[0])) > 0 &&
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

	cliCtx, cancel := context.WithCancel(context.Background())
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

	cliCtx, cancel := context.WithCancel(context.Background())
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
