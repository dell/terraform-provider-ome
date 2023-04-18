package ome

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	//ComplianceViewTypeID - stores the id for the compliance view type.
	ComplianceViewTypeID = 1
	//DeploymentViewTypeID - stores the id for the deployment view type.
	DeploymentViewTypeID = 2
	//ChassisDeviceTypeID - stores the id for the Chassis device type type.
	ChassisDeviceTypeID = 4
	// RetryCount - stores the default value of retry count
	RetryCount = 5
	// SleepInterval - stores the default value of sleep interval
	SleepInterval = 30
	// SleepTimeBeforeJob - wait time in seconds before job tracking
	SleepTimeBeforeJob = 5
	// NicPortDivider - specifies divider used between NIC identifier and port
	NicPortDivider = "/"
	// NicIdentifierAndPort - key for identifying unique NIC in a Vlan
	NicIdentifierAndPort = "%s" + NicPortDivider + "%d"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceTemplate{}
	_ resource.ResourceWithConfigure   = &resourceTemplate{}
	_ resource.ResourceWithImportState = &resourceTemplate{}
)

// NewTemplateResource is new resource for template
func NewTemplateResource() resource.Resource {
	return &resourceTemplate{}
}

type resourceTemplate struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceTemplate) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// panic("unimplemented")
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (*resourceTemplate) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "template"
}

// Order Resource schema
func (r *resourceTemplate) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing template on OpenManage Enterprise.Updates are supported for the following parameters: `name`, `description`, `attributes`, `job_retry_count`, `sleep_interval`, `identity_pool_name`, `vlan`.",
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the template resource.",
				Description:         "ID of the template resource.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the template resource.",
				Description:         "Name of the template resource.",
				Required:            true,
			},
			"fqdds": schema.StringAttribute{
				MarkdownDescription: "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters. This field cannot be updated.",
				Description:         "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters. This field cannot be updated.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					StringDefaultValue(types.StringValue("All")),
				},
				Validators: []validator.String{
					validFqddsValidator{},
				},
			},
			"view_type": schema.StringAttribute{
				MarkdownDescription: "OME template view type, supported types are Deployment, Compliance. This field cannot be updated.",
				Description:         "OME template view type, supported types are Deployment, Compliance. This field cannot be updated.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					StringDefaultValue(types.StringValue("Deployment")),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Deployment",
						"Compliance",
					),
				},
			},
			"view_type_id": schema.Int64Attribute{
				MarkdownDescription: "OME template view type id.",
				Description:         "OME template view type id.",
				Computed:            true,
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: "OME template device type, supported types are Server, Chassis. This field cannot be updated and is applicable only for importing xml.",
				Description:         "OME template device type, supported types are Server, Chassis. This field cannot be updated and is applicable only for importing xml.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					StringDefaultValue(types.StringValue("Server")),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Server",
						"Chassis",
					),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The XML content of template.. This field cannot be updated.",
				Description:         "The XML content of template.. This field cannot be updated.",
				Optional:            true,
				Computed:            true,
			},
			"refdevice_servicetag": schema.StringAttribute{
				MarkdownDescription: "Target device servicetag from which the template needs to be created. This field cannot be updated.",
				Description:         "Target device servicetag from which the template needs to be created. This field cannot be updated.",
				Optional:            true,
				Computed:            true,
			},
			"refdevice_id": schema.Int64Attribute{
				MarkdownDescription: "Target device id from which the template needs to be created. This field cannot be updated.",
				Description:         "Target device id from which the template needs to be created. This field cannot be updated.",
				Optional:            true,
				Computed:            true,
			},
			"reftemplate_name": schema.StringAttribute{
				MarkdownDescription: "Reference Template name from which the template needs to be cloned. This field cannot be updated.",
				Description:         "Reference Template name from which the template needs to be cloned. This field cannot be updated.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the template",
				Description:         "Description of the template",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attributes": schema.ListAttribute{
				MarkdownDescription: "List of attributes associated with a template. This field is ignored while creating a template.",
				Description:         "List of attributes associated with a template. This field is ignored while creating a template.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"attribute_id": types.Int64Type,
						"display_name": types.StringType,
						"value":        types.StringType,
						"is_ignored":   types.BoolType,
					},
				},
			},
			"job_retry_count": schema.Int64Attribute{
				MarkdownDescription: "Number of times the job has to be polled to get the final status of the resource.",
				Description:         "Number of times the job has to be polled to get the final status of the resource.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(RetryCount)),
				},
			},
			"sleep_interval": schema.Int64Attribute{
				MarkdownDescription: "Sleep time interval for job polling in seconds.",
				Description:         "Sleep time interval for job polling in seconds.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(SleepInterval)),
				},
			},
			"identity_pool_name": schema.StringAttribute{
				MarkdownDescription: "Identity Pool name to be attached with template.",
				Description:         "Identity Pool name to be attached with template.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_pool_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the Identity Pool attached with template.",
				Description:         "ID of the Identity Pool attached with template.",
				Computed:            true,
			},
			"vlan": schema.ObjectAttribute{
				MarkdownDescription: "VLAN details to be attached with template.",
				Description:         "VLAN details to be attached with template.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				AttributeTypes: map[string]attr.Type{
					"propogate_vlan":     types.BoolType,
					"bonding_technology": types.StringType,
					"vlan_attributes": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"untagged_network": types.Int64Type,
								"tagged_networks": types.SetType{
									ElemType: types.Int64Type,
								},
								"is_nic_bonded":  types.BoolType,
								"port":           types.Int64Type,
								"nic_identifier": types.StringType,
							},
						},
					},
				},
			},
		},
	}
}

