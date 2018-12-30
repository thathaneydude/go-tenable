package tenable_io

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thathaneydude/go-tenable"
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
	fullUrl := fmt.Sprintf("%v/%v/export/%v/chunks/%v", export.tioClient.basePath, export.ExportType,
		export.ExportUUID, ChunkID)
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

func (tio TenableIOClient) NewExport(exportType string) Export {
	var ret = Export{}
	ret.ExportType = exportType
	ret.tioClient = tio
	return ret
}

type ExportRequestResponse struct {
	ExportUUID string `json:"export_uuid"`
}

type ExportStatusResponse struct {
	Status          string `json:"status"`
	ChunksAvailable []int  `json:"chunks_available"`
}

type AssetChunkDownloadResponse []struct {
	ID                        string    `json:"id,omitempty"`
	HasAgent                  bool      `json:"has_agent,omitempty"`
	HasPluginResults          bool      `json:"has_plugin_results,omitempty"`
	CreatedAt                 time.Time `json:"created_at,omitempty"`
	TerminatedAt              time.Time `json:"terminated_at,omitempty"`
	TerminatedBy              time.Time `json:"terminated_by,omitempty"`
	UpdatedAt                 time.Time `json:"updated_at,omitempty"`
	DeletedAt                 time.Time `json:"deleted_at,omitempty"`
	DeletedBy                 time.Time `json:"deleted_by,omitempty"`
	FirstSeen                 time.Time `json:"first_seen,omitempty"`
	LastSeen                  time.Time `json:"last_seen,omitempty"`
	FirstScanTime             time.Time `json:"first_scan_time,omitempty"`
	LastScanTime              time.Time `json:"last_scan_time,omitempty"`
	LastAuthenticatedScanDate time.Time `json:"last_authenticated_scan_date,omitempty"`
	LastLicensedScanDate      time.Time `json:"last_licensed_scan_date,omitempty"`
	AzureVMID                 string    `json:"azure_vm_id,omitempty"`
	AzureResourceID           string    `json:"azure_resource_id,omitempty"`
	AwsEc2InstanceAmiID       string    `json:"aws_ec2_instance_ami_id,omitempty"`
	AwsEc2InstanceID          string    `json:"aws_ec2_instance_id,omitempty"`
	AgentUUID                 string    `json:"agent_uuid,omitempty"`
	BiosUUID                  string    `json:"bios_uuid,omitempty"`
	EnvironmentID             string    `json:"environment_id,omitempty"`
	AwsOwnerID                string    `json:"aws_owner_id,omitempty"`
	AwsAvailabilityZone       string    `json:"aws_availability_zone,omitempty"`
	AwsRegion                 string    `json:"aws_region,omitempty"`
	AwsVpcID                  string    `json:"aws_vpc_id,omitempty"`
	AwsEc2InstanceGroupName   string    `json:"aws_ec2_instance_group_name,omitempty"`
	AwsEc2InstanceStateName   string    `json:"aws_ec2_instance_state_name,omitempty"`
	AwsEc2InstanceType        string    `json:"aws_ec2_instance_type,omitempty"`
	AwsSubnetID               string    `json:"aws_subnet_id,omitempty"`
	AwsEc2ProductCode         string    `json:"aws_ec2_product_code,omitempty"`
	AwsEc2Name                string    `json:"aws_ec2_name,omitempty"`
	McafeeEpoGUID             string    `json:"mcafee_epo_guid,omitempty"`
	McafeeEpoAgentGUID        string    `json:"mcafee_epo_agent_guid,omitempty"`
	ServicenowSysid           string    `json:"servicenow_sysid,omitempty"`
	AgentNames                []string  `json:"agent_names,omitempty"`
	Ipv4S                     []string  `json:"ipv4s,omitempty"`
	Ipv6S                     []string  `json:"ipv6s,omitempty"`
	Fqdns                     []string  `json:"fqdns,omitempty"`
	MacAddresses              []string  `json:"mac_addresses,omitempty"`
	NetbiosNames              []string  `json:"netbios_names,omitempty"`
	OperatingSystems          []string  `json:"operating_systems,omitempty"`
	SystemTypes               []string  `json:"system_types,omitempty"`
	Hostnames                 []string  `json:"hostnames,omitempty"`
	SSHFingerprints           []string  `json:"ssh_fingerprints,omitempty"`
	QualysAssetIds            []string  `json:"qualys_asset_ids,omitempty"`
	QualysHostIds             []string  `json:"qualys_host_ids,omitempty"`
	ManufacturerTpmIds        []string  `json:"manufacturer_tpm_ids,omitempty"`
	SymantecEpHardwareKeys    []string  `json:"symantec_ep_hardware_keys,omitempty"`
	Sources                   []struct {
		Name      string    `json:"name,omitempty"`
		FirstSeen time.Time `json:"first_seen,omitempty"`
		LastSeen  time.Time `json:"last_seen,omitempty"`
	} `json:"sources,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	NetworkInterfaces []struct {
		Name         string   `json:"name,omitempty"`
		Virtual      bool     `json:"virtual,omitempty"`
		Aliased      bool     `json:"aliased,omitempty"`
		Fqdns        []string `json:"fqdns,omitempty"`
		MacAddresses []string `json:"mac_addresses,omitempty"`
		Ipv4S        []string `json:"ipv4s,omitempty"`
		Ipv6S        []string `json:"ipv6s,omitempty"`
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
