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

	b := &baseClient{
		&http.Client{},
		headers}

	if transport != nil {
		b.HttpClient.Transport = transport
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

	b := &baseClient{
		&http.Client{},
		headers}

	if transport != nil {
		b.HttpClient.Transport = transport
	}

	sc := TenableSC{
		BaseClient: *b,
		BaseURL:    fmt.Sprintf("https://%v/rest", scHost),
	}
	return sc
}

func NewNessusClient(accessKey string, secretKey string, nessusAddress string, port int, transport *http.Transport) Nessus {
	headers := &http.Header{}
	headers.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", accessKey, secretKey))
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "GoTenable")

	b := &baseClient{
		&http.Client{},
		headers}

	if transport != nil {
		b.HttpClient.Transport = transport
	}

	nessus := Nessus{
		BaseClient: *b,
		accessKey:  accessKey,
		secretKey:  secretKey,
		Address:    nessusAddress,
		Port:       port,
		BaseURL:    fmt.Sprintf("https://%v:%v", nessusAddress, port),
	}
	return nessus
}

// Base Client Functions

func (bc baseClient) Get(baseURL string, endpoint string, params string) (*http.Response, error) {
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

	req.Header = *bc.Headers

	httpResp, err := bc.HttpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	return httpResp, nil
}

func (bc baseClient) Post(baseURL string, endpoint string, body []byte) (*http.Response, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)

	if !stringInSlice(strings.ToLower(endpoint), restrictedEndpoints) {
		log.Printf("Requesting POST --> %v : %v\n", fullUrl, string(body))
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build POST request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	req.Header = *bc.Headers

	resp, err := bc.HttpClient.Do(req)
	if err != nil {
		log.Printf("Error while running request \"%v\": %v", fullUrl, err)
		return nil, err
	}
	return resp, err
}

func (bc baseClient) Put(baseURL string, endpoint string, body []byte) (*http.Response, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)
	log.Printf("Requesting PUT --> %v : %v\n", fullUrl, string(body))
	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build PUT request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	req.Header = *bc.Headers

	resp, err := bc.HttpClient.Do(req)
	if err != nil {
		log.Printf("Error while running request \"%v\": %v", fullUrl, err)
		return nil, err
	}
	return resp, err
}

func (bc baseClient) Patch(baseURL string, endpoint string, body []byte) (*http.Response, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)
	log.Printf("Requesting PATCH --> %v : %v\n", fullUrl, string(body))
	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build PATCH request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	req.Header = *bc.Headers

	resp, err := bc.HttpClient.Do(req)
	if err != nil {
		log.Printf("Error while running request \"%v\": %v", fullUrl, err)
		return nil, err
	}
	return resp, err
}

func (bc baseClient) Delete(baseURL string, endpoint string, params string) (*http.Response, error) {
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

	req.Header = *bc.Headers

	httpResp, err := bc.HttpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	return httpResp, nil
}

// TenableIO Client Base Functions

func (io TenableIO) Get(endpoint string, params string) (*http.Response, error) {
	resp, err := io.BaseClient.Get(io.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Post(endpoint string, body []byte) (*http.Response, error) {
	resp, err := io.BaseClient.Post(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Put(endpoint string, body []byte) (*http.Response, error) {
	resp, err := io.BaseClient.Put(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Patch(endpoint string, body []byte) (*http.Response, error) {
	resp, err := io.BaseClient.Patch(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Delete(endpoint string, params string) (*http.Response, error) {
	resp, err := io.BaseClient.Delete(io.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TenableSC Client Base Functions

func (sc TenableSC) Get(endpoint string, params string) (*http.Response, error) {
	resp, err := sc.BaseClient.Get(sc.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Post(endpoint string, body []byte) (*http.Response, error) {
	resp, err := sc.BaseClient.Post(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Put(endpoint string, body []byte) (*http.Response, error) {
	resp, err := sc.BaseClient.Put(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Patch(endpoint string, body []byte) (*http.Response, error) {
	resp, err := sc.BaseClient.Patch(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Delete(endpoint string, params string) (*http.Response, error) {
	resp, err := sc.BaseClient.Delete(sc.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Nessus Client Base Functions

func (n Nessus) Get(endpoint string, params string) (*http.Response, error) {
	resp, err := n.BaseClient.Get(n.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Post(endpoint string, body []byte) (*http.Response, error) {
	resp, err := n.BaseClient.Post(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Put(endpoint string, body []byte) (*http.Response, error) {
	resp, err := n.BaseClient.Put(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Patch(endpoint string, body []byte) (*http.Response, error) {
	resp, err := n.BaseClient.Patch(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Delete(endpoint string, params string) (*http.Response, error) {
	resp, err := n.BaseClient.Delete(n.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Structs

type baseClient struct {
	HttpClient *http.Client
	Headers    *http.Header
}

type TenableIO struct {
	BaseClient baseClient
	accessKey  string
	secretKey  string
	BaseURL    string
}

//
type TenableSC struct {
	BaseClient baseClient
	User       string
	token      int
	session    string
	BaseURL    string
}

type Nessus struct {
	BaseClient baseClient
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
