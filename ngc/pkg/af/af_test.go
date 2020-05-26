//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019-2020 Intel Corporation

package af_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otcshare/epcforedge/ngc/pkg/af"
	config "github.com/otcshare/epcforedge/ngc/pkg/config"
)

type KeyType string

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testingAFClient(fn RoundTripFunc) *http.Client {

	return &http.Client{
		Transport: fn,
	}

}

func genTestConfig(protocol string, protocolVer string) af.GenericCliConfig {

	var (
		cfg     af.Config
		testCfg af.GenericCliConfig
	)

	err := config.LoadJSONConfig(cfgPath+"/af.json", &cfg)
	Expect(err).ShouldNot(HaveOccurred())

	testCfg = *(cfg.CliPcfCfg)
	testCfg.Protocol = protocol
	testCfg.ProtocolVer = protocolVer

	return testCfg
}

var _ = Describe("AF", func() {

	Describe("Utility ", func() {
		Context("HTTP Client generate", func() {
			Specify("Generate http 1.1 client", func() {
				cfg := genTestConfig("http", "1.1")

				By("Create HTTP Client")
				_, err := af.GenHTTPClient(&cfg)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Specify("Generate https 1.1 client", func() {
				cfg := genTestConfig("https", "1.1")

				By("Create HTTP Client")
				_, err := af.GenHTTPClient(&cfg)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Specify("Generate https 1.1 client", func() {
				cfg := genTestConfig("https", "1.1")

				By("Create HTTP Client")
				_, err := af.GenHTTPClient(&cfg)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Specify("Generate http 2.0 client", func() {
				cfg := genTestConfig("http", "2.0")

				By("Create HTTP Client")
				_, err := af.GenHTTPClient(&cfg)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Specify("Generate https 2.0 client", func() {
				cfg := genTestConfig("https", "2.0")

				By("Create HTTP Client")
				_, err := af.GenHTTPClient(&cfg)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Specify("Generate http 3.0 client", func() {
				cfg := genTestConfig("http", "3.0")

				By("Create HTTP Client")
				_, err := af.GenHTTPClient(&cfg)
				Expect(err).Should(HaveOccurred())
			})

			Specify("Generate https 3.0 client", func() {
				cfg := genTestConfig("https", "3.0")

				By("Create HTTP Client")
				_, err := af.GenHTTPClient(&cfg)
				Expect(err).Should(HaveOccurred())
			})
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

			Specify("Sending POST 005 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST005.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request with subID")
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

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending POST 006 request - no location URL", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST006.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request with subID")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/subscriptions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST006.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(201,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("Sending POST 007 request - no SUB-ID in URL", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST006.json")
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
				resBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST006.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				header := make(http.Header)
				header.Set("Location",
					"http://localhost:8080/af/v1/")
				httpclient := createTestHTTPClient(201,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("Sending POST 008 request - invalid json", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/invalid.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request with subID")
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

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

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

			Specify("Read all Subscriptions - 400", func() {

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				reqBody, err := ioutil.ReadFile(
					"./testdata/100_AF_NB_SUB_POST004.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)

				httpclient := createTestHTTPClient(400,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)

				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("Subscription ID GET", func() {
			Specify("VALID SUB-ID", func() {
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

			Specify("INVALID SUB ID", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions/11120",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNotFound))

			})

			Specify("INVALID GET SUBSCRIPTION 501", func() {

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/subscriptions/11111",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				httpclient := createTestHTTPClient(501,
					bytes.NewBufferString(`OK`), make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusNotImplemented))

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

			Specify("SUBSCRIPTION PUT INVALID SUB ID", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/300_AF_NB_SUB_SUBID_PUT001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/subscriptions/11120",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("SUBSCRIPTION PUT Invalid json", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/invalid.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/subscriptions/11111",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("PUT Subscription 501", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/300_AF_NB_SUB_SUBID_PUT001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/subscriptions/11111",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				httpclient := createTestHTTPClient(501,
					bytes.NewBufferString(`OK`), make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusNotImplemented))

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

			Specify("INVALID SUB ID", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/400_AF_NB_SUB_SUBID_PATCH001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/subscriptions/11120",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(
					req.Context(), KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("INVALID json", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/invalid.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/subscriptions/11111",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(
					req.Context(), KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

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
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})

			Specify("INVALID DELETE SUBSCRIPTION 501", func() {

				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/subscriptions/11114",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				httpclient := createTestHTTPClient(501,
					bytes.NewBufferString(`OK`), make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusNotImplemented))

			})

		})

		Describe("Cnca Notify Subscription to AF : ", func() {

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
			})

			Context("Subscription NOTIFY", func() {

				Specify("Sending NOTIFY 001 request", func() {

					By("Preparing Notify request with invalid json")
					ntfBody, err := ioutil.ReadFile(
						"./testdata/AF_SB_NOTIFY_POST005.json")
					Expect(err).ShouldNot(HaveOccurred())

					ntfBodyBytes := bytes.NewReader(ntfBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8081/af/v1/notifications",
						ntfBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

					Expect(resp.Code).To(Equal(http.StatusInternalServerError))

				})

				Specify("Sending NOTIFY 002 request", func() {

					By("Preparing Notify request with non existent Trans ID")
					ntfBody, err := ioutil.ReadFile(
						"./testdata/AF_SB_NOTIFY_POST001.json")
					Expect(err).ShouldNot(HaveOccurred())

					ntfBodyBytes := bytes.NewReader(ntfBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8081/af/v1/notifications",
						ntfBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

					Expect(resp.Code).To(Equal(http.StatusInternalServerError))

				})

				Specify("Sending NOTIFY 003 request", func() {

					By("Preparing Notify request with empty trans ID")
					ntfBody, err := ioutil.ReadFile(
						"./testdata/AF_SB_NOTIFY_POST002.json")
					Expect(err).ShouldNot(HaveOccurred())

					ntfBodyBytes := bytes.NewReader(ntfBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8081/af/v1/notifications",
						ntfBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

					Expect(resp.Code).To(Equal(http.StatusInternalServerError))

				})

				Specify("Sending NOTIFY 004 request", func() {

					By("Preparing Notify request with valid Trans ID")
					ntfBody, err := ioutil.ReadFile(
						"./testdata/AF_SB_NOTIFY_POST003.json")
					Expect(err).ShouldNot(HaveOccurred())

					ntfBodyBytes := bytes.NewReader(ntfBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8081/af/v1/notifications",
						ntfBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

					Expect(resp.Code).To(Equal(http.StatusOK))

				})

				Specify("Sending NOTIFY 005 request", func() {

					By("Preparing Notify request with invalid Trans ID")
					ntfBody, err := ioutil.ReadFile(
						"./testdata/AF_SB_NOTIFY_POST004.json")
					Expect(err).ShouldNot(HaveOccurred())

					ntfBodyBytes := bytes.NewReader(ntfBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8081/af/v1/notifications",
						ntfBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

					Expect(resp.Code).To(Equal(http.StatusInternalServerError))

				})

			})

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

		})

		Context("PFD  GET ALL - NO PFDS ", func() {
			Specify("Read all PFD transactions", func() {

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))
			})
		})

		Context("PFD  GET ALL - 400", func() {
			Specify("Read all PFD transactions", func() {

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_GETALL.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(400,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)

				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

		})

		Context("PFD Transaction POST", func() {
			Specify("Sending PFD POST 001 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusCreated))

			})

		})

		Context("PFD Transaction POST", func() {
			Specify("Sending PFD POST 002 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST002.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusCreated))

			})

		})

		Context("PFD Transaction INVALID POST", func() {
			Specify("Sending PFD POST request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD Transaction INVALID POST", func() {
			Specify("Decode error", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST003.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

		})

		Context("PFD POST Transaction  No Location Header", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(201,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD Transaction  POST SELF APP LINK MISSING", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST_SELF.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				header := make(http.Header)
				header.Set("Location",
					"http://localhost:8080/af/v1/pfd/transactions/10000")
				httpclient := createTestHTTPClient(201,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD Transaction  POST SELF LINK MISSING", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				header := make(http.Header)
				header.Set("Location",
					"http://localhost:8080/af/v1/pfd/transactions/10000")
				httpclient := createTestHTTPClient(201,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("DECODE error in PFD POST TRANS", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8080/af/v1/pfd/transactions",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_invalid.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(201,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})
		Context("PFD  GET ALL", func() {
			Specify("Read all PFD transactions", func() {

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))
			})

			Specify("Read all PFD transactions", func() {
				By("sending wrong url")
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v2/pfd/transactions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNotFound))

			})

			Specify("DECODE error in GET ALL", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(200,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("SELF LINK MISSING IN GET ALL", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_GETALL.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(200,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("SELF APP LINK MISSING IN GET ALL", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_GETALL_SELF.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(200,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD transaction ID GET", func() {
			Specify("", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("SELF LINK MISIING in GET", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(200,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("SELF LINK MISIING in APP", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST_SELF.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(200,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("DECODE error in PFD GET", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_invalid.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(200,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("INVALID GET PFD TRANS", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/11000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNotFound))

			})

		})

		Context("PFD Transcation ID PUT", func() {
			Specify("", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_PUT001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("PUT - SELF LINK MISSING", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("PUT - SELF LAPP INK MISSING", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST_SELF.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("INVALID PUT", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_PUT002.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("INVALID PUT - Decode", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_POST003.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("DECODE error in PFD PUT TRANS", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_PUT001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_invalid.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD Transcation DELETE", func() {
			Specify("DELETE PFD Transaction 02", func() {

				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/pfd/transactions/10001",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})
			Specify("INVALID DELETE PFD Transaction 10", func() {

				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/pfd/transactions/11000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})

			Specify("INVALID DELETE PFD Transaction 400", func() {

				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/pfd/transactions/11000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				httpclient := createTestHTTPClient(400,
					bytes.NewBufferString(`OK`), make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("INVALID DELETE PFD Transaction 451", func() {

				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/pfd/transactions/11000",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				httpclient := createTestHTTPClient(451,
					bytes.NewBufferString(`OK`), make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusUnavailableForLegalReasons))

			})

		})

		Context("PFD transaction Application GET", func() {
			Specify("", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("INVALID GET PFD TRANS 10000 and app10", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app10", nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNotFound))

			})

			Specify("PFD TRANS GET SELF APP LINK MISSING", func() {
				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("DECODE error in PFD APP GET", func() {

				req, err := http.NewRequest(http.MethodGet,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_invalid.json")
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				httpclient := createTestHTTPClient(200,
					reqBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD transaction Application PUT", func() {
			Specify("", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("INVALID PUT FOR PFD TRANS 10000 and app1", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_02.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("INVALID PUT FOR PFD TRANS/ APP - Decode error", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_03.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("PFD APP PUT SELF APP LINK MISSING", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("DECODE error in PFD APP PUT", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_invalid.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD transaction Application PATCH", func() {
			Specify("", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PATCH_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("INVALID PATCH PFD TRANS 10000 and app1", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PATCH_02.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("INVALID PATCH FOR PFD TRANS/ APP - Decode error", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_03.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("PFD APP PATCH SELF APP LINK MISSING", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PATCH_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("DECODE error in PFD APP PATCH", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_APP_PATCH_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				resBody, err := ioutil.ReadFile(
					"./testdata/pfd/AF_NB_PFD_invalid.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				httpclient := createTestHTTPClient(200,
					resBodyBytes, make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

		})

		Context("PFD transaction Application DELETE", func() {
			Specify("", func() {
				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNoContent))

			})

			Specify("INVALID DELETE TRANSACTION APPLICATION", func() {
				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app10", nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNotFound))

			})

			Specify("INVALID DELETE PFD APP 451", func() {

				req, err := http.NewRequest(http.MethodDelete,
					"http://localhost:8080/af/v1/pfd/transactions/10000/"+
						"applications/app1", nil)
				Expect(err).ShouldNot(HaveOccurred())

				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)

				httpclient := createTestHTTPClient(451,
					bytes.NewBufferString(`OK`), make(http.Header))

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusUnavailableForLegalReasons))

			})

		})

	})
})
