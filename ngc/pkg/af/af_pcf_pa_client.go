// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// PolicyAuthAPIClient type
// In most cases there should be only one, shared, PolicyAuthAPIClient.
type PolicyAuthAPIClient struct {
	cfg *GenericCliConfig

	oAuth2Token       string
	httpClient        *http.Client
	rootURI           string
	rootNotifURI      string
	userAgent         string
	locationPrefixURI string
}

// callAPI do the request.
func (c *PolicyAuthAPIClient) callAPI(request *http.Request) (
	*http.Response, error) {

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func getLocationURI(httpResp *http.Response, c *PolicyAuthAPIClient) string {
	var (
		locURL *url.URL
		err    error
	)
	uri := c.locationPrefixURI
	if locURL, err = httpResp.Location(); err != nil {
		httpResp.StatusCode = http.StatusInternalServerError
		log.Errf("Error in getting location header: " + err.Error())
		return ""
	}

	res := strings.Split(locURL.String(), "app-sessions")
	if len(res) == 2 {
		uri += res[1]
	} else {
		log.Errf("Location header returned from PCF is INCORRECT")
	}
	return uri
}

// NewPolicyAuthAPIClient - helper func
/*
 * NewAPIClient creates a new API client. Basically create new http client if
 * not set in client configurations.
 */
func NewPolicyAuthAPIClient(cfg *Config) (*PolicyAuthAPIClient, error) {

	paCfg := cfg.CliPcfCfg
	c := &PolicyAuthAPIClient{}

	httpClient, err := genHTTPClient(paCfg)
	if err != nil {
		log.Errf("Error in generating http client")
		return nil, err
	}
	c.httpClient = httpClient

	c.rootURI = paCfg.Protocol + "://" + paCfg.Hostname + ":" + paCfg.Port +
		paCfg.BasePath
	c.userAgent = cfg.UserAgent

	c.rootNotifURI = "https://" + cfg.SrvCfg.Hostname +
		cfg.SrvCfg.NotifPort + paCfg.NotifURI

	c.locationPrefixURI = "https://" + cfg.SrvCfg.Hostname +
		cfg.SrvCfg.CNCAEndpoint + cfg.LocationPrefixPA

	c.cfg = paCfg

	// API Services
	if paCfg.OAuth2Support {
		token, err := getPcfOAuth2Token()
		if err != nil {
			log.Errf("Pcf OAuth2 Token retrieval error: " +
				err.Error())
			return nil, err
		}
		c.oAuth2Token = token
	}

	return c, nil
}

// prepareRequest build the request
func (c *PolicyAuthAPIClient) prepareRequest(
	ctx context.Context,
	path string, method string,
	reqBody interface{},
	headerParams map[string]string,
) (httpRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect reqBody type and post.
	if reqBody != nil {
		reqContentType := headerParams["Content-Type"]

		body, err = setBody(reqBody, reqContentType)
		if err != nil {
			return nil, err
		}
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		log.Errf("url parsing error, path = %s", path)
		return nil, err
	}

	// Generate a new request
	if body != nil {
		httpRequest, err = http.NewRequest(method, url.String(), body)
	} else {
		httpRequest, err = http.NewRequest(method, url.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		httpRequest.Header = headers
	}

	// Add the user agent to the request.
	httpRequest.Header.Add("User-Agent", c.userAgent)

	if ctx != nil {
		// add context to the request
		httpRequest = httpRequest.WithContext(ctx)

		// Walk through any authentication.
		if c.cfg.OAuth2Support {
			token := c.oAuth2Token
			if token == "" {
				err = errors.New("Nil Ouath2Token in " +
					"PcfApiClient Struct")
				return nil, err
			}
			httpRequest.Header.Add("Authorization", "Bearer "+token)
		}

	}
	return httpRequest, nil
}

// PostAppSessions API handler
/*
 * PostAppSessions Creates a new Individual Application Session Context resource
 * @param ctx context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionContext Contains the information for the creation the
 *   resource
 * @return AppSessionContext
 */
func (c *PolicyAuthAPIClient) PostAppSessions(ctx context.Context,
	appSessionContext AppSessionContext) (PcfPAResponse, error) {

	var (
		httpMethod = http.MethodPost
		postBody   interface{}
		retVal     PcfPAResponse
	)

	// create path and map variables
	path := c.rootURI + "/app-sessions"

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = contentTypeJSON
	headerParams["Accept"] = contentTypeJSON

	postBody = &appSessionContext

	r, err := c.prepareRequest(ctx, path, httpMethod, postBody,
		headerParams)
	if err != nil {
		return retVal, err
	}

	httpResponse, err := c.callAPI(r)
	if err != nil || httpResponse == nil {
		retVal.httpResp = httpResponse
		return retVal, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		retVal.httpResp = httpResponse
		return retVal, err
	}

	switch httpResponse.StatusCode {
	case 201:
		var asc *AppSessionContext = new(AppSessionContext)
		err = json.Unmarshal(respBody, asc)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"PostAppSession: " + err.Error())
			httpResponse.StatusCode = 500
		}
		retVal.locationURI = getLocationURI(httpResponse, c)
		retVal.httpResp = httpResponse
		retVal.appSessCtx = asc
		return retVal, err
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			validatePAAuthToken(c)
		}

		retVal.httpResp = httpResponse

		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"PostAppSession: " + err.Error())
			return retVal, err
		}

		retVal.probDetails = v
		return retVal, err
	}

	err = errors.New(string(respBody))
	log.Errf("PostAppSess- PCF Returned Error: " +
		string(respBody))
	retVal.httpResp = httpResponse
	return retVal, err
}

