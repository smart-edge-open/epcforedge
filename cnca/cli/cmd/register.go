// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var tac int

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register controller to AF services registry",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var s LocationService
		s.TAC = tac
		s.DNAI, _ = cmd.Flags().GetString("dnai")
		s.DNN, _ = cmd.Flags().GetString("dnn")
		s.PriDNS, _ = cmd.Flags().GetString("priDns")
		s.SecDNS, _ = cmd.Flags().GetString("secDns")
		s.UPFIP, _ = cmd.Flags().GetString("upfIp")
		s.SNSSAI, _ = cmd.Flags().GetString("snssai")

		srv, err := json.Marshal(s)
		if err != nil {
			fmt.Println(err)
			return
		}
		// register AF service
		afID, err := OAM5gRegisterAFService(srv)
		if err != nil {
			klog.Info(err)
			return
		}

		fmt.Printf("Service `%s` registered successfully\n", afID)
	},
}

func init() {

	const help = `Register controller to NGC AF services registry

Usage:
  cnca register --dnai=<DNAI> --dnn=<DNN> --tac=<TAC> --priDns=<pri-DNS> 
  	--secDns=<sec-DNS> --upfIp=<UPF-IP> --snssai=<SNSSAI>

Flags:
  -h, --help       help
      --dnai       Identifies DNAI
      --dnn        Identifies data network name
      --tac        Identifies Tracking Area Code (TAC)
      --priDns     Identifies primary DNS
      --secDns     Identifies secondary DNS
      --upfIp      Identifies UPF IP address
      --snssai     Identifies SNSSAI"
`
	// add `register` command
	cncaCmd.AddCommand(registerCmd)
	registerCmd.Flags().String("dnai", "", "Identifies DNAI")
	_ = registerCmd.MarkFlagRequired("dnai")
	registerCmd.Flags().String("dnn", "", "Identifies data network name")
	_ = registerCmd.MarkFlagRequired("dnn")
	registerCmd.Flags().IntVar(&tac, "tac", 0, "Identifies Tracking Area Code (TAC)")
	_ = registerCmd.MarkFlagRequired("tac")
	registerCmd.Flags().String("priDns", "", "Identifies primary DNS")
	_ = registerCmd.MarkFlagRequired("priDns")
	registerCmd.Flags().String("secDns", "", "Identifies secondary DNS")
	_ = registerCmd.MarkFlagRequired("secDns")
	registerCmd.Flags().String("upfIp", "", "Identifies UPF IP address")
	_ = registerCmd.MarkFlagRequired("upfIp")
	registerCmd.Flags().String("snssai", "", "Identifies SNSSAI")
	_ = registerCmd.MarkFlagRequired("snssai")
	registerCmd.SetHelpTemplate(help)
}
