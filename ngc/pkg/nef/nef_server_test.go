package ngcnef_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
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
					time.Sleep(1 * time.Second)
					cancel()
				}()
				err := ngcnef.Run(ctx, validCfgPath)
				Expect(err).To(BeNil())

			})
	})

})
