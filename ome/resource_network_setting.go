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
	if data.OmeSessionSetting != nil {
		if ok, err := isSessionConfigValid(data.OmeSessionSetting); !ok {
			resp.Diagnostics.AddAttributeError(
				path.Root("session_setting"),
				"Attribute Error",
				err.Error(),
			)
		}
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
	// session configuration
	if plan.OmeSessionSetting != nil {
		state.OmeSessionSetting, getErr = getSessionSettingState(omeClient)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Session Get Error", getErr.Error(),
			)
		}
		isChange, critical := updateSessionSettingState(&plan, &state, omeClient)
		if !isChange {
			resp.Diagnostics.AddWarning("No Change Detected.", "No change in session setting on the infrastructure.")
		}
		if critical != nil {
			resp.Diagnostics.AddError(
				"OME Session Create Error", critical.Error(),
			)
		}
	}

	// proxy configuration
	if plan.OmeProxySetting != nil {
		state.OmeProxySetting, getErr = getProxySettingState(omeClient)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Get Error", getErr.Error(),
			)
		}
		isChange, critical := updateProxySettingState(&plan, &state, omeClient)
		if !isChange {
			resp.Diagnostics.AddWarning("No Change Detected.", "No change in proxy setting on the infrastructure.")
		}
		if critical != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Create Error", critical.Error(),
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
	var getErr error
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
	if state.OmeSessionSetting != nil {
		state.OmeSessionSetting, getErr = getSessionSettingState(omeClient)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Proxy Get Error", getErr.Error(),
			)
		}
	}
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
	// session configuration 
	if state.OmeSessionSetting != nil && plan.OmeSessionSetting != nil {
		_, critcal := updateSessionSettingState(&plan, &state, omeClient)
		if critcal != nil {
			resp.Diagnostics.AddError(
				"OME Session Update Error", critcal.Error(),
			)
		}
	} else {
		state.OmeSessionSetting = nil 
	}

	// proxy configuration
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

// =============================== session configuration helper function ===============================

func isSessionConfigValid(planSession *models.OmeSessionSetting) (bool, error) {
	if planSession.EnableUniversalTimeout.ValueBool(){
		if planSession.APITimeout.ValueFloat64() > 0 || planSession.GUITimeout.ValueFloat64() > 0 || planSession.SSHTimeout.ValueFloat64() >0 || planSession.SerialTimeout.ValueFloat64() > 0 {
			return false, fmt.Errorf("please validate that the configuration for api_timeout, gui_timeout, ssh_timeout and serial_timeout are unset when enable_universal_timeout option is active")
		}
		if planSession.UniversalTimeout.ValueFloat64() < 1 {
			return false, fmt.Errorf("please ensure universal_timeout is set when enable_universal_timeout option is active")
		}
	} else {
		if planSession.UniversalTimeout.ValueFloat64() > 0 {
			return false, fmt.Errorf("please ensure universal_timeout is unset when enable_universal_timeout option is disable")
		}
	}
	return true, nil 
}

func getSessionSettingState(omeClient *clients.Client) (*models.OmeSessionSetting, error) {
	sessions, err := omeClient.GetNetworkSessions()
	if err != nil {
		return nil, err
	}
	sessionSettingState := buildSessionSettingState(&sessions)
	return sessionSettingState, nil
}

func isSessionConfigValuesChanged(planSession, stateSession *models.OmeSessionSetting) bool {
	return (!planSession.EnableUniversalTimeout.IsUnknown() && !planSession.EnableUniversalTimeout.Equal(stateSession.EnableUniversalTimeout)) ||
		(!planSession.UniversalTimeout.IsUnknown() && planSession.UniversalTimeout.ValueFloat64() != stateSession.UniversalTimeout.ValueFloat64()) ||
		(!planSession.APISession.IsUnknown() && planSession.APISession.ValueInt64() != stateSession.APISession.ValueInt64()) ||
		(!planSession.APITimeout.IsUnknown() && planSession.APITimeout.ValueFloat64() != stateSession.APITimeout.ValueFloat64()) ||
		(!planSession.GUISession.IsUnknown() && planSession.GUISession.ValueInt64() != stateSession.GUISession.ValueInt64()) ||
		(!planSession.GUITimeout.IsUnknown() && planSession.GUITimeout.ValueFloat64() != stateSession.GUITimeout.ValueFloat64()) ||
		(!planSession.SSHSession.IsUnknown() && planSession.SSHSession.ValueInt64() != stateSession.SSHSession.ValueInt64()) ||
		(!planSession.SSHTimeout.IsUnknown() && planSession.SSHTimeout.ValueFloat64() != stateSession.SSHTimeout.ValueFloat64()) ||
		(!planSession.SerialSession.IsUnknown() && planSession.SerialSession.ValueInt64() != stateSession.SerialSession.ValueInt64()) ||
		(!planSession.SerialTimeout.IsUnknown() && planSession.SerialTimeout.ValueFloat64() != stateSession.SerialTimeout.ValueFloat64())
}

func updateSessionSettingState(plan, state *models.OmeNetworkSetting, omeClient *clients.Client) (bool, error) {
	// var universalSession, guiSession, apiSession,sshSession,serialSession models.SessionInfo
	if isSessionConfigValuesChanged(plan.OmeSessionSetting, state.OmeSessionSetting){
		session, err := omeClient.GetNetworkSessions()
		if err != nil {
			return false, err
		}
		payload, err := buildSessionUpdatePayload(plan.OmeSessionSetting,session.SessionList)
		if err != nil {
			return false, err
		}
		newSession, err := omeClient.UpdateNetworkSessions(payload)
		if err != nil {
			return true, err
		}
		state.OmeSessionSetting = buildSessionSettingState(&models.NetworkSessions{
			SessionList: newSession,
		})
		return true, nil 
	}
	return false, nil 
}

