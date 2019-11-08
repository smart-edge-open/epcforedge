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

package main

// URI :  string formatted accordingding to IETF RFC 3986
type URI string

// Dnai : string identifying the Data Network Area Identifier
type Dnai string

// DnaiChangeType : string identifying the DNAI change type
// Possible values are
// - EARLY: Early notification of UP path reconfiguration.
// - EARLY_LATE: Early and late notification of UP path reconfiguration. This value shall only be present in the subscription to the DNAI change event.
// - LATE: Late notification of UP path reconfiguration.
type DnaiChangeType string

// Dnn : string identify the Data network name
type Dnn string

// ExternalID : string containing a local identifier followed by "@" and a domain identifier. B
// Both the local identifier and the domain identifier shall be encoded as strings that do not contain any "@" characters.
// See Clause 4.6.2 of 3GPP TS 23.682 for more information
type ExternalID string

// ExternalGroupID : string containing a local identifier followed by "@" and a domain identifier.
// Both the local identifier and the domain identifier shall be encoded as strings that do not contain any "@" characters.
// See Clauses 4.6.2 and 4.6.3 of 3GPP TS 23.682 for more information
type ExternalGroupID string

// FlowInfo Flow information struct
type FlowInfo struct {
	// Indicates the IP flow.
	FlowID int32 `json:"flowId"`
	// Indicates the packet filters of the IP flow. Refer to subclause 5.3.8 of 3GPP TS 29.214 for encoding.
	// It shall contain UL and/or DL IP flow description.
	// minItems : 1 maxItems : 2
	FlowDescriptions []string `json:"flowDescriptions,omitempty"`
}

// Supi : Subscription Permanent Identifier
// pattern: '^(imsi-[0-9]{5,15}|nai-.+|.+)$'
type Supi string

// Gpsi : Generic Public Servie Identifiers asssociated wit the UE
// pattern '^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$'
type Gpsi string

// Ipv4Addr : string representing the IPv4 address
// pattern: '^(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])$'
// example: '198.51.100.1'
type Ipv4Addr string

// Ipv6Addr : string representing the IPv6 address
// pattern: '^((:|(0?|([1-9a-f][0-9a-f]{0,3}))):)((0?|([1-9a-f][0-9a-f]{0,3})):){0,6}(:|(0?|([1-9a-f][0-9a-f]{0,3})))$'
// pattern: '^((([^:]+:){7}([^:]+))|((([^:]+:)*[^:]+)?::(([^:]+:)*[^:]+)?))$'
// example: '2001:db8:85a3::8a2e:370:7334'
type Ipv6Addr string

// Ipv6Prefix : string representing the Ipv6 Prefix
// pattern: '^((:|(0?|([1-9a-f][0-9a-f]{0,3}))):)((0?|([1-9a-f][0-9a-f]{0,3})):){0,6}(:|(0?|([1-9a-f][0-9a-f]{0,3})))(\/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))$'
// pattern: '^((([^:]+:){7}([^:]+))|((([^:]+:)*[^:]+)?::(([^:]+:)*[^:]+)?))(\/.+)$'
// example: '2001:db8:abcd:12::0/64'
type Ipv6Prefix string

// Link : string Identifies a referenced resource
type Link URI

// MacAddr48 : Identifies a MAC address
// pattern: '^([0-9a-fA-F]{2})((-[0-9a-fA-F]{2}){5})$'
type MacAddr48 string

// Port : Identifies a port number
// Unsigned integer with valid values between 0 and 65535.
type Port uint32

// RouteInformation Route information struct
type RouteInformation struct {
	// string identifying a Ipv4 address formatted in the \"dotted decimal\"
	// notation as defined in IETF RFC 1166.
	ipv4Addr Ipv4Addr `json:"ipv4Addr,omitempty"`
	// string identifying a Ipv6 address formatted according to clause 4 in
	// IETF RFC 5952.
	// The mixed Ipv4 Ipv6 notation according to clause 5 of IETF RFC 5952 shall
	// not be used.
	ipv6Addr Ipv6Addr `json:"ipv6Addr,omitempty"`
	// Port number
	portNumber uint32 `json:"portNumber"`
}

// RouteToLocation : Describes the traffic routes to the locations of the application
type RouteToLocation struct {
	// Data network access identifier
	dnai Dnai `json:"dnai"`
	// Additional route information about the route to Dnai
	routeInfo RouteInformation `json:"routeInfo,omitempty"`
	// Dnai route profile identifier
	routeProfID string `json:"routeProfId,omitempty"`
}

// Snssai Network slice identifier
type Snssai struct {
	// minimum: 0, 	maximum: 255
	sst uint8 `json:"sst"`
	// pattern: '^[A-Fa-f0-9]{6}$'
	sd string `json:"sd,omitempty"`
}

// SupportedFeatures : A string used to indicate the features supported by an
// API that is used
// (subclause 6.6 in 3GPP TS 29.500).
// The string shall contain a bitmask indicating supported features in
// hexadecimal representation.
// Each character in the string shall take a value of "0" to "9" or "A" to "F"
// and shall represent the support of 4 features as described in table 5.2.2-3.
// The most significant character representing the highest-numbered features
// shall appear first in the string,
// and the character representing features 1 to 4 shall appear last
// in the string.
// The list of features and their numbering (starting with 1)
// are defined separately for each API.
// Possible features for traffic influencing are
// Notification_websocket( takes vlue of 1) and
// Notification_test_event(takes value of 2)
// pattern: '^[A-Fa-f0-9]*$'
type SupportedFeatures string

