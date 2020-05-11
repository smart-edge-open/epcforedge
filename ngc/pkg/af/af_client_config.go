// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019-2020 Intel Corporation

package af

import (
	"net/http"
)

// CliConfig struct
type CliConfig struct {
	Protocol       string `json:"Protocol"`
	NEFHostname    string `json:"NEFHostname"`
	NEFPort        string `json:"NEFPort"`
	NEFBasePath    string `json:"NEFBasePath"`
	NEFPFDBasePath string `json:"NEFPFDBasePath"`
	UserAgent      string `json:"UserAgent"`
	NEFCliCertPath string `json:"NEFCliCertPath"`
	HTTPClient     *http.Client
	OAuth2Support  bool `json:"OAuth2Support"`
}

// NewConfiguration function initializes client configuration
func NewConfiguration(afCtx *Context) *CliConfig {

	cfg := &CliConfig{
		Protocol:       afCtx.cfg.CliCfg.Protocol,
		NEFPort:        afCtx.cfg.CliCfg.NEFPort,
		NEFHostname:    afCtx.cfg.CliCfg.NEFHostname,
		NEFBasePath:    afCtx.cfg.CliCfg.NEFBasePath,
		NEFPFDBasePath: afCtx.cfg.CliCfg.NEFPFDBasePath,
		UserAgent:      afCtx.cfg.CliCfg.UserAgent,
		NEFCliCertPath: afCtx.cfg.CliCfg.NEFCliCertPath,
		OAuth2Support:  afCtx.cfg.CliCfg.OAuth2Support,
	}

	return cfg
}

// GenericCliConfig struct
type GenericCliConfig struct {
	Protocol      string `json:"Protocol"`
	ProtocolVer   string `json:"ProtocolVer"`
	Hostname      string `json:"Hostname"`
	Port          string `json:"Port"`
	BasePath      string `json:"BasePath"`
	CliCertPath   string `json:"CliCertPath"`
	OAuth2Support bool   `json:"OAuth2Support"`
	VerifyCerts   bool   `json:"VerifyCerts"`
	NotifURI      string `json:"NotifURI"`
}

// NewCliPcfConfiguration create new client pcf config struct
func NewCliPcfConfiguration(afCtx *Context) *GenericCliConfig {

	cfg := &GenericCliConfig{
		Protocol:      afCtx.cfg.CliPcfCfg.Protocol,
		ProtocolVer:   afCtx.cfg.CliPcfCfg.ProtocolVer,
		Port:          afCtx.cfg.CliPcfCfg.Port,
		Hostname:      afCtx.cfg.CliPcfCfg.Hostname,
		BasePath:      afCtx.cfg.CliPcfCfg.BasePath,
		CliCertPath:   afCtx.cfg.CliPcfCfg.CliCertPath,
		OAuth2Support: afCtx.cfg.CliPcfCfg.OAuth2Support,
		NotifURI:      afCtx.cfg.CliPcfCfg.NotifURI,
		VerifyCerts:   afCtx.cfg.CliPcfCfg.VerifyCerts,
	}

	return cfg
}
