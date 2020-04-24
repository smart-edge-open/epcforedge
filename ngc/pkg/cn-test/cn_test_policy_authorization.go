/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngccntest

import (
	//"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

var gASCStartID = 5000

// PolicyAuthData : Policy Authorization data
type PolicyAuthData struct {
	asc       map[string]*AppSessionContext
	ascID     int
	locPrefix string
	ascMax    int
}

//IntPolicyAuthorization Init Policy Authorization variables
func IntPolicyAuthorization(cfg Config) {

	NgcData.paData.locPrefix = getLocationURLPrefix(&cfg)
	NgcData.paData.asc = make(map[string]*AppSessionContext)
	NgcData.paData.ascMax = cfg.MaxASCSupport
}

//PolicyAuthorizationAppSessionCreate Post
func PolicyAuthorizationAppSessionCreate(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionCreate -- Entered")

	asc := AppSessionContext{}

	fmt.Println(" len(NgcData.paData.asc) ", len(NgcData.paData.asc))
	fmt.Println("NgcData.paData.ascMax", NgcData.paData.ascMax)
	/*Check if max subscription reached */
	if len(NgcData.paData.asc) > NgcData.paData.ascMax {
		log.Info("Maximum Context creation reached")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Err(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer closeReqBody(r)
	err1 := json.Unmarshal(b, &asc)
	if err1 != nil {
		log.Err(err1)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(asc)

	AscRespData := getAppSessionContextRespData()
	EvsNotif := getEventsNotification()
	EvsNotif.EvSubsURI = asc.AscReqData.EvSubsc.NotifURI

	asc.AscRespData = &AscRespData
	asc.EvsNotif = &EvsNotif

	fmt.Println(asc)

	loc, ascID := genLocURI()
	fmt.Println("Location Header", loc)

	mdata, err2 := json.Marshal(asc)
	if err2 != nil {
		log.Errf("Write Failed: %v", err2)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*Update the data in map*/
	NgcData.paData.asc[ascID] = &asc

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusCreated)

	//Send Success response to AF
	_, err = w.Write(mdata)

	if err != nil {
		log.Errf("Write Failed: %v", err)
		return
	}
	log.Infof("HTTP Response sent: %d", http.StatusCreated)

}

//PolicyAuthorizationAppSessionGet Get
func PolicyAuthorizationAppSessionGet(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionGet -- Entered")
	vars := mux.Vars(r)
	ascID := vars["appSessionId"]
	log.Infof(" APP Session ID  : %s", ascID)

	var err error
	if len(NgcData.paData.asc) > 0 {

		if NgcData.paData.asc[ascID] != nil {
			mdata, err2 := json.Marshal(NgcData.paData.asc[ascID])
			if err2 != nil {
				log.Errf("Write Failed: %v", err2)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			w.WriteHeader(http.StatusOK)

			//Send Success response to AF
			_, err = w.Write(mdata)
			if err != nil {
				log.Errf("Write Failed: %v", err)
				return
			}
			log.Infof("HTTP Response sent: %d", http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func patchAsc(ascp *AppSessionContextUpdateData, AscReqData *AppSessionContextReqData) {

	if len(ascp.AfAppID) > 0 {
		AscReqData.AfAppID = ascp.AfAppID
	}
	if ascp.AfRoutReq != nil {
		AscReqData.AfRoutReq = ascp.AfRoutReq
	}
	if len(ascp.AspID) > 0 {
		AscReqData.AspID = ascp.AspID
	}
	if len(ascp.BdtRefID) > 0 {
		AscReqData.BdtRefID = ascp.BdtRefID
	}
	if ascp.EvSubsc != nil {
		AscReqData.EvSubsc = ascp.EvSubsc
	}
	if ascp.MedComponents != nil {
		AscReqData.MedComponents = ascp.MedComponents
	}
	if len(ascp.MpsID) > 0 {
		AscReqData.MpsID = ascp.MpsID
	}
	if len(ascp.SponID) > 0 {
		AscReqData.SponID = ascp.SponID
	}
	if len(ascp.SponStatus) > 0 {
		AscReqData.SponStatus = ascp.SponStatus
	}

}

//PolicyAuthorizationAppSessionPatch Patch
func PolicyAuthorizationAppSessionPatch(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSession -- Entered")
	vars := mux.Vars(r)

	ascID := vars["appSessionId"]
	log.Infof(" APP Session ID  : %s", ascID)

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Err(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer closeReqBody(r)

	ascp := AppSessionContextUpdateData{}

	err1 := json.Unmarshal(b, &ascp)
	if err1 != nil {
		log.Err(err1)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(NgcData.paData.asc) > 0 {

		if NgcData.paData.asc[ascID] != nil {

			asc := NgcData.paData.asc[ascID]

			patchAsc(&ascp, asc.AscReqData)

			mdata, err2 := json.Marshal(asc)
			if err2 != nil {
				log.Errf("Write Failed: %v", err2)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			w.WriteHeader(http.StatusOK)

			//Send Success response to AF
			_, err = w.Write(mdata)
			if err != nil {
				log.Errf("Write Failed: %v", err)
				return
			}
			log.Infof("HTTP Response sent: %d", http.StatusCreated)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)

}

//PolicyAuthorizationAppSessionDelete Delete
func PolicyAuthorizationAppSessionDelete(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionDelete -- Entered")
	vars := mux.Vars(r)
	ascID := vars["appSessionId"]
	log.Infof(" APP Session ID  : %s", ascID)

	if len(NgcData.paData.asc) > 0 {

		if NgcData.paData.asc[ascID] != nil {

			delete(NgcData.paData.asc, ascID)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)

}

//PolicyAuthorizationAppSessionSubscribe Subscribe
func PolicyAuthorizationAppSessionSubscribe(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionSubscribe -- Entered")
	vars := mux.Vars(r)
	ascID := vars["appSessionId"]
	log.Infof(" APP Session ID  : %s", ascID)

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Err(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer closeReqBody(r)

	evSub := EventsSubscReqData{}

	err1 := json.Unmarshal(b, &evSub)
	if err1 != nil {
		log.Err(err1)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(evSub)

	if len(NgcData.paData.asc) > 0 {

		if NgcData.paData.asc[ascID] != nil {

			asc := NgcData.paData.asc[ascID]

			//delete(asc.AscReqData.EvSubsc)
			asc.AscReqData.EvSubsc = &evSub

			mdata, err2 := json.Marshal(evSub)
			if err2 != nil {
				log.Errf("Write Failed: %v", err2)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			w.WriteHeader(http.StatusCreated)

			//Send Success response to AF
			_, err = w.Write(mdata)
			if err != nil {
				log.Errf("Write Failed: %v", err)
				return
			}
			log.Infof("HTTP Response sent: %d", http.StatusCreated)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)

}

//PolicyAuthorizationAppSessionUnsubscribe Unsubscribe
func PolicyAuthorizationAppSessionUnsubscribe(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionUnsubscribe -- Entered")
	vars := mux.Vars(r)
	ascID := vars["appSessionId"]
	log.Infof(" APP Session ID  : %s", ascID)

	if len(NgcData.paData.asc) > 0 {
		if NgcData.paData.asc[ascID] != nil {
			asc := NgcData.paData.asc[ascID]
			asc.AscReqData.EvSubsc = nil
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)

}

//PolicyAuthorizationAppSessionTestNotify Trigger Notify towards AF
func PolicyAuthorizationAppSessionTestNotify(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("func PolicyAuthorizationAppSessionTestNotify -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])
	w.WriteHeader(http.StatusOK)
}

//PolicyAuthorizationAppSessionTestNotifyTerminate Triggers Terminate towards AF
func PolicyAuthorizationAppSessionTestNotifyTerminate(w http.ResponseWriter,
	r *http.Request) {

	log.Infoln("func PolicyAuthorizationAppSessionTestNotifyTerminate -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])
	w.WriteHeader(http.StatusOK)
}

func closeReqBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		log.Errf("response body was not closed properly")
	}
}

func getAppSessionContextRespData() (AscRespData AppSessionContextRespData) {

	//AscRespData := AppSessionContextRespData{}

	AscRespData.ServAuthInfo = ServAuthNotKnown
	AscRespData.SuppFeat = "0"
	return AscRespData

}
func genLocURI() (string, string) {

	ascIDStr := strconv.Itoa(gASCStartID + NgcData.paData.ascID)
	NgcData.paData.ascID++
	loc := NgcData.paData.locPrefix + ascIDStr

	return loc, ascIDStr
}
func getEventsNotification() (EvsNotif EventsNotification) {

	//EvsNotif := EventsNotification{}

	EvsNotif.AccessType = AccessType3Gpp
	EvsNotif.AnGwAddr = AnGwAddress("192.168.10.11")
	//EvsNotif.EvSubsUri = AscReqData.EvSubsc.NotifUri
	EvsNotif.EvNotifs = make([]PolicyEventNotification, 2)

	afEvnt := PolicyEventNotification{Event: ResourceAllocated}
	afEvnt.Flows = make([]Flows, 2)
	Flows1 := Flows{ContVers: []int32{1, 2}, FNums: []int32{3, 4}, MedCompN: 32}
	afEvnt.Flows[0] = Flows1
	afEvnt.Flows[1] = Flows1
	afEvnt2 := afEvnt
	afEvnt2.Event = FailedResourcesAllocation

	EvsNotif.EvNotifs[0] = afEvnt
	EvsNotif.EvNotifs[1] = afEvnt2

	EvsNotif.FailedResourcAllocReports = make([]ResourcesAllocationInfo, 2)
	ra := ResourcesAllocationInfo{McResourcStatus: MediaComponentResourceActive}

	ra.Flows = make([]Flows, 2)
	ra.Flows[0] = Flows1
	ra.Flows[1] = Flows1
	ra2 := ra
	ra2.McResourcStatus = MediaComponentResourceInActive
	EvsNotif.FailedResourcAllocReports[0] = ra
	EvsNotif.FailedResourcAllocReports[1] = ra2

	plmn := PlmnID{"100", "010"}
	EvsNotif.PlmnID = &plmn

	EvsNotif.QncReports = make([]QosNotificationControlInfo, 2)
	qos := QosNotificationControlInfo{NotifType: QosNotifGuaranteed}
	qos.Flows = make([]Flows, 2)
	qos.Flows[0] = Flows1
	qos.Flows[1] = Flows1

	qos2 := qos
	EvsNotif.QncReports[0] = qos
	EvsNotif.QncReports[0] = qos2

	EvsNotif.RatType = RatType("5G-NR")

	EvsNotif.UsgRep = AccumulatedUsage{10, 20, 30, 40}

	return EvsNotif

}
