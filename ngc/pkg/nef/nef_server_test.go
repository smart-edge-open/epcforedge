/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ngcnef "github.com/open-ness/epcforedge/ngc/pkg/nef"
)

var _ = Describe("NefServer", func() {

	var (
		ctx    context.Context
		cancel func()
		//testErr error
	)

	Describe("NefServer init", func() {
		It("Will init NefServer - Invalid Configurations",
			func() {

				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
				err := ngcnef.Run(ctx, "noconfig")
				Expect(err).NotTo(BeNil())
			})

		It("Will init NefServer - No HTTP endpoints",
			func() {

				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
				err := ngcnef.Run(ctx,
					NefTestCfgBasepath+"invalid_no_eps.json")
				Expect(err).NotTo(BeNil())
			})
		/*
			// Commenting since its covered through other tests
			It("Will init NefServer - Valid Configurations",
				func() {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					go func() {
						testErr = ngcnef.Run(ctx,
							NefTestCfgBasepath+"valid.json")
					}()
					// Send a cancel after 3 seconds
					time.Sleep(3 * time.Second)
					cancel()
					Expect(testErr).To(BeNil())
				})
		*/
	})

})
