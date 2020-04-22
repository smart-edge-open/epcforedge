// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019-2020 Intel Corporation

package af

import (
	"net/http"
)

// CliConfig struct
type CliPcfConfig struct {
	Protocol           string       `json:"Protocol"`
	PcfHostname        string       `json:"PcfHostname"`
	PcfPort            string       `json:"PcfPort"`
	PolicyAuthBasePath string       `json:"PolicyAuthBasePath"`
	UserAgent          string       `json:"UserAgent"`
	CliCertPath        string       `json:"CliCertPath"`
	HTTPClient         *http.Client // Change HTTPClient to HttpClient
	OAuth2Support      bool         `json:"OAuth2Support"`
	Debug              bool         `json:"Debug"`
}

// NewConfiguration function initializes client configuration
func NewCliPcfConfiguration(afCtx *Context) *CliPcfConfig {

	cfg := &CliPcfConfig{
		Protocol:           afCtx.cfg.CliPcfCfg.Protocol,
		PcfPort:            afCtx.cfg.CliPcfCfg.PcfPort,
		PcfHostname:        afCtx.cfg.CliPcfCfg.PcfHostname,
		PolicyAuthBasePath: afCtx.cfg.CliPcfCfg.PolicyAuthBasePath,
		UserAgent:          afCtx.cfg.CliPcfCfg.UserAgent,
		CliCertPath:        afCtx.cfg.CliPcfCfg.CliCertPath,
		OAuth2Support:      afCtx.cfg.CliPcfCfg.OAuth2Support,
		Debug:              afCtx.cfg.CliPcfCfg.Debug,
	}

	return cfg
}
