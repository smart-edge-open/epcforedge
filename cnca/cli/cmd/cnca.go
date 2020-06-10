// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
)

type HTTPProtocol int

const (
	HTTP HTTPProtocol = 1 + iota
	HTTP2
)

// HTTP2/HTTPS constants
const (
	UseHTTPProtocol      = HTTP2
	TLSCAFile            = "root-ca-cert.pem"
	DefaultTLSCAFilePath = "/etc/certs/"
)

var HTTP2ClientTLSCAPath string

// cncaCmd represents the base command when called without any subcommands
var cncaCmd = &cobra.Command{
	Use:          "cnca",
	Long:         "Core Newtwork Configuration Agent (CNCA) command line",
	SilenceUsage: true,
}

// pfdCmd represents the pfd related commands
var pfdCmd = &cobra.Command{
	Use:          "pfd",
	Short:        "Applies/Creates and Manages NGC PFD Transaction and Application",
	Args:         cobra.MaximumNArgs(5),
	SilenceUsage: true,
}

func init() {

	const help = `
  Applies/Creates and Manages NGC PFD Transaction and Application

Usage:
  Create PFD transaction:        
      cnca pfd apply -f <config.yml>
  Get PFD transactions:          
      cnca pfd get transactions
  Get single PFD transaction:    
      cnca pfd get transaction <transaction-id>
  Update single PFD transaction: 
      cnca pfd patch transaction <transaction-id> -f <config.yml>
  Delete single PFD transaction: 
      cnca pfd delete transaction <transaction-id>
  Get single application for a PFD  transaction: 
	  cnca pfd get transaction <transaction-id> application <application-id>
  Update single application in a PFD  transaction:
      cnca pfd patch transaction <transaction-id> application <application-id> -f <config.yml>
  Delete single application in a PFD  transaction:
      cnca pfd delete transaction <transaction-id> -application <application-id>
	  
Flags:
  -h, --help       help
  -f, --filename   YAML configuration file

`
	// add `pfd` command
	cncaCmd.AddCommand(pfdCmd)
	pfdCmd.SetHelpTemplate(help)
}

// Execute CNCA agent
func Execute() error {

	if UseHTTPProtocol == HTTP2 {
		if HTTP2ClientTLSCAPath == "" {
			HTTP2ClientTLSCAPath = DefaultTLSCAFilePath
		}
		http2ClientTLSCAData := HTTP2ClientTLSCAPath + "/" + TLSCAFile
		err := InitHTTP2Client(http2ClientTLSCAData)
		if nil != err {
			fmt.Printf("Failure in Initializing HTTP2 Client: %v\n", err)
			return err
		}
	} else {
		InitHTTPClient()
	}
	return cncaCmd.Execute()
}

func InitHTTP2Client(clientCertData string) error {
	CACert, err := ioutil.ReadFile(clientCertData)
	if err != nil {
		return err
	}

	CACertPool := x509.NewCertPool()
	CACertPool.AppendCertsFromPEM(CACert)

	if UseHTTPProtocol == HTTP2 {
		client = http.Client{
			Timeout: 10 * time.Second,
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: CACertPool,
				},
			},
		}
	} else {
		err = errors.New("Incorrect HTTP Protocol Configured")
		return err
	}
	return nil
}

func InitHTTPClient() {
	client = http.Client{
		Timeout: 10 * time.Second,
	}
}
