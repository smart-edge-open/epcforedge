// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"net/http"
)

func deletePfdAppTransaction(cliCtx context.Context, afCtx *Context,
	pfdID string, appID string) (*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	resp, err := cli.PfdManagementAppDeleteAPI.PfdAppTransactionDelete(cliCtx,
		afCtx.cfg.AfID, pfdID, appID)

	if err != nil {
		return resp, err
	}
	return resp, nil

}

// DeletePfdAppTransaction function
func DeletePfdAppTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		resp     *http.Response
		pfdTrans string
		appID    string
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		errRspHeader(&w, "APP DELETE", "af-ctx retrieved from request is nil",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pfdTrans = getPfdTransIDFromURL(r)

	appID = getPfdAppIDFromURL(r)

	resp, err = deletePfdAppTransaction(cliCtx, afCtx, pfdTrans, appID)
	if err != nil {
		errRspHeader(&w, "APP DELETE", err.Error(), resp.StatusCode)
		return
	}

	w.WriteHeader(resp.StatusCode)
}
