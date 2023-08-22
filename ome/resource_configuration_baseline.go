/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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
	"net/mail"
	"strings"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	//NoOFTries to get the task id
	NoOFTries = 5
	//NotifyNonCompliance to notify on non compliance of device
	NotifyNonCompliance = "NOTIFY_ON_NON_COMPLIANCE" // #nosec G101
	//NotifyOnSchedule to notify on schedule
	NotifyOnSchedule = "NOTIFY_ON_SCHEDULE" // #nosec G101
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceConfigurationBaseline{}
	_ resource.ResourceWithConfigure   = &resourceConfigurationBaseline{}
	_ resource.ResourceWithImportState = &resourceConfigurationBaseline{}
)

// NewConfigurationBaselineResource is new resource for configuration baseline
func NewConfigurationBaselineResource() resource.Resource {
	return &resourceConfigurationBaseline{}
}

type resourceConfigurationBaseline struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceConfigurationBaseline) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (r resourceConfigurationBaseline) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "configuration_baseline"
}

// Template Deployment Resource schema
func (r resourceConfigurationBaseline) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing configuration baselines on OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the configuration baseline resource.",
				Description:         "ID of the configuration baseline resource.",
				Computed:            true,
			},
			"ref_template_id": schema.Int64Attribute{
				MarkdownDescription: "Reference template ID." +
					" Conflicts with `ref_template_name`.",
				Description: "Reference template ID." +
					" Conflicts with 'ref_template_name'.",
				Computed: true,
				Optional: true,
			},
			"ref_template_name": schema.StringAttribute{
				MarkdownDescription: "Reference template name." +
					" Conflicts with `ref_template_id`.",
				Description: "Reference template name." +
					" Conflicts with 'ref_template_id'.",
				Optional: true,
				Computed: true,
			},
			"baseline_name": schema.StringAttribute{
				MarkdownDescription: "Name of the Baseline.",
				Description:         "Name of the Baseline.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the baseline.",
				Description:         "Description of the baseline.",
				Optional:            true,
				Computed:            true,
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "List of the device id on which the baseline compliance needs to be run." +
					" Conflicts with `device_servicetags`.",
				Description: "List of the device id on which the baseline compliance needs to be run." +
					" Conflicts with 'device_servicetags'.",
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"device_servicetags": schema.SetAttribute{
				MarkdownDescription: "List of the device servicetag on which the baseline compliance needs to be run." +
					" Conflicts with `device_ids`.",
				Description: "List of the device servicetag on which the baseline compliance needs to be run." +
					" Conflicts with 'device_ids'.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"schedule": schema.BoolAttribute{
				MarkdownDescription: "Schedule notification via email." +
					" Default value is `false`.",
				Description: "Schedule notification via email." +
					" Default value is 'false'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					BoolDefaultValue(types.BoolValue(false)),
				},
			},
			"notify_on_schedule": schema.BoolAttribute{
				MarkdownDescription: "Schedule notification via cron or any time the baseline becomes non-compliant." +
					" Default value is `false`.",
				Description: "Schedule notification via cron or any time the baseline becomes non-compliant." +
					" Default value is 'false'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					BoolDefaultValue(types.BoolValue(false)),
				},
			},
			"email_addresses": schema.SetAttribute{
				MarkdownDescription: "Email addresses for notification." +
					" Can be set only when `schedule` is `true`.",
				Description: "Email addresses for notification." +
					" Can be set only when 'schedule' is 'true'.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"output_format": schema.StringAttribute{
				MarkdownDescription: "Output format type, the input is case senitive." +
					" Valid values are `html`, `csv`, `pdf`and `xls`. Default value is `html`.",
				Description: "Output format type, the input is case senitive." +
					" Valid values are 'html', 'csv', 'pdf'and 'xls'. Default value is 'html'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					StringDefaultValue(types.StringValue("html")),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(strings.Split(clients.ValidOutputFormat, ",")...),
				},
			},
			"cron": schema.StringAttribute{
				MarkdownDescription: "Cron expression for notification schedule." +
					" Can be set only when both `schedule` and `notify_on_schedule` are set to `true`.",
				Description: "Cron expression for notification schedule." +
					" Can be set only when both 'schedule' and 'notify_on_schedule' are set to 'true'.",
				Optional: true,
			},
			"task_id": schema.Int64Attribute{
				MarkdownDescription: "Task id associated with baseline.",
				Description:         "Task id associated with baseline.",
				Computed:            true,
			},
			"job_retry_count": schema.Int64Attribute{
				MarkdownDescription: "Number of times the job has to be polled to get the final status of the resource." +
					" Default value is `30`.",
				Description: "Number of times the job has to be polled to get the final status of the resource." +
					" Default value is '30'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(30)),
				},
			},
			"sleep_interval": schema.Int64Attribute{
				MarkdownDescription: "Sleep time interval for job polling in seconds." +
					" Default value is `20`.",
				Description: "Sleep time interval for job polling in seconds." +
					" Default value is '20'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(20)),
				},
			},
		},
	}
}

