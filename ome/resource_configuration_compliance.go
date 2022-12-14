package ome

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	//NoOfTriesToGetBaselineStatus to get the task id
	NoOfTriesToGetBaselineStatus = 12
	//NotInventoried to identify is the baseline is inventored
	NotInventoried = "NOT_INVENTORIED"
	//Compliant - status compliant
	Compliant = "Compliant"
	//NonCompliant - status non compliant
	NonCompliant = "Non Compliant"
)

type resourceConfigurationComplianceType struct{}

func (r resourceConfigurationComplianceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Resource for managing configuration baselines remediation. Updates are supported for the following parameters: `target_devices`, `job_retry_count`, `sleep_interval`, `run_later`, `cron`.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the resource.",
				Description:         "ID of the resource.",
				Type:                types.StringType,
				Computed:            true,
			},
			"baseline_name": {
				MarkdownDescription: "Name of the Baseline.",
				Description:         "Name of the Baseline.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"baseline_id": {
				MarkdownDescription: "Id of the Baseline.",
				Description:         "Id of the Baseline.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"target_devices": {
				MarkdownDescription: "Target devices to be remediated.",
				Description:         "Target devices to be remediated.",
				Required:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"device_service_tag": {
						Type:                types.StringType,
						MarkdownDescription: "Target device servicetag to be remediated.",
						Description:         "Target device servicetag to be remediated.",
						Required:            true,
					},
					"compliance_status": {
						Type:                types.StringType,
						MarkdownDescription: "End compliance status of the target device, used to check the drifts in the compliance status.",
						Description:         "End compliance status of the target device, used to check the drifts in the compliance status.",
						Required:            true,
						Validators: []tfsdk.AttributeValidator{
							complianceStateValidator{},
						},
					},
				}),
				Validators: []tfsdk.AttributeValidator{
					sizeAtLeastValidator{min: 1},
				},
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
			"run_later": {
				MarkdownDescription: "Provides options to schedule the remediation task immediately, or at a specified time.",
				Description:         "Provides options to schedule the remediation task immediately, or at a specified time.",
				Type:                types.BoolType,
				Optional:            true,
			},
			"cron": {
				MarkdownDescription: "Cron to schedule the remediation task.",
				Description:         "Cron to schedule the remediation task.",
				Type:                types.StringType,
				Optional:            true,
			},
		},
	}, nil
}

// New resource instance
func (r resourceConfigurationComplianceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceConfigurationCompliance{
		p: *(p.(*provider)),
	}, nil
}

type resourceConfigurationCompliance struct {
	p provider
}

// Create a new resource
func (r resourceConfigurationCompliance) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_configuration_compliance: create started")
	var plan models.ConfigurationRemediation
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		fmt.Println(clients.ErrPlanToTfsdkConversion)
		return
	}
	if plan.RunLater.Value && plan.Cron.Value == "" {
		resp.Diagnostics.AddError(
			clients.ErrGnrBaseLineCreateRemediation,
			clients.ErrCronRequired,
		)
		return
	}

	state := models.ConfigurationRemediation{}

	//Create Session and defer the remove session
	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_compliance: create Creating Session")
	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	baseline, err := checkValidBaseline(omeClient, plan.BaselineName.Value, plan.BaselineID.Value, false)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrBaseLineCreateRemediation,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_compliance create: all baselines data is valid")
	var targetDevices []string
	for _, td := range plan.TargetDevices {
		targetDevices = append(targetDevices, td.DeviceServiceTag.Value)
	}

	targetDeviceIDs, err := checkValidDevices(omeClient, targetDevices, baseline)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrBaseLineCreateRemediation,
			err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "resource_configuration_compliance create: all target devices are valid")

	crp := getRemediationPayload(baseline.ID, targetDeviceIDs, plan.RunLater.Value, plan.Cron.Value)

	tflog.Trace(ctx, "resource_configuration_compliance create: triggered remediation", map[string]interface{}{
		"payload": crp,
	})

	jobID, err := omeClient.RemediateBaseLineDevices(crp)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrBaseLineCreateRemediation,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_compliance create: Job created", map[string]interface{}{
		"jobID": jobID,
	})
	if jobID != 0 && !plan.RunLater.Value {
		tflog.Trace(ctx, "resource_configuration_compliance create: Job track started")
		isSuccess, err := omeClient.TrackJob(jobID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
		if !isSuccess {
			tflog.Trace(ctx, "resource_configuration_compliance create: Job track errored", map[string]interface{}{
				"err": err,
			})
			resp.Diagnostics.AddError(
				clients.ErrGnrBaseLineCreateRemediation,
				err,
			)
		}
	}

	tflog.Trace(ctx, "resource_configuration_compliance create: saving state")
	state = plan
	state.BaselineID = types.Int64{Value: baseline.ID}
	state.BaselineName = types.String{Value: baseline.Name}
	state.ID = types.String{Value: fmt.Sprintf("%d", baseline.ID)}

	//save the data into state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_configuration_compliance create: create finished")
}

