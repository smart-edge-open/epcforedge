// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oam

import (
        "bytes"
        "io/ioutil"
        "net/http"
        "net/http/httptest"
        "github.com/gorilla/mux"
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oam Suite")
}

var _ = Describe("NGC_Proxy", func() {

        Describe("Proxy init", func() {
                It("Will init proxy",
                        func() {
                                testPath := "testdata/testdata_00.json"
                                Expect(InitProxy("valid", "Flexcore", testPath)).NotTo(BeNil())
                                Expect(InitProxy("valid", "APISTUB", testPath)).To(BeNil())
                        })
        })

        Describe("Proxy GetlAll", func() {
                It("Will use proxy to GetAll",
                        func() {
                                testPath := "testdata/testdata_00.json"
                                Expect(InitProxy("valid", "Flexcore", testPath)).NotTo(BeNil())
                                req, err := http.NewRequest("GET", "/services", nil)
                                if err != nil {
                                }
                                rsp := httptest.NewRecorder()
                                ProxyGetAll(rsp,req)
                                Expect(rsp.Code).To(Equal(404))

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
                                Expect(APIStubInit("testdata/testdata_00.json")).To(BeNil())
                                Expect(len(AllRecords)).To(Equal(0))
                                Expect(len(AllRecordsAfId)).To(Equal(0))
                                Expect(NewRecordAfId).To(Equal(0))
                                Expect(APIStubInit("testdata/testdata_01.json")).To(BeNil())
                                Expect(len(AllRecords)).To(Equal(1))
                                Expect(len(AllRecordsAfId)).To(Equal(1))
                                Expect(NewRecordAfId).To(Equal(1))
                        })
        })

        Describe("APISTUB Add", func() {
                It("Will Add new Record",
                        func() {
                                reqBody, err := ioutil.ReadFile("testdata/POST001.json") 
                                if err != nil {
                                }
                                reqBodyBytes := bytes.NewReader(reqBody)
                                req, _ := http.NewRequest(http.MethodPost,"/services",reqBodyBytes)
                                rsp := httptest.NewRecorder()
                                expected := "{\"afid\":\"1\"}"
                                APIStubAdd(rsp, req)
                                Expect(rsp.Code).To(Equal(http.StatusOK))
                                Expect(rsp.Body.String()).To(Equal(expected))
                               
                        })
        })

        Describe("APISTUB Update", func() {
                It("Will Update Record",
                        func() {
                                APIStubInit("testdata/testdata_01.json")
                                reqBody, err := ioutil.ReadFile("testdata/POST001.json") 
                                if err != nil {
                                }
                                reqBodyBytes := bytes.NewReader(reqBody)
                                req, _ := http.NewRequest("PATCH","/services/1",reqBodyBytes)
                                vars := map[string]string{
                                     "afId": "1",
                                }                               
                                req = mux.SetURLVars(req, vars)
                                rsp := httptest.NewRecorder()
                                APIStubUpdate(rsp, req)
                                Expect(rsp.Code).To(Equal(http.StatusOK))
                               
                        })
       })
 
        Describe("APISTUB Del", func() {
                It("Will Delete Record",
                        func() {
                                APIStubInit("testdata/testdata_01.json")
                                reqBody, err := ioutil.ReadFile("testdata/POST001.json") 
                                if err != nil {
                                }
                                reqBodyBytes := bytes.NewReader(reqBody)
                                req, _ := http.NewRequest("DELETE","/services/1",reqBodyBytes)
                                vars := map[string]string{
                                     "afId": "1",
                                }                               
                                req = mux.SetURLVars(req, vars)
                                rsp := httptest.NewRecorder()
                                APIStubDel(rsp, req)
                                Expect(rsp.Code).To(Equal(http.StatusOK))
                               
                        })
       })

       Describe("APISTUB Getll", func() {
                It("Will GetAll Records",
                        func() {
                                APIStubInit("testdata/testdata_01.json")
                                req, err := http.NewRequest("GET", "/services", nil)
                                if err != nil {
                                }
                                rsp := httptest.NewRecorder()
                                expected := "[{\"afInstance\":\"mike\",\"localServices\":"+
                                  "[{\"dnai\":\"mike1_dnai\","+
                                    "\"dnn\":\"mike1_dnn\",\"dns\":\"192.168.9.9\"},"+
                                   "{\"dnai\":\"mike2_dnai\","+
                                    "\"dnn\":\"mike2_dnn\",\"dns\":\"192.168.8.8\"}]}]"
 
                                APIStubGetAll(rsp, req)
                                Expect(rsp.Code).To(Equal(http.StatusOK))
                                Expect(rsp.Body.String()).To(Equal(expected))
                               
                        })
        })
})

