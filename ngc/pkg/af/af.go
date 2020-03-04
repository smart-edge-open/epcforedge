// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019-2020 Intel Corporation

package af

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	oauth2 "github.com/otcshare/epcforedge/ngc/pkg/oauth2"
	"golang.org/x/net/http2"

	logger "github.com/otcshare/common/log"
	config "github.com/otcshare/epcforedge/ngc/pkg/config"
)

const (
	// Enable/Disable HTTP2 Flag
	/* Temp Disabling http2 in AF NB */
	Http2Enabled = false
)

// TransactionIDs type
type TransactionIDs map[int]TrafficInfluSub

// NotifSubscryptions type
type NotifSubscryptions map[string]map[string]TrafficInfluSub

// Store the  Access token
var nefAccessToken string

// ServerConfig struct
type ServerConfig struct {
	CNCAEndpoint   string `json:"CNCAEndpoint"`
	Hostname       string `json:"Hostname"`
	NotifPort      string `json:"NotifPort"`
	UIEndpoint     string `json:"UIEndpoint"`
	ServerCertPath string `json:"ServerCertPath"`
	ServerKeyPath  string `json:"ServerKeyPath"`
}

//Config struct
type Config struct {
	AfID              string       `json:"AfId"`
	AfAPIRoot         string       `json:"AfAPIRoot"`
	LocationPrefixPfd string       `json:"LocationPrefixPfd"`
	SrvCfg            ServerConfig `json:"ServerConfig"`
	CliCfg            CliConfig    `json:"CliConfig"`
}

//Context struct
type Context struct {
	subscriptions NotifSubscryptions
	transactions  TransactionIDs
	cfg           Config
}

var (
	log = logger.DefaultLogger.WithField("ngc-af", nil)
	//AfCtx public var
	AfCtx *Context
	//AfRouter public var
	AfRouter *mux.Router
)

