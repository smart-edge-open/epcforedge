package main_test

import (
	"net/http"
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	af "github.com/otcshare/epcforedge/ngc/pkg/af/lib"	
)

const(
	host = "localhost:8080"
	basePath = "/CNCA/1.0.1"
)

var _ = Describe("CncaClient", func() {

        var reqBody = &af.TrafficInfluSub{
                AfAppID: "app001",
                ExternalGroupID: "",
                AfServiceID: "",
                AfTransID: "",
                Dnn: "edgeLocation001",
                Snssai: af.Snssai{Sst: 0,
                         Sd: "default"},
                SubscribedEvents: []af.SubscribedEvent{""},
                Gpsi: "",
                Ipv4Addr: "",
                Ipv6Addr: "",
                MacAddr: "",
                DnaiChgType: "",
                NotificationDestination: "",
                Self: "",
                //TrafficFilters: []af.FlowInfo{{FlowDescriptions: nil,
                //                            FlowID: nil}},
                //EthTrafficFilters: []af.EthFlowDescription{{""}},
                TrafficRoutes: []af.RouteToLocation{{Dnai: "",
                                                RouteProfID: "",
                                RouteInfo: af.RouteInformation{Ipv4Addr: "", Ipv6Addr: "", PortNumber: 0}}},
                TempValidities: []af.TemporalValidity{{StartTime: "",
                                                      StopTime: "",}},
                ValidGeoZoneIDs: nil,
                SuppFeat: "",
                RequestTestNotification: true,
                WebsockNotifConfig: af.WebsockNotifConfig{WebsocketURI: "",
                                     RequestWebsocketURI: true},
                AppReloInd: false,
                AnyUeInd: true,
                }


	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(reqBody)


	Describe("Cnca sends POST request", func(){
		
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
		
	})
	
})
