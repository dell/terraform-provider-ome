package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// DiscoveryJobPayload will be used in create and update functionality
type DiscoveryJobPayload struct {
	ChassisIdentifier               string                     `json:"ChassisIdentifier,omitempty"`
	CommunityString                 bool                       `json:"CommunityString,omitempty"`
	CreateGroup                     bool                       `json:"CreateGroup,omitempty"`
	DiscoveryConfigGroupDescription string                     `json:"DiscoveryConfigGroupDescription,omitempty"`
	DiscoveryConfigGroupID          int                        `json:"DiscoveryConfigGroupId,omitempty"`
	DiscoveryConfigGroupName        string                     `json:"DiscoveryConfigGroupName,omitempty"`
	DiscoveryConfigModels           []DiscoveryConfigModels    `json:"DiscoveryConfigModels,omitempty"`
	DiscoveryConfigParentGroupID    int                        `json:"DiscoveryConfigParentGroupId,omitempty"`
	DiscoveryConfigTaskParam        []DiscoveryConfigTaskParam `json:"DiscoveryConfigTaskParam,omitempty"`
	DiscoveryConfigTasks            []DiscoveryConfigTasks     `json:"DiscoveryConfigTasks,omitempty"`
	DiscoveryStatusEmailRecipient   string                     `json:"DiscoveryStatusEmailRecipient,omitempty"`
	Schedule                        ScheduleJob                `json:"Schedule,omitempty"`
	TrapDestination                 bool                       `json:"TrapDestination,omitempty"`
	UseAllProfiles                  bool                       `json:"UseAllProfiles,omitempty"`
}

// DiscoveryJobDeletePayload for delete functionality
type DiscoveryJobDeletePayload struct {
	DiscoveryGroupIds []int `json:"DiscoveryGroupIds,omitempty"`
}

// DiscoveryJob will be used in read, create and update
type DiscoveryJob struct {
	DiscoveryConfigGroupID          int                        `json:"DiscoveryConfigGroupId,omitempty"`
	DiscoveryConfigGroupName        string                     `json:"DiscoveryConfigGroupName,omitempty"`
	DiscoveryConfigGroupDescription string                     `json:"DiscoveryConfigGroupDescription,omitempty"`
	DiscoveryStatusEmailRecipient   string                     `json:"DiscoveryStatusEmailRecipient,omitempty"`
	DiscoveryConfigParentGroupID    int                        `json:"DiscoveryConfigParentGroupId,omitempty"`
	CreateGroup                     bool                       `json:"CreateGroup,omitempty"`
	DiscoveryConfigModels           []DiscoveryConfigModels    `json:"DiscoveryConfigModels,omitempty"`
	DiscoveryConfigTaskParam        []DiscoveryConfigTaskParam `json:"DiscoveryConfigTaskParam,omitempty"`
	DiscoveryConfigTasks            []DiscoveryConfigTasks     `json:"DiscoveryConfigTasks,omitempty"`
	Schedule                        ScheduleJob                `json:"Schedule,omitempty"`
	TrapDestination                 bool                       `json:"TrapDestination,omitempty"`
	CommunityString                 bool                       `json:"CommunityString,omitempty"`
	ChassisIdentifier               string                     `json:"ChassisIdentifier,omitempty"`
	UseAllProfiles                  bool                       `json:"UseAllProfiles,omitempty"`
}

// DiscoveryConfigTargets for adding device details
type DiscoveryConfigTargets struct {
	DiscoveryConfigTargetID int    `json:"DiscoveryConfigTargetId,omitempty"`
	NetworkAddressDetail    string `json:"NetworkAddressDetail,omitempty"`
	SubnetMask              string `json:"SubnetMask,omitempty"`
	AddressType             int    `json:"AddressType,omitempty"`
	Disabled                bool   `json:"Disabled,omitempty"`
	Exclude                 bool   `json:"Exclude,omitempty"`
}

// DiscoveryConfigModels for discovery configuration
type DiscoveryConfigModels struct {
	DiscoveryConfigID              int                              `json:"DiscoveryConfigId,omitempty"`
	DiscoveryConfigDescription     string                           `json:"DiscoveryConfigDescription,omitempty"`
	DiscoveryConfigStatus          string                           `json:"DiscoveryConfigStatus,omitempty"`
	DiscoveryConfigTargets         []DiscoveryConfigTargets         `json:"DiscoveryConfigTargets,omitempty"`
	ConnectionProfileID            int                              `json:"ConnectionProfileId,omitempty"`
	ConnectionProfile              string                           `json:"ConnectionProfile,omitempty"`
	DeviceType                     []int                            `json:"DeviceType,omitempty"`
	DiscoveryConfigVendorPlatforms []DiscoveryConfigVendorPlatforms `json:"DiscoveryConfigVendorPlatforms,omitempty"`
}

