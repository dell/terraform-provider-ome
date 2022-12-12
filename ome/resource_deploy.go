package ome

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceDeploymentType struct{}

// Template Deployment Resource schema
func (r resourceDeploymentType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource for managing template deployment on OpenManage Enterprise. Updates are supported for the following parameters: `device_ids`, `device_servicetags`, `boot_to_network_iso`, `forced_shutdown`, `options_time_to_wait_before_shutdown`, `power_state_off`, `options_precheck_only`, `options_strict_checking_vlan`, `options_continue_on_warning`, `run_later`, `cron`, `device_attributes`, `job_retry_count`, `sleep_interval`.",
		Version:             1,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the resource.",
				Description:         "ID of the resource.",
				Type:                types.StringType,
				Computed:            true,
			},
			"template_id": {
				MarkdownDescription: "ID of the existing template.",
				Description:         "ID of the existing template.",
				Type:                types.Int64Type,
				Computed:            true,
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"template_name": {
				MarkdownDescription: "Name of the existing template.",
				Description:         "Name of the existing template.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"device_ids": {
				MarkdownDescription: "List of the device id(s).",
				Description:         "List of the device id(s).",
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
				Optional: true,
			},
			"device_servicetags": {
				MarkdownDescription: "List of the device servicetags.",
				Description:         "List of the device servicetags.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"boot_to_network_iso": {
				MarkdownDescription: "Boot To Network ISO deployment details.",
				Description:         "Boot To Network ISO deployment details.",
				Optional:            true,
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"boot_to_network": types.BoolType,
						"share_type":      types.StringType,
						"iso_timeout":     types.Int64Type,
						"iso_path":        types.StringType,
						"share_detail": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"ip_address": types.StringType,
								"share_name": types.StringType,
								"work_group": types.StringType,
								"user":       types.StringType,
								"password":   types.StringType,
							}},
					},
				},
			},
			"forced_shutdown": {
				MarkdownDescription: "Force shutdown after deployment.",
				Description:         "Force shutdown after deployment.",
				Type:                types.BoolType,
				Optional:            true,
			},
			"options_time_to_wait_before_shutdown": {
				MarkdownDescription: "Option to specify the time to wait before shutdown in seconds. Default and minimum value is 300 and maximum is 3600 seconds respectively.",
				Description:         "Option to specify the time to wait before shutdown in seconds. Default and minimum value is 300 and maximum is 3600 seconds respectively.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 300}),
				},
			},
			"power_state_off": {
				MarkdownDescription: "End power state of a target devices. Default power state is ON. Make it true to switch it to OFF state.",
				Description:         "End power state of a target devices. Default power state is ON. Make it true to switch it to OFF state.",
				Type:                types.BoolType,
				Optional:            true,
			},
			"options_precheck_only": {
				MarkdownDescription: "Option to precheck",
				Description:         "Option to precheck",
				Type:                types.BoolType,
				Optional:            true,
			},
			"options_strict_checking_vlan": {
				MarkdownDescription: "Checks the strict association of vlan.",
				Description:         "Checks the strict association of vlan.",
				Type:                types.BoolType,
				Optional:            true,
			},
			"options_continue_on_warning": {
				MarkdownDescription: "Continue to run the job on warnings.",
				Description:         "Continue to run the job on warnings.",
				Type:                types.BoolType,
				Optional:            true,
			},
			"run_later": {
				MarkdownDescription: "Provides options to schedule the deployment task immediately, or at a specified time.",
				Description:         "Provides options to schedule the deployment task immediately, or at a specified time.",
				Type:                types.BoolType,
				Optional:            true,
			},
			"cron": {
				MarkdownDescription: "Cron to schedule the deployment task. Cron expression should be of future datetime.",
				Description:         "Cron to schedule the deployment task. Cron expression should be of future datetime.",
				Type:                types.StringType,
				Optional:            true,
			},
			"device_attributes": {
				MarkdownDescription: "List of template attributes associated with the target devices for deploymnent.",
				Description:         "List of template attributes associated with the target devices for deploymnent.",
				Optional:            true,
				Type: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"device_servicetags": types.SetType{
								ElemType: types.StringType,
							},
							"attributes": types.ListType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"attribute_id": types.Int64Type,
										"display_name": types.StringType,
										"value":        types.StringType,
										"is_ignored":   types.BoolType,
									},
								},
							},
						},
					},
				},
			},
			"job_retry_count": {
				MarkdownDescription: "Number of times the job has to be polled to get the final status of the resource.",
				Description:         "Number of times the job has to be polled to get the final status of the resource.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 20}),
				},
			},
			"sleep_interval": {
				MarkdownDescription: "Sleep time interval for job polling in seconds.",
				Description:         "Sleep time interval for job polling in seconds.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 60}),
				},
			},
		},
	}, nil
}

