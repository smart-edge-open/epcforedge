//Copyright 2019 Intel Corporation. All rights reserved.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package ngcaf_test

import (
	"context"
	"fmt"
	"testing"

	af "github.com/otcshare/epcforedge/ngc/pkg/af"
	//nef "github.com/otchsare/epcforedge/ngc/pkg/nef"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AF Suite")
}

const (
	cfgPath = "$GOPATH/src/github.com/otcshare/epcforedge/ngc"
)

var (
	srvCtx       context.Context
	srvCancel    context.CancelFunc
	afIsRunning  bool
	nefIsRunning bool
)

func runAf(stopIndication chan bool) error {

	By("Starting AF server")

	srvCtx, srvCancel = context.WithCancel(context.Background())
	_ = srvCancel
	afRunFail := make(chan bool)
	go func() {
		err := af.Run(srvCtx, cfgPath+"/configs/af.json")
		if err != nil {
			fmt.Printf("Run() exited with error: %#v", err)
			afIsRunning = false
			afRunFail <- true
		}
		stopIndication <- true
	}()

	return nil
}

func runNef(stopIndication chan bool) error {
	By("Starting NEF server")
	srvCtx, srvCancel = context.WithCancel(context.Background())
	_ = srvCancel
	nefRunFail := make(chan bool)
	go func() {
		err := nef.Run(srvCtx, cfgPath+"/configs/nef.json")
		if err != nil {
			fmt.Printf("Run() exited with error: %#v", err)
			nefIsRunning = false
			nefRunFail <- true
		}
		stopIndication <- true
	}()

	return nil

}

func stopNef(stopIndication chan bool) int {
	By("Stopping NEF server")
	srvCancel()
	<-stopIndication
	if nefIsRunning == true {
		return 0
	}
	return 1
}

func stopAf(stopIndication chan bool) int {
	By("Stopping AF server")
	srvCancel()
	<-stopIndication
	if afIsRunning == true {
		return 0
	}

	return 1
}
