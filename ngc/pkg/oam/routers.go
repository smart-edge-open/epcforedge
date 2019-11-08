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
		Add,
	},

	Route{
		"Delete",
		strings.ToUpper("Delete"),
		"/oam/v1/af/services/{afId}",
		Delete,
	},

	Route{
		"DeleteDns",
		strings.ToUpper("Delete"),
		"/oam/v1/af/services/{afId}/locationServices/{dnai}",
		DeleteDns,
	},

	Route{
		"Get",
		strings.ToUpper("Get"),
		"/oam/v1/af/services/{afId}",
		Get,
	},

	Route{
		"GetAll",
		strings.ToUpper("Get"),
		"/oam/v1/af/services",
		GetAll,
	},

	Route{
		"Update",
		strings.ToUpper("Patch"),
		"/oam/v1/af/services/{afId}",
		Update,
	},
}
