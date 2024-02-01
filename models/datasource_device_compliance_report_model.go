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

// DeviceComplianceReportModel - Model for Device Compliance Report API
type DeviceComplianceReportModel struct {
	Context string                  `json:"@odata.context"`
	Count   int                     `json:"@odata.count"`
	Value   []DeviceComplianceModel `json:"value"`
}

// DeviceComplianceModel - Model for Device Compliance
type DeviceComplianceModel struct {
	ID int64 `json:"Id"`

	DeviceID int64 `json:"DeviceId"`

	ServiceTag string `json:"ServiceTag"`

	DeviceModel string `json:"DeviceModel"`

	DeviceTypeName string `json:"DeviceTypeName"`

	DeviceName string `json:"DeviceName"`

	FirmwareStatus string `json:"FirmwareStatus"`

	ComplianceStatus string `json:"ComplianceStatus"`

	DeviceTypeID int64 `json:"DeviceTypeId"`

	RebootRequired bool `json:"RebootRequired"`

	DeviceFirmwareUpdateCapable bool `json:"DeviceFirmwareUpdateCapable"`

	DeviceUserFirmwareUpdateCapable bool `json:"DeviceUserFirmwareUpdateCapable"`

	ComponentComplianceReports []ComponentComplianceReportModel `json:"ComponentComplianceReports"`
}

// ComplianceDependenciesModel - Model for Compliance Dependencies
type ComplianceDependenciesModel struct {
	ComplianceDependencyID int64  `json:"ComplianceDependencyId"`
	IsHardDependency       bool   `json:"IsHardDependency"`
	Name                   string `json:"Name"`
	Path                   string `json:"Path"`
	RebootRequired         bool   `json:"RebootRequired"`
	SourceName             string `json:"SourceName"`
	UniqueIdentifier       string `json:"UniqueIdentifier"`
	UpdateAction           string `json:"UpdateAction"`
	URI                    string `json:"Uri"`
	Version                string `json:"Version"`
}

// ComponentComplianceReportModel - Model for Component Compliance Report
type ComponentComplianceReportModel struct {
	ID                        int64                         `json:"Id"`
	Version                   string                        `json:"Version"`
	CurrentVersion            string                        `json:"CurrentVersion"`
	Path                      string                        `json:"Path"`
	Name                      string                        `json:"Name"`
	Criticality               string                        `json:"Criticality"`
	UniqueIdentifier          string                        `json:"UniqueIdentifier"`
	TargetIdentifier          string                        `json:"TargetIdentifier"`
	UpdateAction              string                        `json:"UpdateAction"`
	SourceName                string                        `json:"SourceName"`
	PrerequisiteInfo          string                        `json:"PrerequisiteInfo"`
	ImpactAssessment          string                        `json:"ImpactAssessment"`
	URI                       string                        `json:"Uri"`
	RebootRequired            bool                          `json:"RebootRequired"`
	ComplianceStatus          string                        `json:"ComplianceStatus"`
	ComplianceDependencies    []ComplianceDependenciesModel `json:"ComplianceDependencies"`
	ComponentType             string                        `json:"ComponentType"`
	DependencyUpgradeRequired bool                          `json:"DependencyUpgradeRequired"`
}

// OMEDeviceComplianceData represents the OME Device Compliance
type OMEDeviceComplianceData struct {
	ID                types.Int64            `tfsdk:"id"`
	Reports           []DeviceComplianceData `tfsdk:"device_compliance_reports"`
	DeviceIDs         types.List             `tfsdk:"device_ids"`
	DeviceServiceTags types.List             `tfsdk:"device_service_tags"`
	DeviceGroupNames  types.List             `tfsdk:"device_group_names"`
}

// DeviceComplianceData - The representation of Device Compliance
type DeviceComplianceData struct {
	ComplianceStatus                types.String                     `tfsdk:"compliance_status"`
	ComponentComplianceReport       []ComponentComplianceReportsData `tfsdk:"component_compliance_reports"`
	DeviceID                        types.Int64                      `tfsdk:"device_id"`
	DeviceModel                     types.String                     `tfsdk:"device_model"`
	DeviceName                      types.String                     `tfsdk:"device_name"`
	DeviceTypeID                    types.Int64                      `tfsdk:"device_type_id"`
	DeviceTypeName                  types.String                     `tfsdk:"device_type_name"`
	FirmwareStatus                  types.String                     `tfsdk:"firmware_status"`
	ID                              types.Int64                      `tfsdk:"id"`
	RebootRequired                  types.Bool                       `tfsdk:"reboot_required"`
	ServiceTag                      types.String                     `tfsdk:"service_tag"`
	DeviceFirmwareUpdateCapable     types.Bool                       `tfsdk:"device_firmware_update_capable"`
	DeviceUserFirmwareUpdateCapable types.Bool                       `tfsdk:"device_user_firmware_update_capable"`
}

