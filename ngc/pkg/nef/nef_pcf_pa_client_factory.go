/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

package ngcnef

func initializePcfClient(cfg Config) PcfPolicyAuthorization {

	if cfg.PcfPolicyAuthorizationConfig == nil {
		log.Infof("PcfPolicyAuthorizationConfig is not configured")
		return NewPCFClient(&cfg)
	}
	pcf, err := NewPCFPolicyAuthHTTPClient(&cfg)
	if err == nil {
		return pcf
	}
	return nil
}
