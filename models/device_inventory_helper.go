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

// NewOmeDeviceInventory converts DeviceInventory to OmeDeviceInventory
func NewOmeDeviceInventory(input DeviceInventory) OmeDeviceInventory {
	return OmeDeviceInventory{
		ServerDeviceCards:     newOmeServerDeviceCardInfoList(input.ServerDeviceCards),
		CPUInfo:               newOmeCPUInfoList(input.CPUInfo),
		NICInfo:               newOmeNICInfoList(input.NICInfo),
		FCInfo:                newOmeFCInfoList(input.FCInfo),
		OSInfo:                newOmeOSInfoList(input.OSInfo),
		PowerSupplyInfo:       newOmePowerSupplyInfoList(input.PowerSupplyInfo),
		DiskInfo:              newOmeDiskInfoList(input.DiskInfo),
		RAIDControllerInfo:    newOmeRAIDControllerInfoList(input.RAIDControllerInfo),
		MemoryInfo:            newOmeMemoryInfoList(input.MemoryInfo),
		StorageEnclosureInfo:  newOmeStorageEnclosureInfoList(input.StorageEnclosureInfo),
		ServerPowerStates:     newOmeServerPowerStateList(input.ServerPowerStates),
		DeviceLicenses:        newOmeDeviceLicenseList(input.DeviceLicenses),
		DeviceCapabilities:    newOmeDeviceCapabilityList(input.DeviceCapabilities),
		DeviceFrus:            newOmeDeviceFruList(input.DeviceFrus),
		DeviceLocations:       newOmeDeviceLocationList(input.DeviceLocations),
		DeviceManagement:      newOmeDeviceManagementInfoList(input.DeviceManagement),
		DeviceSoftwares:       newOmeDeviceSoftwareList(input.DeviceSoftwares),
		SubSystemRollupStatus: newOmeSubSystemRollupStatusList(input.SubSystemRollupStatus),
	}
}

// newOmeServerDeviceCardInfoList converts list of ServerDeviceCardInfo to list of OmeServerDeviceCardInfo
func newOmeServerDeviceCardInfoList(inputs []ServerDeviceCardInfo) []OmeServerDeviceCardInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeServerDeviceCardInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeServerDeviceCardInfo(input))
	}
	return out
}

// newOmeCPUInfoList converts list of CPUInfo to list of OmeCPUInfo
func newOmeCPUInfoList(inputs []CPUInfo) []OmeCPUInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeCPUInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeCPUInfo(input))
	}
	return out
}

// newOmeNICInfoList converts list of NICInfo to list of OmeNICInfo
func newOmeNICInfoList(inputs []NICInfo) []OmeNICInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeNICInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeNICInfo(input))
	}
	return out
}

// newOmeFCInfoList converts list of FCInfo to list of OmeFCInfo
func newOmeFCInfoList(inputs []FCInfo) []OmeFCInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeFCInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeFCInfo(input))
	}
	return out
}

// newOmeOSInfoList converts list of OSInfo to list of OmeOSInfo
func newOmeOSInfoList(inputs []OSInfo) []OmeOSInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeOSInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeOSInfo(input))
	}
	return out
}

// newOmePowerSupplyInfoList converts list of PowerSupplyInfo to list of OmePowerSupplyInfo
func newOmePowerSupplyInfoList(inputs []PowerSupplyInfo) []OmePowerSupplyInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmePowerSupplyInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmePowerSupplyInfo(input))
	}
	return out
}

// newOmeDiskInfoList converts list of DiskInfo to list of OmeDiskInfo
func newOmeDiskInfoList(inputs []DiskInfo) []OmeDiskInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeDiskInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeDiskInfo(input))
	}
	return out
}

// newOmeRAIDControllerInfoList converts list of RAIDControllerInfo to list of OmeRAIDControllerInfo
func newOmeRAIDControllerInfoList(inputs []RAIDControllerInfo) []OmeRAIDControllerInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeRAIDControllerInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeRAIDControllerInfo(input))
	}
	return out
}

