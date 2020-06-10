// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use: "apply",
	Short: "Apply LTE CUPS userplane or NGC AF TI subscription using YAML " +
		"configuration file",
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

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

			// create new subscription
			var subLoc string
			subLoc, err = AFCreateSubscription(sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Subscription URI:", subLoc)

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

			// create new LTE userplane
			upID, err := LteCreateUserplane(up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Userplane:", upID)

		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

// pfdApplyCmd represents the apply command
var pfdApplyCmd = &cobra.Command{
	Use: "apply",
	Short: "Apply NGC AF PFD Transaction or NGC AF PFD Application " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

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

		// create new AF PFD Transaction
		pfdData, self, err := AFCreatePfdTransaction(trans)
		if err != nil {
			klog.Info(err)
			if err.Error() == "HTTP failure: 500" && pfdData != nil {
				printPfdReports(pfdData)
			}
			return
		}

		if pfdData != nil {
			printPdfTransStatus(pfdData, self)
		} else {
			fmt.Printf("PFD Transaction URI: %s\n", self)
			fmt.Printf("PFD Transaction ID: %s\n",
				getTransIDFromURL(self))
		}
	},
}

func init() {

	const help = `Apply LTE CUPS userplane or NGC AF TI subscription
using YAML configuration file

Usage:
  cnca apply -f <config.yml>

Example:
  cnca apply -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	const pfdHelp = `Apply NGC AF PFD Transaction or NGC AF PFD Application
using YAML configuration file

Usage:
  cnca pfd apply -f <config.yml>

Example:
  cnca pfd apply -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`
	// add `apply` command
	cncaCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = applyCmd.MarkFlagRequired("filename")
	applyCmd.SetHelpTemplate(help)

	// add pfd `apply` command
	pfdCmd.AddCommand(pfdApplyCmd)
	pfdApplyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = pfdApplyCmd.MarkFlagRequired("filename")
	pfdApplyCmd.SetHelpTemplate(pfdHelp)
}

func readInputData(cmd *cobra.Command) ([]byte, error) {
	ymlFile, _ := cmd.Flags().GetString("filename")
	if ymlFile == "" {
		return nil, errors.New("YAML file missing")
	}

	data, err := ioutil.ReadFile(ymlFile)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getTransIDFromURL(url string) string {
	urlElements := strings.Split(url, "/")
	for index, str := range urlElements {
		if str == "transactions" {
			return urlElements[index+1]
		}
	}
	return ""
}

func getPfdTransData(inputPfdTransData AFPfdManagement) PfdManagement {

	var pfdTransData PfdManagement

	if inputPfdTransData.Policy.SuppFeat != nil {
		*pfdTransData.SuppFeat =
			SupportedFeatures(*inputPfdTransData.Policy.SuppFeat)
	}

	if inputPfdTransData.Policy.PfdDatas != nil {
		pfdTransData.PfdDatas = make(map[string]PfdData)
	}

	for _, inputPfdAppData := range inputPfdTransData.Policy.PfdDatas {
		var pfdAppData PfdData

		pfdAppData.ExternalAppID = inputPfdAppData.ExternalAppID
		pfdAppData.Self = Link(inputPfdAppData.Self)

		if inputPfdAppData.AllowedDelay != nil {
			allowedDelay := DurationSecRm(*inputPfdAppData.AllowedDelay)
			pfdAppData.AllowedDelay = &allowedDelay
		}
		if inputPfdAppData.CachingTime != nil {
			cachingTime := DurationSecRo(*inputPfdAppData.CachingTime)
			pfdAppData.CachingTime = &cachingTime
		}
		if inputPfdAppData.Pfds != nil {
			pfdAppData.Pfds = make(map[string]Pfd)
		}

		for _, inputPfdData := range inputPfdAppData.Pfds {
			pfdAppData.Pfds[inputPfdData.PfdID] = Pfd(inputPfdData)
		}
		pfdTransData.PfdDatas[pfdAppData.ExternalAppID] = pfdAppData
	}
	return pfdTransData
}

func printPdfTransStatus(pfdTransData []byte, transURI string) {

	//Convert the json PFD Transaction data into struct
	pfdTrans := PfdManagement{}
	err := json.Unmarshal(pfdTransData, &pfdTrans)
	if err != nil {
		klog.Info(err)
		return
	}
	if transURI == "" {
		transURI = string(pfdTrans.Self)
	}

	fmt.Printf("PFD Transaction URI: %s\n", transURI)
	fmt.Printf("PFD Transaction ID: %s\n",
		getTransIDFromURL(transURI))
	fmt.Println("    Application IDs:")

	appStatus := make(map[string]string)
	for k := range pfdTrans.PfdDatas {
		appStatus[k] = "Created"
	}
	for _, v := range pfdTrans.PfdReports {
		for _, str := range v.ExternalAppIds {
			appStatus[str] = string(v.FailureCode)
		}
	}
	for k, v := range appStatus {
		if v != "Created" {
			fmt.Printf("      - %s : Failed (Reason: %s)\n", k, v)
		} else {
			fmt.Printf("      - %s : %s\n", k, v)
		}
	}
}
