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

// PfdManagementTransactionGetAPIService type
type PfdManagementTransactionGetAPIService service

func (a *PfdManagementTransactionGetAPIService) handlePfdGetResponse(
	pfdTs *PfdManagement, r *http.Response,
	body []byte) error {

	if r.StatusCode == 200 {
		err := json.Unmarshal(body, pfdTs)
		if err != nil {
			log.Errf("Error decoding response body %s: ", err.Error())
			r.StatusCode = 500
		}
		return err
	}
	return handleGetErrorResp(r, body)
}

/*
PfdTransactionGet Read an active pfd Management
for the AF and the pfd transcation id
Read an active pfd transactions for the AF and the pfd transaction Id
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF
 * @param pfd transaction Identifier of the PFD transaction resource

@return PfdManagenent
*/
func (a *PfdManagementTransactionGetAPIService) PfdTransactionGet(
	ctx context.Context, afID string, pfdTransaction string) (PfdManagement,
	*http.Response, error) {
	var (
		method  = strings.ToUpper("Get")
		getBody interface{}
		ret     PfdManagement
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFPFDBasePath + "/" + afID +
		"/transactions/" + pfdTransaction

	headerParams := make(map[string]string)

	headerParams["Content-Type"] = contentType
	headerParams["Accept"] = contentType

	r, err := a.client.prepareRequest(ctx, path, method,
		getBody, headerParams)
	if err != nil {
		return ret, nil, err
	}

	resp, err := a.client.callAPI(r)
	if err != nil || resp == nil {
		return ret, resp, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errf("response body was not closed properly")
		}
	}()

	if err != nil {
		log.Errf("http response body could not be read")
		return ret, resp, err
	}

	if err = a.handlePfdGetResponse(&ret, resp,
		respBody); err != nil {

		return ret, resp, err
	}

	return ret, resp, nil
}
