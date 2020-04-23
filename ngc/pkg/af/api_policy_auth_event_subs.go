// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	_context "context"
)

// Linger please
var (
	_ _context.Context
)

/*
 * This file will include the code for api related to event subscription
 * and unsublscription.
 */

// PolicyAuthEventSubsAPIService EventsSubscriptionDocumentApi service
type PolicyAuthEventSubsAPIService policyAuthService
