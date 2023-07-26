package ome

import (
	"context"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &discoveryResource{}
)

// NewDiscoveryResource is a helper function to simplify the provider implementation.
func NewDiscoveryResource() resource.Resource {
	return &discoveryResource{}
}

// discoveryResource is the resource implementation.
type discoveryResource struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *discoveryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata returns the resource type name.
func (r *discoveryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "discovery"
}

// Schema defines the schema for the resource.
func (r *discoveryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing discovery on OpenManage Enterprise.",
		Version:             1,
		Attributes:          DiscoveryJobSchema(),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *discoveryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_discovery create : Started")
	//Get Plan Data
	var plan, state models.OmeDiscoveryJob
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "resource_discovery create: updating state finished, saving ...")
	// Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_discovery create: finish")
}

// Read refreshes the Terraform state with the latest data.
func (r *discoveryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_discovery read: started")
	var state models.OmeDiscoveryJob
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	omeClient, d := r.p.createOMESession(ctx, "resource_discovery Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	_, err := omeClient.GetDiscoveryJobByGroupID(state.DiscoveryJobID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrReadDiscovery, err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "resource_discovery read: finished reading state")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_discovery read: finished")
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *discoveryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get state Data
	tflog.Trace(ctx, "resource_discovery update: started")
	var state, plan models.OmeDiscoveryJob
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

	tflog.Trace(ctx, "resource_discovery update: finished state update")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_discovery update: finished")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *discoveryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_discovery delete: started")
	// Get State Data
	var state models.OmeDiscoveryJob
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and differ the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_discovery Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	ddj := models.DiscoveryJobDeletePayload{
		DiscoveryGroupIds: []int{
			int(state.DiscoveryJobID.ValueInt64()),
		},
	}
	tflog.Debug(ctx, "delete group id :", map[string]interface{}{"ids": ddj})
	status, err := omeClient.DeleteDiscoveryJob(ddj)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrDeleteDiscovery,
			err.Error(),
		)
	}
	tflog.Trace(ctx, "resource_discovery delete: finished with status "+status)
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_discovery delete: finished")
}

func (r *discoveryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("discovery_job_id"), req, resp)
}
