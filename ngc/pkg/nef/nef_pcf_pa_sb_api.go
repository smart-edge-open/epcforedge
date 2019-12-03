// Copyright 2019 Intel Corporation. All rights reserved
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

package ngcnef

import "context"

/* The SB interfaces towards the PCF that need to be implemented by
   eith the NEF SB stub / NEF SB client receivers */

// PcfPolicyResponse contains the response from PCF
type PcfPolicyResponse struct {
	// responseCode contains the http response code provided by the PCF
	ResponseCode uint16
	// asc if not nil contains the AppSessionContext data provided by PCF
	Asc *AppSessionContext
	// pd if not not contains the problem infomration from PCF.
	// Valid for 3xx, 4xx, 5xx or 6xx responses
	Pd *ProblemDetails
}

// AppSessionID contains the application session id returned by the PCF
// Its present in the location header as below:
// "{apiRoot}/npcf-policyauthorization/v1/app-sessions/{appSessionId}"
type AppSessionID string

// PcfPolicyAuthorization defines the interfaces that are exposed for
// TrafficInfluence
type PcfPolicyAuthorization interface {
	// PolicyAuthorizationCreate sends POST request to the PCF using the
	// configuration mentioned in the context. Context would have all the
	// informration related to the PCF like the URI, authentication, logging,
	// cancellation It returns the response received from the PCF, the app and
	// any error encountered when sending the request. The contents of the
	// actual response received are part of the PcfPolicyResponse
	PolicyAuthorizationCreate(ctx context.Context,
		body AppSessionContext) (AppSessionID, PcfPolicyResponse, error)

	// PolicyAuthorizationUpdate sends PATCH request to the PCF to the
	// appSessionId using the configuration mentioned in the context. Context
	// would have all the infomration related to the PCF like the URI,
	// authentication, logging, cancellation
	// It returns the response received from the PCF and any error encountered
	// when sending the request. The contents of the actual response received
	// are part of the PcfPolicyResponse
	PolicyAuthorizationUpdate(ctx context.Context,
		body AppSessionContextUpdateData,
		appSessionID AppSessionID) (PcfPolicyResponse, error)

	// PolicyAuthorizationDelete sends, DELETE request to the PCF to the
	// appSessionId using the configuration mentioned in the context. Context
	// would have all the infomration related to the PCF like the URI,
	// authentication, logging, cancellationIt returns the response received
	// from the PCF and any error encountered when sending the
	// request. The contents of the actual response received are part of the
	// PcfPolicyResponse
	PolicyAuthorizationDelete(ctx context.Context, appSessionID AppSessionID) (
		PcfPolicyResponse, error)

	// PolicyAuthorizationGet sends GET request to the PCF to the appSessionId
	// using the configuration mentioned in the context. Context would have all
	// the infomration related to the PCF like the URI, authentication, logging,
	// cancellation It returns the response received from the PCF and any error
	// encountered when sending therequest. The contents of the actual response
	//received are part of the PcfPolicyResponse
	PolicyAuthorizationGet(ctx context.Context, appSessionID AppSessionID) (
		PcfPolicyResponse, error)
}
