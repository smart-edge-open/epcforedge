// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Linger please
var (
	_ context.Context
)

// PfdManagementTransactionPutAPIService type
type PfdManagementTransactionPutAPIService service

func (a *PfdManagementTransactionPutAPIService) handlePfdPutResponse(
	pfdTs *PfdManagement, r *http.Response,
	body []byte) error {

	if r.StatusCode == 200 {
		err := json.Unmarshal(body, pfdTs)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
			r.StatusCode = 500
		}
		return err
	}

	return handlePfdPostPutPatchErrorResp(r, body)

}

/*
PfdTransactionPut updates an existing PFD transaction
source
Updates an existing PFD transaction resource
 * @param ctx context.Context - for authentication, logging, cancellation,
 * 	deadlines, tracing, etc. Passed from http.Request or
 *	context.Background().
 * @param afID Identifier of the AF
 * @param pfdTransactionID Identifier of the pfd Transaction resource
 * @param body Provides the  pfd management structure identified by
 *	pfdTransaction ID

@return PfdManagement
*/
func (a *PfdManagementTransactionPutAPIService) PfdTransactionPut(
	ctx context.Context, afID string, pfdTransaction string,
	body PfdManagement) (PfdManagement, *http.Response, []byte, error) {

	var (
		method   = strings.ToUpper("Put")
		putBody  interface{}
		ret      PfdManagement
		respBody []byte
	)

	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFPFDBasePath + "/" + afID +
		"/transactions/" + pfdTransaction

	headerParams := make(map[string]string)

	headerParams["Content-Type"] = contentType
	headerParams["Accept"] = contentType

	// body params
	putBody = &body
	r, err := a.client.prepareRequest(ctx, path, method,
		putBody, headerParams)
	if err != nil {
		return ret, nil, respBody, err
	}

	resp, err := a.client.callAPI(r)
	if err != nil || resp == nil {
		return ret, resp, respBody, err
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errf("response body was not closed properly")
		}
	}()

	if err != nil {
		log.Errf("http response body could not be read")
		return ret, resp, respBody, err
	}

	if err = a.handlePfdPutResponse(&ret, resp,
		respBody); err != nil {

		return ret, resp, respBody, err
	}

	return ret, resp, respBody, nil
}
