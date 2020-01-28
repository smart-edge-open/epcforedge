/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
