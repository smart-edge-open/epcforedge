// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

// EventSubscResponse struct anyof EventSusbscReqData, EventNotification,
// ProblemDetails
type EventSubscResponse struct {
	eventSubscReq *EventsSubscReqData
	evsNotif      *EventsNotification
	probDetails   *ProblemDetails
}

func handleEventSubscResp(w http.ResponseWriter,
	eventSubscResp EventSubscResponse, httpResp *http.Response) {

	var (
		respBody []byte
		err      error
	)

	if eventSubscResp.eventSubscReq != nil {
		respBody, err = json.Marshal(eventSubscResp.eventSubscReq)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (eventSubsc)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else if eventSubscResp.evsNotif != nil {
		respBody, err = json.Marshal(eventSubscResp.evsNotif)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (evsNotif)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else if eventSubscResp.probDetails != nil {
		respBody, err = json.Marshal(eventSubscResp.probDetails)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (probDetails)"+
				" in PolicyAuthEventSubsc: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(httpResp.StatusCode)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		log.Errf("Response write error in " +
			"PolicyAuthEvemtSubsc: " + err.Error())
	}

	w.WriteHeader(httpResp.StatusCode)
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

	eventSubscReq.NotifURI = afCtx.cfg.CliPcfCfg.PolicyAuthNotifURI
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
			logPolicyRespErr(&w, "Create Policy App Session: "+
				err.Error(), httpResp.StatusCode)
		} else {
			logPolicyRespErr(&w, "Create Policy App Session: "+
				err.Error(), http.StatusInternalServerError)
		}
		return
	}

	handleEventSubscResp(w, eventSubscResp, httpResp)
}

// PolicyAuthEventDelete Event delete request handler
func PolicyAuthEventDelete(w http.ResponseWriter, r *http.Request) {

	var (
		err             error
		httpResp        *http.Response
		probDetails     *ProblemDetails
		probDetailsJSON []byte
	)

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

	if err != nil {
		if httpResp != nil {
			logPolicyRespErr(&w, "Policy auth event delete: "+
				err.Error(), httpResp.StatusCode)
		} else {
			logPolicyRespErr(&w, "Policy auth event delete: "+
				err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJSON, err = json.Marshal(probDetails)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (probDetials)"+
				" in PolicyAuthEventDelete: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		w.WriteHeader(httpResp.StatusCode)
		_, err = w.Write(probDetailsJSON)
		if err != nil {
			log.Errf("Response write error in " +
				" in PolicyAuthEventDelete: " + err.Error())
		}
		return
	}

	w.WriteHeader(httpResp.StatusCode)
}
