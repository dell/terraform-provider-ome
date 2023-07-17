package ome

import (
	"context"
	"reflect"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		Attributes:          OmeDiscoveryJobSchema(),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *discoveryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_discovery create : Started")
	//Get Plan Data
	var plan models.OmeDiscoveryJob
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

	discoveryPayload, err := getDiscoveryPayload(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateDiscovery, err.Error(),
		)
		return
	}

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

	tflog.Trace(ctx, "resource_discovery create: updating state finished, saving ...")
	// Save into State
	state := saveState(cDiscovery)
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

	discovery, err := omeClient.GetDiscoveryJobByGroupID(state.DiscoveryConfigGroupID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrReadDiscovery, err.Error(),
		)
		return
	}

	state = saveState(discovery)
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
		for _, tid := range state.DiscoveryConfigTaskParam {
			tflog.Debug(ctx, "resource_discovery checking job status for", map[string]interface{}{
				"jobid": tid.TaskID.ValueInt64()})
			if tid.TaskID.ValueInt64() != 0 {
				jr, err := omeClient.GetJob(tid.TaskID.ValueInt64())
				if err != nil {
					resp.Diagnostics.AddError(
						clients.ErrGnrUpdateDiscovery,
						err.Error(),
					)
					return
				}
				tflog.Debug(ctx, "resource_discovery update job status is", map[string]interface{}{
					"jobid":  tid.TaskID.ValueInt64(),
					"status": jr.LastRunStatus.ID,
				})

				//if job is running during update, throw error
				if jr.LastRunStatus.ID == clients.RunningStatusID {
					resp.Diagnostics.AddError(
						clients.ErrGnrUpdateDiscovery,
						clients.ErrDiscoveryJobIsRunning,
					)
					return
				}
			}
		}

		discoveryPayload, err := getDiscoveryPayload(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateDiscovery, err.Error(),
			)
			return
		}
		tflog.Trace(ctx, "resource_discovery update Discovery")
		tflog.Debug(ctx, "resource_discovery update Discovery", map[string]interface{}{
			"Create Discovery Request": discoveryPayload,
		})

		discovery, err := omeClient.UpdateDiscoveryJob(discoveryPayload)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateDiscovery, err.Error(),
			)
			return
		}
		state = saveState(discovery)
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
			int(state.DiscoveryConfigGroupID.ValueInt64()),
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

