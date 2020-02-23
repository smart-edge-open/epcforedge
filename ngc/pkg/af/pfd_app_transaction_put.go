// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func putPfdAppTransaction(cliCtx context.Context, pfdTs PfdData,
	afCtx *Context, pfdID string, appID string) (PfdData,
	*http.Response, []byte, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsRet, resp, respBody, err :=
		cli.PfdManagementAppPutAPI.PfdAppTransactionPut(cliCtx,
			afCtx.cfg.AfID, pfdID, appID, pfdTs)

	if err != nil {
		return PfdData{}, resp, respBody, err
	}
	return tsRet, resp, respBody, nil
}

// PutPfdAppTransaction function
func PutPfdAppTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		err              error
		pfdTs            PfdData
		resp             *http.Response
		pfdTransactionID string
		appID            string
		pfdRsp           PfdData
		pfdRespJSON      []byte
		respBody         []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Pfd Management App Put: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&pfdTs); err != nil {
		log.Errf("Pfd Management App Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdTransactionID, err = getPfdTransIDFromURL(r)
	if err != nil {
		log.Errf("Pfd Management App Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appID, err = getPfdAppIDFromURL(r)
	if err != nil {
		log.Errf("Pfd Management App Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdRsp, resp, respBody, err = putPfdAppTransaction(cliCtx, pfdTs, afCtx,
		pfdTransactionID, appID)
	// TBD how to validate the PUT response
	if err != nil {
		log.Errf("Pfd Management App Put : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		if _, err = w.Write(respBody); err != nil {
			log.Errf("Pfd Management App put: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	// Updating the Self Application Link
	self, err := updateAppLink(afCtx.cfg, r, pfdRsp)
	if err != nil {
		log.Errf("Pfd Management Application Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pfdRsp.Self = Link(self)

	pfdRespJSON, err = json.Marshal(pfdRsp)
	if err != nil {
		log.Errf("Pfd Management Application Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(pfdRespJSON); err != nil {
		log.Errf("Pfd Management Application Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
