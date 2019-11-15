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

// PcfClientStub is an implementation of the Pcf Authorization
type PcfClientStub struct {
	pcf string
}

// Configuration to be removed
type Configuration struct {
}

// NewPCFClient creates a new PCF Client
func NewPCFClient(cfg *Configuration) *PcfClientStub {

	c := &PcfClientStub{}
	c.pcf = "PCF Stub"
	return c
}

// PcfPolicyAuthorizationCreate is a stub implementation
func (pcf *PcfClientStub) PcfPolicyAuthorizationCreate(ctx context.Context,
	body AppSessionContext) (AppSessionID, PcfPolicyResponse, error) {

	log.Infof("PcfPolicyAuthorizationCreate Stub Entered")
	_ = ctx
	_ = body
	appSessionID := AppSessionID("appSessionId")
	pcfPr := PcfPolicyResponse{}
	err := errors.New("stub implementation")
	log.Infof("PcfPolicyAuthorizationCreate Stub Exited")
	return appSessionID, pcfPr, err
}

// PolicyAuthorizationUpdate is a stub implementation
func (pcf *PcfClientStub) PolicyAuthorizationUpdate(ctx context.Context,
	body AppSessionContextUpdateData,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PolicyAuthorizationUpdate Stub Entered")
	_ = ctx
	_ = body
	_ = appSessionID
	pcfPr := PcfPolicyResponse{}
	err := errors.New("stub implementation")
	log.Infof("PolicyAuthorizationUpdate Stub Exited")
	return pcfPr, err
}

// PolicyAuthorizationDelete is a stub implementation
func (pcf *PcfClientStub) PolicyAuthorizationDelete(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {

	log.Infof("PolicyAuthorizationDelete Stub Entered")
	_ = ctx
	_ = appSessionID
	pcfPr := PcfPolicyResponse{}
	err := errors.New("stub implementation")
	log.Infof("PolicyAuthorizationDelete Stub Exited")
	return pcfPr, err
}

// PolicyAuthorizationGet is a stub implementation
func (pcf *PcfClientStub) PolicyAuthorizationGet(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PolicyAuthorizationGet Stub Entered")
	_ = ctx
	_ = appSessionID
	pcfPr := PcfPolicyResponse{}
	err := errors.New("stub implementation")
	log.Infof("PolicyAuthorizationGet Stub Exited")
	return pcfPr, err
}
