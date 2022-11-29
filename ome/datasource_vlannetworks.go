package ome

import (
	"context"
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vlanNetowrksDataSourceType struct{}

func (t vlanNetowrksDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Data source to list the vlan networks from OpenManage Enterprise.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID for data source.",
				Description:         "ID for data source.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
			},
			"vlan_networks": {
				MarkdownDescription: "List of vlan networks",
				Description:         "List of vlan networks",
				Computed:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"vlan_id": {
						MarkdownDescription: "Unique ID for the vlan network.",
						Description:         "Unique ID for the vlan network.",
						Type:                types.Int64Type,
						Computed:            true,
					},
					"name": {
						MarkdownDescription: "Name of the vlan network.",
						Description:         "Name of the vlan network.",
						Type:                types.StringType,
						Computed:            true,
					},
					"description": {
						MarkdownDescription: "Description of the vlan network.",
						Description:         "Description of the vlan network.",
						Type:                types.StringType,
						Computed:            true,
					},
					"vlan_maximum": {
						MarkdownDescription: "Vlan maximum.",
						Description:         "Vlan maximum.",
						Type:                types.Int64Type,
						Computed:            true,
					},
					"vlan_minimum": {
						MarkdownDescription: "Vlan minimum.",
						Description:         "Vlan minimum.",
						Type:                types.Int64Type,
						Computed:            true,
					},
					"type": {
						MarkdownDescription: "Type of vlan.",
						Description:         "Type of vlan.",
						Type:                types.Int64Type,
						Computed:            true,
					},
					"internal_ref_nwuu_id": {
						MarkdownDescription: "Reference ID for a vlan.",
						Description:         "Reference ID for a vlan.",
						Type:                types.StringType,
						Computed:            true,
					},
				}),
			},
		},
	}, nil
}

func (t vlanNetowrksDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return vlanNetowrksDataSource{
		p: provider,
	}, diags
}

type vlanNetowrksDataSource struct {
	p provider
}

// Read resource information
func (g vlanNetowrksDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var state models.VLanNetworksTypeTfsdk

	omeClient, err := clients.NewClient(*g.p.clientOpt)
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
			VlanID:            types.Int64{Value: vn.ID},
			Name:              types.String{Value: vn.Name},
			Description:       types.String{Value: vn.Description},
			VLANMaximum:       types.Int64{Value: vn.VLANMaximum},
			VLANMinimum:       types.Int64{Value: vn.VLANMinimum},
			Type:              types.Int64{Value: vn.Type},
			InternalRefNWUUID: types.String{Value: vn.InternalRefNWUUID},
		}
		state.VlanNetworks = append(state.VlanNetworks, vlanNetTsfdk)
	}

	fmt.Println("[DEBUG]-Resource State: ", state)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