// newOmeMemoryInfoList converts list of MemoryInfo to list of OmeMemoryInfo
func newOmeMemoryInfoList(inputs []MemoryInfo) []OmeMemoryInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeMemoryInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeMemoryInfo(input))
	}
	return out
}

// newOmeStorageEnclosureInfoList converts list of StorageEnclosureInfo to list of OmeStorageEnclosureInfo
func newOmeStorageEnclosureInfoList(inputs []StorageEnclosureInfo) []OmeStorageEnclosureInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeStorageEnclosureInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeStorageEnclosureInfo(input))
	}
	return out
}

// newOmeServerPowerStateList converts list of ServerPowerState to list of OmeServerPowerState
func newOmeServerPowerStateList(inputs []ServerPowerState) []OmeServerPowerState {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeServerPowerState, 0)
	for _, input := range inputs {
		out = append(out, newOmeServerPowerState(input))
	}
	return out
}

// newOmeDeviceLicenseList converts list of DeviceLicense to list of OmeDeviceLicense
func newOmeDeviceLicenseList(inputs []DeviceLicense) []OmeDeviceLicense {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeDeviceLicense, 0)
	for _, input := range inputs {
		out = append(out, newOmeDeviceLicense(input))
	}
	return out
}

// newOmeDeviceCapabilityList converts list of DeviceCapability to list of OmeDeviceCapability
func newOmeDeviceCapabilityList(inputs []DeviceCapability) []OmeDeviceCapability {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeDeviceCapability, 0)
	for _, input := range inputs {
		out = append(out, newOmeDeviceCapability(input))
	}
	return out
}

// newOmeDeviceFruList converts list of DeviceFru to list of OmeDeviceFru
func newOmeDeviceFruList(inputs []DeviceFru) []OmeDeviceFru {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeDeviceFru, 0)
	for _, input := range inputs {
		out = append(out, newOmeDeviceFru(input))
	}
	return out
}

// newOmeDeviceLocationList converts list of DeviceLocation to list of OmeDeviceLocation
func newOmeDeviceLocationList(inputs []DeviceLocation) []OmeDeviceLocation {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeDeviceLocation, 0)
	for _, input := range inputs {
		out = append(out, newOmeDeviceLocation(input))
	}
	return out
}

// newOmeDeviceManagementInfoList converts list of DeviceManagementInfo to list of OmeDeviceManagementInfo
func newOmeDeviceManagementInfoList(inputs []DeviceManagementInfo) []OmeDeviceManagementInfo {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeDeviceManagementInfo, 0)
	for _, input := range inputs {
		out = append(out, newOmeDeviceManagementInfo(input))
	}
	return out
}

// newOmeDeviceSoftwareList converts list of DeviceSoftware to list of OmeDeviceSoftware
func newOmeDeviceSoftwareList(inputs []DeviceSoftware) []OmeDeviceSoftware {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeDeviceSoftware, 0)
	for _, input := range inputs {
		out = append(out, newOmeDeviceSoftware(input))
	}
	return out
}

// newOmeSubSystemRollupStatusList converts list of SubSystemRollupStatus to list of OmeSubSystemRollupStatus
func newOmeSubSystemRollupStatusList(inputs []SubSystemRollupStatus) []OmeSubSystemRollupStatus {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeSubSystemRollupStatus, 0)
	for _, input := range inputs {
		out = append(out, newOmeSubSystemRollupStatus(input))
	}
	return out
}

// newOmeSubSystemRollupStatus converts SubSystemRollupStatus to OmeSubSystemRollupStatus
func newOmeSubSystemRollupStatus(input SubSystemRollupStatus) OmeSubSystemRollupStatus {
	return OmeSubSystemRollupStatus{
		ID:            types.Int64Value(input.ID),
		Status:        types.Int64Value(input.Status),
		SubsystemName: types.StringValue(input.SubsystemName),
	}
}

