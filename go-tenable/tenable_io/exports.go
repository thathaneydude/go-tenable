package tenable_io

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (export *Export) GetUnprocessedChunks() []int {
	var UnprocessedChunks []int
	fmt.Printf("Chunks Available: %v ; Chunks Processed: %v\n", export.AvailableChunks, export.ProcessedChunks)
	for _, chunkId := range export.AvailableChunks {
		if !intInSlice(chunkId, export.ProcessedChunks) {
			UnprocessedChunks = append(UnprocessedChunks, chunkId)
		}
	}
	return UnprocessedChunks
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (export *Export) RequestExport(Payload []byte) string {
	fullUrl := fmt.Sprintf("%v/%v/export", export.tioClient.basePath, export.ExportType)
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(Payload))
	if err != nil {
		fmt.Printf("Unable to build %v export request: %v\n", export.ExportType, err)
	}
	resp := export.tioClient.Do(req)

	tmp, _ := ioutil.ReadAll(resp.Body)

	var exportRequestRes ExportRequestResponse
	err = json.Unmarshal(tmp, &exportRequestRes)
	if err != nil {
		fmt.Printf("Unable to read response from %v export: %v\n", export.ExportType, err)
	}
	fmt.Printf("Export Response [%v]: %v \n", resp.StatusCode, exportRequestRes)
	export.ExportUUID = exportRequestRes.ExportUUID
	return exportRequestRes.ExportUUID
}

func (export *Export) RequestStatus() string {
	fullUrl := fmt.Sprintf("%v/%v/export/%v/status", export.tioClient.basePath, export.ExportType, export.ExportUUID)
	fmt.Printf("Requesting Status: %v\n", fullUrl)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Printf("Unable to check status for export %v: %v\n", export.ExportUUID, err)
	}
	resp := export.tioClient.Do(req)
	tmp, _ := ioutil.ReadAll(resp.Body)
	var statusRes = ExportStatusResponse{}
	err = json.Unmarshal(tmp, &statusRes)
	if err != nil {
		fmt.Printf("Unable to unmarshal status response: %v", err)
	}
	export.ExportStatus = statusRes.Status
	export.AvailableChunks = statusRes.ChunksAvailable
	fmt.Printf("Export Status: %v\n", export.ExportStatus)
	return export.ExportStatus
}

func (export *Export) DownloadChunk(ChunkID int) AssetChunkDownloadResponse {
	fullUrl := fmt.Sprintf("%v/%v/export/%v/chunks/%v", export.tioClient.basePath, export.ExportType, export.ExportUUID, ChunkID)
	log.Printf("Requesting Chunk: %v\n", fullUrl)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Printf("Unable to download chunk %v: %v\n", ChunkID, err)
	}
	resp := export.tioClient.Do(req)
	tmp, _ := ioutil.ReadAll(resp.Body)
	var ChunkResponse = AssetChunkDownloadResponse{}
	err = json.Unmarshal(tmp, &ChunkResponse)
	if err != nil {
		log.Printf("Unable to unmarshal asset chunk: %v", err)
	}
	return ChunkResponse
}

type Export struct {
	ExportUUID      string
	ExportType      string
	ExportStatus    string
	AvailableChunks []int
	ProcessedChunks []int
	tioClient       TenableIOClient
}

func (tio TenableIOClient) NewAssetExport() Export {
	var ret = Export{}
	ret.ExportType = "assets"
	ret.tioClient = tio
	return ret
}

func NewVulnExport() Export {
	var ret = Export{}
	ret.ExportType = "vulns"
	return ret
}

type ExportRequestResponse struct {
	ExportUUID string `json:"export_uuid"`
}

type ExportStatusResponse struct {
	Status          string `json:"status"`
	ChunksAvailable []int  `json:"chunks_available"`
}

//type VulnerabilityChunkDownloadResponse struct {
//
//}

