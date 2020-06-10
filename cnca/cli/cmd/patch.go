// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"encoding/json"
	"errors"
	"fmt"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use: "patch",
	Short: "Patch an active LTE CUPS userplane or NGC AF TI subscription " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

		// Read file from the filename provided in command
		data, err := readInputData(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		switch c.Kind {
		case "ngc":

			var s AFTrafficInfluSub
			if err = yaml.Unmarshal(data, &s); err != nil {
				fmt.Println(err)
				return
			}

			var sub []byte
			sub, err = yaml.Marshal(s.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			sub, err = y2j.YAMLToJSON(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			// patch subscription
			err = AFPatchSubscription(args[0], sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])

		case "lte":

			var u LTEUserplane
			if err = yaml.Unmarshal(data, &u); err != nil {
				fmt.Println(err)
				return
			}

			up, err := yaml.Marshal(u.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			up, err = y2j.YAMLToJSON(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			// patch userplane
			err = LtePatchUserplane(args[0], up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])
		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

// pfdPatchCmd represents the patch command
var pfdPatchCmd = &cobra.Command{
	Use: "patch",
	Short: "Patch an active NGC AF PFD Transaction or NGC AF PFD Application " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println(errors.New("Missing input(s)"))
			return
		}

		// Read file from the filename provided in command
		data, err := readInputData(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		if c.Kind != "ngc_pfd" {
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
			return
		}

		if args[0] == "transaction" && args[1] != "" {
			var pfdReportData []byte
			if len(args) > 2 {
				if args[2] == "application" && len(args) > 3 {
					// patch PFD Application

					var s AFPfdData
					if err = yaml.Unmarshal(data, &s); err != nil {
						fmt.Println(err)
						return
					}

					pfdAppData := getPfdAppData(s)

					var app []byte
					app, err = json.Marshal(pfdAppData)
					if err != nil {
						fmt.Println(err)
						return
					}

					pfdReportData, err = AFPatchPfdApplication(args[1], args[3], app)
					if err != nil {
						klog.Info(err)
						if err.Error() == "HTTP failure: 500" && pfdReportData != nil {
							printPfdReport(pfdReportData)
						}
						return
					}
					fmt.Printf("PFD Application %s patched\n", args[3])
					return
				}
			} else {
				// patch PFD Transaction

				var s AFPfdManagement
				if err = yaml.Unmarshal(data, &s); err != nil {
					fmt.Println(err)
					return
				}

				pfdTransData := getPfdTransData(s)

				var trans []byte
				trans, err = json.Marshal(pfdTransData)
				if err != nil {
					fmt.Println(err)
					return
				}
				pfdReportData, err = AFPatchPfdTransaction(args[1], trans)
				if err != nil {
					klog.Info(err)
					if err.Error() == "HTTP failure: 500" && pfdReportData != nil {
						printPfdReports(pfdReportData)
					}
					return
				}
				fmt.Printf("PFD Transaction %s patched\n", args[1])
				return
			}
		}
		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help = `Patch an active LTE CUPS userplane or NGC AF TI subscription 
using YAML configuration file

Usage:
  cnca patch { <userplane-id> | <subscription-id> } -f <config.yml>

Example:
  cnca patch <userplane-id> -f <config.yml>
  cnca patch <subscription-id> -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	const pfdHelp = `Patch an active NGC AF PFD Transaction or NGC AF PFD 
Application using YAML Configuration file

Usage:
  cnca pfd patch { transaction <transaction-id> |
                   transaction <transaction-id> application <application-id>}
                   -f <config.yml>

Example:
  cnca pfd patch transaction <transaction-id> -f <config.yml>
  cnca pfd patch transaction <transaction-id> 
    application <application-id> -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	// add `patch` command
	cncaCmd.AddCommand(patchCmd)
	patchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = patchCmd.MarkFlagRequired("filename")
	patchCmd.SetHelpTemplate(help)

	// add pfd `patch` command
	pfdCmd.AddCommand(pfdPatchCmd)
	pfdPatchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = pfdPatchCmd.MarkFlagRequired("filename")
	pfdPatchCmd.SetHelpTemplate(pfdHelp)
}

func getPfdAppData(inputPfdAppData AFPfdData) PfdData {

	var pfdAppData PfdData

	pfdAppData.ExternalAppID = inputPfdAppData.Policy.ExternalAppID
	pfdAppData.Self = Link(inputPfdAppData.Policy.Self)

	if inputPfdAppData.Policy.AllowedDelay != nil {
		allowedDelay := DurationSecRm(*inputPfdAppData.Policy.AllowedDelay)
		pfdAppData.AllowedDelay = &allowedDelay
	}
	if inputPfdAppData.Policy.CachingTime != nil {
		cachingTime := DurationSecRo(*inputPfdAppData.Policy.CachingTime)
		pfdAppData.CachingTime = &cachingTime
	}
	if inputPfdAppData.Policy.Pfds != nil {
		pfdAppData.Pfds = make(map[string]Pfd)
	}

	for _, inputPfdData := range inputPfdAppData.Policy.Pfds {
		pfdAppData.Pfds[inputPfdData.PfdID] = Pfd(inputPfdData)
	}

	return pfdAppData
}

func printPfdReports(pfdReportData []byte) {
	//Convert the json PFD Transaction Report into struct
	pfdReports := []PfdReport{}
	err := json.Unmarshal(pfdReportData, &pfdReports)
	if err != nil {
		klog.Info(err)
		return
	}
	fmt.Println("PFD Transaction ID:")
	fmt.Println("    Application IDs:")
	appStatus := make(map[string]string)
	for _, v := range pfdReports {
		for _, str := range v.ExternalAppIds {
			appStatus[str] = string(v.FailureCode)
		}
	}
	for k, v := range appStatus {
		fmt.Printf("      - %s : Failed (Reason: %s)\n", k, v)
	}
}

func printPfdReport(pfdReportData []byte) {
	//Convert the json PFD Transaction Report into struct
	pfdReport := PfdReport{}
	err := json.Unmarshal(pfdReportData, &pfdReport)
	if err != nil {
		klog.Info(err)
		return
	}
	fmt.Println("PFD Transaction ID:")
	fmt.Println("    Application IDs:")

	var appIDList []string
	if pfdReport.ExternalAppIds != nil {
		copy(appIDList, pfdReport.ExternalAppIds)
	}

	for _, v := range appIDList {
		fmt.Printf("      - %s : Failed (Reason: %s)\n", v, string(pfdReport.FailureCode))
	}
}
