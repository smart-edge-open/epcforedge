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

import (
	strfmt "github.com/go-openapi/strfmt"
)

// AfAppID Contains an AF application identifier.
type AfAppID string

// Uint32 uint32
type Uint32 int32

// Uint32Rm uint32 rm
type Uint32Rm int32

// ChargingID charging ID
type ChargingID = Uint32

// AfRequestedData af requested data
type AfRequestedData string

const (
	// AfRequestedDataUEIDENTITY captures enum value "UE_IDENTITY"
	AfRequestedDataUEIDENTITY AfRequestedData = "UE_IDENTITY"
)

// EutraCellID eutra cell Id
type EutraCellID string

// Mcc mcc
type Mcc string

// Mnc mnc
type Mnc string

// PlmnID plmn Id
type PlmnID struct {

	// mcc
	// Required: true
	Mcc Mcc `json:"mcc"`

	// mnc
	// Required: true
	Mnc Mnc `json:"mnc"`
}

// Ecgi ecgi
type Ecgi struct {

	// eutra cell Id
	// Required: true
	EutraCellID EutraCellID `json:"eutraCellId"`

	// plmn Id
	// Required: true
	PlmnID *PlmnID `json:"plmnId"`
}

// GNbID g nb Id
type GNbID struct {

	// bit length
	// Required: true
	// Maximum: 32
	// Minimum: 22
	BitLength *int64 `json:"bitLength"`

	// g n b value
	// Required: true
	// Pattern: ^[A-Fa-f0-9]{6,8}$
	GNBValue *string `json:"gNBValue"`
}

// N3IwfID n3 iwf Id
type N3IwfID string

// NgeNbID nge nb Id
type NgeNbID string

// GlobalRanNodeID global ran node Id
type GlobalRanNodeID struct {

	// g nb Id
	GNbID *GNbID `json:"gNbId,omitempty"`

	// n3 iwf Id
	N3IwfID N3IwfID `json:"n3IwfId,omitempty"`

	// nge nb Id
	NgeNbID NgeNbID `json:"ngeNbId,omitempty"`

	// plmn Id
	// Required: true
	PlmnID *PlmnID `json:"plmnId"`
}

// NrCellID nr cell Id
type NrCellID string

// Ncgi ncgi
type Ncgi struct {

	// nr cell Id
	// Required: true
	NrCellID NrCellID `json:"nrCellId"`

	// plmn Id
	// Required: true
	PlmnID *PlmnID `json:"plmnId"`
}

// PresenceState presence state
type PresenceState string

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

// Tac tac
type Tac string

// Tai tai
type Tai struct {

	// plmn Id
	// Required: true
	PlmnID *PlmnID `json:"plmnId"`

	// tac
	// Required: true
	Tac Tac `json:"tac"`
}

// PresenceInfo presence info
type PresenceInfo struct {

	// ecgi list
	// Min Items: 1
	EcgiList []*Ecgi `json:"ecgiList"`

	// global ran node Id list
	// Min Items: 1
	GlobalRanNodeIDList []*GlobalRanNodeID `json:"globalRanNodeIdList"`

	// ncgi list
	// Min Items: 1
	NcgiList []*Ncgi `json:"ncgiList"`

	// pra Id
	PraID string `json:"praId,omitempty"`

	// presence state
	PresenceState PresenceState `json:"presenceState,omitempty"`

	// tracking area list
	// Min Items: 1
	TrackingAreaList []*Tai `json:"trackingAreaList"`
}

// SpatialValidity describes explicitly the route to an Application location
type SpatialValidity struct {

	// presence info list
	// Required: true
	PresenceInfoList map[string]PresenceInfo `json:"presenceInfoList"`
}

// UpPathChgEvent up path chg event
type UpPathChgEvent struct {

	// dnai chg type
	// Enum: [EARLY EARLY_LATE LATE]
	DnaiChgType string `json:"dnaiChgType,omitempty"`

	// It is used to set the value of Notification Correlation ID in the
	// notification sent by the SMF.
	NotifCorreID string `json:"notifCorreId,omitempty"`

	// notification Uri
	NotificationURI string `json:"notificationUri,omitempty"`
}

// AfRoutingRequirement describes the event information delivered in the
// subscription
type AfRoutingRequirement struct {

	// addr preser ind
	AddrPreserInd bool `json:"addrPreserInd,omitempty"`

	// app reloc
	AppReloc bool `json:"appReloc,omitempty"`

	// route to locs
	// Min Items: 1
	RouteToLocs []*RouteToLocation `json:"routeToLocs"`

	// sp val
	SpVal *SpatialValidity `json:"spVal,omitempty"`

	// temp vals
	// Min Items: 1
	TempVals []*TemporalValidity `json:"tempVals"`

	// up path chg sub
	UpPathChgSub *UpPathChgEvent `json:"upPathChgSub,omitempty"`
}

// AppSessionContext Represents an Individual Application Session Context
// resource.
type AppSessionContext struct {

	// asc req data
	AscReqData *AppSessionContextReqData `json:"ascReqData,omitempty"`

	// asc resp data
	AscRespData *AppSessionContextRespData `json:"ascRespData,omitempty"`

	// evs notif
	EvsNotif *EventsNotification `json:"evsNotif,omitempty"`
}

// AspID Contains an identity of an application service provider.
type AspID string

// BdtReferenceID string identifying a BDT Reference ID as defined in subclause
// 5.3.3 of 3GPP TS 29.154.
type BdtReferenceID string

