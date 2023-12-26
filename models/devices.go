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
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ome/utils"
)

// Devices - list of device response from on OME
type Devices struct {
	Value    []Device `json:"value"`
	NextLink string   `json:"@odata.nextLink"`
}

// Device - embedded device response from the Devices
type Device struct {
	ID                            int64                       `json:"Id"`
	Type                          int64                       `json:"Type"`
	Identifier                    string                      `json:"Identifier"`
	DeviceServiceTag              string                      `json:"DeviceServiceTag"`
	ChassisServiceTag             string                      `json:"ChassisServiceTag"`
	Model                         string                      `json:"Model"`
	PowerState                    int64                       `json:"PowerState"`
	ManagedState                  int64                       `json:"ManagedState"`
	Status                        int64                       `json:"Status"`
	ConnectionState               bool                        `json:"ConnectionState"`
	AssetTag                      *string                     `json:"AssetTag"`
	SystemID                      int64                       `json:"SystemId"`
	DeviceName                    string                      `json:"DeviceName"`
	LastInventoryTime             string                      `json:"LastInventoryTime"`
	LastStatusTime                string                      `json:"LastStatusTime"`
	DeviceSubscription            *string                     `json:"DeviceSubscription"`
	DeviceCapabilities            []int64                     `json:"DeviceCapabilities"`
	SlotConfiguration             SlotConfiguration           `json:"SlotConfiguration"`
	DeviceManagement              []DeviceManagement          `json:"DeviceManagement"`
	Enabled                       bool                        `json:"Enabled"`
	ConnectionStateReason         int64                       `json:"ConnectionStateReason"`
	ChassisIP                     string                      `json:"ChassisIp"`
	DiscoveryConfigurationJobInfo []DiscoveryConfigurationJob `json:"DiscoveryConfigurationJobInformation"`
}

// BelongsToPool - method to check if a device belongs to that ip pool
func (d *Device) BelongsToPool(pool utils.IPSet) bool {
	for _, devM := range d.DeviceManagement {
		if pool.Contains(devM.NetworkAddress) {
			return true
		}
	}
	return false
}

// DiscoveryConfigurationJob - DiscoveryConfigurationJob
type DiscoveryConfigurationJob struct {
	GroupID          string `json:"GroupId"`
	CreatedBy        string `json:"CreatedBy"`
	DiscoveryJobName string `json:"DiscoveryJobName"`
}

// SlotConfiguration - SlotConfiguration
type SlotConfiguration struct {
	ChassisName *string `json:"ChassisName"`
}

// DeviceManagement - embedded device management response from the Devices
type DeviceManagement struct {
	ManagementID        int64               `json:"ManagementId"`
	NetworkAddress      net.IP              `json:"NetworkAddress"`
	MacAddress          string              `json:"MacAddress"`
	ManagementType      int64               `json:"ManagementType"`
	InstrumentationName string              `json:"InstrumentationName"`
	DNSName             string              `json:"DnsName"`
	ManagementProfile   []ManagementProfile `json:"ManagementProfile"`
}

// ManagementProfile - ManagementProfile
type ManagementProfile struct {
	ManagementProfileID int64  `json:"ManagementProfileId"`
	ProfileID           string `json:"ProfileId"`
	ManagementID        int64  `json:"ManagementId"`
	AgentName           string `json:"AgentName"`
	Version             string `json:"Version"`
	ManagementURL       string `json:"ManagementURL"`
	HasCreds            int64  `json:"HasCreds"`
	Status              int64  `json:"Status"`
	StatusDateTime      string `json:"StatusDateTime"`
}

// ######## tfsdk models

// DevicesResModel - Tfsdk model for devices resource
type DevicesResModel struct {
	ID      types.String `tfsdk:"id"`
	Devices types.List   `tfsdk:"devices"` // []DeviceItemModel
}

// DeviceItemModel - Tfsdk model for each device in devices resource
type DeviceItemModel struct {
	ID         types.Int64  `tfsdk:"id"`
	ServiceTag types.String `tfsdk:"service_tag"`
	IPs        types.List   `tfsdk:"management_ips"` // []string
}

func (DeviceItemModel) getType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.Int64Type,
		"service_tag": types.StringType,
		"management_ips": types.ListType{
			ElemType: types.StringType,
		},
	}
}

