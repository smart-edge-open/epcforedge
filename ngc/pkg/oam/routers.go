package oam 

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Do not steal my cake!")
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
		"/afRegisters",
		Add,
	},

	Route{
		"Delete",
		strings.ToUpper("Delete"),
		"/afRegisters/{afId}",
		Delete,
	},

	Route{
		"DeleteDns",
		strings.ToUpper("Delete"),
		"/afRegisters/{afId}/locationServices/{dnai}",
		DeleteDns,
	},

	Route{
		"Get",
		strings.ToUpper("Get"),
		"/afRegisters/{afId}",
		Get,
	},

	Route{
		"GetAll",
		strings.ToUpper("Get"),
		"/afRegisters",
		GetAll,
	},

	Route{
		"Update",
		strings.ToUpper("Patch"),
		"/afRegisters/{afId}",
		Update,
	},
}
