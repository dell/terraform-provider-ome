package ome

import (
	"context"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &configurationReportDataSource{}
	_ datasource.DataSourceWithConfigure = &configurationReportDataSource{}
)

// NewConfigurationReportDataSource is a new datasource for configuration report
func NewConfigurationReportDataSource() datasource.DataSource {
	return &configurationReportDataSource{}
}

type configurationReportDataSource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *configurationReportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*configurationReportDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "configuration_report_info"
}

func (*configurationReportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to list the compliance configuration report of a baseline from OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "ID for data source.",
				Description:         "ID for data source.",
				Computed:            true,
				Optional:            true,
			},
			"baseline_name": schema.StringAttribute{
				MarkdownDescription: "Name of the Baseline.",
				Description:         "Name of the Baseline.",
				Required:            true,
			},
			"fetch_attributes": schema.BoolAttribute{
				MarkdownDescription: "Fetch  device compliance attribute report.",
				Description:         "Fetch  device compliance attribute report.",
				Computed:            true,
				Optional:            true,
			},
			"compliance_report_device": schema.ListNestedAttribute{
				MarkdownDescription: "Device complaince report.",
				Description:         "Device complaince report.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.Int64Attribute{
							MarkdownDescription: "Device ID",
							Description:         "Device ID",
							Computed:            true,
						},
						"device_servicetag": schema.StringAttribute{
							MarkdownDescription: "Device servicetag.",
							Description:         "Device servicetag.",
							Computed:            true,
						},
						"device_name": schema.StringAttribute{
							MarkdownDescription: "Device Name.",
							Description:         "Device Name.",
							Computed:            true,
						},
						"model": schema.StringAttribute{
							MarkdownDescription: "Device model.",
							Description:         "Device model.",
							Computed:            true,
						},
						"compliance_status": schema.StringAttribute{
							MarkdownDescription: "Device compliance status.",
							Description:         "Device compliance status.",
							Computed:            true,
						},
						"device_type": schema.Int64Attribute{
							MarkdownDescription: "Device type",
							Description:         "Device type",
							Computed:            true,
						},
						"inventory_time": schema.StringAttribute{
							MarkdownDescription: "Inventory Time.",
							Description:         "Inventory Time.",
							Computed:            true,
						},
						"device_compliance_details": schema.StringAttribute{
							MarkdownDescription: "Device compliance details.",
							Description:         "Device compliance details.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Read resource information
func (g *configurationReportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	baseline, err := omeClient.GetBaselineByName(config.BaseLineName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrConfigurationReport, err.Error(),
		)
		return
	}
	if baseline.ConfigComplianceSummary.ComplianceStatus != "NotInventored" {

		baselineID := baseline.ID

		complianceReports, err := omeClient.GetBaselineDevComplianceReportsByID(baselineID)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrConfigurationReport, err.Error(),
			)
			return
		}
		if !config.BaseLineName.IsUnknown() {
			state.BaseLineName = config.BaseLineName
		}
		if !config.ID.IsUnknown() {
			state.ID = config.BaseLineName
		}

		for _, cr := range complianceReports {
			compStatus := "Compliant"
			if cr.ComplianceStatus != 1 {
				compStatus = "Non Compliant"
			}
			crd := models.ComplianceReportDevice{
				DeviceID:         types.Int64Value(cr.ID),
				DeviceServiceTag: types.StringValue(cr.ServiceTag),
				DeviceName:       types.StringValue(cr.DeviceName),
				Model:            types.StringValue(cr.Model),
				ComplianceStatus: types.StringValue(compStatus),
				DeviceType:       types.Int64Value(cr.DeviceType),
				InventoryTime:    types.StringValue(cr.InventoryTime),
			}
			if config.FetchAttributes.ValueBool() {
				attrResp, err := omeClient.GetBaselineDevAttrComplianceReportsByID(baselineID, cr.ID)
				if err != nil {
					resp.Diagnostics.AddError(
						clients.ErrGnrConfigurationReport, err.Error(),
					)
					return
				}
				crd.DeviceComplianceDetails = types.StringValue(attrResp)
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
