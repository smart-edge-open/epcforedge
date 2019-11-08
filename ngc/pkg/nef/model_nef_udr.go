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

// TrafficInfluData traffic influ data
type TrafficInfluData struct {
	// Contains the Notification Correlation Id allocated by the NEF for the UP path
	// change notification. It shall be included when the NEF requests the UP path
	// change notification
	UpPathChgNotifURI URI `json:"upPathChgNotifUri,omitempty"`
	// Identifies whether an application can be relocated once a location of the application has been selected.
	AppReloInd bool `json:"appReloInd,omitempty"`
	// Identifies an application.
	// Required: true
	AfAppID string `json:"afAppId,omitempty"`
	// Identifies a dnn
	Dnn Dnn `json:"dnn,omitempty"`
	// Identifies Ethernet packet filters. Either "trafficFilters" or "ethTrafficFilters" shall be included if applicable.
	EthTrafficFilters []EthFlowDescription `json:"ethTrafficFilters,omitempty"`
	// The identification of slice
	Snssai Snssai `json:"snssai,omitempty"`
	// Identifies a group of users.
	InterGroupID string `json:"interGroupId"`
	// supi Identifies a suer
	Supi Supi `json:"supi"`
	// Identifies IP packet filters. Either "trafficFilters" or "ethTrafficFilters" shall be included if applicable.
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
	// Identifies a network area information that the request applies only to the traffic of
	// UE(s) located in this specific zone
	NwAreaInfo NetworkAreaInfo `json:"nwAreaInfo,omitempty"`
	// Contains the Notification Correlation Id allocated by the NEF for the UP path change notification.
	UpPathChgNotifCorreID URI `json:"upPathChgNotifCorreId,omitempty"`
}

// TrafficInfluDataPatch traffic influ data patch
type TrafficInfluDataPatch struct {
	// Contains the Notification Correlation Id allocated by the NEF for the UP path change notification.
	UpPathChgNotifCorreID string `json:"upPathChgNotifCorreId,omitempty"`
	// Identifies whether an application can be relocated once a location of the application has been selected.
	AppReloInd bool `json:"appReloInd,omitempty"`
	// dnn
	Dnn Dnn `json:"dnn,omitempty"`
	// snssai
	Snssai Snssai `json:"snssai,omitempty"`
	// Identifies a group of users.
	InternalGroupID string `json:"internalGroupId,omitempty"`
	// Identifies Ethernet packet filters. Either "trafficFilters" or "ethTrafficFilters" shall be included if applicable.
	// Min Items: 1
	EthTrafficFilters []EthFlowDescription `json:"ethTrafficFilters"`
	// supi
	Supi Supi `json:"supi,omitempty"`
	// Identifies IP packet filters. Either "trafficFilters" or "ethTrafficFilters" shall be included if applicable.
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

// NetworkAreaInfo Describes a network area information in which the NF service consumer requests the number of UEs.
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
