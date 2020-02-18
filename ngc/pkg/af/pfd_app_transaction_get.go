// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func getPfdAppTransaction(cliCtx context.Context, afCtx *Context,
	pfdTransID string, appID string) (PfdData, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	pfdTs, resp, err := cli.PfdManagementAppGetAPI.PfdAppTransactionGet(
		cliCtx, afCtx.cfg.AfID, pfdTransID, appID)

	if err != nil {
		return PfdData{}, resp, err
	}
	return pfdTs, resp, nil
}

// GetPfdAppTransaction function
func GetPfdAppTransaction(w http.ResponseWriter, r *http.Request) {

	var (
		err              error
		pfdResp          PfdData
		appID            string
		resp             *http.Response
		pfdTransactionID string
		pfdRespJSON      []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Pfd Management Application get: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pfdTransactionID, err = getPfdTransIDFromURL(r)

	if err != nil {
		log.Errf("Pfd Management  Application get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appID, err = getPfdAppIDFromURL(r)
	log.Infof("Application ID from URL is %s", appID)
	if err != nil {
		log.Errf("Pfd Management Application  get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pfdResp, resp, err = getPfdAppTransaction(cliCtx, afCtx, pfdTransactionID,
		appID)
	if err != nil {
		log.Errf("Pfd Management Application get : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	pfdRespJSON, err = json.Marshal(pfdResp)
	if err != nil {
		log.Errf("Pfd Management Application get : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(pfdRespJSON); err != nil {
		log.Errf("Pfd Management Application get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
