package ome

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
)

type configurationReportDataSourceType struct{}

func (t configurationReportDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Data source to list the compliance configuration report of a baseline from OpenManage Enterprise.",
		Attributes: map[string]tfsdk.Attribute{

			"id": {
				MarkdownDescription: "ID for data source.",
				Description:         "ID for data source.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
			},
			"baseline_name": {
				MarkdownDescription: "Name of the Baseline.",
				Description:         "Name of the Baseline.",
				Type:                types.StringType,
				Required:            true,
			},
			"fetch_attributes": {
				MarkdownDescription: "Fetch  device compliance attribute report.",
				Description:         "Fetch  device compliance attribute report.",
				Type:                types.BoolType,
				Computed:            true,
				Optional:            true,
			},
			"compliance_report_device": {
				MarkdownDescription: "Device complaince report.",
				Description:         "Device complaince report.",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"device_id": {
						MarkdownDescription: "Device ID",
						Description:         "Device ID",
						Type:                types.Int64Type,
						Computed:            true,
					},
					"device_servicetag": {
						MarkdownDescription: "Device servicetag.",
						Description:         "Device servicetag.",
						Type:                types.StringType,
						Computed:            true,
					},
					"device_name": {
						MarkdownDescription: "Device Name.",
						Description:         "Device Name.",
						Type:                types.StringType,
						Computed:            true,
					},
					"model": {
						MarkdownDescription: "Device model.",
						Description:         "Device model.",
						Type:                types.StringType,
						Computed:            true,
					},
					"compliance_status": {
						MarkdownDescription: "Device compliance status.",
						Description:         "Device compliance status.",
						Type:                types.StringType,
						Computed:            true,
					},
					"device_type": {
						MarkdownDescription: "Device type",
						Description:         "Device type",
						Type:                types.Int64Type,
						Computed:            true,
					},
					"inventory_time": {
						MarkdownDescription: "Inventory Time.",
						Description:         "Inventory Time.",
						Type:                types.StringType,
						Computed:            true,
					},
					"device_compliance_details": {
						MarkdownDescription: "Device compliance details.",
						Description:         "Device compliance details.",
						Type:                types.StringType,
						Computed:            true,
					},
				}),
				Computed: true,
			},
		},
	}, nil
}

func (t configurationReportDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return configurationReportDataSource{
		p: provider,
	}, diags
}

type configurationReportDataSource struct {
	p provider
}

// Read resource information
func (g configurationReportDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config models.ConfigurationReports
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, err := clients.NewClient(*g.p.clientOpt)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return
	}

	_, err = omeClient.CreateSession()
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return
	}
	defer omeClient.RemoveSession()

	var state models.ConfigurationReports

	baseline, err := omeClient.GetBaselineByName(config.BaseLineName.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrConfigurationReport, err.Error(),
		)
		return
	}
	if baseline.ConfigComplianceSummary.ComplianceStatus != "NotInventored" {

		baselineId := baseline.ID

		complianceReports, err := omeClient.GetBaselineDevComplianceReportsByID(baselineId)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrConfigurationReport, err.Error(),
			)
			return
		}
		state.BaseLineName.Value = config.BaseLineName.Value
		state.ID.Value = config.BaseLineName.Value

		for _, cr := range complianceReports {
			compStatus := "Compliant"
			if cr.ComplianceStatus != 1 {
				compStatus = "Non Compliant"
			}
			crd := models.ComplianceReportDevice{
				DeviceID:         types.Int64{Value: cr.ID},
				DeviceServiceTag: types.String{Value: cr.ServiceTag},
				DeviceName:       types.String{Value: cr.DeviceName},
				Model:            types.String{Value: cr.Model},
				ComplianceStatus: types.String{Value: compStatus},
				DeviceType:       types.Int64{Value: cr.DeviceType},
				InventoryTime:    types.String{Value: cr.InventoryTime},
			}
			if config.FetchAttributes.Value {
				attrResp, err := omeClient.GetBaselineDevAttrComplianceReportsByID(baselineId, cr.ID)
				if err != nil {
					resp.Diagnostics.AddError(
						clients.ErrGnrConfigurationReport, err.Error(),
					)
					return
				}
				crd.DeviceComplianceDetails = types.String{Value: attrResp}
			}
			state.ComplianceReportDevice = append(state.ComplianceReportDevice, crd)
		}
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
