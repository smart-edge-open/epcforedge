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

// PolicyAuthIndividualAppSessAPIService PolicyAuthService type
type PolicyAuthIndividualAppSessAPIService policyAuthService

// DeleteAppSession API handler
/*
 * DeleteAppSession Deletes an existing Individual Application Session Context
 * @param ctx _context.Context - for authentication, logging, cancellation,
 *    deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the Individual Application Session
 *    Context resource
 * @param optional nil or *DeleteAppSessionOpts - Optional Parameters:
 * @param "EventsSubscReqData" (optional.Interface of EventsSubscReqData) -
 *   deletion of the Individual Application Session Context resource, req
 *   notification
 * @return AppSessionContext
 */
func (a *PolicyAuthIndividualAppSessAPIService) DeleteAppSession(
	ctx _context.Context, appSessionID string,
	eventSubscReq *EventsSubscReqData) (
	AppSessionContext, *ProblemDetails, *_nethttp.Response, error) {

	var (
		httpMethod = _nethttp.MethodPost
		reqBody    interface{}
		retVal     AppSessionContext
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.PcfHostname +
		a.client.cfg.PcfPort + a.client.cfg.PolicyAuthBasePath +
		"/app-sessions/" + appSessionID + "/delete"

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = contentTypeJSON
	headerParams["Accept"] = contentTypeJSON

	// body params
	if eventSubscReq != nil {
		reqBody = eventSubscReq
	}

	r, err := a.client.prepareRequest(ctx, path, httpMethod, reqBody,
		headerParams)
	if err != nil {
		return retVal, nil, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, nil, httpResponse, err
	}

	respBody, err := _ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, nil, httpResponse, err
	}

	if httpResponse.StatusCode == 200 {
		err = json.Unmarshal(respBody, &retVal)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"DeleteAppSession: " + err.Error())
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
				"DeleteAppSession: " + err.Error())
			httpResponse.StatusCode = 500
			return retVal, nil, httpResponse, err
		}
		return retVal, v, httpResponse, err
	}

	err = errors.New(string(respBody))
	log.Errf("DeleteAppSess- PCF Returned Error: " +
		string(respBody))
	return retVal, nil, httpResponse, err

}

// GetAppSession API Handler
/*
 * GetAppSession Reads an existing Individual Application Session Context
 * @param ctx _context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the resource
 * @return AppSessionContext for 200 resp, otherwise ProblemDetails
 */
func (a *PolicyAuthIndividualAppSessAPIService) GetAppSession(
	ctx _context.Context, appSessionID string) (
	AppSessionContext, *ProblemDetails, *_nethttp.Response, error) {
	var (
		httpMethod = _nethttp.MethodGet
		reqBody    interface{}
		retVal     AppSessionContext
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.PcfHostname +
		a.client.cfg.PcfPort + a.client.cfg.PolicyAuthBasePath +
		"/app-sessions/" + appSessionID

	headerParams := make(map[string]string)
	headerParams["Accept"] = contentTypeJSON

	r, err := a.client.prepareRequest(ctx, path, httpMethod, reqBody,
		headerParams)
	if err != nil {
		return retVal, nil, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, nil, httpResponse, err
	}

	respBody, err := _ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, nil, httpResponse, err
	}

	if httpResponse.StatusCode == 200 {
		err = json.Unmarshal(respBody, &retVal)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"GetAppSession: " + err.Error())
			httpResponse.StatusCode = 500
		}
		return retVal, nil, httpResponse, err
	}

	switch httpResponse.StatusCode {
	case 400, 401, 403, 404, 406, 429, 500, 503:
		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"GetAppSession: " + err.Error())
			httpResponse.StatusCode = 500
			return retVal, nil, httpResponse, err
		}
		return retVal, v, httpResponse, err
	}

	err = errors.New(string(respBody))
	log.Errf("GetAppSess- PCF Returned Error: " +
		string(respBody))
	return retVal, nil, httpResponse, err
}

// ModAppSession API handler
/*
 * ModAppSession Modifies an existing Individual Application Session Context
 * @param ctx _context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the resource
 * @param appSessionContextUpdateData modification of the resource.
 * @return AppSessionContext on 200 resp otherwise ProbleDetails
 */
func (a *PolicyAuthIndividualAppSessAPIService) ModAppSession(
	ctx _context.Context, appSessionID string,
	appSessionContextUpdateData AppSessionContextUpdateData) (
	AppSessionContext, *ProblemDetails, *_nethttp.Response, error) {

	var (
		httpMethod = _nethttp.MethodPatch
		patchBody  interface{}
		retVal     AppSessionContext
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.PcfHostname +
		a.client.cfg.PcfPort + a.client.cfg.PolicyAuthBasePath +
		"/app-sessions/" + appSessionID

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = "application/json"
	headerParams["Accept"] = "application/json"

	// body params
	patchBody = &appSessionContextUpdateData
	r, err := a.client.prepareRequest(ctx, path, httpMethod, patchBody,
		headerParams)
	if err != nil {
		return retVal, nil, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, nil, httpResponse, err
	}

	respBody, err := _ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, nil, httpResponse, err
	}

	if httpResponse.StatusCode == 200 {
		err = json.Unmarshal(respBody, &retVal)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"ModAppSession: " + err.Error())
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
				"ModAppSession: " + err.Error())
			httpResponse.StatusCode = 500
			return retVal, nil, httpResponse, err
		}
		return retVal, v, httpResponse, err
	}

	err = errors.New(string(respBody))
	log.Errf("ModAppSess- PCF Returned Error: " +
		string(respBody))
	return retVal, nil, httpResponse, err
}
