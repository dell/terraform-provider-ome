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
	"terraform-provider-ome/helper"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &firmwareCatalogDataSource{}
	_ datasource.DataSourceWithConfigure = &firmwareCatalogDataSource{}
)

// NewFirmwareCatalogDataSource creates a new FirmwareCatalogDataSource.
func NewFirmwareCatalogDataSource() datasource.DataSource {
	return &firmwareCatalogDataSource{}
}

type firmwareCatalogDataSource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *firmwareCatalogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*firmwareCatalogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "firmware_catalog"
}

// Schema implements datasource.DataSource
func (*firmwareCatalogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query firmware catalog from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Description: "This Terraform DataSource is used to query firmware catalog from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Attributes: omeFirmwareCatalogDataSchema(),
	}
}

// Read implements datasource.DataSource
func (g *firmwareCatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan models.OMECatalogData
	diag := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, d := g.p.createOMESession(ctx, "datasource_firmware_catalog Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	cat, err := helper.GetAllCatalogFirmware(omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching firmware catalogs",
			err.Error(),
		)
		return
	}
	var filteredNames []string
	diags := plan.Names.ElementsAs(ctx, &filteredNames, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	vals, filterErr := helper.FilterCatalogFirmware(ctx, filteredNames, cat)
	if filterErr != nil {
		resp.Diagnostics.AddError(
			"Error processing firmware catalogs",
			filterErr.Error(),
		)
		return
	}

	if plan.ID.IsNull() {
		plan.ID = types.Int64Value(0)
	}
	plan.Catalog = vals

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
