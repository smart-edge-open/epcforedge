/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
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

		postbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PFD_POST_01.json")
		putbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PFD_PUT_01.json")
		putappbody, _ := ioutil.ReadFile(testJSONPath +
			"AF_NEF_PFD_APP_PUT_01.json")
		patchappbody, _ := ioutil.ReadFile(testJSONPath +
			"AF_NEF_PFD_APP_PATCH_01.json")

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

		It("Will Send a valid PFD GET ALL", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))

		})

		It("Will Send a valid PFD GET for PFD TRANS 10000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "10000", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send a valid PUT for PFD TRANS 10000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "", putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send a valid PFD GET for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "GET", "10000", "app1", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send a valid PFD PUT for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PUT", "10000", "app1",
				putappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send a valid PFD PATCH for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "PATCH", "10000", "app1",
				patchappbody)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})

		It("Will Send a valid PFD DELETE for PFD TRANS 10000 and app1", func() {

			rr, req := CreatePFDReqForNEF(ctx, "DELETE", "10000", "app1", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})
		It("Will Send a valid DELETE for PFD TRANS 10000", func() {

			rr, req := CreatePFDReqForNEF(ctx, "DELETE", "10000", "", nil)
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
