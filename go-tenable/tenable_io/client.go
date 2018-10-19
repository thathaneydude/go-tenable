package tenable_io

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const basePath = "https://cloud.tenable.com"

type SecurityCenterClient struct {
	client  *http.Client
	baseURL *url.URL
	headers map[string]string
	token   string
}

type TenableIOClient struct {
	client    *http.Client
	basePath  string
	accessKey string
	secretKey string
}

func NewTenableIOClient(accessKey string, secretKey string) TenableIOClient {
	client := &http.Client{}
	tio := &TenableIOClient{
		client,
		basePath,
		accessKey,
		secretKey,
	}
	return *tio
}

func (tio *TenableIOClient) Do(req *http.Request) http.Response {
	req.Header.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", tio.accessKey, tio.secretKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoTenable")
	fmt.Printf("Requesting \"%v\": %v\n", req.URL, req.Body)
	resp, err := tio.client.Do(req)
	if err != nil {
		fmt.Printf("Unable to run request: %v\n", err)
	}

	return *resp
}

func (tio *TenableIOClient) NewRequest(method string, endpoint string, body []byte) *http.Request {
	fullUrl := fmt.Sprintf("%v/%v", basePath, endpoint)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build request [%v] %v request\n", method, endpoint)
	}
	return req
}

//func NewClient(httpClient *http.Client, baseURL string) (*SecurityCenterClient, error){
//	if httpClient == nil {
//		httpClient = http.DefaultClient
//	}
//	baseURL = fmt.Sprintf("%s/rest/", baseURL)
//	parsedBaseURL, err := url.Parse(baseURL)
//	if err != nil {
//		return nil, err
//	}
//
//	c := &SecurityCenterClient{
//		client: httpClient,
//		baseURL:parsedBaseURL,
//	}
//
//	return c, nil
//}
//func (c *SecurityCenterClient) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
//	rel, err := url.Parse(urlStr)
//	if err != nil {
//		return nil, err
//	}
//
//	u := c.baseURL.ResolveReference(rel)
//
//	var buf io.ReadWriter
//	if body != nil {
//		buf = new(bytes.Buffer)
//		err = json.NewEncoder(buf).Encode(body)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	req, err := http.NewRequest(method, u.String(), buf)
//	if err != nil {
//		return nil, err
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	if c.token != "" {
//		req.Header.Add("X-SecurityCenter", c.token)
//	}
//
//	return req, nil
//}
//func (c *SecurityCenterClient) Do(req *http.Request) (*http.Response, error) {
//	httpResp, err := c.client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	err = CheckResponse(httpResp)
//	if err != nil {
//		return httpResp, err
//	}
//
//	return httpResp, err
//}
//
//func CheckResponse(r *http.Response) error {
//	if c := r.StatusCode; 200 <= c && c <= 299 {
//		return nil
//	}
//
//	err := fmt.Errorf("Request failed with status code: %d", r.StatusCode)
//	return err
//}
//
//func (c *SecurityCenterClient) Login(username string, password string) *SecurityCenterClient {
//	// Make POST request to SC and store token returned in headers
//	data := url.Values{}
//	data.Set("username", username)
//	data.Set("password", password)
//	req, err := c.NewRequest("POST", "token", bytes.NewBufferString(data.Encode()))
//	if err != nil{
//		fmt.Printf("Error logging into to SecurityCenter: %v", err)
//	}
//	resp, err := c.Do(req)
//	json.NewDecoder(resp.Body).Decode(&SCTokenResponse)
//	fmt.Printf("Token response: %v\n", resp)
//	return c
//}
//
//func (c *SecurityCenterClient) Logout() {
//	// Make DELETE request to SC with headers
//	req, err := c.NewRequest("DELETE", "token", nil)
//	if err != nil {
//		fmt.Printf("Unable to log out with headers: %v", c.headers)
//	}
//	fmt.Printf("Request: %v", req)
//}
