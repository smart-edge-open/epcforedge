/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */
package ngcnef_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testingPCFClient(fn RoundTripFunc) *http.Client {

	return &http.Client{
		Transport: fn,
	}

}

var _ = Describe("NefPcfPaRestClient", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	Describe("client request methods to PCF", func() {
		Context("Initializing PCF client with HTTP 2.0/https", func() {
			It("Will init NefServer",
				func() {
					ctx, cancel = context.WithCancel(context.Background())

					defer cancel()
					go func() {
						ngcnef.TestPcf = true
						err := ngcnef.Run(ctx, NefTestCfgBasepath+"valid_pcf.json")
						Expect(err).To(BeNil())

					}()
					time.Sleep(2 * time.Second)
				})
		})
		Context("Initializing PCF client with HTTP 2.0/http", func() {
			It("Will init NefServer",
				func() {
					ctx, cancel = context.WithCancel(context.Background())

					defer cancel()
					go func() {

						err := ngcnef.Run(ctx, NefTestCfgBasepath+"valid_pcf_2.json")
						Expect(err).To(BeNil())

					}()
					time.Sleep(2 * time.Second)
				})
		})
		Context("Initializing PCF client with HTTP 1.1/https", func() {
			It("Will init NefServer",
				func() {
					ctx, cancel = context.WithCancel(context.Background())

					defer cancel()
					go func() {

						err := ngcnef.Run(ctx, NefTestCfgBasepath+"valid_pcf_1.json")
						Expect(err).To(BeNil())

					}()
					time.Sleep(2 * time.Second)
				})
		})

	})
	Describe("client request methods to PCF", func() {
		Specify("Sending valid POST request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_post_req.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContext
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_post_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Location", "1234test")
					return &http.Response{
						StatusCode: 201,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			appid, pcr, err1 := pcfc.PolicyAuthorizationCreate(ctx, ascreq)
			Expect(err1).ShouldNot(HaveOccurred())
			Expect(string(appid)).To(Equal("1234test"))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusCreated))
			Expect(pcr.Asc).ToNot(Equal(nil))
		})

		Specify("Sending valid PATCH request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_patch_req.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContextUpdateData
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_patch_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters

					return &http.Response{
						StatusCode: 200,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationUpdate(ctx, ascreq, ngcnef.AppSessionID("test1234"))
			Expect(err1).ShouldNot(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusOK))

		})
		Specify("Sending valid GET request ", func() {

			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_get_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{
						StatusCode: 200,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationGet(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).ShouldNot(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusOK))

		})
		Specify("Sending valid DELETE request ", func() {

			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{
						StatusCode: 204,
						// Send response to be tested
						Body: ioutil.NopCloser(nil),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).ShouldNot(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNoContent))

		})
		Specify("Sending valid 200 DELETE request ", func() {
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_del_200_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{
						StatusCode: 200,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).ShouldNot(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusOK))

		})
		Specify("Sending invalid POST request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_post_req.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContext
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_post_err_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Location", "1234test")
					respHeader.Set("Content-Type", "application/problem+json")
					return &http.Response{
						StatusCode: 403,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			appid, pcr, err1 := pcfc.PolicyAuthorizationCreate(ctx, ascreq)
			Expect(err1).Should(HaveOccurred())
			Expect(string(appid)).To(Equal(""))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusForbidden))
			Expect(pcr.Pd).ToNot(Equal(nil))
		})
		Specify("Sending invalid POST request for 401 response", func() {

			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_post_err_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Location", "1234test")
					respHeader.Set("Content-Type", "application/problem+json")
					return &http.Response{
						StatusCode: 401,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			appid, pcr, err1 := pcfc.PolicyAuthorizationCreate(ctx, ngcnef.AppSessionContext{})
			Expect(err1).Should(HaveOccurred())
			Expect(string(appid)).To(Equal(""))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusUnauthorized))
			Expect(pcr.Pd).ToNot(Equal(nil))
		})
		Specify("Sending invalid POST request with blank oath2 token", func() {

			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{}, nil

				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: ""}
			appid, _, err1 := pcfc.PolicyAuthorizationCreate(ctx, ngcnef.AppSessionContext{})
			Expect(err1).Should(HaveOccurred())
			Expect(string(appid)).To(Equal(""))

		})

		Specify("Sending invalid PATCH request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_patch_req.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContextUpdateData
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_patch_err_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Content-Type", "application/problem+json")
					return &http.Response{
						StatusCode: 404,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}, nil
				})
			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationUpdate(ctx, ascreq, ngcnef.AppSessionID("test1234"))
			Expect(err1).Should(HaveOccurred())
			Expect(pcr.Pd).ToNot(Equal(nil))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNotFound))

		})
		Specify("Sending invalid PATCH request for 401 response", func() {

			resBodyBytes := bytes.NewReader([]byte("test"))
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{
						StatusCode: 401,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationUpdate(ctx, ngcnef.AppSessionContextUpdateData{},
				ngcnef.AppSessionID(""))
			Expect(err1).Should(HaveOccurred())
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusUnauthorized))

		})
		Specify("Sending invalid PATCH request with blank oath2 token", func() {

			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{}, nil

				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: ""}
			_, err1 := pcfc.PolicyAuthorizationUpdate(ctx, ngcnef.AppSessionContextUpdateData{},
				ngcnef.AppSessionID(""))
			Expect(err1).Should(HaveOccurred())

		})
		Specify("Sending invalid GET request ", func() {

			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_get_err_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Content-Type", "application/problem+json")
					return &http.Response{
						StatusCode: 404,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationGet(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).Should(HaveOccurred())
			Expect(pcr.Pd).ToNot(Equal(nil))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNotFound))
		})
		Specify("Sending invalid GET request for 401 response", func() {

			resBodyBytes := bytes.NewReader([]byte("test"))
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{
						StatusCode: 401,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationGet(ctx, ngcnef.AppSessionID(""))
			Expect(err1).Should(HaveOccurred())
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusUnauthorized))

		})
		Specify("Sending invalid GET request with blank oath2 token", func() {

			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{}, nil

				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: ""}
			_, err1 := pcfc.PolicyAuthorizationGet(ctx, ngcnef.AppSessionID(""))
			Expect(err1).Should(HaveOccurred())

		})
		Specify("Sending invalid DELETE request ", func() {
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				NefTestCfgBasepath + "pcf_pa_get_err_resp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Content-Type", "application/problem+json")
					return &http.Response{
						StatusCode: 404,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).Should(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNotFound))

		})
		Specify("Sending invalid DELETE request for 401 response", func() {

			resBodyBytes := bytes.NewReader([]byte("test"))
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{
						StatusCode: 401,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}, nil
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			pcr, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID(""))
			Expect(err1).Should(HaveOccurred())
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusUnauthorized))

		})
		Specify("Sending invalid DELETE request with blank oath2 token", func() {

			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					return &http.Response{}, nil

				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: ""}
			_, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID(""))
			Expect(err1).Should(HaveOccurred())

		})
		Specify("Checking server timeout for GET request ", func() {
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					err := errors.New("no response from server")
					return nil, err
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			_, err1 := pcfc.PolicyAuthorizationGet(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).Should(HaveOccurred())

		})
		Specify("Checking server timeout for POST request ", func() {
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					err := errors.New("no response from server")
					return nil, err
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			_, _, err1 := pcfc.PolicyAuthorizationCreate(ctx, ngcnef.AppSessionContext{})
			Expect(err1).Should(HaveOccurred())

		})
		Specify("Checking server timeout for PATCH request ", func() {
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					err := errors.New("no response from server")
					return nil, err
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			appid := ngcnef.AppSessionID("")
			_, err1 := pcfc.PolicyAuthorizationUpdate(ctx, ngcnef.AppSessionContextUpdateData{}, appid)
			Expect(err1).Should(HaveOccurred())

		})
		Specify("Checking server timeout for DELETE request ", func() {
			httpclient :=
				testingPCFClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					err := errors.New("no response from server")
					return nil, err
				})

			pcfc := ngcnef.PcfClient{HTTPClient: httpclient, RootURI: "https://localhost:29507",
				ResourceURI: "/npcf-policyauthorization/v1/app-sessions/",
				Pcfcfg:      &ngcnef.PcfPolicyAuthorizationConfig{OAuth2Support: true},
				OAuth2Token: "teststring"}
			_, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).Should(HaveOccurred())

		})
	})
})
