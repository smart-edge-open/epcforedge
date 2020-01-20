// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func createSubscription(cliCtx context.Context, ts TrafficInfluSub,
	afCtx *Context) (TrafficInfluSub, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsResp, resp, err := cli.TrafficInfluSubPostAPI.SubscriptionPost(cliCtx,
		afCtx.cfg.AfID, ts)

	if err != nil {
		return TrafficInfluSub{}, resp, err
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

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Traffic Influance Subscription create: " +
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
		log.Errf("AF context transactions map been initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	//store transaction ID to a list of currently used transaction IDs
	afCtx.transactions[transID] = TrafficInfluSub{}
	log.Infof("Saving transaction ID %d", transID)

	ts.AFTransID = strconv.Itoa(transID)
	if len(tsResp.SubscribedEvents) == 0 {
		ts.Self = Link("https://" + afCtx.cfg.SrvCfg.Hostname +
			afCtx.cfg.SrvCfg.NotifPort + DefaultNotifURL)
	}
	tsResp, resp, err = createSubscription(cliCtx, ts, afCtx)
	if err != nil {
		log.Errf("Traffic Influence Subscription create : %s", err.Error())
		delete(afCtx.transactions, transID)
		log.Infof("Deleted transaction ID %v", transID)
		w.WriteHeader(getStatusCode(resp))
		return
	}

	if url, err = resp.Location(); err != nil {
		log.Errf("Traffic Influence Subscription create: %s", err.Error())
		delete(afCtx.transactions, transID)
		log.Infof("Deleted transaction ID %v", transID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if subscriptionID, err = getSubsIDFromURL(url); err != nil {
		delete(afCtx.transactions, transID)
		log.Infof("Deleted transaction ID %v", transID)
		log.Errf("Traffic Influence Subscription create: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url.String())
	if len(tsResp.SubscribedEvents) == 0 {
		//keep in memory only subscriptions to notifications
		delete(afCtx.transactions, transID)
		log.Infof("Deleted transaction ID %v", transID)
	} else {
		afCtx.transactions[transID] = tsResp
		log.Infof("Saving subscription %s.",
			subscriptionID)
		afCtx.subscriptions[subscriptionID] =
			map[string]TrafficInfluSub{
				strconv.Itoa(transID): afCtx.transactions[transID]}

	}
	w.WriteHeader(resp.StatusCode)
}
