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

// VLanNetworks of OME
type VLanNetworks struct {
	ID                int64  `json:"Id"`
	Name              string `json:"Name"`
	Description       string `json:"Description"`
	VLANMaximum       int64  `json:"VlanMaximum"`
	VLANMinimum       int64  `json:"VlanMinimum"`
	Type              int64  `json:"Type"`
	InternalRefNWUUID string `json:"InternalRefNWUUId"`
}

// VLanNetworksTypeTfsdk is used to hold the config data
type VLanNetworksTypeTfsdk struct {
	ID           types.String        `tfsdk:"id"`
	VlanNetworks []VLanNetworksTfsdk `tfsdk:"vlan_networks"`
}

// VLanNetworksTfsdk is used to hold the vlan network data
type VLanNetworksTfsdk struct {
	VlanID            types.Int64  `tfsdk:"vlan_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	VLANMaximum       types.Int64  `tfsdk:"vlan_maximum"`
	VLANMinimum       types.Int64  `tfsdk:"vlan_minimum"`
	Type              types.Int64  `tfsdk:"type"`
	InternalRefNWUUID types.String `tfsdk:"internal_ref_nwuu_id"`
}
