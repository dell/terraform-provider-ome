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

// OmeDeviceInventory - OmeDeviceInventory
type OmeDeviceInventory struct {
	ServerDeviceCards     []OmeServerDeviceCardInfo  `tfsdk:"server_device_cards"`
	CPUInfo               []OmeCPUInfo               `tfsdk:"cpus"`
	NICInfo               []OmeNICInfo               `tfsdk:"nics"`
	FCInfo                []OmeFCInfo                `tfsdk:"fcis"`
	OSInfo                []OmeOSInfo                `tfsdk:"os"`
	PowerSupplyInfo       []OmePowerSupplyInfo       `tfsdk:"power_supply"`
	DiskInfo              []OmeDiskInfo              `tfsdk:"disks"`
	RAIDControllerInfo    []OmeRAIDControllerInfo    `tfsdk:"raid_controllers"`
	MemoryInfo            []OmeMemoryInfo            `tfsdk:"memory"`
	StorageEnclosureInfo  []OmeStorageEnclosureInfo  `tfsdk:"storage_enclosures"`
	ServerPowerStates     []OmeServerPowerState      `tfsdk:"power_state"`
	DeviceLicenses        []OmeDeviceLicense         `tfsdk:"licenses"`
	DeviceCapabilities    []OmeDeviceCapability      `tfsdk:"capabilities"`
	DeviceFrus            []OmeDeviceFru             `tfsdk:"frus"`
	DeviceLocations       []OmeDeviceLocation        `tfsdk:"locations"`
	DeviceManagement      []OmeDeviceManagementInfo  `tfsdk:"management_info"`
	DeviceSoftwares       []OmeDeviceSoftware        `tfsdk:"softwares"`
	SubSystemRollupStatus []OmeSubSystemRollupStatus `tfsdk:"subsytem_rollup_status"`
}

// OmeSubSystemRollupStatus - OmeSubSystemRollupStatus
type OmeSubSystemRollupStatus struct {
	ID            types.Int64  `tfsdk:"id"`
	Status        types.Int64  `tfsdk:"status"`
	SubsystemName types.String `tfsdk:"subsystem_name"`
}

// OmeDeviceSoftware - OmeDeviceSoftware
type OmeDeviceSoftware struct {
	Version           types.String `tfsdk:"version"`
	InstallationDate  types.String `tfsdk:"installation_date"`
	Status            types.String `tfsdk:"status"`
	SoftwareType      types.String `tfsdk:"software_type"`
	VendorID          types.String `tfsdk:"vendor_id"`
	SubDeviceID       types.String `tfsdk:"sub_device_id"`
	SubVendorID       types.String `tfsdk:"sub_vendor_id"`
	ComponentID       types.String `tfsdk:"component_id"`
	PciDeviceID       types.String `tfsdk:"pci_device_id"`
	DeviceDescription types.String `tfsdk:"device_description"`
	InstanceID        types.String `tfsdk:"instance_id"`
}

// OmeDeviceManagementInfo - OmeDeviceManagementInfo
type OmeDeviceManagementInfo struct {
	ManagementID        types.Int64        `tfsdk:"management_id"`
	IPAddress           types.String       `tfsdk:"ip_address"`
	MACAddress          types.String       `tfsdk:"mac_address"`
	InstrumentationName types.String       `tfsdk:"instrumentation_name"`
	DNSName             types.String       `tfsdk:"dns_name"`
	ManagementType      OmeManagementType  `tfsdk:"management_type"`
	EndPointAgents      []OmeEndPointAgent `tfsdk:"end_point_agents"`
}

