package ngcnef_test

import (
	"bytes"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

const validCfgPath = "../../configs/nef.json"

var _ = Describe("NefServer", func() {
	Describe("NefServer init", func() {
		It("Will init NefServer - Invalid Configurations",
			func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				err := ngcnef.Run(ctx, "noconfig")
				Expect(err).NotTo(BeNil())
			})
		It("Will init NefServer - Valid Configurations",
			func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				go func() {
					/* Send a cancel after 5 seconds */
					time.Sleep(3 * time.Second)
					cancel()
				}()
				err := ngcnef.Run(ctx, validCfgPath)
				Expect(err).To(BeNil())

			})
	})

	Describe("NefServer NB POST", func() {
		It("Will Send a valid POST",
			func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				go func() {
					_ = ngcnef.Run(ctx, validCfgPath)
				}()
				time.Sleep(3 * time.Second)

				postbody, _ := ioutil.ReadFile("../../test/nef/nef-cli-scripts/json/AF_NEF_POST_01.json")
				req, _ := http.NewRequest("POST",
					"http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions",
					bytes.NewBuffer(postbody))
				req.Header.Set("Content-Type", "application/json")
				ctx = context.WithValue(
					req.Context(),
					"nefCtx",
					ngcnef.NefAppG.NefCtx)
				rr := httptest.NewRecorder()
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr, req.WithContext(ctx))
				time.Sleep(3 * time.Second)
				cancel()

			})
	})

})