// OmeDeviceData - schema for device data source
type OmeDeviceData struct {
	ID             types.Int64           `tfsdk:"id"`
	Filters        types.Object          `tfsdk:"filters"`
	Devices        []OmeSingleDeviceData `tfsdk:"devices"`
	InventoryTypes []string              `tfsdk:"inventory_types"`
}

// OmeDeviceDataFilters - schema for device data source filters
type OmeDeviceDataFilters struct {
	IDs        types.List   `tfsdk:"ids"`
	SvcTags    types.List   `tfsdk:"device_service_tags"`
	IPExprs    types.List   `tfsdk:"ip_expressions"`
	FilterExpr types.String `tfsdk:"filter_expression"`
}

// OmeSingleDeviceData is the tfsdk version of Device
type OmeSingleDeviceData struct {
	ID                            types.Int64                        `tfsdk:"id"`
	Type                          types.Int64                        `tfsdk:"type"`
	Identifier                    types.String                       `tfsdk:"identifier"`
	DeviceServiceTag              types.String                       `tfsdk:"device_service_tag"`
	ChassisServiceTag             types.String                       `tfsdk:"chassis_service_tag"`
	Model                         types.String                       `tfsdk:"model"`
	PowerState                    types.Int64                        `tfsdk:"power_state"`
	ManagedState                  types.Int64                        `tfsdk:"managed_state"`
	Status                        types.Int64                        `tfsdk:"status"`
	ConnectionState               types.Bool                         `tfsdk:"connection_state"`
	AssetTag                      types.String                       `tfsdk:"asset_tag"`
	SystemID                      types.Int64                        `tfsdk:"system_id"`
	DeviceName                    types.String                       `tfsdk:"device_name"`
	LastInventoryTime             types.String                       `tfsdk:"last_inventory_time"`
	LastStatusTime                types.String                       `tfsdk:"last_status_time"`
	DeviceSubscription            types.String                       `tfsdk:"device_subscription"`
	DeviceCapabilities            types.List                         `tfsdk:"device_capabilities"`
	SlotConfiguration             OmeSlotConfigurationData           `tfsdk:"slot_configuration"`
	DeviceManagement              []OmeDeviceManagementData          `tfsdk:"device_management"`
	Enabled                       types.Bool                         `tfsdk:"enabled"`
	ConnectionStateReason         types.Int64                        `tfsdk:"connection_state_reason"`
	ChassisIP                     types.String                       `tfsdk:"chassis_ip"`
	DiscoveryConfigurationJobInfo []OmeDiscoveryConfigurationJobData `tfsdk:"discovery_configuration_job_information"`
	Inventory                     OmeDeviceInventory                 `tfsdk:"detailed_inventory"`
}

// OmeDiscoveryConfigurationJobData is the tfsdk version of DiscoveryConfigurationJob
type OmeDiscoveryConfigurationJobData struct {
	GroupID          types.String `tfsdk:"group_id"`
	CreatedBy        types.String `tfsdk:"created_by"`
	DiscoveryJobName types.String `tfsdk:"discovery_job_name"`
}

// OmeSlotConfigurationData is the tfsdk version of SlotConfiguration
type OmeSlotConfigurationData struct {
	ChassisName types.String `tfsdk:"chassis_name"`
}

// OmeDeviceManagementData is the tfsdk version of DeviceManagement
type OmeDeviceManagementData struct {
	ManagementID        types.Int64                `tfsdk:"management_id"`
	NetworkAddress      types.String               `tfsdk:"network_address"`
	MacAddress          types.String               `tfsdk:"mac_address"`
	ManagementType      types.Int64                `tfsdk:"management_type"`
	InstrumentationName types.String               `tfsdk:"instrumentation_name"`
	DNSName             types.String               `tfsdk:"dns_name"`
	ManagementProfile   []OmeManagementProfileData `tfsdk:"management_profile"`
}