// New resource instance
func (r resourceDeploymentType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDeployment{
		p: *(p.(*provider)),
	}, nil
}

type resourceDeployment struct {
	p provider
}

// Create a new resource
func (r resourceDeployment) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Trace(ctx, "resource_deploy create : Started")
	//Get Plan Data
	var plan models.TemplateDeployment
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateDeploymentState := models.TemplateDeployment{}

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

	tflog.Trace(ctx, "resource_deploy create: session created")

	omeTemplate, err := omeClient.GetTemplateByIDOrName(plan.TemplateID.Value, plan.TemplateName.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrInvalidTemplate,
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "resource_deploy create: ome Template", map[string]interface{}{
		"id":   omeTemplate.ID,
		"name": omeTemplate.Name,
	})

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

	usedDeviceInput, err := clients.DeviceMutuallyExclusive(serviceTags, devIDs)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentCreate, err.Error(),
		)
		return
	}

	devices, err := omeClient.GetDevices(serviceTags, devIDs, []string{})
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentCreate, err.Error(),
		)
		return
	}

	_, deviceIDs, _ := omeClient.GetUniqueDevicesIdsAndServiceTags(devices)

	options := getOptions(plan)

	deploymentRequest := models.OMETemplateDeployRequest{
		ID:        omeTemplate.ID,
		TargetIDS: deviceIDs,
		Options:   options,
	}

	//Boot to Network ISO starts
	if !plan.BootToNetworkISO.Null {
		bootToNetworkISOModel, diags, err := getBootToNetworkISO(ctx, plan)
		if err != nil {
			resp.Diagnostics.Append(diags...)
			resp.Diagnostics.AddWarning(
				clients.ErrUnableToParseData,
				err.Error(),
			)
		}
		deploymentRequest.NetworkBootISOModel = bootToNetworkISOModel
	}
	//Boot to Network ISO Ends
	// Schedule Starts
	if plan.RunLater.Value {
		deploymentRequest.Schedule = getSchedule(plan)
	}
	//Schedule ends
	// Device Attrs
	if len(plan.DeviceAttributes.Elems) > 0 {
		deploymentRequest.Attributes = getDeviceAttributes(ctx, devices, plan)
	}

	tflog.Trace(ctx, "resource_deploy create: started creating deployment job")

	deploymentJobID, err := omeClient.CreateDeployment(deploymentRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentCreate, err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "resource_deploy create: finished creating deployment job", map[string]interface{}{
		"deploymentJobID": deploymentJobID,
	})

	if !plan.RunLater.Value {
		tflog.Trace(ctx, "resource_deploy create: started job tracking")
		isSuccess, message := omeClient.TrackJob(deploymentJobID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
		if !isSuccess {
			resp.Diagnostics.AddWarning(
				clients.ErrTemplateDeploymentCreate, message,
			)
		}
	}

	tflog.Trace(ctx, "resource_deploy create: updating state started")

	_ = updateDeploymentState(&templateDeploymentState, &plan, omeTemplate.ID, omeTemplate.Name, omeClient, usedDeviceInput)

	tflog.Trace(ctx, "resource_deploy create: updating state finished, saving ...")
	// Save into State
	diags = resp.State.Set(ctx, &templateDeploymentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_deploy create: finish")
}

// Read resource information
func (r resourceDeployment) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	//Get State Data
	tflog.Trace(ctx, "resource_deploy read: started")
	var stateTemplateDeployment models.TemplateDeployment
	diags := req.State.Get(ctx, &stateTemplateDeployment)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := stateTemplateDeployment.TemplateID.Value
	templateName := stateTemplateDeployment.TemplateName.Value

	tflog.Debug(ctx, "resource_deploy read: reading a template", map[string]interface{}{
		"id":   templateID,
		"name": templateName,
	})
	var serviceTags []string
	var devIDs []int64

	diags = stateTemplateDeployment.DeviceServicetags.ElementsAs(ctx, &serviceTags, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = stateTemplateDeployment.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var usedDeviceInput string
	if len(serviceTags) > 0 {
		usedDeviceInput = clients.ServiceTags
	} else if len(devIDs) > 0 {
		usedDeviceInput = clients.DeviceIDs
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

	tflog.Trace(ctx, "resource_deploy read: client created started updating state")
	_ = updateDeploymentState(&stateTemplateDeployment, &stateTemplateDeployment, templateID, templateName, omeClient, usedDeviceInput)

	tflog.Trace(ctx, "resource_deploy read: finished reading state")
	//Save into State
	diags = resp.State.Set(ctx, &stateTemplateDeployment)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_deploy read: finished")
}

// Update resource
func (r resourceDeployment) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	//Get state Data
	tflog.Trace(ctx, "resource_deploy update: started")
	var state models.TemplateDeployment
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan Data
	var plan models.TemplateDeployment
	diags = req.Plan.Get(ctx, &plan)
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

	if (plan.TemplateID.Value != 0 && plan.TemplateID.Value != state.TemplateID.Value) || (plan.TemplateName.Value != "" && plan.TemplateName.Value != state.TemplateName.Value) {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentUpdate,
			clients.ErrTemplateChanges,
		)
		return
	}
	tflog.Debug(ctx, "resource_deploy update: started with template", map[string]interface{}{
		"id":   plan.TemplateID.Value,
		"name": plan.TemplateName.Value,
	})

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

	usedDeviceInput, err := clients.DeviceMutuallyExclusive(serviceTags, devIDs)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentUpdate, err.Error(),
		)
		return
	}

	planDevices, err := omeClient.GetDevices(serviceTags, devIDs, []string{})
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentUpdate, err.Error(),
		)
		return
	}

	_, planDeviceIDs, _ := omeClient.GetUniqueDevicesIdsAndServiceTags(planDevices)
	//get the state devicds's

	var stateServiceTags []string
	var stateDevIDs []int64

	diags = state.DeviceServicetags.ElementsAs(ctx, &stateServiceTags, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = state.DeviceIDs.ElementsAs(ctx, &stateDevIDs, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	serverProfiles, err := omeClient.GetServerProfileInfoByTemplateName(state.TemplateName.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentUpdate, err.Error(),
		)
		return
	}

	var stateDeviceIDs []int64

	for _, serverProfile := range serverProfiles.Value {
		stateDeviceIDs = append(stateDeviceIDs, serverProfile.TargetID)
	}

	newDeployDevIDs := compare(planDeviceIDs, stateDeviceIDs)

	//remove
	removeDeployDevIDs := compare(stateDeviceIDs, planDeviceIDs)

	tflog.Debug(ctx, "resource_deploy update: Target ids", map[string]interface{}{
		"newTargets":    newDeployDevIDs,
		"removeTargets": removeDeployDevIDs,
	})

	options := getOptions(plan)

	deploymentRequest := models.OMETemplateDeployRequest{
		ID:        state.TemplateID.Value,
		TargetIDS: newDeployDevIDs,
		Options:   options,
	}

	//Boot to Network ISO starts
	if !plan.BootToNetworkISO.Null {
		bootToNetworkISOModel, diags, err := getBootToNetworkISO(ctx, plan)
		if err != nil {
			resp.Diagnostics.Append(diags...)
			resp.Diagnostics.AddWarning(
				clients.ErrUnableToParseData,
				err.Error(),
			)
		}
		deploymentRequest.NetworkBootISOModel = bootToNetworkISOModel
	}
	//Boot to Network ISO Ends
	// Schedule Starts
	if plan.RunLater.Value {
		deploymentRequest.Schedule = getSchedule(plan)
	}
	//Schedule ends
	// Device Attrs
	if len(plan.DeviceAttributes.Elems) > 0 {
		deploymentRequest.Attributes = getDeviceAttributes(ctx, planDevices, plan)
	}

	if len(newDeployDevIDs) > 0 {
		tflog.Trace(ctx, "resource_deploy update: started deployment")
		deploymentJobID, err := omeClient.CreateDeployment(deploymentRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrTemplateDeploymentUpdate, err.Error(),
			)
			return
		}

		if !plan.RunLater.Value {
			tflog.Trace(ctx, "resource_deploy update: started job tracking")
			isSuccess, message := omeClient.TrackJob(deploymentJobID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
			if !isSuccess {
				resp.Diagnostics.AddWarning(
					"unable to complete the deployment for the template: ", message,
				)
			}
		}
	}

	if len(removeDeployDevIDs) > 0 {
		profileArr := make([]int64, len(removeDeployDevIDs))
		for i, removeDevID := range removeDeployDevIDs {
			for _, serverProfile := range serverProfiles.Value {
				if serverProfile.TargetID == removeDevID {
					profileArr[i] = serverProfile.ID
				}
			}
		}
		tflog.Debug(ctx, "resource_deploy update: deleting server profiles", map[string]interface{}{
			"profileIds": profileArr,
		})
		err = deleteProfiles(ctx, omeClient, profileArr)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrTemplateDeploymentUpdate,
				err.Error(),
			)
			return
		}
	}

	tflog.Trace(ctx, "resource_deploy update: started state update")

	_ = updateDeploymentState(&state, &plan, state.TemplateID.Value, state.TemplateName.Value, omeClient, usedDeviceInput)

	tflog.Trace(ctx, "resource_deploy update: finished state update")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_deploy update: finished")
}

