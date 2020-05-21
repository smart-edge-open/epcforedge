// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package main

import (
	"cnca"
	"os"
)

func main() {

	if err := cnca.Execute(); err != nil {
		os.Exit(1)
	}
}
