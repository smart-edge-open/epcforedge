// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

// AccessType the model 'AccessType'
type AccessType string

// List of AccessType
const (
	AccessType3Gpp    AccessType = "3GPP_ACCESS"
	AccessTypeNon3Gpp AccessType = "NON_3GPP_ACCESS"
)

// Event type
type Event string

// list of Event
const (
	AccessTypeChange          Event = "ACCESS_TYPE_CHANGE"
	FailedResourcesAllocation Event = "FAILED_RESOURCES_ALLOCATION"
	PlmnChg                   Event = "PLMN_CHG"
	QosNotif                  Event = "QOS_NOTIF"
	ResourceAllocated         Event = "SUCCESSFUL_RESOURCES_ALLOCATION"
	UsageReport               Event = "USAGE_REPORT"
)

// NotifMethod type
type NotifMethod string

// List of NotifMethod
const (
	EventDetection NotifMethod = "EVENT_DETECTION"
	OneTime        NotifMethod = "ONE_TIME"
)

// AnGwAddress type
/*
 *AnGwAddress describes the address of the access network gateway control node
 * It can be an Ipv4 or Ipv6 Address
 */
type AnGwAddress string

// FlowDirection type
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

// List of FlowDirection
const (
	DLFlow           FlowDirection = "DOWNLINK"
	ULFlow           FlowDirection = "UPLINK"
	BiDirectionFlow  FlowDirection = "BIDIRECTIONAL"
	FlowNotSpecified FlowDirection = "UNSPECIFIED"
)

// FlowStatus type
type FlowStatus string

// list of FlowStatus
const (
	ULFlowEnabled FlowStatus = "ENABLED-UPLINK"
	DLFlowEnabled FlowStatus = "ENABLED-DOWNLINK"
	FlowEnabled   FlowStatus = "ENABLED"
	FlowDisabled  FlowStatus = "DISABLED"
	FlowRemoved   FlowStatus = "REMOVED"
)

// FlowUsage type
type FlowUsage string

// List of FlowUsage
const (
	FlowUsageNotSpecified FlowUsage = "NO_INFO"
	RTCPFlow              FlowUsage = "RTCP"
)

// MediaComponentResourcesStatus type
type MediaComponentResourcesStatus string

// list of MediaComponentResourceStatus
const (
	MediaComponentResourceActive   MediaComponentResourcesStatus = "ACTIVE"
	MediaComponentResourceInActive MediaComponentResourcesStatus = "INACTIVE"
)

// MediaType type
type MediaType string

// list of MediaType
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

// PresenceState type
type PresenceState string

// List of PresenceState
const (
	PresenceStateInArea    PresenceState = "IN_AREA"
	PresenceStateOutOfArea PresenceState = "OUT_OF_AREA"
	PresenceStateUnknown   PresenceState = "UNKNOWN"
	PresenceStateInactive  PresenceState = "INACTIVE"
)

// QosNotifType model
type QosNotifType string

// List of QosNotifType
const (
	QosNotifGuaranteed    QosNotifType = "GUARANTEED"
	QosNotifNotGuaranteed QosNotifType = "NOT_GUARANTEED"
)

// RatType type
type RatType string

// List of RatType
const (
	RatTypeNR      RatType = "NR"
	RatTypeEUTRA   RatType = "EUTRA"
	RatTypeWLAN    RatType = "WLAN"
	RatTypeVIRTUAL RatType = "VIRTUAL"
)

// ReservPriority type
type ReservPriority string

// list of ReservPriority
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

// ServAuthInfo type
type ServAuthInfo string

// List of ServAuthInfos
const (
	ServAuthNotKnown    ServAuthInfo = "TP_NOT_KNOWN"
	ServAuthExpired     ServAuthInfo = "TP_EXPIRED"
	ServAuthNotOccurred ServAuthInfo = "TP_NOT_YET_OCCURRED"
)

// SponsoringStatus type
type SponsoringStatus string

// List of SponsoringStatus
const (
	SponsorEnabled  SponsoringStatus = "SPONSOR_ENABLED"
	SponsorDisabled SponsoringStatus = "SPONSOR_DISABLED"
)

// TerminationCause type
type TerminationCause string

// List of TerminationCause
const (
	AllSDFDeactivated    TerminationCause = "ALL_SDF_DEACTIVATION"
	PDUSessionTerminated TerminationCause = "PDU_SESSION_TERMINATION"
)

// AccumulatedUsage struct
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

// PolicyEventNotification describes the event info delivered in the notification
type PolicyEventNotification struct {
	Event Event   `json:"event"`
	Flows []Flows `json:"flows,omitempty"`
}

