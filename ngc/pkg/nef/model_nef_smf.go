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
	SourceUeIpv4Addr Ipv4Addr `json:"sourceUeIpv4Addr,omitempty"`
	// The Ipv6 Address Prefix of the served UE for the source DNAI. May
	// be included for event "UP_PATH_CH".
	SourceUeIpv6Prefix Ipv6Prefix `json:"sourceUeIpv6Prefix,omitempty"`
	// The IPv4 Address of the served UE for the target DNAI. May be included
	// for event "UP_PATH_CH".
	TargetUeIpv4Addr Ipv4Addr `json:"targetUeIpv4Addr,omitempty"`
	// The Ipv6 Address Prefix of the served UE for the target DNAI. May
	// be included for event "UP_PATH_CH".
	TargetUeIpv6Prefix Ipv6Prefix `json:"targetUeIpv6Prefix,omitempty"`
	// N6 traffic routing information for the source DNAI. Shall be included
	// for event "UP_PATH_CH".
	SourceTraRouting RouteToLocation `json:"sourceTraRouting,omitempty"`
	// N6 traffic routing information for the target DNAI. Shall be included
	// for event "UP_PATH_CH".
	TargetTraRouting RouteToLocation `json:"targetTraRouting,omitempty"`
	// UE MAC address. May be included for event "UP_PATH_CH".
	UeMac MacAddr48 `json:"ueMac,omitempty"`
	// Added IPv4 Address(es). May be included for event "UE_IP_CH".
	AdIpv4Addr Ipv4Addr `json:"adIpv4Addr,omitempty"`
	// Added Ipv6 Address Prefix(es). May be included for event "UE_IP_CH".
	AdIpv6Prefix Ipv6Prefix `json:"adIpv6Prefix,omitempty"`
	// Removed IPv4 Address(es). May be included for event "UE_IP_CH".
	ReIpv4Addr Ipv4Addr `json:"reIpv4Addr,omitempty"`
	// Removed Ipv6 Address Prefix(es). May be included for event "UE_IP_CH".
	ReIpv6Prefix Ipv6Prefix `json:"reIpv6Prefix,omitempty"`
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
