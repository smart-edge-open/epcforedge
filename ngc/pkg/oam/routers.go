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
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Route : route handler structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes : slice for route
type Routes []Route

// NewRouter : function of mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// Index : function for index
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello From OpenNESS!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"Add",
		strings.ToUpper("Post"),
		"/oam/v1/af/services",
		add,
	},

	Route{
		"Delete",
		strings.ToUpper("Delete"),
		"/oam/v1/af/services/{afServiceId}",
		delete,
	},

	Route{
		"Get",
		strings.ToUpper("Get"),
		"/oam/v1/af/services/{afServiceId}",
		get,
	},

	Route{
		"GetAll",
		strings.ToUpper("Get"),
		"/oam/v1/af/services",
		getAll,
	},

	Route{
		"Update",
		strings.ToUpper("Patch"),
		"/oam/v1/af/services/{afServiceId}",
		update,
	},
}
