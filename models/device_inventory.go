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
	"encoding/json"
	"fmt"
)

// DeviceInventory - device inventory response from OME
type DeviceInventory struct {
	ServerDeviceCards     []ServerDeviceCardInfo
	CPUInfo               []CPUInfo
	NICInfo               []NICInfo
	FCInfo                []FCInfo
	OSInfo                []OSInfo
	PowerSupplyInfo       []PowerSupplyInfo
	DiskInfo              []DiskInfo
	RAIDControllerInfo    []RAIDControllerInfo
	MemoryInfo            []MemoryInfo
	StorageEnclosureInfo  []StorageEnclosureInfo
	ServerPowerStates     []ServerPowerState
	DeviceLicenses        []DeviceLicense
	DeviceCapabilities    []DeviceCapability
	DeviceFrus            []DeviceFru
	DeviceLocations       []DeviceLocation
	DeviceManagement      []DeviceManagementInfo
	DeviceSoftwares       []DeviceSoftware
	SubSystemRollupStatus []SubSystemRollupStatus
}

// DeviceInventoryInfo - A temporary struct to hold the JSON data for device inventory
type DeviceInventoryInfo struct {
	Type string          `json:"InventoryType"`
	Info json.RawMessage `json:"InventoryInfo"`
}

// UnmarshalJSON - Implement custom json unmarshalling for DeviceInventory
func (d *DeviceInventory) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON data into the temporary struct
	temp := make([]DeviceInventoryInfo, 0)
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("error unmarshalling device type: %w", err)
	}

	// Assign the values from the temporary struct to the actual Person struct
	for _, inv := range temp {
		err := d.AddInfo(inv)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddInfo - Append a DeviceInventoryInfo into a DeviceInventory
func (d *DeviceInventory) AddInfo(inv DeviceInventoryInfo) error {
	var err error
	switch inv.Type {
	case "serverDeviceCards":
		err = json.Unmarshal(inv.Info, &d.ServerDeviceCards)
	case "serverProcessors":
		err = json.Unmarshal(inv.Info, &d.CPUInfo)
	case "serverNetworkInterfaces":
		err = json.Unmarshal(inv.Info, &d.NICInfo)
	case "serverFcCards":
		err = json.Unmarshal(inv.Info, &d.FCInfo)
	case "serverOperatingSystems":
		err = json.Unmarshal(inv.Info, &d.OSInfo)
	case "serverPowerSupplies":
		err = json.Unmarshal(inv.Info, &d.PowerSupplyInfo)
	case "serverArrayDisks":
		err = json.Unmarshal(inv.Info, &d.DiskInfo)
	case "serverRaidControllers":
		err = json.Unmarshal(inv.Info, &d.RAIDControllerInfo)
	case "serverMemoryDevices":
		err = json.Unmarshal(inv.Info, &d.MemoryInfo)
	case "serverStorageEnclosures":
		err = json.Unmarshal(inv.Info, &d.StorageEnclosureInfo)
	case "serverSupportedPowerStates":
		err = json.Unmarshal(inv.Info, &d.ServerPowerStates)
	case "deviceLicense":
		err = json.Unmarshal(inv.Info, &d.DeviceLicenses)
	case "deviceCapabilities":
		err = json.Unmarshal(inv.Info, &d.DeviceCapabilities)
	case "deviceFru":
		err = json.Unmarshal(inv.Info, &d.DeviceFrus)
	case "deviceLocation":
		err = json.Unmarshal(inv.Info, &d.DeviceLocations)
	case "deviceManagement":
		err = json.Unmarshal(inv.Info, &d.DeviceManagement)
	case "deviceSoftware":
		err = json.Unmarshal(inv.Info, &d.DeviceSoftwares)
	case "subsystemRollupStatus":
		err = json.Unmarshal(inv.Info, &d.SubSystemRollupStatus)
	}
	if err != nil {
		err = fmt.Errorf("error unmarshalling inventory value: %w", err)
	}

	return err
}

// NewDeviceInventory - Creates blank inventory struct
func NewDeviceInventory() DeviceInventory {
	return DeviceInventory{
		ServerDeviceCards:     make([]ServerDeviceCardInfo, 0),
		CPUInfo:               make([]CPUInfo, 0),
		NICInfo:               make([]NICInfo, 0),
		FCInfo:                make([]FCInfo, 0),
		OSInfo:                make([]OSInfo, 0),
		PowerSupplyInfo:       make([]PowerSupplyInfo, 0),
		DiskInfo:              make([]DiskInfo, 0),
		RAIDControllerInfo:    make([]RAIDControllerInfo, 0),
		MemoryInfo:            make([]MemoryInfo, 0),
		StorageEnclosureInfo:  make([]StorageEnclosureInfo, 0),
		ServerPowerStates:     make([]ServerPowerState, 0),
		DeviceLicenses:        make([]DeviceLicense, 0),
		DeviceCapabilities:    make([]DeviceCapability, 0),
		DeviceFrus:            make([]DeviceFru, 0),
		DeviceLocations:       make([]DeviceLocation, 0),
		DeviceManagement:      make([]DeviceManagementInfo, 0),
		DeviceSoftwares:       make([]DeviceSoftware, 0),
		SubSystemRollupStatus: make([]SubSystemRollupStatus, 0),
	}
}

// AddInventory - Merges given inventory to self
func (d *DeviceInventory) AddInventory(dd DeviceInventory) {
	d.ServerDeviceCards = append(d.ServerDeviceCards, dd.ServerDeviceCards...)
	d.CPUInfo = append(d.CPUInfo, dd.CPUInfo...)
	d.NICInfo = append(d.NICInfo, dd.NICInfo...)
	d.FCInfo = append(d.FCInfo, dd.FCInfo...)
	d.OSInfo = append(d.OSInfo, dd.OSInfo...)
	d.PowerSupplyInfo = append(d.PowerSupplyInfo, dd.PowerSupplyInfo...)
	d.DiskInfo = append(d.DiskInfo, dd.DiskInfo...)
	d.RAIDControllerInfo = append(d.RAIDControllerInfo, dd.RAIDControllerInfo...)
	d.MemoryInfo = append(d.MemoryInfo, dd.MemoryInfo...)
	d.StorageEnclosureInfo = append(d.StorageEnclosureInfo, dd.StorageEnclosureInfo...)
	d.ServerPowerStates = append(d.ServerPowerStates, dd.ServerPowerStates...)
	d.DeviceLicenses = append(d.DeviceLicenses, dd.DeviceLicenses...)
	d.DeviceCapabilities = append(d.DeviceCapabilities, dd.DeviceCapabilities...)
	d.DeviceFrus = append(d.DeviceFrus, dd.DeviceFrus...)
	d.DeviceLocations = append(d.DeviceLocations, dd.DeviceLocations...)
	d.DeviceManagement = append(d.DeviceManagement, dd.DeviceManagement...)
	d.DeviceSoftwares = append(d.DeviceSoftwares, dd.DeviceSoftwares...)
	d.SubSystemRollupStatus = append(d.SubSystemRollupStatus, dd.SubSystemRollupStatus...)
}

// SubSystemRollupStatus - SubSystemRollupStatus
type SubSystemRollupStatus struct {
	ID            int64  `json:"Id"`
	Status        int64  `json:"Status"`
	SubsystemName string `json:"SubsystemName"`
}

// DeviceSoftware - DeviceSoftware
type DeviceSoftware struct {
	Version           string `json:"Version"`
	InstallationDate  string `json:"InstallationDate"`
	Status            string `json:"Status"`
	SoftwareType      string `json:"SoftwareType"`
	VendorID          string `json:"VendorId"`
	SubDeviceID       string `json:"SubDeviceId"`
	SubVendorID       string `json:"SubVendorId"`
	ComponentID       string `json:"ComponentId"`
	PciDeviceID       string `json:"PciDeviceId"`
	DeviceDescription string `json:"DeviceDescription"`
	InstanceID        string `json:"InstanceId"`
}

// DeviceManagementInfo - DeviceManagementInfo
type DeviceManagementInfo struct {
	ManagementID        int64           `json:"ManagementId"`
	IPAddress           string          `json:"IpAddress"`
	MACAddress          string          `json:"MacAddress"`
	InstrumentationName string          `json:"InstrumentationName"`
	DNSName             string          `json:"DnsName"`
	ManagementType      ManagementType  `json:"ManagementType"`
	EndPointAgents      []EndPointAgent `json:"EndPointAgents"`
}

// ManagementType - ManagementType
type ManagementType struct {
	ManagementType int64  `json:"ManagementType"`
	Name           string `json:"Name"`
	Description    string `json:"Description"`
}

// EndPointAgent - EndPointAgent
type EndPointAgent struct {
	ManagementProfileID int64  `json:"ManagementProfileId"`
	ProfileID           string `json:"ProfileId"`
	AgentName           string `json:"AgentName"`
	Version             string `json:"Version"`
	ManagementURL       string `json:"ManagementURL"`
	HasCreds            int64  `json:"HasCreds"`
	Status              int64  `json:"Status"`
	StatusDateTime      string `json:"StatusDateTime"`
}

// DeviceLocation - DeviceLocation
type DeviceLocation struct {
	ID                   int64  `json:"Id"`
	Room                 string `json:"Room"`
	Rack                 string `json:"Rack"`
	Aisle                string `json:"Aisle"`
	Datacenter           string `json:"Datacenter"`
	Rackslot             string `json:"Rackslot"`
	ManagementSystemUnit int64  `json:"ManagementSystemUnit"`
}

// DeviceFru - Device FRU
type DeviceFru struct {
	Revision     string `json:"Revision"`
	ID           int64  `json:"Id"`
	Manufacturer string `json:"Manufacturer"`
	Name         string `json:"Name"`
	PartNumber   string `json:"PartNumber"`
	SerialNumber string `json:"SerialNumber"`
}

// DeviceCapability - Device Capability
type DeviceCapability struct {
	ID             int64                `json:"Id"`
	CapabilityType DeviceCapabilityType `json:"CapabilityType"`
}

// DeviceCapabilityType - Device Capability Type
type DeviceCapabilityType struct {
	CapabilityID int64  `json:"CapabilityId"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	IDOwner      int64  `json:"IdOwner"`
}

// ServerDeviceCardInfo - Server Device Card Info
type ServerDeviceCardInfo struct {
	ID           int64  `json:"Id"`
	SlotNumber   string `json:"SlotNumber"`
	Manufacturer string `json:"Manufacturer"`
	Description  string `json:"Description"`
	DatabusWidth string `json:"DatabusWidth"`
	SlotLength   string `json:"SlotLength"`
	SlotType     string `json:"SlotType"`
}

// CPUInfo - CPU Info
type CPUInfo struct {
	ID                   int64  `json:"Id"`
	Family               string `json:"Family"`
	MaxSpeed             int64  `json:"MaxSpeed"`
	CurrentSpeed         int64  `json:"CurrentSpeed"`
	SlotNumber           string `json:"SlotNumber"`
	Status               int64  `json:"Status"`
	NumberOfCores        int64  `json:"NumberOfCores"`
	NumberOfEnabledCores int64  `json:"NumberOfEnabledCores"`
	BrandName            string `json:"BrandName"`
	ModelName            string `json:"ModelName"`
	InstanceID           string `json:"InstanceId"`
	Voltage              string `json:"Voltage"`
}

// Partition - Partition
type Partition struct {
	Fqdd                     string `json:"Fqdd"`
	CurrentMacAddress        string `json:"CurrentMacAddress"`
	PermanentMacAddress      string `json:"PermanentMacAddress"`
	PermanentIscsiMacAddress string `json:"PermanentIscsiMacAddress"`
	PermanentFcoeMacAddress  string `json:"PermanentFcoeMacAddress"`
	Wwn                      string `json:"Wwn"`
	Wwpn                     string `json:"Wwpn"`
	VirtualWwn               string `json:"VirtualWwn"`
	VirtualWwpn              string `json:"VirtualWwpn"`
	VirtualMacAddress        string `json:"VirtualMacAddress"`
	NicMode                  string `json:"NicMode"`
	FcoeMode                 string `json:"FcoeMode"`
	IscsiMode                string `json:"IscsiMode"`
	MinBandwidth             int64  `json:"MinBandwidth"`
	MaxBandwidth             int64  `json:"MaxBandwidth"`
}

// Port - Port
type Port struct {
	PortID      string      `json:"PortId"`
	ProductName string      `json:"ProductName"`
	LinkStatus  string      `json:"LinkStatus"`
	LinkSpeed   int64       `json:"LinkSpeed"`
	Partitions  []Partition `json:"Partitions"`
}

// NICInfo - NIC Info
type NICInfo struct {
	NicID      string `json:"NicId"`
	VendorName string `json:"VendorName"`
	Ports      []Port `json:"Ports"`
}

// FCInfo - FC Info
type FCInfo struct {
	ID                 int64  `json:"Id"`
	Fqdd               string `json:"Fqdd"`
	DeviceDescription  string `json:"DeviceDescription"`
	DeviceName         string `json:"DeviceName"`
	FirstFctargetLun   string `json:"FirstFctargetLun"`
	FirstFctargetWwpn  string `json:"FirstFctargetWwpn"`
	PortNumber         int64  `json:"PortNumber"`
	PortSpeed          string `json:"PortSpeed"`
	SecondFctargetLun  string `json:"SecondFctargetLun"`
	SecondFctargetWwpn string `json:"SecondFctargetWwpn"`
	VendorName         string `json:"VendorName"`
	Wwn                string `json:"Wwn"`
	Wwpn               string `json:"Wwpn"`
	LinkStatus         string `json:"LinkStatus"`
	VirtualWwn         string `json:"VirtualWwn"`
	VirtualWwpn        string `json:"VirtualWwpn"`
}

// OSInfo - OS Info
type OSInfo struct {
	ID        int64  `json:"Id"`
	OsName    string `json:"OsName"`
	OsVersion string `json:"OsVersion"`
	Hostname  string `json:"Hostname"`
}

// PowerSupplyInfo - Power Supply Info
type PowerSupplyInfo struct {
	ID                                  int64  `json:"Id"`
	Name                                string `json:"Name"`
	PowerSupplyType                     int64  `json:"PowerSupplyType"`
	OutputWatts                         int64  `json:"OutputWatts"`
	Location                            string `json:"Location"`
	RedundancyState                     string `json:"RedundancyState"`
	Status                              int64  `json:"Status"`
	State                               string `json:"State"`
	FirmwareVersion                     string `json:"FirmwareVersion"`
	InputVoltage                        int64  `json:"InputVoltage"`
	Model                               string `json:"Model"`
	Manufacturer                        string `json:"Manufacturer"`
	Range1MaxInputPowerWatts            int64  `json:"Range1MaxInputPowerWatts"`
	SerialNumber                        string `json:"SerialNumber"`
	ActiveInputVoltage                  string `json:"ActiveInputVoltage"`
	InputPowerUnits                     string `json:"InputPowerUnits"`
	OperationalStatus                   string `json:"OperationalStatus"`
	Range1MaxInputVoltageHighMilliVolts int64  `json:"Range1MaxInputVoltageHighMilliVolts"`
	RatedMaxOutputPower                 int64  `json:"RatedMaxOutputPower"`
	RequestedState                      int64  `json:"RequestedState"`
	AcInput                             bool   `json:"AcInput"`
	AcOutput                            bool   `json:"AcOutput"`
	SwitchingSupply                     bool   `json:"SwitchingSupply"`
}

// DiskInfo - Disk Info
type DiskInfo struct {
	ID                          int64  `json:"Id"`
	DiskNumber                  string `json:"DiskNumber"`
	VendorName                  string `json:"VendorName"`
	Status                      int64  `json:"Status"`
	StatusString                string `json:"StatusString"`
	ModelNumber                 string `json:"ModelNumber"`
	SerialNumber                string `json:"SerialNumber"`
	SasAddress                  string `json:"SasAddress"`
	Revision                    string `json:"Revision"`
	ManufacturedDay             int64  `json:"ManufacturedDay"`
	ManufacturedWeek            int64  `json:"ManufacturedWeek"`
	ManufacturedYear            int64  `json:"ManufacturedYear"`
	EncryptionAbility           bool   `json:"EncryptionAbility"`
	FormFactor                  string `json:"FormFactor"`
	PartNumber                  string `json:"PartNumber"`
	PredictiveFailureState      string `json:"PredictiveFailureState"`
	EnclosureID                 string `json:"EnclosureId"`
	Channel                     int64  `json:"Channel"`
	Size                        string `json:"Size"`
	FreeSpace                   string `json:"FreeSpace"`
	UsedSpace                   string `json:"UsedSpace"`
	BusType                     string `json:"BusType"`
	SlotNumber                  int64  `json:"SlotNumber"`
	MediaType                   string `json:"MediaType"`
	RemainingReadWriteEndurance string `json:"RemainingReadWriteEndurance"`
	SecurityState               string `json:"SecurityState"`
	RaidStatus                  string `json:"RaidStatus"`
}

// ServerVirtualDisk - ServerVirtualDisk
type ServerVirtualDisk struct {
	ID               int64  `json:"Id"`
	RaidControllerID int64  `json:"RaidControllerId"`
	DeviceID         int64  `json:"DeviceId"`
	Fqdd             string `json:"Fqdd"`
	State            string `json:"State"`
	RollupStatus     int64  `json:"RollupStatus"`
	Status           int64  `json:"Status"`
	Layout           string `json:"Layout"`
	MediaType        string `json:"MediaType"`
	Name             string `json:"Name"`
	ReadPolicy       string `json:"ReadPolicy"`
	WritePolicy      string `json:"WritePolicy"`
	CachePolicy      string `json:"CachePolicy"`
	StripeSize       string `json:"StripeSize"`
	Size             string `json:"Size"`
	TargetID         int64  `json:"TargetId"`
	LockStatus       string `json:"LockStatus"`
}

// RAIDControllerInfo - RAID Controller Info
type RAIDControllerInfo struct {
	ID                       int64               `json:"Id"`
	Name                     string              `json:"Name"`
	Fqdd                     string              `json:"Fqdd"`
	DeviceDescription        string              `json:"DeviceDescription"`
	Status                   int64               `json:"Status"`
	StatusTypeString         string              `json:"StatusTypeString"`
	RollupStatus             int64               `json:"RollupStatus"`
	RollupStatusString       string              `json:"RollupStatusString"`
	FirmwareVersion          string              `json:"FirmwareVersion"`
	CacheSizeInMb            int64               `json:"CacheSizeInMb"`
	PciSlot                  string              `json:"PciSlot"`
	DriverVersion            string              `json:"DriverVersion"`
	StorageAssignmentAllowed string              `json:"StorageAssignmentAllowed"`
	ServerVirtualDisks       []ServerVirtualDisk `json:"ServerVirtualDisks"`
}

// MemoryInfo - Memory Info
type MemoryInfo struct {
	ID                    int64  `json:"Id"`
	Name                  string `json:"Name"`
	BankName              string `json:"BankName"`
	Size                  int64  `json:"Size"`
	Status                int64  `json:"Status"`
	Manufacturer          string `json:"Manufacturer"`
	PartNumber            string `json:"PartNumber"`
	SerialNumber          string `json:"SerialNumber"`
	TypeDetails           string `json:"TypeDetails"`
	ManufacturerDate      string `json:"ManufacturerDate"`
	Speed                 int64  `json:"Speed"`
	CurrentOperatingSpeed int64  `json:"CurrentOperatingSpeed"`
	Rank                  string `json:"Rank"`
	InstanceID            string `json:"InstanceId"`
	DeviceDescription     string `json:"DeviceDescription"`
}

// StorageEnclosureInfo - Storage Enclosure Info
type StorageEnclosureInfo struct {
	ID               int64  `json:"Id"`
	Name             string `json:"Name"`
	Status           int64  `json:"Status"`
	StatusTypeString string `json:"StatusTypeString"`
	ChannelNumber    string `json:"ChannelNumber"`
	BackplanePartNum string `json:"BackplanePartNum"`
	NumberOfFanPacks int64  `json:"NumberOfFanPacks"`
	Version          string `json:"Version"`
	RollupStatus     int64  `json:"RollupStatus"`
	SlotCount        int64  `json:"SlotCount"`
}

// ServerPowerState - ServerPowerState
type ServerPowerState struct {
	ID         int64 `json:"Id"`
	PowerState int64 `json:"PowerState"`
}

// DeviceLicense - Device License
type DeviceLicense struct {
	SoldDate           string      `json:"SoldDate"`
	LicenseBound       int64       `json:"LicenseBound"`
	EvalTimeRemaining  int64       `json:"EvalTimeRemaining"`
	AssignedDevices    string      `json:"AssignedDevices"`
	LicenseStatus      int64       `json:"LicenseStatus"`
	EntitlementID      string      `json:"EntitlementId"`
	LicenseDescription string      `json:"LicenseDescription"`
	LicenseType        LicenseType `json:"LicenseType"`
}

// LicenseType - License Type
type LicenseType struct {
	Name      string `json:"Name"`
	LicenseID int64  `json:"LicenseId"`
}