// newOmeDeviceSoftware converts DeviceSoftware to OmeDeviceSoftware
func newOmeDeviceSoftware(input DeviceSoftware) OmeDeviceSoftware {
	return OmeDeviceSoftware{
		Version:           types.StringValue(input.Version),
		InstallationDate:  types.StringValue(input.InstallationDate),
		Status:            types.StringValue(input.Status),
		SoftwareType:      types.StringValue(input.SoftwareType),
		VendorID:          types.StringValue(input.VendorID),
		SubDeviceID:       types.StringValue(input.SubDeviceID),
		SubVendorID:       types.StringValue(input.SubVendorID),
		ComponentID:       types.StringValue(input.ComponentID),
		PciDeviceID:       types.StringValue(input.PciDeviceID),
		DeviceDescription: types.StringValue(input.DeviceDescription),
		InstanceID:        types.StringValue(input.InstanceID),
	}
}

// newOmeDeviceManagementInfo converts DeviceManagementInfo to OmeDeviceManagementInfo
func newOmeDeviceManagementInfo(input DeviceManagementInfo) OmeDeviceManagementInfo {
	return OmeDeviceManagementInfo{
		ManagementID:        types.Int64Value(input.ManagementID),
		IPAddress:           types.StringValue(input.IPAddress),
		MACAddress:          types.StringValue(input.MACAddress),
		InstrumentationName: types.StringValue(input.InstrumentationName),
		DNSName:             types.StringValue(input.DNSName),
		ManagementType:      newOmeManagementType(input.ManagementType),
		EndPointAgents:      newOmeEndPointAgentList(input.EndPointAgents),
	}
}

// newOmeEndPointAgentList converts list of EndPointAgent to list of OmeEndPointAgent
func newOmeEndPointAgentList(inputs []EndPointAgent) []OmeEndPointAgent {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeEndPointAgent, 0)
	for _, input := range inputs {
		out = append(out, newOmeEndPointAgent(input))
	}
	return out
}

// newOmeManagementType converts ManagementType to OmeManagementType
func newOmeManagementType(input ManagementType) OmeManagementType {
	return OmeManagementType{
		ManagementType: types.Int64Value(input.ManagementType),
		Name:           types.StringValue(input.Name),
		Description:    types.StringValue(input.Description),
	}
}

// newOmeEndPointAgent converts EndPointAgent to OmeEndPointAgent
func newOmeEndPointAgent(input EndPointAgent) OmeEndPointAgent {
	return OmeEndPointAgent{
		ManagementProfileID: types.Int64Value(input.ManagementProfileID),
		ProfileID:           types.StringValue(input.ProfileID),
		AgentName:           types.StringValue(input.AgentName),
		Version:             types.StringValue(input.Version),
		ManagementURL:       types.StringValue(input.ManagementURL),
		HasCreds:            types.Int64Value(input.HasCreds),
		Status:              types.Int64Value(input.Status),
		StatusDateTime:      types.StringValue(input.StatusDateTime),
	}
}

// newOmeDeviceLocation converts DeviceLocation to OmeDeviceLocation
func newOmeDeviceLocation(input DeviceLocation) OmeDeviceLocation {
	return OmeDeviceLocation{
		ID:                   types.Int64Value(input.ID),
		Room:                 types.StringValue(input.Room),
		Rack:                 types.StringValue(input.Rack),
		Aisle:                types.StringValue(input.Aisle),
		Datacenter:           types.StringValue(input.Datacenter),
		Rackslot:             types.StringValue(input.Rackslot),
		ManagementSystemUnit: types.Int64Value(input.ManagementSystemUnit),
	}
}

// newOmeDeviceFru converts DeviceFru to OmeDeviceFru
func newOmeDeviceFru(input DeviceFru) OmeDeviceFru {
	return OmeDeviceFru{
		Revision:     types.StringValue(input.Revision),
		ID:           types.Int64Value(input.ID),
		Manufacturer: types.StringValue(input.Manufacturer),
		Name:         types.StringValue(input.Name),
		PartNumber:   types.StringValue(input.PartNumber),
		SerialNumber: types.StringValue(input.SerialNumber),
	}
}

