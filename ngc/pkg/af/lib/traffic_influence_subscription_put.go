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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func modifySubscriptionByPut(ts TrafficInfluSub, afCtx *afContext,
	subscriptionId string, cliCtx context.Context) (TrafficInfluSub,
	*http.Response, error) {

	cliCfg := NewConfiguration()
	cli := NewAFClient(cliCfg)

	tsResp, resp, err := cli.TrafficInfluSubPutApi.SubscriptionPut(cliCtx,
		afCtx.cfg.AfId, subscriptionId, ts)

	if err != nil {

		log.Errf("AF Traffic Influance Subscription PUT: %s", err.Error())
		return TrafficInfluSub{}, nil, err
	}
	return tsResp, resp, nil
}

func ModifySubscriptionPut(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		ts             TrafficInfluSub
		tsResp         TrafficInfluSub
		resp           *http.Response
		subscriptionId string
		transId        int
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
		log.Errf("Traffic Influance Subscription PUT: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transId, err = genTransactionId(afCtx)
	if err != nil {

		log.Errf("Traffic Influance Subscription PUT %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts.AfTransId = strconv.Itoa(transId)
	fmt.Printf("TransID: %s, %d. ", ts.AfTransId, transId)
	subscriptionId, err = getSubsIdFromUrl(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription PUT: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	afCtx.transactions[transId] = TrafficInfluSub{}
	tsResp, resp, err = modifySubscriptionByPut(ts, afCtx, subscriptionId,
		cliCtx)
	if err != nil {
		log.Errf("Traffic Influence Subscription create : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return

	} else {

		if ts.SubscribedEvents == nil {
			delete(afCtx.transactions, transId)
		} else {
			afCtx.transactions[transId] = tsResp
		}

		if _, ok :=
			afCtx.subscriptions[subscriptionId][strconv.Itoa(transId)]; ok {
			delete(afCtx.transactions, transId)
		}

		if afCtx.subscriptions[subscriptionId] == nil {
			afCtx.subscriptions[subscriptionId] =
				make(map[string]TrafficInfluSub)
		}

		afCtx.subscriptions[subscriptionId] =
			map[string]TrafficInfluSub{ts.AfTransId: tsResp}

		if resp != nil {
			w.WriteHeader(resp.StatusCode)
		}
	}
}
