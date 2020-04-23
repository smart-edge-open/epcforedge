/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

package ngcnef

func initializePcfClient(cfg Config) PcfPolicyAuthorization {

	return NewPCFClient(&cfg)

}
