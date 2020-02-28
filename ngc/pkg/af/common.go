// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
)

//TransIDMax var
var TransIDMax = math.MaxInt32

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
		err := json.Unmarshal(body, &v)
		if err != nil {
			newErr.error = err.Error()
			return newErr
		}
		newErr.model = v
		if r.StatusCode == 401 {
			if fetchNEFAuthorizationToken() != nil {
				log.Infoln("Token refresh failed")
			}
		}

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
		err := json.Unmarshal(body, &v)
		if err != nil {
			newErr.error = err.Error()
			return newErr
		}
		newErr.model = v
		log.Errf("NEF returned error - %s", r.Status)
		if r.StatusCode == 401 {
			if fetchNEFAuthorizationToken() != nil {
				log.Infoln("Token refresh failed")
			}
		}

		return newErr

	default:
		b, _ := ioutil.ReadAll(r.Body)
		err := fmt.Errorf("NEF returned error - %s, %s", r.Status, string(b))
		return err
	}
}
