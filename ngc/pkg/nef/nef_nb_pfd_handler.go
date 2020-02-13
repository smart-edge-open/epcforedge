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
	"time"

	//"strconv"

	"github.com/gorilla/mux"
)

func createNewPFDTrans(nefCtx *nefContext, afID string,
	trans PfdManagement) (loc string, rsp map[string]nefPFDSBRspData,
	err error) {

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
	var rsp nefPFDSBRspData
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
			rsp1 := nefSBRspData{errorCode: rsp.result.errorCode}
			sendErrorResponseToAF(w, rsp1)
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
	nef := &nefCtx.nef

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
	pfdBody.PfdReports = make(map[string]PfdReport)

	// Validate the mandatory parameters
	resRsp, status := validateAFPfdManagementData(pfdBody)
	if !status {
		log.Err(resRsp.result.pd.Title)
		rsp1 := nefSBRspData{errorCode: resRsp.result.errorCode}
		sendErrorResponseToAF(w, rsp1)
		return
	}

	//Validation of parameters in PfdDatas where AppID is the key
	for key, pfdData := range pfdBody.PfdDatas {

		if len(pfdData.Pfds) == 0 {

			// prepare pfd report with OTHER REASON
			generatePfdReport(key, "OTHER_REASON", pfdBody.PfdReports)
			delete(pfdBody.PfdDatas, key)
		}

		// TBD Validate for Duplicate Application ID
		if nef.nefCheckPfdAppIDExists(key) {

			// Prepare Pfd report APP ID DUPLICATED
			generatePfdReport(key, "APP_ID_DUPLICATED", pfdBody.PfdReports)
			delete(pfdBody.PfdDatas, key)
		}

		for _, pfd := range pfdData.Pfds {
			rspPfd, status := validateAFPfdData(pfd)
			if !status {
				// Prepare Pfd report OTHER REASON
				// delete app and break
				_ = key //app ID
				_ = rspPfd
				//delete(pfdBody.PfdDatas, key)

			}
		}

	}

	// TBD check if all the apps are deleted then send 500 response with
	// pFD reports otherwise create

	loc, rsp, err3 := createNewPFDTrans(nefCtx, vars["scsAsId"], pfdBody)

	if err3 != nil {
		log.Err(err3)

		// TBD to update the PFD report on rsp
		_ = rsp

		// we return bad request here since we have reached the max
		rsp1 := nefSBRspData{errorCode: 400}
		/*
			if rsperrorCode == 500 {
				send500PFDResponseToAF(w, rsp, pfdBody.PfdReports)

			} else {

		*/

		sendErrorResponseToAF(w, rsp1)

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
	logNef(nef)

}

// ReadPFDManagementTransaction : Read a particular PFD transaction details
func ReadPFDManagementTransaction(w http.ResponseWriter, r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])
	log.Infof(" TRANSACTION ID  : %s", vars["transactionId"])

	af, ok := nef.nefGetAf(vars["scsAsId"])

	if ok != nil {
		sendCustomeErrorRspToAF(w, 404, "Failed to find AF records")
		return
	}

	rsp, pfdTrans, err := af.afGetPfdTransaction(nefCtx, vars["transactionId"])

	if err != nil {
		log.Err(err)
		rsp1 := nefSBRspData{errorCode: rsp.result.errorCode}
		sendErrorResponseToAF(w, rsp1)
		return
	}

	mdata, err2 := json.Marshal(pfdTrans)
	if err2 != nil {
		log.Err(err2)
		sendCustomeErrorRspToAF(w, 400, "Failed to Marshal GETPFDresponse data")
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

// ReadPFDManagementApplication : Read a particular PFD transaction details of
// and external application identifier
func ReadPFDManagementApplication(w http.ResponseWriter, r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])
	log.Infof(" PFD TRANSACTION ID  : %s", vars["transactionId"])
	log.Infof(" PFD APPLICATION ID : %s", vars["appId"])

	af, ok := nef.nefGetAf(vars["scsAsId"])

	if ok != nil {
		sendCustomeErrorRspToAF(w, 404, "Failed to find AF records")
		return
	}

	pfdTransID := vars["transactionId"]
	appID := vars["appId"]
	rsp, pfdData, err := af.afGetPfdApplication(nefCtx, pfdTransID, appID)

	if err != nil {
		log.Err(err)
		sendErrorResponseToAF(w, rsp)
		return
	}

	mdata, err2 := json.Marshal(pfdData)
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

