package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ConfigureBaselines to hold planned and state data
type ConfigureBaselines struct {
	ID                     types.Int64  `tfsdk:"id"`
	RefTemplateID          types.Int64  `tfsdk:"ref_template_id"`
	RefTemplateName        types.String `tfsdk:"ref_template_name"`
	Description            types.String `tfsdk:"description"`
	BaselineName           types.String `tfsdk:"baseline_name"`
	DeviceIDs              types.List   `tfsdk:"device_ids"`
	DeviceServicetags      types.List   `tfsdk:"device_servicetags"`
	ScheduleNotification   types.Bool   `tfsdk:"schedule_notification"`
	NotificationOnSchedule types.Bool   `tfsdk:"notification_on_schedule"`
	EmailAddresses         types.List   `tfsdk:"email_addresses"`
	OutputFormat           types.String `tfsdk:"output_format"`
	Cron                   types.String `tfsdk:"cron"`
	TaskID                 types.Int64  `tfsdk:"task_id"`
	JobRetryCount          types.Int64  `tfsdk:"job_retry_count"`
	SleepInterval          types.Int64  `tfsdk:"sleep_interval"`
}

// ConfigurationBaselinePayload - payload to create a baseline
type ConfigurationBaselinePayload struct {
	ID                   int64                 `json:"Id,omitempty"`
	Name                 string                `json:"Name"`
	Description          string                `json:"Description"`
	TemplateID           int64                 `json:"TemplateId"`
	BaselineTargets      []BaselineTarget      `json:"BaselineTargets"`
	NotificationSettings *NotificationSettings `json:"NotificationSettings,omitempty"`
}

// OmeBaseline - resembles the OME response of GET Baseline By ID
type OmeBaseline struct {
	ID                      int64                   `json:"Id"`
	Name                    string                  `json:"Name"`
	Description             string                  `json:"Description"`
	TemplateID              int64                   `json:"TemplateId"`
	TemplateName            string                  `json:"TemplateName"`
	TemplateType            int64                   `json:"TemplateType"`
	TaskID                  int64                   `json:"TaskId"`
	PercentageComplete      string                  `json:"PercentageComplete"`
	TaskStatus              int64                   `json:"TaskStatus"`
	LastRun                 string                  `json:"LastRun"`
	BaselineTargets         []BaselineTarget        `json:"BaselineTargets"`
	ConfigComplianceSummary ConfigComplianceSummary `json:"ConfigComplianceSummary"`
	NotificationSettings    *NotificationSettings   `json:"NotificationSettings,omitempty"`
}

// OmeBaselines contains the list of ome baseline details in the response
type OmeBaselines struct {
	Value    []OmeBaseline `json:"value"`
	NextLink string        `json:"@odata.nextLink"`
}

// ConfigComplianceSummary - holds compliance summary returned in GET Baseline by ID
type ConfigComplianceSummary struct {
	ComplianceStatus   string `json:"ComplianceStatus"`
	NumberOfCritical   int64  `json:"NumberOfCritical"`
	NumberOfWarning    int64  `json:"NumberOfWarning"`
	NumberOfNormal     int64  `json:"NumberOfNormal"`
	NumberOfIncomplete int64  `json:"NumberOfIncomplete"`
}

// BaselineTarget - contains details about baseline target device
type BaselineTarget struct {
	ID   int64              `json:"Id"`
	Type BaselineTargetType `json:"Type"`
}

// BaselineTargetType - contains details about device type
type BaselineTargetType struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name"`
}

// NotificationSettings - contains details about baseline notification settings
type NotificationSettings struct {
	NotificationType string                       `json:"NotificationType"`
	EmailAddresses   []string                     `json:"EmailAddresses"`
	Schedule         BaselineNotificationSchedule `json:"Schedule"`
	OutputFormat     string                       `json:"OutputFormat"`
}

// BaselineNotificationSchedule - contains cron expression for baseline notification schedule
type BaselineNotificationSchedule struct {
	Cron string `json:"Cron"`
}

// BaseLineIDsData holds the baseline ids
type BaseLineIDsData struct {
	BaselineIDs []int64 `json:"BaselineIds"`
}

// ConfigurationRemediation holds the plan data
type ConfigurationRemediation struct {
	ID            types.String    `tfsdk:"id"`
	BaselineName  types.String    `tfsdk:"baseline_name"`
	BaselineID    types.Int64     `tfsdk:"baseline_id"`
	TargetDevices []TargetDevices `tfsdk:"target_devices"`
	JobRetryCount types.Int64     `tfsdk:"job_retry_count"`
	SleepInterval types.Int64     `tfsdk:"sleep_interval"`
	RunLater      types.Bool      `tfsdk:"run_later"`
	Cron          types.String    `tfsdk:"cron"`
}

// TargetDevices -  holds the plan data
type TargetDevices struct {
	DeviceServiceTag types.String `tfsdk:"device_service_tag"`
	ComplianceStatus types.String `tfsdk:"compliance_status"`
}

// ConfigurationRemediationPayload - payload for remediation
type ConfigurationRemediationPayload struct {
	ID        int64       `json:"Id"`
	DeviceIDS []int64     `json:"DeviceIds"`
	Schedule  OMESchedule `json:"Schedule"`
}

// OMEDeviceComplianceReports compliance reports
type OMEDeviceComplianceReports struct {
	Value []OMEDeviceComplianceReport `json:"value"`
}

// OMEDeviceComplianceReport - reports fo devices
type OMEDeviceComplianceReport struct {
	ID                      int64                   `json:"Id"`
	DeviceName              string                  `json:"DeviceName"`
	IPAddress               string                  `json:"IpAddress"`
	IPAddresses             []string                `json:"IpAddresses"`
	Model                   string                  `json:"Model"`
	ServiceTag              string                  `json:"ServiceTag"`
	ComplianceStatus        int64                   `json:"ComplianceStatus"`
	DeviceType              int64                   `json:"DeviceType"`
	InventoryTime           string                  `json:"InventoryTime"`
	DeviceComplianceDetails DeviceComplianceDetails `json:"DeviceComplianceDetails"`
}

// DeviceComplianceDetails - details deive reports
type DeviceComplianceDetails struct {
	OdataID string `json:"@odata.id"`
}
