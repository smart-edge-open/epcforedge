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

package main

import (
	"context"
	"net/http"
	"time"

	"github.com/otcshare/edgenode/pkg/config"
)

// Config : NEF Module Configuration Data Structure
type Config struct {
	Endpoint       string `json:"endpoint"`
	LocationPrefix string `json:"locationPrefix"`
	MaxSubSupport  int    `json:"maxSubSupport"`
	MaxAFSupport   int    `json:"maxAFSupport"`
	SubStartID     int    `json:"subStartID"`
}

// NEF Module Context Data Structure
type nefContext struct {
	cfg Config
	nef nefData
}

// runServer : This function cretaes a Router object and starts a HTTP Server
//             in a separate go routine. Also it listens for NEF module
//             running context cancellation event in another go routine. If
//             cancellation event occurs, it shutdowns the HTTP Server.
// Input Args:
//   - ctx:    NEF Module Running context
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
// Output Args:
//    - error: retruns no error. It only logs the error if any happened while
//             starting the HTTP Server
func runServer(ctx context.Context, nefCtx *nefContext) error {

	var err error

	/* NEFRouter obeject is created. After creation this object contains all
	 * the HTTP Service Handlers. These hanlders will be called when HTTP
	 * server receives any HTTP Request */
	nefRouter := NewNEFRouter(nefCtx)

	/* HTTP Server object is created */
	server := &http.Server{
		Addr:           nefCtx.cfg.Endpoint,
		Handler:        nefRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	stopServerCh := make(chan bool, 2)

	/* Go Routine is spawned here for listening for cancellation event on
	 * context */
	go func(stopServerCh chan bool) {
		<-ctx.Done()
		log.Info("Executing graceful stop")
		if err = server.Close(); err != nil {
			log.Errf("Could not close NEF server: %#v", err)
		}

		log.Info("NEF server stopped")
		stopServerCh <- true
	}(stopServerCh)

	/* Go Routine is spawned here for starting HTTP Server */
	go func(stopServerCh chan bool) {

		log.Infof("NEF listening on %s", server.Addr)
		if err = server.ListenAndServe(); err != nil {
			log.Errf("NEF server error: " + err.Error())
		}
		log.Info("Exiting")

		stopServerCh <- true
	}(stopServerCh)

	/* This self go routine is waiting for the receive events from the spawned
	 * go routines */
	<-stopServerCh
	<-stopServerCh

	return nil
}

// Run : This function reads the NEF Module configuration file and stores in
//       NEF Module Context. This also calls the Initialization/Creation of
//       NEF Data. Also it  calls runServer function for starting HTTP Server.
// Input Args:
//    - ctx:     NEF Module Running context
//    - cfgPath: This is NEF Module Configuration file path
// Output Args:
//     - error: retruns error in case any error occurred in reading NEF
//              configuration file or any error occurred in starting server
func Run(ctx context.Context, cfgPath string) error {

	var nefCtx nefContext

	/* Reads NEF Configuration file which is json format. Also it converts
	 * configuration data from json format to structure data */
	err := config.LoadJSONConfig(cfgPath, &nefCtx.cfg)
	if err != nil {
		log.Errf("Failed to load NEF configuration: %v", err)
		return err
	}

	/* Creates/Initializes NEF Data */
	nefCtx.nef.nefCreate()

	return runServer(ctx, &nefCtx)
}
