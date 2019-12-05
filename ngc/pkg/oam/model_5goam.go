// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package oam

// AFServiceID is ID of service
type AFServiceID struct {
	AFServiceID string `json:"afServiceId,omitempty"`
}

// LocationService JSON struct
type LocationService struct {
	DNAI   string `json:"dnai,omitempty"` //DNAI value
	DNN    string `json:"dnn,omitempty"`  //DNN value
	TAC    int    `json:"tac,omitempty"`
	PriDNS string `json:"priDns,omitempty"`
	SecDNS string `json:"secDns,omitempty"`
	UPFIP  string `json:"upfIp,omitempty"`
	SNSSAI string `json:"snssai,omitempty"`
}

// AFService JSON struct
type AFService struct {
	AFServiceID     string          `json:"afServiceId,omitempty"`
	LocationService LocationService `json:"locationService,omitempty"`
}

// AFServiceList JSON struct
type AFServiceList struct {
	AfServiceList []AFService `json:"afServiceList,omitempty"`
}
