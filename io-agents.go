package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func (tio *TenableIOClient) ListAgents() []AgentResponse {
	log.Printf("Fetching all agent information from Tenable.io\n")
	var agentResponses []AgentResponse
	const limit = 5000
	offset := 0
	for {

		agentRes := fetchAgentBatch(tio, limit, offset)
		if len(agentRes.Agents) > 0 {
			agentResponses = append(agentResponses, agentRes)
		}

		if len(agentRes.Agents) < limit {
			break
		}
		offset += limit
	}
	return agentResponses
}

func fetchAgentBatch(tio *TenableIOClient, limit int, offset int) AgentResponse {
	log.Printf("* Fetching agents [%v - %v]\n", offset, offset+limit)
	fullUrl := fmt.Sprintf("scanners/tenable/agents?offset=%v&limit=%v", offset, limit)
	req := tio.NewRequest("GET", fullUrl, nil)
	resp := tio.Do(req)
	tmp, _ := ioutil.ReadAll(resp.Body)
	var agentResponse AgentResponse
	unmarshalError := json.Unmarshal(tmp, &agentResponse)
	if unmarshalError != nil {
		fmt.Println("There was an error:", unmarshalError)
	}

	return agentResponse
}

type AgentResponse struct {
	Agents []struct {
		ID           int    `json:"id"`
		UUID         string `json:"uuid"`
		Name         string `json:"name"`
		Platform     string `json:"platform"`
		Distro       string `json:"distro"`
		IP           string `json:"ip"`
		LastScanned  int    `json:"last_scanned"`
		PluginFeedID string `json:"plugin_feed_id"`
		CoreBuild    string `json:"core_build"`
		CoreVersion  string `json:"core_version"`
		LinkedOn     int    `json:"linked_on"`
		LastConnect  int    `json:"last_connect"`
		Status       string `json:"status"`
		Groups       []struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"groups"`
	} `json:"agents"`
	Pagination struct {
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Sort   []struct {
			Name  string `json:"name"`
			Order string `json:"order"`
		} `json:"sort"`
	} `json:"pagination"`
}
