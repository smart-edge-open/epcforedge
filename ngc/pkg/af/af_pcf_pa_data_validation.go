// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019-2020 Intel Corporation

package af

import (
	"encoding/json"
	"errors"
	"net/http"
)

type paSuppFeat struct {
	inflOnTrafficRouting bool
	sponsConnectivity    bool
	medComnVersioning    bool
}

/* This function Decodes the Event Subscription rqeuest*/
func decodeValidateEventSubscReq(r *http.Request, w http.ResponseWriter,
	evsReqData *EventsSubscReqData) (err error) {

	if r.Body != nil && r.ContentLength > 0 {
		err = json.NewDecoder(r.Body).Decode(evsReqData)
		if err != nil {
			logPolicyRespErr(&w, "Json decode error in "+
				"DeletePolicyAuthAppSession: "+err.Error(),
				http.StatusBadRequest)
			return err
		}

		err = validateEventSubsc(evsReqData)
		if err != nil {
			return err
		}
	}
	return nil
}

/* This function validates Event Subscription*/
func validateEventSubsc(evsReqData *EventsSubscReqData) (err error) {

	if evsReqData == nil {
		return nil
	}

	length := len(evsReqData.Events)
	if length == 0 {
		err = errors.New("validateEventSubsc: events subscription " +
			" array empty")
		return err
	}

	for i := 0; i < length; i++ {
		event := evsReqData.Events[i].Event
		switch event {
		case AccessTypeChange, FailedResourcesAllocation, PlmnChg,
			QosNotif, ResourceAllocated, UsageReport:
		default:
			err = errors.New("validateEventSubsc: Invalid event ")
			return err
		}
	}
	return nil
}

/* This function performs the following functionality:
- Validates the AppSessionContext for the mandatory parameters
- Updates the notificationURI/NotifURI in ascReqData and stores Notification
   params*/

func validateAndSetupAppSessCtx(appSess *AppSessionContext, evInfo *EventInfo,
	afCtx *Context) (err error) {

	if err = validateAppSessCtx(appSess); err == nil {
		err = setAppSessNotifParams(appSess, evInfo, afCtx)
	}

	return err

}

/* This function valides the NotifURI, suppFeat, EvSubsc and Medcomponents*/
func validateAppSessCtx(appSess *AppSessionContext) (err error) {

	var suppFeatures paSuppFeat
	ascReqData := appSess.AscReqData
	if ascReqData == nil {
		err = errors.New("validateAppSessCtx: nil ascReqData")
		return err
	}

	if ascReqData.NotifURI == "" {
		err = errors.New("validateAppSessCtx: nil notifURI in req data")
		return err
	}

	if ascReqData.UeIpv4 == "" && ascReqData.UeIpv6 == "" &&
		ascReqData.UeMac == "" {
		err = errors.New("validateAppSessCtx: UE Addr info not present" +
			"in asc Req data")
		return err
	}

	suppFeat := ascReqData.SuppFeat
	if len(suppFeat) == 0 {
		err = errors.New("validateAppSessCtx: empty SuppFeat in" +
			" req data")
		return err
	}

	suppFeatures, err = decodeSuppFeat(suppFeat)
	if err != nil {
		return err
	}

	err = validateSuppFeatAscReqData(ascReqData, &suppFeatures)
	if err != nil {
		return err
	}

	err = validateEventSubsc(ascReqData.EvSubsc)
	if err != nil {
		return err
	}

	err = validateMedComponents(ascReqData.MedComponents,
		&suppFeatures)
	if err != nil {
		return err
	}

	return nil
}

