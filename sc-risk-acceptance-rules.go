package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (sc *TenableSC) ListRiskAcceptanceRules() AcceptRiskRuleResponse {
	var params = "id,repository,organization,User,plugin,hostType,hostValue,Port,protocol,expires,status,comments," +
		"createdTime,modifiedTime"
	resp, err := sc.Get("acceptRiskRule", params)
	tmp, _ := ioutil.ReadAll(resp.Body)
	var Rules = AcceptRiskRuleResponse{}
	err = json.Unmarshal(tmp, &Rules)

	if err != nil {
		fmt.Printf("Unable to unmarshal Risk Acceptance Rules: %v\n", err)
	}

	return Rules
}

type AcceptRiskRuleResponse struct {
	Type            string           `json:"type"`
	AcceptRiskRules []AcceptRiskRule `json:"response"`
	ErrorCode       int              `json:"error_code"`
	ErrorMsg        string           `json:"error_msg"`
	Warnings        []interface{}    `json:"warnings"`
	Timestamp       int              `json:"timestamp"`
}

type AcceptRiskRule struct {
	HostValue    interface{} `json:"hostValue,omitempty"`
	HostType     string      `json:"hostType"`
	Port         string      `json:"Port,omitempty"`
	Protocol     string      `json:"protocol,omitempty"`
	Expires      string      `json:"expires"`
	CreatedTime  string      `json:"createdTime"`
	ModifiedTime string      `json:"modifiedTime"`
	Status       string      `json:"status"`
	ID           string      `json:"id"`
	Comments     string      `json:"comments"`
	Plugin       struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	} `json:"plugin"`
	Repository struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	} `json:"repository"`
	User struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Firstname string `json:"firstname,omitempty"`
		Lastname  string `json:"lastname,omitempty"`
	} `json:"User"`
}
