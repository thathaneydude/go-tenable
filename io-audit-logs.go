package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

func (tio *TenableIOClient) ListEvents(filter EventFilter) ([]Event, IOError) {
	const endpoint = "audit-log/v1/events"
	req := tio.NewRequest("GET", endpoint, nil)
	var funcError IOError

	if filter != (EventFilter{}) {
		GetParams := req.URL.Query()
		filterString := fmt.Sprintf("%v.%v:%v", filter.Filter, filter.Operator, filter.Value)
		GetParams.Add("f", filterString)
		req.URL.RawQuery = GetParams.Encode()
	}

	// Request Logs
	LogResponse := tio.Do(req)
	ResponseBytes, _ := ioutil.ReadAll(LogResponse.Body)
	var Logs AuditLogResponse

	// Unmarshal API response to AuditLogResponse struct
	responseError := json.Unmarshal(ResponseBytes, &Logs)
	if responseError != nil {
		fmt.Printf("Unable to unmarshal Audit Log Response: %v\n", responseError)
	}

	return Logs.Events, funcError
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