// EventSubscription describes the event info delivered in the subscription
type EventSubscription struct {
	Event       Event       `json:"event"`
	NotifMethod NotifMethod `json:"notifMethod,omitempty"`
}

// RoutingRequirement describes the event info delivered in the subscription
type RoutingRequirement struct {
	AppReloc     bool               `json:"appReloc,omitempty"`
	RouteToLocs  []RouteToLocation  `json:"routeToLocs,omitempty"`
	SpVal        *SpatialValidity   `json:"spVal,omitempty"`
	TempVals     []TemporalValidity `json:"tempVals,omitempty"`
	UpPathChgSub *UpPathChgEvent    `json:"upPathChgSub,omitempty"`
}

// AppSessionContextReqData type
/*
 * AppSessionContextReqData Identifies the service requirements of an Individual
 * Application Session Context.
 */
type AppSessionContextReqData struct {
	AfAppID   string              `json:"afAppId,omitempty"`
	AfRoutReq *RoutingRequirement `json:"afRoutReq,omitempty"`
	// Contains an identity of an application service provider.
	AspID string `json:"aspId,omitempty"`
	/*
	 * string identifying a BDT Reference ID as defined in subclause
	 * 5.3.3 of 3GPP TS 29.154.
	 */
	BdtRefID      string                    `json:"bdtRefId,omitempty"`
	Dnn           string                    `json:"dnn,omitempty"`
	EvSubsc       *EventsSubscReqData       `json:"evSubsc,omitempty"`
	MedComponents map[string]MediaComponent `json:"medComponents,omitempty"`
	IPDomain      string                    `json:"ipDomain,omitempty"`
	// indication of MPS service request
	MpsID      string           `json:"mpsId,omitempty"`
	NotifURI   string           `json:"notifUri"`
	SliceInfo  *SNSSAI          `json:"sliceInfo,omitempty"`
	SponID     string           `json:"sponId,omitempty"`
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
	// This parameter is present only if websocket delivery of notification
	// is requested by consumer for events
	AfwebsockNotifConfig *AfwebsockNotifConfig `json:"afwebsockNotifConfig,omitempty"`
}

// AppSessionContextRespData struct
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

	// The Websocket Uri used for Notification delivery.
	// This is set by AF and is present in the response from AF to consumer
	WebsocketURI string `json:"websocketUri,omitempty"`
}

// AppSessionContextUpdateData struct
/*
 * AppSessionContextUpdateData Identifies the modifications to an Individual
 * Application Session Context and may include the modifications to the
 * sub-resource Events Subscription.
 */
type AppSessionContextUpdateData struct {
	// Contains an AF application identifier.
	AfAppID   string              `json:"afAppId,omitempty"`
	AfRoutReq *RoutingRequirement `json:"afRoutReq,omitempty"`
	// Contains an identity of an application service provider.
	AspID string `json:"aspId,omitempty"`
	/*
	 * string identifying a BDT Reference ID as defined in subclause
	 * 5.3.3 of 3GPP TS 29.154.
	 */
	BdtRefID      string                    `json:"bdtRefId,omitempty"`
	EvSubsc       *EventsSubscReqData       `json:"evSubsc,omitempty"`
	MedComponents map[string]MediaComponent `json:"medComponents,omitempty"`
	// indication of MPS service request
	MpsID string `json:"mpsId,omitempty"`
	// Contains an identity of a sponsor.
	SponID     string           `json:"sponId,omitempty"`
	SponStatus SponsoringStatus `json:"sponStatus,omitempty"`
	// This parameter is present only if websocket delivery of notification
	// is requested by consumer for events
	AfwebsockNotifConfig *AfwebsockNotifConfig `json:"afwebsockNotifConfig,omitempty"`
}

// AppSessionContext Represents an Individual Application Session
// Context resource
type AppSessionContext struct {
	AscReqData  *AppSessionContextReqData  `json:"ascReqData,omitempty"`
	AscRespData *AppSessionContextRespData `json:"ascRespData,omitempty"`
	EvsNotif    *EventsNotification        `json:"evsNotif,omitempty"`
}

// Ecgi Struct
type Ecgi struct {
	PlmnID      PlmnID `json:"plmnId"`
	EutraCellID string `json:"eutraCellId"`
}

