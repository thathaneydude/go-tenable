package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (sc *TenableSC) ListAssets() AssetListResponse {
	resp, err := sc.Get("asset", "canUse,canManage,owner,groups,ownerGroup,status,name,type,"+
		"description,createdTime,modifiedTime,ipCount,repositories,tags,creator,targetGroup,template")
	tmp, _ := ioutil.ReadAll(resp.Body)
	var Assets = AssetListResponse{}
	err = json.Unmarshal(tmp, &Assets)

	if err != nil {
		log.Printf("Unable to unmarshal Asset Response: %v", err)
	}

	return Assets

}

type AssetResponse struct {
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
}

type AssetListResponse struct {
	Type     string `json:"type"`
	Response struct {
		Usable []AssetResponse `json:"usable"`
	} `json:"response"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}

// Static IP Asset

func (sc TenableSC) NewStaticIPAsset(Name string, IPAddresses string, Description string, Tag string) StaticIPAsset {
	asset := StaticIPAsset{
		ID:          0,
		Name:        Name,
		Type:        "static",
		DefinedIPs:  IPAddresses,
		Description: Description,
		Tag:         Tag,
	}
	return asset
}

type StaticIPAsset struct {
	ID          int
	Name        string
	Type        string
	DefinedIPs  string
	Description string
	Tag         string
}

func (asset *StaticIPAsset) Create(sc TenableSC) StaticIPCreateResponse {
	log.Printf("Creating %v Asset %v\n", asset.Type, asset.Name)
	payload := make(map[string]string)
	payload["name"] = asset.Name
	payload["definedIPs"] = asset.DefinedIPs
	payload["description"] = asset.Description
	payload["tags"] = asset.Tag
	payload["type"] = asset.Type

	bPayload, _ := json.Marshal(payload)

	resp, err := sc.Post("asset", bPayload)
	if err != nil {
		log.Printf("Unable to make asset creation request: %v\n", err)
	}
	var createResponse StaticIPCreateResponse
	bBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		err = json.Unmarshal(bBody, &createResponse)
		asset.ID, _ = strconv.Atoi(createResponse.Response.ID)
	} else {
		log.Printf("Failed to create Asset: %v\n", string(bBody))
	}

	return createResponse
}

type StaticIPCreateResponse struct {
	Type     string `json:"type"`
	Response struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Description  string `json:"description"`
		Tags         string `json:"tags"`
		Context      string `json:"context"`
		Status       string `json:"status"`
		CreatedTime  string `json:"createdTime"`
		ModifiedTime string `json:"modifiedTime"`
		TypeFields   struct {
			DefinedIPs string `json:"definedIPs"`
		} `json:"typeFields"`
		Repositories []struct {
			IPCount    string `json:"ipCount"`
			Repository struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"repository"`
		} `json:"repositories"`
		IPCount         int           `json:"ipCount"`
		Groups          []interface{} `json:"groups"`
		AssetDataFields []interface{} `json:"assetDataFields"`
		CanUse          string        `json:"canUse"`
		CanManage       string        `json:"canManage"`
		Creator         struct {
			ID        string `json:"id"`
			Username  string `json:"username"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
		} `json:"creator"`
		Owner struct {
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
		TargetGroup struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"targetGroup"`
		Template struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"template"`
	} `json:"response"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}

func (asset StaticIPAsset) Edit(sc TenableSC) *http.Response {
	log.Printf("Updating Asset %v (%v)\n", asset.Name, asset.ID)
	if asset.ID == 0 {
		log.Printf("No ID set in struct. Unable to update Asset without its unique identifier\n")
		return nil
	}

	tmpPayload := make(map[string]string)

	if asset.Name != "" {
		tmpPayload["name"] = asset.Name
	}

	if asset.Description != "" {
		tmpPayload["description"] = asset.Description
	}

	if asset.Tag != "" {
		tmpPayload["tags"] = asset.Tag
	}

	if asset.DefinedIPs != "" {
		tmpPayload["definedIPs"] = asset.DefinedIPs
	}

	bPayload, _ := json.Marshal(tmpPayload)

	resp, _ := sc.Patch(fmt.Sprintf("asset/%v", asset.ID), bPayload)

	if resp.StatusCode == 200 {
		log.Printf("Successfully updated asset %v (%v)\n", asset.Name, asset.ID)
	} else {
		bBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error updating asset: %v\n", bBody)
	}
	return resp
}

