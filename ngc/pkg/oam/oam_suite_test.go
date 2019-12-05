// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package oam

import (
	"bytes"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testdataBasepath = "../../test/oam/ngc-apistub-testdata/"
const postdataBasepath = "../../test/oam/cnca-cli-scripts/json/"

func TestOam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oam Suite")
}

var _ = Describe("NGC_Proxy", func() {

	Describe("Proxy init", func() {
		It("Will init proxy",
			func() {
				testPath := testdataBasepath + "testdata_00.json"
				Expect(InitProxy("v", "Flexcore", testPath)).NotTo(BeNil())
				Expect(InitProxy("v", "APISTUB", testPath)).To(BeNil())
			})
	})

	Describe("Proxy GetlAll", func() {
		It("Will use proxy to GetAll",
			func() {
				testPath := testdataBasepath + "testdata_00.json"
				Expect(InitProxy("valid", "Flexcore", testPath)).NotTo(BeNil())
				req, err := http.NewRequest("GET", "/services", nil)
				Expect(err).ShouldNot(HaveOccurred())
				rsp := httptest.NewRecorder()
				ProxyGetAll(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusNotFound))

			})
	})

	Describe("Proxy Get", func() {
		It("Will use proxy to Get",
			func() {
				testPath := testdataBasepath + "testdata_00.json"
				Expect(InitProxy("valid", "Flexcore", testPath)).NotTo(BeNil())
				req, err := http.NewRequest("GET", "/services", nil)
				Expect(err).ShouldNot(HaveOccurred())
				rsp := httptest.NewRecorder()
				ProxyGet(rsp, req)
				Expect(rsp.Code).NotTo(Equal(http.StatusOK))

			})
	})

	Describe("Proxy Del", func() {
		It("Will use proxy to Del",
			func() {
				testPath := testdataBasepath + "testdata_00.json"
				Expect(InitProxy("valid", "Flexcore", testPath)).NotTo(BeNil())
				req, err := http.NewRequest("DEL", "/services/1", nil)
				Expect(err).ShouldNot(HaveOccurred())
				rsp := httptest.NewRecorder()
				ProxyDel(rsp, req)
				Expect(rsp.Code).NotTo(Equal(http.StatusOK))

			})
	})

	Describe("Proxy Add", func() {
		It("Will use proxy to Add",
			func() {
				testPath := testdataBasepath + "testdata_00.json"
				Expect(InitProxy("valid", "Flexcore", testPath)).NotTo(BeNil())
				req, err := http.NewRequest("POST", "/services/1", nil)
				Expect(err).ShouldNot(HaveOccurred())
				rsp := httptest.NewRecorder()
				ProxyAdd(rsp, req)
				Expect(rsp.Code).NotTo(Equal(http.StatusOK))

			})
	})

	Describe("Proxy Update", func() {
		It("Will use proxy to Update",
			func() {
				testPath := testdataBasepath + "testdata_00.json"
				Expect(InitProxy("valid", "Flexcore", testPath)).NotTo(BeNil())
				req, err := http.NewRequest("PATCH", "/services/1", nil)
				Expect(err).ShouldNot(HaveOccurred())
				rsp := httptest.NewRecorder()
				ProxyUpdate(rsp, req)
				Expect(rsp.Code).NotTo(Equal(http.StatusOK))

			})
	})
})

