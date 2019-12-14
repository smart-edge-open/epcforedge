// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func getAllSubscriptions(cliCtx context.Context, afCtx *Context) (
	[]TrafficInfluSub, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tSubs, resp, err := cli.TrafficInfluSubGetAllAPI.SubscriptionsGetAll(
		cliCtx, afCtx.cfg.AfID)

	if err != nil {
		return nil, nil, err
	}
	return tSubs, resp, nil

}

//GetAllSubscriptions function
func GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		tsResp     []TrafficInfluSub
		resp       *http.Response
		tsRespJSON []byte
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
	tsResp, resp, err = getAllSubscriptions(cliCtx, afCtx)
	if err != nil {
		log.Errf("Traffic Influence Subscriptions get all : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tsRespJSON, err = json.Marshal(tsResp)
	if err != nil {
		log.Errf("Traffic Influence Subscriptions get all: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(tsRespJSON); err != nil {
		log.Errf("Traffic Influance Subscription get all: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