// newOmeDeviceCapability converts DeviceCapability to OmeDeviceCapability
func newOmeDeviceCapability(input DeviceCapability) OmeDeviceCapability {
	return OmeDeviceCapability{
		ID:             types.Int64Value(input.ID),
		CapabilityType: newOmeDeviceCapabilityType(input.CapabilityType),
	}
}

// newOmeDeviceCapabilityType converts DeviceCapabilityType to OmeDeviceCapabilityType
func newOmeDeviceCapabilityType(input DeviceCapabilityType) OmeDeviceCapabilityType {
	return OmeDeviceCapabilityType{
		CapabilityID: types.Int64Value(input.CapabilityID),
		Name:         types.StringValue(input.Name),
		Description:  types.StringValue(input.Description),
		IDOwner:      types.Int64Value(input.IDOwner),
	}
}

// newOmeServerDeviceCardInfo converts ServerDeviceCardInfo to OmeServerDeviceCardInfo
func newOmeServerDeviceCardInfo(input ServerDeviceCardInfo) OmeServerDeviceCardInfo {
	return OmeServerDeviceCardInfo{
		ID:           types.Int64Value(input.ID),
		SlotNumber:   types.StringValue(input.SlotNumber),
		Manufacturer: types.StringValue(input.Manufacturer),
		Description:  types.StringValue(input.Description),
		DatabusWidth: types.StringValue(input.DatabusWidth),
		SlotLength:   types.StringValue(input.SlotLength),
		SlotType:     types.StringValue(input.SlotType),
	}
}

// newOmeCPUInfo converts CPUInfo to OmeCPUInfo
func newOmeCPUInfo(input CPUInfo) OmeCPUInfo {
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

// newOmePartition converts Partition to OmePartition
func newOmePartition(input Partition) OmePartition {
	return OmePartition{
		Fqdd:                     types.StringValue(input.Fqdd),
		CurrentMacAddress:        types.StringValue(input.CurrentMacAddress),
		PermanentMacAddress:      types.StringValue(input.PermanentMacAddress),
		PermanentIscsiMacAddress: types.StringValue(input.PermanentIscsiMacAddress),
		PermanentFcoeMacAddress:  types.StringValue(input.PermanentFcoeMacAddress),
		Wwn:                      types.StringValue(input.Wwn),
		Wwpn:                     types.StringValue(input.Wwpn),
		VirtualWwn:               types.StringValue(input.VirtualWwn),
		VirtualWwpn:              types.StringValue(input.VirtualWwpn),
		VirtualMacAddress:        types.StringValue(input.VirtualMacAddress),
		NicMode:                  types.StringValue(input.NicMode),
		FcoeMode:                 types.StringValue(input.FcoeMode),
		IscsiMode:                types.StringValue(input.IscsiMode),
		MinBandwidth:             types.Int64Value(input.MinBandwidth),
		MaxBandwidth:             types.Int64Value(input.MaxBandwidth),
	}
}

// newOmePort converts Port to OmePort
func newOmePort(input Port) OmePort {
	return OmePort{
		PortID:      types.StringValue(input.PortID),
		ProductName: types.StringValue(input.ProductName),
		LinkStatus:  types.StringValue(input.LinkStatus),
		LinkSpeed:   types.Int64Value(input.LinkSpeed),
		Partitions:  newOmePartitionList(input.Partitions),
	}
}

// newOmePartitionList converts list of Partition to list of OmePartition
func newOmePartitionList(inputs []Partition) []OmePartition {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmePartition, 0)
	for _, input := range inputs {
		out = append(out, newOmePartition(input))
	}
	return out
}

// newOmeNICInfo converts NICInfo to OmeNICInfo
func newOmeNICInfo(input NICInfo) OmeNICInfo {
	return OmeNICInfo{
		NicID:      types.StringValue(input.NicID),
		VendorName: types.StringValue(input.VendorName),
		Ports:      newOmePortList(input.Ports),
	}
}

