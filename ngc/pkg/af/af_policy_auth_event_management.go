// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// EventSubscResponse struct anyof EventSusbscReqData, EventNotification,
// ProblemDetails
type EventSubscResponse struct {
	eventSubscReq *EventsSubscReqData
	evsNotif      *EventsNotification
	probDetails   *ProblemDetails
}

func handleEventSubscResp(w *http.ResponseWriter, locationPrefixURI string,
	eventSubscResp EventSubscResponse, httpResp *http.Response) {

	var (
		respBody []byte
		err      error
		url      *url.URL
	)

	if eventSubscResp.eventSubscReq != nil {
		respBody, err = json.Marshal(eventSubscResp.eventSubscReq)
		if err != nil {
			logPolicyRespErr(w, "Json marshal error (eventSubsc)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
		}
	} else if eventSubscResp.evsNotif != nil {
		respBody, err = json.Marshal(eventSubscResp.evsNotif)
		if err != nil {
			logPolicyRespErr(w, "Json marshal error (evsNotif)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
		}
	} else if eventSubscResp.probDetails != nil {
		respBody, err = json.Marshal(eventSubscResp.probDetails)
		if err != nil {
			logPolicyRespErr(w, "Json marshal error (probDetails)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
		}
	} else {
		(*w).WriteHeader(httpResp.StatusCode)
		return
	}

	_, err = (*w).Write(respBody)
	if err != nil {
		log.Errf("Response write error in " +
			"PolicyAuthEvemtSubsc: " + err.Error())
	}

	if httpResp.StatusCode == 201 {
		uri := locationPrefixURI
		if url, err = httpResp.Location(); err != nil {
			logPolicyRespErr(w, "PolicyAuthEventSubsc: "+
				err.Error(), httpResp.StatusCode)
			return
		}

		res := strings.Split(url.String(), "app-sessions")
		if len(res) == 2 {
			uri += res[1]
		} else {
			log.Errf("Location URI from PCF is INCORRECT")
		}

		(*w).Header().Set("Location", uri)
	}

	(*w).WriteHeader(httpResp.StatusCode)
}

// PolicyAuthEventSubsc Event susbscription request handler
func PolicyAuthEventSubsc(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		eventSubscReq  EventsSubscReqData
		httpResp       *http.Response
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

	err = json.NewDecoder(r.Body).Decode(&eventSubscReq)
	if err != nil {
		logPolicyRespErr(&w, "Json Decode error in "+
			"PolicyAuthAppEventSubs: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	eventSubscReq.NotifURI = afCtx.cfg.CliPcfCfg.NotifURI
	appSessionID := getAppSessionID(r)

	apiClient := afCtx.cfg.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"PolicyAuthAppEventSubs",
			http.StatusInternalServerError)
		return
	}

	eventSubscResp, httpResp, err =
		apiClient.PolicyAuthEventSubsAPI.UpdateEventsSubsc(cliCtx,
			appSessionID, &eventSubscReq)

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

	handleEventSubscResp(&w, apiClient.locationPrefixURI, eventSubscResp,
		httpResp)
}

// PolicyAuthEventDelete Event delete request handler
func PolicyAuthEventDelete(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		httpResp    *http.Response
		probDetails *ProblemDetails
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

	apiClient := afCtx.cfg.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"PolicyAuthAppEventDelete",
			http.StatusInternalServerError)
		return
	}

	probDetails, httpResp, err =
		apiClient.PolicyAuthEventSubsAPI.DeleteEventsSubsc(cliCtx,
			appSessionID)

	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
}
