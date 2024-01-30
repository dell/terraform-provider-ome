/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ComplianceReport - ComplianceReport responses from api.
type ComplianceReport struct {
	Context string                   `json:"@odata.context"`
	Count   int                      `json:"@odata.count"`
	Value   []DeviceComplianceReport `json:"value"`
}

// DeviceComplianceReport struct for DeviceComplianceReport
type DeviceComplianceReport struct {
	Type                            string                      `json:"@odata.type"`
	DataID                          string                      `json:"@odata.id"`
	ComplianceStatus                *ComplianceStatusType       `json:"ComplianceStatus,omitempty"`
	ComponentComplianceReports      []ComponentComplianceReport `json:"ComponentComplianceReports,omitempty"`
	DeviceFirmwareUpdateCapable     bool                        `json:"DeviceFirmwareUpdateCapable,omitempty"`
	DeviceID                        int32                       `json:"DeviceId,omitempty"`
	DeviceModel                     string                      `json:"DeviceModel,omitempty"`
	DeviceName                      string                      `json:"DeviceName,omitempty"`
	DeviceTypeID                    int32                       `json:"DeviceTypeId,omitempty"`
	DeviceTypeName                  string                      `json:"DeviceTypeName,omitempty"`
	DeviceUserFirmwareUpdateCapable bool                        `json:"DeviceUserFirmwareUpdateCapable,omitempty"`
	FirmwareStatus                  string                      `json:"FirmwareStatus,omitempty"`
	ID                              int32                       `json:"Id,omitempty"`
	RebootRequired                  bool                        `json:"RebootRequired,omitempty"`
	ServiceTag                      string                      `json:"ServiceTag,omitempty"`
}

// ComplianceSummaryComplianceStatus struct for ComplianceSummaryComplianceStatus
type ComplianceSummaryComplianceStatus struct {
	ComplianceStatusType *ComplianceStatusType
}

// ComplianceStatusType the model 'ComplianceStatusType'
type ComplianceStatusType string

// List of ComplianceStatusType
const (
	OK        ComplianceStatusType = "OK"
	WARNING   ComplianceStatusType = "WARNING"
	CRITICAL  ComplianceStatusType = "CRITICAL"
	DOWNGRADE ComplianceStatusType = "DOWNGRADE"
	UNKNOWN   ComplianceStatusType = "UNKNOWN"
)

// ComponentComplianceReport struct for ComponentComplianceReport
type ComponentComplianceReport struct {
	Type                      string                 `json:"@odata.type"`
	ComplianceDependencies    []ComplianceDependency `json:"ComplianceDependencies,omitempty"`
	ComplianceStatus          string                 `json:"ComplianceStatus,omitempty"`
	ComponentType             string                 `json:"ComponentType,omitempty"`
	Criticality               string                 `json:"Criticality,omitempty"`
	CurrentVersion            string                 `json:"CurrentVersion,omitempty"`
	DependencyUpgradeRequired bool                   `json:"DependencyUpgradeRequired,omitempty"`
	ID                        int32                  `json:"Id,omitempty"`
	ImpactAssessment          string                 `json:"ImpactAssessment,omitempty"`
	Name                      string                 `json:"Name,omitempty"`
	Path                      string                 `json:"Path,omitempty"`
	PrerequisiteInfo          string                 `json:"PrerequisiteInfo,omitempty"`
	RebootRequired            bool                   `json:"RebootRequired,omitempty"`
	SourceName                string                 `json:"SourceName,omitempty"`
	TargetIdentifier          string                 `json:"TargetIdentifier,omitempty"`
	UniqueIdentifier          string                 `json:"UniqueIdentifier,omitempty"`
	UpdateAction              string                 `json:"UpdateAction,omitempty"`
	URI                       string                 `json:"Uri,omitempty"`
	Version                   string                 `json:"Version,omitempty"`
}

// ComplianceDependency struct for ComplianceDependency
type ComplianceDependency struct {
	Type                   string `json:"@odata.type"`
	ComplianceDependencyID int32  `json:"ComplianceDependencyId,omitempty"`
	IsHardDependency       bool   `json:"IsHardDependency,omitempty"`
	Name                   string `json:"Name,omitempty"`
	Path                   string `json:"Path,omitempty"`
	RebootRequired         bool   `json:"RebootRequired,omitempty"`
	SourceName             string `json:"SourceName,omitempty"`
	UniqueIdentifier       string `json:"UniqueIdentifier,omitempty"`
	UpdateAction           string `json:"UpdateAction,omitempty"`
	URI                    string `json:"Uri,omitempty"`
	Version                string `json:"Version,omitempty"`
}

