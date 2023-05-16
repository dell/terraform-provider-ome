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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceDeployment{}
	_ resource.ResourceWithConfigure   = &resourceDeployment{}
	_ resource.ResourceWithImportState = &resourceDeployment{}
)

// NewDeploymentResource is a new resource for deployment
func NewDeploymentResource() resource.Resource {
	return &resourceDeployment{}
}

type resourceDeployment struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceDeployment) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (resourceDeployment) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "deployment"
}

// Template Deployment Resource schema
func (r resourceDeployment) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing template deployment on OpenManage Enterprise.",
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the deploy resource.",
				Description:         "ID of the deploy resource.",
				Computed:            true,
			},
			"template_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the existing template." +
					" If a template with this ID is found, `template_name` will be ignored." +
					" Cannot be updated.",
				Description: "ID of the existing template." +
					" If a template with this ID is found, 'template_name' will be ignored." +
					" Cannot be updated.",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"template_name": schema.StringAttribute{
				MarkdownDescription: "Name of the existing template." +
					" Cannot be updated.",
				Description: "Name of the existing template." +
					" Cannot be updated.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "List of the device id(s)." +
					" Conflicts with `device_servicetags`.",
				Description: "List of the device id(s)." +
					" Conflicts with 'device_servicetags.'",
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"device_servicetags": schema.SetAttribute{
				MarkdownDescription: "List of the device servicetags." +
					" Conflicts with `device_ids`.",
				Description: "List of the device servicetags." +
					" Conflicts with 'device_ids.'",
				ElementType: types.StringType,
				Optional:    true,
			},
			"boot_to_network_iso": schema.ObjectAttribute{
				MarkdownDescription: "Boot To Network ISO deployment details.",
				Description:         "Boot To Network ISO deployment details.",
				Optional:            true,
				AttributeTypes: map[string]attr.Type{
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
			"forced_shutdown": schema.BoolAttribute{
				MarkdownDescription: "Force shutdown after deployment.",
				Description:         "Force shutdown after deployment.",
				Optional:            true,
			},
			"options_time_to_wait_before_shutdown": schema.Int64Attribute{
				MarkdownDescription: "Option to specify the time to wait before shutdown in seconds. Default and minimum value is 300 and maximum is 3600 seconds respectively." +
					" Default value is `300`.",
				Description: "Option to specify the time to wait before shutdown in seconds. Default and minimum value is 300 and maximum is 3600 seconds respectively." +
					" Default value is '300'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(300)),
				},
			},
			"power_state_off": schema.BoolAttribute{
				MarkdownDescription: "End power state of a target devices. Default power state is ON. Make it true to switch it to OFF state.",
				Description:         "End power state of a target devices. Default power state is ON. Make it true to switch it to OFF state.",
				Optional:            true,
			},
			"options_precheck_only": schema.BoolAttribute{
				MarkdownDescription: "Option to precheck",
				Description:         "Option to precheck",
				Optional:            true,
			},
			"options_strict_checking_vlan": schema.BoolAttribute{
				MarkdownDescription: "Checks the strict association of vlan.",
				Description:         "Checks the strict association of vlan.",
				Optional:            true,
			},
			"options_continue_on_warning": schema.BoolAttribute{
				MarkdownDescription: "Continue to run the job on warnings.",
				Description:         "Continue to run the job on warnings.",
				Optional:            true,
			},
			"run_later": schema.BoolAttribute{
				MarkdownDescription: "Provides options to schedule the deployment task immediately, or at a specified time.",
				Description:         "Provides options to schedule the deployment task immediately, or at a specified time.",
				Optional:            true,
			},
			"cron": schema.StringAttribute{
				MarkdownDescription: "Cron to schedule the deployment task. Cron expression should be of future datetime.",
				Description:         "Cron to schedule the deployment task. Cron expression should be of future datetime.",
				Optional:            true,
			},
			"device_attributes": schema.ListAttribute{
				MarkdownDescription: "List of template attributes associated with the target devices for deploymnent.",
				Description:         "List of template attributes associated with the target devices for deploymnent.",
				Optional:            true,
				ElementType: types.ObjectType{
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
			"job_retry_count": schema.Int64Attribute{
				MarkdownDescription: "Number of times the job has to be polled to get the final status of the resource." +
					" Default value is `20`.",
				Description: "Number of times the job has to be polled to get the final status of the resource." +
					" Default value is '20'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(20)),
				},
			},
			"sleep_interval": schema.Int64Attribute{
				MarkdownDescription: "Sleep time interval for job polling in seconds." +
					" Default value is `60`.",
				Description: "Sleep time interval for job polling in seconds." +
					" Default value is '60'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(60)),
				},
			},
		},
	}
}

