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
