/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import "context"

/* The SB interfaces towards the UDR that need to be implemented by
   eith the NEF SB stub / NEF SB client receivers */

// UdrInfluenceResponse contains the response information from UDR
type UdrInfluenceResponse struct {
	// responseCode contains the http response code provided by the UDR
	ResponseCode uint16
	// tid if not nil contains the TrafficInfluData data provided by UDR
	Tid *TrafficInfluData
	// pd if not not contains the problem infomration from UDR.
	// Valid for 3xx, 4xx, 5xx or 6xx responses
	Pd *ProblemDetails
}

// UdrPfdResponse contains the response information from UDR
type UdrPfdResponse struct {
	// responseCode contains the http response code provided by the UDR
	ResponseCode uint16
	// AppPfd if not nil contains the PfdDataForApp data provided by UDR
	AppPfd *PfdDataForApp
	// pd if not not contains the problem infomration from UDR.
	// Valid for 3xx, 4xx, 5xx or 6xx responses
	Pd *ProblemDetails
}

// InfluenceID contains the application session id returned by the UDR
// Its present in the location header as below:
// "{apiRoot}/nudr-dr/v1/aapplication-data/influenceData/{influenceId}"
type InfluenceID string

// UdrAppID is the application identifier of the request PFD(s)
// "{apiRoot}/nudr-dr/v1/application-data/pfds/{appId}"
type UdrAppID string

// UdrData defines the interfaces that are exposed for
// TrafficInfluence,
type UdrData interface {
	// UdrInfluenceDataCreate sends PUT request to the UDR using the
	// configuration mentioned in the context. Context would have all the
	// inforrmation related to the UDR like the URI, authentication, logging,
	// cancellation It returns the response received from the UDR, the app and
	// any error encountered when sending the request. The contents of the
	// actual response received are part of the UDRPolicyResponse
	UdrInfluenceDataCreate(ctx context.Context,
		body TrafficInfluData,
		iid InfluenceID) (UdrInfluenceResponse, error)

	// UdrInfluenceDataUpdate sends PATCH request to the UDR to the
	// appSessionId using the configuration mentioned in the context. Context
	// would have all the infomration related to the UDR like the URI,
	// authentication, logging, cancellation
	// It returns the response received from the UDR and any error encountered
	// when sending the request. The contents of the actual response received
	// are part of the UdrPolicyResponse
	UdrInfluenceDataUpdate(ctx context.Context,
		body TrafficInfluDataPatch,
		iid InfluenceID) (UdrInfluenceResponse, error)

	// UdrInfluenceDataDelete sends, DELETE request to the UDR to the
	// appSessionId using the configuration mentioned in the context. Context
	// would have all the infomration related to the UDR like the URI,
	// authentication, logging, cancellationIt returns the response received
	// from the UDR and any error encountered when sending the
	// request. The contents of the actual response received are part of the
	// UdrPolicyResponse
	UdrInfluenceDataDelete(ctx context.Context, iid InfluenceID) (
		UdrInfluenceResponse, error)

	// UdrInfluenceDataGet sends GET request to the UDR
	// using the configuration mentioned in the context. Context would have all
	// the infomration related to the UDR like the URI, authentication, logging,
	// cancellation. It returns the response received from the UDR and any error
	// encountered when sending therequest. The contents of the actual response
	// received are part of the UdrPolicyResponse
	UdrInfluenceDataGet(ctx context.Context) (
		UdrInfluenceResponse, error)

	// UdrPfdDataCreate sends PUT request to the UDR using the
	// configuration mentioned in the context. Context would have all the
	// inforrmation related to the UDR like the URI, authentication, logging,
	// cancellation It returns the response received from the UDR, the app and
	// any error encountered when sending the request. The contents of the
	// actual response received are part of the UDRPfdResponse
	UdrPfdDataCreate(ctx context.Context, body PfdDataForApp) (
		UdrPfdResponse, error)

	// UdrPfdDataGet sends GET request to the UDR
	// using the configuration mentioned in the context. Context would have all
	// the infomration related to the UDR like the URI, authentication, logging,
	// cancellation. It returns the response received from the UDR and any error
	// encountered when sending therequest. The contents of the actual response
	// received are part of the UDRPfdResponse
	UdrPfdDataGet(ctx context.Context, appID UdrAppID) (
		UdrPfdResponse, error)

	// UdrPfdDataDelete sends, DELETE request to the UDR to the
	// appId using the configuration mentioned in the context. Context
	// would have all the infomration related to the UDR like the URI,
	// authentication, logging, cancellationIt returns the response received
	// from the UDR and any error encountered when sending the
	// request. The contents of the actual response received are part of the
	// UDRPfdResponse
	UdrPfdDataDelete(ctx context.Context, appID UdrAppID) (
		UdrPfdResponse, error)
}
