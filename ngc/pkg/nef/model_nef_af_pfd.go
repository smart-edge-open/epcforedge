/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

// Pfd is the structure of Packet Flow Description for an external Application
// Identifier
type Pfd struct {
	// Identifies a PDF of an application identifier.
	PfdID string `json:"pfdID"`
	// Represents a 3-tuple with protocol, server ip and server port for UL/DL
	// application traffic. The content of the string has the same encoding as
	// the IPFilterRule AVP value as defined in IETFÂ RFCÂ 6733.
	FlowDescriptions []string `json:"flowDescriptions,omitempty"`
	// Indicates a URL or a regular expression which is used to match the
	// significant parts of the URL.
	Urls []string `json:"urls,omitempty"`
	// Indicates an FQDN or a regular expression as a domain name matching
	// criteria.
	DomainNames []string `json:"domainNames,omitempty"`
}

// PfdData is the type that represents a PFD request to add, update or remove
// PFD(s) for one external application identifier provided by AF
type PfdData struct {
	// Each element uniquely identifies external application identifier
	ExternalAppID string `json:"externalAppID"`
	// Link to the resource. This parameter shall be supplied by the NEF in
	// HTTP responses that include an object of PfdData type
	Self Link `json:"self,omitempty"`
	// Contains the PFDs of the external application identifier. Each PFD is
	// identified in the map via a key containing the PFD identifier.
	Pfds map[string]Pfd `json:"pfds"`
	// Indicates that the list of PFDs in this request should be deployed
	// within the time interval indicated by the Allowed Delay
	AllowedDelay DurationSecRm `json:"allowedDelay,omitempty"`
	// SCEF supplied property, inclusion of this property means the allowed
	// delayed cannot be satisfied, i.e. it is smaller than the caching time,
	// but the PFD data is still stored.
	CachingTime DurationSecRo `json:"cachingTime,omitempty"`
}

// PfdManagement resource for a PFD management request
type PfdManagement struct {
	// Link to the resource "Individual PFD Management Transaction".
	// This parameter shall be supplied by the NEF in HTTP responses.
	Self Link `json:"self,omitempty"`
	// String identifying supported optional features of PFD Management
	// This attribute shall be provided in the POST request and in the
	// response of successful resource creation.
	SuppFeat SupportedFeatures `json:"suppFeat,omitempty"`
	// Each element uniquely identifies the PFDs for an external application
	// identifier. Each element is identified in the map via an external
	// application identifier as key. The response shall include successfully
	// provisioned PFD data of application(s).
	PfdDatas map[string]PfdData `json:"pfdDatas"`
	// Supplied by the AF and contains the external application identifiers
	// for which PFD(s) are not added or modified successfully. The failure
	// reason is also included. Each element provides the related information
	// for one or more external application identifier(s) and is identified in
	// the map via the failure identifier as key.
	PfdReports map[string]PfdReport `json:"pfdReports,omitempty"`
}

// FailureCode represents the failure reason of the PFD management
type FailureCode string

// Possible values of FailureCode
const (
	// This value indicates that something functions wrongly in PFD
	// provisioning or the PFD provisioning does not function at all.
	Malfunction FailureCode = "MALFUNCTION"
	// This value indicates there is resource limitation for PFD storage.
	ResourceLimitation = "RESOURCE_LIMITATION"
	// This value indicates that the allowed delay is too short and PFD(s) are
	// not stored
	ShortDelay = "SHORT_DELAY"
	// The received external application identifier(s) are already provisioned
	AppIDDuplicated = "APP_ID_DUPLICATED"
	// Other reason specified
	OtherReason = "OTHER_REASON"
)

// PfdReport is the type that represents a PFD report to indicate the
// external application identifier(s) which PFD(s) are not added or
// modified successfully and corresponding failure reason.
type PfdReport struct {
	// Identifies the external application identifier(s) which PFD(s) are not
	// added or modified successfully
	ExternalAppIds []string `json:"externalAppIds"`
	// Identifies the failure reason
	FailureCode FailureCode `json:"failureCode"`
	// It shall be included when the allowed delayed cannot be satisfied, i.e.
	// it is smaller than the caching time configured in fetching PFD.
	CachingTime DurationSec `json:"cachingTime,omitempty"`
}
