package go_tenable

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

// With the provided user name and password, attempts to create an authenticated session with Tenable.sc using the
// token endpoint. The function needs to be executed prior to any other SC client request as it sets the token and
func (sc *TenableSC) Login(scUser string, scPassword string) (*TokenResponse, error) {
	// Read in the SC username and password
	payload := TokenRequest{
		scUser,
		scPassword,
	}

	// Make POST request to token endpoint
	resp, err := sc.Post("token", payload.ToBytes())
	if err != nil {
		log.Printf("Unable to request a new token: %v\n", err)
		return nil, err
	}

	// Unmarshal response
	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read login response body: %v\n", err)
		return nil, err
	}

	// Check to see if a cookie was returned with the token
	var cookieFound string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "TNS_SESSIONID" {
			cookieFound = cookie.Value
		}
	}

	// If no cookie was found then return error
	if cookieFound == "" {
		var errorResp = "Unable to find \"TNS_SESSIONID\" cookie in login response\n"
		log.Print(errorResp)
		return nil, errors.New(errorResp)
	}

	// Attempt to unmarshal to a successful response
	var tokenResponse = TokenResponse{}
	err = json.Unmarshal(tmp, &tokenResponse)
	if err != nil {
		// If the response is unsuccessful attmpt to unmarshal to an error response
		// log.Printf("Unable to unmarshal successful token request response: %v", err)
		tokenErrorResponse := TokenError{}
		err = json.Unmarshal(tmp, &tokenErrorResponse)
		if err != nil {
			log.Printf("Unable to unmarshal error response either. %v", err)
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("%v", tokenResponse.ErrorMsg))
	}
	sc.token = tokenResponse.Response.Token
	sc.session = cookieFound

	sc.BaseClient.Headers.Add("X-SecurityCenter", fmt.Sprintf("%v", sc.token))
	sc.BaseClient.Headers.Add("Cookie", fmt.Sprintf("TNS_SESSIONID=%v", sc.session))
	return &tokenResponse, nil
}

func (sc *TenableSC) Logout() error {
	_, err := sc.Delete("token", "")
	if err != nil {
		log.Printf("Unable to log out of Tenable.SC: %v\n", err)
		return err
	}
	return nil
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

type TokenError struct {
	Type      string        `json:"type"`
	Response  []interface{} `json:"response"`
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