// DeletePFDManagementApplication deletes the existing PFD transaction of
// an application identifier
func DeletePFDManagementApplication(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])
	log.Infof(" PFD TRANSACTION ID  : %s", vars["transactionId"])
	log.Infof(" PFD APPLICATION ID  : %s", vars["appId"])

	af, err := nef.nefGetAf(vars["scsAsId"])

	if err != nil {
		log.Err(err)
		sendCustomeErrorRspToAF(w, 404, "Failed to find AF entry")
		return
	}

	pfdTransID := vars["transactionId"]
	appID := vars["appId"]

	rsp, err := af.afDeletePfdApplication(nefCtx, pfdTransID, appID)

	if err != nil {
		log.Err(err)
		sendErrorResponseToAF(w, rsp)
		return
	}
	// Response should be 204 as per 3GPP 29.522
	w.WriteHeader(http.StatusNoContent)

	log.Infof("HTTP Response sent: %d", http.StatusNoContent)

	logNef(nef)
}

// DeletePFDManagementTransaction deletes the existing PFD transaction
func DeletePFDManagementTransaction(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])
	log.Infof(" PFD TRANSACTION ID  : %s", vars["transactionId"])

	af, err := nef.nefGetAf(vars["scsAsId"])

	if err != nil {
		log.Err(err)
		sendCustomeErrorRspToAF(w, 404, "Failed to find AF entry")
		return
	}
	rsp, err := af.afDeletePfdTransaction(nefCtx, vars["transactionId"])

	if err != nil {
		log.Err(err)
		rsp1 := nefSBRspData{errorCode: rsp.result.errorCode}
		sendErrorResponseToAF(w, rsp1)
		return
	}

	// Response should be 204 as per 3GPP 29.522
	w.WriteHeader(http.StatusNoContent)

	log.Infof("HTTP Response sent: %d", http.StatusNoContent)

	// If the AF subcount and transaction count is 0 delete the AF

	logNef(nef)
}