// Delete resource
func (r resourceDeployment) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Trace(ctx, "resource_deploy delete: started")
	// Get State Data
	var statetemplateDeployment models.TemplateDeployment
	diags := req.State.Get(ctx, &statetemplateDeployment)
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

	tflog.Debug(ctx, "resource_deploy delete: started with template", map[string]interface{}{
		"id":   statetemplateDeployment.TemplateID.Value,
		"name": statetemplateDeployment.TemplateName.Value,
	})

	serverProfiles, err := omeClient.GetServerProfileInfoByTemplateName(statetemplateDeployment.TemplateName.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentDelete, err.Error(),
		)
		return
	}

	profileArr := make([]int64, len(serverProfiles.Value))

	for i, serverProfile := range serverProfiles.Value {
		profileArr[i] = serverProfile.ID
	}

	tflog.Debug(ctx, "resource_deploy delete: deleting server profiles", map[string]interface{}{
		"profileIds": profileArr,
	})

	err = deleteProfiles(ctx, omeClient, profileArr)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentDelete,
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_deploy delete: finished")
}

func deleteProfiles(ctx context.Context, omeClient *clients.Client, profileArr []int64) error {
	pdr := models.ProfileDeleteRequest{
		ProfileIds: profileArr,
	}
	err := omeClient.DeleteDeployment(pdr)
	if err != nil {
		return err
	}
	return nil
}

