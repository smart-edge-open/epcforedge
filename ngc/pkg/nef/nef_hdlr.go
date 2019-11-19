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
	//	"encoding/json"
	"errors"
	//	"io/ioutil"
	//	"net/http"
	"strconv"
	//	"github.com/gorilla/mux"
)

//NEF context data
type nefData struct {
	//nefport   string
	//location  string
	afcount int
	//subIdnum  int
	//maxSubSup int
	//maxAfSup  int
	afs map[string]*afData

	//Member functions
	//nefCreate
	//nefAddAf
	//GetAf
	//DeleteAf
	//Destroy
}

//AF Subscription data
type afSubscription struct {
	subid string
	loc   string
	ti    TrafficInfluSub
}


//AF data
type afData struct {
	afid       string
	subIdnum   int
	maxSubSupp int
	subs       map[string]*afSubscription
	//Member functions
	//afCreate
	//afAddSubscription
	//afGetSubscription
	//afUpdateSubscription
	//afDeleteSubscription
	//afDestroy
}

//Creates a AF instance
func (af *afData) afCreate(nefCtx *nefContext, afid string) error {

	//Validate afid ??

	af.afid = afid
	af.subIdnum = nefCtx.cfg.SubStartID //Start Number
	af.maxSubSupp = nefCtx.cfg.MaxSubSupport
	af.subs = make(map[string]*afSubscription)
	return nil
}

//Creates a new subscription
func (af *afData) afAddSubscription(nefCtx *nefContext, ti TrafficInfluSub) (loc string, err error) {

	if af.subIdnum >= nefCtx.cfg.MaxSubSupport+nefCtx.cfg.SubStartID {
		return "", errors.New("MAX SUBS Created")
	}

	//Generate a unique subscription ID string
	subIDStr := strconv.Itoa(af.subIdnum)
	af.subIdnum++
	log.Infoln(af.subIdnum)	

	//Create the Subscription info
	afsub := afSubscription{subid: subIDStr, ti: ti}

	//Create Location URI
	afsub.loc = /*TODO - Add Local IP:Port */ nefCtx.cfg.LocationPrefix + af.afid + "/subscriptions/" + subIDStr

	//Link the subscription with the AF
	af.subs[subIDStr] = &afsub

	log.Infoln(" SUBSCRIPTION ADDED ")
	log.Infoln(len(af.subs))
	//log.Infoln(af.subs)

	return afsub.loc, nil
}

func (af *afData) afUpdateSubscription(nefCtx *nefContext, subID string, ti TrafficInfluSub) (err error) {

	sub, ok := af.subs[subID]
	if ok == false {
		sub.ti = ti
		return errors.New("Subscription Not Found")
	}
	sub.ti = ti
	return
}

func (af *afData) afGetSubscription(nefCtx *nefContext, subID string) (ti TrafficInfluSub, err error) {

	_, ok := af.subs[subID]

	if ok == true {
		return af.subs[subID].ti, nil
	}

	return ti, errors.New("SubscriptionId Not found")

}

func (af *afData) afGetSubscriptionList(nefCtx *nefContext) (subslist []TrafficInfluSub, err error) {

	if len(af.subs) > 0 {
		for _, value := range af.subs {
			subslist = append(subslist, value.ti)
		}
		return subslist, nil
	}

	return nil, errors.New("No Subscriptions present")
}

func (af *afData) afDeleteSubscription(nefCtx *nefContext, subID string) error {
	//Check if AF is already present
	_, ok := af.subs[subID]

	if ok == true {
		delete(af.subs, subID)
		af.subIdnum--
		return nil
	}

	return errors.New("SubscriotionId not found")

}
func (af *afData) afDestroy(afid string) error {

	//Todo delete all subscriptions, needed in go ??
	//Needed for gracefully disconnecting
	return errors.New("AF data cleaned")
}

//Initialize the NEF component
func (nef *nefData) nefCreate() error {

	nef.afcount = 0
	nef.afs = make(map[string]*afData)

	return nil
}

/*
func NEFInit() error {

	return nef.nefCreate()
}*/

func (nef *nefData) nefAddAf(nefCtx *nefContext, afID string) (af *afData, err error) {

	var afe afData

	//Check if AF is already present
	_, ok := nef.afs[afID]

	if ok == true {
		return nef.afs[afID], errors.New("AF already present")
	}

	//Create a new entry of AF

	afe.afCreate(nefCtx, afID)
	nef.afs[afID] = &afe
	nef.afcount++

	return &afe, nil
}

func (nef *nefData) nefGetAf(afID string) (af *afData, err error) {

	//Check if AF is already present
	afe, ok := nef.afs[afID]

	if ok == true {
		return afe, nil
	}
	err = errors.New("AF entry not present")
	return afe, err
}

func (nef *nefData) nefDeleteAf(afID string) (err error) {

	//Check if AF is already present
	_, ok := nef.afs[afID]

	if ok == true {
		delete(nef.afs, afID)
		nef.afcount--
		return nil
	}

	err = errors.New("AF entry not present")
	return err
}

func (nef *nefData) nefDestroy() {

	// Todo
}
