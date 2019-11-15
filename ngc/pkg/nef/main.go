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

	logtool "github.com/otcshare/common/log"
)

/* Log handler initialized. This is to be used throughout the nef module for
 * logging */
var log = logtool.DefaultLogger.WithField("NEF", nil)

/* Path for NEF Configuration file */
const cfgPath string = "../configs/nef.json"

/* Function: main
 * Description: Entry point for NEF Module Execution
 * Input Args: None
 * Output Args: None */
func main() {

	unusedlint()

	/* Opening a file for Logging and setting it to Logger Module */
	file, err1 := os.OpenFile("nef.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		log.Errf("Failed to open NEF log file: %s", err1.Error())
		os.Exit(1)
	}
	defer file.Close()
	logtool.SetOutput(file)

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
	 * 2. To notify in case context is cancelled. */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* Subscribing to os Interrupt/Signal - SIGTERM and waiting for
	 * notification in a separate go routine. When the notification is received
	 * the created context will be cancelled */
	osSignalCh := make(chan os.Signal, 1)
	signal.Notify(osSignalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-osSignalCh
		log.Infof("Received signal: %#v", sig)
		cancel()
	}()

	log.Infof("Starting NEF server ...")
	Run(ctx, cfgPath)
}

func unusedlint() {
	/* For unused variables lint warning to be  removed later */
	ti := TrafficInfluSub{}
	_ = ti

	tis := TrafficInfluSubPatch{}
	_ = tis

	upc := UpPathChange
	_ = upc

	pd := ProblemDetails{}
	_ = pd

	var ipx Ipv6Prefix = " "
	_ = ipx

	var pt Port = 8080
	_ = pt

	var apsd AppSessionID = "empty"
	_ = apsd

	een := NsmfEventExposureNotification{}
	_ = een

	pcfpr := PcfPolicyResponse{}
	_ = pcfpr

	udrpr := UdrInfluenceResponse{}
	_ = udrpr

	ac := AppSessionContext{}
	acu := AppSessionContextUpdateData{}

	// Avoid lint unused warning :  PCF client stub invocation

	var pcfClient PcfPolicyAuthorization = NewPCFClient(nil)
	ctx := context.Background()
	asd := AppSessionID("dummy")
	_, _, _ = pcfClient.PcfPolicyAuthorizationCreate(ctx, ac)
	_, _ = pcfClient.PolicyAuthorizationDelete(ctx, asd)
	_, _ = pcfClient.PolicyAuthorizationGet(ctx, asd)
	_, _ = pcfClient.PolicyAuthorizationUpdate(ctx, acu, asd)

	// Avoid lint unused warning :  UDR client stub invocation
	tid := TrafficInfluData{}
	tids := TrafficInfluDataPatch{}
	iid := InfluenceID("empty")

	var udrClient UdrInfluenceData = NewUDRClient(nil)
	_, _ = udrClient.UdrInfluenceDataCreate(ctx, tid, iid)
	_, _ = udrClient.UdrInfluenceDataUpdate(ctx, tids, iid)
	_, _ = udrClient.UdrInfluenceDataDelete(ctx, iid)
	_, _ = udrClient.UdrInfluenceDataGet(ctx)

	ev := EventNotification{}
	_ = ev

}
