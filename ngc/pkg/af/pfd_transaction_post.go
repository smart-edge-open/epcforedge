// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

func createPfdTransaction(cliCtx context.Context, pfdTrans PfdManagement,
	afCtx *Context) (PfdManagement, *http.Response, []byte, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	pfdResp, resp, respBody, err :=
		cli.PfdManagementPostAPI.PfdTransactionPost(cliCtx,
			afCtx.cfg.AfID, pfdTrans)

	if err != nil {
		return PfdManagement{}, resp, respBody, err
	}
	return pfdResp, resp, respBody, nil
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
		respBody    []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		errRspHeader(&w, "POST", "af-ctx retrieved from request is nil",
			http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&pfdTrans); err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}

	pfdRsp, resp, respBody, err = createPfdTransaction(cliCtx, pfdTrans, afCtx)

	if err != nil {
		log.Errf("Pfd Management Transaction create: %s", err.Error())
		w.WriteHeader(getStatusCode(resp))

		if _, err = w.Write(respBody); err != nil {
			log.Errf("PFD create error handling error reposnse")
			return
		}

		return
	}

	if url, err = resp.Location(); err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}

	// Updating the location url, Self Link and Application Self Link in AF
	afURL := updatePfdURL(afCtx.cfg, r, url.String())

	self, err := updateSelfLink(afCtx.cfg, r, pfdRsp)
	if err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}
	pfdRsp.Self = Link(self)
	err = updateAppsLink(afCtx.cfg, r, pfdRsp)
	if err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}
	pfdRespJSON, err = json.Marshal(pfdRsp)
	if err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", afURL)
	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(pfdRespJSON); err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}

}
