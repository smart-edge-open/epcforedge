/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

package ngcnef

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// PcfClient is an implementation of the Pcf Authorization
type PcfClient struct {
	Pcf        string
	HTTPClient *http.Client
	PcfRootURI string
	PcfURI     string
}

//NewPCFRestClient creates a new PCF Client
func NewPCFRestClient(cfg *Config) *PcfClient {

	c := &PcfClient{}
	c.Pcf = "PCF freegc"

	c.HTTPClient = &http.Client{
		Timeout: 15 * time.Second,
	}

	CACert, err1 := ioutil.ReadFile(cfg.PcfPolicyAuthorizationConfig.ClientCert)
	if err1 != nil {
		fmt.Printf("NEF Certification loading Error: %v", err1)
		return nil

	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(CACert)

	tlsConfig := &tls.Config{
		//RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	c.HTTPClient.Transport = &http2.Transport{
		TLSClientConfig: tlsConfig,
	}

	c.PcfRootURI = cfg.PcfPolicyAuthorizationConfig.Scheme + "://" + cfg.PcfPolicyAuthorizationConfig.APIRoot
	c.PcfURI = cfg.PcfPolicyAuthorizationConfig.URI

	log.Infof("PCF Client created")
	return c
}

//PolicyAuthorizationCreate is a actual implementation
// Successful response : 201 and body contains AppSessionContext
func (pcf *PcfClient) PolicyAuthorizationCreate(ctx context.Context,
	body AppSessionContext) (AppSessionID, PcfPolicyResponse, error) {

	log.Infof("PCFs PolicyAuthorizationCreate Entered")
	_ = ctx

	pcfPr := PcfPolicyResponse{}
	apiURL := pcf.PcfRootURI + pcf.PcfURI
	var appsesid string
	var req *http.Request
	var res *http.Response

	var appSessionContext AppSessionContext
	var problemDetails ProblemDetails

	log.Infof("Triggering PCF /* POST */ :" + apiURL)

	postbody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	req, err = http.NewRequest("POST", apiURL, bytes.NewBuffer(postbody))

	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}
	res, err = pcf.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}
	appSessionID := AppSessionID("")
	log.Infof("Body in the response =>")
	resbody, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(resbody))

	defer res.Body.Close()

	if res.StatusCode == 201 {
		appsesid = res.Header.Get("Location")
		log.Infof("appsessionid" + appsesid)
		appSessionID = AppSessionID(appsesid)
		appSessionContext = AppSessionContext{}
		err = json.Unmarshal(resbody, &appSessionContext)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}
		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Asc = &appSessionContext
		pcfPr.Pd = nil

	} else {
		problemDetails = ProblemDetails{}

		e := json.Unmarshal(resbody, &problemDetails)
		if e != nil {
			fmt.Printf("Failed go error :%s", e)
		}
		log.Infof("PCFs PolicyAuthorizationCreate failed ")
		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Pd = &problemDetails
		if err == nil {
			err = errors.New("failed post")
		}

	}

	return appSessionID, pcfPr, err
}

