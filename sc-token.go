package go_tenable

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

func (sc *TenableSCClient) Login(scUser string, scPassword string) {

	// Make POST request to SC and store token returned in headers
	payload, err := json.Marshal(&TokenRequest{Username: scUser, Password: scPassword})
	if err != nil {
		log.Printf("Unable to marshal login request: %v\n", err)
	}
	req := sc.NewRequest("POST", "token", payload)
	resp := sc.Do(req)
	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read response body: %v\n", err)
	}

	requestCookie := strings.Split(strings.Split(resp.Cookies()[1].Raw, ";")[0], "=")[1]

	var tokenResponse = TokenResponse{}
	err = json.Unmarshal(tmp, &tokenResponse)

	sc.token = tokenResponse.Response.Token
	sc.session = requestCookie

}

func (sc *TenableSCClient) Logout() {
	req := sc.NewRequest("DELETE", "token", nil)
	sc.Do(req)
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
