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
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &deviceDatasource{}
	_ datasource.DataSourceWithConfigure = &deviceDatasource{}
)

// NewDeviceDatasource is new datasource for group devices
func NewDeviceDatasource() datasource.DataSource {
	return &deviceDatasource{}
}

type deviceDatasource struct {
	p *omeProvider
}

// Configure implements datasource.DataSourceWithConfigure
func (g *deviceDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	g.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*deviceDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "device"
}

// Schema implements datasource.DataSource
func (*deviceDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query devices from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Attributes: omeDeviceDataSchema(),
	}
}

// Read implements datasource.DataSource
func (g *deviceDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan models.OmeDeviceData
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var filters models.OmeDeviceDataFilters
	diags = plan.Filters.As(ctx, &filters, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty: true,
	})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, d := g.p.createOMESession(ctx, "datasource_device Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	devs, err := g.ReadDevices(ctx, omeClient, filters)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching devices",
			err.Error(),
		)
		return
	}

	vals := make([]models.OmeSingleDeviceData, 0)
	for _, dev := range devs {
		val := models.NewSingleOmeDeviceData(dev)
		vals = append(vals, val)
	}

	// If at least one of the filters are set then do detailed inventory
	if !filters.FilterExpr.IsNull() ||
		len(filters.IDs.Elements()) > 0 ||
		len(filters.SvcTags.Elements()) > 0 ||
		len(filters.IPExprs.Elements()) > 0 {
		for i, dev := range devs {
			id := dev.ID
			inv, err2 := g.ReadDeviceInventory(ctx, omeClient, id, plan.InventoryTypes)
			if err2 != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error getting detailed inventory by id: %d", id),
					err2.Error(),
				)
				return
			}
			tflog.Info(ctx, fmt.Sprint(inv))
			vals[i].Inventory = inv
		}
	}
	g.WriteState(ctx, plan, vals, resp)
}

// Read implements datasource.DataSource
func (g *deviceDatasource) WriteState(ctx context.Context, plan models.OmeDeviceData,
	vals []models.OmeSingleDeviceData, resp *datasource.ReadResponse) {

	// needed for acceptance testing - setting the id and then writing plan to state
	if plan.ID.IsNull() {
		plan.ID = types.Int64Value(0)
	}
	plan.Devices = vals
	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read implements datasource.DataSource
func (g *deviceDatasource) ReadDevices(ctx context.Context, client *clients.Client,
	filters models.OmeDeviceDataFilters) ([]models.Device, error) {

	var (
		err error
		ret []models.Device
	)
	if !filters.FilterExpr.IsNull() {
		devs, err2 := client.GetAllDevices(map[string]string{
			"$filter": filters.FilterExpr.ValueString(),
		})
		ret, err = devs.Value, err2
	} else if !filters.IDs.IsNull() {
		inputs := make([]int64, 0)
		_ = filters.IDs.ElementsAs(ctx, &inputs, false)
		ret, err = client.GetDevices(nil, inputs, nil)
	} else if !filters.SvcTags.IsNull() {
		inputs := make([]string, 0)
		_ = filters.SvcTags.ElementsAs(ctx, &inputs, false)
		ret, err = client.GetDevices(inputs, nil, nil)
	} else {
		devs, err2 := client.GetAllDevices(nil)
		ret, err = devs.Value, err2
	}

	if err != nil {
		return nil, err
	}

	if !filters.IPExprs.IsNull() {
		inputs := make([]string, 0)
		_ = filters.IPExprs.ElementsAs(ctx, &inputs, false)
		return clients.FilterDeviceByIps(ret, inputs)
	}
	return ret, err
}

// Read implements datasource.DataSource
func (g *deviceDatasource) ReadDeviceInventory(ctx context.Context, client *clients.Client,
	id int64, itypes []string) (models.OmeDeviceInventory, error) {

	var err error
	retv := models.NewDeviceInventory()

	if itypes == nil {
		retv, err = client.GetDeviceInventory(id)
	} else {
		for _, t := range itypes {
			reta, err2 := client.GetDeviceInventoryByType(id, t)
			if err2 != nil {
				return models.NewOmeDeviceInventory(retv), err2
			}
			retv.AddInventory(reta)
		}
	}

	return models.NewOmeDeviceInventory(retv), err
}
