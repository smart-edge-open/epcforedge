// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"net/http"
	"time"
)

func afLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Infoln("")
		log.Infof("%s %s %s %s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
		inner.ServeHTTP(w, r)
	})
}
