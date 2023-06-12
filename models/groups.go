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

import "github.com/hashicorp/terraform-plugin-framework/types"

// Groups - list of group response from OME
type Groups struct {
	Value []Group `json:"value"`
}

// Group - embedded group response from the groups
type Group struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name"`
}

// GroupDevicesData - schema for data source groupdevices
type GroupDevicesData struct {
	ID                types.String `tfsdk:"id"`
	DeviceIDs         types.List   `tfsdk:"device_ids"`
	DeviceServicetags types.List   `tfsdk:"device_servicetags"`
	DeviceGroupNames  types.Set    `tfsdk:"device_group_names"`
}
