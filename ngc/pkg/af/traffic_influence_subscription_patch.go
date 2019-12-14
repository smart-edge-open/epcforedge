// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func modifySubscriptionByPatch(cliCtx context.Context, ts TrafficInfluSubPatch,
	afCtx *Context, sID string) (TrafficInfluSub,
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

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Traffic Influance Subscription create: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	}
	w.WriteHeader(resp.StatusCode)
}
