//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019-2020 Intel Corporation

package af_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otcshare/epcforedge/ngc/pkg/af"
)

func TestAf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AF Suite")
}

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

func NotificationPost(w http.ResponseWriter, r *http.Request) {

	defer GinkgoRecover()
	Expect(r.Body).ShouldNot(Equal(nil))
	log.Println("Protocol:     Method", r.Proto, r.Method)
	defer r.Body.Close()
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")

}

/*
func startServer(server *httptest.Server, stopServerCh chan bool) {

	server.Start()
	stopServerCh <- true
}*/

// connectConsumer sends a consumer notifications GET request to the appliance
func connectWsAF(socket *websocket.Dialer, header *http.Header) *websocket.Conn {
	By("Sending consumer notification GET request")
	conn, resp, err := socket.Dial("wss://localhost:8082/af/v1/af-notifications", *header)
	Expect(err).ShouldNot(HaveOccurred())

	By("Comparing GET response code")
	defer resp.Body.Close()
	Expect(resp.Status).To(Equal("101 Switching Protocols"))

	return conn
}

// connectConsumer sends a consumer notifications GET request to the appliance
func connectWsAFForbidden(socket *websocket.Dialer, header *http.Header) {
	By("Sending consumer notification GET request")
	_, resp, _ := socket.Dial("wss://localhost:8082/af/v1/af-notifications", *header)
	//Expect(err).ShouldNot(HaveOccurred())

	By("Comparing GET response code")
	defer resp.Body.Close()
	Expect(resp.StatusCode).To(Equal(403))

}

func getNotifyFromConn(conn *websocket.Conn, response *af.Afnotification,
	corrID string) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	By("Reading message from web socket connection")
	err := conn.ReadJSON(response)
	Expect(err).ShouldNot(HaveOccurred())

	By("Received notification struct decoding")
	if response.Event == af.UPPathChangeEvent {
		var ev af.NotificationUpPathChg

		err = json.Unmarshal(response.Payload, &ev)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ev.NotifyID).To(Equal(corrID))
	}

}

