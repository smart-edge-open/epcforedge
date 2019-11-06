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
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Linger please
var (
	_ context.Context
)

type TrafficInfluenceSubscriptionPostApiService service

/*
TrafficInfluenceSubscriptionPostApiService Creates a new subscription resource
Creates a new subscription resource
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afId Identifier of the AF
 * @param body Request to create a new subscription resource

@return TrafficInfluSub
*/
func (a *TrafficInfluenceSubscriptionPostApiService) SubscriptionPost(
	ctx context.Context, afId string,
	body TrafficInfluSub) (TrafficInfluSub, *http.Response, error) {

	var (
		localVarHttpMethod  = strings.ToUpper("Post")
		localVarPostBody    interface{}
		localVarFileName    string
		localVarFileBytes   []byte
		localVarReturnValue TrafficInfluSub
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/{afId}/subscriptions"
	localVarPath = strings.Replace(localVarPath,
		"{"+"afId"+"}", fmt.Sprintf("%v", afId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType :=
		selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept :=
		selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	// body params
	localVarPostBody = &body
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod,
		localVarPostBody, localVarHeaderParams, localVarQueryParams,
		localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {

		// If we succeed, return the data, otherwise pass on to
		// decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody,
			localVarHttpResponse.Header.Get("Content-Type"))
		if err == nil {
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericError{
			body:  localVarBody,
			error: localVarHttpResponse.Status,
		}

		if localVarHttpResponse.StatusCode == 201 {
			var v TrafficInfluSub
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 400 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 401 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 403 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 404 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 411 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHttpResponse,
					newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 413 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 415 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 429 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 500 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 503 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		if localVarHttpResponse.StatusCode == 0 {
			var v interface{}
			err = a.client.decode(&v, localVarBody,
				localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue,
					localVarHttpResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHttpResponse, newErr
		}

		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}
