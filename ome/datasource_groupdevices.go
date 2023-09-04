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
	"context"
	"fmt"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &groupDevicesDatasource{}
	_ datasource.DataSourceWithConfigure = &groupDevicesDatasource{}
)

// NewGroupDevicesDatasource is new datasource for group devices
func NewGroupDevicesDatasource() datasource.DataSource {
	return &groupDevicesDatasource{}
}

type groupDevicesDatasource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *groupDevicesDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*groupDevicesDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "groupdevices_info"
}

// Schema implements datasource.DataSource
func (*groupDevicesDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to list groups and their devices from OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID for group devices data source.",
				Description:         "ID for group devices data source.",
				Computed:            true,
				Optional:            true,
			},
			"device_ids": schema.ListAttribute{
				MarkdownDescription: "List of the device id(s) associated with any of the groups.",
				Description:         "List of the device id(s) associated with any of the groups.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"device_servicetags": schema.ListAttribute{
				MarkdownDescription: "List of the device servicetags associated with any of the groups.",
				Description:         "List of the device servicetags associated with any of the groups.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"device_group_names": schema.SetAttribute{
				MarkdownDescription: "List of the device group names.",
				Description:         "List of the device group names.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"device_groups": schema.MapNestedAttribute{
				MarkdownDescription: "Map of the groups fetched keyed by its name.",
				Description:         "Map of the groups fetched keyed by its name.",
				NestedObject:        schema.NestedAttributeObject{Attributes: groupSchema()},
				Computed:            true,
			},
		},
	}
}

func groupSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of the group.",
			Description:         "ID of the group.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Name of the group.",
			Description:         "Name of the group.",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Description of the group.",
			Description:         "Description of the group.",
			Computed:            true,
		},
		"membership_type_id": schema.Int64Attribute{
			MarkdownDescription: "Membership Type ID of the group.",
			Description:         "Membership Type ID of the group.",
			Computed:            true,
		},
		"parent_id": schema.Int64Attribute{
			MarkdownDescription: "Parent ID of the group.",
			Description:         "Parent ID of the group.",
			Computed:            true,
		},
		"global_status": schema.Int64Attribute{
			MarkdownDescription: "global_status of the group.",
			Description:         "global_status of the group.",
			Computed:            true,
		},
		"id_owner": schema.Int64Attribute{
			MarkdownDescription: "ID Owner of the group.",
			Description:         "ID Owner of the group.",
			Computed:            true,
		},
		"creation_time": schema.StringAttribute{
			MarkdownDescription: "Creation time of the group.",
			Description:         "Creation time of the group.",
			Computed:            true,
		},
		"updated_time": schema.StringAttribute{
			MarkdownDescription: "Last updation time of the group.",
			Description:         "Last updation time of the group.",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the group.",
			Description:         "The user who created the group.",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who updated the group.",
			Description:         "The user who updated the group.",
			Computed:            true,
		},
		"visible": schema.BoolAttribute{
			MarkdownDescription: "If the group is visible or not.",
			Description:         "If the group is visible or not.",
			Computed:            true,
		},
		"definition_id": schema.Int64Attribute{
			MarkdownDescription: "Definition ID of the group.",
			Description:         "Definition ID of the group.",
			Computed:            true,
		},
		"definition_description": schema.StringAttribute{
			MarkdownDescription: "Definition description of the group.",
			Description:         "Definition description of the group.",
			Computed:            true,
		},
		"type_id": schema.Int64Attribute{
			MarkdownDescription: "Type ID of the group.",
			Description:         "Type ID of the group.",
			Computed:            true,
		},
		"has_attributes": schema.BoolAttribute{
			MarkdownDescription: "If the group has attributes.",
			Description:         "If the group has attributes.",
			Computed:            true,
		},
		"is_access_allowed": schema.BoolAttribute{
			MarkdownDescription: "If access of this group is allowed.",
			Description:         "If access of this group is allowed.",
			Computed:            true,
		},
		"devices": schema.SetNestedAttribute{
			MarkdownDescription: "Devices of the group.",
			Description:         "Devices of the group.",
			NestedObject:        schema.NestedAttributeObject{Attributes: deviceInputSchema()},
			Computed:            true,
		},
		"sub_groups": schema.SetNestedAttribute{
			MarkdownDescription: "Sub Groups of the group.",
			Description:         "Sub Groups of the group.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: subGroupSchema()},
		},
	}
}

func subGroupSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of the sub group.",
			Description:         "ID of the sub group.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Name of the sub group.",
			Description:         "Name of the sub group.",
			Computed:            true,
		},
	}
}

func deviceInputSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of the device.",
			Description:         "ID of the device.",
			Computed:            true,
		},
		"servicetag": schema.StringAttribute{
			MarkdownDescription: "Service Tag of the device",
			Description:         "Service Tag of the device",
			Computed:            true,
		},
	}
}

// Read implements datasource.DataSource
func (g *groupDevicesDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var groupDevices models.GroupDevicesData
	diags := req.Config.Get(ctx, &groupDevices)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	groupDevices.ID = types.StringValue("dummy")
	groupNames := []string{}
	resp.Diagnostics.Append(groupDevices.DeviceGroupNames.ElementsAs(ctx, &groupNames, true)...)
	if diags.HasError() {
		return
	}

	omeClient, d := g.p.createOMESession(ctx, "datasource_vlannetworks_info Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	allDevices := make([]models.Device, 0)
	for _, groupName := range groupNames {
		group, err := omeClient.GetExpandedGroupByName(groupName, "")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error getting group by name: %s", groupName),
				err.Error(),
			)
			continue
		}

		devices, err := omeClient.GetDevicesByGroupID(group.ID)
		if err != nil {
			if len(devices.Value) != 0 {
				resp.Diagnostics.AddWarning(
					"Unable to fetch devices during pagination",
					err.Error(),
				)
			} else {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Unable to fetch devices for group %s", group.Name),
					err.Error(),
				)
				return
			}
			continue
		}
		allDevices = append(allDevices, devices.Value...)
		groupDevices.SetGroup(group, devices.Value)
	}

	uniqueDevices := omeClient.GetUniqueDevices(allDevices)
	groupDevices.SetDevices(uniqueDevices)

	diags = resp.State.Set(ctx, &groupDevices)
	resp.Diagnostics.Append(diags...)
}
