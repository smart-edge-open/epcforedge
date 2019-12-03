package ngcnef_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
	"time"
)

var _ = Describe("NefServer", func() {

	var (
		ctx     context.Context
		cancel  func()
		testErr error
	)

	Describe("NefServer init", func() {
		It("Will init NefServer - Invalid Configurations",
			func() {

				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
				err := ngcnef.Run(ctx, "noconfig")
				Expect(err).NotTo(BeNil())
			})
		It("Will init NefServer - Valid Configurations",
			func() {
				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
				go func() {
					testErr = ngcnef.Run(ctx,
						NefTestCfgBasepath+"valid.json")
				}()
				/* Send a cancel after 5 seconds */
				time.Sleep(3 * time.Second)
				cancel()
				Expect(testErr).To(BeNil())
			})
	})

})
