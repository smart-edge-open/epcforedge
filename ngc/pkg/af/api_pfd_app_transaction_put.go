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

// PfdManagementTransactionAppPutAPIService type
type PfdManagementTransactionAppPutAPIService service

func (a *PfdManagementTransactionAppPutAPIService) handlePfdAppPutResponse(
	pfdTs *PfdData, r *http.Response,
	body []byte) error {

	if r.StatusCode == 200 {
		err := json.Unmarshal(body, pfdTs)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
		}
		return err
	}

	return handlePfdPostPutPatchErrorResp(r, body)

}

/*
PfdAppTransactionPut updates an existing PFD transaction for an external app ID
source
Updates an existing pfd transaction resource for an external application id
 * @param ctx context.Context - for authentication, logging, cancellation,
 * 	deadlines, tracing, etc. Passed from http.Request or
 *	context.Background().
 * @param afID Identifier of the AF
 * @param pfdTransactionID Identifier of the pfd Transaction resource
 * @param appID Identifier of the external application id
 * @param body Provides the  PfdData structure identified by
 *	pfdTransaction ID and appID

@return PfdData
*/
func (a *PfdManagementTransactionAppPutAPIService) PfdAppTransactionPut(
	ctx context.Context, afID string, pfdTransaction string, appID string,
	body PfdData) (PfdData, *http.Response, []byte, error) {

	var (
		method   = strings.ToUpper("Put")
		putBody  interface{}
		ret      PfdData
		respBody []byte
	)
	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFPFDBasePath + "/" + afID +
		"/transactions/" + pfdTransaction + "/applications/" + appID

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

	if err = a.handlePfdAppPutResponse(&ret, resp,
		respBody); err != nil {

		return ret, resp, respBody, err
	}

	return ret, resp, respBody, nil
}
