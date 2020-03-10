/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngcnef

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	oauth2 "github.com/otcshare/epcforedge/ngc/pkg/oauth2"
)

// Route : Structure which describes HTTP Request Handler type and other
//         details like name, method, path, etc
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc /*    */
}

// NEFRoutes : NEF Routes lists which contains Routes with different HTTP
//             Request handlers for NEF
var NEFRoutes = []Route{
	// Traffic Influence Routes
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
	// PFD Management Routes
	{
		"ReadAllPFDManagementTransaction",
		strings.ToUpper("Get"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions",
		ReadAllPFDManagementTransaction,
	},

	{
		"CreatePFDManagementTransaction",
		strings.ToUpper("Post"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions",
		CreatePFDManagementTransaction,
	},

	{
		"ReadPFDManagementTransaction",
		strings.ToUpper("Get"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions/{transactionId}",
		ReadPFDManagementTransaction,
	},

	{
		"UpdatePutPFDManagementTransaction",
		strings.ToUpper("Put"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions/{transactionId}",
		UpdatePutPFDManagementTransaction,
	},

	{
		"DeletePFDManagementTransaction",
		strings.ToUpper("Delete"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions/{transactionId}",
		DeletePFDManagementTransaction,
	},

	{
		"ReadPFDManagementApplication",
		strings.ToUpper("Get"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions/{transactionId}/" +
			"applications/{appId}",
		ReadPFDManagementApplication,
	},

	{
		"DeletePFDManagementApplication",
		strings.ToUpper("Delete"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions/{transactionId}/" +
			"applications/{appId}",
		DeletePFDManagementApplication,
	},

	{
		"UpdatePutPFDManagementApplication",
		strings.ToUpper("Put"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions/{transactionId}/" +
			"applications/{appId}",
		UpdatePutPFDManagementApplication,
	},

	{
		"PatchPFDManagementApplication",
		strings.ToUpper("Patch"),
		"/3gpp-pfd-management/v1/{scsAsId}/transactions/{transactionId}/" +
			"applications/{appId}",
		PatchPFDManagementApplication,
	},
}

type nefCtxKey string

// NewNEFRouter : This function creates and initializes a NEF Router with all
//                the available routes for NEF Module. This router object is
//                defined in "github.com/gorilla/mux" package.
//  Input Args:
//     - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
//  Output Args:
//     - error: retruns pointer to created mux.Router object
func NewNEFRouter(nefCtx *nefContext) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// smf upf notification route
	smfNotif := Route{}
	smfNotif.Name = "NotifySmfUPFEvent"
	smfNotif.Method = strings.ToUpper("Post")
	smfNotif.Handler = NotifySmfUPFEvent
	smfNotif.Pattern = nefCtx.cfg.UpfNotificationResURIPath
	NEFRoutes = append(NEFRoutes, smfNotif)

	for _, route := range NEFRoutes {

		var handler http.Handler = route.Handler
		handler = nefRouteLogger(handler, route.Name)

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
				nefCtxKey("nefCtx"),
				nefCtx)

			if nefCtx.cfg.OAuth2Support {
				if nefValidateAccessToken(w, r) {
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

func nefValidateAccessToken(w http.ResponseWriter, r *http.Request) bool {

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

// nefRouteLogger : This function logs data received in HTTP request.
// Input Args:
//    - httpHandler: This is HTTP handler function pointer for HTTP Request
//                   Received
//    - name: This is route name.
// Output Args:
//    - httpHandler: This is HTTP handler function pointer for HTTP Request
//                   Received. This HTTP Handler is actually the updated HTTP
//                   Handler. The updated HTTP Handler now can logging of HTTP
//                   Request Info also
func nefRouteLogger(httpHandler http.Handler, name string) http.Handler {
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