// newOmePortList converts list of Port to list of OmePort
func newOmePortList(inputs []Port) []OmePort {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmePort, 0)
	for _, input := range inputs {
		out = append(out, newOmePort(input))
	}
	return out
}

// newOmeFCInfo converts FCInfo to OmeFCInfo
func newOmeFCInfo(input FCInfo) OmeFCInfo {
	return OmeFCInfo{
		ID:                 types.Int64Value(input.ID),
		Fqdd:               types.StringValue(input.Fqdd),
		DeviceDescription:  types.StringValue(input.DeviceDescription),
		DeviceName:         types.StringValue(input.DeviceName),
		FirstFctargetLun:   types.StringValue(input.FirstFctargetLun),
		FirstFctargetWwpn:  types.StringValue(input.FirstFctargetWwpn),
		PortNumber:         types.Int64Value(input.PortNumber),
		PortSpeed:          types.StringValue(input.PortSpeed),
		SecondFctargetLun:  types.StringValue(input.SecondFctargetLun),
		SecondFctargetWwpn: types.StringValue(input.SecondFctargetWwpn),
		VendorName:         types.StringValue(input.VendorName),
		Wwn:                types.StringValue(input.Wwn),
		Wwpn:               types.StringValue(input.Wwpn),
		LinkStatus:         types.StringValue(input.LinkStatus),
		VirtualWwn:         types.StringValue(input.VirtualWwn),
		VirtualWwpn:        types.StringValue(input.VirtualWwpn),
	}
}

// newOmeOSInfo converts OSInfo to OmeOSInfo
func newOmeOSInfo(input OSInfo) OmeOSInfo {
	return OmeOSInfo{
		ID:        types.Int64Value(input.ID),
		OsName:    types.StringValue(input.OsName),
		OsVersion: types.StringValue(input.OsVersion),
		Hostname:  types.StringValue(input.Hostname),
	}
}

// newOmePowerSupplyInfo converts PowerSupplyInfo to OmePowerSupplyInfo
func newOmePowerSupplyInfo(input PowerSupplyInfo) OmePowerSupplyInfo {
	return OmePowerSupplyInfo{
		ID:                                  types.Int64Value(input.ID),
		Name:                                types.StringValue(input.Name),
		PowerSupplyType:                     types.Int64Value(input.PowerSupplyType),
		OutputWatts:                         types.Int64Value(input.OutputWatts),
		Location:                            types.StringValue(input.Location),
		RedundancyState:                     types.StringValue(input.RedundancyState),
		Status:                              types.Int64Value(input.Status),
		State:                               types.StringValue(input.State),
		FirmwareVersion:                     types.StringValue(input.FirmwareVersion),
		InputVoltage:                        types.Int64Value(input.InputVoltage),
		Model:                               types.StringValue(input.Model),
		Manufacturer:                        types.StringValue(input.Manufacturer),
		Range1MaxInputPowerWatts:            types.Int64Value(input.Range1MaxInputPowerWatts),
		SerialNumber:                        types.StringValue(input.SerialNumber),
		ActiveInputVoltage:                  types.StringValue(input.ActiveInputVoltage),
		InputPowerUnits:                     types.StringValue(input.InputPowerUnits),
		OperationalStatus:                   types.StringValue(input.OperationalStatus),
		Range1MaxInputVoltageHighMilliVolts: types.Int64Value(input.Range1MaxInputVoltageHighMilliVolts),
		RatedMaxOutputPower:                 types.Int64Value(input.RatedMaxOutputPower),
		RequestedState:                      types.Int64Value(input.RequestedState),
		AcInput:                             types.BoolValue(input.AcInput),
		AcOutput:                            types.BoolValue(input.AcOutput),
		SwitchingSupply:                     types.BoolValue(input.SwitchingSupply),
	}
}

