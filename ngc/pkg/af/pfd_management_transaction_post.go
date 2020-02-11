// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

func createPfdTransaction(cliCtx context.Context, pfdTrans PfdManagement,
	afCtx *Context) (PfdManagement, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	pfdResp, resp, err := cli.PfdManagementPostAPI.PfdTransactionPost(cliCtx,
		afCtx.cfg.AfID, pfdTrans)

	if err != nil {
		return PfdManagement{}, resp, err
	}
	return pfdResp, resp, nil
}

// CreatePfdTransaction function - Create a PFD transaction
func CreatePfdTransaction(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		pfdTrans    PfdManagement
		resp        *http.Response
		url         *url.URL
		pfdRespJSON []byte
		pfdRsp      PfdManagement
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Pfd Management Transaction create: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&pfdTrans); err != nil {
		log.Errf("Pfd Management Transcation create: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdRsp, resp, err = createPfdTransaction(cliCtx, pfdTrans, afCtx)
	// TBD what to validate in response
	if err != nil {
		log.Errf("Pfd Management Transaction create: %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	if url, err = resp.Location(); err != nil {
		log.Errf("Pfd Management Transaction create: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdRespJSON, err = json.Marshal(pfdRsp)
	if err != nil {
		log.Errf("Pfd Management create : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url.String())
	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(pfdRespJSON); err != nil {
		log.Errf("Pfd Management create: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
