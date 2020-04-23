// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

// In most cases there should be only one, shared, PolicyAuthApiClient.
type PolicyAuthApiClient struct {
	cfg *CliPcfConfig
	/*
	 * Reuse a single struct instead of allocating one for each
	 * policyAuthService on the heap.
	 */
	common policyAuthService

	// API Services

	PolicyAuthAppSessionApi *PolicyAuthAppSessionApiService

	PolicyAuthEventSubsApi *PolicyAuthEventSubsApiService

	PolicyAuthIndividualAppSessApi *PolicyAuthIndividualAppSessApiService
}

type policyAuthService struct {
	client *PolicyAuthApiClient
}

/*
 * NewAPIClient creates a new API client. Basically create new http client if
 * not set in client configurations.
 */
func NewPolicyAuthApiClient(cfg *CliPcfConfig) *PolicyAuthApiClient {
	if cfg.HTTPClient == nil {
		CACert, err := ioutil.ReadFile(cfg.CliCertPath)
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

	c := &PolicyAuthApiClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.PolicyAuthAppSessionApi = (*PolicyAuthAppSessionApiService)(&c.common)
	c.PolicyAuthEventSubsApi = (*PolicyAuthEventSubsApiService)(&c.common)
	c.PolicyAuthIndividualAppSessApi =
		(*PolicyAuthIndividualAppSessApiService)(&c.common)

	return c
}

// contains is a case insenstive match, finding needle in a haystack
func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.ToLower(a) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}

// callAPI do the request.
func (c *PolicyAuthApiClient) callAPI(request *http.Request) (
	*http.Response, error) {

	if c.cfg.Debug {
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		//log.Printf("\n%s\n", string(dump))
		_ = dump
	}

	resp, err := c.cfg.HTTPClient.Do(request)
	if err != nil {
		return resp, err
	}

	if c.cfg.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, err
		}
		//log.Printf("\n%s\n", string(dump))
		_ = dump
	}

	return resp, err
}

// prepareRequest build the request
func (c *PolicyAuthApiClient) prepareRequest(
	ctx context.Context,
	path string, method string,
	reqBody interface{},
	headerParams map[string]string,
) (httpRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect reqBody type and post.
	if reqBody != nil {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(reqBody)
			headerParams["Content-Type"] = contentType
		}

		body, err = setBody(reqBody, contentType)
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
	httpRequest.Header.Add("User-Agent", c.cfg.UserAgent)

	if ctx != nil {
		// add context to the request
		httpRequest = httpRequest.WithContext(ctx)

		// Walk through any authentication.

	}
	return httpRequest, nil
}

/*
 * detectContentType method is used to figure out `Request.Body` content type
 * for request header.
 */
func detectContentType(body interface{}) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}
