package af

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getAppSessionId(r *http.Request) string {
	vars := mux.Vars(r)
	retVal := vars["appSessionId"]
	return retVal
}

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
		probDetails *ProblemDetails
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

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthAppSessionApi.PostAppSessions(cliCtx,
			appSess)
	if err != nil {
		if httpResp != nil {
			logErr(&w, "Create Policy App Session"+err.Error(),
				httpResp.StatusCode)
		} else {
			logErr(&w, "Create Policy App Session"+err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJson, err := json.Marshal(probDetails)
		if err != nil {
			logErr(&w, "Json marshal error (probDetials) in "+
				"CreatePolicyAuthAppSessions"+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJson)
		if err != nil {
			log.Errf("Response write error in " +
				"CreatePolicyAuthAppSessions" + err.Error())
			return
		}
		return
	}

	appSessJson, err := json.Marshal(appSessResp)
	if err != nil {
		logErr(&w, "Json marshal error in CreatePolicyAuthAppSessions"+
			err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJson)
	if err != nil {
		log.Errf("Response write error in CreatePolicyAuthAppSessions")
		return
	}
}

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
		logErr(&w, "nil afCtx in DeletePolicyAuthAppSessions",
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
			logErr(&w, "Json decode error in "+
				"DeletePolicyAuthAppSession: "+err.Error(),
				http.StatusBadRequest)
			return
		}
	}

	apiClient := AfCtx.cfg.policyAuthApiClient
	if apiClient == nil {
		logErr(&w, "nil policyAuthApiClient in "+
			"DeletePolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	appSessionId := getAppSessionId(r)

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessApi.
			DeleteAppSession(cliCtx, appSessionId, &eventSubscReq)
	if err != nil {
		if httpResp != nil {
			logErr(&w, "Delete Policy App Session"+err.Error(),
				httpResp.StatusCode)
		} else {
			logErr(&w, "Delete Policy App Session"+err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJson, err := json.Marshal(probDetails)
		if err != nil {
			logErr(&w, "Json marshal error (probDetials) in "+
				"DeletePolicyAuthAppSessions"+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJson)
		if err != nil {
			log.Errf("Response write error in " +
				"DeletePolicyAuthAppSessions")
			return
		}
		return
	}

	appSessJson, err := json.Marshal(appSessResp)
	if err != nil {
		logErr(&w, "Json marshal error in DeletePolicyAuthAppSessions"+
			err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJson)
	if err != nil {
		log.Errf("Response write error in DeletePolicyAuthAppSessions")
		return
	}
}

func GetPolicyAuthAppSession(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		appSessResp AppSessionContext
		httpResp    *http.Response
		probDetails *ProblemDetails
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logErr(&w, "nil afCtx in GetPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiClient := AfCtx.cfg.policyAuthApiClient
	if apiClient == nil {
		logErr(&w, "nil policyAuthApiClient in "+
			"GetPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appSessionId := getAppSessionId(r)

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessApi.GetAppSession(cliCtx,
			appSessionId)

	if err != nil {
		if httpResp != nil {
			logErr(&w, "Get Policy App Session"+err.Error(),
				httpResp.StatusCode)
		} else {
			logErr(&w, "Get Policy App Session"+err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJson, err := json.Marshal(probDetails)
		if err != nil {
			logErr(&w, "Json marshal error (probDetials) in "+
				"GetPolicyAuthAppSessions"+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJson)
		if err != nil {
			log.Errf("Response write error in " +
				"GetPolicyAuthAppSessions")
			return
		}
		return
	}

	appSessRespJson, err := json.Marshal(appSessResp)
	if err != nil {
		logErr(&w, "Json marshal error in GetPolicyAuthAppSessions"+
			err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessRespJson)
	if err != nil {
		log.Errf("Response write error in GetPolicyAuthAppSessions")
		return
	}
}

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
		logErr(&w, "nil afCtx in ModifyPolicyAuthAppSessions",
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

	apiClient := AfCtx.cfg.policyAuthApiClient
	if apiClient == nil {
		logErr(&w, "nil policyAuthApiClient in "+
			"ModifyPolicyAuthAppSessions",
			http.StatusInternalServerError)
		return
	}

	appSessionId := getAppSessionId(r)

	appSessResp, probDetails, httpResp, err =
		apiClient.PolicyAuthIndividualAppSessApi.ModAppSession(cliCtx,
			appSessionId, appSessUpdateData)
	if err != nil {
		if httpResp != nil {
			logErr(&w, "Modify Policy App Session"+err.Error(),
				httpResp.StatusCode)
		} else {
			logErr(&w, "Modify Policy App Session"+err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	if probDetails != nil {
		probDetailsJson, err := json.Marshal(probDetails)
		if err != nil {
			logErr(&w, "Json marshal error (probDetials) in "+
				"ModifyPolicyAuthAppSessions"+err.Error(),
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(probDetailsJson)
		if err != nil {
			log.Errf("Response write error in " +
				"ModifyPolicyAuthAppSessions")
			return
		}
		return
	}

	appSessJson, err := json.Marshal(appSessResp)
	if err != nil {
		logErr(&w, "Json marshal error in ModifyPolicyAuthAppSessions"+
			err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpResp.StatusCode)
	_, err = w.Write(appSessJson)
	if err != nil {
		log.Errf("Response write error in ModifyPolicyAuthAppSessions")
		return
	}
}
