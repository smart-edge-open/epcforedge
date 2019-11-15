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

var AllRecords []AfService
var AllRecordsAfId []string
var NewRecordAfId  int

func APIStubInit(apistub_testdatapath string) error {
    cfgData, err := ioutil.ReadFile(filepath.Clean(apistub_testdatapath))
    if err != nil {
       return err
    }
    err = json.Unmarshal(cfgData, &AllRecords) 
    if err != nil {
       return err
    }
    log.Printf("[APISTUB MODE] Init AllRecords with num %d: \n", len(AllRecords))
    for _, a := range AllRecords {
       log.Println(a)
       NewRecordAfId++;
       newAfId   := []string { strconv.Itoa(NewRecordAfId) }
       AllRecordsAfId = append(AllRecordsAfId, newAfId...)
    }
    
    return nil

}

func APIStubReset() error {
    AllRecords = nil
    AllRecordsAfId = nil
    NewRecordAfId = 0
    return nil

}

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
    w.WriteHeader(404)
}

func APIStubAdd(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Add: %s\n", r.URL.Path)

    body, _ := ioutil.ReadAll(r.Body)
    log.Printf("HTTPRequest Body: %s\n", string(body))

    var newRecord []AfService
    if err := json.Unmarshal(body, &newRecord); err == nil {
             // Append the new record.
             AllRecords = append(AllRecords, newRecord...)
             
             // Append the new AFID for the new record
             NewRecordAfId ++
             newAfId   := []string { strconv.Itoa(NewRecordAfId) }
             AllRecordsAfId = append(AllRecordsAfId, newAfId...)
             
             // Print all records information
             log.Printf("[APISTUB MODE] NewRecords AfId: %d\n", NewRecordAfId)
             log.Printf("[APISTUB MODE] AllRecords with num: %d\n", len(AllRecords))
             log.Println(AllRecords)
             log.Println(AllRecordsAfId)
             
             // Respons Body.
             var rspData AfId
             rspData.AfId = strconv.Itoa(NewRecordAfId)
             jData, err := json.Marshal(rspData)
             if err != nil {
                 w.WriteHeader(404)
                 log.Println(err)
                 return;
             }
             w.Header().Set("Content-Type", "application/json; charset=UTF-8")
             w.WriteHeader(http.StatusOK)
             w.Write(jData)
             return
    } else {
             log.Println(err)
    }

    log.Printf("Add Failed\n")
    w.WriteHeader(404)
}

func APIStubDel(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Del: %s\n", URLBase + r.URL.Path)

    // get AFID
    vars := mux.Vars(r)
        
    // loop recorded AFID
    j := -1
    found := 0
    for _, a := range AllRecordsAfId {
         j++;
         if a == vars["afId"] {
               found = 1
               break
         }
    }
    if found == 0 {
         log.Printf("Not found in the AllRecords\n")
         w.WriteHeader(404)
         return
    }

    AllRecords = append(AllRecords[:j], AllRecords[j+1:]...)
    AllRecordsAfId = append(AllRecordsAfId[:j], AllRecordsAfId[j+1:]...)
    if len(AllRecordsAfId) == 0 { NewRecordAfId = 0}

    log.Printf("[APISTUB MODE] AllRecords with num: %d\n", len(AllRecords))
    log.Println(AllRecords)
    log.Println(AllRecordsAfId)
    w.WriteHeader(http.StatusOK)

}


func APIStubDelDnn(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Del: %s\n", URLBase + r.URL.Path)

    // get AFID
    vars := mux.Vars(r)
    afId := vars["afId"]
    dnai := vars["dnai"]
    log.Printf("[APISTUB MODE] DelDnn afId %s, dnai %s\n", afId, dnai)

    
    j := -1
    found := 0
    for _, a := range AllRecordsAfId {
         j++;
         if a == afId {
               found = 1
               break
         }
    }
    if found == 0 {
         log.Printf("Not found afId %s in the AllRecords\n", afId)
         w.WriteHeader(404)
         return
    }
    record := AllRecords[j]

    k := -1
    found = 0
    for _, b := range record.LocalServices {
         k++;
         if b.Dnai == dnai {
               found = 1
               break
         }
    }
    if found == 0 {
         log.Printf("Not found dnai %s in the afId %s\n", dnai, afId)
         w.WriteHeader(404)
         return
    }

    AllRecords[j].LocalServices = append(AllRecords[j].LocalServices[:k], AllRecords[j].LocalServices[k+1:]...)
   
    /*
    AllRecords = append(AllRecords[:j], AllRecords[j+1:]...)
    AllRecordsAfId = append(AllRecordsAfId[:j], AllRecordsAfId[j+1:]...)
    if len(AllRecordsAfId) == 0 { NewRecordAfId = 0}

    log.Printf("[APISTUB MODE] AllRecords with num: %d\n", len(AllRecords))
    log.Println(AllRecords)
    log.Println(AllRecordsAfId)
    */
    w.WriteHeader(http.StatusOK)

}



func APIStubGet(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Get: %s\n", r.URL.Path)

    // afId check
    vars := mux.Vars(r)
    j := -1
    found := 0
    for _, a := range AllRecordsAfId {
        j++;
        if a == vars["afId"] {
            found = 1
            break
        }
    }

    if found == 0 {
        log.Printf("Not found in the AllRecords\n")
        w.WriteHeader(404)
        return
    }

    log.Printf("[APISTUB MODE] GetRecord with index: %d\n", j)
    log.Println(AllRecords[j])
   
    // Respons Body.
    jData, err := json.Marshal(AllRecords[j])
    if err != nil {
        w.WriteHeader(404)
        log.Println(err)
        return;
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    w.Write(jData)
}


func APIStubUpdate(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Update: %s\n", r.URL.Path)

    // afId Check
    vars := mux.Vars(r)
    j := -1
    found := 0
    for _, a := range AllRecordsAfId {
        j++;
        if a == vars["afId"] {
             found = 1
             break
        }
    }
    if found == 0 {
        log.Printf("Not found in the AllRecords\n")
        w.WriteHeader(404)
        return
    }
    log.Printf("[APISTUB MODE] GetRecord with index: %d\n", j)
    body, _ := ioutil.ReadAll(r.Body)
    log.Printf("HTTPRequest Body: %s\n", string(body))
        
    //insert and delete 
    var newRecord []AfService
    if err := json.Unmarshal(body, &newRecord); err == nil {
        AllRecords[j] = newRecord[0]
        w.WriteHeader(http.StatusOK)
        return        
        
    } else {
        log.Println(err)
    }
    
    log.Printf("Update Failed")
    w.WriteHeader(404)

}
