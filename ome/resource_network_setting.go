package ome

import (
	"context"
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
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

func (r *networkSettingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data models.OmeNetworkSetting
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.OmeProxySetting != nil {
		if ok, err := isProxyConfigValid(data.OmeProxySetting); !ok {
			resp.Diagnostics.AddAttributeError(
				path.Root("proxy_setting"),
				"Attribute Error",
				err.Error(),
			)
		}
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *networkSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_network_setting create : Started")
	//Get Plan Data
	var plan, state models.OmeNetworkSetting
	var getErr error
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

	// proxy configuration
	if plan.OmeProxySetting != nil {
		state.OmeProxySetting, getErr = getProxySettingState(omeClient)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Get Error", getErr.Error(),
			)
		}
		isChange, critcal := updateProxySettingState(&plan, &state, omeClient)
		if !isChange {
			resp.Diagnostics.AddWarning("No Change Detected.", "No change in proxy setting on the infrastructure.")
		}
		if critcal != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Create Error", critcal.Error(),
			)
		}
	}

	if state.ID.IsNull() {
		state.ID = types.StringValue("placeholder")
	}
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
	// proxy configuration
	if state.OmeProxySetting != nil {
		proxySettingState, err := getProxySettingState(omeClient)
		if err != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Get Error", err.Error(),
			)
		}
		// save proxy config into terraform state
		proxySettingState.Password = state.OmeProxySetting.Password
		state.OmeProxySetting = proxySettingState
	}

	tflog.Trace(ctx, "resource_network_setting read: finished reading state")
	//Save into State
	if state.ID.IsNull() {
		state.ID = types.StringValue("placeholder")
	}
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
	if state.OmeProxySetting != nil && plan.OmeProxySetting != nil {
		_, critcal := updateProxySettingState(&plan, &state, omeClient)
		if critcal != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Update Error", critcal.Error(),
			)
		}
	} else {
		state.OmeProxySetting = nil
	}

	tflog.Trace(ctx, "resource_network_setting update: finished state update")
	//Save into State
	if state.ID.IsNull() {
		state.ID = types.StringValue("placeholder")
	}
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

func isProxyConfigValid(planProxy *models.OmeProxySetting) (bool, error) {
	if planProxy.EnableProxy.ValueBool() {
		if planProxy.IPAddress.ValueString() == "" || planProxy.ProxyPort.ValueInt64() == 0 {
			return false, fmt.Errorf("please ensure that you set both the IP address and port when enabling the proxy")
		}
		if planProxy.EnableAuthentication.ValueBool() {
			if planProxy.Username.ValueString() == "" || planProxy.Password.ValueString() == "" {
				return false, fmt.Errorf("please ensure that you set both the username and password when enabling the proxy authentication")
			}
		} else if planProxy.Username.ValueString() != "" || planProxy.Password.ValueString() != "" {
			return false, fmt.Errorf("please ensure enable authentication should be set to true before setting username and password")
		}
	} else if planProxy.IPAddress.ValueString() != "" || planProxy.ProxyPort.ValueInt64() > 0 || planProxy.EnableAuthentication.ValueBool() || planProxy.Username.ValueString() != "" || planProxy.Password.ValueString() != "" {
		return false, fmt.Errorf("please ensure enable proxy should be set to true before setting any ome proxy configuration")
	}

	return true, nil
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

func updateProxySettingState(plan, state *models.OmeNetworkSetting, omeClient *clients.Client) (bool, error) {
	if isProxyConfigValuesChanged(plan.OmeProxySetting, state.OmeProxySetting) {
		payloadProxy := models.PayloadProxyConfiguration{
			IPAddress:            plan.OmeProxySetting.IPAddress.ValueString(),
			PortNumber:           int(plan.OmeProxySetting.ProxyPort.ValueInt64()),
			EnableAuthentication: plan.OmeProxySetting.EnableAuthentication.ValueBool(),
			EnableProxy:          plan.OmeProxySetting.EnableProxy.ValueBool(),
			Username:             plan.OmeProxySetting.Username.ValueString(),
			Password:             plan.OmeProxySetting.Password.ValueString(),
		}
		newProxy, err := omeClient.UpdateProxyConfig(payloadProxy)
		if err != nil {
			return true, err
		}
		state.OmeProxySetting = buildProxySettingState(&newProxy)
		if !plan.OmeProxySetting.Password.IsUnknown() {
			state.OmeProxySetting.Password = plan.OmeProxySetting.Password
		}
		return true, nil
	}
	return false, nil
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
