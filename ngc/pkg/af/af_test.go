//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019 Intel Corporation

package ngcaf_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	//"time"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otcshare/tls-af-stub/ngc/pkg/af"
)

const (
	host     = "localhost:8080"
	basePath = "/af/v1"
)

func TestAf(t *testing.T) {
        RegisterFailHandler(Fail)
        RunSpecs(t, "AF Suite")
}

var _ = Describe("AF", func() {

	var (
		ctx context.Context
		srvCancel context.CancelFunc
		afIsRunning bool
		)

	ctx, _ = context.WithCancel(context.Background())

	Describe("Cnca client request methods to AF : ", func() {

		Context("Subscription GET ALL", func() {

			By("Starting AF server")
			ctx, srvCancel = context.WithCancel(context.Background())
			_ = srvCancel
			afRunFail := make(chan bool)
			go func() {
				err := ngcaf.Run(ctx, "../../configs/af.json")
				Expect(err).ShouldNot(HaveOccurred())
				if err != nil {
					fmt.Printf("Run() exited with error: %#v", err)
					afIsRunning = false
					afRunFail <- true
				}
			}()
			_=afIsRunning
		})
	})

	//var log = logger.DefaultLogger.WithField("ngc-af Test Suite", nil)
	Describe("Cnca client request methods to AF : ", func() {

		Context("Subscription GET ALL", func() {
			Specify("", func() {
				//time.Sleep(3 * time.Second)

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), string("af-ctx"), ngcaf.AfCtx_g)
				ngcaf.AfRouter_g.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))
			})
		})

		Context("Subscription GET ALL - 2", func() {
			Specify("", func() {
				//time.Sleep(3 * time.Second)

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), string("af-ctx"), ngcaf.AfCtx_g)
				ngcaf.AfRouter_g.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusOK))
			})
		})

		Context("Subscription POST", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile("./testdata/100_AF_NB_SUB_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/af/v1/subscriptions", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), string("af-ctx"), ngcaf.AfCtx_g)
				ngcaf.AfRouter_g.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusCreated))

			})
		})

		Context("Subscription ID GET", func() {
			Specify("", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions/1",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), string("af-ctx"), ngcaf.AfCtx_g)
				ngcaf.AfRouter_g.ServeHTTP(resp, req.WithContext(ctx))
		
				Expect(resp.Code).To(Equal(http.StatusOK))
				Expect(resp.Location).NotTo(Equal(""))
			
			})
		})

		Context("Subscription ID PUT", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile("./testdata/300_AF_NB_SUB_SUBID_PUT001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/af/v1/subscriptions/1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), string("af-ctx"), ngcaf.AfCtx_g)
				ngcaf.AfRouter_g.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

		})

		Context("Subscription ID PATCH", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile("./testdata/400_AF_NB_SUB_SUBID_PATCH001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch, "http://localhost:8080/af/v1/subscriptions/1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), string("af-ctx"), ngcaf.AfCtx_g)
				ngcaf.AfRouter_g.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

		})

		Context("Subscription ID DELETE", func() {
			Specify("", func() {
				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/subscriptions/1",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(), string("af-ctx"), ngcaf.AfCtx_g)
				ngcaf.AfRouter_g.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})

		})
	})

	//another describe to stop
	Describe("Stop the AF Server", func() {
		It("Disconnect AF fServer", func() {
			srvCancel()
			//time.Sleep(2 * time.Second)
		})
	})
})
