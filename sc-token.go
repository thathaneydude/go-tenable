package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func (sc *TenableSCClient) Login(scUser string, scPassword string) {

	// Make POST request to SC and store token returned in headers
	payload, err := json.Marshal(&TokenRequest{Username: scUser, Password: scPassword})
	if err != nil {
		fmt.Printf("Unable to marshal login request: %v\n", err)
	}
	req, err := sc.NewRequest("POST", "token", payload)
	if err != nil {
		fmt.Printf("Error logging into to SecurityCenter: %v\n", err)
	}
	resp, err := sc.Do(req)
	if err != nil {
		fmt.Printf("Unable to request login token: %v\n", err)
	}
	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read response body: %v\n", err)
	}

	requestCookie := strings.Split(strings.Split(resp.Cookies()[1].Raw, ";")[0], "=")[1]

	var tokenResponse = TokenResponse{}
	err = json.Unmarshal(tmp, &tokenResponse)

	sc.token = tokenResponse.Response.Token
	sc.session = requestCookie

}

func (sc *TenableSCClient) Logout() {
	req, err := sc.NewRequest("DELETE", "token", nil)
	if err != nil {
		fmt.Printf("Unable to log out of Tenable.SC: %v\n", err)
	}
	sc.Do(req)
}
