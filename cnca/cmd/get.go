// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"errors"
	"fmt"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get active LTE CUPS userplane(s) or NGC AF TI subscription(s)",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

		if args[0] == "subscription" {

			if len(args) < 2 {
				fmt.Println(errors.New("Missing input"))
				return
			}

			// get subscription
			sub, err := AFGetSubscription(args[1])
			if err != nil {
				klog.Info(err)
				return
			}

			sub, err = y2j.JSONToYAML(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active AF Subscription:\n%s", string(sub))
			return
		} else if args[0] == "subscriptions" {

			// get subscriptions
			sub, err := AFGetSubscription("all")
			if err != nil {
				klog.Info(err)
				return
			}

			if string(sub) == "[]" {
				sub = []byte("none")
			}

			sub, err = y2j.JSONToYAML(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active AF Subscriptions:\n%s", string(sub))
			return
		} else if args[0] == "userplane" {

			if len(args) < 2 {
				fmt.Println(errors.New("Missing input"))
				return
			}

			// get userplane
			up, err := LteGetUserplane(args[1])
			if err != nil {
				klog.Info(err)
				return
			}

			up, err = y2j.JSONToYAML(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active LTE CUPS Userplane:\n%s", string(up))
			return
		} else if args[0] == "userplanes" {

			// get userplanes
			up, err := LteGetUserplane("all")
			if err != nil {
				klog.Info(err)
				return
			}

			if string(up) == "[]" {
				up = []byte("none")
			}

			up, err = y2j.JSONToYAML(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active LTE CUPS Userplanes:\n%s", string(up))
			return
		}

		fmt.Println(errors.New("Invalid input(s)"))
	},
}

// pfdGetCmd represents the get command
var pfdGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get active NGC AF PFD Transaction(s) or NGC AF PFD Application(s)",
	Args:  cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

		if args[0] == "transactions" || args[0] == "transaction" {
			var transID string
			var appID string
			var pfdData []byte
			var err error

			if args[0] == "transaction" && len(args) > 1 {
				transID = args[1]
				if len(args) > 2 && args[2] != "" {
					if args[2] == "application" && len(args) > 3 {
						appID = args[3]
					} else {
						fmt.Println(errors.New("Invalid input(s)"))
						return
					}
				}
			} else if args[0] == "transactions" {
				transID = "all"
			} else {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			if appID != "" {
				// get PFD application
				pfdData, err = AFGetPfdApplication(transID, appID)
				if err != nil {
					klog.Info(err)
					return
				}
			} else {
				// get PFD transaction
				pfdData, err = AFGetPfdTransaction(transID)
				if err != nil {
					klog.Info(err)
					return
				}
			}

			if args[0] == "transactions" && string(pfdData) == "[]" {
				pfdData = []byte("none")
			}

			pfdData, err = y2j.JSONToYAML(pfdData)
			if err != nil {
				fmt.Println(err)
				return
			}

			if appID != "" {
				fmt.Printf("PFD Application: %s\n%s", appID, string(pfdData))
			} else {
				fmt.Printf("PFD Transaction: %s\n%s", transID, string(pfdData))
			}
			return
		}

		fmt.Println(errors.New("Invalid input(s)"))
	},
}

//paGetCmd represents get command
var paGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an active NGC AF Application Session",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

		var appSessionID string
		var appSession []byte
		var err error

		appSessionID = args[0]

		if appSessionID != "" {
			appSession, err = AFGetPaAppSession(appSessionID)
			if err != nil {
				klog.Info(err)
				return
			}

			if string(appSession) == "[]" {
				appSession = []byte("none")
			}

			appSession, err = y2j.JSONToYAML(appSession)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("App session context Data: %s\n%s", appSessionID, string(appSession))
			return
		}
		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help = `Get active LTE CUPS userplane(s) or NGC AF TI subscription(s)

Usage:
  cnca get { userplanes | 
             subscriptions | 
             userplane <userplane-id> | 
             subscription <subscription-id>}


Example:
  cnca get userplane <userplane-id>
  cnca get subscription <subscription-id>
  cnca get userplanes
  cnca get subscriptions

Flags:
  -h, --help   help
`

	const pfdHelp = `Get active NGC AF PFD Transaction(s) or NGC AF PFD 
Application(s)

Usage:
  cnca pfd get { transactions | 
                 transaction <transaction-id> |
                 transaction <transaction-id> application <application-id>}

Example:
  cnca pfd get transactions
  cnca pfd get transaction <transaction-id>
  cnca pfd get transaction <transaction-id> application <application-id>

Flags:
  -h, --help   help
`

	const paHelp = `Get an active NGC AF Application Session

Usage:
  cnca policy-authorization get <appSession-id>

Example:
  cnca policy-authorization get <appSession-id>

Flags:
  -h, --help   help
`

	// add `get` command
	cncaCmd.AddCommand(getCmd)
	getCmd.SetHelpTemplate(help)

	// add pfd `get` command
	pfdCmd.AddCommand(pfdGetCmd)
	pfdGetCmd.SetHelpTemplate(pfdHelp)

	//add pa `get` command
	paCmd.AddCommand(paGetCmd)
	paGetCmd.SetHelpTemplate(paHelp)
}
