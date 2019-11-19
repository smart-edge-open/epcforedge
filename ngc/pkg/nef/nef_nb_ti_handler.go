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

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	//"strconv"

	"github.com/gorilla/mux"
)


func createNewSub(nefCtx *nefContext, afID string, ti TrafficInfluSub) (loc string, err error) {

	var af *afData

	nef := &nefCtx.nef

	//Validate the Traffic Influence
	err = validateTIS(nefCtx, ti)
	if err != nil {
		log.Infoln(err)
		return "", err
	}

	af, err = nef.nefGetAf(afID)

	if err != nil {
		log.Infoln("NO AF PRESENT CREATE AF")
		af, _ = nef.nefAddAf(nefCtx, afID)
	} else {
		log.Infoln("AF PRESENT AF")
		log.Infoln(af)
	}

	loc, err = af.afAddSubscription(nefCtx, ti)

	if err != nil {
		log.Infoln(err)
		return loc, err
	}

	//log.Infoln(nef)
	//log.Infof("AF COUNT: %+v", nef.afcount)

	return loc, nil
}

//Validate the Traffic influence data received from AF
func validateTIS(nefCtx *nefContext, ti TrafficInfluSub) (err error) {

	nef := &nefCtx.nef
	//Check if we have crossed max supported AF
	if nef.afcount >= nefCtx.cfg.MaxAFSupport {
		log.Infoln("MAX AF exceeded ")
		return errors.New("MAX AF exceeded")
		//return err
	}
	return nil
}

