package ome

import (
	"context"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &groupDevicesDatasource{}
	_ datasource.DataSourceWithConfigure = &groupDevicesDatasource{}
)

// NewGroupDevicesDatasource is new datasource for group devices
func NewGroupDevicesDatasource() datasource.DataSource {
	return &groupDevicesDatasource{}
}

type groupDevicesDatasource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *groupDevicesDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*groupDevicesDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "groupdevices_info"
}

// Schema implements datasource.DataSource
func (*groupDevicesDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to list the devices in the group from OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "ID for group devices data source.",
				Description:         "ID for group devices data source.",
				Computed:            true,
				Optional:            true,
			},

			"device_ids": schema.ListAttribute{
				MarkdownDescription: "List of the device id(s) associated with a group",
				Description:         "List of the device id(s) associated with a group",
				ElementType:         types.Int64Type,
				Computed:            true,
			},

			"device_servicetags": schema.ListAttribute{
				MarkdownDescription: "List of the device servicetags associated with a group",
				Description:         "List of the device servicetags associated with a group",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"device_group_names": schema.SetAttribute{
				MarkdownDescription: "List of the device group names.",
				Description:         "List of the device group names.",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
}

// Read implements datasource.DataSource
func (g *groupDevicesDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var groupDevices models.GroupDevicesData
	diags := req.Config.Get(ctx, &groupDevices)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	groupNames := []string{}
	resp.Diagnostics.Append(groupDevices.DeviceGroupNames.ElementsAs(ctx, &groupNames, true)...)
	if diags.HasError() {
		return
	}

	omeClient, d := g.p.createOMESession(ctx, "datasource_vlannetworks_info Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	devices, err := omeClient.GetDevicesByGroups(groupNames)
	devIDs := []attr.Value{}
	devSvcTags := []attr.Value{}

	if err != nil {
		if len(devices) != 0 {
			resp.Diagnostics.AddWarning(
				"Unable to fetch devices during pagination: ",
				err.Error(),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to fetch devices for groups: ",
				err.Error(),
			)
			return
		}
	}
	devices = omeClient.GetUniqueDevices(devices)
	for _, device := range devices {
		devIDs = append(devIDs, types.Int64Value(device.ID))
		devSvcTags = append(devSvcTags, types.StringValue(device.DeviceServiceTag))
	}

	devIDsTfsdk, _ := types.ListValue(
		types.Int64Type,
		devIDs,
	)

	groupDevices.DeviceIDs = devIDsTfsdk

	devSTsTfsdk, _ := types.ListValue(
		types.StringType,
		devSvcTags,
	)

	groupDevices.DeviceServicetags = devSTsTfsdk

	diags = resp.State.Set(ctx, &groupDevices)
	resp.Diagnostics.Append(diags...)
}
