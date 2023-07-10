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
type DiscoveryConfigTargets struct {
	DiscoveryConfigTargetID int    `json:"DiscoveryConfigTargetId,omitempty"`
	NetworkAddressDetail    string `json:"NetworkAddressDetail,omitempty"`
	SubnetMask              string `json:"SubnetMask,omitempty"`
	AddressType             int    `json:"AddressType,omitempty"`
	Disabled                bool   `json:"Disabled,omitempty"`
	Exclude                 bool   `json:"Exclude,omitempty"`
}
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
type DiscoveryConfigTaskParam struct {
	TaskID            int `json:"TaskId,omitempty"`
	TaskTypeID        int `json:"TaskTypeId,omitempty"`
	ExecutionSequence int `json:"ExecutionSequence,omitempty"`
}
type ScheduleJob struct {
	RunNow    bool      `json:"RunNow,omitempty"`
	RunLater  bool      `json:"RunLater,omitempty"`
	Recurring Recurring `json:"Recurring,omitempty"`
	Cron      string    `json:"Cron,omitempty"`
	StartTime string    `json:"StartTime,omitempty"`
	EndTime   string    `json:"EndTime,omitempty"`
}

type DiscoveryConfigTasks struct {
	DiscoveryConfigDescription           string `json:"DiscoveryConfigDescription,omitempty"`
	DiscoveryConfigEmailRecipient        string `json:"DiscoveryConfigEmailRecipient,omitempty"`
	DiscoveryConfigDiscoveredDeviceCount string `json:"DiscoveryConfigDiscoveredDeviceCount,omitempty"`
	DiscoveryConfigRequestId             int    `json:"DiscoveryConfigRequestId,omitempty"`
	DiscoveryConfigExpectedDeviceCount   string `json:"DiscoveryConfigExpectedDeviceCount,omitempty"`
	DiscoveryConfigName                  string `json:"DiscoveryConfigName,omitempty"`
}

type DiscoveryConfigVendorPlatforms struct {
	VendorPlatformId                int `json:"VendorPlatformId,omitempty"`
	DiscoveryConfigVendorPlatformId int `json:"DiscoveryConfigVendorPlatformId,omitempty"`
}

type Recurring struct {
	Hourly  Hourly  `json:"Hourly,omitempty"`
	Daily   Daily   `json:"Daily,omitempty"`
	Weekley Weekley `json:"Weekley,omitempty"`
}

type Hourly struct {
	Frequency int `json:"Frequency,omitempty"`
}

type Daily struct {
	Frequency int  `json:"Frequency,omitempty"`
	Time      Time `json:"Time,omitempty"`
}

type Weekley struct {
	Day  string `json:"Day,omitempty"`
	Time Time   `json:"Time,omitempty"`
}

type Time struct {
	Minutes int `json:"Minutes,omitempty"`
	Hour    int `json:"Hour,omitempty"`
}
