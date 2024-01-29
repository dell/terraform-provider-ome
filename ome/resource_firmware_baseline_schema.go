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

package ome

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FirmwareBaselineSchema returns the schema of the firmware baseline
func FirmwareBaselineSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"catalog_id": schema.Int64Attribute{
			MarkdownDescription: "ID of the catalog.",
			Description:         "ID of the catalog.",
			Computed:            true,
		},
		"compliance_summary": schema.ObjectAttribute{
			MarkdownDescription: "Compliance Summary",
			Description:         "Compliance Summary",
			Computed:            true,
			AttributeTypes:      SingleComplianceSummarySchema(),
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Description of the firmware baseline",
			Description:         "Description of the firmware baseline",
			Computed:            true,
			Optional:            true,
		},
		"downgrade_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates if the firmware can be downgraded",
			Description:         "Indicates if the firmware can be downgraded",
			Computed:            true,
		},
		"filter_no_reboot_required": schema.BoolAttribute{
			MarkdownDescription: "Filters applicable updates where no reboot is required during create baseline for firmware updates.",
			Description:         "Filters applicable updates where no reboot is required during create baseline for firmware updates.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of the firmware baseline.",
			Description:         "ID of the firmware baseline.",
			Computed:            true,
			Optional:            true,
		},
		"is_64_bit": schema.BoolAttribute{
			MarkdownDescription: "This must always be set to true. The size of the DUP files used is 64 bits.",
			Description:         "This must always be set to true. The size of the DUP files used is 64 bits.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(true)),
			},
		},
		"last_run": schema.StringAttribute{
			MarkdownDescription: "Last Run Time for the firmware baseline",
			Description:         "Last Run Time for the firmware baseline",
			Computed:            true,
			Optional:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Name of the firmware baseline",
			Description:         "Name of the firmware baseline",
			Required:            true,
			Validators:          []validator.String{stringvalidator.LengthBetween(1, 256)},
		},
		"repository_id": schema.Int64Attribute{
			MarkdownDescription: "ID of the repository. Derived from the catalog response",
			Description:         "ID of the repository. Derived from the catalog response",
			Computed:            true,
		},
		"repository_name": schema.StringAttribute{
			MarkdownDescription: "Name of the repository",
			Description:         "Name of the repository",
			Computed:            true,
			Optional:            true,
		},
		"repository_type": schema.StringAttribute{
			MarkdownDescription: "Type of the repository",
			Description:         "Type of the repository",
			Computed:            true,
		},
		"targets": schema.ListNestedAttribute{
			MarkdownDescription: "The DeviceID, if the baseline is being created for devices or, the GroupID, if the baseline is being created for a group of devices.",
			Description:         "The DeviceID, if the baseline is being created for devices or, the GroupID, if the baseline is being created for a group of devices.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SingleTargetSchema()},
		},
		"task_id": schema.Int64Attribute{
			MarkdownDescription: "Identifier of task which created this baseline.",
			Description:         "Identifier of task which created this baseline.",
			Computed:            true,
		},
		"task_status": schema.StringAttribute{
			MarkdownDescription: "Task status.",
			Description:         "Task status.",
			Computed:            true,
		},
		"catalog_name": schema.StringAttribute{
			MarkdownDescription: "Name of the catalog",
			Description:         "Name of the catalog",
			Required:            true,
		},
		"device_names": schema.ListAttribute{
			MarkdownDescription: "Device names is the list of device names that you want to add to the firmware baseline being created. One of DeviceNames or DeviceServiceTags or GroupNames is required",
			Description:         "Device names is the list of device names that you want to add to the firmware baseline being created. One of DeviceNames or DeviceServiceTags or GroupNames is required",
			ElementType:         types.StringType,
			Optional:            true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_service_tags")),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_names")),
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
				),
			},
		},
		"device_service_tags": schema.ListAttribute{
			MarkdownDescription: "Device service tags is the list of device service tags that you want to add to the firmware baseline being created.One of DeviceNames or DeviceServiceTags or GroupNames is required",
			Description:         "Device service tags is the list of device service tags that you want to add to the firmware baseline being created.One of DeviceNames or DeviceServiceTags or GroupNames is required",
			ElementType:         types.StringType,
			Optional:            true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_names")),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_names")),
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
				),
			},
		},
		"group_names": schema.ListAttribute{
			MarkdownDescription: "Group names is the list of group names that you want to add to the firmware baseline being created.One of DeviceNames or DeviceServiceTags or GroupNames is required",
			Description:         "Group names is the list of group names that you want to add to the firmware baseline being created.One of DeviceNames or DeviceServiceTags or GroupNames is required",
			ElementType:         types.StringType,
			Optional:            true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_names")),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_service_tags")),
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
				),
			},
		},
	}
}

// SingleComplianceSummarySchema returns a map of attribute types for the compliance summary.
func SingleComplianceSummarySchema() map[string]attr.Type {
	return map[string]attr.Type{
		"compliance_status":   types.StringType,
		"number_of_critical":  types.Int64Type,
		"number_of_downgrade": types.Int64Type,
		"number_of_normal":    types.Int64Type,
		"number_of_warning":   types.Int64Type,
		"number_of_unknown":   types.Int64Type,
	}
}

// SingleTargetSchema returns a map of attribute types for the target.
func SingleTargetSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of device associated with the firmware baseline.",
			Description:         "ID of device associated with the firmware baseline.",
			Required:            true,
		},
		"type": schema.ObjectAttribute{
			MarkdownDescription: "Type of device associated with the firmware baseline..",
			Description:         "Type of device associated with the firmware baseline..",
			Required:            true,
			AttributeTypes:      SingleTargetTypeSchema(),
		},
	}
}

// SingleTargetTypeSchema returns a map of attribute types for the target type.
func SingleTargetTypeSchema() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.Int64Type,
		"name": types.StringType,
	}
}