// Dnn dnn
type Dnn string

// AfNotifMethod af notif method
type AfNotifMethod string

// AfEvent af event
type AfEvent string

const (

	// AfEventACCESSTYPECHANGE captures enum value "ACCESS_TYPE_CHANGE"
	AfEventACCESSTYPECHANGE AfEvent = "ACCESS_TYPE_CHANGE"

	// AfEventANIREPORT captures enum value "ANI_REPORT"
	AfEventANIREPORT AfEvent = "ANI_REPORT"

	// AfEventCHARGINGCORRELATION captures enum value "CHARGING_CORRELATION"
	AfEventCHARGINGCORRELATION AfEvent = "CHARGING_CORRELATION"

	// AfEventFAILEDRESOURCESALLOCATION captures enum value
	// "FAILED_RESOURCES_ALLOCATION"
	AfEventFAILEDRESOURCESALLOCATION AfEvent = "FAILED_RESOURCES_ALLOCATION"

	// AfEventOUTOFCREDIT captures enum value "OUT_OF_CREDIT"
	AfEventOUTOFCREDIT AfEvent = "OUT_OF_CREDIT"

	// AfEventPLMNCHG captures enum value "PLMN_CHG"
	AfEventPLMNCHG AfEvent = "PLMN_CHG"

	// AfEventQOSNOTIF captures enum value "QOS_NOTIF"
	AfEventQOSNOTIF AfEvent = "QOS_NOTIF"

	// AfEventSUCCESSFULRESOURCESALLOCATION captures enum value
	// "SUCCESSFUL_RESOURCES_ALLOCATION"
	AfEventSUCCESSFULRESOURCESALLOCATION AfEvent = "SUCCESSFUL_RESOURCES_ALLOCATION"

	// AfEventUSAGEREPORT captures enum value "USAGE_REPORT"
	AfEventUSAGEREPORT AfEvent = "USAGE_REPORT"
)

// AfEventSubscription describes the event information delivered in the
// subscription
type AfEventSubscription struct {

	// event
	// Required: true
	Event AfEvent `json:"event"`

	// notif method
	NotifMethod AfNotifMethod `json:"notifMethod,omitempty"`
}

// URI string providing an URI formatted according to IETF RFC 3986.
type URI string

// RequiredAccessInfo required access info
type RequiredAccessInfo string

// Volume Unsigned integer identifying a volume in units of bytes.
type Volume int64

// DurationSec Unsigned integer identifying a period of time in units of
// seconds.
type DurationSec int64

// UsageThreshold usage threshold
type UsageThreshold struct {

	// downlink volume
	DownlinkVolume Volume `json:"downlinkVolume,omitempty"`

	// duration
	Duration DurationSec `json:"duration,omitempty"`

	// total volume
	TotalVolume Volume `json:"totalVolume,omitempty"`

	// uplink volume
	UplinkVolume Volume `json:"uplinkVolume,omitempty"`
}

// EventsSubscReqData Identifies the events the application subscribes to.
type EventsSubscReqData struct {

	// events
	// Required: true
	// Min Items: 1
	Events []*AfEventSubscription `json:"events"`

	// notif Uri
	NotifURI URI `json:"notifUri,omitempty"`

	// req ani
	ReqAni RequiredAccessInfo `json:"reqAni,omitempty"`

	// usg thres
	UsgThres *UsageThreshold `json:"usgThres,omitempty"`
}

// Gpsi gpsi
type Gpsi string

// CodecData Contains codec related information.
type CodecData string

// ContentVersion Represents the content version of some content.
type ContentVersion int64

// FlowStatus flow status
type FlowStatus string

// BitRate bit rate
type BitRate string

// FlowDescription Defines a packet filter of an IP flow.
type FlowDescription string

// FlowUsage flow usage
type FlowUsage string

// TosTrafficClass 2-octet string, where each octet is encoded in hexadecimal
// representation. The first octet contains the IPv4 Type-of-Service or the
// IPv6 Traffic-Class field and the second octet contains the ToS/Traffic
// Class mask field.
type TosTrafficClass string

// AfSigProtocol af sig protocol
type AfSigProtocol string

const (

	// AfSigProtocolNOINFORMATION captures enum value "NO_INFORMATION"
	AfSigProtocolNOINFORMATION AfSigProtocol = "NO_INFORMATION"

	// AfSigProtocolSIP captures enum value "SIP"
	AfSigProtocolSIP AfSigProtocol = "SIP"
)

// MediaSubComponent Identifies a media subcomponent
type MediaSubComponent struct {

	// af sig protocol
	AfSigProtocol AfSigProtocol `json:"afSigProtocol,omitempty"`

	// ethf descs
	// Max Items: 2
	// Min Items: 1
	EthfDescs []*EthFlowDescription `json:"ethfDescs"`

	// f descs
	// Max Items: 2
	// Min Items: 1
	FDescs []FlowDescription `json:"fDescs"`

	// f num
	// Required: true
	FNum int64 `json:"fNum"`

	// f status
	FStatus FlowStatus `json:"fStatus,omitempty"`

	// flow usage
	FlowUsage FlowUsage `json:"flowUsage,omitempty"`

	// mar bw dl
	MarBwDl BitRate `json:"marBwDl,omitempty"`

	// mar bw ul
	MarBwUl BitRate `json:"marBwUl,omitempty"`

	// tos tr cl
	TosTrCl TosTrafficClass `json:"tosTrCl,omitempty"`
}

// MediaType media type
type MediaType string

