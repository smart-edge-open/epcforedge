/*
 * Npcf_PolicyAuthorization Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package af

import (
	_context "context"
	_ioutil "io/ioutil"
	_nethttp "net/http"
	/*
		_neturl "net/url"
			"strings"
	*/
	"github.com/antihax/optional"
)

// Linger please
var (
	_ _context.Context
)

// PolicyAuthIndividualAppSessApiService IndividualApplicationSessionContextDocumentApi service
type PolicyAuthIndividualAppSessApiService policyAuthService

// DeleteAppSessionOpts Optional parameters for the method 'DeleteAppSession'
type DeleteAppSessionOpts struct {
	EventsSubscReqData optional.Interface
}

/*
DeleteAppSession Deletes an existing Individual Application Session Context
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionId string identifying the Individual Application Session Context resource
 * @param optional nil or *DeleteAppSessionOpts - Optional Parameters:
 * @param "EventsSubscReqData" (optional.Interface of EventsSubscReqData) -  deletion of the Individual Application Session Context resource, req notification
@return AppSessionContext
*/
func (a *PolicyAuthIndividualAppSessApiService) DeleteAppSession(ctx _context.Context, appSessionId string, eventSubscReq *EventsSubscReqData) (AppSessionContext, *_nethttp.Response, error) {
	var (
		httpMethod = _nethttp.MethodPost
		reqBody    interface{}
		retVal     AppSessionContext
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.PcfHostname +
		a.client.cfg.PcfPort + a.client.cfg.PolicyAuthBasePath +
		"/app-sessions/" + appSessionId + "/delete"

	headerParams := make(map[string]string)

	// to determine the Content-Type header
	httpContentTypes := []string{"application/json"}

	// set Content-Type header
	httpContentType := selectHeaderContentType(httpContentTypes)
	if httpContentType != "" {
		headerParams["Content-Type"] = httpContentType
	}

	// to determine the Accept header
	httpHeaderAccepts := []string{"application/json", "application/problem+json"}

	// set Accept header
	httpHeaderAccept := selectHeaderAccept(httpHeaderAccepts)
	if httpHeaderAccept != "" {
		headerParams["Accept"] = httpHeaderAccept
	}
	// body params
	if eventSubscReq != nil {
		reqBody = eventSubscReq
	}

	r, err := a.client.prepareRequest(ctx, path, httpMethod, reqBody, headerParams)
	if err != nil {
		return retVal, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, httpResponse, err
	}

	respBody, err := _ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		return retVal, httpResponse, err
	}

	newErr := GenericError{
		body:  respBody,
		error: httpResponse.Status,
	}

	if httpResponse.StatusCode == 200 {
		var v AppSessionContext
		err = a.client.decode(&v, respBody, httpResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return retVal, httpResponse, newErr
		}
		newErr.model = v
		return retVal, httpResponse, newErr
	}

	switch httpResponse.StatusCode {
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			var v ProblemDetails
			err = a.client.decode(&v, respBody, httpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return retVal, httpResponse, newErr
			}
			newErr.model = v
			return retVal, httpResponse, newErr
		}
	}

	err = a.client.decode(&retVal, respBody, httpResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr.error = err.Error()
		return retVal, httpResponse, newErr
	}

	return retVal, httpResponse, nil
}

