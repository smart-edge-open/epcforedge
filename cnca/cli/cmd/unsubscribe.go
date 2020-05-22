// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// paUnsubscribeCmd represents the delete command
var paUnsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe",
	Short: "Remove subscription to all subscribed events for app session",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input(s)"))
			return
		}

		if args[0] != "" {
			err := AFPaEventUnsubscribe(args[0])
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription to events of app session %s removed\n", args[0])
			return
		}

		fmt.Println(errors.New("Invalid input(s)"))

	},
}

func init() {
	const help = `Remove subscription to all subscribed events for app session

Usage:
  cnca policy-authorization unsubscribe <appSessionId>

Example:
  cnca policy-authorization unsubscribe <appSessionId>

Flags:
  -h, --help   help
`
	//add pa `delete` command
	paCmd.AddCommand(paUnsubscribeCmd)
	paUnsubscribeCmd.SetHelpTemplate(help)

}