type AssetChunkDownloadResponse []struct {
	ID                        string        `json:"id"`
	HasAgent                  bool          `json:"has_agent"`
	HasPluginResults          bool          `json:"has_plugin_results"`
	CreatedAt                 time.Time     `json:"created_at"`
	TerminatedAt              interface{}   `json:"terminated_at"`
	TerminatedBy              interface{}   `json:"terminated_by"`
	UpdatedAt                 time.Time     `json:"updated_at"`
	DeletedAt                 interface{}   `json:"deleted_at"`
	DeletedBy                 interface{}   `json:"deleted_by"`
	FirstSeen                 time.Time     `json:"first_seen"`
	LastSeen                  time.Time     `json:"last_seen"`
	FirstScanTime             time.Time     `json:"first_scan_time"`
	LastScanTime              time.Time     `json:"last_scan_time"`
	LastAuthenticatedScanDate time.Time     `json:"last_authenticated_scan_date"`
	LastLicensedScanDate      time.Time     `json:"last_licensed_scan_date"`
	AzureVMID                 interface{}   `json:"azure_vm_id"`
	AzureResourceID           interface{}   `json:"azure_resource_id"`
	AwsEc2InstanceAmiID       interface{}   `json:"aws_ec2_instance_ami_id"`
	AwsEc2InstanceID          interface{}   `json:"aws_ec2_instance_id"`
	AgentUUID                 string        `json:"agent_uuid"`
	BiosUUID                  string        `json:"bios_uuid"`
	EnvironmentID             interface{}   `json:"environment_id"`
	AwsOwnerID                interface{}   `json:"aws_owner_id"`
	AwsAvailabilityZone       interface{}   `json:"aws_availability_zone"`
	AwsRegion                 interface{}   `json:"aws_region"`
	AwsVpcID                  interface{}   `json:"aws_vpc_id"`
	AwsEc2InstanceGroupName   interface{}   `json:"aws_ec2_instance_group_name"`
	AwsEc2InstanceStateName   interface{}   `json:"aws_ec2_instance_state_name"`
	AwsEc2InstanceType        interface{}   `json:"aws_ec2_instance_type"`
	AwsSubnetID               interface{}   `json:"aws_subnet_id"`
	AwsEc2ProductCode         interface{}   `json:"aws_ec2_product_code"`
	AwsEc2Name                interface{}   `json:"aws_ec2_name"`
	McafeeEpoGUID             interface{}   `json:"mcafee_epo_guid"`
	McafeeEpoAgentGUID        interface{}   `json:"mcafee_epo_agent_guid"`
	ServicenowSysid           interface{}   `json:"servicenow_sysid"`
	AgentNames                []string      `json:"agent_names"`
	Ipv4S                     []string      `json:"ipv4s"`
	Ipv6S                     []string      `json:"ipv6s"`
	Fqdns                     []string      `json:"fqdns"`
	MacAddresses              []string      `json:"mac_addresses"`
	NetbiosNames              []interface{} `json:"netbios_names"`
	OperatingSystems          []string      `json:"operating_systems"`
	SystemTypes               []string      `json:"system_types"`
	Hostnames                 []string      `json:"hostnames"`
	SSHFingerprints           []interface{} `json:"ssh_fingerprints"`
	QualysAssetIds            []interface{} `json:"qualys_asset_ids"`
	QualysHostIds             []interface{} `json:"qualys_host_ids"`
	ManufacturerTpmIds        []interface{} `json:"manufacturer_tpm_ids"`
	SymantecEpHardwareKeys    []interface{} `json:"symantec_ep_hardware_keys"`
	Sources                   []struct {
		Name      string    `json:"name"`
		FirstSeen time.Time `json:"first_seen"`
		LastSeen  time.Time `json:"last_seen"`
	} `json:"sources"`
	Tags              []interface{} `json:"tags"`
	NetworkInterfaces []struct {
		Name         string        `json:"name"`
		Virtual      interface{}   `json:"virtual"`
		Aliased      bool          `json:"aliased"`
		Fqdns        []interface{} `json:"fqdns"`
		MacAddresses []string      `json:"mac_addresses"`
		Ipv4S        []string      `json:"ipv4s"`
		Ipv6S        []string      `json:"ipv6s"`
	} `json:"network_interfaces"`
}

type AssetRequestBody struct {
	ChunkSize                 int      `json:"chunk_size"`
	CreatedAt                 float32  `json:"filters.created_at,omitempty"`
	UpdatedAt                 float32  `json:"filters.updated_at,omitempty"`
	TerminatedAt              float32  `json:"filters.terminated_at,omitempty"`
	DeletedAt                 float32  `json:"filters.deleted_at,omitempty"`
	FirstScanTime             float32  `json:"filters.first_scan_time,omitempty"`
	LastAuthenticatedScanTime float32  `json:"filters.last_authenticated_scan_time,omitempty"`
	LastAssessed              float32  `json:"filters.last_assessed,omitempty"`
	ServiceNowSysID           bool     `json:"filters.servicenow_sysid,omitempty"`
	Sources                   []string `json:"filters.sources,omitempty"`
	HasPluginResults          bool     `json:"filters.has_plugin_results,omitempty"`
}

func (req AssetRequestBody) ToBytes() []byte {
	ret, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Unable to marshal request body")
	}
	return ret
}
