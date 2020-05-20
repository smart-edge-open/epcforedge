// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Connectivity constants
const (
	NgcOAMServiceEndpoint      = "http://127.0.0.1:30070/ngcoam/v1/af"
	NgcAFServiceEndpoint       = "http://127.0.0.1:30050/af/v1"
	LteOAMServiceEndpoint      = "http://127.0.0.1:8082"
	NgcOAMServiceHTTP2Endpoint = "https://127.0.0.1:30070/ngcoam/v1/af"
	NgcAFServiceHTTP2Endpoint  = "https://127.0.0.1:30050/af/v1"
	LteOAMServiceHTTP2Endpoint = "https://127.0.0.1:8082"
)

// HTTP client
var client http.Client

func getNgcOAMServiceURL() string {
	if UseHTTPProtocol == HTTP2 {
		return NgcOAMServiceHTTP2Endpoint + "/services"
	}
	return NgcOAMServiceEndpoint + "/services"
}

func getNgcAFServiceURL() string {
	if UseHTTPProtocol == HTTP2 {
		return NgcAFServiceHTTP2Endpoint + "/subscriptions"
	}
	return NgcAFServiceEndpoint + "/subscriptions"
}

func getNgcAFPfdServiceURL() string {
	if UseHTTPProtocol == HTTP2 {
		return NgcAFServiceHTTP2Endpoint + "/pfd/transactions"
	}
	return NgcAFServiceEndpoint + "/pfd/transactions"
}

func getNgcAFPaServiceURL() string {
	if UseHTTPProtocol == HTTP2 {
		return NgcAFServiceHTTP2Endpoint + "/policy-authorization/app-sessions"
	}
	return NgcAFServiceEndpoint + "/policy-authorization/app-sessions"
}

func getLteOAMServiceURL() string {
	if UseHTTPProtocol == HTTP2 {
		return LteOAMServiceHTTP2Endpoint + "/userplanes"
	}
	return LteOAMServiceEndpoint + "/userplanes"
}

// OAM5gRegisterAFService register controller to AF services registry
func OAM5gRegisterAFService(locService []byte) (string, error) {

	url := getNgcOAMServiceURL()

	req, err := http.NewRequest("POST", url, bytes.NewReader(locService))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var s AFServiceID
		err = json.Unmarshal(b, &s)
		if err != nil {
			return "", err
		}
		return s.AFServiceID, nil
	}
	return "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
}

