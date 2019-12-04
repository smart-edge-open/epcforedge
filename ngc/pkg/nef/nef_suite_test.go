/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const NefTestCfgBasepath = "../../test/nef/configs/"
const NefTestJSONBasepath = "../../test/nef/nef-cli-scripts/json/"
const NefTIFApiPrefix = "http://localhost:8091/3gpp-traffic-influence/v1/"
const NefTIFApiPrefixHTTP2 = "https://localhost:8090/3gpp-traffic-influence/v1/"

func TestNef(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nef Suite")
}
