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
	"errors"
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &resourceDevices{}
	_ resource.ResourceWithConfigure = &resourceDevices{}
)

// NewDevicesResource is new resource for application csr
func NewDevicesResource() resource.Resource {
	return &resourceDevices{}
}

type resourceDevices struct {
	p *omeProvider
	c *clients.Client
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceDevices) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (r resourceDevices) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "devices"
}

// Devices Resource schema
func (r resourceDevices) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing devices on OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID for devices resource.",
				Description:         "ID for devices resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"devices": schema.ListNestedAttribute{
				MarkdownDescription: "List of devices to be managed by this resource.",
				Description:         "List of devices to be managed by this resource.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.singleDeviceSchema(),
				},
			},
		},
	}
}

func (*resourceDevices) singleDeviceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of device.",
			Description:         "ID of device.",
			Optional:            true,
			Computed:            true,
		},
		"service_tag": schema.StringAttribute{
			MarkdownDescription: "Service tag of device.",
			Description:         "Service tag of device.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("id")),
			},
		},
		"management_ips": schema.ListAttribute{
			MarkdownDescription: "List of management IPs of device.",
			Description:         "List of management IPs of device.",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// Create a new resource
func (r resourceDevices) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	//Get Plan Data
	tflog.Trace(ctx, "resource_devices create: started")
	var plan models.DevicesResModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_devices Configure")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	r.c = omeClient
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_devices getting current infrastructure state")

	state, err := r.getState(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting current infrastructure state.",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDevices) getState(ctx context.Context, plan models.DevicesResModel) (models.DevicesResModel, error) {
	var (
		state models.DevicesResModel
		err   error
		devs  []models.Device
	)

	if plan.Devices.IsUnknown() {
		devM, err2 := r.c.GetAllDevices(nil)
		devs = devM.Value
		err = err2
	} else {
		pdevs := make([]models.DeviceItemModel, 0)
		diags := plan.Devices.ElementsAs(ctx, &pdevs, false)
		if diags.HasError() {
			err := fmt.Errorf("")
			for _, d := range diags.Errors() {
				err = fmt.Errorf("%w, {%s,%s}", err, d.Summary(), d.Detail())
			}
			return state, err
		}
		for _, v := range pdevs {
			var (
				dev  models.Device
				err2 = fmt.Errorf("neither id nor service tag provided for device")
			)
			if !v.ID.IsUnknown() {
				dev, err2 = r.c.GetDevice("", v.ID.ValueInt64())
			} else if !v.ServiceTag.IsUnknown() {
				dev, err2 = r.c.GetDevice(v.ServiceTag.ValueString(), 0)
			}
			if err2 != nil {
				if !v.ID.IsUnknown() && !v.ServiceTag.IsUnknown() && errors.Is(err2, clients.ErrItemNotFound) {
					// this condition means that this device is not found during read refresh
					continue
				}
				err = r.join(err, err2)
				continue
			}
			devs = append(devs, dev)
		}
	}
	if err != nil {
		return state, err
	}

	// for _, v := range devs {
	// 	state.Devices = append(state.Devices, models.NewDeviceItemModel(v))
	// }
	state = models.NewDevicesResModel(devs)
	return state, nil
}

func (r *resourceDevices) join(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	return fmt.Errorf("%w, %s", err1, err2.Error())
}

// Read resource information
func (r resourceDevices) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	//Get Plan Data
	tflog.Trace(ctx, "resource_devices create: started")
	var state models.DevicesResModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_devices Configure")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	r.c = omeClient
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_devices getting current infrastructure state")

	state, err := r.getState(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting current infrastructure state.",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceDevices) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_devices update: started")
	var plan models.DevicesResModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_devices Configure")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	r.c = omeClient
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_devices getting current infrastructure state")

	state, err := r.getState(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting current infrastructure state.",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r resourceDevices) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Just remove State Data
	resp.State.RemoveResource(ctx)
}
