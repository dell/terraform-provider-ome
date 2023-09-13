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
	"regexp"
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
				Config:      testGetDevicesWithEmptyIP,
				ExpectError: regexp.MustCompile(".*filters.ip_expressions.*string length must be at least 1.*"),
				PlanOnly:    true,
			},
			{
				Config:      testGetDevicesWithEmptySvcTag,
				ExpectError: regexp.MustCompile(".*filters.device_service_tags.*string length must be at least 1.*"),
				PlanOnly:    true,
			},
			{
				Config:      testGetDevicesWithEmptyQuerys,
				ExpectError: regexp.MustCompile(".*Attribute filters.filter_expression string length must be at least 1.*"),
				PlanOnly:    true,
			},
			{
				Config:      testGetDevicesWithNoIDs,
				ExpectError: regexp.MustCompile(".*Attribute filters.ids list must contain at least 1 elements.*"),
				PlanOnly:    true,
			},
			{
				Config:      testGetDevicesWithNoIPs,
				ExpectError: regexp.MustCompile(".*Attribute filters.ip_expressions list must contain at least 1 elements.*"),
				PlanOnly:    true,
			},
			{
				Config:      testGetDevicesWithNoSvcTags,
				ExpectError: regexp.MustCompile(".*Attribute filters.device_service_tags list must contain at least 1 elements.*"),
				PlanOnly:    true,
			},
			{
				Config:      testGetDevicesWithInvalidInventory,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Match.*"),
			},
			{
				Config:      testGetInvalidDevices,
				ExpectError: regexp.MustCompile(".*Error fetching devices.*"),
			},
			{
				Config: testGetAllDevicesWithoutFilters + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetAllDevices + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetDevicesWithIDs + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetDevicesWithSTags + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetDevicesWithIPs + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "false"),
					resource.TestCheckOutput("fetched_inventory", "true"),
				),
			},
			{
				Config: testGetDevicesWithFilterExp + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "true"),
					resource.TestCheckOutput("fetched_inventory", "false"),
				),
			},
			{
				Config: testGetDevicesWithIPAndInvType + devDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_multiple", "false"),
					resource.TestCheckOutput("fetched_inventory", "true"),
					resource.TestCheckOutput("fetched_selected_inventory", "true"),
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
var testGetDevicesWithNoIDs = testProvider + `
data "ome_device" "devs" {
	filters = {
		ids = []
	}
}
`

var testGetDevicesWithNoSvcTags = testProvider + `
data "ome_device" "devs" {
	filters = {
		device_service_tags = []
	}
}
`

var testGetDevicesWithNoIPs = testProvider + `
data "ome_device" "devs" {
	filters = {
		ip_expressions = []
	}
}
`

var testGetDevicesWithEmptySvcTag = testProvider + `
data "ome_device" "devs" {
	filters = {
		device_service_tags = [""]
	}
}
`

var testGetDevicesWithEmptyIP = testProvider + `
data "ome_device" "devs" {
	filters = {
		ip_expressions = [""]
	}
}
`

var testGetDevicesWithEmptyQuerys = testProvider + `
data "ome_device" "devs" {
	filters = {
		filter_expression = ""
	}
}
`

var testGetAllDevices = testProvider + `
data "ome_device" "devs" {
	filters = {}
}
`

var testGetAllDevicesWithoutFilters = testProvider + `
data "ome_device" "devs" {
}
`

var testGetDevicesWithIDs = testProvider + `
data "ome_device" "devs" {
	filters = {
		ids = [` + DeviceID1 + `, ` + DeviceID2 + `]
	}
}
`

var testGetDevicesWithInvalidInventory = testProvider + `
data "ome_device" "devs" {
	filters = {
		device_service_tags = ["` + DeviceSvcTag1 + `", "` + DeviceSvcTag2 + `"]
	}	
	inventory_types = ["invalid"]
}
`

var testGetDevicesWithSTags = testProvider + `
data "ome_device" "devs" {
	filters = {
		device_service_tags = ["` + DeviceSvcTag1 + `", "` + DeviceSvcTag2 + `"]
	}
}
`

var testGetInvalidDevices = testProvider + `
data "ome_device" "devs" {
	filters = {
		device_service_tags = ["invalid"]
	}
}
`

var testGetDevicesWithIPs = testProvider + `
data "ome_device" "devs" {
	filters = {
		ip_expressions = ["` + DeviceIP1 + `"]
	}
}
`

var testGetDevicesWithFilterExp = testProvider + `
data "ome_device" "devs" {
	filters = {
		filter_expression = "Model eq 'PowerEdge MX840c'"
	}
}
`
var testGetDevicesWithIPAndInvType = testProvider + `
data "ome_device" "devs" {
	filters = {
		ip_expressions = ["` + DeviceIP1 + `"]
	}
	inventory_types = ["serverNetworkInterfaces","serverArrayDisks"]
}
output "fetched_selected_inventory" {
	value = alltrue(flatten([for dev in data.ome_device.devs.devices: 
		[
			dev.detailed_inventory.server_device_cards == null,
			dev.detailed_inventory.cpus == null,
			dev.detailed_inventory.nics != null,
			dev.detailed_inventory.fcis == null,
			dev.detailed_inventory.os == null,
			dev.detailed_inventory.power_supply == null,
			dev.detailed_inventory.disks != null,
			dev.detailed_inventory.raid_controllers == null,
			dev.detailed_inventory.memory == null,
			dev.detailed_inventory.storage_enclosures == null,
			dev.detailed_inventory.power_state == null,
			dev.detailed_inventory.licenses == null,
			dev.detailed_inventory.capabilities == null,
			dev.detailed_inventory.frus == null,
			dev.detailed_inventory.locations == null,
			dev.detailed_inventory.management_info == null,
			dev.detailed_inventory.softwares == null,
			dev.detailed_inventory.subsytem_rollup_status  == null
		]
	]))
}
`