// Create a new resource
func (r *resourceTemplate) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Read the data from Plan
	tflog.Trace(ctx, "resource_template create: started")
	var plan models.Template
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "resource_template create: reference data", map[string]interface{}{
		"refdeviceid":         plan.RefdeviceID.ValueInt64(),
		"refdeviceservicetag": plan.RefdeviceServicetag,
		"refTemplate":         plan.ReftemplateName,
	})

	template := models.Template{}

	err := validateCreate(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateTemplate, err.Error(),
		)
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_template Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	viewTypeID, err := omeClient.GetViewTypeID(plan.ViewType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateTemplate, err.Error(),
		)
		return
	}

	deviceTypeID, err := omeClient.GetDeviceTypeID(plan.DeviceType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateTemplate, err.Error(),
		)
		return
	}

	omeTemplateData := models.OMETemplate{}
	var templateID int64

	if plan.ReftemplateName.ValueString() != "" {
		tflog.Info(ctx, "resource_template create: creating a template from a reference template")

		if !plan.Description.IsUnknown() {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, "description cannot be modified while cloning from a reference template.",
			)
			return
		}

		// The identity pool and Vlans does not get cloned into the new template in OME.
		sourceTemplate, err := omeClient.GetTemplateByName(plan.ReftemplateName.ValueString())
		if err != nil || sourceTemplate.Name == "" {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, "Unable to clone the template because Source template does not exist.",
			)
			return
		}

		if sourceTemplate.ViewTypeID == ComplianceViewTypeID && viewTypeID == DeploymentViewTypeID {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, "cannot clone compliance template as deployment template.",
			)
			return
		}

		cloneTemplateRequest := models.OMECloneTemplate{
			SourceTemplateID: sourceTemplate.ID,
			NewTemplateName:  plan.Name.ValueString(),
			ViewTypeID:       viewTypeID,
		}
		templateID, err = omeClient.CloneTemplateByRefTemplateID(cloneTemplateRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}

		omeTemplateData, err = omeClient.GetTemplateByID(templateID)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}
	} else if plan.Content.ValueString() != "" { // template import
		tflog.Info(ctx, "resource_template create: creating a template from a xml content")

		if !plan.Description.IsUnknown() {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, "description is not supported for template import operation.",
			)
			return
		}

		importTemplateRequest := models.OMEImportTemplate{
			ViewTypeID: viewTypeID,
			Type:       deviceTypeID,
			Name:       plan.Name.ValueString(),
			Content:    plan.Content.ValueString(),
		}

		templateID, err = omeClient.ImportTemplate(importTemplateRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}

		omeTemplateData, err = omeClient.GetTemplateByID(templateID)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}

	} else {
		tflog.Info(ctx, "resource_template create: creating a template from a reference device")
		deviceID, err := omeClient.ValidateDevice(plan.RefdeviceServicetag.ValueString(), plan.RefdeviceID.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}

		ct := models.CreateTemplate{
			Fqdds:          strings.ReplaceAll(plan.FQDDS.ValueString(), " ", ""),
			ViewTypeID:     viewTypeID,
			SourceDeviceID: deviceID,
			Name:           plan.Name.ValueString(),
			Description:    plan.Description.ValueString(),
		}

		templateID, err = omeClient.CreateTemplate(ct)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("template created with id %d", templateID))
		time.Sleep(SleepTimeBeforeJob * time.Second)
		omeTemplateData, err = omeClient.GetTemplateByID(templateID)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}

		isSuccess, message := omeClient.TrackJob(omeTemplateData.TaskID, plan.JobRetryCount.ValueInt64(), plan.SleepInterval.ValueInt64())
		if !isSuccess {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, message,
			)
			_, err = omeClient.Delete(fmt.Sprintf(clients.TemplateAPI+"(%d)", templateID), nil, nil)
			if err != nil {
				resp.Diagnostics.AddError(
					clients.ErrCreateTemplate,
					err.Error(),
				)
				return
			}
			return
		}
	}

	tflog.Trace(ctx, "resource_template create: fetching template attributes")

	omeAttributes, err := omeClient.GetTemplateAttributes(omeTemplateData.ID, []models.Attribute{}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateTemplate,
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_template create: fetching template valn data")
	omeVlan, err := omeClient.GetSchemaVlanData(templateID)
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate,
				err.Error(),
			)
			return
		}
	}
	//PropogateVlan is default true and is not available in the response from OME,hence setting it to true here to persist in state
	omeVlan.PropagateVLAN = true
	if !plan.RefdeviceServicetag.IsUnknown() {
		template.RefdeviceServicetag = plan.RefdeviceServicetag
	}

	if !plan.ReftemplateName.IsUnknown() {
		template.ReftemplateName = plan.ReftemplateName
	}

	if !plan.RefdeviceID.IsUnknown() {
		template.RefdeviceID = plan.RefdeviceID
	}

	if plan.ReftemplateName.ValueString() != "" {
		template.RefdeviceID = types.Int64Value(omeTemplateData.SourceDeviceID)
	}

	if !plan.Content.IsUnknown() {
		template.Content = plan.Content
	}
	if !plan.ViewType.IsUnknown() {
		template.ViewType = plan.ViewType
	}
	if !plan.DeviceType.IsUnknown() {
		template.DeviceType = plan.DeviceType
	}
	if !plan.FQDDS.IsUnknown() {
		template.FQDDS = plan.FQDDS
	} // The default value of fqdds is set to `All`. So if the config doesn't have any value specified, the default value in the plan is `All`.
	if !plan.JobRetryCount.IsUnknown() {
		template.JobRetryCount = plan.JobRetryCount
	}
	if !plan.SleepInterval.IsUnknown() {
		template.SleepInterval = plan.SleepInterval
	}
	if !plan.IdentityPoolName.IsUnknown() {
		template.IdentityPoolName = plan.IdentityPoolName
	}

	tflog.Trace(ctx, "resource_template create: started updating state")

	updateState(&template, []models.VlanAttributes{}, &omeTemplateData, omeAttributes, omeVlan)

	tflog.Trace(ctx, "resource_template create: finished updating state")

	//Save into State if template creation is successful
	diags = resp.State.Set(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_template create: finished")

}

