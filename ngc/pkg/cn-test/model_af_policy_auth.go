// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package ngccntest

// AccessType the model 'AccessType'
type AccessType string

// List of AccessType
const (
	_3_GPP_ACCESS    AccessType = "3GPP_ACCESS"
	NON_3_GPP_ACCESS AccessType = "NON_3GPP_ACCESS"
)

// MacAddr type
type MacAddr string

type AfEvent string

const (
	AccessTypeChange          AfEvent = "ACCESS_TYPE_CHANGE"
	FailedResourcesAllocation AfEvent = "FAILED_RESOURCES_ALLOCATION"
	PlmnChg                   AfEvent = "PLMN_CHG"
	QosNotif                  AfEvent = "QOS_NOTIF"
	ResourceAllocated         AfEvent = "SUCCESSFUL_RESOURCES_ALLOCATION"
	UsageReport               AfEvent = "USAGE_REPORT"
)

type AfNotifMethod string

const (
	EventDetection AfNotifMethod = "EVENT_DETECTION"
	OneTime        AfNotifMethod = "ONE_TIME"
)

/*
 *AnGwAddress describes the address of the access network gateway control node
 * It can be an Ipv4 or Ipv6 Address
 */
type AnGwAddress string

/*
 * FlowDirection Possible values are -
 * -DOWNLINK: The corresponding filter applies for traffic to the UE.
 * -UPLINK: The corresponding filter applies for traffic from the UE.
 * -BIDIRECTIONAL: The corresponding filter applies for traffic both to and from
 *    the UE.
 * -UNSPECIFIED: The corresponding filter applies for traffic to the UE
 *   (downlink), but has no specific direction declared. The service data flow
 *   detection shall apply the filter for uplink traffic as if the filter was
 *   bidirectional.
 * The PCF shall not use the value UNSPECIFIED in filters created by the network
 * in NW-initiated procedures. The PCF shall only include the value UNSPECIFIED
 * in filters in UE-initiated procedures if the same value is received from the
 * SMF.
 */
type FlowDirection string

const (
	DLFlow           FlowDirection = "DOWNLINK"
	ULFlow           FlowDirection = "UPLINK"
	BiDirectionFlow  FlowDirection = "BIDIRECTIONAL"
	FlowNotSpecified FlowDirection = "UNSPECIFIED"
)

type FlowStatus string

const (
	ULFlowEnabled FlowStatus = "ENABLED-UPLINK"
	DLFlowEnabled FlowStatus = "ENABLED-UPLINK"
	FlowEnabled   FlowStatus = "ENABLED"
	FlowDisabled  FlowStatus = "DISABLED"
	FlowRemoved   FlowStatus = "REMOVED"
)

type FlowUsage string

const (
	FlowUsageNotSpecified FlowUsage = "NO_INFO"
	RTCPFlow              FlowUsage = "RTCP"
)

type MediaComponentResourcesStatus string

const (
	MediaComponentResourceActive   MediaComponentResourcesStatus = "ACTIVE"
	MediaComponentResourceInActive MediaComponentResourcesStatus = "INACTIVE"
)

type MediaType string

const (
	MediaTypeAudio       MediaType = "AUDIO"
	MediaTypeVideo       MediaType = "VIDEO"
	MediaTypeData        MediaType = "DATA"
	MediaTypeApplication MediaType = "APPLICATION"
	MediaTypeControl     MediaType = "CONTROL"
	MediaTypeText        MediaType = "TEXT"
	MediaTypeMessage     MediaType = "MESSAGE"
	MediaTypeMisc        MediaType = "OTHER"
)

type PresenceState string

const (
	PresenceStateInArea    PresenceState = "IN_AREA"
	PresenceStateOutOfArea PresenceState = "OUT_OF_AREA"
	PresenceStateUnknown   PresenceState = "UNKNOWN"
	PresenceStateInactive  PresenceState = "INACTIVE"
)

