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

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

// Route describes traffic routing
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// AppPolicyIDs stores App Policy IDs with contents
var appPolicyIDs map[string]string

// OAMRoutes lists handlers for OAM routes
var AFRoutes = []Route{
	Route{
		"CreateDNSSubscription",
		strings.ToUpper("Post"),
		"/OAMTransactions",
		CreateDNSSubscription,
	},
	Route{
		"UpdateDNSSubscription",
		strings.ToUpper("Post"),
		"/OAMTransactions/{AppPolicyID}",
		UpdateDNSSubscription,
	},
	Route{
		"DeleteDNSSubscription",
		strings.ToUpper("Delete"),
		"/OAMTransactions/{AppPolicyID}",
		DeleteDNSSubscription,
	},
	Route{
		"GetDNSSubscription",
		strings.ToUpper("Get"),
		"/OAMTransactions/{AppPolicyID}",
		GetDNSSubscription,
	},
}

// CreateDNSSubscription Triggers creation of a new Traffic
// Influence Subscription at AF that will provide the arguments to NEF in
// another POST message. 
func CreateDNSSubscription(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("POST /OAMTransactions `%s`", b)

	_, exists := appPolicyIDs[string(b)]
	if exists {
		log.Printf("AppPolicyID `%s` already exists!", string(b))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
        OAMSendToFlexCore()
	appPolicyIDs[string(b)] = "{"+string(b)+"-json-contents}"
	w.WriteHeader(http.StatusOK)
}

// UpdateDNSSubscription Triggers update of an existing DNS Subscription
func UpdateDNSSubscription(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	appPolicyID := mux.Vars(r)["AppPolicyID"]
	log.Printf("POST /OAMTransactions/%s `%s`", appPolicyID, b)

	_, exists := appPolicyIDs[appPolicyID]
	if !exists {
		log.Printf("AppPolicyID `%s` does not exist!", appPolicyID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	appPolicyIDs[appPolicyID] = string(b)
	log.Println("Updated AppPolicyID:", appPolicyID)
	w.WriteHeader(http.StatusOK)
}

// DeleteDNSSubscription Triggers deletion of an existing DNS Subscription
func DeleteDNSSubscription(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appPolicyID := mux.Vars(r)["AppPolicyID"]
	log.Printf("DELETE /OAMTransactions/%s `%s`", appPolicyID, b)

	_, exists := appPolicyIDs[appPolicyID]
	if !exists {
		log.Printf("AppPolicyID `%s` does not exist!", appPolicyID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	delete(appPolicyIDs, appPolicyID)
	log.Println("Deleted AppPolicyID:", appPolicyID)
	w.WriteHeader(http.StatusOK)
}

// GetDNSSubscription Retrieves an existing DNS Subscription
func GetDNSSubscription(w http.ResponseWriter, r *http.Request) {
	appPolicyID := mux.Vars(r)["AppPolicyID"]
	log.Printf("GET /OAMTransactions/%s", appPolicyID)

	_, exists := appPolicyIDs[appPolicyID]
	if !exists {
		log.Printf("AppPolicyID `%s` does not exist!", appPolicyID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Write([]byte(appPolicyIDs[appPolicyID]))
	//w.WriteHeader(http.StatusOK)
}

// NewOAMRouter initializes OAM router
func NewOAMRouter() *mux.Router {
	appPolicyIDs = make(map[string]string)
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range AFRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}
	return router
}