// Read resource information
func (r *resourceTemplate) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_template read: started")
	var template models.Template
	diags := req.State.Get(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID, _ := strconv.ParseInt(template.ID.ValueString(), 10, 64)

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_template Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	stateAttributes := []models.Attribute{}
	stateAttributeObjects := []types.Object{}
	template.Attributes.ElementsAs(ctx, &stateAttributeObjects, true)

	for _, stateAttrObject := range stateAttributeObjects {
		stateAttribute := models.Attribute{}
		stateAttrObject.As(ctx, &stateAttribute, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
		stateAttributes = append(stateAttributes, stateAttribute)
	}

	tflog.Debug(ctx, "resource_template read: Template id", map[string]interface{}{
		"templateid": templateID,
	})

	omeTemplateData, err := omeClient.GetTemplateByID(templateID)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrReadTemplate, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_template read: fetching template attributes")

	omeAttributes, err := omeClient.GetTemplateAttributes(templateID, stateAttributes, true)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrReadTemplate,
			fmt.Sprintf("Unable to refresh template attributes: %s", err.Error()),
		)
		return
	}

	stateVlan := models.Vlan{}
	diags = template.Vlan.As(ctx, &stateVlan, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
	if diags.HasError() {
		resp.Diagnostics.AddError(
			clients.ErrReadTemplate,
			"Unable to fetch Vlan from state. Hence, Cannot refresh the template resource",
		)
		return
	}
	tflog.Trace(ctx, "resource_template read: fetching template vlan data")

	omeVlan, err := omeClient.GetSchemaVlanData(templateID)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrReadTemplate,
			fmt.Sprintf("Unable to refresh vlan attributes: %s", err.Error()),
		)
		return
	}

	omeVlan.PropagateVLAN = stateVlan.PropogateVlan.ValueBool()

	if omeTemplateData.IdentityPoolID != 0 {
		identityPool, err := omeClient.GetIdentityPoolByID(omeTemplateData.IdentityPoolID)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrReadTemplate,
				fmt.Sprintf("Unable to fetch the identity pools: %s", err.Error()),
			)
			return
		}
		template.IdentityPoolName = types.StringValue(identityPool.Name)
	}

	tflog.Trace(ctx, "resource_template read: updating state started")

	vlanAttrs := []models.VlanAttributes{}
	stateVlan.VlanAttributes.ElementsAs(ctx, &vlanAttrs, true)

	updateState(&template, vlanAttrs, &omeTemplateData, omeAttributes, omeVlan)

	tflog.Trace(ctx, "resource_template read: updating state finished")

	diags = resp.State.Set(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_template read: finished")
}

