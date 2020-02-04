// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func getPfdTransaction(cliCtx context.Context, afCtx *Context,
	subscriptionID string) (PfdManagement, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	pfdTs, resp, err := cli.PfdManagementGetAPI.PfdTransactionGet(
		cliCtx, afCtx.cfg.AfID, subscriptionID)

	if err != nil {
		return PfdManagement{}, resp, err
	}
	return pfdTs, resp, nil
}

// GetPfdTransaction function
func GetPfdTransaction(w http.ResponseWriter, r *http.Request) {

	var (
		err              error
		pfdResp          PfdManagement
		resp             *http.Response
		pfdTransactionID string
		pfdRespJSON      []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Pfd Management get: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pfdTransactionID, err = getPfdTransIDFromURL(r.URL)
	if err != nil {
		log.Errf("Pfd Management get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdResp, resp, err = getPfdTransaction(cliCtx, afCtx, pfdTransactionID)
	if err != nil {
		log.Errf("Pfd Management get : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	pfdRespJSON, err = json.Marshal(pfdResp)
	if err != nil {
		log.Errf("Pfd Management get : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(pfdRespJSON); err != nil {
		log.Errf("Pfd Management get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
