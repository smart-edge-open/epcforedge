// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func patchPfdAppTransaction(cliCtx context.Context, pfdData PfdData,
	afCtx *Context, pfdID string, appID string) (PfdData,
	*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsRet, resp, err := cli.PfdManagementAppPatchAPI.PfdAppTransactionPatch(
		cliCtx, afCtx.cfg.AfID, pfdID, appID, pfdData)

	if err != nil {
		return PfdData{}, resp, err
	}
	return tsRet, resp, nil
}

// PatchPfdAppTransaction function
func PatchPfdAppTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		pfdData     PfdData
		resp        *http.Response
		pfdTransID  string
		appID       string
		pfdRsp      PfdData
		pfdRespJSON []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Pfd Management Application patch: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&pfdData); err != nil {
		log.Errf("Pfd Management Application patch %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdTransID, err = getPfdTransIDFromURL(r)
	if err != nil {
		log.Errf("Pfd Management Application patch %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appID, err = getPfdAppIDFromURL(r)
	if err != nil {
		log.Errf("Pfd Management Application patch %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdRsp, resp, err = patchPfdAppTransaction(cliCtx, pfdData, afCtx,
		pfdTransID, appID)
	if err != nil {
		log.Errf("Pfd Management Application patch : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	// Updating the Self Application Link
	self, err := updateAppLink(afCtx.cfg, r, pfdRsp)
	if err != nil {
		log.Errf("Pfd Management Application Patch: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pfdRsp.Self = Link(self)

	pfdRespJSON, err = json.Marshal(pfdRsp)
	if err != nil {
		log.Errf("Pfd Management Application Patch: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(pfdRespJSON); err != nil {
		log.Errf("Pfd Management Application Patch: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