// Update resource
func (r resourceTemplate) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "resource_template update: started")
	var planTemplate models.Template
	planDiags := req.Plan.Get(ctx, &planTemplate)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateTemplate models.Template
	stateDiags := req.State.Get(ctx, &stateTemplate)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID, _ := strconv.ParseInt(stateTemplate.ID.ValueString(), 10, 64)

	if isConfigValuesChanged(planTemplate, stateTemplate) {
		resp.Diagnostics.AddError(
			clients.ErrUpdateTemplate,
			"cannot update the following fields : `refdevice_servicetag`,`refdevice_id`,`view_type`, `reftemplate_name`, `content` and `fqdds`",
		)
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_template Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Debug(ctx, "resource_template update: Template id", map[string]interface{}{
		"templateid": templateID,
	})
	var (
		identityPool models.IdentityPool
		err          error
	)
	if planTemplate.IdentityPoolName.ValueString() != "" {
		identityPool, err = validateIOPoolName(omeClient, planTemplate.IdentityPoolName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrUpdateTemplate,
				err.Error(),
			)
			return
		}
	}
	planVlan := getVlanForTemplate(ctx, resp, planTemplate)

	if !planTemplate.Vlan.IsUnknown() {
		err := validateVlanNetworkData(omeClient, templateID, planVlan)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrUpdateTemplate,
				err.Error(),
			)
			return
		}

	}

	updatePayload := models.UpdateTemplate{
		Name:        planTemplate.Name.ValueString(),
		ID:          templateID,
		Description: stateTemplate.Description.ValueString(),
	}

	if planTemplate.Description != stateTemplate.Description {
		updatePayload.Description = planTemplate.Description.ValueString()
	}

	stateAttributes := getTfsdkStateAttributes(ctx, stateTemplate)
	// Terraform compares the list elements based on order, hence it is expected that the practitioner gives all attributes
	// along with the attribute for which modification is expected.
	da, _ := getDeltaAttributes(ctx, planTemplate, stateAttributes)

	tflog.Trace(ctx, "resource_template update: finished fetching delta attributes")

	if len(da) != 0 {
		tflog.Info(ctx, "resource_template update: delta attributes exists")
		updatePayload.Attributes = da
	}

	tflog.Trace(ctx, "resource_template update: started a call to update template")
	err = omeClient.UpdateTemplate(updatePayload)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrUpdateTemplate,
			err.Error(),
		)
		return
	}

	// Updating networkConfig now
	var updateIdentityPool = planTemplate.IdentityPoolName != stateTemplate.IdentityPoolName
	updateVlan := !reflect.DeepEqual(planTemplate.Vlan.Attributes(), stateTemplate.Vlan.Attributes())
	updateNetworkConfig := updateIdentityPool || updateVlan

	if updateNetworkConfig {
		tflog.Info(ctx, "resource_template update: updating network config")
		nwConfig := &models.UpdateNetworkConfig{
			TemplateID: templateID,
		}
		if updateIdentityPool {
			tflog.Info(ctx, "resource_template update: updating identity pool")
			// When IdentityPool is assigned to a  template for the first time, OME modifies the is_ignored parameter
			// of the below parameters from True to False,
			// ['IOIDOpt 1 Initiator Persistence Policy', 'IOIDOpt 1 Storage Target Persistence Policy',
			// 'OIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered',
			// 'IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered']. This will cause an inconsistency between plan and state.
			// Hence, before IO pool is assigned, these attributes will have to be modified.
			if planTemplate.IdentityPoolName.ValueString() != "" {
				nwConfig.IdentityPoolID = identityPool.ID
			} else {
				nwConfig.IdentityPoolID = 0
			}
		} else {
			nwConfig.IdentityPoolID = stateTemplate.IdentityPoolID.ValueInt64()
		}
		if planTemplate.IdentityPoolName.ValueString() != "" {
			identityPool, err := omeClient.GetIdentityPoolByName(planTemplate.IdentityPoolName.ValueString())
			if err != nil {
				resp.Diagnostics.AddWarning(
					fmt.Sprintf("Unable to update IdentityPool parameters to the template: %d", templateID),
					err.Error(),
				)
			} else {
				nwConfig.IdentityPoolID = identityPool.ID
			}
		} else {
			nwConfig.IdentityPoolID = 0
		}
		if updateVlan {
			tflog.Info(ctx, "resource_template update: updating vlan attrs")
			vlanNetworkView, err := omeClient.GetVlanNetworkModel(templateID)
			if err != nil {
				resp.Diagnostics.AddWarning(
					clients.ErrUpdateTemplate,
					fmt.Sprintf("unable to fetch network view from OME to update vlan to the template: %d, Error: %s", templateID, err.Error()),
				)
			} else {
				tflog.Info(ctx, "resource_template update: fetching vlan attrs")
				stateVlan := getVlanForTemplate(ctx, resp, stateTemplate)
				// when tagged networks are to be removed, its expected that the practitioner would give input as [0],
				// Provider will convert this to [] and send to API. [0] will be written back to state file to make plan
				//and state consistent.
				tflog.Info(ctx, "resource_template update: updating vlan attrs has no errors")
				deltaVlan := getDeltaVlan(ctx, planVlan, stateVlan)
				nwConfig.BondingTechnology = deltaVlan.BondingTechnology
				nwConfig.PropagateVLAN = deltaVlan.PropagateVLAN
				payloadVlanAttributes := []models.PayloadVlanAttribute{}
				for _, deltaVlanAttr := range deltaVlan.OMEVlanAttributes {
					payloadVlanAttr, _ := omeClient.GetPayloadVlanAttribute(vlanNetworkView, deltaVlanAttr.NicIdentifier, deltaVlanAttr.Port)
					payloadVlanAttr.Tagged = deltaVlanAttr.Tagged
					payloadVlanAttr.Untagged = deltaVlanAttr.Untagged
					payloadVlanAttr.IsNICBonded = deltaVlanAttr.IsNICBonded
					payloadVlanAttributes = append(payloadVlanAttributes, payloadVlanAttr)
				}
				nwConfig.VLANAttributes = payloadVlanAttributes
			}

		}
		tflog.Info(ctx, "resource_template update: triggering update netowrk config")
		err = omeClient.UpdateNetworkConfig(nwConfig)
		if err != nil {
			resp.Diagnostics.AddWarning(
				clients.ErrUpdateTemplate,
				fmt.Sprintf("unable to update network configuration to the template: %d, Error: %s", templateID, err.Error()),
			)
		}

	}

	tflog.Trace(ctx, "resource_template update: fetching template by id")
	omeTemplateData, err := omeClient.GetTemplateByID(templateID)

	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrUpdateTemplate, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_template update: fetching template attributes")

	omeAttributes, err := omeClient.GetTemplateAttributes(templateID, stateAttributes, true)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrUpdateTemplate,
			fmt.Sprintf("unable to refresh template attributes: %s", err.Error()),
		)
		return
	}

	tflog.Trace(ctx, "resource_template update: fetching vlan data")

	updatedVlan, err := omeClient.GetSchemaVlanData(templateID)
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrUpdateTemplate,
				fmt.Sprintf("unable to refresh vlan attributes: %s", err.Error()),
			)
			return
		}
	}

	updatedVlan.PropagateVLAN = planVlan.PropagateVLAN

	if !planTemplate.IdentityPoolName.IsUnknown() {
		stateTemplate.IdentityPoolName = planTemplate.IdentityPoolName
	}

	if !planTemplate.SleepInterval.IsUnknown() {
		stateTemplate.SleepInterval = planTemplate.SleepInterval
	}

	if !planTemplate.ViewType.IsUnknown() {
		stateTemplate.ViewType = planTemplate.ViewType
	}

	if !planTemplate.FQDDS.IsUnknown() {
		stateTemplate.FQDDS = planTemplate.FQDDS
	}

	if !planTemplate.JobRetryCount.IsUnknown() {
		stateTemplate.JobRetryCount = planTemplate.JobRetryCount
	}

	tflog.Trace(ctx, "resource_template update: updating state data started")

	tfsdkVlan := models.Vlan{}
	planTemplate.Vlan.As(ctx, &tfsdkVlan, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})

	vlanAttrs := []models.VlanAttributes{}
	tfsdkVlan.VlanAttributes.ElementsAs(ctx, &vlanAttrs, true)

	updateState(&stateTemplate, vlanAttrs, &omeTemplateData, omeAttributes, updatedVlan)

	tflog.Trace(ctx, "resource_template update: updating state data finished")
	//Save into State if template update is successful
	diags := resp.State.Set(ctx, &stateTemplate)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_template update: finished")
}

