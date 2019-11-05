/*
 */

package main

import (
	"os"
	"log"
	"net/http"
	oam "github.com/otcshare/epcforedge/ngc/pkg/oam"
)

type oamCfg struct {
        TLSEndpoint        string        `json:"TlsEndpoint"`
        OpenEndpoint       string        `json:"OpenEndpoint"`
        NgcEndpoint        string        `json:"NgcEndpoint"`
        NgcTarget          string        `json:"NgcTarget"`
        APIStubPath        string        `json:"APIStubPath"`
}

func main() {
        log.SetPrefix("[5goam]")
        log.SetFlags(log.Ldate|log.Ltime |log.LUTC |log.Lshortfile)
	log.Printf("Server started")
        
        var cfg oamCfg
        err := oam.LoadJSONConfig("./oam.json", &cfg)
        if err != nil {
                log.Printf("Failed to load config: %#v", err)
                os.Exit(1)
        }
        log.Printf("Configuration data: %s, %s, %s, %s\n", 
               cfg.TLSEndpoint, 
               cfg.OpenEndpoint, 
               cfg.NgcEndpoint, 
               cfg.APIStubPath)

        // New Http Router
        err = oam.InitProxy(cfg.NgcEndpoint, cfg.NgcTarget, cfg.APIStubPath)
        if err != nil {
                log.Printf("Failed to init proxy: %#v", err)
                os.Exit(1)
        }

	router := oam.NewRouter()
	http.ListenAndServe(cfg.OpenEndpoint, router)
	//log.Fatal(http.ListenAndServe(":8080", router))
	//log.Fatal(http.ListenAndServe(cfg.OpenEndpoint, router))
}