type QosNotifType string

const (
	QosNotifGuaranteed    QosNotifType = "GUARANTEED"
	QosNotifNotGuaranteed QosNotifType = "NOT_GUARANTEED"
)

type RatType string

const (
	RatTypeNR      RatType = "NR"
	RatTypeEUTRA   RatType = "EUTRA"
	RatTypeWLAN    RatType = "WLAN"
	RatTypeVIRTUAL RatType = "VIRTUAL"
)

type ReservPriority string

const (
	ReservPriority1  ReservPriority = "PRIO_1"
	ReservPriority2  ReservPriority = "PRIO_2"
	ReservPriority3  ReservPriority = "PRIO_3"
	ReservPriority4  ReservPriority = "PRIO_4"
	ReservPriority5  ReservPriority = "PRIO_5"
	ReservPriority6  ReservPriority = "PRIO_6"
	ReservPriority7  ReservPriority = "PRIO_7"
	ReservPriority8  ReservPriority = "PRIO_8"
	ReservPriority9  ReservPriority = "PRIO_9"
	ReservPriority10 ReservPriority = "PRIO_10"
	ReservPriority11 ReservPriority = "PRIO_11"
	ReservPriority12 ReservPriority = "PRIO_12"
	ReservPriority13 ReservPriority = "PRIO_13"
	ReservPriority14 ReservPriority = "PRIO_14"
	ReservPriority15 ReservPriority = "PRIO_15"
	ReservPriority16 ReservPriority = "PRIO_16"
)

type ServAuthInfo string

const (
	ServAuthNotKnown   ServAuthInfo = "TP_NOT_KNOWN"
	ServAuthExpired    ServAuthInfo = "TP_EXPIRED"
	ServAuthNotOcurred ServAuthInfo = "TP_NOT_YET_OCURRED"
)

type SponsoringStatus string

const (
	SponsorEnabled  SponsoringStatus = "SPONSOR_EnABLED"
	SponsorDisabled SponsoringStatus = "SPONSOR_DISABLED"
)

type TerminationCause string

const (
	AllSDFDeactivated    TerminationCause = "ALL_SDF_DEACTIVATION"
	PDUSessionTerminated TerminationCause = "PDU_SESSION_TERMINATION"
)

type AccumulatedUsage struct {
	// Unsigned integer identifying a period of time in units of seconds.
	Duration int32 `json:"duration,omitempty"`
	// Unsigned integer identifying a volume in units of bytes.
	TotalVolume int64 `json:"totalVolume,omitempty"`
	// Unsigned integer identifying a volume in units of bytes.
	DownlinkVolume int64 `json:"downlinkVolume,omitempty"`
	// Unsigned integer identifying a volume in units of bytes.
	UplinkVolume int64 `json:"uplinkVolume,omitempty"`
}

// AfEventNotification describes the event info delivered in the notification
type AfEventNotification struct {
	Event AfEvent `json:"event"`
	Flows []Flows `json:"flows,omitempty"`
}

// AfEventSubscription describes the event info delivered in the subscription
type AfEventSubscription struct {
	Event       AfEvent       `json:"event"`
	NotifMethod AfNotifMethod `json:"notifMethod,omitempty"`
}

// AfRoutingRequirement describes the event info delivered in the subscription
type AfRoutingRequirement struct {
	AppReloc     bool               `json:"appReloc,omitempty"`
	RouteToLocs  []RouteToLocation  `json:"routeToLocs,omitempty"`
	SpVal        *SpatialValidity   `json:"spVal,omitempty"`
	TempVals     []TemporalValidity `json:"tempVals,omitempty"`
	UpPathChgSub *UpPathChgEvent    `json:"upPathChgSub,omitempty"`
}

/*
 * AppSessionContextReqData Identifies the service requirements of an Individual
 * Application Session Context.
 */