// DiscoveryConfigTaskParam to config task execution
type DiscoveryConfigTaskParam struct {
	TaskID            int `json:"TaskId,omitempty"`
	TaskTypeID        int `json:"TaskTypeId,omitempty"`
	ExecutionSequence int `json:"ExecutionSequence,omitempty"`
}

// ScheduleJob Schedule of job execution.
type ScheduleJob struct {
	RunNow    bool      `json:"RunNow,omitempty"`
	RunLater  bool      `json:"RunLater,omitempty"`
	Recurring Recurring `json:"Recurring,omitempty"`
	Cron      string    `json:"Cron,omitempty"`
	StartTime string    `json:"StartTime,omitempty"`
	EndTime   string    `json:"EndTime,omitempty"`
}

// DiscoveryConfigTasks to configure discovery task
type DiscoveryConfigTasks struct {
	DiscoveryConfigDescription           string `json:"DiscoveryConfigDescription,omitempty"`
	DiscoveryConfigEmailRecipient        string `json:"DiscoveryConfigEmailRecipient,omitempty"`
	DiscoveryConfigDiscoveredDeviceCount string `json:"DiscoveryConfigDiscoveredDeviceCount,omitempty"`
	DiscoveryConfigRequestId             int    `json:"DiscoveryConfigRequestId,omitempty"`
	DiscoveryConfigExpectedDeviceCount   string `json:"DiscoveryConfigExpectedDeviceCount,omitempty"`
	DiscoveryConfigName                  string `json:"DiscoveryConfigName,omitempty"`
}

// DiscoveryConfigVendorPlatforms to provider vendor platform details.
type DiscoveryConfigVendorPlatforms struct {
	VendorPlatformId                int `json:"VendorPlatformId,omitempty"`
	DiscoveryConfigVendorPlatformId int `json:"DiscoveryConfigVendorPlatformId,omitempty"`
}

// Recurring for schedule job
type Recurring struct {
	Hourly  Hourly  `json:"Hourly,omitempty"`
	Daily   Daily   `json:"Daily,omitempty"`
	Weekley Weekley `json:"Weekley,omitempty"`
}

// Hourly for setting hourly recurring job schedule
type Hourly struct {
	Frequency int `json:"Frequency,omitempty"`
}

// Daily for setting daily recurring job schedule
type Daily struct {
	Frequency int  `json:"Frequency,omitempty"`
	Time      Time `json:"Time,omitempty"`
}

// Weekley for setting daily recurring job schedule
type Weekley struct {
	Day  string `json:"Day,omitempty"`
	Time Time   `json:"Time,omitempty"`
}

// Time for setting minutes and hours.
type Time struct {
	Minutes int `json:"Minutes,omitempty"`
	Hour    int `json:"Hour,omitempty"`
}

// golang struct to convert terraform schema to golang type with tfsdk tagged struct 

// DiscoveryJob will be used in read, create and update
type OmeDiscoveryJob struct {
	DiscoveryConfigGroupID          types.Int64                   `tfsdk:"discovery_config_group_id"`
	DiscoveryConfigGroupName        types.String                  `tfsdk:"discovery_config_group_name"`
	DiscoveryConfigGroupDescription types.String                  `tfsdk:"discovery_config_group_description"`
	DiscoveryStatusEmailRecipient   types.String                  `tfsdk:"discovery_status_email_recipient"`
	DiscoveryConfigParentGroupID    types.Int64                   `tfsdk:"discovery_config_parent_group_id"`
	CreateGroup                     types.Bool                    `tfsdk:"create_group"`
	DiscoveryConfigModels           []OmeDiscoveryConfigModels    `tfsdk:"discovery_config_models"`
	DiscoveryConfigTaskParam        []OmeDiscoveryConfigTaskParam `tfsdk:"discovery_config_task_param"`
	DiscoveryConfigTasks            []OmeDiscoveryConfigTasks     `tfsdk:"discovery_config_tasks"`
	Schedule                        ScheduleJob                   `tfsdk:"schedule"`
	TrapDestination                 types.Bool                    `tfsdk:"trap_destination"`
	CommunityString                 types.Bool                    `tfsdk:"community_string"`
	ChassisIdentifier               types.String                  `tfsdk:"chassis_identifier"`
	UseAllProfiles                  types.Bool                    `tfsdk:"use_all_profiles"`
}