// Create a new resource
func (r resourceDeployment) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_deploy create : Started")
	//Get Plan Data
	var plan models.TemplateDeployment
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get plan devices
	var serviceTags []string
	var devIDs []int64

	diags = plan.DeviceServicetags.ElementsAs(ctx, &serviceTags, true)
	resp.Diagnostics.Append(diags...)

	diags = plan.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	templateDeploymentState := models.TemplateDeployment{}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_deploy Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Trace(ctx, "resource_deploy create: session created")

	omeTemplate, err := omeClient.GetTemplateByIDOrName(plan.TemplateID.ValueInt64(), plan.TemplateName.ValueString())
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
	if !plan.BootToNetworkISO.IsNull() {
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
	if plan.RunLater.ValueBool() {
		deploymentRequest.Schedule = getSchedule(plan)
	}
	//Schedule ends
	// Device Attrs
	if len(plan.DeviceAttributes.Elements()) > 0 {
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

	if !plan.RunLater.ValueBool() {
		tflog.Trace(ctx, "resource_deploy create: started job tracking")
		isSuccess, message := omeClient.TrackJob(deploymentJobID, plan.JobRetryCount.ValueInt64(), plan.SleepInterval.ValueInt64())
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
	tflog.Trace(ctx, "resource_deploy create: finish")
}

// Read resource information
func (r resourceDeployment) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//Get State Data
	tflog.Trace(ctx, "resource_deploy read: started")
	var stateTemplateDeployment models.TemplateDeployment
	diags := req.State.Get(ctx, &stateTemplateDeployment)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := stateTemplateDeployment.TemplateID.ValueInt64()
	templateName := stateTemplateDeployment.TemplateName.ValueString()

	tflog.Debug(ctx, "resource_deploy read: reading a template", map[string]interface{}{
		"id":   templateID,
		"name": templateName,
	})
	var serviceTags []string
	var devIDs []int64

	diags = stateTemplateDeployment.DeviceServicetags.ElementsAs(ctx, &serviceTags, true)
	resp.Diagnostics.Append(diags...)

	diags = stateTemplateDeployment.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var usedDeviceInput string
	if len(serviceTags) > 0 {
		usedDeviceInput = clients.ServiceTags
	} else if len(devIDs) > 0 {
		usedDeviceInput = clients.DeviceIDs
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_deploy Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Trace(ctx, "resource_deploy read: client created started updating state")
	_ = updateDeploymentState(&stateTemplateDeployment, &stateTemplateDeployment, templateID, templateName, omeClient, usedDeviceInput)

	tflog.Trace(ctx, "resource_deploy read: finished reading state")
	//Save into State
	diags = resp.State.Set(ctx, &stateTemplateDeployment)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_deploy read: finished")
}

// Update resource
func (r resourceDeployment) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_deploy Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	if (plan.TemplateID.ValueInt64() != 0 && plan.TemplateID.ValueInt64() != state.TemplateID.ValueInt64()) || (plan.TemplateName.ValueString() != "" && plan.TemplateName.ValueString() != state.TemplateName.ValueString()) {
		resp.Diagnostics.AddError(
			clients.ErrTemplateDeploymentUpdate,
			clients.ErrTemplateChanges,
		)
		return
	}
	tflog.Debug(ctx, "resource_deploy update: started with template", map[string]interface{}{
		"id":   plan.TemplateID.ValueInt64(),
		"name": plan.TemplateName.ValueString(),
	})

	// get the plan device ids
	var serviceTags []string
	var devIDs []int64

	diags = plan.DeviceServicetags.ElementsAs(ctx, &serviceTags, true)
	resp.Diagnostics.Append(diags...)

	diags = plan.DeviceIDs.ElementsAs(ctx, &devIDs, true)
	resp.Diagnostics.Append(diags...)

	// get the state devicds's

	var stateServiceTags []string
	var stateDevIDs []int64

	diags = state.DeviceServicetags.ElementsAs(ctx, &stateServiceTags, true)
	resp.Diagnostics.Append(diags...)

	diags = state.DeviceIDs.ElementsAs(ctx, &stateDevIDs, true)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
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

	serverProfiles, err := omeClient.GetServerProfileInfoByTemplateName(state.TemplateName.ValueString())
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
		ID:        state.TemplateID.ValueInt64(),
		TargetIDS: newDeployDevIDs,
		Options:   options,
	}

	//Boot to Network ISO starts
	if !plan.BootToNetworkISO.IsNull() {
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
	if plan.RunLater.ValueBool() {
		deploymentRequest.Schedule = getSchedule(plan)
	}
	//Schedule ends
	// Device Attrs
	if len(plan.DeviceAttributes.Elements()) > 0 {
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

		if !plan.RunLater.ValueBool() {
			tflog.Trace(ctx, "resource_deploy update: started job tracking")
			isSuccess, message := omeClient.TrackJob(deploymentJobID, plan.JobRetryCount.ValueInt64(), plan.SleepInterval.ValueInt64())
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

	_ = updateDeploymentState(&state, &plan, state.TemplateID.ValueInt64(), state.TemplateName.ValueString(), omeClient, usedDeviceInput)

	tflog.Trace(ctx, "resource_deploy update: finished state update")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_deploy update: finished")
}

// Delete resource
func (r resourceDeployment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_deploy delete: started")
	// Get State Data
	var statetemplateDeployment models.TemplateDeployment
	diags := req.State.Get(ctx, &statetemplateDeployment)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_deploy Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Debug(ctx, "resource_deploy delete: started with template", map[string]interface{}{
		"id":   statetemplateDeployment.TemplateID.ValueInt64(),
		"name": statetemplateDeployment.TemplateName.ValueString(),
	})

	serverProfiles, err := omeClient.GetServerProfileInfoByTemplateName(statetemplateDeployment.TemplateName.ValueString())
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
func (r resourceDeployment) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "resource_deploy import: started")
	// Save the import identifier in the id attribute
	var stateTemplateDeployment models.TemplateDeployment
	templateName := req.ID

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_deploy ImportState")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	omeTemplate, err := omeClient.GetTemplateByName(templateName)
	if err != nil {
		resp.Diagnostics.AddError(clients.ErrImportDeployment, err.Error())
		return
	}
	templateID := omeTemplate.ID

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
		deviceSTVal := types.StringValue(device.DeviceServiceTag)
		profileDevSTVals = append(profileDevSTVals, deviceSTVal)
	}
	stateTemplateDeployment.ID = types.StringValue(strconv.FormatInt(templateID, 10))
	stateTemplateDeployment.TemplateID = types.Int64Value(templateID)
	stateTemplateDeployment.TemplateName = types.StringValue(templateName)
	devSTsTfsdk, _ := types.SetValue(
		types.StringType,
		profileDevSTVals,
	)
	if !devSTsTfsdk.IsUnknown() {
		stateTemplateDeployment.DeviceServicetags = devSTsTfsdk
	}
	devIDsTfsdk, _ := types.SetValue(types.Int64Type, []attr.Value{})
	if !devIDsTfsdk.IsUnknown() {
		stateTemplateDeployment.DeviceIDs = devIDsTfsdk
	}
	// set empty device attributes
	attributesObjects := []attr.Value{}
	attributesObject, _ := types.ObjectValue(
		map[string]attr.Type{
			"attribute_id": types.Int64Type,
			"display_name": types.StringType,
			"value":        types.StringType,
			"is_ignored":   types.BoolType,
		},
		map[string]attr.Value{
			"attribute_id": types.Int64Value(0),
			"display_name": types.StringValue(""),
			"value":        types.StringValue(""),
			"is_ignored":   types.BoolValue(false),
		},
	)
	attributesObjects = append(attributesObjects, attributesObject)

	attributesTfsdk, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_id": types.Int64Type,
				"display_name": types.StringType,
				"value":        types.StringType,
				"is_ignored":   types.BoolType,
			},
		},
		attributesObjects,
	)

	deviceAttributeObjects := []attr.Value{}
	deviceAttributeObject, _ := types.ObjectValue(
		map[string]attr.Type{
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
		map[string]attr.Value{
			"device_servicetags": devSTsTfsdk,
			"attributes":         attributesTfsdk,
		},
	)

	deviceAttributeObjects = append(deviceAttributeObjects, deviceAttributeObject)

	deviceAttributeTfsdk, _ := types.ListValue(
		types.ObjectType{
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
		deviceAttributeObjects,
	)

	if !deviceAttributeTfsdk.IsUnknown() {
		stateTemplateDeployment.DeviceAttributes = deviceAttributeTfsdk
	}
	shareDetailsTfsdk, _ := types.ObjectValue(
		map[string]attr.Type{
			"ip_address": types.StringType,
			"share_name": types.StringType,
			"work_group": types.StringType,
			"user":       types.StringType,
			"password":   types.StringType,
		},
		map[string]attr.Value{
			"ip_address": types.StringValue(""),
			"share_name": types.StringValue(""),
			"work_group": types.StringValue(""),
			"user":       types.StringValue(""),
			"password":   types.StringValue(""),
		},
	)

	bootToNetworkISOTfsdk, _ := types.ObjectValue(
		map[string]attr.Type{
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
		map[string]attr.Value{
			"boot_to_network": types.BoolValue(false),
			"share_type":      types.StringValue(""),
			"iso_timeout":     types.Int64Value(0),
			"iso_path":        types.StringValue(""),
			"share_detail":    shareDetailsTfsdk,
		},
	)

	if !bootToNetworkISOTfsdk.IsUnknown() {
		stateTemplateDeployment.BootToNetworkISO = bootToNetworkISOTfsdk
	}
	//Save into State
	diags := resp.State.Set(ctx, &stateTemplateDeployment)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_deploy import: finished")
}

func updateDeploymentState(stateTemplateDeployment, planTemplateDeployment *models.TemplateDeployment, templateID int64, templateName string, omeClient *clients.Client, usedDeviceInput string) error {
	stateTemplateDeployment.ID = types.StringValue(strconv.FormatInt(templateID, 10))
	stateTemplateDeployment.TemplateID = types.Int64Value(templateID)
	stateTemplateDeployment.TemplateName = types.StringValue(templateName)
	if !planTemplateDeployment.JobRetryCount.IsUnknown() {
		stateTemplateDeployment.JobRetryCount = planTemplateDeployment.JobRetryCount
	}
	if !planTemplateDeployment.SleepInterval.IsUnknown() {
		stateTemplateDeployment.SleepInterval = planTemplateDeployment.SleepInterval
	}
	devIDList := planTemplateDeployment.DeviceIDs.Elements()
	devSTList := planTemplateDeployment.DeviceServicetags.Elements()
	profileDevSTVals := []attr.Value{}
	profileDevIDVals := []attr.Value{}
	serverProfiles, err := omeClient.GetServerProfileInfoByTemplateName(templateName)
	if err != nil {
		return err
	}
	for _, serverProfile := range serverProfiles.Value {
		device, _ := omeClient.GetDevice("", serverProfile.TargetID)
		deviceSTVal := types.StringValue(device.DeviceServiceTag)
		profileDevSTVals = append(profileDevSTVals, deviceSTVal)
		deviceIDVal := types.Int64Value(serverProfile.TargetID)
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
		devSTsTfsdk, _ := types.SetValue(
			types.StringType,
			filteredDevSTList,
		)
		stateTemplateDeployment.DeviceServicetags = devSTsTfsdk
		stateTemplateDeployment.DeviceIDs = types.SetNull(types.Int64Type)
		//planTemplateDeployment.DeviceIDs
	case clients.DeviceIDs:
		devIDsTfsdk, _ := types.SetValue(
			types.Int64Type,
			filteredDevIDList,
		)
		stateTemplateDeployment.DeviceIDs = devIDsTfsdk
		stateTemplateDeployment.DeviceServicetags = types.SetNull(types.StringType)
	}
	if !planTemplateDeployment.BootToNetworkISO.IsUnknown() {
		stateTemplateDeployment.BootToNetworkISO = planTemplateDeployment.BootToNetworkISO
	}
	if !planTemplateDeployment.DeviceAttributes.IsUnknown() {
		stateTemplateDeployment.DeviceAttributes = planTemplateDeployment.DeviceAttributes
	}
	if !planTemplateDeployment.OptionsContinueOnWarning.IsUnknown() {
		stateTemplateDeployment.OptionsContinueOnWarning = planTemplateDeployment.OptionsContinueOnWarning
	}
	if !planTemplateDeployment.PowerStateOff.IsUnknown() {
		stateTemplateDeployment.PowerStateOff = planTemplateDeployment.PowerStateOff
	}
	if !planTemplateDeployment.OptionsPrecheckOnly.IsUnknown() {
		stateTemplateDeployment.OptionsPrecheckOnly = planTemplateDeployment.OptionsPrecheckOnly
	}
	if !planTemplateDeployment.ForcedShutdown.IsUnknown() {
		stateTemplateDeployment.ForcedShutdown = planTemplateDeployment.ForcedShutdown
	}
	if !planTemplateDeployment.OptionsStrictCheckingVlan.IsUnknown() {
		stateTemplateDeployment.OptionsStrictCheckingVlan = planTemplateDeployment.OptionsStrictCheckingVlan
	}
	if !planTemplateDeployment.OptionsTimeToWaitBeforeShutdown.IsUnknown() {
		stateTemplateDeployment.OptionsTimeToWaitBeforeShutdown = planTemplateDeployment.OptionsTimeToWaitBeforeShutdown
	}
	if !planTemplateDeployment.RunLater.IsUnknown() {
		stateTemplateDeployment.RunLater = planTemplateDeployment.RunLater
	}
	if !planTemplateDeployment.Cron.IsUnknown() {
		stateTemplateDeployment.Cron = planTemplateDeployment.Cron
	}
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
		Cron:      plan.Cron.ValueString(),
		StartTime: "",
		EndTime:   "",
	}
	return schedule
}

func getOptions(plan models.TemplateDeployment) models.OMEOptions {
	options := models.OMEOptions{
		ShutdownType:             0,
		TimeToWaitBeforeShutdown: plan.OptionsTimeToWaitBeforeShutdown.ValueInt64(),
		EndHostPowerState:        1,
		PrecheckOnly:             plan.OptionsPrecheckOnly.ValueBool(),
		ContinueOnWarning:        plan.OptionsContinueOnWarning.ValueBool(),
		StrictCheckingVLAN:       plan.OptionsStrictCheckingVlan.ValueBool(),
	}

	if plan.ForcedShutdown.ValueBool() {
		options.ShutdownType = 1
	}

	if plan.PowerStateOff.ValueBool() {
		options.EndHostPowerState = 0
	}
	return options
}

func getBootToNetworkISO(ctx context.Context, plan models.TemplateDeployment) (models.OMENetworkBootISOModel, []diag.Diagnostic, error) {
	bootToNetworkISO := models.BootToNetworkISO{}
	diags := plan.BootToNetworkISO.As(ctx, &bootToNetworkISO, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
	if diags.HasError() {
		return models.OMENetworkBootISOModel{}, diags, fmt.Errorf(clients.ErrUnableToParseBootToNetISO)
	}

	shareDetail := models.ShareDetail{}
	diags = bootToNetworkISO.ShareDetail.As(ctx, &shareDetail, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
	if diags.HasError() {
		return models.OMENetworkBootISOModel{}, diags, fmt.Errorf(clients.ErrUnableToParseBootToNetISO)
	}

	bootToNetworkISOModel := models.OMENetworkBootISOModel{
		BootToNetwork:  bootToNetworkISO.BootToNetwork.ValueBool(),
		ISOPath:        bootToNetworkISO.IsoPath.ValueString(),
		ISOTimeout:     bootToNetworkISO.IsoTimeout.ValueInt64(),
		ISOTimeoutUnit: 2,
		ShareType:      bootToNetworkISO.ShareType.ValueString(),
		ShareDetail: models.OMEShareDetail{
			IPAddress: shareDetail.IPAddress.ValueString(),
			ShareName: shareDetail.ShareName.ValueString(),
			WorkGroup: shareDetail.WorkGroup.ValueString(),
			User:      shareDetail.User.ValueString(),
			Password:  shareDetail.Password.ValueString(),
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
				ID:        attribute.AttributeID.ValueInt64(),
				Value:     attribute.Value.ValueString(),
				IsIgnored: attribute.IsIgnored.ValueBool(),
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
