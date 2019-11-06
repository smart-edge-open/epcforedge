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

func getSubscription(afCtx *afContext, subscriptionId string,
	cliCtx context.Context) (TrafficInfluSub, *http.Response, error) {

	cliCfg := NewConfiguration()
	cli := NewAFClient(cliCfg)

	tsResp, resp, err := cli.TrafficInfluSubGetApi.SubscriptionGet(
		cliCtx, afCtx.cfg.AfId, subscriptionId)

	if err != nil {

		log.Errf("AF Traffic Influance Subscription GET: %s", err.Error())
		return TrafficInfluSub{}, nil, err
	}
	return tsResp, resp, nil
}

func GetSubscription(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		tsResp         TrafficInfluSub
		resp           *http.Response
		subscriptionId string
		transId        int
		tsRespJson     []byte
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

	transId, err = genTransactionId(afCtx)
	if err != nil {

		log.Errf("Traffic Influance Subscription PUT %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subscriptionId, err = getSubsIdFromUrl(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription PUT: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	afCtx.transactions[transId] = TrafficInfluSub{}
	tsResp, resp, err = getSubscription(afCtx, subscriptionId, cliCtx)
	if err != nil {
		log.Errf("Traffic Influence Subscription create : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return

	} else {

		if resp != nil {
			w.WriteHeader(resp.StatusCode)
		}

		tsRespJson, err = json.Marshal(tsResp)
		w.Write(tsRespJson)
	}
}
