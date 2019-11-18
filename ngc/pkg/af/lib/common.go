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
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/url"
	"strings"
)

func genAFTransID(trans TransactionIDs) int {
	var (
		num   int
		min   = 1
		found = true
	)
	for max := range trans {
		num =
			max
		break
	}
	for max := range trans {
		if max > num {
			num = max
		}
	}

	if num == math.MaxInt32 {
		num = min
	}
	//look for a free ID until it is <= math.MaxInt32 is achieved again
	for found && num < math.MaxInt32 {
		num++
		//check if the ID is in use, if not - return the ID
		if _, found = trans[num]; !found {
			trans[num] = TrafficInfluSub{}
			return num
		}
	}
	return 0
}

func getSubsIDFromURL(url *url.URL) (string, error) {

	subsURL := url.String()
	if url == nil {
		return "", errors.New("empty URL in the request message")
	}
	s := strings.Split(subsURL, "/")
	return s[len(s)-1], nil
}

func genTransactionID(afCtx *afContext) (int, error) {

	afTransID := genAFTransID(afCtx.transactions)
	if afTransID == 0 {
		return 0, errors.New("the pool of AF Transaction IDs is already used")
	}

	return afTransID, nil

}

func handleGetErrorResp(localVarHTTPResponse *http.Response,
	localVarBody []byte) error {

	newErr := GenericError{
		body:  localVarBody,
		error: localVarHTTPResponse.Status,
	}
	switch localVarHTTPResponse.StatusCode {
	case 400, 401, 403, 404, 406, 429, 500, 503:

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

func handlePostPutPatchErrorResp(localVarHTTPResponse *http.Response,
	localVarBody []byte) error {

	newErr := GenericError{
		body:  localVarBody,
		error: localVarHTTPResponse.Status,
	}
	switch localVarHTTPResponse.StatusCode {
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:

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
