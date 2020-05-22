//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019-2020 Intel Corporation

package af_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otcshare/epcforedge/ngc/pkg/af"
)

var cfgPath string = "./testdata/testconfigs/af.json"

type srvData struct {
	ctx         context.Context
	srvCancel   context.CancelFunc
	afIsRunning bool
}

var testSrvData srvData

func TestAf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Af Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := af.Run(ctx, cfgPath)
		Expect(err).To(BeNil())
		testSrvData.afIsRunning = true
	}()
	time.Sleep(2 * time.Second)
	testSrvData.ctx = ctx
	testSrvData.srvCancel = cancel
})

var _ = AfterSuite(func() {
	testSrvData.srvCancel()
})
