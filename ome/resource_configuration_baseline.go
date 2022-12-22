package ome

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	//NoOFTries to get the task id
	NoOFTries = 5
)

type resourceConfigurationBaselineType struct{}

// Template Deployment Resource schema
func (r resourceConfigurationBaselineType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource for managing configuration baselines on OpenManage Enterprise. Updates are supported for all the parameters. When `schedule` is `true`, following parameters are considered: `notify_on_schedule`, `cron`, `email_addresses`, `output_format`",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the resource.",
				Description:         "ID of the resource.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"ref_template_id": {
				MarkdownDescription: "Reference template ID.",
				Description:         "Reference template ID.",
				Type:                types.Int64Type,
				Computed:            true,
				Optional:            true,
			},
			"ref_template_name": {
				MarkdownDescription: "Reference template name.",
				Description:         "Reference template name.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"baseline_name": {
				MarkdownDescription: "Name of the Baseline.",
				Description:         "Name of the Baseline.",
				Type:                types.StringType,
				Required:            true,
			},
			"description": {
				MarkdownDescription: "Description of the baseline.",
				Description:         "Description of the baseline.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"device_ids": {
				MarkdownDescription: "List of the device id on which the baseline compliance needs to be run.",
				Description:         "List of the device id on which the baseline compliance needs to be run.",
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
				Optional: true,
			},
			"device_servicetags": {
				MarkdownDescription: "List of the device servicetag on which the baseline compliance needs to be run.",
				Description:         "List of the device servicetag on which the baseline compliance needs to be run.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"schedule": {
				MarkdownDescription: "Schedule notification via email.",
				Description:         "Schedule notification via email.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Bool{Value: false}),
				},
			},
			"notify_on_schedule": {
				MarkdownDescription: "Schedule notification via cron or any time the baseline becomes non-compliant.",
				Description:         "Schedule notification via cron or any time the baseline becomes non-compliant.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Bool{Value: false}),
				},
			},
			"email_addresses": {
				MarkdownDescription: "Email addresses for notification.",
				Description:         "Email addresses for notification.",
				Type: types.SetType{
					ElemType: types.StringType},
				Optional: true,
			},
			"output_format": {
				MarkdownDescription: "Output format type, the input is case senitive.",
				Description:         "Output format type, the input is case senitive.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.String{Value: "html"}),
				},
				Validators: []tfsdk.AttributeValidator{
					outputFormatValidator{},
				},
			},
			"cron": {
				MarkdownDescription: "Cron expression for notification schedule.",
				Description:         "Cron expression for notification schedule.",
				Type:                types.StringType,
				Optional:            true,
			},
			"task_id": {
				MarkdownDescription: "Task id associated with baseline.",
				Description:         "Task id associated with baseline.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"job_retry_count": {
				MarkdownDescription: "Number of times the job has to be polled to get the final status of the resource.",
				Description:         "Number of times the job has to be polled to get the final status of the resource.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 30}),
				},
			},
			"sleep_interval": {
				MarkdownDescription: "Sleep time interval for job polling in seconds.",
				Description:         "Sleep time interval for job polling in seconds.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 20}),
				},
			},
		},
	}, nil
}

// New resource instance
func (r resourceConfigurationBaselineType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceConfigurationBaseline{
		p: *(p.(*provider)),
	}, nil
}

type resourceConfigurationBaseline struct {
	p provider
}

// Create a new resource
func (r resourceConfigurationBaseline) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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
	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_baseline create Creating Session")
	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_configuration_baseline create Validating Template Details")
	omeTemplate, err := validateRefTemplateDetails(plan.RefTemplateID.Value, plan.RefTemplateName.Value, omeClient)
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
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = plan.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
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

	isSuccess, message := omeClient.TrackJob(baseline.TaskID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
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
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceConfigurationBaseline) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	//Get State Data
	var state models.ConfigureBaselines
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and differ the remove session
	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	var usedDeviceInput string

	if len(state.DeviceIDs.Elems) > 0 {
		usedDeviceInput = clients.DeviceIDs
	} else if len(state.DeviceServicetags.Elems) > 0 {
		usedDeviceInput = clients.ServiceTags
	}
	baseline, err := omeClient.GetBaselineByID(state.ID.Value)
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
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceConfigurationBaseline) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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
	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	tflog.Trace(ctx, "resource_configuration_baseline update checking the job status")
	tflog.Debug(ctx, "resource_configuration_baseline checking job status for", map[string]interface{}{
		"jobid": state.TaskID.Value,
	})

	if state.TaskID.Value != 0 {
		jr, err := omeClient.GetJob(state.TaskID.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateBaseline,
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, "resource_configuration_baseline update job status is", map[string]interface{}{
			"jobid":  state.TaskID.Value,
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
	omeTemplate, err := validateRefTemplateDetails(plan.RefTemplateID.Value, plan.RefTemplateName.Value, omeClient)
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
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = plan.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
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
	cb.ID = state.ID.Value // For update case

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

	isSuccess, message := omeClient.TrackJob(baseline.TaskID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
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
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceConfigurationBaseline) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Get State Data
	var state models.ConfigureBaselines
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and differ the remove session
	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	err = omeClient.DeleteBaseline([]int64{state.ID.Value})

	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrDeleteBaseline,
			err.Error(),
		)
	}
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceConfigurationBaseline) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	var state models.ConfigureBaselines
	baselineName := req.ID

	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return
	}

	defer omeClient.RemoveSession()

	baseline, err := omeClient.GetBaselineByName(baselineName)
	if err != nil {
		resp.Diagnostics.AddError(clients.ErrImportDeployment, err.Error())
		return
	}
	updateBaselineState(ctx, &state, &state, baseline, clients.ServiceTags, omeClient)
	//Save into State
	state.EmailAddresses.ElemType = types.StringType
	state.DeviceIDs.ElemType = types.Int64Type
	state.OutputFormat = types.String{Value: "html"}
	state.JobRetryCount = types.Int64{Value: 30}
	state.SleepInterval = types.Int64{Value: 20}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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