/* This function validates the Supported Features and fields w.r.t the
individual feature*/
func validateSuppFeatAscReqData(ascReqData *AppSessionContextReqData,
	suppFeatures *paSuppFeat) (err error) {

	if suppFeatures.inflOnTrafficRouting {
		if ascReqData.AfRoutReq == nil {
			err = errors.New("validateAppSessCtx: " +
				"nil AfRoutReq")
			return err
		}

		if ascReqData.Dnn == "" {
			err = errors.New("validateAppSessCtx: " +
				"Dnn is nil in Req data")
			return err
		}
	}

	if suppFeatures.sponsConnectivity {
		if ascReqData.AspID == "" {
			err = errors.New("validateAppSessCtx: " +
				"AspID is nil in Req data")
			return err
		}

		if ascReqData.SponID == "" {
			err = errors.New("validateAppSessCtx: " +
				"SponID is nil in Req data")
			return err
		}

		if !(ascReqData.SponStatus == SponsorEnabled ||
			ascReqData.SponStatus == SponsorDisabled) {
			err = errors.New("validateAppSessCtx: " +
				"SponStatus is invalid in Req data")
			return err
		}

		if ascReqData.EvSubsc != nil &&
			ascReqData.EvSubsc.UsgThres == nil {
			err = errors.New("validateAppSessCtx: " +
				"UsgThres is nil in EventsSubscReqData")
			return err
		}
	}
	return nil
}

/* This function validates the AppSessionContextUpdateData mandatory
parameters */
func validateAscUpdateData(ascUpdateData *AppSessionContextUpdateData) (
	err error) {

	err = validateEventSubsc(ascUpdateData.EvSubsc)
	if err != nil {
		return err
	}

	err = validateMedComponents(ascUpdateData.MedComponents, nil)
	if err != nil {
		return err
	}
	return nil
}

/* This function decodes the Supported Features bitmap */
func decodeSuppFeat(suppFeat string) (retVal paSuppFeat, err error) {

	retVal.inflOnTrafficRouting = false
	retVal.sponsConnectivity = false
	retVal.medComnVersioning = false

	var parsedFeat byte
	length := len(suppFeat)
	if length == 1 {
		parsedFeat = suppFeat[0]
	} else if length > 1 {
		for i := 0; i < length-1; i++ {
			if suppFeat[i] != '0' {
				err = errors.New("Invalid supported feature")
				return retVal, err
			}
		}
		index := length - 1
		parsedFeat = suppFeat[index]
	}

	if parsedFeat < '0' || parsedFeat > '7' {
		err = errors.New("Invalid supported feature")
		return retVal, err
	}

	bitMask := parsedFeat - '0'
	if bitMask&1 == 1 {
		retVal.inflOnTrafficRouting = true
	}

	if bitMask&2 == 2 {
		retVal.sponsConnectivity = true
	}

	if bitMask&4 == 4 {
		retVal.medComnVersioning = true
	}

	return retVal, nil
}

/* This function validates the Media Components w.r.t features enabled*/
func validateMedComponents(medCompns map[string]MediaComponent,
	feat *paSuppFeat) (err error) {

	if len(medCompns) == 0 {
		return nil
	}

	for _, medComp := range medCompns {
		if medComp.MedCompN == 0 {
			err = errors.New("validateMedCompn: MedCompN = 0")
			return err
		}

		if len(medComp.MedSubComps) != 0 {
			for _, medSubComp := range medComp.MedSubComps {
				if medSubComp.FNum == 0 {
					err = errors.New("validateMedCompn: " +
						"fNum = 0")
					return err
				}
			}
		}

		if feat == nil {
			continue
		}

		if feat.inflOnTrafficRouting {
			if medComp.AfRoutReq == nil {
				err = errors.New("validateMedCompn: " +
					"nil AfRoutReq")
				return err
			}
		}

		if feat.medComnVersioning {
			if medComp.ContVer == 0 {
				err = errors.New("validateMedCompn: " +
					"ContVer = 0 ")
				return err
			}
		}
	}
	return nil
}

/* This function validates the Termination Request from PCF */
func validateTermInfo(termInfo *TerminationInfo) (err error) {

	if len(termInfo.ResURI) == 0 {
		err = errors.New("ResURI is nil")
		return err
	}

	if termInfo.TermCause == AllSDFDeactivated ||
		termInfo.TermCause == PDUSessionTerminated {
		err = errors.New("Invalid TermCause")
		return err
	}
	return nil
}
