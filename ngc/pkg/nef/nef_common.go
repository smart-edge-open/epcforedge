/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngcnef

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
)

const contentType string = "application/json"

var (
	jsonCheck = regexp.MustCompile("(?i:[application|text]/json)")
	xmlCheck  = regexp.MustCompile("(?i:[application|text]/xml)")
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
	return json.Unmarshal(cfgData, &config)
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer,
	err error) {

	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if jsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	} else if xmlCheck.MatchString(contentType) {
		err = xml.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("invalid body type %s", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

//prepareRequest build the request
func prepareRequest(
	ctx context.Context,
	path string, method string,
	reqBody interface{},
	headerParams map[string]string,
) (httpRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect reqBody type and post.
	if reqBody != nil {
		reqContentType := headerParams["Content-Type"]

		body, err = setBody(reqBody, reqContentType)
		if err != nil {
			return nil, err
		}
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		log.Errf("url parsing error, path = %s", path)
		return nil, err
	}

	// Generate a new request
	if body != nil {
		httpRequest, err = http.NewRequest(method, url.String(), body)
	} else {
		httpRequest, err = http.NewRequest(method, url.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		httpRequest.Header = headers
	}

	// Add the user agent to the request.
	//httpRequest.Header.Add("User-Agent", c.UserAgent)

	if ctx != nil {
		// add context to the request
		httpRequest = httpRequest.WithContext(ctx)
	}
	return httpRequest, nil
}
func closeRespBody(r *http.Response) {
	err := r.Body.Close()
	if err != nil {
		log.Errf("response body was not closed properly")
	}
}
