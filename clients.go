package go_tenable

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var restrictedEndpoints = []string{"token"}

// Creating a New Clients
func NewTenableIOClient(accessKey string, secretKey string, transport *http.Transport) TenableIO {
	headers := &http.Header{}
	headers.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", accessKey, secretKey))
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "GoTenable")

	b := &BaseClient{
		&http.Client{},
		headers}

	if transport != nil {
		b.httpClient.Transport = transport
	}
	client := TenableIO{
		*b,
		accessKey,
		secretKey,
		"https://cloud.tenable.com"}
	return client
}

func NewTenableSCClient(scHost string, transport *http.Transport) TenableSC {
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "GoTenable")

	b := &BaseClient{
		&http.Client{},
		headers}

	if transport != nil {
		b.httpClient.Transport = transport
	}

	sc := TenableSC{
		baseClient: *b,
		BaseURL:    fmt.Sprintf("https://%v/rest", scHost),
	}
	return sc
}

func NewNessusClient(accessKey string, secretKey string, nessusAddress string, port int, transport *http.Transport) Nessus {
	headers := &http.Header{}
	headers.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", accessKey, secretKey))
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "GoTenable")

	b := &BaseClient{
		&http.Client{},
		headers}

	if transport != nil {
		b.httpClient.Transport = transport
	}

	nessus := Nessus{
		baseClient: *b,
		accessKey:  accessKey,
		secretKey:  secretKey,
		Address:    nessusAddress,
		Port:       port,
		BaseURL:    fmt.Sprintf("https://%v:%v", nessusAddress, port),
	}
	return nessus
}

// Base Client Functions

func (bc BaseClient) AddHeader(header string, header_value string) {
	bc.headers.Set(header, header_value)
}


func (bc BaseClient) Get(baseURL string, endpoint string, params string) (*http.Response, error) {
	var fullURL string
	if params != "" {
		fullURL = fmt.Sprintf("%v/%v?%v", baseURL, endpoint, params)
	} else {
		fullURL = fmt.Sprintf("%v/%v", baseURL, endpoint)
	}
	log.Printf("Requesting GET --> %v\n", fullURL)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Unable to build GET request \"%v\": %v\n", fullURL, err)
		return nil, err
	}

	req.Header = *bc.headers

	httpResp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	return httpResp, nil
}

func (bc BaseClient) Post(baseURL string, endpoint string, body []byte) (*http.Response, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)

	if !stringInSlice(strings.ToLower(endpoint), restrictedEndpoints) {
		log.Printf("Requesting POST --> %v : %v\n", fullUrl, string(body))
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build POST request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	req.Header = *bc.headers

	resp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Error while running request \"%v\": %v", fullUrl, err)
		return nil, err
	}
	return resp, err
}

func (bc BaseClient) Put(baseURL string, endpoint string, body []byte) (*http.Response, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)
	log.Printf("Requesting PUT --> %v : %v\n", fullUrl, string(body))
	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build PUT request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	req.Header = *bc.headers

	resp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Error while running request \"%v\": %v", fullUrl, err)
		return nil, err
	}
	return resp, err
}

func (bc BaseClient) Patch(baseURL string, endpoint string, body []byte) (*http.Response, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)
	log.Printf("Requesting PATCH --> %v : %v\n", fullUrl, string(body))
	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build PATCH request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	req.Header = *bc.headers

	resp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Error while running request \"%v\": %v", fullUrl, err)
		return nil, err
	}
	return resp, err
}

func (bc BaseClient) Delete(baseURL string, endpoint string, params string) (*http.Response, error) {
	var fullURL string
	if params != "" {
		fullURL = fmt.Sprintf("%v/%v?%v", baseURL, endpoint, params)
	} else {
		fullURL = fmt.Sprintf("%v/%v", baseURL, endpoint)
	}
	log.Printf("Requesting DELETE --> %v\n", fullURL)
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		log.Printf("Unable to build DELETE request \"%v\": %v\n", fullURL, err)
		return nil, err
	}

	req.Header = *bc.headers

	httpResp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	return httpResp, nil
}

// TenableIO Client Base Functions

func (io TenableIO) Get(endpoint string, params string) (*http.Response, error) {
	resp, err := io.baseClient.Get(io.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Post(endpoint string, body []byte) (*http.Response, error) {
	resp, err := io.baseClient.Post(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Put(endpoint string, body []byte) (*http.Response, error) {
	resp, err := io.baseClient.Put(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Patch(endpoint string, body []byte) (*http.Response, error) {
	resp, err := io.baseClient.Patch(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Delete(endpoint string, params string) (*http.Response, error) {
	resp, err := io.baseClient.Delete(io.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TenableSC Client Base Functions

func (sc TenableSC) Get(endpoint string, params string) (*http.Response, error) {
	resp, err := sc.baseClient.Get(sc.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Post(endpoint string, body []byte) (*http.Response, error) {
	resp, err := sc.baseClient.Post(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Put(endpoint string, body []byte) (*http.Response, error) {
	resp, err := sc.baseClient.Put(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Patch(endpoint string, body []byte) (*http.Response, error) {
	resp, err := sc.baseClient.Patch(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Delete(endpoint string, params string) (*http.Response, error) {
	resp, err := sc.baseClient.Delete(sc.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Nessus Client Base Functions

func (n Nessus) Get(endpoint string, params string) (*http.Response, error) {
	resp, err := n.baseClient.Get(n.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Post(endpoint string, body []byte) (*http.Response, error) {
	resp, err := n.baseClient.Post(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Put(endpoint string, body []byte) (*http.Response, error) {
	resp, err := n.baseClient.Put(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Patch(endpoint string, body []byte) (*http.Response, error) {
	resp, err := n.baseClient.Patch(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Delete(endpoint string, params string) (*http.Response, error) {
	resp, err := n.baseClient.Delete(n.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Structs

type BaseClient struct {
	httpClient *http.Client
	headers    *http.Header
}

type TenableIO struct {
	baseClient BaseClient
	accessKey  string
	secretKey  string
	BaseURL    string
}

//
type TenableSC struct {
	baseClient BaseClient
	User       string
	token      int
	session    string
	BaseURL    string
}

type Nessus struct {
	baseClient BaseClient
	accessKey  string
	secretKey  string
	Address    string
	Port       int
	BaseURL    string
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
