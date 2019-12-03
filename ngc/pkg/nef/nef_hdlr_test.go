package ngcnef_test

import (
	"bytes"
	"context"
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
const baseAPIURL = "http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions"

func CreateReqForNEF(method string, subID string,
	ctx context.Context, body []byte) (*httptest.ResponseRecorder,
	*http.Request) {
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

	ctx = context.WithValue(
		req.Context(),
		"nefCtx",
		ngcnef.NefAppG.NefCtx)
	rr := httptest.NewRecorder()
	return rr, req
}

var _ = Describe("Test NEF Server NB API's ", func() {
	var ctx context.Context
	//var cancel context.CancelFunc

	ctx, _ = context.WithCancel(context.Background())
	//defer cancel()

	Describe("Start the NEF Server: To be done to start NEF API testing",
		func() {
			It("Will init NefServer",
				func() {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					go func() {
						/* Send a cancel after 5 seconds */
						time.Sleep(1 * time.Second)
						cancel()
					}()
					err := ngcnef.Run(ctx, validCfgPath)
					Expect(err).To(BeNil())
				})
		})

	Describe("REQ towards PCF (POST/PUT/PATCH/DELETE)", func() {

		postbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_POST_01.json")
		putbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PUT_01.json")
		patchbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PATCH_01.json")

		It("Send valid POST to NEF towards PCF ", func() {
			rr, req := CreateReqForNEF("POST", "", ctx, postbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))

			Expect(rr.Code).Should(Equal(http.StatusCreated))
			//Validate Body of TI
			//Expect(bytes.Equal(rr.Body.Bytes(), postbody)).Should(BeTrue())
		})
		It("Will Send a valid GET all towards PCF", func() {

			rr, req := CreateReqForNEF("GET", "", ctx, nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
			//Validate TI
		})
		It("Will Send a valid GET towards PCF", func() {

			rr, req := CreateReqForNEF("GET", "11111", ctx, nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid PUT towards PCF", func() {

			rr, req := CreateReqForNEF("PUT", "11111", ctx, putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNotFound))
		})
		It("Will Send a valid PATCH towards PCF", func() {

			rr, req := CreateReqForNEF("PATCH", "11111", ctx, patchbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid DELETE towards PCF", func() {

			rr, req := CreateReqForNEF("DELETE", "11111", ctx, nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})
	})

	Describe("REQ towards UDR(POST/PUT/PATCH/DELETE)", func() {
		postbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_POST_UDR_01.json")
		putbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PUT_UDR_01.json")
		patchbody, _ := ioutil.ReadFile(testJSONPath + "AF_NEF_PATCH_UDR_01.json")

		It("Send valid POST to NEF towards UDR ", func() {
			rr, req := CreateReqForNEF("POST", "", ctx, postbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusCreated))
		})
		It("Will Send a valid GET all towards UDR", func() {

			rr, req := CreateReqForNEF("GET", "", ctx, nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid GET towards UDR", func() {

			rr, req := CreateReqForNEF("GET", "11111", ctx, nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid PUT towards UDR", func() {

			rr, req := CreateReqForNEF("PUT", "11111", ctx, putbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid PATCH towards UDR", func() {

			rr, req := CreateReqForNEF("PATCH", "11111", ctx, patchbody)
			req.Header.Set("Content-Type", "application/json")
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusOK))
		})
		It("Will Send a valid DELETE towards UDR", func() {

			rr, req := CreateReqForNEF("DELETE", "11111", ctx, nil)
			ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
			Expect(rr.Code).Should(Equal(http.StatusNoContent))
		})
	})

})
