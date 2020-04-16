package ngcnef

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// TestPcf boolean to be set to true in PCF UT
var TestPcf bool = false

// HTTPClient to be setup in PCF UT
var HTTPClient *http.Client

// SetHTTPClient Function to setup a httpClient for testing if required
func SetHTTPClient(httpClient *http.Client) {

	HTTPClient = httpClient

}

func sendreq(ctx context.Context, client *http.Client, method string, apiURL string, asc AppSessionContext, ascup AppSessionContextUpdateData, appSessionID AppSessionID) (AppSessionID, AppSessionContext, int, error, ProblemDetails) {

	switch method {
	case "get":
		sessid := string(appSessionID)
		log.Infof("Triggering PCF /* GET */ : " + apiURL + sessid)
		res, err := client.Get(apiURL + sessid)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)

		}
		log.Infof("Body in the response =>")
		body, _ := ioutil.ReadAll(res.Body)
		log.Infof(string(body))
		defer res.Body.Close()
		fmt.Println("pcfclient get")
		fmt.Println(res.StatusCode)
		if res.StatusCode == 200 {
			var appSessionContext AppSessionContext = AppSessionContext{}
			e := json.Unmarshal(body, &appSessionContext)
			if e != nil {
				fmt.Printf("Failed go error :%s", e)
			}
			return appSessionID, appSessionContext, res.StatusCode, nil, ProblemDetails{}
		} else {
			var problemDetails ProblemDetails = ProblemDetails{}
			e := json.Unmarshal(body, &problemDetails)
			if e != nil {
				fmt.Printf("Failed go error :%s", e)
			}
			return appSessionID, AppSessionContext{}, res.StatusCode, nil, problemDetails
		}

	case "post":
		var appsesid string
		var req *http.Request
		var err error

		log.Infof("Triggering PCF /* POST */ :" + apiURL)
		//res, err := client.Post(apiURL, "application/json", bytes.NewBuffer(b))

		log.Infof("post triggered")
		b, err := json.Marshal(asc)
		if err != nil {
			fmt.Println(err)
		}
		req, err = http.NewRequest("POST", apiURL, bytes.NewBuffer(b))

		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed go error :%s", err)
		}

		appsesid = res.Header.Get("Location")

		log.Infof("appsessionid" + appsesid)
		log.Infof("Body in the response =>")
		body, _ := ioutil.ReadAll(res.Body)
		log.Infof(string(body))

		appSessionID = AppSessionID(appsesid)
		defer res.Body.Close()
		fmt.Println("pcfclient post")
		fmt.Println(res.StatusCode)
		if res.StatusCode == 201 {
			var appSessionContext AppSessionContext = AppSessionContext{}
			e := json.Unmarshal(body, &appSessionContext)
			if e != nil {
				fmt.Printf("Failed go error :%s", e)
			}
			return appSessionID, appSessionContext, res.StatusCode, nil, ProblemDetails{}
		} else {
			var problemDetails ProblemDetails = ProblemDetails{}

			e := json.Unmarshal(body, &problemDetails)
			if e != nil {
				fmt.Printf("Failed go error :%s", e)
			}
			return appSessionID, AppSessionContext{}, res.StatusCode, nil, problemDetails
		}

	case "patch":
		sessid := string(appSessionID)
		b, err := json.Marshal(ascup)
		if err != nil {
			fmt.Println("parsing")
			fmt.Println(err)

		}
		fmt.Println(sessid)
		log.Infof("Triggering PCF /* PATCH */ :" + apiURL + sessid)
		req, err := http.NewRequest("PATCH", apiURL+sessid, bytes.NewBuffer(b))
		if err != nil {
			fmt.Println("req")
			log.Infof("Failed go error :%s", err)
		}
		res, err := client.Do(req)

		if err != nil {
			fmt.Println("res")
			fmt.Printf("Failed go error :%s", err)

		}

		log.Infof("Body in the response =>")
		body, _ := ioutil.ReadAll(res.Body)
		log.Infof(string(body))

		appSessionID = AppSessionID(sessid)
		defer res.Body.Close()
		fmt.Println("pcfclient patch")
		fmt.Println(res.StatusCode)
		if res.StatusCode == 200 {
			var appSessionContext AppSessionContext = AppSessionContext{}
			e := json.Unmarshal(body, &appSessionContext)
			if e != nil {
				fmt.Printf("Failed go error :%s", e)
			}
			return appSessionID, appSessionContext, res.StatusCode, nil, ProblemDetails{}
		} else {
			var problemDetails ProblemDetails = ProblemDetails{}
			e := json.Unmarshal(body, &problemDetails)
			if e != nil {
				fmt.Printf("Failed go error :%s", e)
			}
			return appSessionID, AppSessionContext{}, res.StatusCode, nil, problemDetails
		}

	case "delete":
		log.Infof("Triggering PCF /* DELETE */ : " + apiURL)
		req, err := http.NewRequest("POST", apiURL, nil)
		if err != nil {
			fmt.Println("Failed go error :", err)
		}
		res, err := client.Do(req)

		if err != nil {
			fmt.Println("Failed go error :", err)
		}

		fmt.Println("pcfclient delete")
		fmt.Println(res.StatusCode)
		if res.StatusCode == 204 {

			return appSessionID, AppSessionContext{}, res.StatusCode, nil, ProblemDetails{}
		}
		/*  else if res.StatusCode == 200 {
			var eventnoti EventsNotification
			log.Infof("Body in the response =>")
			body, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			log.Infof(string(body))
			e := json.Unmarshal(body, &eventnoti)
			if e != nil {
				fmt.Printf("Failed go error :%s", e)
			}
			return appSessionID, AppSessionContext{}, res.StatusCode, nil, ProblemDetails{}, eventnoti
		} */
		return appSessionID, AppSessionContext{}, res.StatusCode, nil, ProblemDetails{}

	}
	return appSessionID, AppSessionContext{}, 0, nil, ProblemDetails{}
}

