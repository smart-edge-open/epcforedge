/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngccntest

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	oauth2 "github.com/otcshare/epcforedge/ngc/pkg/oauth2"
)

type cnTestCtxKey string

// Route : Structure which describes HTTP Request Handler type and other
//         details like name, method, path, etc
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc /*    */
}

// CNTESTRoutes : CN-TEST Routes lists which contains Routes with different HTTP
//             Request handlers for CN-TEST
var CNTESTRoutes = []Route{
	// Policy Authorization Routes
	{
		"PolicyAuthorizationAppSessionCreate",
		strings.ToUpper("Post"),
		"/npcf-policyauthorization/v1/app-sessions",
		PolicyAuthorizationAppSessionCreate,
	},
	{
		"PolicyAuthorizationAppSessionGet",
		strings.ToUpper("Get"),
		"/npcf-policyauthorization/v1/app-sessions/{appSessionId}",
		PolicyAuthorizationAppSessionGet,
	},
	{
		"PolicyAuthorizationAppSessionPatch",
		strings.ToUpper("Patch"),
		"/npcf-policyauthorization/v1/app-sessions/{appSessionId}",
		PolicyAuthorizationAppSessionPatch,
	},
	{
		"PolicyAuthorizationAppSessionDelete",
		strings.ToUpper("Delete"),
		"/npcf-policyauthorization/v1/app-sessions/{appSessionId}/delete",
		PolicyAuthorizationAppSessionDelete,
	},
	{
		"PolicyAuthorizationAppSessionSubscribe",
		strings.ToUpper("Put"),
		"/npcf-policyauthorization/v1/app-sessions/{appSessionId}/events-subscription",
		PolicyAuthorizationAppSessionSubscribe,
	},
	{
		"PolicyAuthorizationAppSessionUnsubscribe",
		strings.ToUpper("Delete"),
		"/npcf-policyauthorization/v1/app-sessions/{appSessionId}/events-subscription",
		PolicyAuthorizationAppSessionUnsubscribe,
	},
}

// NewCnTestRouter : This function creates and initializes a CN-TEST Router with all
//                the available routes for CN-TEST Module. This router object is
//                defined in "github.com/gorilla/mux" package.
//  Input Args:
//     - cnTestCtx: This is CN-TEST Module Context. This contains the CN-TEST Module Data.
//  Output Args:
//     - error: retruns pointer to created mux.Router object
func NewCnTestRouter(cnTestCtx *cnTestContext) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	//Subscription for notification to be added

	for _, route := range CNTESTRoutes {

		var handler http.Handler = route.Handler
		handler = cnTestRouteLogger(handler, route.Name)

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
				cnTestCtxKey("cnTestCtx"),
				cnTestCtx)

			if cnTestCtx.cfg.OAuth2Support {
				if cnTestValidateAccessToken(w, r) {
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			} else {
				//OAuth2 disabled
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	})
	return router
}

func cnTestValidateAccessToken(w http.ResponseWriter, r *http.Request) bool {

	reqToken := r.Header.Get("Authorization")

	if len(reqToken) == 0 {
		log.Info("Authorization header missing")
		//Authorization header is not present
		w.Header().Set("WWW-Authenticate", "Bearer realm="+r.RequestURI)

		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	//Get the token
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	status, err := oauth2.ValidateAccessToken(reqToken)

	if err != nil {
		log.Infoln("Token Validation failed")
		if status == oauth2.StatusInvalidToken {
			w.Header().Set("WWW-Authenticate", "Bearer realm="+r.RequestURI)
			w.WriteHeader(http.StatusUnauthorized)
			return false
		} else if status == oauth2.StatusBadRequest {
			w.WriteHeader(http.StatusBadRequest)
			return false
		}
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}

// cnTestRouteLogger : This function logs data received in HTTP request.
// Input Args:
//    - httpHandler: This is HTTP handler function pointer for HTTP Request
//                   Received
//    - name: This is route name.
// Output Args:
//    - httpHandler: This is HTTP handler function pointer for HTTP Request
//                   Received. This HTTP Handler is actually the updated HTTP
//                   Handler. The updated HTTP Handler now can logging of HTTP
//                   Request Info also
func cnTestRouteLogger(httpHandler http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Infof("HTTP Request Received :%s", r.Proto)
		log.Infof("===============================================")
		log.Infof(" Method : %s ", r.Method)
		log.Infof(" URL PATH : %s", r.RequestURI)
		log.Infof(" Route Name : %s", name)
		log.Infof("===============================================")
		log.Infof("HTTP Request Handling -- STARTS")

		httpHandler.ServeHTTP(w, r)

		log.Infof("HTTP Request Handling -- ENDS. Time Taken: %s",
			time.Since(start))
	})
}

//PolicyAuthorizationAppSessionGet Get
func PolicyAuthorizationAppSessionGet(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionGet -- Entered")
	vars := mux.Vars(r)
	log.Infof(" APP Session ID  : %s", vars["appSessionId"])

}

//PolicyAuthorizationAppSessionCreate Post
func PolicyAuthorizationAppSessionCreate(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionCreate -- Entered")
}

//PolicyAuthorizationAppSessionPatch Patch
func PolicyAuthorizationAppSessionPatch(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSession -- Entered")

}

//PolicyAuthorizationAppSessionDelete Delete
func PolicyAuthorizationAppSessionDelete(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionDelete -- Entered")

}

//PolicyAuthorizationAppSessionSubscribe Subscribe
func PolicyAuthorizationAppSessionSubscribe(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionSubscribe -- Entered")

}

//PolicyAuthorizationAppSessionUnsubscribe Unsubscribe
func PolicyAuthorizationAppSessionUnsubscribe(w http.ResponseWriter,
	r *http.Request) {
	log.Infoln("PolicyAuthorizationAppSessionUnsubscribe -- Entered")

}