// PreemptionCapability preemption capability
type PreemptionCapability string

const (

	// PreemptionCapabilityNOTPREEMPT captures enum value "NOT_PREEMPT"
	PreemptionCapabilityNOTPREEMPT PreemptionCapability = "NOT_PREEMPT"

	// PreemptionCapabilityMAYPREEMPT captures enum value "MAY_PREEMPT"
	PreemptionCapabilityMAYPREEMPT PreemptionCapability = "MAY_PREEMPT"
)

// PreemptionVulnerability preemption vulnerability
type PreemptionVulnerability string

const (

	// PreemptionVulnerabilityNOTPREEMPTABLE captures enum value "NOT_PREEMPTABLE"
	PreemptionVulnerabilityNOTPREEMPTABLE PreemptionVulnerability = "NOT_PREEMPTABLE"

	// PreemptionVulnerabilityPREEMPTABLE captures enum value "PREEMPTABLE"
	PreemptionVulnerabilityPREEMPTABLE PreemptionVulnerability = "PREEMPTABLE"
)

// PrioritySharingIndicator priority sharing indicator
type PrioritySharingIndicator string

const (

	// PrioritySharingIndicatorENABLED captures enum value "ENABLED"
	PrioritySharingIndicatorENABLED PrioritySharingIndicator = "ENABLED"

	// PrioritySharingIndicatorDISABLED captures enum value "DISABLED"
	PrioritySharingIndicatorDISABLED PrioritySharingIndicator = "DISABLED"
)

// ReservPriority reserv priority
type ReservPriority string

const (

	// ReservPriorityPRIO1 captures enum value "PRIO_1"
	ReservPriorityPRIO1 ReservPriority = "PRIO_1"

	// ReservPriorityPRIO2 captures enum value "PRIO_2"
	ReservPriorityPRIO2 ReservPriority = "PRIO_2"

	// ReservPriorityPRIO3 captures enum value "PRIO_3"
	ReservPriorityPRIO3 ReservPriority = "PRIO_3"

	// ReservPriorityPRIO4 captures enum value "PRIO_4"
	ReservPriorityPRIO4 ReservPriority = "PRIO_4"

	// ReservPriorityPRIO5 captures enum value "PRIO_5"
	ReservPriorityPRIO5 ReservPriority = "PRIO_5"

	// ReservPriorityPRIO6 captures enum value "PRIO_6"
	ReservPriorityPRIO6 ReservPriority = "PRIO_6"

	// ReservPriorityPRIO7 captures enum value "PRIO_7"
	ReservPriorityPRIO7 ReservPriority = "PRIO_7"

	// ReservPriorityPRIO8 captures enum value "PRIO_8"
	ReservPriorityPRIO8 ReservPriority = "PRIO_8"

	// ReservPriorityPRIO9 captures enum value "PRIO_9"
	ReservPriorityPRIO9 ReservPriority = "PRIO_9"

	// ReservPriorityPRIO10 captures enum value "PRIO_10"
	ReservPriorityPRIO10 ReservPriority = "PRIO_10"

	// ReservPriorityPRIO11 captures enum value "PRIO_11"
	ReservPriorityPRIO11 ReservPriority = "PRIO_11"

	// ReservPriorityPRIO12 captures enum value "PRIO_12"
	ReservPriorityPRIO12 ReservPriority = "PRIO_12"

	// ReservPriorityPRIO13 captures enum value "PRIO_13"
	ReservPriorityPRIO13 ReservPriority = "PRIO_13"

	// ReservPriorityPRIO14 captures enum value "PRIO_14"
	ReservPriorityPRIO14 ReservPriority = "PRIO_14"

	// ReservPriorityPRIO15 captures enum value "PRIO_15"
	ReservPriorityPRIO15 ReservPriority = "PRIO_15"

	// ReservPriorityPRIO16 captures enum value "PRIO_16"
	ReservPriorityPRIO16 ReservPriority = "PRIO_16"
)

// MediaComponent Identifies a media component.
type MediaComponent struct {

	// af app Id
	AfAppID AfAppID `json:"afAppId,omitempty"`

	// af rout req
	AfRoutReq *AfRoutingRequirement `json:"afRoutReq,omitempty"`

	// codecs
	// Max Items: 2
	// Min Items: 1
	Codecs []CodecData `json:"codecs"`

	// cont ver
	ContVer ContentVersion `json:"contVer,omitempty"`

	// f status
	FStatus FlowStatus `json:"fStatus,omitempty"`

	// mar bw dl
	MarBwDl BitRate `json:"marBwDl,omitempty"`

	// mar bw ul
	MarBwUl BitRate `json:"marBwUl,omitempty"`

	// med comp n
	// Required: true
	MedCompN int64 `json:"medCompN"`

	// med sub comps
	MedSubComps map[string]MediaSubComponent `json:"medSubComps,omitempty"`

	// med type
	MedType MediaType `json:"medType,omitempty"`

	// mir bw dl
	MirBwDl BitRate `json:"mirBwDl,omitempty"`

	// mir bw ul
	MirBwUl BitRate `json:"mirBwUl,omitempty"`

	// preempt cap
	PreemptCap PreemptionCapability `json:"preemptCap,omitempty"`

	// preempt vuln
	PreemptVuln PreemptionVulnerability `json:"preemptVuln,omitempty"`

	// prio sharing ind
	PrioSharingInd PrioritySharingIndicator `json:"prioSharingInd,omitempty"`

	// res prio
	ResPrio ReservPriority `json:"resPrio,omitempty"`

	// rr bw
	RrBw BitRate `json:"rrBw,omitempty"`

	// rs bw
	RsBw BitRate `json:"rsBw,omitempty"`

	// sharing key dl
	SharingKeyDl Uint32 `json:"sharingKeyDl,omitempty"`

	// sharing key ul
	SharingKeyUl Uint32 `json:"sharingKeyUl,omitempty"`
}