func handleApSessDeleteErrResp(v *ProblemDetails,
	respBody []byte) error {

	err := json.Unmarshal(respBody, v)
	if err != nil {
		log.Errf("Error in unmarshalling response body, " +
			"DeleteAppSession: " + err.Error())
		return err
	}
	return nil
}

// DeleteAppSession API handler
/*
 * DeleteAppSession Deletes an existing Individual Application Session Context
 * @param ctx context.Context - for authentication, logging, cancellation,
 *    deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the Individual Application Session
 *    Context resource
 * @param optional nil or *DeleteAppSessionOpts - Optional Parameters:
 * @param "EventsSubscReqData" (optional.Interface of EventsSubscReqData) -
 *   deletion of the Individual Application Session Context resource, req
 *   notification
 * @return AppSessionContext
 */
func (c *PolicyAuthAPIClient) DeleteAppSession(
	ctx context.Context, appSessionID string,
	eventSubscReq *EventsSubscReqData) (
	PcfPAResponse, error) {

	var (
		httpMethod = http.MethodPost
		reqBody    interface{}
		retVal     PcfPAResponse
	)

	// create path and map variables
	path := c.rootURI + "/app-sessions/" + appSessionID +
		"/delete"

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = contentTypeJSON
	headerParams["Accept"] = contentTypeJSON

	// body params
	if eventSubscReq != nil {
		reqBody = eventSubscReq
	}

	r, err := c.prepareRequest(ctx, path, httpMethod, reqBody,
		headerParams)
	if err != nil {
		return retVal, err
	}

	httpResponse, err := c.callAPI(r)
	retVal.httpResp = httpResponse
	if err != nil || httpResponse == nil {
		return retVal, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, err
	}

	switch httpResponse.StatusCode {
	case 200:
		var asc *AppSessionContext = new(AppSessionContext)
		err = json.Unmarshal(respBody, &asc)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"DeleteAppSession: " + err.Error())
			httpResponse.StatusCode = 500
		}
		retVal.appSessCtx = asc
		return retVal, err
	case 204:
		return retVal, nil
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			validatePAAuthToken(c)
		}
		var v *ProblemDetails = new(ProblemDetails)
		err = handleApSessDeleteErrResp(v, respBody)
		retVal.probDetails = v
		return retVal, err
	}

	err = errors.New(string(respBody))
	log.Errf("DeleteAppSess- PCF Returned Error: " +
		string(respBody))
	return retVal, err
}

