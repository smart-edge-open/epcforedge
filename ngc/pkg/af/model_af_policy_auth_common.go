// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

/*
 * This file include structs which are common to traffic_influence and policy
 * auth model. These struct are commented to pass compilation.
 * If using AF Policy authorization as Saperate entity, then remove commented
 * struct.
 */
/* --Remove-this-comment
// DNAIChangeType type
type DNAIChangeType string

// List of DnaiChangeType
const (
	Early     DNAIChangeType = "EARLY"
	EarlyLate DNAIChangeType = "EARLY_LATE"
	Late      DNAIChangeType = "LATE"
)

type IPv4Addr string

type IPv6Addr string

// EthFlowDescription Identifies an Ethernet flow
type EthFlowDescription struct {
	DestMacAddr string `json:"destMacAddr,omitempty"`
	//EtherType number
	EthType string `json:"ethType"`
	// Defines a packet filter of an IP flow.
	FDesc string `json:"fDesc,omitempty"`
	// Possible values are DOWNLINK - The corresponding filter applies for
	// traffic to the UE.
	// UPLINK - The corresponding filter applies for traffic from the UE.
	// BIDIRECTIONAL The corresponding filter applies for traffic both to
	// and from the UE.
	// UNSPECIFIED - The corresponding filter applies for traffic to the UE
	// (downlink), but has no specific direction declared.
	// The service data flow detection shall apply the filter for uplink
	// traffic as if the filter was bidirectional.
	FDir string `json:"fDir,omitempty"`
	// Source mac address
	SourceMacAddr string `json:"sourceMacAddr,omitempty"`
	// Vlan tags
	VLANTags []string `json:"vlanTags,omitempty"`
}

// InvalidParam Invalid Parameters struct
type InvalidParam struct {
	// Attribute''s name encoded as a JSON Pointer, or header''s name.
	Param string `json:"param"`
	// A human-readable reason, e.g. \"must be a positive integer\".
	Reason string `json:"reason,omitempty"`
}

type ProblemDetails struct {
	// string providing an URI formatted according to IETF RFC 3986.
	Type string `json:"type,omitempty"`
	/*
	 * A short, human-readable summary of the problem type. It should not
	 * change from occurrence to occurrence of the problem.
	 --Add-multiline-close
	Title string `json:"title,omitempty"`
	// The HTTP status code for this occurrence of the problem.
	Status int32 `json:"status,omitempty"`
	/*
	 * A human-readable explanation specific to this occurrence of the
	 * problem.
	 --Add-multiline-close
	Detail string `json:"detail,omitempty"`
	// string providing an URI formatted according to IETF RFC 3986.
	Instance string `json:"instance,omitempty"`
	/*
	 * A machine-readable application error cause specific to this
	 * occurrence of the problem. This IE should be present and provide
	 * application-related error information, if available.
	 --Add-multiline-close
	Cause string `json:"cause,omitempty"`
	/*
	 * Description of invalid parameters, for a request rejected due to
	 * invalid parameters.
	 --Add-multiline-close
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

// RouteInformation Route information struct
type RouteInformation struct {
	// string identifying a Ipv4 address formatted in the \"dotted decimal\"
	// notation as defined in IETF RFC 1166.
	IPv4Addr IPv4Addr `json:"ipv4Addr,omitempty"`
	// string identifying a Ipv6 address formatted according to clause 4 in
	// IETF RFC 5952.
	// The mixed Ipv4 Ipv6 notation according to clause 5 of IETF RFC 5952
	// shall not be used.
	IPv6Addr IPv6Addr `json:"ipv6Addr,omitempty"`
	// Port number
	PortNumber int32 `json:"portNumber"`
}

// RouteToLocation Route to location structure
type RouteToLocation struct {
	// Data network access identifier
	DNAI DNAI `json:"dnai"`
	// Dnai route profile identifier
	RouteProfID string `json:"routeProfId,omitempty"`
	// Additional route information about the route to Dnai
	RouteInfo RouteInformation `json:"routeInfo,omitempty"`
}
--Remove-this-comment*/

type Snssai struct {
	Sst int32  `json:"sst"`
	Sd  string `json:"sd,omitempty"`
}

/* -- Additional
// TemporalValidity Indicates the time interval(s) during which the AF request
// is to be applied
type TemporalValidity struct {
	// string with format \"date-time\" as defined in OpenAPI.
	StartTime string `json:"startTime,omitempty"`
	// string with format \"date-time\" as defined in OpenAPI.
	StopTime string `json:"stopTime,omitempty"`
}
 -- Additional */
