package backup

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Backup struct {
	Checkpoint     ProtectResp    `json:"checkpoint"`
	CheckpointItem CheckpointItem `json:"checkpoint_item"`
}

type ProtectResp struct {
	Status         string   `json:"status"`
	CreatedAt      string   `json:"created_at"`
	Id             string   `json:"id"`
	ResourceGraph  string   `json:"resource_graph"`
	ProjectId      string   `json:"project_id"`
	ProtectionPlan PlanResp `json:"protection_plan"`
}

type PlanResp struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	ExtraInfo string `json:"extra_info"`
}

func (r commonResult) Extract() (*Backup, error) {
	var response Backup
	err := r.ExtractInto(&response)
	return &response, err
}

type QueryResponse struct {
	Protectable []CheckResp `json:"protectable"`
}

type CheckResp struct {
	Result       bool   `json:"result"`
	ResourceType string `json:"resource_type"`
	ErrorCode    string `json:"error_code"`
	ErrorMsg     string `json:"error_msg"`
	ResourceId   string `json:"resource_id"`
}

func (r commonResult) ExtractQueryResponse() (*QueryResponse, error) {
	var response QueryResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type CheckpointItem struct {
	CheckpointId string        `json:"checkpoint_id"`
	CreatedAt    string        `json:"created_at"`
	ExtendInfo   ExtendInfo    `json:"extend_info"`
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	ResourceId   string        `json:"resource_id"`
	Status       string        `json:"status"`
	UpdatedAt    string        `json:"updated_at"`
	BackupData   BackupData    `json:"backup_data"`
	Description  string        `json:"description"`
	Tags         []ResourceTag `json:"tags"`
	ResourceType string        `json:"resource_type"`
}

type ExtendInfo struct {
	AutoTrigger          bool            `json:"auto_trigger"`
	AverageSpeed         int             `json:"average_speed"`
	CopyFrom             string          `json:"copy_from"`
	CopyStatus           string          `json:"copy_status"`
	FailCode             FailCode        `json:"fail_code"`
	FailOp               string          `json:"fail_op"`
	FailReason           string          `json:"fail_reason"`
	ImageType            string          `json:"image_type"`
	Incremental          bool            `json:"incremental"`
	Progress             int             `json:"progress"`
	ResourceAz           string          `json:"resource_az"`
	ResourceName         string          `json:"resource_name"`
	ResourceType         string          `json:"resource_type"`
	Size                 int             `json:"size"`
	SpaceSavingRatio     int             `json:"space_saving_ratio"`
	VolumeBackups        []VolumeBackups `json:"volume_backups"`
	FinishedAt           string          `json:"finished_at"`
	TaskId               string          `json:"taskid"`
	HypervisorType       string          `json:"hypervisor_type"`
	SupportedRestoreMode string          `json:"supported_restore_mode"`
	Supportlld           bool            `json:"support_lld"`
}

type BackupData struct {
	RegionName       string `json:"__openstack_region_name"`
	CloudServiceType string `json:"cloudservicetype"`
	Disk             int    `json:"disk"`
	ImageType        string `json:"imagetype"`
	Ram              int    `json:"ram"`
	Vcpus            int    `json:"vcpus"`
	Eip              string `json:"eip"`
	PrivateIp        string `json:"private_ip"`
}

type FailCode struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type VolumeBackups struct {
	AverageSpeed     int    `json:"average_speed"`
	Bootable         bool   `json:"bootable"`
	Id               string `json:"id"`
	ImageType        string `json:"image_type"`
	Incremental      bool   `json:"incremental"`
	Name             string `json:"name"`
	Size             int    `json:"size"`
	SourceVolumeId   string `json:"source_volume_id"`
	SourceVolumeSize int    `json:"source_volume_size"`
	SpaceSavingRatio int    `json:"space_saving_ratio"`
	Status           string `json:"status"`
	SourceVolumeName string `json:"source_volume_name"`
}

type ResourceTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ExtractBackup will get the backup object from the commonResult
func (r commonResult) ExtractBackup() (*CheckpointItem, error) {
	var s struct {
		Backup *CheckpointItem `json:"checkpoint_item"`
	}

	err := r.ExtractInto(&s)
	return s.Backup, err
}

// BackupPage is the page returned by a pager when traversing over a
// collection of backups.
type BackupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of backups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r BackupPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"checkpoint_items_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a BackupPage struct is empty.
func (r BackupPage) IsEmpty() (bool, error) {
	is, err := ExtractBackups(r)
	return len(is) == 0, err
}

// ExtractBackups accepts a Page struct, specifically a BackupPage struct,
// and extracts the elements into a slice of Backup structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBackups(r pagination.Page) ([]CheckpointItem, error) {
	var s struct {
		Vpcs []CheckpointItem `json:"checkpoint_items"`
	}
	err := (r.(BackupPage)).ExtractInto(&s)
	return s.Vpcs, err
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type QueryResult struct {
	commonResult
}
