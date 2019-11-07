package oam 

import (
	"net/http"
)

func Add(w http.ResponseWriter, r *http.Request) {
        ProxyAdd(w, r)
}

func Delete(w http.ResponseWriter, r *http.Request) {
        ProxyDel(w, r)
}

func DeleteDns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        //TBD ....
	w.WriteHeader(http.StatusOK)
}

func Get(w http.ResponseWriter, r *http.Request) {
        ProxyGet(w, r)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
        ProxyGetAll(w, r)
}

func Update(w http.ResponseWriter, r *http.Request) {
        ProxyUpdate(w, r)
}
