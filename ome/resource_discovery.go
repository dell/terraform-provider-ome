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
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-ome/clients"
	"terraform-provider-ome/helper"
	"terraform-provider-ome/models"
	"terraform-provider-ome/utils"
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
		MarkdownDescription: "This terraform resource is used to manage Discovery entity on OME." +
			"We can Create, Update and Delete OME Discoveries using this resource. We can also do an 'Import' an existing 'Discovery' from OME .",
		Version:    1,
		Attributes: DiscoveryJobSchema(),
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
		if data.Timeout.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("timeout"),
				"Attribute Error",
				"With Schedule as RunNow, Timeout must be set.",
			)
		}
		if data.PartialFailure.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("partial_failure"),
				"Attribute Error",
				"With Schedule as RunNow, Partial Failure must be set.",
			)
		}
	}

	if data.Schedule.ValueString() == "RunLater" {
		if !data.Timeout.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("timeout"),
				"Attribute Error",
				"With Schedule as RunLater, Timeout can't be set.",
			)
		}
		if !data.PartialFailure.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("partial_failure"),
				"Attribute Error",
				"With Schedule as RunLater, Partial Failure can't be set.",
			)
		}
		if data.Cron.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("cron"),
				"Attribute Error",
				"With Schedule as RunLater, cron must be set.",
			)
		}
	}

	if len(data.DiscoveryConfigTargets) > 0 {
		for idx, dct := range data.DiscoveryConfigTargets {
			idxError := "Inappropriate value for attribute \"discovery_config_targets\": element: " + strconv.Itoa(idx) + "\n"
			if dct.Redfish == nil && dct.SNMP == nil && dct.SSH == nil && dct.WSMAN == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("discovery_config_targets").AtListIndex(idx),
					"Attribute Error",
					idxError+"Atleast one of protocol should be configured for the discovery targets.")
			}
		}
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *discoveryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_discovery create : Started")
	//Get Plan Data
	req.Plan.SetAttribute(ctx, path.Root("job_tracking"), &models.OmeJobTracking{})
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
	// if schedule is set to RunNow, we will track the job till it times out.
	err = jobTrackState(ctx, state, plan, omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateDiscovery, err.Error(),
		)
	}
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
	req.Plan.SetAttribute(ctx, path.Root("job_tracking"), &models.OmeJobTracking{})
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if !reflect.DeepEqual(state, plan) {
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
	// }
	err = jobTrackState(ctx, state, plan, omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateDiscovery, err.Error(),
		)
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
	var state models.OmeDiscoveryJob
	tflog.Trace(ctx, "resource_discovery import: started")
	omeClient, d := r.p.createOMESession(ctx, "resource_discovery Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	id, _ := strconv.Atoi(req.ID)
	respDiscovery, err := omeClient.GetDiscoveryJobByGroupID(int64(id))
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrReadDiscovery, err.Error(),
		)
		return
	}
	state = discoveryState(ctx, respDiscovery, state)
	tflog.Trace(ctx, "resource_discovery import: finished reading state")
	//Save into State
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_discovery import: finished")
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

			networAddress := utils.ConvertListValueToStringSlice(dct.NetworkAddressDetail)

			for _, networkAddress := range networAddress {
				network := models.DiscoveryConfigTargets{
					NetworkAddressDetail: networkAddress,
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
		if len(plan.DiscoveryConfigTargets) == 0 {
			state.DiscoveryConfigTargets = append(state.DiscoveryConfigTargets, getOmeDiscoveryConfigTargets(ctx, dct, models.OmeDiscoveryConfigTargets{}))
			continue
		}
		state.DiscoveryConfigTargets = append(state.DiscoveryConfigTargets, getOmeDiscoveryConfigTargets(ctx, dct, plan.DiscoveryConfigTargets[idx]))
	}
	if len(resp.DiscoveryConfigTaskParam) == 1 {
		state.JobID = types.Int64Value(int64(resp.DiscoveryConfigTaskParam[0].TaskID))
	}
	state.Timeout = plan.Timeout
	state.PartialFailure = plan.PartialFailure
	state.JobTracking = plan.JobTracking
	return
}

