package go_tenable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func (n *Nessus) GetHealthStats(count int) ScannerSettingsResponse {
	if count == 0 {
		log.Printf("Zero is an invalid number of records to fetch. Setting to 1")
	}

	resp, err := n.Get("settings/health/stats", fmt.Sprintf("count=%v", count))
	if err != nil {
		log.Printf("Unable to fetch health stats for %v: %v", n.Address, err)

	}
	var healthResponse ScannerSettingsResponse
	tmp, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(tmp, &healthResponse)
	if err != nil {
		log.Printf("Unable to unmarshal Health Seetings Response: %v", err)
	}
	return healthResponse
}

type ScannerSettingsResponse struct {
	PerfStatsHistory []struct {
		KbytesReceived      int         `json:"kbytes_received"`
		KbytesSent          int         `json:"kbytes_sent"`
		AvgDNSLookupTime    int         `json:"avg_dns_lookup_time"`
		NumDNSLookups       int         `json:"num_dns_lookups"`
		AvgRdnsLookupTime   int         `json:"avg_rdns_lookup_time"`
		NumRdnsLookups      int         `json:"num_rdns_lookups"`
		CPULoadAvg          int         `json:"cpu_load_avg"`
		NessusLogDiskFree   int         `json:"nessus_log_disk_free"`
		NessusLogDiskTotal  int         `json:"nessus_log_disk_total"`
		NessusDataDiskFree  int         `json:"nessus_data_disk_free"`
		NessusDataDiskTotal int         `json:"nessus_data_disk_total"`
		TempDiskFree        int         `json:"temp_disk_free"`
		TempDiskTotal       int         `json:"temp_disk_total"`
		NumTCPSessions      int         `json:"num_tcp_sessions"`
		NessusVmem          int         `json:"nessus_vmem"`
		NessusMem           int         `json:"nessus_mem"`
		SysRAMUsed          interface{} `json:"sys_ram_used"`
		SysRAM              int         `json:"sys_ram"`
		SysCores            int         `json:"sys_cores"`
		NumHosts            int         `json:"num_hosts"`
		NumScans            int         `json:"num_scans"`
		Timestamp           int         `json:"timestamp"`
	} `json:"perf_stats_history"`
	PerfStatsCurrent struct {
		KbytesReceived      int         `json:"kbytes_received"`
		KbytesSent          int         `json:"kbytes_sent"`
		AvgDNSLookupTime    int         `json:"avg_dns_lookup_time"`
		NumDNSLookups       int         `json:"num_dns_lookups"`
		AvgRdnsLookupTime   int         `json:"avg_rdns_lookup_time"`
		NumRdnsLookups      int         `json:"num_rdns_lookups"`
		CPULoadAvg          int         `json:"cpu_load_avg"`
		NessusLogDiskFree   int         `json:"nessus_log_disk_free"`
		NessusLogDiskTotal  int         `json:"nessus_log_disk_total"`
		NessusDataDiskFree  int         `json:"nessus_data_disk_free"`
		NessusDataDiskTotal int         `json:"nessus_data_disk_total"`
		TempDiskFree        int         `json:"temp_disk_free"`
		TempDiskTotal       int         `json:"temp_disk_total"`
		NumTCPSessions      int         `json:"num_tcp_sessions"`
		NessusVmem          int         `json:"nessus_vmem"`
		NessusMem           int         `json:"nessus_mem"`
		SysRAMUsed          interface{} `json:"sys_ram_used"`
		SysRAM              int         `json:"sys_ram"`
		SysCores            int         `json:"sys_cores"`
		NumHosts            int         `json:"num_hosts"`
		NumScans            int         `json:"num_scans"`
		Timestamp           int         `json:"timestamp"`
	} `json:"perf_stats_current"`
}
