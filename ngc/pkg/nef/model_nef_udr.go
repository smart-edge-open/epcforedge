/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

// TrafficInfluData traffic influ data
type TrafficInfluData struct {
	// Contains the Notification Correlation Id allocated by the NEF for the
	// UP path  change notification. It shall be included when the NEF
	// requests the UP path change notification
	UpPathChgNotifURI URI `json:"upPathChgNotifUri,omitempty"`
	// Identifies whether an application can be relocated once a location of
	// the application has been selected.
	AppReloInd bool `json:"appReloInd,omitempty"`
	// Identifies an application.
	// Required: true
	AfAppID string `json:"afAppId,omitempty"`
	// Identifies a dnn
	Dnn Dnn `json:"dnn,omitempty"`
	// Identifies Ethernet packet filters. Either "trafficFilters" or
	// "ethTrafficFilters" shall be included if applicable.
	EthTrafficFilters []EthFlowDescription `json:"ethTrafficFilters,omitempty"`
	// The identification of slice
	Snssai Snssai `json:"snssai,omitempty"`
	// Identifies a group of users.
	InterGroupID string `json:"interGroupId"`
	// supi Identifies a suer
	Supi Supi `json:"supi"`
	// Identifies IP packet filters. Either "trafficFilters" or
	// "ethTrafficFilters" shall be included if applicable.
	// Required: true
	// Min Items: 1
	TrafficFilters []FlowInfo `json:"trafficFilters"`
	// Identifies the N6 traffic routing requirement.
	// Required: true
	// Min Items: 1
	TrafficRoutes []RouteToLocation `json:"trafficRoutes"`
	// valid end time
	// Format: date-time
	ValidEndTime DateTime `json:"validEndTime,omitempty"`
	// valid start time
	// Format: date-time
	ValidStartTime DateTime `json:"validStartTime,omitempty"`
	// Identifies a network area information that the request applies only to
	// the traffic of UE(s) located in this specific zone
	NwAreaInfo NetworkAreaInfo `json:"nwAreaInfo,omitempty"`
	// Contains the Notification Correlation Id allocated by the NEF for the
	// UP path change notification.
	UpPathChgNotifCorreID string `json:"upPathChgNotifCorreId,omitempty"`
}

// TrafficInfluDataPatch traffic influ data patch
type TrafficInfluDataPatch struct {
	// Contains the Notification Correlation Id allocated by the NEF for the
	// UP path change notification.
	UpPathChgNotifCorreID string `json:"upPathChgNotifCorreId,omitempty"`
	// Identifies whether an application can be relocated once a location of
	// the application has been selected.
	AppReloInd bool `json:"appReloInd,omitempty"`
	// dnn
	Dnn Dnn `json:"dnn,omitempty"`
	// snssai
	Snssai Snssai `json:"snssai,omitempty"`
	// Identifies a group of users.
	InternalGroupID string `json:"internalGroupId,omitempty"`
	// Identifies Ethernet packet filters. Either "trafficFilters" or
	// "ethTrafficFilters" shall be included if applicable.
	// Min Items: 1
	EthTrafficFilters []EthFlowDescription `json:"ethTrafficFilters"`
	// supi
	Supi Supi `json:"supi,omitempty"`
	// Identifies IP packet filters. Either "trafficFilters" or
	// "ethTrafficFilters" shall be included if applicable.
	// Min Items: 1
	TrafficFilters []FlowInfo `json:"trafficFilters"`
	// Identifies the N6 traffic routing requirement.
	// Min Items: 1
	TrafficRoutes []RouteToLocation `json:"trafficRoutes"`
	// valid end time
	// Format: date-time
	ValidEndTime DateTime `json:"validEndTime,omitempty"`
	// valid start time
	// Format: date-time
	ValidStartTime DateTime `json:"validStartTime,omitempty"`
	// nw area info
	NwAreaInfo NetworkAreaInfo `json:"nwAreaInfo,omitempty"`
	// up path chg notif Uri
	UpPathChgNotifURI URI `json:"upPathChgNotifUri,omitempty"`
}

// NetworkAreaInfo Describes a network area information in which the NF service
// consumer requests the number of UEs.
// To be moved to the BDT policy mdoel in future
type NetworkAreaInfo struct {
	// Contains a list of E-UTRA cell identities.
	// Min Items: 1
	Ecgis []Ecgi `json:"ecgis"`
	// Contains a list of NR cell identities.
	// Min Items: 1
	Ncgis []Ncgi `json:"ncgis"`
	// Contains a list of NG RAN nodes.
	// Min Items: 1
	GRanNodeIds []GlobalRanNodeID `json:"gRanNodeIds"`
	// Contains a list of tracking area identities.
	// Min Items: 1
	Tais []Tai `json:"tais"`
}

// PfdContent represents the content of a PFD for an application identifier.
type PfdContent struct {
	// Identifies a PDF of an application identifier.
	PfdId string `json:"pfdId,omitempty"`
	// Represents a 3-tuple with protocol, server ip and server port for
	// UL/DL application traffic.
	FlowDescriptions []string `json:"flowDescriptions,omitempty"`
	// Indicates a URL or a regular expression which is used to match the
	// significant parts of the URL.
	Urls []string `json:"urls,omitempty"`
	// Indicates an FQDN or a regular expression as a domain name matching
	// criteria.
	DomainNames []string `json:"domainNames,omitempty"`
}

// PfdDataForApp represents the PFDs for an application identifier
type PfdDataForApp struct {
	// Identifier of an application.
	AppId ApplicationId `json:"appId"`
	// PFDs for the application identifier.
	Pfds []PfdContent `json:"pfds"`
	// Caching time for an application identifier.
	CachingTime DateTime `json:"cachingTime,omitempty"`
}
