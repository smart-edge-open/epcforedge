// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
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

/* Client implementation of the pcf stub */

package main

import (
	"context"
	"errors"
)

// AfClient is an implementation of the Af Notification
type AfClient struct {
	af string
}

// NewAfClient creates a new Udr Client
func NewAfClient(cfg *Configuration) *AfClient {

	c := &AfClient{}
	c.af = "Af Notification Client"
	return c
}

// UdrInfluencAfNotification is an implementation for sending upf event
func (af *AfClient) AfNotificationUpfEvent(ctx context.Context,
	afURI URI, body EventNotification) error {

	log.Infof("AfNotificationUpfEvent Stub Entered")
	_ = ctx
	_ = body
	_ = afURI

	err := errors.New("af stub implementation")
	log.Infof("AfNotificationUpfEvent Stub Exited")
	return err
}