type AppSessionContextReqData struct {
	AfAppId   string                `json:"afAppId,omitempty"`
	AfRoutReq *AfRoutingRequirement `json:"afRoutReq,omitempty"`
	// Contains an identity of an application service provider.
	AspId string `json:"aspId,omitempty"`
	/*
	 * string identifying a BDT Reference ID as defined in subclause
	 * 5.3.3 of 3GPP TS 29.154.
	 */
	BdtRefId      string                    `json:"bdtRefId,omitempty"`
	Dnn           string                    `json:"dnn,omitempty"`
	EvSubsc       *EventsSubscReqData       `json:"evSubsc,omitempty"`
	MedComponents map[string]MediaComponent `json:"medComponents,omitempty"`
	IpDomain      string                    `json:"ipDomain,omitempty"`
	// indication of MPS service request
	MpsId      string           `json:"mpsId,omitempty"`
	NotifUri   string           `json:"notifUri"`
	SliceInfo  *Snssai          `json:"sliceInfo,omitempty"`
	SponId     string           `json:"sponId,omitempty"`
	SponStatus SponsoringStatus `json:"sponStatus,omitempty"`
	Supi       string           `json:"supi,omitempty"`
	Gpsi       string           `json:"gpsi,omitempty"`
	/*
	 * A string used to indicate the features supported by an API that is
	 * used as defined in subclause 6.6 in 3GPP TS 29.500 [1]. The string
	 * shall contain a bitmask indicating supported features in hexadecimal
	 * representation. Each character in the string shall take a value of
	 * \"0\" to \"9\" or \"A\" to \"F\" and shall represent the support of 4
	 * features as described in table 5.2.2-3. The most significant
	 * character representing the highest-numbered features shall appear
	 * first in the string, and the character representing features 1 to 4
	 * shall appear last in the string. The list of features and their
	 * numbering (starting with 1) are defined separately for each API.
	 * Possible features for traffic influencing are
	 * Notification_websocket(1), Notification_test_event(2)
	 */
	SuppFeat string   `json:"suppFeat"`
	UeIpv4   IPv4Addr `json:"ueIpv4,omitempty"`
	UeIpv6   IPv6Addr `json:"ueIpv6,omitempty"`
	UeMac    MacAddr  `json:"ueMac,omitempty"`
}

/*
 * AppSessionContextRespData Describes the authorization data of an Individual
 * Application Session Context created by the PCF.
 */
type AppSessionContextRespData struct {
	ServAuthInfo ServAuthInfo `json:"servAuthInfo,omitempty"`
	/*
	 * A string used to indicate the features supported by an API that is
	 * used as defined in subclause 6.6 in 3GPP TS 29.500 [1]. The string
	 * shall contain a bitmask indicating supported features in hexadecimal
	 * representation. Each character in the string shall take a value of
	 * \"0\" to \"9\" or \"A\" to \"F\" and shall represent the support of 4
	 * features as described in table 5.2.2-3. The most significant
	 * character representing the highest-numbered features shall appear
	 * first in the string, and the character representing features 1 to 4
	 * shall appear last in the string. The list of features and their
	 * numbering (starting with 1) are defined separately for each API.
	 * Possible features for traffic influencing are
	 * Notification_websocket(1), Notification_test_event(2)
	 */
	SuppFeat string `json:"suppFeat,omitempty"`
}

/*
 * AppSessionContextUpdateData Identifies the modifications to an Individual
 * Application Session Context and may include the modifications to the
 * sub-resource Events Subscription.
 */