// OmeManagementType - OmeManagementType
type OmeManagementType struct {
	ManagementType types.Int64  `tfsdk:"management_type"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
}

// OmeEndPointAgent - OmeEndPointAgent
type OmeEndPointAgent struct {
	ManagementProfileID types.Int64  `tfsdk:"management_profile_id"`
	ProfileID           types.String `tfsdk:"profile_id"`
	AgentName           types.String `tfsdk:"agent_name"`
	Version             types.String `tfsdk:"version"`
	ManagementURL       types.String `tfsdk:"management_url"`
	HasCreds            types.Int64  `tfsdk:"has_creds"`
	Status              types.Int64  `tfsdk:"status"`
	StatusDateTime      types.String `tfsdk:"status_date_time"`
}

// OmeDeviceLocation - OmeDeviceLocation
type OmeDeviceLocation struct {
	ID                   types.Int64  `tfsdk:"id"`
	Room                 types.String `tfsdk:"room"`
	Rack                 types.String `tfsdk:"rack"`
	Aisle                types.String `tfsdk:"aisle"`
	Datacenter           types.String `tfsdk:"datacenter"`
	Rackslot             types.String `tfsdk:"rackslot"`
	ManagementSystemUnit types.Int64  `tfsdk:"management_system_unit"`
}

// OmeDeviceFru - OmeDeviceFru
type OmeDeviceFru struct {
	Revision     types.String `tfsdk:"revision"`
	ID           types.Int64  `tfsdk:"id"`
	Manufacturer types.String `tfsdk:"manufacturer"`
	Name         types.String `tfsdk:"name"`
	PartNumber   types.String `tfsdk:"part_number"`
	SerialNumber types.String `tfsdk:"serial_number"`
}

// OmeDeviceCapability - OmeDeviceCapability
type OmeDeviceCapability struct {
	ID             types.Int64             `tfsdk:"id"`
	CapabilityType OmeDeviceCapabilityType `tfsdk:"capability_type"`
}

// OmeDeviceCapabilityType - OmeDeviceCapabilityType
type OmeDeviceCapabilityType struct {
	CapabilityID types.Int64  `tfsdk:"capability_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	IDOwner      types.Int64  `tfsdk:"id_owner"`
}

// OmeServerDeviceCardInfo - OmeServerDeviceCardInfo
type OmeServerDeviceCardInfo struct {
	ID           types.Int64  `tfsdk:"id"`
	SlotNumber   types.String `tfsdk:"slot_number"`
	Manufacturer types.String `tfsdk:"manufacturer"`
	Description  types.String `tfsdk:"description"`
	DatabusWidth types.String `tfsdk:"databus_width"`
	SlotLength   types.String `tfsdk:"slot_length"`
	SlotType     types.String `tfsdk:"slot_type"`
}

// OmeCPUInfo - OmeCPUInfo
type OmeCPUInfo struct {
	ID                   types.Int64  `tfsdk:"id"`
	Family               types.String `tfsdk:"family"`
	MaxSpeed             types.Int64  `tfsdk:"max_speed"`
	CurrentSpeed         types.Int64  `tfsdk:"current_speed"`
	SlotNumber           types.String `tfsdk:"slot_number"`
	Status               types.Int64  `tfsdk:"status"`
	NumberOfCores        types.Int64  `tfsdk:"number_of_cores"`
	NumberOfEnabledCores types.Int64  `tfsdk:"number_of_enabled_cores"`
	BrandName            types.String `tfsdk:"brand_name"`
	ModelName            types.String `tfsdk:"model_name"`
	InstanceID           types.String `tfsdk:"instance_id"`
	Voltage              types.String `tfsdk:"voltage"`
}

// NewOmeCPUInfo - NewOmeCPUInfo
func NewOmeCPUInfo(input CPUInfo) OmeCPUInfo {
	return OmeCPUInfo{
		ID:                   types.Int64Value(input.ID),
		Family:               types.StringValue(input.Family),
		MaxSpeed:             types.Int64Value(input.MaxSpeed),
		CurrentSpeed:         types.Int64Value(input.CurrentSpeed),
		SlotNumber:           types.StringValue(input.SlotNumber),
		Status:               types.Int64Value(input.Status),
		NumberOfCores:        types.Int64Value(input.NumberOfCores),
		NumberOfEnabledCores: types.Int64Value(input.NumberOfEnabledCores),
		BrandName:            types.StringValue(input.BrandName),
		ModelName:            types.StringValue(input.ModelName),
		InstanceID:           types.StringValue(input.InstanceID),
		Voltage:              types.StringValue(input.Voltage),
	}
}

