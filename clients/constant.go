/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://mozilla.org/MPL/2.0/
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	// RunningStatusID - job success status ID
	RunningStatusID = 2050
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
	// BaselineAPI - api used to manage baselines
	BaselineAPI = "/api/TemplateService/Baselines"
	// BaselineByIDAPI - api used to manage baseline by ID
	BaselineByIDAPI = "/api/TemplateService/Baselines(%d)"
	// BaselineDeviceComplianceReportsAPI - api used to fetch baseline device compliance report
	BaselineDeviceComplianceReportsAPI = "/api/TemplateService/Baselines(%d)/DeviceConfigComplianceReports"
	// BaselineDeviceAttrComplianceReportsAPI - api used to fetch baseline device attributes compliance report
	BaselineDeviceAttrComplianceReportsAPI = "/api/TemplateService/Baselines(%d)/DeviceConfigComplianceReports(%d)/DeviceComplianceDetails"
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
	// GroupAPI - api for fetching group details
	GroupAPI = "/api/GroupService/Groups"
	// GroupServiceAPI - api to get and delete group by id
	GroupServiceAPI = GroupAPI + "(%d)"
	// GroupServiceActionsAPI - api to create and modify a group
	GroupServiceActionsAPI = "/api/GroupService/Actions/GroupService.%sGroup"
	// GroupServiceDeviceActionsAPI - api to add and remove devices from a group
	GroupServiceDeviceActionsAPI = "/api/GroupService/Actions/GroupService.%sMemberDevices"
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
	//BaseLineRemoveAPI - api to remove a baseline
	BaseLineRemoveAPI = "/api/TemplateService/Actions/TemplateService.RemoveBaseline"
	//BaseLineConfigRemediationAPI - api to remediate a baseline
	BaseLineConfigRemediationAPI = "/api/TemplateService/Actions/TemplateService.Remediate"
	//BaseLineConfigDeviceCompReport - api to get device compliance report of a baseline
	BaseLineConfigDeviceCompReport = "/api/TemplateService/Baselines(%d)/DeviceConfigComplianceReports"
	//VlanNetworksAPI - api to vlan networks
	VlanNetworksAPI = "/api/NetworkConfigurationService/Networks"
	//ImportTemplateAPI - api to import a template
	ImportTemplateAPI = "/api/TemplateService/Actions/TemplateService.Import"
	// TemplateNameContainsAPI - api to fetch templates by name
	TemplateNameContainsAPI = "/api/TemplateService/Templates?$filter=contains(Name, '%s')"
	//UserAPI - api to manage users
	UserAPI = "/api/AccountService/Accounts"
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
	// ErrComplianceTemplateIDOrName - error message when either compliance template ID or name is not given
	ErrComplianceTemplateIDOrName = "either compliance template ID or template name is expected"
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
	// ErrScheduleNotification - message returned when email address are not provided when schedule notification is true
	ErrScheduleNotification = "please provide a valid email address, when schedule notification is set to true"
	// ErrGnrCreateBaseline - summary returned when failed to create baseline
	ErrGnrCreateBaseline = "error creating a baseline"
	// ErrGnrUpdateBaseline - summary returned when failed to update baseline
	ErrGnrUpdateBaseline = "error updating a baseline"
	// ErrGnrDeleteBaseline - summary returned when failed to delete baseline
	ErrGnrDeleteBaseline = "error deleting a baseline"
	// ErrGnrReadBaseline - summary returned when failed to read baseline
	ErrGnrReadBaseline = "error reading a baseline"
	// ErrGnrImportBaseline - message returned when import baseline fails
	ErrGnrImportBaseline = "Unable to import baseline"
	// ErrInvalidEmailAddress - message returned when invalid email address is provided
	ErrInvalidEmailAddress = "invalid email address %s"
	// ErrInvalidCronExpression - message returned when invalid cron expression is not provided
	ErrInvalidCronExpression = "cron value is required when notification is scheduled"
	// ErrInvalidRefTemplateNameorID - refernce template name required
	ErrInvalidRefTemplateNameorID = "either reference template Id or template name is required"
	// ErrInvalidRefTemplateType - template type is not comliance
	ErrInvalidRefTemplateType = "reference template id or name should be of type compliance"
	// ErrBaselineCreationTask - baseline report generation task failed
	ErrBaselineCreationTask = "baseline report generation task failed"
	// ErrDeviceNotCapable - device capablity
	ErrDeviceNotCapable = "device with %v are not capable for deployments"
	// ErrBaseLineJobIsRunning - device capablity
	ErrBaseLineJobIsRunning = "job with id %d is already running please wait for sometime and try again"
	// WarningBaselineDeviceCapability - message returned when create baseline has incompatible devices
	WarningBaselineDeviceCapability = "%v devices are not valid to create baseline"
	// ErrBaselineNameNotFound - message returned when provided baseline name does not exist
	ErrBaselineNameNotFound = "baseline not found: %s"
	// ErrGnrBaseLineCreateRemediation - message returned when there is a error in baseline remediation for configuration
	ErrGnrBaseLineCreateRemediation = "baseline configuration remediation create error"
	// ErrGnrBaseLineRemediation - message returned when there is a error in baseline remediation for configuration
	ErrGnrBaseLineReadRemediation = "baseline configuration remediation read error"
	// ErrGnrBaseLineRemediation - message returned when there is a baseline report generation in progress
	ErrBaseLineReportInProgress = "inventory update is in progress, retry after some time"
	// ErrBaseLineInvalidDevices - message returned when baseline has invalid devices
	ErrBaseLineInvalidDevices = "devices %v are not part of a baseline"
	// ErrBaseLineInvalid - message returned when baseline name or id is invalid
	ErrBaseLineInvalid = "either baseline name or id is required"
	// ErrBaseLineUpdateRemediation - message returned when there is a error in baseline remediation for configuration
	ErrBaseLineUpdateRemediation = "baseline configuration remediation update error"
	// ErrBaseLineUpdateRemediation - message returned when baseline name or id is changed
	ErrBaseLineModified = "baseline name or id cannot be modified"
	// ErrBaseLineTargetsSize - message returned when min length is not satisfied
	ErrBaseLineTargetsSize = "list must contain at least %d elements"
	// ErrBaseLineTargetsSize - message returned when min length is not satisfied
	ErrBaseLineComplianceStatus = "supported value is %s"
	// ErrBaselineReportForDevice - message returned when device report is not avaiable for a servicetag
	ErrBaselineReportForDevice = "device reports not found for baseline %d and device %s"
	// ErrCreateTemplate - message returned when template creation fails
	// ErrInvalidTemplateDeviceType - error message for invalid template device type
	ErrInvalidTemplateDeviceType = "Invalid template device type for template creation"
	// ErrTemplateDeploymentCreate - message returned when template deployment fails
	ErrTemplateDeploymentCreate = "unable to create template deployment resource"
	// ErrTemplateDeploymentUpdate - message returned when template deployment fails
	ErrTemplateDeploymentUpdate = "unable to update template deployment resource"
	// ErrTemplateDeploymentRead - message returned when template deployment fails
	ErrTemplateDeploymentRead = "unable to read template deployment resource"
	// ErrTemplateDeploymentDelete - message returned when template deployment fails
	ErrTemplateDeploymentDelete = "unable to delete template deployment resource"
	// ErrCreateTemplate - message returned when template creation fails
	ErrCreateTemplate = "Unable to create template"
	// ErrReadTemplate - message returned when template read fails
	ErrReadTemplate = "Unable to read template"
	// ErrUpdateTemplate - message returned when update template fails
	ErrUpdateTemplate = "Unable to update template"
	// ErrDeleteTemplate - message returned when template delete fails
	ErrDeleteTemplate = "Unable to delete template"
	// ErrImportTemplate - message returned when import template fails
	ErrImportTemplate = "Unable to import template"
	// ErrGnrConfigurationReport - message returned when report could not be fetched
	ErrGnrConfigurationReport = "unable to fetch the report"
	// ErrCronRequired - message returned when run_later is true but cron is not provided
	ErrCronRequired = "cron is required when run_later is true"
	//ErrBaseLineScheduleValid
	ErrBaseLineScheduleValid = "attributes `cron` and `email_addresses` are accepted only when `schedule` is true"
	//ErrBaseLineNotifyValid
	ErrBaseLineNotifyValid = "attributes `cron` is not accepted only when `schedule` is true and `notify_on_schedule` is false"
	// ErrGnrCreateUser - summary returned when failed to create User
	ErrGnrCreateUser = "error creating a User"
	// ErrGnrUpdateUser - summary returned when failed to update User
	ErrGnrUpdateUser = "error updating a User"
	// ErrGnrDeleteUser - summary returned when failed to delete User
	ErrGnrDeleteUser = "error deleting a User"
	// ErrGnrReadUser - summary returned when failed to read User
	ErrGnrReadUser = "error reading a User"
	// ErrGnrImportUser - message returned when import User fails
	ErrGnrImportUser = "Unable to import User"
)

// FailureStatusIDs - list of failure status IDs from OME for a job
var FailureStatusIDs = []any{2070, 2090, 2100, 2101, 2102, 2103}

const (
	// ValidFQDDS = Valid FQDDS supported in template creation
	ValidFQDDS string = "All,iDRAC,System,BIOS,NIC,LifeCycleController,RAID,EventFilters"
	// ValidOutputFormat - valid output formats
	ValidOutputFormat string = "html,csv,pdf,xls"
	// ValidTemplateViewTypes = Valid template view types supported in template creation
	ValidTemplateViewTypes string = "Deployment,Compliance"
	// ValidComplainceStatus = Valid compliance status supported
	ValidComplainceStatus string = "Compliant"
	// ValidTemplateDeviceTypes = Valid template device types supported in template creation
	ValidTemplateDeviceTypes string = "Server,Chassis"
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
