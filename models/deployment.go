package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// TemplateDeployment to hold planned and state data
type TemplateDeployment struct {
	ID                              types.String `tfsdk:"id"`
	TemplateID                      types.Int64  `tfsdk:"template_id"`
	TemplateName                    types.String `tfsdk:"template_name"`
	DeviceIDs                       types.Set    `tfsdk:"device_ids"`
	DeviceServicetags               types.Set    `tfsdk:"device_servicetags"`
	BootToNetworkISO                types.Object `tfsdk:"boot_to_network_iso"`
	DeviceAttributes                types.List   `tfsdk:"device_attributes"`
	JobRetryCount                   types.Int64  `tfsdk:"job_retry_count"`
	SleepInterval                   types.Int64  `tfsdk:"sleep_interval"`
	ForcedShutdown                  types.Bool   `tfsdk:"forced_shutdown"`
	OptionsTimeToWaitBeforeShutdown types.Int64  `tfsdk:"options_time_to_wait_before_shutdown"`
	PowerStateOff                   types.Bool   `tfsdk:"power_state_off"`
	OptionsPrecheckOnly             types.Bool   `tfsdk:"options_precheck_only"`
	OptionsStrictCheckingVlan       types.Bool   `tfsdk:"options_strict_checking_vlan"`
	OptionsContinueOnWarning        types.Bool   `tfsdk:"options_continue_on_warning"`
	RunLater                        types.Bool   `tfsdk:"run_later"`
	Cron                            types.String `tfsdk:"cron"`
}

// BootToNetworkISO to hold planned and state data for boot info
type BootToNetworkISO struct {
	BootToNetwork types.Bool   `tfsdk:"boot_to_network"`
	ShareType     types.String `tfsdk:"share_type"`
	IsoPath       types.String `tfsdk:"iso_path"`
	IsoTimeout    types.Int64  `tfsdk:"iso_timeout"`
	ShareDetail   types.Object `tfsdk:"share_detail"`
}

// ShareDetail to hold planned and state data for share details
type ShareDetail struct {
	IPAddress types.String `tfsdk:"ip_address"`
	ShareName types.String `tfsdk:"share_name"`
	WorkGroup types.String `tfsdk:"work_group"`
	User      types.String `tfsdk:"user"`
	Password  types.String `tfsdk:"password"`
}

// Options to hold planned and state data for power opttions
type Options struct {
	ShutdownType             types.Int64 `tfsdk:"shutdown_type"`
	TimeToWaitBeforeShutdown types.Int64 `tfsdk:"time_to_wait_before_shutdown"`
	EndHostPowerState        types.Int64 `tfsdk:"end_host_power_state"`
	PrecheckOnly             types.Bool  `tfsdk:"precheck_only"`
	ContinueOnWarning        types.Bool  `tfsdk:"continue_on_warning"`
	StrictCheckingVlan       types.Bool  `tfsdk:"strict_checking_vlan"`
}

// Schedule to hold planned and state data deploy schedule
type Schedule struct {
	RunNow    types.Bool   `tfsdk:"run_now"`
	RunLater  types.Bool   `tfsdk:"run_later"`
	Cron      types.String `tfsdk:"cron"`
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
}

// DeviceAttributes to hold planned and state data
type DeviceAttributes struct {
	DeviceServiceTags types.Set  `tfsdk:"device_servicetags"`
	Attributes        types.List `tfsdk:"attributes"`
}

// OMETemplateDeployRequest to form a request to deploy template
type OMETemplateDeployRequest struct {
	ID                  int64                  `json:"Id"`
	TargetIDS           []int64                `json:"TargetIds"`
	Schedule            OMESchedule            `json:"Schedule"`
	Attributes          []OMEDeviceAttributes  `json:"Attributes"`
	Options             OMEOptions             `json:"Options"`
	NetworkBootISOModel OMENetworkBootISOModel `json:"NetworkBootIsoModel"`
}

// OMEDeviceAttributes to form a request to deploy template
type OMEDeviceAttributes struct {
	DeviceID   int64          `json:"DeviceId"`
	Attributes []OMEAttribute `json:"Attributes"`
}

// OMEAttribute to form a request to deploy template
type OMEAttribute struct {
	ID        int64  `json:"Id"`
	Value     string `json:"Value"`
	IsIgnored bool   `json:"IsIgnored"`
}

// OMENetworkBootISOModel to form a request to deploy template
type OMENetworkBootISOModel struct {
	BootToNetwork  bool           `json:"BootToNetwork"`
	ShareType      string         `json:"ShareType"`
	ISOPath        string         `json:"IsoPath"`
	ISOTimeout     int64          `json:"IsoTimeout"`
	ISOTimeoutUnit int64          `json:"IsoTimeoutUnit"`
	ShareDetail    OMEShareDetail `json:"ShareDetail"`
}

// OMEShareDetail to form a request to deploy template
type OMEShareDetail struct {
	IPAddress string `json:"IpAddress"`
	ShareName string `json:"ShareName"`
	WorkGroup string `json:"WorkGroup"`
	User      string `json:"User"`
	Password  string `json:"Password"`
}

// OMEOptions to form a request to deploy template
type OMEOptions struct {
	ShutdownType             int64 `json:"ShutdownType"`
	TimeToWaitBeforeShutdown int64 `json:"TimeToWaitBeforeShutdown"`
	EndHostPowerState        int64 `json:"EndHostPowerState"`
	PrecheckOnly             bool  `json:"PrecheckOnly"`
	ContinueOnWarning        bool  `json:"ContinueOnWarning"`
	StrictCheckingVLAN       bool  `json:"StrictCheckingVlan"`
}

// OMESchedule to form a request to deploy template
type OMESchedule struct {
	RunNow    bool   `json:"RunNow"`
	RunLater  bool   `json:"RunLater"`
	Cron      string `json:"Cron"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
}

// OMEServerProfiles to form list of profiles.
type OMEServerProfiles struct {
	Value []OMEServerProfile `json:"value"`
}

// OMEServerProfile to form a profile type object.
type OMEServerProfile struct {
	ID           int64  `json:"Id,omitempty"`
	ProfileName  string `json:"ProfileName,omitempty"`
	TemplateID   int64  `json:"TemplateId,omitempty"`
	TemplateName string `json:"TemplateName,omitempty"`
	TargetID     int64  `json:"TargetId,omitempty"`
}

// ProfileDeleteRequest to delete profiles type request.
type ProfileDeleteRequest struct {
	ProfileIds []int64 `json:"ProfileIds"`
}
