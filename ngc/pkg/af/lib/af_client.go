// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package af

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	jsonCheck = regexp.MustCompile("(?i:[application|text]/json)")
	xmlCheck  = regexp.MustCompile("(?i:[application|text]/xml)")
)

// Client manages communication with the NEF Northbound API API v1.0.1
// In most cases there should be only one, shared, Client.
type Client struct {
	cfg *Configuration
	// Reuse a single struct instead of allocating one for each service on
	// the heap.
	common service

	// API Services
	TrafficInfluSubGetAllAPI *TrafficInfluenceSubscriptionGetAllAPIService
	TrafficInfluSubDeleteAPI *TrafficInfluenceSubscriptionDeleteAPIService
	TrafficInfluSubGetAPI    *TrafficInfluenceSubscriptionGetAPIService
	TrafficInfluSubPatchAPI  *TrafficInfluenceSubscriptionPatchAPIService
	TrafficInfluSubPostAPI   *TrafficInfluenceSubscriptionPostAPIService
	TrafficInfluSubPutAPI    *TrafficInfluenceSubscriptionPutAPIService
}

type service struct {
	client *Client
}

// NewClient creates a new API client. Requires a userAgent string describing
// your application.
// optionally a custom http.Client to allow for advanced features
// such as caching.
func NewClient(cfg *Configuration) *Client {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
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

	return c
}

// selectHeaderContentType select a content type from the available list.
func selectHeaderContentType(contentTypes []string) string {
	if len(contentTypes) == 0 {
		return ""
	}
	if contains(contentTypes, "application/json") {
		return "application/json"
	}
	// use the first content type specified in 'consumes'
	return contentTypes[0]
}

// selectHeaderAccept join all accept types and return
func selectHeaderAccept(accepts []string) string {
	if len(accepts) == 0 {
		return ""
	}

	if contains(accepts, "application/json") {
		return "application/json"
	}

	return strings.Join(accepts, ",")
}

// contains is a case insenstive match, finding needle in a haystack
func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.EqualFold(strings.ToLower(a), strings.ToLower(needle)) {
			//strings.ToLower(a) == strings.ToLower(needle) {
			return true
		}
	}
	return false
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
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}

		body, err = setBody(postBody, contentType)
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
	postBody interface{},
	headerParams map[string]string,
) (localVarRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect postBody type and post.
	body, err = genBody(postBody, headerParams)
	if err != nil {
		return nil, err
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Generate a new request
	localVarRequest, err = genNewRequest(body, url.String(), method)
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

// detectContentType method is used to figure out `Request.Body` content
// type for request header
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

// Ripped from https://github.com/gregjones/httpcache/blob/master/httpcache.go
type cacheControl map[string]string

func parseCacheControl(headers http.Header) cacheControl {
	cc := cacheControl{}
	ccHeader := headers.Get("Cache-Control")
	for _, part := range strings.Split(ccHeader, ",") {
		part = strings.Trim(part, " ")
		if part == "" {
			continue
		}
		if strings.ContainsRune(part, '=') {
			keyval := strings.Split(part, "=")
			cc[strings.Trim(keyval[0], " ")] =
				strings.Trim(keyval[1], ",")
		} else {
			cc[part] = ""
		}
	}
	return cc
}

// CacheExpires helper function to determine remaining time before repeating a
// request.
func CacheExpires(r *http.Response) time.Time {
	// Figure out when the cache expires.
	var (
		expires  time.Time
		lifetime time.Duration
	)
	now, err := time.Parse(time.RFC1123, r.Header.Get("date"))
	if err != nil {
		return time.Now()
	}
	respCacheControl := parseCacheControl(r.Header)

	if maxAge, ok := respCacheControl["max-age"]; ok {
		lifetime, err = time.ParseDuration(maxAge + "s")
		if err != nil {
			log.Errf("error parsing time duration")
			expires = now
		} else {
			expires = now.Add(lifetime)
		}
	} else {
		expiresHeader := r.Header.Get("Expires")
		if expiresHeader != "" {
			expires, err = time.Parse(time.RFC1123, expiresHeader)
			if err != nil {
				expires = now
			}
		}
	}
	return expires
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

// Body returns the raw bytes of the response
func (e GenericError) Body() []byte {
	return e.body
}

// Model returns the unpacked model of the error
func (e GenericError) Model() interface{} {
	return e.model
}
