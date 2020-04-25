// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	oauth2 "github.com/otcshare/epcforedge/ngc/pkg/oauth2"
	"golang.org/x/net/http2"
)

var contentTypeJSON string = "application/json"

// PolicyAuthAPIClient type
// In most cases there should be only one, shared, PolicyAuthAPIClient.
type PolicyAuthAPIClient struct {
	cfg *CliPcfConfig
	/*
	 * Reuse a single struct instead of allocating one for each
	 * policyAuthService on the heap.
	 */
	common policyAuthService

	oAuth2Token string

	// API Services
	PolicyAuthAppSessionAPI *PolicyAuthAppSessionAPIService

	PolicyAuthEventSubsAPI *PolicyAuthEventSubsAPIService

	PolicyAuthIndividualAppSessAPI *PolicyAuthIndividualAppSessAPIService
}

type policyAuthService struct {
	client *PolicyAuthAPIClient
}

func getPcfOAuth2Token() (token string, err error) {

	token, err = oauth2.GetAccessToken()
	if err == nil {
		log.Infoln("Got Pcf OAuth2 Access Token: " + token)
	}
	return token, err
}

// NewPolicyAuthAPIClient - helper func
/*
 * NewAPIClient creates a new API client. Basically create new http client if
 * not set in client configurations.
 */
func NewPolicyAuthAPIClient(cfg *CliPcfConfig) (*PolicyAuthAPIClient, error) {
	if cfg.HTTPClient == nil {
		CACert, err := ioutil.ReadFile(cfg.CliCertPath)
		if err != nil {
			log.Errf("Error: %v", err)
			return nil, err
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

	c := &PolicyAuthAPIClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.PolicyAuthAppSessionAPI = (*PolicyAuthAppSessionAPIService)(&c.common)
	c.PolicyAuthEventSubsAPI = (*PolicyAuthEventSubsAPIService)(&c.common)
	c.PolicyAuthIndividualAppSessAPI =
		(*PolicyAuthIndividualAppSessAPIService)(&c.common)
	if cfg.OAuth2Support {
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

// callAPI do the request.
func (c *PolicyAuthAPIClient) callAPI(request *http.Request) (
	*http.Response, error) {

	resp, err := c.cfg.HTTPClient.Do(request)
	if err != nil {
		return resp, err
	}

	return resp, err
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
		if reqContentType == "" {
			reqContentType = detectContentType(reqBody)
			headerParams["Content-Type"] = reqContentType
		}

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
	httpRequest.Header.Add("User-Agent", c.cfg.UserAgent)

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

/*
 * detectContentType method is used to figure out `Request.Body` content type
 * for request header.
 */
func detectContentType(body interface{}) string {
	reqContentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		reqContentType = "application/json; charset=utf-8"
	case reflect.String:
		reqContentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			reqContentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			reqContentType = "application/json; charset=utf-8"
		}
	}

	return reqContentType
}
