/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

package af

// Dnai : string identifying the Data Network Area Identifier
type Dnai string

// DnaiChangeType : string identifying the DNAI change type
// Possible values are
// - EARLY: Early notification of UP path reconfiguration.
// - EARLY_LATE: Early and late notification of UP path reconfiguration. This
// value shall only be present in the subscription to the DNAI change event.
// - LATE: Late notification of UP path reconfiguration.
type DnaiChangeType string

// DateTime is in the date-time format
type DateTime string

// PduSessionID Valid values are 0 to 255
type PduSessionID uint8

// Supi : Subscription Permanent Identifier
// pattern: '^(imsi-[0-9]{5,15}|nai-.+|.+)$'
type Supi string

// Gpsi : Generic Public Servie Identifiers asssociated wit the UE
// pattern '^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$'
type Gpsi string

// IPv6Prefix : string representing the Ipv6 Prefix
// pattern: '^((:|(0?|([1-9a-f][0-9a-f]{0,3}))):)((0?|([1-9a-f][0-9a-f]{0,3}))
// :){0,6}(:|(0?|([1-9a-f][0-9a-f]{0,3})))(\/(([0-9])|([0-9]{2})|(1[0-1][0-9])
//|(12[0-8])))$'
// pattern: '^((([^:]+:){7}([^:]+))|((([^:]+:)*[^:]+)?::(([^:]+:)*[^:]+)?))
// (\/.+)$'
// example: '2001:db8:abcd:12::0/64'
type IPv6Prefix string

// MacAddr48 : Identifies a MAC address
// pattern: '^([0-9a-fA-F]{2})((-[0-9a-fA-F]{2}){5})$'
type MacAddr48 string

// NsmfEventExposureNotification Provides Information about observed events
type NsmfEventExposureNotification struct {
	// Notification correlation ID used to identify the subscription which the
	// notification is corresponding to. It shall be set to the same value as
	// the notifId attribute of NsmfEventExposure data type or the value of
	// "notifiCorreId" within the UpPathChgEvent data type defined in
	// 3GPP TS 29.512 [14].
	NotifID string `json:"notifId"`
	// Notifications about Individual Events
	// Note 3GPP 29508 defines this as EventNotification but this conflicts with
	//  3GPP 29522 EventNotification so rename it
	EventNotifs []NsmEventNotification `json:"eventNotifs"`
}

// NsmEventNotification describes about the Notifications about Individual
// Events from SMF
type NsmEventNotification struct {
	// Event that is notified.
	Event SmfEvent `json:"event"`
	// Time at which the event is observed.
	TimeStamp DateTime `json:"timeStamp"`
	// Subscription Permanent Identifier. It is included when the subscription
	// applies to a group of UE(s) or any UE.
	Supi Supi `json:"supi,omitempty"`
	// Identifies a GPSI. It shall contain an MSISDN. It is included when it is
	// available and the subscription applies to a group of UE(s) or any UE.
	Gpsi Gpsi `json:"gpsi,omitempty"`
	// Source DN Access Identifier. Shall be included for event "UP_PATH_CH"
	// if the DNAI changed (NOTE).
	SourceDnai Dnai `json:"sourceDnai,omitempty"`
	// Target DN Access Identifier. Shall be included for event "UP_PATH_CH"
	// if the DNAI changed (NOTE).
	TargetDnai Dnai `json:"targetDnai,omitempty"`
	// DNAI Change Type. Shall be included for event "UP_PATH_CH".
	DnaiChgType DnaiChangeType `json:"dnaiChgType,omitempty"`
	// The IPv4 Address of the served UE for the source DNAI. May be included
	// for event "UP_PATH_CH".
	SourceUeIpv4Addr IPv4Addr `json:"sourceUeIpv4Addr,omitempty"`
	// The Ipv6 Address Prefix of the served UE for the source DNAI. May
	// be included for event "UP_PATH_CH".
	SourceUeIpv6Prefix IPv6Prefix `json:"sourceUeIpv6Prefix,omitempty"`
	// The IPv4 Address of the served UE for the target DNAI. May be included
	// for event "UP_PATH_CH".
	TargetUeIpv4Addr IPv4Addr `json:"targetUeIpv4Addr,omitempty"`
	// The Ipv6 Address Prefix of the served UE for the target DNAI. May
	// be included for event "UP_PATH_CH".
	TargetUeIpv6Prefix IPv6Prefix `json:"targetUeIpv6Prefix,omitempty"`
	// N6 traffic routing information for the source DNAI. Shall be included
	// for event "UP_PATH_CH".
	SourceTraRouting RouteToLocation `json:"sourceTraRouting,omitempty"`
	// N6 traffic routing information for the target DNAI. Shall be included
	// for event "UP_PATH_CH".
	TargetTraRouting RouteToLocation `json:"targetTraRouting,omitempty"`
	// UE MAC address. May be included for event "UP_PATH_CH".
	UeMac MacAddr48 `json:"ueMac,omitempty"`
	// Added IPv4 Address(es). May be included for event "UE_IP_CH".
	AdIpv4Addr IPv4Addr `json:"adIpv4Addr,omitempty"`
	// Added Ipv6 Address Prefix(es). May be included for event "UE_IP_CH".
	AdIpv6Prefix IPv6Prefix `json:"adIpv6Prefix,omitempty"`
	// Removed IPv4 Address(es). May be included for event "UE_IP_CH".
	ReIpv4Addr IPv4Addr `json:"reIpv4Addr,omitempty"`
	// Removed Ipv6 Address Prefix(es). May be included for event "UE_IP_CH".
	ReIpv6Prefix IPv6Prefix `json:"reIpv6Prefix,omitempty"`
	// New PLMN ID. Shall be included for event "PLMN_CH".
	PlmnID PlmnID `json:"plmnId,omitempty"`
	// New Access Type. Shall be included for event "AC_TY_CH".
	AccType AccessType `json:"accType,omitempty"`
	// PDU session ID. Shall be included for event "PDU_SES_REL".
	PduSeID PduSessionID `json:"pduSeId,omitempty"`
}

// SmfEvent This string provides forward-compatibility with future
//          extensions to the enumeration but is not used to encode
//          content defined in the present version of this API
// Values currently defined are:
// - AC_TY_CH
// - UP_PATH_CH
// - PDU_SES_REL
// - PLMN_CH
// - UE_IP_CH
type SmfEvent string
