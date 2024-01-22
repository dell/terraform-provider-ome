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

package ome

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// FwBaseComplianceReportDataSchema returns a map of string to schema.Attribute
// that represents the schema for the Firmware Base Compliance Report Data.
//
// It includes the following attributes:
//   - "id": A string attribute that represents the compliance status.
//   - "firmware_compliance_reports": A list nested attribute that represents
//     the firmware baseline compliance reports. It is computed and contains
//     attributes from the FwBaseComplianceReportSchema.
//
// The function does not take any parameters and returns a map[string]schema.Attribute.
func FwBaseComplianceReportDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "Compliance Status",
			Description:         "Compliance Status",
			Computed:            true,
		},
		"baseline_name": schema.StringAttribute{
			MarkdownDescription: "Name of the Baseline.",
			Description:         "Name of the Baseline.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"firmware_compliance_reports": schema.ListNestedAttribute{
			MarkdownDescription: "Firmware Baseline Compliance Reports",
			Description:         "Firmware Baseline Compliance Reports",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: FwBaseComplianceReportSchema()},
		},
	}
}

// FwBaseComplianceReportSchema is a function that returns the schema for OmeFwBaseComplianceReportSchema
func FwBaseComplianceReportSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"compliance_status": schema.StringAttribute{
			MarkdownDescription: "Compliance Status",
			Description:         "Compliance Status",
			Computed:            true,
		},
		"component_compliance_reports": schema.ListNestedAttribute{
			MarkdownDescription: "Component Compliance Reports",
			Description:         "Component Compliance Reports",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ComponentComplianceReportSchema()},
		},
		"device_firmware_update_capable": schema.BoolAttribute{
			MarkdownDescription: "Device Firmware Update Capable",
			Description:         "Device Firmware Update Capable",
			Computed:            true,
		},
		"device_id": schema.Int64Attribute{
			MarkdownDescription: "Device ID",
			Description:         "Device ID",
			Computed:            true,
		},
		"device_model": schema.StringAttribute{
			MarkdownDescription: "Device Model",
			Description:         "Device Model",
			Computed:            true,
		},
		"device_name": schema.StringAttribute{
			MarkdownDescription: "Device Name",
			Description:         "Device Name",
			Computed:            true,
		},
		"device_type_id": schema.Int64Attribute{
			MarkdownDescription: "Device Type ID",
			Description:         "Device Type ID",
			Computed:            true,
		},
		"device_type_name": schema.StringAttribute{
			MarkdownDescription: "Device Type Name",
			Description:         "Device Type Name",
			Computed:            true,
		},
		"device_user_firmware_update_capable": schema.BoolAttribute{
			MarkdownDescription: "Device User Firmware Update Capable",
			Description:         "Device User Firmware Update Capable",
			Computed:            true,
		},
		"firmware_status": schema.StringAttribute{
			MarkdownDescription: "Firmware Status",
			Description:         "Firmware Status",
			Computed:            true,
		},
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},
		"reboot_required": schema.BoolAttribute{
			MarkdownDescription: "Reboot Required",
			Description:         "Reboot Required",
			Computed:            true,
		},
		"service_tag": schema.StringAttribute{
			MarkdownDescription: "Service Tag",
			Description:         "Service Tag",
			Computed:            true,
		},
	}
}