// PolicyAuthorizationUpdate is a actual implementation
// Successful response : 200 and body contains AppSessionContext
func (pcf *PcfClient) PolicyAuthorizationUpdate(ctx context.Context,
	body AppSessionContextUpdateData,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PCFs PolicyAuthorizationUpdate Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx

	pcfPr := PcfPolicyResponse{}
	apiURL := pcf.PcfRootURI + pcf.PcfURI
	var req *http.Request
	var res *http.Response

	var appSessionContext AppSessionContext
	var problemDetails ProblemDetails
	sessid := string(appSessionID)
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(sessid)
	log.Infof("Triggering PCF /* PATCH */ :" + apiURL + sessid)
	req, err = http.NewRequest("PATCH", apiURL+sessid, bytes.NewBuffer(b))
	if err != nil {
		log.Infof("Failed go error :%s", err)
	}
	res, err = pcf.HTTPClient.Do(req)

	if err != nil {
		fmt.Printf("Failed go error :%s", err)

	}

	log.Infof("Body in the response =>")
	resbody, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(resbody))

	appSessionID = AppSessionID(sessid)
	defer res.Body.Close()

	if res.StatusCode == 200 {
		appSessionContext = AppSessionContext{}
		err = json.Unmarshal(resbody, &appSessionContext)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}
		log.Infof("PCFs PolicyAuthorizationUpdate AppSessionID %s updated",
			string(appSessionID))

		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Asc = &appSessionContext

	} else {
		problemDetails = ProblemDetails{}
		err = json.Unmarshal(resbody, &problemDetails)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}
		log.Infof("PCFs PolicyAuthorizationUpdate AppSessionID %s not found",
			string(appSessionID))
		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Pd = &problemDetails
		if err == nil {
			err = errors.New("failed patch")
		}
	}

	log.Infof("PCFs PolicyAuthorizationUpdate Exited for AppSessionID %s",
		string(appSessionID))
	return pcfPr, err
}

// PolicyAuthorizationDelete is a actual implementation
// Successful response : 204 and empty body
func (pcf *PcfClient) PolicyAuthorizationDelete(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {

	log.Infof("PCFs PolicyAuthorizationDelete Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx

	pcfPr := PcfPolicyResponse{}
	sessid := string(appSessionID)
	var req *http.Request
	var res *http.Response
	apiURL := pcf.PcfRootURI + pcf.PcfURI + sessid + "/delete"

	var err error
	log.Infof("Triggering PCF /* DELETE */ : " + apiURL)
	req, err = http.NewRequest("POST", apiURL, nil)
	if err != nil {
		fmt.Println("Failed go error :", err)
	}
	res, err = pcf.HTTPClient.Do(req)

	if err != nil {
		fmt.Println("Failed go error :", err)
	}

	if res.StatusCode == 204 {
		log.Infof("PCFs PolicyAuthorizationDelete AppSessionID %s found",
			sessid)
		pcfPr.ResponseCode = uint16(res.StatusCode)

	} else if res.StatusCode == 200 {
		var eventnoti EventsNotification
		log.Infof("Body in the response =>")
		body, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		log.Infof(string(body))
		err = json.Unmarshal(body, &eventnoti)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}
		pcfPr.ResponseCode = uint16(res.StatusCode)

	} else {
		log.Infof("PCFs PolicyAuthorizationDelete AppSessionID %s not found",
			sessid)
		if err == nil {
			err = errors.New("failed delete")
		}
		pcfPr.ResponseCode = uint16(res.StatusCode)
	}
	log.Infof("PCFs PolicyAuthorizationDelete Exited for AppSessionID %s",
		sessid)

	return pcfPr, err
}

// PolicyAuthorizationGet is a actual implementation
// Successful response : 204 and empty body
func (pcf *PcfClient) PolicyAuthorizationGet(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PCFs PolicyAuthorizationGet Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx

	apiURL := pcf.PcfRootURI + pcf.PcfURI
	pcfPr := PcfPolicyResponse{}
	sessid := string(appSessionID)
	var res *http.Response

	var appSessionContext AppSessionContext
	var problemDetails ProblemDetails

	log.Infof("Triggering PCF /* GET */ : " + apiURL + sessid)
	res, err := pcf.HTTPClient.Get(apiURL + sessid)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)

	}
	log.Infof("Body in the response =>")
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
	defer res.Body.Close()

	if res.StatusCode == 200 {
		appSessionContext = AppSessionContext{}
		err = json.Unmarshal(body, &appSessionContext)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s found",
			string(appSessionID))

		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Asc = &appSessionContext

	} else {
		problemDetails = ProblemDetails{}
		err = json.Unmarshal(body, &problemDetails)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s not found",
			string(appSessionID))
		if err == nil {
			err = errors.New("failed get")
		}

		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Pd = &problemDetails
	}

	return pcfPr, err
}
