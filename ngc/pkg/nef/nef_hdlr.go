/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import (
	"context"
	"errors"
)

const correlationIDOffset = 20
const subNotFound string = "Subscription Not Found"
const pfdNotFound string = "PFD transaction Not Found"
const appNotFound string = "Application in PFD transaction Not Found"

//NEF context data
type nefData struct {
	ctx                  context.Context
	afCount              int
	locationURLPrefix    string
	locationURLPrefixPfd string
	pcfClient            PcfPolicyAuthorization
	udrClient            UdrData
	corrID               uint
	afs                  map[string]*afData
	upfNotificationURL   URI
}

//NEFSBGetFn is the callback for SB API
type NEFSBGetFn func(subData *afSubscription, nefCtx *nefContext) (
	sub TrafficInfluSub, rsp nefSBRspData, err error)

//NEFSBPutFn is the callback for SB API
type NEFSBPutFn func(subData *afSubscription, nefCtx *nefContext,
	ti TrafficInfluSub) (rsp nefSBRspData, err error)

//NEFSBPatchFn is the callback for SB API
type NEFSBPatchFn func(subData *afSubscription, nefCtx *nefContext,
	tisp TrafficInfluSubPatch) (rsp nefSBRspData, err error)

//NEFSBDeleteFn is the callback for SB API
type NEFSBDeleteFn func(subData *afSubscription, nefCtx *nefContext) (
	rsp nefSBRspData, err error)

/**********************************/
/* PFD handler function type  	  */
/**********************************/

// NEFSBGetPfdFn is the callback for SB API to get PFD transaction
type NEFSBGetPfdFn func(transData *afPfdTransaction, nefCtx *nefContext) (
	trans PfdManagement, rsp nefPFDSBRspData, err error)

// NEFSBPutPfdFn is the callback for SB API to put PFD transaction
type NEFSBPutPfdFn func(transData *afPfdTransaction, nefCtx *nefContext,
	trans PfdManagement) (rsp map[string]nefPFDSBRspData, err error)

// NEFSBAppPutPfdFn is the callback for SB API to put PFD transaction
type NEFSBAppPutPfdFn func(transData *afPfdTransaction, nefCtx *nefContext,
	app PfdData) (rsp nefPFDSBRspData, err error)

// NEFSBDeletePfdFn is the callback for SB API to delete PFD transaction
type NEFSBDeletePfdFn func(transData *afPfdTransaction, nefCtx *nefContext) (
	rsp nefPFDSBRspData, err error)

/**********************************/

//PCF Subscription data
type afSubscription struct {
	subid string
	ti    TrafficInfluSub

	//Applicable in case of single UE case only
	appSessionID AppSessionID

	//Applicable in case of Multiple UE only
	iid                       InfluenceID
	NotifCorreID              string
	afNotificationDestination Link
	NEFSBGet                  NEFSBGetFn
	NEFSBPut                  NEFSBPutFn
	NEFSBPatch                NEFSBPatchFn
	NEFSBDelete               NEFSBDeleteFn
}

//PFD transaction data
type afPfdTransaction struct {
	transID       string
	pfdManagement PfdManagement

	NEFSBPfdGet    NEFSBGetPfdFn
	NEFSBPfdPut    NEFSBPutPfdFn
	NEFSBAppPfdPut NEFSBAppPutPfdFn
	NEFSBPfdDelete NEFSBDeletePfdFn
}

//AF data
type afData struct {
	afID       string
	subIDnum   int
	transIDnum int
	maxSubSupp int
	subs       map[string]*afSubscription
	pfdtrans   map[string]*afPfdTransaction
}

type nefSBRspData struct {
	errorCode int
	pd        ProblemDetails
}

type nefPFDSBRspData struct {
	result nefSBRspData
	//fc     *FailureCode
}

//Creates a AF instance
func (af *afData) afCreate(nefCtx *nefContext, afID string) error {

	//Validate afid ??

	af.afID = afID
	af.subIDnum = nefCtx.cfg.SubStartID //Start Number
	af.transIDnum = nefCtx.cfg.PfdTransStartID
	af.maxSubSupp = nefCtx.cfg.MaxSubSupport
	af.subs = make(map[string]*afSubscription)
	//PFD transaction
	af.pfdtrans = make(map[string]*afPfdTransaction)
	return nil
}

//Initialize the NEF component
func (nef *nefData) nefCreate(ctx context.Context, cfg Config) error {

	nef.ctx = ctx
	nef.afCount = 0
	nef.pcfClient = NewPCFClient(&cfg)
	if nef.pcfClient == nil {
		return errors.New("PCF Client creation failed")
	}
	nef.udrClient = NewUDRClient(&cfg)
	if nef.udrClient == nil {
		return errors.New("PCF Client creation failed")
	}
	nef.afs = make(map[string]*afData)
	nef.corrID = uint(cfg.SubStartID + correlationIDOffset)

	if cfg.NefAPIRoot == "" {
		return errors.New("NefAPIRoot is empty")
	}

	if cfg.LocationPrefix == "" {
		return errors.New("NEF LocationPrefix is empty")
	}

	// Generate the location url prefix
	nef.locationURLPrefix = getNefLocationURLPrefix(&cfg)
	log.Infof("NEF Location URL Prefix :%s", nef.locationURLPrefix)

	// Generate the location url prefix for PFD
	nef.locationURLPrefixPfd = getNefLocationURLPrefixPfd(&cfg)
	log.Infof("NEF Location URL Prefix :%s", nef.locationURLPrefixPfd)

	// Genereate the notification url
	if cfg.UpfNotificationResURIPath == "" {
		return errors.New("UpfNotificationResURIPath is empty")
	}
	nef.upfNotificationURL = getNefNotificationURI(&cfg)
	log.Infof("SMF UPF Notification URL :%s", nef.upfNotificationURL)
	return nil
}

/*
func NEFInit() error {

	return nef.nefCreate()
}*/

func (nef *nefData) nefAddAf(nefCtx *nefContext, afID string) (af *afData,
	err error) {

	var afe afData

	if len(nef.afs) >= nefCtx.cfg.MaxAFSupport {
		log.Infoln("MAX AF exceeded ")
		return af, errors.New("MAX AF exceeded")
	}

	//Check if AF is already present
	_, ok := nef.afs[afID]

	if ok {
		return nef.afs[afID], errors.New("AF already present")
	}

	//Create a new entry of AF

	_ = afe.afCreate(nefCtx, afID)
	nef.afs[afID] = &afe
	nef.afCount++

	return &afe, nil
}

func (nef *nefData) nefCheckPfdAppIDExists(appID string) bool {

	for _, v := range nef.afs {
		for _, trans := range v.pfdtrans {
			for key := range trans.pfdManagement.PfdDatas {
				if key == appID {
					return true
				}

			}
		}
	}
	return false

}

func (nef *nefData) nefGetAf(afID string) (af *afData, err error) {

	//Check if AF is already present
	afe, ok := nef.afs[afID]

	if ok {
		return afe, nil
	}
	err = errors.New("AF entry not present")
	return afe, err
}

func (nef *nefData) nefDeleteAf(afID string) (err error) {

	//Check if AF is already present
	_, ok := nef.afs[afID]

	if ok {
		delete(nef.afs, afID)
		//nef.afCount--
		return nil
	}

	err = errors.New("AF entry not present")
	return err
}

func (nef *nefData) nefDestroy() {

	// Todo
}