// OAM5gUnregisterAFService unregister controller from AF services registry
func OAM5gUnregisterAFService(serviceID string) error {

	url := getNgcOAMServiceURL() + "/" + serviceID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// AFCreateSubscription create new Traffic Influence Subscription at AF
func AFCreateSubscription(sub []byte) (string, error) {

	url := getNgcAFServiceURL()

	req, err := http.NewRequest("POST", url, bytes.NewReader(sub))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	// retrieve URI of the newly created subscription from response header
	subLoc := resp.Header.Get("Location")
	if subLoc == "" {
		return "", fmt.Errorf("Empty subscription URI returned from AF")
	}
	return subLoc, nil
}

// AFPatchSubscription update an active subscription for the AF
func AFPatchSubscription(subID string, sub []byte) error {

	url := getNgcAFServiceURL() + "/" + subID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(sub))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// AFGetSubscription get the active Traffic Influence Subscription for the AF
func AFGetSubscription(subID string) ([]byte, error) {
	var sub []byte
	var req *http.Request
	var err error
	var url string

	if subID == "all" {
		url = getNgcAFServiceURL()
	} else {
		url = getNgcAFServiceURL() + "/" + subID
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return sub, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return sub, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return sub, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	sub, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return sub, err
	}
	return sub, nil
}

// AFDeleteSubscription delete an active Traffic Influence Subscription for AF
func AFDeleteSubscription(subID string) error {

	url := getNgcAFServiceURL() + "/" + subID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// LteCreateUserplane create new LTE userplane
func LteCreateUserplane(up []byte) (string, error) {

	url := getLteOAMServiceURL()

	req, err := http.NewRequest("POST", url, bytes.NewReader(up))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var u CupsUserplaneID
	err = json.Unmarshal(b, &u)
	if err != nil {
		return "", err
	}

	return u.ID, nil
}

// LtePatchUserplane update an active LTE CUPS userplane
func LtePatchUserplane(upID string, up []byte) error {

	url := getLteOAMServiceURL() + "/" + upID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(up))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// LteGetUserplane get the active CUPS userplane
func LteGetUserplane(upID string) ([]byte, error) {
	var up []byte
	var req *http.Request
	var err error
	var url string

	if upID == "all" {
		url = getLteOAMServiceURL()
	} else {
		url = getLteOAMServiceURL() + "/" + upID
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return up, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return up, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return up, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	up, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return up, err
	}
	return up, nil
}

// LteDeleteUserplane delete an active LTE CUPS userplane
func LteDeleteUserplane(upID string) error {

	url := getLteOAMServiceURL() + "/" + upID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// AFCreatePfdTransaction create new PFD transaction at AF
func AFCreatePfdTransaction(trans []byte) ([]byte, string, error) {

	var pfdData []byte

	url := getNgcAFPfdServiceURL()

	req, err := http.NewRequest("POST", url, bytes.NewReader(trans))
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusInternalServerError {
		return nil, "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	// retrieve URI of the newly created transaction from response header
	self := resp.Header.Get("Self")
	if resp.Body != nil {
		pfdData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, "", err
		}
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return pfdData, self, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return pfdData, self, nil
}

// AFGetPfdTransaction get the active PFD Transaction for the AF
func AFGetPfdTransaction(transID string) ([]byte, error) {
	var trans []byte
	var req *http.Request
	var err error
	var url string

	if transID == "all" {
		url = getNgcAFPfdServiceURL()
	} else {
		url = getNgcAFPfdServiceURL() + "/" + transID
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return trans, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return trans, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return trans, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	trans, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return trans, err
	}
	return trans, nil
}

// AFPatchPfdTransaction update an active PFD Transaction for the AF
func AFPatchPfdTransaction(transID string, trans []byte) ([]byte, error) {

	var pfdReports []byte

	url := getNgcAFPfdServiceURL() + "/" + transID

	req, err := http.NewRequest("PUT", url, bytes.NewReader(trans))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusInternalServerError {
		return nil, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	if resp.Body != nil {
		pfdReports, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return pfdReports, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return pfdReports, nil
}

// AFDeletePfdTransaction delete an active PFD Transaction for the AF
func AFDeletePfdTransaction(transID string) error {

	url := getNgcAFPfdServiceURL() + "/" + transID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// AFGetPfdApplication get the active PFD Application for the AF
func AFGetPfdApplication(transID string, appID string) ([]byte, error) {
	var trans []byte
	var req *http.Request
	var err error

	url := getNgcAFPfdServiceURL() + "/" + transID + "/applications/" + appID

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return trans, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return trans, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return trans, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	trans, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return trans, err
	}
	return trans, nil
}

// AFPatchPfdApplication update an active PFD Application for the AF
func AFPatchPfdApplication(transID string, appID string, trans []byte) ([]byte, error) {

	var pfdReports []byte
	url := getNgcAFPfdServiceURL() + "/" + transID + "/applications/" + appID

	req, err := http.NewRequest("PUT", url, bytes.NewReader(trans))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusInternalServerError {
		return nil, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	if resp.Body != nil {
		pfdReports, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return pfdReports, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return pfdReports, nil
}

// AFDeletePfdApplication delete an active PFD Application for the AF
func AFDeletePfdApplication(transID string, appID string) error {

	url := getNgcAFPfdServiceURL() + "/" + transID + "/applications/" + appID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// AFCreatePaAppSession creates new application session at AF
func AFCreatePaAppSession(appSession []byte) ([]byte, string, error) {

	var appSessionRespData []byte
	url := getNgcAFPaServiceURL()

	req, err := http.NewRequest("POST", url, bytes.NewReader(appSession))
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	// retrieve URI of the newly created transaction from response header
	appSessionLoc := resp.Header.Get("Location")
	if resp.Body != nil {
		appSessionRespData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, "", err
		}
	}

	return appSessionRespData, appSessionLoc, nil
}

// AFGetPaAppSession gets the active application session at AF
func AFGetPaAppSession(appSessionID string) ([]byte, error) {
	var appSession []byte
	var req *http.Request
	var err error

	url := getNgcAFPaServiceURL() + "/" + appSessionID

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return appSession, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return appSession, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return appSession, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	appSession, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return appSession, err
	}
	return appSession, nil
}

// AFPatchPaAppSession update an active application session at AF
func AFPatchPaAppSession(appSessionID string, appSession []byte) ([]byte, error) {

	var appSessionResp []byte
	url := getNgcAFPaServiceURL() + "/" + appSessionID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(appSession))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	if resp.Body != nil {
		appSessionResp, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return appSessionResp, nil
}

// AFPatchPaAppSession delete an active application session at AF
func AFDeletePaAppSession(appSessionID string) error {

	url := getNgcAFPaServiceURL() + "/" + appSessionID + "/delete"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}
	return nil
}

// AFPaEventSubscribe create or modify event subscription sub resource at AF
func AFPaEventSubscribe(appSessionID string, evSubscReqData []byte) ([]byte, error) {
	var appSession []byte
	url := getNgcAFPaServiceURL() + "/" + appSessionID + "/" + "events-subscription"

	req, err := http.NewRequest("PUT", url, bytes.NewReader(evSubscReqData))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	if resp.Body != nil {
		appSession, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return appSession, nil
}

//AFPaEventUnsubscribe deletes an active event subscription sub resource at AF
func AFPaEventUnsubscribe(appSessionID string) error {

	url := getNgcAFPaServiceURL() + "/" + appSessionID + "/events-subscription"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}
	return nil
}
