package edgeaf

import (
	"log"
	"net/http"
	"time"
)

// Connectivity constants
const (
	AFServerPort = "80"
)

func main() {

	AFRouter := NewAFRouter()
	s := &http.Server{
		Addr:           ":"+AFServerPort,
		Handler:        AFRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("AF listening on", s.Addr)
	log.Fatal(s.ListenAndServe())
}
