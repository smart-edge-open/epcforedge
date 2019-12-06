//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019 Intel Corporation

package af_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otcshare/epcforedge/ngc/pkg/af"
)

func TestAf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AF Suite")
}

type KeyType string

var _ = Describe("AF", func() {

	var (
		ctx         context.Context
		srvCancel   context.CancelFunc
		afIsRunning bool
	)

	ctx, srvCancel = context.WithCancel(context.Background())

	Describe("Cnca client request methods to AF : ", func() {

		Context("Subscription GET ALL", func() {

			By("Starting AF server")
			ctx, srvCancel = context.WithCancel(context.Background())
			_ = srvCancel
			afRunFail := make(chan bool)
			go func() {
				err := af.Run(ctx, "testdata/testconfigs/af.json")
				Expect(err).ShouldNot(HaveOccurred())
				if err != nil {
					fmt.Printf("Run() exited with error: %#v", err)
					afIsRunning = false
					afRunFail <- true
				}
			}()
			_ = afIsRunning
		})
	})

	Describe("Cnca client request methods to AF : ", func() {

		Context("Subscription POST", func() {
			Specify("Sending POST 001 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/subscriptions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusCreated))

			})
			Specify("Sending POST 002 Request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST002.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/subscriptions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusCreated))

			})

			Specify("Sending POST 003 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST003.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/subscriptions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusCreated))

			})
			Specify("Sending POST 004 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST004.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request with subID")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/subscriptions/1000",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusMethodNotAllowed))

			})

		})

		Context("Subscription GET ALL", func() {
			Specify("Read all subscriptions", func() {

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))
			})

			Specify("Read all subscriptions", func() {
				By("sending wrong url")
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v2/subscriptions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNotFound))

			})
		})

		Context("Subscription ID GET", func() {
			Specify("", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions/11112",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})
		})

		Context("Subscription ID PUT", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/300_AF_NB_SUB_SUBID_PUT001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/subscriptions/11113",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

		})

		Context("Subscription ID PATCH", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/400_AF_NB_SUB_SUBID_PATCH001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/subscriptions/11112",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(
					req.Context(), KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

		})

		Context("Subscription ID DELETE", func() {
			Specify("DELETE Subcription 01", func() {

				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/subscriptions/11111",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})

			Specify("DELETE Subcription 02", func() {
				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/subscriptions/11112",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})
			Specify("DELETE Subcription 03", func() {
				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/subscriptions/11113",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})
			Specify("DELETE Subcription 04", func() {
				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/subscriptions/11114",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Stop the AF Server", func() {
		It("Disconnect AF Server", func() {
			srvCancel()
		})
	})
})
