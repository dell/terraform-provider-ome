package models

import "github.com/hashicorp/terraform-plugin-framework/types"

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
	CreateGroup                     bool                       `json:"CreateGroup"`
	DiscoveryConfigModels           []DiscoveryConfigModels    `json:"DiscoveryConfigModels,omitempty"`
	DiscoveryConfigTaskParam        []DiscoveryConfigTaskParam `json:"DiscoveryConfigTaskParam,omitempty"`
	DiscoveryConfigTasks            []DiscoveryConfigTasks     `json:"DiscoveryConfigTasks,omitempty"`
	Schedule                        ScheduleJob                `json:"Schedule,omitempty"`
	TrapDestination                 bool                       `json:"TrapDestination"`
	CommunityString                 bool                       `json:"CommunityString"`
	ChassisIdentifier               string                     `json:"ChassisIdentifier,omitempty"`
	UseAllProfiles                  bool                       `json:"UseAllProfiles"`
}

// DiscoveryConfigTargets for adding device details
type DiscoveryConfigTargets struct {
	DiscoveryConfigTargetID int    `json:"DiscoveryConfigTargetId"`
	NetworkAddressDetail    string `json:"NetworkAddressDetail"`
	SubnetMask              string `json:"SubnetMask"`
	AddressType             int    `json:"AddressType"`
	Disabled                bool   `json:"Disabled"`
	Exclude                 bool   `json:"Exclude"`
}

// DiscoveryConfigModels for discovery configuration
type DiscoveryConfigModels struct {
	DiscoveryConfigID              int                              `json:"DiscoveryConfigId"`
	DiscoveryConfigDescription     string                           `json:"DiscoveryConfigDescription"`
	DiscoveryConfigStatus          string                           `json:"DiscoveryConfigStatus"`
	DiscoveryConfigTargets         []DiscoveryConfigTargets         `json:"DiscoveryConfigTargets"`
	ConnectionProfileID            int                              `json:"ConnectionProfileId"`
	ConnectionProfile              string                           `json:"ConnectionProfile"`
	DeviceType                     []int                            `json:"DeviceType"`
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
	RunNow    bool   `json:"RunNow"`
	RunLater  bool   `json:"RunLater"`
	Cron      string `json:"Cron"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
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
	Credentials        []Protocols `json:"credentials"`
}

// CredSNMP to get the credential of the SNMP protocol.
type CredSNMP struct {
	Community  string `json:"community"`
	EnableV1V2 bool   `json:"enableV1V2"`
	EnableV3   bool   `json:"enableV3"`
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
type Protocols struct {
	ID         int         `json:"id"`
	Type       string      `json:"type"`
	AuthType   string      `json:"authType"`
	Modified   bool        `json:"modified"`
	Credential interface{} `json:"credentials,omitempty"`
}

// tfsdk struct

// OmeDiscoveryJob will be used in read, create and update
type OmeDiscoveryJob struct {
	DiscoveryJobID         types.String                `tfsdk:"id"`
	DiscoveryJobName       types.String                `tfsdk:"name"`
	EmailRecipient         types.String                `tfsdk:"email_recipient"`
	DiscoveryConfigTargets []OmeDiscoveryConfigTargets `tfsdk:"discovery_config_targets"`
	Schedule               types.String                `tfsdk:"schedule"`
	Cron                   types.String                `tfsdk:"cron"`
	TrapDestination        types.Bool                  `tfsdk:"trap_destination"`
	CommunityString        types.Bool                  `tfsdk:"community_string"`
}

// OmeDiscoveryConfigTargets for discovery configuration
type OmeDiscoveryConfigTargets struct {
	NetworkAddressDetail []types.String `tfsdk:"network_address_detail"`
	DeviceType           []types.String `tfsdk:"device_type"`
	Redfish              *OmeRedfish     `tfsdk:"redfish"`
	SNMP                 *OmeSNMP        `tfsdk:"snmp"`
	SSH                  *OmeSSH         `tfsdk:"ssh"`
}

// OmeRedfish for discovery configuration target REDFISH protocol.
type OmeRedfish struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	Retries  types.Int64  `tfsdk:"retries"`
	Timeout  types.Int64  `tfsdk:"timeout"`
	CnCheck  types.Bool   `tfsdk:"cn_check"`
	CaCheck  types.Bool   `tfsdk:"ca_check"`
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
