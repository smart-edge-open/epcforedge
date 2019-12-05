// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/otcshare/common/log"
	"github.com/otcshare/epcforedge/ngc/pkg/af"
)

var log = logger.DefaultLogger.WithField("main", nil)

const cfgPath = "configs/af.json"

func main() {

	lvl, err := logger.ParseLevel("info")
	if err != nil {
		log.Errf("Failed to parse log level: %s", err.Error())
		os.Exit(1)
	}
	logger.SetLevel(lvl)
	parenCtx, cancel := context.WithCancel(context.Background())
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-osSignals
		log.Infof("Received signal: %#v", sig)
		cancel()
	}()

	log.Infof("Starting NGC AF servers..")
	if err = af.Run(parenCtx, cfgPath); err != nil {
		log.Errf("AF finished with error: %v", err)
	}
}
