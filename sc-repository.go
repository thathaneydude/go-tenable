package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (sc *TenableSC) ListRepositories() RepoListResponse {
	var params = "fields=id,name"
	resp, err := sc.Get("repository", params)
	tmp, _ := ioutil.ReadAll(resp.Body)
	var RepoList = RepoListResponse{}
	err = json.Unmarshal(tmp, &RepoList)
	if err != nil {
		fmt.Printf("Unable to unmarshal the list of Repositories: %v\n", err)
	}

	return RepoList
}



func (sc *TenableSC) RepositoryDetail(RepoId int) RepoDetailResponse {
        var params = "fields=fields=name,description,type,dataFormat,organizations,createdTime,modifiedTime,vulnCount," + 
                             "running,lastSyncTime,lastVulnUpdate,typeFields,correlation"
        resp, err := sc.Get(fmt.Sprintf("repository/%s",RepoId), params)
        tmp, _ := ioutil.ReadAll(resp.Body)
        var RepoDetail = RepoDetailResponse{}
        err = json.Unmarshal(tmp, &RepoDetail)
        if err != nil {
                fmt.Printf("Unable to unmarshal the list of Repositories: %v\n", err)
        }

        return RepoDetail
}




type RepoList struct {
	Name string `json:"name"`
	Id string `json:"id"`
}



type RepoListResponse struct {
	Type   string    `json:"type"`
	Repos []RepoList `json:"response"`
	Error_code int   `json:"error_code"`
	Error_msg string `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}




type RepoDetailResponse struct {
	Type     string `json:"type"`
	Response struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		Type           string `json:"type"`
		DataFormat     string `json:"dataFormat"`
		DownloadFormat string `json:"downloadFormat"`
		CreatedTime    string `json:"createdTime"`
		ModifiedTime   string `json:"modifiedTime"`
		VulnCount      string `json:"vulnCount"`
		Running        string `json:"running"`
		LastSyncTime   string `json:"lastSyncTime"`
		LastVulnUpdate string `json:"lastVulnUpdate"`
		ID             string `json:"id"`
		Organizations  []struct {
			ID          string `json:"id"`
			GroupAssign string `json:"groupAssign"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"organizations"`
		TypeFields struct {
			NessusSchedule struct {
			Type       string `json:"type"`
			Start      string `json:"start"`
			RepeatRule string `json:"repeatRule"`
		} `json:"nessusSchedule"`
	Correlation            []interface{} `json:"correlation"`
	IPRange                string        `json:"ipRange"`
	IPCount                string        `json:"ipCount"`
	RunningNessus          string        `json:"runningNessus"`
	LastGenerateNessusTime string        `json:"lastGenerateNessusTime"`
	TrendingDays           string        `json:"trendingDays"`
	TrendWithRaw           string        `json:"trendWithRaw"`
	LastTrendUpdate        string        `json:"lastTrendUpdate"`
} `json:"typeFields"`
} `json:"response"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}