// newOmeDiskInfo converts DiskInfo to OmeDiskInfo
func newOmeDiskInfo(input DiskInfo) OmeDiskInfo {
	return OmeDiskInfo{
		ID:                          types.Int64Value(input.ID),
		DiskNumber:                  types.StringValue(input.DiskNumber),
		VendorName:                  types.StringValue(input.VendorName),
		Status:                      types.Int64Value(input.Status),
		StatusString:                types.StringValue(input.StatusString),
		ModelNumber:                 types.StringValue(input.ModelNumber),
		SerialNumber:                types.StringValue(input.SerialNumber),
		SasAddress:                  types.StringValue(input.SasAddress),
		Revision:                    types.StringValue(input.Revision),
		ManufacturedDay:             types.Int64Value(input.ManufacturedDay),
		ManufacturedWeek:            types.Int64Value(input.ManufacturedWeek),
		ManufacturedYear:            types.Int64Value(input.ManufacturedYear),
		EncryptionAbility:           types.BoolValue(input.EncryptionAbility),
		FormFactor:                  types.StringValue(input.FormFactor),
		PartNumber:                  types.StringValue(input.PartNumber),
		PredictiveFailureState:      types.StringValue(input.PredictiveFailureState),
		EnclosureID:                 types.StringValue(input.EnclosureID),
		Channel:                     types.Int64Value(input.Channel),
		Size:                        types.StringValue(input.Size),
		FreeSpace:                   types.StringValue(input.FreeSpace),
		UsedSpace:                   types.StringValue(input.UsedSpace),
		BusType:                     types.StringValue(input.BusType),
		SlotNumber:                  types.Int64Value(input.SlotNumber),
		MediaType:                   types.StringValue(input.MediaType),
		RemainingReadWriteEndurance: types.StringValue(input.RemainingReadWriteEndurance),
		SecurityState:               types.StringValue(input.SecurityState),
		RaidStatus:                  types.StringValue(input.RaidStatus),
	}
}

// newOmeServerVirtualDisk converts ServerVirtualDisk to OmeServerVirtualDisk
func newOmeServerVirtualDisk(input ServerVirtualDisk) OmeServerVirtualDisk {
	return OmeServerVirtualDisk{
		ID:               types.Int64Value(input.ID),
		RaidControllerID: types.Int64Value(input.RaidControllerID),
		DeviceID:         types.Int64Value(input.DeviceID),
		Fqdd:             types.StringValue(input.Fqdd),
		State:            types.StringValue(input.State),
		RollupStatus:     types.Int64Value(input.RollupStatus),
		Status:           types.Int64Value(input.Status),
		Layout:           types.StringValue(input.Layout),
		MediaType:        types.StringValue(input.MediaType),
		Name:             types.StringValue(input.Name),
		ReadPolicy:       types.StringValue(input.ReadPolicy),
		WritePolicy:      types.StringValue(input.WritePolicy),
		CachePolicy:      types.StringValue(input.CachePolicy),
		StripeSize:       types.StringValue(input.StripeSize),
		Size:             types.StringValue(input.Size),
		TargetID:         types.Int64Value(input.TargetID),
		LockStatus:       types.StringValue(input.LockStatus),
	}
}

// newOmeRAIDControllerInfo converts RAIDControllerInfo to OmeRAIDControllerInfo
func newOmeRAIDControllerInfo(input RAIDControllerInfo) OmeRAIDControllerInfo {
	return OmeRAIDControllerInfo{
		ID:                       types.Int64Value(input.ID),
		Name:                     types.StringValue(input.Name),
		Fqdd:                     types.StringValue(input.Fqdd),
		DeviceDescription:        types.StringValue(input.DeviceDescription),
		Status:                   types.Int64Value(input.Status),
		StatusType:               types.StringValue(input.StatusTypeString),
		RollupStatus:             types.Int64Value(input.RollupStatus),
		RollupStatusString:       types.StringValue(input.RollupStatusString),
		FirmwareVersion:          types.StringValue(input.FirmwareVersion),
		CacheSizeInMb:            types.Int64Value(input.CacheSizeInMb),
		PciSlot:                  types.StringValue(input.PciSlot),
		DriverVersion:            types.StringValue(input.DriverVersion),
		StorageAssignmentAllowed: types.StringValue(input.StorageAssignmentAllowed),
		ServerVirtualDisks:       newOmeServerVirtualDiskList(input.ServerVirtualDisks),
	}
}

