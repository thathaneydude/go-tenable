package go_tenable

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func (sc *TenableSCClient) ListAssets() AssetResponse {

	req := sc.NewRequest("GET", "asset", nil)
	resp := sc.Do(req)
	tmp, _ := ioutil.ReadAll(resp.Body)
	var assetResponse = AssetResponse{}
	err := json.Unmarshal(tmp, &assetResponse)

	if err != nil {
		log.Printf("Unable to unmarshal Asset Response: %v", err)
	}

	return assetResponse

}

type AssetResponse struct {
	Type     string `json:"type"`
	Response struct {
		Usable []struct {
			Status       string `json:"status"`
			Name         string `json:"name"`
			Type         string `json:"type"`
			Description  string `json:"description"`
			CreatedTime  string `json:"createdTime"`
			ModifiedTime string `json:"modifiedTime"`
			Tags         string `json:"tags"`
			ID           string `json:"id"`
			Repositories []struct {
				IPCount    string `json:"ipCount"`
				Repository struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"repository"`
			} `json:"repositories"`
			IPCount   int           `json:"ipCount"`
			Groups    []interface{} `json:"groups"`
			CanUse    string        `json:"canUse"`
			CanManage string        `json:"canManage"`
			Owner     struct {
				ID        string `json:"id"`
				Username  string `json:"username"`
				Firstname string `json:"firstname"`
				Lastname  string `json:"lastname"`
			} `json:"owner"`
			OwnerGroup struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"ownerGroup"`
			Template struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"template"`
			TargetGroup struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"targetGroup"`
			Creator struct {
				ID        string `json:"id"`
				Username  string `json:"username"`
				Firstname string `json:"firstname"`
				Lastname  string `json:"lastname"`
			} `json:"creator"`
		} `json:"usable"`
	} `json:"response"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}