func getDiscoveryPayload(ctx context.Context, plan *models.OmeDiscoveryJob) (models.DiscoveryJobPayload, error) {
	dcm := make([]models.DiscoveryConfigModels, 0)
	for _, dcmPlan := range plan.DiscoveryConfigModels {
		dctarget := make([]models.DiscoveryConfigTargets, 0)
		for _, dctargetPlan := range dcmPlan.DiscoveryConfigTargets {
			uniqueDCTarget := models.DiscoveryConfigTargets{
				DiscoveryConfigTargetID: int(dctargetPlan.DiscoveryConfigTargetID.ValueInt64()),
				NetworkAddressDetail:    dctargetPlan.NetworkAddressDetail.ValueString(),
				SubnetMask:              dctargetPlan.SubnetMask.ValueString(),
				AddressType:             int(dctargetPlan.AddressType.ValueInt64()),
				Disabled:                dctargetPlan.Disabled.ValueBool(),
				Exclude:                 dctargetPlan.Exclude.ValueBool(),
			}
			dctarget = append(dctarget, uniqueDCTarget)
		}
		dtype := make([]int, 0)
		for _, deviceTypePlan := range dcmPlan.DeviceType {
			dtype = append(dtype, int(deviceTypePlan.ValueInt64()))
		}
		dcvp := make([]models.DiscoveryConfigVendorPlatforms, 0)
		for _, dcvpPlan := range dcmPlan.DiscoveryConfigVendorPlatforms {
			uniqueDCVP := models.DiscoveryConfigVendorPlatforms{
				VendorPlatformId:                int(dcvpPlan.VendorPlatformId.ValueInt64()),
				DiscoveryConfigVendorPlatformId: int(dcvpPlan.DiscoveryConfigVendorPlatformId.ValueInt64()),
			}
			dcvp = append(dcvp, uniqueDCVP)
		}
		uniqueDCM := models.DiscoveryConfigModels{
			DiscoveryConfigID:              int(dcmPlan.DiscoveryConfigID.ValueInt64()),
			DiscoveryConfigDescription:     dcmPlan.DiscoveryConfigDescription.ValueString(),
			DiscoveryConfigStatus:          dcmPlan.DiscoveryConfigStatus.ValueString(),
			DiscoveryConfigTargets:         dctarget,
			ConnectionProfileID:            int(dcmPlan.ConnectionProfileID.ValueInt64()),
			ConnectionProfile:              dcmPlan.ConnectionProfile.ValueString(),
			DeviceType:                     dtype,
			DiscoveryConfigVendorPlatforms: dcvp,
		}
		dcm = append(dcm, uniqueDCM)
	}
	tflog.Debug(ctx, "resource_discovery create Creating Discovery Config Model", map[string]interface{}{
		"Create Discovery Config Model Request": dcm,
	})
	dctp := make([]models.DiscoveryConfigTaskParam, 0)
	for _, dctpPlan := range plan.DiscoveryConfigTaskParam {
		uniqueDCTP := models.DiscoveryConfigTaskParam{
			TaskID:            int(dctpPlan.TaskID.ValueInt64()),
			TaskTypeID:        int(dctpPlan.TaskTypeID.ValueInt64()),
			ExecutionSequence: int(dctpPlan.ExecutionSequence.ValueInt64()),
		}
		dctp = append(dctp, uniqueDCTP)
	}
	dctask := make([]models.DiscoveryConfigTasks, 0)
	for _, dctPlan := range plan.DiscoveryConfigTasks {
		uniqueDCT := models.DiscoveryConfigTasks{
			DiscoveryConfigDescription:           dctPlan.DiscoveryConfigDescription.ValueString(),
			DiscoveryConfigEmailRecipient:        dctPlan.DiscoveryConfigEmailRecipient.ValueString(),
			DiscoveryConfigDiscoveredDeviceCount: dctPlan.DiscoveryConfigDiscoveredDeviceCount.ValueString(),
			DiscoveryConfigRequestId:             int(dctPlan.DiscoveryConfigRequestId.ValueInt64()),
			DiscoveryConfigExpectedDeviceCount:   dctPlan.DiscoveryConfigExpectedDeviceCount.ValueString(),
			DiscoveryConfigName:                  dctPlan.DiscoveryConfigName.ValueString(),
		}
		dctask = append(dctask, uniqueDCT)
	}
	sj := getScheduleJob(plan)
	payload := models.DiscoveryJobPayload{
		ChassisIdentifier:               plan.ChassisIdentifier.ValueString(),
		CommunityString:                 plan.CommunityString.ValueBool(),
		CreateGroup:                     plan.CreateGroup.ValueBool(),
		DiscoveryConfigGroupDescription: plan.DiscoveryConfigGroupDescription.ValueString(),
		DiscoveryConfigGroupID:          int(plan.DiscoveryConfigGroupID.ValueInt64()),
		DiscoveryConfigGroupName:        plan.DiscoveryConfigGroupName.ValueString(),
		DiscoveryConfigModels:           dcm,
		DiscoveryConfigParentGroupID:    int(plan.DiscoveryConfigParentGroupID.ValueInt64()),
		DiscoveryConfigTaskParam:        dctp,
		DiscoveryConfigTasks:            dctask,
		DiscoveryStatusEmailRecipient:   plan.DiscoveryStatusEmailRecipient.ValueString(),
		Schedule:                        sj,
		TrapDestination:                 plan.TrapDestination.ValueBool(),
		UseAllProfiles:                  plan.UseAllProfiles.ValueBool(),
	}
	return payload, nil
}