// Delete resource
func (r resourceTemplate) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_template delete: started")
	var template models.Template
	diags := resp.State.Get(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_template Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Trace(ctx, "resource_template delete: started delete")
	tflog.Debug(ctx, "resource_template delete: started delete for template", map[string]interface{}{
		"templateid": template.ID.ValueString(),
	})

	_, err := omeClient.Delete(fmt.Sprintf(clients.TemplateAPI+"(%s)", template.ID.ValueString()), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrDeleteTemplate,
			err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "resource_template delete: finished")
}

// Import resource
func (r resourceTemplate) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "resource_template import: started")
	var template models.Template
	template.Name = types.StringValue(req.ID)

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_template Import")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	omeTemplateData, err := omeClient.GetTemplateByName(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrImportTemplate, err.Error(),
		)
		return
	}

	if omeTemplateData.ID == 0 {
		resp.Diagnostics.AddError(
			clients.ErrImportTemplate, "invalid template name",
		)
		return
	}

	omeAttributes, err := omeClient.GetTemplateAttributes(omeTemplateData.ID, []models.Attribute{}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrImportTemplate,
			fmt.Sprintf("unable to get template attributes: %s", err.Error()),
		)
		return
	}

	omeVlan, err := omeClient.GetSchemaVlanData(omeTemplateData.ID)
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrImportTemplate,
				fmt.Sprintf("Unable to refresh vlan attributes: %s", err.Error()),
			)
			return
		}
	}
	tflog.Trace(ctx, "resource_template import: started state update")

	updateState(&template, []models.VlanAttributes{}, &omeTemplateData, omeAttributes, omeVlan)
	tflog.Trace(ctx, "resource_template import: finished state update")

	viewType := "Deployment"
	if omeTemplateData.ViewTypeID == ComplianceViewTypeID {
		viewType = "Compliance"
	}

	deviceType := "Server"
	if omeTemplateData.TypeID == ChassisDeviceTypeID {
		deviceType = "Chassis"
	}

	template.RefdeviceID = types.Int64Value(omeTemplateData.SourceDeviceID)
	template.RefdeviceServicetag = types.StringValue("NA")
	template.ReftemplateName = types.StringValue("NA")
	template.Content = types.StringValue("NA")
	template.ViewType = types.StringValue(viewType)
	template.DeviceType = types.StringValue(deviceType)
	template.JobRetryCount = types.Int64Value(RetryCount)
	template.SleepInterval = types.Int64Value(SleepInterval)
	template.FQDDS = types.StringValue("All")
	diags := resp.State.Set(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_template import: finished")
}

