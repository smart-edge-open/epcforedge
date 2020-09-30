// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package main

import (
	"os"

	cnca "github.com/open-ness/epcforedge/cnca/cli/cmd"
)

func main() {

	if err := cnca.Execute(); err != nil {
		os.Exit(1)
	}
}
