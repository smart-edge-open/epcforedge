/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import "context"

/* The SB interface towards the AF for sending the notifications received
   from different NF's */

// AfNotification definesthe interfaces that are exposed for sending
// nofitifications towards the AF
type AfNotification interface {

	// AAfNotificationUpfEvent sends the UPF event through POST method
	// towards the AF
	AfNotificationUpfEvent(ctx context.Context,
		afURI URI,
		body EventNotification) error
}
