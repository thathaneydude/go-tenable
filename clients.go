package go_tenable

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var restrictedEndpoints = []string{"token"}

// Structs & interface

type TenableClient interface {
	Get() (map[string]interface{}, error)
	Post() (map[string]interface{}, error)
	Put() (map[string]interface{}, error)
	Patch() (map[string]interface{}, error)
	Delete() (map[string]interface{}, error)
}

type baseClient struct {
	httpClient http.Client
	headers    http.Header
	transport http.Transport
}

type TenableIO struct {
	baseClient baseClient
	accessKey  string
	secretKey  string
	BaseURL    string
}

type TenableSC struct {
	baseClient baseClient
	User       string
	token      int
	session    string
	BaseURL    string
}

type Nessus struct {
	baseClient baseClient
	accessKey  string
	secretKey  string
	Address    string
	Port       int
	BaseURL    string
}

// Client Constructors
func NewTenableIOClient(accessKey string, secretKey string) TenableIO {
	headers := http.Header{}
	headers.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", accessKey, secretKey))
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "GoTenable")

	b := newBaseClient(headers)

	io := TenableIO{
		b,
		accessKey,
		secretKey,
		"https://cloud.tenable.com"}
	return io
}

func NewTenableSCClient(scHost string) TenableSC {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "GoTenable")

	b := newBaseClient(headers)

	sc := TenableSC{
		baseClient: b,
		BaseURL:    fmt.Sprintf("https://%v/rest", scHost),
	}
	return sc
}

func NewNessusClient(accessKey string, secretKey string, nessusAddress string, port int, transport *http.Transport) Nessus {
	headers := http.Header{}
	headers.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", accessKey, secretKey))
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "GoTenable")

	b := newBaseClient(headers)

	nessus := Nessus{
		baseClient: b,
		accessKey:  accessKey,
		secretKey:  secretKey,
		Address:    nessusAddress,
		Port:       port,
		BaseURL:    fmt.Sprintf("https://%v:%v", nessusAddress, port),
	}
	return nessus
}

func newBaseClient(headers http.Header) baseClient {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	b := baseClient{
		http.Client{},
		headers,
		transport,
	}
	return b
}

// Base Client Functions
func (bc baseClient) Get(baseURL string, endpoint string, params string) (map[string]interface{}, error) {
	var fullURL string
	// Build full URL with GET Params if they exist
	if params != "" {
		fullURL = fmt.Sprintf("%v/%v?%v", baseURL, endpoint, params)
	} else {
		fullURL = fmt.Sprintf("%v/%v", baseURL, endpoint)
	}
	log.Printf("Requesting GET --> %v\n", fullURL)

	// Build the HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Unable to build GET request \"%v\": %v\n", fullURL, err)
		return nil, err
	}

	// Add the base client's headers to the request for authentication / content type
	req.Header = bc.headers

	// Send the request and do error handling
	httpResp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	// Unmarshal the HTTP response to a generic map string interface
	marshaledResponse, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("Unable to read HTTP Response: %v\n", err)
		return nil, err
	}
	var returnMap map[string]interface{}
	err = json.Unmarshal(marshaledResponse, &returnMap)

	return returnMap, nil
}

func (bc baseClient) Post(baseURL string, endpoint string, body []byte) (map[string]interface{}, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)

	if !stringInSlice(strings.ToLower(endpoint), restrictedEndpoints) {
		log.Printf("Requesting POST --> %v : %v\n", fullUrl, string(body))
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build POST request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	// Add the base client's headers to the request for authentication / content type
	req.Header = bc.headers

	// Send the request and do error handling
	httpResp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	// Unmarshal the HTTP response to a generic map string interface
	marshaledResponse, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("Unable to read HTTP Response: %v\n", err)
		return nil, err
	}
	var returnMap map[string]interface{}
	err = json.Unmarshal(marshaledResponse, &returnMap)

	return returnMap, nil
}

func (bc baseClient) Put(baseURL string, endpoint string, body []byte) (map[string]interface{}, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)
	log.Printf("Requesting PUT --> %v : %v\n", fullUrl, string(body))
	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build PUT request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	// Add the base client's headers to the request for authentication / content type
	req.Header = bc.headers

	// Send the request and do error handling
	httpResp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	// Unmarshal the HTTP response to a generic map string interface
	marshaledResponse, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("Unable to read HTTP Response: %v\n", err)
		return nil, err
	}
	var returnMap map[string]interface{}
	err = json.Unmarshal(marshaledResponse, &returnMap)

	return returnMap, nil
}

func (bc baseClient) Patch(baseURL string, endpoint string, body []byte) (map[string]interface{}, error) {
	fullUrl := fmt.Sprintf("%v/%v", baseURL, endpoint)
	log.Printf("Requesting PATCH --> %v : %v\n", fullUrl, string(body))
	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build PATCH request \"%v\": %v \n", fullUrl, err)
		return nil, err
	}

	// Add the base client's headers to the request for authentication / content type
	req.Header = bc.headers

	// Send the request and do error handling
	httpResp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	// Unmarshal the HTTP response to a generic map string interface
	marshaledResponse, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("Unable to read HTTP Response: %v\n", err)
		return nil, err
	}
	var returnMap map[string]interface{}
	err = json.Unmarshal(marshaledResponse, &returnMap)

	return returnMap, nil
}

func (bc baseClient) Delete(baseURL string, endpoint string, params string) (map[string]interface{}, error) {
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

	// Add the base client's headers to the request for authentication / content type
	req.Header = bc.headers

	// Send the request and do error handling
	httpResp, err := bc.httpClient.Do(req)
	if err != nil {
		log.Printf("Unable to run request: %v\n", err)
		return nil, err
	}
	// Unmarshal the HTTP response to a generic map string interface
	marshaledResponse, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("Unable to read HTTP Response: %v\n", err)
		return nil, err
	}
	var returnMap map[string]interface{}
	err = json.Unmarshal(marshaledResponse, &returnMap)

	return returnMap, nil
}

// TenableIO Client Base Functions
func (io TenableIO) SetTransport(transport http.Transport) {
	io.baseClient.transport = transport
}

func (io TenableIO) Get(endpoint string, params string) (map[string]interface{}, error) {
	resp, err := io.baseClient.Get(io.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Post(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := io.baseClient.Post(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Put(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := io.baseClient.Put(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Patch(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := io.baseClient.Patch(io.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (io TenableIO) Delete(endpoint string, params string) (map[string]interface{}, error) {
	resp, err := io.baseClient.Delete(io.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TenableSC Client Base Functions

func (sc TenableSC) Get(endpoint string, params string) (map[string]interface{}, error) {
	resp, err := sc.baseClient.Get(sc.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Post(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := sc.baseClient.Post(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Put(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := sc.baseClient.Put(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Patch(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := sc.baseClient.Patch(sc.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc TenableSC) Delete(endpoint string, params string) (map[string]interface{}, error) {
	resp, err := sc.baseClient.Delete(sc.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Nessus Client Base Functions

func (n Nessus) Get(endpoint string, params string) (map[string]interface{}, error) {
	resp, err := n.baseClient.Get(n.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Post(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := n.baseClient.Post(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Put(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := n.baseClient.Put(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Patch(endpoint string, body []byte) (map[string]interface{}, error) {
	resp, err := n.baseClient.Patch(n.BaseURL, endpoint, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n Nessus) Delete(endpoint string, params string) (map[string]interface{}, error) {
	resp, err := n.baseClient.Delete(n.BaseURL, endpoint, params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
