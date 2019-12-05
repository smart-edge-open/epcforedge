// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package oam

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
)

// AllRecords store all the AFService
var AllRecords []AFService

// NewRecordAFServiceID is id to allocate
var NewRecordAFServiceID int

// AFServiceIDBaseValue is base value for the id to allocate
const AFServiceIDBaseValue = 123456

// APIStubInit : stub init
func APIStubInit(apistubTestdatapath string) error {
	// Init value
	NewRecordAFServiceID = AFServiceIDBaseValue // BaseValue for new AFServiceID
	if 0 == len(apistubTestdatapath) {
		return nil 
	}
	// Read records from test stub file
	cfgData, err := ioutil.ReadFile(filepath.Clean(apistubTestdatapath))
	if err != nil {
		return err
	}
	err = json.Unmarshal(cfgData, &AllRecords)
	if err != nil {
		return err
	}
	log.Infof("[APISTUB MODE] Init with num %d: \n", len(AllRecords))
	for _, a := range AllRecords {
		// ignore serviceID in the test, allocate new serviceID
		a.AFServiceID = APIStubNewAFServiceID()
	}

	return nil
}

// APIStubReset : stub reset
func APIStubReset() error {

	AllRecords = nil
	NewRecordAFServiceID = AFServiceIDBaseValue // BaseValue for new AFServiceID
	return nil

}

// APIStubPrintAll : stub print all
func APIStubPrintAll() {
	// Print all records
	log.Infof("[APISTUB MODE] NewAFServiceID: %d\n", NewRecordAFServiceID)
	log.Infof("[APISTUB MODE] AllRecords num is: %d\n", len(AllRecords))
}

// APIStubNewAFServiceID : allocate new service id
func APIStubNewAFServiceID() string {
	NewRecordAFServiceID++
	return strconv.Itoa(NewRecordAFServiceID)
}

// APIStubGetRecordIndex : get record by service id
func APIStubGetRecordIndex(serviceID string) int {

	log.Infof("[APISTUB MODE]  Searching: %s\n", serviceID)
	// loop recorded AFID
	for i, a := range AllRecords {
		if a.AFServiceID == serviceID {
			return i
		}
	}
	return -1
}

// APIStubGetAll : get all records from stub
func APIStubGetAll(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL GetAll: %s\n", r.URL.Path)
	log.Infof("Number of All Records is: %d", len(AllRecords))
	ret, _ := json.Marshal(AllRecords)
	if ret != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(ret))
		return
	}

	log.Errf("GetAll Failed")
	w.WriteHeader(http.StatusNotFound)
}

// APIStubAdd : add one to records
func APIStubAdd(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Add: %s\n", r.URL.Path)
	body, _ := ioutil.ReadAll(r.Body)
	log.Infof("HTTPRequest Body: %s\n", string(body))

	//var httpBody     LocationService
	// create and append the new record.
	newRecord := make([]AFService, 1)
	//newRecord[0].LocationService  = httpBody
	err := json.Unmarshal(body, &(newRecord[0].LocationService))
	if err != nil {
		log.Errf("Add Failed: %s\n", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newRecord[0].AFServiceID = APIStubNewAFServiceID()
	AllRecords = append(AllRecords, newRecord...)
	APIStubPrintAll()

	// Respons Body.
	var rspData AFServiceID
	rspData.AFServiceID = newRecord[0].AFServiceID
	jData, err := json.Marshal(rspData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Errf("Add Failed: %s\n", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jData)
}

// APIStubDel : delete one from the records
func APIStubDel(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Del: %s\n", URLBase+r.URL.Path)

	// get AFID
	vars := mux.Vars(r)

	// get recorded AFService
	j := APIStubGetRecordIndex(vars["afServiceId"])
	if j == -1 {
		log.Errf("Not found in the AllRecords\n")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	AllRecords = append(AllRecords[:j], AllRecords[j+1:]...)
	if len(AllRecords) == 0 {
		_ = APIStubReset()
	}
	APIStubPrintAll()
	w.WriteHeader(http.StatusOK)
}

// APIStubGet : get one from the records
func APIStubGet(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Get: %s\n", r.URL.Path)

	// afId check
	vars := mux.Vars(r)
	// get recorded AFService
	j := APIStubGetRecordIndex(vars["afServiceId"])
	if j == -1 {
		log.Errf("Not found in the AllRecords\n")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Infof("[APISTUB MODE] GetRecord with index: %d\n", j)

	// Respons Body.
	rspBody := AllRecords[j].LocationService
	jData, err := json.Marshal(rspBody)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Errf("err: %s\n", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jData)
}

// APIStubUpdate : update one from records
func APIStubUpdate(w http.ResponseWriter, r *http.Request) {

	log.Infof("URL Update: %s\n", r.URL.Path)

	// afId Check
	vars := mux.Vars(r)
	// get recorded AFService
	j := APIStubGetRecordIndex(vars["afServiceId"])
	if j == -1 {
		log.Errf("Not found in the AllRecords\n")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	log.Infof("[APISTUB MODE] GetRecord with index: %d\n", j)
	log.Infof("HTTPRequest Body: %s\n", string(body))

	if err := json.Unmarshal(body, 
		&(AllRecords[j].LocationService)); err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	//insert and delete
	//var newRecord []AfService
	//if err := json.Unmarshal(body, &newRecord); err == nil {
	//     AllRecords[j] = newRecord[0]
	//     w.WriteHeader(http.StatusOK)
	//     return
	//}

	log.Errf("Update Failed")
	w.WriteHeader(http.StatusNotFound)

}
