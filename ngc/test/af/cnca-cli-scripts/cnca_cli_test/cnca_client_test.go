package main_test

import (
	"net/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const(
	host = "localhost:8080"
	basePath = "/AF/v1"
)

var _ = Describe("CncaClient", func() {

	Describe("Cnca sends GET request", func() {
		
		Context("with correct url", func() {
			Specify("will return no error and valid status code", func() {
		
				client := &http.Client{} 
				resp, err := client.Get("http://" + host + basePath + "/subscriptions")
				Expect(err).ShouldNot(HaveOccurred())
				//defer resp.Body.Close()
				Expect(resp.Status).To(Equal("200 OK"))

			})
		})
		
		Context("with incomplete url", func() {

			Specify("will return no error and invalid status code", func(){
				
				client := &http.Client{}
				resp, err := client.Get("http://" + host + "/subscriptions")	
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.Status).To(Equal("404 Not Found"))
			})
		})		
	})
})
