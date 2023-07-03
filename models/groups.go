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

// GroupDevicesData - schema for data source groupdevices
type GroupDevicesData struct {
	ID                types.String `tfsdk:"id"`
	DeviceIDs         types.List   `tfsdk:"device_ids"`
	DeviceServicetags types.List   `tfsdk:"device_servicetags"`
	DeviceGroupNames  types.Set    `tfsdk:"device_group_names"`
}

type ObjectId struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// type DeviceGroup struct {
// 	ID   int64  `json:"Id"`
// 	Name string `json:"Name"`
// 	Description string `json:"Description"`
// 	MembershipTypeId int `json:"MembershipTypeId"`
// 	ParentId int `json:"ParentId"`
// }

// GroupDevicesRes - schema for resource groupdevices
type GroupDeviceRes struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	// DeviceGroupNames  types.Set    `tfsdk:"device_group_names"`
	MembershipTypeId types.Int64 `tfsdk:"membership_type_id"`
	ParentId         types.Int64 `tfsdk:"parent_id"`
}

func NewGroupDeviceRes(g Group) (GroupDeviceRes, error) {
	return GroupDeviceRes{
		ID:               types.Int64Value(g.ID),
		Name:             types.StringValue(g.Name),
		Description:      types.StringValue(g.Description),
		MembershipTypeId: types.Int64Value(g.MembershipTypeId),
		ParentId:         types.Int64Value(g.ParentId),
	}, nil
}

func (g *GroupDeviceRes) GetPayload() Group {
	return Group{
		ID:               g.ID.ValueInt64(),
		Name:             g.Name.ValueString(),
		Description:      g.Description.ValueString(),
		MembershipTypeId: 12,
		ParentId:         g.ParentId.ValueInt64(),
	}
}

func (g *GroupDeviceRes) GetUpdatePayload() Group {
	return Group{
		ID:               g.ID.ValueInt64(),
		Name:             g.Name.ValueString(),
		Description:      g.Description.ValueString(),
		MembershipTypeId: 12,
		ParentId:         g.ParentId.ValueInt64(),
	}
}
