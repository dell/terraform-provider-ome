package clients

import "time"

// General constants
const (
	// SessionType - defines the source of session
	SessionType = "API"
	// AuthTokenHeader - key for fetching auth token from header
	AuthTokenHeader = "x-auth-token" // #nosec G101
	// SuccessStatusID - job success status ID
	SuccessStatusID = 2060
	// waitTime - sleep interval between retries
	waitTime = 5 * time.Second
	// Retries - Number of http retries
	Retries = 3
	//ServiceTags - constant servivetags to identify the input
	ServiceTags = "servicetags"
	//DeviceIDs - constant deviceids to identify the input
	DeviceIDs = "deviceids"
)

// API's constants
const (
	// SessionAPI - api used to manage session
	SessionAPI = "/api/SessionService/Sessions"
	// TemplateAPI - api used to manage templates
	TemplateAPI = "/api/TemplateService/Templates"
	// JobAPI - api used to manage jobs
	JobAPI = "/api/JobService/Jobs"
	// IdentityPoolAPI - api used to manage IdentityPools
	IdentityPoolAPI = "/api/IdentityPoolService/IdentityPools"
	// UpdateNetworkConfigAPI - api used to update network configuration to the template
	UpdateNetworkConfigAPI = "/api/TemplateService/Actions/TemplateService.UpdateNetworkConfig"
	// LastExecDetailAPI - api used to get last execution details
	LastExecDetailAPI = "/api/JobService/Jobs(%d)/LastExecutionDetail"
	//DeviceAPI - api for managing devices
	DeviceAPI = "/api/DeviceService/Devices"
	// TemplateViewTypeAPI - api to get view type
	TemplateViewTypeAPI = "/api/TemplateService/TemplateViewTypes"
	// TemplateDeviceTypeAPI - api to get view type
	TemplateDeviceTypeAPI = "/api/TemplateService/TemplateTypes"
	//GroupAPI - api for fetching group details
	GroupAPI = "/api/GroupService/Groups"
	//GroupServiceDevicesAPI - api for fetching device ids from a group id
	GroupServiceDevicesAPI = GroupAPI + "(%d)/Devices"
	//DeployAPI - api to deploy a template on the given devices
	DeployAPI = "/api/TemplateService/Actions/TemplateService.Deploy"
	//ProfileAPI - api to manage profiles
	ProfileAPI = "/api/ProfileService/Profiles"
	//UnAssignProfileAPI - api to unassign profile
	UnAssignProfileAPI = "/api/ProfileService/Actions/ProfileService.UnassignProfiles"
	//DeleteProfileAPI - api to delete profile
	DeleteProfileAPI = "/api/ProfileService/Actions/ProfileService.Delete"
	//CloneTemplateAPI - api to clone a template
	CloneTemplateAPI = "/api/TemplateService/Actions/TemplateService.Clone"
)

