// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package main

import (
        "os"
        "net/http"
        logger "github.com/otcshare/common/log"
        oam "github.com/otcshare/epcforedge/ngc/pkg/oam"
        config "github.com/otcshare/epcforedge/ngc/pkg/config"
)

type oamCfg struct {
        TLSEndpoint        string        `json:"TlsEndpoint"`
        OpenEndpoint       string        `json:"OpenEndpoint"`
        NgcEndpoint        string        `json:"NgcEndpoint"`
        NgcType            string        `json:"NgcType"`
        NgcTestData        string        `json:"NgcTestData"`
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
        log.Infof("LocalConfig: %s, %s, %s, %s, %s\n", 
               cfg.TLSEndpoint, 
               cfg.OpenEndpoint, 
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
        log.Infof("OAM HTTP Server Listening on:  %s\n",cfg.OpenEndpoint); 
	http.ListenAndServe(cfg.OpenEndpoint, router)
}
