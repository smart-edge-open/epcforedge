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
    "errors"
    "log"
    "net/http"
)

// NGCType : type of ngc
var NGCType   string  // APISTUB or 5GFLEXCORE

// URLBase : base path
var URLBase   string
const apiStub = "APISTUB"

// InitProxy : init proxy
// The proxy acts as reverse proxy to handle request from CNCA 
// and forward it to the target. 
// The target can be API_STUB_TEST or flexcore. 
//NOTE: current version only support API_STUB.
func InitProxy(npcEndpoint string, redirectTarget string, path string) error {
    URLBase = "http://" + npcEndpoint
    NGCType = redirectTarget
    if NGCType == apiStub {
        if nil != APIStubInit(path) {
           return errors.New("init error")
        }
    } else {
       return errors.New("can't not support flexcore")
    }
    
    return nil

}

// ProxyGetAll : get all by proxy
func ProxyGetAll(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL GetAll: %s\n", URLBase + r.URL.Path)

    if NGCType == apiStub {
        APIStubGetAll(w, r)
    } else {
        log.Printf("GetAll Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(http.StatusNotFound)
    }
}

// ProxyAdd : add by proxy
func ProxyAdd(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Add: %s\n", URLBase + r.URL.Path)

    if NGCType == apiStub {
        APIStubAdd(w, r)
    } else {
        log.Printf("Add Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(http.StatusNotFound)
    }
}

// ProxyDel : del by proxy
func ProxyDel(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Del: %s\n", URLBase + r.URL.Path)

    if NGCType == apiStub {
        APIStubDel(w, r)
    } else {
        log.Printf("Del Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(http.StatusNotFound)
    }

}

// ProxyGet : get by proxy
func ProxyGet(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Get: %s\n", URLBase + r.URL.Path)

    if NGCType == apiStub {
        APIStubGet(w, r)
    } else {
        log.Printf("Get Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(http.StatusNotFound)
    }

}

// ProxyUpdate : udpate byproxy
func ProxyUpdate(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Update: %s\n", URLBase + r.URL.Path)

    if NGCType == apiStub {
        APIStubUpdate(w, r)
    } else {
        log.Printf("Update Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(http.StatusNotFound)
    }

}
