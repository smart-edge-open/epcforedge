/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */
package ngcnef_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
)

const basePFDAPIURL = "http://localhost:8091/3gpp-pfd-management/" +
	"v1/AF_01/transactions"
const baseInvalidPFDAPIURL = "http://localhost:8091/3gpp-pfd-management/" +
	"v1/AF_02/transactions"

const testJSONPFDPath = "../../test/nef/nef-cli-scripts/pfd/json/"

func CreatePFDReqForNEF(ctx context.Context, method string, pfdTrans string,
	appID string, body []byte) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request
	if len(pfdTrans) > 0 {
		if len(appID) > 0 {

			if body != nil {
				//PUT/POST

				req, _ = http.NewRequest(method,
					basePFDAPIURL+"/"+pfdTrans+"/applications/"+appID,
					bytes.NewBuffer(body))
			} else {
				fmt.Println("Application ID is ", appID)
				//GET DELETE
				req, _ = http.NewRequest(method,
					basePFDAPIURL+"/"+pfdTrans+"/applications/"+appID,
					nil)
			}

		} else {

			if body != nil {
				//PUT
				req, _ = http.NewRequest(method, basePFDAPIURL+"/"+pfdTrans,
					bytes.NewBuffer(body))
			} else {
				//GET PFD/ DELETE

				req, _ = http.NewRequest(method, basePFDAPIURL+"/"+pfdTrans,
					nil)
			}

		}

	} else {
		if body != nil {
			//POST
			req, _ = http.NewRequest(method, basePFDAPIURL,
				bytes.NewBuffer(body))
		} else {
			//GET ALL
			req, _ = http.NewRequest(method, basePFDAPIURL, nil)
		}
	}

	rr := httptest.NewRecorder()
	return rr, req
}

func CreateInvalidPFDReqForNEF(ctx context.Context, method string,
	pfdTrans string,
	appID string, body []byte) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request
	if len(pfdTrans) > 0 {
		if len(appID) > 0 {

			if body != nil {
				//PUT/POST

				req, _ = http.NewRequest(method,
					baseInvalidPFDAPIURL+"/"+pfdTrans+"/applications/"+appID,
					bytes.NewBuffer(body))
			} else {
				fmt.Println("Application ID is ", appID)
				//GET DELETE
				req, _ = http.NewRequest(method,
					baseInvalidPFDAPIURL+"/"+pfdTrans+"/applications/"+appID,
					nil)
			}

		} else {

			if body != nil {
				//PUT
				req, _ = http.NewRequest(method, baseInvalidPFDAPIURL+
					"/"+pfdTrans,
					bytes.NewBuffer(body))
			} else {
				//GET PFD/ DELETE

				req, _ = http.NewRequest(method, baseInvalidPFDAPIURL+"/"+
					pfdTrans,
					nil)
			}

		}

	} else {
		if body != nil {
			//POST
			req, _ = http.NewRequest(method, baseInvalidPFDAPIURL,
				bytes.NewBuffer(body))
		} else {
			//GET ALL
			req, _ = http.NewRequest(method, baseInvalidPFDAPIURL, nil)
		}
	}

	rr := httptest.NewRecorder()
	return rr, req
}