// PreemptionControlInformation preemption control information
type PreemptionControlInformation string

const (

	// PreemptionControlInformationMOSTRECENT captures enum value "MOST_RECENT"
	PreemptionControlInformationMOSTRECENT PreemptionControlInformation = "MOST_RECENT"

	// PreemptionControlInformationLEASTRECENT captures enum value "LEAST_RECENT"
	PreemptionControlInformationLEASTRECENT PreemptionControlInformation = "LEAST_RECENT"

	// PreemptionControlInformationHIGHESTBW captures enum value "HIGHEST_BW"
	PreemptionControlInformationHIGHESTBW PreemptionControlInformation = "HIGHEST_BW"
)

// ServiceInfoStatus service info status
type ServiceInfoStatus string

const (

	// ServiceInfoStatusFINAL captures enum value "FINAL"
	ServiceInfoStatusFINAL ServiceInfoStatus = "FINAL"

	// ServiceInfoStatusPRELIMINARY captures enum value "PRELIMINARY"
	ServiceInfoStatusPRELIMINARY ServiceInfoStatus = "PRELIMINARY"
)

// ServiceUrn Contains values of the service URN and may include subservices.
type ServiceUrn string

// SponID Contains an identity of a sponsor.
type SponID string

// SponsoringStatus sponsoring status
type SponsoringStatus string

const (

	// SponsoringStatusSPONSORDISABLED captures enum value "SPONSOR_DISABLED"
	SponsoringStatusSPONSORDISABLED SponsoringStatus = "SPONSOR_DISABLED"

	// SponsoringStatusSPONSORENABLED captures enum value "SPONSOR_ENABLED"
	SponsoringStatusSPONSORENABLED SponsoringStatus = "SPONSOR_ENABLED"
)

// Supi supi
type Supi string

// IPV4Addr Ipv4 addr
type IPV4Addr string

// IPV6Addr Ipv6 addr
type IPV6Addr string

// MacAddr48 mac addr48
type MacAddr48 string

// AppSessionContextReqData Identifies the service requirements of an
// Individual Application Session Context.
type AppSessionContextReqData struct {

	// af app Id
	AfAppID AfAppID `json:"afAppId,omitempty"`

	// af charging identifier
	AfChargingIdentifier ChargingID `json:"afChargingIdentifier,omitempty"`

	// af req data
	AfReqData AfRequestedData `json:"afReqData,omitempty"`

	// af rout req
	AfRoutReq *AfRoutingRequirement `json:"afRoutReq,omitempty"`

	// asp Id
	AspID AspID `json:"aspId,omitempty"`

	// bdt ref Id
	BdtRefID BdtReferenceID `json:"bdtRefId,omitempty"`

	// dnn
	Dnn Dnn `json:"dnn,omitempty"`

	// ev subsc
	EvSubsc *EventsSubscReqData `json:"evSubsc,omitempty"`

	// gpsi
	Gpsi Gpsi `json:"gpsi,omitempty"`

	// ip domain
	IPDomain string `json:"ipDomain,omitempty"`

	// indication of MCVideo service request
	McVideoID string `json:"mcVideoId,omitempty"`

	// indication of MCPTT service request
	McpttID string `json:"mcpttId,omitempty"`

	// med components
	MedComponents map[string]MediaComponent `json:"medComponents,omitempty"`

	// indication of MPS service request
	MpsID string `json:"mpsId,omitempty"`

	// notif Uri
	// Required: true
	NotifURI URI `json:"notifUri"`

	// preempt controlinfo
	PreemptControlinfo PreemptionControlInformation `json:"preemptControlinfo,omitempty"`

	// res prio
	ResPrio ReservPriority `json:"resPrio,omitempty"`

	// serv inf status
	ServInfStatus ServiceInfoStatus `json:"servInfStatus,omitempty"`

	// serv urn
	ServUrn ServiceUrn `json:"servUrn,omitempty"`

	// slice info
	SliceInfo Snssai `json:"sliceInfo,omitempty"`

	// spon Id
	SponID SponID `json:"sponId,omitempty"`

	// spon status
	SponStatus SponsoringStatus `json:"sponStatus,omitempty"`

	// supi
	Supi Supi `json:"supi,omitempty"`

	// supp feat
	// Required: true
	SuppFeat SupportedFeatures `json:"suppFeat"`

	// ue Ipv4
	UeIPV4 IPV4Addr `json:"ueIpv4,omitempty"`

	// ue Ipv6
	UeIPV6 IPV6Addr `json:"ueIpv6,omitempty"`

	// ue mac
	UeMac MacAddr48 `json:"ueMac,omitempty"`
}

// ServAuthInfo serv auth info
type ServAuthInfo string

// UeIdentityInfo Represents 5GS-Level UE identities.
type UeIdentityInfo struct {

	// gpsi
	// Pattern: ^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$
	Gpsi string `json:"gpsi,omitempty"`

	// pei
	// Pattern: ^(imei-[0-9]{15}|imeisv-[0-9]{16}|.+)$
	Pei string `json:"pei,omitempty"`

	// supi
	// Pattern: ^(imsi-[0-9]{5,15}|nai-.+|.+)$
	Supi string `json:"supi,omitempty"`
}

