/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019-2020 Intel Corporation
 */

package ngcnef

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	oauth2 "github.com/otcshare/epcforedge/ngc/pkg/oauth2"
	"golang.org/x/net/http2"
)

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

	mreqbody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Failed marshaling req body :%s", err)
		return nil, err

	}
	var body *bytes.Buffer

	// Detect reqBody type and post.
	if reqBody != nil {
		reqContentType := headerParams["Content-Type"]

		body, err = setBody(mreqbody, reqContentType)
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

//HTTPclient creates a new HTTP Client
func genHTTPClient(cfg *PcfPolicyAuthorizationConfig) (*http.Client, error) {

	HTTPClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	if cfg.Protocol == "https" {
		CACert, err1 := ioutil.ReadFile(cfg.ClientCert)
		if err1 != nil {
			fmt.Printf("NEF Certification loading Error: %v", err1)
			return nil, err1

		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(CACert)
		var tlsConfig *tls.Config
		if !cfg.VerifyCerts {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}

		} else {
			tlsConfig = &tls.Config{
				RootCAs: caCertPool,
			}

		}
		if cfg.ProtocolVer == "1.1" {
			HTTPClient.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
			return HTTPClient, nil
		} else if cfg.ProtocolVer == "2.0" {

			HTTPClient.Transport = &http2.Transport{
				TLSClientConfig: tlsConfig,
			}
			return HTTPClient, nil
		} else {

			return nil, errors.New("Unsupported HTTP version provided")
		}

	} else if cfg.Protocol == "http" {

		if cfg.ProtocolVer == "1.1" {
			return HTTPClient, nil
		} else if cfg.ProtocolVer == "2.0" {

			HTTPClient.Transport = &http2.Transport{
				AllowHTTP: true,
				DialTLS: func(network, addr string,
					cfg *tls.Config) (
					net.Conn, error) {
					return net.Dial(network, addr)
				},
			}
			return HTTPClient, nil
		} else {

			return nil, errors.New("Unsupported HTTP version provided")
		}
	} else {
		return nil, errors.New("Unsupported HTTP protocol provided")
	}

}
func getPcfOAuth2Token() (token string, err error) {

	token, err = oauth2.GetAccessToken()
	if err == nil {
		log.Infoln("Got Pcf OAuth2 Access Token: " + token)
	}
	return token, err
}
func validatePAAuthToken(a *PcfClient) {

	var err error
	if a.Pcfcfg.OAuth2Support {
		a.OAuth2Token, err = getPcfOAuth2Token()
		if err != nil {
			log.Errf("Oauth2 token refresh error")
		}
	}
}
