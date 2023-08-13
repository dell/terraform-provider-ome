package ome

import (
	"context"
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &networkSettingResource{}
)

// NewNetworkSettingResource is a helper function to simplify the provider implementation.
func NewNetworkSettingResource() resource.Resource {
	return &networkSettingResource{}
}

// networkSettingResource is the resource implementation.
type networkSettingResource struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *networkSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata returns the resource type name.
func (r *networkSettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "network_setting"
}

// Schema defines the schema for the resource.
func (r *networkSettingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing network_setting on OpenManage Enterprise.",
		Version:             1,
		Attributes:          NetworkSettingSchema(),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *networkSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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
	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_discovery Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	// get proxy config
	proxySettingState, err := getProxySettingState(omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"OME Proxy Create Error", err.Error(),
		)
	}
	if isProxyConfigValuesChanged(plan.OmeProxySetting, proxySettingState) {
		tflog.Trace(ctx, "change detected for proxy setting : "+fmt.Sprintf("%#v", proxySettingState)+fmt.Sprintf("%#v", plan.OmeProxySetting))
		updateProxySettingState, err := updateProxySettingState(&plan, omeClient)
		if err != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Create Error", err.Error(),
			)
		}
		proxySettingState = updateProxySettingState
	} else {
		resp.Diagnostics.AddWarning("No Change Detected.","No change in proxy setting on the infrastructure.")
	}
	// save proxy config into terraform state
	state.OmeProxySetting = proxySettingState
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_setting create: finish")
}

// Read refreshes the Terraform state with the latest data.
func (r *networkSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_network_setting read: started")
	var state models.OmeNetworkSetting
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_discovery Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	// get proxy config
	proxySettingState, err := getProxySettingState(omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"OME Proxy Get Error", err.Error(),
		)
	}
	// save proxy config into terraform state
	state.OmeProxySetting = proxySettingState
	tflog.Trace(ctx, "resource_network_setting read: finished reading state")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_setting read: finished")
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_discovery Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	if isProxyConfigValuesChanged(plan.OmeProxySetting, state.OmeProxySetting) {
		updateProxySettingState, err := updateProxySettingState(&plan, omeClient)
		if err != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Update Error", err.Error(),
			)
		}
		// save proxy config into terraform state
		state.OmeProxySetting = updateProxySettingState
	}
	tflog.Trace(ctx, "resource_network_setting update: finished state update")
	//Save into State

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_setting update: finished")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func getProxySettingState(omeClient *clients.Client) (*models.OmeProxySetting, error) {
	proxy, err := omeClient.GetProxyConfig()
	if err != nil {
		return nil, err
	}
	proxySettingState := buildProxySettingState(&proxy)
	return proxySettingState, err
}

func isProxyConfigValuesChanged(planProxy, stateProxy *models.OmeProxySetting) bool {
	return (!planProxy.EnableProxy.IsUnknown() && !stateProxy.EnableProxy.Equal(planProxy.EnableProxy)) ||
		(!planProxy.ProxyPort.IsUnknown() && planProxy.ProxyPort.ValueInt64() != stateProxy.ProxyPort.ValueInt64()) ||
		(!planProxy.IPAddress.IsUnknown() && planProxy.IPAddress.ValueString() != stateProxy.IPAddress.ValueString()) ||
		(!planProxy.Username.IsUnknown() && planProxy.Username.ValueString() != stateProxy.Username.ValueString()) ||
		(!planProxy.Password.IsUnknown() && planProxy.Password.ValueString() != stateProxy.Password.ValueString()) ||
		(!planProxy.EnableAuthentication.IsUnknown() && !stateProxy.EnableAuthentication.Equal(planProxy.EnableAuthentication))
}

func updateProxySettingState(plan *models.OmeNetworkSetting, omeClient *clients.Client) (*models.OmeProxySetting, error) {
	payloadProxy := models.PayloadProxyConfiguration{
		IPAddress:            plan.OmeProxySetting.IPAddress.ValueString(),
		PortNumber:           int(plan.OmeProxySetting.ProxyPort.ValueInt64()),
		EnableAuthentication: plan.OmeProxySetting.EnableAuthentication.ValueBool(),
		EnableProxy:          plan.OmeProxySetting.EnableProxy.ValueBool(),
		Username:             plan.OmeProxySetting.Username.ValueString(),
		Password:             plan.OmeProxySetting.Password.ValueString(),
	}
	newProxy, err := omeClient.UpdateProxyConfig(payloadProxy)
	proxySettingState := buildProxySettingState(&newProxy)
	if err != nil {
		return &models.OmeProxySetting{}, err
	}
	return proxySettingState, nil
}

func buildProxySettingState(proxy *models.ProxyConfiguration) *models.OmeProxySetting {
	proxySettingState := &models.OmeProxySetting{
		EnableProxy:          types.BoolValue(proxy.EnableProxy),
		IPAddress:            types.StringValue(proxy.IPAddress),
		ProxyPort:            types.Int64Value(int64(proxy.PortNumber)),
		EnableAuthentication: types.BoolValue(proxy.EnableAuthentication),
		Username:             types.StringValue(proxy.Username),
		Password:             types.StringValue(proxy.Password),
	}
	return proxySettingState
}