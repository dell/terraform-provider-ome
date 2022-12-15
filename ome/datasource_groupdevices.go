package ome

import (
	"context"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupDevicesDataSourceType struct{}

func (t groupDevicesDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Data source to list the devices in the group from OpenManage Enterprise.",
		Attributes: map[string]tfsdk.Attribute{

			"id": {
				MarkdownDescription: "ID for data source.",
				Description:         "ID for data source.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
			},

			"device_ids": {
				MarkdownDescription: "List of the device id(s) associated with a group",
				Description:         "List of the device id(s) associated with a group",
				Type: types.ListType{
					ElemType: types.Int64Type,
				},
				Computed: true,
			},

			"device_servicetags": {
				MarkdownDescription: "List of the device servicetags associated with a group",
				Description:         "List of the device servicetags associated with a group",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed: true,
			},
			"device_group_names": {
				MarkdownDescription: "List of the device group names.",
				Description:         "List of the device group names.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Required: true,
			},
		},
	}, nil
}

func (t groupDevicesDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return groupDevicesDataSource{
		p: provider,
	}, diags
}

type groupDevicesDataSource struct {
	p provider
}

// Read resource information
func (g groupDevicesDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var groupDevices models.GroupDevicesData
	diags := req.Config.Get(ctx, &groupDevices)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	groupNames := []string{}
	diags = groupDevices.DeviceGroupNames.ElementsAs(ctx, &groupNames, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

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

	devices, err := omeClient.GetDevicesByGroups(groupNames)
	if err != nil && len(devices) != 0 {
		resp.Diagnostics.AddWarning(
			"Unable to fetch devices during pagination: ",
			err.Error(),
		)
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch devices for groups: ",
			err.Error(),
		)
		return
	}
	devIDs := []attr.Value{}
	devSvcTags := []attr.Value{}
	devices = omeClient.GetUniqueDevices(devices)

	for _, device := range devices {
		devIDs = append(devIDs, types.Int64{Value: device.ID})
		devSvcTags = append(devSvcTags, types.String{Value: device.DeviceServiceTag})
	}

	devIDsTfsdk := types.List{
		ElemType: types.Int64Type,
	}
	devIDsTfsdk.Elems = devIDs
	groupDevices.DeviceIDs = devIDsTfsdk

	devSTsTfsdk := types.List{
		ElemType: types.StringType,
	}
	devSTsTfsdk.Elems = devSvcTags
	groupDevices.DeviceServicetags = devSTsTfsdk

	diags = resp.State.Set(ctx, &groupDevices)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
