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
	ID                 int64              `json:"Id"`
	DeviceServiceTag   string             `json:"DeviceServiceTag"`
	DeviceCapabilities []int64            `json:"DeviceCapabilities,omitempty"`
	DeviceManagement   []DeviceManagement `json:"DeviceManagement"`
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

// DeviceManagement - embedded device management response from the Devices
type DeviceManagement struct {
	NetworkAddress net.IP `json:"NetworkAddress"`
}