// WebsockNotifConfig Websocket noticcation configuration
type WebsockNotifConfig struct {
	// string formatted according to IETF RFC 3986 identifying a
	// referenced resource.
	websocketURI Link `json:"websocketUri,omitempty"`
	// Set by the AF to indicate that the Websocket delivery is requested.
	requestWebsocketURI bool `json:"requestWebsocketUri,omitempty"`
}

// ProblemDetails Problem details struct
type ProblemDetails struct {
	// problem type
	Type Link `json:"type,omitempty"`
	// A short, human-readable summary of the problem type.
	// It should not change from occurrence to occurrence of the problem.
	Title string `json:"title,omitempty"`
	// A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`
	// URL to problem instance
	Instance Link `json:"instance,omitempty"`
	// A machine-readable application error cause specific to this occurrence
	// of the problem.
	// This IE should be present and provide application-related error
	// information, if available.
	Cause string `json:"cause,omitempty"`
	// Description of invalid parameters, for a request rejected due to
	// invalid parameters.
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
	// The HTTP status code for this occurrence of the problem.
	Status int32 `json:"status,omitempty"`
}

// InvalidParam Invalid Parameters struct
type InvalidParam struct {
	// Attribute''s name encoded as a JSON Pointer, or header''s name.
	Param string `json:"param"`
	// A human-readable reason, e.g. \"must be a positive integer\".
	Reason string `json:"reason,omitempty"`
}

// PresenceState presence state
type PresenceState string

// Possible values of Presence State
const (

	// PresenceStateINAREA captures enum value "IN_AREA"
	PresenceStateINAREA PresenceState = "IN_AREA"

	// PresenceStateOUTOFAREA captures enum value "OUT_OF_AREA"
	PresenceStateOUTOFAREA PresenceState = "OUT_OF_AREA"

	// PresenceStateUNKNOWN captures enum value "UNKNOWN"
	PresenceStateUNKNOWN PresenceState = "UNKNOWN"

	// PresenceStateINACTIVE captures enum value "INACTIVE"
	PresenceStateINACTIVE PresenceState = "INACTIVE"
)

// Mcc mcc
type Mcc string

// Mnc mnc
type Mnc string

// PlmnID plmn Id
type PlmnID struct {

	// mcc
	// Required: true
	mcc Mcc `json:"mcc"`

	// mnc
	// Required: true
	mnc Mnc `json:"mnc"`
}

// Tac tac
type Tac string

// Tai tai
type Tai struct {

	// plmn Id
	// Required: true
	plmnID PlmnID `json:"plmnId"`

	// tac
	// Required: true
	tac Tac `json:"tac"`
}

// EutraCellID eutra cell Id
type EutraCellID string

// Ecgi ecgi
type Ecgi struct {
	// eutra cell Id
	// Required: true
	eutraCellID EutraCellID `json:"eutraCellId"`
	// plmn Id
	// Required: true
	plmnID PlmnID `json:"plmnId"`
}

// NrCellID nr cell Id
type NrCellID string

// Ncgi ncgi
type Ncgi struct {
	// nr cell Id
	// Required: true
	nrCellID NrCellID `json:"nrCellId"`
	// plmn Id
	// Required: true
	plmnID PlmnID `json:"plmnId"`
}

// GNbID g nb Id
type GNbID struct {
	// bit length
	// Required: true
	// Maximum: 32
	// Minimum: 22
	bitLength uint8 `json:"bitLength"`
	// g n b value
	// Required: true
	// Pattern: ^[A-Fa-f0-9]{6,8}$
	gNBValue string `json:"gNBValue"`
}

// N3IwfID n3 iwf Id
type N3IwfID string

// NgeNbID nge nb Id
type NgeNbID string

// GlobalRanNodeID global ran node Id
type GlobalRanNodeID struct {
	// plmn Id
	// Required: true
	plmnID PlmnID `json:"plmnId"`
	// n3 iwf Id
	N3IwfID N3IwfID `json:"n3IwfId,omitempty"`
	// g nb Id
	gNbID GNbID `json:"gNbId,omitempty"`
	// nge nb Id
	NgeNbID NgeNbID `json:"ngeNbId,omitempty"`
}

// PresenceInfo presence info
type PresenceInfo struct {
	// pra Id
	praID string `json:"praId,omitempty"`
	// presence state
	presenceState PresenceState `json:"presenceState,omitempty"`
	// ecgi list
	// Min Items: 1
	EcgiList []Ecgi `json:"ecgiList"`
	// ncgi list
	// Min Items: 1
	NcgiList []Ncgi `json:"ncgiList"`
	// global ran node Id list
	// Min Items: 1
	globalRanNodeIDList []GlobalRanNodeID `json:"globalRanNodeIdList"`
}

// SpatialValidity Describes the spatial validity of an AF request for influencing traffic routing
type SpatialValidity struct {
	presenceInfoList PresenceInfo `json:"presenceInfoList"`
}

// DateTime is in the date-time format
type DateTime string
