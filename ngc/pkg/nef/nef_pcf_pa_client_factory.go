/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

package ngcnef

func initializePcfClient(cfg Config) PcfPolicyAuthorization {

	if cfg.PcfPolicyAuthorizationConfig == nil {
		return NewPCFClient(&cfg)
	}
	return NewPCFPolicyAuthHTTPClient(&cfg)
}
