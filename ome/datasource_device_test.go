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

package ome

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSource_ReadDevice(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testGetAllDevicesWithoutFilters + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ome_device.devs", "devices"),
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetAllDevices + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ome_device.devs", "devices"),
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetDevicesWithIDs + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ome_device.devs", "devices"),
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetDevicesWithSTags + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ome_device.devs", "devices"),
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetDevicesWithIPs + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ome_device.devs", "devices"),
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "false"),
					resource.TestCheckOutput("fetched_inventory", "true"),
				),
			},
			{
				Config: testGetDevicesWithFilterExp + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ome_device.devs", "devices"),
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
		},
	})
}

var devDataOut = `
output "fetched_any" {
	value = length(data.ome_device.devs.devices) != 0
}
output "fetched_multiple" {
	value = length(data.ome_device.devs.devices) > 1
}
output "fetched_inventory" {
	value = anytrue(flatten([for dev in data.ome_device.devs.devices: 
		[
		dev.detailed_inventory.server_device_cards != null,
		dev.detailed_inventory.cpus != null,
		dev.detailed_inventory.nics != null,
		dev.detailed_inventory.fcis != null,
		dev.detailed_inventory.os != null,
		dev.detailed_inventory.power_supply != null,
		dev.detailed_inventory.disks != null,
		dev.detailed_inventory.raid_controllers != null,
		dev.detailed_inventory.memory != null,
		dev.detailed_inventory.storage_enclosures != null,
		dev.detailed_inventory.power_state != null,
		dev.detailed_inventory.licenses != null,
		dev.detailed_inventory.capabilities != null,
		dev.detailed_inventory.frus != null,
		dev.detailed_inventory.locations != null,
		dev.detailed_inventory.management_info != null,
		dev.detailed_inventory.softwares != null,
		dev.detailed_inventory.subsytem_rollup_status != null
		]
	]))
}
`

var old = `
output "omitted_inventory" {
	value = anytrue([for dev in data.ome_device.devs.devices: 
		dev.detailed_inventory.server_device_cards == null &&
		dev.detailed_inventory.cpus == null &&
		dev.detailed_inventory.nics == null &&
		dev.detailed_inventory.fcis == null &&
		dev.detailed_inventory.os == null &&
		dev.detailed_inventory.power_supply == null &&
		dev.detailed_inventory.disks == null &&
		dev.detailed_inventory.raid_controllers == null &&
		dev.detailed_inventory.memory == null &&
		dev.detailed_inventory.storage_enclosures == null &&
		dev.detailed_inventory.power_state == null &&
		dev.detailed_inventory.licenses == null &&
		dev.detailed_inventory.capabilities == null &&
		dev.detailed_inventory.frus == null &&
		dev.detailed_inventory.locations == null &&
		dev.detailed_inventory.management_info == null &&
		dev.detailed_inventory.softwares == null &&
		dev.detailed_inventory.subsytem_rollup_status  == null
		])
}
`

var testGetAllDevices = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}
data "ome_device" "devs" {
	filters = {}
}
`

var testGetAllDevicesWithoutFilters = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}
data "ome_device" "devs" {
}
`

var testGetDevicesWithIDs = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}
data "ome_device" "devs" {
	filters = {
		ids = [10112, 13528]
	}
}
`

var testGetDevicesWithSTags = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}
data "ome_device" "devs" {
	filters = {
		device_service_tags = ["CZNF1T2", "CZMC1T2"]
	}
}
`

var testGetDevicesWithIPs = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}
data "ome_device" "devs" {
	filters = {
		ip_expressions = ["10.226.197.113"]
	}
}
`

var testGetDevicesWithFilterExp = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}
data "ome_device" "devs" {
	filters = {
		filter_expression = "Model eq 'PowerEdge MX840c'"
	}
}
`
