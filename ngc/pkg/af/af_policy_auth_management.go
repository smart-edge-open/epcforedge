package af

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func logErr(w *http.ResponseWriter, err string, statusCode int) {
	log.Errf("%s")
	(*w).WriteHeader(statusCode)
	return
}
func CreatePolicyAuthAppSessions(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSess     AppSessionContext
		appSessResp AppSessionContext
		httpResp    *http.Response
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logErr(&w, "nil afCtx in CreatePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewDecoder(r.Body).Decode(&appSess)
	if err != nil {
		logErr(&w, "Json Decode error in CreatePolicyAuthAppSessions: "+
			err.Error(), http.StatusBadRequest)
		return
	}

	apiClient := AfCtx.cfg.policyAuthApiClient
	if apiClient == nil {
		logErr(&w, "nil policyAuthApiClient in "+
			"CreatePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	/*
		cliCfg := AfCtx.cfg.CliPcfCfg
		if cliCfg == nil {

			log.Errf("CreatePolicyAuthAppSessions nil client config")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/

	appSessResp, httpResp, err =
		apiClient.PolicyAuthAppSessionApi.PostAppSessions(cliCtx,
			appSess)
	if err != nil {
		// Handle error
		if httpResp != nil {
			logErr(&w, "Create Policy App Session"+err.Error(),
				httpResp.StatusCode)
		} else {
			logErr(&w, "Create Policy App Session"+err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	appSessJson, err := json.Marshal(appSessResp)
	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJson)
	if err != nil {
		// Handle Error
		return
	}
}

func DeletePolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err           error
		eventSubscReq EventsSubscReqData
		appSessResp   AppSessionContext
		httpResp      *http.Response
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("CreatePolicyAuthAppSessions nil afCtx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if body is not null then decode
	err = json.NewDecoder(r.Body).Decode(&eventSubscReq)
	if err != nil {
		log.Errf("Policy auth app sess create error in json decode %s",
			err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiClient := AfCtx.cfg.policyAuthApiClient
	if apiClient == nil {
		log.Errf("CreatePolicyAuthAppSessions nil policyAuthApiClient")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Do we need this?
	cliCfg := AfCtx.cfg.CliPcfCfg
	if cliCfg == nil {
		log.Errf("CreatePolicyAuthAppSessions nil client config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appSessionId := getAppSessionId(r) // use inline or call func?

	appSessResp, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessApi.DeleteAppSession(cliCtx,
			appSessionId, &eventSubscReq)
	if err != nil {
		// Handle error
		return
	}

	appSessJson, err := json.Marshal(appSessResp)
	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJson)
	if err != nil {
		// Handle Error
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func getAppSessionId(r *http.Request) string {
	vars := mux.Vars(r)
	retVal := vars["appSessionId"]
	return retVal
}

func GetPolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSessResp AppSessionContext
		httpResp    *http.Response
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("CreatePolicyAuthAppSessions nil afCtx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiClient := AfCtx.cfg.policyAuthApiClient
	if apiClient == nil {
		log.Errf("CreatePolicyAuthAppSessions nil policyAuthApiClient")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appSessionId := getAppSessionId(r) // use inline or call func?

	appSessResp, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessApi.GetAppSession(cliCtx,
			appSessionId)

	if err != nil {
		// Handle Error
		return
	}

	appSessRespJson, err := json.Marshal(appSessResp)
	if err != nil {
		// Handle Error
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessRespJson)
	if err != nil {
		// Handle Error
		return
	}
}

func ModifyPolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err               error
		appSessUpdateData AppSessionContextUpdateData
		appSessResp       AppSessionContext
		httpResp          *http.Response
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("CreatePolicyAuthAppSessions nil afCtx")
		w.WriteHeader(http.StatusInternalServerError)
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

	apiClient := AfCtx.cfg.policyAuthApiClient
	if apiClient == nil {
		log.Errf("CreatePolicyAuthAppSessions nil policyAuthApiClient")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Do we need this?
	cliCfg := AfCtx.cfg.CliPcfCfg
	if cliCfg == nil {
		log.Errf("CreatePolicyAuthAppSessions nil client config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appSessionId := getAppSessionId(r) // use inline or call func?

	appSessResp, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessApi.ModAppSession(cliCtx,
			appSessionId, appSessUpdateData)
	if err != nil {
		// Handle error
		return
	}

	appSessJson, err := json.Marshal(appSessResp)
	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJson)
	if err != nil {
		// Handle Error
		return
	}
}
