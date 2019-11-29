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

package ngcnef

import (
	"context"
	"encoding/json"
	logtool "github.com/otcshare/common/log"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

// Log handler initialized. This is to be used throughout the nef module for
// logging
var log = logtool.DefaultLogger.WithField("NEF", nil)

//HTTPConfig contains the configuration for the HTTP 1.1
type HTTPConfig struct {
	Endpoint string `json:"endpoint"`
}

//HTTP2Config Contains the configuration for the HTTP2
type HTTP2Config struct {
	Endpoint      string `json:"endpoint"`
	NefServerCert string `json:"NefServerCert"`
	NefServerKey  string `json:"NefServerKey"`
	AfServerCert  string `json:"AfServerCert"`
}

// Config contains NEF Module Configuration Data Structure
type Config struct {
	// API Root for the NEF
	NefAPIRoot                string `json:"nefAPIRoot`
	LocationPrefix            string `json:"locationPrefix"`
	MaxSubSupport             int    `json:"maxSubSupport"`
	MaxAFSupport              int    `json:"maxAFSupport"`
	SubStartID                int    `json:"subStartID"`
	UpfNotificationResURIPath string `json:"UpfNotificationResUriPath"`
	UserAgent                 string `json:"UserAgent"`
	HTTPConfig                HTTPConfig
	HTTP2Config               HTTP2Config
	AfServiceIDs              []interface{} `json:"afServiceIDs"`
}

// NEF Module Context Data Structure
type nefContext struct {
	cfg Config
	nef nefData
}

/* Go Routine is spawned here for starting HTTP Server */
func startHTTPServer(server *http.Server,
	stopServerCh chan bool) {
	if server != nil {
		log.Infof("HTTP 1.1 listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Errf("HTTP server error: " + err.Error())
		}
	}
	stopServerCh <- true
}

/* Go Routine is spawned here for starting HTTP-2 Server */
func startHTTP2Server(serverHTTP2 *http.Server, nefCtx *nefContext,
	stopServerCh chan bool) {
	if serverHTTP2 != nil {
		log.Infof("HTTP 2.0 listening on %s", serverHTTP2.Addr)
		if err := serverHTTP2.ListenAndServeTLS(
			nefCtx.cfg.HTTP2Config.NefServerCert,
			nefCtx.cfg.HTTP2Config.NefServerKey); err != nil {
			log.Errf("HTTP2server error: " + err.Error())
		}
	}
	stopServerCh <- true
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
	var server, serverHTTP2 *http.Server

	/* NEFRouter obeject is created. After creation this object contains all
	 * the HTTP Service Handlers. These hanlders will be called when HTTP
	 * server receives any HTTP Request */
	nefRouter := NewNEFRouter(nefCtx)

	// 1 for http2, 1 for http and 1 for the os signal
	numchannels := 3

	// Check if http and http 2 are both configured to determine number
	// of channels

	if nefCtx.cfg.HTTPConfig.Endpoint == "" {
		log.Info("HTTP Server not configured")
		numchannels--
	} else {
		// HTTP Server object is created
		server = &http.Server{
			Addr:           nefCtx.cfg.HTTPConfig.Endpoint,
			Handler:        nefRouter,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}

	if nefCtx.cfg.HTTP2Config.Endpoint == "" {
		log.Info("HTTP 2 Server not configured")
		numchannels--
	} else {
		serverHTTP2 = &http.Server{
			Addr:           nefCtx.cfg.HTTP2Config.Endpoint,
			Handler:        nefRouter,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		if err = http2.ConfigureServer(serverHTTP2,
			&http2.Server{}); err != nil {
			log.Errf("failed at configuring HTTP2 server")
			return err
		}
	}
	if server == nil && serverHTTP2 == nil {
		log.Err("HTTP Servers are not configured")
		return err
	}

	stopServerCh := make(chan bool, numchannels)

	/* Go Routine is spawned here for listening for cancellation event on
	 * context */
	go func(stopServerCh chan bool) {
		<-ctx.Done()
		if server != nil {
			log.Info("Executing graceful stop for HTTP Server")
			if err = server.Close(); err != nil {
				log.Errf("Could not close HTTP server: %#v", err)
			}
			log.Info("HTTP server stopped")
		}

		if serverHTTP2 != nil {

			if err = serverHTTP2.Close(); err != nil {
				log.Errf("Could not close HTTP2 server: %#v", err)
			}
			log.Info("HTTP2 server stopped")
		}

		/* De-initializes NEF Data */
		nefCtx.nef.nefDestroy()

		stopServerCh <- true
	}(stopServerCh)

	/* Go Routine is spawned here for starting HTTP Server */
	go startHTTPServer(server, stopServerCh)
	/* Go Routine is spawned here for starting HTTP-2 Server */
	go startHTTP2Server(serverHTTP2, nefCtx, stopServerCh)
	/* This self go routine is waiting for the receive events from the spawned
	 * go routines */
	<-stopServerCh
	<-stopServerCh
	if numchannels == 3 {
		<-stopServerCh
	}
	log.Info("Exiting NEF server")
	return nil

}

// LoadJSONConfig reads a file located at configPath and unmarshals it to
// config structure
func loadJSONConfig(configPath string, config interface{}) error {
	cfgData, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return err
	}
	return json.Unmarshal(cfgData, config)
}

// Run : This function reads the NEF Module configuration file and stores in
//       NEF Module Context. This also calls the Initialization/Creation of
//       NEF Data. Also it  calls runServer function for starting HTTP Server.
// Input Args:
//    - ctx:     NEF Module Running context
//    - cfgPath: This is NEF Module Configuration file path
// Output Args:
//     - error: returns error in case any error occurred in reading NEF
//              configuration file, NEF create error or any error occurred in
//              starting server
func Run(ctx context.Context, cfgPath string) error {

	var nefCtx nefContext

	/* Reads NEF Configuration file which is json format. Also it converts
	 * configuration data from json format to structure data */
	err := loadJSONConfig(cfgPath, &nefCtx.cfg)
	if err != nil {
		log.Errf("Failed to load NEF configuration: %v", err)
		return err

	}

	printConfig(nefCtx.cfg)

	/* Creates/Initializes NEF Data */
	err = nefCtx.nef.nefCreate(ctx, nefCtx.cfg)
	if err != nil {
		log.Errf("NEF Create Failed: %v", err)
		return err
	}

	return runServer(ctx, &nefCtx)
}

func printConfig(cfg Config) {

	log.Infoln("********************* NGC NEF CONFIGURATION ******************")
	log.Infoln("APIRoot: ", cfg.NefAPIRoot)
	log.Infoln("LocationPrefix: ", cfg.LocationPrefix)
	log.Infoln("UpfNotificationResUriPath:", cfg.UpfNotificationResURIPath)
	log.Infoln("UserAgent:", cfg.UserAgent)
	log.Infoln("-------------------------- NEF SERVER ----------------------")
	log.Infoln("EndPoint(HTTP): ", cfg.HTTPConfig.Endpoint)
	log.Infoln("EndPoint(HTTP2): ", cfg.HTTP2Config.Endpoint)
	log.Infoln("ServerCert(HTTP2): ", cfg.HTTP2Config.NefServerCert)
	log.Infoln("ServerKey(HTTP2): ", cfg.HTTP2Config.NefServerKey)
	log.Infoln("AFServerCert(HTTP2): ", cfg.HTTP2Config.AfServerCert)
	log.Infoln("*************************************************************")

}
