/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

/* Client implementation of the pcf stub */

package ngcnef

import (
	"context"
	"math/rand"
	"strconv"
)

// PcfClientStub is an implementation of the Pcf Authorization
type PcfClientStub struct {
	pcf       string
	initialID int
	// database to store the contents of the app session contexts created
	paDb map[int]AppSessionContext
}

// NewPCFClient creates a new PCF Client
func NewPCFClient(cfg *Config) *PcfClientStub {

	c := &PcfClientStub{}
	c.pcf = "PCF Stub"
	/* Generate a randome number for currSessionId */
	c.initialID = rand.Intn(10000)
	c.paDb = make(map[int]AppSessionContext)
	log.Infof("PCF Stub Client created with initial session id: %d",
		c.initialID)
	return c
}

// genAppSessionID - creates a new session id to be used
func genAppSessionID(pcf *PcfClientStub) int {

	size := len(pcf.paDb)
	log.Infof("PCFs Policy Authorization DB size : %d", size)
	sessid := pcf.initialID
	for i := 0; i < size; i++ {
		_, prs := pcf.paDb[sessid]
		if !prs {
			break
		}
		sessid++
	}
	log.Infof("PCFs Policy Authorization AppSessionId created: %d", sessid)
	return sessid
}

// PolicyAuthorizationCreate is a stub implementation
// Successful response : 201 and body contains AppSessionContext
func (pcf *PcfClientStub) PolicyAuthorizationCreate(ctx context.Context,
	body AppSessionContext) (AppSessionID, PcfPolicyResponse, error) {

	log.Infof("PCFs PolicyAuthorizationCreate Entered")
	_ = ctx

	var err error
	pcfPr := PcfPolicyResponse{}
	// generated a session id return the same body as provided in the request
	sessid := genAppSessionID(pcf)
	Asc := body
	pcf.paDb[sessid] = Asc
	pcfPr.ResponseCode = 201
	pcfPr.Asc = &Asc
	pcfPr.Pd = nil
	appSessionID := AppSessionID(strconv.Itoa(sessid))
	log.Infof("PCFs PolicyAuthorizationCreate [CorrId,NotifUrl,DnaiChgType]"+
		" => [%s,%s,%s]", body.AscReqData.AfRoutReq.UpPathChgSub.NotifCorreID,
		body.AscReqData.AfRoutReq.UpPathChgSub.NotificationURI,
		body.AscReqData.AfRoutReq.UpPathChgSub.DnaiChgType)

	log.Infof("PCFs PolicyAuthorizationCreate Exited successfully with "+
		"sessid: %s", appSessionID)
	return appSessionID, pcfPr, err
}

// PolicyAuthorizationUpdate is a stub implementation
// Successful response : 200 and body contains AppSessionContext
func (pcf *PcfClientStub) PolicyAuthorizationUpdate(ctx context.Context,
	body AppSessionContextUpdateData,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PCFs PolicyAuthorizationUpdate Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx

	var err error
	pcfPr := PcfPolicyResponse{}
	// convert the appsession id to integer
	sessid, _ := strconv.Atoi(string(appSessionID))
	// check for the presence of the sessid in the database
	asc, prs := pcf.paDb[sessid]
	// if not found return an error i.e 404
	if !prs {
		log.Infof("PCFs PolicyAuthorizationUpdate AppSessionID %s not found",
			string(appSessionID))
		pcfPr.ResponseCode = 404
	} else {
		log.Infof("PCFs PolicyAuthorizationUpdate AppSessionID %s updated",
			string(appSessionID))

		asc.AscReqData.AfAppID = body.AfAppID
		asc.AscReqData.AfRoutReq = body.AfRoutReq
		pcf.paDb[sessid] = asc
		pcfPr.ResponseCode = 204
		pcfPr.Asc = &asc

	}
	log.Infof("PCFs PolicyAuthorizationUpdate Exited for AppSessionID %s",
		string(appSessionID))
	return pcfPr, err
}

// PolicyAuthorizationDelete is a stub implementation
// Successful response : 204 and empty body
func (pcf *PcfClientStub) PolicyAuthorizationDelete(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {

	log.Infof("PCFs PolicyAuthorizationDelete Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx

	var err error
	pcfPr := PcfPolicyResponse{}
	// convert the appsession id to integer
	sessid, _ := strconv.Atoi(string(appSessionID))
	// check for the presence of the sessid in the database
	_, prs := pcf.paDb[sessid]
	// if not found return an error i.e 404
	if !prs {
		log.Infof("PCFs PolicyAuthorizationDelete AppSessionID %s not found",
			string(appSessionID))
		pcfPr.ResponseCode = 404
	} else {
		log.Infof("PCFs PolicyAuthorizationDelete AppSessionID %s found",
			string(appSessionID))
		delete(pcf.paDb, sessid)
		log.Infof("PCFs Policy Authorization DB size : %d", len(pcf.paDb))
		pcfPr.ResponseCode = 204

	}
	log.Infof("PCFs PolicyAuthorizationDelete Stub Exited for AppSessionID %s",
		string(appSessionID))
	return pcfPr, err
}

// PolicyAuthorizationGet is a stub implementation
// Successful response : 204 and empty body
func (pcf *PcfClientStub) PolicyAuthorizationGet(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PCFs PolicyAuthorizationGet Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx

	var err error
	pcfPr := PcfPolicyResponse{}
	// convert the appsession id to integer
	sessid, _ := strconv.Atoi(string(appSessionID))
	// check for the presence of the sessid in the database
	asc, prs := pcf.paDb[sessid]
	// if not found return an error i.e 404
	if !prs {
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s not found",
			string(appSessionID))
		pcfPr.ResponseCode = 404
	} else {
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s found",
			string(appSessionID))
		pcfPr.ResponseCode = 200
		pcfPr.Asc = &asc

	}
	log.Infof("PCFs PolicyAuthorizationGet Exited for AppSessionID %s",
		string(appSessionID))
	return pcfPr, err
}