// newOmeServerVirtualDiskList converts list of ServerVirtualDisk to list of OmeServerVirtualDisk
func newOmeServerVirtualDiskList(inputs []ServerVirtualDisk) []OmeServerVirtualDisk {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]OmeServerVirtualDisk, 0)
	for _, input := range inputs {
		out = append(out, newOmeServerVirtualDisk(input))
	}
	return out
}

// newOmeMemoryInfo converts MemoryInfo to OmeMemoryInfo
func newOmeMemoryInfo(input MemoryInfo) OmeMemoryInfo {
	return OmeMemoryInfo{
		ID:                    types.Int64Value(input.ID),
		Name:                  types.StringValue(input.Name),
		BankName:              types.StringValue(input.BankName),
		Size:                  types.Int64Value(input.Size),
		Status:                types.Int64Value(input.Status),
		Manufacturer:          types.StringValue(input.Manufacturer),
		PartNumber:            types.StringValue(input.PartNumber),
		SerialNumber:          types.StringValue(input.SerialNumber),
		TypeDetails:           types.StringValue(input.TypeDetails),
		ManufacturerDate:      types.StringValue(input.ManufacturerDate),
		Speed:                 types.Int64Value(input.Speed),
		CurrentOperatingSpeed: types.Int64Value(input.CurrentOperatingSpeed),
		Rank:                  types.StringValue(input.Rank),
		InstanceID:            types.StringValue(input.InstanceID),
		DeviceDescription:     types.StringValue(input.DeviceDescription),
	}
}

// newOmeStorageEnclosureInfo converts StorageEnclosureInfo to OmeStorageEnclosureInfo
func newOmeStorageEnclosureInfo(input StorageEnclosureInfo) OmeStorageEnclosureInfo {
	return OmeStorageEnclosureInfo{
		ID:               types.Int64Value(input.ID),
		Name:             types.StringValue(input.Name),
		Status:           types.Int64Value(input.Status),
		StatusTypeString: types.StringValue(input.StatusTypeString),
		ChannelNumber:    types.StringValue(input.ChannelNumber),
		BackplanePartNum: types.StringValue(input.BackplanePartNum),
		NumberOfFanPacks: types.Int64Value(input.NumberOfFanPacks),
		Version:          types.StringValue(input.Version),
		RollupStatus:     types.Int64Value(input.RollupStatus),
		SlotCount:        types.Int64Value(input.SlotCount),
	}
}

// newOmeServerPowerState converts ServerPowerState to OmeServerPowerState
func newOmeServerPowerState(input ServerPowerState) OmeServerPowerState {
	return OmeServerPowerState{
		ID:         types.Int64Value(input.ID),
		PowerState: types.Int64Value(input.PowerState),
	}
}

// newOmeDeviceLicense converts DeviceLicense to OmeDeviceLicense
func newOmeDeviceLicense(input DeviceLicense) OmeDeviceLicense {
	return OmeDeviceLicense{
		SoldDate:           types.StringValue(input.SoldDate),
		LicenseBound:       types.Int64Value(input.LicenseBound),
		EvalTimeRemaining:  types.Int64Value(input.EvalTimeRemaining),
		AssignedDevices:    types.StringValue(input.AssignedDevices),
		LicenseStatus:      types.Int64Value(input.LicenseStatus),
		EntitlementID:      types.StringValue(input.EntitlementID),
		LicenseDescription: types.StringValue(input.LicenseDescription),
		LicenseType:        newOmeLicenseType(input.LicenseType),
	}
}

// newOmeLicenseType converts LicenseType to OmeLicenseType
func newOmeLicenseType(input LicenseType) OmeLicenseType {
	return OmeLicenseType{
		Name:      types.StringValue(input.Name),
		LicenseID: types.Int64Value(input.LicenseID),
	}
}
