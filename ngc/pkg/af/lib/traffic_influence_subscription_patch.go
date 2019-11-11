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
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func modifySubscriptionByPatch(cliCtx context.Context, ts TrafficInfluSubPatch,
	afCtx *afContext, subscriptionID string) (TrafficInfluSub,
	*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsResp, resp, err := cli.TrafficInfluSubPatchAPI.SubscriptionPatch(cliCtx,
		afCtx.cfg.AfID, subscriptionID, ts)

	if err != nil {

		log.Errf("AF Traffic Influance Subscription PUT: %s", err.Error())
		return TrafficInfluSub{}, nil, err
	}
	return tsResp, resp, nil
}

// ModifySubscriptionPatch function
func ModifySubscriptionPatch(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		tsPatch        TrafficInfluSubPatch
		tsResp         TrafficInfluSub
		resp           *http.Response
		subscriptionID string
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*afContext)
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

	subscriptionID, err = getSubsIDFromURL(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription PATCH: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tsResp, resp, err = modifySubscriptionByPatch(cliCtx, tsPatch, afCtx,
		subscriptionID)
	if err != nil {
		log.Errf("Traffic Influence Subscription PATCH : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if interMap, ok := afCtx.subscriptions[subscriptionID]; ok {

		for transID := range interMap { //there's only one entry in the map
			afCtx.subscriptions[subscriptionID][(transID)] = tsResp
		}

	} else {

		log.Info("Traffic Influence Subscription: "+
			"subscriptionID %s not found in local memory", subscriptionID)
	}
	w.WriteHeader(resp.StatusCode)
}
