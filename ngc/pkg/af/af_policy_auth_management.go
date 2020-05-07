// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"

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
}

func setAppSessNotifURICorreID(appSess *AppSessionContext, afCtx *Context) (
	err error) {

	ascReqData := appSess.AscReqData
	if ascReqData == nil {
		err = errors.New("Nil AppSessionContextReqData")
		log.Errf("%s", err.Error())
		return err
	}

	ascReqData.NotifURI = afCtx.cfg.CliPcfCfg.NotifURI

	afRoutReq := ascReqData.AfRoutReq
	if afRoutReq != nil && afRoutReq.UpPathChgSub != nil {
		id := atomic.AddInt32(&notifCorreID, 1)
		afRoutReq.UpPathChgSub.NotifCorreID = strconv.Itoa(int(id))
	} else {
		log.Errf("notif correl id is not set due to wrong req data")
	}
	return nil
}

// CreatePolicyAuthAppSessions func create one or more App Session Ctx
func CreatePolicyAuthAppSessions(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSess     AppSessionContext
		appSessResp AppSessionContext
		httpResp    *http.Response
		probDetails *ProblemDetails
		appSessJSON []byte
		url         *url.URL
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

	err = setAppSessNotifURICorreID(&appSess, afCtx)
	if err != nil {
		logPolicyRespErr(&w, "CreatePolicyAuthAppSessions: "+
			err.Error(), http.StatusBadRequest)
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

	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	appSessJSON, err = json.Marshal(appSessResp)
	if err != nil {
		logPolicyRespErr(&w, "Json marshal error in "+
			"CreatePolicyAuthAppSessions: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	uri := apiClient.locationPrefixURI
	if url, err = httpResp.Location(); err != nil {
		logPolicyRespErr(&w, "CreatePolicyAuthAppSessions: "+
			err.Error(), http.StatusInternalServerError)
		return
	}

	res := strings.Split(url.String(), "app-sessions")
	if len(res) == 2 {
		uri += res[1]
	} else {
		log.Errf("Location header returned from PCF is INCORRECT")
	}

	w.Header().Set("Location", uri)
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
		log.Infoln(w)
		_, err = (*w).Write(probDetailsJSON)
		if err != nil {
			log.Errf("Response write error in " + funcName +
				err.Error())
			return
		}
		(*w).WriteHeader(httpResp.StatusCode)
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
		appSessResp   AppSessionContext
		httpResp      *http.Response
		probDetails   *ProblemDetails
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

	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	if httpResp.StatusCode == 204 {
		return
	}

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
		err         error
		appSessResp AppSessionContext
		httpResp    *http.Response
		probDetails *ProblemDetails
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

	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
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
		appSessJSON       []byte
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

	err = json.NewDecoder(r.Body).Decode(&appSessUpdateData)
	if err != nil {
		log.Errf("Policy auth app sess create error in json decode %s",
			err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if appSessUpdateData.EvSubsc != nil {
		appSessUpdateData.EvSubsc.NotifURI = afCtx.cfg.CliPcfCfg.
			NotifURI
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

	if err != nil || probDetails != nil {
		handlePAErrorResp(probDetails, err, &w, httpResp, funcName)
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