func getScheduleJob(plan *models.OmeDiscoveryJob) models.ScheduleJob {
	recurring := models.Recurring{}
	if reflect.ValueOf(plan.Schedule.Recurring).IsValid() {
		hourly := models.Hourly{}
		if !plan.Schedule.Recurring.Hourly.Frequency.IsUnknown() {
			hourly.Frequency = int(plan.Schedule.Recurring.Hourly.Frequency.ValueInt64())
		}
		daily := models.Daily{}
		if !plan.Schedule.Recurring.Daily.Frequency.IsUnknown() && !plan.Schedule.Recurring.Daily.Time.Minutes.IsUnknown() && !plan.Schedule.Recurring.Daily.Time.Hour.IsUnknown() {
			daily.Time.Minutes = int(plan.Schedule.Recurring.Daily.Time.Minutes.ValueInt64())
			daily.Time.Hour = int(plan.Schedule.Recurring.Daily.Time.Hour.ValueInt64())
			daily.Frequency = int(plan.Schedule.Recurring.Daily.Frequency.ValueInt64())
		}
		weekley := models.Weekley{}
		if !plan.Schedule.Recurring.Weekley.Time.Minutes.IsUnknown() && !plan.Schedule.Recurring.Weekley.Time.Hour.IsUnknown() && !plan.Schedule.Recurring.Weekley.Day.IsUnknown() {
			weekley.Time.Minutes = int(plan.Schedule.Recurring.Weekley.Time.Minutes.ValueInt64())
			weekley.Time.Hour = int(plan.Schedule.Recurring.Weekley.Time.Hour.ValueInt64())
			weekley.Day = plan.Schedule.Recurring.Weekley.Day.ValueString()
		}

		recurring.Hourly = hourly
		recurring.Daily = daily
		recurring.Weekley = weekley

	}
	sj := models.ScheduleJob{
		RunNow:    plan.Schedule.RunNow.ValueBool(),
		RunLater:  plan.Schedule.RunLater.ValueBool(),
		Recurring: recurring,
		Cron:      plan.Schedule.Cron.ValueString(),
		StartTime: plan.Schedule.StartTime.ValueString(),
		EndTime:   plan.Schedule.EndTime.ValueString(),
	}
	return sj
}

