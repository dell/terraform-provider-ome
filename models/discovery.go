package models 

// DiscoveryJobPayload will be used in create and update functionality
type DiscoveryJobPayload struct {
	ChassisIdentifier               any    `json:"ChassisIdentifier"`
	CommunityString                 bool   `json:"CommunityString"`
	CreateGroup                     bool   `json:"CreateGroup"`
	DiscoveryConfigGroupDescription string `json:"DiscoveryConfigGroupDescription"`
	DiscoveryConfigGroupID          int    `json:"DiscoveryConfigGroupId"`
	DiscoveryConfigGroupName        string `json:"DiscoveryConfigGroupName"`
	DiscoveryConfigModels           []struct {
		ConnectionProfile          string `json:"ConnectionProfile"`
		ConnectionProfileID        int    `json:"ConnectionProfileId"`
		DeviceType                 []int  `json:"DeviceType"`
		DiscoveryConfigDescription string `json:"DiscoveryConfigDescription"`
		DiscoveryConfigID          int    `json:"DiscoveryConfigId"`
		DiscoveryConfigStatus      string `json:"DiscoveryConfigStatus"`
		DiscoveryConfigTargets     []struct {
			AddressType             int    `json:"AddressType"`
			Disabled                bool   `json:"Disabled"`
			DiscoveryConfigTargetID int    `json:"DiscoveryConfigTargetId"`
			Exclude                 bool   `json:"Exclude"`
			NetworkAddressDetail    string `json:"NetworkAddressDetail"`
			SubnetMask              any    `json:"SubnetMask"`
		} `json:"DiscoveryConfigTargets"`
		DiscoveryConfigVendorPlatforms []any `json:"DiscoveryConfigVendorPlatforms"`
	} `json:"DiscoveryConfigModels"`
	DiscoveryConfigParentGroupID int `json:"DiscoveryConfigParentGroupId"`
	DiscoveryConfigTaskParam     []struct {
		ExecutionSequence int `json:"ExecutionSequence"`
		TaskID            int `json:"TaskId"`
		TaskTypeID        int `json:"TaskTypeId"`
	} `json:"DiscoveryConfigTaskParam"`
	DiscoveryConfigTasks          []any `json:"DiscoveryConfigTasks"`
	DiscoveryStatusEmailRecipient any   `json:"DiscoveryStatusEmailRecipient"`
	Schedule                      struct {
		Cron      string `json:"Cron"`
		EndTime   any    `json:"EndTime"`
		Recurring any    `json:"Recurring"`
		RunLater  bool   `json:"RunLater"`
		RunNow    bool   `json:"RunNow"`
		StartTime any    `json:"StartTime"`
	} `json:"Schedule"`
	TrapDestination bool `json:"TrapDestination"`
	UseAllProfiles  any  `json:"UseAllProfiles"`
}

// DiscoveryJobResponse will be used in read, create and update
type DiscoveryJobResponse struct {
	ChassisIdentifier               any    `json:"ChassisIdentifier"`
	CommunityString                 bool   `json:"CommunityString"`
	CreateGroup                     bool   `json:"CreateGroup"`
	DiscoveryConfigGroupDescription string `json:"DiscoveryConfigGroupDescription"`
	DiscoveryConfigGroupID          int    `json:"DiscoveryConfigGroupId"`
	DiscoveryConfigGroupName        string `json:"DiscoveryConfigGroupName"`
	DiscoveryConfigModels           []struct {
		ConnectionProfile          string `json:"ConnectionProfile"`
		ConnectionProfileID        int    `json:"ConnectionProfileId"`
		DeviceType                 []int  `json:"DeviceType"`
		DiscoveryConfigDescription string `json:"DiscoveryConfigDescription"`
		DiscoveryConfigID          int    `json:"DiscoveryConfigId"`
		DiscoveryConfigStatus      string `json:"DiscoveryConfigStatus"`
		DiscoveryConfigTargets     []struct {
			AddressType             int    `json:"AddressType"`
			Disabled                bool   `json:"Disabled"`
			DiscoveryConfigTargetID int    `json:"DiscoveryConfigTargetId"`
			Exclude                 bool   `json:"Exclude"`
			NetworkAddressDetail    string `json:"NetworkAddressDetail"`
			SubnetMask              any    `json:"SubnetMask"`
		} `json:"DiscoveryConfigTargets"`
		DiscoveryConfigVendorPlatforms []any `json:"DiscoveryConfigVendorPlatforms"`
	} `json:"DiscoveryConfigModels"`
	DiscoveryConfigParentGroupID int `json:"DiscoveryConfigParentGroupId"`
	DiscoveryConfigTaskParam     []struct {
		ExecutionSequence int `json:"ExecutionSequence"`
		TaskID            int `json:"TaskId"`
		TaskTypeID        int `json:"TaskTypeId"`
	} `json:"DiscoveryConfigTaskParam"`
	DiscoveryConfigTasks          []any `json:"DiscoveryConfigTasks"`
	DiscoveryStatusEmailRecipient any   `json:"DiscoveryStatusEmailRecipient"`
	Schedule                      struct {
		Cron      string `json:"Cron"`
		EndTime   any    `json:"EndTime"`
		Recurring any    `json:"Recurring"`
		RunLater  bool   `json:"RunLater"`
		RunNow    bool   `json:"RunNow"`
		StartTime any    `json:"StartTime"`
	} `json:"Schedule"`
	TrapDestination bool `json:"TrapDestination"`
	UseAllProfiles  any  `json:"UseAllProfiles"`
}

type DiscoveryJobDeletePayload struct {
	DiscoveryGroupIds []int `json:"DiscoveryGroupIds"`
}
