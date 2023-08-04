package ome

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/mapstructure"
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

func (r *discoveryResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data models.OmeDiscoveryJob
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Schedule.ValueString() == "RunNow" {
		if !data.Cron.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("cron"),
				"Attribute Error",
				"With Schedule as RunNow, CRON can't be set.",
			)
		}
	}

	if len(data.DiscoveryConfigTargets) > 0 {
		for idx, dct := range data.DiscoveryConfigTargets {
			idxError := "Inappropriate value for attribute `discovery_config_targets`: element: " + strconv.Itoa(idx) + "\n"
			if dct.Redfish == nil && dct.SNMP == nil && dct.SSH == nil && dct.WSMAN == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("discovery_config_targets").AtListIndex(idx),
					"Attribute Error",
					idxError+"Atleast one of protocol should be configured for the discovery targets.")
			}
			if len(dct.DeviceType) > 0 {
				for idx, dt := range dct.DeviceType {
					currDT := dt.ValueString()
					if !isDeviceType(currDT) {
						resp.Diagnostics.AddAttributeError(
							path.Root("discovery_config_targets").AtListIndex(idx).AtName("device_type"),
							"Attribute Error",
							idxError+"The device type list should contain the following values: `SERVER`, `CHASSIS`, `NETWORK SWITCH`, and `STORAGE`.")
					}
				}
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("discovery_config_targets").AtListIndex(idx).AtName("device_type"),
					"Attribute Error",
					idxError+"Atleast one of device type should be configured. ")
			}
			if len(dct.NetworkAddressDetail) == 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("discovery_config_targets").AtListIndex(idx).AtName("network_address_detail"),
					"Attribute Error",
					idxError+"Atleast one of network address detail should be configured. ")
			}
		}
	} else {
		resp.Diagnostics.AddAttributeError(
			path.Root("discovery_config_targets"),
			"Attribute Error",
			"Define at least one discovery configuration target in the list.",
		)
	}
}

