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

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &resourceDeviceAction{}
	_ resource.ResourceWithConfigure = &resourceDeviceAction{}
)

const (
	defaultJobTimeout int64 = 10
	interval                = 5
	name                    = "Just-trying-out"
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
		Description: "This terraform resource is used to run actions on devices managed by OME." +
			" The only supported action, for now, is refreshing inventory." +
			" This resource creates a job in OME to run the actions and does not support updating in-place." +
			" The resource generates a recreation plan instead for any necessary update action.",
		MarkdownDescription: "This terraform resource is used to run actions on devices managed by OME." +
			" The only supported action, for now, is refreshing inventory." +
			" This resource creates a job in OME to run the actions and does not support updating in-place." +
			" The resource generates a recreation plan instead for any necessary update action.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the job created on OME appliance for carrying out the action.",
				Description:         "ID of the job created on OME appliance for carrying out the action.",
				Computed:            true,
			},
			"device_ids": schema.ListAttribute{
				MarkdownDescription: "List of id of devices on whom the action would be carried out.",
				Description:         "List of id of devices on whom the action would be carried out.",
				Required:            true,
				ElementType:         types.Int64Type,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
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
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cron": schema.StringAttribute{
				MarkdownDescription: "Cron expression to schedule an action in the future." +
					" If not specified, the action runs immediately on apply." +
					" Conflicts with `timeout`.",
				Description: "Cron expression to schedule an action in the future." +
					" If not specified, the action runs immediately on apply." +
					" Conflicts with 'timeout'.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("timeout")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout, in minutes, for monitoring an immediately running action." +
					" Conflicts with `cron`." +
					" Default value is `10`.",
				Description: "Timeout, in minutes, for monitoring an immediately running action." +
					" Conflicts with 'cron'." +
					" Default value is '10'.",
				Optional: true,
			},
			"job_name": schema.StringAttribute{
				MarkdownDescription: "Name of the job to be created on the OME appliance that will run the action.",
				Description:         "Name of the job to be created on the OME appliance that will run the action.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"job_description": schema.StringAttribute{
				MarkdownDescription: "Description of the job to be created on the OME appliance that will run the action.",
				Description:         "Description of the job to be created on the OME appliance that will run the action.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_run_time": schema.StringAttribute{
				MarkdownDescription: "Last run time of the job.",
				Description:         "Last run time of the job.",
				Computed:            true,
			},
			"next_run_time": schema.StringAttribute{
				MarkdownDescription: "Next run time of the job.",
				Description:         "Next run time of the job.",
				Computed:            true,
			},
			"last_run_status": schema.StringAttribute{
				MarkdownDescription: "Last run status of the job.",
				Description:         "Last run status of the job.",
				Computed:            true,
			},
			"current_status": schema.StringAttribute{
				MarkdownDescription: "Current status of the job.",
				Description:         "Current status of the job.",
				Computed:            true,
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Start time of the job.",
				Description:         "Start time of the job.",
				Computed:            true,
			},
			"end_time": schema.StringAttribute{
				MarkdownDescription: "End time of the job.",
				Description:         "End time of the job.",
				Computed:            true,
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
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Cron.IsNull() {
		return
	}

	timeout := defaultJobTimeout
	if !plan.Timeout.IsNull() {
		timeout = plan.Timeout.ValueInt64()
	}
	retries := timeout * 60 / interval

	if ok, message := omeClient.TrackJob(state.ID.ValueInt64(), retries, interval); !ok {
		resp.Diagnostics.AddError(
			"Refresh Job could not complete.",
			message,
		)
	} else {
		tflog.Info(ctx, "Refresh job completed successfully. "+message)
	}

	state, dgs = r.read(ctx, state)
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
	jobResp, err := r.c.RefreshDeviceInventory(plan.DeviceIDs, clients.JobOpts{
		Name:        plan.JobName.ValueString(),
		Description: plan.JobDescription.ValueString(),
		Schedule:    plan.Cron.ValueString(),
	})
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

func (r resourceDeviceAction) convertJobRespToTfsdk(ctx context.Context, resp clients.JobResp,
	pstate models.DeviceActionModel) models.DeviceActionModel {
	ret := r.getJobModel(resp, pstate)
	ret.ID = types.Int64Value(resp.ID)
	ret.DeviceIDs = pstate.DeviceIDs
	ret.Action = pstate.Action
	return ret
}

func (r resourceDeviceAction) getJobModel(resp clients.JobResp, pstate models.DeviceActionModel) models.DeviceActionModel {
	cron := types.StringNull()
	if !(resp.Schedule == clients.RunNowSchedule || resp.Schedule == "") {
		cron = types.StringValue(resp.Schedule)
	}
	return models.DeviceActionModel{
		Cron:           cron,
		Timeout:        pstate.Timeout,
		JobName:        types.StringValue(resp.JobName),
		JobDescription: types.StringValue(resp.JobDescription),
		NextRunTime:    types.StringValue(resp.NextRun),
		LastRunTime:    types.StringValue(resp.LastRun),
		JobStatus:      types.StringValue(resp.JobStatus.Name),
		LastRunStatus:  types.StringValue(resp.LastRunStatus.Name),
		StartTime:      types.StringValue(resp.StartTime),
		EndTime:        types.StringValue(resp.EndTime),
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
	// update ONLY happens if someone ONLY changes timeout
	// so set state timeout as plan
	var plan, state models.DeviceActionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Timeout = plan.Timeout
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
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
