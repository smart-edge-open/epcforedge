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
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func createSubscription(ts TrafficInfluSub, afCtx *afContext,
	cliCtx context.Context) (TrafficInfluSub, *http.Response, error) {

	cliCfg := NewConfiguration()
	cli := NewClient(cliCfg)

	tsResp, resp, err := cli.TrafficInfluSubPostAPI.SubscriptionPost(cliCtx,
		afCtx.cfg.AfId, ts)

	if err != nil {

		log.Errf("AF Traffic Influance Subscription POST: %s", err.Error())
		return TrafficInfluSub{}, nil, err
	}
	return tsResp, resp, nil
}

// CreateSubscription function
func CreateSubscription(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		ts             TrafficInfluSub
		tsResp         TrafficInfluSub
		resp           *http.Response
		url            *url.URL
		subscriptionID string
		transID        int
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

	if err = json.NewDecoder(r.Body).Decode(&ts); err != nil {
		log.Errf("Traffic Influance Subscription create: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transID, err = genTransactionID(afCtx)
	if err != nil {

		log.Errf("Traffic Influance Subscription create %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//fmt.Printf("TransID: %s, %d. ", ts.AfTransID, afTransId )
	//store transaction ID to a list of currently used transaction IDs
	afCtx.transactions[transID] = TrafficInfluSub{}

	ts.AfTransID = strconv.Itoa(transID)
	tsResp, resp, err = createSubscription(ts, afCtx, cliCtx)
	if err != nil {
		log.Errf("Traffic Influence Subscription create : %s", err.Error())
		delete(afCtx.transactions, transID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if url, err = resp.Location(); err != nil {
		log.Errf("Traffic Influence Subscription create: %s", err.Error())
		delete(afCtx.transactions, transID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if subscriptionID, err = getSubsIDFromURL(url); err != nil {
		delete(afCtx.transactions, transID)
		log.Errf("Traffic Influence Subscription create: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url.String())

	if len(tsResp.SubscribedEvents) == 0 {
		delete(afCtx.transactions, transID)
	} else {
		afCtx.transactions[transID] = tsResp
		log.Infof("Saving subscription ID : %s to local memory.",
			subscriptionID)
		afCtx.subscriptions[subscriptionID] =
			map[string]TrafficInfluSub{strconv.Itoa(transID): afCtx.transactions[transID]}

	}
	w.WriteHeader(resp.StatusCode)
}