// Import resource
func (r resourceDeployment) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Trace(ctx, "resource_deploy import: started")
	// Save the import identifier in the id attribute
	var stateTemplateDeployment models.TemplateDeployment
	templateName := req.ID
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

	omeTemplate, err := omeClient.GetTemplateByName(templateName)
	if err != nil {
		resp.Diagnostics.AddError(clients.ErrImportDeployment, err.Error())
		return
	}
	templateID := omeTemplate.ID

	//Create Session and differ the remove session
	defer omeClient.RemoveSession()

	profileDevSTVals := []attr.Value{}
	serverProfiles, err := omeClient.GetServerProfileInfoByTemplateName(templateName)
	if err != nil {
		resp.Diagnostics.AddError(clients.ErrImportDeployment, err.Error())
		return
	}

	if len(serverProfiles.Value) == 0 {
		resp.Diagnostics.AddError(clients.ErrImportDeployment, fmt.Sprintf(clients.ErrImportNoProfiles, templateName))
		return
	}

	for _, serverProfile := range serverProfiles.Value {
		device, _ := omeClient.GetDevice("", serverProfile.TargetID)
		deviceSTVal := types.String{Value: device.DeviceServiceTag}
		profileDevSTVals = append(profileDevSTVals, deviceSTVal)
	}
	stateTemplateDeployment.ID.Value = strconv.FormatInt(templateID, 10)
	stateTemplateDeployment.TemplateID.Value = templateID
	stateTemplateDeployment.TemplateName.Value = templateName
	devSTsTfsdk := types.Set{
		ElemType: types.StringType,
	}
	devSTsTfsdk.Elems = profileDevSTVals
	stateTemplateDeployment.DeviceServicetags = devSTsTfsdk
	devIDsTfsdk := types.Set{
		ElemType: types.Int64Type,
	}
	devIDsTfsdk.Elems = []attr.Value{}
	stateTemplateDeployment.DeviceIDs = devIDsTfsdk

	// set empty device attributes
	attributesObjects := []attr.Value{}
	attributesObject := types.Object{
		AttrTypes: map[string]attr.Type{
			"attribute_id": types.Int64Type,
			"display_name": types.StringType,
			"value":        types.StringType,
			"is_ignored":   types.BoolType,
		},
		Attrs: map[string]attr.Value{
			"attribute_id": types.Int64{Value: 0},
			"display_name": types.String{Value: ""},
			"value":        types.String{Value: ""},
			"is_ignored":   types.Bool{Value: false},
		},
	}
	attributesObjects = append(attributesObjects, attributesObject)

	attributesTfsdk := types.List{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_id": types.Int64Type,
				"display_name": types.StringType,
				"value":        types.StringType,
				"is_ignored":   types.BoolType,
			},
		},
		Elems: attributesObjects,
	}

	deviceAttributeObjects := []attr.Value{}
	deviceAttributeObject := types.Object{
		AttrTypes: map[string]attr.Type{
			"device_servicetags": types.SetType{
				ElemType: types.StringType,
			},
			"attributes": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"attribute_id": types.Int64Type,
						"display_name": types.StringType,
						"value":        types.StringType,
						"is_ignored":   types.BoolType,
					},
				},
			},
		},
		Attrs: map[string]attr.Value{
			"device_servicetags": devSTsTfsdk,
			"attributes":         attributesTfsdk,
		},
	}

	deviceAttributeObjects = append(deviceAttributeObjects, deviceAttributeObject)

	deviceAttributeTfsdk := types.List{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"device_servicetags": types.SetType{
					ElemType: types.StringType,
				},
				"attributes": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"attribute_id": types.Int64Type,
							"display_name": types.StringType,
							"value":        types.StringType,
							"is_ignored":   types.BoolType,
						},
					},
				},
			},
		},
		Elems: deviceAttributeObjects,
	}
	stateTemplateDeployment.DeviceAttributes = deviceAttributeTfsdk

	shareDetailsTfsdk := types.Object{
		AttrTypes: map[string]attr.Type{
			"ip_address": types.StringType,
			"share_name": types.StringType,
			"work_group": types.StringType,
			"user":       types.StringType,
			"password":   types.StringType,
		},
		Attrs: map[string]attr.Value{
			"ip_address": types.String{Value: ""},
			"share_name": types.String{Value: ""},
			"work_group": types.String{Value: ""},
			"user":       types.String{Value: ""},
			"password":   types.String{Value: ""},
		},
	}

	bootToNetworkISOTfsdk := types.Object{
		AttrTypes: map[string]attr.Type{
			"boot_to_network": types.BoolType,
			"share_type":      types.StringType,
			"iso_timeout":     types.Int64Type,
			"iso_path":        types.StringType,
			"share_detail": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"ip_address": types.StringType,
					"share_name": types.StringType,
					"work_group": types.StringType,
					"user":       types.StringType,
					"password":   types.StringType,
				},
			},
		},
		Attrs: map[string]attr.Value{
			"boot_to_network": types.Bool{Value: false},
			"share_type":      types.String{Value: ""},
			"iso_timeout":     types.Int64{Value: 0},
			"iso_path":        types.String{Value: ""},
			"share_detail":    shareDetailsTfsdk,
		},
	}

	stateTemplateDeployment.BootToNetworkISO = bootToNetworkISOTfsdk

	//Save into State
	diags := resp.State.Set(ctx, &stateTemplateDeployment)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_deploy import: finished")
}

