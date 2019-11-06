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

package af

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
)

func genAFTransId(trans AFTransactionIDs) int {
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

	if num == math.MaxInt32 {
		num = min
	}
	//look for a free ID until it is <= math.MaxInt64 is achieved again
	for found && num < math.MaxInt32 {
		//if below the max value, increment it by 1
		num = num + 1
		//check if the ID is in use, if not - return the ID
		if _, found = trans[num]; !found {
			trans[num] = TrafficInfluSub{}
			fmt.Printf("Num :%d", num)
			return num
		} else {
			num = num + 1
		}
	}
	return 0
}

func getSubsIdFromUrl(url *url.URL) (string, error) {

	subsUrl := url.String()
	if url == nil {
		return "", errors.New("empty URL in the request message")
	}
	s := strings.Split(subsUrl, "/")
	return s[len(s)-1], nil
}

func genTransactionId(afCtx *afContext) (int, error) {

	afTransId := genAFTransId(afCtx.transactions)
	if afTransId == 0 {
		return 0, errors.New("the pool of AF Transaction IDs is already used")
	}

	return afTransId, nil

}
