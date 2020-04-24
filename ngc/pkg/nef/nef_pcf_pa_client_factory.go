/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

package ngcnef

func initializePcfClient(cfg Config) PcfPolicyAuthorization {

	if cfg.PcfPolicyAuthorizationConfig == nil {
		return NewPCFClient(&cfg)
	} /*else if TestPcf == true {
		return NewPCFTestClient(&cfg)
	} else {
		return NewPCFClientF(&cfg)
	}*/
	return NewPCFClientF(&cfg)
}

//NewPCFTestClient is for UT PCF client
/* func NewPCFTestClient(cfg *Config) *PcfClient {

	c := &PcfClient{}
	c.Pcf = "PCF test client freegc"

	c.HTTPClient = &http.Client{}

	c.PcfRootURI = cfg.PcfPolicyAuthorizationConfig.Scheme + "://" + cfg.PcfPolicyAuthorizationConfig.APIRoot
	c.PcfURI = cfg.PcfPolicyAuthorizationConfig.URI

	log.Infof("PCF test Client created")
	return c
} */