// OmeManagementProfileData is the tfsdk version of ManagementProfile
type OmeManagementProfileData struct {
	ManagementProfileID types.Int64  `tfsdk:"management_profile_id"`
	ProfileID           types.String `tfsdk:"profile_id"`
	ManagementID        types.Int64  `tfsdk:"management_id"`
	AgentName           types.String `tfsdk:"agent_name"`
	Version             types.String `tfsdk:"version"`
	ManagementURL       types.String `tfsdk:"management_url"`
	HasCreds            types.Int64  `tfsdk:"has_creds"`
	Status              types.Int64  `tfsdk:"status"`
	StatusDateTime      types.String `tfsdk:"status_date_time"`
}

// ################### tfsdk converters

// NewDeviceItemModel - converts Device to DeviceItemModel
func NewDeviceItemModel(input Device) DeviceItemModel {
	ips := make([]string, 0)
	for _, m := range input.DeviceManagement {
		ips = append(ips, m.NetworkAddress.String())
	}
	return DeviceItemModel{
		ID:         types.Int64Value(input.ID),
		ServiceTag: types.StringValue(input.DeviceServiceTag),
		IPs:        stringListValue(ips),
	}
}

// NewDevicesResModel - converts list of json devices to DevicesResModel
func NewDevicesResModel(input []Device) (DevicesResModel, diag.Diagnostics) {
	devs := make([]DeviceItemModel, 0)
	for _, v := range input {
		devs = append(devs, NewDeviceItemModel(v))
	}
	tfdevs, dgs := objListValue(DeviceItemModel{}.getType(), devs)
	return DevicesResModel{
		ID:      types.StringValue("dummy"),
		Devices: tfdevs,
	}, dgs
}

// NewDevicesResModelFromID - converts a string id into a terraform plan
func NewDevicesResModelFromID(id string) (types.List, diag.Diagnostics) {
	var dgs diag.Diagnostics
	if id == "" {
		return types.ListUnknown(types.ObjectType{
			AttrTypes: DeviceItemModel{}.getType(),
		}), dgs
	}
	var plan types.List
	vals := strings.Split(id, ":")
	var did, didst string
	if l := len(vals); l == 1 {
		did = "id"
		didst = id
	} else if l > 2 {
		dgs.AddError(
			"Wrong \"id\" format.",
			fmt.Sprintf("Expecting a string having maximum 1 semicolon, found %d", l-1),
		)
		return plan, dgs
	} else {
		did = vals[0]
		didst = vals[1]
	}
	dids := strings.Split(didst, ",")
	devs := make([]DeviceItemModel, 0)
	ips := make([]string, 0)
	if did == "id" {
		for _, did := range dids {
			didInt, err := strconv.ParseInt(did, 10, 64)
			if err != nil {
				dgs.AddError(
					"Wrong \"id\" format. ID could not be converted to an Int64.",
					err.Error(),
				)
			}
			devs = append(devs, DeviceItemModel{
				ServiceTag: types.StringUnknown(),
				ID:         types.Int64Value(didInt),
				IPs:        stringListValue(ips),
			})
		}
	} else if did == "svc_tag" {
		for _, did := range dids {
			devs = append(devs, DeviceItemModel{
				ServiceTag: types.StringValue(did),
				ID:         types.Int64Unknown(),
				IPs:        stringListValue(ips),
			})
		}
	} else {
		dgs.AddError(
			"Wrong \"id\" format.",
			fmt.Sprintf("Identifier of type %s is not recognised, valid values are [id, svc_tag].", did),
		)
	}
	if dgs.HasError() {
		return plan, dgs
	}
	return types.ListValueFrom(context.TODO(), types.ObjectType{
		AttrTypes: DeviceItemModel{}.getType(),
	}, devs)
}

