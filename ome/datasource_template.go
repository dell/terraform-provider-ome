package ome

import (
	"context"
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type templateDataSourceType struct{}

func (t templateDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Data Source to list the Template details from OpenManage Enterprise",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Id of the template.",
				Description:         "Id of the template.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
			},
			"name": {
				MarkdownDescription: "Name of the template.",
				Description:         "Name of the template.",
				Type:                types.StringType,
				Required:            true,
			},
			"view_type_id": {
				MarkdownDescription: "OME template view type id.",
				Description:         "OME template view type id.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"device_type_id": {
				MarkdownDescription: "Template type ID, indicating the type of device for which configuration is supported, current supported device is server",
				Description:         "Template type ID, indicating the type of device for which configuration is supported, current supported device is server",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"refdevice_id": {
				MarkdownDescription: "Target device id from which the template is created.",
				Description:         "Target device id from which the template is created.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"content": {
				MarkdownDescription: "The XML content of template from which the template will be created",
				Description:         "The XML content of template from which the template will be created",
				Type:                types.StringType,
				Optional:            true,
			},
			"description": {
				MarkdownDescription: "Description for the template.",
				Description:         "Description for the template.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"attributes": {
				MarkdownDescription: "List of attributes associated with template.",
				Description:         "List of attributes associated with template.",
				Optional:            true,
				Computed:            true,
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

func (t templateDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return templateDataSource{
		p: provider,
	}, diags
}

type templateDataSource struct {
	p provider
}

// Read resource information
func (t templateDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var template models.TemplateDataSource
	diags := req.Config.Get(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateName := template.Name.Value
	omeClient, err := clients.NewClient(*t.p.clientOpt)
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

	omeTemplateData, err := omeClient.GetTemplateByName(templateName)
	if err == nil && omeTemplateData.Name == "" {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading the template", err.Error(),
		)
		return
	}

	omeAttributes, err := omeClient.GetTemplateAttributes(omeTemplateData.ID, stateAttributes, true)
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

	omeVlan.PropagateVLAN = stateVlan.PropogateVlan.Value
	updateDataSourceState(&template, &omeTemplateData, omeAttributes, omeVlan)

	diags = resp.State.Set(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func updateDataSourceState(template *models.TemplateDataSource, omeTemplateData *models.OMETemplate, omeTemplateAttributes []models.OmeAttribute, omeVlan models.OMEVlan) {

	template.ID = types.String{Value: fmt.Sprintf("%d", omeTemplateData.ID)}
	template.Name = types.String{Value: omeTemplateData.Name}
	template.Description = types.String{Value: omeTemplateData.Description}
	template.ViewTypeID = types.Int64{Value: omeTemplateData.ViewTypeID}
	template.DeviceTypeID = types.Int64{Value: omeTemplateData.TypeID}
	template.RefdeviceID = types.Int64{Value: omeTemplateData.SourceDeviceID}
	template.IdentityPoolID = types.Int64{Value: omeTemplateData.IdentityPoolID}
	template.Content = types.String{Value: omeTemplateData.Content}

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
