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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceTemplateType struct{}

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

// Order Resource schema
func (r resourceTemplateType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource for managing template on OpenManage Enterprise.Updates are supported for the following parameters: `name`, `description`, `attributes`, `job_retry_count`, `sleep_interval`, `identity_pool_name`, `vlan`.",
		Version:             1,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Id of the template.",
				Description:         "Template ID",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				MarkdownDescription: "Name of the template.",
				Description:         "Name of the template.",
				Type:                types.StringType,
				Required:            true,
			},
			"fqdds": {
				MarkdownDescription: "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters. This field cannot be updated.",
				Description:         "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters. This field cannot be updated.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.String{Value: "All"}),
				},
				Validators: []tfsdk.AttributeValidator{
					validFqddsValidator{},
				},
			},
			"view_type": {
				MarkdownDescription: "OME template view type, supported types are Deployment, Compliance. This field cannot be updated.",
				Description:         "OME template view type, supported types are Deployment, Compliance. This field cannot be updated.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.String{Value: "Deployment"}),
				},
				Validators: []tfsdk.AttributeValidator{
					validTemplateViewTypeValidator{},
				},
			},
			"view_type_id": {
				MarkdownDescription: "OME template view type id.",
				Description:         "OME template view type id.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"device_type": {
				MarkdownDescription: "OME template device type, supported types are Server, Chassis. This field cannot be updated and is applicable only for importing xml.",
				Description:         "OME template device type, supported types are Server, Chassis. This field cannot be updated and is applicable only for importing xml.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.String{Value: "Server"}),
				},
				Validators: []tfsdk.AttributeValidator{
					validTemplateDeviceTypeValidator{},
				},
			},
			"content": {
				MarkdownDescription: "The XML content of template.. This field cannot be updated.",
				Description:         "The XML content of template.. This field cannot be updated.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"refdevice_servicetag": {
				MarkdownDescription: "Target device servicetag from which the template needs to be created. This field cannot be updated.",
				Description:         "Target device servicetag from which the template needs to be created. This field cannot be updated.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"refdevice_id": {
				MarkdownDescription: "Target device id from which the template needs to be created. This field cannot be updated.",
				Description:         "Target device id from which the template needs to be created. This field cannot be updated.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"reftemplate_name": {
				MarkdownDescription: "Reference Template name from which the template needs to be cloned. This field cannot be updated.",
				Description:         "Reference Template name from which the template needs to be cloned. This field cannot be updated.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"description": {
				MarkdownDescription: "Description of the template",
				Description:         "Description of the template",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"attributes": {
				MarkdownDescription: "List of attributes associated with a template. This field is ignored while creating a template.",
				Description:         "List of attributes associated with a template. This field is ignored while creating a template.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.ListType{
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
			"job_retry_count": {
				MarkdownDescription: "Number of times the job has to be polled to get the final status of the resource.",
				Description:         "Number of times the job has to be polled to get the final status of the resource.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: RetryCount}),
				},
			},
			"sleep_interval": {
				MarkdownDescription: "Sleep time interval for job polling in seconds.",
				Description:         "Sleep time interval for job polling in seconds.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: SleepInterval}),
				},
			},
			"identity_pool_name": {
				MarkdownDescription: "Identity Pool name to be attached with template.",
				Description:         "Identity Pool name to be attached with template.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"identity_pool_id": {
				MarkdownDescription: "ID of the Identity Pool attached with template.",
				Description:         "ID of the Identity Pool attached with template.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"vlan": {
				MarkdownDescription: "VLAN details to be attached with template.",
				Description:         "VLAN details to be attached with template.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
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
		},
	}, nil
}

// New resource instance
func (r resourceTemplateType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTemplate{
		p: *(p.(*provider)),
	}, nil
}

type resourceTemplate struct {
	p provider
}

