/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngcnef

// AppSessionContext Represents an Individual Application Session Context
// resource.
type AppSessionContext struct {
	// Identifies the service requirements of an Individual Application Session
	// Context.
	// It shall be present in HTTP POST request messages and may be included
	// in the HTTP response messages
	AscReqData AppSessionContextReqData `json:"ascReqData,omitempty"`
	// The following fields will not be supported as its not required for
	// Traffic Influence
	AscRespData AppSessionContextRespData `json:"ascRespData,omitempty"`
	EvsNotif    EventsNotification        `json:"evsNotif,omitempty"`
}

//AppSessionContextRespData added
type AppSessionContextRespData struct {
	ServAuthInfo ServAuthInfo `json:"servAuthInfo,omitempty" bson:"servAuthInfo"`
	SuppFeat     string       `json:"suppFeat,omitempty" bson:"suppFeat"`
}

//ServAuthInfo added
type ServAuthInfo string

//EventsNotification describes the notification of a matched event
type EventsNotification struct {
	AccessType AccessType `json:"accessType,omitempty" bson:"accessType"`
	//AnGwAddr                  *AnGwAddress                 `json:"anGwAddr,omitempty" bson:"anGwAddr"`
	EvSubsURI string                `json:"evSubsUri" bson:"evSubsUri"`
	EvNotifs  []AfEventNotification `json:"evNotifs" bson:"evNotifs"`
	//FailedResourcAllocReports []ResourcesAllocationInfo    `json:"failedResourcAllocReports,omitempty" bson:"failedResourcAllocReports"`
	//PlmnId                    *PlmnId                      `json:"plmnId,omitempty" bson:"plmnId"`
	//QncReports                []QosNotificationControlInfo `json:"qncReports,omitempty" bson:"qncReports"`
	//RatType                   RatType                      `json:"ratType,omitempty" bson:"ratType"`
	//UsgRep                    *AccumulatedUsage            `json:"usgRep,omitempty" bson:"usgRep"`
}

//AfEventNotification describes the event information delivered in the notification
type AfEventNotification struct {
	Event AfEvent `json:"event" bson:"event"`
	//Flows []Flows `json:"flows,omitempty" bson:"flows"`
}

//AfEvent string
type AfEvent string

// AppSessionContextReqData Identifies the service requirements of an
// Individual Application Session Context.
type AppSessionContextReqData struct {
	// AF application identifier
	AfAppID AfAppID `json:"afAppId,omitempty"`
	// Indicates the AF traffic routing requirements. It shall be included if
	//  Influence on Traffic Routing feature is supported
	AfRoutReq AfRoutingRequirement `json:"afRoutReq,omitempty"`
	// dnn Data Network Name. It shall be present when the "afRoutReq"
	// attribute is  present
	Dnn Dnn `json:"dnn,omitempty"`
	// slice info
	SliceInfo Snssai `json:"sliceInfo,omitempty"`
	// Subscription Permanent Identifier
	Supi Supi `json:"supi,omitempty"`
	// Subscription Permanent Identifier
	Gpsi Gpsi `json:"gpsi,omitempty"`
	// supp feat
	// Required: true
	SuppFeat SupportedFeatures `json:"suppFeat"`
	// ue Ipv4
	UeIpv4 Ipv4Addr `json:"ueIpv4,omitempty"`
	// ue Ipv6
	UeIpv6 Ipv6Addr `json:"ueIpv6,omitempty"`
	// ue mac
	UeMac    MacAddr48 `json:"ueMac,omitempty"`
	NotifURI string    `json:"notifUri" bson:"notifUri"`
	// The following fields have been omitted as they are not required for
	// Traffic Influ feature
	// evSubsc and notifUri - Not Required
	// AspId - Required when Sponspored Connnectivity is supported
	// bdtRefId - Required when BDT is supported
	// ipDomain - Required when Qos is supported
	// medComponents - Required when Qos is supported
	// mpsId - Required when Multimedia Priority Service is supported
	// sponId - Required when Sponspored Connnectivity is supported
	// sponStatus - Required when Sponspored Connnectivity is supported
}

