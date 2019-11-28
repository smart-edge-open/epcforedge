//Copyright 2019 Intel Corporation. All rights reserved.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/otcshare/epcforedge/ngc/pkg/af"
)

const (
	host     = "localhost:8080"
	basePath = "/af/v1"
)

var _ = Describe("AF", func() {

	Describe("Create New Subscription", func() {
		Context("create 1 subscription",
			func() {
				Specify("Will return no error and valid status code", func() {
					client := &http.Client{}
					reqBody, err := ioutil.ReadFile("./testdata/100_AF_NB_SUB_POST001.json")
					Expect(err).ShouldNot(HaveOccurred())
					reqBodyBytes := bytes.NewReader(reqBody)
					resp, err := client.Post("http://"+host+basePath+"/subscriptions", "application/json", reqBodyBytes)
					Expect(err).ShouldNot(HaveOccurred())
					defer resp.Body.Close()
					Expect(resp.Status).To(Equal("201 Created"))
				})
			})
	})

	/*Describe("Cnca sends POST request", func(){

	                Context("with correct parameters in json body", func() {
	                        Specify("will return no error and valid status code", func() {

					client := &http.Client{}
					resp, err := client.Post("http://" + host + basePath + "/subscriptions", "application/json", buffer)
					Expect(err).ShouldNot(HaveOccurred())
					defer resp.Body.Close()
					Expect(resp.Status).To(Equal("201 Created"))
				})
	                })
			Context("With no parameters in request body", func() {
				Specify("will return no error and invalid status code", func() {
					client := &http.Client{}
					resp, err := client.Post("http://"+host+basePath+"/subscriptions", "application/json", nil)
					Expect(err).ShouldNot(HaveOccurred())
					defer resp.Body.Close()
					Expect(resp.Status).To(Equal("500 Internal Server Error"))
				})
			})
	        })


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

		})*/

})
