package go_tenable

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func (io *TenableIO) ListAgents(limit int) []map[string]interface{} {
	log.Info("Fetching all agent information from Tenable.io\n")

	offset := 0
	var fullAgentList []map[string]interface{}
	for {
		agentResponse, err := fetchAgentBatch(io, limit, offset)
		if err != nil {
			// Error already logged when fetching the batch
			continue
		}
		currentAgentList := agentResponse["agents"]
		if len(currentAgentList.([]interface{})) > 0 {
			fullAgentList = append(fullAgentList, currentAgentList.(map[string]interface{}))
		}

		if len(currentAgentList.([]interface{})) < limit {
			break
		}

		offset += limit
	}
	return fullAgentList
}

func fetchAgentBatch(io *TenableIO, limit int, offset int) (map[string]interface{}, error) {
	log.Printf("Fetching agents [%v - %v]\n", offset, offset+limit)

	resp, err := io.Get("scanners/1/agents",
		fmt.Sprintf("offset=%v&limit=%v", offset, offset+limit))
	if err != nil {
		log.Printf("Unable to fetch Agent batch: %v\n", err)
		return nil, err
	}
	return resp, nil
}