// NullableComplianceDependencyUpdateAction struct with ComplianceDependencyUpdateAction
type NullableComplianceDependencyUpdateAction struct {
	Value *ComplianceDependencyUpdateAction
	IsSet bool
}

// ComplianceDependencyUpdateAction struct with UpdateAction
type ComplianceDependencyUpdateAction struct {
	UpdateAction *UpdateAction
}

// UpdateAction the model 'UpdateAction'
type UpdateAction string

// List of UpdateAction
const (
	UPGRADE UpdateAction = "UPGRADE"
	EQUAL   UpdateAction = "EQUAL"
)

// OmeFwComplianceReportData represents the OME Firmware Compliance Report
type OmeFwComplianceReportData struct {
	ID           types.Int64             `tfsdk:"id"`
	BaseLineName types.String            `tfsdk:"baseline_name"`
	Report       []OmeFwComplianceReport `tfsdk:"firmware_compliance_reports"`
	//filter
	CrFilter *OmeFwComplianceReportFilter `tfsdk:"filter"`
}

// OmeFwComplianceReportFilter the model 'OmeFwComplianceReportFilter'
type OmeFwComplianceReportFilter struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// OmeFwComplianceReport is the tfsdk model of OmeFwComplianceReport
type OmeFwComplianceReport struct {
	ComplianceStatus                types.String                   `tfsdk:"compliance_status"`
	ComponentComplianceReports      []OmeComponentComplianceReport `tfsdk:"component_compliance_reports"`
	DeviceFirmwareUpdateCapable     types.Bool                     `tfsdk:"device_firmware_update_capable"`
	DeviceID                        types.Int64                    `tfsdk:"device_id"`
	DeviceModel                     types.String                   `tfsdk:"device_model"`
	DeviceName                      types.String                   `tfsdk:"device_name"`
	DeviceTypeID                    types.Int64                    `tfsdk:"device_type_id"`
	DeviceTypeName                  types.String                   `tfsdk:"device_type_name"`
	DeviceUserFirmwareUpdateCapable types.Bool                     `tfsdk:"device_user_firmware_update_capable"`
	FirmwareStatus                  types.String                   `tfsdk:"firmware_status"`
	ID                              types.Int64                    `tfsdk:"id"`
	RebootRequired                  types.Bool                     `tfsdk:"reboot_required"`
	ServiceTag                      types.String                   `tfsdk:"service_tag"`
}

// OmeComponentComplianceReport is the tfsdk model of OmeComponentComplianceReport
type OmeComponentComplianceReport struct {
	//Type                      string                    `tfsdk:"@odata.type"`
	ComplianceDependencies    []OmeComplianceDependency `tfsdk:"compliance_dependencies"`
	ComplianceStatus          types.String              `tfsdk:"compliance_status"`
	ComponentType             types.String              `tfsdk:"component_type"`
	Criticality               types.String              `tfsdk:"criticality"`
	CurrentVersion            types.String              `tfsdk:"current_version"`
	DependencyUpgradeRequired types.Bool                `tfsdk:"dependency_upgrade_required"`
	ID                        types.Int64               `tfsdk:"id"`
	ImpactAssessment          types.String              `tfsdk:"impact_assessment"`
	Name                      types.String              `tfsdk:"name"`
	Path                      types.String              `tfsdk:"path"`
	PrerequisiteInfo          types.String              `tfsdk:"prerequisite_info"`
	RebootRequired            types.Bool                `tfsdk:"reboot_required"`
	SourceName                types.String              `tfsdk:"source_name"`
	TargetIdentifier          types.String              `tfsdk:"target_identifier"`
	UniqueIdentifier          types.String              `tfsdk:"unique_identifier"`
	UpdateAction              types.String              `tfsdk:"update_action"`
	URI                       types.String              `tfsdk:"uri"`
	Version                   types.String              `tfsdk:"version"`
}

// OmeComplianceDependency is the tfsdk model of OmeComplianceDependency
type OmeComplianceDependency struct {
	ComplianceDependencyID types.Int64  `tfsdk:"compliance_dependency_id"`
	IsHardDependency       types.Bool   `tfsdk:"is_hard_dependency"`
	Name                   types.String `tfsdk:"name"`
	Path                   types.String `tfsdk:"path"`
	RebootRequired         types.Bool   `tfsdk:"reboot_required"`
	SourceName             types.String `tfsdk:"source_name"`
	UniqueIdentifier       types.String `tfsdk:"unique_identifier"`
	UpdateAction           types.String `tfsdk:"update_action"`
	URI                    types.String `tfsdk:"uri"`
	Version                types.String `tfsdk:"version"`
}
