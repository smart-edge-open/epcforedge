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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Linger please
var (
	_ context.Context
)

// TrafficInfluenceSubscriptionDeleteAPIService type
type TrafficInfluenceSubscriptionDeleteAPIService service

func (a *TrafficInfluenceSubscriptionDeleteAPIService) handleDeleteResponse(
	localVarHTTPResponse *http.Response,
	localVarBody []byte) error {

	if localVarHTTPResponse.StatusCode > 300 {

		newErr := GenericError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}

		switch localVarHTTPResponse.StatusCode {

		case 400, 401, 403, 404, 429, 500, 503:

			var v ProblemDetails
			err := json.Unmarshal(localVarBody, &v)
			if err != nil {
				newErr.error = err.Error()
				return newErr
			}
			newErr.model = v
			return newErr

		default:
			var v interface{}
			err := json.Unmarshal(localVarBody, &v)
			if err != nil {
				newErr.error = err.Error()
				return newErr
			}
			newErr.model = v
			return newErr
		}
	}

	return nil
}

/*
SubscriptionDelete Deletes an already
existing subscription
Deletes an already existing subscription
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF
 * @param subscriptionID Identifier of the subscription resource


*/
func (a *TrafficInfluenceSubscriptionDeleteAPIService) SubscriptionDelete(
	ctx context.Context, afID string, subscriptionID string) (*http.Response,
	error) {
	var (
		localVarHTTPMethod = strings.ToUpper("Delete")
		localVarPostBody   interface{}
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath +
		"/{afId}/subscriptions/{subscriptionId}"
	localVarPath = strings.Replace(localVarPath,
		"{"+"afId"+"}", fmt.Sprintf("%v", afID), -1)
	localVarPath = strings.Replace(localVarPath,
		"{"+"subscriptionId"+"}", fmt.Sprintf("%v", subscriptionID), -1)

	localVarHeaderParams := make(map[string]string)

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod,
		localVarPostBody, localVarHeaderParams)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(r)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	if err != nil {
		if err = localVarHTTPResponse.Body.Close(); err != nil {
			log.Errf("response body could not be closed properly")
		}
		return localVarHTTPResponse, err
	}
	if err = localVarHTTPResponse.Body.Close(); err != nil {
		log.Errf("response body could not be closed properly")
	}

	if err = a.handleDeleteResponse(localVarHTTPResponse,
		localVarBody); err != nil {

		return localVarHTTPResponse, err
	}

	return localVarHTTPResponse, nil
}
