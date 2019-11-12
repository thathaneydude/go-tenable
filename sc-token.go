package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func (sc *TenableSC) Login(scUser string, scPassword string) {
	payload := TokenRequest{
		scUser,
		scPassword,
	}
	resp, err := sc.Post("token", payload.ToBytes())
	if err != nil {
		log.Printf("Unable to request a new token: %v\n", err)
	}

	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read response body: %v\n", err)
	}

	requestCookie := strings.Split(strings.Split(resp.Cookies()[1].Raw, ";")[0], "=")[1]

	var tokenResponse = TokenResponse{}
	err = json.Unmarshal(tmp, &tokenResponse)

	sc.token = tokenResponse.Response.Token
	sc.session = requestCookie

	sc.BaseClient.Headers.Add("X-SecurityCenter", fmt.Sprintf("%v", sc.token))
	sc.BaseClient.Headers.Add("Cookie", fmt.Sprintf("TNS_SESSIONID=%v", sc.session))

}

func (sc *TenableSC) Logout() {
	_, err := sc.Delete("token", "")
	if err != nil {
		log.Printf("Unable to log out of Tenable.SC: %v\n", err)
	}
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

func (req TokenRequest) ToBytes() []byte {
	ret, err := json.Marshal(req)
	if err != nil {
		log.Printf("Unable to marshal token request body")
	}
	return ret
}
