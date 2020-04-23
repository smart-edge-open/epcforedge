// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	_context "context"
	"encoding/json"
	"errors"
	_ioutil "io/ioutil"
	_nethttp "net/http"
)

// Linger please
var (
	_ _context.Context
)

// PolicyAuthAppSessionAPIService policyAuthService type
type PolicyAuthAppSessionAPIService policyAuthService

// PostAppSessions API handler
/*
 * PostAppSessions Creates a new Individual Application Session Context resource
 * @param ctx _context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionContext Contains the information for the creation the
 *   resource
 * @return AppSessionContext
 */
func (a *PolicyAuthAppSessionAPIService) PostAppSessions(ctx _context.Context,
	appSessionContext AppSessionContext) (
	AppSessionContext, *ProblemDetails, *_nethttp.Response, error) {

	var (
		httpMethod = _nethttp.MethodPost
		postBody   interface{}
		retVal     AppSessionContext
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.PcfHostname +
		a.client.cfg.PcfPort + a.client.cfg.PolicyAuthBasePath +
		"/app-sessions"

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = "application/json"
	headerParams["Accept"] = "appication/json"

	postBody = &appSessionContext

	r, err := a.client.prepareRequest(ctx, path, httpMethod, postBody,
		headerParams)
	if err != nil {
		return retVal, nil, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, nil, httpResponse, err
	}

	respBody, err := _ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		return retVal, nil, httpResponse, err
	}

	if httpResponse.StatusCode == 201 {
		err = json.Unmarshal(respBody, &retVal)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"PostAppSession: " + err.Error())
			httpResponse.StatusCode = 500
		}
		return retVal, nil, httpResponse, err
	}

	switch httpResponse.StatusCode {
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"PostAppSession: " + err.Error())
			httpResponse.StatusCode = 500
			return retVal, nil, httpResponse, err
		}
		return retVal, v, httpResponse, err
	}

	err = errors.New(string(respBody))
	log.Errf("PostAppSess- PCF Returned Error: " +
		string(respBody))
	return retVal, nil, httpResponse, err
}