type AppSessionContextUpdateData struct {
	// Contains an AF application identifier.
	AfAppId   string                `json:"afAppId,omitempty"`
	AfRoutReq *AfRoutingRequirement `json:"afRoutReq,omitempty"`
	// Contains an identity of an application service provider.
	AspId string `json:"aspId,omitempty"`
	/*
	 * string identifying a BDT Reference ID as defined in subclause
	 * 5.3.3 of 3GPP TS 29.154.
	 */
	BdtRefId      string                    `json:"bdtRefId,omitempty"`
	EvSubsc       *EventsSubscReqData       `json:"evSubsc,omitempty"`
	MedComponents map[string]MediaComponent `json:"medComponents,omitempty"`
	// indication of MPS service request
	MpsId string `json:"mpsId,omitempty"`
	// Contains an identity of a sponsor.
	SponId     string           `json:"sponId,omitempty"`
	SponStatus SponsoringStatus `json:"sponStatus,omitempty"`
}

/*
 * AppSessionContext Represents an Individual Application Session Context
 * resource.
 */
type AppSessionContext struct {
	AscReqData  *AppSessionContextReqData  `json:"ascReqData,omitempty"`
	AscRespData *AppSessionContextRespData `json:"ascRespData,omitempty"`
	EvsNotif    *EventsNotification        `json:"evsNotif,omitempty"`
}

type Ecgi struct {
	PlmnId      PlmnId `json:"plmnId"`
	EutraCellId string `json:"eutraCellId"`
}

// EventsNotification describes the notification of a matched event
type EventsNotification struct {
	AccessType AccessType  `json:"accessType,omitempty"`
	AnGwAddr   AnGwAddress `json:"anGwAddr,omitempty"`
	// string providing an URI formatted according to IETF RFC 3986.
	EvSubsUri                 string                       `json:"evSubsUri"`
	EvNotifs                  []AfEventNotification        `json:"evNotifs"`
	FailedResourcAllocReports []ResourcesAllocationInfo    `json:"failedResourcAllocReports,omitempty"`
	PlmnId                    *PlmnId                      `json:"plmnId,omitempty"`
	QncReports                []QosNotificationControlInfo `json:"qncReports,omitempty"`
	RatType                   RatType                      `json:"ratType,omitempty"`
	UsgRep                    AccumulatedUsage             `json:"usgRep,omitempty"`
}

// EventsSubscReqData Identifies the events the application subscribes to.
type EventsSubscReqData struct {
	Events []AfEventSubscription `json:"events"`
	// string providing an URI formatted according to IETF RFC 3986.
	NotifUri string          `json:"notifUri,omitempty"`
	UsgThres *UsageThreshold `json:"usgThres,omitempty"`
}

// Flows Identifies the flows
type Flows struct {
	ContVers []int32 `json:"contVers,omitempty"`
	FNums    []int32 `json:"fNums,omitempty"`
	MedCompN int32   `json:"medCompN"`
}

// GlobalRanNodeId struct for GlobalRanNodeId
type GlobalRanNodeId struct {
	PlmnId  *PlmnId `json:"plmnId"`
	N3IwfId string  `json:"n3IwfId,omitempty"`
	GNbId   *GNbId  `json:"gNbId,omitempty"`
	NgeNbId string  `json:"ngeNbId,omitempty"`
}

type GNbId struct {
	BitLength int32  `json:"bitLength"`
	GNBValue  string `json:"gNBValue"`
}

// MediaComponent Identifies a media component.
type MediaComponent struct {
	// Contains an AF application identifier.
	AfAppId   string                `json:"afAppId,omitempty"`
	AfRoutReq *AfRoutingRequirement `json:"afRoutReq,omitempty"`
	// Represents the content version of some content.
	ContVer     int32                        `json:"contVer,omitempty"`
	Codecs      []string                     `json:"codecs,omitempty"`
	FStatus     FlowStatus                   `json:"fStatus,omitempty"`
	MarBwDl     string                       `json:"marBwDl,omitempty"`
	MarBwUl     string                       `json:"marBwUl,omitempty"`
	MedCompN    int32                        `json:"medCompN"`
	MedSubComps map[string]MediaSubComponent `json:"medSubComps,omitempty"`
	MedType     MediaType                    `json:"medType,omitempty"`
	MirBwDl     string                       `json:"mirBwDl,omitempty"`
	MirBwUl     string                       `json:"mirBwUl,omitempty"`
	ResPrio     ReservPriority               `json:"resPrio,omitempty"`
}

