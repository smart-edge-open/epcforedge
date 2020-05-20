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

	err = json.NewDecoder(r.Body).Decode(&termInfo)
	if err != nil {
		logPolicyRespErr(&w, "Json Decode error in "+
			"PolicyAuthEventNotifTerminate: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	err = validateTermInfo(&termInfo)
	if err != nil {
		logPolicyRespErr(&w, "PolicyAuthEventNotifTerminate:"+
			err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
}

// PolicyAuthSMFNotify Handler for SMF UP_PATH_CH notifications
func PolicyAuthSMFNotify(w http.ResponseWriter,
	r *http.Request) {

	var (
		smfEv        NsmfEventExposureNotification
		nsmEvNo      NsmfEventNotification
		upEventFound bool
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		logPolicyRespErr(&w, "nil afCtx in PolicyAuthSMFNotify",
			http.StatusInternalServerError)
		return
	}

	if r.Body == nil {
		log.Errf("PolicyAuthSMFNotify Empty Body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Retrieve the event notification information from the request
	if err := json.NewDecoder(r.Body).Decode(&smfEv); err != nil {
		logPolicyRespErr(&w, "Json Decode error in "+
			"PolicyAuthSMFNotify:",
			http.StatusBadRequest)
		return
	}

	// Validate the content of the NsmfEventExposureNotification
	// Check if notification id is present
	if smfEv.NotifID == "" {
		logPolicyRespErr(&w, "Missing notif id in "+
			"PolicyAuthSMFNotify: ",
			http.StatusBadRequest)
		return
	}

	// Check if notification events with UP_PATH_CH is present
	if len(smfEv.EventNotifs) == 0 {
		logPolicyRespErr(&w, "Missing event notifs in "+
			"PolicyAuthSMFNotify: ",
			http.StatusBadRequest)
		return
	}

	for _, nsmEvNo = range smfEv.EventNotifs {
		if nsmEvNo.Event == "UP_PATH_CH" {
			log.Infof("PolicyAuthSMFNotify found an entry for UP_PATH_CH")
			upEventFound = true
			break
		}

	}

	if upEventFound {

		sendUpPathEventNotification(smfEv.NotifID, afCtx, nsmEvNo)
	}

	w.WriteHeader(http.StatusNoContent)

}
