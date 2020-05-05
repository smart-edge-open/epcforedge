// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019-2020 Intel Corporation

package af

// GenericCliConfig struct
type GenericCliConfig struct {
	Protocol          string `json:"Protocol"`
	ProtocolVer       string `json:"ProtocolVer"`
	Hostname          string `json:"Hostname"`
	Port              string `json:"Port"`
	BasePath          string `json:"BasePath"`
	LocationPrefixURI string `json:"LocationPrefixURI"`
	CliCertPath       string `json:"CliCertPath"`
	OAuth2Support     bool   `json:"OAuth2Support"`
	VerifyCerts       bool   `json:"VerifyCerts"`
	NotifURI          string `json:"NotifURI"`
}

// NewCliPcfConfiguration create new client pcf config struct
func NewCliPcfConfiguration(afCtx *Context) *GenericCliConfig {

	cfg := &GenericCliConfig{
		Protocol:          afCtx.cfg.CliPcfCfg.Protocol,
		ProtocolVer:       afCtx.cfg.CliPcfCfg.ProtocolVer,
		Port:              afCtx.cfg.CliPcfCfg.Port,
		Hostname:          afCtx.cfg.CliPcfCfg.Hostname,
		BasePath:          afCtx.cfg.CliPcfCfg.BasePath,
		LocationPrefixURI: afCtx.cfg.CliPcfCfg.LocationPrefixURI,
		CliCertPath:       afCtx.cfg.CliPcfCfg.CliCertPath,
		OAuth2Support:     afCtx.cfg.CliPcfCfg.OAuth2Support,
		NotifURI:          afCtx.cfg.CliPcfCfg.NotifURI,
		VerifyCerts:       afCtx.cfg.CliPcfCfg.VerifyCerts,
	}

	return cfg
}
