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

package af

import (
	"context"
	"encoding/json"
	//"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func modifySubscriptionByPatch(ts TrafficInfluSubPatch, afCtx *afContext,
	subscriptionId string, cliCtx context.Context) (TrafficInfluSub,
	*http.Response, error) {

	cliCfg := NewConfiguration()
	cli := NewAFClient(cliCfg)

	tsResp, resp, err := cli.TrafficInfluSubPatchApi.SubscriptionPatch(cliCtx,
		afCtx.cfg.AfId, subscriptionId, ts)

	if err != nil {

		log.Errf("AF Traffic Influance Subscription PUT: %s", err.Error())
		return TrafficInfluSub{}, nil, err
	}
	return tsResp, resp, nil
}

func ModifySubscriptionPatch(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		tsPatch        TrafficInfluSubPatch
		tsResp         TrafficInfluSub
		resp           *http.Response
		subscriptionId string
	)

	afCtx := r.Context().Value(string("af-ctx")).(*afContext)
	cliCtx, cancel := context.WithCancel(context.Background())

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-osSignals
		log.Infof("Received signal: %#v", sig)
		cancel()
	}()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&tsPatch); err != nil {
		log.Errf("Traffic Influance Subscription PATCH: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subscriptionId, err = getSubsIdFromUrl(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription PATCH: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tsResp, resp, err = modifySubscriptionByPatch(tsPatch, afCtx,
		subscriptionId, cliCtx)
	if err != nil {
		log.Errf("Traffic Influence Subscription PATCH : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return

	} else {

		//fmt.Printf("Patch subscriptions %d, %d, s%", len(afCtx.subscriptions),
		//len(afCtx.subscriptions[subscriptionId]))

		if interMap, ok := afCtx.subscriptions[subscriptionId]; ok {

			for transId := range interMap { //there's only one entry in the map
				//fmt.Println("TransId:", transId, "Value:", ts)
				afCtx.subscriptions[subscriptionId][(transId)] = tsResp
			}
		} else {

			log.Info("Traffic Influence Subscription: "+
				"subscriptionId %s not found in local memory", subscriptionId)
		}
		//fmt.Println(afCtx.subscriptions)
		if resp != nil {
			w.WriteHeader(resp.StatusCode)
		}
	}
}
