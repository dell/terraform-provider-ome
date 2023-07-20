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
		MarkdownDescription: "Data source to list the devices in the group from OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "ID for group devices data source.",
				Description:         "ID for group devices data source.",
				Computed:            true,
				Optional:            true,
			},

			"device_ids": schema.ListAttribute{
				MarkdownDescription: "List of the device id(s) associated with a group",
				Description:         "List of the device id(s) associated with a group",
				ElementType:         types.Int64Type,
				Computed:            true,
			},

			"device_servicetags": schema.ListAttribute{
				MarkdownDescription: "List of the device servicetags associated with a group",
				Description:         "List of the device servicetags associated with a group",
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
				MarkdownDescription: "List of the device group names.",
				Description:         "List of the device group names.",
				NestedObject:        schema.NestedAttributeObject{Attributes: OmeGroupSchema()},
				Computed:            true,
			},
		},
	}
}

func OmeGroupSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},

		"description": schema.StringAttribute{
			MarkdownDescription: "description",
			Description:         "description",
			Computed:            true,
		},

		"membership_type_id": schema.Int64Attribute{
			MarkdownDescription: "Membership Type ID",
			Description:         "Membership Type ID",

			Computed: true,
			// ElementType: types.Int64Type,
		},

		"parent_id": schema.Int64Attribute{
			MarkdownDescription: "Parent ID",
			Description:         "Parent ID",

			Computed: true,
			// ElementType: types.Int64Type,
		},

		"global_status": schema.Int64Attribute{
			MarkdownDescription: "global_status",
			Description:         "global_status",

			Computed: true,
		},

		"id_owner": schema.Int64Attribute{
			MarkdownDescription: "IDOwner",
			Description:         "IDOwner",

			Computed: true,
		},

		"creation_time": schema.StringAttribute{
			MarkdownDescription: "creation_time",
			Description:         "creation_time",

			Computed: true,
		},

		"updated_time": schema.StringAttribute{
			MarkdownDescription: "updated_time",
			Description:         "updated_time",

			Computed: true,
		},

		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",

			Computed: true,
		},

		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",

			Computed: true,
		},

		"visible": schema.BoolAttribute{
			MarkdownDescription: "visible",
			Description:         "visible",

			Computed: true,
		},

		"definition_id": schema.Int64Attribute{
			MarkdownDescription: "Definition ID",
			Description:         "Definition ID",

			Computed: true,
		},

		"definition_description": schema.StringAttribute{
			MarkdownDescription: "definition_description",
			Description:         "definition_description",

			Computed: true,
		},

		"type_id": schema.Int64Attribute{
			MarkdownDescription: "Type ID",
			Description:         "Type ID",

			Computed: true,
		},

		"has_attributes": schema.BoolAttribute{
			MarkdownDescription: "has_attributes",
			Description:         "has_attributes",

			Computed: true,
		},

		"is_access_allowed": schema.BoolAttribute{
			MarkdownDescription: "is_access_allowed",
			Description:         "is_access_allowed",

			Computed: true,
		},
		"devices": schema.SetNestedAttribute{
			MarkdownDescription: "device ids",
			Description:         "device ids",
			NestedObject:        schema.NestedAttributeObject{Attributes: OmeDeviceInputSchema()},
			Computed:            true,
		},

		"sub_groups": schema.SetNestedAttribute{
			MarkdownDescription: "Sub Groups",
			Description:         "Sub Groups",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: OmeSubGroupSchema()},
		},
	}
}

func OmeSubGroupSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},
	}
}

func OmeDeviceInputSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},
		"servicetag": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
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