// DeviceComplianceDatasource - The representation of Device Compliance
type DeviceComplianceDatasource struct {
	CatalogID               types.Int64            `tfsdk:"catalog_id"`
	ComplianceSummary       ComplianceSummary      `tfsdk:"compliance_summary"`
	Description             types.String           `tfsdk:"description"`
	DowngradeEnabled        types.Bool             `tfsdk:"downgrade_enabled"`
	FilterNoRebootRequired  types.Bool             `tfsdk:"filter_no_reboot_required"`
	ID                      types.Int64            `tfsdk:"id"`
	Is64Bit                 types.Bool             `tfsdk:"is_64_bit"`
	LastRun                 types.String           `tfsdk:"last_run"`
	Name                    types.String           `tfsdk:"name"`
	RepositoryID            types.Int64            `tfsdk:"repository_id"`
	RepositoryName          types.String           `tfsdk:"repository_name"`
	RepositoryType          types.String           `tfsdk:"repository_type"`
	Targets                 Targets                `tfsdk:"targets"`
	TaskID                  types.Int64            `tfsdk:"task_id"`
	TaskStatus              types.String           `tfsdk:"task_status"`
	CatalogName             types.String           `tfsdk:"catalog_name"`
	DeviceNames             types.List             `tfsdk:"device_names"`
	DeviceServiceTags       types.List             `tfsdk:"device_service_tags"`
	GroupNames              types.List             `tfsdk:"group_names"`
	DeviceComplianceReports []DeviceComplianceData `tfsdk:"device_compliance_reports"`
}

// ComplianceSummary struct for ComplianceSummary
type ComplianceSummary struct {
	ComplianceStatus  types.String `tfsdk:"compliance_status"`
	NumberOfCritical  types.Int64  `tfsdk:"number_of_critical"`
	NumberOfDowngrade types.Int64  `tfsdk:"number_of_downgrade"`
	NumberOfNormal    types.Int64  `tfsdk:"number_of_normal"`
	NumberOfWarning   types.Int64  `tfsdk:"number_of_warning"`
	NumberOfUnknown   types.Int64  `tfsdk:"number_of_unknown"`
}

// Targets struct for Targets
type Targets struct {
	// DeviceIDs can be determined through /api/DeviceService/Devices and GroupIDs can be determined through /api/GroupService/Groups
	ID   types.Int64        `tfsdk:"id"`
	Type TargetTypeBaseline `tfsdk:"type"`
}

// TargetTypeBaseline struct for TargetTypeModel
type TargetTypeBaseline struct {
	// Device type ID. DeviceType IDs can be determined through /api/DeviceService/DeviceType
	ID types.Int64 `tfsdk:"id"`
	// Type of the target (DEVICE or GROUP)
	Name string `tfsdk:"name"`
}

// FirmwareBaselinesGetModel - The representation of Firmware Baseline
type FirmwareBaselinesGetModel struct {
	// ID of the catalog. Users must enumerate all catalogs and match the 'Name' of the repository with the input provided while creating the catalog
	CatalogID         int64                  `json:"CatalogId"`
	ComplianceSummary ComplianceSummaryModel `json:"ComplianceSummary,omitempty"`
	// Description of the baseline
	Description *string `json:"Description,omitempty"`
	// Indicates if the firmware can be downgraded
	DowngradeEnabled bool `json:"DowngradeEnabled"`
	// Filters applicable updates where no reboot is required during create baseline for firmware updates.
	FilterNoRebootRequired *bool `json:"FilterNoRebootRequired,omitempty"`
	// ID of the baseline. For PUT is required and for POST is optional.
	ID *int32 `json:"Id,omitempty"`
	// This must always be set to true. The size of the DUP files used is 64 bits.
	Is64Bit bool    `json:"Is64Bit"`
	LastRun *string `json:"LastRun,omitempty"`
	// Name of the baseline
	Name string `json:"Name"`
	// ID of the repository. Derived from the catalog response
	RepositoryID   int64   `json:"RepositoryId"`
	RepositoryName *string `json:"RepositoryName,omitempty"`
	RepositoryType *string `json:"RepositoryType,omitempty"`
	// The DeviceID, if the baseline is being created for devices or, the GroupID, if the baseline is being created for a group of devices.
	Targets []TargetModel `json:"Targets"`
	// Identifier of task which created this baseline.
	TaskID                 *int64                  `json:"TaskId,omitempty"`
	TaskStatusID           *int32                  `json:"TaskStatusId,omitempty"`
	DeviceComplianceReport []DeviceComplianceModel `json:"DeviceComplianceReports,omitempty"`
}

// ComplianceDependencies - The representation of Compliance Dependencies
type ComplianceDependencies struct {
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

// ComponentComplianceReportsData - The representation of Component Compliance Report
type ComponentComplianceReportsData struct {
	ID                        types.Int64              `tfsdk:"id"`
	Version                   types.String             `tfsdk:"version"`
	CurrentVersion            types.String             `tfsdk:"current_version"`
	Path                      types.String             `tfsdk:"path"`
	Name                      types.String             `tfsdk:"name"`
	Criticality               types.String             `tfsdk:"criticality"`
	UniqueIdentifier          types.String             `tfsdk:"unique_identifier"`
	TargetIdentifier          types.String             `tfsdk:"target_identifier"`
	UpdateAction              types.String             `tfsdk:"update_action"`
	SourceName                types.String             `tfsdk:"source_name"`
	PrerequisiteInfo          types.String             `tfsdk:"prerequisite_info"`
	ImpactAssessment          types.String             `tfsdk:"impact_assessment"`
	URI                       types.String             `tfsdk:"uri"`
	RebootRequired            types.Bool               `tfsdk:"reboot_required"`
	ComplianceStatus          types.String             `tfsdk:"compliance_status"`
	ComplianceDependencies    []ComplianceDependencies `tfsdk:"compliance_dependencies"`
	ComponentType             types.String             `tfsdk:"component_type"`
	DependencyUpgradeRequired types.Bool               `tfsdk:"dependency_upgrade_required"`
}
