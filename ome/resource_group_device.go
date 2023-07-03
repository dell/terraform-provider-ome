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
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceDeviceGroup{}
	_ resource.ResourceWithConfigure   = &resourceDeviceGroup{}
	_ resource.ResourceWithImportState = &resourceDeviceGroup{}
)

// NewDeviceGroupResource is a new resource for deployment
func NewDeviceGroupResource() resource.Resource {
	return &resourceDeviceGroup{}
}

type resourceDeviceGroup struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceDeviceGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (resourceDeviceGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "group_device"
}

// Template DeviceGroup Resource schema
func (r resourceDeviceGroup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for static Device Groups on OpenManage Enterprise.",
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the device group resource.",
				Description:         "ID of the device group resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the template resource.",
				Description:         "Name of the template resource.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the template",
				Description:         "Description of the template",
				Required:            true,
			},
			"parent_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the parent group of the template." +
					" If the group is to be a root group, this field should be put to `0`." +
					" Defaults to `0`.",
				Description: "ID of the parent group of the template." +
					" If the group is to be a root group, this field should be put to `0`." +
					" Defaults to `0`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"membership_type_id": schema.Int64Attribute{
				MarkdownDescription: "Membership type of the template",
				Description:         "Membership type of the template",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create a new resource
func (r resourceDeviceGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_group_device create : Started")
	//Get Plan Data
	var plan models.GroupDeviceRes
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_group_device Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	var (
		id    int64
		err   error
		group models.Group
	)

	if id, err = omeClient.CreateGroupDevice(plan.GetPayload()); err != nil {
		resp.Diagnostics.AddError(
			"Error while creation", err.Error(),
		)
		return
	}

	group, err = omeClient.GetGroupById(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching group after creation", err.Error(),
		)
		return
	}

	state, _ := models.NewGroupDeviceRes(group)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read resource information
func (r resourceDeviceGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get State Data
	tflog.Trace(ctx, "resource_group_device read: started")
	var state models.GroupDeviceRes
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_group_device Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Trace(ctx, "resource_group_device read: client created started updating state")

	group, err := omeClient.GetGroupById(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching group after creation", err.Error(),
		)
	}

	state, _ = models.NewGroupDeviceRes(group)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "resource_group_device read: finished")
}

// Update resource
func (r resourceDeviceGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get state Data
	tflog.Trace(ctx, "resource_group_device update: started")
	var state models.GroupDeviceRes
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan Data
	var plan models.GroupDeviceRes
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_group_device Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	if err := omeClient.UpdateGroupDevice(plan.GetPayload()); err != nil {
		resp.Diagnostics.AddError(
			"Error while updation", err.Error(),
		)
		return
	}

	group, err := omeClient.GetGroupById(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching group after updation", err.Error(),
		)
		return
	}

	state, _ = models.NewGroupDeviceRes(group)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete resource
func (r resourceDeviceGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_group_device delete: started")
	// Get State Data
	var state models.GroupDeviceRes
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_group_device Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	err := omeClient.DeleteGroupDevice(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting device group",
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_group_device delete: finished")
}

// Import resource
func (r resourceDeviceGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "resource_group_device import: started")
	// Save the import identifier in the id attribute
	// var state models.GroupDeviceRes
	groupName := req.ID

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_group_device ImportState")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	group, err := omeClient.GetSingleGroupByName(groupName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing group", err.Error(),
		)
		return
	}

	state, _ := models.NewGroupDeviceRes(group)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "resource_group_device import: finished")
}
