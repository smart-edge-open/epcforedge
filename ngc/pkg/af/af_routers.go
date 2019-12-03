// Copyright 2019 Intel Corporation, Inc. All rights reserved
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

package ngcaf

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// DefaultNotifURL const
const DefaultNotifURL = "/af/v1/notifications"

type keyType string

// Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes type
type Routes []Route

// NewAFRouter function
func NewAFRouter(afCtx *AFContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range afRoutes {
		var handler http.Handler = route.HandlerFunc
		handler = afLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				keyType("af-ctx"),
				afCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	return router
}

// NewNotifRouter function
func NewNotifRouter(afCtx *AFContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range notifRoutes {
		var handler http.Handler = route.HandlerFunc
		handler = afLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				keyType("af-ctx"),
				afCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	return router
}

var notifRoutes = Routes{
	Route{
		"NotificationPost",
		strings.ToUpper("Post"),
		DefaultNotifURL,
		NotificationPost,
	},
}

var afRoutes = Routes{

	Route{
		"GetAllSubscriptions",
		strings.ToUpper("Get"),
		"/af/v1/subscriptions",
		GetAllSubscriptions,
	},

	Route{
		"GetSubscription",
		strings.ToUpper("Get"),
		"/af/v1/subscriptions/{subscriptionId}",
		GetSubscription,
	},

	Route{
		"DeleteSubscription",
		strings.ToUpper("Delete"),
		"/af/v1/subscriptions/{subscriptionId}",
		DeleteSubscription,
	},

	Route{
		"SubscriptionPatch",
		strings.ToUpper("Patch"),
		"/af/v1/subscriptions/{subscriptionId}",
		ModifySubscriptionPatch,
	},

	Route{
		"CreateSubscription",
		strings.ToUpper("Post"),
		"/af/v1/subscriptions",
		CreateSubscription,
	},

	Route{
		"SubscriptionPut",
		strings.ToUpper("Put"),
		"/af/v1/subscriptions/{subscriptionId}",
		ModifySubscriptionPut,
	},
}