func validateCreate(plan models.Template) error {
	// all references cannot be empty
	if plan.ReftemplateName.ValueString() == "" && plan.RefdeviceID.ValueInt64() == 0 && plan.RefdeviceServicetag.ValueString() == "" && plan.Content.ValueString() == "" {
		return fmt.Errorf("either reftemplate_name or refdevice_id or refdevice_servicetag or content is required")
	}

	// any two references given results in error
	if (plan.ReftemplateName.ValueString() != "" && (plan.RefdeviceID.ValueInt64() != 0 || plan.RefdeviceServicetag.ValueString() != "" || plan.Content.ValueString() != "")) ||
		(plan.RefdeviceID.ValueInt64() != 0 && (plan.ReftemplateName.ValueString() != "" || plan.RefdeviceServicetag.ValueString() != "" || plan.Content.ValueString() != "")) ||
		(plan.RefdeviceServicetag.ValueString() != "" && (plan.RefdeviceID.ValueInt64() != 0 || plan.ReftemplateName.ValueString() != "" || plan.Content.ValueString() != "")) ||
		(plan.Content.ValueString() != "" && (plan.RefdeviceID.ValueInt64() != 0 || plan.ReftemplateName.ValueString() != "" || plan.RefdeviceServicetag.ValueString() != "")) {
		return fmt.Errorf("either reftemplate_name or refdevice_id or refdevice_servicetag or content is required")
	}

	// Identity Pool name and Vlan is supported only during update
	if !plan.IdentityPoolName.IsUnknown() || !plan.Vlan.IsUnknown() {
		return fmt.Errorf("attributes identity_pool_name and vlan cannot be associated during create")
	}

	// Attributes is part of plan during create
	if !plan.Attributes.IsUnknown() {
		return fmt.Errorf("attributes cannot be modified during create")
	}
	return nil
}

func validateVlanNetworkData(omeClient *clients.Client, templateID int64, planVlan models.OMEVlan) error {
	remoteVlanFromTemplate, err := omeClient.GetSchemaVlanData(templateID)
	if err != nil {
		return err
	}

	vlanNetworks, err := omeClient.GetAllVlanNetworks()
	if err != nil {
		return err
	}

	err = validateVlan(planVlan, remoteVlanFromTemplate, vlanNetworks)
	if err != nil {
		return err
	}

	return nil
}

func validateIOPoolName(omeClient *clients.Client, name string) (models.IdentityPool, error) {
	identityPool, err := omeClient.GetIdentityPoolByName(name)
	if err != nil {
		return identityPool, err
	}
	return identityPool, nil
}

func isConfigValuesChanged(planTemplate, stateTemplate models.Template) bool {
	return (!planTemplate.RefdeviceID.IsUnknown() && stateTemplate.RefdeviceID.ValueInt64() != planTemplate.RefdeviceID.ValueInt64()) ||
		(!planTemplate.RefdeviceServicetag.IsUnknown() && stateTemplate.RefdeviceServicetag.ValueString() != planTemplate.RefdeviceServicetag.ValueString()) ||
		(!planTemplate.ViewType.IsUnknown() && stateTemplate.ViewType.ValueString() != planTemplate.ViewType.ValueString()) ||
		(!planTemplate.FQDDS.IsUnknown() && stateTemplate.FQDDS.ValueString() != planTemplate.FQDDS.ValueString()) ||
		(!planTemplate.ReftemplateName.IsUnknown() && stateTemplate.ReftemplateName.ValueString() != planTemplate.ReftemplateName.ValueString()) ||
		(!planTemplate.Content.IsUnknown() && stateTemplate.Content.ValueString() != planTemplate.Content.ValueString())
}

func validateVlan(planVlan, remoteVlan models.OMEVlan, vlanNetworks []models.VLanNetworks) error {
	if len(remoteVlan.OMEVlanAttributes) == 0 {
		return fmt.Errorf("vlan attributes are not available in the template")
	}
	if len(planVlan.OMEVlanAttributes) != len(remoteVlan.OMEVlanAttributes) {
		return fmt.Errorf("number of port and nic identifier is inconsistent with the template")
	}
	remoteVlanIdentifiers := make(map[string]bool)
	var invalidNetworkIDs []int64
	for _, vlanAttr := range remoteVlan.OMEVlanAttributes {
		key := fmt.Sprintf(NicIdentifierAndPort, vlanAttr.NicIdentifier, vlanAttr.Port)
		remoteVlanIdentifiers[key] = true
	}
	var remoteVlanNetworkIDs = map[int64]bool{0: true}
	for _, vn := range vlanNetworks {
		remoteVlanNetworkIDs[vn.ID] = true
	}

	dupTagNetworkIDs := map[string][]int64{}
	networkKeys := []string{}

	for _, planVlanAttr := range planVlan.OMEVlanAttributes {
		key := fmt.Sprintf(NicIdentifierAndPort, planVlanAttr.NicIdentifier, planVlanAttr.Port)
		if _, ok := remoteVlanIdentifiers[key]; !ok {
			return fmt.Errorf("invalid combination of Nic Identifier and Port %s", key)
		}
		networkKeys = append(networkKeys, key)
		untagged := planVlanAttr.Untagged
		if !isValidNetworkID(remoteVlanNetworkIDs, untagged) {
			invalidNetworkIDs = append(invalidNetworkIDs, untagged)
		}
		taggedNetworkMap := map[int64]bool{}
		dupNetworkTag := []int64{}

		for _, tag := range planVlanAttr.Tagged {
			if _, ok := taggedNetworkMap[tag]; ok {
				dupNetworkTag = append(dupNetworkTag, tag)
			}
			taggedNetworkMap[tag] = true
			if !isValidNetworkID(remoteVlanNetworkIDs, tag) {
				invalidNetworkIDs = append(invalidNetworkIDs, tag)
			}
		}
		if len(dupNetworkTag) != 0 {
			fmtKey := strings.Replace(key, NicPortDivider, ", Port: ", 1)
			dupTagNetworkIDs[fmtKey] = dupNetworkTag
		}
	}

	isDuplicate := hasDuplicates(networkKeys)
	if len(networkKeys) != len(isDuplicate) {
		return fmt.Errorf("duplicate combination of Nic Identifier/Port %v ", isDuplicate)
	}

	if len(dupTagNetworkIDs) != 0 {
		return fmt.Errorf("duplicate vlan network IDs %v ", dupTagNetworkIDs)
	}

	if len(invalidNetworkIDs) != 0 {
		return fmt.Errorf("invalid vlan network IDs %v ", unique(invalidNetworkIDs))
	}
	return nil
}