func runServer(ctx context.Context, AfCtx *Context) error {

	var err error

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With",
		"Content-Type", "Authorization"})
	originsOK := handlers.AllowedOrigins(
		[]string{AfCtx.cfg.SrvCfg.UIEndpoint})
	methodsOK := handlers.AllowedMethods([]string{"GET", "HEAD",
		"POST", "PUT", "PATCH", "OPTIONS", "DELETE"})

	AfCtx.transactions = make(TransactionIDs)
	AfCtx.subscriptions = make(NotifSubscryptions)
	AfRouter = NewAFRouter(AfCtx)
	NotifRouter := NewNotifRouter(AfCtx)

	serverCNCA := &http.Server{
		Addr:         AfCtx.cfg.SrvCfg.CNCAEndpoint,
		Handler:      handlers.CORS(headersOK, originsOK, methodsOK)(AfRouter),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if Http2Enabled == true {
		if err = http2.ConfigureServer(
			serverCNCA, &http2.Server{}); err != nil {
			log.Errf("AF failed at configuring HTTP2 server (CNCA Server)")
			return err
		}
	}
	serverNotif := &http.Server{
		Addr:         AfCtx.cfg.SrvCfg.NotifPort,
		Handler:      NotifRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err = http2.ConfigureServer(serverNotif, &http2.Server{}); err != nil {
		log.Errf("AF failed at configuring HTTP2 server (NEF Server)")
		return err
	}

	if AfCtx.cfg.CliCfg.OAuth2Support {
		log.Infoln("Fetching NEF access token")
		if fetchNEFAuthorizationToken() != nil {
			log.Infoln("Failed to get access token")
			return err
		}
	} else {
		log.Infoln("OAuth2 DISABLED")
	}

	stopServerCh := make(chan bool, 2)
	go func(stopServerCh chan bool) {
		<-ctx.Done()
		log.Info("Executing graceful stop")
		if err = serverCNCA.Close(); err != nil {
			log.Errf("Could not close AF CNCA server: %#v", err)
		}
		log.Info("AF CNCA server stopped")
		if err = serverNotif.Close(); err != nil {
			log.Errf("Could not close AF Notifications server: %#v", err)
		}
		log.Info("AF Notification server stopped")
		stopServerCh <- true
	}(stopServerCh)

	go func(stopServerCh chan bool) {
		log.Infof("Serving AF Notifications on: %s",
			AfCtx.cfg.SrvCfg.NotifPort)
		if err = serverNotif.ListenAndServeTLS(
			AfCtx.cfg.SrvCfg.ServerCertPath,
			AfCtx.cfg.SrvCfg.ServerKeyPath); err != http.ErrServerClosed {

			log.Errf("AF Notifications server error: " + err.Error())
		}
		stopServerCh <- true
	}(stopServerCh)

	if Http2Enabled == true {
		log.Infof("Serving AF (CNCA HTTP2 Requests) on: %s",
			AfCtx.cfg.SrvCfg.CNCAEndpoint)
		err = serverCNCA.ListenAndServeTLS(AfCtx.cfg.SrvCfg.ServerCertPath,
			AfCtx.cfg.SrvCfg.ServerKeyPath)
	} else {
		log.Infof("Serving AF (CNCA HTTP Requests) on: %s",
			AfCtx.cfg.SrvCfg.CNCAEndpoint)
		err = serverCNCA.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		log.Errf("AF CNCA server error: " + err.Error())
		return err
	}

	<-stopServerCh
	<-stopServerCh
	return nil
}

func printConfig(cfg Config) {

	log.Infoln("********************* NGC AF CONFIGURATION ******************")
	log.Infoln("AfID: ", cfg.AfID)
	log.Infoln("LocationPrefixPfd ", cfg.LocationPrefixPfd)
	log.Infoln("-------------------------- CNCA SERVER ----------------------")
	log.Infoln("CNCAEndpoint: ", cfg.SrvCfg.CNCAEndpoint)
	log.Infoln("-------------------- NEF NOTIFICATIONS SERVER ---------------")
	log.Infoln("Hostname: ", cfg.SrvCfg.Hostname)
	log.Infoln("NotifPort: ", cfg.SrvCfg.NotifPort)
	log.Infoln("ServerCertPath: ", cfg.SrvCfg.ServerCertPath)
	log.Infoln("ServerKeyPath: ", cfg.SrvCfg.ServerKeyPath)
	log.Infoln("UIEndpoint: ", cfg.SrvCfg.UIEndpoint)
	log.Infoln("------------------------- CLIENT TO NEF ---------------------")
	log.Infoln("Protocol: ", cfg.CliCfg.Protocol)
	log.Infoln("NEFPort: ", cfg.CliCfg.NEFPort)
	log.Infoln("NEFBasePath: ", cfg.CliCfg.NEFBasePath)
	log.Infoln("NEFPFDBasePath: ", cfg.CliCfg.NEFPFDBasePath)
	log.Infoln("UserAgent: ", cfg.CliCfg.UserAgent)
	log.Infoln("NEFCliCertPath: ", cfg.CliCfg.NEFCliCertPath)
	log.Infoln("OAuth2Support: ", cfg.CliCfg.OAuth2Support)
	log.Infoln("*************************************************************")

}

// Run function
func Run(parentCtx context.Context, cfgPath string) error {

	var AfCtx Context

	// load AF configuration from file
	err := config.LoadJSONConfig(cfgPath, &AfCtx.cfg)

	if err != nil {
		log.Errf("Failed to load AF configuration: %v", err)
		return err
	}
	printConfig(AfCtx.cfg)

	return runServer(parentCtx, &AfCtx)
}

func fetchNEFAuthorizationToken() error {

	var err error

	nefAccessToken, err = oauth2.GetAccessToken()
	if err != nil {
		log.Errf("Failed to Fetch Access Token ")
		return err
	}
	log.Errf("Got Access Token ", nefAccessToken)
	return nil
}

func getNEFAuthorizationToken() (token string, err error) {

	return nefAccessToken, nil
}
