// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func getSubscription(cliCtx context.Context, afCtx *Context,
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

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Traffic Influance Subscription get: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		log.Errf("Traffic Influance Subscription get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