func updateDeploymentState(stateTemplateDeployment, planTemplateDeployment *models.TemplateDeployment, templateID int64, templateName string, omeClient *clients.Client, usedDeviceInput string) error {
	stateTemplateDeployment.ID.Value = strconv.FormatInt(templateID, 10)
	stateTemplateDeployment.TemplateID.Value = templateID
	stateTemplateDeployment.TemplateName.Value = templateName
	stateTemplateDeployment.JobRetryCount = planTemplateDeployment.JobRetryCount
	stateTemplateDeployment.SleepInterval = planTemplateDeployment.SleepInterval

	devIDsTfsdk := types.Set{
		ElemType: types.Int64Type,
	}
	devIDList := planTemplateDeployment.DeviceIDs.Elems

	devSTsTfsdk := types.Set{
		ElemType: types.StringType,
	}
	devSTList := planTemplateDeployment.DeviceServicetags.Elems
	profileDevSTVals := []attr.Value{}
	profileDevIDVals := []attr.Value{}
	serverProfiles, err := omeClient.GetServerProfileInfoByTemplateName(templateName)
	if err != nil {
		return err
	}
	for _, serverProfile := range serverProfiles.Value {
		device, _ := omeClient.GetDevice("", serverProfile.TargetID)
		deviceSTVal := types.String{Value: device.DeviceServiceTag}
		profileDevSTVals = append(profileDevSTVals, deviceSTVal)
		deviceIDVal := types.Int64{Value: serverProfile.TargetID}
		profileDevIDVals = append(profileDevIDVals, deviceIDVal)
	}

	filteredDevSTList := []attr.Value{}
	for _, planDevST := range devSTList {
		if tfsdkContains(profileDevSTVals, planDevST) {
			filteredDevSTList = append(filteredDevSTList, planDevST)
		}
	}

	filteredDevIDList := []attr.Value{}
	for _, planDevID := range devIDList {
		if tfsdkContains(profileDevIDVals, planDevID) {
			filteredDevIDList = append(filteredDevIDList, planDevID)
		}
	}

	for _, profileDevID := range profileDevIDVals {
		if !tfsdkContains(devIDList, profileDevID) {
			filteredDevIDList = append(filteredDevIDList, profileDevID)
		}
	}

	for _, profileDevST := range profileDevSTVals {
		if !tfsdkContains(devSTList, profileDevST) {
			filteredDevSTList = append(filteredDevSTList, profileDevST)
		}
	}

	switch usedDeviceInput {
	case clients.ServiceTags:
		devSTsTfsdk.Elems = filteredDevSTList
		stateTemplateDeployment.DeviceServicetags = devSTsTfsdk
		stateTemplateDeployment.DeviceIDs = planTemplateDeployment.DeviceIDs
	case clients.DeviceIDs:
		devIDsTfsdk.Elems = filteredDevIDList
		stateTemplateDeployment.DeviceIDs = devIDsTfsdk
		stateTemplateDeployment.DeviceServicetags = planTemplateDeployment.DeviceServicetags
	}

	stateTemplateDeployment.BootToNetworkISO = planTemplateDeployment.BootToNetworkISO
	stateTemplateDeployment.DeviceAttributes = planTemplateDeployment.DeviceAttributes
	stateTemplateDeployment.OptionsContinueOnWarning = planTemplateDeployment.OptionsContinueOnWarning
	stateTemplateDeployment.PowerStateOff = planTemplateDeployment.PowerStateOff
	stateTemplateDeployment.OptionsPrecheckOnly = planTemplateDeployment.OptionsPrecheckOnly
	stateTemplateDeployment.ForcedShutdown = planTemplateDeployment.ForcedShutdown
	stateTemplateDeployment.OptionsStrictCheckingVlan = planTemplateDeployment.OptionsStrictCheckingVlan
	stateTemplateDeployment.OptionsTimeToWaitBeforeShutdown = planTemplateDeployment.OptionsTimeToWaitBeforeShutdown
	stateTemplateDeployment.RunLater = planTemplateDeployment.RunLater
	stateTemplateDeployment.Cron = planTemplateDeployment.Cron
	return nil
}

