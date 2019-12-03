// Copyright 2019 Intel Corporation. All rights reserved
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
