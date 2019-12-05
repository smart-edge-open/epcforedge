// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package oam

import (
	"net/http"
)

func add(w http.ResponseWriter, r *http.Request) {
	ProxyAdd(w, r)
}

func delete(w http.ResponseWriter, r *http.Request) {
	ProxyDel(w, r)
}

func get(w http.ResponseWriter, r *http.Request) {
	ProxyGet(w, r)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	ProxyGetAll(w, r)
}

func update(w http.ResponseWriter, r *http.Request) {
	ProxyUpdate(w, r)
}
