package ome

import (
	"context"

	// "math/big"
	// "strconv"
	// "time"

	// "github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceTemplateType struct{}

// Order Resource schema
func (r resourceTemplateType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Schema for Template on OpenManage Enterprise",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Name of the template",
				Description:         "Name of the template",
				Type:                types.StringType,
				Required:            true,
			},
			"fqdds": {
				MarkdownDescription: "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters",
				Description:         "Comma seperated values of components from a specified server, should be one of these iDRAC, System, BIOS, NIC, LifeCycleController, RAID, and EventFilters",
				Type:                types.StringType,
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.String{Value: "All"}),
				},
			},
			"viewtype": {
				MarkdownDescription: "OME template view type, should be one Deployment, Compliance, Sample",
				Description:         "OME template view type, should be one Deployment, Compliance, Sample",
				Type:                types.StringType,
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.String{Value: "Deployment"}),
				},
			},
			"type": {
				MarkdownDescription: "Template type ID, indicating the type of device for which configuration is supported, current supported device is server",
				Description:         "Template type ID, indicating the type of device for which configuration is supported, current supported device is server",
				Type:                types.StringType,
				Computed:            true,
			},
			"refdevice_servicetag": {
				MarkdownDescription: "Target device servicetag from which the template needs to be created",
				Description:         "Target device servicetag from which the template needs to be created",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"refdevice_id": {
				MarkdownDescription: "Target device id from which the template needs to be created",
				Description:         "Target device id from which the template needs to be created",
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
				MarkdownDescription: "Description for the template",
				Description:         "Description for the template",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"template_id": {
				MarkdownDescription: "Template ID",
				Description:         "Template ID",
				Type:                types.StringType,
				Computed:            true,
			},
			"attributes": {
				MarkdownDescription: "List of template attributes",
				Description:         "List of template attributes",
				Optional:            true,
				Computed:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"attribute_id": {
						MarkdownDescription: "Unique identifier of an attribute",
						Description:         "Unique identifier of an attribute",
						Type:                types.Int64Type,
						Computed:            true,
					},
					"display_name": {
						MarkdownDescription: "Display Name of the attribute as per the GUI",
						Description:         "Display Name of the attribute as per the GUI",
						Type:                types.StringType,
						Required:            true,
					},
					"value": {
						MarkdownDescription: "Value of an attribute",
						Description:         "Value of an attribute",
						Type:                types.StringType,
						Required:            true,
					},
					"is_ignored": {
						MarkdownDescription: "Indicates whether the attribute should be ignored or included when the template is deployed to another device",
						Description: "Indicates whether the attribute should be ignored or included	when the template is deployed to another device",
						Type:     types.BoolType,
						Computed: true,
					},
				}),
			},
			"identity_pool_name": {
				MarkdownDescription: "Identity Pool name to be attached with template",
				Description:         "Identity Pool name to be attached with template",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"vlan": {
				MarkdownDescription: "VLAN details to be attached with template",
				Description:         "VLAN details to be attached with template",
				Computed:            true,
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"propogate_vlan": {
						MarkdownDescription: "To deploy the modified VLAN settings immediately without rebooting the server",
						Description:         "To deploy the modified VLAN settings immediately without rebooting the server",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"bonding_technology": {
						MarkdownDescription: "Identity Pool name to be attached with template",
						Description:         "Identity Pool name to be attached with template",
						Type:                types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"vlan_attributes": {
						MarkdownDescription: "Identity Pool name to be attached with template",
						Description:         "Identity Pool name to be attached with template",
						Optional:            true,
						Computed:            true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"component_id": {
								MarkdownDescription: "Identity Pool name to be attached with template",
								Description:         "Identity Pool name to be attached with template",
								Type:                types.Int64Type,
								Computed:            true,
							},
							"untagged_network": {
								MarkdownDescription: "untagged network to be associated with the specified nic identifier and port",
								Description:         "untagged network to be associated with the specified nic identifier and port",
								Type:                types.Int64Type,
								Optional:            true,
								Computed:            true,
							},
							"tagged_networks": {
								MarkdownDescription: "tagged networks to be associated woith the specified nic identifier and port",
								Description:         "tagged networks to be associated woith the specified nic identifier and port",
								Type: types.ListType{
									ElemType: types.Int64Type,
								},
								Optional: true,
								Computed: true,
							},
							"is_nic_bonded": {
								MarkdownDescription: "Is Nic bonded",
								Description:         "Is Nic bonded",
								Type:                types.BoolType,
								Optional:            true,
								Computed:            true,
							},
							"port": {
								MarkdownDescription: "NIC port",
								Description:         "NIC port",
								Type:                types.Int64Type,
								Optional:            true,
								Computed:            true,
							},
							"nic_identifier": {
								MarkdownDescription: "Display name of NIC port in the template for VLAN configuration",
								Description:         "Display name of NIC port in the template for VLAN configuration",
								Type:                types.StringType,
								Optional:            true,
								Computed:            true,
							},
						}),
					},
				}),
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
func (resourceTemplate) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
}

// Read resource information
func (resourceTemplate) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
}

// Update resource
func (resourceTemplate) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// Delete resource
func (resourceTemplate) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}

// Import resource
func (resourceTemplate) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
