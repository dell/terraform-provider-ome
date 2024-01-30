/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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
	"strconv"
	"terraform-provider-ome/helper"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &fwBaselineCompReportDatasource{}
	_ datasource.DataSourceWithConfigure = &fwBaselineCompReportDatasource{}
)

// NewfwBaselineCompReportDatasource is new datasource for firmware baseline compliance report
func NewfwBaselineCompReportDatasource() datasource.DataSource {
	tflog.Info(context.TODO(), "Initializing fwBaselineCompReportDatasource")
	return &fwBaselineCompReportDatasource{}
}

type fwBaselineCompReportDatasource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *fwBaselineCompReportDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*fwBaselineCompReportDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "fw_baseline_compliance_report_info"
}

// Schema implements datasource.DataSource
func (*fwBaselineCompReportDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query firmware baseline compliance report of OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Description: "This Terraform DataSource is used to query firmware baseline compliance report from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Attributes: FwBaseComplianceReportDataSchema(),
		Blocks:     FwBaseComplianceReportFilterBlockSchema(),
	}
}

// Read reads the firmware baseline compliance report from the OME API and
// populates the datasource with the retrieved data.
//
// ctx: The context.Context object for the function.
// req: The datasource.ReadRequest object for the function.
// resp: The datasource.ReadResponse object for the function.
func (g *fwBaselineCompReportDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var plan models.OmeFwComplianceReportData //models.OmeFwComplianceReportData //
	diag := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, d := g.p.createOMESession(ctx, "datasource_fw_compliance_report Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		tflog.Debug(ctx, strconv.Itoa(d.ErrorsCount()))
		return
	}
	defer omeClient.RemoveSession()
	baselineID, err := omeClient.GetUpdateServiceBaselineIDByName(plan.BaseLineName.ValueString())
	if err != nil || baselineID == -1 {
		resp.Diagnostics.AddError(
			"Error fetching baseline", err.Error(),
		)
		return
	}
	if plan.CrFilter == nil {
		plan.CrFilter = &models.OmeFwComplianceReportFilter{
			Key:   types.StringValue(""),
			Value: types.StringValue(""),
		}
	}
	if plan.CrFilter.Key.IsNull() {
		plan.CrFilter.Key = types.StringValue("")
	}
	if plan.CrFilter.Value.IsNull() {
		plan.CrFilter.Value = types.StringValue("")
	}
	filterKey := plan.CrFilter.Key.ValueString()
	filterVal := plan.CrFilter.Value.ValueString()
	tflog.Debug(ctx, "got map value", map[string]interface{}{
		"key":   filterKey,
		"value": filterVal,
	})
	report, err := helper.GetFwBaselineComplianceReport(ctx, omeClient, baselineID, filterKey, filterVal)
	if err != nil || report == nil {
		resp.Diagnostics.AddError(
			"Error fetching firmware baseline compliance report",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "got report", map[string]interface{}{
		"report": report.Value[0].DeviceName,
	})
	if plan.ID.IsNull() {
		plan.ID = types.Int64Value(1)
	}

	plan.Report = helper.NewOmeFwComplianceReportList(report.Value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
