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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceStaticGroup{}
	_ resource.ResourceWithConfigure   = &resourceStaticGroup{}
	_ resource.ResourceWithImportState = &resourceStaticGroup{}
)

// NewDeviceGroupResource is a new resource for deployment
func NewStaticGroupResource() resource.Resource {
	return &resourceStaticGroup{}
}

type resourceStaticGroup struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceStaticGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (resourceStaticGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "static_group"
}

// Template DeviceGroup Resource schema
func (r resourceStaticGroup) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for static Device Groups on OpenManage Enterprise.",
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the static group resource.",
				Description:         "ID of the static group resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the static group resource.",
				Description:         "Name of the static group resource.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the static group",
				Description:         "Description of the static group",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the parent group of the static group.",
				Description:         "ID of the parent group of the static group.",
				Required:            true,
			},
			"membership_type_id": schema.Int64Attribute{
				MarkdownDescription: "Membership type of the static group",
				Description:         "Membership type of the static group",
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
func (r resourceStaticGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_static_group create : Started")
	//Get Plan Data
	var plan models.StaticGroup
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_static_group Create")
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
	if id, err = omeClient.CreateGroup(createPayload); err != nil {
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
func (r resourceStaticGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get State Data
	tflog.Trace(ctx, "resource_static_group read: started")
	var state models.StaticGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_static_group Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Trace(ctx, "resource_static_group read: client created started updating state")

	finalState, dgs := r.ReadRes(omeClient, state.ID.ValueInt64())
	resp.Diagnostics.Append(dgs...)
	if dgs.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalState)...)

	tflog.Trace(ctx, "resource_static_group read: finished")
}

// Update resource
func (r resourceStaticGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get state Data
	tflog.Trace(ctx, "resource_static_group update: started")
	var state models.StaticGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan Data
	var plan models.StaticGroup
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_static_group Update")
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

func (r resourceStaticGroup) UpdateRes(omeClient *clients.Client, plan, state models.StaticGroup, ctx context.Context) (models.StaticGroup, diag.Diagnostics) {

	var d diag.Diagnostics
	if payload, ok := plan.GetPayload(state); !ok {
		if err := omeClient.UpdateGroup(payload); err != nil {
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
		if err := omeClient.AddGroupMembers(payloadAdd); err != nil {
			d.AddError(
				"Error while adding group devices", err.Error(),
			)
			return state, d
		}
	}
	if len(payloadRmv.DeviceIds) != 0 {
		if err := omeClient.RemoveGroupMembers(payloadRmv); err != nil {
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

func (r resourceStaticGroup) ReadRes(omeClient *clients.Client, id int64) (ret models.StaticGroup, d diag.Diagnostics) {
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

	ret, _ = models.NewStaticGroup(group, devs)
	return ret, d
}

// Delete resource
func (r resourceStaticGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_static_group delete: started")
	// Get State Data
	var state models.StaticGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_static_group Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	err := omeClient.DeleteGroup(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting static group",
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_static_group delete: finished")
}

// Import resource
func (r resourceStaticGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "resource_static_group import: started")
	// Save the import identifier in the id attribute
	// var state models.StaticGroup
	groupName := req.ID

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_static_group ImportState")
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
	state, _ := models.NewStaticGroup(group, devs)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "resource_static_group import: finished")
}
