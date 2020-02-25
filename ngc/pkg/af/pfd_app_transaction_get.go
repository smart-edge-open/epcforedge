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
		errRspHeader(&w, "APP GET", "af-ctx retrieved from request is nil",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pfdTransactionID = getPfdTransIDFromURL(r)

	appID = getPfdAppIDFromURL(r)

	pfdResp, resp, err = getPfdAppTransaction(cliCtx, afCtx, pfdTransactionID,
		appID)
	if err != nil {
		errRspHeader(&w, "APP GET", err.Error(), getStatusCode(resp))
		return
	}

	// Updating the Self Application Link
	self, err := updateAppLink(afCtx.cfg, r, pfdResp)
	if err != nil {
		errRspHeader(&w, "APP GET", err.Error(), http.StatusInternalServerError)
		return
	}
	pfdResp.Self = Link(self)

	pfdRespJSON, err = json.Marshal(pfdResp)
	if err != nil {
		errRspHeader(&w, "APP GET", err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(pfdRespJSON); err != nil {
		errRspHeader(&w, "APP GET", err.Error(), http.StatusInternalServerError)
		return
	}
}
