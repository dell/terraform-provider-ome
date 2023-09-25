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
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// NewDevicesResource is new resource for devices
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
		MarkdownDescription: "This terraform resource is used to manage a set of device entities on OME.",
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

	state, dgs := r.getState(ctx, plan.Devices)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDevices) getState(ctx context.Context, tfDevices types.List) (
	models.DevicesResModel, diag.Diagnostics) {
	var (
		state models.DevicesResModel
		devs  []models.Device
		pdevs []models.DeviceItemModel
	)

	if !tfDevices.IsUnknown() {
		pdevs = make([]models.DeviceItemModel, 0)
		diags := tfDevices.ElementsAs(ctx, &pdevs, false)
		if diags.HasError() {
			return state, diags
		}
	}

	var dgs diag.Diagnostics

	if pdevs == nil {
		devM, err := r.c.GetAllDevices(nil)
		devs = devM.Value
		if err != nil {
			dgs.AddError(
				"Error getting devices.",
				err.Error(),
			)
		}
	} else {
		for _, v := range pdevs {
			var (
				dev models.Device
				err = fmt.Errorf("neither id nor service tag provided for device")
			)
			if !v.ID.IsUnknown() {
				dev, err = r.c.GetDevice("", v.ID.ValueInt64())
			} else if !v.ServiceTag.IsUnknown() {
				dev, err = r.c.GetDevice(v.ServiceTag.ValueString(), 0)
			}
			if err != nil {
				if !v.ID.IsUnknown() && !v.ServiceTag.IsUnknown() && errors.Is(err, clients.ErrItemNotFound) {
					// this condition means that this device is not found during read refresh
					continue
				}
				dgs.AddError(
					"Error getting device.",
					err.Error(),
				)
				continue
			}
			devs = append(devs, dev)
		}
	}
	if dgs.HasError() {
		return state, dgs
	}
	return models.NewDevicesResModel(devs)
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

	state, dgs := r.getState(ctx, state.Devices)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceDevices) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_devices update: started")
	var oldState, plan models.DevicesResModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &oldState)
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

	resp.Diagnostics.Append(r.updateDevs(ctx, oldState, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, dgs := r.getState(ctx, plan.Devices)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDevices) getDevIDsToRmv(ctx context.Context, state, plan models.DevicesResModel) (
	[]int64, diag.Diagnostics) {
	var (
		dgs      diag.Diagnostics
		idsToRmv = make([]int64, 0)
	)

	if plan.Devices.IsUnknown() {
		return idsToRmv, dgs
	}
	pdevs := make([]models.DeviceItemModel, 0)
	dgs.Append(plan.Devices.ElementsAs(ctx, &pdevs, false)...)
	sdevs := make([]models.DeviceItemModel, 0)
	dgs.Append(state.Devices.ElementsAs(ctx, &sdevs, false)...)
	if dgs.HasError() {
		return idsToRmv, dgs
	}
	mid, mstag := make(map[int64]bool), make(map[string]bool)
	for _, v := range pdevs {
		if !v.ID.IsUnknown() {
			mid[v.ID.ValueInt64()] = true
		}
		if !v.ServiceTag.IsUnknown() {
			mstag[v.ServiceTag.ValueString()] = true
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("Device IDs present in plan: %v", mid))
	tflog.Debug(ctx, fmt.Sprintf("Device ServiceTags present in plan: %v", mstag))
	for _, v := range sdevs {
		_, oki := mid[v.ID.ValueInt64()]
		_, okt := mstag[v.ServiceTag.ValueString()]
		if !(oki || okt) {
			// this condition means that device v must be removed
			idsToRmv = append(idsToRmv, v.ID.ValueInt64())
		}
	}
	return idsToRmv, dgs
}

func (r resourceDevices) updateDevs(ctx context.Context, state, plan models.DevicesResModel) diag.Diagnostics {
	idsToRmv, dgs := r.getDevIDsToRmv(ctx, state, plan)
	tflog.Info(ctx, fmt.Sprintf("resource_devices removing devices with IDs %v", idsToRmv))
	err := r.c.RemoveDevices(idsToRmv)
	if err != nil {
		dgs.AddError("Could not remove devices", err.Error())
	}
	return dgs
}

// Delete resource
func (r resourceDevices) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Just remove State Data
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceDevices) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "resource_devices import: started")
	planDevs, dgs := models.NewDevicesResModelFromID(req.ID)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// defer the remove session
	omeClient, ds := r.p.createOMESession(ctx, "resource_devices Import")
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	r.c = omeClient
	defer omeClient.RemoveSession()

	state, dgs := r.getState(ctx, planDevs)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_devices import: finished")
}
