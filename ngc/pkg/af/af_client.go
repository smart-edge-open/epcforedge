// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"golang.org/x/net/http2"
)

const contentType string = "application/json"

var (
	jsonCheck = regexp.MustCompile("(?i:[application|text]/json)")
	xmlCheck  = regexp.MustCompile("(?i:[application|text]/xml)")
)

// Client manages communication with the NEF Northbound API API v1.0.1
type Client struct {
	cfg *CliConfig
	// Reuse a single struct instead of allocating one for each service
	// on the heap.
	common service
	// API Services
	TrafficInfluSubGetAllAPI  *TrafficInfluenceSubscriptionGetAllAPIService
	TrafficInfluSubDeleteAPI  *TrafficInfluenceSubscriptionDeleteAPIService
	TrafficInfluSubGetAPI     *TrafficInfluenceSubscriptionGetAPIService
	TrafficInfluSubPatchAPI   *TrafficInfluenceSubscriptionPatchAPIService
	TrafficInfluSubPostAPI    *TrafficInfluenceSubscriptionPostAPIService
	TrafficInfluSubPutAPI     *TrafficInfluenceSubscriptionPutAPIService
	PfdManagementGetAllAPI    *PfdManagementTransactionGetAllAPIService
	PfdManagementPostAPI      *PfdManagementTransactionPostAPIService
	PfdManagementGetAPI       *PfdManagementTransactionGetAPIService
	PfdManagementDeleteAPI    *PfdManagementTransactionDeleteAPIService
	PfdManagementPutAPI       *PfdManagementTransactionPutAPIService
	PfdManagementAppGetAPI    *PfdManagementTransactionAppGetAPIService
	PfdManagementAppDeleteAPI *PfdManagementTransactionAppDeleteAPIService
	PfdManagementAppPutAPI    *PfdManagementTransactionAppPutAPIService
	PfdManagementAppPatchAPI  *PfdManagementTransactionAppPatchAPIService
}

type service struct {
	client *Client
}

// TestAf boolean to be set to true in AF UT
var TestAf bool = false

// HTTPClient to be setup in AF UT
var HTTPClient *http.Client

// SetHTTPClient Function to setup a httpClient for testing if required
func SetHTTPClient(httpClient *http.Client) {

	HTTPClient = httpClient

}

// NewClient creates a new API client.
func NewClient(cfg *CliConfig) *Client {

	if TestAf == true {

		cfg.HTTPClient = HTTPClient

	}
	if cfg.HTTPClient == nil || TestAf == false {

		CACert, err := ioutil.ReadFile(cfg.NEFCliCertPath)
		if err != nil {
			log.Errf("Error: %v", err)
		}

		CACertPool := x509.NewCertPool()
		CACertPool.AppendCertsFromPEM(CACert)

		cfg.HTTPClient = &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: CACertPool,
				},
			},
		}
	}

	c := &Client{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.TrafficInfluSubGetAllAPI =
		(*TrafficInfluenceSubscriptionGetAllAPIService)(&c.common)
	c.TrafficInfluSubDeleteAPI =
		(*TrafficInfluenceSubscriptionDeleteAPIService)(&c.common)
	c.TrafficInfluSubGetAPI =
		(*TrafficInfluenceSubscriptionGetAPIService)(&c.common)
	c.TrafficInfluSubPatchAPI =
		(*TrafficInfluenceSubscriptionPatchAPIService)(&c.common)
	c.TrafficInfluSubPostAPI =
		(*TrafficInfluenceSubscriptionPostAPIService)(&c.common)
	c.TrafficInfluSubPutAPI =
		(*TrafficInfluenceSubscriptionPutAPIService)(&c.common)
	c.PfdManagementGetAllAPI =
		(*PfdManagementTransactionGetAllAPIService)(&c.common)
	c.PfdManagementPostAPI =
		(*PfdManagementTransactionPostAPIService)(&c.common)
	c.PfdManagementGetAPI =
		(*PfdManagementTransactionGetAPIService)(&c.common)
	c.PfdManagementDeleteAPI =
		(*PfdManagementTransactionDeleteAPIService)(&c.common)
	c.PfdManagementPutAPI =
		(*PfdManagementTransactionPutAPIService)(&c.common)
	c.PfdManagementAppGetAPI =
		(*PfdManagementTransactionAppGetAPIService)(&c.common)
	c.PfdManagementAppDeleteAPI =
		(*PfdManagementTransactionAppDeleteAPIService)(&c.common)
	c.PfdManagementAppPutAPI =
		(*PfdManagementTransactionAppPutAPIService)(&c.common)
	c.PfdManagementAppPatchAPI =
		(*PfdManagementTransactionAppPatchAPIService)(&c.common)

	return c
}

// callAPI do the request.
func (c *Client) callAPI(request *http.Request) (*http.Response, error) {
	return c.cfg.HTTPClient.Do(request)
}

func genNewRequest(body io.Reader, url string,
	method string) (*http.Request, error) {

	var (
		localVarRequest *http.Request
		err             error
	)

	localVarRequest, err = http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	return localVarRequest, nil
}

func genBody(postBody interface{},
	headerParams map[string]string) (*bytes.Buffer, error) {

	var (
		body = &bytes.Buffer{}
		err  error
	)

	if postBody != nil {
		contenttype := headerParams["Content-Type"]
		body, err = setBody(postBody, contenttype)
		if err != nil {
			return nil, err
		}

	}
	return body, nil
}

// prepareRequest build the request
func (c *Client) prepareRequest(
	ctx context.Context,
	path string, method string,
	body interface{},
	headerParams map[string]string,
) (localVarRequest *http.Request, err error) {

	var b *bytes.Buffer

	// Detect body type and post.
	b, err = genBody(body, headerParams)
	if err != nil {
		return nil, err
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Generate a new request
	localVarRequest, err = genNewRequest(b, url.String(), method)
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		localVarRequest.Header = headers
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.cfg.UserAgent)

	if ctx != nil {
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)

		// Walk through any authentication here.

	}

	if c.cfg.OAuth2Support {
		auth, err := getNEFAuthorizationToken()
		if err != nil {
			return nil, err
		}
		// Add the Authorization header to the request.
		localVarRequest.Header.Add("Authorization", "Bearer "+auth)
	}

	return localVarRequest, nil
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer,
	err error) {

	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if jsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	} else if xmlCheck.MatchString(contentType) {
		err = xml.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("invalid body type %s", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

// GenericError Provides access to the body,
// error and model on returned errors.
type GenericError struct {
	body  []byte
	error string
	model interface{}
}

// Error returns non-empty string if there was an error.
func (e GenericError) Error() string {
	return e.error
}