// GetAppSession API Handler
/*
 * GetAppSession Reads an existing Individual Application Session Context
 * @param ctx context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the resource
 * @return AppSessionContext for 200 resp, otherwise ProblemDetails
 */
func (c *PolicyAuthAPIClient) GetAppSession(
	ctx context.Context, appSessionID string) (PcfPAResponse, error) {
	var (
		httpMethod = http.MethodGet
		reqBody    interface{}
		retVal     PcfPAResponse
	)

	// create path and map variables
	path := c.rootURI + "/app-sessions/" + appSessionID

	headerParams := make(map[string]string)
	headerParams["Accept"] = contentTypeJSON

	r, err := c.prepareRequest(ctx, path, httpMethod, reqBody,
		headerParams)
	if err != nil {
		return retVal, err
	}

	httpResponse, err := c.callAPI(r)
	retVal.httpResp = httpResponse
	if err != nil || httpResponse == nil {
		return retVal, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, err
	}

	switch httpResponse.StatusCode {
	case 200:
		var asc *AppSessionContext = new(AppSessionContext)
		err = json.Unmarshal(respBody, &asc)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"GetAppSession: " + err.Error())
			httpResponse.StatusCode = 500
		}
		retVal.appSessCtx = asc
		return retVal, err
	case 400, 401, 403, 404, 406, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			validatePAAuthToken(c)
		}
		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		retVal.httpResp = httpResponse
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"GetAppSession: " + err.Error())
			return retVal, err
		}
		retVal.probDetails = v
		return retVal, err
	}

	err = errors.New(string(respBody))
	log.Errf("GetAppSess- PCF Returned Error: " +
		string(respBody))
	return retVal, err
}

// ModAppSession API handler
/*
 * ModAppSession Modifies an existing Individual Application Session Context
 * @param ctx context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the resource
 * @param ascUpdateData modification of the resource.
 * @return AppSessionContext on 200 resp otherwise ProbleDetails
 */
func (c *PolicyAuthAPIClient) ModAppSession(
	ctx context.Context, appSessionID string,
	ascUpdateData AppSessionContextUpdateData) (PcfPAResponse, error) {

	var (
		httpMethod = http.MethodPatch
		patchBody  interface{}
		retVal     PcfPAResponse
	)

	// create path and map variables
	path := c.rootURI + "/app-sessions/" + appSessionID

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = "application/json"
	headerParams["Accept"] = "application/json"

	// body params
	patchBody = &ascUpdateData
	r, err := c.prepareRequest(ctx, path, httpMethod, patchBody,
		headerParams)
	if err != nil {
		return retVal, err
	}

	httpResponse, err := c.callAPI(r)
	retVal.httpResp = httpResponse
	if err != nil || httpResponse == nil {
		return retVal, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, err
	}

	switch httpResponse.StatusCode {
	case 200:
		var asc *AppSessionContext = new(AppSessionContext)
		err = json.Unmarshal(respBody, &asc)
		if err != nil {
			log.Errf("Error in unmarshalling json, " +
				"ModAppSession: " + err.Error())
			httpResponse.StatusCode = 500
		}
		retVal.appSessCtx = asc
		return retVal, err
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			validatePAAuthToken(c)
		}
		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"ModAppSession: " + err.Error())
			return retVal, err
		}
		retVal.probDetails = v
		return retVal, err
	}

	err = errors.New(string(respBody))
	log.Errf("ModAppSess- PCF Returned Error: " +
		string(respBody))
	return retVal, err
}

