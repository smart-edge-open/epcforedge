//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019-2020 Intel Corporation

package af_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otcshare/epcforedge/ngc/pkg/af"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var cfgPath string = "./testdata/testconfigs"

type srvData struct {
	ctx         context.Context
	srvCancel   context.CancelFunc
	afIsRunning bool
	notifServer *http.Server
}

var testSrvData srvData

func TestAf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Af Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := af.Run(ctx, cfgPath+"/af.json")
		Expect(err).To(BeNil())
		testSrvData.afIsRunning = true
	}()

	go func() {
		ngcnef.TestPcf = true
		err := ngcnef.Run(ctx, cfgPath+"/nef.json")
		Expect(err).To(BeNil())
	}()

	time.Sleep(2 * time.Second)
	testSrvData.ctx = ctx
	testSrvData.srvCancel = cancel

	// Start the Notify Server
	go func() {

		h2s := &http2.Server{}
		http.HandleFunc("/notification", NotificationPost)
		handler := http.HandlerFunc(NotificationPost)

		testSrvData.notifServer = &http.Server{
			Addr:         ":8450",
			Handler:      h2c.NewHandler(handler, h2s),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		testSrvData.notifServer.ListenAndServe()

	}()

})

var _ = AfterSuite(func() {
	testSrvData.srvCancel()
	testSrvData.notifServer.Close()
})
