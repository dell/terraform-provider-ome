package ome

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
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
)

// Order Resource schema
func (r resourceTemplateType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Resource for managing template on OpenManage Enterprise. Updates are supported for the following parameters: `name`, `description`, `attributes`, `job_retry_count`, `sleep_interval`, `identity_pool_name`, `vlan`.",
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
				MarkdownDescription: "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters.",
				Description:         "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters.",
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
				MarkdownDescription: "OME template view type, supported types are Deployment, Compliance.",
				Description:         "OME template view type, supported types are Deployment, Compliance.",
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
			"refdevice_servicetag": {
				MarkdownDescription: "Target device servicetag from which the template needs to be created.",
				Description:         "Target device servicetag from which the template needs to be created.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"refdevice_id": {
				MarkdownDescription: "Target device id from which the template needs to be created.",
				Description:         "Target device id from which the template needs to be created.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"reftemplate_name": {
				MarkdownDescription: "Reference Template name from which the template needs to be cloned.",
				Description:         "Reference Template name from which the template needs to be cloned.",
				Type:                types.StringType,
				Optional:            true,
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
					DefaultAttribute(types.Int64{Value: 5}),
				},
			},
			"sleep_interval": {
				MarkdownDescription: "Sleep time interval for job polling in seconds.",
				Description:         "Sleep time interval for job polling in seconds.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 30}),
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
									"tagged_networks": types.ListType{
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
	if plan.ReftemplateName.Value != "" && (plan.RefdeviceID.Value != 0 || plan.RefdeviceServicetag.Value != "") {
		resp.Diagnostics.AddError(
			"error creating/cloning the template", "please provide either reftemplate_name or refdevice_id/refdevice_servicetag",
		)
		return
	}

	//Create a Template
	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create OME session: ",
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	viewTypeID, err := omeClient.GetViewTypeID(plan.ViewType.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"error creating the template", err.Error(),
		)
		return
	}

	omeTemplateData := models.OMETemplate{}
	var templateID int64

	if plan.ReftemplateName.Value != "" {

		tflog.Info(ctx, "resource_template create: creating a template from a reference template")

		// The identity pool and Vlans does not get cloned into the new template in OME.
		sourceTemplate, err := omeClient.GetTemplateByName(plan.ReftemplateName.Value)
		if err != nil || sourceTemplate.Name == "" {
			resp.Diagnostics.AddError(
				"error cloning the template with given reference template name", "Unable to clone the template because Source template does not exist.",
			)
			return
		}

		if sourceTemplate.ViewTypeID == ComplianceViewTypeID && plan.ViewType.Value == "Deployment" {
			resp.Diagnostics.AddError(
				"error cloning the template", "cannot clone compliance template as deployment template.",
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
				"error cloning the template", err.Error(),
			)
			return
		}

		omeTemplateData, err = omeClient.GetTemplateByID(templateID)
		if err != nil {
			resp.Diagnostics.AddError(
				"error fetching the cloned template", err.Error(),
			)
			return
		}
	} else {
		tflog.Info(ctx, "resource_template create: creating a template from a reference device")

		deviceID, err := omeClient.ValidateDevice(plan.RefdeviceServicetag.Value, plan.RefdeviceID.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"error creating the template: ", err.Error(),
			)
			return
		}

		ct := models.CreateTemplate{
			Fqdds:          plan.FQDDS.Value,
			ViewTypeID:     viewTypeID,
			SourceDeviceID: deviceID,
			Name:           plan.Name.Value,
			Description:    plan.Description.Value,
		}

		templateID, err = omeClient.CreateTemplate(ct)
		if err != nil {
			resp.Diagnostics.AddError(
				"error creating the template", err.Error(),
			)
			return
		}

		log.Printf("template created with id %d", templateID)
		time.Sleep(2 * time.Second)
		omeTemplateData, err = omeClient.GetTemplateByID(templateID)
		if err != nil {
			resp.Diagnostics.AddError(
				"error creating the template", err.Error(),
			)
			return
		}

		isSuccess, message := omeClient.TrackJob(omeTemplateData.TaskID, plan.JobRetryCount.Value, plan.SleepInterval.Value)
		if !isSuccess {
			resp.Diagnostics.AddError(
				"template creation failed with status error ", message,
			)
			//TBD : Delete a template
			return
		}
	}

	tflog.Trace(ctx, "resource_template create: fetching template attributes")

	omeAttributes, err := omeClient.GetTemplateAttributes(omeTemplateData.ID, []models.Attribute{}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh template attributes:",
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_template create: fetching template valn data")
	omeVlan, err := omeClient.GetSchemaVlanData(templateID)
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to refresh vlan attributes:",
				err.Error(),
			)
			return
		}
	}
	//PropogateVlan is default true and is not available in the response from OME,hence setting it to true here to persist in state
	omeVlan.PropagateVLAN = true
	template.RefdeviceServicetag.Value = plan.RefdeviceServicetag.Value
	template.ViewType.Value = plan.ViewType.Value
	template.FQDDS.Value = plan.FQDDS.Value // The default value of fqdds is set to `All`. So if the config doesn't have any value specified, the default value in the plan is `All`.
	template.JobRetryCount = plan.JobRetryCount
	template.SleepInterval = plan.SleepInterval
	template.ReftemplateName = plan.ReftemplateName

	tflog.Trace(ctx, "resource_template create: started updating state")

	updateState(&template, &omeTemplateData, omeAttributes, omeVlan)

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
			"Unable to create client",
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create OME session: ",
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
			"error reading the template", err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_template read: fetching template attributes")

	omeAttributes, err := omeClient.GetTemplateAttributes(templateID, stateAttributes, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh template attributes:",
			err.Error(),
		)
		return
	}

	stateVlan := models.Vlan{}
	diags = template.Vlan.As(ctx, &stateVlan, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Unable to fetch Vlan from state ",
			"Hence, Cannot refresh the template resource",
		)
		return
	}
	tflog.Trace(ctx, "resource_template read: fetching template vlan data")

	omeVlan, err := omeClient.GetSchemaVlanData(templateID)
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to refresh vlan attributes:",
				err.Error(),
			)
			return
		}
	}

	omeVlan.PropagateVLAN = stateVlan.PropogateVlan.Value
	tflog.Trace(ctx, "resource_template read: updating state started")

	updateState(&template, &omeTemplateData, omeAttributes, omeVlan)

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
	omeClient, err := clients.NewClient(*r.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create OME session: ",
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	tflog.Debug(ctx, "resource_template update: Template id", map[string]interface{}{
		"templateid": templateID,
	})

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
			"Unable to update the template",
			err.Error(),
		)
		return
	}

	// Updating networkConfig now
	var updateIdentityPool = planTemplate.IdentityPoolName.Value != stateTemplate.IdentityPoolName.Value

	var updateVlan = false
	planVlan, err := getVlanForTemplate(ctx, resp, planTemplate)
	if err == nil {
		tflog.Info(ctx, "resource_template update:checking if VLAN attrs are equal")
		updateVlan = !reflect.DeepEqual(planTemplate.Vlan.Attrs, stateTemplate.Vlan.Attrs)
	}

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
		}
		if updateVlan {
			tflog.Info(ctx, "resource_template update: updating vlan attrs")
			vlanNetworkView, err := omeClient.GetVlanNetworkModel(templateID)
			if err != nil {
				resp.Diagnostics.AddWarning(
					fmt.Sprintf("Unable to fetch network view from OME to update vlan to the template: %d", templateID),
					err.Error(),
				)
			} else {
				tflog.Info(ctx, "resource_template update: fetching vlan attrs")
				stateVlan, err := getVlanForTemplate(ctx, resp, stateTemplate)
				if err != nil {
					resp.Diagnostics.AddWarning(
						"Unable to fetch vlan data from state to create payload for vlan: ",
						err.Error(),
					)
				} else {
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

		}
		tflog.Info(ctx, "resource_template update: triggering update netowrk config")
		err = omeClient.UpdateNetworkConfig(nwConfig)
		if err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Unable to update network configuration to the template: %d", templateID),
				err.Error(),
			)
		}

	}

	tflog.Trace(ctx, "resource_template update: fetching template by id")
	omeTemplateData, err := omeClient.GetTemplateByID(templateID)

	if err != nil {
		resp.Diagnostics.AddError(
			"error creating the template", err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_template update: fetching template attributes")

	omeAttributes, err := omeClient.GetTemplateAttributes(templateID, stateAttributes, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh template attributes:",
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_template update: fetching vlan data")

	updatedVlan, err := omeClient.GetSchemaVlanData(templateID)
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to refresh vlan attributes:",
				err.Error(),
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

	updateState(&stateTemplate, &omeTemplateData, omeAttributes, updatedVlan)

	tflog.Trace(ctx, "resource_template update: updating state data finished")
	//Save into State if template update is successful
	diags := resp.State.Set(ctx, &stateTemplate)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_template update: finished")
}

func getVlanForTemplate(ctx context.Context, resp *tfsdk.UpdateResourceResponse, Template models.Template) (models.OMEVlan, error) {
	omeVlan := models.OMEVlan{}
	vlan := models.Vlan{}

	vlanDiags := Template.Vlan.As(ctx, &vlan, types.ObjectAsOptions{UnhandledNullAsEmpty: true})
	if vlanDiags.HasError() {
		resp.Diagnostics.Append(vlanDiags...)
		resp.Diagnostics.AddWarning(
			clients.ErrUnableToParseVlan,
			"Vlan attributes update cannot be done",
		)
		return omeVlan, fmt.Errorf(clients.ErrUnableToParseVlan)
	}

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
	return omeVlan, nil
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
			"Unable to create client",
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create OME session: ",
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
			fmt.Sprintf("Unable to delete template: %s", template.ID.Value),
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
			"Unable to create client",
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create OME session: ",
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	omeTemplateData, err := omeClient.GetTemplateByName(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading the template", err.Error(),
		)
		return
	}

	if omeTemplateData.ID == 0 {
		resp.Diagnostics.AddError(
			"unable to get template", "invalid template name",
		)
		return
	}

	omeAttributes, err := omeClient.GetTemplateAttributes(omeTemplateData.ID, []models.Attribute{}, true)
	if err != nil {
		fmt.Printf("Error in get template attr: %s", err.Error())
		resp.Diagnostics.AddError(
			"unable to get template attributes:",
			err.Error(),
		)
		return
	}

	omeVlan, err := omeClient.GetSchemaVlanData(omeTemplateData.ID)
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to refresh vlan attributes:",
				err.Error(),
			)
			return
		}
	}
	tflog.Trace(ctx, "resource_template import: started state update")
	updateState(&template, &omeTemplateData, omeAttributes, omeVlan)
	tflog.Trace(ctx, "resource_template import: finished state update")
	diags := resp.State.Set(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_template import: finished")
}