func handleUpdateEventResp(respBody []byte, httpResp *http.Response,
	c *PolicyAuthAPIClient) (
	retVal EventSubscResponse, err error) {

	var (
		eventSubscResp EventsSubscReqData
		evsNotifResp   EventsNotification
	)

	retVal.httpResp = httpResp
	switch httpResp.StatusCode {
	case 200, 201:
		retVal.locationURI = getLocationURI(httpResp, c)

		err = json.Unmarshal(respBody, &eventSubscResp)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"UpdateEventSubsc: " + err.Error())
			httpResp.StatusCode = 500
			return retVal, err
		}
		retVal.eventSubscReq = &eventSubscResp

		err = json.Unmarshal(respBody, &evsNotifResp)
		if err == nil {
			retVal.evsNotif = &evsNotifResp
		}
		return retVal, nil

	case 204:
		return retVal, nil

	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:
		if httpResp.StatusCode == 401 {
			validatePAAuthToken(c)
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
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the Events Subscription resource
 * @param eventsSubscReqData Creation or modification of an Events Subscription
 * resource.
@return AnyOfEventsSubscReqDataEventsNotification
*/
func (c *PolicyAuthAPIClient) UpdateEventsSubsc(ctx context.Context,
	appSessionID string, eventSubscReq *EventsSubscReqData) (
	EventSubscResponse, error) {

	var (
		httpMethod = http.MethodPut
		reqBody    interface{}
		retVal     EventSubscResponse
	)

	// create path and map variables
	path := c.rootURI + "/app-sessions/" + appSessionID +
		"/events-subscription"

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = contentTypeJSON
	headerParams["Accept"] = contentTypeJSON

	// body params
	if eventSubscReq != nil {
		reqBody = eventSubscReq
	}

	r, err := c.prepareRequest(ctx, path, httpMethod, reqBody,
		headerParams)
	if err != nil {
		return retVal, err
	}

	httpResponse, err := c.callAPI(r)
	retVal.httpResp = httpResponse
	if err != nil || httpResponse == nil {
		return retVal, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, err
	}

	retVal, err = handleUpdateEventResp(respBody, httpResponse,
		c)
	return retVal, err

}

/*
DeleteEventsSubsc deletes the Events Subscription subresource
 * @param ctx context.Context - for authentication, logging, cancellation,
 *   deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param appSessionID string identifying the Individual Application Session
 * Context resource
*/
func (c *PolicyAuthAPIClient) DeleteEventsSubsc(ctx context.Context,
	appSessionID string) (EventSubscResponse, error) {
	var (
		httpMethod = http.MethodDelete
		reqBody    interface{}
		retVal     EventSubscResponse
	)

	// create path and map variables
	path := c.rootURI + "/app-sessions/" + appSessionID +
		"/events-subscription"

	headerParams := make(map[string]string)
	headerParams["Accept"] = contentTypeJSON

	r, err := c.prepareRequest(ctx, path, httpMethod, reqBody,
		headerParams)
	if err != nil {
		return retVal, err
	}

	httpResponse, err := c.callAPI(r)
	retVal.httpResp = httpResponse
	if err != nil || httpResponse == nil {
		return retVal, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	defer func() {
		err = httpResponse.Body.Close()
		if err != nil {
			log.Errf("Resp Body wasn't closed properly" +
				err.Error())
		}
	}()
	if err != nil {
		return retVal, err
	}

	switch httpResponse.StatusCode {
	case 204:
		return retVal, nil

	case 400, 401, 403, 404, 429, 500, 503:
		if httpResponse.StatusCode == 401 {
			validatePAAuthToken(c)
		}
		var v *ProblemDetails = new(ProblemDetails)
		err = json.Unmarshal(respBody, v)
		if err != nil {
			log.Errf("Error in unmarshalling response body, " +
				"DeleteEventSubsc: " + err.Error())
			return retVal, err
		}
		retVal.probDetails = v
		return retVal, err
	}

	err = errors.New(string(respBody))
	if err == nil {
		err = errors.New("Empty response from PCF")
	}

	log.Errf("DeleteEventSubs- PCF Returned Error: " +
		string(respBody))
	return retVal, err
}
