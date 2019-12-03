// Copyright 2019 Intel Corporation, Inc. All rights reserved
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

package ngcaf

import (
	"context"
	"encoding/json"
	"net/http"
)

func getSubscription(cliCtx context.Context, afCtx *AFContext,
	subscriptionID string) (TrafficInfluSub, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	ts, resp, err := cli.TrafficInfluSubGetAPI.SubscriptionGet(
		cliCtx, afCtx.cfg.AfID, subscriptionID)

	if err != nil {
		return TrafficInfluSub{}, nil, err
	}
	return ts, resp, nil
}

// GetSubscription function
func GetSubscription(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		tsResp         TrafficInfluSub
		resp           *http.Response
		subscriptionID string
		tsRespJSON     []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*AFContext)
	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	subscriptionID, err = getSubsIDFromURL(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tsResp, resp, err = getSubscription(cliCtx, afCtx, subscriptionID)
	if err != nil {
		log.Errf("Traffic Influence Subscription get : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tsRespJSON, err = json.Marshal(tsResp)
	if err != nil {
		log.Errf("Traffic Influence Subscription get : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(tsRespJSON); err != nil {
		log.Errf("Traffic Influance Subscription get %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
