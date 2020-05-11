// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"errors"
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

func setAppSessNotifURI(appSess *AppSessionContext, afCtx *Context) {

	ascReqData := appSess.AscReqData
	ascReqData.NotifURI = pcfPANotifURI

	afRoutReq := ascReqData.AfRoutReq
	if afRoutReq != nil && afRoutReq.UpPathChgSub != nil {
		afRoutReq.UpPathChgSub.NotificationURI = smfPANotifURI
	}

	for _, medCompn := range ascReqData.MedComponents {
		afRoutReq := medCompn.AfRoutReq
		if afRoutReq != nil && afRoutReq.UpPathChgSub != nil {
			afRoutReq.UpPathChgSub.NotificationURI = smfPANotifURI
		}
	}

	if ascReqData.EvSubsc != nil {
		ascReqData.EvSubsc.NotifURI = pcfPANotifURI
	}
}

func validateAppSessCtx(appSess *AppSessionContext) (err error) {

	ascReqData := appSess.AscReqData
	if ascReqData == nil {
		err = errors.New("nil ascReqData")
		return err
	}

	return nil
	/*
		switch ascReqData.SuppFeat {
		case "InfluenceOnTrafficRouting":
			return nil
	*/

}

// CreatePolicyAuthAppSessions func create one or more App Session Ctx
func CreatePolicyAuthAppSessions(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSess     AppSessionContext
		appSessJSON []byte
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

	err = validateAppSessCtx(&appSess)
	if err != nil {
		logPolicyRespErr(&w, "CreatePolicyAuthAppSessions: "+
			err.Error(), http.StatusBadRequest)
		return
	}

	setAppSessNotifURI(&appSess, afCtx)

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"CreatePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	//	appSessResp, probDetails, httpResp, err =
	apiResp, err := apiClient.PostAppSessions(cliCtx,
		appSess)

	probDetails := apiResp.probDetails
	httpResp := apiResp.httpResp
	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	appSessResp := apiResp.appSessCtx
	appSessJSON, err = json.Marshal(appSessResp)
	if err != nil {
		logPolicyRespErr(&w, "Json marshal error in "+
			"CreatePolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", apiResp.locationURI)
	w.WriteHeader(httpResp.StatusCode)

	_, err = w.Write(appSessJSON)
	if err != nil {
		log.Errf("Response write error in CreatePolicyAuthAppSessions")
		return
	}
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

func decodeEventSubscReq(r *http.Request, w http.ResponseWriter,
	eventSubscReq *EventsSubscReqData) (err error) {

	if r.Body != nil && r.ContentLength > 0 {
		err = json.NewDecoder(r.Body).Decode(eventSubscReq)
		if err != nil {
			logPolicyRespErr(&w, "Json decode error in "+
				"DeletePolicyAuthAppSession: "+err.Error(),
				http.StatusBadRequest)
			return err
		}
	}
	return nil
}

// DeletePolicyAuthAppSession func deletes an App Session Ctx
func DeletePolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err           error
		eventSubscReq EventsSubscReqData
		apiResp       PcfPAResponse
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
	err = decodeEventSubscReq(r, w, &eventSubscReq)
	if err != nil {
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
		appSessionID, &eventSubscReq)

	httpResp := apiResp.httpResp
	probDetails := apiResp.probDetails
	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

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

func setAppSessUpdateDataNotifURI(ascUpdateData *AppSessionContextUpdateData,
	afCtx *Context) {

	afRoutReq := ascUpdateData.AfRoutReq
	if afRoutReq != nil && afRoutReq.UpPathChgSub != nil {
		afRoutReq.UpPathChgSub.NotificationURI = smfPANotifURI
	}

	if ascUpdateData.EvSubsc != nil {
		ascUpdateData.EvSubsc.NotifURI = pcfPANotifURI
	}

	for _, medCompn := range ascUpdateData.MedComponents {
		afRoutReq := medCompn.AfRoutReq
		if afRoutReq != nil && afRoutReq.UpPathChgSub != nil {
			afRoutReq.UpPathChgSub.NotificationURI = smfPANotifURI
		}
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

	setAppSessUpdateDataNotifURI(&ascUpdateData, afCtx)

	apiClient := afCtx.data.policyAuthAPIClient
	if apiClient == nil {
		logPolicyRespErr(&w, "nil policyAuthAPIClient in "+
			"ModifyPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	appSessionID := getAppSessionID(r)

	apiResp, err = apiClient.ModAppSession(cliCtx, appSessionID,
		ascUpdateData)

	probDetails := apiResp.probDetails
	httpResp := apiResp.httpResp
	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	appSessResp := apiResp.appSessCtx
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