// AppSessionContextRespData Describes the authorization data of an Individual
// Application Session Context created by the PCF.
type AppSessionContextRespData struct {

	// serv auth info
	ServAuthInfo ServAuthInfo `json:"servAuthInfo,omitempty"`

	// supp feat
	SuppFeat SupportedFeatures `json:"suppFeat,omitempty"`

	// ue ids
	// Min Items: 1
	UeIds []*UeIdentityInfo `json:"ueIds"`
}

// AccessType access type
type AccessType string

const (

	// AccessTypeNr3GPPACCESS captures enum value "3GPP_ACCESS"
	AccessTypeNr3GPPACCESS AccessType = "3GPP_ACCESS"

	// AccessTypeNON3GPPACCESS captures enum value "NON_3GPP_ACCESS"
	AccessTypeNON3GPPACCESS AccessType = "NON_3GPP_ACCESS"
)

// Flows Identifies the flows
type Flows struct {

	// cont vers
	// Min Items: 1
	ContVers []ContentVersion `json:"contVers"`

	// f nums
	// Min Items: 1
	FNums []int64 `json:"fNums"`

	// med comp n
	// Required: true
	MedCompN *int64 `json:"medCompN"`
}

// AccessNetChargingIdentifier Describes the access network charging identifier.
type AccessNetChargingIdentifier struct {

	// acc net cha Id value
	// Required: true
	AccNetChaIDValue *int32 `json:"accNetChaIdValue"`

	// flows
	// Min Items: 1
	Flows []*Flows `json:"flows"`
}

// AccNetChargingAddress Describes the network entity within the access network
// performing charging
type AccNetChargingAddress struct {

	// an charg Ipv4 addr
	AnChargIPV4Addr IPV4Addr `json:"anChargIpv4Addr,omitempty"`

	// an charg Ipv6 addr
	AnChargIPV6Addr IPV6Addr `json:"anChargIpv6Addr,omitempty"`
}

// AnGwAddress describes the address of the access network gateway control node
type AnGwAddress struct {

	// an gw Ipv4 addr
	AnGwIPV4Addr string `json:"anGwIpv4Addr,omitempty"`

	// an gw Ipv6 addr
	AnGwIPV6Addr string `json:"anGwIpv6Addr,omitempty"`
}

// AfEventNotification describes the event information delivered in the
// notification
type AfEventNotification struct {

	// event
	// Required: true
	Event AfEvent `json:"event"`

	// flows
	// Min Items: 1
	Flows []*Flows `json:"flows"`
}

// MediaComponentResourcesStatus media component resources status
type MediaComponentResourcesStatus string

const (

	// MediaComponentResourcesStatusACTIVE captures enum value "ACTIVE"
	MediaComponentResourcesStatusACTIVE MediaComponentResourcesStatus = "ACTIVE"

	// MediaComponentResourcesStatusINACTIVE captures enum value "INACTIVE"
	MediaComponentResourcesStatusINACTIVE MediaComponentResourcesStatus = "INACTIVE"
)

// ResourcesAllocationInfo describes the status of the PCC rule(s) related to
// certain media components.
type ResourcesAllocationInfo struct {

	// flows
	// Min Items: 1
	Flows []*Flows `json:"flows"`

	// mc resourc status
	// Required: true
	McResourcStatus MediaComponentResourcesStatus `json:"mcResourcStatus"`
}

// FinalUnitAction final unit action
type FinalUnitAction string

const (

	// FinalUnitActionTERMINATE captures enum value "TERMINATE"
	FinalUnitActionTERMINATE FinalUnitAction = "TERMINATE"

	// FinalUnitActionREDIRECT captures enum value "REDIRECT"
	FinalUnitActionREDIRECT FinalUnitAction = "REDIRECT"

	// FinalUnitActionRESTRICTACCESS captures enum value "RESTRICT_ACCESS"
	FinalUnitActionRESTRICTACCESS FinalUnitAction = "RESTRICT_ACCESS"
)

// OutOfCreditInformation Indicates the SDFs without available credit and the
// corresponding termination action.
type OutOfCreditInformation struct {

	// fin unit act
	// Required: true
	FinUnitAct FinalUnitAction `json:"finUnitAct"`

	// flows
	// Min Items: 1
	Flows []*Flows `json:"flows"`
}

// QosNotifType qos notif type
type QosNotifType string

const (

	// QosNotifTypeGUARANTEED captures enum value "GUARANTEED"
	QosNotifTypeGUARANTEED QosNotifType = "GUARANTEED"

	// QosNotifTypeNOTGUARANTEED captures enum value "NOT_GUARANTEED"
	QosNotifTypeNOTGUARANTEED QosNotifType = "NOT_GUARANTEED"
)

// QosNotificationControlInfo Indicates whether the QoS targets for a GRB flow
// are not  guaranteed or guaranteed again
type QosNotificationControlInfo struct {

	// flows
	// Min Items: 1
	Flows []*Flows `json:"flows"`

	// notif type
	// Required: true
	NotifType QosNotifType `json:"notifType"`
}

// RatType rat type
type RatType string

