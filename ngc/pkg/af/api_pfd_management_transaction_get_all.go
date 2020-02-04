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

// PfdManagementTransactionGetAllAPIService type
type PfdManagementTransactionGetAllAPIService service

func (a *PfdManagementTransactionGetAllAPIService) handlePfdGetAllResponse(
	ts *[]PfdManagement, r *http.Response, body []byte) error {

	if r.StatusCode == 200 {
		err := json.Unmarshal(body, ts)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
		}
		return err
	}

	return handleGetErrorResp(r, body)
}

/*
PfdTransactionsGetAll read all of the active
Pfd transactions for the AF
read all of the active pfd transactions for the AF
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF

@return []PfdManagement
*/
func (a *PfdManagementTransactionGetAllAPIService) PfdTransactionsGetAll(
	ctx context.Context, afID string) ([]PfdManagement,
	*http.Response, error) {
	var (
		method  = strings.ToUpper("Get")
		getBody interface{}
		ret     []PfdManagement
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFPFDPath +
		"/{afId}/transactions"

	path = strings.Replace(path, "{"+"afId"+"}",
		fmt.Sprintf("%v", afID), -1)

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

	if err = a.handlePfdGetAllResponse(&ret, resp,
		respBody); err != nil {
		return ret, resp, err
	}

	return ret, resp, nil
}
