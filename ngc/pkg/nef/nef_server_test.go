package ngcnef_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
	"os"
	"os/signal"
	"syscall"
)

const validCfgPath = "../../configs/nef.json"

var _ = Describe("NefServer", func() {
	Describe("NefServer init", func() {
		It("Will init NefServer",
			func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				err := ngcnef.Run(ctx, "noconfig")
				Expect(err).NotTo(BeNil())

				go func() {
					sig := <-osSignalCh
					cancel()
				}()
				err = ngcnef.Run(ctx, "validCfgPath")
				Expect(err).To(BeNil())

			})
	})

})
