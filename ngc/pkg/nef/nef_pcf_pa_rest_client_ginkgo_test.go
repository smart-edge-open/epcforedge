package ngcnef_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
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
		Context("Initializing PCF client", func() {
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
	})
	Describe("client request methods to PCF", func() {
		Specify("Sending valid POST request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextpostreq.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContext
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextpostresp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Location", "1234test")
					return &http.Response{
						StatusCode: 201,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			appid, pcr, err1 := pcfc.PolicyAuthorizationCreate(ctx, ascreq)
			Expect(err1).ShouldNot(HaveOccurred())
			Expect(string(appid)).To(Equal("1234test"))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusCreated))
			Expect(pcr.Asc).ToNot(Equal(nil))
		})

		Specify("Sending valid PATCH request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextpatchreq.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContextUpdateData
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextpatchresp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters

					return &http.Response{
						StatusCode: 200,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			pcr, err1 := pcfc.PolicyAuthorizationUpdate(ctx, ascreq, ngcnef.AppSessionID("test1234"))
			Expect(err1).ShouldNot(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusOK))

		})
		Specify("Sending valid GET request ", func() {

			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextgetresp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters
					return &http.Response{
						StatusCode: 200,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			pcr, err1 := pcfc.PolicyAuthorizationGet(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).ShouldNot(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusOK))

		})
		Specify("Sending valid DELETE request ", func() {

			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters
					return &http.Response{
						StatusCode: 204,
						// Send response to be tested
						Body: ioutil.NopCloser(nil),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			pcr, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).ShouldNot(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNoContent))

		})
		Specify("Sending invalid POST request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextpostreq.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContext
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextposterrresp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters
					respHeader := make(http.Header)
					respHeader.Set("Location", "1234test")
					return &http.Response{
						StatusCode: 403,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: respHeader,
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			appid, pcr, err1 := pcfc.PolicyAuthorizationCreate(ctx, ascreq)
			Expect(err1).Should(HaveOccurred())
			Expect(string(appid)).To(Equal(""))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusForbidden))
			Expect(pcr.Pd).ToNot(Equal(nil))
		})

		Specify("Sending invalid PATCH request ", func() {
			By("Reading json file")
			reqBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextpatchreq.json")
			Expect(err).ShouldNot(HaveOccurred())

			By("Preparing request")
			var ascreq ngcnef.AppSessionContextUpdateData
			err = json.Unmarshal(reqBody, &ascreq)
			Expect(err).ShouldNot(HaveOccurred())
			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextpatcherrresp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters

					return &http.Response{
						StatusCode: 404,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			pcr, err1 := pcfc.PolicyAuthorizationUpdate(ctx, ascreq, ngcnef.AppSessionID("test1234"))
			Expect(err1).Should(HaveOccurred())
			Expect(pcr.Pd).ToNot(Equal(nil))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNotFound))

		})
		Specify("Sending invalid GET request ", func() {

			By("Preparing response")
			resBody, err := ioutil.ReadFile(
				testJSONPath + "appcontextgeterrresp.json")
			Expect(err).ShouldNot(HaveOccurred())
			resBodyBytes := bytes.NewReader(resBody)
			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters
					return &http.Response{
						StatusCode: 404,
						// Send response to be tested
						Body: ioutil.NopCloser(resBodyBytes),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			pcr, err1 := pcfc.PolicyAuthorizationGet(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).Should(HaveOccurred())
			Expect(pcr.Pd).ToNot(Equal(nil))
			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNotFound))

		})
		Specify("Sending invalid DELETE request ", func() {

			httpclient :=
				testingPCFClient(func(req *http.Request) *http.Response {
					// Test request parameters
					return &http.Response{
						StatusCode: 404,
						// Send response to be tested
						Body: ioutil.NopCloser(nil),
						// Must be set to non-nil value or it panics
						Header: make(http.Header),
					}
				})

			ngcnef.SetHTTPClient(httpclient)
			pcfc := ngcnef.PcfClient{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
			pcr, err1 := pcfc.PolicyAuthorizationDelete(ctx, ngcnef.AppSessionID("1234test"))
			Expect(err1).Should(HaveOccurred())

			Expect(int(pcr.ResponseCode)).To(Equal(http.StatusNotFound))

		})
	})
})
