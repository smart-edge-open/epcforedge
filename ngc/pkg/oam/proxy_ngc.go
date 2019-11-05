package oam 

import (
//    "bytes"
    "encoding/json"
//    "os"
    "log"
//    "fmt"
    "path/filepath"
    "io/ioutil"
//    "net/url"
    "net/http"
//    "net/http/httputil"
//      sw "../swagger"
)


var targetNgc string  // APISTUB or 5GFLEXCORE
var targetUrl string
var testData  AfRegisterList

// Init Proxy
func InitProxy(npcEndpoint string, redirectTarget string, apistub_testdatapath string) error {
    targetUrl = "http://" + npcEndpoint
    targetNgc = redirectTarget
    if targetNgc == "APISTUB" {

        cfgData, err := ioutil.ReadFile(filepath.Clean(apistub_testdatapath))
        if err != nil {
             return err
        }
        err = json.Unmarshal(cfgData, &testData) 
        if err != nil {
             return err
        }
        log.Printf("APISTUB MODE with TestData: \n")
        log.Println(testData)
    }
    
    return nil

}



// Reversy Proxy to retrieve resource from NGC servers
// ToDo: Need to convert 5goam json format to the data model used for ngc
func ProxyGetAll(w http.ResponseWriter, r *http.Request) {

    // Change URL 
    url := targetUrl + r.URL.Path
    log.Printf("GetAll From:  %s\n", url)

    if targetNgc == "APISTUB" {
        ret, _ := json.Marshal(testData)
        if ret != nil {
             w.Header().Set("Content-Type", "application/json; charset=UTF-8")
             w.WriteHeader(http.StatusOK)
             w.Write([]byte(ret))
             return
        }
    }

    log.Printf("GetAll Failed with TargetNGC %s\n", targetNgc)
    w.WriteHeader(404)
}





   /*
    if r.Method == http.MethodGet {
        resp, err := http.Get(url)
        if err != nil {
	   // handle error
        }
        
        defer resp.Body.Close()
        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        log.Printf("%s\n", string(contents))

    } else if r.Method == http.MethodDelete {
    } else if r.Method == http.MethodPost {
    } else if r.Method == http.MethodPatch {
    } else {
        fmt.Printf("This is a " + r.Method + " request\n")
    }    

    // ngc url 
    //url := "http://127.0.0.1:12345/postdata"
    //contentType := "application/json;charset=utf-8"

    // 
    //var httpdata HttpData
    //httpdata.Flag = 1
    //httpdata.Msg = "terrychow"

    //b, err := json.Marshal(w)
    //if err != nil {
    //    fmt.Println("json format error:", err)
    //    return
    //}

    //body := bytes.NewBuffer(b)

    //resp, err := http.Post(url, contentType, body)
    //if err != nil {
    //    fmt.Println("Post failed:", err)
    //    return
    //}

    //defer resp.Body.Close()

    //content, err := ioutil.ReadAll(resp.Body)
    //if err != nil {
    //    fmt.Println("Read failed:", err)
    //    return
    //}

    //fmt.Println("header:", resp.Header)
    //fmt.Println("content:", string(content))

    */
