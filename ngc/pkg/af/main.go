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
	"os"
	"os/signal"
	"syscall"

	logger "github.com/otcshare/common/log"

	 af "./lib"
)

var log = logger.DefaultLogger.WithField("main", nil)

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

	log.Infof("Starting server ...")
	af.Run(parenCtx, "../../configs/af.json")
}
