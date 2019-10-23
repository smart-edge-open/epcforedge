package edgeoam

import (
	"log"
	"net/http"
	"time"
)

// Connectivity constants
const (
	OAMServerPort = "80"
)

func main() {

	OAMRouter := NewOAMRouter()
	s := &http.Server{
		Addr:           ":"+OAMServerPort,
		Handler:        OAMRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("OAM listening on", s.Addr)
	log.Fatal(s.ListenAndServe())
}