// Create a new resource
func (r resourceTemplate) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	//Read the data from Plan
	tflog.Trace(ctx, "resource_template create: started")
	var plan models.Template
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "resource_template create: reference data", map[string]interface{}{
		"refdeviceid":         plan.RefdeviceID.Value,
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

	viewTypeID, err := omeClient.GetViewTypeID(plan.ViewType.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateTemplate, err.Error(),
		)
		return
	}

	deviceTypeID, err := omeClient.GetDeviceTypeID(plan.DeviceType.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateTemplate, err.Error(),
		)
		return
	}

	omeTemplateData := models.OMETemplate{}
	var templateID int64

	if plan.ReftemplateName.Value != "" {
		tflog.Info(ctx, "resource_template create: creating a template from a reference template")

		if !plan.Description.Unknown {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, "description cannot be modified while cloning from a reference template.",
			)
			return
		}

		// The identity pool and Vlans does not get cloned into the new template in OME.
		sourceTemplate, err := omeClient.GetTemplateByName(plan.ReftemplateName.Value)
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
			NewTemplateName:  plan.Name.Value,
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
	} else if plan.Content.Value != "" { // template import
		tflog.Info(ctx, "resource_template create: creating a template from a xml content")

		if !plan.Description.Unknown {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, "description is not supported for template import operation.",
			)
			return
		}

		importTemplateRequest := models.OMEImportTemplate{
			ViewTypeID: viewTypeID,
			Type:       deviceTypeID,
			Name:       plan.Name.Value,
			Content:    plan.Content.Value,
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
		deviceID, err := omeClient.ValidateDevice(plan.RefdeviceServicetag.Value, plan.RefdeviceID.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrCreateTemplate, err.Error(),
			)
			return
		}

		ct := models.CreateTemplate{
			Fqdds:          strings.ReplaceAll(plan.FQDDS.Value, " ", ""),
			ViewTypeID:     viewTypeID,
			SourceDeviceID: deviceID,
			Name:           plan.Name.Value,
			Description:    plan.Description.Value,
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

		isSuccess, message := omeClient.TrackJob(omeTemplateData.TaskID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
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
	template.RefdeviceServicetag.Value = plan.RefdeviceServicetag.Value
	template.ReftemplateName.Value = plan.ReftemplateName.Value
	template.RefdeviceID.Value = plan.RefdeviceID.Value
	if plan.ReftemplateName.Value != "" {
		template.RefdeviceID.Value = omeTemplateData.SourceDeviceID
	}
	template.Content.Value = plan.Content.Value
	template.ViewType.Value = plan.ViewType.Value
	template.DeviceType.Value = plan.DeviceType.Value
	template.FQDDS.Value = plan.FQDDS.Value // The default value of fqdds is set to `All`. So if the config doesn't have any value specified, the default value in the plan is `All`.
	template.JobRetryCount = plan.JobRetryCount
	template.SleepInterval = plan.SleepInterval
	template.IdentityPoolName.Value = plan.IdentityPoolName.Value

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
func (r resourceTemplate) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Trace(ctx, "resource_template read: started")
	var template models.Template
	diags := req.State.Get(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID, _ := strconv.ParseInt(template.ID.Value, 10, 64)
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

	stateAttributes := []models.Attribute{}
	stateAttributeObjects := []types.Object{}
	template.Attributes.ElementsAs(ctx, &stateAttributeObjects, true)

	for _, stateAttrObject := range stateAttributeObjects {
		stateAttribute := models.Attribute{}
		stateAttrObject.As(ctx, &stateAttribute, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
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
	diags = template.Vlan.As(ctx, &stateVlan, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
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

	omeVlan.PropagateVLAN = stateVlan.PropogateVlan.Value

	if omeTemplateData.IdentityPoolID != 0 {
		identityPool, err := omeClient.GetIdentityPoolByID(omeTemplateData.IdentityPoolID)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrReadTemplate,
				fmt.Sprintf("Unable to fetch the identity pools: %s", err.Error()),
			)
			return
		}
		template.IdentityPoolName = types.String{Value: identityPool.Name}
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
func (r resourceTemplate) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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
	templateID, _ := strconv.ParseInt(stateTemplate.ID.Value, 10, 64)

	if isConfigValuesChanged(planTemplate, stateTemplate) {
		resp.Diagnostics.AddError(
			clients.ErrUpdateTemplate,
			"cannot update the following fields : `refdevice_servicetag`,`refdevice_id`,`view_type`, `reftemplate_name`, `content` and `fqdds`",
		)
		return
	}

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

	tflog.Debug(ctx, "resource_template update: Template id", map[string]interface{}{
		"templateid": templateID,
	})
	var identityPool models.IdentityPool
	if planTemplate.IdentityPoolName.Value != "" {
		identityPool, err = validateIOPoolName(omeClient, planTemplate.IdentityPoolName.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrUpdateTemplate,
				err.Error(),
			)
			return
		}
	}
	planVlan := getVlanForTemplate(ctx, resp, planTemplate)

	if !planTemplate.Vlan.Unknown {
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
		Name:        planTemplate.Name.Value,
		ID:          templateID,
		Description: stateTemplate.Description.Value,
	}

	if planTemplate.Description.Value != stateTemplate.Description.Value {
		updatePayload.Description = planTemplate.Description.Value
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
	var updateIdentityPool = planTemplate.IdentityPoolName.Value != stateTemplate.IdentityPoolName.Value
	updateVlan := !reflect.DeepEqual(planTemplate.Vlan.Attrs, stateTemplate.Vlan.Attrs)
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
			if planTemplate.IdentityPoolName.Value != "" {
				nwConfig.IdentityPoolID = identityPool.ID
			} else {
				nwConfig.IdentityPoolID = 0
			}
		} else {
			nwConfig.IdentityPoolID = stateTemplate.IdentityPoolID.Value
		}
		if planTemplate.IdentityPoolName.Value != "" {
			identityPool, err := omeClient.GetIdentityPoolByName(planTemplate.IdentityPoolName.Value)
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

	stateTemplate.IdentityPoolName.Value = planTemplate.IdentityPoolName.Value
	stateTemplate.SleepInterval.Value = planTemplate.SleepInterval.Value
	stateTemplate.ViewType.Value = planTemplate.ViewType.Value
	stateTemplate.FQDDS.Value = planTemplate.FQDDS.Value
	stateTemplate.JobRetryCount.Value = planTemplate.JobRetryCount.Value

	tflog.Trace(ctx, "resource_template update: updating state data started")

	tfsdkVlan := models.Vlan{}
	planTemplate.Vlan.As(ctx, &tfsdkVlan, types.ObjectAsOptions{UnhandledNullAsEmpty: true})

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
func (r resourceTemplate) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Trace(ctx, "resource_template delete: started")
	var template models.Template
	diags := resp.State.Get(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	tflog.Trace(ctx, "resource_template delete: started delete")
	tflog.Debug(ctx, "resource_template delete: started delete for template", map[string]interface{}{
		"templateid": template.ID.Value,
	})

	_, err = omeClient.Delete(fmt.Sprintf(clients.TemplateAPI+"(%s)", template.ID.Value), nil, nil)
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
func (r resourceTemplate) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Trace(ctx, "resource_template import: started")
	var template models.Template
	template.Name.Value = req.ID
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

	template.RefdeviceID = types.Int64{Value: omeTemplateData.SourceDeviceID}
	template.RefdeviceServicetag = types.String{Value: "NA"}
	template.ReftemplateName = types.String{Value: "NA"}
	template.Content = types.String{Value: "NA"}
	template.ViewType = types.String{Value: viewType}
	template.DeviceType = types.String{Value: deviceType}
	template.JobRetryCount = types.Int64{Value: RetryCount}
	template.SleepInterval = types.Int64{Value: SleepInterval}
	template.FQDDS = types.String{Value: "All"}
	diags := resp.State.Set(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_template import: finished")
}

func validateCreate(plan models.Template) error {
	// all references cannot be empty
	if plan.ReftemplateName.Value == "" && plan.RefdeviceID.Value == 0 && plan.RefdeviceServicetag.Value == "" && plan.Content.Value == "" {
		return fmt.Errorf("either reftemplate_name or refdevice_id or refdevice_servicetag or content is required")
	}

	// any two references given results in error
	if (plan.ReftemplateName.Value != "" && (plan.RefdeviceID.Value != 0 || plan.RefdeviceServicetag.Value != "" || plan.Content.Value != "")) ||
		(plan.RefdeviceID.Value != 0 && (plan.ReftemplateName.Value != "" || plan.RefdeviceServicetag.Value != "" || plan.Content.Value != "")) ||
		(plan.RefdeviceServicetag.Value != "" && (plan.RefdeviceID.Value != 0 || plan.ReftemplateName.Value != "" || plan.Content.Value != "")) ||
		(plan.Content.Value != "" && (plan.RefdeviceID.Value != 0 || plan.ReftemplateName.Value != "" || plan.RefdeviceServicetag.Value != "")) {
		return fmt.Errorf("either reftemplate_name or refdevice_id or refdevice_servicetag or content is required")
	}

	// Identity Pool name and Vlan is supported only during update
	if !plan.IdentityPoolName.Unknown || !plan.Vlan.Unknown {
		return fmt.Errorf("attributes identity_pool_name and vlan cannot be associated during create")
	}

	// Attributes is part of plan during create
	if !plan.Attributes.Unknown {
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
	return (!planTemplate.RefdeviceID.Unknown && stateTemplate.RefdeviceID.Value != planTemplate.RefdeviceID.Value) ||
		(!planTemplate.RefdeviceServicetag.Unknown && stateTemplate.RefdeviceServicetag.Value != planTemplate.RefdeviceServicetag.Value) ||
		(!planTemplate.ViewType.Unknown && stateTemplate.ViewType.Value != planTemplate.ViewType.Value) ||
		(!planTemplate.FQDDS.Unknown && stateTemplate.FQDDS.Value != planTemplate.FQDDS.Value) ||
		(!planTemplate.ReftemplateName.Unknown && stateTemplate.ReftemplateName.Value != planTemplate.ReftemplateName.Value) ||
		(!planTemplate.Content.Unknown && stateTemplate.Content.Value != planTemplate.Content.Value)
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

func getVlanForTemplate(ctx context.Context, resp *tfsdk.UpdateResourceResponse, Template models.Template) models.OMEVlan {
	omeVlan := models.OMEVlan{}
	vlan := models.Vlan{}

	Template.Vlan.As(ctx, &vlan, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
	omeVlan.BondingTechnology = vlan.BondingTechnology.Value
	omeVlan.PropagateVLAN = vlan.PropogateVlan.Value
	omeVlanAttrs := []models.OMEVlanAttribute{}
	vlanAttrs := []models.VlanAttributes{}
	vlan.VlanAttributes.ElementsAs(ctx, &vlanAttrs, true)
	for _, vlanAttr := range vlanAttrs {
		taggedNetworks := []int64{}
		vlanAttr.TaggedNetworks.ElementsAs(ctx, &taggedNetworks, true)
		omeVlanAttr := models.OMEVlanAttribute{
			Untagged:      vlanAttr.UntaggedNetwork.Value,
			Tagged:        taggedNetworks,
			IsNICBonded:   vlanAttr.IsNicBonded.Value,
			Port:          vlanAttr.Port.Value,
			NicIdentifier: vlanAttr.NicIdentifier.Value,
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
		stateAttrObject.As(ctx, &stateAttribute, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
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
		planUpdateAttrObject.As(ctx, &planUpdateAttribute, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
		planUpdateAttributes = append(planUpdateAttributes, planUpdateAttribute)
	}
	if !reflect.DeepEqual(planUpdateAttributes, stateAttributes) {
		for index, attribute := range planUpdateAttributes {
			if attribute.Value != stateAttributes[index].Value {
				updateAttribute := models.UpdateAttribute{
					ID:        attribute.AttributeID.Value,
					IsIgnored: attribute.IsIgnored.Value,
					Value:     attribute.Value.Value,
				}
				updatedAttributes = append(updatedAttributes, updateAttribute)
			}
		}
	}

	return updatedAttributes, nil
}

func updateState(stateTemplate *models.Template, planVlanAttributes []models.VlanAttributes, omeTemplateData *models.OMETemplate, omeTemplateAttributes []models.OmeAttribute, omeVlan models.OMEVlan) {

	stateTemplate.ID.Value = fmt.Sprintf("%d", omeTemplateData.ID)
	stateTemplate.Name.Value = omeTemplateData.Name
	stateTemplate.Description.Value = omeTemplateData.Description
	stateTemplate.ViewTypeID.Value = omeTemplateData.ViewTypeID
	stateTemplate.IdentityPoolID.Value = omeTemplateData.IdentityPoolID

	attributesTfsdk := types.List{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"attribute_id": types.Int64Type,
				"display_name": types.StringType,
				"value":        types.StringType,
				"is_ignored":   types.BoolType,
			},
		},
	}
	attributeObjects := []attr.Value{}

	for _, attribute := range omeTemplateAttributes {
		attributeDetails := map[string]attr.Value{}
		attributeDetails["attribute_id"] = types.Int64{Value: attribute.AttributeID}
		attributeDetails["display_name"] = types.String{Value: attribute.DisplayName}
		attributeDetails["value"] = types.String{Value: attribute.Value}
		attributeDetails["is_ignored"] = types.Bool{Value: attribute.IsIgnored}
		attributeObject := types.Object{
			Attrs: attributeDetails,
			AttrTypes: map[string]attr.Type{
				"attribute_id": types.Int64Type,
				"display_name": types.StringType,
				"value":        types.StringType,
				"is_ignored":   types.BoolType,
			},
		}
		attributeObjects = append(attributeObjects, attributeObject)
	}
	attributesTfsdk.Elems = attributeObjects
	stateTemplate.Attributes = attributesTfsdk

	omeVlanMap := map[string]models.OMEVlanAttribute{}

	for _, vlanAttr := range omeVlan.OMEVlanAttributes {
		key := fmt.Sprintf(NicIdentifierAndPort, vlanAttr.NicIdentifier, vlanAttr.Port)
		omeVlanMap[key] = vlanAttr
	}

	vlanAttrsObjects := []attr.Value{}

	for _, planVlanAttr := range planVlanAttributes {
		key := fmt.Sprintf(NicIdentifierAndPort, planVlanAttr.NicIdentifier.Value, planVlanAttr.Port.Value)
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

	vlanAttrList := types.List{
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
		Elems: vlanAttrsObjects,
	}

	vlanTfsdk := types.Object{
		AttrTypes: map[string]attr.Type{
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
		Attrs: map[string]attr.Value{
			"propogate_vlan":     types.Bool{Value: omeVlan.PropagateVLAN},
			"bonding_technology": types.String{Value: omeVlan.BondingTechnology},
			"vlan_attributes":    vlanAttrList,
		},
	}

	stateTemplate.Vlan = vlanTfsdk
}

func getVlanAtrrObject(omeVlanAttr models.OMEVlanAttribute) types.Object {
	vlanAttrObject := types.Object{
		AttrTypes: map[string]attr.Type{
			"untagged_network": types.Int64Type,
			"tagged_networks": types.SetType{
				ElemType: types.Int64Type,
			},
			"is_nic_bonded":  types.BoolType,
			"port":           types.Int64Type,
			"nic_identifier": types.StringType,
		},
	}
	vlanAttrMap := map[string]attr.Value{}
	vlanAttrMap["untagged_network"] = types.Int64{Value: omeVlanAttr.Untagged}
	taggedNetworks := []attr.Value{}
	for _, tn := range omeVlanAttr.Tagged {
		taggedNetworks = append(taggedNetworks, types.Int64{Value: tn})
	}

	vlanAttrMap["tagged_networks"] = types.Set{
		ElemType: types.Int64Type,
		Elems:    taggedNetworks,
	}
	vlanAttrMap["is_nic_bonded"] = types.Bool{Value: omeVlanAttr.IsNICBonded}
	vlanAttrMap["port"] = types.Int64{Value: omeVlanAttr.Port}
	vlanAttrMap["nic_identifier"] = types.String{Value: omeVlanAttr.NicIdentifier}
	vlanAttrObject.Attrs = vlanAttrMap
	return vlanAttrObject
}
