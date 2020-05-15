// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

// AfwebsockNotifConfig Websocket configuration for delivering AF notifications
// for all Policy Authorization Events
type AfwebsockNotifConfig struct {
	// The Websocket Uri used for Notification delivery.
	// This is set by AF and is present in the response from AF to consumer
	WebsocketURI string `json:"websocketUri,omitempty"`
	// Set by the consumer to indicate that the Websocket delivery is requested.
	RequestWebsocketURI bool `json:"requestWebsocketUri,omitempty"`
	// Unique Identification of the consumer. Set by consumer
	ConsumerID string `json:"consumerID,omitempty"`
}

// NotifyEvent represents the list of events for which the consumer can register
// with AF for notifications
type NotifyEvent string

// Possible values of NotifyEvent
const (
	// Event when the UP path changes for the PDU session.
	UPPathChangeEvent NotifyEvent = "UP_PATH_CHANGE"
)

// Afnotification sent by AF to consumer when the event is receievd from
// NEF/CoreNetwork. This is sent over Websocket/HTTP2.0
type Afnotification struct {
	Event NotifyEvent `json:"event"`
	// Payload is the Notification based on the specific NotifyEvent.
	// For UP_PATH_CHANGE, it is NotificationUpPathChg
	Payload []byte `json:"payload"`
}

// NotificationUpPathChg structure for the notification of
// UP_PATH_CHANGE event
type NotificationUpPathChg struct {
	// NotifyID Identifies the notification correlation
	// In PolicyAuth this is appSessionId and in Traffic Influence
	// this is afTransId
	NotifyID string `json:"notifyID"`
	// DnaiChgType
	DNAIChgType DNAIChangeType `json:"dnaiChgType,omitempty"`
	// SourceTrafficRoute
	SourceTrafficRoute RouteToLocation `json:"sourceTrafficRoute,omitempty"`
	// SubscribedEvent
	SubscribedEvent SubscribedEvent `json:"subscribedEvent,omitempty"`
	// TargetTrafficRoute
	TargetTrafficRoute RouteToLocation `json:"targetTrafficRoute,omitempty"`
	// Gpsi
	GPSI string `json:"gpsi,omitempty"`
	// SrcUeIpv4Addr
	SrcUEIPv4Addr SrcUEIPv4Addr `json:"srcUeIpv4Addr,omitempty"`
	// SrcUeIpv6Prefix
	SrcUEIPv6Prefix SrcUEIPv6Prefix `json:"srcUeIpv6Prefix,omitempty"`
	// TgtUeIpv4Addr
	TgtUEIP4Addr TgtUEIPv4Addr `json:"tgtUeIpv4Addr,omitempty"`
	// TgtUeIpv6Prefix
	TgtUEIPv6Prefix TgtUEIPv6Prefix `json:"tgtUeIpv6Prefix,omitempty"`
	// UeMac
	UEMac UEMac `json:"ueMac,omitempty"`
}