// NewSingleOmeDeviceData converts DeviceData to OmeDeviceData
func NewSingleOmeDeviceData(input Device) OmeSingleDeviceData {
	return OmeSingleDeviceData{
		ID:                            types.Int64Value(input.ID),
		Type:                          types.Int64Value(input.Type),
		Identifier:                    types.StringValue(input.Identifier),
		DeviceServiceTag:              types.StringValue(input.DeviceServiceTag),
		ChassisServiceTag:             types.StringValue(input.ChassisServiceTag),
		Model:                         types.StringValue(input.Model),
		PowerState:                    types.Int64Value(input.PowerState),
		ManagedState:                  types.Int64Value(input.ManagedState),
		Status:                        types.Int64Value(input.Status),
		ConnectionState:               types.BoolValue(input.ConnectionState),
		AssetTag:                      stringPointerValue(input.AssetTag),
		SystemID:                      types.Int64Value(input.SystemID),
		DeviceName:                    types.StringValue(input.DeviceName),
		LastInventoryTime:             types.StringValue(input.LastInventoryTime),
		LastStatusTime:                types.StringValue(input.LastStatusTime),
		DeviceSubscription:            stringPointerValue(input.DeviceSubscription),
		DeviceCapabilities:            int64ListValue(input.DeviceCapabilities),
		SlotConfiguration:             newOmeSlotConfigurationData(input.SlotConfiguration),
		DeviceManagement:              newDeviceManagementList(input.DeviceManagement),
		Enabled:                       types.BoolValue(input.Enabled),
		ConnectionStateReason:         types.Int64Value(input.ConnectionStateReason),
		ChassisIP:                     types.StringValue(input.ChassisIP),
		DiscoveryConfigurationJobInfo: newDiscoveryConfigurationJobList(input.DiscoveryConfigurationJobInfo),
		// Inventory: NewOmeDeviceInventory(input.),
	}
}

func newDeviceManagementList(inputs []DeviceManagement) []OmeDeviceManagementData {
	ret := make([]OmeDeviceManagementData, 0)
	for _, input := range inputs {
		ret = append(ret, newOmeDeviceManagementData(input))
	}
	return ret
}

func newDiscoveryConfigurationJobList(inputs []DiscoveryConfigurationJob) []OmeDiscoveryConfigurationJobData {
	ret := make([]OmeDiscoveryConfigurationJobData, 0)
	for _, input := range inputs {
		ret = append(ret, newOmeDiscoveryConfigurationJobData(input))
	}
	return ret
}

// newOmeDiscoveryConfigurationJobData converts DiscoveryConfigurationJobData to OmeDiscoveryConfigurationJobData
func newOmeDiscoveryConfigurationJobData(input DiscoveryConfigurationJob) OmeDiscoveryConfigurationJobData {
	return OmeDiscoveryConfigurationJobData{
		GroupID:          types.StringValue(input.GroupID),
		CreatedBy:        types.StringValue(input.CreatedBy),
		DiscoveryJobName: types.StringValue(input.DiscoveryJobName),
	}
}

// newOmeSlotConfigurationData converts SlotConfigurationData to OmeSlotConfigurationData
func newOmeSlotConfigurationData(input SlotConfiguration) OmeSlotConfigurationData {
	return OmeSlotConfigurationData{
		ChassisName: stringPointerValue(input.ChassisName),
	}
}

// newOmeDeviceManagementData converts DeviceManagementData to OmeDeviceManagementData
func newOmeDeviceManagementData(input DeviceManagement) OmeDeviceManagementData {
	a := OmeDeviceManagementData{
		ManagementID:        types.Int64Value(input.ManagementID),
		NetworkAddress:      types.StringValue(input.NetworkAddress.String()),
		MacAddress:          types.StringValue(input.MacAddress),
		ManagementType:      types.Int64Value(input.ManagementType),
		InstrumentationName: types.StringValue(input.InstrumentationName),
		DNSName:             types.StringValue(input.DNSName),
	}
	mProfs := make([]OmeManagementProfileData, 0)
	for _, mProf := range input.ManagementProfile {
		mProfs = append(mProfs, newOmeManagementProfileData(mProf))
	}
	a.ManagementProfile = mProfs
	return a
}

// newOmeManagementProfileData converts ManagementProfileData to OmeManagementProfileData
func newOmeManagementProfileData(input ManagementProfile) OmeManagementProfileData {
	return OmeManagementProfileData{
		ManagementProfileID: types.Int64Value(input.ManagementProfileID),
		ProfileID:           types.StringValue(input.ProfileID),
		ManagementID:        types.Int64Value(input.ManagementID),
		AgentName:           types.StringValue(input.AgentName),
		Version:             types.StringValue(input.Version),
		ManagementURL:       types.StringValue(input.ManagementURL),
		HasCreds:            types.Int64Value(input.HasCreds),
		Status:              types.Int64Value(input.Status),
		StatusDateTime:      types.StringValue(input.StatusDateTime),
	}
}
