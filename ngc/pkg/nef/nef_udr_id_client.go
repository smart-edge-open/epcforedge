/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

/* Client implementation of the UDR Stub */

package ngcnef

import (
	"context"
)

// UdrClientStub is an implementation of the Udr Influence data
type UdrClientStub struct {
	udr string
	// database to store the contents of the udr influence data
	tidDb map[string]TrafficInfluData
}

// NewUDRClient creates a new Udr Client
func NewUDRClient(cfg *Config) *UdrClientStub {

	c := &UdrClientStub{}
	c.udr = "UDR Stub"
	c.tidDb = make(map[string]TrafficInfluData)
	log.Info("UDR Stub Client created")
	return c
}

// UdrInfluenceDataCreate is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataCreate(ctx context.Context,
	body TrafficInfluData, iid InfluenceID) (UdrInfluenceResponse, error) {

	log.Infof("UDRs InfluenceDataCreate Entered for %s", string(iid))
	_ = ctx

	var err error
	udrPr := UdrInfluenceResponse{}
	// generated a session id return the same body as provided in the request
	Tid := body
	_, psc := udr.tidDb[string(iid)]
	if !psc {
		udrPr.ResponseCode = 201
		log.Infof("UDRs InfluenceDataCreate entry not found and creating")
	} else {
		udrPr.ResponseCode = 200
		log.Infof("UDRs InfluenceDataCreate entry found and updating")
	}
	udr.tidDb[string(iid)] = Tid

	udrPr.Tid = &Tid
	udrPr.Pd = nil

	log.Infof("UDRs InfluenceDataCreate [CorrId,NotifUrl]  => [%s,%s]",
		body.UpPathChgNotifCorreID, body.UpPathChgNotifURI)
	log.Infof("UDRs UdrInfluenceDataCreate Exited for %s", string(iid))
	return udrPr, err
}

func updateTidWithTidPatch(tid *TrafficInfluData,
	body *TrafficInfluDataPatch) {
	if body.UpPathChgNotifCorreID != "" {
		tid.UpPathChgNotifCorreID = body.UpPathChgNotifCorreID
	}
	tid.AppReloInd = body.AppReloInd
	if body.Dnn != "" {
		tid.Dnn = body.Dnn
	}
	if body.Snssai.Sd != "" {
		tid.Snssai = body.Snssai
	}
	if body.InternalGroupID != "" {
		tid.InterGroupID = body.InternalGroupID
	}
	if len(body.EthTrafficFilters) > 0 {
		tid.EthTrafficFilters = body.EthTrafficFilters
	}
	if body.Supi != "" {
		tid.Supi = body.Supi
	}
	if len(body.TrafficFilters) > 0 {
		tid.TrafficFilters = body.TrafficFilters
	}
	if len(body.TrafficRoutes) > 0 {
		tid.TrafficRoutes = body.TrafficRoutes
	}
	if body.ValidEndTime != "" {
		tid.ValidEndTime = body.ValidEndTime
	}
	if body.ValidStartTime != "" {
		tid.ValidStartTime = body.ValidStartTime
	}
	tid.NwAreaInfo = body.NwAreaInfo
	if body.UpPathChgNotifURI != "" {
		tid.UpPathChgNotifURI = body.UpPathChgNotifURI
	}

}

// UdrInfluenceDataUpdate is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataUpdate(ctx context.Context,
	body TrafficInfluDataPatch, iid InfluenceID) (UdrInfluenceResponse,
	error) {
	log.Infof("UDRs InfluenceDataUpdate Entered for %s", string(iid))
	_ = ctx

	var err error
	udrPr := UdrInfluenceResponse{}
	// check for the presence of the sessid in the database
	tid, prs := udr.tidDb[string(iid)]
	// if not found return an error i.e 404
	if !prs {
		log.Infof("UDRs InfluenceDataUpdate InfluenceID %s not found",
			string(iid))
		udrPr.ResponseCode = 404
	} else {
		log.Infof("UDRs InfluenceDataUpdate InfluenceID %s updated",
			string(iid))
		updateTidWithTidPatch(&tid, &body)
		udr.tidDb[string(iid)] = tid
		udrPr.ResponseCode = 204
		udrPr.Tid = &tid

	}
	log.Infof("UDRs UdrInfluenceDataUpdate Exited for %s", string(iid))
	return udrPr, err
}

// UdrInfluenceDataDelete is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataDelete(ctx context.Context,
	iid InfluenceID) (UdrInfluenceResponse, error) {

	log.Infof("UDRs InfluenceDataDelete for %s", string(iid))
	_ = ctx

	var err error
	udrPr := UdrInfluenceResponse{}
	// check for the presence of the sessid in the database
	_, prs := udr.tidDb[string(iid)]
	// if not found return an error i.e 404
	if !prs {
		log.Infof("UDRs InfluenceDataDelete InfluenceID %s not found",
			string(iid))
		udrPr.ResponseCode = 404
	} else {
		log.Infof("UDRs InfluenceDataDelete InfluenceID %s found",
			string(iid))
		delete(udr.tidDb, string(iid))
		log.Infof("UDRs Influence Data DB size : %d", len(udr.tidDb))
		udrPr.ResponseCode = 204

	}
	log.Infof("UDRs InfluenceDataDelete Exited for %s", string(iid))
	return udrPr, err
}

// UdrInfluenceDataGet is a stub implementation
func (udr *UdrClientStub) UdrInfluenceDataGet(ctx context.Context) (
	UdrInfluenceResponse, error) {
	log.Infof("UdrInfluenceDataGet Stub Entered")
	_ = ctx
	udrPr := UdrInfluenceResponse{}
	var err error
	log.Infof("UdrInfluenceDataGet Stub Exited")
	return udrPr, err
}
