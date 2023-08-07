package ome

import (
	"context"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &network_settingResource{}
)

// Newnetwork_settingResource is a helper function to simplify the provider implementation.
func Newnetwork_settingResource() resource.Resource {
	return &network_settingResource{}
}

// network_settingResource is the resource implementation.
type network_settingResource struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *network_settingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata returns the resource type name.
func (r *network_settingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_setting"
}

// Schema defines the schema for the resource.
func (r *network_settingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing network_setting on OpenManage Enterprise.",
		Version:             1,
		Attributes:          OmeNetworkSettingSchema(),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *network_settingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_network_setting create : Started")
	//Get Plan Data
	var plan, state models.OmeNetworkSetting
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "resource_network_setting create: updating state finished, saving ...")
	// Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_setting create: finish")
}

// Read refreshes the Terraform state with the latest data.
func (r *network_settingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_network_setting read: started")
	var state models.OmeNetworkSetting
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "resource_network_setting read: finished reading state")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_setting read: finished")
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *network_settingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get state Data
	tflog.Trace(ctx, "resource_network_setting update: started")
	var state, plan models.OmeNetworkSetting
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan Data
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "resource_network_setting update: finished state update")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_setting update: finished")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *network_settingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_network_setting delete: started")
	// Get State Data
	var state models.OmeNetworkSetting
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_network_setting delete: finished")
}