func buildSessionUpdatePayload(plan *models.OmeSessionSetting, curr []models.SessionInfo) ([]models.SessionInfo, error) {
	var payload []models.SessionInfo
	const MinuteToMilliSecond = 60000
	SessionTypeMap := map[string]bool{
		"API": false,
		"GUI": false,
		"SSH": false,
		"Serial": false,
	}
	for _, session := range curr {
		if session.SessionType == "API" {
			if plan.APISession.ValueInt64() > 0 {
				session.MaxSessions = int(plan.APISession.ValueInt64())
			}
			if plan.APITimeout.ValueFloat64() > 0 {
				session.SessionTimeout = int(plan.APITimeout.ValueFloat64() * MinuteToMilliSecond)
			}
			SessionTypeMap["API"] = true
		}
		if session.SessionType == "GUI" {
			if plan.GUISession.ValueInt64() > 0 {
				session.MaxSessions = int(plan.GUISession.ValueInt64())
			}
			if plan.GUITimeout.ValueFloat64() > 0 {
				session.SessionTimeout = int(plan.GUITimeout.ValueFloat64() * MinuteToMilliSecond)
			}
			SessionTypeMap["GUI"] = true
		}
		if session.SessionType == "SSH" {
			if plan.SSHSession.ValueInt64() > 0 {
				session.MaxSessions = int(plan.SSHSession.ValueInt64())
			}
			if plan.SSHTimeout.ValueFloat64() > 0 {
				session.SessionTimeout = int(plan.SSHTimeout.ValueFloat64() * MinuteToMilliSecond)
			}
			SessionTypeMap["SSH"] = true
		}
		if session.SessionType == "Serial" {
			if plan.SerialSession.ValueInt64() > 0 {
				session.MaxSessions = int(plan.SerialSession.ValueInt64())
			}
			if plan.SerialTimeout.ValueFloat64() > 0 {
				session.SessionTimeout = int(plan.SerialTimeout.ValueFloat64() * MinuteToMilliSecond)
			}
			SessionTypeMap["Serial"] = true
		}
		if plan.EnableUniversalTimeout.ValueBool() {
			session.SessionTimeout = int(plan.UniversalTimeout.ValueFloat64() * MinuteToMilliSecond)
		} else {
			if session.SessionType == "UniversalTimeout" {
				session.SessionTimeout = -1
			}
		}
		payload = append(payload, session)
	}
	for sessionType, found := range SessionTypeMap {
		if !found && sessionType == "SSH" && (plan.SSHSession.ValueInt64() != 0 || plan.SSHTimeout.ValueFloat64() != 0) {
			return nil, fmt.Errorf("please verify that the SSH Session is unset, as the infrastructure does not provide support for SSH sessions")
		}
		if !found && sessionType == "Serial" && (plan.SerialSession.ValueInt64() != 0 || plan.SerialTimeout.ValueFloat64() != 0){
			return nil, fmt.Errorf("please verify that the Serial Session is unset, as the infrastructure does not provide support for Serial sessions")
		}
		if !found && sessionType == "API" && (plan.APISession.ValueInt64() != 0 || plan.APITimeout.ValueFloat64() != 0){
			return nil, fmt.Errorf("please verify that the API Session is unset, as the infrastructure does not provide support for API sessions")
		}
		if !found && sessionType == "GUI" && (plan.GUISession.ValueInt64() != 0 || plan.GUITimeout.ValueFloat64() != 0){
			return nil, fmt.Errorf("please verify that the GUI Session is unset, as the infrastructure does not provide support for GUI sessions")
		}
	}
	return payload, nil
}

func buildSessionSettingState(sessions *models.NetworkSessions) *models.OmeSessionSetting {
	const MilliSecondToMinute = 60000
	sessionState := &models.OmeSessionSetting{}
	for _, session := range sessions.SessionList {
		if session.SessionType == "UniversalTimeout" {
			if session.SessionTimeout != -1 {
				sessionState.EnableUniversalTimeout = types.BoolValue(true)
				sessionState.UniversalTimeout = types.Float64Value(float64(session.SessionTimeout / MilliSecondToMinute))
			} else {
				sessionState.EnableUniversalTimeout = types.BoolValue(false)
			}
		} 
		// convert the session timeout in millisecond to minute 
		session.SessionTimeout /= MilliSecondToMinute
		if session.SessionType == "GUI" {
			sessionState.GUISession = types.Int64Value(int64(session.MaxSessions))
			sessionState.GUITimeout = types.Float64Value(float64(session.SessionTimeout))
		} 
		if session.SessionType == "API" {
			sessionState.APISession = types.Int64Value(int64(session.MaxSessions))
			sessionState.APITimeout = types.Float64Value(float64(session.SessionTimeout))
		} 
		if session.SessionType == "SSH" {
			sessionState.SSHSession = types.Int64Value(int64(session.MaxSessions))
			sessionState.SSHTimeout = types.Float64Value(float64(session.SessionTimeout))
		} 
		if session.SessionType == "Serial" {
			sessionState.SerialSession = types.Int64Value(int64(session.MaxSessions))
			sessionState.SerialTimeout = types.Float64Value(float64(session.SessionTimeout))
		}
	}
	return sessionState
}

// ============================= proxy configuration helper function =================================

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