// OmePartition - OmePartition
type OmePartition struct {
	Fqdd                     types.String `tfsdk:"fqdd"`
	CurrentMacAddress        types.String `tfsdk:"current_mac_address"`
	PermanentMacAddress      types.String `tfsdk:"permanent_mac_address"`
	PermanentIscsiMacAddress types.String `tfsdk:"permanent_iscsi_mac_address"`
	PermanentFcoeMacAddress  types.String `tfsdk:"permanent_fcoe_mac_address"`
	Wwn                      types.String `tfsdk:"wwn"`
	Wwpn                     types.String `tfsdk:"wwpn"`
	VirtualWwn               types.String `tfsdk:"virtual_wwn"`
	VirtualWwpn              types.String `tfsdk:"virtual_wwpn"`
	VirtualMacAddress        types.String `tfsdk:"virtual_mac_address"`
	NicMode                  types.String `tfsdk:"nic_mode"`
	FcoeMode                 types.String `tfsdk:"fcoe_mode"`
	IscsiMode                types.String `tfsdk:"iscsi_mode"`
	MinBandwidth             types.Int64  `tfsdk:"min_bandwidth"`
	MaxBandwidth             types.Int64  `tfsdk:"max_bandwidth"`
}

// OmePort - OmePort
type OmePort struct {
	PortID      types.String   `tfsdk:"port_id"`
	ProductName types.String   `tfsdk:"product_name"`
	LinkStatus  types.String   `tfsdk:"link_status"`
	LinkSpeed   types.Int64    `tfsdk:"link_speed"`
	Partitions  []OmePartition `tfsdk:"partitions"`
}

// OmeNICInfo - OmeNICInfo
type OmeNICInfo struct {
	NicID      types.String `tfsdk:"nic_id"`
	VendorName types.String `tfsdk:"vendor_name"`
	Ports      []OmePort    `tfsdk:"ports"`
}

// OmeFCInfo - OmeFCInfo
type OmeFCInfo struct {
	ID                 types.Int64  `tfsdk:"id"`
	Fqdd               types.String `tfsdk:"fqdd"`
	DeviceDescription  types.String `tfsdk:"device_description"`
	DeviceName         types.String `tfsdk:"device_name"`
	FirstFctargetLun   types.String `tfsdk:"first_fctarget_lun"`
	FirstFctargetWwpn  types.String `tfsdk:"first_fctarget_wwpn"`
	PortNumber         types.Int64  `tfsdk:"port_number"`
	PortSpeed          types.String `tfsdk:"port_speed"`
	SecondFctargetLun  types.String `tfsdk:"second_fctarget_lun"`
	SecondFctargetWwpn types.String `tfsdk:"second_fctarget_wwpn"`
	VendorName         types.String `tfsdk:"vendor_name"`
	Wwn                types.String `tfsdk:"wwn"`
	Wwpn               types.String `tfsdk:"wwpn"`
	LinkStatus         types.String `tfsdk:"link_status"`
	VirtualWwn         types.String `tfsdk:"virtual_wwn"`
	VirtualWwpn        types.String `tfsdk:"virtual_wwpn"`
}

// OmeOSInfo - OmeOSInfo
type OmeOSInfo struct {
	ID        types.Int64  `tfsdk:"id"`
	OsName    types.String `tfsdk:"os_name"`
	OsVersion types.String `tfsdk:"os_version"`
	Hostname  types.String `tfsdk:"hostname"`
}

