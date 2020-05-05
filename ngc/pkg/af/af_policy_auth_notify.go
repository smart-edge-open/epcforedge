// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"encoding/json"
	"net/http"
)

// PolicyAuthEventNotify Event notification handler
func PolicyAuthEventNotify(w http.ResponseWriter, r *http.Request) {

	var (
		err        error
		eventNotif EventsNotification
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in PolicyAuthEventNotify",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewDecoder(r.Body).Decode(&eventNotif)
	if err != nil {
		logPolicyRespErr(&w, "Json Decode error in "+
			"PolicyAuthEventNotify: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
}

// PolicyAuthEventNotifTerminate Event notification termination handler
func PolicyAuthEventNotifTerminate(w http.ResponseWriter, r *http.Request) {

	var (
		err      error
		termInfo TerminationInfo
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in PolicyAuthEventNotifTermin",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewDecoder(r.Body).Decode(&termInfo)
	if err != nil {
		logPolicyRespErr(&w, "Json Decode error in "+
			"PolicyAuthEventNotifTerminate: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
}
