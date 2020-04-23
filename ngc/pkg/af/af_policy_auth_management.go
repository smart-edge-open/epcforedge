// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getAppSessionID(r *http.Request) string {
	vars := mux.Vars(r)
	retVal := vars["appSessionId"]
	return retVal
}

func logPolicyRespErr(w *http.ResponseWriter, err string, statusCode int) {
	log.Errf("%s", err)
	(*w).WriteHeader(statusCode)
	return
}

// CreatePolicyAuthAppSessions func create one or more App Session Ctx
func CreatePolicyAuthAppSessions(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSess     AppSessionContext
		appSessResp AppSessionContext
		httpResp    *http.Response
		probDetails *ProblemDetails
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in CreatePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewDecoder(r.Body).Decode(&appSess)
	if err != nil {
		logPolicyRespErr(&w, "Json Decode error in "+
			"CreatePolicyAuthAppSessions: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	apiClient := afCtx.cfg.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"CreatePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthAppSessionAPI.PostAppSessions(cliCtx,
			appSess)
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

	if probDetails != nil {
		probDetailsJSON, err := json.Marshal(probDetails)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (probDetials)"+
				" in CreatePolicyAuthAppSessions: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJSON)
		if err != nil {
			log.Errf("Response write error in " +
				"CreatePolicyAuthAppSessions: " + err.Error())
			return
		}
		return
	}

	appSessJSON, err := json.Marshal(appSessResp)
	if err != nil {
		logPolicyRespErr(&w, "Json marshal error in "+
			"CreatePolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJSON)
	if err != nil {
		log.Errf("Response write error in CreatePolicyAuthAppSessions")
		return
	}
}

// DeletePolicyAuthAppSession func deletes an App Session Ctx
func DeletePolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err           error
		eventSubscReq EventsSubscReqData
		appSessResp   AppSessionContext
		httpResp      *http.Response
		probDetails   *ProblemDetails
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in DeletePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if body is not null then decode
	if r.Body != nil && r.ContentLength > 0 {
		err = json.NewDecoder(r.Body).Decode(&eventSubscReq)
		if err != nil {
			logPolicyRespErr(&w, "Json decode error in "+
				"DeletePolicyAuthAppSession: "+err.Error(),
				http.StatusBadRequest)
			return
		}
	}

	apiClient := afCtx.cfg.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"DeletePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	appSessionID := getAppSessionID(r)

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessAPI.
			DeleteAppSession(cliCtx, appSessionID, &eventSubscReq)
	if err != nil {
		if httpResp != nil {
			logPolicyRespErr(&w, "Delete Policy App Session: "+
				err.Error(), httpResp.StatusCode)
		} else {
			logPolicyRespErr(&w, "Delete Policy App Session: "+
				err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJSON, err := json.Marshal(probDetails)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (probDetials)"+
				" in DeletePolicyAuthAppSessions: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJSON)
		if err != nil {
			log.Errf("Response write error in " +
				"DeletePolicyAuthAppSessions")
			return
		}
		return
	}

	appSessJSON, err := json.Marshal(appSessResp)
	if err != nil {
		logPolicyRespErr(&w, "Json marshal error in "+
			"DeletePolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJSON)
	if err != nil {
		log.Errf("Response write error in DeletePolicyAuthAppSessions")
		return
	}
}

// GetPolicyAuthAppSession func retreives App Session Ctx from PCF server
func GetPolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSessResp AppSessionContext
		httpResp    *http.Response
		probDetails *ProblemDetails
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in GetPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiClient := afCtx.cfg.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"GetPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appSessionID := getAppSessionID(r)

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessAPI.GetAppSession(cliCtx,
			appSessionID)

	if err != nil {
		if httpResp != nil {
			logPolicyRespErr(&w, "Get Policy App Session: "+
				err.Error(), httpResp.StatusCode)
		} else {
			logPolicyRespErr(&w, "Get Policy App Session: "+
				err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJSON, err := json.Marshal(probDetails)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (probDetials)"+
				" in GetPolicyAuthAppSessions "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJSON)
		if err != nil {
			log.Errf("Response write error in " +
				"GetPolicyAuthAppSessions")
			return
		}
		return
	}

	appSessRespJSON, err := json.Marshal(appSessResp)
	if err != nil {
		logPolicyRespErr(&w, "Json marshal error in "+
			"GetPolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessRespJSON)
	if err != nil {
		log.Errf("Response write error in GetPolicyAuthAppSessions")
		return
	}
}

// ModifyPolicyAuthAppSession func modifies App Session Ctx
func ModifyPolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err               error
		appSessUpdateData AppSessionContextUpdateData
		appSessResp       AppSessionContext
		httpResp          *http.Response
		probDetails       *ProblemDetails
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in ModifyPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewDecoder(r.Body).Decode(&appSessUpdateData)
	if err != nil {
		log.Errf("Policy auth app sess create error in json decode %s",
			err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiClient := afCtx.cfg.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"ModifyPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	appSessionID := getAppSessionID(r)

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessAPI.ModAppSession(cliCtx,
			appSessionID, appSessUpdateData)
	if err != nil {
		if httpResp != nil {
			logPolicyRespErr(&w, "Modify Policy App Session: "+
				err.Error(), httpResp.StatusCode)
		} else {
			logPolicyRespErr(&w, "Modify Policy App Session: "+
				err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJSON, err := json.Marshal(probDetails)
		if err != nil {
			logPolicyRespErr(&w, "Json marshal error (probDetials)"+
				" in ModifyPolicyAuthAppSessions: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJSON)
		if err != nil {
			log.Errf("Response write error in " +
				"ModifyPolicyAuthAppSessions")
			return
		}
		return
	}

	appSessJSON, err := json.Marshal(appSessResp)
	if err != nil {
		logPolicyRespErr(&w, "Json marshal error in "+
			"ModifyPolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJSON)
	if err != nil {
		log.Errf("Response write error in ModifyPolicyAuthAppSessions")
		return
	}
}
