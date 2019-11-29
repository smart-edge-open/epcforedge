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

package oam 

import (
    "encoding/json"
    "log"
    "strconv"
    "path/filepath"
    "io/ioutil"
    "net/http"
    "github.com/gorilla/mux"

)

// AllRecords store all the AFService
var AllRecords             []AFService

// NewRecordAFServiceID is id to allocate
var NewRecordAFServiceID   int

// AFServiceIDBaseValue is base value for the id to allocate
const AFServiceIDBaseValue = 123456

// APIStubInit : stub init
func APIStubInit(apistubTestdatapath string) error {
    // Init value     
    NewRecordAFServiceID = AFServiceIDBaseValue // BaseValue for new AFServiceID
    
    // Read records from test stub file
    cfgData, err := ioutil.ReadFile(filepath.Clean(apistubTestdatapath))
    if err != nil {
       return err
    }
    err = json.Unmarshal(cfgData, &AllRecords) 
    if err != nil {
       return err
    }
    log.Printf("[APISTUB MODE] Init with num %d: \n", len(AllRecords))
    for _, a := range AllRecords {
       log.Println(a)
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
    log.Printf("[APISTUB MODE] NewAFServiceID: %d\n", NewRecordAFServiceID)
    log.Printf("[APISTUB MODE] AllRecords num is: %d\n", len(AllRecords))
    log.Println(AllRecords)
}

// APIStubNewAFServiceID : allocate new service id
func APIStubNewAFServiceID() string {
     NewRecordAFServiceID++
     return   strconv.Itoa(NewRecordAFServiceID) 
}

// APIStubGetRecordIndex : get record by service id
func APIStubGetRecordIndex(serviceID string) int {

    log.Printf("[APISTUB MODE]  Searching: %s\n", serviceID)
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

    log.Printf("URL GetAll: %s\n", r.URL.Path)
    log.Printf("Number of All Records is: %d", len(AllRecords))
    ret, _ := json.Marshal(AllRecords)
    if ret != nil {
         w.Header().Set("Content-Type", "application/json; charset=UTF-8")
         w.WriteHeader(http.StatusOK)
         w.Write([]byte(ret))
         return
    }

    log.Printf("GetAll Failed")
    w.WriteHeader(http.StatusNotFound)
}

// APIStubAdd : add one to records
func APIStubAdd(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Add: %s\n", r.URL.Path)
    body, _ := ioutil.ReadAll(r.Body)
    log.Printf("HTTPRequest Body: %s\n", string(body))

    //var httpBody     LocationService
    // create and append the new record.
    newRecord := make([]AFService,1)
    //newRecord[0].LocationService  = httpBody
    err := json.Unmarshal(body, &(newRecord[0].LocationService))
    if  err != nil {
        log.Println(err)
        log.Printf("Add Failed\n")
        w.WriteHeader(http.StatusNotFound)
        return
    }

    newRecord[0].AFServiceID      = APIStubNewAFServiceID()
    AllRecords = append(AllRecords, newRecord...)
    APIStubPrintAll()
         
    // Respons Body.
    var rspData AFServiceID
    rspData.AFServiceID = newRecord[0].AFServiceID 
    jData, err := json.Marshal(rspData)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        log.Println(err)
        return;
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    w.Write(jData)
}

// APIStubDel : delete one from the records
func APIStubDel(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Del: %s\n", URLBase + r.URL.Path)

    // get AFID
    vars := mux.Vars(r)
        
    // get recorded AFService
    j := APIStubGetRecordIndex(vars["afServiceId"])
    if j == -1 {
       log.Printf("Not found in the AllRecords\n")
       w.WriteHeader(http.StatusNotFound) 
       return
    }

    AllRecords = append(AllRecords[:j], AllRecords[j+1:]...)
    if len(AllRecords) == 0 { 
       APIStubReset()
    }
    APIStubPrintAll()
    w.WriteHeader(http.StatusOK)
}

// APIStubGet : get one from the records
func APIStubGet(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Get: %s\n", r.URL.Path)

    // afId check
    vars := mux.Vars(r)
    // get recorded AFService
    j := APIStubGetRecordIndex(vars["afServiceId"])
    if j == -1 {
       log.Printf("Not found in the AllRecords\n")
       w.WriteHeader(http.StatusNotFound) 
       return
    }


    log.Printf("[APISTUB MODE] GetRecord with index: %d\n", j)
    log.Println(AllRecords[j])

    // Respons Body.
    rspBody := AllRecords[j].LocationService
    jData, err := json.Marshal(rspBody)
    if err != nil {
       w.WriteHeader(http.StatusNotFound)
       log.Println(err)
       return
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    w.Write(jData)
}

// APIStubUpdate : update one from records
func APIStubUpdate(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Update: %s\n", r.URL.Path)

    // afId Check
    vars := mux.Vars(r)
    // get recorded AFService
    j := APIStubGetRecordIndex(vars["afServiceId"])
    if j == -1 {
       log.Printf("Not found in the AllRecords\n")
       w.WriteHeader(http.StatusNotFound) 
       return
    }
    
    body, _ := ioutil.ReadAll(r.Body)
    log.Printf("[APISTUB MODE] GetRecord with index: %d\n", j)
    log.Printf("HTTPRequest Body: %s\n", string(body))

    if err := json.Unmarshal(body, &(AllRecords[j].LocationService)); err == nil {
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

    log.Printf("Update Failed")
    w.WriteHeader(http.StatusNotFound)

}
