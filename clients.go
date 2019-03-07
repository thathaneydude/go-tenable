package go_tenable

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

// Public Functions

func NewTenableSCClient(scHost string) TenableSCClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sc := &TenableSCClient{
		client:  &http.Client{Transport: tr},
		baseURL: fmt.Sprintf("https://%v/rest", scHost),
	}
	return *sc
}

func NewTenableIOClient(accessKey string, secretKey string) TenableIOClient {
	client := &http.Client{}
	tio := &TenableIOClient{
		client,
		tenableIOBasePath,
		accessKey,
		secretKey,
	}
	return *tio
}

// Tenable.sc Client Section

type TenableSCClient struct {
	client  *http.Client
	baseURL string
	user    string
	token   int
	session string
}

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Type     string `json:"type"`
	Response struct {
		LastLogin        string `json:"lastLogin"`
		LastLoginIP      string `json:"lastLoginIP"`
		FailedLogins     string `json:"failedLogins"`
		FailedLoginIP    string `json:"failedLoginIP"`
		LastFailedLogin  string `json:"lastFailedLogin"`
		Token            int    `json:"token"`
		UnassociatedCert string `json:"unassociatedCert"`
	} `json:"response"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}

type TokenCookie struct {
	HTTPOnly bool   `json:"httpOnly"`
	Path     string `json:"path"`
	Secure   bool   `json:"secure"`
	Value    string `json:"value"`
}

func (sc *TenableSCClient) NewRequest(method, endpoint string, body []byte) (*http.Request, error) {

	fullURL := fmt.Sprintf("%v/%v", sc.baseURL, endpoint)

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if sc.token != 0 {
		req.Header.Set("X-SecurityCenter", fmt.Sprintf("%v", sc.token))
		req.Header.Set("Cookie", fmt.Sprintf("TNS_SESSIONID=%v", sc.session))
	}

	return req, err
}

func (sc *TenableSCClient) Do(req *http.Request) (*http.Response, error) {
	httpResp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	return httpResp, err
}

// Tenable.io Client Section

const tenableIOBasePath = "https://cloud.tenable.com"

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
	fullUrl := fmt.Sprintf("%v/%v", tenableIOBasePath, endpoint)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to build request [%v] %v request\n", method, endpoint)
	}
	return req
}