// OmePowerSupplyInfo - OmePowerSupplyInfo
type OmePowerSupplyInfo struct {
	ID                                  types.Int64  `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	PowerSupplyType                     types.Int64  `tfsdk:"power_supply_type"`
	OutputWatts                         types.Int64  `tfsdk:"output_watts"`
	Location                            types.String `tfsdk:"location"`
	RedundancyState                     types.String `tfsdk:"redundancy_state"`
	Status                              types.Int64  `tfsdk:"status"`
	State                               types.String `tfsdk:"state"`
	FirmwareVersion                     types.String `tfsdk:"firmware_version"`
	InputVoltage                        types.Int64  `tfsdk:"input_voltage"`
	Model                               types.String `tfsdk:"model"`
	Manufacturer                        types.String `tfsdk:"manufacturer"`
	Range1MaxInputPowerWatts            types.Int64  `tfsdk:"range1_max_input_power_watts"`
	SerialNumber                        types.String `tfsdk:"serial_number"`
	ActiveInputVoltage                  types.String `tfsdk:"active_input_voltage"`
	InputPowerUnits                     types.String `tfsdk:"input_power_units"`
	OperationalStatus                   types.String `tfsdk:"operational_status"`
	Range1MaxInputVoltageHighMilliVolts types.Int64  `tfsdk:"range1_max_input_voltage_high_milli_volts"`
	RatedMaxOutputPower                 types.Int64  `tfsdk:"rated_max_output_power"`
	RequestedState                      types.Int64  `tfsdk:"requested_state"`
	AcInput                             types.Bool   `tfsdk:"ac_input"`
	AcOutput                            types.Bool   `tfsdk:"ac_output"`
	SwitchingSupply                     types.Bool   `tfsdk:"switching_supply"`
}

// OmeDiskInfo - OmeDiskInfo
type OmeDiskInfo struct {
	ID                          types.Int64  `tfsdk:"id"`
	DiskNumber                  types.String `tfsdk:"disk_number"`
	VendorName                  types.String `tfsdk:"vendor_name"`
	Status                      types.Int64  `tfsdk:"status"`
	StatusString                types.String `tfsdk:"status_string"`
	ModelNumber                 types.String `tfsdk:"model_number"`
	SerialNumber                types.String `tfsdk:"serial_number"`
	SasAddress                  types.String `tfsdk:"sas_address"`
	Revision                    types.String `tfsdk:"revision"`
	ManufacturedDay             types.Int64  `tfsdk:"manufactured_day"`
	ManufacturedWeek            types.Int64  `tfsdk:"manufactured_week"`
	ManufacturedYear            types.Int64  `tfsdk:"manufactured_year"`
	EncryptionAbility           types.Bool   `tfsdk:"encryption_ability"`
	FormFactor                  types.String `tfsdk:"form_factor"`
	PartNumber                  types.String `tfsdk:"part_number"`
	PredictiveFailureState      types.String `tfsdk:"predictive_failure_state"`
	EnclosureID                 types.String `tfsdk:"enclosure_id"`
	Channel                     types.Int64  `tfsdk:"channel"`
	Size                        types.String `tfsdk:"size"`
	FreeSpace                   types.String `tfsdk:"free_space"`
	UsedSpace                   types.String `tfsdk:"used_space"`
	BusType                     types.String `tfsdk:"bus_type"`
	SlotNumber                  types.Int64  `tfsdk:"slot_number"`
	MediaType                   types.String `tfsdk:"media_type"`
	RemainingReadWriteEndurance types.String `tfsdk:"remaining_read_write_endurance"`
	SecurityState               types.String `tfsdk:"security_state"`
	RaidStatus                  types.String `tfsdk:"raid_status"`
}

// OmeServerVirtualDisk - OmeServerVirtualDisk
type OmeServerVirtualDisk struct {
	ID               types.Int64  `tfsdk:"id"`
	RaidControllerID types.Int64  `tfsdk:"raid_controller_id"`
	DeviceID         types.Int64  `tfsdk:"device_id"`
	Fqdd             types.String `tfsdk:"fqdd"`
	State            types.String `tfsdk:"state"`
	RollupStatus     types.Int64  `tfsdk:"rollup_status"`
	Status           types.Int64  `tfsdk:"status"`
	Layout           types.String `tfsdk:"layout"`
	MediaType        types.String `tfsdk:"media_type"`
	Name             types.String `tfsdk:"name"`
	ReadPolicy       types.String `tfsdk:"read_policy"`
	WritePolicy      types.String `tfsdk:"write_policy"`
	CachePolicy      types.String `tfsdk:"cache_policy"`
	StripeSize       types.String `tfsdk:"stripe_size"`
	Size             types.String `tfsdk:"size"`
	TargetID         types.Int64  `tfsdk:"target_id"`
	LockStatus       types.String `tfsdk:"lock_status"`
}

// OmeRAIDControllerInfo - OmeRAIDControllerInfo
type OmeRAIDControllerInfo struct {
	ID                       types.Int64            `tfsdk:"id"`
	Name                     types.String           `tfsdk:"name"`
	Fqdd                     types.String           `tfsdk:"fqdd"`
	DeviceDescription        types.String           `tfsdk:"device_description"`
	Status                   types.Int64            `tfsdk:"status"`
	StatusType               types.String           `tfsdk:"status_type"`
	RollupStatus             types.Int64            `tfsdk:"rollup_status"`
	RollupStatusString       types.String           `tfsdk:"rollup_status_string"`
	FirmwareVersion          types.String           `tfsdk:"firmware_version"`
	CacheSizeInMb            types.Int64            `tfsdk:"cache_size_in_mb"`
	PciSlot                  types.String           `tfsdk:"pci_slot"`
	DriverVersion            types.String           `tfsdk:"driver_version"`
	StorageAssignmentAllowed types.String           `tfsdk:"storage_assignment_allowed"`
	ServerVirtualDisks       []OmeServerVirtualDisk `tfsdk:"server_virtual_disks"`
}

// OmeMemoryInfo - OmeMemoryInfo
type OmeMemoryInfo struct {
	ID                    types.Int64  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	BankName              types.String `tfsdk:"bank_name"`
	Size                  types.Int64  `tfsdk:"size"`
	Status                types.Int64  `tfsdk:"status"`
	Manufacturer          types.String `tfsdk:"manufacturer"`
	PartNumber            types.String `tfsdk:"part_number"`
	SerialNumber          types.String `tfsdk:"serial_number"`
	TypeDetails           types.String `tfsdk:"type_details"`
	ManufacturerDate      types.String `tfsdk:"manufacturer_date"`
	Speed                 types.Int64  `tfsdk:"speed"`
	CurrentOperatingSpeed types.Int64  `tfsdk:"current_operating_speed"`
	Rank                  types.String `tfsdk:"rank"`
	InstanceID            types.String `tfsdk:"instance_id"`
	DeviceDescription     types.String `tfsdk:"device_description"`
}

// OmeStorageEnclosureInfo - OmeStorageEnclosureInfo
type OmeStorageEnclosureInfo struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Status           types.Int64  `tfsdk:"status"`
	StatusTypeString types.String `tfsdk:"status_type"`
	ChannelNumber    types.String `tfsdk:"channel_number"`
	BackplanePartNum types.String `tfsdk:"backplane_part_num"`
	NumberOfFanPacks types.Int64  `tfsdk:"number_of_fan_packs"`
	Version          types.String `tfsdk:"version"`
	RollupStatus     types.Int64  `tfsdk:"rollup_status"`
	SlotCount        types.Int64  `tfsdk:"slot_count"`
}

// OmeServerPowerState - OmeServerPowerState
type OmeServerPowerState struct {
	ID         types.Int64 `tfsdk:"id"`
	PowerState types.Int64 `tfsdk:"power_state"`
}

// OmeDeviceLicense - OmeDeviceLicense
type OmeDeviceLicense struct {
	SoldDate           types.String   `tfsdk:"sold_date"`
	LicenseBound       types.Int64    `tfsdk:"license_bound"`
	EvalTimeRemaining  types.Int64    `tfsdk:"eval_time_remaining"`
	AssignedDevices    types.String   `tfsdk:"assigned_devices"`
	LicenseStatus      types.Int64    `tfsdk:"license_status"`
	EntitlementId      types.String   `tfsdk:"entitlement_id"`
	LicenseDescription types.String   `tfsdk:"license_description"`
	LicenseType        OmeLicenseType `tfsdk:"license_type"`
}

// OmeLicenseType - OmeLicenseType
type OmeLicenseType struct {
	Name      types.String `tfsdk:"name"`
	LicenseID types.Int64  `tfsdk:"license_id"`
}
