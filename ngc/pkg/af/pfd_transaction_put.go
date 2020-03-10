// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func putPfdTransaction(cliCtx context.Context, pfdTs PfdManagement,
	afCtx *Context, sID string) (PfdManagement,
	*http.Response, []byte, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsRet, resp, respBody, err :=
		cli.PfdManagementPutAPI.PfdTransactionPut(cliCtx,
			afCtx.cfg.AfID, sID, pfdTs)

	if err != nil {
		return PfdManagement{}, resp, respBody, err
	}
	return tsRet, resp, respBody, nil
}

// PutPfdTransaction function - To update the PFD transcation
func PutPfdTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		err              error
		pfdTs            PfdManagement
		resp             *http.Response
		pfdTransactionID string
		pfdRsp           PfdManagement
		pfdRespJSON      []byte
		resBody          []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		errRspHeader(&w, "PUT", "af-ctx retrieved from request is nil",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&pfdTs); err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusBadRequest)
		return
	}

	pfdTransactionID = getPfdTransIDFromURL(r)

	pfdRsp, resp, resBody, err = putPfdTransaction(cliCtx, pfdTs, afCtx,
		pfdTransactionID)
	// TBD how to validate the PUT response
	if err != nil {
		log.Errf("Pfd Management Put : %s", err.Error())
		if resp != nil {
			w.WriteHeader(getStatusCode(resp))
			if _, err = w.Write(resBody); err != nil {
				errRspHeader(&w, "PUT", err.Error(),
					http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	// Updating the Self Link and Applications Self Link in AF

	self, err := updateSelfLink(afCtx.cfg, r, pfdRsp)
	if err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusInternalServerError)
		return
	}
	pfdRsp.Self = Link(self)
	err = updateAppsLink(afCtx.cfg, r, pfdRsp)
	if err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusInternalServerError)
		return
	}

	pfdRespJSON, err = json.Marshal(pfdRsp)
	if err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(pfdRespJSON); err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusInternalServerError)
		return
	}

}
