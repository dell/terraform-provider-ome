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
	"context"
	"strconv"
	"terraform-provider-ome/helper"
	"terraform-provider-ome/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &networkSettingResource{}
)

// NewFirmwareCatalogResource is a helper function to simplify the provider implementation.
func NewFirmwareCatalogResource() resource.Resource {
	return &firmwareCatalogResource{}
}

// firmwareCatalogResource is the resource implementation.
type firmwareCatalogResource struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *firmwareCatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata returns the resource type name.
func (r *firmwareCatalogResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "firmware_catalog"
}

// Schema implements resource.Resource.
func (r *firmwareCatalogResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This terraform resource is used to manage firmware catalogs entity on OME." +
			"We can Create, Update and Delete OME firmware catalogs using this resource. We can also do an 'Import' an existing 'firmware catalog' from OME .",
		Version:    1,
		Attributes: FirmwareCatalogSchema(),
	}
}

// Create implements resource.Resource.
func (r *firmwareCatalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "firmwareCatalogResource: create started")
	var plan models.OmeSingleCatalogResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_firmware_catalog Configure")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}

	defer omeClient.RemoveSession()

	valError := helper.ValidateCatalogCreate(plan)

	if valError != nil {
		resp.Diagnostics.AddError(
			`Unable to create catalog, validation error: `, valError.Error(),
		)
		return
	}

	createModel := helper.MakeCatalogJSONModel(0, 0, plan)
	cat, err := helper.CreateCatalogFirmware(omeClient, createModel)

	if err != nil {
		resp.Diagnostics.AddError(
			`Unable to create catalog: `+plan.Name.ValueString()+``, err.Error(),
		)
		return
	}

	// Adding small timeout because catalog ID is not available in read operation otherwise, so AT and FT were failing.
	time.Sleep(5 * time.Second)

	// Set the tf state after create
	state, mapErr := helper.SetStateCatalogFirmware(ctx, cat, plan)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			`Unable to process catalog after create: `+plan.Name.ValueString()+`.`, mapErr.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "firmwareCatalogResource: create end")
}

// Read implements resource.Resource.
func (r *firmwareCatalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "firmwareCatalogResource: read started")
	var currentState models.OmeSingleCatalogResource
	diags := req.State.Get(ctx, &currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_firmware_catalog Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	// Get the ID after create, for whatever reason the create api does not return the actual ID
	// Instead it returns 0.
	// The only way to get the true id is to get all of the catalogs and find the one that matches by name (names are required to be unique for catalogs)
	if currentState.ID.ValueInt64() == 0 {
		id, idErr := helper.GetIDFromNameFirmwareCatalog(omeClient, currentState.Name.ValueString())
		if idErr != nil {
			resp.Diagnostics.AddError(
				`Unable to read catalog id after create: `+currentState.Name.ValueString()+`.`, idErr.Error(),
			)
			return
		}
		currentState.ID = types.Int64Value(id)
	}

	cat, err := helper.GetSpecificCatalogFirmware(omeClient, currentState.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			`Unable to read specific firmware catalog: `+currentState.Name.ValueString()+``, err.Error(),
		)
	}

	// Set the tf state after read
	state, mapErr := helper.SetStateCatalogFirmware(ctx, cat, currentState)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			`Unable to process catalog after read: `+state.Name.ValueString()+`.`, mapErr.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "firmwareCatalogResource: read end")
}

// Delete implements resource.Resource.
func (r *firmwareCatalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "firmwareCatalogResource: delete started")

	var state models.OmeSingleCatalogResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If status is currently "Running" we can't delete.
	// We should show error and wait for it to finish before deleting.
	// Add this error message on the tf size for a more clear error
	if state.Status.ValueString() == "Running" {
		resp.Diagnostics.AddError(
			`Unable to delete catalog, the catalog is currently running an update.`,
			"Once the update of the catalog is complete the catalog can be deleted.",
		)
		return
	}

	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_firmware_catalog Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	err := helper.DeleteCatalogFirmware(omeClient, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			`Unable to delete firmware catalog: `+state.Name.ValueString()+``, err.Error(),
		)
	}
	tflog.Trace(ctx, "firmwareCatalogResource: delete end")
}

// Update implements resource.Resource.
func (r *firmwareCatalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "firmwareCatalogResource: update started")
	var plan models.OmeSingleCatalogResource
	var state models.OmeSingleCatalogResource

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

	valError := helper.ValidateCatalogUpdate(plan, state)

	if valError != nil {
		resp.Diagnostics.AddError(
			`Unable to update catalog, validation error: `, valError.Error(),
		)
		return
	}

	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_firmware_catalog Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	repo := models.CatalogRepository{}
	repoDiags := state.Repository.As(ctx, &repo, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
	resp.Diagnostics.Append(repoDiags...)
	if repoDiags.HasError() {
		return
	}

	updateModel := helper.MakeCatalogJSONModel(state.ID.ValueInt64(), repo.ID.ValueInt64(), plan)

	cat, err := helper.UpdateCatalogFirmware(omeClient, state.ID.ValueInt64(), updateModel)
	if err != nil {
		resp.Diagnostics.AddError(
			`Unable to update catalog: `+state.Name.ValueString()+``, err.Error(),
		)
		return
	}

	// Update tf state after update of catalog
	finalState, mapErr := helper.SetStateCatalogFirmware(ctx, cat, plan)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			`Unable to process catalog after update: `+state.Name.ValueString()+`.`, mapErr.Error(),
		)
		return
	}
	diags := resp.State.Set(ctx, &finalState)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "firmwareCatalogResource: update end")
}

// ImportState implements resource.Resource.
func (r *firmwareCatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var importState models.OmeSingleCatalogResource
	tflog.Trace(ctx, "firmwareCatalogResource: import state started")
	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_firmware_catalog Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	id, coversionErr := strconv.Atoi(req.ID)
	if coversionErr != nil {
		resp.Diagnostics.AddError(
			`Unable to import firmware catalog, id must be an integer: `+req.ID+``, coversionErr.Error(),
		)
	}
	cat, err := helper.GetSpecificCatalogFirmware(omeClient, int64(id))
	if err != nil {
		resp.Diagnostics.AddError(
			`Unable to import firmware catalog: `+req.ID+``, err.Error(),
		)
	}

	// Set the tf state after read
	state, mapErr := helper.SetStateCatalogFirmware(ctx, cat, importState)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			`Unable to process catalog after import: `+state.Name.ValueString()+`.`, mapErr.Error(),
		)
		return
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "firmwareCatalogResource: import state end")
}
