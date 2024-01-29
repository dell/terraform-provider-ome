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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// FirmwareBaselinesModel struct for FirmwareBaselinesModel
type FirmwareBaselinesModel struct {
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
	TaskID       *int64 `json:"TaskId,omitempty"`
	TaskStatusID *int32 `json:"TaskStatusId,omitempty"`
}

// ComplianceSummaryModel struct for ComplianceSummaryModel
type ComplianceSummaryModel struct {
	ComplianceStatus  *string `json:"ComplianceStatus,omitempty"`
	NumberOfCritical  *int64  `json:"NumberOfCritical,omitempty"`
	NumberOfNormal    *int64  `json:"NumberOfNormal,omitempty"`
	NumberOfWarning   *int64  `json:"NumberOfWarning,omitempty"`
	NumberOfDowngrade *int64  `json:"NumberOfDowngrade"`
	NumberOfUnknown   *int64  `json:"NumberOfUnknown"`
}

// TargetModel struct for TargetModel
type TargetModel struct {
	// DeviceIDs can be determined through /api/DeviceService/Devices and GroupIDs can be determined through /api/GroupService/Groups
	ID   int32           `json:"Id"`
	Type TargetTypeModel `json:"Type"`
}

// TargetTypeModel struct for TargetTypeModel
type TargetTypeModel struct {
	// Device type ID. DeviceType IDs can be determined through /api/DeviceService/DeviceType
	ID int32 `json:"Id"`
	// Type of the target (DEVICE or GROUP)
	Name string `json:"Name"`
}

// FirmwareBaselineResource represents the Firmware Baseline resource model
type FirmwareBaselineResource struct {
	CatalogID              types.Int64  `tfsdk:"catalog_id"`
	ComplianceSummary      types.Object `tfsdk:"compliance_summary"`
	Description            types.String `tfsdk:"description"`
	DowngradeEnabled       types.Bool   `tfsdk:"downgrade_enabled"`
	FilterNoRebootRequired types.Bool   `tfsdk:"filter_no_reboot_required"`
	ID                     types.Int64  `tfsdk:"id"`
	Is64Bit                types.Bool   `tfsdk:"is_64_bit"`
	LastRun                types.String `tfsdk:"last_run"`
	Name                   types.String `tfsdk:"name"`
	RepositoryID           types.Int64  `tfsdk:"repository_id"`
	RepositoryName         types.String `tfsdk:"repository_name"`
	RepositoryType         types.String `tfsdk:"repository_type"`
	Targets                types.List   `tfsdk:"targets"`
	TaskID                 types.Int64  `tfsdk:"task_id"`
	TaskStatus             types.String `tfsdk:"task_status"`
	CatalogName            types.String `tfsdk:"catalog_name"`
	DeviceNames            types.List   `tfsdk:"device_names"`
	DeviceServiceTags      types.List   `tfsdk:"device_service_tags"`
	GroupNames             types.List   `tfsdk:"group_names"`
}

// CreateUpdateFirmwareBaseline - payload to create/update a firmware baseline
type CreateUpdateFirmwareBaseline struct {
	ID                     int64         `json:"Id"`
	CatalogID              int64         `json:"CatalogId"`
	Description            string        `json:"Description"`
	FilterNoRebootRequired bool          `json:"FilterNoRebootRequired"`
	Is64Bit                bool          `json:"Is64Bit"`
	Name                   string        `json:"Name"`
	RepositoryID           int64         `json:"RepositoryId"`
	Targets                []TargetModel `json:"Targets"`
}