var _ = Describe("AF", func() {

	var (
		ctx         context.Context
		srvCancel   context.CancelFunc
		afIsRunning bool
		server      *httptest.Server
	)

	Describe("Cnca client request methods to AF : ", func() {

		Context("Subscription GET ALL", func() {

			By("Starting AF server")
			var err error
			ctx, srvCancel = context.WithCancel(context.Background())
			_ = srvCancel
			afRunFail := make(chan bool)
			go func() {

				err = af.Run(ctx, "./testdata/testconfigs/af.json")

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 201,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters

						return &http.Response{
							StatusCode: 201,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: header,
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 400,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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

				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 501,
							// Send response to be tested
							Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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

				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 501,
							// Send response to be tested
							Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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

				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 501,
							// Send response to be tested
							Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 400,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 201,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters

						return &http.Response{
							StatusCode: 201,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: header,
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters

						return &http.Response{
							StatusCode: 201,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: header,
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 201,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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

				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 400,
							// Send response to be tested
							Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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

				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 451,
							// Send response to be tested
							Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(reqBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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
				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(resBodyBytes),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

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

				httpclient :=
					testingAFClient(func(req *http.Request) *http.Response {
						// Test request parameters
						return &http.Response{
							StatusCode: 451,
							// Send response to be tested
							Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					})

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusUnavailableForLegalReasons))

			})

		})
		Describe("Policy Authorization SMF Notification", func() {

			It("POST SMF notification for missing body", func() {

				By("Preparing request")
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8081/af/v1/policy-authorization/"+
						"smfnotify",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			It("POST SMF notification json parsing fialed", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/policy_auth/SMF_AF_NOTIF_err.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8081/af/v1/policy-authorization/"+
						"smfnotify",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			It("POST SMF notification notifid missing", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/policy_auth/SMF_AF_NOTIF_no_id.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8081/af/v1/policy-authorization/"+
						"smfnotify",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			It("POST SMF notification event notifs missing", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/policy_auth/SMF_AF_NOTIF_no_evts.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8081/af/v1/policy-authorization/"+
						"smfnotify",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("POST SMF  notification up_path_ch events missing", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(
					"./testdata/policy_auth/SMF_AF_NOTIF_no_upfs.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					"http://localhost:8081/af/v1/policy-authorization/"+
						"smfnotify",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})

			It("POST SMF notification for invalid correlation id",
				func() {

					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/SMF_AF_NOTIF_01.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8081/af/v1/policy-authorization/"+
							"smfnotify",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

					Expect(resp.Code).To(Equal(http.StatusNoContent))

				})

			Context("PolicyAuth POST/UPDATE/DELETE", func() {

				Specify("Sending POST 001 request", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_01.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					resBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_01.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					header := make(http.Header)
					header.Set("Location",
						"https://localhost:8095/af/v1/policy-authorization/"+
							"app-sessions/5001")
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 201,
								// Send response to be tested
								Body: ioutil.NopCloser(resBodyBytes),
								// Must be set to non-nil value or it panics
								Header: header,
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusCreated))

				})

				It("POST SMF notification with correct correlationId",
					func() {

						By("Reading json file")
						reqBody, err := ioutil.ReadFile(
							"./testdata/policy_auth/SMF_AF_NOTIF_01.json")
						Expect(err).ShouldNot(HaveOccurred())

						By("Preparing request")
						reqBodyBytes := bytes.NewReader(reqBody)
						req, err := http.NewRequest(http.MethodPost,
							"http://localhost:8081/af/v1/policy-authorization/"+
								"smfnotify",
							reqBodyBytes)
						Expect(err).ShouldNot(HaveOccurred())

						By("Sending request")
						resp := httptest.NewRecorder()
						ctx := context.WithValue(req.Context(),
							KeyType("af-ctx"), af.AfCtx)
						af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

						Expect(resp.Code).To(Equal(http.StatusNoContent))

					})

				Specify("Sending PATCH 001 request", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPATCH_01.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPatch,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					resBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_01.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 200,
								// Send response to be tested
								Body: ioutil.NopCloser(resBodyBytes),
								// Must be set to non-nil value or it panics
								Header: make(http.Header),
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusOK))

				})

				Specify("Sending PATCH 003 request - neither ws and notifyURI", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPATCH_no_ws_no_uri.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPatch,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					Expect(resp.Code).To(Equal(http.StatusBadRequest))

				})

				Specify("Sending DELETE request", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XDELETE_01.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001/delete",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 204,
								// Send response to be tested
								Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
								// Must be set to non-nil value or it panics
								Header: make(http.Header),
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusNoContent))

				})

				Specify("Sending POST 002 request - ws", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_ws.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					resBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_ws.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					header := make(http.Header)
					header.Set("Location",
						"https://localhost:8095/af/v1/"+
							"policy-authorization/app-sessions/5001")
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 201,
								// Send response to be tested
								Body: ioutil.NopCloser(resBodyBytes),
								// Must be set to non-nil value or it panics
								Header: header,
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusCreated))

				})

				Specify("Sending PATCH 002 request - ws already established", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPATCH_ws.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPatch,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)

					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					Expect(resp.Code).To(Equal(http.StatusBadRequest))

				})

				Specify("Sending PATCH 001 request - Both ws and notifyURI", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPATCH_01.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPatch,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001",
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

			Context("PolicyAuth GET", func() {
				Specify("Sending GET request", func() {
					req, err := http.NewRequest(http.MethodGet,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001",
						nil)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					resBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_ws.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					header := make(http.Header)
					header.Set("Location",
						"https://localhost:8095/af/v1/"+
							"policy-authorization/app-sessions/5001")
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 200,
								// Send response to be tested
								Body: ioutil.NopCloser(resBodyBytes),
								// Must be set to non-nil value or it panics
								Header: header,
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusOK))

				})

				Specify("Sending DELETE request", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XDELETE_01.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001/delete",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 204,
								// Send response to be tested
								Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
								// Must be set to non-nil value or it panics
								Header: make(http.Header),
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusNoContent))

				})
			})
			Context("PolicyAuth POST - Both ws and notificationURI", func() {
				Specify("Sending POST 003 request", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_both_ws_uri.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)

					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					Expect(resp.Code).To(Equal(http.StatusBadRequest))

				})

				Specify("Sending POST 004 request", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_no_ws_no_uri.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions",
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

			Context("PolicyAuth POST/UPDATE/DELETE - Websocket", func() {
				Specify("Sending POST request - ws", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_1_corrID_ws.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					resBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_1_corrID_ws.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					header := make(http.Header)
					header.Set("Location",
						"https://localhost:8095/af/v1/"+
							"policy-authorization/app-sessions/5001")
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 201,
								// Send response to be tested
								Body: ioutil.NopCloser(resBodyBytes),
								// Must be set to non-nil value or it panics
								Header: header,
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusCreated))
					var appSess af.AppSessionContext
					// Decode response to check for websocketURI
					err = json.NewDecoder(resp.Body).Decode(&appSess)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(appSess.AscRespData.WebsocketURI).ShouldNot(Equal(""))

				})

				Specify("Websocket Connect and Notification Receive",
					func() {

						// Connect to Websocket

						CACert, err := ioutil.ReadFile("/etc/certs/root-ca-cert.pem")
						Expect(err).ShouldNot(HaveOccurred())

						CACertPool := x509.NewCertPool()
						CACertPool.AppendCertsFromPEM(CACert)

						var socket = websocket.Dialer{
							ReadBufferSize:  512,
							WriteBufferSize: 512,
							TLSClientConfig: &tls.Config{
								RootCAs: CACertPool,
							},
						}

						var header = http.Header{}
						header["Origin"] = []string{"ConsumerID1"}

						conn := connectWsAF(&socket, &header)
						defer conn.Close()

						By("Reading json file")
						reqBody, err := ioutil.ReadFile(
							"./testdata/policy_auth/SMF_AF_NOTIF_ws.json")
						Expect(err).ShouldNot(HaveOccurred())

						By("Preparing request")
						reqBodyBytes := bytes.NewReader(reqBody)
						req, err := http.NewRequest(http.MethodPost,
							"http://localhost:8081/af/v1/policy-authorization/"+
								"smfnotify",
							reqBodyBytes)
						Expect(err).ShouldNot(HaveOccurred())

						By("Sending request")
						resp := httptest.NewRecorder()
						ctx := context.WithValue(req.Context(),
							KeyType("af-ctx"), af.AfCtx)
						af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

						Expect(resp.Code).To(Equal(http.StatusNoContent))
						var afEvent af.Afnotification
						getNotifyFromConn(conn, &afEvent, "1240")
						// Testing second connection, old connection is closed
						conn2 := connectWsAF(&socket, &header)
						defer conn2.Close()

					})

				Specify("Sending DELETE request", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XDELETE_01.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001/delete",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 204,
								// Send response to be tested
								Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
								// Must be set to non-nil value or it panics
								Header: make(http.Header),
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusNoContent))

				})

				Specify("Sending POST request - http uri", func() {
					By("Reading json file")

					// Start a local HTTP server for notifications
					http.HandleFunc("/notification", NotificationPost)
					handler := http.HandlerFunc(NotificationPost)

					server = httptest.NewUnstartedServer(handler)
					server.Config.ReadHeaderTimeout = 10 * time.Second
					server.Config.WriteTimeout = 10 * time.Second
					server.Start()

					fmt.Printf("Test Server : %s\n", server.URL)
					// Close the server when test finishes
					defer server.Close()

					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_1_corrID_https_uri.json")
					Expect(err).ShouldNot(HaveOccurred())

					var asc af.AppSessionContext
					err = json.Unmarshal(reqBody, &asc)
					Expect(err).ShouldNot(HaveOccurred())

					asc.AscReqData.AfRoutReq.UpPathChgSub.NotificationURI = server.URL + "/notification"

					reqBody, err = json.Marshal(asc)
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPost,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					resBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_1_corrID_https_uri.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					header := make(http.Header)
					header.Set("Location",
						"https://localhost:8095/af/v1/"+
							"policy-authorization/app-sessions/5001")
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 201,
								// Send response to be tested
								Body: ioutil.NopCloser(resBodyBytes),
								// Must be set to non-nil value or it panics
								Header: header,
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusCreated))

					// generate Notification

					By("Reading json file")
					reqBody, err = ioutil.ReadFile(
						"./testdata/policy_auth/SMF_AF_NOTIF_ws.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes = bytes.NewReader(reqBody)
					req, err = http.NewRequest(http.MethodPost,
						"http://localhost:8081/af/v1/policy-authorization/"+
							"smfnotify",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp = httptest.NewRecorder()
					ctx = context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

					Expect(resp.Code).To(Equal(http.StatusNoContent))

				})

				Specify("Sending PATCH update to ws ", func() {
					By("Reading json file")
					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPATCH_1_corrID_ws.json")
					Expect(err).ShouldNot(HaveOccurred())

					By("Preparing request")
					reqBodyBytes := bytes.NewReader(reqBody)
					req, err := http.NewRequest(http.MethodPatch,
						"http://localhost:8080/af/v1/policy-authorization/"+
							"app-sessions/5001",
						reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())

					By("Sending request")
					resp := httptest.NewRecorder()
					ctx := context.WithValue(req.Context(),
						KeyType("af-ctx"), af.AfCtx)
					resBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_1_corrID_ws.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					httpclient :=
						testingAFClient(func(req *http.Request) *http.Response {
							// Test request parameters
							return &http.Response{
								StatusCode: 200,
								// Send response to be tested
								Body: ioutil.NopCloser(resBodyBytes),
								// Must be set to non-nil value or it panics
								Header: make(http.Header),
							}
						})

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusOK))
					var appSess af.AppSessionContext
					// Decode response to check for websocketURI
					err = json.NewDecoder(resp.Body).Decode(&appSess)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(appSess.AscRespData.WebsocketURI).ShouldNot(Equal(""))
				})

				Specify("Incorrect Websocket Connect - 403",
					func() {

						// Connect to Websocket

						CACert, err := ioutil.ReadFile("/etc/certs/root-ca-cert.pem")
						Expect(err).ShouldNot(HaveOccurred())

						CACertPool := x509.NewCertPool()
						CACertPool.AppendCertsFromPEM(CACert)

						var socket = websocket.Dialer{
							ReadBufferSize:  512,
							WriteBufferSize: 512,
							TLSClientConfig: &tls.Config{
								RootCAs: CACertPool,
							},
						}

						var header = http.Header{}
						header["Origin"] = []string{"ConsumerID10"}

						connectWsAFForbidden(&socket, &header)

					})

			})

		})
	})

	Describe("Stop the AF Server", func() {
		It("Disconnect AF Server", func() {
			srvCancel()
			server.Close()
		})
	})
})
