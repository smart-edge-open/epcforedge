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

// PfdManagementTransactionPostAPIService type
type PfdManagementTransactionPostAPIService service

func (a *PfdManagementTransactionPostAPIService) handlePfdPostResponse(
	pfdTrans *PfdManagement, r *http.Response,
	body []byte) error {

	if r.StatusCode == 201 {

		err := json.Unmarshal(body, pfdTrans)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
		}
		return err
	}

	return handlePfdPostPutPatchErrorResp(r, body)
}

/*
PfdTransactionPost Creates a new PFD management resource
Creates a new PFD management resource
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF
 * @param body Request to create a new PFD management resource

@return PfdManagement
*/
func (a *PfdManagementTransactionPostAPIService) PfdTransactionPost(
	ctx context.Context, afID string, body PfdManagement) (PfdManagement,
	*http.Response, []byte, error) {

	var (
		method   = strings.ToUpper("Post")
		postBody interface{}
		ret      PfdManagement
		respBody []byte
	)

	// create path and map variables

	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFPFDBasePath + "/" + afID +
		"/transactions"

	headerParams := make(map[string]string)

	headerParams["Content-Type"] = contentTypePfd
	headerParams["Accept"] = contentTypePfd

	// body params
	postBody = &body
	r, err := a.client.prepareRequest(ctx, path, method,
		postBody, headerParams)
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

	if err = a.handlePfdPostResponse(&ret, resp,
		respBody); err != nil {

		return ret, resp, respBody, err
	}

	return ret, resp, respBody, err
}