func compare(comparing, comparedTo []int64) []int64 {
	compareToMap := make(map[int64]int64)
	for _, val := range comparedTo {
		compareToMap[val]++
	}

	var devIDs []int64
	for _, val := range comparing {
		if compareToMap[val] > 0 {
			compareToMap[val]--
			continue
		}
		devIDs = append(devIDs, val)
	}
	return devIDs
}

func tfsdkContains(toCheckList []attr.Value, toCheck attr.Value) bool {
	for _, val := range toCheckList {
		if reflect.DeepEqual(val, toCheck) {
			return true
		}
	}
	return false
}

func getSchedule(plan models.TemplateDeployment) models.OMESchedule {
	schedule := models.OMESchedule{
		RunNow:    false,
		RunLater:  true,
		Cron:      plan.Cron.Value,
		StartTime: "",
		EndTime:   "",
	}
	return schedule
}

func getOptions(plan models.TemplateDeployment) models.OMEOptions {
	options := models.OMEOptions{
		ShutdownType:             0,
		TimeToWaitBeforeShutdown: plan.OptionsTimeToWaitBeforeShutdown.Value,
		EndHostPowerState:        1,
		PrecheckOnly:             plan.OptionsPrecheckOnly.Value,
		ContinueOnWarning:        plan.OptionsContinueOnWarning.Value,
		StrictCheckingVLAN:       plan.OptionsStrictCheckingVlan.Value,
	}

	if plan.ForcedShutdown.Value {
		options.ShutdownType = 1
	}

	if plan.PowerStateOff.Value {
		options.EndHostPowerState = 0
	}
	return options
}

