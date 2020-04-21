/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngccntest

import (
	//"context"
	"github.com/gorilla/mux"
	"net/http"
)

//PolicyAuthorizationAppSessionGet Get
func PolicyAuthorizationAppSessionGet(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionGet -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])

	w.WriteHeader(http.StatusOK)

}

//PolicyAuthorizationAppSessionCreate Post
func PolicyAuthorizationAppSessionCreate(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionCreate -- Entered")

	w.WriteHeader(http.StatusCreated)

	/*
		asc := AppSessionContext{}

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Err(err)
			sendCustomeErrorRspToAF(w, 400, "Failed to read HTTP POST Body")
			return
		}

		defer closeReqBody(r)
		err1 := json.Unmarshal(b, &asc)
		if err1 != nil {
			log.Err(err1)
			sendCustomeErrorRspToAF(w, 400, "Failed UnMarshal PATCH data")
			return
		}


		AscReqData  *AppSessionContextReqData  `json:"ascReqData,omitempty"`
		AscRespData *AppSessionContextRespData `json:"ascRespData,omitempty"`
		EvsNotif    *EventsNotification        `json:"evsNotif,omitempty"`
	*/

}

//PolicyAuthorizationAppSessionPatch Patch
func PolicyAuthorizationAppSessionPatch(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSession -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])
	w.WriteHeader(http.StatusNoContent)

}

//PolicyAuthorizationAppSessionDelete Delete
func PolicyAuthorizationAppSessionDelete(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionDelete -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])
	w.WriteHeader(http.StatusNoContent)

}

//PolicyAuthorizationAppSessionSubscribe Subscribe
func PolicyAuthorizationAppSessionSubscribe(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionSubscribe -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])
	w.WriteHeader(http.StatusNoContent)

}

//PolicyAuthorizationAppSessionUnsubscribe Unsubscribe
func PolicyAuthorizationAppSessionUnsubscribe(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionUnsubscribe -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])
	w.WriteHeader(http.StatusNoContent)

}

func closeReqBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		log.Errf("response body was not closed properly")
	}
}
