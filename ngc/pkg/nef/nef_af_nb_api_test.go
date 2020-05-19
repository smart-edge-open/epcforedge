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

//const validCfgPath = "../../configs/nef.json"
const testJSONPath = "../../test/nef/nef-cli-scripts/json/"
const baseAPIURL = "http://localhost:8091/3gpp-traffic-influence/" +
	"v1/AF_01/subscriptions"

func CreateReqForNEF(ctx context.Context, method string, subID string,
	body []byte) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request
	if len(subID) > 0 {
		if body != nil {
			//PUT, PATCH
			req, _ = http.NewRequest(method, baseAPIURL+"/"+subID,
				bytes.NewBuffer(body))
		} else {
			//GET SUB
			req, _ = http.NewRequest(method, baseAPIURL+"/"+subID, nil)
		}
	} else {
		if body != nil {
			//POST
			req, _ = http.NewRequest(method, baseAPIURL, bytes.NewBuffer(body))
		} else {
			//GET ALL
			req, _ = http.NewRequest(method, baseAPIURL, nil)
		}
	}
	/*
		ctx = context.WithValue(
			req.Context(),
			"nefCtx",
			ngcnef.NefAppG.NefCtx)
	*/
	rr := httptest.NewRecorder()
	return rr, req
}

var _ = Describe("Test NEF Server NB API's ", func() {
	var ctx context.Context
	var cancel func()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	Describe("Start the NEF Server: To be done to start NEF API testing",
		func() {
			It("Will init NefServer",
				func() {
					ctx, cancel = context.WithCancel(context.Background())
					defer cancel()
					go func() {
						generateCerts()
						err := ngcnef.Run(ctx, NefTestCfgBasepath+"valid.json")
						Expect(err).To(BeNil())
						removeCerts()
					}()
					time.Sleep(2 * time.Second)
				})
		})

	Describe("REQ towards PCF (POST/PUT/PATCH/DELETE)", func() {

		postbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_POST_01.json")
		putbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PUT_01.json")
		patchbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PATCH_01.json")

		It("Send valid GET all to NEF -No TI Data in response as no Sub exists",
			func() {
				rr, req := CreateReqForNEF(ctx, "GET", "", nil)
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				Expect(rr.Code).Should(Equal(http.StatusOK))

				//Validate TI
				//Read Body from response
				resp := rr.Result()
				b, _ := ioutil.ReadAll(resp.Body)

				//Convert the body(json data) into Traffic Influence Struct data
				var trInBody []ngcnef.TrafficInfluSub
				err := json.Unmarshal(b, &trInBody)
				Expect(err).Should(BeNil())
				fmt.Print("Body Received: ")
				fmt.Println(trInBody)
				resp.Body.Close()
				Expect(len(trInBody)).Should(Equal(0))
			})
		It("Send valid POST to NEF towards PCF ", func() {
			rr, req := CreateReqForNEF(ctx, "POST", "", postbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusCreated))
			//Validate Body of TI
			//Read Body from response
			resp := rr.Result()
			b, _ := ioutil.ReadAll(resp.Body)

			//Convert the body(json data) into Traffic Influence Struct data
			var trInBody ngcnef.TrafficInfluSub
			err := json.Unmarshal(b, &trInBody)
			Expect(err).Should(BeNil())

			fmt.Print("Self in TI Received: ")
			fmt.Println(trInBody.Self)
			resp.Body.Close()
			Expect(trInBody.Self).ShouldNot(Equal(""))
		})
		It("Will Send a valid GET all towards PCF", func() {

			rr, req := CreateReqForNEF(ctx, "GET", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
			//Validate TI
		})
		It("Will Send a valid GET towards PCF", func() {

			rr, req := CreateReqForNEF(ctx, "GET", "11111", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid PUT towards PCF", func() {

			rr, req := CreateReqForNEF(ctx, "PUT", "11111", putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})
		It("Will Send a valid PATCH towards PCF", func() {

			rr, req := CreateReqForNEF(ctx, "PATCH", "11111", patchbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid DELETE towards PCF", func() {

			rr, req := CreateReqForNEF(ctx, "DELETE", "11111", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})
	})

	Describe("REQ towards UDR(POST/PUT/PATCH/DELETE)", func() {
		postbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_POST_UDR_01.json")
		putbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PUT_UDR_01.json")
		patchbody, _ := ioutil.ReadFile(testJSONPath +
			"AF_NEF_PATCH_UDR_01.json")

		It("Send valid POST to NEF towards UDR ", func() {
			rr, req := CreateReqForNEF(ctx, "POST", "", postbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusCreated))
		})
		It("Will Send a valid GET all towards UDR", func() {

			rr, req := CreateReqForNEF(ctx, "GET", "", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid GET towards UDR", func() {

			rr, req := CreateReqForNEF(ctx, "GET", "11111", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid PUT towards UDR", func() {

			rr, req := CreateReqForNEF(ctx, "PUT", "11111", putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))

			//Validate Body of TI
			//Read Body from response
			resp := rr.Result()
			b, _ := ioutil.ReadAll(resp.Body)

			//Convert the body(json data) into Traffic Influence Struct data
			var trInBody ngcnef.TrafficInfluSub
			err := json.Unmarshal(b, &trInBody)
			Expect(err).Should(BeNil())

			fmt.Print("Self in TI Received: ")
			fmt.Println(trInBody.Self)
			resp.Body.Close()
			Expect(trInBody.Self).ShouldNot(Equal(""))
		})
		It("Will Send a valid PATCH towards UDR", func() {

			rr, req := CreateReqForNEF(ctx, "PATCH", "11111", patchbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid DELETE towards UDR", func() {

			rr, req := CreateReqForNEF(ctx, "DELETE", "11111", nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})
	})

	Describe("End the NEF Server: To be done to end NEF API testing",
		func() {
			It("Will stop NefServer", func() {
				cancel()
				time.Sleep(2 * time.Second)
			})
		})

})
