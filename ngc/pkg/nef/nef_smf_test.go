package ngcnef_test

import (
	"bytes"
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("NefSmf", func() {
	var (
		ctx    context.Context
		cancel func()
	)

	Describe("NefServer SMF Functionality", func() {
		It("Starting the NEF server", func() {
			ctx, cancel = context.WithCancel(context.Background())
			defer cancel()
			go func() {
				fmt.Println("** Starting the NEF server ***")
				err := ngcnef.Run(ctx, NefTestCfgBasepath+"valid.json")
				Expect(err).To(BeNil())
			}()
			time.Sleep(2 * time.Second)
		})

		It("POST an UPF notification",
			func() {
				postbody, _ := ioutil.ReadFile(NefTestJSONBasepath +
					"SMF_NEF_NOTIF_01.json")
				req, _ := http.NewRequest("POST", NefTIFApiPrefixHTTP2+
					"notification/upf", bytes.NewBuffer(postbody))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				ngcnef.NefAppG.NefRouter.ServeHTTP(rr,
					req.WithContext(ctx))
				Expect(rr.Code == http.StatusNotFound).To(BeTrue())
			})

		It("Stopping the NEF server", func() {
			cancel()
			time.Sleep(2 * time.Second)
			fmt.Print("** Stopping the NEF server ** ")
		})

	})

})
