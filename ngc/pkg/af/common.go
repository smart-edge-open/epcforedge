// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019-2020 Intel Corporation

package af

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	oauth2 "github.com/otcshare/epcforedge/ngc/pkg/oauth2"
	"golang.org/x/net/http2"
)

//TransIDMax var
var TransIDMax = math.MaxInt32

func getPcfOAuth2Token() (token string, err error) {

	token, err = oauth2.GetAccessToken()
	if err == nil {
		log.Infoln("Got Pcf OAuth2 Access Token: " + token)
	}
	return token, err
}

func getStatusCode(r *http.Response) int {
	if r != nil {
		return r.StatusCode
	}
	return http.StatusInternalServerError
}

func genAFTransID(trans TransactionIDs) int {
	var (
		num   int
		min   = 1
		found = true
	)

	for max := range trans {
		num =
			max
		break
	}
	for max := range trans {
		if max > num {
			num = max
		}
	}

	if num == TransIDMax {
		num = min
	}
	//look for a free ID until it is <= math.MaxInt32 is achieved again
	for found && num < TransIDMax {
		num++
		//check if the ID is in use, if not - return the ID
		if _, found = trans[num]; !found {
			trans[num] = TrafficInfluSub{}
			return num
		}
	}
	return 0
}

func genHTTPClient(cfg *GenericCliConfig) (*http.Client, error) {

	if cfg.Protocol == "https" {
		var tlsCfg *tls.Config

		if len(cfg.CliCertPath) == 0 {
			return nil, errors.New("Certificate Path Nil")
		}
		CACert, err := ioutil.ReadFile(cfg.CliCertPath)
		if err != nil {
			log.Errf("Error: %v", err)
			return nil, err
		}

		CACertPool := x509.NewCertPool()
		CACertPool.AppendCertsFromPEM(CACert)

		tlsCfg = &tls.Config{
			RootCAs: CACertPool,
		}
		// Below commented code is for Debug purpose only. Uncomment it
		// to enable skipping certificate verification.
		/* ------Linter warning disable------
		if cfg.VerifyCerts {
			tlsCfg = &tls.Config{
				RootCAs: CACertPool,
			}
		} else {
			tlsCfg = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		*/

		switch cfg.ProtocolVer {
		case "1.1":
			httpClient := &http.Client{
				Timeout: 15 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig: tlsCfg,
				},
			}
			return httpClient, nil
		case "2.0":
			httpClient := &http.Client{
				Timeout: 15 * time.Second,
				Transport: &http2.Transport{
					TLSClientConfig: tlsCfg,
				},
			}
			return httpClient, nil
		default:
			err = errors.New("Unsupported protocol version" +
				cfg.ProtocolVer)
			log.Errf("%s", err.Error())
			return nil, err
		}
	} else if cfg.Protocol == "http" {
		switch cfg.ProtocolVer {
		case "1.1":
			httpClient := &http.Client{
				Timeout: 15 * time.Second,
			}
			return httpClient, nil
		case "2.0":
			httpClient := &http.Client{
				Timeout: 15 * time.Second,
				Transport: &http2.Transport{
					AllowHTTP: true,
					DialTLS: func(network, addr string,
						cfg *tls.Config) (
						net.Conn, error) {
						return net.Dial(network, addr)
					},
				},
			}
			return httpClient, nil
		default:
			err := errors.New("Unsupported protocol version" +
				cfg.ProtocolVer)
			log.Errf("%s", err.Error())
			return nil, err
		}
	}

	err := errors.New("Not recognizable Protocol")
	log.Errf("%s", err.Error())
	return nil, err
}

// prepareRequest build the request
func (c *GenericCliConfig) prepareRequest(
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
			if h == "Authorization" && c.OAuth2Support {
				v = "Bearer " + v
				headers.Set(h, v)
			} else {
				headers.Set(h, v)
			}
		}
		httpRequest.Header = headers
	}

	// Add the user agent to the request.
	httpRequest.Header.Add("User-Agent", userAgent)

	if ctx != nil {
		// add context to the request
		httpRequest = httpRequest.WithContext(ctx)
	}
	return httpRequest, nil
}

func genTransactionID(afCtx *Context) (int, error) {

	tID := genAFTransID(afCtx.transactions)
	if tID == 0 {
		return 0, errors.New("the pool of AF Transaction IDs is already used")
	}
	return tID, nil
}

func getSubsIDFromURL(u *url.URL) (string, error) {

	if u == nil {
		return "", errors.New("empty URL in the request message")
	}

	sURL := u.String()

	// It is assumed the URL address
	// ends with  "/subscriptions/{subscriptionID}"
	s := strings.Split(sURL, "subscriptions")
	switch len(s) {
	case 1:
		return "", errors.New("subscriptionID was not found " +
			"in the URL string")
	case 2:
		sID := strings.Replace(s[1], "/", "", -1)
		return sID, nil

	default:
		return "", errors.New("wrong URL")
	}
}

func getPfdTransIDFromURL(u *http.Request) string {

	vars := mux.Vars(u)
	pID := vars["transactionId"]
	return pID

}

func getPfdAppIDFromURL(u *http.Request) string {

	vars := mux.Vars(u)
	aID := vars["appId"]
	return aID

}

