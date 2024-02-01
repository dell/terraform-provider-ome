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
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/helper"
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
	resp.TypeName = req.ProviderTypeName + "appliance_network"
}

// Schema defines the schema for the resource.
func (r *networkSettingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This terraform resource is used to manage Appliance Network Settings on OME." +
			"We can Create, Update and Delete OME Appliance Network Settings using this resource.",
		Version:    1,
		Attributes: NetworkSettingSchema(),
	}
}

func (r *networkSettingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data models.OmeNetworkSetting
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.OmeTimeSetting != nil {
		if ok, err := isTimeConfigValid(data.OmeTimeSetting); !ok {
			resp.Diagnostics.AddAttributeError(
				path.Root("time_setting"),
				"Attribute Error",
				err.Error(),
			)
		}
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
	if data.OmeAdapterSetting != nil {
		if err := isAdapterConfigValid(data.OmeAdapterSetting); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("adapter_setting"),
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
	// time configuration
	if plan.OmeTimeSetting != nil {
		state.OmeTimeSetting, getErr = getTimeSettingState(omeClient)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Time Get Error", getErr.Error(),
			)
		}
		isChange, critical := updateTimeSettingState(&plan, &state, omeClient)
		if !isChange {
			resp.Diagnostics.AddWarning("No Change Detected.", "No change in time setting on the infrastructure.")
		}
		if critical != nil {
			resp.Diagnostics.AddError(
				"OME Time Create Error", critical.Error(),
			)
		}
	}

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

	// adapter configuration
	if plan.OmeAdapterSetting != nil {
		state.OmeAdapterSetting, getErr = getAdapterSettingState(omeClient, plan.OmeAdapterSetting)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Adapter Get Error", getErr.Error(),
			)
		}
		err := updateAdapterSettingState(ctx, &plan, &state, omeClient)
		if err != nil {
			resp.Diagnostics.AddError(
				"OME Adapter Create Error", err.Error(),
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

	if state.OmeTimeSetting != nil {
		state.OmeTimeSetting, getErr = getTimeSettingState(omeClient)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Time Get Error", getErr.Error(),
			)
		}
	}

	// session configuration
	if state.OmeSessionSetting != nil {
		state.OmeSessionSetting, getErr = getSessionSettingState(omeClient)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Session Get Error", getErr.Error(),
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

	// adapter configuration
	if state.OmeAdapterSetting != nil {
		state.OmeAdapterSetting, getErr = getAdapterSettingState(omeClient, state.OmeAdapterSetting)
		if getErr != nil {
			resp.Diagnostics.AddError(
				"OME Adapter Get Error", getErr.Error(),
			)
		}
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
	// time configuration
	if plan.OmeTimeSetting != nil {
		_, critical := updateTimeSettingState(&plan, &state, omeClient)
		if critical != nil {
			resp.Diagnostics.AddError(
				"OME Time Update Error", critical.Error(),
			)
		}
	} else {
		state.OmeTimeSetting = nil
	}

	// session configuration
	if plan.OmeSessionSetting != nil {
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

	// adapter configuration
	if plan.OmeAdapterSetting != nil {
		err := updateAdapterSettingState(ctx, &plan, &state, omeClient)
		if err != nil {
			resp.Diagnostics.AddError(
				"OME Adapter Update Error", err.Error(),
			)
		}
	} else {
		state.OmeAdapterSetting = nil
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

// ============================== adapter configuration helper function ================================

func isAdapterConfigValid(plan *models.OmeAdapterSetting) error {
	if plan.IPV4Config != nil {
		if !plan.IPV4Config.EnableIPv4.IsUnknown() && plan.IPV4Config.EnableIPv4.ValueBool() {
			if !plan.IPV4Config.EnableDHCP.IsUnknown() {
				check := plan.IPV4Config.StaticIPAddress.IsNull() && plan.IPV4Config.StaticSubnetMask.IsNull() && plan.IPV4Config.StaticGateway.IsNull()
				if plan.IPV4Config.EnableDHCP.ValueBool() {
					if !check {
						return fmt.Errorf("static_ip_address / static_subnet_mask / static_gateway should not be set when enable_dhcp is active")
					}
				} else {
					if check {
						return fmt.Errorf("static_ip_address / static_subnet_mask / static_gateway are required when enable_dhcp is inactive")
					}
				}
			}
			if !plan.IPV4Config.UseDHCPforDNSServerNames.IsUnknown() {
				check := plan.IPV4Config.StaticPreferredDNSServer.IsNull() && plan.IPV4Config.StaticAlternateDNSServer.IsNull()
				if plan.IPV4Config.UseDHCPforDNSServerNames.ValueBool() && !check {
					return fmt.Errorf("static_preferred_dns_server / static_alternate_dns_server should not be set when use_dhcp_for_dns_server_names is active")
				}
			}
		}
	}

	if plan.IPV6Config != nil {
		if !plan.IPV6Config.EnableIPv6.IsUnknown() && plan.IPV6Config.EnableIPv6.ValueBool() {
			if !plan.IPV6Config.EnableAutoConfiguration.IsUnknown() {
				check := plan.IPV6Config.StaticIPAddress.IsNull() && plan.IPV6Config.StaticPrefixLength.IsNull() && plan.IPV6Config.StaticGateway.IsNull()
				if plan.IPV6Config.EnableAutoConfiguration.ValueBool() {
					if !(check) {
						return fmt.Errorf("static_ip_address / static_prefix_length / static_gateway should not be set when enable_auto_configuration is disable")
					}
				} else {
					if check {
						return fmt.Errorf("static_ip_address / static_prefix_length / static_gateway are required when enable_auto_configuration is disable")
					}
				}
			}

			if !plan.IPV6Config.UseDHCPforDNSServerNames.IsUnknown() {
				check := plan.IPV6Config.StaticPreferredDNSServer.IsNull() && plan.IPV6Config.StaticAlternateDNSServer.IsNull()
				if plan.IPV6Config.UseDHCPforDNSServerNames.ValueBool() && !(check) {
					return fmt.Errorf("static_preferred_dns_server / static_alternate_dns_server should not be set when use_dhcp_for_dns_servers_name is active")
				}
			}
		}
	}

	if plan.ManagementVLAN != nil {
		if !plan.ManagementVLAN.EnableVLAN.IsUnknown() && !plan.ManagementVLAN.EnableVLAN.ValueBool() && !plan.ManagementVLAN.ID.IsNull() {
			return fmt.Errorf("please validate enable_vlan is true, when id is set")
		}
	}

	if plan.DNSConfig != nil {
		if !plan.DNSConfig.RegisterWithDNS.IsUnknown() && !plan.DNSConfig.RegisterWithDNS.ValueBool() && !plan.DNSConfig.DNSName.IsNull() {
			return fmt.Errorf("please validate register_with_dn is true, when dns_name is set")
		}

		if !plan.DNSConfig.RegisterWithDNS.IsUnknown() && plan.DNSConfig.UseDHCPforDNSServerNames.ValueBool() && !plan.DNSConfig.DNSDomainName.IsNull() {
			return fmt.Errorf("please validate use_dhcp_for_dns_server_names is false, when dns_domain_name is set")
		}
	}

	return nil
}

func updateAdapterSettingState(ctx context.Context, plan, state *models.OmeNetworkSetting, omeClient *clients.Client) error {
	var newOmeIP string
	currentAdapter, err := omeClient.GetNetworkAdapterConfigByInterface(state.OmeAdapterSetting.InterfaceName.ValueString())
	if err != nil {
		return err
	}
	adapter := models.UpdateNetworkAdapterSetting{
		InterfaceName:    plan.OmeAdapterSetting.InterfaceName.ValueString(),
		ProfileName:      plan.OmeAdapterSetting.InterfaceName.ValueString(),
		EnableNIC:        plan.OmeAdapterSetting.EnableNic.ValueBool(),
		Delay:            int(plan.OmeAdapterSetting.RebootDelay.ValueInt64()),
		PrimaryInterface: true,
	}
	if plan.OmeAdapterSetting.IPV4Config != nil {
		adapter.Ipv4Configuration = models.Ipv4Configuration{
			Enable:                   plan.OmeAdapterSetting.IPV4Config.EnableIPv4.ValueBool(),
			EnableDHCP:               plan.OmeAdapterSetting.IPV4Config.EnableDHCP.ValueBool(),
			StaticIPAddress:          plan.OmeAdapterSetting.IPV4Config.StaticIPAddress.ValueString(),
			StaticSubnetMask:         plan.OmeAdapterSetting.IPV4Config.StaticSubnetMask.ValueString(),
			StaticGateway:            plan.OmeAdapterSetting.IPV4Config.StaticGateway.ValueString(),
			UseDHCPForDNSServerNames: plan.OmeAdapterSetting.IPV4Config.UseDHCPforDNSServerNames.ValueBool(),
			StaticPreferredDNSServer: plan.OmeAdapterSetting.IPV4Config.StaticPreferredDNSServer.ValueString(),
			StaticAlternateDNSServer: plan.OmeAdapterSetting.IPV4Config.StaticAlternateDNSServer.ValueString(),
		}
		newOmeIP = plan.OmeAdapterSetting.IPV4Config.StaticIPAddress.ValueString()
	} else {
		adapter.Ipv4Configuration = currentAdapter.Ipv4Configuration
	}

	if plan.OmeAdapterSetting.IPV6Config != nil {
		adapter.Ipv6Configuration = models.Ipv6Configuration{
			Enable:                   plan.OmeAdapterSetting.IPV6Config.EnableIPv6.ValueBool(),
			EnableAutoConfiguration:  plan.OmeAdapterSetting.IPV6Config.EnableAutoConfiguration.ValueBool(),
			StaticIPAddress:          plan.OmeAdapterSetting.IPV6Config.StaticIPAddress.ValueString(),
			StaticPrefixLength:       int(plan.OmeAdapterSetting.IPV6Config.StaticPrefixLength.ValueInt64()),
			StaticGateway:            plan.OmeAdapterSetting.IPV6Config.StaticGateway.ValueString(),
			UseDHCPForDNSServerNames: plan.OmeAdapterSetting.IPV6Config.UseDHCPforDNSServerNames.ValueBool(),
			StaticPreferredDNSServer: plan.OmeAdapterSetting.IPV6Config.StaticPreferredDNSServer.ValueString(),
			StaticAlternateDNSServer: plan.OmeAdapterSetting.IPV6Config.StaticAlternateDNSServer.ValueString(),
		}
	} else {
		adapter.Ipv6Configuration = currentAdapter.Ipv6Configuration
	}

	if plan.OmeAdapterSetting.DNSConfig != nil {
		adapter.DNSConfiguration = models.DNSConfiguration{
			RegisterWithDNS:         plan.OmeAdapterSetting.DNSConfig.RegisterWithDNS.ValueBool(),
			DNSName:                 plan.OmeAdapterSetting.DNSConfig.DNSName.ValueString(),
			UseDHCPForDNSDomainName: plan.OmeAdapterSetting.DNSConfig.UseDHCPforDNSServerNames.ValueBool(),
			DNSDomainName:           plan.OmeAdapterSetting.DNSConfig.DNSDomainName.ValueString(),
		}
	} else {
		adapter.DNSConfiguration = currentAdapter.DNSConfiguration
	}

	if plan.OmeAdapterSetting.ManagementVLAN != nil {
		adapter.ManagementVLAN = models.ManagementVLAN{
			EnableVLAN: plan.OmeAdapterSetting.ManagementVLAN.EnableVLAN.ValueBool(),
			ID:         int(plan.OmeAdapterSetting.ManagementVLAN.ID.ValueInt64()),
		}
	} else {
		adapter.ManagementVLAN = currentAdapter.ManagementVLAN
	}

	newJob, err := omeClient.UpdateNetworkAdapterConfig(adapter)
	if err != nil {
		return err
	}
	err = helper.NetworkJobRunner(ctx, omeClient, newJob.ID)
	if err != nil {
		if newOmeIP != "" {
			currentURL := omeClient.GetURL()
			omeClient.SetURL(fmt.Sprintf("https://%s:%d", newOmeIP, 443))
			err = helper.NetworkJobRunner(ctx, omeClient, newJob.ID)
			if err != nil {
				return err
			}
			state.OmeAdapterSetting, err = getAdapterSettingState(omeClient, plan.OmeAdapterSetting)
			if err != nil {
				return err
			}
			omeClient.SetURL(currentURL)
			return nil
		}
		return err
	}
	state.OmeAdapterSetting, err = getAdapterSettingState(omeClient, plan.OmeAdapterSetting)
	if err != nil {
		return err
	}
	return nil
}

func getAdapterSettingState(omeClient *clients.Client, plan *models.OmeAdapterSetting) (*models.OmeAdapterSetting, error) {
	if plan.InterfaceName.ValueString() == "" {
		return nil, fmt.Errorf("interface not found")
	}
	adapter, err := omeClient.GetNetworkAdapterConfigByInterface(plan.InterfaceName.ValueString())
	if err != nil {
		return nil, err
	}
	adapterSettingState := &models.OmeAdapterSetting{
		EnableNic:     types.BoolValue(adapter.EnableNIC),
		InterfaceName: types.StringValue(adapter.InterfaceName),
		RebootDelay:   types.Int64Value(int64(adapter.Delay)),
	}
	if plan.IPV4Config != nil {
		adapterSettingState.IPV4Config = &models.OmeIPv4Config{
			EnableIPv4:               types.BoolValue(adapter.Ipv4Configuration.Enable),
			EnableDHCP:               types.BoolValue(adapter.Ipv4Configuration.EnableDHCP),
			StaticIPAddress:          types.StringValue(adapter.Ipv4Configuration.StaticIPAddress),
			StaticSubnetMask:         types.StringValue(adapter.Ipv4Configuration.StaticSubnetMask),
			StaticGateway:            types.StringValue(adapter.Ipv4Configuration.StaticGateway),
			UseDHCPforDNSServerNames: types.BoolValue(adapter.Ipv4Configuration.UseDHCPForDNSServerNames),
			StaticPreferredDNSServer: types.StringValue(adapter.Ipv4Configuration.StaticPreferredDNSServer),
			StaticAlternateDNSServer: types.StringValue(adapter.Ipv4Configuration.StaticAlternateDNSServer),
		}
	} else {
		adapterSettingState.IPV4Config = nil
	}
	if plan.IPV6Config != nil {
		adapterSettingState.IPV6Config = &models.OmeIPv6Config{
			EnableIPv6:               types.BoolValue(adapter.Ipv6Configuration.Enable),
			EnableAutoConfiguration:  types.BoolValue(adapter.Ipv6Configuration.EnableAutoConfiguration),
			StaticIPAddress:          types.StringValue(adapter.Ipv6Configuration.StaticIPAddress),
			StaticPrefixLength:       types.Int64Value(int64(adapter.Ipv6Configuration.StaticPrefixLength)),
			StaticGateway:            types.StringValue(adapter.Ipv6Configuration.StaticGateway),
			UseDHCPforDNSServerNames: types.BoolValue(adapter.Ipv6Configuration.UseDHCPForDNSServerNames),
			StaticPreferredDNSServer: types.StringValue(adapter.Ipv6Configuration.StaticPreferredDNSServer),
			StaticAlternateDNSServer: types.StringValue(adapter.Ipv6Configuration.StaticAlternateDNSServer),
		}
	} else {
		adapterSettingState.IPV6Config = nil
	}
	if plan.DNSConfig != nil {
		adapterSettingState.DNSConfig = &models.OmeDNSConfig{
			RegisterWithDNS:          types.BoolValue(adapter.DNSConfiguration.RegisterWithDNS),
			UseDHCPforDNSServerNames: types.BoolValue(adapter.DNSConfiguration.UseDHCPForDNSDomainName),
			DNSName:                  types.StringValue(adapter.DNSConfiguration.DNSName),
			DNSDomainName:            types.StringValue(adapter.DNSConfiguration.DNSDomainName),
		}
	} else {
		adapterSettingState.DNSConfig = nil
	}
	if plan.ManagementVLAN != nil {
		adapterSettingState.ManagementVLAN = &models.OmeManagementVLAN{
			EnableVLAN: types.BoolValue(adapter.ManagementVLAN.EnableVLAN),
			ID:         types.Int64Value(int64(adapter.ManagementVLAN.ID)),
		}
	} else {
		adapterSettingState.ManagementVLAN = nil
	}
	return adapterSettingState, nil
}

// =============================== time configuration helper function ==================================
func isTimeConfigValid(planTime *models.OmeTimeSetting) (bool, error) {
	if planTime.EnableNTP.IsUnknown() {
		return true, nil
	}
	if planTime.EnableNTP.ValueBool() {
		if !planTime.SystemTime.IsNull() {
			return false, fmt.Errorf("system_time should not be set when enable_ntp is active")
		}
		if planTime.PrimaryNTPAddress.IsNull() {
			return false, fmt.Errorf("primary_ntp_address should be set when enable_ntp is active")
		}
		return true, nil
	}

	if !(planTime.PrimaryNTPAddress.IsNull() && planTime.SecondaryNTPAddress1.IsNull() && planTime.SecondaryNTPAddress2.IsNull()) {
		return false, fmt.Errorf("primary_ntp_address, secondary_ntp_address1 and secondary_ntp_address2 should not be set when enable_ntp is disable")
	}

	if planTime.SystemTime.IsNull() {
		return false, fmt.Errorf("system_time should be set when enable_ntp is disable")
	}

	return true, nil
}

func getTimeSettingState(omeClient *clients.Client) (*models.OmeTimeSetting, error) {
	time, err := omeClient.GetTimeConfiguration()
	if err != nil {
		return nil, err
	}
	timeSettingState := buildTimeSettingState(&time)
	return timeSettingState, nil
}

func isTimeConfigValuesChanged(planTime, stateTime *models.OmeTimeSetting) bool {
	return (!planTime.EnableNTP.IsUnknown() && !planTime.EnableNTP.Equal(stateTime.EnableNTP)) ||
		(!planTime.TimeZone.IsUnknown() && planTime.TimeZone.ValueString() != stateTime.SystemTime.ValueString()) ||
		(!planTime.SystemTime.IsUnknown() && planTime.SystemTime.ValueString() != stateTime.SystemTime.ValueString()) ||
		(!planTime.PrimaryNTPAddress.IsUnknown() && planTime.PrimaryNTPAddress.ValueString() != stateTime.PrimaryNTPAddress.ValueString()) ||
		(!planTime.SecondaryNTPAddress1.IsUnknown() && planTime.SecondaryNTPAddress1.ValueString() != stateTime.SecondaryNTPAddress1.ValueString()) ||
		(!planTime.SecondaryNTPAddress2.IsUnknown() && planTime.SecondaryNTPAddress2.ValueString() != stateTime.SecondaryNTPAddress2.ValueString())
}

func updateTimeSettingState(plan, state *models.OmeNetworkSetting, omeClient *clients.Client) (bool, error) {
	if isTimeConfigValuesChanged(plan.OmeTimeSetting, state.OmeTimeSetting) {
		payload := models.TimeConfig{
			TimeZone:             plan.OmeTimeSetting.TimeZone.ValueString(),
			EnableNTP:            plan.OmeTimeSetting.EnableNTP.ValueBool(),
			PrimaryNTPAddress:    plan.OmeTimeSetting.PrimaryNTPAddress.ValueString(),
			SecondaryNTPAddress1: plan.OmeTimeSetting.SecondaryNTPAddress1.ValueString(),
			SecondaryNTPAddress2: plan.OmeTimeSetting.SecondaryNTPAddress2.ValueString(),
			SystemTime:           plan.OmeTimeSetting.SystemTime.ValueString(),
		}
		newTime, err := omeClient.UpdateTimeConfiguration(payload)
		if err != nil {
			return true, err
		}
		state.OmeTimeSetting = buildTimeSettingState(&newTime)
		if !plan.OmeTimeSetting.SystemTime.IsNull() {
			state.OmeTimeSetting.SystemTime = plan.OmeTimeSetting.SystemTime
		}
		if plan.OmeTimeSetting.SystemTime.IsUnknown() {
			state.OmeTimeSetting.SystemTime = types.StringValue("")
		}
		return true, nil
	}
	return false, nil
}

func buildTimeSettingState(time *models.TimeConfig) *models.OmeTimeSetting {
	timeState := &models.OmeTimeSetting{
		EnableNTP:            types.BoolValue(time.EnableNTP),
		SystemTime:           types.StringValue(time.SystemTime),
		TimeZone:             types.StringValue(time.TimeZone),
		PrimaryNTPAddress:    types.StringValue(time.PrimaryNTPAddress),
		SecondaryNTPAddress1: types.StringValue(time.SecondaryNTPAddress1),
		SecondaryNTPAddress2: types.StringValue(time.SecondaryNTPAddress2),
	}
	return timeState
}

// =============================== session configuration helper function ===============================

func isSessionConfigValid(planSession *models.OmeSessionSetting) (bool, error) {
	if planSession.EnableUniversalTimeout.IsUnknown() {
		return true, nil
	}
	if planSession.EnableUniversalTimeout.ValueBool() {
		if !(planSession.APITimeout.IsNull() && planSession.GUITimeout.IsNull() && planSession.SSHTimeout.IsNull() && planSession.SerialTimeout.IsNull()) {
			return false, fmt.Errorf("api_timeout, gui_timeout, ssh_timeout and serial_timeout should not be set when enable_universal_timeout option is active")
		}
		if planSession.UniversalTimeout.IsNull() {
			return false, fmt.Errorf("universal_timeout should be set when enable_universal_timeout option is active")
		}
		return true, nil
	}
	if !planSession.UniversalTimeout.IsNull() {
		return false, fmt.Errorf("universal_timeout should not be set when enable_universal_timeout option is disable")
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
	if isSessionConfigValuesChanged(plan.OmeSessionSetting, state.OmeSessionSetting) {
		session, err := omeClient.GetNetworkSessions()
		if err != nil {
			return false, err
		}
		payload, err := buildSessionUpdatePayload(plan.OmeSessionSetting, session.SessionList)
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
		"API":    false,
		"GUI":    false,
		"SSH":    false,
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
		if !found && sessionType == "Serial" && (plan.SerialSession.ValueInt64() != 0 || plan.SerialTimeout.ValueFloat64() != 0) {
			return nil, fmt.Errorf("please verify that the Serial Session is unset, as the infrastructure does not provide support for Serial sessions")
		}
		if !found && sessionType == "API" && (plan.APISession.ValueInt64() != 0 || plan.APITimeout.ValueFloat64() != 0) {
			return nil, fmt.Errorf("please verify that the API Session is unset, as the infrastructure does not provide support for API sessions")
		}
		if !found && sessionType == "GUI" && (plan.GUISession.ValueInt64() != 0 || plan.GUITimeout.ValueFloat64() != 0) {
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
	if planProxy.EnableProxy.IsUnknown() {
		return true, nil
	}
	if planProxy.EnableProxy.ValueBool() {
		if planProxy.IPAddress.IsNull() || planProxy.ProxyPort.IsNull() {
			return false, fmt.Errorf("both IP address and port are required when enabling proxy")
		}

		if !(planProxy.EnableAuthentication.IsNull() || planProxy.EnableAuthentication.IsUnknown()) {
			if planProxy.EnableAuthentication.ValueBool() {
				if planProxy.Username.IsNull() || planProxy.Password.IsNull() {
					return false, fmt.Errorf("both username and password are required when enabling proxy authentication")
				}
				return true, nil
			}
			if !(planProxy.Username.IsNull() && planProxy.Password.IsNull()) {
				return false, fmt.Errorf("enable authentication should be set to true before setting username and password")
			}
		}
		return true, nil
	}
	if !(planProxy.IPAddress.IsNull() && planProxy.ProxyPort.IsNull() && planProxy.EnableAuthentication.IsNull() && planProxy.Username.IsNull() && planProxy.Password.IsNull()) {
		return false, fmt.Errorf("enable proxy should be set to true before setting any ome proxy configuration")
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