func (asset *StaticIPAsset) View(sc TenableSC) {
	log.Printf("Fetching latest information on Asset %v (%v)\n", asset.Name, asset.ID)

	resp, _ := sc.Get(fmt.Sprintf("asset/%v", asset.ID), "id,name,description,typeFields,tags")
	var viewRes StaticIPViewResponse
	tmp, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(tmp, &viewRes)
	if err != nil {
		log.Printf("Unable to unmarshal Static IP view response: %v\n", err)
	}

	asset.ID, _ = strconv.Atoi(viewRes.Response.ID)
	asset.Name = viewRes.Response.Name
	asset.Description = viewRes.Response.Description
	asset.DefinedIPs = viewRes.Response.TypeFields.DefinedIPs
	asset.Tag = viewRes.Response.Tags
}

type StaticIPViewResponse struct {
	Type     string `json:"type"`
	Response struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
		TypeFields  struct {
			DefinedIPs string `json:"definedIPs"`
		} `json:"typeFields"`
	} `json:"response"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}

func (asset StaticIPAsset) Calculating(sc TenableSC) bool {
	log.Printf("Fetching info for Asset %v (%v)", asset.Name, asset.ID)
	if asset.ID == 0 {
		log.Printf("No ID set in struct. Unable to get Asset info without its unique identifier")
		return false
	}

	resp, _ := sc.Get(fmt.Sprintf("asset/%v", asset.ID), "id,name,ipCount")
	var statusResponse AssetStatusResponse
	tmp, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(tmp, &statusResponse)

	if err != nil {
		log.Printf("Unable to fetch calculation status for Asset %v (%v): %v", asset.Name, asset.ID, err)
	}

	var ret bool

	if statusResponse.Response.IPCount == -1 {
		ret = true
	} else {
		ret = false
	}

	return ret
}

type AssetStatusResponse struct {
	Type     string `json:"type"`
	Response struct {
		ID      string      `json:"id"`
		Name    string      `json:"name"`
		IPCount interface{} `json:"ipCount"`
	} `json:"response"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
	Warnings  []interface{} `json:"warnings"`
	Timestamp int           `json:"timestamp"`
}

func (asset StaticIPAsset) Delete(sc TenableSC) *http.Response {
	log.Printf("Deleting Asset %v (%v)\n", asset.Name, asset.ID)
	if asset.ID == 0 {
		log.Printf("No ID set in struct. Unable to delete Asset without its unique identifier\n")
		return nil
	}

	resp, _ := sc.Delete(fmt.Sprintf("asset/%v", asset.ID), "")
	return resp
}

//// DNS Asset
//
//type DNSAsset struct {
//	Name            string `json:"name"`
//	DefinedDNSNames string `json:"definedDNSNames"`
//	Description     string `json:"description,omitempty"`
//	Tags            string `json:"tags,omitempty"`
//	Type            string `json:"type"`
//}
//
//func (asset *DNSAsset) ToBytes () []byte {
//	ret, err := json.Marshal(asset)
//	if err != nil {
//		fmt.Printf("Unable to marshal Asset body")
//	}
//	return ret
//}
//
//func (asset *DNSAsset) Create(sc *TenableSC) DNSAssetCreationResponse {
//	req := sc.NewRequest("POST", "asset", asset.ToBytes())
//	resp := sc.Do(req)
//	var createResponse DNSAssetCreationResponse
//	tmp, _ := ioutil.ReadAll(resp.Body)
//	err := json.Unmarshal(tmp, &createResponse)
//	if err != nil {
//		log.Printf("Unable to marshal DNS Asset Create Response: %v", err)
//	}
//	return createResponse
//}
//
//func (asset *DNSAsset) View(sc *TenableSC) {
//
//}
//
//func (asset *DNSAsset) Edit(sc *TenableSC) {
//
//}
//
//func (asset *DNSAsset) Delete(sc *TenableSC) {
//
//}
//
//// Structs
//

