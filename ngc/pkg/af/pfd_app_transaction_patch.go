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
	*http.Response, []byte, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsRet, resp, respBody, err :=
		cli.PfdManagementAppPatchAPI.PfdAppTransactionPatch(
			cliCtx, afCtx.cfg.AfID, pfdID, appID, pfdData)

	if err != nil {
		return PfdData{}, resp, respBody, err
	}
	return tsRet, resp, respBody, nil
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
		respBody    []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		errRspHeader(&w, "APP PATCH", "af-ctx retrieved from request is nil",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&pfdData); err != nil {
		errRspHeader(&w, "APP-PATCH", err.Error(),
			http.StatusBadRequest)
		return
	}

	pfdTransID = getPfdTransIDFromURL(r)

	appID = getPfdAppIDFromURL(r)

	pfdRsp, resp, respBody, err = patchPfdAppTransaction(cliCtx, pfdData, afCtx,
		pfdTransID, appID)
	if err != nil {
		log.Errf("Pfd Management Application patch : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		if _, err = w.Write(respBody); err != nil {
			errRspHeader(&w, "APP-PATCH", err.Error(),
				http.StatusInternalServerError)
			return
		}
		return
	}

	// Updating the Self Application Link
	self, err := updateAppLink(afCtx.cfg, r, pfdRsp)
	if err != nil {
		errRspHeader(&w, "APP-PATCH", err.Error(),
			http.StatusInternalServerError)
		return
	}
	pfdRsp.Self = Link(self)

	pfdRespJSON, err = json.Marshal(pfdRsp)
	if err != nil {
		errRspHeader(&w, "APP-PATCH", err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(pfdRespJSON); err != nil {
		errRspHeader(&w, "APP-PATCH", err.Error(),
			http.StatusInternalServerError)
		return
	}

}