const (

	// RatTypeNR captures enum value "NR"
	RatTypeNR RatType = "NR"

	// RatTypeEUTRA captures enum value "EUTRA"
	RatTypeEUTRA RatType = "EUTRA"

	// RatTypeWLAN captures enum value "WLAN"
	RatTypeWLAN RatType = "WLAN"

	// RatTypeVIRTUAL captures enum value "VIRTUAL"
	RatTypeVIRTUAL RatType = "VIRTUAL"
)

// DateTime date time
type DateTime strfmt.DateTime

// EutraLocation eutra location
type EutraLocation struct {

	// age of location information
	// Maximum: 32767
	// Minimum: 0
	AgeOfLocationInformation int64 `json:"ageOfLocationInformation,omitempty"`

	// ecgi
	// Required: true
	Ecgi Ecgi `json:"ecgi"`

	// geodetic information
	// Pattern: ^[0-9A-F]{20}$
	GeodeticInformation string `json:"geodeticInformation,omitempty"`

	// geographical information
	// Pattern: ^[0-9A-F]{16}$
	GeographicalInformation string `json:"geographicalInformation,omitempty"`

	// global ngenb Id
	GlobalNgenbID *GlobalRanNodeID `json:"globalNgenbId,omitempty"`

	// tai
	// Required: true
	Tai Tai `json:"tai"`

	// ue location timestamp
	// Format: date-time
	UeLocationTimestamp DateTime `json:"ueLocationTimestamp,omitempty"`
}

// N3gaLocation n3ga location
type N3gaLocation struct {

	// n3 iwf Id
	// Pattern: ^[A-Fa-f0-9]+$
	N3IwfID string `json:"n3IwfId,omitempty"`

	// n3gpp tai
	N3gppTai *Tai `json:"n3gppTai,omitempty"`

	// port number
	PortNumber int64 `json:"portNumber,omitempty"`

	// ue Ipv4 addr
	UeIPV4Addr IPV4Addr `json:"ueIpv4Addr,omitempty"`

	// ue Ipv6 addr
	UeIPV6Addr IPV6Addr `json:"ueIpv6Addr,omitempty"`
}

// NrLocation nr location
type NrLocation struct {

	// age of location information
	// Maximum: 32767
	// Minimum: 0
	AgeOfLocationInformation *int64 `json:"ageOfLocationInformation,omitempty"`

	// geodetic information
	// Pattern: ^[0-9A-F]{20}$
	GeodeticInformation string `json:"geodeticInformation,omitempty"`

	// geographical information
	// Pattern: ^[0-9A-F]{16}$
	GeographicalInformation string `json:"geographicalInformation,omitempty"`

	// global gnb Id
	GlobalGnbID *GlobalRanNodeID `json:"globalGnbId,omitempty"`

	// ncgi
	// Required: true
	Ncgi *Ncgi `json:"ncgi"`

	// tai
	// Required: true
	Tai *Tai `json:"tai"`

	// ue location timestamp
	// Format: date-time
	UeLocationTimestamp DateTime `json:"ueLocationTimestamp,omitempty"`
}

// UserLocation user location
type UserLocation struct {

	// eutra location
	EutraLocation *EutraLocation `json:"eutraLocation,omitempty"`

	// n3ga location
	N3gaLocation *N3gaLocation `json:"n3gaLocation,omitempty"`

	// nr location
	NrLocation *NrLocation `json:"nrLocation,omitempty"`
}

// TimeZone time zone
type TimeZone string

// AccumulatedUsage accumulated usage
type AccumulatedUsage struct {

	// downlink volume
	DownlinkVolume Volume `json:"downlinkVolume,omitempty"`

	// duration
	Duration DurationSec `json:"duration,omitempty"`

	// total volume
	TotalVolume Volume `json:"totalVolume,omitempty"`

	// uplink volume
	UplinkVolume Volume `json:"uplinkVolume,omitempty"`
}

// EventsNotification describes the notification of a matched event
type EventsNotification struct {

	// access type
	AccessType AccessType `json:"accessType,omitempty"`

	// an charg addr
	AnChargAddr *AccNetChargingAddress `json:"anChargAddr,omitempty"`

	// an charg ids
	// Min Items: 1
	AnChargIds []*AccessNetChargingIdentifier `json:"anChargIds"`

	// an gw addr
	AnGwAddr *AnGwAddress `json:"anGwAddr,omitempty"`

	// ev notifs
	// Required: true
	// Min Items: 1
	EvNotifs []*AfEventNotification `json:"evNotifs"`

	// ev subs Uri
	// Required: true
	EvSubsURI *string `json:"evSubsUri"`

	// failed resourc alloc reports
	// Min Items: 1
	FailedResourcAllocReports []*ResourcesAllocationInfo `json:"failedResourcAllocReports"`

	// no net loc supp
	NoNetLocSupp bool `json:"noNetLocSupp,omitempty"`

	// out of cred reports
	// Min Items: 1
	OutOfCredReports []*OutOfCreditInformation `json:"outOfCredReports"`

	// plmn Id
	PlmnID *PlmnID `json:"plmnId,omitempty"`

	// qnc reports
	// Min Items: 1
	QncReports []*QosNotificationControlInfo `json:"qncReports"`

	// rat type
	RatType RatType `json:"ratType,omitempty"`

	// ue loc
	UeLoc *UserLocation `json:"ueLoc,omitempty"`

	// ue time zone
	UeTimeZone TimeZone `json:"ueTimeZone,omitempty"`

	// usg rep
	UsgRep *AccumulatedUsage `json:"usgRep,omitempty"`
}

