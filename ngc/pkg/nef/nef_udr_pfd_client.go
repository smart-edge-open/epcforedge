/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

/* Client implementation of the UDR Stub */

package ngcnef

import (
	"context"
	"errors"
)

// TestClient variable is only for UnitTesting purpose to inject errors in stub
var TestClient = false

// UdrPfdClientStub is an implementation of the Udr Influence data
type UdrPfdClientStub struct {
	udr string
	//database to store content of udr PFD data
	appPfd map[string]*PfdDataForApp
}

// NewUDRPfdClient creates a new Udr Client
func NewUDRPfdClient(cfg *Config) *UdrPfdClientStub {

	c := &UdrPfdClientStub{}
	c.udr = "UDR Pfd Stub"
	c.appPfd = make(map[string]*PfdDataForApp)
	log.Info("UDR PFD Stub Client created")
	return c
}

// UdrPfdDataCreate is a stub implementation
func (udr *UdrPfdClientStub) UdrPfdDataCreate(ctx context.Context,
	body PfdDataForApp) (rsp UdrPfdResponse, err error) {

	log.Infof("UdrPfdDataCreate Stub Entered")
	_ = ctx

	log.Info("UdrPfdDataCreate: Invoke UDR SB PUT -> ")
	udr.appPfd[string(body.AppID)] = &body
	log.Infof("UdrPfdDataCreate Stub Exited")

	if TestClient {
		rsp.ResponseCode = 400
		return rsp, errors.New("Error in UDR SB PUT")
	}
	return rsp, err
}

// UdrPfdDataGet is a stub implementation
func (udr *UdrPfdClientStub) UdrPfdDataGet(ctx context.Context,
	appID UdrAppID) (rsp UdrPfdResponse, err error) {
	log.Infof("UdrPfdDataGet Stub Entered")
	_ = ctx

	rsp.AppPfd = udr.appPfd[string(appID)]
	log.Info("Get PFD for AppId : ", appID)
	log.Info("UdrPfdDataGet: Invoke UDR SB GET -> ")

	log.Infof("UdrPfdDataGet Stub Exited")

	if TestClient {
		rsp.ResponseCode = 404
		return rsp, errors.New("Error in UDR SB GET")
	}
	return rsp, err
}

// UdrPfdDataDelete is a stub implementation
func (udr *UdrPfdClientStub) UdrPfdDataDelete(ctx context.Context,
	appID UdrAppID) (rsp UdrPfdResponse, err error) {
	log.Infof("UdrPfdDataDelete Stub Entered")
	_ = ctx
	log.Info("Deleted PFD AppId : ", appID)
	log.Info("UdrPfdDataDelete: Invoke UDR SB DELETE -> ")

	delete(udr.appPfd, string(appID))
	log.Infof("UdrPfdDataDelete Stub Exited")

	if TestClient {
		rsp.ResponseCode = 400
		return rsp, errors.New("Error in UDR SB DELETE")
	}

	return rsp, err
}
