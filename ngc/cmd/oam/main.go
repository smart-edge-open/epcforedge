/*
 */

package main

import (
	"os"
	"log"
	"net/http"
	oam "github.com/otcshare/epcforedge/ngc/pkg/oam"
        config "github.com/otcshare/epcforedge/ngc/pkg/config"
)

type oamCfg struct {
        TLSEndpoint        string        `json:"TlsEndpoint"`
        OpenEndpoint       string        `json:"OpenEndpoint"`
        NgcEndpoint        string        `json:"NgcEndpoint"`
        NgcType            string        `json:"NgcType"`
        NgcTestData        string        `json:"NgcTestData"`
}

func main() {
        log.SetPrefix("[5goam]")
        log.SetFlags(log.Ldate|log.Ltime |log.LUTC |log.Lshortfile)
	log.Printf("Server started")
        
        var cfg oamCfg
        err := config.LoadJSONConfig("./configs/oam.json", &cfg)
        if err != nil {
                log.Printf("Failed to load config: %#v", err)
                os.Exit(1)
        }
        log.Printf("LocalConfig: %s, %s, %s, %s, %s\n", 
               cfg.TLSEndpoint, 
               cfg.OpenEndpoint, 
               cfg.NgcEndpoint, 
               cfg.NgcType, 
               cfg.NgcTestData)

        // New Http Router
        err = oam.InitProxy(cfg.NgcEndpoint, cfg.NgcType, cfg.NgcTestData)
        if err != nil {
                log.Printf("Failed to init proxy: %#v", err)
                os.Exit(1)
        }

	router := oam.NewRouter()
	http.ListenAndServe(cfg.OpenEndpoint, router)
}