// AppSessionContextUpdateData Contains the modification(s) to apply to the
// Individual Application Session Context resource
type AppSessionContextUpdateData struct {
	// AF application identifier
	AfAppID AfAppID `json:"afAppId,omitempty"`
	// Indicates the AF traffic routing requirements. It shall be included if
	//  Influence on Traffic Routing feature is supported
	AfRoutReq AfRoutingRequirement `json:"afRoutReq,omitempty"`

	// The following fields have been omitted as they are not required for
	// Traffic Influ feature
	// evSubsc and notifUri - Not Required
	// AspId - Required when Sponspored Connnectivity is supported
	// bdtRefId - Required when BDT is supported
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
	// A list of traffic routes to applications locations
	// Min Items: 1
	RouteToLocs []RouteToLocation `json:"routeToLocs"`
	// sp val
	SpVal *SpatialValidity `json:"spVal,omitempty"`
	// temp vals
	// Min Items: 1
	TempVals []TemporalValidity `json:"tempVals"`
	// up path chg sub
	UpPathChgSub UpPathChgEvent `json:"upPathChgSub,omitempty"`
	// Indication of application relocation possibility
	AppReloc bool `json:"appReloc,omitempty"`
	// addr preser ind
	AddrPreserInd bool `json:"addrPreserInd,omitempty"`
}

//FlowDescription : Defines a packet filter of an IP flow.
type FlowDescription string

// FlowDirection : Possible values are
//  - DOWNLINK: The corresponding filter applies for traffic to the UE.
//  - UPLINK: The corresponding filter applies for traffic from the UE.
//  - BIDIRECTIONAL: The corresponding filter applies for traffic both to and
// from the UE.
//  - UNSPECIFIED: The corresponding filter applies for traffic to the UE
// (downlink), but has no specific direction declared.
//                 The service data flow detection shall apply the filter for
// uplink traffic as if the filter was bidirectional.
//                 The PCF shall not use the value UNSPECIFIED in filters
// created by the network in NW-initiated procedures.
//                 The PCF shall only include the value UNSPECIFIED in filters
// in UE-initiated procedures if the same value is received from the SMF.
// TBD: To be moved to the SMF policy related file later on
type FlowDirection string

// EthFlowDescription :  Identifies an Ethernet flow
type EthFlowDescription struct {
	// destination Mac address
	DestMacAddr MacAddr48 `json:"destMacAddr,omitempty"`
	// A two-octet string that represents the Ethertype, as described in
	// IEEE 802.3 [16] and IETF RFC 7042 [18]
	// in hexadecimal representation.  Each character in the string shall take a
	// value of "0" to "9" or "A" to "F" and shall represent 4 bits. The most
	// significant
	// character representing the 4 most  significant bits of the ethType shall
	// appear first in the string, and the character representing the 4 least
	// significant bits of the ethType shall appear last in the string
	EthType string `json:"ethType,omitempty"`
	// Contains the flow description for the Uplink or Downlink IP flow. It
	// shall be present when the Ethertype is IP.
	FDesc FlowDescription `json:"fDesc,omitempty"`
	// Contains the packet filter direction.
	FDir FlowDirection `json:"fDir,omitempty"`
	// Source MAC address
	SourceMacAddr MacAddr48 `json:"sourceMacAddr,omitempty"`
	//  minItems: 1, maxItems: 2
	// Customer-VLAN and/or Service-VLAN tags containing the VID, PCP/DEI fields
	// as defined in IEEE 802.1Q [17] and IETF RFC 7042 [18].
	// Each field is encoded as a two-octet string in hexadecimal
	// representation.
	// Each character in the string shall take a value of "0" to "9" or "A"
	// to "F" and shall
	// represent 4 bits. The most significant character representing the 4 most
	// significant bits of the VID or PCF/DEI field shall appear first in the
	// string, and the character representing the 4 least significant bits of
	// the VID or PCF/DEI field shall appear last in the string
	VlanTags []string `json:"vlanTags,omitempty"`
}

// UpPathChgEvent : UP path management events
// To be moved to SMPolicy file when available
type UpPathChgEvent struct {
	// notification URI
	NotificationURI URI `json:"notificationUri"`
	// It is used to set the value of Notification Correlation ID in the n
	// otification sent by the SMF.
	NotifCorreID string `json:"notifCorreId"`
	// DNAI change type to be notified
	DnaiChgType DnaiChangeType `json:"dnaiChgType"`
}