var _ = Describe("NGC_APIStub", func() {

	BeforeEach(func() {
		APIStubReset()
	})

	AfterEach(func() {
		APIStubReset()
	})

	Describe("APISTUB init", func() {
		It("Will init APSTUB",
			func() {
				Expect(APIStubInit("nonexistent-file")).NotTo(BeNil())
				Expect(APIStubInit("conf")).NotTo(BeNil())
				tmp := testdataBasepath + "testdata_00.json"
				Expect(APIStubInit(tmp)).To(BeNil())
				Expect(len(AllRecords)).To(Equal(0))
				Expect(NewRecordAFServiceID).To(Equal(AFServiceIDBaseValue))
				tmp = testdataBasepath + "testdata_01.json"
				Expect(APIStubInit(tmp)).To(BeNil())
				Expect(len(AllRecords)).To(Equal(1))
				Expect(NewRecordAFServiceID).To(Equal(AFServiceIDBaseValue + 1))
			})
	})

	Describe("APISTUB reset", func() {
		It("Will reset APSTUB",
			func() {
				Expect(APIStubReset()).To(BeNil())
			})
	})

	Describe("APISTUB Add", func() {
		It("Will Add new Record",
			func() {
				tmp := postdataBasepath + "POST001.json"
				reqBody, err := ioutil.ReadFile(tmp)
				Expect(err).ShouldNot(HaveOccurred())
				dats := bytes.NewReader(reqBody)
				req, _ := http.NewRequest(http.MethodPost, "/services", dats)
				rsp := httptest.NewRecorder()
				expected := "{\"afServiceId\":\"123457\"}"
				APIStubAdd(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusOK))
				Expect(rsp.Body.String()).To(Equal(expected))

				tmp = postdataBasepath + "POST002.json"
				reqBody, err = ioutil.ReadFile(tmp)
				Expect(err).ShouldNot(HaveOccurred())
				dats = bytes.NewReader(reqBody)
				req, _ = http.NewRequest(http.MethodPost, "/services", dats)
				rsp = httptest.NewRecorder()
				expected = "{\"afServiceId\":\"123458\"}"
				APIStubAdd(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusOK))
				Expect(rsp.Body.String()).To(Equal(expected))

				tmp = postdataBasepath + "POST003.json"
				reqBody, err = ioutil.ReadFile(tmp)
				Expect(err).ShouldNot(HaveOccurred())
				dats = bytes.NewReader(reqBody)
				req, _ = http.NewRequest(http.MethodPost, "/services", dats)
				rsp = httptest.NewRecorder()
				expected = "{\"afServiceId\":\"123459\"}"
				APIStubAdd(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusOK))
				Expect(rsp.Body.String()).To(Equal(expected))

			})
	})

	Describe("APISTUB Update", func() {
		It("Will Update Record",
			func() {
				tmp := testdataBasepath + "testdata_01.json"
				tmpP := postdataBasepath + "POST001.json"
				APIStubInit(tmp)
				reqBody, err := ioutil.ReadFile(tmpP)
				Expect(err).ShouldNot(HaveOccurred())
				reqBodyBytes := bytes.NewReader(reqBody)
				req, _ := http.NewRequest("PATCH", "/services/1", reqBodyBytes)
				vars := map[string]string{
					"afServiceId": "123457",
				}
				req = mux.SetURLVars(req, vars)
				rsp := httptest.NewRecorder()
				APIStubUpdate(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusOK))

				req, _ = http.NewRequest("PATCH", "/services/2", nil)
				vars = map[string]string{
					"afServiceId": "123458",
				}
				req = mux.SetURLVars(req, vars)
				rsp = httptest.NewRecorder()
				APIStubDel(rsp, req)
				Expect(rsp.Code).NotTo(Equal(http.StatusOK))

			})
	})

	Describe("APISTUB Del", func() {
		It("Will Delete Record",
			func() {
				tmp := testdataBasepath + "testdata_01.json"
				APIStubInit(tmp)
				req, _ := http.NewRequest("DELETE", "/services/123457", nil)
				vars := map[string]string{
					"afServiceId": "123457",
				}
				req = mux.SetURLVars(req, vars)
				rsp := httptest.NewRecorder()
				APIStubDel(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusOK))
				Expect(len(AllRecords)).To(Equal(0))

				req, _ = http.NewRequest("DELETE", "/services/123457", nil)
				vars = map[string]string{
					"afServiceId": "123457",
				}
				req = mux.SetURLVars(req, vars)
				rsp = httptest.NewRecorder()
				APIStubDel(rsp, req)
				Expect(rsp.Code).NotTo(Equal(http.StatusOK))

			})
	})

	Describe("APISTUB Get", func() {
		It("Will Get one Record",
			func() {
				tmp := testdataBasepath + "testdata_01.json"
				APIStubInit(tmp)
				req, _ := http.NewRequest("GET", "/services/123457", nil)
				vars := map[string]string{
					"afServiceId": "123457",
				}
				req = mux.SetURLVars(req, vars)
				rsp := httptest.NewRecorder()
				APIStubGet(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusOK))

				req, _ = http.NewRequest("GET", "/services/123458", nil)
				vars = map[string]string{
					"afServiceId": "123458",
				}
				req = mux.SetURLVars(req, vars)
				rsp = httptest.NewRecorder()
				APIStubGet(rsp, req)
				Expect(rsp.Code).NotTo(Equal(http.StatusOK))

			})
	})

	Describe("APISTUB Getll", func() {
		It("Will GetAll Records",
			func() {
				tmp := testdataBasepath + "testdata_01.json"
				APIStubInit(tmp)
				req, err := http.NewRequest("GET", "/services", nil)
				Expect(err).ShouldNot(HaveOccurred())
				rsp := httptest.NewRecorder()
				APIStubGetAll(rsp, req)
				Expect(rsp.Code).To(Equal(http.StatusOK))

			})
	})
})
