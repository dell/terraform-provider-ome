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
	RunNow    bool   `json:"RunNow,omitempty"`
	RunLater  bool   `json:"RunLater,omitempty"`
	Cron      string `json:"Cron,omitempty"`
	StartTime string `json:"StartTime,omitempty"`
	EndTime   string `json:"EndTime,omitempty"`
}

// DiscoveryConfigTasks to configure discovery task
type DiscoveryConfigTasks struct {
	DiscoveryConfigDescription           string `json:"DiscoveryConfigDescription,omitempty"`
	DiscoveryConfigEmailRecipient        string `json:"DiscoveryConfigEmailRecipient,omitempty"`
	DiscoveryConfigDiscoveredDeviceCount string `json:"DiscoveryConfigDiscoveredDeviceCount,omitempty"`
	DiscoveryConfigRequestID             int    `json:"DiscoveryConfigRequestId,omitempty"`
	DiscoveryConfigExpectedDeviceCount   string `json:"DiscoveryConfigExpectedDeviceCount,omitempty"`
	DiscoveryConfigName                  string `json:"DiscoveryConfigName,omitempty"`
}

// DiscoveryConfigVendorPlatforms to provider vendor platform details.
type DiscoveryConfigVendorPlatforms struct {
	VendorPlatformID                int `json:"VendorPlatformId,omitempty"`
	DiscoveryConfigVendorPlatformID int `json:"DiscoveryConfigVendorPlatformId,omitempty"`
}

// ConnectionProfiles to get the credentials for different protocols
type ConnectionProfiles struct {
	ProfileName        string        `json:"profileName"`
	ProfileDescription string        `json:"profileDescription"`
	Type               string        `json:"type"`
	Credentials        []Credentials `json:"credentials"`
}

// CredSNMP to get the credential of the SNMP protocol.
type CredSNMP struct {
	Community  string `json:"community"`
	EnableV1V2 bool   `json:"enableV1V2"`
	Port       int    `json:"port"`
	Retries    int    `json:"retries"`
	Timeout    int    `json:"timeout"`
}

// CredSSH to get the credential of the SSH protocol.
type CredSSH struct {
	Username        string `json:"username"`
	IsSudoUser      bool   `json:"isSudoUser"`
	Password        string `json:"password"`
	Port            int    `json:"port"`
	UseKey          bool   `json:"useKey"`
	Retries         int    `json:"retries"`
	Timeout         int    `json:"timeout"`
	CheckKnownHosts bool   `json:"checkKnownHosts"`
}

// CredREDFISH to get the credential of the REDFISH protocol.
type CredREDFISH struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CaCheck   bool   `json:"caCheck"`
	CnCheck   bool   `json:"cnCheck"`
	Port      int    `json:"port"`
	Retries   int    `json:"retries"`
	Timeout   int    `json:"timeout"`
	IsHTTP    bool   `json:"isHttp"`
	KeepAlive bool   `json:"keepAlive"`
}

// Credentials to branch out the different protocol based on credentials attribute with interface{} type.
type Credentials struct {
	ID          int         `json:"id"`
	Type        string      `json:"type"`
	AuthType    string      `json:"authType"`
	Modified    bool        `json:"modified"`
	Credential interface{} `json:"credentials,omitempty"`
}

// tfsdk struct

// OmeDiscoveryJob will be used in read, create and update
type OmeDiscoveryJob struct {
	DiscoveryJobID         types.String                `tfsdk:"id"`
	DiscoveryJobName       types.String                `tfsdk:"name"`
	EmailRecipient         types.String                `tfsdk:"email_recipient"`
	DiscoveryConfigTargets []OmeDiscoveryConfigTargets `tfsdk:"discovery_config_targets"`
	JobWait                types.Bool                  `tfsdk:"job_wait"`
	JobWaitTimeout         types.Int64                 `tfsdk:"job_wait_timeout"`
	Schedule               types.String                `tfsdk:"schedule"`
	Cron                   types.String                `tfsdk:"cron"`
	IgnorePartialFailure   types.Bool                  `tfsdk:"ignore_partial_failure"`
	TrapDestination        types.Bool                  `tfsdk:"trap_destination"`
	CommunityString        types.Bool                  `tfsdk:"community_string"`
}

// OmeDiscoveryConfigTargets for discovery configuration
type OmeDiscoveryConfigTargets struct {
	NetworkAddressDetail []types.String `tfsdk:"network_address_detail"`
	DeviceType           []types.String `tfsdk:"device_type"`
	Redfish              OmeRedfish     `tfsdk:"redfish"`
	SNMP                 OmeSNMP        `tfsdk:"snmp"`
	SSH                  OmeSSH         `tfsdk:"ssh"`
}

// OmeRedfish for discovery configuration target REDFISH protocol.
type OmeRedfish struct {
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	Port            types.Int64  `tfsdk:"port"`
	Retries         types.Int64  `tfsdk:"retries"`
	Timeout         types.Int64  `tfsdk:"timeout"`
	CnCheck         types.Bool   `tfsdk:"cn_check"`
	CaCheck         types.Bool   `tfsdk:"ca_check"`
	// Domain          types.String `tfsdk:"domain"`
	// CertificateData types.String `tfsdk:"certificate_data"`
}

// OmeSNMP for discovery configuration target REDFISH protocol.
type OmeSNMP struct {
	Community types.String `tfsdk:"community"`
	Port      types.Int64  `tfsdk:"port"`
	Retries   types.Int64  `tfsdk:"retries"`
	Timeout   types.Int64  `tfsdk:"timeout"`
}

// OmeSSH for discovery configuration target REDFISH protocol.
type OmeSSH struct {
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	Port            types.Int64  `tfsdk:"port"`
	Retries         types.Int64  `tfsdk:"retries"`
	Timeout         types.Int64  `tfsdk:"timeout"`
	CheckKnownHosts types.Bool   `tfsdk:"check_known_hosts"`
	IsSudoUser      types.Bool   `tfsdk:"is_sudo_user"`
}
