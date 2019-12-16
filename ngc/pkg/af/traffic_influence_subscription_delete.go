// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"net/http"
	"strconv"
)

func deleteSubscription(cliCtx context.Context, afCtx *Context,
	sID string) (*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	resp, err := cli.TrafficInfluSubDeleteAPI.SubscriptionDelete(cliCtx,
		afCtx.cfg.AfID, sID)

	if err != nil {
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

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Traffic Influance Subscription delete: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if afCtx.subscriptions == nil {
		log.Errf("AF context subscriptions map has not been initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if afCtx.transactions == nil {
		log.Errf("AF context  transactions map has not been initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	subscriptionID, err = getSubsIDFromURL(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription delete: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err = deleteSubscription(cliCtx, afCtx, subscriptionID)
	if err != nil {
		log.Errf("Traffic Influence Subscription delete: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if interMap, ok := afCtx.subscriptions[subscriptionID]; ok {
		for transID := range interMap {
			var i int
			if i, err = strconv.Atoi(transID); err != nil {
				log.Errf("Error converting transID to integer: %v", err)
			} else {
				delete(afCtx.transactions, i)
				log.Infof("Deleted transaction ID %v", i)
			}
		}
		delete(afCtx.subscriptions, subscriptionID)
	}

	w.WriteHeader(resp.StatusCode)
}
