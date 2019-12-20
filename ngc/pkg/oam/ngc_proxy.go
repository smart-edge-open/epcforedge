// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package oam

import (
	"errors"
	logger "github.com/open-ness/common/log"
	"net/http"
)

// NGCType : type of ngc
var NGCType string // APISTUB or 5GFLEXCORE

// URLBase : base path
var URLBase string

const apiStub = "APISTUB"

var log = logger.DefaultLogger.WithField("oam", nil)

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

	log.Infof("URL GetAll: %s\n", URLBase+r.URL.Path)

	if NGCType == apiStub {
		APIStubGetAll(w, r)
	} else {
		log.Errf("GetAll Failed with TargetNGC %s\n", NGCType)
		w.WriteHeader(http.StatusNotFound)
	}
}

// ProxyAdd : add by proxy
func ProxyAdd(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Add: %s\n", URLBase+r.URL.Path)

	if NGCType == apiStub {
		APIStubAdd(w, r)
	} else {
		log.Errf("Add Failed with TargetNGC %s\n", NGCType)
		w.WriteHeader(http.StatusNotFound)
	}
}

// ProxyDel : del by proxy
func ProxyDel(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Del: %s\n", URLBase+r.URL.Path)

	if NGCType == apiStub {
		APIStubDel(w, r)
	} else {
		log.Errf("Del Failed with TargetNGC %s\n", NGCType)
		w.WriteHeader(http.StatusNotFound)
	}

}

// ProxyGet : get by proxy
func ProxyGet(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Get: %s\n", URLBase+r.URL.Path)

	if NGCType == apiStub {
		APIStubGet(w, r)
	} else {
		log.Errf("Get Failed with TargetNGC %s\n", NGCType)
		w.WriteHeader(http.StatusNotFound)
	}

}

// ProxyUpdate : udpate byproxy
func ProxyUpdate(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Update: %s\n", URLBase+r.URL.Path)

	if NGCType == apiStub {
		APIStubUpdate(w, r)
	} else {
		log.Errf("Update Failed with TargetNGC %s\n", NGCType)
		w.WriteHeader(http.StatusNotFound)
	}

}
