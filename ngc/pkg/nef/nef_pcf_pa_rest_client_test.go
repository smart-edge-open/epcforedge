package ngcnef_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	ngcnef "github.com/otcshare/epcforedge/ngc/pkg/nef"
	"golang.org/x/net/http2"
)

func TestPcfClient_PolicyAuthorizationCreate(t *testing.T) {
	type fields struct {
		Pcf        string
		HTTPClient *http.Client
		PcfRootURI string
		PcfURI     string
	}
	type args struct {
		ctx  context.Context
		body ngcnef.AppSessionContext
	}
	ctx, _ := context.WithCancel(context.Background())
	resBody, err := ioutil.ReadFile(testJSONPath + "appcontextpostresp.json")
	if err != nil {
		fmt.Println(err)
	}
	resBodyBytes := bytes.NewReader(resBody)
	httpclient :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters
			respHeader := make(http.Header)
			respHeader.Set("Location", "1234test")
			return &http.Response{
				StatusCode: 201,
				// Send response to be tested
				Body: ioutil.NopCloser(resBodyBytes),
				// Must be set to non-nil value or it panics
				Header: respHeader,
			}
		})
	resBody2, err1 := ioutil.ReadFile(testJSONPath + "appcontextposterrresp.json")
	if err1 != nil {
		fmt.Println(err1)
	}
	resBodyBytes2 := bytes.NewReader(resBody2)
	httpclient2 :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters

			return &http.Response{
				StatusCode: 403,
				// Send response to be tested
				Body: ioutil.NopCloser(resBodyBytes2),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})
	reqbody, err := ioutil.ReadFile(testJSONPath + "appcontextpostreq.json")
	if err != nil {
		fmt.Println(err)
	}
	var ascreq ngcnef.AppSessionContext
	err = json.Unmarshal(reqbody, &ascreq)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}
	var ascresp ngcnef.AppSessionContext
	err = json.Unmarshal(resBody, &ascresp)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}

	f := fields{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
	f2 := fields{Pcf: "test", HTTPClient: httpclient2, PcfRootURI: "testuri", PcfURI: "test"}
	a := args{ctx: ctx, body: ascreq}
	a2 := args{ctx: ctx, body: ngcnef.AppSessionContext{}}
	pcr1 := ngcnef.PcfPolicyResponse{ResponseCode: 201, Asc: &ascresp}
	problemDetails := ngcnef.ProblemDetails{Cause: "REQUESTED_SERVICE_NOT_AUTHORIZED", Status: 403}
	pcr2 := ngcnef.PcfPolicyResponse{ResponseCode: 403, Pd: &problemDetails}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ngcnef.AppSessionID
		want1   ngcnef.PcfPolicyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{"post_success", f, a, ngcnef.AppSessionID("1234test"), pcr1, false},
		{"post_fail", f2, a2, ngcnef.AppSessionID(""), pcr2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcf := &ngcnef.PcfClient{
				Pcf:        tt.fields.Pcf,
				HTTPClient: tt.fields.HTTPClient,
				PcfRootURI: tt.fields.PcfRootURI,
				PcfURI:     tt.fields.PcfURI,
			}
			got, got1, err := pcf.PolicyAuthorizationCreate(tt.args.ctx, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PcfClient.PolicyAuthorizationCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PcfClient.PolicyAuthorizationCreate() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PcfClient.PolicyAuthorizationCreate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPcfClient_PolicyAuthorizationUpdate(t *testing.T) {
	type fields struct {
		Pcf        string
		HTTPClient *http.Client
		PcfRootURI string
		PcfURI     string
	}
	type args struct {
		ctx          context.Context
		body         ngcnef.AppSessionContextUpdateData
		appSessionID ngcnef.AppSessionID
	}
	ctx, _ := context.WithCancel(context.Background())
	resBody, err := ioutil.ReadFile(testJSONPath + "appcontextpatchresp.json")
	if err != nil {
		fmt.Println(err)
	}
	resBodyBytes := bytes.NewReader(resBody)
	httpclient :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(resBodyBytes),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})
	resBody2, err2 := ioutil.ReadFile(testJSONPath + "appcontextpatcherrresp.json")
	if err2 != nil {
		fmt.Println(err2)
	}
	resBodyBytes2 := bytes.NewReader(resBody2)
	httpclient2 :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters

			return &http.Response{
				StatusCode: 404,
				// Send response to be tested
				Body: ioutil.NopCloser(resBodyBytes2),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})
	reqbody, err := ioutil.ReadFile(testJSONPath + "appcontextpatchreq.json")
	if err != nil {
		fmt.Println(err)
	}
	var ascreq ngcnef.AppSessionContextUpdateData
	err = json.Unmarshal(reqbody, &ascreq)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}
	var ascresp ngcnef.AppSessionContext
	err = json.Unmarshal(resBody, &ascresp)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}
	f := fields{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
	f2 := fields{Pcf: "test", HTTPClient: httpclient2, PcfRootURI: "testuri", PcfURI: "test"}
	a := args{ctx: ctx, body: ascreq, appSessionID: "test1234"}
	a2 := args{ctx: ctx, body: ascreq, appSessionID: ngcnef.AppSessionID("")}
	pcr1 := ngcnef.PcfPolicyResponse{ResponseCode: 200, Asc: &ascresp}
	problemDetails := ngcnef.ProblemDetails{Cause: "APPLICATION_SESSION_CONTEXT_NOT_FOUND", Status: 404}
	pcr2 := ngcnef.PcfPolicyResponse{ResponseCode: 404, Pd: &problemDetails}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ngcnef.PcfPolicyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{"patch_success", f, a, pcr1, false},
		{"patch_fail", f2, a2, pcr2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcf := &ngcnef.PcfClient{
				Pcf:        tt.fields.Pcf,
				HTTPClient: tt.fields.HTTPClient,
				PcfRootURI: tt.fields.PcfRootURI,
				PcfURI:     tt.fields.PcfURI,
			}
			got, err := pcf.PolicyAuthorizationUpdate(tt.args.ctx, tt.args.body, tt.args.appSessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PcfClient.PolicyAuthorizationUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PcfClient.PolicyAuthorizationUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPcfClient_PolicyAuthorizationDelete(t *testing.T) {
	type fields struct {
		Pcf        string
		HTTPClient *http.Client
		PcfRootURI string
		PcfURI     string
	}
	type args struct {
		ctx          context.Context
		appSessionID ngcnef.AppSessionID
	}
	ctx, _ := context.WithCancel(context.Background())

	httpclient :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters

			return &http.Response{
				StatusCode: 204,
				// Send response to be tested
				Body: ioutil.NopCloser(nil),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})

	httpclient2 :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters

			return &http.Response{
				StatusCode: 404,
				// Send response to be tested
				Body: ioutil.NopCloser(nil),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})

	f := fields{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
	f2 := fields{Pcf: "test", HTTPClient: httpclient2, PcfRootURI: "testuri", PcfURI: "test"}
	a := args{ctx: ctx, appSessionID: ngcnef.AppSessionID("1234test")}
	a2 := args{ctx: ctx, appSessionID: ngcnef.AppSessionID("test")}
	pcr1 := ngcnef.PcfPolicyResponse{ResponseCode: 204}

	pcr2 := ngcnef.PcfPolicyResponse{ResponseCode: 404}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ngcnef.PcfPolicyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{"delete_success", f, a, pcr1, false},
		{"delete_fail", f2, a2, pcr2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcf := &ngcnef.PcfClient{
				Pcf:        tt.fields.Pcf,
				HTTPClient: tt.fields.HTTPClient,
				PcfRootURI: tt.fields.PcfRootURI,
				PcfURI:     tt.fields.PcfURI,
			}
			got, err := pcf.PolicyAuthorizationDelete(tt.args.ctx, tt.args.appSessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PcfClient.PolicyAuthorizationDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PcfClient.PolicyAuthorizationDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPcfClient_PolicyAuthorizationGet(t *testing.T) {
	type fields struct {
		Pcf        string
		HTTPClient *http.Client
		PcfRootURI string
		PcfURI     string
	}
	type args struct {
		ctx          context.Context
		appSessionID ngcnef.AppSessionID
	}
	ctx, _ := context.WithCancel(context.Background())
	resBody, err := ioutil.ReadFile(testJSONPath + "appcontextgetresp.json")
	if err != nil {
		fmt.Println(err)
	}
	resBodyBytes := bytes.NewReader(resBody)
	httpclient :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(resBodyBytes),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})
	resBody2, err := ioutil.ReadFile(testJSONPath + "appcontextgeterrresp.json")
	if err != nil {
		fmt.Println(err)
	}
	resBodyBytes2 := bytes.NewReader(resBody2)
	httpclient2 :=
		testingPCFClient(func(req *http.Request) *http.Response {
			// Test request parameters

			return &http.Response{
				StatusCode: 404,
				// Send response to be tested
				Body: ioutil.NopCloser(resBodyBytes2),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})

	var ascresp ngcnef.AppSessionContext
	err = json.Unmarshal(resBody, &ascresp)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}
	f := fields{Pcf: "test", HTTPClient: httpclient, PcfRootURI: "testuri", PcfURI: "test"}
	f2 := fields{Pcf: "test", HTTPClient: httpclient2, PcfRootURI: "testuri", PcfURI: "test"}
	a := args{ctx: ctx, appSessionID: ngcnef.AppSessionID("1234test")}
	a2 := args{ctx: ctx, appSessionID: ngcnef.AppSessionID("test")}
	pcr1 := ngcnef.PcfPolicyResponse{ResponseCode: 200, Asc: &ascresp}
	problemDetails := ngcnef.ProblemDetails{Cause: "CONTEXT_NOT_FOUND", Status: 404}
	pcr2 := ngcnef.PcfPolicyResponse{ResponseCode: 404, Pd: &problemDetails}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ngcnef.PcfPolicyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{"get_success", f, a, pcr1, false},
		{"get_fail", f2, a2, pcr2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcf := &ngcnef.PcfClient{
				Pcf:        tt.fields.Pcf,
				HTTPClient: tt.fields.HTTPClient,
				PcfRootURI: tt.fields.PcfRootURI,
				PcfURI:     tt.fields.PcfURI,
			}
			got, err := pcf.PolicyAuthorizationGet(tt.args.ctx, tt.args.appSessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PcfClient.PolicyAuthorizationGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PcfClient.PolicyAuthorizationGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPCFClientF(t *testing.T) {
	type args struct {
		cfg *ngcnef.Config
	}
	cfgData, err := ioutil.ReadFile(filepath.Clean(NefTestCfgBasepath + "valid_pcf.json"))
	var cfg ngcnef.Config
	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		fmt.Printf("Failed go error :%s", err)
	}
	p := &ngcnef.PcfClient{Pcf: "PCF freegc",
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		PcfRootURI: cfg.PcfPolicyAuthorizationConfig.Scheme + "://" + cfg.PcfPolicyAuthorizationConfig.APIRoot,
		PcfURI:     cfg.PcfPolicyAuthorizationConfig.URI,
	}
	CACert, err1 := ioutil.ReadFile(cfg.PcfPolicyAuthorizationConfig.ClientCert)
	if err1 != nil {
		fmt.Printf("NEF Certification loading Error: %v", err1)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(CACert)

	tlsConfig := &tls.Config{
		//RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	p.HTTPClient.Transport = &http2.Transport{
		TLSClientConfig: tlsConfig,
	}
	a := args{cfg: &cfg}
	tests := []struct {
		name string
		args args
		want *ngcnef.PcfClient
	}{
		// TODO: Add test cases.
		{"client initialize", a, p},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ngcnef.NewPCFClientF(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPCFClientF() = %v, want %v", got, tt.want)
			}
		})
	}
}
