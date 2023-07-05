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
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				MarkdownDescription: "ID of the parent group of the template.",
				Description:         "ID of the parent group of the template.",
				Required:            true,
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
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "Device ids in the group",
				Description:         "Device ids in the group",
				Required:            true,
				ElementType:         types.Int64Type,
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
		id  int64
		err error
	)

	createPayload, _ := plan.GetPayload(plan)
	if id, err = omeClient.CreateGroupDevice(createPayload); err != nil {
		resp.Diagnostics.AddError(
			"Error while creation", err.Error(),
		)
		return
	}

	initialState, dgs := r.ReadRes(omeClient, id)
	resp.Diagnostics.Append(dgs...)
	if dgs.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, initialState)...)
	if dgs.HasError() {
		return
	}

	finalState, dgs2 := r.UpdateRes(omeClient, plan, initialState, ctx)
	resp.Diagnostics.Append(dgs2...)
	if dgs2.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalState)...)
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

	finalState, dgs := r.ReadRes(omeClient, state.ID.ValueInt64())
	resp.Diagnostics.Append(dgs...)
	if dgs.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalState)...)

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

	finalState, dgs := r.UpdateRes(omeClient, plan, state, ctx)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalState)...)
}

func (r resourceDeviceGroup) UpdateRes(omeClient *clients.Client, plan, state models.GroupDeviceRes, ctx context.Context) (models.GroupDeviceRes, diag.Diagnostics) {

	var d diag.Diagnostics
	if payload, ok := plan.GetPayload(state); !ok {
		if err := omeClient.UpdateGroupDevice(payload); err != nil {
			d.AddError(
				"Error while updation", err.Error(),
			)
			return state, d
		}
	}

	payloadAdd, payloadRmv, dgs := plan.GetMemberPayload(ctx, state)
	d.Append(dgs...)
	if d.HasError() {
		return state, d
	}
	if len(payloadAdd.DeviceIds) != 0 {
		if err := omeClient.AddGroupDeviceMembers(payloadAdd); err != nil {
			d.AddError(
				"Error while adding group devices", err.Error(),
			)
			return state, d
		}
	}
	if len(payloadRmv.DeviceIds) != 0 {
		if err := omeClient.RemoveGroupDeviceMembers(payloadRmv); err != nil {
			d.AddError(
				"Error while removing group devices", err.Error(),
			)
			return state, d
		}
	}

	ret, dgs2 := r.ReadRes(omeClient, state.ID.ValueInt64())
	d.Append(dgs2...)
	return ret, d
}

func (r resourceDeviceGroup) ReadRes(omeClient *clients.Client, id int64) (ret models.GroupDeviceRes, d diag.Diagnostics) {
	group, err := omeClient.GetGroupById(id)
	if err != nil {
		d.AddError(
			"Error fetching group by id", err.Error(),
		)
		return ret, d
	}

	devs, err2 := omeClient.GetDevicesByGroupID(id)
	if err2 != nil {
		d.AddError(
			"Error reading devices of group", err.Error(),
		)
		return
	}

	ret, _ = models.NewGroupDeviceRes(group, devs)
	return ret, d
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

	devs, err2 := omeClient.GetDevicesByGroupID(group.ID)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error importing devices of group", err.Error(),
		)
		return
	}
	state, _ := models.NewGroupDeviceRes(group, devs)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "resource_group_device import: finished")
}
