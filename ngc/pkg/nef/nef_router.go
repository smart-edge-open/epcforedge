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
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

/* Route structure which describes HTTP Request Handler type and other details
 * like name, method, path, etc */
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

/* NEF Routes lists which contains Routes with different HTTP Request handlers
 * for NEF */
var NEFRoutes = []Route{
	{
		"ReadAllTrafficInfluenceSubscription",
		strings.ToUpper("Get"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions",
		ReadAllTrafficInfluenceSubscription,
	},
	{
		"CreateTrafficInfluenceSubscription",
		strings.ToUpper("Post"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions",
		CreateTrafficInfluenceSubscription,
	},
	{
		"ReadTrafficInfluenceSubscription",
		strings.ToUpper("Get"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		ReadTrafficInfluenceSubscription,
	},
	{
		"UpdatePutTrafficInfluenceSubscription",
		strings.ToUpper("Put"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		UpdatePutTrafficInfluenceSubscription,
	},
	{
		"UpdatePatchTrafficInfluenceSubscription",
		strings.ToUpper("Patch"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		UpdatePatchTrafficInfluenceSubscription,
	},
	{
		"DeleteTrafficInfluenceSubscription",
		strings.ToUpper("Delete"),
		"/3gpp-traffic-influence/v1/{afId}/subscriptions/{subscriptionId}",
		DeleteTrafficInfluenceSubscription,
	},
}


/* Function: NewNEFRouter
*  Description: This function creates and initializes a NEF Router with all the
*               available routes for NEF Module. This router object is defined
*               in "github.com/gorilla/mux" package.
*  Input Args:
*     - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
*  Output Args:
*     - error: retruns pointer to created mux.Router object */
func NewNEFRouter(nefCtx *nefContext) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range NEFRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				string("nefCtx"),
				nefCtx)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	})

	return router
}
