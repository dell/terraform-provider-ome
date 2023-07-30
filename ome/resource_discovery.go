package ome

import (
	"context"
	"encoding/json"
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

	discoveryPayload := getDiscoveryPayload(ctx, &plan)

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
	state = discoveryState(ctx, cDiscovery, state)
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

	did, _ := strconv.Atoi(state.DiscoveryJobID.ValueString())
	respDiscovery, err := omeClient.GetDiscoveryJobByGroupID(int64(did))
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
	did, _ := strconv.Atoi(state.DiscoveryJobID.ValueString())
	ddj := models.DiscoveryJobDeletePayload{
		DiscoveryGroupIds: []int{
			did,
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getDiscoveryPayload(ctx context.Context, plan *models.OmeDiscoveryJob) (payload models.DiscoveryJobPayload) {
	return payload
}

func discoveryState(ctx context.Context, resp models.DiscoveryJob, plan models.OmeDiscoveryJob) (state models.OmeDiscoveryJob) {
	state = models.OmeDiscoveryJob{
		DiscoveryJobID:   types.StringValue(strconv.Itoa(resp.DiscoveryConfigGroupID)),
		DiscoveryJobName: types.StringValue(resp.DiscoveryConfigGroupName),
		EmailRecipient:   types.StringValue(resp.DiscoveryStatusEmailRecipient),
		TrapDestination:  types.BoolValue(resp.TrapDestination),
		CommunityString:  types.BoolValue(resp.CommunityString),
	}
	if plan.Schedule.ValueString() == "RunNow" || resp.Schedule.RunNow || resp.Schedule.Cron == "startnow" {
		state.Schedule = types.StringValue("RunNow")
		if !plan.JobWait.IsNull() {
			state.JobWait = types.BoolValue(plan.JobWait.ValueBool())
		} else {
			state.JobWait = types.BoolValue(true)
		}
		if !plan.JobWaitTimeout.IsNull() {
			state.JobWaitTimeout = types.Int64Value(plan.JobWaitTimeout.ValueInt64())
		} else {
			state.JobWaitTimeout = types.Int64Value(1200)
		}
		if !plan.IgnorePartialFailure.IsNull() {
			state.IgnorePartialFailure = types.BoolValue(plan.IgnorePartialFailure.ValueBool())
		} else {
			state.IgnorePartialFailure = types.BoolValue(false)
		}
	} else if plan.Schedule.ValueString() == "RunLater" || resp.Schedule.RunLater || len(resp.Schedule.Cron) > 0 {
		state.Schedule = types.StringValue("RunLater")
		state.Cron = types.StringValue(resp.Schedule.Cron)
	}
	for _, dct := range resp.DiscoveryConfigModels {
		state.DiscoveryConfigTargets = append(state.DiscoveryConfigTargets, getOmeDiscoveryConfigTargets(ctx, dct))
	}
	return
}

func getOmeDiscoveryConfigTargets(ctx context.Context, resp models.DiscoveryConfigModels) (state models.OmeDiscoveryConfigTargets) {
	var connectionProfiles models.ConnectionProfiles
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
		return
	}
	for _, creds := range connectionProfiles.Credentials {
		if credMap, ok := creds.Credential.(map[string]interface{}); ok {
			if creds.Type == "REDFISH" {
				var cred models.CredREDFISH
				err := mapstructure.Decode(credMap,&cred)
				if err != nil{
					continue
				}
				state.Redfish.Username = types.StringValue(cred.Username)
				state.Redfish.Password = types.StringValue(cred.Password)
				state.Redfish.Port = types.Int64Value(int64(cred.Port))
				state.Redfish.Retries = types.Int64Value(int64(cred.Retries))
				state.Redfish.Timeout = types.Int64Value(int64(cred.Timeout))
				state.Redfish.CaCheck = types.BoolValue(cred.CaCheck)
				state.Redfish.CnCheck = types.BoolValue(cred.CnCheck)
			} else if creds.Type == "SNMP" {
				var cred models.CredSNMP
				err := mapstructure.Decode(credMap,&cred)
				if err != nil{
					continue
				}
				state.SNMP.Community = types.StringValue(cred.Community)
				state.SNMP.Port = types.Int64Value(int64(cred.Port))
				state.SNMP.Retries = types.Int64Value(int64(cred.Retries))
				state.SNMP.Timeout = types.Int64Value(int64(cred.Timeout))
			} else if creds.Type == "SSH" {
				var cred models.CredSSH
				err := mapstructure.Decode(credMap,&cred)
				if err != nil{
					continue
				}
				state.SSH.Username = types.StringValue(cred.Username)
				state.SSH.Password = types.StringValue(cred.Password)
				state.SSH.Port = types.Int64Value(int64(cred.Port))
				state.SSH.Retries = types.Int64Value(int64(cred.Retries))
				state.SSH.Timeout = types.Int64Value(int64(cred.Timeout))
				state.SSH.IsSudoUser = types.BoolValue(cred.IsSudoUser)
				state.SSH.CheckKnownHosts = types.BoolValue(cred.CheckKnownHosts)
			}
		}
	}
	return
}
