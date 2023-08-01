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
	"net"

	"github.com/netdata/go.d.plugin/pkg/iprange"
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
func (d *Device) BelongsToPool(pool iprange.Pool) bool {
	for _, devM := range d.DeviceManagement {
		if pool.Contains(devM.NetworkAddress) {
			return true
		}
	}
	return false
}

type DiscoveryConfigurationJob struct {
	GroupID          string `json:"GroupId"`
	CreatedBy        string `json:"CreatedBy"`
	DiscoveryJobName string `json:"DiscoveryJobName"`
}

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
