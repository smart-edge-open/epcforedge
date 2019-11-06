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
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func deleteSubscription(afCtx *afContext, subscriptionID string,
	cliCtx context.Context) (*http.Response, error) {

	cliCfg := NewConfiguration()
	cli := NewClient(cliCfg)

	resp, err := cli.TrafficInfluSubDeleteAPI.SubscriptionDelete(cliCtx,
		afCtx.cfg.AfId, subscriptionID)

	if err != nil {

		log.Errf("AF Traffic Influance Subscription DELETE: %s", err.Error())

		return nil, err
	}
	return resp, nil

}

// DeleteSubscription function
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		resp           *http.Response
		subscriptionID string
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

	subscriptionID, err = getSubsIDFromURL(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription PUT: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err = deleteSubscription(afCtx, subscriptionID, cliCtx)
	if err != nil {
		log.Errf("Traffic Influence Subscription DELETE : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if interMap, ok := afCtx.subscriptions[subscriptionID]; ok {

		for transID := range interMap {
			var i int
			if i, err = strconv.Atoi(transID); err == nil {
				delete(afCtx.transactions, i)
			}
		}
		delete(afCtx.subscriptions, subscriptionID)
	}

	w.WriteHeader(resp.StatusCode)
}
