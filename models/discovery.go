package models

// DiscoveryJobPayload will be used in create and update functionality
type DiscoveryJobPayload struct {
	ChassisIdentifier               any                        `json:"ChassisIdentifier,omitempty"`
	CommunityString                 bool                       `json:"CommunityString,omitempty"`
	CreateGroup                     bool                       `json:"CreateGroup,omitempty"`
	DiscoveryConfigGroupDescription string                     `json:"DiscoveryConfigGroupDescription,omitempty"`
	DiscoveryConfigGroupID          int                        `json:"DiscoveryConfigGroupId,omitempty"`
	DiscoveryConfigGroupName        string                     `json:"DiscoveryConfigGroupName,omitempty"`
	DiscoveryConfigModels           []DiscoveryConfigModels    `json:"DiscoveryConfigModels,omitempty"`
	DiscoveryConfigParentGroupID    int                        `json:"DiscoveryConfigParentGroupId,omitempty"`
	DiscoveryConfigTaskParam        []DiscoveryConfigTaskParam `json:"DiscoveryConfigTaskParam,omitempty"`
	DiscoveryConfigTasks            []any                      `json:"DiscoveryConfigTasks,omitempty"`
	DiscoveryStatusEmailRecipient   any                        `json:"DiscoveryStatusEmailRecipient,omitempty"`
	Schedule                        ScheduleJob                `json:"Schedule,omitempty"`
	TrapDestination                 bool                       `json:"TrapDestination,omitempty"`
	UseAllProfiles                  any                        `json:"UseAllProfiles,omitempty"`
}

type DiscoveryJobDeletePayload struct {
	DiscoveryGroupIds []int `json:"DiscoveryGroupIds,omitempty"`
}

// DiscoveryJob will be used in read, create and update
type DiscoveryJob struct {
	DiscoveryConfigGroupID          int                        `json:"DiscoveryConfigGroupId,omitempty"`
	DiscoveryConfigGroupName        string                     `json:"DiscoveryConfigGroupName,omitempty"`
	DiscoveryConfigGroupDescription string                     `json:"DiscoveryConfigGroupDescription,omitempty"`
	DiscoveryStatusEmailRecipient   interface{}                `json:"DiscoveryStatusEmailRecipient,omitempty"`
	DiscoveryConfigParentGroupID    int                        `json:"DiscoveryConfigParentGroupId,omitempty"`
	CreateGroup                     bool                       `json:"CreateGroup,omitempty"`
	DiscoveryConfigModels           []DiscoveryConfigModels    `json:"DiscoveryConfigModels,omitempty"`
	DiscoveryConfigTaskParam        []DiscoveryConfigTaskParam `json:"DiscoveryConfigTaskParam,omitempty"`
	DiscoveryConfigTasks            []interface{}              `json:"DiscoveryConfigTasks,omitempty"`
	Schedule                        ScheduleJob                `json:"Schedule,omitempty"`
	TrapDestination                 bool                       `json:"TrapDestination,omitempty"`
	CommunityString                 bool                       `json:"CommunityString,omitempty"`
	ChassisIdentifier               interface{}                `json:"ChassisIdentifier,omitempty"`
	UseAllProfiles                  interface{}                `json:"UseAllProfiles,omitempty"`
}
type DiscoveryConfigTargets struct {
	DiscoveryConfigTargetID int         `json:"DiscoveryConfigTargetId,omitempty"`
	NetworkAddressDetail    string      `json:"NetworkAddressDetail,omitempty"`
	SubnetMask              interface{} `json:"SubnetMask,omitempty"`
	AddressType             int         `json:"AddressType,omitempty"`
	Disabled                bool        `json:"Disabled,omitempty"`
	Exclude                 bool        `json:"Exclude,omitempty"`
}
type DiscoveryConfigModels struct {
	DiscoveryConfigID              int                      `json:"DiscoveryConfigId,omitempty"`
	DiscoveryConfigDescription     string                   `json:"DiscoveryConfigDescription,omitempty"`
	DiscoveryConfigStatus          string                   `json:"DiscoveryConfigStatus,omitempty"`
	DiscoveryConfigTargets         []DiscoveryConfigTargets `json:"DiscoveryConfigTargets,omitempty"`
	ConnectionProfileID            int                      `json:"ConnectionProfileId,omitempty"`
	ConnectionProfile              string                   `json:"ConnectionProfile,omitempty"`
	DeviceType                     []int                    `json:"DeviceType,omitempty"`
	DiscoveryConfigVendorPlatforms []interface{}            `json:"DiscoveryConfigVendorPlatforms,omitempty"`
}
type DiscoveryConfigTaskParam struct {
	TaskID            int `json:"TaskId,omitempty"`
	TaskTypeID        int `json:"TaskTypeId,omitempty"`
	ExecutionSequence int `json:"ExecutionSequence,omitempty"`
}
type ScheduleJob struct {
	RunNow    bool        `json:"RunNow,omitempty"`
	RunLater  bool        `json:"RunLater,omitempty"`
	Recurring interface{} `json:"Recurring,omitempty"`
	Cron      string      `json:"Cron,omitempty"`
	StartTime string      `json:"StartTime,omitempty"`
	EndTime   string      `json:"EndTime,omitempty"`
}
