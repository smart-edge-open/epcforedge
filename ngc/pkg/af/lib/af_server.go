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
	"net/http"

	logger "github.com/otcshare/common/log"
	"github.com/otcshare/edgenode/pkg/config"
)

// TransactionIDs type
type TransactionIDs map[int]TrafficInfluSub

// NotifSubscryptions type
type NotifSubscryptions map[string]map[string]TrafficInfluSub

//Config structure
type Config struct {
	Endpoint  string `json:"Endpoint"`
	AfID      string `json:"AfId"`
	BasePath  string `json:"BasePath"`
	UserAgent string `json:"UserAgent"`
}

type afContext struct {
	subscriptions NotifSubscryptions
	transactions  TransactionIDs
	cfg           Config
}

var log = logger.DefaultLogger.WithField("af-ctx", nil)

func runServer(ctx context.Context, afCtx *afContext) error {

	var err error

	afCtx.transactions = make(TransactionIDs)
	afCtx.subscriptions = make(NotifSubscryptions)
	afRouter := NewAFRouter(afCtx)
	server := &http.Server{
		Addr:    afCtx.cfg.Endpoint,
		Handler: afRouter,
	}
	stopServerCh := make(chan bool, 2)

	go func(stopServerCh chan bool) {
		<-ctx.Done()
		log.Info("Executing graceful stop")
		if err = server.Close(); err != nil {
			log.Errf("Could not close AF server: %#v", err)
		}

		log.Info("AF server stopped")
		stopServerCh <- true
	}(stopServerCh)

	go func(stopServerCh chan bool) {
		log.Infof("Serving AF on: %s", afCtx.cfg.Endpoint)
		if err = server.ListenAndServe(); err != nil {
			log.Info("AF server error: " + err.Error())
		}
		log.Errf("Stopped AF serving")
		stopServerCh <- true
	}(stopServerCh)

	<-stopServerCh
	<-stopServerCh
	return nil
}

// Run function
func Run(parentCtx context.Context, cfgPath string) error {

	var afCtx afContext

	// load AF configuration from file
	err := config.LoadJSONConfig(cfgPath, &afCtx.cfg)

	if err != nil {
		log.Errf("Failed to load AF configuration: %v", err)
		return err
	}

	return runServer(parentCtx, &afCtx)

}
