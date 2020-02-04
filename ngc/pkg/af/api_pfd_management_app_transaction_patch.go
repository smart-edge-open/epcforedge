// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Linger please
var (
	_ context.Context
)

// PfdManagementTransactionAppPatchAPIService type
type PfdManagementTransactionAppPatchAPIService service

func (a *PfdManagementTransactionAppPatchAPIService) handleAppPatchResponse(
	ts *PfdData, r *http.Response,
	body []byte) error {

	if r.StatusCode == 200 {
		err := json.Unmarshal(body, ts)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
		}
		return err
	}

	return handlePostPutPatchErrorResp(r, body)

}

/*
PfdAppTransactionPatch Updates an existing pfd transaction for an app ID
 * @param ctx context.Context - for authentication, logging, cancellation,
 * 	deadlines, tracing, etc. Passed from http.Request or
 *	context.Background().
 * @param afID Identifier of the AF
 * @param pfdTransId Identifier of the pfdData
 * @param appID Identifier of the external application
 * @param body Provides a pfd data for pfd management identified by
 *	transaction id and appID

@return PfdData
*/
func (a *PfdManagementTransactionAppPatchAPIService) PfdAppTransactionPatch(
	ctx context.Context, afID string, pfdTransID string, appID string,
	body PfdData) (PfdData, *http.Response, error) {

	var (
		method    = strings.ToUpper("Patch")
		patchBody interface{}
		ret       PfdData
	)

	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFPFDPath +
		"/{afId}/transactions/{transactionId}/application/{applicationId}"

	path = strings.Replace(path,
		"{"+"afId"+"}", fmt.Sprintf("%v", afID), -1)
	path = strings.Replace(path,
		"{"+"transactionId"+"}", fmt.Sprintf("%v", pfdTransID), -1)
	path = strings.Replace(path,
		"{"+"applicationId"+"}", fmt.Sprintf("%v", appID), -1)

	headerParams := make(map[string]string)

	// to determine the Content-Type header
	contentTypes := []string{"application/json"}

	// set Content-Type header
	contentType := selectHeaderContentType(contentTypes)
	if contentType != "" {
		headerParams["Content-Type"] = contentType
	}

	// to determine the Accept header
	headerAccepts := []string{"application/json"}

	// set Accept header
	headerAccept := selectHeaderAccept(headerAccepts)
	if headerAccept != "" {
		headerParams["Accept"] = headerAccept
	}
	// body params
	patchBody = &body
	r, err := a.client.prepareRequest(ctx, path, method,
		patchBody, headerParams)
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

	if err = a.handleAppPatchResponse(&ret, resp,
		respBody); err != nil {

		return ret, resp, err
	}

	return ret, resp, nil
}
