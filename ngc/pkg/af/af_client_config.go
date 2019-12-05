// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"net/http"
)

// CliConfig struct
type CliConfig struct {
	NEFBasePath    string `json:"NEFBasePath"`
	UserAgent      string `json:"UserAgent"`
	NEFCliCertPath string `json:"NEFCliCertPath"`
	HTTPClient     *http.Client
}

// NewConfiguration function initializes client configuration
func NewConfiguration(afCtx *Context) *CliConfig {

	cfg := &CliConfig{
		NEFBasePath:    afCtx.cfg.CliCfg.NEFBasePath,
		UserAgent:      afCtx.cfg.CliCfg.UserAgent,
		NEFCliCertPath: afCtx.cfg.CliCfg.NEFCliCertPath,
	}

	return cfg
}
