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
	// HTTP2Enabled flag is for enabling/disabling HTTP2
	HTTP2Enabled = true
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
	AfID              string            `json:"AfId"`
	AfAPIRoot         string            `json:"AfAPIRoot"`
	LocationPrefixPfd string            `json:"LocationPrefixPfd"`
	LocationPrefixPA  string            `json:"LocationPrefixPA"`
	UserAgent         string            `json:"UserAgent"`
	SrvCfg            ServerConfig      `json:"ServerConfig"`
	CliCfg            CliConfig         `json:"CliConfig"`
	CliPcfCfg         *GenericCliConfig `json:"CliPAConfig"`
}

type afData struct {
	policyAuthAPIClient pcfPolicyAuthAPI
	// TODO websocket connections of all consumers, consumerID is the key
	//consumerConns      ConsumerConns
}

//Context struct
type Context struct {
	subscriptions NotifSubscryptions
	transactions  TransactionIDs
	appSessionsEv AppSessEv
	cfg           Config
	data          afData
}

var (
	log = logger.DefaultLogger.WithField("ngc-af", nil)
	//AfCtx public var
	AfCtx *Context
	//AfRouter public var
	AfRouter *mux.Router
	//NotifRouter public var
	NotifRouter *mux.Router
)

func runServer(ctx context.Context, afCtx *Context) error {

	var err error

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With",
		"Content-Type", "Authorization"})
	originsOK := handlers.AllowedOrigins(
		[]string{afCtx.cfg.SrvCfg.UIEndpoint})
	methodsOK := handlers.AllowedMethods([]string{"GET",
		"POST", "PUT", "PATCH", "DELETE"})

	afCtx.transactions = make(TransactionIDs)
	afCtx.subscriptions = make(NotifSubscryptions)

	err = initAFData(afCtx)
	if err != nil {
		return err
	}

	AfRouter = NewAFRouter(afCtx)
	NotifRouter = NewNotifRouter(afCtx)

	serverCNCA := &http.Server{
		Addr:         afCtx.cfg.SrvCfg.CNCAEndpoint,
		Handler:      handlers.CORS(headersOK, originsOK, methodsOK)(AfRouter),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if HTTP2Enabled == true {
		if err = http2.ConfigureServer(
			serverCNCA, &http2.Server{}); err != nil {
			log.Errf("AF failed at configuring HTTP2 server (CNCA Server)")
			return err
		}
	}
	serverNotif := &http.Server{
		Addr:         afCtx.cfg.SrvCfg.NotifPort,
		Handler:      NotifRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err = http2.ConfigureServer(serverNotif, &http2.Server{}); err != nil {
		log.Errf("AF failed at configuring HTTP2 server (NEF Server)")
		return err
	}

	if afCtx.cfg.CliCfg.OAuth2Support {
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
			afCtx.cfg.SrvCfg.NotifPort)
		if err = serverNotif.ListenAndServeTLS(
			afCtx.cfg.SrvCfg.ServerCertPath,
			afCtx.cfg.SrvCfg.ServerKeyPath); err != http.ErrServerClosed {

			log.Errf("AF Notifications server error: " + err.Error())
		}
		stopServerCh <- true
	}(stopServerCh)

	if HTTP2Enabled == true {
		log.Infof("Serving AF (CNCA HTTP2 Requests) on: %s",
			afCtx.cfg.SrvCfg.CNCAEndpoint)
		err = serverCNCA.ListenAndServeTLS(afCtx.cfg.SrvCfg.ServerCertPath,
			afCtx.cfg.SrvCfg.ServerKeyPath)
	} else {
		log.Infof("Serving AF (CNCA HTTP Requests) on: %s",
			afCtx.cfg.SrvCfg.CNCAEndpoint)
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

func initAFData(afCtx *Context) (err error) {
	if err = initPACfg(afCtx); err == nil {
		initNotify(afCtx)
	}
	return err
}

/*
 * function to initialize different variable specific to Policy Authorization
 * - Initiate Policy auth api client which is reused for connecting to PCF
 * - Initiate Notification URI to be used while sending req to PCF
 */
func initPACfg(afCtx *Context) (err error) {

	paCfg := afCtx.cfg.CliPcfCfg
	cfg := afCtx.cfg
	err = validateCliPACfg(paCfg)
	if err != nil {
		log.Errf("Policy Auth client configuration invalid")
		return err
	}

	afCtx.data.policyAuthAPIClient, err =
		NewPolicyAuthAPIClient(&cfg)
	if err != nil {
		log.Errf("Unable to create policy auth api client")
		return err
	}

	pcfPANotifURI = "https://" + cfg.SrvCfg.Hostname +
		cfg.SrvCfg.NotifPort + paCfg.NotifURI

	smfPANotifURI = "https://" + cfg.SrvCfg.Hostname +
		cfg.SrvCfg.NotifPort + paCfg.NotifURI + "/smfnotify"

	return nil
}

func printGenericClientConfig(cfg *GenericCliConfig) {
	log.Infoln("Protocol: ", cfg.Protocol)
	log.Infoln("ProtocolVer: ", cfg.ProtocolVer)
	log.Infoln("Hostname: ", cfg.Hostname)
	log.Infoln("Port: ", cfg.Port)
	log.Infoln("BasePath: ", cfg.BasePath)
	log.Infoln("CliCertPath: ", cfg.CliCertPath)
	log.Infoln("OAuth2Support: ", cfg.OAuth2Support)
	log.Infoln("NotifURI: ", cfg.NotifURI)
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
	log.Infoln("NotifyClientCertPath: ", cfg.CliCfg.NotifyClientCertPath)
	log.Infoln("OAuth2Support: ", cfg.CliCfg.OAuth2Support)
	log.Infoln("--------------- CLIENT TO PCF (Policy Auth)---------------")
	printGenericClientConfig(cfg.CliPcfCfg)
	log.Infoln("**********************************************************")

}

// Run function
func Run(parentCtx context.Context, cfgPath string) error {

	var afCtx Context

	// load AF configuration from file
	err := config.LoadJSONConfig(cfgPath, &afCtx.cfg)

	if err != nil {
		log.Errf("Failed to load AF configuration: %v", err)
		return err
	}
	printConfig(afCtx.cfg)

	return runServer(parentCtx, &afCtx)
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