// Create a new resource
func (r resourceConfigurationBaseline) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_configuration_baseline create: started")
	var plan models.ConfigureBaselines
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := models.ConfigureBaselines{}

	err := validateNotification(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}
	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_baseline Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_configuration_baseline create Validating Template Details")
	omeTemplate, err := validateRefTemplateDetails(plan.RefTemplateID.ValueInt64(), plan.RefTemplateName.ValueString(), omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateBaseline,
			err.Error(),
		)
		return
	}

	var serviceTags []string
	var devIDs []int64

	diags = plan.DeviceServicetags.ElementsAs(ctx, &serviceTags, true)
	resp.Diagnostics.Append(diags...)

	diags = plan.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "resource_configuration_baseline create Validating device details")
	targetDevices, usedDeviceInput, err := getValidTargetDevices(omeClient, serviceTags, devIDs)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateBaseline,
			err.Error(),
		)
		return
	}

	cb, err := getPayload(ctx, &plan, omeTemplate.ID, targetDevices)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateBaseline, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_baseline create Creating Baseline")
	tflog.Debug(ctx, "resource_configuration_baseline create Creating Baseline", map[string]interface{}{
		"Create Baseline Request": cb,
	})

	cBaseline, err := omeClient.CreateBaseline(cb)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateBaseline, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_baseline : create Finished creating Baseline")

	tflog.Trace(ctx, "resource_configuration_baseline : create Fetching task id for a baseline")

	baseline, err := getLatestBaseline(omeClient, cBaseline.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateBaseline, err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "resource_configuration_baseline : create Baseline created ", map[string]interface{}{
		"baselineid": baseline.ID,
		"taskid":     baseline.TaskID,
	})

	isSuccess, message := omeClient.TrackJob(baseline.TaskID, plan.JobRetryCount.ValueInt64(), plan.SleepInterval.ValueInt64())
	if !isSuccess {
		resp.Diagnostics.AddWarning(
			clients.ErrBaselineCreationTask, message,
		)
	}

	tflog.Trace(ctx, "resource_configuration_baseline Done Track Job", map[string]interface{}{
		"TaskID": fmt.Sprint(baseline.TaskID),
	})

	// Save into State
	updateBaselineState(ctx, &state, &plan, baseline, usedDeviceInput, omeClient)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read resource information
