/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSource_ReadGroupDevices(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testgroupDeviceDS,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ome_groupdevices_info.gd", "id"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_ids.#", "2"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_servicetags.#", "2"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_groups.%", "1"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_groups.test_device_group.name", "test_device_group"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_groups.test_device_group.devices.#", "2"),
					resource.TestCheckResourceAttrSet("data.ome_groupdevices_info.gd", "device_groups.test_device_group.id"),
					resource.TestCheckResourceAttrSet("data.ome_groupdevices_info.gd", "device_groups.test_device_group.parent_id"),
					resource.TestCheckResourceAttrSet("data.ome_groupdevices_info.gd", "device_groups.test_device_group.description"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_groups.test_device_group.visible", "true"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_groups.test_device_group.sub_groups.#", "0"),
				),
			},

			{
				// create subGroup, but it wont reflect in datasource yet
				Config: testgroupDeviceDSWithSubGroups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ome_groupdevices_info.gd", "id"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_groups.test_device_group.sub_groups.#", "0"),
				),
			},

			{
				// now it will reflect
				Config: testgroupDeviceDSWithSubGroups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ome_groupdevices_info.gd", "id"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_groups.test_device_group.sub_groups.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.ome_groupdevices_info.gd",
						"device_groups.test_device_group.sub_groups.*",
						map[string]string{
							"name": "Linux",
						},
					),
				),
			},
			{
				Config:      testgroupDeviceDSEmptyGroup,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(".*at least 1 elements, got: 0.*"),
			},
			{
				Config:      testgroupDeviceDSInvalidGroup,
				ExpectError: regexp.MustCompile(".*no items found, expecting one.*"),
			},
		},
	})
}

var testgroupPreReq = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_groupdevices_info" "all_father" {
		id = "0"
		device_group_names = ["Static Groups"]
	}

	data "ome_device" "devs" {
		filters = {
			device_service_tags = ["` + DeviceSvcTag1 + `", "` + DeviceSvcTag2 + `"]
		}
	}

	resource "ome_static_group" "test_group" {
		name        = "test_device_group"
		description = "Group for Group Device DataSource Acceptance Test"
		parent_id   = data.ome_groupdevices_info.all_father.device_groups["Static Groups"].id
		device_ids = data.ome_device.devs.devices[*].id
	}
`

var testgroupDeviceDS = testgroupPreReq + `
	data "ome_groupdevices_info" "gd" {
		depends_on = [ ome_static_group.test_group ]
		id = "0"
		device_group_names = ["test_device_group"]
	}
`

var testgroupDeviceDSWithSubGroups = testgroupPreReq + `
	data "ome_groupdevices_info" "gd" {
		depends_on = [ ome_static_group.test_group ]
		id = "0"
		device_group_names = ["test_device_group"]
	}

	resource "ome_static_group" "linux-group" {
		name        = "Linux"
		description = "All linux servers"
		parent_id   = data.ome_groupdevices_info.gd.device_groups["test_device_group"].id
		device_ids = []
	}
`

var testgroupDeviceDSInvalidGroup = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_groupdevices_info" "gd" {
		id = "0"
		device_group_names = ["NO_GROUP"]
	}
`

var testgroupDeviceDSEmptyGroup = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_groupdevices_info" "gd" {
		id = "0"
		device_group_names = []
	}
`
