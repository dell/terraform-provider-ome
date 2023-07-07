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

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Groups - list of group response from OME
type Groups struct {
	Value []Group `json:"value"`
}

// Group - embedded group response from the groups
type Group struct {
	ID               int64  `json:"Id,omitempty"`
	Name             string `json:"Name"`
	Description      string `json:"Description"`
	MembershipTypeId int64  `json:"MembershipTypeId"`
	ParentId         int64  `json:"ParentId"`
}

type GroupMemberPayload struct {
	GroupId   int64   `json:"GroupId"`
	DeviceIds []int64 `json:"MemberDeviceIds,omitempty"`
}

func NewGroupMemberPayload(gid int64) GroupMemberPayload {
	return GroupMemberPayload{
		GroupId:   gid,
		DeviceIds: make([]int64, 0),
	}
}

func (g *GroupMemberPayload) RegisterDevice(id int64) {
	g.DeviceIds = append(g.DeviceIds, id)
}

// GroupDevicesData - schema for data source groupdevices
type GroupDevicesData struct {
	ID                types.String `tfsdk:"id"`
	DeviceIDs         types.List   `tfsdk:"device_ids"`
	DeviceServicetags types.List   `tfsdk:"device_servicetags"`
	DeviceGroupNames  types.Set    `tfsdk:"device_group_names"`
}

// StaticGroup - schema for resource static group
type StaticGroup struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	MembershipTypeId types.Int64  `tfsdk:"membership_type_id"`
	ParentId         types.Int64  `tfsdk:"parent_id"`
	DeviceIds        types.Set    `tfsdk:"device_ids"`
}

func NewStaticGroup(g Group, devs Devices) (StaticGroup, diag.Diagnostics) {
	devidVals := make([]attr.Value, 0)
	for _, device := range devs.Value {
		devidVals = append(devidVals, types.Int64Value(device.ID))
	}
	deviceIds, dgs := types.SetValue(types.Int64Type, devidVals)
	return StaticGroup{
		ID:               types.Int64Value(g.ID),
		Name:             types.StringValue(g.Name),
		Description:      types.StringValue(g.Description),
		MembershipTypeId: types.Int64Value(g.MembershipTypeId),
		ParentId:         types.Int64Value(g.ParentId),
		DeviceIds:        deviceIds,
	}, dgs
}

func (g *StaticGroup) GetPayload(state StaticGroup) (Group, bool) {
	ret := Group{
		ID:               state.ID.ValueInt64(),
		Name:             g.Name.ValueString(),
		Description:      g.Description.ValueString(),
		MembershipTypeId: 12,
		ParentId:         g.ParentId.ValueInt64(),
	}
	return ret, ret.Name == state.Name.ValueString() && ret.Description == state.Description.ValueString()
}

func (plan *StaticGroup) GetMemberPayload(ctx context.Context, state StaticGroup) (GroupMemberPayload,
	GroupMemberPayload, diag.Diagnostics) {
	var d diag.Diagnostics
	toAdd, toRmv := NewGroupMemberPayload(state.ID.ValueInt64()), NewGroupMemberPayload(state.ID.ValueInt64())
	if plan.DeviceIds.Equal(state.DeviceIds) {
		return toAdd, toRmv, d
	}
	planDevIds, dgs1 := plan.GetDeviceIdMap(ctx)
	d.Append(dgs1...)
	stateDevIds, dgs2 := state.GetDeviceIdMap(ctx)
	d.Append(dgs2...)

	// Loop over all devices in state
	for sid := range stateDevIds {
		if _, ok := planDevIds[sid]; !ok {
			// Register all devices to remove
			toRmv.RegisterDevice(sid)
		} else {
			// Flag all devices that are already in state
			planDevIds[sid] = false
		}
	}

	// Register all devices to add (ie, devices not flagged)
	for pid, toAddFlag := range planDevIds {
		if toAddFlag {
			toAdd.RegisterDevice(pid)
		}
	}

	return toAdd, toRmv, d
}

func (g *StaticGroup) GetDeviceIdMap(ctx context.Context) (map[int64]bool, diag.Diagnostics) {
	var d diag.Diagnostics
	ret, devIds := make(map[int64]bool), make([]int64, 0)
	d.Append(g.DeviceIds.ElementsAs(ctx, &devIds, false)...)
	for _, id := range devIds {
		ret[id] = true
	}
	return ret, d
}
