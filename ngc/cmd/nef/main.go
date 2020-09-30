// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	logtool "github.com/open-ness/common/log"
	ngcnef "github.com/open-ness/epcforedge/ngc/pkg/nef"
)

// Log handler initialized. This is to be used for NEF Main
var log = logtool.DefaultLogger.WithField("NEF-MAIN", nil)

// Path for NEF Configuration file
const cfgPath string = "configs/nef.json"

// main: Entry point for NEF Module Execution
// Input Args: None
// Output Args: None
func main() {

	/* Reading Log Level and and set it to logger, As of now it is hardcoded to
	 * info */
	lvl, err := logtool.ParseLevel("info")
	if err != nil {
		log.Errf("Failed to parse log level: %s", err.Error())
		os.Exit(1)
	}
	logtool.SetLevel(lvl)
	log.Infof("Logger Level: %d", lvl)

	/* Creating a context. This context will be used for following:
	 * 1. To store the NEF Module Context data and other module related data.
	 * 2. To notify in case context is canceled. */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* Subscribing to os Interrupt/Signal - SIGTERM and waiting for
	 * notification in a separate go routine. When the notification is received
	 * the created context will be canceled */
	osSignalCh := make(chan os.Signal, 1)
	signal.Notify(osSignalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-osSignalCh
		log.Infof("Received signal: %#v", sig)
		cancel()
	}()

	log.Infof("Starting NEF server ...")
	_ = ngcnef.Run(ctx, cfgPath)

}
