/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

/* Client implementation of the UDR Stub */

package ngcnef

import (
	"context"
)

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

	return rsp, err
}
