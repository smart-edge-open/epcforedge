/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	//"strconv"

	"github.com/gorilla/mux"
)

func createNewPFDTrans(nefCtx *nefContext, afID string,
	trans PfdManagement) (loc string, rsp nefSBRspData, err error) {

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

	loc, rsp, err = af.afAddPFDTransaction(nefCtx, trans)

	if err != nil {
		return loc, rsp, err
	}

	return loc, rsp, nil
}

// ReadAllPFDManagementTransaction : API to read all the PFD Transactions
func ReadAllPFDManagementTransaction(w http.ResponseWriter,
	r *http.Request) {

	var pfdTrans []PfdManagement
	var rsp nefSBRspData
	var err error

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID : %s", vars["scsAsId"])

	af, err := nef.nefGetAf(vars["scsAsId"])

	if err != nil {
		/* Failure in getting AF with afId received. In this case no
		 * transaction data will be returned to AF */
		log.Infoln(err)
	} else {
		rsp, pfdTrans, err = af.afGetPfdTransactionList(nefCtx)
		if err != nil {
			log.Err(err)
			sendErrorResponseToAF(w, rsp)
			return
		}
	}

	mdata, err2 := json.Marshal(pfdTrans)

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

// CreatePFDManagementTransaction  Handles the PFD Management requested
// by AF
func CreatePFDManagementTransaction(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])

	b, err := ioutil.ReadAll(r.Body)
	defer closeReqBody(r)

	if err != nil {
		sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP POST Body")
		return
	}

	//Pfd Management data
	pfdBody := PfdManagement{}

	//Convert the json Traffic Influence data into struct
	err1 := json.Unmarshal(b, &pfdBody)

	if err1 != nil {
		log.Err(err1)
		sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal POST data")
		return
	}
	/*
		//validate the mandatory parameters
		resRsp, status := validateAFTrafficInfluenceData(trInBody)
		if !status {
			log.Err(resRsp.pd.Title)
			sendErrorResponseToAF(w, resRsp)
			return
		}
	*/
	loc, rsp, err3 := createNewPFDTrans(nefCtx, vars["scsAsId"], pfdBody)

	if err3 != nil {
		log.Err(err3)
		// we return bad request here since we have reached the max
		rsp.errorCode = 400
		sendErrorResponseToAF(w, rsp)
		return
	}
	log.Infoln(loc)

	pfdBody.Self = Link(loc)

	//Martshal data and send into the body
	mdata, err2 := json.Marshal(pfdBody)

	if err2 != nil {
		log.Err(err2)
		sendCustomeErrorRspToAF(w, 400, "Failed to Marshal GET response data")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Location", loc)

	// Response should be 201 Created as per 3GPP 29.522
	w.WriteHeader(http.StatusCreated)
	log.Infof("CreatePFDManagementresponses => %d",
		http.StatusCreated)
	_, err = w.Write(mdata)
	if err != nil {
		log.Errf("Write Failed: %v", err)
		return
	}
	nef := &nefCtx.nef
	logNef(nef)

}

//PFD Management functions
func (af *afData) afGetPfdTransaction(nefCtx *nefContext,
	transID string) (rsp nefSBRspData, trans PfdManagement, err error) {

	transPfd, ok := af.pfdtrans[transID]

	if !ok {
		rsp.errorCode = 404
		rsp.pd.Title = subNotFound
		return rsp, trans, errors.New(subNotFound)
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
	return rsp, transPfd.pfdManagement, err
}

func (af *afData) afGetPfdTransactionList(nefCtx *nefContext) (rsp nefSBRspData,
	transList []PfdManagement, err error) {

	var transPfd PfdManagement

	if len(af.pfdtrans) > 0 {

		for key := range af.pfdtrans {

			rsp, transPfd, err = af.afGetPfdTransaction(nefCtx, key)

			if err != nil {
				return rsp, transList, err
			}
			transList = append(transList, transPfd)
		}
	}
	return rsp, transList, err
}

//Creates a new subscription
func (af *afData) afAddPFDTransaction(nefCtx *nefContext,
	trans PfdManagement) (loc string, rsp nefSBRspData, err error) {

	/*Check if max subscription reached */
	if len(af.pfdtrans) >= nefCtx.cfg.MaxSubSupport {

		rsp.errorCode = 400
		rsp.pd.Title = "MAX Transaction Reached"
		return "", rsp, errors.New("MAX TRANS Created")
	}
	//Generate a unique subscription ID string
	transIDStr := strconv.Itoa(af.transIDnum)
	af.transIDnum++

	//Create PFD transaction data
	aftrans := afPfdTransaction{transID: transIDStr, pfdManagement: trans}

	/*rsp, err = nefSBUDRPost(&afsub, nefCtx, ti)

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
	*/
	//Link the subscription with the AF
	af.pfdtrans[transIDStr] = &aftrans

	//Create Location URI
	loc = nefCtx.nef.locationURLPrefixPfd + af.afID + "/transactions/" +
		transIDStr

	aftrans.pfdManagement.Self = Link(loc)

	log.Infoln(" NEW AF PFD transaction added " + transIDStr)

	return loc, rsp, nil
}

// Generate the notification uri for PFD
func getNefLocationURLPrefixPfd(cfg *Config) string {

	var uri string
	// If http2 port is configured use it else http port
	if cfg.HTTP2Config.Endpoint != "" {
		uri = "https://" + cfg.NefAPIRoot +
			cfg.HTTP2Config.Endpoint
	} else {
		uri = "http://" + cfg.NefAPIRoot +
			cfg.HTTPConfig.Endpoint
	}
	uri += cfg.LocationPrefixPfd
	return uri

}
