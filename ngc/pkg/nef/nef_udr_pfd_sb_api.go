/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import "context"

/* The SB interfaces towards the UDR that need to be implemented by
   eith the NEF SB stub / NEF SB client receivers */

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

// UdrAppID is the application identifier of the request PFD(s)
// "{apiRoot}/nudr-dr/v1/application-data/pfds/{appId}"
type UdrAppID string

// UdrPfdData defines the interfaces that are exposed for
// TrafficInfluence,
type UdrPfdData interface {

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
