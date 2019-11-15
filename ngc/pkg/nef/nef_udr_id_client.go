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
	"log"
)

// UdrClientStub is an implementation of the Udr Influence data
type UdrClientStub struct {
	udr string
}

// NewUDRClient creates a new Udr Client
func NewUDRClient(cfg *Configuration) *UdrClientStub {

	c := &UdrClientStub{}
	c.udr = "UDR Stub"
	return c
}

// UdrInfluenceDataCreate is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataCreate(ctx context.Context,
	body TrafficInfluData, iid InfluenceID) (UdrInfluenceResponse, error) {

	log.Print("UdrInfluenceDataCreate Stub Entered")
	_ = ctx
	_ = body
	_ = iid
	udrPr := UdrInfluenceResponse{}
	err := errors.New("stub implementation")
	log.Print("UdrInfluenceDataCreate Stub Exited")
	return udrPr, err
}

// UdrInfluenceDataUpdate is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataUpdate(ctx context.Context,
	body TrafficInfluDataPatch, iid InfluenceID) (UdrInfluenceResponse,
	error) {
	log.Print("UdrInfluenceDataUpdate Stub Entered")
	_ = ctx
	_ = body
	_ = iid
	udrPr := UdrInfluenceResponse{}
	err := errors.New("stub implementation")
	log.Print("UdrInfluenceDataUpdate Stub Exited")
	return udrPr, err
}

// UdrInfluenceDataDelete is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataDelete(ctx context.Context,
	iid InfluenceID) (UdrInfluenceResponse, error) {

	log.Print("UdrInfluenceDataDelete Stub Entered")
	_ = ctx
	_ = iid
	udrPr := UdrInfluenceResponse{}
	err := errors.New("stub implementation")
	log.Print("UdrInfluenceDataDelete Stub Exited")
	return udrPr, err
}

// UdrInfluenceDataGet is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataGet(ctx context.Context) (
	UdrInfluenceResponse, error) {
	log.Print("UdrInfluenceDataGet Stub Entered")
	_ = ctx
	udrPr := UdrInfluenceResponse{}
	err := errors.New("stub implementation")
	log.Print("UdrInfluenceDataGet Stub Exited")
	return udrPr, err
}
