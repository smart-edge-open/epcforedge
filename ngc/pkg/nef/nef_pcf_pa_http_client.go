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

//NewPCFPolicyAuthHTTPClient creates a new PCF Client
func NewPCFPolicyAuthHTTPClient(cfg *Config) *PcfClient {

	c := &PcfClient{}
	c.Pcf = "PCF PA Client"

	c.HTTPClient = &http.Client{
		Timeout: 15 * time.Second,
	}
	if cfg.PcfPolicyAuthorizationConfig.Scheme == "https" {
		CACert, err1 := ioutil.ReadFile(cfg.PcfPolicyAuthorizationConfig.ClientCert)
		if err1 != nil {
			fmt.Printf("NEF Certification loading Error: %v", err1)
			return nil

		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(CACert)
		var tlsConfig *tls.Config
		if cfg.PcfPolicyAuthorizationConfig.InsecureSkipVerify == true {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}

		} else {
			tlsConfig = &tls.Config{
				RootCAs: caCertPool,
			}

		}

		c.HTTPClient.Transport = &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	c.PcfRootURI = cfg.PcfPolicyAuthorizationConfig.Scheme + "://" + cfg.PcfPolicyAuthorizationConfig.APIRoot
	c.PcfURI = cfg.PcfPolicyAuthorizationConfig.URI

	log.Infoln("PCF Client created with the following configuration:")
	log.Infoln("Scheme: ", cfg.PcfPolicyAuthorizationConfig.Scheme)
	log.Infoln("ApiRoot: ", cfg.PcfPolicyAuthorizationConfig.APIRoot)
	log.Infoln("Uri: ", cfg.PcfPolicyAuthorizationConfig.URI)
	log.Infoln("InsecureSkipVerify: ", cfg.PcfPolicyAuthorizationConfig.InsecureSkipVerify)
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
	appSessionID := AppSessionID("")
	var resbody []byte
	log.Infof("Triggering PCF /* POST */ :" + apiURL)

	postbody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed go error in marshalling POST body error:%s", err)
		goto END

	}
	req, err = http.NewRequest("POST", apiURL, bytes.NewBuffer(postbody))
	// Add user-agent header and content-type header
	req.Header.Set("User-Agent", "NEF-OPENNESS-2006")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	if err != nil {
		fmt.Printf("Failed go error in creating POST request:%s", err)
		goto END
	}
	res, err = pcf.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("Failed go error in receiving POST response:%s", err)
		goto END
	}

	log.Infof("Body in the response =>")
	resbody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed go error in reading POST response body:%s", err)
		goto END

	}
	log.Infof(string(resbody))

	defer res.Body.Close()

	if res.StatusCode == 201 {
		appsesid = res.Header.Get("Location")
		log.Infof("appsessionid" + appsesid)
		appSessionID = AppSessionID(appsesid)
		appSessionContext = AppSessionContext{}
		err = json.Unmarshal(resbody, &appSessionContext)
		if err != nil {
			fmt.Printf("Failed go error in unmarshalling POST response body:%s", err)
			goto END
		}
		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Asc = &appSessionContext
		pcfPr.Pd = nil

	} else {
		problemDetails = ProblemDetails{}
		contentType := res.Header.Get("Content-type")
		if contentType == "application/problem+json" {
			e := json.Unmarshal(resbody, &problemDetails)
			if e != nil {
				fmt.Printf("Failed go error in unmarshalling POST response body:%s", e)
				goto END
			}
		}
		log.Infof("PCFs PolicyAuthorizationCreate failed ")
		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Pd = &problemDetails
		if err == nil {
			err = errors.New("failed post")
		}

	}
