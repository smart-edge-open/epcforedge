//SPDX-License-Identifier: Apache-2.0
//Copyright © 2019-2020 Intel Corporation

package af_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

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
	})
})
