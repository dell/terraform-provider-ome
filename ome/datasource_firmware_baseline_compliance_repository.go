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
	"terraform-provider-ome/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &firmwareBaselineComplianceRepositoryDatasource{}
	_ datasource.DataSourceWithConfigure = &firmwareBaselineComplianceRepositoryDatasource{}
)

// NewFirmwareBaselineComplianceRepositoryDatasource is new datasource for Repository FBC report
func NewFirmwareBaselineComplianceRepositoryDatasource() datasource.DataSource {
	return &firmwareBaselineComplianceRepositoryDatasource{}
}

type firmwareBaselineComplianceRepositoryDatasource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *firmwareBaselineComplianceRepositoryDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*firmwareBaselineComplianceRepositoryDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "firmware_repository"
}

// Schema implements datasource.DataSource
func (*firmwareBaselineComplianceRepositoryDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query the firmware baseline compliance repository from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Attributes: omeFirmwareBaselineComplianceRepositoryDataSchema(),
	}
}

// Read implements datasource.DataSource
func (g *firmwareBaselineComplianceRepositoryDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan models.OMERepositoryData
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, d := g.p.createOMESession(ctx, "datasource_firmware_repository Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	repositories, errGet := helper.GetAllRepositories(omeClient)
	if errGet != nil {
		resp.Diagnostics.AddError(
			"Error Reading Repositories",
			errGet.Error(),
		)
		return
	}
	var repos []models.RepositoryModel
	var filterErr error

	if len(plan.Names) == 0 {
		// get all Repositories
		repos = repositories
	} else {
		// get filtered Repositories
		repos, filterErr = helper.GetFilteredRepositoriesByName(ctx, repositories, plan)
	}

	if filterErr != nil {
		resp.Diagnostics.AddError(
			"Error Filtering Repositories",
			filterErr.Error(),
		)
		return
	}

	vals := make([]models.CatalogRepository, 0)
	for _, repo := range repos {
		val := models.CatalogRepository{}
		err := utils.CopyFields(ctx, repo, &val)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Copying values for repositories",
				err.Error(),
			)
			return
		}
		vals = append(vals, val)
	}

	if plan.ID.IsNull() {
		plan.ID = types.Int64Value(0)
	}
	plan.Repositories = vals

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
