package go_tenable

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Public Functions
func (nessus *NessusClient) GetStatus() StatusResponse {
	log.Printf("Fetching Nessus Scanner %v status\n", nessus.address)
	var statusResponse StatusResponse

	req := nessus.NewRequest("GET", "server/status", nil)
	resp := nessus.Do(req)
	tmp, _ := ioutil.ReadAll(resp.Body)
	unmarshalError := json.Unmarshal(tmp, &statusResponse)
	if unmarshalError != nil {
		log.Printf("Unable to unmarshal status response: %v", unmarshalError)
	}
	return statusResponse
}

func (nessus *NessusClient) GetProperties() PropertiesResponse {
	log.Printf("Fetching Nessus Scanner %v Properties", nessus.address)
	var propertiesResponse PropertiesResponse

	req := nessus.NewRequest("GET", "server/properties", nil)
	resp := nessus.Do(req)
	tmp, _ := ioutil.ReadAll(resp.Body)
	unmarshalError := json.Unmarshal(tmp, &propertiesResponse)
	if unmarshalError != nil {
		log.Printf("Unable to unmarshal properties response: %v", unmarshalError)
	}
	return propertiesResponse
}

// Nessus Structs

type StatusResponse struct {
	Code     int         `json:"code"`
	Progress interface{} `json:"progress"`
	Status   string      `json:"status"`
}

type PropertiesResponse struct {
	Installers struct {
	} `json:"installers"`
	LoadedPluginSet string `json:"loaded_plugin_set"`
	Features        struct {
		Report       bool `json:"report"`
		RemoteLink   bool `json:"remote_link"`
		Users        bool `json:"users"`
		PluginRules  bool `json:"plugin_rules"`
		API          bool `json:"api"`
		ScanAPI      bool `json:"scan_api"`
		LocalScanner bool `json:"local_scanner"`
		Logs         bool `json:"logs"`
		SMTP         bool `json:"smtp"`
	} `json:"features"`
	ServerUUID string `json:"server_uuid"`
	Update     struct {
		Href       interface{} `json:"href"`
		NewVersion int         `json:"new_version"`
		Restart    int         `json:"restart"`
	} `json:"update"`
	RestartPending struct {
		Reason interface{} `json:"reason"`
		Type   interface{} `json:"type"`
	} `json:"restart_pending"`
	NessusUIVersion string        `json:"nessus_ui_version"`
	NessusType      string        `json:"nessus_type"`
	Notifications   []interface{} `json:"notifications"`
	License         struct {
		Features struct {
			Report       bool `json:"report"`
			RemoteLink   bool `json:"remote_link"`
			Users        bool `json:"users"`
			PluginRules  bool `json:"plugin_rules"`
			API          bool `json:"api"`
			ScanAPI      bool `json:"scan_api"`
			LocalScanner bool `json:"local_scanner"`
			Logs         bool `json:"logs"`
			SMTP         bool `json:"smtp"`
		} `json:"features"`
		Type           string `json:"type"`
		ExpirationDate int    `json:"expiration_date"`
		Ips            int64  `json:"ips"`
		Restricted     bool   `json:"restricted"`
		Agents         int    `json:"agents"`
		Mode           int    `json:"mode"`
		Scanners       int    `json:"scanners"`
		ScannersUsed   int    `json:"scanners_used"`
		AgentsUsed     int    `json:"agents_used"`
		Name           string `json:"name"`
	} `json:"license"`
	FeedError              int         `json:"feed_error"`
	RestartNeeded          interface{} `json:"restart_needed"`
	ServerBuild            string      `json:"server_build"`
	ShowNpv7WhatsNew       int         `json:"show_npv7_whats_new"`
	Npv7DowngradeAvailable int         `json:"npv7_downgrade_available"`
	Capabilities           struct {
		ScanVulnerabilityGroups      bool `json:"scan_vulnerability_groups"`
		ReportEmailConfig            bool `json:"report_email_config"`
		ScanVulnerabilityGroupsMixed bool `json:"scan_vulnerability_groups_mixed"`
	} `json:"capabilities"`
	PluginSet               string      `json:"plugin_set"`
	IdleTimeout             string      `json:"idle_timeout"`
	NessusUIBuild           string      `json:"nessus_ui_build"`
	Npv7UpgradeRequired     bool        `json:"npv7_upgrade_required"`
	ScannerBoottime         int         `json:"scanner_boottime"`
	Npv7                    int         `json:"npv7"`
	LoginBanner             interface{} `json:"login_banner"`
	Npv7UpgradeNotification int         `json:"npv7_upgrade_notification"`
	Platform                string      `json:"platform"`
	ServerVersion           string      `json:"server_version"`
}
