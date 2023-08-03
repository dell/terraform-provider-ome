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
	"fmt"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
		MarkdownDescription: "Data source to list devices from OpenManage Enterprise.",
		Attributes:          OmeDeviceDataSchema(),
	}
}

// Read implements datasource.DataSource
func (g *deviceDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan models.OmeDeviceData
	diags := req.Config.Get(ctx, &plan)
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

	id := plan.ID.ValueInt64()

	devs, err := omeClient.GetDevices(nil, []int64{id}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error device by id: %d", id),
			err.Error(),
		)
		return
	}
	inv, err2 := omeClient.GetDeviceInventory(id)
	if err2 != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting inv by id: %d", id),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprint(inv))
	state := plan
	vals := make([]models.OmeSingleDeviceData, 0)
	for _, dev := range devs {
		val := models.NewSingleOmeDeviceData(dev)
		val.Inventory = models.NewOmeDeviceInventory(inv)
		vals = append(vals, val)
	}
	state.Devices = vals

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