func jobTrackState(ctx context.Context, state models.OmeDiscoveryJob, plan models.OmeDiscoveryJob, omeClient *clients.Client) error {
	if plan.Timeout.ValueInt64() > 0 && !plan.PartialFailure.IsUnknown() {
		results, err := helper.DiscoverJobRunner(ctx, omeClient, state.JobID.ValueInt64(), plan.Timeout.ValueInt64(), plan.PartialFailure.ValueBool())
		if err != nil && !plan.PartialFailure.ValueBool() {
			return err
		}
		JobExecutionResults := make([]basetypes.StringValue, 0)
		DiscoveredIPResults := make([]basetypes.StringValue, 0)
		UnDiscoveredIPResults := make([]basetypes.StringValue, 0)
		for _, jer := range results {
			reIP := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
			ip := reIP.FindString(jer)
			reComp := regexp.MustCompile(".*Completed$")
			isDiscoverd := reComp.MatchString(jer)
			if isDiscoverd {
				DiscoveredIPResults = append(DiscoveredIPResults, types.StringValue(ip))
			} else {
				UnDiscoveredIPResults = append(UnDiscoveredIPResults, types.StringValue(ip))
			}
			JobExecutionResults = append(JobExecutionResults, types.StringValue(jer))
		}
		state.JobTracking.JobExecutionResults = JobExecutionResults
		state.JobTracking.DiscoveredIPs = DiscoveredIPResults
		state.JobTracking.UnDiscoveredIPs = UnDiscoveredIPResults
		if len(state.JobTracking.UnDiscoveredIPs) > 0 && !plan.PartialFailure.ValueBool() {
			return fmt.Errorf("discovery job completed with partial failure")
		}
	}
	return nil
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
	var networkAddress []string
	for _, network := range resp.DiscoveryConfigTargets {
		networkAddress = append(networkAddress, network.NetworkAddressDetail)
	}

	state.NetworkAddressDetail = utils.ConvertStringListValue(networkAddress)

	err := json.Unmarshal([]byte(resp.ConnectionProfile), &connectionProfiles)
	if err != nil {
		tflog.Debug(ctx, "Error unmarshaling JSON: "+err.Error())
		return models.OmeDiscoveryConfigTargets{}
	}
	for _, creds := range connectionProfiles.Credentials {
		if credMap, ok := creds.Credential.(map[string]interface{}); ok {
			tflog.Info(ctx, fmt.Sprintf("Creds %v, Type: %s", credMap, creds.Type))
			if creds.Type == "REDFISH" && plan.Redfish != nil {
				cred := &models.CredREDFISH{}
				state.Redfish = &models.OmeRedfish{}
				bytes := getCredMapByteArray(ctx, credMap)
				unmarshalErr := json.Unmarshal(bytes, cred)
				if unmarshalErr != nil {
					tflog.Error(ctx, unmarshalErr.Error())
					continue
				}
				state.Redfish.Username = types.StringValue(cred.Username)
				if plan.Redfish != nil {
					state.Redfish.Password = plan.Redfish.Password
				}

				state.Redfish.Port = types.Int64Value(int64(cred.Port))
				state.Redfish.Retries = types.Int64Value(int64(cred.Retries))
				state.Redfish.Timeout = types.Int64Value(int64(cred.Timeout))
				state.Redfish.CaCheck = types.BoolValue(cred.CaCheck)
				state.Redfish.CnCheck = types.BoolValue(cred.CnCheck)
			} else if creds.Type == "WSMAN" && plan.WSMAN != nil {
				cred := &models.CredWSMAN{}
				state.WSMAN = &models.OmeWSMAN{}
				bytes := getCredMapByteArray(ctx, credMap)
				unmarshalErr := json.Unmarshal(bytes, cred)
				if unmarshalErr != nil {
					tflog.Error(ctx, unmarshalErr.Error())
					continue
				}
				state.WSMAN.Username = types.StringValue(cred.Username)
				if plan.WSMAN != nil {
					state.WSMAN.Password = plan.WSMAN.Password
				}
				state.WSMAN.Port = types.Int64Value(int64(cred.Port))
				state.WSMAN.Retries = types.Int64Value(int64(cred.Retries))
				state.WSMAN.Timeout = types.Int64Value(int64(cred.Timeout))
				state.WSMAN.CaCheck = types.BoolValue(cred.CaCheck)
				state.WSMAN.CnCheck = types.BoolValue(cred.CnCheck)
			} else if creds.Type == "SNMP" && plan.SNMP != nil {
				cred := &models.CredSNMP{}
				state.SNMP = &models.OmeSNMP{}
				bytes := getCredMapByteArray(ctx, credMap)
				unmarshalErr := json.Unmarshal(bytes, cred)
				if unmarshalErr != nil {
					tflog.Error(ctx, unmarshalErr.Error())
					continue
				}
				state.SNMP.Community = types.StringValue(cred.Community)
				state.SNMP.Port = types.Int64Value(int64(cred.Port))
				state.SNMP.Retries = types.Int64Value(int64(cred.Retries))
				state.SNMP.Timeout = types.Int64Value(int64(cred.Timeout))
			} else if creds.Type == "SSH" && plan.SSH != nil {
				cred := &models.CredSSH{}
				state.SSH = &models.OmeSSH{}
				bytes := getCredMapByteArray(ctx, credMap)
				unmarshalErr := json.Unmarshal(bytes, cred)
				if unmarshalErr != nil {
					tflog.Error(ctx, unmarshalErr.Error())
					continue
				}
				state.SSH.Username = types.StringValue(cred.Username)
				if plan.SSH != nil {
					state.SSH.Password = plan.SSH.Password
				}
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

func getCredMapByteArray(ctx context.Context, credMap map[string]interface{}) []byte {
	bytes, errMarshal := json.Marshal(credMap)
	if errMarshal != nil {
		tflog.Error(ctx, errMarshal.Error())
	}
	return bytes
}
