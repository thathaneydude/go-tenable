package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func (io *TenableIO) ListEvents(filter EventFilter) ([]Event, error) {

	var GetParams string
	if filter != (EventFilter{}) {
		filterString := fmt.Sprintf("%v.%v:%v", filter.Filter, filter.Operator, filter.Value)
		GetParams = fmt.Sprintf("f=%v", filterString)
	} else {
		GetParams = ""
	}
	resp, err := io.Get("audit-log/v1/events", GetParams)
	if err != nil {
		log.Printf("Unable to request audit log events: %v\n", err)
		return nil, err
	}

	ResponseBytes, _ := ioutil.ReadAll(resp.Body)
	var Logs AuditLogResponse

	// Unmarshal API response to AuditLogResponse struct
	responseError := json.Unmarshal(ResponseBytes, &Logs)
	if responseError != nil {
		log.Printf("Unable to unmarshal Audit Log Response: %v\n", responseError)
	}

	return Logs.Events, nil
}

type EventFilter struct {
	Filter   string
	Operator string
	Value    string
}

type Event struct {
	ID          string      `json:"id"`
	Action      string      `json:"action"`
	Crud        string      `json:"crud"`
	IsFailure   bool        `json:"is_failure"`
	Received    time.Time   `json:"received"`
	Description interface{} `json:"description"`
	Actor       struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"actor"`
	IsAnonymous interface{} `json:"is_anonymous"`
	Target      struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"target"`
	Fields []interface{} `json:"fields"`
}

type AuditLogResponse struct {
	Events     []Event `json:"events"`
	Pagination struct {
		Total int `json:"total"`
		Limit int `json:"limit"`
	} `json:"pagination"`
}