// SpatialValidityRm this data type is defined in the same way as the
// SpatialValidity data type, but with the OpenAPI nullable property set to
// true
type SpatialValidityRm struct {

	// presence info list
	// Required: true
	PresenceInfoList map[string]PresenceInfo `json:"presenceInfoList"`
}

// AfRoutingRequirementRm this data type is defined in the same way as the
// AfRoutingRequirement data type, but with the OpenAPI nullable property set to true and the spVal and tempVals attributes defined as removable.
type AfRoutingRequirementRm struct {

	// addr preser ind
	AddrPreserInd bool `json:"addrPreserInd,omitempty"`

	// app reloc
	AppReloc bool `json:"appReloc,omitempty"`

	// route to locs
	// Min Items: 1
	RouteToLocs []*RouteToLocation `json:"routeToLocs"`

	// sp val
	SpVal *SpatialValidityRm `json:"spVal,omitempty"`

	// temp vals
	// Min Items: 1
	TempVals []*TemporalValidity `json:"tempVals"`

	// up path chg sub
	UpPathChgSub *UpPathChgEvent `json:"upPathChgSub,omitempty"`
}

// VolumeRm Unsigned integer identifying a volume in units of bytes with
// "nullable=true" property.
type VolumeRm int64

// DurationSecRm Unsigned integer identifying a period of time in units of
// seconds with "nullable=true" property.
type DurationSecRm int64

// UsageThresholdRm usage threshold rm
type UsageThresholdRm struct {

	// downlink volume
	DownlinkVolume VolumeRm `json:"downlinkVolume,omitempty"`

	// duration
	Duration DurationSecRm `json:"duration,omitempty"`

	// total volume
	TotalVolume VolumeRm `json:"totalVolume,omitempty"`

	// uplink volume
	UplinkVolume VolumeRm `json:"uplinkVolume,omitempty"`
}

// EventsSubscReqDataRm this data type is defined in the same way as the
// EventsSubscReqData data type, but with the OpenAPI nullable property set to
// true.
type EventsSubscReqDataRm struct {

	// events
	// Required: true
	Events []*AfEventSubscription `json:"events"`

	// notif Uri
	NotifURI URI `json:"notifUri,omitempty"`

	// req ani
	ReqAni RequiredAccessInfo `json:"reqAni,omitempty"`

	// usg thres
	UsgThres *UsageThresholdRm `json:"usgThres,omitempty"`
}

// BitRateRm bit rate rm
type BitRateRm string

// TosTrafficClassRm this data type is defined in the same way as the
// TosTrafficClass data type, but with the OpenAPI nullable property set to
// true
type TosTrafficClassRm string

// MediaSubComponentRm This data type is defined in the same way as the
// MediaSubComponent data type, but with the OpenAPI nullable property set to
// true. Removable attributes marBwDland marBwUl are defined with the
// corresponding removable data type.
type MediaSubComponentRm struct {

	// af sig protocol
	AfSigProtocol AfSigProtocol `json:"afSigProtocol,omitempty"`

	// ethf descs
	// Max Items: 2
	// Min Items: 1
	EthfDescs []*EthFlowDescription `json:"ethfDescs"`

	// f descs
	// Max Items: 2
	// Min Items: 1
	FDescs []FlowDescription `json:"fDescs"`

	// f num
	// Required: true
	FNum *int64 `json:"fNum"`

	// f status
	FStatus FlowStatus `json:"fStatus,omitempty"`

	// flow usage
	FlowUsage FlowUsage `json:"flowUsage,omitempty"`

	// mar bw dl
	MarBwDl BitRateRm `json:"marBwDl,omitempty"`

	// mar bw ul
	MarBwUl BitRateRm `json:"marBwUl,omitempty"`

	// tos tr cl
	TosTrCl TosTrafficClassRm `json:"tosTrCl,omitempty"`
}

// PreemptionCapabilityRm preemption capability rm
type PreemptionCapabilityRm string

const (

	// PreemptionCapabilityRmNOTPREEMPT captures enum value "NOT_PREEMPT"
	PreemptionCapabilityRmNOTPREEMPT PreemptionCapabilityRm = "NOT_PREEMPT"

	// PreemptionCapabilityRmMAYPREEMPT captures enum value "MAY_PREEMPT"
	PreemptionCapabilityRmMAYPREEMPT PreemptionCapabilityRm = "MAY_PREEMPT"
)

// PreemptionVulnerabilityRm preemption vulnerability rm
type PreemptionVulnerabilityRm string

const (

	// PreemptionVulnerabilityRmNOTPREEMPTABLE captures enum value "NOT_PREEMPTABLE"
	PreemptionVulnerabilityRmNOTPREEMPTABLE PreemptionVulnerabilityRm = "NOT_PREEMPTABLE"

	// PreemptionVulnerabilityRmPREEMPTABLE captures enum value "PREEMPTABLE"
	PreemptionVulnerabilityRmPREEMPTABLE PreemptionVulnerabilityRm = "PREEMPTABLE"
)

