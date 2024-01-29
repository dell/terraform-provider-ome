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
	"strconv"
	"terraform-provider-ome/helper"
	"terraform-provider-ome/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &resourceFirmwareBaseline{}
	_ resource.ResourceWithConfigure = &resourceFirmwareBaseline{}
)

// NewFirmwareBaselineResource is new resource for firmware baseline
func NewFirmwareBaselineResource() resource.Resource {
	return &resourceFirmwareBaseline{}
}

type resourceFirmwareBaseline struct {
	p *omeProvider
}

const (
	// BaselineRetryCount - stores the default value of retry count
	BaselineRetryCount = 5
	// BaselineSleepInterval - stores the default value of sleep interval
	BaselineSleepInterval = 30
	// BaselineSleepTimeBeforeJob - wait time in seconds before job tracking
	BaselineSleepTimeBeforeJob = 5
)

// Configure implements resource.ResourceWithConfigure
func (r *resourceFirmwareBaseline) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (r resourceFirmwareBaseline) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "firmware_baseline"
}

// Devices Resource schema
func (r resourceFirmwareBaseline) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This terraform resource is used to manage firmware baseline entity on OME." +
			"We can Create, Update and Delete OME firmware baseline using this resource. We can also do an 'Import' an existing 'firmware baseline' from OME .",
		Version:    1,
		Attributes: FirmwareBaselineSchema(),
	}
}

// Create a new Firmware Baseline resource
func (r resourceFirmwareBaseline) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get Plan Data
	tflog.Trace(ctx, "resource_firmware_baseline create: started")
	var plan models.FirmwareBaselineResource
	diags := req.Plan.Get(ctx, &plan)
	// Read Terraform plan into the model
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_firmware_baseline Configure")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}

	defer omeClient.RemoveSession()

	var payload models.CreateUpdateFirmwareBaseline

	targets, err := helper.CreateTargetModel(omeClient, plan)

	if err != nil {
		resp.Diagnostics.AddError(
			`Unable to create target model for: `+plan.Name.ValueString()+``, err.Error(),
		)
		return
	}
	if len(targets) == 0 {
		resp.Diagnostics.AddError(
			`Unable to create target model for: `+plan.Name.ValueString()+``, "No targets found. Please check device names / service tags / group names",
		)
		return
	}
	payload.Targets = targets

	if plan.CatalogName.ValueString() != "" {
		catalog, err := helper.GetCatalogFirmwareByName(omeClient, plan.CatalogName.ValueString())
		if err != nil || catalog == nil {
			resp.Diagnostics.AddError("Not Found", "Catalog details not found")
			return
		}
		payload.CatalogID = catalog.ID
		payload.RepositoryID = catalog.Repository.ID
	}
	if plan.Name.ValueString() != "" {
		payload.Name = plan.Name.ValueString()
	}
	payload.Is64Bit = plan.Is64Bit.ValueBool()
	payload.FilterNoRebootRequired = plan.FilterNoRebootRequired.ValueBool()
	payload.Description = plan.Description.ValueString()
	jobID, errCreate := helper.CreateFirmwareBaseline(omeClient, payload)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			`Unable to create Baseline: `+plan.Name.ValueString()+``, errCreate.Error(),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Baseline created with id %d", jobID))
	// Wait for the job to finish
	time.Sleep(BaselineSleepTimeBeforeJob * time.Second)

	if jobID != 0 {
		isSuccess, message := omeClient.TrackJob(jobID, BaselineRetryCount, BaselineSleepInterval)
		if !isSuccess {
			resp.Diagnostics.AddError(
				"Create Baseline job for: "+plan.Name.ValueString()+" has some errors",
				message,
			)
			return
		}
	}

	// Get Firmware Baseline Data
	omeBaselineData, errGet := helper.GetFirmwareBaselineWithName(*omeClient, plan.Name.ValueString())
	if errGet != nil {
		resp.Diagnostics.AddError(
			`Could not get Baseline after create: `+plan.Name.ValueString()+``, errGet.Error(),
		)
		return
	}
	// Set the tf state after Read
	state, errCopy := helper.SetStateBaseline(ctx, omeBaselineData, plan)
	state.DeviceNames = plan.DeviceNames
	state.DeviceServiceTags = plan.DeviceServiceTags
	state.GroupNames = plan.GroupNames
	if errCopy != nil {
		resp.Diagnostics.AddError(
			"Could not copy Baseline data", errCopy.Error(),
		)
		return
	}

	//Save into State if template creation is successful
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_firmware_baseline create: finished")

}

