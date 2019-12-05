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


package main

import (
        "os"
        "log"
        "net/http"
	"time"
        oam "github.com/otcshare/epcforedge/ngc/pkg/oam"
        config "github.com/otcshare/epcforedge/ngc/pkg/config"
	"github.com/gorilla/handlers"
)

type oamCfg struct {
        TLSEndpoint        string        `json:"TlsEndpoint"`
        OpenEndpoint       string        `json:"OpenEndpoint"`
	UIEndpoint	   string	 `json:"UIEndpoint"`
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
               cfg.UIEndpoint,
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

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOK := handlers.AllowedOrigins([]string{cfg.UIEndpoint})
	methodsOK := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"})

	serverOAM := &http.Server{
		Addr:		cfg.OpenEndpoint,
		Handler:	handlers.CORS(headersOK, originsOK, methodsOK)(router),
		ReadTimeout:	5 * time.Second,
		WriteTimeout:	10 * time.Second,
	}

	log.Printf("Serving OAM on: %s", cfg.OpenEndpoint)
	serverOAM.ListenAndServe()
}
