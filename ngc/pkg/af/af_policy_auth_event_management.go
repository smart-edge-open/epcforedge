// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func handleEventSubscResp(w *http.ResponseWriter,
	eventSubscResp EventSubscResponse) {

	var (
		respBody []byte
		err      error
	)

	httpResp := eventSubscResp.httpResp
	if httpResp.StatusCode == 201 {
		(*w).Header().Set("Location", eventSubscResp.locationURI)
	}

	if eventSubscResp.eventSubscReq != nil {
		respBody, err = json.Marshal(eventSubscResp.eventSubscReq)
		if err != nil {
			logPolicyRespErr(w, "Json marshal error (eventSubsc)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
	}

	if eventSubscResp.evsNotif != nil {
		log.Infoln("Notification Received as part of event subscribe")
	} else if eventSubscResp.probDetails != nil {
		respBody, err = json.Marshal(eventSubscResp.probDetails)
		if err != nil {
			logPolicyRespErr(w, "Json marshal error (probDetails)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
		}
	}

	(*w).WriteHeader(httpResp.StatusCode)
	_, err = (*w).Write(respBody)
	if err != nil {
		log.Errf("Response write error in " +
			"PolicyAuthEvemtSubsc: " + err.Error())
	}

}

// PolicyAuthEventSubsc Event susbscription request handler
func PolicyAuthEventSubsc(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		evsReqData     EventsSubscReqData
		eventSubscResp EventSubscResponse
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in PolicyAuthAppEventSubs",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewDecoder(r.Body).Decode(&evsReqData)
	if err != nil {
		logPolicyRespErr(&w, "Json Decode error in "+
			"PolicyAuthAppEventSubs: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	err = validateEventSubsc(&evsReqData)
	if err != nil {
		logPolicyRespErr(&w, "PolicyAuthAppEventSubs : "+
			err.Error(), http.StatusBadRequest)
		return
	}

	evsReqData.NotifURI = pcfPANotifURI
	appSessionID := getAppSessionID(r)

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"PolicyAuthAppEventSubs",
			http.StatusInternalServerError)
		return
	}

	eventSubscResp, err = apiClient.UpdateEventsSubsc(cliCtx, appSessionID,
		&evsReqData)

	httpResp := eventSubscResp.httpResp
	if err != nil {
		if httpResp != nil {
			logPolicyRespErr(&w, "PolicyAuthAppEventSubs: "+
				err.Error(), httpResp.StatusCode)
		} else {
			logPolicyRespErr(&w, "PolicyAuthAppEventSubs: "+
				err.Error(), http.StatusInternalServerError)
		}
		return
	}

	handleEventSubscResp(&w, eventSubscResp)
}

// PolicyAuthEventDelete Event delete request handler
func PolicyAuthEventDelete(w http.ResponseWriter, r *http.Request) {

	var (
		err     error
		evsResp EventSubscResponse
	)

	funcName := "PolicyAuthEventDelete: "
	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in PolicyAuthAppEventDelete",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	appSessionID := getAppSessionID(r)

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"PolicyAuthAppEventDelete",
			http.StatusInternalServerError)
		return
	}

	evsResp, err = apiClient.DeleteEventsSubsc(cliCtx, appSessionID)

	probDetails := evsResp.probDetails
	httpResp := evsResp.httpResp
	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
}
