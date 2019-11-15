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

import (
	"net/http"
)

func Add(w http.ResponseWriter, r *http.Request) {
        ProxyAdd(w, r)
}

func Delete(w http.ResponseWriter, r *http.Request) {
        ProxyDel(w, r)
}

func DeleteDns(w http.ResponseWriter, r *http.Request) {
        ProxyDelDnn(w, r)
}

func Get(w http.ResponseWriter, r *http.Request) {
        ProxyGet(w, r)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
        ProxyGetAll(w, r)
}

func Update(w http.ResponseWriter, r *http.Request) {
        ProxyUpdate(w, r)
}
