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

// PolicyAuthEventSubsAPIService EventsSubscriptionDocumentApi service
type PolicyAuthEventSubsAPIService policyAuthService

func handleUpdateEventResp(respBody []byte, statusCode int,
	a *PolicyAuthEventSubsAPIService) (
	retVal EventSubscResponse, err error) {

	var (
		eventSubscResp EventsSubscReqData
		evsNotifResp   EventsNotification
	)

	switch statusCode {
	case 200, 201:
		err = json.Unmarshal(respBody, &eventSubscResp)
		if err == nil {
			retVal.eventSubscReq = &eventSubscResp
			return retVal, nil
		}
		err = json.Unmarshal(respBody, &evsNotifResp)
		if err == nil {
			retVal.evsNotif = &evsNotifResp
			return retVal, nil
		}
		log.Errf("Error in unmarshalling response body, " +
			"UpdateEventSubsc: " + err.Error())
		return retVal, err

	case 204:
		return retVal, nil

	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		if statusCode == 401 {
			validatePAAuthToken(a.client)
		}
		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"UpdateEventSubsc: " + err.Error())
			return retVal, err
		}
		retVal.probDetails = v
		return retVal, err
	}

	err = errors.New(string(respBody))
	if err == nil {
		err = errors.New("No response from PCF")
	}

	log.Errf("UpdateEventSubsc- PCF Returned Error: " +
		string(respBody))
	return retVal, err
}

/*
UpdateEventsSubsc creates or modifies an Events Subscription subresource
 * @param ctx _context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the Events Subscription resource
 * @param eventsSubscReqData Creation or modification of an Events Subscription
 * resource.
@return AnyOfEventsSubscReqDataEventsNotification
*/
func (a *PolicyAuthEventSubsAPIService) UpdateEventsSubsc(ctx _context.Context,
	appSessionID string, eventSubscReq *EventsSubscReqData) (
	EventSubscResponse, *_nethttp.Response, error) {

	var (
		httpMethod = _nethttp.MethodPut
		reqBody    interface{}
		retVal     EventSubscResponse
	)

	// create path and map variables
	path := a.client.rootURI + "/app-sessions/" + appSessionID +
		"/events-subscription"

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
		return retVal, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, httpResponse, err
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
		return retVal, httpResponse, err
	}

	retVal, err = handleUpdateEventResp(respBody, httpResponse.StatusCode, a)
	return retVal, httpResponse, err

}

/*
DeleteEventsSubsc deletes the Events Subscription subresource
 * @param ctx _context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the Individual Application Session
 * Context resource
*/
func (a *PolicyAuthEventSubsAPIService) DeleteEventsSubsc(ctx _context.Context,
	appSessionID string) (*ProblemDetails, *_nethttp.Response, error) {
	var (
		httpMethod = _nethttp.MethodDelete
		reqBody    interface{}
		retVal     *ProblemDetails
	)

	// create path and map variables
	path := a.client.rootURI + "/app-sessions/" + appSessionID +
		"/events-subscription"

	headerParams := make(map[string]string)
	headerParams["Accept"] = contentTypeJSON

	r, err := a.client.prepareRequest(ctx, path, httpMethod, reqBody,
		headerParams)
	if err != nil {
		return retVal, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, httpResponse, err
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
		return retVal, httpResponse, err
	}

	switch httpResponse.StatusCode {
	case 204:
		return retVal, httpResponse, nil

	case 400, 401, 403, 404, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			validatePAAuthToken(a.client)
		}
		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"DeleteEventSubsc: " + err.Error())
			return retVal, httpResponse, err
		}
		retVal = v
		return retVal, httpResponse, err
	}

	err = errors.New(string(respBody))
	if err == nil {
		err = errors.New("Empty response from PCF")
	}

	log.Errf("DeleteEventSubs- PCF Returned Error: " +
		string(respBody))
	return retVal, httpResponse, err
}
