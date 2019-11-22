package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"bytes"
	af "github.com/otcshare/epcforedge/ngc/pkg/af/lib"
)


func main(){	

	client := &http.Client{}

	var ReqBody = &af.TrafficInfluSub{
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
		//			      FlowID: nil}},
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

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(ReqBody)

	
	//fmt.Println("Sending GET request")
	fmt.Println("Sending POST request")

	//request, err := http.NewRequest("GET", "http://localhost:8080/CNCA/1.0.1/subscriptions", nil)
	request, err := http.NewRequest("POST", "http://localhost:8080/CNCA/1.0.1/subscriptions", buf)
	if err != nil {
		fmt.Println("Error preparing request")
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	
	if err != nil {
		fmt.Println("Could not send request")
		log.Fatal(err)
	}
	
	defer response.Body.Close()
	fmt.Printf("Response status : %v \n", response.Status)	
}
