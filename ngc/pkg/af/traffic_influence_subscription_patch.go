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

func modifySubscriptionByPatch(cliCtx context.Context, ts TrafficInfluSubPatch,
	afCtx *AFContext, sID string) (TrafficInfluSub,
	*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsRet, resp, err := cli.TrafficInfluSubPatchAPI.SubscriptionPatch(cliCtx,
		afCtx.cfg.AfID, sID, ts)

	if err != nil {
		return TrafficInfluSub{}, nil, err
	}
	return tsRet, resp, nil
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

	afCtx := r.Context().Value(keyType("af-ctx")).(*AFContext)
	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&tsPatch); err != nil {
		log.Errf("Traffic Influance Subscription modify: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subscriptionID, err = getSubsIDFromURL(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription modify: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tsResp, resp, err = modifySubscriptionByPatch(cliCtx, tsPatch, afCtx,
		subscriptionID)
	if err != nil {
		log.Errf("Traffic Influence Subscription modify : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if interMap, ok := afCtx.subscriptions[subscriptionID]; ok {

		for transID := range interMap {
			afCtx.subscriptions[subscriptionID][(transID)] = tsResp
		}

	} else {

		log.Info("Traffic Influence Subscription: "+
			"subscriptionID %s not found in local memory", subscriptionID)
	}
	w.WriteHeader(resp.StatusCode)
}
