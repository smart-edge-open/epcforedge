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

// AFRoutes lists handlers for AF routes
var NEFRoutes = []Route{
	Route{
		//Read all of the active subscriptions for the AF
		"ReadAllTrafficInfluenceSubscription",
		strings.ToUpper("Get"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions",
		ReadAllTrafficInfluenceSubscription,
	},
	Route{
		"CreateTrafficInfluenceSubscription",
		strings.ToUpper("Post"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions",
		CreateTrafficInfluenceSubscription,
	},
	Route{
		"ReadTrafficInfluenceSubscription",
		strings.ToUpper("Get"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		ReadTrafficInfluenceSubscription,
	},
	Route{
		"UpdatePutTrafficInfluenceSubscription",
		strings.ToUpper("Put"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		UpdatePutTrafficInfluenceSubscription,
	},
	Route{
		"UpdatePatchTrafficInfluenceSubscription",
		strings.ToUpper("Patch"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		UpdatePatchTrafficInfluenceSubscription,
	},
	Route{
		"DeleteTrafficInfluenceSubscription",
		strings.ToUpper("Delete"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		DeleteTrafficInfluenceSubscription,
	},
}

func ReadAllTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	log.Printf("===============================================")
	log.Printf(" Method : GET ")
	log.Printf(" URL PATH : " + r.URL.Path[1:])

	vars := mux.Vars(r)

	log.Printf(" AFID  : %s", vars["afId"])
	w.WriteHeader(http.StatusOK)
}

func CreateTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	log.Printf("===============================================")
	log.Printf(" Method : POST ")
	log.Printf(" URL PATH : "+ r.URL.Path[1:])
	vars := mux.Vars(r)
	log.Printf(" AFID  : %s", vars["afId"])


	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf(" BODY : `%s`", b)
//	log.Printf(r.URL.Path[1:])

	w.WriteHeader(http.StatusOK)
}

func ReadTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	log.Printf("===============================================")
	log.Printf(" Method : GET ")
	log.Printf(" URL PATH : " + r.URL.Path[1:])

	vars := mux.Vars(r)

	log.Printf(" AFID  : %s", vars["afId"])
	log.Printf(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])
	w.WriteHeader(http.StatusOK)
}
func UpdatePutTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	log.Printf("===============================================")
	log.Printf(" Method : PUT ")
	log.Printf(" URL PATH : " + r.URL.Path[1:])

	vars := mux.Vars(r)

	log.Printf(" AFID  : %s", vars["afId"])
	log.Printf(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])
	w.WriteHeader(http.StatusOK)
}
func UpdatePatchTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	log.Printf("===============================================")
	log.Printf(" Method : PATCH ")
	log.Printf(" URL PATH : " + r.URL.Path[1:])

	vars := mux.Vars(r)

	log.Printf(" AFID  : %s", vars["afId"])
	log.Printf(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])
	w.WriteHeader(http.StatusOK)
}
func DeleteTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	log.Printf("===============================================")
	log.Printf(" Method : DELETE ")
	log.Printf(" URL PATH : " + r.URL.Path[1:])

	vars := mux.Vars(r)

	log.Printf(" AFID  : %s", vars["afId"])
	log.Printf(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])
	w.WriteHeader(http.StatusOK)
}

// NewNEFRouter initializes NEF router
func NewNEFRouter() *mux.Router {
	
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range NEFRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}
	return router
}