/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

package ngcnef

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	//"strconv"
)

func closeReqBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		log.Errf("response body was not closed properly")
	}
}

func sendCustomeErrorRspToAF(w http.ResponseWriter, eCode int,
	custTitleString string) {

	eRsp := nefSBRspData{errorCode: eCode,
		pd: ProblemDetails{Title: custTitleString}}

	sendErrorResponseToAF(w, eRsp)

}
func sendErrorResponseToAF(w http.ResponseWriter, rsp nefSBRspData) {

	mdata, eCode := createErrorJSON(rsp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(eCode)
	_, err := w.Write(mdata)

	if err != nil {
		log.Err("NEF ERROR : Failed to send response to AF !!!")
	}
	log.Infof("HTTP Response sent: %d", eCode)
}

/*
func send500PFDResponseToAF(w http.ResponseWriter, rsp nefSBRspData,
	pfdReportList map[string]PfdReport) {

	mdata, eCode := createErrorJSON(rsp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(eCode)
	_, err := w.Write(mdata)
	if err != nil {
		log.Err("NEF ERROR : Failed to send response to AF !!!")
	}

	if pfdReportList != nil {
		body, e := json.Marshal(pfdReportList)
		if e != nil {
			log.Err("NEF ERROR : Failed to send response to AF !!!")
		}
		_, err = w.Write(body)
		if err != nil {
			log.Err("NEF ERROR : Failed to send response to AF !!!")
		}
	}

	log.Infof("HTTP Response sent: %d", eCode)
}

*/
func createErrorJSON(rsp nefSBRspData) (mdata []byte, statusCode int) {

	var err error
	statusCode = 404

	/*
		TBD for future: Removed check for 401, 403, 413 and 429
		due cyclometrix complexity lint warning. Once a better mechanism
		is found to be added back. Anyhow currently these errors are not
		supported
	*/

	if rsp.errorCode == 400 || rsp.errorCode == 404 || rsp.errorCode == 411 ||
		rsp.errorCode == 415 || rsp.errorCode == 500 || rsp.errorCode == 503 {
		statusCode = rsp.errorCode
		mdata, err = json.Marshal(rsp.pd)

		if err == nil {
			/*No return */
			log.Info(statusCode)
			return mdata, statusCode
		}
	}
	/*Send default error */
	pd := ProblemDetails{Title: " NEF Error "}

	mdata, err = json.Marshal(pd)

	if err != nil {
		return mdata, statusCode
	}
	/*Any case return mdata */
	return mdata, statusCode
}

func logNef(nef *nefData) {

	log.Infof("AF count %+v", len(nef.afs))
	if len(nef.afs) > 0 {
		for key, value := range nef.afs {
			log.Infof(" AF ID : %+v, Sub Registered Count %+v",
				key, len(value.subs))
			for _, vs := range value.subs {
				log.Infof("   SubId : %+v, ServiceId: %+v", vs.subid,
					vs.ti.AfServiceID)
			}

		}
	}

}

// LoadJSONConfig reads a file located at configPath and unmarshals it to
// config structure
func loadJSONConfig(configPath string, config interface{}) error {
	cfgData, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return err
	}
	return json.Unmarshal(cfgData, config)
}

func generatePfdReport(appID string,
	failureReason string, pfdReportList map[string]PfdReport) {

	// TBD Prepare the pfd report
	/*
		switch failureReason {

		case "APP_ID_DUPLICATED":

			if val, ok := pfdReportList["APP_ID_DUPLICATED"]; !ok {
				// Create the first PFD report
				var appIds []string
				appIds = append(appIds, appId)
				pfdReport := PfdReport{ExternalAppIds: appIds,
					FailureCode: AppIDDuplicated}
				PfdReport["APP_ID_DUPLICATED"] = pfdReport

			} else {
				pfdReportList["APP_ID_DUPLICATED"].ExternalAppIds =
					append(pfdReportList["APP_ID_DUPLICATED"].ExternalAppIds,
					 appId)
			}

		}
	*/

}
