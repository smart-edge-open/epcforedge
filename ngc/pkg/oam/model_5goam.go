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

package oam

// AFServiceID is ID of service
type AFServiceID struct {
     AFServiceID   string   `json:"afServiceId,omitempty"` // AF service ID
}

// LocationService JSON struct
type LocationService struct {
     DNAI      string     `json:"dnai,omitempty"`  //DNAI value
     DNN       string     `json:"dnn,omitempty"`   //DNN value
     TAC       int        `json:"tac,omitempty"`
     PriDNS    string     `json:"priDns,omitempty"`
     SecDNS    string     `json:"secDns,omitempty"`
     UPFIP     string     `json:"upfIp,omitempty"`
     SNSSAI    string     `json:"snssai,omitempty"`
}


// AFService JSON struct
type AFService struct {
     AFServiceID      string          `json:"afServiceId,omitempty"`
     LocationService  LocationService `json:"locationService,omitempty"`
}

// AFServiceList JSON struct
type AFServiceList struct {
     AfServiceList []AFService      `json:"afServiceList,omitempty"`
}