//
//type StaticIPAssetViewResponse struct {
//	Type     string `json:"type"`
//	Response struct {
//		ID           string `json:"id"`
//		Name         string `json:"name"`
//		Type         string `json:"type"`
//		Description  string `json:"description"`
//		Tags         string `json:"tags"` // Broken
//		Context      string `json:"context"`
//		Status       string `json:"status"`
//		CreatedTime  string `json:"createdTime"`
//		ModifiedTime string `json:"modifiedTime"`
//		TypeFields   struct {
//			DefinedIPs string `json:"definedIPs"`
//		} `json:"typeFields"`
//		Repositories []struct {
//			IPCount    int `json:"ipCount,string"`
//			Repository struct {
//				ID          string `json:"id"`
//				Name        string `json:"name"`
//				Description string `json:"description"`
//			} `json:"repository"`
//		} `json:"repositories"`
//		IPCount         int           `json:"ipCount,string"`
//		Groups          []interface{} `json:"groups"`
//		AssetDataFields []interface{} `json:"assetDataFields"`
//		CanUse          string        `json:"canUse"`
//		CanManage       string        `json:"canManage"`
//		Creator         struct {
//			ID        string `json:"id"`
//			Username  string `json:"username"`
//			Firstname string `json:"firstname"`
//			Lastname  string `json:"lastname"`
//		} `json:"creator"`
//		Owner struct {
//			ID        string `json:"id"`
//			Username  string `json:"username"`
//			Firstname string `json:"firstname"`
//			Lastname  string `json:"lastname"`
//		} `json:"owner"`
//		OwnerGroup struct {
//			ID          string `json:"id"`
//			Name        string `json:"name"`
//			Description string `json:"description"`
//		} `json:"ownerGroup"`
//		TargetGroup struct {
//			ID          int    `json:"id"`
//			Name        string `json:"name"`
//			Description string `json:"description"`
//		} `json:"targetGroup"`
//		Template struct {
//			ID          int    `json:"id"`
//			Name        string `json:"name"`
//			Description string `json:"description"`
//		} `json:"template"`
//	} `json:"response"`
//	ErrorCode int           `json:"error_code"`
//	ErrorMsg  string        `json:"error_msg"`
//	Warnings  []interface{} `json:"warnings"`
//	Timestamp int           `json:"timestamp"`
//}
//
//type DNSAssetCreationResponse struct {
//	Type     string `json:"type"`
//	Response struct {
//		ID           string `json:"id"`
//		Name         string `json:"name"`
//		Type         string `json:"type"`
//		Description  string `json:"description"`
//		Tags         string `json:"tags"`
//		Context      string `json:"context"`
//		Status       string `json:"status"`
//		CreatedTime  string `json:"createdTime"`
//		ModifiedTime string `json:"modifiedTime"`
//		TypeFields   struct {
//			DefinedDNSNames string `json:"definedDNSNames"`
//		} `json:"typeFields"`
//		Repositories []struct {
//			IPCount    string `json:"ipCount"`
//			Repository struct {
//				ID          string `json:"id"`
//				Name        string `json:"name"`
//				Description string `json:"description"`
//			} `json:"repository"`
//		} `json:"repositories"`
//		IPCount         int           `json:"ipCount"`
//		Groups          []interface{} `json:"groups"`
//		AssetDataFields []interface{} `json:"assetDataFields"`
//		CanUse          string        `json:"canUse"`
//		CanManage       string        `json:"canManage"`
//		Creator         struct {
//			ID        string `json:"id"`
//			Username  string `json:"username"`
//			Firstname string `json:"firstname"`
//			Lastname  string `json:"lastname"`
//		} `json:"creator"`
//		Owner struct {
//			ID        string `json:"id"`
//			Username  string `json:"username"`
//			Firstname string `json:"firstname"`
//			Lastname  string `json:"lastname"`
//		} `json:"owner"`
//		OwnerGroup struct {
//			ID          string `json:"id"`
//			Name        string `json:"name"`
//			Description string `json:"description"`
//		} `json:"ownerGroup"`
//		TargetGroup struct {
//			ID          int    `json:"id"`
//			Name        string `json:"name"`
//			Description string `json:"description"`
//		} `json:"targetGroup"`
//		Template struct {
//			ID          int    `json:"id"`
//			Name        string `json:"name"`
//			Description string `json:"description"`
//		} `json:"template"`
//	} `json:"response"`
//	ErrorCode int           `json:"error_code"`
//	ErrorMsg  string        `json:"error_msg"`
//	Warnings  []interface{} `json:"warnings"`
//	Timestamp int           `json:"timestamp"`
//}