// MediaComponentRm This data type is defined in the same way as the
// MediaComponent data type, but with the OpenAPI nullable property set to true
type MediaComponentRm struct {

	// af app Id
	AfAppID AfAppID `json:"afAppId,omitempty"`

	// af rout req
	AfRoutReq *AfRoutingRequirementRm `json:"afRoutReq,omitempty"`

	// codecs
	// Max Items: 2
	// Min Items: 1
	Codecs []CodecData `json:"codecs"`

	// cont ver
	ContVer ContentVersion `json:"contVer,omitempty"`

	// f status
	FStatus FlowStatus `json:"fStatus,omitempty"`

	// mar bw dl
	MarBwDl BitRateRm `json:"marBwDl,omitempty"`

	// mar bw ul
	MarBwUl BitRateRm `json:"marBwUl,omitempty"`

	// med comp n
	// Required: true
	MedCompN *int64 `json:"medCompN"`

	// med sub comps
	MedSubComps map[string]MediaSubComponentRm `json:"medSubComps,omitempty"`

	// med type
	MedType MediaType `json:"medType,omitempty"`

	// mir bw dl
	MirBwDl BitRateRm `json:"mirBwDl,omitempty"`

	// mir bw ul
	MirBwUl BitRateRm `json:"mirBwUl,omitempty"`

	// preempt cap
	PreemptCap PreemptionCapabilityRm `json:"preemptCap,omitempty"`

	// preempt vuln
	PreemptVuln PreemptionVulnerabilityRm `json:"preemptVuln,omitempty"`

	// prio sharing ind
	PrioSharingInd PrioritySharingIndicator `json:"prioSharingInd,omitempty"`

	// res prio
	ResPrio ReservPriority `json:"resPrio,omitempty"`

	// rr bw
	RrBw BitRateRm `json:"rrBw,omitempty"`

	// rs bw
	RsBw BitRateRm `json:"rsBw,omitempty"`

	// sharing key dl
	SharingKeyDl Uint32Rm `json:"sharingKeyDl,omitempty"`

	// sharing key ul
	SharingKeyUl Uint32Rm `json:"sharingKeyUl,omitempty"`
}

// PreemptionControlInformationRm preemption control information rm
type PreemptionControlInformationRm string

const (

	// PreemptionControlInformationRmMOSTRECENT captures enum value "MOST_RECENT"
	PreemptionControlInformationRmMOSTRECENT PreemptionControlInformationRm = "MOST_RECENT"

	// PreemptionControlInformationRmLEASTRECENT captures enum value "LEAST_RECENT"
	PreemptionControlInformationRmLEASTRECENT PreemptionControlInformationRm = "LEAST_RECENT"

	// PreemptionControlInformationRmHIGHESTBW captures enum value "HIGHEST_BW"
	PreemptionControlInformationRmHIGHESTBW PreemptionControlInformationRm = "HIGHEST_BW"
)

// SipForkingIndication sip forking indication
type SipForkingIndication string

const (

	// SipForkingIndicationSINGLEDIALOGUE captures enum value "SINGLE_DIALOGUE"
	SipForkingIndicationSINGLEDIALOGUE SipForkingIndication = "SINGLE_DIALOGUE"

	// SipForkingIndicationSEVERALDIALOGUES captures enum value "SEVERAL_DIALOGUES"
	SipForkingIndicationSEVERALDIALOGUES SipForkingIndication = "SEVERAL_DIALOGUES"
)

// AppSessionContextUpdateData Identifies the modifications to an Individual
// Application Session Context and may include the modifications to the
// sub-resource Events Subscription.
type AppSessionContextUpdateData struct {

	// af app Id
	AfAppID AfAppID `json:"afAppId,omitempty"`

	// af rout req
	AfRoutReq *AfRoutingRequirementRm `json:"afRoutReq,omitempty"`

	// asp Id
	AspID AspID `json:"aspId,omitempty"`

	// bdt ref Id
	BdtRefID BdtReferenceID `json:"bdtRefId,omitempty"`

	// ev subsc
	EvSubsc *EventsSubscReqDataRm `json:"evSubsc,omitempty"`

	// indication of modification of MCVideo service
	McVideoID string `json:"mcVideoId,omitempty"`

	// indication of MCPTT service request
	McpttID string `json:"mcpttId,omitempty"`

	// med components
	MedComponents map[string]MediaComponentRm `json:"medComponents,omitempty"`

	// indication of MPS service request
	MpsID string `json:"mpsId,omitempty"`

	// preempt controlinfo
	PreemptControlinfo PreemptionControlInformationRm `json:"preemptControlinfo,omitempty"`

	// res prio
	ResPrio ReservPriority `json:"resPrio,omitempty"`

	// serv inf status
	ServInfStatus ServiceInfoStatus `json:"servInfStatus,omitempty"`

	// sip fork ind
	SipForkInd SipForkingIndication `json:"sipForkInd,omitempty"`

	// spon Id
	SponID SponID `json:"sponId,omitempty"`

	// spon status
	SponStatus SponsoringStatus `json:"sponStatus,omitempty"`
}

// AcceptableServiceInfo Indicates the maximum bandwidth that shall be
// authorized by the PCF.
type AcceptableServiceInfo struct {

	// acc bw med comps
	AccBwMedComps map[string]MediaComponent `json:"accBwMedComps,omitempty"`

	// mar bw dl
	// Pattern: ^\d+(\.\d+)? (bps|Kbps|Mbps|Gbps|Tbps)$
	MarBwDl string `json:"marBwDl,omitempty"`

	// mar bw ul
	// Pattern: ^\d+(\.\d+)? (bps|Kbps|Mbps|Gbps|Tbps)$
	MarBwUl string `json:"marBwUl,omitempty"`
}

// ExtendedProblemDetails Extends ProblemDetails to also include the acceptable
// service info.
type ExtendedProblemDetails struct {
	ProblemDetails

	// acceptable serv info
	AcceptableServInfo *AcceptableServiceInfo `json:"acceptableServInfo,omitempty"`
}
