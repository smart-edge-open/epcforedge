// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// unregisterCmd represents the unregister command
var unregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Un-register controller from AF services registry",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("AF service ID missing"))
			return
		}
		// unregister AF service
		err := OAM5gUnregisterAFService(args[0])
		if err != nil {
			klog.Info(err)
			return
		}
		fmt.Printf("Service ID `%s` unregistered successfully\n", args[0])
	},
}

func init() {

	const help = `Unregister controller from NGC AF services registry

Usage:
  cnca unregister <service-id>

Flags:
  -h, --help       help
`
	// add `register` command
	cncaCmd.AddCommand(unregisterCmd)
	unregisterCmd.SetHelpTemplate(help)
}
