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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &resourceDeviceAction{}
	_ resource.ResourceWithConfigure = &resourceDeviceAction{}
)

// NewDeviceActionResource is new resource for device_action
func NewDeviceActionResource() resource.Resource {
	return &resourceDeviceAction{}
}

type resourceDeviceAction struct {
	p *omeProvider
	c *clients.Client
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceDeviceAction) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (r resourceDeviceAction) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "device_action"
}

// Devices Resource schema
func (r resourceDeviceAction) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing device_action on OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the job created for carrying out the action.",
				Description:         "ID of the job created for carrying out the action.",
				Computed:            true,
			},
			"device_ids": schema.ListAttribute{
				MarkdownDescription: "List of device_action to be managed by this resource.",
				Description:         "List of device_action to be managed by this resource.",
				Required:            true,
				ElementType:         types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action to be performed on the devices." +
					" Accepted values are [`inventory_refresh`]." +
					" Default value is `inventory_refresh`.",
				Description: "Action to be performed on the devices." +
					" Accepted values are ['inventory_refresh']." +
					" Default value is 'inventory_refresh'.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("inventory_refresh"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

// Create a new resource
func (r resourceDeviceAction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_device_action create: started")
	var plan models.DeviceActionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_device_action Configure")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	r.c = omeClient
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_device_action getting current infrastructure state")

	state, dgs := r.create(ctx, plan)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDeviceAction) create(ctx context.Context, plan models.DeviceActionModel) (
	models.DeviceActionModel, diag.Diagnostics) {
	var dgs diag.Diagnostics
	jobResp, err := r.c.RefreshDeviceInventory(plan.DeviceIDs)
	if err != nil {
		dgs.AddError("Error creating job.", err.Error())
		return plan, dgs
	}
	return r.convertJobRespToTfsdk(ctx, jobResp, plan), dgs
}

func (r resourceDeviceAction) read(ctx context.Context, pstate models.DeviceActionModel) (
	models.DeviceActionModel, diag.Diagnostics) {
	var (
		state models.DeviceActionModel
		dgs   diag.Diagnostics
	)

	id := pstate.ID.ValueInt64()
	jobResp, err := r.c.GetJob(id)
	if err != nil {
		dgs.AddError("Job not found.", err.Error())
		return state, dgs
	}
	return r.convertJobRespToTfsdk(ctx, jobResp, pstate), dgs
}

func (r resourceDeviceAction) convertJobRespToTfsdk(ctx context.Context, jobResp clients.JobResp,
	pstate models.DeviceActionModel) models.DeviceActionModel {
	return models.DeviceActionModel{
		ID:        types.Int64Value(jobResp.ID),
		DeviceIDs: pstate.DeviceIDs,
		Action:    pstate.Action,
	}
}

// Read resource information
func (r resourceDeviceAction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get State Data
	tflog.Trace(ctx, "resource_device_action read: started")
	var state models.DeviceActionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_device_action Read")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	r.c = omeClient
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_device_action getting current job state")

	state, dgs := r.read(ctx, state)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceDeviceAction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update is not expected by this resource.",
		"An update plan should never be generated. Please report this bug to the developers.",
	)
}

// Delete resource
func (r resourceDeviceAction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//Get State Data
	tflog.Trace(ctx, "resource_device_action delete: started")
	var state models.DeviceActionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_device_action Delete")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_device_action deleting job")
	err := omeClient.DeleteJob(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting job", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