func unique(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func hasDuplicates(strArr []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueNames []string
	for _, str := range strArr {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = true
			uniqueNames = append(uniqueNames, str)
		}
	}
	return uniqueNames
}

func isValidNetworkID(remoteVlanIds map[int64]bool, vlanID int64) bool {
	if _, ok := remoteVlanIds[vlanID]; !ok {
		return false
	}
	return true
}

func getVlanForTemplate(ctx context.Context, resp *resource.UpdateResponse, Template models.Template) models.OMEVlan {
	omeVlan := models.OMEVlan{}
	vlan := models.Vlan{}

	Template.Vlan.As(ctx, &vlan, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
	omeVlan.BondingTechnology = vlan.BondingTechnology.ValueString()
	omeVlan.PropagateVLAN = vlan.PropogateVlan.ValueBool()
	omeVlanAttrs := []models.OMEVlanAttribute{}
	vlanAttrs := []models.VlanAttributes{}
	vlan.VlanAttributes.ElementsAs(ctx, &vlanAttrs, true)
	for _, vlanAttr := range vlanAttrs {
		taggedNetworks := []int64{}
		vlanAttr.TaggedNetworks.ElementsAs(ctx, &taggedNetworks, true)
		omeVlanAttr := models.OMEVlanAttribute{
			Untagged:      vlanAttr.UntaggedNetwork.ValueInt64(),
			Tagged:        taggedNetworks,
			IsNICBonded:   vlanAttr.IsNicBonded.ValueBool(),
			Port:          vlanAttr.Port.ValueInt64(),
			NicIdentifier: vlanAttr.NicIdentifier.ValueString(),
		}
		omeVlanAttrs = append(omeVlanAttrs, omeVlanAttr)
	}
	omeVlan.OMEVlanAttributes = omeVlanAttrs
	return omeVlan
}

func getTfsdkStateAttributes(ctx context.Context, stateTemplate models.Template) []models.Attribute {
	stateAttributes := []models.Attribute{}
	stateAttributeObjects := []types.Object{}
	stateTemplate.Attributes.ElementsAs(ctx, &stateAttributeObjects, true)

	for _, stateAttrObject := range stateAttributeObjects {
		stateAttribute := models.Attribute{}
		stateAttrObject.As(ctx, &stateAttribute, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
		stateAttributes = append(stateAttributes, stateAttribute)
	}
	return stateAttributes
}

func getDeltaVlan(ctx context.Context, planVlan, stateVlan models.OMEVlan) models.OMEVlan {
	deltaVlan := models.OMEVlan{}
	deltaVlan.BondingTechnology = planVlan.BondingTechnology
	deltaVlan.PropagateVLAN = planVlan.PropagateVLAN
	deltaVlanAttributes := []models.OMEVlanAttribute{}
	for index, vlanAttr := range planVlan.OMEVlanAttributes {
		if !reflect.DeepEqual(vlanAttr, stateVlan.OMEVlanAttributes[index]) {
			if len(vlanAttr.Tagged) == 1 && vlanAttr.Tagged[0] == 0 {
				vlanAttr.Tagged = []int64{}
			}
			deltaVlanAttributes = append(deltaVlanAttributes, vlanAttr)
		}

	}
	deltaVlan.OMEVlanAttributes = deltaVlanAttributes
	return deltaVlan
}

func getDeltaAttributes(ctx context.Context, planTemplate models.Template, stateAttributes []models.Attribute) ([]models.UpdateAttribute, error) {
	updatedAttributes := []models.UpdateAttribute{}

	planUpdateAttributes := []models.Attribute{}
	planUpdateAttributeObjects := []types.Object{}
	planTemplate.Attributes.ElementsAs(ctx, &planUpdateAttributeObjects, true)

	for _, planUpdateAttrObject := range planUpdateAttributeObjects {
		planUpdateAttribute := models.Attribute{}
		planUpdateAttrObject.As(ctx, &planUpdateAttribute, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
		planUpdateAttributes = append(planUpdateAttributes, planUpdateAttribute)
	}
	if !reflect.DeepEqual(planUpdateAttributes, stateAttributes) {
		for index, attribute := range planUpdateAttributes {
			if attribute.Value != stateAttributes[index].Value {
				updateAttribute := models.UpdateAttribute{
					ID:        attribute.AttributeID.ValueInt64(),
					IsIgnored: attribute.IsIgnored.ValueBool(),
					Value:     attribute.Value.ValueString(),
				}
				updatedAttributes = append(updatedAttributes, updateAttribute)
			}
		}
	}

	return updatedAttributes, nil
}

func updateState(stateTemplate *models.Template, planVlanAttributes []models.VlanAttributes, omeTemplateData *models.OMETemplate, omeTemplateAttributes []models.OmeAttribute, omeVlan models.OMEVlan) {

	stateTemplate.ID = types.StringValue(fmt.Sprintf("%d", omeTemplateData.ID))
	stateTemplate.Name = types.StringValue(omeTemplateData.Name)
	stateTemplate.Description = types.StringValue(omeTemplateData.Description)
	stateTemplate.ViewTypeID = types.Int64Value(omeTemplateData.ViewTypeID)
	stateTemplate.IdentityPoolID = types.Int64Value(omeTemplateData.IdentityPoolID)

	attributeObjects := []attr.Value{}

	for _, attribute := range omeTemplateAttributes {
		attributeDetails := map[string]attr.Value{}
		attributeDetails["attribute_id"] = types.Int64Value(attribute.AttributeID)
		attributeDetails["display_name"] = types.StringValue(attribute.DisplayName)
		attributeDetails["value"] = types.StringValue(attribute.Value)
		attributeDetails["is_ignored"] = types.BoolValue(attribute.IsIgnored)
		attributeObject, _ := types.ObjectValue(
			map[string]attr.Type{
				"attribute_id": types.Int64Type,
				"display_name": types.StringType,
				"value":        types.StringType,
				"is_ignored":   types.BoolType,
			}, attributeDetails)
		attributeObjects = append(attributeObjects, attributeObject)
	}
	attributesTfsdk, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_id": types.Int64Type,
				"display_name": types.StringType,
				"value":        types.StringType,
				"is_ignored":   types.BoolType,
			},
		}, attributeObjects)
	stateTemplate.Attributes = attributesTfsdk

	omeVlanMap := map[string]models.OMEVlanAttribute{}

	for _, vlanAttr := range omeVlan.OMEVlanAttributes {
		key := fmt.Sprintf(NicIdentifierAndPort, vlanAttr.NicIdentifier, vlanAttr.Port)
		omeVlanMap[key] = vlanAttr
	}

	vlanAttrsObjects := []attr.Value{}

	for _, planVlanAttr := range planVlanAttributes {
		key := fmt.Sprintf(NicIdentifierAndPort, planVlanAttr.NicIdentifier.ValueString(), planVlanAttr.Port.ValueInt64())
		if omeVlanAttr, ok := omeVlanMap[key]; ok {
			vlanAttrObject := getVlanAtrrObject(omeVlanAttr)
			vlanAttrsObjects = append(vlanAttrsObjects, vlanAttrObject)

			delete(omeVlanMap, key)
		}
	}

	if len(omeVlanMap) != 0 {
		for _, omeVlanAttr := range omeVlanMap {
			vlanAttrObject := getVlanAtrrObject(omeVlanAttr)
			vlanAttrsObjects = append(vlanAttrsObjects, vlanAttrObject)
		}
	}

	vlanAttrList, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"untagged_network": types.Int64Type,
				"tagged_networks": types.SetType{
					ElemType: types.Int64Type,
				},
				"is_nic_bonded":  types.BoolType,
				"port":           types.Int64Type,
				"nic_identifier": types.StringType,
			},
		}, vlanAttrsObjects)

	vlanTfsdk, _ := types.ObjectValue(
		map[string]attr.Type{
			"propogate_vlan":     types.BoolType,
			"bonding_technology": types.StringType,
			"vlan_attributes": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"untagged_network": types.Int64Type,
						"tagged_networks": types.SetType{
							ElemType: types.Int64Type,
						},
						"is_nic_bonded":  types.BoolType,
						"port":           types.Int64Type,
						"nic_identifier": types.StringType,
					},
				},
			},
		},
		map[string]attr.Value{
			"propogate_vlan":     types.BoolValue(omeVlan.PropagateVLAN),
			"bonding_technology": types.StringValue(omeVlan.BondingTechnology),
			"vlan_attributes":    vlanAttrList,
		},
	)

	stateTemplate.Vlan = vlanTfsdk
}

