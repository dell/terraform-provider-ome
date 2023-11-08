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
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &vlanNetworksDataSource{}
	_ datasource.DataSourceWithConfigure = &vlanNetworksDataSource{}
)

// NewVlanNetworkDataSource is a new datasource for VlanNetwork
func NewVlanNetworkDataSource() datasource.DataSource {
	return &vlanNetworksDataSource{}
}

type vlanNetworksDataSource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *vlanNetworksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*vlanNetworksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "vlannetworks_info"
}

func (g vlanNetworksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query templates from OME." +
			" The information fetched from this data source can be used for getting the details for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID for vlan networks data source.",
				Description:         "ID for vlan networks data source.",
				Computed:            true,
				// Optional:            true,
			},
			"vlan_networks": schema.ListNestedAttribute{
				MarkdownDescription: "List of vlan networks",
				Description:         "List of vlan networks",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vlan_id": schema.Int64Attribute{
							MarkdownDescription: "Unique ID for the vlan network.",
							Description:         "Unique ID for the vlan network.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the vlan network.",
							Description:         "Name of the vlan network.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the vlan network.",
							Description:         "Description of the vlan network.",
							Computed:            true,
						},
						"vlan_maximum": schema.Int64Attribute{
							MarkdownDescription: "Vlan maximum.",
							Description:         "Vlan maximum.",
							Computed:            true,
						},
						"vlan_minimum": schema.Int64Attribute{
							MarkdownDescription: "Vlan minimum.",
							Description:         "Vlan minimum.",
							Computed:            true,
						},
						"type": schema.Int64Attribute{
							MarkdownDescription: "Type of vlan.",
							Description:         "Type of vlan.",
							Computed:            true,
						},
						"internal_ref_nwuu_id": schema.StringAttribute{
							MarkdownDescription: "Reference ID for a vlan.",
							Description:         "Reference ID for a vlan.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Read resource information
func (g vlanNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.VLanNetworksTypeTfsdk
	state.ID = types.StringValue("0")

	omeClient, d := g.p.createOMESession(ctx, "datasource_vlannetworks_info Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	vlanNetworksOme, err := omeClient.GetAllVlanNetworks()
	if err != nil {
		resp.Diagnostics.AddWarning(
			"unable to get the vlan netowrk details",
			err.Error(),
		)
		return
	}

	for _, vn := range vlanNetworksOme {
		vlanNetTsfdk := models.VLanNetworksTfsdk{
			VlanID:            types.Int64Value(vn.ID),
			Name:              types.StringValue(vn.Name),
			Description:       types.StringValue(vn.Description),
			VLANMaximum:       types.Int64Value(vn.VLANMaximum),
			VLANMinimum:       types.Int64Value(vn.VLANMinimum),
			Type:              types.Int64Value(vn.Type),
			InternalRefNWUUID: types.StringValue(vn.InternalRefNWUUID),
		}
		state.VlanNetworks = append(state.VlanNetworks, vlanNetTsfdk)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
