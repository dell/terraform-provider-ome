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
