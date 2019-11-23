// Copyright 2019 Intel Corporation, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ngcaf

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
func NewConfiguration(afCtx *afContext) *CliConfig {

	cfg := &CliConfig{
		NEFBasePath:    afCtx.cfg.CliCfg.NEFBasePath,
		UserAgent:      afCtx.cfg.CliCfg.UserAgent,
		NEFCliCertPath: afCtx.cfg.CliCfg.NEFCliCertPath,
	}

	return cfg
}