// UpdatePutPFDManagementTransaction updates an existing PFD transaction
func UpdatePutPFDManagementTransaction(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])
	log.Infof(" PFD TRANSACTION ID  : %s", vars["transactionId"])

	af, ok := nef.nefGetAf(vars["scsAsId"])
	if ok == nil {

		b, err := ioutil.ReadAll(r.Body)
		defer closeReqBody(r)

		if err != nil {
			log.Err(err)
			sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP PUT Body")
			return
		}

		//PFD Transaction data
		pfdTrans := PfdManagement{}

		//Convert the json Traffic Influence data into struct
		err1 := json.Unmarshal(b, &pfdTrans)

		if err1 != nil {
			log.Err(err1)
			sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal PUT data")
			return
		}

		rsp, newPfdTrans, err := af.afUpdatePutPfdTransaction(nefCtx,
			vars["transactionId"], pfdTrans)

		// TBD to update the pfd report on rsp
		_ = rsp
		if err != nil {
			// TBD how to send the error code
			rsp1 := nefSBRspData{errorCode: 400}
			sendErrorResponseToAF(w, rsp1)
			return
		}

		mdata, err2 := json.Marshal(newPfdTrans)

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

// UpdatePutPFDManagementApplication updates an existing PFD transaction
func UpdatePutPFDManagementApplication(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])
	log.Infof(" PFD TRANSACTION ID  : %s", vars["transactionId"])
	log.Infof(" PFD APPLICATION ID  : %s", vars["appId"])

	af, ok := nef.nefGetAf(vars["scsAsId"])
	if ok == nil {

		b, err := ioutil.ReadAll(r.Body)
		defer closeReqBody(r)

		if err != nil {
			log.Err(err)
			sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP PUT Body")
			return
		}

		//PFD Transaction data
		pfdData := PfdData{}

		//Convert the json PFD Management data into struct
		err1 := json.Unmarshal(b, &pfdData)

		if err1 != nil {
			log.Err(err1)
			sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal PUT data")
			return
		}

		rsp, newPfdData, err := af.afUpdatePutPfdApplication(nefCtx,
			vars["transactionId"], vars["appId"], pfdData)

		if err != nil {

			rsp1 := nefSBRspData{errorCode: rsp.result.errorCode}
			sendErrorResponseToAF(w, rsp1)
			return
		}

		mdata, err2 := json.Marshal(newPfdData)

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

// PatchPFDManagementApplication patches the PFD application PFDs
func PatchPFDManagementApplication(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(nefCtxKey("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["scsAsId"])
	log.Infof(" PFD TRANSACTION ID  : %s", vars["transactionId"])
	log.Infof(" PFD APPLICATION ID  : %s", vars["appId"])

	af, ok := nef.nefGetAf(vars["scsAsId"])
	if ok == nil {

		b, err := ioutil.ReadAll(r.Body)
		defer closeReqBody(r)

		if err != nil {
			log.Err(err)
			sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP PUT Body")
			return
		}

		//PFD Transaction data
		pfdData := PfdData{}

		//Convert the json PFD Management data into struct
		err1 := json.Unmarshal(b, &pfdData)

		if err1 != nil {
			log.Err(err1)
			sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal PUT data")
			return
		}

		rsp, newPfdData, err := af.afUpdatePatchPfdApplication(nefCtx,
			vars["transactionId"], vars["appId"], pfdData)

		if err != nil {

			rsp1 := nefSBRspData{errorCode: rsp.result.errorCode}
			sendErrorResponseToAF(w, rsp1)

			return
		}

		mdata, err2 := json.Marshal(newPfdData)

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

//PFD Management functions

func (af *afData) afUpdatePutPfdApplication(nefCtx *nefContext, transID string,
	appID string, pfdData PfdData) (rsp nefPFDSBRspData, updPfd PfdData,
	err error) {

	pfdTrans, ok := af.pfdtrans[transID]

	if !ok {
		rsp.result.errorCode = 400
		rsp.result.pd.Title = pfdNotFound

		return rsp, updPfd, errors.New(pfdNotFound)
	}
	trans, ok := pfdTrans.pfdManagement.PfdDatas[appID]

	if !ok {
		rsp.result.errorCode = 404
		rsp.result.pd.Title = appNotFound
		return rsp, trans, errors.New(appNotFound)
	}

	rsp, err = pfdTrans.NEFSBAppPfdPut(pfdTrans, nefCtx, pfdData)

	if err != nil {
		log.Err("Failed to Update the PFD Application")
		return rsp, pfdData, err
	}

	pfdTrans.pfdManagement.PfdDatas[appID] = pfdData
	updPfd = pfdData

	log.Infoln("Update PFD transaction Successful")
	return rsp, updPfd, err
}

func (af *afData) afUpdatePatchPfdApplication(nefCtx *nefContext,
	transID string, appID string, pfdData PfdData) (rsp nefPFDSBRspData,
	updPfd PfdData, err error) {

	pfdTrans, ok := af.pfdtrans[transID]

	if !ok {
		rsp.result.errorCode = 400
		rsp.result.pd.Title = pfdNotFound

		return rsp, updPfd, errors.New(pfdNotFound)
	}

	/*rsp, err = sub.NEFSBPut(sub, nefCtx, ti)

	if err != nil {
		log.Err("Failed to Update Subscription")
		return rsp, updtTI, err
	}*/

	trans, ok := pfdTrans.pfdManagement.PfdDatas[appID]

	if !ok {
		rsp.result.errorCode = 404
		rsp.result.pd.Title = appNotFound
		return rsp, trans, errors.New(appNotFound)

	}

	rsp, err = pfdTrans.NEFSBAppPfdPut(pfdTrans, nefCtx, pfdData)

	if err != nil {
		log.Err("Failed to Update the PFD Application")
		return rsp, pfdData, err
	}
	// Updating the PFDs present in the
	for key := range pfdData.Pfds {

		pfd, ok := trans.Pfds[key]
		if ok {
			pfd = pfdData.Pfds[key]
			trans.Pfds[key] = pfd
			log.Infof("PFD id %s updated by PATCH ", pfd.PfdID)
		}

	}

	updPfd = trans
	log.Infoln("Patch PFD Application PFDs Successful")
	return rsp, updPfd, err
}

func (af *afData) afUpdatePutPfdTransaction(nefCtx *nefContext, transID string,
	trans PfdManagement) (rsp map[string]nefPFDSBRspData, updPfd PfdManagement,
	err error) {

	pfdTrans, ok := af.pfdtrans[transID]

	if !ok {

		// TBD how to send the error code
		//rsp.result.errorCode = 400
		//rsp.result.pd.Title = pfdNotFound

		return rsp, updPfd, errors.New(pfdNotFound)
	}

	rsp, err = pfdTrans.NEFSBPfdPut(pfdTrans, nefCtx, trans)

	if err != nil {

		//Return error
		return rsp, trans, err
	}

	/*rsp, err = sub.NEFSBPut(sub, nefCtx, ti)

	if err != nil {
		log.Err("Failed to Update Subscription")
		return rsp, updtTI, err
	}*/

	updPfd = trans
	updPfd.Self = pfdTrans.pfdManagement.Self
	pfdTrans.pfdManagement = updPfd

	log.Infoln("Update PFD transaction Successful")
	return rsp, updPfd, err
}

func (af *afData) afDeletePfdTransaction(nefCtx *nefContext,
	pfdTrans string) (rsp nefPFDSBRspData, err error) {

	//Check if PFD transaction is already present
	trans, ok := af.pfdtrans[pfdTrans]

	if !ok {
		rsp.result.errorCode = 404
		rsp.result.pd.Title = pfdNotFound
		return rsp, errors.New(pfdNotFound)
	}

	rsp, err = trans.NEFSBPfdDelete(trans, nefCtx)
	if err != nil {
		log.Err("Failed to Delete PFD transaction")
		rsp.result.errorCode = 400
		return rsp, err
	}

	//Delete local entry in map of pfd transactions
	delete(af.pfdtrans, pfdTrans)

	// TBD check if all trans and sub deleted for AF then delete AF

	return rsp, err
}

func (af *afData) afGetPfdApplication(nefCtx *nefContext,
	transID string, appID string) (rsp nefSBRspData, trans PfdData, err error) {

	transPfd, ok := af.pfdtrans[transID]

	if !ok {
		rsp.errorCode = 404
		rsp.pd.Title = pfdNotFound
		return rsp, trans, errors.New(pfdNotFound)
	}

	trans, ok = transPfd.pfdManagement.PfdDatas[appID]

	if !ok {
		rsp.errorCode = 404
		rsp.pd.Title = appNotFound
		return rsp, trans, errors.New(appNotFound)
	}

	//Return locally
	return rsp, trans, err
}

func (af *afData) afDeletePfdApplication(nefCtx *nefContext,
	transID string, appID string) (rsp nefSBRspData, err error) {

	transPfd, ok := af.pfdtrans[transID]

	if !ok {
		rsp.errorCode = 404
		rsp.pd.Title = pfdNotFound
		return rsp, errors.New(pfdNotFound)
	}

	_, ok = transPfd.pfdManagement.PfdDatas[appID]

	if !ok {
		rsp.errorCode = 404
		rsp.pd.Title = appNotFound
		return rsp, errors.New(appNotFound)
	}

	delete(transPfd.pfdManagement.PfdDatas, appID)

	// TBD check if all app deleted for trans then delete trans
	// check if all trans and sub deleted for AF, delete AF
	//Return locally
	return rsp, err
}

func (af *afData) afGetPfdTransaction(nefCtx *nefContext,
	transID string) (rsp nefPFDSBRspData, trans PfdManagement, err error) {

	transPfd, ok := af.pfdtrans[transID]

	if !ok {
		rsp.result.errorCode = 404
		rsp.result.pd.Title = pfdNotFound
		return rsp, trans, errors.New(pfdNotFound)
	}

	_, rsp, err = transPfd.NEFSBPfdGet(transPfd, nefCtx)
	if err != nil {
		log.Infoln("Failed to Get PFD transaction")
		return rsp, transPfd.pfdManagement, err
	}

	//Return locally
	return rsp, transPfd.pfdManagement, err
}

func (af *afData) afGetPfdTransactionList(nefCtx *nefContext) (
	rsp nefPFDSBRspData, transList []PfdManagement, err error) {

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
	trans PfdManagement) (loc string, rsp map[string]nefPFDSBRspData,
	err error) {

	/*Check if max subscription reached */
	if len(af.pfdtrans) >= nefCtx.cfg.MaxPfdTransSupport {

		//rsp.errorCode = 400
		//rsp.pd.Title = "MAX Transaction Reached"
		return "", rsp, errors.New("MAX TRANS Created")
	}
	//Generate a unique transaction ID string
	transIDStr := strconv.Itoa(af.transIDnum)
	af.transIDnum++

	//Create PFD transaction data
	aftrans := afPfdTransaction{transID: transIDStr, pfdManagement: trans}

	aftrans.NEFSBPfdGet = nefSBUDRPFDGet
	aftrans.NEFSBAppPfdPut = nefSBUDRAPPPFDPut
	aftrans.NEFSBPfdPut = nefSBUDRPFDPut
	aftrans.NEFSBPfdDelete = nefSBUDRPFDDelete

	rsp, err = aftrans.NEFSBPfdPut(&aftrans, nefCtx, trans)

	if err != nil {

		//Return error
		return "", rsp, err
	}

	//Store Notification Destination URI
	//afsub.afNotificationDestination = ti.NotificationDestination

	//Link the PFD transaction with the AF
	af.pfdtrans[transIDStr] = &aftrans

	//Create Location URI
	loc = nefCtx.nef.locationURLPrefixPfd + af.afID + "/transactions/" +
		transIDStr

	af.pfdtrans[transIDStr].pfdManagement.Self = Link(loc)

	//Also update the self link in each application
	for k, v := range af.pfdtrans[transIDStr].pfdManagement.PfdDatas {

		/*Assign the application ID in the link */
		v.Self = Link(loc) + "/applications/" + Link(k)
		log.Infof("Application ID is %s", k)
		af.pfdtrans[transIDStr].pfdManagement.PfdDatas[k] = v

	}

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

// validateAFPfdManagementData Function to validate mandatory parameters of
// PFD Management received from AF
func validateAFPfdManagementData(pfdTrans PfdManagement) (rsp nefPFDSBRspData,
	status bool) {

	if len(pfdTrans.PfdDatas) == 0 {
		rsp.result.errorCode = 400
		rsp.result.pd.Title = "Missing PFD Data"
		return rsp, false
	}

	//Validation of Individual PFD
	for _, v := range pfdTrans.PfdDatas {
		if len(v.Pfds) == 0 {
			rsp.result.errorCode = 400
			rsp.result.pd.Title = "Missing PFD Data"
			return rsp, false
		}

		if len(v.ExternalAppID) == 0 {
			rsp.result.errorCode = 400
			rsp.result.pd.Title = "Missing Application ID"
			return rsp, false
		}

	}
	return rsp, true
}

// validateAFPfdData Function to validate mandatory parameters of
// PFD  received from AF
func validateAFPfdData(pfd Pfd) (rsp nefPFDSBRspData,
	status bool) {

	if len(pfd.PfdID) == 0 {
		rsp.result.errorCode = 400
		rsp.result.pd.Title = "PFD ID missing"
		return rsp, false
	}
	if len(pfd.DomainNames) == 0 && len(pfd.FlowDescriptions) == 0 &&
		len(pfd.Urls) == 0 {
		rsp.result.errorCode = 400
		rsp.result.pd.Title = "No domainNames" +
			"FlowDescriptions and Urls present for PFD"
		return rsp, false
	}
	rspPfd, result := validatePfd(pfd)
	return rspPfd, result
}

// validateAFPfd Function to validate mandatory parameters of
// PFD
func validatePfd(pfd Pfd) (rsp nefPFDSBRspData,
	status bool) {
	if len(pfd.DomainNames) != 0 {
		if len(pfd.FlowDescriptions) != 0 || len(pfd.Urls) != 0 {
			rsp.result.errorCode = 400
			rsp.result.pd.Title = "only one of domainNames" +
				"FlowDescriptions and Urls present for PFD"
			return rsp, false
		}
	}
	if len(pfd.Urls) != 0 {
		if len(pfd.FlowDescriptions) != 0 || len(pfd.DomainNames) != 0 {
			rsp.result.errorCode = 400
			rsp.result.pd.Title = "only one of domainNames" +
				"FlowDescriptions and Urls present for PFD"
			return rsp, false
		}
	}
	if len(pfd.FlowDescriptions) != 0 {
		if len(pfd.Urls) != 0 || len(pfd.DomainNames) != 0 {
			rsp.result.errorCode = 400
			rsp.result.pd.Title = "only one of domainNames" +
				"FlowDescriptions and Urls present for PFD"
			return rsp, false
		}
	}

	return rsp, true
}

func nefSBUDRAPPPFDGet(transData *afPfdTransaction, nefCtx *nefContext,
	appID ApplicationID) (appPfd PfdData, rsp nefPFDSBRspData, err error) {

	log.Info("nefSBUDRAPPPFDGet Entered ")
	nef := &nefCtx.nef

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	r, e := nef.udrPfdClient.UdrPfdDataGet(cliCtx, UdrAppID(appID))
	if e != nil {

		return appPfd, rsp, e
	}
	appPfd.ExternalAppID = string(r.AppPfd.AppID)

	for _, ele := range r.AppPfd.Pfds {
		var conData = Pfd{}

		conData.PfdID = ele.PfdID
		conData.DomainNames = ele.DomainNames
		conData.FlowDescriptions = ele.FlowDescriptions
		conData.Urls = ele.Urls

		appPfd.Pfds[ele.PfdID] = conData
	}
	log.Info("nefSBUDRAPPPFDGet Exited ")

	return appPfd, rsp, err
}

//NEFSBGetPfdFn is the callback for SB API to get PFD transaction
func nefSBUDRPFDGet(transData *afPfdTransaction, nefCtx *nefContext) (
	trans PfdManagement, rsp nefPFDSBRspData, err error) {

	for k, v := range transData.pfdManagement.PfdDatas {
		appPfd, r, e := nefSBUDRAPPPFDGet(transData, nefCtx,
			ApplicationID(v.ExternalAppID))
		if e != nil {
			return trans, r, e
		}
		trans.PfdDatas[k] = appPfd
	}

	return trans, rsp, nil
}

//NEFSBPutPfdFn is the callback for SB API to put PFD transaction
func nefSBUDRAPPPFDPut(transData *afPfdTransaction, nefCtx *nefContext,
	app PfdData) (rsp nefPFDSBRspData, err error) {

	nef := &nefCtx.nef
	var pfdApp PfdDataForApp

	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	/*Timer for pfdApp.AllowedDelay can be started here, not supported
	currently */
	pfdApp.AppID = ApplicationID(app.ExternalAppID)

	log.Info("nefSBUDRAPPPFDPut ->  ")
	if app.CachingTime != nil {

		i := time.Duration(*app.CachingTime)
		timeLater := DateTime(time.Now().Add(time.Second * i).String())
		pfdApp.CachingTime = &timeLater
	}

	for _, pfdv := range app.Pfds {

		var c PfdContent
		c.PfdID = pfdv.PfdID
		c.DomainNames = pfdv.DomainNames
		c.FlowDescriptions = pfdv.FlowDescriptions
		c.Urls = pfdv.Urls
		pfdApp.Pfds = append(pfdApp.Pfds, c)
	}

	_, e := nef.udrPfdClient.UdrPfdDataCreate(cliCtx, pfdApp)

	if e != nil {
		rsp.result.errorCode = 400
		return rsp, e
	}

	rsp.result.errorCode = 200

	return rsp, err

}

// nefSBUDRPFDPut is the callback for SB API to put PFD transaction
func nefSBUDRPFDPut(transData *afPfdTransaction, nefCtx *nefContext,
	trans PfdManagement) (rsp map[string]nefPFDSBRspData, err error) {

	rspDetails := make(map[string]nefPFDSBRspData)

	for k, v := range trans.PfdDatas {

		r, e := nefSBUDRAPPPFDPut(transData, nefCtx, v)
		if e != nil {
			//fatal error return
			return rspDetails, e
		}
		/*Update the map with the response*/
		rspDetails[k] = r
	}
	return rspDetails, nil
}

//NEFSBDeletePfdFn is the callback for SB API to delete PFD transaction
func nefSBUDRAPPPFDDelete(transData *afPfdTransaction, nefCtx *nefContext,
	appID ApplicationID) (rsp nefPFDSBRspData, err error) {

	nef := &nefCtx.nef
	cliCtx, cancel := context.WithCancel(nef.ctx)
	defer cancel()

	_, e := nef.udrPfdClient.UdrPfdDataDelete(cliCtx, UdrAppID(appID))

	if e != nil {
		return rsp, e
	}

	return rsp, nil
}

//NEFSBDeletePfdFn is the callback for SB API to delete PFD transaction
func nefSBUDRPFDDelete(transData *afPfdTransaction, nefCtx *nefContext) (
	rsp nefPFDSBRspData, err error) {

	for _, v := range transData.pfdManagement.PfdDatas {
		r, e := nefSBUDRAPPPFDDelete(transData, nefCtx,
			ApplicationID(v.ExternalAppID))
		if e != nil {
			return r, e
		}
	}
	return rsp, nil
}