func updateBaselineState(ctx context.Context, state *models.ConfigureBaselines, plan *models.ConfigureBaselines, omeBaseline models.OmeBaseline, usedDeviceInput string, omeClient *clients.Client) {
	state.ID = types.Int64{Value: omeBaseline.ID}
	state.RefTemplateID = types.Int64{Value: omeBaseline.TemplateID}
	state.RefTemplateName = types.String{Value: omeBaseline.TemplateName}
	state.Description = types.String{Value: omeBaseline.Description}
	state.BaselineName = types.String{Value: omeBaseline.Name}

	if usedDeviceInput == clients.ServiceTags {
		devSTsTfsdk := types.Set{
			ElemType: types.StringType,
		}
		apiDeviceIDs := map[string]models.Device{}
		devSts := []string{}
		deviceStVals := []attr.Value{}
		for _, bTarget := range omeBaseline.BaselineTargets {
			device, _ := omeClient.GetDevice("", bTarget.ID)
			apiDeviceIDs[device.DeviceServiceTag] = device
		}

		plan.DeviceServicetags.ElementsAs(ctx, &devSts, true)

		for _, devSt := range devSts {
			if val, ok := apiDeviceIDs[devSt]; ok {
				deviceStVals = append(deviceStVals, types.String{Value: val.DeviceServiceTag})
				delete(apiDeviceIDs, devSt)
			}
		}

		if len(apiDeviceIDs) != 0 {
			for _, val := range apiDeviceIDs {
				deviceStVals = append(deviceStVals, types.String{Value: val.DeviceServiceTag})
			}
		}

		devSTsTfsdk.Elems = deviceStVals

		state.DeviceServicetags = devSTsTfsdk
		state.DeviceIDs = plan.DeviceIDs
	} else {
		devIDsTfsdk := types.Set{
			ElemType: types.Int64Type,
		}
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
				deviceIDVals = append(deviceIDVals, types.Int64{Value: val.ID})
				delete(apiDeviceIDs, devID)
			}
		}

		if len(apiDeviceIDs) != 0 {
			for _, val := range apiDeviceIDs {
				deviceIDVals = append(deviceIDVals, types.Int64{Value: val.ID})
			}
		}

		devIDsTfsdk.Elems = deviceIDVals
		state.DeviceIDs = devIDsTfsdk
		state.DeviceServicetags = plan.DeviceServicetags
	}

	notificationSettings := omeBaseline.NotificationSettings
	if notificationSettings != nil {
		state.Schedule = types.Bool{Value: true}
		state.OutputFormat = types.String{Value: notificationSettings.OutputFormat}

		emailAddress := []attr.Value{}
		for _, v := range notificationSettings.EmailAddresses {
			emailAddress = append(emailAddress, types.String{Value: v})
		}

		emailTfsdk := types.Set{
			ElemType: types.StringType,
		}
		emailTfsdk.Elems = emailAddress
		state.EmailAddresses = emailTfsdk

		if notificationSettings.NotificationType == "NOTIFY_ON_NON_COMPLIANCE" {
			state.NotifyOnSchedule = types.Bool{Value: false}
		} else {
			state.NotifyOnSchedule = types.Bool{Value: true}
		}

		if notificationSettings.Schedule.Cron != "" {
			state.Cron = types.String{Value: notificationSettings.Schedule.Cron}
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
	state.TaskID = types.Int64{Value: omeBaseline.TaskID}
	state.JobRetryCount = plan.JobRetryCount
	state.SleepInterval = plan.SleepInterval
}

func getPayload(ctx context.Context, plan *models.ConfigureBaselines, templateID int64, targetDevices []models.Device) (models.ConfigurationBaselinePayload, error) {
	cbp := models.ConfigurationBaselinePayload{
		Name:        plan.BaselineName.Value,
		Description: plan.Description.Value,
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

	if plan.Schedule.Value {
		if len(plan.EmailAddresses.Elems) == 0 {
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

		if plan.NotifyOnSchedule.Value && plan.Cron.Value == "" {
			return models.ConfigurationBaselinePayload{}, fmt.Errorf(clients.ErrInvalidCronExpression)
		}

		notificationType := "NOTIFY_ON_NON_COMPLIANCE"
		if plan.NotifyOnSchedule.Value {
			notificationType = "NOTIFY_ON_SCHEDULE"
		}

		notificationSettings := models.NotificationSettings{
			NotificationType: notificationType,
			EmailAddresses:   emailaddressList,
			Schedule: models.BaselineNotificationSchedule{
				Cron: plan.Cron.Value,
			},
			OutputFormat: strings.ToUpper(plan.OutputFormat.Value),
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
	if !plan.Schedule.Value {
		if !plan.Cron.Null || !plan.EmailAddresses.Null {
			return fmt.Errorf(clients.ErrBaseLineScheduleValid)
		}
	} else {
		if !plan.NotifyOnSchedule.Value && !plan.Cron.Null {
			return fmt.Errorf(clients.ErrBaseLineNotifyValid)
		}
	}
	return nil
}