func saveState(resp models.DiscoveryJob) (state models.OmeDiscoveryJob) {
	state.ChassisIdentifier = types.StringValue(resp.ChassisIdentifier)
	state.CommunityString = types.BoolValue(resp.CommunityString)
	state.CreateGroup = types.BoolValue(resp.CreateGroup)
	state.DiscoveryConfigGroupDescription = types.StringValue(resp.DiscoveryConfigGroupDescription)
	state.DiscoveryConfigGroupID = types.Int64Value(int64(resp.DiscoveryConfigGroupID))
	state.DiscoveryConfigGroupName = types.StringValue(resp.DiscoveryConfigGroupName)
	state.DiscoveryConfigParentGroupID = types.Int64Value(int64(resp.DiscoveryConfigParentGroupID))
	state.DiscoveryStatusEmailRecipient = types.StringValue(resp.DiscoveryStatusEmailRecipient)
	state.TrapDestination = types.BoolValue(resp.TrapDestination)
	state.UseAllProfiles = types.BoolValue(resp.UseAllProfiles)
	state.DiscoveryConfigModels = []models.OmeDiscoveryConfigModels{}
	for _, dcm := range resp.DiscoveryConfigModels {
		lodct := []models.OmeDiscoveryConfigTargets{}
		for _, dodct := range dcm.DiscoveryConfigTargets {
			odct := models.OmeDiscoveryConfigTargets{
				DiscoveryConfigTargetID: types.Int64Value(int64(dodct.DiscoveryConfigTargetID)),
				NetworkAddressDetail:    types.StringValue(dodct.NetworkAddressDetail),
				SubnetMask:              types.StringValue(dodct.SubnetMask),
				AddressType:             types.Int64Value(int64(dodct.AddressType)),
				Disabled:                types.BoolValue(dodct.Disabled),
				Exclude:                 types.BoolValue(dodct.Exclude),
			}
			lodct = append(lodct, odct)
		}
		ldt := make([]types.Int64, 0)
		for _, ddt := range dcm.DeviceType {
			ldt = append(ldt, types.Int64Value(int64(ddt)))
		}
		lodcvp := []models.OmeDiscoveryConfigVendorPlatforms{}
		for _, ddcvp := range dcm.DiscoveryConfigVendorPlatforms {
			odcvp := models.OmeDiscoveryConfigVendorPlatforms{
				VendorPlatformId:                types.Int64Value(int64(ddcvp.VendorPlatformId)),
				DiscoveryConfigVendorPlatformId: types.Int64Value(int64(ddcvp.DiscoveryConfigVendorPlatformId)),
			}
			lodcvp = append(lodcvp, odcvp)
		}
		odcm := models.OmeDiscoveryConfigModels{
			DiscoveryConfigID:              types.Int64Value(int64(dcm.DiscoveryConfigID)),
			DiscoveryConfigDescription:     types.StringValue(dcm.DiscoveryConfigDescription),
			DiscoveryConfigStatus:          types.StringValue(dcm.DiscoveryConfigStatus),
			DiscoveryConfigTargets:         lodct,
			ConnectionProfileID:            types.Int64Value(int64(dcm.ConnectionProfileID)),
			ConnectionProfile:              types.StringValue(dcm.ConnectionProfile),
			DeviceType:                     ldt,
			DiscoveryConfigVendorPlatforms: lodcvp,
		}
		state.DiscoveryConfigModels = append(state.DiscoveryConfigModels, odcm)
	}
	state.DiscoveryConfigTaskParam = []models.OmeDiscoveryConfigTaskParam{}
	for _, ddctp := range resp.DiscoveryConfigTaskParam {
		odctp := models.OmeDiscoveryConfigTaskParam{
			TaskID:            types.Int64Value(int64(ddctp.TaskID)),
			TaskTypeID:        types.Int64Value(int64(ddctp.TaskTypeID)),
			ExecutionSequence: types.Int64Value(int64(ddctp.ExecutionSequence)),
		}
		state.DiscoveryConfigTaskParam = append(state.DiscoveryConfigTaskParam, odctp)
	}
	state.DiscoveryConfigTasks = []models.OmeDiscoveryConfigTasks{}
	for _, ddct := range resp.DiscoveryConfigTasks {
		odct := models.OmeDiscoveryConfigTasks{
			DiscoveryConfigDescription:           types.StringValue(ddct.DiscoveryConfigDescription),
			DiscoveryConfigEmailRecipient:        types.StringValue(ddct.DiscoveryConfigEmailRecipient),
			DiscoveryConfigDiscoveredDeviceCount: types.StringValue(ddct.DiscoveryConfigDiscoveredDeviceCount),
			DiscoveryConfigRequestId:             types.Int64Value(int64(ddct.DiscoveryConfigRequestId)),
			DiscoveryConfigExpectedDeviceCount:   types.StringValue(ddct.DiscoveryConfigExpectedDeviceCount),
			DiscoveryConfigName:                  types.StringValue(ddct.DiscoveryConfigName),
		}
		state.DiscoveryConfigTasks = append(state.DiscoveryConfigTasks, odct)
	}
	state.Schedule = models.OmeScheduleJob{
		RunNow:   types.BoolValue(resp.Schedule.RunNow),
		RunLater: types.BoolValue(resp.Schedule.RunLater),
		Recurring: models.OmeRecurring{
			Hourly: models.OmeHourly{
				Frequency: types.Int64Value(int64(resp.Schedule.Recurring.Hourly.Frequency)),
			},
			Daily: models.OmeDaily{
				Frequency: types.Int64Value(int64(resp.Schedule.Recurring.Daily.Frequency)),
				Time: models.OmeTime{
					Minutes: types.Int64Value(int64(resp.Schedule.Recurring.Daily.Time.Minutes)),
					Hour:    types.Int64Value(int64(resp.Schedule.Recurring.Daily.Time.Hour)),
				},
			},
			Weekley: models.OmeWeekley{
				Day: types.StringValue(resp.Schedule.Recurring.Weekley.Day),
				Time: models.OmeTime{
					Minutes: types.Int64Value(int64(resp.Schedule.Recurring.Weekley.Time.Minutes)),
					Hour:    types.Int64Value(int64(resp.Schedule.Recurring.Weekley.Time.Hour)),
				},
			},
		},
		Cron:      types.StringValue(resp.Schedule.Cron),
		StartTime: types.StringValue(resp.Schedule.StartTime),
		EndTime:   types.StringValue(resp.Schedule.EndTime),
	}
	return
}