// DiscoveryConfigTargets for adding device details
type OmeDiscoveryConfigTargets struct {
	DiscoveryConfigTargetID types.Int64  `tfsdk:"discovery_config_target_id"`
	NetworkAddressDetail    types.String `tfsdk:"network_address_detail"`
	SubnetMask              types.String `tfsdk:"subnet_mask"`
	AddressType             types.Int64  `tfsdk:"address_type"`
	Disabled                types.Bool   `tfsdk:"disabled"`
	Exclude                 types.Bool   `tfsdk:"exclude"`
}

// DiscoveryConfigModels for discovery configuration
type OmeDiscoveryConfigModels struct {
	DiscoveryConfigID              types.Int64                         `tfsdk:"discovery_config_id"`
	DiscoveryConfigDescription     types.String                        `tfsdk:"discovery_config_description"`
	DiscoveryConfigStatus          types.String                        `tfsdk:"discovery_config_status"`
	DiscoveryConfigTargets         []OmeDiscoveryConfigTargets         `tfsdk:"discovery_config_targets"`
	ConnectionProfileID            types.Int64                         `tfsdk:"connection_profile_id"`
	ConnectionProfile              types.String                        `tfsdk:"connection_profile"`
	DeviceType                     []types.Int64                       `tfsdk:"device_type"`
	DiscoveryConfigVendorPlatforms []OmeDiscoveryConfigVendorPlatforms `tfsdk:"discovery_config_vendor_platforms"`
}

// DiscoveryConfigTaskParam to config task execution
type OmeDiscoveryConfigTaskParam struct {
	TaskID            types.Int64 `tfsdk:"task_id"`
	TaskTypeID        types.Int64 `tfsdk:"task_type_id"`
	ExecutionSequence types.Int64 `tfsdk:"execution_sequence"`
}

// ScheduleJob Schedule of job execution.
type OmeScheduleJob struct {
	RunNow    types.Bool   `tfsdk:"run_now"`
	RunLater  types.Bool   `tfsdk:"run_later"`
	Recurring OmeRecurring `tfsdk:"recurring"`
	Cron      types.String `tfsdk:"cron"`
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
}

// DiscoveryConfigTasks to configure discovery task
type OmeDiscoveryConfigTasks struct {
	DiscoveryConfigDescription           types.String `tfsdk:"discovery_config_description"`
	DiscoveryConfigEmailRecipient        types.String `tfsdk:"discovery_config_email_recipient"`
	DiscoveryConfigDiscoveredDeviceCount types.String `tfsdk:"discovery_config_discovered_device_count"`
	DiscoveryConfigRequestId             types.Int64  `tfsdk:"discovery_config_request_id"`
	DiscoveryConfigExpectedDeviceCount   types.String `tfsdk:"discovery_config_expected_device_count"`
	DiscoveryConfigName                  types.String `tfsdk:"discovery_config_name"`
}

// DiscoveryConfigVendorPlatforms to provider vendor platform details.
type OmeDiscoveryConfigVendorPlatforms struct {
	VendorPlatformId                types.Int64 `tfsdk:"vendor_platform_id"`
	DiscoveryConfigVendorPlatformId types.Int64 `tfsdk:"discovery_config_vendor_platform_id"`
}

// Recurring for schedule job
type OmeRecurring struct {
	Hourly  OmeHourly  `tfsdk:"hourly"`
	Daily   OmeDaily   `tfsdk:"daily"`
	Weekley OmeWeekley `tfsdk:"weekley"`
}

// Hourly for setting hourly recurring job schedule
type OmeHourly struct {
	Frequency types.Int64 `tfsdk:"frequency"`
}

// Daily for setting daily recurring job schedule
type OmeDaily struct {
	Frequency types.Int64 `tfsdk:"frequency"`
	Time      OmeTime     `tfsdk:"time"`
}

// Weekley for setting daily recurring job schedule
type OmeWeekley struct {
	Day  types.String `tfsdk:"day"`
	Time OmeTime      `tfsdk:"time"`
}

// Time for setting minutes and hours.
type OmeTime struct {
	Minutes types.Int64 `tfsdk:"minutes"`
	Hour    types.Int64 `tfsdk:"hour"`
}
