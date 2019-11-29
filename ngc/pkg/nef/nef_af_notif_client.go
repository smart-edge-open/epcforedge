// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/* Client implementation of the pcf stub */

package ngcnef

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

/*  The AF Client is an implemenation of the AF Notification.
Refer 3GPP 29500 "5.2.2.2-1: Mandatory to support HTTP request standard
headers" for the list of headers to be sent in the request.
The following HTTP headers are supported in the requests sent and here's
the role of the AF Client for these headers

Accept : application/json
Accept-Encoding : Not supported and not added. Ex: compress, gzip
Content-Length : This would be added by the GO HTTP Stack
Content-Type: application/json
Content-Encoding: Not supported and not added as we do not specify any special
content encodings
User-Agent: Would be read from the configuration file. It's of the format
NEF-xxxx eg: NEF-OPENNESS-1912
Cache-Control: Not supported and not added
If-Module-Since: Not supported and not added
If-Match : Not supported and not added
Via : This is added by proxies and managed by the GO HTTP stack
Authorization: Not supported and not added. Might be needed in future when
OAuth 2.0 would be supported
*/

// AfClient is an implementation of the Af Notification
type AfClient struct {
	af string
}

// NewAfClient creates a new Udr Client
func NewAfClient(cfg *Config) *AfClient {

	c := &AfClient{}
	c.af = "Af Notification Client"
	return c
}

// AfNotificationUpfEvent is an implementation for sending upf event
func (af *AfClient) AfNotificationUpfEvent(ctx context.Context,
	afURI URI, body EventNotification) error {

	var client http.Client

	nefCtx := ctx.Value(nefCtxKey("nefCtx")).(*nefContext)

	log.Infof("AfNotificationUpfEvent uri :%s", afURI)

	/* Check the url type - if its https or http */
	u, err := url.Parse(string(afURI))
	if err != nil {
		log.Errf("AfNotification URl error :%v", err)
		return err
	}

	// If https then load the certificate
	if u.Scheme == "https" {
		CACert, err1 := ioutil.ReadFile(nefCtx.cfg.HTTP2Config.AfServerCert)
		if err1 != nil {
			log.Errf("Af Certification loading Error: %v", err)
			return err1
		}

		CACertPool := x509.NewCertPool()
		CACertPool.AppendCertsFromPEM(CACert)

		client = http.Client{
			Timeout: 15 * time.Second,
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: CACertPool,
				},
			},
		}
	} else if u.Scheme == "http" {
		client = http.Client{Timeout: 15 * time.Second}
	} else {
		log.Errf("Unsupported url scheme: %s", u.Scheme)
		return errors.New("Unsupported url scheme")

	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		log.Err(err)
		return err
	}
	//log.Infof("POST body ==> \n %s", string(requestBody))
	// Set request type as POST
	req, _ := http.NewRequest("POST", string(afURI),
		bytes.NewBuffer(requestBody))
	// Add user-agent header and content-type header
	req.Header.Set("User-Agent", "NEF-OPENNESS-1912")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	log.Info("Sending a request to the server")
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err)
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errf("response body was not closed properly")
		}
	}()

	log.Infof("Headers in the response %d =>", resp.StatusCode)
	for k, v := range resp.Header {
		log.Infof("%q:%q\n", k, v)
	}
	log.Info("Body in the response =>")
	respbody, err := ioutil.ReadAll(resp.Body)
	log.Infof(string(respbody))
	return err
}
