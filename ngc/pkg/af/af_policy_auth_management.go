// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var contentTypeJSON string = "application/json"

var smfPANotifURI string

var pcfPANotifURI string

func getAppSessionID(r *http.Request) string {
	vars := mux.Vars(r)
	retVal := vars["appSessionId"]
	return retVal
}

func logPolicyRespErr(w *http.ResponseWriter, err string, statusCode int) {
	log.Errf("%s", err)
	(*w).WriteHeader(statusCode)
}

func handlePAErrorResp(probDetails *ProblemDetails, err error,
	w *http.ResponseWriter, httpResp *http.Response, funcName string) {

	var probDetailsJSON []byte
	if err != nil {
		if httpResp != nil {
			logPolicyRespErr(w, funcName+err.Error(),
				httpResp.StatusCode)
		} else {
			logPolicyRespErr(w, funcName+err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJSON, err = json.Marshal(probDetails)
		if err != nil {
			logPolicyRespErr(w, "Json marshal error (probDetials)"+
				" in "+funcName+err.Error(),
				http.StatusInternalServerError)
			return
		}

		(*w).WriteHeader(httpResp.StatusCode)
		_, err = (*w).Write(probDetailsJSON)
		if err != nil {
			log.Errf("Response write error in " + funcName +
				err.Error())
			return
		}
		return
	}

}

// CreatePolicyAuthAppSessions func create one or more App Session Ctx
func CreatePolicyAuthAppSessions(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSess     AppSessionContext
		appSessJSON []byte
		evInfo      EventInfo
	)

	funcName := "CreatePolicyAuthAppSessions: "
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

	err = validateAndSetupAppSessCtx(&appSess, &evInfo, afCtx)

	if err != nil {
		logPolicyRespErr(&w, "CreatePolicyAuthAppSessions: "+
			err.Error(), http.StatusBadRequest)
		return
	}

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"CreatePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	apiResp, err := apiClient.PostAppSessions(cliCtx,
		appSess)

	probDetails := apiResp.probDetails
	httpResp := apiResp.httpResp
	if err != nil || probDetails != nil {
		if len(apiResp.retryAfter) > 0 {
			w.Header().Set("Retry-After", apiResp.retryAfter)
		}
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	w.Header().Set("Location", apiResp.locationURI)

	if httpResp.StatusCode == 303 {
		w.WriteHeader(httpResp.StatusCode)
		return
	}

	appSessResp := apiResp.appSessCtx
	err = setAppSessInfo(apiResp.locationURI, &evInfo, appSessResp, afCtx)
	if err != nil {
		logPolicyRespErr(&w, "CreatePolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	appSessJSON, err = json.Marshal(appSessResp)
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
		err        error
		evsReqData EventsSubscReqData
		apiResp    PcfPAResponse
	)

	funcName := "DeletePolicyAuthAppSession: "
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
	err = decodeValidateEventSubscReq(r, w, &evsReqData)
	if err != nil {
		logPolicyRespErr(&w, "DeletePolicyAuthAppSessions: "+
			err.Error(), http.StatusBadRequest)
		return
	}

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"DeletePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	appSessionID := getAppSessionID(r)

	apiResp, err = apiClient.DeleteAppSession(cliCtx,
		appSessionID, &evsReqData)

	httpResp := apiResp.httpResp
	probDetails := apiResp.probDetails
	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	// Close the websocket if no other notifyId
	log.Infoln("Deleting the appSessionsEv for appSessionID", appSessionID)
	evInfo := afCtx.appSessionsEv[appSessionID]
	err = chkRemoveWSConn(evInfo, appSessionID, afCtx)
	if err != nil {
		log.Errf(err.Error())
	}
	delete(afCtx.appSessionsEv, appSessionID)

	w.WriteHeader(httpResp.StatusCode)
	if httpResp.StatusCode == 204 {
		return
	}

	appSessResp := apiResp.appSessCtx
	appSessJSON, err := json.Marshal(appSessResp)
	if err != nil {
		logPolicyRespErr(&w, "Json marshal error in "+
			"DeletePolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	_, err = w.Write(appSessJSON)
	if err != nil {
		log.Errf("Response write error in DeletePolicyAuthAppSessions")
		return
	}
}

// GetPolicyAuthAppSession func retreives App Session Ctx from PCF server
func GetPolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err     error
		apiResp PcfPAResponse
	)

	funcName := "GetPolicyAuthAppSession: "
	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in GetPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"GetPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appSessionID := getAppSessionID(r)

	apiResp, err = apiClient.GetAppSession(cliCtx, appSessionID)

	probDetails := apiResp.probDetails
	httpResp := apiResp.httpResp
	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	appSessResp := apiResp.appSessCtx
	err = updateAppSessInResp(appSessResp, appSessionID, afCtx)
	if err != nil {
		logPolicyRespErr(&w, "Updating the response "+
			"GetPolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
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
		err           error
		ascUpdateData AppSessionContextUpdateData
		appSessJSON   []byte
		apiResp       PcfPAResponse
	)

	funcName := "ModifyPolicyAuthAppSession: "
	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in ModifyPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewDecoder(r.Body).Decode(&ascUpdateData)
	if err != nil {
		log.Errf("Policy auth app sess create error in json decode %s",
			err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	appSessionID := getAppSessionID(r)

	log.Infoln("App SessionID received in PATCH ", appSessionID)

	err = validateAscUpdateData(&ascUpdateData)
	if err != nil {
		logPolicyRespErr(&w, funcName+err.Error(),
			http.StatusBadRequest)
		return
	}

	err = modifyAppSessNotifParams(&ascUpdateData, appSessionID, afCtx)
	if err != nil {
		logPolicyRespErr(&w, "ModifyPolicyAuthAppSession: "+
			err.Error(), http.StatusBadRequest)
		return
	}

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"ModifyPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	apiResp, err = apiClient.ModAppSession(cliCtx, appSessionID,
		ascUpdateData)

	probDetails := apiResp.probDetails
	httpResp := apiResp.httpResp
	if err != nil || probDetails != nil {
		if len(apiResp.retryAfter) > 0 {
			w.Header().Set("Retry-After", apiResp.retryAfter)
		}

		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	appSessResp := apiResp.appSessCtx
	err = updateAppSessInResp(appSessResp, appSessionID, afCtx)
	if err != nil {
		logPolicyRespErr(&w, "Updating the response "+
			"ModifyPolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	appSessJSON, err = json.Marshal(appSessResp)
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
