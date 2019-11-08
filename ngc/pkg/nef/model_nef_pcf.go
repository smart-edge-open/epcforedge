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

// AppSessionContext Represents an Individual Application Session Context
// resource.
type AppSessionContext struct {
	// Identifies the service requirements of an Individual Application Session Context.
	// It shall be present in HTTP POST request messages and may be included
	// in the HTTP response messages
	AscReqData AppSessionContextReqData `json:"ascReqData,omitempty"`
	// The following fields will not be supported as its not required for Traffic Influence
	// AscRespData AppSessionContextRespData `json:"ascRespData,omitempty"`
	// EvsNotif EventsNotification `json:"evsNotif,omitempty"`

}

// AppSessionContextReqData Identifies the service requirements of an
// Individual Application Session Context.
type AppSessionContextReqData struct {
	// AF application identifier
	afAppID AfAppID `json:"afAppId,omitempty"`
	// Indicates the AF traffic routing requirements. It shall be included if Influence on Traffic Routing feature is supported
	afRoutReq AfRoutingRequirement `json:"afRoutReq,omitempty"`
	// dnn Data Network Name. It shall be present when the "afRoutReq" attribute is  present
	dnn Dnn `json:"dnn,omitempty"`
	// slice info
	sliceInfo Snssai `json:"sliceInfo,omitempty"`
	// Subscription Permanent Identifier
	supi Supi `json:"supi,omitempty"`
	// Subscription Permanent Identifier
	gpsi Gpsi `json:"gpsi,omitempty"`
	// supp feat
	// Required: true
	SuppFeat SupportedFeatures `json:"suppFeat"`
	// ue Ipv4
	ueIpv4 Ipv4Addr `json:"ueIpv4,omitempty"`
	// ue Ipv6
	ueIpv6 Ipv6Addr `json:"ueIpv6,omitempty"`
	// ue mac
	ueMac MacAddr48 `json:"ueMac,omitempty"`

	// TBD : check if evSubsc and notifUri are required for TrafficInflu
	// The following fields have been ommitted as they are not required for Traffic Influ feature
	// AspId - Required when Sponspored Connnectivity is supported
	// bdtRefId - Required when BDT is supported
	// ipDomain - Required when Qos is supported
	// medComponents - Required when Qos is supported
	// mpsId - Required when Multimedia Priority Service is supported
	// sponId - Required when Sponspored Connnectivity is supported
	// sponStatus - Required when Sponspored Connnectivity is supported
}

// AfAppID Contains an AF application identifier.
type AfAppID string

// AfRoutingRequirement describes the event information delivered in the
// subscription
type AfRoutingRequirement struct {
	// Indication of application relocation possibility
	appReloc bool `json:"appReloc,omitempty"`
	// A list of traffic routes to applications locations
	// Min Items: 1
	routeToLocs []RouteToLocation `json:"routeToLocs"`
	// sp val
	spVal SpatialValidity `json:"spVal,omitempty"`
	// addr preser ind
	AddrPreserInd bool `json:"addrPreserInd,omitempty"`
	// temp vals
	// Min Items: 1
	TempVals []TemporalValidity `json:"tempVals"`
	// up path chg sub
	UpPathChgSub UpPathChgEvent `json:"upPathChgSub,omitempty"`
}

//FlowDescription : Defines a packet filter of an IP flow.
type FlowDescription string

// FlowDirection : Possible values are
//  - DOWNLINK: The corresponding filter applies for traffic to the UE.
//  - UPLINK: The corresponding filter applies for traffic from the UE.
//  - BIDIRECTIONAL: The corresponding filter applies for traffic both to and from the UE.
//  - UNSPECIFIED: The corresponding filter applies for traffic to the UE (downlink), but has no specific direction declared.
//                 The service data flow detection shall apply the filter for uplink traffic as if the filter was bidirectional.
//                 The PCF shall not use the value UNSPECIFIED in filters created by the network in NW-initiated procedures.
//                 The PCF shall only include the value UNSPECIFIED in filters in UE-initiated procedures if the same value is received from the SMF.
// TBD: To be moved to the SMF policy related file later on
type FlowDirection string

// EthFlowDescription :  Identifies an Ethernet flow
type EthFlowDescription struct {
	// destination Mac address
	destMacAddr MacAddr48 `json:"destMacAddr,omitempty"`
	// A two-octet string that represents the Ethertype, as described in IEEE 802.3 [16] and IETF RFC 7042 [18]
	// in hexadecimal representation.  Each character in the string shall take a
	// value of "0" to "9" or "A" to "F" and shall represent 4 bits. The most significant
	// character representing the 4 most  significant bits of the ethType shall
	// appear first in the string, and the character representing the 4 least
	// significant bits of the ethType shall appear last in the string
	ethType string `json:"ethType, omitempty"`
	// Contains the flow description for the Uplink or Downlink IP flow. It shall be
	// present when the Ethertype is IP.
	fDesc FlowDescription `json:"fDesc, omitempty"`
	// Contains the packet filter direction.
	fDir FlowDirection `json:"fDir, omitempty"`
	// Source MAC address
	sourceMacAddr MacAddr48 `json:"sourceMacAddr, omitempty"`
	//  minItems: 1, maxItems: 2
	// Customer-VLAN and/or Service-VLAN tags containing the VID, PCP/DEI fields
	// as defined in IEEE 802.1Q [17] and IETF RFC 7042 [18].
	// Each field is encoded as a two-octet string in hexadecimal representation.
	// Each character in the string shall take a value of "0" to "9" or "A" to "F" and shall
	// represent 4 bits. The most significant character representing the 4 most
	// significant bits of the VID or PCF/DEI field shall appear first in the string, and
	// the character representing the 4 least significant bits of the VID or PCF/DEI
	// field shall appear last in the string
	vlanTags []string `json:"vlanTags, omitempty"`
}

// UpPathChgEvent : UP path management events
type UpPathChgEvent struct {
	// notification URI
	notificationURI URI `json:"notificationUri`
	// It is used to set the value of Notification Correlation ID in the notification sent by the SMF.
	notifCorreID string `json:"notifCorreId"`
	// DNAI change type to be notified
	dnaiChgType DnaiChangeType
}