// Read resource information
func (r resourceConfigurationCompliance) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	//Get State Data
	tflog.Trace(ctx, "resource_configuration_compliance: read started")
	var state models.ConfigurationRemediation
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

	tflog.Trace(ctx, "resource_configuration_compliance: read checking status report")
	//check the compliance status to check if the reports are generated
	err = checkReportsStatus(omeClient, state.BaselineID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrBaseLineReadRemediation,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_compliance: read checking status finshed")

	for i, td := range state.TargetDevices {
		deviceReport, err := omeClient.GetConfiBaselineDeviceReport(state.BaselineID.Value, td.DeviceServiceTag.Value)
		if err != nil {
			if err != nil {
				resp.Diagnostics.AddError(
					clients.ErrGnrBaseLineReadRemediation,
					err.Error(),
				)
				return
			}
		}
		compliantStatus := Compliant
		if deviceReport.ComplianceStatus != 1 {
			compliantStatus = NonCompliant
		}
		state.TargetDevices[i] = models.TargetDevices{
			DeviceServiceTag: td.DeviceServiceTag,
			ComplianceStatus: types.String{Value: compliantStatus},
		}
	}

	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_configuration_compliance: read finished")
}

// Update resource
func (r resourceConfigurationCompliance) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	//Get state Data
	tflog.Trace(ctx, "resource_configuration_compliance: update started")
	var state models.ConfigurationRemediation
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan Data
	var plan models.ConfigurationRemediation
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.RunLater.Value && plan.Cron.Value == "" {
		resp.Diagnostics.AddError(
			clients.ErrBaseLineUpdateRemediation,
			clients.ErrCronRequired,
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

	tflog.Trace(ctx, "resource_configuration_compliance: update checking if baseline name or id is changed")
	if (plan.BaselineID.Value != 0 && plan.BaselineID.Value != state.BaselineID.Value) || (plan.BaselineName.Value != "" && plan.BaselineName.Value != state.BaselineName.Value) {
		resp.Diagnostics.AddError(
			clients.ErrBaseLineUpdateRemediation,
			clients.ErrBaseLineModified,
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_compliance: update checking for valid baseline")
	baseline, err := checkValidBaseline(omeClient, plan.BaselineName.Value, plan.BaselineID.Value, true)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrBaseLineCreateRemediation,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_compliance: update checking for terget devices")

	var targetDevices []string
	for _, td := range plan.TargetDevices {
		targetDevices = append(targetDevices, td.DeviceServiceTag.Value)
	}

	targetDeviceIDs, err := checkValidDevices(omeClient, targetDevices, baseline)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrBaseLineCreateRemediation,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_compliance: target devices are valid")
	crp := getRemediationPayload(baseline.ID, targetDeviceIDs, plan.RunLater.Value, plan.Cron.Value)

	tflog.Trace(ctx, "resource_configuration_compliance: update remidiation started", map[string]interface{}{
		"payload": crp,
	})
	jobID, err := omeClient.RemediateBaseLineDevices(crp)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrBaseLineUpdateRemediation,
			err.Error(),
		)
	}

	tflog.Trace(ctx, "resource_configuration_compliance: update remidiation job created", map[string]interface{}{
		"jobID": jobID,
	})
	if jobID != 0 && !plan.RunLater.Value {
		isSuccess, err := omeClient.TrackJob(jobID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
		if !isSuccess {
			tflog.Trace(ctx, "resource_configuration_compliance: update remidiation job failed", map[string]interface{}{
				"err": err,
			})
			resp.Diagnostics.AddError(
				clients.ErrBaseLineUpdateRemediation,
				err,
			)
		}
	}

	tflog.Trace(ctx, "resource_configuration_compliance: update remidiation state updating")
	state = plan
	state.BaselineID = types.Int64{Value: baseline.ID}
	state.BaselineName = types.String{Value: baseline.Name}
	state.ID = types.String{Value: fmt.Sprintf("%d", baseline.ID)}

	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_configuration_compliance: update finished")
}

// Delete resource
func (r resourceConfigurationCompliance) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceConfigurationCompliance) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	var state models.ConfigurationRemediation
	_ = req.ID

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

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// validate the Devices against the baseline
// validate the baseline
// non compliance devices is handled by remediation API
func checkValidDevices(omeClient *clients.Client, targetDevices []string, baseline models.OmeBaseline) ([]int64, error) {

	var baselineDevices []int64
	var targetDeviceIDs []int64
	deviceIDServiceTagsMap := map[int64]string{}

	for _, bt := range baseline.BaselineTargets {
		baselineDevices = append(baselineDevices, bt.ID)
	}

	for _, st := range targetDevices {
		device, err := omeClient.GetDevice(st, 0)
		if err != nil {
			return []int64{}, err
		}
		targetDeviceIDs = append(targetDeviceIDs, device.ID)
		deviceIDServiceTagsMap[device.ID] = device.DeviceServiceTag
	}

	diffDeviceIDs := clients.CompareInt64(targetDeviceIDs, baselineDevices)
	diffDeviceServiceTags := []string{}
	for _, diffDeviceID := range diffDeviceIDs {
		diffDeviceServiceTags = append(diffDeviceServiceTags, deviceIDServiceTagsMap[diffDeviceID])
	}
	if len(diffDeviceIDs) != 0 {
		return []int64{}, fmt.Errorf(clients.ErrBaseLineInvalidDevices, diffDeviceServiceTags)
	}

	return targetDeviceIDs, nil
}

func checkValidBaseline(omeClient *clients.Client, baselineName string, baseLineID int64, checkreportStatus bool) (models.OmeBaseline, error) {
	var baseline models.OmeBaseline
	var err error
	if baseLineID != 0 && baselineName != "" {
		return models.OmeBaseline{}, fmt.Errorf(clients.ErrBaseLineInvalid)
	}
	if baseLineID == 0 && baselineName == "" {
		return models.OmeBaseline{}, fmt.Errorf(clients.ErrBaseLineInvalid)
	}
	if baseLineID != 0 {
		baseline, err = omeClient.GetBaselineByID(baseLineID)
	} else {
		baseline, err = omeClient.GetBaselineByName(baselineName)
	}
	if err != nil {
		return models.OmeBaseline{}, err
	}

	if checkreportStatus {
		if strings.ToUpper(baseline.ConfigComplianceSummary.ComplianceStatus) == NotInventoried {
			return models.OmeBaseline{}, fmt.Errorf(clients.ErrBaseLineReportInProgress)
		}
	}
	return baseline, nil
}

func checkReportsStatus(omeClient *clients.Client, baselineID int64) error {
	var baseline models.OmeBaseline
	var err error
	var complianceStatus string
	tries := 1
	baseline, err = omeClient.GetBaselineByID(baselineID)
	if err != nil {
		return err
	}
	complianceStatus = baseline.ConfigComplianceSummary.ComplianceStatus
	for strings.ToUpper(complianceStatus) == NotInventoried && NoOfTriesToGetBaselineStatus != tries {
		tries++
		time.Sleep(10 * time.Second) // sleep for 10 secs
		baseline, err = omeClient.GetBaselineByID(baselineID)
		if err != nil {
			return err
		}
		complianceStatus = baseline.ConfigComplianceSummary.ComplianceStatus
	}
	if strings.ToUpper(baseline.ConfigComplianceSummary.ComplianceStatus) == NotInventoried {
		return fmt.Errorf(clients.ErrBaseLineReportInProgress)
	}
	return nil
}

func getRemediationPayload(baselineID int64, targetDeviceIDs []int64, runLater bool, cron string) models.ConfigurationRemediationPayload {
	crp := models.ConfigurationRemediationPayload{
		ID:        baselineID,
		DeviceIDS: targetDeviceIDs,
		Schedule: models.OMESchedule{
			RunNow:   true,
			RunLater: false,
		},
	}
	if runLater {
		crp.Schedule.RunNow = false
		crp.Schedule.RunLater = true
		crp.Schedule.Cron = cron
		crp.Schedule.StartTime = ""
		crp.Schedule.EndTime = ""
	}
	return crp
}
