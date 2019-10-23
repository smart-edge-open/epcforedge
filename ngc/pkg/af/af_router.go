package edgeaf

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

// AFRoutes lists handlers for AF routes
var AFRoutes = []Route{
	Route{
		"CreateTrafficInfluenceSubscription",
		strings.ToUpper("Post"),
		"/AFTransactions",
		CreateTrafficInfluenceSubscription,
	},
	Route{
		"UpdateTrafficInfluenceSubscription",
		strings.ToUpper("Post"),
		"/AFTransactions/{AppPolicyID}",
		UpdateTrafficInfluenceSubscription,
	},
	Route{
		"DeleteTrafficInfluenceSubscription",
		strings.ToUpper("Delete"),
		"/AFTransactions/{AppPolicyID}",
		DeleteTrafficInfluenceSubscription,
	},
	Route{
		"GetTrafficInfluenceSubscription",
		strings.ToUpper("Get"),
		"/AFTransactions/{AppPolicyID}",
		GetTrafficInfluenceSubscription,
	},
}

// CreateTrafficInfluenceSubscription Triggers creation of a new Traffic
// Influence Subscription at AF that will provide the arguments to NEF in
// another POST message. 
func CreateTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("POST /AFTransactions `%s`", b)

	_, exists := appPolicyIDs[string(b)]
	if exists {
		log.Printf("AppPolicyID `%s` already exists!", string(b))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appPolicyIDs[string(b)] = "{"+string(b)+"-json-contents}"
	w.WriteHeader(http.StatusOK)
}

// UpdateTrafficInfluenceSubscription Triggers update of an existing Traffic
// Influence Subscription
func UpdateTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	appPolicyID := mux.Vars(r)["AppPolicyID"]
	log.Printf("POST /AFTransactions/%s `%s`", appPolicyID, b)

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

// DeleteTrafficInfluenceSubscription Triggers deletion of an existing Traffic
// Influence Subscription
func DeleteTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appPolicyID := mux.Vars(r)["AppPolicyID"]
	log.Printf("DELETE /AFTransactions/%s `%s`", appPolicyID, b)

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

// GetTrafficInfluenceSubscription Retrieves an existing Traffic Influence
// Subscription
func GetTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {
	appPolicyID := mux.Vars(r)["AppPolicyID"]
	log.Printf("GET /AFTransactions/%s", appPolicyID)

	_, exists := appPolicyIDs[appPolicyID]
	if !exists {
		log.Printf("AppPolicyID `%s` does not exist!", appPolicyID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Write([]byte(appPolicyIDs[appPolicyID]))
	//w.WriteHeader(http.StatusOK)
}

// NewAFRouter initializes AF router
func NewAFRouter() *mux.Router {
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