var _ = Describe("Test NEF Server PFD NB API's ", func() {
	var ctx context.Context
	var cancel func()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	Describe("Start the NEF Server: To be done to start NEF PFD API testing",
		func() {
			It("Will init NefServer",
				func() {
					ctx, cancel = context.WithCancel(context.Background())
					defer cancel()
					go func() {
						err := ngcnef.Run(ctx, NefTestCfgBasepath+"valid.json")
						Expect(err).To(BeNil())
					}()
					time.Sleep(2 * time.Second)
				})
		})

	Describe("VALID REQ to NEF GET/POST/PUT/PATCH/DELETE", func() {

		postbody, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_POST_001.json")
		postbody201, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_POST_003.json")
		postbodyUDR, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_POST_004.json")
		postbody400, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_POST_005.json")
		postbodyMissingAppID, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_POST_006.json")
		postbodyUnmarshall, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_POST_007.json")
		putbody, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_PUT_001.json")
		putbody400, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_PUT_002.json")
		putbodyPfdDataMissing, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_PUT_003.json")
		putbodyPfdMissing, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_PUT_004.json")
		putappbody, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_APP_PUT_001.json")
		putappbodyPfdMissing, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_APP_PUT_003.json")
		patchappbody, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_APP_PATCH_001.json")
		patchappbody400, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_APP_PATCH_002.json")
		patchappbodyMissingPfd, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_APP_PATCH_003.json")
		patchappbodyUnmarshall, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_APP_PATCH_004.json")
		putappbody400, _ := ioutil.ReadFile(testJSONPFDPath +
			"AF_NEF_PFD_APP_PUT_002.json")

		It("Send valid GET all to NEF -No Data as no PFD exists",
			func() {
				rr, req := CreatePFDReqForNEF(ctx, "GET", "", "", nil)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusOK))

				//Validate PFD
				//Read Body from response
				resp := rr.Result()
				b, _ := ioutil.ReadAll(resp.Body)

				//Convert the body(json data) into PFD Management Struct data
				var pfdBody []ngcnef.PfdManagement
				err := json.Unmarshal(b, &pfdBody)
				Expect(err).Should(BeNil())
				fmt.Print("Body Received: ")
				fmt.Println(pfdBody)
				resp.Body.Close()
				Expect(len(pfdBody)).Should(Equal(0))
			})

		It("Send valid POST to NEF", func() {
			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "", postbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusCreated))
			//Validate Body of Trans
			//Read Body from response
			resp := rr.Result()
			b, _ := ioutil.ReadAll(resp.Body)

			//Convert the body(json data) into PFD Management Struct data
			var pfdBody ngcnef.PfdManagement
			err := json.Unmarshal(b, &pfdBody)
			Expect(err).Should(BeNil())

			fmt.Print("Self in PFD manageemnt Received: ")
			fmt.Println(pfdBody.Self)
			Expect(pfdBody.Self).ShouldNot(Equal(""))
			for _, v := range pfdBody.PfdDatas {
				fmt.Println(v.Self)
				Expect(v.Self).ShouldNot(Equal(""))
			}
			resp.Body.Close()

		})

		It("Send invalid POST to NEF (500 error response)", func() {
			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "", postbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusInternalServerError))
			//Validate Body in response ?

			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send invalid POST to NEF appID missing)", func() {
			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "",
				postbodyMissingAppID)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			//Validate Body in response ?

			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send invalid POST to NEF (400 error response from UDR)", func() {
			ngcnef.TestClient = true
			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "", postbodyUDR)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			//Validate Body in response ?
			ngcnef.TestClient = false
			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send invalid POST to NEF (Invalid PFDs)", func() {

			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "", postbody400)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			//Validate Body in response ?

			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send invalid POST to NEF (Unmarshall error)", func() {

			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "",
				postbodyUnmarshall)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			//Validate Body in response ?

			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send invalid PUT to NEF (Unmarshall error)", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "",
				postbodyUnmarshall)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			//Validate Body in response ?

			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send invalid POST to NEF (400 error response from UDR)", func() {
			ngcnef.TestClient = true
			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "", postbodyUDR)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			//Validate Body in response ?
			ngcnef.TestClient = false
			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send invalid POST to NEF (500 error response from UDR)", func() {
			ngcnef.TestNEFSB = true
			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "", postbodyUDR)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusInternalServerError))
			//Validate Body in response ?
			ngcnef.TestNEFSB = false
			resp := rr.Result()
			resp.Body.Close()

		})

		It("Send POST to NEF with one application invalid", func() {
			rr, req := CreatePFDReqForNEF(ctx, "POST", "", "", postbody201)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusCreated))
			//Validate Body in response ?

			resp := rr.Result()
			b, _ := ioutil.ReadAll(resp.Body)

			//Convert the body(json data) into PFD Management Struct data
			var pfdBody ngcnef.PfdManagement

			err := json.Unmarshal(b, &pfdBody)
			Expect(err).Should(BeNil())

			fmt.Print("Self in PFD manageemnt Received: ")
			fmt.Println(pfdBody.Self)
			Expect(pfdBody.Self).ShouldNot(Equal(""))

			for _, v := range pfdBody.PfdReports {
				fmt.Printf("Failure Code is %s", v.FailureCode)
				Expect(string(v.FailureCode)).Should(Equal("APP_ID_DUPLICATED"))
			}

			resp.Body.Close()

		})

		It("Will Send a valid PFD GET ALL", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))

		})

		It("Will Send a invalid PFD GET ALL UDR CLIENT", func() {

			ngcnef.TestClient = true
			rr, req := CreatePFDReqForNEF(ctx, "GET", "", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
			ngcnef.TestClient = false

		})

		It("Will Send a valid PFD GET for PFD TRANS 10000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "10000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send an invalid PFD GET - INVALID AF", func() {

			rr, req := CreateInvalidPFDReqForNEF(ctx, "GET", "10000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))

		})

		It("Will Send n invalid PFD GET for PFD TRANS 11000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "11000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})

		It("Will Send n invalid PFD GET for PFD TRANS 10000"+
			"400 from UDR", func() {

			ngcnef.TestClient = true
			rr, req := CreatePFDReqForNEF(ctx, "GET", "10000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
			ngcnef.TestClient = false
		})

		It("Will Send a valid PUT for PFD TRANS 10000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "", putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send a invalid PUT for PFD TRANS (INVALID AF)", func() {

			rr, req := CreateInvalidPFDReqForNEF(ctx, "PUT", "10000",
				"", putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})

		It("Will Send a invalid PUT for PFD TRANS 10000(400 from UDR)", func() {

			ngcnef.TestClient = true
			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "", putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			ngcnef.TestClient = false
		})

		It("Will Send a invalid PUT for PFD TRANS 10000(500 from UDR)", func() {

			ngcnef.TestNEFSB = true
			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "", putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusInternalServerError))
			ngcnef.TestNEFSB = false
		})

		It("Will Send an invalid PUT for PFD TRANS 10000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "", putbody400)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
		})

		It("Will Send an invalid PUT for PFD TRANS 10000 (PFDData missing)",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "",
					putbodyPfdDataMissing)
				req.Header.Set("Content-Type", "application/json")
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			})

		It("Will Send an invalid PUT for PFD TRANS 10000 (PFDs missing)",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "",
					putbodyPfdMissing)
				req.Header.Set("Content-Type", "application/json")
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			})
		It("Will Send a valid PFD GET for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "10000", "app1", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send an invalid PFD PUT for PFD TRANS/APP - Unmarshall error",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app1",
					postbodyUnmarshall)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			})

		It("Will Send invalid PFD PATCH for PFD TRANS/APP - Unmarshall error",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app1",
					postbodyUnmarshall)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			})

		It("Will Send a invalid PFD GET for PFDTRANS and app1 - INVALID AF",
			func() {

				rr, req := CreateInvalidPFDReqForNEF(ctx, "GET", "10000",
					"app1", nil)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusNotFound))
			})

		It("Will Send an invalid PFD GET for PFD TRANS 10000 and app10",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "GET", "10000", "app10", nil)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusNotFound))
			})
		It("Will Send a invalid PFD GET for PFD TRANS 11000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "11000", "app1", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})

		It("Will Send a valid PFD PUT for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app1",
				putappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send a invalid PFD PUT for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app1",
				putappbody400)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
		})

		It("Will Send a invalid PFD PUT for PFD TRANS 10000 and app1"+
			"PFD missing", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app1",
				putappbodyPfdMissing)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
		})

		It("Will Send a invalid PFD PUT for PFD TRANS 10000 and app1"+
			"(400 from UDR)", func() {

			ngcnef.TestClient = true
			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app1",
				putappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			ngcnef.TestClient = false
		})

		It("Send invalid APP PUT to NEF (500 error response from UDR)", func() {
			ngcnef.TestNEFSB = true
			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app1",
				putappbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusInternalServerError))

			ngcnef.TestNEFSB = false
			resp := rr.Result()
			resp.Body.Close()

		})

		It("Will Send a invalid PFD PUT for PFD TRANS 11000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "11000", "app1",
				putappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})

		It("Will Send a invalid PFD PUT for PFD TRANS 10000 and app10", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app10",
				putappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})

		It("Will Send a valid PFD PATCH for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app1",
				patchappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send an invalid PFD PATCH for TRANS/APP - INVALID AF", func() {

			rr, req := CreateInvalidPFDReqForNEF(ctx, "PATCH", "10000", "app1",
				patchappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})
		It("Will Send a invalid PFD PATCH for PFD TRANS 10000 and app1",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app1",
					patchappbody400)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			})
		It("Will Send a invalid PFD PATCH for PFD TRANS 10000 and app10",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app10",
					patchappbody)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusNotFound))
			})
		It("Will Send a invalid PFD PATCH for PFD TRANS 11000 and app1",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "11000", "app1",
					patchappbody)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusNotFound))
			})
		It("Will Send a invalid PFD PATCH for PFD TRANS 10000 and app1",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "11000", "app1",
					patchappbodyMissingPfd)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			})
		It("Will Send a invalid PFD PATCH - Unmarsahll error",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app1",
					patchappbodyUnmarshall)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			})

		It("Will Send a invalid PFD PATCH for PFD TRANS 10000 and app1"+
			"(400 from UDR)",
			func() {

				ngcnef.TestClient = true
				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app1",
					patchappbody)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusBadRequest))
				ngcnef.TestClient = false
			})

		It("Send invalid APP PATCH to NEF (500 error response from UDR)",
			func() {
				ngcnef.TestNEFSB = true
				rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app1",
					patchappbody)
				req.Header.Set("Content-Type", "application/json")
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

				Expect(rr.Code).Should(Equal(http.StatusInternalServerError))

				ngcnef.TestNEFSB = false
				resp := rr.Result()
				resp.Body.Close()

			})

		It("Will Send a valid PFD DELETE for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "DELETE", "10000", "app1", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})

		It("Will Send an invalid PFD DELETE for TRANS/APP - INVALID AF",
			func() {

				rr, req := CreateInvalidPFDReqForNEF(ctx, "DELETE", "10000",
					"app1", nil)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusNotFound))
			})

		It("Will Send an invalid PFD DELETE for PFD TRANS 10000 and app10",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "DELETE", "10000", "app10",
					nil)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusNotFound))
			})

		It("Will Send an invalid PFD DELETE for PFD TRANS 11000 and app1",
			func() {

				rr, req := CreatePFDReqForNEF(ctx, "DELETE", "11000", "app1",
					nil)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusNotFound))
			})
		It("Will Send a invalid DELETE for PFD TRANS (400 from UDR)", func() {
			ngcnef.TestClient = true
			rr, req := CreatePFDReqForNEF(ctx, "DELETE", "10000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusBadRequest))
			ngcnef.TestClient = false
		})

		It("Will Send a valid DELETE for PFD TRANS 10000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "DELETE", "10000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})

		It("Will Send an invalid DELETE for PFD TRANS- INVALID AF", func() {

			rr, req := CreateInvalidPFDReqForNEF(ctx, "DELETE", "10000", "",
				nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})

		It("Will Send a invalid DELETE for PFD TRANS 11000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "DELETE", "11000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})

		It("Will Send a valid DELETE for PFD TRANS 10004", func() {

			rr, req := CreatePFDReqForNEF(ctx, "DELETE", "10004", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})

	})

	Describe("End the NEF Server: To be done to end NEF PFD API testing",
		func() {
			It("Will stop NefServer", func() {
				cancel()
				time.Sleep(2 * time.Second)
			})
		})

})
