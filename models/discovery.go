package models

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
	DiscoveryConfigRequestId             int    `json:"DiscoveryConfigRequestId,omitempty"`
	DiscoveryConfigExpectedDeviceCount   string `json:"DiscoveryConfigExpectedDeviceCount,omitempty"`
	DiscoveryConfigName                  string `json:"DiscoveryConfigName,omitempty"`
}

// DiscoveryConfigVendorPlatforms to provider vendor platform details.
type DiscoveryConfigVendorPlatforms struct {
	VendorPlatformId                int `json:"VendorPlatformId,omitempty"`
	DiscoveryConfigVendorPlatformId int `json:"DiscoveryConfigVendorPlatformId,omitempty"`
}
