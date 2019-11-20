package main

import (
	"fmt"
	"log"
	"net/http"
	//"io/ioutil"
	//"strings"
)

//func SendGetRequest(c http.Client, 

func main(){	

	client := &http.Client{}
	fmt.Println("Sending GET request")
	request, err := http.NewRequest("GET", "http://localhost:8080/CNCA/1.0.1/subscriptions", nil)
	if err != nil {
		fmt.Println("Status code 0")
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	
	if err != nil {
		fmt.Println("Could not send request")
		log.Fatal(err)
	}
	
	//body, err := ioutil.ReadAll(response.body)
	
	//if err != nil {
	//	fmt.Printf("Status code 2 %v", resp.StatusCode)
	//	log.Fatal(err)
	//}

	//fmt.Println(string(body))
	fmt.Printf("Response status : %v", response.Status)	
}
