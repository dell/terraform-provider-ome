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
	_ datasource.DataSource              = &deviceComplianceReportDatasource{}
	_ datasource.DataSourceWithConfigure = &deviceComplianceReportDatasource{}
)

// NewDeviceComplianceReportDataSource - Creates a new device compliance report datasource.
func NewDeviceComplianceReportDataSource() datasource.DataSource {
	return &deviceComplianceReportDatasource{}
}

type deviceComplianceReportDatasource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *deviceComplianceReportDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*deviceComplianceReportDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "device_compliance_report"
}

// Schema implements datasource.DataSource
func (*deviceComplianceReportDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query compliance configuration report of a compliance template baseline data from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Description: "This Terraform DataSource is used to query device compliance report of a compliance template baseline data from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Attributes: omeDeviceReportSchema(),
	}
}

// Read implements datasource.DataSource
func (g *deviceComplianceReportDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan models.OMEDeviceComplianceData
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, d := g.p.createOMESession(ctx, "device_compliance_report Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	baselineID, err := omeClient.GetUpdateServiceBaselineIDByName(plan.BaselineName.ValueString())
	if err != nil || baselineID == -1 {
		resp.Diagnostics.AddError(
			"Error fetching baseline", err.Error(),
		)
		return
	}

	complianceReports, err := helper.GetAllDeviceComplianceReport(omeClient, baselineID)

	if err != nil {
		resp.Diagnostics.AddError("Error reading device compliance report", err.Error())
		return
	}

	vals, stateErr := helper.SetStateDeviceComplianceReport(ctx, complianceReports)
	if stateErr != nil {
		resp.Diagnostics.AddError(
			"Error processing device compliance report",
			stateErr.Error(),
		)
		return
	}

	if plan.ID.IsNull() {
		plan.ID = types.Int64Value(0)
	}
	plan.Reports = vals

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
