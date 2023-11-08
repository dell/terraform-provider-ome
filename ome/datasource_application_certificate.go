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
	_ datasource.DataSource              = &appCertDatasource{}
	_ datasource.DataSourceWithConfigure = &appCertDatasource{}
)

// NewAppCertDataSource is new datasource for application certificate
func NewAppCertDataSource() datasource.DataSource {
	return &appCertDatasource{}
}

type appCertDatasource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *appCertDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*appCertDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "application_certificate"
}

// Schema implements datasource.DataSource
func (g *appCertDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This terraform DataSource is used to query the existing application certificate data from OME." +
			" The information fetched from this data source can be used for getting the details for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID for application certificate data source.",
				Description:         "ID for application certificate data source.",
				Computed:            true,
			},
			"issued_to": schema.SingleNestedAttribute{
				MarkdownDescription: "List of the device id(s) associated with any of the groups.",
				Description:         "List of the device id(s) associated with any of the groups.",
				Attributes:          g.infoSchema(),
				Computed:            true,
			},
			"issued_by": schema.SingleNestedAttribute{
				MarkdownDescription: "List of the device servicetags associated with any of the groups.",
				Description:         "List of the device servicetags associated with any of the groups.",
				Attributes:          g.infoSchema(),
				Computed:            true,
			},
			"valid_to": schema.StringAttribute{
				MarkdownDescription: "List of the device group names.",
				Description:         "List of the device group names.",
				Computed:            true,
			},
			"valid_from": schema.StringAttribute{
				MarkdownDescription: "Map of the groups fetched keyed by its name.",
				Description:         "Map of the groups fetched keyed by its name.",
				Computed:            true,
			},
		},
	}
}

func (*appCertDatasource) infoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"distinguished_name": schema.StringAttribute{
			MarkdownDescription: "Distinguished Name.",
			Description:         "Distinguished Name.",
			Computed:            true,
			Optional:            true,
		},
		"department_name": schema.StringAttribute{
			MarkdownDescription: "Department Name.",
			Description:         "Department Name.",
			Computed:            true,
		},
		"business_name": schema.StringAttribute{
			MarkdownDescription: "Business Name.",
			Description:         "Business Name.",
			Computed:            true,
		},
		"locality": schema.StringAttribute{
			MarkdownDescription: "Locality of the business.",
			Description:         "Locality of the business.",
			Computed:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "State of the business.",
			Description:         "State of the business.",
			Computed:            true,
		},
		"country": schema.StringAttribute{
			MarkdownDescription: "Country of the business.",
			Description:         "Country of the business.",
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "Email address.",
			Description:         "Email address.",
			Computed:            true,
		},
		"subject_alternate_names": schema.ListAttribute{
			MarkdownDescription: "Subject Alternate names.",
			Description:         "Subject Alternate names.",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// Read implements datasource.DataSource
func (g *appCertDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	omeClient, d := g.p.createOMESession(ctx, "datasource_vlannetworks_info Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	info, err := omeClient.GetCert()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch certificate information.",
			err.Error(),
		)
	}

	appCert := models.NewCertInfoModel(info)
	resp.Diagnostics.Append(resp.State.Set(ctx, &appCert)...)
}
