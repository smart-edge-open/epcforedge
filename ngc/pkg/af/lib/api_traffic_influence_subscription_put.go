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

//TrafficInfluenceSubscriptionPutAPIService type
type TrafficInfluenceSubscriptionPutAPIService service

func (a *TrafficInfluenceSubscriptionPutAPIService) handlePutResponse(
	localVarReturnValue *TrafficInfluSub, localVarHTTPResponse *http.Response,
	localVarBody []byte) error {

	if localVarHTTPResponse.StatusCode == 200 {

		err := json.Unmarshal(localVarBody, localVarReturnValue)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
		}
		return err
	}

	return handlePostPutPatchErrorResp(localVarHTTPResponse, localVarBody)

}

/*
SubscriptionPut Replaces an existing
subscription resource
Replaces an existing subscription resource
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF
 * @param subscriptionID Identifier of the subscription resource
 * @param body Parameters to replace the existing subscription
@return TrafficInfluSub, *http.Response, error
*/
func (a *TrafficInfluenceSubscriptionPutAPIService) SubscriptionPut(
	ctx context.Context, afID string, subscriptionID string,
	body TrafficInfluSub) (TrafficInfluSub, *http.Response, error) {

	var (
		localVarHTTPMethod  = strings.ToUpper("Put")
		localVarPostBody    interface{}
		localVarReturnValue TrafficInfluSub
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
	localVarHTTPContentType :=
		selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept :=
		selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = &body
	r, err := a.client.prepareRequest(ctx, localVarPath,
		localVarHTTPMethod, localVarPostBody, localVarHeaderParams)

	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(r)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	defer func() {
		err = localVarHTTPResponse.Body.Close()
		if err != nil {
			log.Errf("response body was not closed properly")
		}
	}()

	if err != nil {
		log.Errf("http response body could not be read")
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if err = a.handlePutResponse(&localVarReturnValue, localVarHTTPResponse,
		localVarBody); err != nil {
		log.Errf("Handle Put response")
		return localVarReturnValue, localVarHTTPResponse, err
	}
	return localVarReturnValue, localVarHTTPResponse, nil
}