END:
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
	var resbody []byte
	var appSessionContext AppSessionContext
	var problemDetails ProblemDetails
	sessid := string(appSessionID)
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed go error in marshalling PATCH body error:%s", err)
		goto END
	}
	fmt.Println(sessid)
	log.Infof("Triggering PCF /* PATCH */ :" + apiURL + sessid)
	req, err = http.NewRequest("PATCH", apiURL+sessid, bytes.NewBuffer(b))
	req.Header.Set("User-Agent", "NEF-OPENNESS-2006")
	req.Header.Set("Content-Type", "application/merge-patch+json")
	req = req.WithContext(ctx)
	if err != nil {
		fmt.Printf("Failed go error in creating PATCH request:%s", err)
		goto END
	}
	res, err = pcf.HTTPClient.Do(req)

	if err != nil {
		fmt.Printf("Failed go error in receiving PATCH response:%s", err)
		goto END
	}

	log.Infof("Body in the response =>")
	resbody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed go error in reading PATCH response body:%s", err)
		goto END

	}
	log.Infof(string(resbody))

	appSessionID = AppSessionID(sessid)
	defer res.Body.Close()

	if res.StatusCode == 200 {
		appSessionContext = AppSessionContext{}
		err = json.Unmarshal(resbody, &appSessionContext)
		if err != nil {
			fmt.Printf("Failed go error in unmarshalling PATCH response body:%s", err)
			goto END
		}
		log.Infof("PCFs PolicyAuthorizationUpdate AppSessionID %s updated",
			string(appSessionID))

		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Asc = &appSessionContext

	} else {
		problemDetails = ProblemDetails{}
		contentType := res.Header.Get("Content-type")
		if contentType == "application/problem+json" {
			err = json.Unmarshal(resbody, &problemDetails)
			if err != nil {
				fmt.Printf("Failed go error in unmarshalling PATCH response body:%s", err)
				goto END
			}
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
END:
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
	var resbody []byte
	apiURL := pcf.PcfRootURI + pcf.PcfURI + sessid + "/delete"

	var err error
	log.Infof("Triggering PCF /* DELETE */ : " + apiURL)
	req, err = http.NewRequest("POST", apiURL, nil)
	req.Header.Set("User-Agent", "NEF-OPENNESS-2006")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	if err != nil {
		fmt.Printf("Failed go error in creating DELETE request:%s", err)
		goto END
	}
	res, err = pcf.HTTPClient.Do(req)

	if err != nil {
		fmt.Printf("Failed go error in receiving DELETE response:%s", err)
		goto END
	}

	if res.StatusCode == 204 {
		log.Infof("PCFs PolicyAuthorizationDelete AppSessionID %s found",
			sessid)
		pcfPr.ResponseCode = uint16(res.StatusCode)

	} else if res.StatusCode == 200 {
		//var eventnoti EventsNotification
		log.Infof("Body in the response =>")
		resbody, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Failed go error in reading DELETE response body:%s", err)
			goto END

		}
		defer res.Body.Close()
		log.Infof(string(resbody))
		/* err = json.Unmarshal(body, &eventnoti)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		} */
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
END:
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
	var req *http.Request
	var appSessionContext AppSessionContext
	var problemDetails ProblemDetails
	var err error
	var resbody []byte
	log.Infof("Triggering PCF /* GET */ : " + apiURL + sessid)
	req, err = http.NewRequest("GET", apiURL+sessid, nil)
	if err != nil {
		fmt.Printf("Failed go error in creating GET request:%s", err)
		goto END

	}
	req.Header.Set("User-Agent", "NEF-OPENNESS-2006")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	res, err = pcf.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("Failed go error in creating GET response:%s", err)
		goto END

	}
	log.Infof("Body in the response =>")
	resbody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed go error in reading GET response body:%s", err)
		goto END

	}
	log.Infof(string(resbody))
	defer res.Body.Close()

	if res.StatusCode == 200 {
		appSessionContext = AppSessionContext{}
		err = json.Unmarshal(resbody, &appSessionContext)
		if err != nil {
			fmt.Printf("Failed go error in unmarshalling POST response body:%s", err)
			goto END
		}
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s found",
			string(appSessionID))

		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Asc = &appSessionContext

	} else {
		problemDetails = ProblemDetails{}
		contentType := res.Header.Get("Content-type")
		if contentType == "application/problem+json" {
			err = json.Unmarshal(resbody, &problemDetails)
			if err != nil {
				fmt.Printf("Failed go error in unmarshalling POST response body:%s", err)
				goto END
			}
		}
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s not found",
			string(appSessionID))
		if err == nil {
			err = errors.New("failed get")
		}

		pcfPr.ResponseCode = uint16(res.StatusCode)
		pcfPr.Pd = &problemDetails
	}
END:
	return pcfPr, err
}