// MediaSubComponent Identifies a media subcomponent
type MediaSubComponent struct {
	EthfDescs []EthFlowDescription `json:"ethfDescs,omitempty"`
	FNum      int32                `json:"fNum"`
	FDescs    []string             `json:"fDescs,omitempty"`
	FStatus   FlowStatus           `json:"fStatus,omitempty"`
	MarBwDl   string               `json:"marBwDl,omitempty"`
	MarBwUl   string               `json:"marBwUl,omitempty"`
	/*
	 * 2-octet string, where each octet is encoded in hexadecimal
	 * representation. The first octet contains the IPv4 Type-of-Service or
	 * the IPv6 Traffic-Class field and the second octet contains the
	 * ToS/Traffic Class mask field.
	 */
	TosTrCl   string    `json:"tosTrCl,omitempty"`
	FlowUsage FlowUsage `json:"flowUsage,omitempty"`
}

type Ncgi struct {
	PlmnId   PlmnId `json:"plmnId"`
	NrCellId string `json:"nrCellId"`
}

type PlmnId struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"mnc"`
}

type PresenceInfo struct {
	PraId               string            `json:"praId,omitempty"`
	PresenceState       PresenceState     `json:"presenceState,omitempty"`
	TrackingAreaList    []Tai             `json:"trackingAreaList,omitempty"`
	EcgiList            []Ecgi            `json:"ecgiList,omitempty"`
	NcgiList            []Ncgi            `json:"ncgiList,omitempty"`
	GlobalRanNodeIdList []GlobalRanNodeId `json:"globalRanNodeIdList,omitempty"`
}

/*
 * QosNotificationControlInfo Indicates whether the QoS targets for a GRB flow
 * are not  guaranteed or guaranteed again
 */
type QosNotificationControlInfo struct {
	NotifType QosNotifType `json:"notifType"`
	Flows     []Flows      `json:"flows,omitempty"`
}

/*
 * ResourcesAllocationInfo describes the status of the PCC rule(s) related to
 * certain media components.
 */
type ResourcesAllocationInfo struct {
	McResourcStatus MediaComponentResourcesStatus `json:"mcResourcStatus"`
	Flows           []Flows                       `json:"flows,omitempty"`
}

// SpatialValidity describes explicitly the route to an Application location
type SpatialValidity struct {
	PresenceInfoList map[string]PresenceInfo `json:"presenceInfoList"`
}

type Tai struct {
	PlmnId PlmnId `json:"plmnId"`
	Tac    string `json:"tac"`
}

/*
 * TerminationInfo indicates the cause for requesting the deletion of the
 * Individual Application Session Context resource
 */
type TerminationInfo struct {
	TermCause TerminationCause `json:"termCause"`
	// string providing an URI formatted according to IETF RFC 3986.
	ResUri string `json:"resUri"`
}

type UpPathChgEvent struct {
	// string providing an URI formatted according to IETF RFC 3986.
	NotificationUri string `json:"notificationUri"`
	/*
	 * It is used to set the value of Notification Correlation ID in the
	 * notification sent by the SMF.
	 */
	NotifCorreId string         `json:"notifCorreId"`
	DnaiChgType  DNAIChangeType `json:"dnaiChgType"`
}

type UsageThreshold struct {
	// Unsigned integer identifying a period of time in units of seconds.
	Duration int32 `json:"duration,omitempty"`
	// Unsigned integer identifying a volume in units of bytes.
	TotalVolume int64 `json:"totalVolume,omitempty"`
	// Unsigned integer identifying a volume in units of bytes.
	DownlinkVolume int64 `json:"downlinkVolume,omitempty"`
	// Unsigned integer identifying a volume in units of bytes.
	UplinkVolume int64 `json:"uplinkVolume,omitempty"`
}