func updateState(template *models.Template, omeTemplateData *models.OMETemplate, omeTemplateAttributes []models.OmeAttribute, omeVlan models.OMEVlan) {

	template.ID.Value = fmt.Sprintf("%d", omeTemplateData.ID)
	template.Name.Value = omeTemplateData.Name
	template.Description.Value = omeTemplateData.Description
	template.ViewTypeID.Value = omeTemplateData.ViewTypeID
	template.RefdeviceID.Value = omeTemplateData.SourceDeviceID
	template.IdentityPoolID.Value = omeTemplateData.IdentityPoolID

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
	template.Attributes = attributesTfsdk

	var vlanTfsdk types.Object
	vlanAttrsObjects := []attr.Value{}

	for _, omeVlanAttr := range omeVlan.OMEVlanAttributes {
		vlanAttrObject := types.Object{
			AttrTypes: map[string]attr.Type{
				"untagged_network": types.Int64Type,
				"tagged_networks": types.ListType{
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

		vlanAttrMap["tagged_networks"] = types.List{
			ElemType: types.Int64Type,
			Elems:    taggedNetworks,
		}
		vlanAttrMap["is_nic_bonded"] = types.Bool{Value: omeVlanAttr.IsNICBonded}
		vlanAttrMap["port"] = types.Int64{Value: omeVlanAttr.Port}
		vlanAttrMap["nic_identifier"] = types.String{Value: omeVlanAttr.NicIdentifier}
		vlanAttrObject.Attrs = vlanAttrMap
		vlanAttrsObjects = append(vlanAttrsObjects, vlanAttrObject)
	}

	vlanAttrList := types.List{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"untagged_network": types.Int64Type,
				"tagged_networks": types.ListType{
					ElemType: types.Int64Type,
				},
				"is_nic_bonded":  types.BoolType,
				"port":           types.Int64Type,
				"nic_identifier": types.StringType,
			},
		},
		Elems: vlanAttrsObjects,
	}
	vlanTfsdk = types.Object{
		AttrTypes: map[string]attr.Type{
			"propogate_vlan":     types.BoolType,
			"bonding_technology": types.StringType,
			"vlan_attributes": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"untagged_network": types.Int64Type,
						"tagged_networks": types.ListType{
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

	template.Vlan = vlanTfsdk
}
