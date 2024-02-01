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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func omeDeviceReportSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "Dummy ID of the datasource.",
			Description:         "Dummy ID of the datasource.",
			Computed:            true,
		},
		"device_compliance_reports": schema.ListNestedAttribute{
			MarkdownDescription: "Reports fetched.",
			Description:         "Reports fetched.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeSingleDeviceComplianceReportDataSchema()},
		},
		"device_ids": schema.ListAttribute{
			MarkdownDescription: "Device Ids is the list of device ids that you want get the compliance report for. One of device_ids or device_service_tags or device_group_names is required",
			Description:         "Device Ids is the list of device ids that you want get the compliance report for. One of device_ids or device_service_tags or device_group_names is required",
			ElementType:         types.Int64Type,
			Optional:            true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_service_tags")),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_group_names")),
				listvalidator.SizeAtLeast(1),
				// listvalidator.ValueStringsAre(
				// 	int64validator.LengthAtLeast(1),
				// ),
			},
		},
		"device_service_tags": schema.ListAttribute{
			MarkdownDescription: "Device service tags is the list of device service tags that you want get the compliance report for. One of device_ids or device_service_tags or device_group_names is required",
			Description:         "Device service tags is the list of device service tags that you want get the compliance report for. One of device_ids or device_service_tags or device_group_names is required",
			ElementType:         types.StringType,
			Optional:            true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_ids")),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_group_names")),
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
				),
			},
		},
		"device_group_names": schema.ListAttribute{
			MarkdownDescription: "Group names is the list of group names that you want to that you want get the compliance report for. One of device_ids or device_service_tags or device_group_names is required",
			Description:         "Group names is the list of group names that you want get the compliance report for. One of device_ids or device_service_tags or device_group_names is required",
			ElementType:         types.StringType,
			Optional:            true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_ids")),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_service_tags")),
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
				),
			},
		},
	}
}

// The schema for OME Device Compliance Report
func omeSingleDeviceComplianceReportDataSchema() map[string]schema.Attribute {
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
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
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
					"id": schema.Int64Attribute{
						MarkdownDescription: "Id",
						Description:         "Id",
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
						MarkdownDescription: "Target Identifer",
						Description:         "Target Identifer",
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
					"dependency_upgrade_required": schema.BoolAttribute{
						MarkdownDescription: "Dependency Upgrade Required",
						Description:         "Dependency Upgrade Required",
						Computed:            true,
					},
					"version": schema.StringAttribute{
						MarkdownDescription: "Version",
						Description:         "Version",
						Computed:            true,
					},
					"compliance_dependencies": schema.ListNestedAttribute{
						MarkdownDescription: "Compliance Dependencies",
						Description:         "Compliance Dependencies",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"compliance_dependency_id": schema.Int64Attribute{
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
							},
						},
					},
				},
			},
		},
		"device_id": schema.Int64Attribute{
			MarkdownDescription: "Device Id",
			Description:         "Device Id",
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
			MarkdownDescription: "Device Type Id",
			Description:         "Device Type Id",
			Computed:            true,
		},
		"device_type_name": schema.StringAttribute{
			MarkdownDescription: "Device Type Name",
			Description:         "Device Type Name",
			Computed:            true,
		},
		"firmware_status": schema.StringAttribute{
			MarkdownDescription: "Firmware Status",
			Description:         "Firmware Status",
			Computed:            true,
		},
		"id": schema.Int64Attribute{
			MarkdownDescription: "Id",
			Description:         "Id",
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
		"device_firmware_update_capable": schema.BoolAttribute{
			MarkdownDescription: "Device Firmware Update Capable",
			Description:         "Device Firmware Update Capable",
			Computed:            true,
		},
		"device_user_firmware_update_capable": schema.BoolAttribute{
			MarkdownDescription: "Device User Firmware Update Capable",
			Description:         "Device User Firmware Update Capable",
			Computed:            true,
		},
	}
}
