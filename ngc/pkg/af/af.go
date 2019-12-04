// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package ngcaf

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	logger "github.com/otcshare/common/log"
	config "github.com/otcshare/epcforedge/ngc/pkg/config"
)

// TransactionIDs type
type TransactionIDs map[int]TrafficInfluSub

// NotifSubscryptions type
type NotifSubscryptions map[string]map[string]TrafficInfluSub

// ServerConfig struct
type ServerConfig struct {
	CNCAEndpoint        string `json:"CNCAEndpoint"`
        Hostname	    string `json:"Hostname"`
	NotifPort       string `json:"NotifPort"`
	NotifServerCertPath string `json:"NotifServerCertPath"`
	NotifServerKeyPath  string `json:"NotifServerKeyPath"`
}

//Config struct
type Config struct {
	AfID   string       `json:"AfId"`
	SrvCfg ServerConfig `json:"ServerConfig"`
	CliCfg CliConfig    `json:"CliConfig"`
}

//AFContext struct
type AFContext struct {
	subscriptions NotifSubscryptions
	transactions  TransactionIDs
	cfg           Config
}

var log = logger.DefaultLogger.WithField("ngc-af", nil)

func runServer(ctx context.Context, afCtx *AFContext) error {

	var err error

	afCtx.transactions = make(TransactionIDs)
	afCtx.subscriptions = make(NotifSubscryptions)
	afRouter := NewAFRouter(afCtx)
	nRouter := NewNotifRouter(afCtx)

	serverCNCA := &http.Server{
		Addr:         afCtx.cfg.SrvCfg.CNCAEndpoint,
		Handler:      afRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverNotif := &http.Server{
		Addr:         afCtx.cfg.SrvCfg.NotifPort,
		Handler:      nRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err = http2.ConfigureServer(serverNotif, &http2.Server{}); err != nil {
		log.Errf("AF failed at configuring HTTP2 server")
		return err
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
			afCtx.cfg.SrvCfg.NotifServerCertPath,
			afCtx.cfg.SrvCfg.NotifServerKeyPath); err != http.ErrServerClosed {

			log.Errf("AF Notifications server error: " + err.Error())
		}
		stopServerCh <- true
	}(stopServerCh)

	log.Infof("Serving AF on: %s", afCtx.cfg.SrvCfg.CNCAEndpoint)
	if err = serverCNCA.ListenAndServe(); err != http.ErrServerClosed {
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
	log.Infoln("-------------------------- CNCA SERVER ----------------------")
	log.Infoln("CNCAEndpoint: ", cfg.SrvCfg.CNCAEndpoint)
	log.Infoln("-------------------- NEF NOTIFICATIONS SERVER ---------------")
        log.Infoln("Hostname: ", cfg.SrvCfg.Hostname)
	log.Infoln("NotifPort: ", cfg.SrvCfg.NotifPort)
	log.Infoln("NotifServerCertPath: ", cfg.SrvCfg.NotifServerCertPath)
	log.Infoln("NotifServerKeyPath: ", cfg.SrvCfg.NotifServerKeyPath)
	log.Infoln("------------------------- CLIENT TO NEF ---------------------")
	log.Infoln("NEFBasePath: ", cfg.CliCfg.NEFBasePath)
	log.Infoln("UserAgent: ", cfg.CliCfg.UserAgent)
	log.Infoln("NEFCliCertPath: ", cfg.CliCfg.NEFCliCertPath)
	log.Infoln("*************************************************************")

}

// Run function
func Run(parentCtx context.Context, cfgPath string) error {

	var afCtx AFContext

	// load AF configuration from file
	err := config.LoadJSONConfig(cfgPath, &afCtx.cfg)

	if err != nil {
		log.Errf("Failed to load AF configuration: %v", err)
		return err
	}
	printConfig(afCtx.cfg)

	return runServer(parentCtx, &afCtx)
}