func getVlanAtrrObject(omeVlanAttr models.OMEVlanAttribute) types.Object {
	vlanAttrMap := map[string]attr.Value{}
	vlanAttrMap["untagged_network"] = types.Int64Value(omeVlanAttr.Untagged)
	taggedNetworks := []attr.Value{}
	for _, tn := range omeVlanAttr.Tagged {
		taggedNetworks = append(taggedNetworks, types.Int64Value(tn))
	}

	vlanAttrMap["tagged_networks"], _ = types.SetValue(
		types.Int64Type,
		taggedNetworks,
	)
	vlanAttrMap["is_nic_bonded"] = types.BoolValue(omeVlanAttr.IsNICBonded)
	vlanAttrMap["port"] = types.Int64Value(omeVlanAttr.Port)
	vlanAttrMap["nic_identifier"] = types.StringValue(omeVlanAttr.NicIdentifier)
	vlanAttrObject, _ := types.ObjectValue(
		map[string]attr.Type{
			"untagged_network": types.Int64Type,
			"tagged_networks": types.SetType{
				ElemType: types.Int64Type,
			},
			"is_nic_bonded":  types.BoolType,
			"port":           types.Int64Type,
			"nic_identifier": types.StringType,
		}, vlanAttrMap)
	return vlanAttrObject
}
