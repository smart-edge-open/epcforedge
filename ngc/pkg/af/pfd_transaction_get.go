// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func getPfdTransaction(cliCtx context.Context, afCtx *Context,
	pfdTrans string) (PfdManagement, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	pfdTs, resp, err := cli.PfdManagementGetAPI.PfdTransactionGet(
		cliCtx, afCtx.cfg.AfID, pfdTrans)

	if err != nil {
		return PfdManagement{}, resp, err
	}
	return pfdTs, resp, nil
}

// GetPfdTransaction function - Read a particular PFD transaction
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
		errRspHeader(&w, "GET", "af-ctx retrieved from request is nil",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pfdTransactionID = getPfdTransIDFromURL(r)

	pfdResp, resp, err = getPfdTransaction(cliCtx, afCtx, pfdTransactionID)
	if err != nil {
		if resp != nil {
			errRspHeader(&w, "GET", err.Error(), getStatusCode(resp))
		} else {
			errRspHeader(&w, "GET", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Updating the Self Link and Applications Self Link in AF

	self, err := updateSelfLink(afCtx.cfg, r, pfdResp)
	if err != nil {
		errRspHeader(&w, "GET", err.Error(),
			http.StatusInternalServerError)
		return
	}
	pfdResp.Self = Link(self)
	err = updateAppsLink(afCtx.cfg, r, pfdResp)
	if err != nil {
		errRspHeader(&w, "GET", err.Error(), http.StatusInternalServerError)
		return
	}

	pfdRespJSON, err = json.Marshal(pfdResp)
	if err != nil {
		errRspHeader(&w, "GET", err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(pfdRespJSON); err != nil {
		errRspHeader(&w, "GET", err.Error(), http.StatusInternalServerError)
		return
	}
}