// ReadAllTrafficInfluenceSubscription : API to read all the subscritions
func ReadAllTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(string("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])

	af, err := nef.nefGetAf(vars["afId"])

	if err != nil {
		log.Infof("Error: AF ID not found ")
		log.Infoln(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}

	subslist, _ := af.afGetSubscriptionList(nefCtx)

	mdata, err2 := json.Marshal(subslist)
	if err2 != nil {
		log.Infof("Error:  Failed to marshal the json data")
		log.Infoln(err2)
		//panic(err2)
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(mdata)

	w.WriteHeader(http.StatusOK)
	log.Infof("HTTP Response sent: %d", http.StatusOK)
}

// CreateTrafficInfluenceSubscription : Handles the traffic influence requested
// by AF
func CreateTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(string("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])

	b, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}

	//Traffic Influence data
	TrInBody := TrafficInfluSub{}

	//Convert the json Traffic Influence data into struct
	err1 := json.Unmarshal(b, &TrInBody)

	//Print
	//log.Infof("\n Traffic Influence data from AF\n%+v\n\n", TrInBody)

	if err1 != nil {
		log.Infof("Error: Failed to UNmarshal POST req ")
		log.Infoln(err1)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}
	//Update the data respose of POST

	//Martshal data and send into the body
	mdata, err2 := json.Marshal(TrInBody)

	if err2 != nil {
		log.Infof("Error:  Failed to marshal the json data")
		log.Infoln(err2)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
	}

	//loc, err3 := createNewSubscription(vars["afId"], TrInBody)
	loc, err3 := createNewSub(nefCtx, vars["afId"], TrInBody)
	log.Infoln(loc)

	logNef(nef)

	if err3 != nil {
		log.Infof("Error:  Failed to Create AF data")
		log.Infoln(err3)
		//panic(err3)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Location", loc)

	w.Write(mdata)

	w.WriteHeader(http.StatusOK)
	log.Infof("HTTP Response sent: %d", http.StatusOK)
}

// ReadTrafficInfluenceSubscription : Read a particular subscription details
func ReadTrafficInfluenceSubscription(w http.ResponseWriter, r *http.Request) {

	nefCtx := r.Context().Value(string("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, ok := nef.nefGetAf(vars["afId"])

	if ok != nil {
		log.Infoln(ok)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}
	log.Infoln("AF Found ")

	substi, err := af.afGetSubscription(nefCtx, vars["subscriptionId"])

	if err != nil {
		log.Infoln(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}

	mdata, err2 := json.Marshal(substi)
	if err2 != nil {
		log.Infof("Error:  Failed to marshal the json data")
		log.Infoln(err2)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(mdata)

	w.WriteHeader(http.StatusOK)
	log.Infof("HTTP Response sent: %d", http.StatusOK)
}

// UpdatePutTrafficInfluenceSubscription : Updates a traffic influence created
// earlier (PUT Req)
func UpdatePutTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(string("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, ok := nef.nefGetAf(vars["afId"])
	if ok == nil {

		b, err := ioutil.ReadAll(r.Body)

		defer r.Body.Close()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Traffic Influence data
		TrInBody := TrafficInfluSub{}

		//Convert the json Traffic Influence data into struct
		err1 := json.Unmarshal(b, &TrInBody)

		if err1 != nil {
			log.Infof("Error: Failed to UNmarshal POST req ")
			log.Infoln(err1)
			w.WriteHeader(http.StatusInternalServerError)
			log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
			return
		}

		ok := af.afUpdateSubscription(nefCtx, vars["subscriptionId"], TrInBody)
		if ok != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
			return
		}

		mdata, err2 := json.Marshal(TrInBody)

		if err2 != nil {
			log.Infof("Error:  Failed to marshal the json data")
			log.Infoln(err2)
			w.WriteHeader(http.StatusInternalServerError)
			log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.Write(mdata)

	}

	w.WriteHeader(http.StatusOK)
	log.Infof("HTTP Response sent: %d", http.StatusOK)
}

// UpdatePatchTrafficInfluenceSubscription : Updates a traffic influence created
//  earlier (PATCH Req)
func UpdatePatchTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(string("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, ok := nef.nefGetAf(vars["afId"])
	if ok == nil {

		b, err := ioutil.ReadAll(r.Body)

		defer r.Body.Close()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Traffic Influence data
		TrInBody := TrafficInfluSub{}

		//Convert the json Traffic Influence data into struct
		err1 := json.Unmarshal(b, &TrInBody)

		if err1 != nil {
			log.Infof("Error: Failed to UNmarshal POST req ")
			log.Infoln(err1)
			w.WriteHeader(http.StatusInternalServerError)
			log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
			return
		}

		ok := af.afUpdateSubscription(nefCtx, vars["subscriptionId"], TrInBody)
		if ok != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
			return
		}

		mdata, err2 := json.Marshal(TrInBody)

		if err2 != nil {
			log.Infof("Error:  Failed to marshal the json data")
			log.Infoln(err2)
			w.WriteHeader(http.StatusInternalServerError)
			log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.Write(mdata)

	}

	w.WriteHeader(http.StatusOK)
	log.Infof("HTTP Response sent: %d", http.StatusOK)
}

// DeleteTrafficInfluenceSubscription : Deletes a traffic influence created by
//  AF
func DeleteTrafficInfluenceSubscription(w http.ResponseWriter,
	r *http.Request) {

	nefCtx := r.Context().Value(string("nefCtx")).(*nefContext)
	nef := &nefCtx.nef

	vars := mux.Vars(r)
	log.Infof(" AFID  : %s", vars["afId"])
	log.Infof(" SUBSCRIPTION ID  : %s", vars["subscriptionId"])

	af, err := nef.nefGetAf(vars["afId"])

	if err != nil {
		log.Infoln(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}
	err = af.afDeleteSubscription(nefCtx, vars["subscriptionId"])
	if err != nil {
		log.Infoln(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("HTTP Response sent: %d", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Infof("HTTP Response sent: %d", http.StatusOK)

	logNef(nef)
}

func logNef(nef *nefData) {

	log.Infof("AF Count %+v", len(nef.afs))
	if len(nef.afs) > 0 {
		for key, value := range nef.afs {
			log.Infof(" AFKey : %+v, valAF_Id : %+v", key, value.afid)

			log.Infof("SUB Count for AF: %+v is %+v", value.afid, len(value.subs))
			for ks, vs := range value.subs {
				log.Infof("   SubKey : %+v, valSub_Id : %+v, ServiceId: %+v", ks, vs.subid, vs.ti.AfServiceID)
			}
		}
	}

}
