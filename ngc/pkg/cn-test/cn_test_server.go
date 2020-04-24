/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngccntest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	logtool "github.com/otcshare/common/log"
	"golang.org/x/net/http2"
)

// CnTestApp structure to store the variables/contexts for access in UT
type CnTestApp struct {
	cnTestRouter *mux.Router
	//cnTestCtx    *cnTestContext
}

// CnTestAppG is the NEF App variable which can be used for accessing the
// global contexts
var CnTestAppG CnTestApp

// Log handler initialized. This is to be used throughout the CN-TEST module for
// logging
var log = logtool.DefaultLogger.WithField("CN-TEST", nil)

//HTTPConfig contains the configuration for the HTTP 1.1
type HTTPConfig struct {
	Endpoint string `json:"endpoint"`
}

//HTTP2Config Contains the configuration for the HTTP2
type HTTP2Config struct {
	Endpoint         string `json:"endpoint"`
	CnTestServerCert string `json:"CnTestServerCert"`
	CnTestServerKey  string `json:"CnTestServerKey"`
	AfClientCert     string `json:"AfClientCert"`
}

// Config contains NEF Module Configuration Data Structure
type Config struct {
	HTTPConfig     HTTPConfig
	HTTP2Config    HTTP2Config
	LocationPrefix string `json:"locationPrefix"`
	CnTestAPIRoot  string `json:"CnTestAPIRoot"`
	MaxASCSupport  int    `json:"MaxASCSupport"`
	OAuth2Support  bool   `json:"OAuth2Support"`
}

// CN-TEST Module Context Data Structure
type cnTestContext struct {
	cfg Config
	//cntest cnTestData
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
func startHTTP2Server(serverHTTP2 *http.Server, cnTestCtx *cnTestContext,
	stopServerCh chan bool) {
	if serverHTTP2 != nil {
		log.Infof("HTTP 2.0 listening on %s", serverHTTP2.Addr)
		if err := serverHTTP2.ListenAndServeTLS(
			cnTestCtx.cfg.HTTP2Config.CnTestServerCert,
			cnTestCtx.cfg.HTTP2Config.CnTestServerKey); err != nil {
			log.Errf("HTTP2server error: " + err.Error())
		}
	}
	stopServerCh <- true
}

// runServer : This function creates a Router object and starts a HTTP Server
//             in a separate go routine. Also it listens for NEF module
//             running context cancellation event in another go routine. If
//             cancellation event occurs, it shutdowns the HTTP Server.
// Input Args:
//   - ctx:    NEF Module Running context
//   - nefCtx: This is NEF Module Context. This contains the NEF Module Data.
// Output Args:
//    - error: retruns no error. It only logs the error if any happened while
//             starting the HTTP Server
func runServer(ctx context.Context, cnTestCtx *cnTestContext) error {

	var err error
	var server, serverHTTP2 *http.Server

	/* CnTestRouter obeject is created. After creation this object contains all
	 * the HTTP Service Handlers. These handlers will be called when HTTP
	 * server receives any HTTP Request */
	cnTestRouter := NewCnTestRouter(cnTestCtx)
	CnTestAppG.cnTestRouter = cnTestRouter

	// 1 for http2, 1 for http and 1 for the os signal
	numchannels := 3

	// Check if http and http 2 are both configured to determine number
	// of channels

	if cnTestCtx.cfg.HTTPConfig.Endpoint == "" {
		log.Info("HTTP Server not configured")
		numchannels--
	} else {
		// HTTP Server object is created
		server = &http.Server{
			Addr:           cnTestCtx.cfg.HTTPConfig.Endpoint,
			Handler:        cnTestRouter,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}

	if cnTestCtx.cfg.HTTP2Config.Endpoint == "" {
		log.Info("HTTP 2 Server not configured")
		numchannels--
	} else {
		serverHTTP2 = &http.Server{
			Addr:           cnTestCtx.cfg.HTTP2Config.Endpoint,
			Handler:        cnTestRouter,
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
		return errors.New("HTTP Endpoints config missing")
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

		stopServerCh <- true
	}(stopServerCh)

	/* Go Routine is spawned here for starting HTTP Server */
	go startHTTPServer(server, stopServerCh)
	/* Go Routine is spawned here for starting HTTP-2 Server */
	go startHTTP2Server(serverHTTP2, cnTestCtx, stopServerCh)
	/* This self go routine is waiting for the receive events from the spawned
	 * go routines */
	<-stopServerCh
	<-stopServerCh
	if numchannels == 3 {
		<-stopServerCh
	}
	log.Info("Exiting CN-TEST server")
	return nil

}

// Run : This function creates CN-TEST Server
func Run(ctx context.Context, cfgPath string) error {
	var cnTestCtx cnTestContext

	err := loadJSONConfig(cfgPath, &cnTestCtx.cfg)
	if err != nil {
		log.Errf("Failed to load NEF configuration: %v", err)
		return err

	}
	printConfig(cnTestCtx.cfg)
	IntPolicyAuthorization(cnTestCtx.cfg)
	return runServer(ctx, &cnTestCtx)
}

func printConfig(cfg Config) {

	log.Infoln("********************* NGC CN-TEST CONFIGURATION ******************")
	log.Infoln("APIRoot: ", cfg.CnTestAPIRoot)
	log.Infoln("MaxASCSupport: ", cfg.MaxASCSupport)
	log.Infoln("OAuth2Support:", cfg.OAuth2Support)
	log.Infoln("LocationPrefix:", cfg.LocationPrefix)

	log.Infoln("-------------------------- CN-TEST SERVER ----------------------")
	log.Infoln("EndPoint(HTTP): ", cfg.HTTPConfig.Endpoint)
	log.Infoln("EndPoint(HTTP2): ", cfg.HTTP2Config.Endpoint)
	log.Infoln("ServerCert(HTTP2): ", cfg.HTTP2Config.CnTestServerCert)
	log.Infoln("ServerKey(HTTP2): ", cfg.HTTP2Config.CnTestServerKey)
	log.Infoln("AFClientCert(HTTP2): ", cfg.HTTP2Config.AfClientCert)
	log.Infoln("*************************************************************")

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

func getLocationURLPrefix(cfg *Config) string {

	var uri string
	// If http2 port is configured use it else http port
	if cfg.HTTP2Config.Endpoint != "" {
		uri = "https://" + cfg.CnTestAPIRoot +
			cfg.HTTP2Config.Endpoint
	} else {
		uri = "http://" + cfg.CnTestAPIRoot +
			cfg.HTTPConfig.Endpoint
	}
	uri += cfg.LocationPrefix
	return uri

}
