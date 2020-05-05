// SPDX-License-Identifier: Apache-2.0
// Copyright © 2020 Intel Corporation

package af

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	oauth2 "github.com/otcshare/epcforedge/ngc/pkg/oauth2"
	"golang.org/x/net/http2"
)

var contentTypeJSON string = "application/json"

// Notif correlId in upPatchChgEvent struct
var notifCorreID int32 = 1

// PolicyAuthAPIClient type
// In most cases there should be only one, shared, PolicyAuthAPIClient.
type PolicyAuthAPIClient struct {
	cfg *GenericCliConfig
	/*
	 * Reuse a single struct instead of allocating one for each
	 * policyAuthService on the heap.
	 */
	common policyAuthService

	oAuth2Token       string
	httpClient        *http.Client
	rootURI           string
	rootNotifURI      string
	userAgent         string
	locationPrefixURI string

	// API Services
	PolicyAuthAppSessionAPI *PolicyAuthAppSessionAPIService

	PolicyAuthEventSubsAPI *PolicyAuthEventSubsAPIService

	PolicyAuthIndividualAppSessAPI *PolicyAuthIndividualAppSessAPIService
}

type policyAuthService struct {
	client *PolicyAuthAPIClient
}

// callAPI do the request.
func (c *PolicyAuthAPIClient) callAPI(request *http.Request) (
	*http.Response, error) {

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func getLocationPrefixURI(srvCfg *ServerConfig, cfg *GenericCliConfig) (
	string, error) {

	uri := srvCfg.Hostname + srvCfg.CNCAEndpoint +
		cfg.LocationPrefixURI
	return uri, nil
}

func getPcfOAuth2Token() (token string, err error) {

	token, err = oauth2.GetAccessToken()
	if err == nil {
		log.Infoln("Got Pcf OAuth2 Access Token: " + token)
	}
	return token, err
}

func getRootNotfiURI(notifSrvCfg *notifSrvConfig, cfg *GenericCliConfig) (
	string, error) {

	var (
		uri string
		err error
	)

	if notifSrvCfg == nil {
		err = errors.New("Nil Notification server config")
		log.Errf("%s", err.Error())
		return uri, err
	}

	uri = notifSrvCfg.Protocol + "://" + notifSrvCfg.Hostname +
		":" + notifSrvCfg.Port + cfg.NotifURI
	return uri, nil
}

/*
 * function to initialize different variable specific to Policy Authorization
 * - Initiate Policy auth api client which is reused for connecting to PCF
 * - Initiate Notification URI to be used while sending req to PCF
 */
func initPACfg(afCtx *Context) (err error) {

	paCfg := afCtx.cfg.CliPcfCfg
	err = validateCliPACfg(paCfg)
	if err != nil {
		log.Errf("Policy Auth client configuration invalid")
		return err
	}

	afCtx.cfg.policyAuthAPIClient, err =
		NewPolicyAuthAPIClient(&afCtx.cfg)
	if err != nil {
		log.Errf("Unable to create policy auth api client")
		return err
	}

	paCfg.NotifURI, err = getRootNotfiURI(afCtx.cfg.NotifSrvCfg, paCfg)
	if err != nil {
		return err
	}

	return nil
}

// NewPolicyAuthAPIClient - helper func
/*
 * NewAPIClient creates a new API client. Basically create new http client if
 * not set in client configurations.
 */
func NewPolicyAuthAPIClient(cfg *Config) (*PolicyAuthAPIClient, error) {

	paCfg := cfg.CliPcfCfg
	c := &PolicyAuthAPIClient{}

	httpClient, err := genHTTPClient(paCfg)
	if err != nil {
		log.Errf("Error in generating http client")
		return nil, err
	}
	c.httpClient = httpClient

	c.rootURI = paCfg.Protocol + "://" + paCfg.Hostname + ":" + paCfg.Port +
		paCfg.BasePath
	c.userAgent = cfg.UserAgent

	c.rootNotifURI, err = getRootNotfiURI(cfg.NotifSrvCfg, paCfg)
	if err != nil {
		return nil, err
	}

	c.locationPrefixURI, err = getLocationPrefixURI(&cfg.SrvCfg, paCfg)
	if err != nil {
		return nil, err
	}

	c.cfg = paCfg
	c.common.client = c

	// API Services
	c.PolicyAuthAppSessionAPI = (*PolicyAuthAppSessionAPIService)(&c.common)
	c.PolicyAuthEventSubsAPI = (*PolicyAuthEventSubsAPIService)(&c.common)
	c.PolicyAuthIndividualAppSessAPI =
		(*PolicyAuthIndividualAppSessAPIService)(&c.common)
	if paCfg.OAuth2Support {
		token, err := getPcfOAuth2Token()
		if err != nil {
			log.Errf("Pcf OAuth2 Token retrieval error: " +
				err.Error())
			return nil, err
		}
		c.oAuth2Token = token
	}

	return c, nil
}

// prepareRequest build the request
func (c *PolicyAuthAPIClient) prepareRequest(
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
	httpRequest.Header.Add("User-Agent", c.userAgent)

	if ctx != nil {
		// add context to the request
		httpRequest = httpRequest.WithContext(ctx)

		// Walk through any authentication.
		if c.cfg.OAuth2Support {
			token := c.oAuth2Token
			if token == "" {
				err = errors.New("Nil Ouath2Token in " +
					"PcfApiClient Struct")
				return nil, err
			}
			httpRequest.Header.Add("Authorization", "Bearer "+token)
		}

	}
	return httpRequest, nil
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

func genHTTPClient(cfg *GenericCliConfig) (*http.Client, error) {

	if cfg.Protocol == "https" {
		CACert, err := ioutil.ReadFile(cfg.CliCertPath)
		if err != nil {
			log.Errf("Error: %v", err)
			return nil, err
		}

		CACertPool := x509.NewCertPool()
		CACertPool.AppendCertsFromPEM(CACert)

		var tlsCfg *tls.Config
		if cfg.VerifyCerts {
			tlsCfg = &tls.Config{
				RootCAs: CACertPool,
			}
		} else {
			tlsCfg = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

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
