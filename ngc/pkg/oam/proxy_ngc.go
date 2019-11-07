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


var NGCType   string  // APISTUB or 5GFLEXCORE
var URLBase   string
var AllRecords []AfRegister
var AllRecordsAfId []string
var NewRecordAfId  int

// Init Proxy
// The proxy acts as reverse proxy to handle request from CNCA and forward it to the target. 
// The target can be API_STUB_TEST or flexcore. NOTE: current version only support API_STUB.
func InitProxy(npcEndpoint string, redirectTarget string, apistub_testdatapath string) error {
    URLBase = "http://" + npcEndpoint
    NGCType = redirectTarget
    if NGCType == "APISTUB" {
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
    }
    
    return nil

}

func ProxyGetAll(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL GetAll: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        ret, _ := json.Marshal(AllRecords)
        if ret != nil {
             w.Header().Set("Content-Type", "application/json; charset=UTF-8")
             w.WriteHeader(http.StatusOK)
             w.Write([]byte(ret))
             return
        }
    }

    log.Printf("GetAll Failed with TargetNGC %s\n", NGCType)
    w.WriteHeader(404)
}

func ProxyAdd(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Add: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
        body, _ := ioutil.ReadAll(r.Body)
        log.Printf("HTTPRequest Body: %s\n", string(body))

        var newRecord []AfRegister
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
    }

    log.Printf("Add Failed with TargetNGC %s\n", NGCType)
    w.WriteHeader(404)
}


func ProxyDel(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Del: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
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
        log.Printf("[APISTUB MODE] AllRecords with num: %d\n", len(AllRecords))
        log.Println(AllRecords)
        log.Println(AllRecordsAfId)
        w.WriteHeader(http.StatusOK)
        return
    }

    w.WriteHeader(404)

}

func ProxyGet(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Get: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
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
        log.Printf("[APISTUB MODE] GetRecord with index: %d\n", j)
        log.Println(AllRecords[j])
   
        // Respons Body.
        //var rspData AfRegister
        jData, err := json.Marshal(AllRecords[j])
        if err != nil {
                 w.WriteHeader(404)
                 log.Println(err)
                 return;
        }
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        w.Write(jData)
        return
    }
    
    w.WriteHeader(404)

}


func ProxyUpdate(w http.ResponseWriter, r *http.Request) {

    log.Printf("URL Update: %s\n", URLBase + r.URL.Path)

    if NGCType == "APISTUB" {
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
        log.Printf("[APISTUB MODE] GetRecord with index: %d\n", j)
        body, _ := ioutil.ReadAll(r.Body)
        log.Printf("HTTPRequest Body: %s\n", string(body))
        
        //insert and delete 
        var newRecord []AfRegister
        if err := json.Unmarshal(body, &newRecord); err == nil {
             AllRecords[j] = newRecord[0]
             // insert  the new record.
             //AllRecords = append(AllRecords, 0)
             //copy(s[i+1:], s[i:])
             // Append the new AFID for the new record
             //NewRecordAfId ++
             //newAfId   := []string { strconv.Itoa(NewRecordAfId) }
             //AllRecordsAfId = append(AllRecordsAfId, newAfId...)
             w.WriteHeader(http.StatusOK)
             return        
        
        } else {
            log.Println(err)
        }
        
    }

    log.Printf("Update Failed with TargetNGC %s\n", NGCType)
    w.WriteHeader(404)

}
