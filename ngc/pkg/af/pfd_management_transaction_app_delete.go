// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

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
		log.Errf("Pfd App Transaction delete: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pfdTrans, err = getPfdTransIDFromURL(r)
	if err != nil {
		log.Errf("Pfd App Transaction delete %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appID, err = getPfdAppIDFromURL(r)
	if err != nil {
		log.Errf("Pfd App Transaction delete %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err = deletePfdAppTransaction(cliCtx, afCtx, pfdTrans, appID)
	if err != nil {
		log.Errf("Pfd App Transaction delete %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	w.WriteHeader(resp.StatusCode)
}