// EventsNotification describes the notification of a matched event
type EventsNotification struct {
	AccessType AccessType  `json:"accessType,omitempty"`
	AnGwAddr   AnGwAddress `json:"anGwAddr,omitempty"`
	// string providing an URI formatted according to IETF RFC 3986.
	EvSubsURI                 string                       `json:"evSubsUri"`
	EvNotifs                  []PolicyEventNotification    `json:"evNotifs"`
	FailedResourcAllocReports []ResourcesAllocationInfo    `json:"failedResourcAllocReports,omitempty"`
	PlmnID                    *PlmnID                      `json:"plmnId,omitempty"`
	QncReports                []QosNotificationControlInfo `json:"qncReports,omitempty"`
	RatType                   RatType                      `json:"ratType,omitempty"`
	UsgRep                    AccumulatedUsage             `json:"usgRep,omitempty"`
}

// EventsSubscReqData Identifies the events the application subscribes to.
type EventsSubscReqData struct {
	Events []EventSubscription `json:"events"`
	// string providing an URI formatted according to IETF RFC 3986.
	NotifURI string          `json:"notifUri,omitempty"`
	UsgThres *UsageThreshold `json:"usgThres,omitempty"`
}

// Flows Identifies the flows
type Flows struct {
	ContVers []int32 `json:"contVers,omitempty"`
	FNums    []int32 `json:"fNums,omitempty"`
	MedCompN int32   `json:"medCompN"`
}

// GlobalRanNodeID struct for GlobalRanNodeId
type GlobalRanNodeID struct {
	PlmnID  *PlmnID `json:"plmnId"`
	N3IwfID string  `json:"n3IwfId,omitempty"`
	GnbID   *GnbID  `json:"gNbId,omitempty"`
	NgeNbID string  `json:"ngeNbId,omitempty"`
}

// GnbID struct
type GnbID struct {
	BitLength int32  `json:"bitLength"`
	GNBValue  string `json:"gNBValue"`
}

// MediaComponent Identifies a media component.
type MediaComponent struct {
	ContVer  int32 `json:"contVer,omitempty"`
	MedCompN int32 `json:"medCompN"`
	// Contains an AF application identifier.
	AfAppID   string              `json:"afAppId,omitempty"`
	AfRoutReq *RoutingRequirement `json:"afRoutReq,omitempty"`
	// Represents the content version of some content.
	Codecs      []string                     `json:"codecs,omitempty"`
	FStatus     FlowStatus                   `json:"fStatus,omitempty"`
	MarBwDl     string                       `json:"marBwDl,omitempty"`
	MarBwUl     string                       `json:"marBwUl,omitempty"`
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

// Ncgi struct
type Ncgi struct {
	PlmnID   PlmnID `json:"plmnId"`
	NrCellID string `json:"nrCellId"`
}

// PlmnID struct
type PlmnID struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"mnc"`
}

// PresenceInfo  struct
type PresenceInfo struct {
	PraID               string            `json:"praId,omitempty"`
	PresenceState       PresenceState     `json:"presenceState,omitempty"`
	TrackingAreaList    []Tai             `json:"trackingAreaList,omitempty"`
	EcgiList            []Ecgi            `json:"ecgiList,omitempty"`
	NcgiList            []Ncgi            `json:"ncgiList,omitempty"`
	GlobalRanNodeIDList []GlobalRanNodeID `json:"globalRanNodeIdList,omitempty"`
}

// QosNotificationControlInfo struct
/*
 * QosNotificationControlInfo Indicates whether the QoS targets for a GRB flow
 * are not  guaranteed or guaranteed again
 */
type QosNotificationControlInfo struct {
	NotifType QosNotifType `json:"notifType"`
	Flows     []Flows      `json:"flows,omitempty"`
}

// ResourcesAllocationInfo struct
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

// Tai Struct
type Tai struct {
	PlmnID PlmnID `json:"plmnId"`
	Tac    string `json:"tac"`
}

// TerminationInfo struct
/*
 * TerminationInfo indicates the cause for requesting the deletion of the
 * Individual Application Session Context resource
 */
type TerminationInfo struct {
	TermCause TerminationCause `json:"termCause"`
	// string providing an URI formatted according to IETF RFC 3986.
	ResURI string `json:"resUri"`
}

// UpPathChgEvent struct
type UpPathChgEvent struct {
	// string providing an URI formatted according to IETF RFC 3986.
	NotificationURI string `json:"notificationUri,omitempty"`
	/*
	 * It is used to set the value of Notification Correlation ID in the
	 * notification sent by the SMF.
	 */
	NotifCorreID string         `json:"notifCorreId"`
	DnaiChgType  DNAIChangeType `json:"dnaiChgType"`
}

// UsageThreshold struct
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
