package go_tenable

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Public Functions

func NewTenableSCClient(scHost string, transport *http.Transport) TenableSCClient {
	sc := &TenableSCClient{
		client:  &http.Client{Transport: transport},
		baseURL: fmt.Sprintf("https://%v/rest", scHost),
	}
	return *sc
}

func NewTenableIOClient(accessKey string, secretKey string, transport *http.Transport) TenableIOClient {
	client := &http.Client{Transport: transport}

	tio := &TenableIOClient{
		client,
		"https://cloud.tenable.com",
		accessKey,
		secretKey,
	}
	return *tio
}

func NewNessusClient(accessKey string, secretKey string, nessusAddress string, port int, transport *http.Transport) NessusClient {
	client := &http.Client{Transport: transport}

	nessus := &NessusClient{
		client,
		nessusAddress,
		port,
		accessKey,
		secretKey,
	}
	return *nessus
}

// Nessus Client Section

type NessusClient struct {
	client    *http.Client
	address   string
	port      int
	accessKey string
	secretKey string
}

func (tio *NessusClient) Do(req *http.Request) http.Response {
	req.Header.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", tio.accessKey, tio.secretKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoTenable")
	log.Printf("Requesting \"%v\": %v\n", req.URL, req.Body)
	resp, err := tio.client.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
	}
	log.Printf("Response Status [%v]", resp.Status)
	return *resp
}

func (nessus *NessusClient) NewRequest(method string, endpoint string, body []byte) *http.Request {
	fullUrl := fmt.Sprintf("https://%v:%v/%v", nessus.address, nessus.port, endpoint)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build request [%v] %v request\n", method, endpoint)
	}
	return req
}

// Tenable.sc Client Section

type TenableSCClient struct {
	client  *http.Client
	baseURL string
	user    string
	token   int
	session string
}

func (sc *TenableSCClient) NewRequest(method, endpoint string, body []byte) *http.Request {

	fullURL := fmt.Sprintf("%v/%v", sc.baseURL, endpoint)

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	if sc.token != 0 {
		req.Header.Set("X-SecurityCenter", fmt.Sprintf("%v", sc.token))
		req.Header.Set("Cookie", fmt.Sprintf("TNS_SESSIONID=%v", sc.session))
	}

	return req
}

func (sc *TenableSCClient) Do(req *http.Request) *http.Response {
	httpResp, err := sc.client.Do(req)
	if err != nil {
		return nil
	}
	return httpResp
}

// Tenable.io Client Section

type TenableIOClient struct {
	client    *http.Client
	basePath  string
	accessKey string
	secretKey string
}

func (tio *TenableIOClient) Do(req *http.Request) http.Response {
	req.Header.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", tio.accessKey, tio.secretKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoTenable")
	log.Printf("Requesting \"%v\": %v\n", req.URL, req.Body)
	resp, err := tio.client.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
	}

	return *resp
}

func (tio *TenableIOClient) NewRequest(method string, endpoint string, body []byte) *http.Request {
	fullUrl := fmt.Sprintf("%v/%v", tio.basePath, endpoint)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build request [%v] %v request\n", method, endpoint)
	}
	return req
}