func (r resourceConfigurationBaseline) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get State Data
	var state models.ConfigureBaselines
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and differ the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_baseline Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	var usedDeviceInput string

	if len(state.DeviceIDs.Elements()) > 0 {
		usedDeviceInput = clients.DeviceIDs
	} else if len(state.DeviceServicetags.Elements()) > 0 {
		usedDeviceInput = clients.ServiceTags
	}
	baseline, err := omeClient.GetBaselineByID(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrReadBaseline, err.Error(),
		)
		return
	}

	updateBaselineState(ctx, &state, &state, baseline, usedDeviceInput, omeClient)
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceConfigurationBaseline) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get state Data
	var state models.ConfigureBaselines
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan Data
	var plan models.ConfigureBaselines
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateNotification(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	//Create Session and differ the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_baseline Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Trace(ctx, "resource_configuration_baseline update checking the job status")
	tflog.Debug(ctx, "resource_configuration_baseline checking job status for", map[string]interface{}{
		"jobid": state.TaskID.ValueInt64(),
	})

	if state.TaskID.ValueInt64() != 0 {
		jr, err := omeClient.GetJob(state.TaskID.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateBaseline,
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, "resource_configuration_baseline update job status is", map[string]interface{}{
			"jobid":  state.TaskID.ValueInt64(),
			"status": jr.LastRunStatus.ID,
		})

		//if job is running during update, throw error
		if jr.LastRunStatus.ID == clients.RunningStatusID {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateBaseline,
				clients.ErrBaseLineJobIsRunning,
			)
			return
		}
	}

	tflog.Info(ctx, "resource_configuration_baseline update Validating Template Details")
	omeTemplate, err := validateRefTemplateDetails(plan.RefTemplateID.ValueInt64(), plan.RefTemplateName.ValueString(), omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrUpdateBaseline,
			err.Error(),
		)
		return
	}

	var serviceTags []string
	var devIDs []int64

	diags = plan.DeviceServicetags.ElementsAs(ctx, &serviceTags, true)
	resp.Diagnostics.Append(diags...)

	diags = plan.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "resource_configuration_baseline update Validating device details")

	targetDevices, usedDeviceInput, err := getValidTargetDevices(omeClient, serviceTags, devIDs)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrUpdateBaseline,
			err.Error(),
		)
		return
	}

	cb, err := getPayload(ctx, &plan, omeTemplate.ID, targetDevices)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrUpdateBaseline, err.Error(),
		)
		return
	}
	cb.ID = state.ID.ValueInt64() // For update case

	tflog.Trace(ctx, "resource_configuration_baseline update Creating Baseline")
	tflog.Debug(ctx, "resource_configuration_baseline update Creating Baseline", map[string]interface{}{
		"Create Baseline Request": cb,
	})

	uBaseline, err := omeClient.UpdateBaseline(cb)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrUpdateBaseline, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_baseline : update Finished creating Baseline")

	tflog.Trace(ctx, "resource_configuration_baseline : update Fetching task id for a baseline")

	baseline, err := getLatestBaseline(omeClient, uBaseline.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrUpdateBaseline, err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "resource_configuration_baseline : update Baseline created ", map[string]interface{}{
		"baselineid": baseline.ID,
		"taskid":     baseline.TaskID,
	})

	isSuccess, message := omeClient.TrackJob(baseline.TaskID, plan.JobRetryCount.ValueInt64(), plan.SleepInterval.ValueInt64())
	if !isSuccess {
		resp.Diagnostics.AddWarning(
			clients.ErrBaselineCreationTask, message,
		)
	}

	tflog.Trace(ctx, "resource_configuration_baseline update Done Track Job", map[string]interface{}{
		"TaskID": fmt.Sprint(baseline.TaskID),
	})

	// Save into State
	updateBaselineState(ctx, &state, &plan, baseline, usedDeviceInput, omeClient)

	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r resourceConfigurationBaseline) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get State Data
	var state models.ConfigureBaselines
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and differ the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_baseline Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	err := omeClient.DeleteBaseline([]int64{state.ID.ValueInt64()})

	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrDeleteBaseline,
			err.Error(),
		)
	}
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceConfigurationBaseline) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Save the import identifier in the id attribute
	var state models.ConfigureBaselines
	baselineName := req.ID

	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_baseline ImportState")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	baseline, err := omeClient.GetBaselineByName(baselineName)
	if err != nil {
		resp.Diagnostics.AddError(clients.ErrImportDeployment, err.Error())
		return
	}
	dia := updateBaselineState(ctx, &state, &state, baseline, clients.ServiceTags, omeClient)
	resp.Diagnostics.Append(dia...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Save into State
	if len(state.EmailAddresses.Elements()) == 0 {
		state.EmailAddresses, _ = types.SetValue(types.StringType, nil)
	}
	if len(state.DeviceIDs.Elements()) == 0 {
		state.DeviceIDs, _ = types.SetValue(types.Int64Type, nil)
	}
	state.OutputFormat = types.StringValue("html")
	state.JobRetryCount = types.Int64Value(30)
	state.SleepInterval = types.Int64Value(20)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func validateRefTemplateDetails(refTemplateID int64, refTemplateName string, omeClient *clients.Client) (models.OMETemplate, error) {
	if refTemplateID == 0 && refTemplateName == "" {
		return models.OMETemplate{}, fmt.Errorf(clients.ErrInvalidRefTemplateNameorID)
	}
	if refTemplateID > 0 && refTemplateName != "" {
		return models.OMETemplate{}, fmt.Errorf(clients.ErrInvalidRefTemplateNameorID)
	}
	omeTemplate, err := omeClient.GetTemplateByIDOrName(refTemplateID, refTemplateName)
	if err != nil {
		return models.OMETemplate{}, err
	}
	if omeTemplate.ViewTypeID != ComplianceViewTypeID {
		return models.OMETemplate{}, fmt.Errorf(clients.ErrInvalidRefTemplateType)
	}
	return omeTemplate, nil
}

func validateDevicesCapablity(deviceIds []int64, deviceServiceTags []string, omeClient *clients.Client) ([]models.Device, error) {
	var invalidDevices []models.Device
	devices, err := omeClient.GetDevices(deviceServiceTags, deviceIds, []string{})
	if err != nil {
		return []models.Device{}, err
	}
	for _, device := range devices {
		validDevice := false
		for _, deviceCapability := range device.DeviceCapabilities {
			if deviceCapability == 33 {
				validDevice = true
				break
			}
		}
		if !validDevice {
			invalidDevices = append(invalidDevices, device)
		}
	}
	if len(invalidDevices) != 0 {
		return []models.Device{}, fmt.Errorf(clients.ErrDeviceNotCapable, invalidDevices)
	}
	return devices, nil
}

func updateBaselineState(ctx context.Context, state *models.ConfigureBaselines, plan *models.ConfigureBaselines, omeBaseline models.OmeBaseline, usedDeviceInput string, omeClient *clients.Client) (dia diag.Diagnostics) {
	state.ID = types.Int64Value(omeBaseline.ID)
	state.RefTemplateID = types.Int64Value(omeBaseline.TemplateID)
	state.RefTemplateName = types.StringValue(omeBaseline.TemplateName)
	state.Description = types.StringValue(omeBaseline.Description)
	state.BaselineName = types.StringValue(omeBaseline.Name)

	if usedDeviceInput == clients.ServiceTags {
		apiDeviceIDs := map[string]models.Device{}
		devSts := []string{}
		deviceStVals := []attr.Value{}
		for _, bTarget := range omeBaseline.BaselineTargets {
			device, _ := omeClient.GetDevice("", bTarget.ID)
			apiDeviceIDs[device.DeviceServiceTag] = device
		}

		if len(plan.DeviceServicetags.Elements()) > 0 {
			plan.DeviceServicetags.ElementsAs(ctx, &devSts, true)
		}

		for _, devSt := range devSts {
			if val, ok := apiDeviceIDs[devSt]; ok {
				deviceStVals = append(deviceStVals, types.StringValue(val.DeviceServiceTag))
				delete(apiDeviceIDs, devSt)
			}
		}

		if len(apiDeviceIDs) != 0 {
			for _, val := range apiDeviceIDs {
				deviceStVals = append(deviceStVals, types.StringValue(val.DeviceServiceTag))
			}
		}

		devSTsTfsdk, _ := types.SetValue(
			types.StringType,
			deviceStVals,
		)
		state.DeviceServicetags = devSTsTfsdk
		state.DeviceIDs = plan.DeviceIDs
	} else {
		apiDeviceIDs := map[int64]models.Device{}
		devIDs := []int64{}
		deviceIDVals := []attr.Value{}
		for _, bTarget := range omeBaseline.BaselineTargets {
			device, _ := omeClient.GetDevice("", bTarget.ID)
			apiDeviceIDs[device.ID] = device
		}

		plan.DeviceIDs.ElementsAs(ctx, &devIDs, true)

		for _, devID := range devIDs {
			if val, ok := apiDeviceIDs[devID]; ok {
				deviceIDVals = append(deviceIDVals, types.Int64Value(val.ID))
				delete(apiDeviceIDs, devID)
			}
		}

		if len(apiDeviceIDs) != 0 {
			for _, val := range apiDeviceIDs {
				deviceIDVals = append(deviceIDVals, types.Int64Value(val.ID))
			}
		}

		devIDsTfsdk, _ := types.SetValue(types.Int64Type, deviceIDVals)
		state.DeviceIDs = devIDsTfsdk
		state.DeviceServicetags = plan.DeviceServicetags
	}

	notificationSettings := omeBaseline.NotificationSettings
	if notificationSettings != nil {
		state.Schedule = types.BoolValue(true)
		state.OutputFormat = types.StringValue(notificationSettings.OutputFormat)

		emailAddress := []attr.Value{}
		for _, v := range notificationSettings.EmailAddresses {
			emailAddress = append(emailAddress, types.StringValue(v))
		}

		emailTfsdk, emailerror := types.SetValue(types.StringType, emailAddress)
		dia.Append(emailerror...)
		state.EmailAddresses = emailTfsdk

		if notificationSettings.NotificationType == NotifyNonCompliance {
			state.NotifyOnSchedule = types.BoolValue(false)
		} else {
			state.NotifyOnSchedule = types.BoolValue(true)
		}

		if notificationSettings.Schedule.Cron != "" {
			state.Cron = types.StringValue(notificationSettings.Schedule.Cron)
		} else {
			state.Cron = plan.Cron
		}
	} else {
		state.Schedule = plan.Schedule
		state.NotifyOnSchedule = plan.NotifyOnSchedule
		state.EmailAddresses = plan.EmailAddresses
		state.OutputFormat = plan.OutputFormat
		state.Cron = plan.Cron

	}
	state.TaskID = types.Int64Value(omeBaseline.TaskID)
	state.JobRetryCount = plan.JobRetryCount
	state.SleepInterval = plan.SleepInterval
	return
}

func getPayload(ctx context.Context, plan *models.ConfigureBaselines, templateID int64, targetDevices []models.Device) (models.ConfigurationBaselinePayload, error) {
	cbp := models.ConfigurationBaselinePayload{
		Name:        plan.BaselineName.ValueString(),
		Description: plan.Description.ValueString(),
		TemplateID:  templateID,
	}

	var baselineTargets []models.BaselineTarget
	for _, tDevices := range targetDevices {
		baselineTarget := models.BaselineTarget{
			ID: tDevices.ID,
			Type: models.BaselineTargetType{
				ID:   1,
				Name: "DEVICE",
			},
		}
		baselineTargets = append(baselineTargets, baselineTarget)
	}
	cbp.BaselineTargets = baselineTargets

	if plan.Schedule.ValueBool() {
		if len(plan.EmailAddresses.Elements()) == 0 {
			return models.ConfigurationBaselinePayload{}, fmt.Errorf(clients.ErrScheduleNotification)
		}
		var emailaddressList []string
		plan.EmailAddresses.ElementsAs(ctx, &emailaddressList, true)
		for _, emailAddress := range emailaddressList {
			_, err := mail.ParseAddress(emailAddress)
			if err != nil {
				return models.ConfigurationBaselinePayload{}, fmt.Errorf(clients.ErrInvalidEmailAddress, emailAddress)
			}
		}

		if plan.NotifyOnSchedule.ValueBool() && plan.Cron.ValueString() == "" {
			return models.ConfigurationBaselinePayload{}, fmt.Errorf(clients.ErrInvalidCronExpression)
		}

		notificationType := NotifyNonCompliance
		if plan.NotifyOnSchedule.ValueBool() {
			notificationType = NotifyOnSchedule
		}

		notificationSettings := models.NotificationSettings{
			NotificationType: notificationType,
			EmailAddresses:   emailaddressList,
			Schedule: models.BaselineNotificationSchedule{
				Cron: plan.Cron.ValueString(),
			},
			OutputFormat: strings.ToUpper(plan.OutputFormat.ValueString()),
		}
		cbp.NotificationSettings = &notificationSettings
	}

	return cbp, nil
}

func getLatestBaseline(omeClient *clients.Client, baselineID int64) (models.OmeBaseline, error) {
	retries := 1
	var taskID int64 = 0
	var baseline models.OmeBaseline
	var err error
	for taskID == 0 && retries != NoOFTries {
		time.Sleep(3 * time.Second)
		retries = retries + 1
		baseline, err = omeClient.GetBaselineByID(baselineID)
		if err != nil {
			return models.OmeBaseline{}, err
		}
		taskID = baseline.TaskID
	}
	return baseline, nil
}

func getValidTargetDevices(omeClient *clients.Client, serviceTags []string, devIDs []int64) ([]models.Device, string, error) {
	usedDeviceInput, err := clients.DeviceMutuallyExclusive(serviceTags, devIDs)
	if err != nil {
		return []models.Device{}, "", err
	}

	//Check if the Devices has a capablity 33 (ome advance license)
	targetDevices, err := validateDevicesCapablity(devIDs, serviceTags, omeClient)
	if err != nil {
		return []models.Device{}, "", err
	}
	return targetDevices, usedDeviceInput, err
}

func validateNotification(plan models.ConfigureBaselines) error {
	if !plan.Schedule.ValueBool() {
		if !plan.Cron.IsNull() || !plan.EmailAddresses.IsNull() {
			return fmt.Errorf(clients.ErrBaseLineScheduleValid)
		}
	} else {
		if !(plan.NotifyOnSchedule.ValueBool() || plan.Cron.IsNull()) {
			return fmt.Errorf(clients.ErrBaseLineNotifyValid)
		}
	}
	return nil
}