// PcfClient is an implementation of the Pcf Authorization
type PcfClient struct {
	pcf string

	// database to store the contents of the app session contexts created
	paDb       map[string]AppSessionContext
	httpClient *http.Client
	pcfRootURI string
	pcfURI     string
}

//NewPCFClientF creates a new PCF Client
func NewPCFClientF(cfg *Config) *PcfClient {

	c := &PcfClient{}
	c.pcf = "PCF freegc"

	c.paDb = make(map[string]AppSessionContext)
	if TestPcf {

		c.httpClient = HTTPClient

	} else {
		c.httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}

		CACert, err1 := ioutil.ReadFile(cfg.PcfPolicyAuthorizationConfig.ClientCert)
		if err1 != nil {
			fmt.Printf("NEF Certification loading Error: %v", err1)
			return nil

		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(CACert)

		tlsConfig := &tls.Config{
			//RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}

		c.httpClient.Transport = &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	c.pcfRootURI = cfg.PcfPolicyAuthorizationConfig.Scheme + "://" + cfg.PcfPolicyAuthorizationConfig.APIRoot
	c.pcfURI = cfg.PcfPolicyAuthorizationConfig.URI

	log.Infof("PCF Client created")
	return c
}

// PolicyAuthorizationCreate is a stub implementation
// Successful response : 201 and body contains AppSessionContext
func (pcf *PcfClient) PolicyAuthorizationCreate(ctx context.Context,
	body AppSessionContext) (AppSessionID, PcfPolicyResponse, error) {

	log.Infof("PCFs PolicyAuthorizationCreate Entered")
	_ = ctx
	if TestPcf {

		pcf.httpClient = HTTPClient

	}
	var err error
	pcfPr := PcfPolicyResponse{}

	Asc := body

	appSessionID, asc, code, err, pd := sendreq(ctx, pcf.httpClient, "post", pcf.pcfRootURI+pcf.pcfURI, Asc, AppSessionContextUpdateData{}, AppSessionID(""))

	if code != 201 {
		log.Infof("PCFs PolicyAuthorizationCreate failed ")
		pcfPr.ResponseCode = uint16(code)
		pcfPr.Pd = &pd
		err = errors.New("failed")

		//pcfPr.Asc=nil
	} else {
		sessid := string(appSessionID)
		//appSessionID := AppSessionID(sessid)

		pcf.paDb[sessid] = asc
		pcfPr.ResponseCode = uint16(code)
		pcfPr.Asc = &asc
		pcfPr.Pd = nil

		log.Infof("PCFs PolicyAuthorizationCreate Exited successfully with "+
			"sessid: %s", sessid)
	}
	/* log.Infof("response sent:")
	b1, err1 := json.Marshal(pcfPr)
	if err1 != nil {
		fmt.Println(err)
	}
	log.Infof(string(b1)) */
	return appSessionID, pcfPr, err
}

// PolicyAuthorizationUpdate is a stub implementation
// Successful response : 200 and body contains AppSessionContext
func (pcf *PcfClient) PolicyAuthorizationUpdate(ctx context.Context,
	body AppSessionContextUpdateData,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PCFs PolicyAuthorizationUpdate Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx
	if TestPcf {

		pcf.httpClient = HTTPClient

	}
	var err error
	pcfPr := PcfPolicyResponse{}

	appSessionID, Asc, code, err, pd := sendreq(ctx, pcf.httpClient, "patch", pcf.pcfRootURI+pcf.pcfURI, AppSessionContext{}, body, appSessionID)

	sessid := string(appSessionID)
	//fmt.Println(code)

	_, prs := pcf.paDb[string(sessid)]
	// if not found return an error i.e 404
	//if !prs && code != 200 {
	if code != 200 {
		log.Infof("PCFs PolicyAuthorizationUpdate AppSessionID %s not found",
			string(appSessionID))
		pcfPr.ResponseCode = uint16(code)
		pcfPr.Pd = &pd
		err = errors.New("failed")

	} else {
		log.Infof("PCFs PolicyAuthorizationUpdate AppSessionID %s updated",
			string(appSessionID))
		if prs {
			pcf.paDb[string(sessid)] = Asc
			pcfPr.ResponseCode = uint16(code)
			pcfPr.Asc = &Asc
		}

	}
	log.Infof("PCFs PolicyAuthorizationUpdate Exited for AppSessionID %s",
		string(appSessionID))
	return pcfPr, err
}

// PolicyAuthorizationDelete is a stub implementation
// Successful response : 204 and empty body
func (pcf *PcfClient) PolicyAuthorizationDelete(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {

	log.Infof("PCFs PolicyAuthorizationDelete Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx
	if TestPcf {

		pcf.httpClient = HTTPClient

	}
	pcfPr := PcfPolicyResponse{}
	sessid := string(appSessionID)
	_, prs := pcf.paDb[sessid]
	var err1 error
	if prs {
		_, _, code, err, _ := sendreq(ctx, pcf.httpClient, "delete", pcf.pcfRootURI+pcf.pcfURI+string(appSessionID)+"/delete", AppSessionContext{}, AppSessionContextUpdateData{}, appSessionID)

		//fmt.Println(code)
		err1 = err
		// convert the appsession id to integer
		//sessid, _ := strconv.Atoi(string(appSessionID))

		// check for the presence of the sessid in the database

		// if not found return an error i.e 404
		//if !prs && (code != 200 && code != 204)
		if code != 200 && code != 204 {
			log.Infof("PCFs PolicyAuthorizationDelete AppSessionID %s not found",
				string(appSessionID))
			err1 = errors.New("failed")
			pcfPr.Pd = &ProblemDetails{}
			pcfPr.Pd.Cause = "wrong appsession id"
			pcfPr.Pd.Status = 404
			pcfPr.ResponseCode = 404
		} else {
			log.Infof("PCFs PolicyAuthorizationDelete AppSessionID %s found",
				string(appSessionID))
			if code == 200 {
				/* log.Infof("response received")
				fmt.Println(even)
				log.Infof("delete response received from pcf:")
				b1, err1 := json.Marshal(even)
				if err1 != nil {
					fmt.Println(err)
				}
				ioutil.WriteFile("appcontextdeleteres.json", b1, 0644) */
			}

			//log.Infof("PCFs Policy Authorization DB size : %d", len(pcf.paDb))
			if prs {
				delete(pcf.paDb, sessid)
				log.Infof("PCFs Policy Authorization DB size : %d", len(pcf.paDb))
			}

			pcfPr.ResponseCode = uint16(code)

		}
		log.Infof("PCFs PolicyAuthorizationDelete Stub Exited for AppSessionID %s",
			string(appSessionID))
	} else {
		err1 = errors.New("failed")
		pcfPr.Pd = &ProblemDetails{}
		pcfPr.Pd.Cause = "wrong appsession id"
		pcfPr.Pd.Status = 404
		pcfPr.ResponseCode = 404
	}

	return pcfPr, err1
}

// PolicyAuthorizationGet is a stub implementation
// Successful response : 204 and empty body
func (pcf *PcfClient) PolicyAuthorizationGet(ctx context.Context,
	appSessionID AppSessionID) (PcfPolicyResponse, error) {
	log.Infof("PCFs PolicyAuthorizationGet Entered for AppSessionID %s",
		string(appSessionID))
	_ = ctx
	if TestPcf {

		pcf.httpClient = HTTPClient

	}
	_, asc, code, err, pd := sendreq(ctx, pcf.httpClient, "get", pcf.pcfRootURI+pcf.pcfURI, AppSessionContext{}, AppSessionContextUpdateData{}, appSessionID)
	//fmt.Println(code)
	pcfPr := PcfPolicyResponse{}

	// check for the presence of the sessid in the database
	_, prs := pcf.paDb[string(appSessionID)]
	// if not found return an error i.e 404
	//if !prs && code != 200 {
	if code != 200 {
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s not found",
			string(appSessionID))
		pcfPr.ResponseCode = uint16(code)
		pcfPr.Pd = &pd
		err = errors.New("failed")

	} else {
		log.Infof("PCFs PolicyAuthorizationGet AppSessionID %s found",
			string(appSessionID))
		if prs {
			pcfPr.ResponseCode = uint16(code)
			pcfPr.Asc = &asc
		}

	}
	return pcfPr, err
}