func handleGetErrorResp(r *http.Response,
	body []byte) error {

	newErr := GenericError{
		body:  body,
		error: r.Status,
	}

	switch r.StatusCode {
	case 400, 401, 403, 404, 406, 429, 500, 503:

		var v ProblemDetails
		log.Errf("Error from NEF server - %s", r.Status)

		if r.StatusCode == 401 {

			if fetchNEFAuthorizationToken() != nil {
				log.Infoln("Token refresh failed")
			}
		}

		err := json.Unmarshal(body, &v)
		if err != nil {
			newErr.error = err.Error()
			return newErr
		}
		newErr.model = v
		return newErr

	default:
		b, _ := ioutil.ReadAll(r.Body)
		err := fmt.Errorf("NEF returned error - %s, %s", r.Status, string(b))
		return err
	}
}

func handlePostPutPatchErrorResp(r *http.Response,
	body []byte) error {

	newErr := GenericError{
		body:  body,
		error: r.Status,
	}

	switch r.StatusCode {
	case 400, 401, 403, 404, 411, 413, 415, 429, 500, 503:

		var v ProblemDetails
		if r.StatusCode == 401 {
			if fetchNEFAuthorizationToken() != nil {
				log.Infoln("Token refresh failed")
			}
		}
		err := json.Unmarshal(body, &v)
		if err != nil {
			newErr.error = err.Error()
			return newErr
		}
		newErr.model = v
		log.Errf("NEF returned error - %s", r.Status)

		return newErr

	default:
		b, _ := ioutil.ReadAll(r.Body)
		err := fmt.Errorf("NEF returned error - %s, %s", r.Status, string(b))
		return err
	}
}

func handlePfdPostPutPatchErrorResp(r *http.Response,
	body []byte) error {

	newErr := GenericError{
		body:  body,
		error: r.Status,
	}

	switch r.StatusCode {
	case 400, 401, 403, 404, 411, 413, 415, 429, 503:

		var v ProblemDetails

		if r.StatusCode == 401 {
			if fetchNEFAuthorizationToken() != nil {
				log.Infoln("Token refresh failed")
			}
		}
		err := json.Unmarshal(body, &v)
		if err != nil {
			newErr.error = err.Error()
			return newErr
		}
		newErr.model = v
		log.Errf("NEF returned error - %s", r.Status)
		return newErr
	case 500:
		return newErr
	default:
		b, _ := ioutil.ReadAll(r.Body)
		err := fmt.Errorf("NEF returned error - %s, %s", r.Status, string(b))
		return err
	}
}

func updatePfdURL(cfg Config, r *http.Request, resURL string) string {

	res := strings.Split(resURL, "transactions")

	var afURL string
	if HTTP2Enabled == true {
		afURL = "https" + "://" + cfg.SrvCfg.Hostname +
			cfg.SrvCfg.CNCAEndpoint + cfg.LocationPrefixPfd +
			"transactions" + res[1]
	} else {
		afURL = "http" + "://" + cfg.SrvCfg.Hostname +
			cfg.SrvCfg.CNCAEndpoint + cfg.LocationPrefixPfd +
			"transactions" + res[1]
	}
	return afURL

}

func updateSelfLink(cfg Config, r *http.Request,
	pfdTrans PfdManagement) (string, error) {

	nefSelf := pfdTrans.Self

	if nefSelf == "" {
		return "", errors.New("NEF Self Link Not Present")
	}

	res := strings.Split(string(nefSelf), "transactions")
	pID := strings.Split(res[1], "/")

	var afSelf string

	if HTTP2Enabled == true {

		afSelf = "https" + "://" + cfg.SrvCfg.Hostname +
			cfg.SrvCfg.CNCAEndpoint + cfg.LocationPrefixPfd +
			"transactions/" + pID[1]
	} else {

		afSelf = "http" + "://" + cfg.SrvCfg.Hostname +
			cfg.SrvCfg.CNCAEndpoint + cfg.LocationPrefixPfd +
			"transactions/" + pID[1]
	}
	return afSelf, nil
}

func updateAppsLink(cfg Config, r *http.Request,
	pfdTrans PfdManagement) error {
	for key, v := range pfdTrans.PfdDatas {

		appSelf, err := updateAppLink(cfg, r, v)
		if err != nil {
			return err
		}
		v.Self = Link(appSelf)
		pfdTrans.PfdDatas[key] = v
	}
	return nil
}

func updateAppLink(cfg Config, r *http.Request,
	pfdData PfdData) (string, error) {

	self := pfdData.Self
	if self == "" {
		return "", errors.New("NEF App Self Link Not Present")
	}
	res := strings.Split(string(self), "transactions")
	pID := strings.Split(res[1], "/")
	app := strings.Split(string(self), "applications")

	var appSelf string
	if HTTP2Enabled == true {
		appSelf = "https" + "://" + cfg.SrvCfg.Hostname +
			cfg.SrvCfg.CNCAEndpoint + cfg.LocationPrefixPfd +
			"transactions/" + pID[1] + "/applications" + app[1]
	} else {
		appSelf = "http" + "://" + cfg.SrvCfg.Hostname +
			cfg.SrvCfg.CNCAEndpoint + cfg.LocationPrefixPfd +
			"transactions/" + pID[1] + "/applications" + app[1]
	}
	return appSelf, nil

}

func validateCliPACfg(paCfg *GenericCliConfig) (err error) {

	if paCfg == nil {
		err = errors.New("Nil policy auth cli configuration")
		log.Errf("%s", err.Error())
		return err
	}
	return nil
}

func validatePAAuthToken(a *PolicyAuthAPIClient) {

	var err error
	if a.cfg.OAuth2Support {
		a.oAuth2Token, err = getPcfOAuth2Token()
		if err != nil {
			log.Errf("Oauth2 token refresh error")
		}
	}
}

func errRspHeader(w *http.ResponseWriter, method string,
	errString string, statusCode int) {
	log.Errf("Pfd Management %s : %s", method, errString)
	(*w).WriteHeader(statusCode)

}
