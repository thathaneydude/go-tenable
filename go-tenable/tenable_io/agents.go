package tenable_io

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (tio TenableIOClient) ListAgents() []AgentResponse {
	fmt.Printf("Fetching all agent information from Tenable.io\n")
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

func fetchAgentBatch(tio TenableIOClient, limit int, offset int) AgentResponse {
	fmt.Printf("* Fetching agents [%v - %v]\n", offset, offset+limit)
	fullUrl := fmt.Sprintf("%v/scanners/tenable/agents?offset=%v&limit=%v", tio.basePath, offset, limit)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Printf("Unable to build agent list request: %v", err)
	}
	req.Header.Set("X-ApiKeys", fmt.Sprintf("accessKey=%v; secretKey=%v;", tio.accessKey, tio.secretKey))
	req.Header.Set("Content-Type", "application/json")
	resp, err := tio.client.Do(req)
	if err != nil {
		fmt.Printf("Unable to request agent list from tenable.io: %v", err)
	}

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