// Read resource information
func (r resourceFirmwareBaseline) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Trace(ctx, "resource_firmware_baseline Read: started")
	var curState models.FirmwareBaselineResource
	diags := req.State.Get(ctx, &curState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_firmware_baseline Configure")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}

	defer omeClient.RemoveSession()

	omeBaselineData, err := helper.GetFirmwareBaselineWithID(*omeClient, curState.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			`Could not Read Baseline: `+curState.Name.ValueString()+``, err.Error(),
		)
		return
	}
	// Set the tf state after Read
	state, errCopy := helper.SetStateBaseline(ctx, omeBaselineData, curState)
	state.DeviceNames = curState.DeviceNames
	state.DeviceServiceTags = curState.DeviceServiceTags
	state.GroupNames = curState.GroupNames
	if errCopy != nil {
		resp.Diagnostics.AddError(
			"Could not copy Baseline data", errCopy.Error(),
		)
		return
	}

	//Save into State if template creation is successful
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_firmware_baseline Read: finished")

}

// Update resource
func (r resourceFirmwareBaseline) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "resource_firmware_baseline: update started")
	var plan models.FirmwareBaselineResource
	var state models.FirmwareBaselineResource

	diagState := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diagState...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagsPlan := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diagsPlan...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_firmware_baseline Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	// Update Firmware Baseline based on the plan
	jobID, errUpd := helper.UpdateFirmwareBaseline(*omeClient, state, plan)
	if errUpd != nil {
		resp.Diagnostics.AddError(
			`Unable to Update Baseline: `+plan.Name.ValueString()+``, errUpd.Error(),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Baseline Updated with id %d", jobID))
	// Wait for the job to finish
	time.Sleep(BaselineSleepTimeBeforeJob * time.Second)

	if jobID != 0 {
		isSuccess, message := omeClient.TrackJob(jobID, BaselineRetryCount, BaselineSleepInterval)
		if !isSuccess {
			resp.Diagnostics.AddError(
				"Update Baseline job for: "+plan.Name.ValueString()+" has some errors",
				message,
			)
			return
		}
	}

	omeBaselineData, err := helper.GetFirmwareBaselineWithName(*omeClient, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			`Could not get Baseline after update: `+plan.Name.ValueString()+``, err.Error(),
		)
		return
	}

	// Set the tf state after update
	updState, errCopy := helper.SetStateBaseline(ctx, omeBaselineData, plan)
	updState.DeviceNames = plan.DeviceNames
	updState.DeviceServiceTags = plan.DeviceServiceTags
	updState.GroupNames = plan.GroupNames
	if errCopy != nil {
		resp.Diagnostics.AddError(
			"Could not copy Baseline data after update", errCopy.Error(),
		)
		return
	}

	diags := resp.State.Set(ctx, &updState)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "resource_firmware_baseline: update end")
}

// Delete resource
func (r resourceFirmwareBaseline) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_firmware_baseline delete: started")
	var state models.FirmwareBaselineResource
	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_firmware_baseline Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_firmware_baseline delete: started delete for baseline", map[string]interface{}{
		"baselineId": state.ID.ValueInt64(),
	})

	err := helper.DeleteFirmwareBaseline(*omeClient, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete Baseline",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "resource_firmware_baselinee delete: finished")
}

// ImportState imports an existing Resource
func (r *resourceFirmwareBaseline) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var importState models.FirmwareBaselineResource
	tflog.Info(ctx, "resource_firmware_baseline: import state started")
	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_firmware_baseline Import")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	id, coversionErr := strconv.Atoi(req.ID)
	if coversionErr != nil {
		resp.Diagnostics.AddError(
			`Unable to import firmware baseline, id must be an integer: `+req.ID+``, coversionErr.Error(),
		)
	}
	tflog.Trace(ctx, fmt.Sprintf(" Firmware Baseline: import state id is %d", id))

	baseline, err := helper.GetFirmwareBaselineWithID(*omeClient, int64(id))
	if err != nil {
		resp.Diagnostics.AddError(
			`Unable to import firmware baseline: `+req.ID+``, err.Error(),
		)
	}

	// Set the tf state after read
	state, mapErr := helper.SetStateBaseline(ctx, baseline, importState)
	state.DeviceNames = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("")})
	state.DeviceServiceTags = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("")})
	state.GroupNames = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("")})
	if mapErr != nil {
		resp.Diagnostics.AddError(
			`Unable to process map state for  firmware baseline: `+state.Name.ValueString()+`.`, mapErr.Error(),
		)
		return
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Info(ctx, "resource_firmware_baseline: import state end")
}
