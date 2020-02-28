// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	logger "github.com/otcshare/common/log"
	config "github.com/otcshare/epcforedge/ngc/pkg/config"
	oam "github.com/otcshare/epcforedge/ngc/pkg/oam"
)

type oamCfg struct {
	TLSEndpoint  string `json:"TlsEndpoint"`
	OpenEndpoint string `json:"OpenEndpoint"`
	UIEndpoint   string `json:"UIEndpoint"`
	NgcEndpoint  string `json:"NgcEndpoint"`
	NgcType      string `json:"NgcType"`
	NgcTestData  string `json:"NgcTestData"`
}

var log = logger.DefaultLogger.WithField("oam-main", nil)

func main() {

	lvl, err := logger.ParseLevel("info")
	if err != nil {
		log.Errf("Failed to parse log level: %s", err.Error())
		os.Exit(1)
	}
	logger.SetLevel(lvl)

	var cfg oamCfg
	err = config.LoadJSONConfig("./configs/oam.json", &cfg)
	if err != nil {
		log.Errf("Failed to load config: %s", err.Error())
		os.Exit(1)
	}
	log.Infof("LocalConfig: %s, %s, %s, %s, %s, %s\n",
		cfg.TLSEndpoint,
		cfg.OpenEndpoint,
		cfg.UIEndpoint,
		cfg.NgcEndpoint,
		cfg.NgcType,
		cfg.NgcTestData)

	// New Http Router
	err = oam.InitProxy(cfg.NgcEndpoint, cfg.NgcType, cfg.NgcTestData)
	if err != nil {
		log.Infof("Failed to init proxy: %s", err.Error())
		os.Exit(1)
	}

	router := oam.NewRouter()

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With",
		"Content-Type", "Authorization"})
	originsOK := handlers.AllowedOrigins([]string{cfg.UIEndpoint})
	methodsOK := handlers.AllowedMethods([]string{"GET", "HEAD", "POST",
		"PATCH", "OPTIONS", "DELETE"})

	serverOAM := &http.Server{
		Addr:         cfg.OpenEndpoint,
		Handler:      handlers.CORS(headersOK, originsOK, methodsOK)(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Infof("OAM HTTP Server Listening on:  %s\n", cfg.OpenEndpoint)
	serverOAM.ListenAndServe()
}