func getBootToNetworkISO(ctx context.Context, plan models.TemplateDeployment) (models.OMENetworkBootISOModel, []diag.Diagnostic, error) {
	bootToNetworkISO := models.BootToNetworkISO{}
	diags := plan.BootToNetworkISO.As(ctx, &bootToNetworkISO, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
	if diags.HasError() {
		return models.OMENetworkBootISOModel{}, diags, fmt.Errorf(clients.ErrUnableToParseBootToNetISO)
	}

	shareDetail := models.ShareDetail{}
	diags = bootToNetworkISO.ShareDetail.As(ctx, &shareDetail, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
	if diags.HasError() {
		return models.OMENetworkBootISOModel{}, diags, fmt.Errorf(clients.ErrUnableToParseBootToNetISO)
	}

	bootToNetworkISOModel := models.OMENetworkBootISOModel{
		BootToNetwork:  bootToNetworkISO.BootToNetwork.Value,
		ISOPath:        bootToNetworkISO.IsoPath.Value,
		ISOTimeout:     bootToNetworkISO.IsoTimeout.Value,
		ISOTimeoutUnit: 2,
		ShareType:      bootToNetworkISO.ShareType.Value,
		ShareDetail: models.OMEShareDetail{
			IPAddress: shareDetail.IPAddress.Value,
			ShareName: shareDetail.ShareName.Value,
			WorkGroup: shareDetail.WorkGroup.Value,
			User:      shareDetail.User.Value,
			Password:  shareDetail.Password.Value,
		},
	}
	return bootToNetworkISOModel, nil, nil
}

func getDeviceAttributes(ctx context.Context, devices []models.Device, plan models.TemplateDeployment) []models.OMEDeviceAttributes {
	omeDeviceAttributes := []models.OMEDeviceAttributes{}
	deviceAttributes := []models.DeviceAttributes{}

	deviceMap := map[string]int64{}
	for _, d := range devices {
		deviceMap[d.DeviceServiceTag] = d.ID
	}

	plan.DeviceAttributes.ElementsAs(ctx, &deviceAttributes, true)
	for _, deviceAttribute := range deviceAttributes {
		omeDeviceAttribute := models.OMEDeviceAttributes{}
		attributeList := []models.Attribute{}
		deviceServiceTags := []string{}
		deviceAttribute.Attributes.ElementsAs(ctx, &attributeList, true)
		deviceAttribute.DeviceServiceTags.ElementsAs(ctx, &deviceServiceTags, true)
		omeAttributes := []models.OMEAttribute{}
		for _, attribute := range attributeList {
			omeAttribute := models.OMEAttribute{
				ID:        attribute.AttributeID.Value,
				Value:     attribute.Value.Value,
				IsIgnored: attribute.IsIgnored.Value,
			}
			omeAttributes = append(omeAttributes, omeAttribute)
		}
		for _, deviceServiceTag := range deviceServiceTags {
			if val, ok := deviceMap[deviceServiceTag]; ok {
				omeDeviceAttribute.DeviceID = val
				omeDeviceAttribute.Attributes = omeAttributes
				omeDeviceAttributes = append(omeDeviceAttributes, omeDeviceAttribute)
			}
		}
	}
	return omeDeviceAttributes
}
