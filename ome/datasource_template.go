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
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ datasource.DataSource              = &templateDataSource{}
	_ datasource.DataSourceWithConfigure = &templateDataSource{}
)

// NewTemplateDataSource is a new datasource for template
func NewTemplateDataSource() datasource.DataSource {
	return &templateDataSource{}
}

type templateDataSource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (t *templateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	t.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*templateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "template_info"
}

func (t templateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data Source to list the Template details from OpenManage Enterprise",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the template data source.",
				Description:         "ID of the template data source.",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the template.",
				Description:         "Name of the template.",
				Required:            true,
			},
			"view_type_id": schema.Int64Attribute{
				MarkdownDescription: "OME template view type id.",
				Description:         "OME template view type id.",
				Computed:            true,
			},
			"device_type_id": schema.Int64Attribute{
				MarkdownDescription: "Template type ID, indicating the type of device for which configuration is supported, current supported device is server",
				Description:         "Template type ID, indicating the type of device for which configuration is supported, current supported device is server",
				Computed:            true,
			},
			"refdevice_id": schema.Int64Attribute{
				MarkdownDescription: "Target device id from which the template is created.",
				Description:         "Target device id from which the template is created.",
				Optional:            true,
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The XML content of template from which the template will be created",
				Description:         "The XML content of template from which the template will be created",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the template.",
				Description:         "Description for the template.",
				Optional:            true,
				Computed:            true,
			},
			"attributes": schema.ListAttribute{
				MarkdownDescription: "List of attributes associated with template.",
				Description:         "List of attributes associated with template.",
				Optional:            true,
				Computed:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"attribute_id": types.Int64Type,
						"display_name": types.StringType,
						"value":        types.StringType,
						"is_ignored":   types.BoolType,
					},
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
				AttributeTypes: map[string]attr.Type{
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
	}
}

// Read resource information
func (t templateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var template models.TemplateDataSource
	diags := req.Config.Get(ctx, &template)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateName := template.Name.ValueString()
	omeClient, d := t.p.createOMESession(ctx, "datasource_template Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
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
	}
	stateVlan := models.Vlan{}
	diags = template.Vlan.As(ctx, &stateVlan, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
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

	omeVlan.PropagateVLAN = stateVlan.PropogateVlan.ValueBool()
	updateDataSourceState(&template, &omeTemplateData, omeAttributes, omeVlan)

	diags = resp.State.Set(ctx, &template)
	resp.Diagnostics.Append(diags...)
}

func updateDataSourceState(template *models.TemplateDataSource, omeTemplateData *models.OMETemplate, omeTemplateAttributes []models.OmeAttribute, omeVlan models.OMEVlan) {

	template.ID = types.StringValue(fmt.Sprintf("%d", omeTemplateData.ID))
	template.Name = types.StringValue(omeTemplateData.Name)
	template.Description = types.StringValue(omeTemplateData.Description)
	template.ViewTypeID = types.Int64Value(omeTemplateData.ViewTypeID)
	template.DeviceTypeID = types.Int64Value(omeTemplateData.TypeID)
	template.RefdeviceID = types.Int64Value(omeTemplateData.SourceDeviceID)
	template.IdentityPoolID = types.Int64Value(omeTemplateData.IdentityPoolID)
	template.Content = types.StringValue(omeTemplateData.Content)
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
			}, attributeDetails,
		)
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
		},
		attributeObjects,
	)

	template.Attributes = attributesTfsdk

	var vlanTfsdk types.Object
	vlanAttrsObjects := []attr.Value{}

	for _, omeVlanAttr := range omeVlan.OMEVlanAttributes {

		vlanAttrMap := map[string]attr.Value{}
		vlanAttrMap["untagged_network"] = types.Int64Value(omeVlanAttr.Untagged)
		taggedNetworks := []attr.Value{}
		for _, tn := range omeVlanAttr.Tagged {
			taggedNetworks = append(taggedNetworks, types.Int64Value(tn))
		}

		vlanAttrMap["tagged_networks"], _ = types.ListValue(
			types.Int64Type,
			taggedNetworks,
		)
		vlanAttrMap["is_nic_bonded"] = types.BoolValue(omeVlanAttr.IsNICBonded)
		vlanAttrMap["port"] = types.Int64Value(omeVlanAttr.Port)
		vlanAttrMap["nic_identifier"] = types.StringValue(omeVlanAttr.NicIdentifier)
		vlanAttrObject, _ := types.ObjectValue(
			map[string]attr.Type{
				"untagged_network": types.Int64Type,
				"tagged_networks": types.ListType{
					ElemType: types.Int64Type,
				},
				"is_nic_bonded":  types.BoolType,
				"port":           types.Int64Type,
				"nic_identifier": types.StringType,
			}, vlanAttrMap,
		)
		vlanAttrsObjects = append(vlanAttrsObjects, vlanAttrObject)
	}

	vlanAttrList, _ := types.ListValue(
		types.ObjectType{
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
		vlanAttrsObjects,
	)
	vlanTfsdk, _ = types.ObjectValue(
		map[string]attr.Type{
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
		map[string]attr.Value{
			"propogate_vlan":     types.BoolValue(omeVlan.PropagateVLAN),
			"bonding_technology": types.StringValue(omeVlan.BondingTechnology),
			"vlan_attributes":    vlanAttrList,
		},
	)

	template.Vlan = vlanTfsdk
}
