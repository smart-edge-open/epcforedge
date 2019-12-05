// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

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
		"/ngcoam/v1/af/services",
		add,
	},

	Route{
		"Delete",
		strings.ToUpper("Delete"),
		"/ngcoam/v1/af/services/{afServiceId}",
		delete,
	},

	Route{
		"Get",
		strings.ToUpper("Get"),
		"/ngcoam/v1/af/services/{afServiceId}",
		get,
	},

	Route{
		"GetAll",
		strings.ToUpper("Get"),
		"/ngcoam/v1/af/services",
		getAll,
	},

	Route{
		"Update",
		strings.ToUpper("Patch"),
		"/ngcoam/v1/af/services/{afServiceId}",
		update,
	},
}
