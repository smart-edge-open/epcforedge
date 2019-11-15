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

// ServiceList JSON struct
type AfServiceList struct {
        AfServices []AfService      `json:"afServices,omitempty"`
}

// AF Service JSON struct
type AfService struct {
        AfInstance    string         `json:"afInstance,omitempty"`
        LocalServices []LocalService `json:"localServices,omitempty"`
}


// local Service JSON struct
type LocalService struct {
        Dnai   string                   `json:"dnai,omitempty"`
        Dnn    string                   `json:"dnn,omitempty"`
        Dns    string                   `json:"dns,omitempty"`
}

//  AfId struct
type AfId struct {
        AfId   string                   `json:"afid,omitempty"`
}