// ComponentComplianceReportSchema is a function that returns the schema for OmeComponentComplianceReport
func ComponentComplianceReportSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"compliance_dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "Compliance Dependencies",
			Description:         "Compliance Dependencies",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ComplianceDependencySchema()},
		},
		"compliance_status": schema.StringAttribute{
			MarkdownDescription: "Compliance Status",
			Description:         "Compliance Status",
			Computed:            true,
		},
		"component_type": schema.StringAttribute{
			MarkdownDescription: "Component Type",
			Description:         "Component Type",
			Computed:            true,
		},
		"criticality": schema.StringAttribute{
			MarkdownDescription: "Criticality",
			Description:         "Criticality",
			Computed:            true,
		},
		"current_version": schema.StringAttribute{
			MarkdownDescription: "Current Version",
			Description:         "Current Version",
			Computed:            true,
		},
		"dependency_upgrade_required": schema.BoolAttribute{
			MarkdownDescription: "Dependency Upgrade Required",
			Description:         "Dependency Upgrade Required",
			Computed:            true,
		},
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},
		"impact_assessment": schema.StringAttribute{
			MarkdownDescription: "Impact Assessment",
			Description:         "Impact Assessment",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "Path",
			Description:         "Path",
			Computed:            true,
		},
		"prerequisite_info": schema.StringAttribute{
			MarkdownDescription: "Prerequisite Info",
			Description:         "Prerequisite Info",
			Computed:            true,
		},
		"reboot_required": schema.BoolAttribute{
			MarkdownDescription: "Reboot Required",
			Description:         "Reboot Required",
			Computed:            true,
		},
		"source_name": schema.StringAttribute{
			MarkdownDescription: "Source Name",
			Description:         "Source Name",
			Computed:            true,
		},
		"target_identifier": schema.StringAttribute{
			MarkdownDescription: "Target Identifier",
			Description:         "Target Identifier",
			Computed:            true,
		},
		"unique_identifier": schema.StringAttribute{
			MarkdownDescription: "Unique Identifier",
			Description:         "Unique Identifier",
			Computed:            true,
		},
		"update_action": schema.StringAttribute{
			MarkdownDescription: "Update Action",
			Description:         "Update Action",
			Computed:            true,
		},
		"uri": schema.StringAttribute{
			MarkdownDescription: "URI",
			Description:         "URI",
			Computed:            true,
		},
		"version": schema.StringAttribute{
			MarkdownDescription: "Version",
			Description:         "Version",
			Computed:            true,
		},
	}
}

// ComplianceDependencySchema is a function that returns the schema for OmeComplianceDependency
func ComplianceDependencySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"compliance_dependency_id": schema.StringAttribute{
			MarkdownDescription: "Compliance Dependency Id",
			Description:         "Compliance Dependency Id",
			Computed:            true,
		},
		"is_hard_dependency": schema.BoolAttribute{
			MarkdownDescription: "Is Hard Dependency",
			Description:         "Is Hard Dependency",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "Path",
			Description:         "Path",
			Computed:            true,
		},
		"reboot_required": schema.BoolAttribute{
			MarkdownDescription: "Reboot Required",
			Description:         "Reboot Required",
			Computed:            true,
		},
		"source_name": schema.StringAttribute{
			MarkdownDescription: "Source Name",
			Description:         "Source Name",
			Computed:            true,
		},
		"unique_identifier": schema.StringAttribute{
			MarkdownDescription: "Unique Identifier",
			Description:         "Unique Identifier",
			Computed:            true,
		},
		"update_action": schema.StringAttribute{
			MarkdownDescription: "Update Action",
			Description:         "Update Action",
			Computed:            true,
		},
		"uri": schema.StringAttribute{
			MarkdownDescription: "Uri",
			Description:         "Uri",
			Computed:            true,
		},
		"version": schema.StringAttribute{
			MarkdownDescription: "Version",
			Description:         "Version",
			Computed:            true,
		},
	}
}

// FwBaseComplianceReportFilterBlockSchema is a function that returns the block schema for OmeFwBaseComplianceReportFilter
func FwBaseComplianceReportFilterBlockSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("DeviceName", "DeviceModel", "ServiceTag"),
					},
					MarkdownDescription: "Firmware Baseline Compliance Reports with filter key and value pair. Supported filter keys are: DeviceName, DeviceModel, ServiceTag",
					Description:         "Firmware Baseline Compliance Reports with filter key and value pair. Supported filters are: DeviceName, DeviceModel, ServiceTag",
				},
				"value": schema.StringAttribute{
					Optional:            true,
					Description:         "The value for the filter key.",
					MarkdownDescription: "The value for the filter key",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
			},
		},
	}
}