func isDeviceType(str string) bool {
	list := []string{"SERVER", "CHASSIS", "NETWORK SWITCH", "STORAGE"}
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
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

	// Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_discovery Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	discoveryPayload := getDiscoveryPayload(ctx, &plan, nil)

	tflog.Debug(ctx, "resource_discovery create Creating Discovery Job", map[string]interface{}{
		"Create Discovery Request": discoveryPayload,
	})

	cDiscovery, err := omeClient.CreateDiscoveryJob(discoveryPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateDiscovery, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_discovery : create Finished creating discovery")
	tflog.Trace(ctx, "resource_discovery : create Fetching discovery id for a discovery")
	state = discoveryState(ctx, cDiscovery, plan)
	// Save into State
	diags = resp.State.Set(ctx, &state)
	tflog.Trace(ctx, "resource_discovery create: updating state finished, saving ...")
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
	id, _ := strconv.Atoi(state.DiscoveryJobID.ValueString())
	respDiscovery, err := omeClient.GetDiscoveryJobByGroupID(int64(id))
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrReadDiscovery, err.Error(),
		)
		return
	}
	state = discoveryState(ctx, respDiscovery, state)
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

	if !reflect.DeepEqual(state, plan) {
		//Create Session and differ the remove session
		omeClient, d := r.p.createOMESession(ctx, "resource_discovery Update")
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}
		defer omeClient.RemoveSession()
		discoveryPayload := getDiscoveryPayload(ctx, &plan, &state)
		tflog.Trace(ctx, "resource_discovery update Discovery")
		tflog.Debug(ctx, "resource_discovery update Discovery", map[string]interface{}{
			"Create Discovery Request": discoveryPayload,
		})
		respDiscovery, err := omeClient.UpdateDiscoveryJob(discoveryPayload)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateDiscovery, err.Error(),
			)
			return
		}
		state = discoveryState(ctx, respDiscovery, plan)
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
	id, _ := strconv.Atoi(state.DiscoveryJobID.ValueString())
	ddj := models.DiscoveryJobDeletePayload{
		DiscoveryGroupIds: []int{id},
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getDiscoveryPayload(ctx context.Context, plan *models.OmeDiscoveryJob, state *models.OmeDiscoveryJob) (payload models.DiscoveryJob) {
	var id int
	if state != nil {
		id, _ = strconv.Atoi(state.DiscoveryJobID.ValueString())
	}
	payload = models.DiscoveryJob{
		CommunityString:               plan.CommunityString.ValueBool(),
		DiscoveryConfigGroupID:        id,
		DiscoveryConfigGroupName:      plan.DiscoveryJobName.ValueString(),
		DiscoveryStatusEmailRecipient: plan.EmailRecipient.ValueString(),
		TrapDestination:               plan.TrapDestination.ValueBool(),
	}

	if plan.Schedule.ValueString() == "RunLater" {
		payload.Schedule = models.ScheduleJob{
			RunNow:    false,
			RunLater:  true,
			Cron:      plan.Cron.ValueString(),
			StartTime: "",
			EndTime:   "",
		}
	} else {
		payload.Schedule = models.ScheduleJob{
			RunNow:    true,
			RunLater:  false,
			Cron:      "startnow",
			StartTime: "",
			EndTime:   "",
		}
	}
	if len(plan.DiscoveryConfigTargets) > 0 {
		for _, dct := range plan.DiscoveryConfigTargets {
			var dcm models.DiscoveryConfigModels
			for _, dt := range dct.DeviceType {
				deviceMap := map[string]int{
					"SERVER":         1000,
					"NETWORK SWITCH": 7000,
					"CHASSIS":        2000,
					"STORAGE":        5000,
				}
				dcm.DeviceType = append(dcm.DeviceType, deviceMap[dt.ValueString()])
			}
			for _, networkAddress := range dct.NetworkAddressDetail {
				network := models.DiscoveryConfigTargets{
					NetworkAddressDetail: networkAddress.ValueString(),
					AddressType:          30,
					Disabled:             false,
					Exclude:              false,
				}
				dcm.DiscoveryConfigTargets = append(dcm.DiscoveryConfigTargets, network)
			}
			dcm.ConnectionProfile = getConnectionProfile(ctx, dct)
			payload.DiscoveryConfigModels = append(payload.DiscoveryConfigModels, dcm)
		}
	}
	return payload
}

func getConnectionProfile(ctx context.Context, plan models.OmeDiscoveryConfigTargets) (connectionProfile string) {
	connections := models.ConnectionProfiles{
		ProfileName:        "",
		ProfileDescription: "",
		Type:               "DISCOVERY",
	}
	if plan.Redfish != nil {
		protocolRedfish := models.Protocols{
			ID:       0,
			Type:     "REDFISH",
			AuthType: "Basic",
			Modified: false,
		}
		redfish := models.CredREDFISH{
			Username:  plan.Redfish.Username.ValueString(),
			Password:  plan.Redfish.Password.ValueString(),
			CaCheck:   plan.Redfish.CaCheck.ValueBool(),
			CnCheck:   plan.Redfish.CnCheck.ValueBool(),
			Port:      int(plan.Redfish.Port.ValueInt64()),
			Retries:   int(plan.Redfish.Retries.ValueInt64()),
			Timeout:   int(plan.Redfish.Timeout.ValueInt64()),
			IsHTTP:    false,
			KeepAlive: false,
		}
		protocolRedfish.Credential = redfish
		connections.Credentials = append(connections.Credentials, protocolRedfish)
	}
	if plan.WSMAN != nil {
		protocolWsman := models.Protocols{
			ID:       0,
			Type:     "WSMAN",
			AuthType: "Basic",
			Modified: false,
		}
		wsman := models.CredWSMAN{
			Username:  plan.WSMAN.Username.ValueString(),
			Password:  plan.WSMAN.Password.ValueString(),
			CaCheck:   plan.WSMAN.CaCheck.ValueBool(),
			CnCheck:   plan.WSMAN.CnCheck.ValueBool(),
			Port:      int(plan.WSMAN.Port.ValueInt64()),
			Retries:   int(plan.WSMAN.Retries.ValueInt64()),
			Timeout:   int(plan.WSMAN.Timeout.ValueInt64()),
			IsHTTP:    false,
			KeepAlive: false,
		}
		protocolWsman.Credential = wsman
		connections.Credentials = append(connections.Credentials, protocolWsman)
	}
	if plan.SNMP != nil {
		protocolSNMP := models.Protocols{
			ID:       0,
			Type:     "SNMP",
			AuthType: "Basic",
			Modified: false,
		}
		snmp := models.CredSNMP{
			Community:  plan.SNMP.Community.ValueString(),
			EnableV1V2: true,
			EnableV3:   false,
			Port:       int(plan.SNMP.Port.ValueInt64()),
			Retries:    int(plan.SNMP.Retries.ValueInt64()),
			Timeout:    int(plan.SNMP.Timeout.ValueInt64()),
		}
		protocolSNMP.Credential = snmp
		connections.Credentials = append(connections.Credentials, protocolSNMP)
	}

	if plan.SSH != nil {
		protocolSSH := models.Protocols{
			ID:       0,
			Type:     "SSH",
			AuthType: "Basic",
			Modified: false,
		}
		ssh := models.CredSSH{
			Username:        plan.SSH.Username.ValueString(),
			IsSudoUser:      plan.SSH.IsSudoUser.ValueBool(),
			Password:        plan.SSH.Password.ValueString(),
			Port:            int(plan.SSH.Port.ValueInt64()),
			UseKey:          false,
			Retries:         int(plan.SSH.Retries.ValueInt64()),
			Timeout:         int(plan.SSH.Timeout.ValueInt64()),
			CheckKnownHosts: plan.SSH.CheckKnownHosts.ValueBool(),
		}
		protocolSSH.Credential = ssh
		connections.Credentials = append(connections.Credentials, protocolSSH)
	}

	jsonData, err := json.Marshal(connections)
	if err != nil {
		tflog.Debug(ctx, "Error marshaling JSON: "+err.Error())
		return
	}
	connectionProfile = string(jsonData)
	tflog.Debug(ctx, "connection profile constructed string: "+connectionProfile)
	return
}

func discoveryState(ctx context.Context, resp models.DiscoveryJob, plan models.OmeDiscoveryJob) (state models.OmeDiscoveryJob) {
	state = models.OmeDiscoveryJob{
		DiscoveryJobID:   types.StringValue(strconv.Itoa(resp.DiscoveryConfigGroupID)),
		DiscoveryJobName: types.StringValue(resp.DiscoveryConfigGroupName),
		TrapDestination:  types.BoolValue(resp.TrapDestination),
		CommunityString:  types.BoolValue(resp.CommunityString),
	}
	if resp.DiscoveryStatusEmailRecipient != "" {
		state.EmailRecipient = types.StringValue(resp.DiscoveryStatusEmailRecipient)
	}
	if plan.Schedule.ValueString() == "RunNow" || resp.Schedule.RunNow || resp.Schedule.Cron == "startnow" {
		state.Schedule = types.StringValue("RunNow")
	} else if plan.Schedule.ValueString() == "RunLater" || resp.Schedule.RunLater || len(resp.Schedule.Cron) > 0 {
		state.Schedule = types.StringValue("RunLater")
		state.Cron = types.StringValue(resp.Schedule.Cron)
	}
	for idx, dct := range resp.DiscoveryConfigModels {
		state.DiscoveryConfigTargets = append(state.DiscoveryConfigTargets, getOmeDiscoveryConfigTargets(ctx, dct, plan.DiscoveryConfigTargets[idx]))
	}
	if len(resp.DiscoveryConfigTaskParam) == 1 {
		state.JobID = types.Int64Value(int64(resp.DiscoveryConfigTaskParam[0].TaskID))
	}
	return
}

func getOmeDiscoveryConfigTargets(ctx context.Context, resp models.DiscoveryConfigModels, plan models.OmeDiscoveryConfigTargets) models.OmeDiscoveryConfigTargets {
	connectionProfiles := models.ConnectionProfiles{}
	state := models.OmeDiscoveryConfigTargets{}
	deviceMap := map[int]string{
		1000: "SERVER",
		7000: "NETWORK SWITCH",
		2000: "CHASSIS",
		5000: "STORAGE",
	}
	for _, did := range resp.DeviceType {
		if val, ok := deviceMap[did]; ok {
			state.DeviceType = append(state.DeviceType, types.StringValue(val))
		}
	}
	for _, network := range resp.DiscoveryConfigTargets {
		state.NetworkAddressDetail = append(state.NetworkAddressDetail, types.StringValue(network.NetworkAddressDetail))
	}
	err := json.Unmarshal([]byte(resp.ConnectionProfile), &connectionProfiles)
	if err != nil {
		tflog.Debug(ctx, "Error unmarshaling JSON: "+err.Error())
		return models.OmeDiscoveryConfigTargets{}
	}
	for _, creds := range connectionProfiles.Credentials {
		if credMap, ok := creds.Credential.(map[string]interface{}); ok {
			if creds.Type == "REDFISH" && plan.Redfish != nil {
				cred := models.CredREDFISH{}
				state.Redfish = &models.OmeRedfish{}
				err := mapstructure.Decode(credMap, &cred)
				if err != nil {
					continue
				}
				state.Redfish.Username = types.StringValue(cred.Username)
				state.Redfish.Password = plan.Redfish.Password
				state.Redfish.Port = types.Int64Value(int64(cred.Port))
				state.Redfish.Retries = types.Int64Value(int64(cred.Retries))
				state.Redfish.Timeout = types.Int64Value(int64(cred.Timeout))
				state.Redfish.CaCheck = types.BoolValue(cred.CaCheck)
				state.Redfish.CnCheck = types.BoolValue(cred.CnCheck)
			} else if creds.Type == "WSMAN" && plan.WSMAN != nil {
				cred := models.CredWSMAN{}
				state.WSMAN = &models.OmeWSMAN{}
				err := mapstructure.Decode(credMap, &cred)
				if err != nil {
					continue
				}
				state.WSMAN.Username = types.StringValue(cred.Username)
				state.WSMAN.Password = plan.WSMAN.Password
				state.WSMAN.Port = types.Int64Value(int64(cred.Port))
				state.WSMAN.Retries = types.Int64Value(int64(cred.Retries))
				state.WSMAN.Timeout = types.Int64Value(int64(cred.Timeout))
				state.WSMAN.CaCheck = types.BoolValue(cred.CaCheck)
				state.WSMAN.CnCheck = types.BoolValue(cred.CnCheck)
			} else if creds.Type == "SNMP" {
				cred := models.CredSNMP{}
				state.SNMP = &models.OmeSNMP{}
				err := mapstructure.Decode(credMap, &cred)
				if err != nil {
					continue
				}
				state.SNMP.Community = types.StringValue(cred.Community)
				state.SNMP.Port = types.Int64Value(int64(cred.Port))
				state.SNMP.Retries = types.Int64Value(int64(cred.Retries))
				state.SNMP.Timeout = types.Int64Value(int64(cred.Timeout))
			} else if creds.Type == "SSH" && plan.SSH != nil {
				cred := models.CredSSH{}
				state.SSH = &models.OmeSSH{}
				err := mapstructure.Decode(credMap, &cred)
				if err != nil {
					continue
				}
				state.SSH.Username = types.StringValue(cred.Username)
				state.SSH.Password = plan.SSH.Password
				state.SSH.Port = types.Int64Value(int64(cred.Port))
				state.SSH.Retries = types.Int64Value(int64(cred.Retries))
				state.SSH.Timeout = types.Int64Value(int64(cred.Timeout))
				state.SSH.IsSudoUser = types.BoolValue(cred.IsSudoUser)
				state.SSH.CheckKnownHosts = types.BoolValue(cred.CheckKnownHosts)
			}
		}
	}
	return state
}
