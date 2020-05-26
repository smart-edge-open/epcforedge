//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019-2020 Intel Corporation

package af_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otcshare/epcforedge/ngc/pkg/af"
)

var rootURL string = "https://localhost:8080/af/v1/policy-authorization/" +
	"app-sessions"

var notifRootURL string = "https://localhost:8081/af/v1/policy-authorization"

var jsonFilePath string = "./testdata/policy_auth/"

var locationPrefix string = "http://localhost:8080/npcf-policyauthorization/" +
	"v1/app-sessions/"

func createTestHTTPClient(statusCode int, body io.Reader,
	header http.Header) *http.Client {

	return testingAFClient(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: statusCode,
			// Send response to be tested
			Body: ioutil.NopCloser(body),
			// Must be set to non-nil value or it panics
			Header: header,
		}
	})
}

func NotificationPost(w http.ResponseWriter, r *http.Request) {

	defer GinkgoRecover()
	Expect(r.Body).ShouldNot(Equal(nil))
	log.Println("Notification Protocol:     Method", r.Proto, r.Method)
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// connectConsumer sends a consumer notifications GET request to the appliance
func connectWsAF(socket *websocket.Dialer, header *http.Header) *websocket.Conn {
	By("Sending consumer notification GET request")
	conn, resp, err := socket.Dial(testSrvData.wsURI, *header)
	Expect(err).ShouldNot(HaveOccurred())

	By("Comparing GET response code")
	defer resp.Body.Close()
	Expect(resp.Status).To(Equal("101 Switching Protocols"))

	return conn
}

// connectConsumer sends a consumer notifications GET request to the appliance
func connectWsAFForbidden(socket *websocket.Dialer, header *http.Header) {
	By("Sending consumer notification GET request")
	_, resp, _ := socket.Dial(testSrvData.wsURI, *header)
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

var _ = Describe("AF PA", func() {

	Describe("Cnca  AF PA : ", func() {

		Context("PA AppSessionCtx Create", func() {
			Specify("Sending Invalid POST 001 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 002 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_002.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 003 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_003.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 004 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_004.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 005 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_005.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 006 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_006.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 007 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_007.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 008 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_008.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 009 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_009.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 010 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_010.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 010_1 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_010_1.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 011 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_011.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 012 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_012.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 013 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_013.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 014 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_014.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending Invalid POST 015 request", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_POST_015.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending POST valid request - [Resp] 201: no Location URL", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBodyBytes := reqBodyBytes
				header := make(http.Header)
				httpclient := createTestHTTPClient(201,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("Sending POST valid request - [Resp] 201: valid", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				header.Set("Location", locationPrefix+"5000")
				httpclient := createTestHTTPClient(201,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusCreated))

			})

			Specify("Sending POST valid request - [Resp] valid 303", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				header.Set("Location", locationPrefix+"5000")
				httpclient := createTestHTTPClient(303,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusSeeOther))

			})

			Specify("Sending POST valid request - [Resp] valid 401", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL,
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"problem_details.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(401,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusUnauthorized))

			})

		})

		Context("Application session context GET", func() {
			Specify("Sending PA GET request, [Resp] 200", func() {
				By("Preparing request")
				req, err := http.NewRequest(http.MethodGet,
					rootURL+"/5001",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("Sending GET request: [Resp] 401", func() {
				By("Preparing request")
				req, err := http.NewRequest(http.MethodGet,
					rootURL+"/4999",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"problem_details.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(401,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusUnauthorized))

			})

			Specify("Sending GET request: [Resp] 408", func() {
				By("Preparing request")
				req, err := http.NewRequest(http.MethodGet,
					rootURL+"/4999",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(408,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusRequestTimeout))

			})

		})

		Context("PA AppSessionCtx Patch", func() {
			Specify("Sending PA PATCH 001 Request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_PATCH_001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"https://localhost:8080/af/v1/policy-authorization/app-sessions/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending PA PATCH 002 Request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_PATCH_002.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"https://localhost:8080/af/v1/policy-authorization/app-sessions/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending PA PATCH 003 Request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_PATCH_003.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"https://localhost:8080/af/v1/policy-authorization/app-sessions/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending PA PATCH 004 Request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_PATCH_004.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					"https://localhost:8080/af/v1/policy-authorization/app-sessions/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending PATCH valid request - [Resp] 200: valid", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPATCH_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					rootURL+"/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("Sending PATCH valid request - [Resp] 403: no retry after header", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPATCH_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					rootURL+"/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"problem_details.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(403,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusInternalServerError))

			})

			Specify("Sending PATCH valid request - [Resp] 403: valid", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPATCH_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					rootURL+"/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"problem_details.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				header.Set("Retry-After", "5000")
				httpclient := createTestHTTPClient(403,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusForbidden))

			})

			Specify("Sending PATCH valid request - [Resp] 408", func() {
				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPATCH_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPatch,
					rootURL+"/5001",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				header.Set("Retry-After", "5000")
				httpclient := createTestHTTPClient(408,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false
				Expect(resp.Code).To(Equal(http.StatusRequestTimeout))

			})
		})

		Context("PA AppSessionCtx Delete", func() {
			Specify("Sending PA DELETE 001 request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_DELETE_001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL+"/5002/delete",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending PA DELETE 002 request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_DELETE_002.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL+"/5002/delete",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending valid DELETE request, [Resp] 200", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XDELETE_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL+"/5001/delete",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(200,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusOK))

			})

			Specify("Sending valid DELETE request, [Resp] 204", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XDELETE_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL+"/5001/delete",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(204,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusNoContent))

			})

			Specify("Sending valid DELETE request, [Resp] 401: empty resp body", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XDELETE_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL+"/5001/delete",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(401,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusUnauthorized))

			})

			Specify("Sending valid DELETE request, [Resp] 408", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_XDELETE_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					rootURL+"/5001/delete",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(408,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusRequestTimeout))

			})
		})

		Context("PA event subscribe", func() {
			Specify("Sending EVENT PUT 001 request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EVENT_PUT_001.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					rootURL+"/5001/events-subscription",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending EVENT PUT 002 request", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EVENT_PUT_002.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					rootURL+"/5001/events-subscription",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending EVENT PUT valid request, [Resp] 201", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EVENT_XPUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					rootURL+"/5001/events-subscription",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_SB_EVENTS_PUT_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				header.Set("Location", locationPrefix+"5000/events-subscription")
				httpclient := createTestHTTPClient(201,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusCreated))

			})

			Specify("Sending EVENT PUT valid request, [Resp] 204", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EVENT_XPUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					rootURL+"/5001/events-subscription",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(204,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusNoContent))

			})

			Specify("Sending EVENT PUT valid request, [Resp] 401", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EVENT_XPUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					rootURL+"/5001/events-subscription",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"problem_details.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(401,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusUnauthorized))

			})

			Specify("Sending EVENT PUT valid request, [Resp] 408", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EVENT_XPUT_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPut,
					rootURL+"/5001/events-subscription",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(408,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusRequestTimeout))

			})

		})

		Context("PA event unsubscribe", func() {
			Specify("Sending valid event DELETE req, [Resp] 204", func() {

				By("Preparing request")
				req, err := http.NewRequest(http.MethodDelete,
					rootURL+"/5001/events-subscription",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(204,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})

			Specify("Sending valid event DELETE req, [Resp] 401", func() {

				By("Preparing request")
				req, err := http.NewRequest(http.MethodDelete,
					rootURL+"/5001/events-subscription",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"problem_details.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(401,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			})

			Specify("Sending valid event DELETE req, [Resp] 408", func() {

				By("Preparing request")
				req, err := http.NewRequest(http.MethodDelete,
					rootURL+"/5001/events-subscription",
					nil)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				resBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_NB_PA_EMPTY_RESP.json")
				Expect(err).ShouldNot(HaveOccurred())
				resBodyBytes := bytes.NewReader(resBody)

				header := make(http.Header)
				httpclient := createTestHTTPClient(408,
					resBodyBytes, header)

				af.TestAf = true
				af.SetHTTPClient(httpclient)
				af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
				af.TestAf = false

				Expect(resp.Code).To(Equal(http.StatusRequestTimeout))
			})

		})

		Context("PA event notify", func() {

			Specify("Sending invalid Notify POST request ", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_SB_PCF_NOTIF_POST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					notifRootURL+"/notify",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending valid Notify POST request, [Resp] 204 ", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_SB_PCF_NOTIF_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					notifRootURL+"/notify",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNoContent))

			})

			Specify("Sending invalid Notif Terminate POST 01 request ", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_SB_PCF_NOTIF_TERM_POST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					notifRootURL+"/terminate",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending invalid Notif Terminate POST 02 request ", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_SB_PCF_NOTIF_TERM_POST_02.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					notifRootURL+"/terminate",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

			})

			Specify("Sending valid Notify Term POST 01 request, [Resp] 204 ", func() {

				By("Reading json file")
				reqBody, err := ioutil.ReadFile(jsonFilePath +
					"AF_SB_PCF_NOTIF_TERM_XPOST_01.json")
				Expect(err).ShouldNot(HaveOccurred())

				By("Preparing request")
				reqBodyBytes := bytes.NewReader(reqBody)
				req, err := http.NewRequest(http.MethodPost,
					notifRootURL+"/terminate",
					reqBodyBytes)
				Expect(err).ShouldNot(HaveOccurred())

				By("Sending request")
				resp := httptest.NewRecorder()
				ctx := context.WithValue(req.Context(),
					KeyType("af-ctx"), af.AfCtx)
				af.NotifRouter.ServeHTTP(resp, req.WithContext(ctx))

				Expect(resp.Code).To(Equal(http.StatusNoContent))

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

					httpclient := createTestHTTPClient(201, resBodyBytes, header)

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
					httpclient := createTestHTTPClient(200, resBodyBytes,
						make(http.Header))

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

					httpclient := createTestHTTPClient(204,
						bytes.NewBufferString(`OK`), make(http.Header))

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

					httpclient := createTestHTTPClient(201,
						resBodyBytes, header)

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
					httpclient := createTestHTTPClient(200,
						resBodyBytes, make(http.Header))

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

					httpclient := createTestHTTPClient(204,
						bytes.NewBufferString(`OK`), make(http.Header))

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
					httpclient := createTestHTTPClient(201,
						resBodyBytes, header)

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

					testSrvData.wsURI = appSess.AscRespData.WebsocketURI
					log.Println("Websocket URI received is", testSrvData.wsURI)

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
					httpclient := createTestHTTPClient(204,
						bytes.NewBufferString(`OK`), make(http.Header))

					af.TestAf = true
					af.SetHTTPClient(httpclient)
					af.AfRouter.ServeHTTP(resp, req.WithContext(ctx))
					af.TestAf = false
					Expect(resp.Code).To(Equal(http.StatusNoContent))

				})

				Specify("Sending POST request - http uri", func() {
					By("Reading json file")

					reqBody, err := ioutil.ReadFile(
						"./testdata/policy_auth/AF_NB_PA_XPOST_1_corrID_http_uri.json")
					Expect(err).ShouldNot(HaveOccurred())

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
						"./testdata/policy_auth/AF_NB_PA_XPOST_1_corrID_http_uri.json")
					Expect(err).ShouldNot(HaveOccurred())
					resBodyBytes := bytes.NewReader(resBody)
					header := make(http.Header)
					header.Set("Location",
						"https://localhost:8095/af/v1/"+
							"policy-authorization/app-sessions/5001")
					httpclient := createTestHTTPClient(201,
						resBodyBytes, header)
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
					httpclient := createTestHTTPClient(200,
						resBodyBytes, make(http.Header))

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
})
