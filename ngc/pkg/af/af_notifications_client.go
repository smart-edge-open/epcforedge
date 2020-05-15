/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2019 Intel Corporation
 */

/* Client implementation of the pcf stub */

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

// NotificationAPIClient type
// In most cases there should be only one, shared, NotificationAPIClient.
type NotificationAPIClient struct {
	cfg        *GenericCliConfig
	httpClient *http.Client
	userAgent  string
}

type notifyClientAPI interface {
	NotificationUpPathEvent(notificationURI NotificationURI,
		body NotificationUpPathChg) error
}

// NotificationUpPathEvent is an implementation for sending UP_PATH_CHANGE
// event
func (c *NotificationAPIClient) NotificationUpPathEvent(notifURI NotificationURI,
	body NotificationUpPathChg) error {

	var (
		httpMethod = http.MethodPost
		postBody   interface{}
		afEvent    Afnotification
	)
	// create path and map variables
	path := string(notifURI)

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = contentTypeJSON
	headerParams["Accept"] = contentTypeJSON

	afEvent.Event = UPPathChangeEvent

	payload, err := json.Marshal(body)
	if err != nil {
		log.Err(err)
		return err
	}
	afEvent.Payload = payload

	postBody = &afEvent

	r, err := c.cfg.prepareRequest(cliCtx, path, httpMethod, postBody,
		headerParams)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil || resp == nil {
		return err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errf("response body was not closed properly")
		}
	}()

	log.Infof("Status in the response %d =>", resp.StatusCode)

	return err
}

// NewAFNotifyAPIClient To generate a new Notification Client
func NewAFNotifyAPIClient(afCtx *Context) (*NotificationAPIClient, error) {

	cfg := &GenericCliConfig{
		Protocol:      "http",
		ProtocolVer:   "2.0",
		OAuth2Support: false,
		VerifyCerts:   false,
	}

	c := &NotificationAPIClient{}

	httpClient, err := genHTTPClient(cfg)
	if err != nil {
		log.Errf("Error in generating http client")
		return nil, err
	}
	c.httpClient = httpClient

	c.userAgent = afCtx.cfg.UserAgent
	c.cfg = cfg

	return c, nil
}