/*
GetAppSession Reads an existing Individual Application Session Context
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionId string identifying the resource
@return AppSessionContext
*/
func (a *PolicyAuthIndividualAppSessApiService) GetAppSession(ctx _context.Context, appSessionId string) (AppSessionContext, *_nethttp.Response, error) {
	var (
		httpMethod = _nethttp.MethodGet
		reqBody    interface{}
		retVal     AppSessionContext
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.PcfHostname +
		a.client.cfg.PcfPort + a.client.cfg.PolicyAuthBasePath +
		"/app-sessions/" + appSessionId

	headerParams := make(map[string]string)

	// to determine the Content-Type header
	httpContentTypes := []string{}

	// set Content-Type header
	httpContentType := selectHeaderContentType(httpContentTypes)
	if httpContentType != "" {
		headerParams["Content-Type"] = httpContentType
	}

	// to determine the Accept header
	httpHeaderAccepts := []string{"application/json", "application/problem+json"}

	// set Accept header
	httpHeaderAccept := selectHeaderAccept(httpHeaderAccepts)
	if httpHeaderAccept != "" {
		headerParams["Accept"] = httpHeaderAccept
	}

	r, err := a.client.prepareRequest(ctx, path, httpMethod, reqBody, headerParams)
	if err != nil {
		return retVal, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, httpResponse, err
	}

	respBody, err := _ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		return retVal, httpResponse, err
	}

	newErr := GenericError{
		body:  respBody,
		error: httpResponse.Status,
	}
	if httpResponse.StatusCode == 200 {
		var v AppSessionContext
		err = a.client.decode(&v, respBody, httpResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return retVal, httpResponse, newErr
		}
		newErr.model = v
		return retVal, httpResponse, newErr
	}

	switch httpResponse.StatusCode {
	case 400, 401, 403, 404, 406, 429, 500, 503:
		var v ProblemDetails
		err = a.client.decode(&v, respBody, httpResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return retVal, httpResponse, newErr
		}
		newErr.model = v
		return retVal, httpResponse, newErr
	default:
		err = a.client.decode(&retVal, respBody, httpResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return retVal, httpResponse, newErr
		}
	}
	return retVal, httpResponse, nil
}

/*
 */

/*
ModAppSession Modifies an existing Individual Application Session Context
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionId string identifying the resource
 * @param appSessionContextUpdateData modification of the resource.
@return AppSessionContext
*/
func (a *PolicyAuthIndividualAppSessApiService) ModAppSession(ctx _context.Context, appSessionId string, appSessionContextUpdateData AppSessionContextUpdateData) (AppSessionContext, *_nethttp.Response, error) {
	var (
		httpMethod = _nethttp.MethodPatch
		patchBody  interface{}
		retVal     AppSessionContext
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.PcfHostname +
		a.client.cfg.PcfPort + a.client.cfg.PolicyAuthBasePath +
		"/app-sessions/" + appSessionId

	headerParams := make(map[string]string)

	// to determine the Content-Type header
	httpContentTypes := []string{"application/merge-patch+json"}

	// set Content-Type header
	httpContentType := selectHeaderContentType(httpContentTypes)
	if httpContentType != "" {
		headerParams["Content-Type"] = httpContentType
	}

	// to determine the Accept header
	httpHeaderAccepts := []string{"application/json", "application/problem+json"}

	// set Accept header
	httpHeaderAccept := selectHeaderAccept(httpHeaderAccepts)
	if httpHeaderAccept != "" {
		headerParams["Accept"] = httpHeaderAccept
	}
	// body params
	patchBody = &appSessionContextUpdateData
	r, err := a.client.prepareRequest(ctx, path, httpMethod, patchBody, headerParams)
	if err != nil {
		return retVal, nil, err
	}

	httpResponse, err := a.client.callAPI(r)
	if err != nil || httpResponse == nil {
		return retVal, httpResponse, err
	}

	respBody, err := _ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		return retVal, httpResponse, err
	}

	newErr := GenericError{
		body:  respBody,
		error: httpResponse.Status,
	}
	if httpResponse.StatusCode == 200 {
		var v AppSessionContext
		err = a.client.decode(&v, respBody, httpResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return retVal, httpResponse, newErr
		}
		newErr.model = v
		return retVal, httpResponse, newErr
	}
	switch httpResponse.StatusCode {
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			var v ProblemDetails
			err = a.client.decode(&v, respBody, httpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return retVal, httpResponse, newErr
			}
			newErr.model = v
			return retVal, httpResponse, newErr
		}
	}

	err = a.client.decode(&retVal, respBody, httpResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr.error = err.Error()
		return retVal, httpResponse, newErr
	}

	return retVal, httpResponse, nil
}