// Messages constants
const (
	// ErrRetryTimeoutMsg - retry timeout error message
	ErrRetryTimeoutMsg = "request time out after retrying %d times"
	// ErrResponseMsg - error response message
	ErrResponseMsg = "status: %d, body: %s"
	// ErrEmptyBodyMsg - error empty body message
	ErrEmptyBodyMsg = "body cannot be empty"
	// ErrInvalidDeviceIdentifiers - error message for invalid device service tag
	ErrInvalidDeviceIdentifiers = "invalid device servicetag or id"
	// ErrEmptyDeviceDetails - error message when both device service tag and device id not given
	ErrEmptyDeviceDetails = "either Device ID or Servicetag is required"
	// ErrInvalidFqdds = error message for invalid fqdds
	ErrInvalidFqdds = "Invalid FQDDS for template creation"
	// ErrInvalidTemplateViewType - error message for invalid template view type
	ErrInvalidTemplateViewType = "Invalid template view type for template creation"
	// SuccessMsg - job success message
	SuccessMsg = "Successfully completed the job"
	// JobIncompleteMsg - job incomplete message after retries
	JobIncompleteMsg = "Job %d incomplete after polling %d times...Check status in console"
	// SuccessTemplateMessage - message returned on sucessful creation of template
	SuccessTemplateMessage = "template created successfully"
	// ErrTemplateMessage - message returned when error encountered on creation of template
	ErrTemplateMessage = "error occurred while creating a template"
	// IncompleteTemplateMsg - template incomplete message after retries
	IncompleteTemplateMsg = "Status of template with name %s and id %d could not be checked after %d times...Check status in console"
	// ErrInvalidNetworkDetails - message returned when error encountered on creation of template
	ErrInvalidNetworkDetails = "unable to find the combination of nic identifier and port in template nic model"
	//ErrInvalidIdentityPool - message returned when the given identityPool name is invalid
	ErrInvalidIdentityPool = "IdentityPool: '%s' is not available in the appliance"
	// ErrUnableToParseVlan - message returned when Vlan could not be parsed from plan or state
	ErrUnableToParseVlan = "unable to parse vlan data for the template from plan/state"
	// ErrUnableToParseData - message returned when unable to parse the plan/state
	ErrUnableToParseData = "unable to parse data from plan/state"
	// ErrUnableToParseBootToNetISO - message returned when unable to parse the boot to netowrk iso data
	ErrUnableToParseBootToNetISO = "failed to parse boot_to_network_iso attributes"
	// ErrDeviceRequired - message returned when device data is not specified
	ErrDeviceRequired = "please provide device IDs or service tags"
	// ErrDeviceMutuallyExclusive - message returned when device data is not specified
	ErrDeviceMutuallyExclusive = "please provide one of the device IDs or service tags"
	// ErrInvalidTemplate - message returned when invalid template id or name
	ErrInvalidTemplate = "please provide a valid template ID or name"
	// ErrPlanToTfsdkConversion - message returned when converting the plan to tfsdk
	ErrPlanToTfsdkConversion = "Error occured converting the plan values to tfsdk struct"
	// ErrStateToTfsdkConversion - message returned when converting the state to tfsdk
	ErrStateToTfsdkConversion = "Error occured converting the state values to tfsdk struct"
	// ErrStateToTfsdkConversion - message returned when template id or name changed
	ErrTemplateChanges = "template id or name cannot be changed"
	// ErrTemplateDeploymentGeneral - message returned when template deployment fails
	ErrTemplateDeploymentGeneral = "unable to create or update or delete the template deployment resource"
	// ErrCreateClient - message returned when client creation fails
	ErrCreateClient = "Unable to create client"
	// ErrCreateSession - message returned when session creation fails
	ErrCreateSession = "Unable to create OME session"
	// ErrImportDeployment - message returned when import deployment fails
	ErrImportDeployment = "Unable to import deployment"
	// ErrImportNoProfiles - message returned when import deployment fails for no existing profile
	ErrImportNoProfiles = "no deployment profiles exist for the template - %s"
)

// FailureStatusIDs - list of failure status IDs from OME for a job
var FailureStatusIDs = []any{2070, 2090, 2100, 2101, 2102, 2103}

const (
	// ValidFQDDS = Valid FQDDS supported in template creation
	ValidFQDDS string = "iDRAC,System,BIOS,NIC,LifeCycleController,RAID,EventFilters"
	// ValidTemplateViewTypes = Valid template view types supported in template creation
	ValidTemplateViewTypes string = "Deployment,Compliance"
)

// constants for Vlan attributes
const (
	// NICBondingEnabled
	NICBondingEnabled = "NIC Bonding Enabled"
	// VlanTagged
	VlanTagged = "Vlan Tagged"
	// VlanUntagged
	VlanUntagged = "Vlan UnTagged"
	// Port
	Port = "Port "
	// NICModel
	NICModel = "NICModel"
	// NICBondingTechnologyAttributeGrp
	NICBondingTechnologyAttributeGrp = "NicBondingTechnology"
	// NICBondingTechnologyAttribute
	NICBondingTechnologyAttribute = "Nic Bonding Technology"
)
