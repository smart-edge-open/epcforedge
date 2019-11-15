package oam 

import (
    "errors"
    "log"
    "net/http"
)


var NGCType   string  // APISTUB or 5GFLEXCORE
var URLBase   string

// Init Proxy
// The proxy acts as reverse proxy to handle request from CNCA and forward it to the target. 
// The target can be API_STUB_TEST or flexcore. NOTE: current version only support API_STUB.
func InitProxy(npcEndpoint string, redirectTarget string, apistub_testdatapath string) error {
    URLBase = "http://" + npcEndpoint
    NGCType = redirectTarget
    if NGCType == "APISTUB" {
        APIStubInit(apistub_testdatapath)
    } else {
       return errors.New("can't not support flexcore")
    }
    
    return nil

}

func ProxyGetAll(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL GetAll: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        APIStubGetAll(w, r)
    } else {
        log.Printf("GetAll Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(404)
    }
}

func ProxyAdd(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Add: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        APIStubAdd(w, r)
    } else {
        log.Printf("Add Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(404)
    }
}


func ProxyDel(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Del: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        APIStubDel(w, r)
    } else {
        log.Printf("Del Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(404)
    }

}

func ProxyDelDnn(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL DelDnn: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        APIStubDelDnn(w, r)
    } else {
        log.Printf("Del Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(404)
    }
}

func ProxyGet(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Get: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        APIStubGet(w, r)
    } else {
        log.Printf("Get Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(404)
    }

}


func ProxyUpdate(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Update: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        APIStubUpdate(w, r)
    } else {
        log.Printf("Update Failed with TargetNGC %s\n", NGCType)
        w.WriteHeader(404)
    }

}
