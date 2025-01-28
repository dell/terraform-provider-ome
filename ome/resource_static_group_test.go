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
	"fmt"
	"regexp"
	"terraform-provider-ome/clients"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	DeviceGroup1       = "test_acc_group_device_1"
	DeviceGroup1Update = "test_acc_group_device_1_updated"
)

func TestAccStaticGroup(t *testing.T) {

	testAccProvider := testProvider

	preReqs := `
	data "ome_device" "devs" {
		filters = {
			device_service_tags = ["` + DeviceSvcTag1 + `", "` + DeviceSvcTag2 + `"]
		}
	}

	data "ome_groupdevices_info" "ome_root" {
		device_group_names = ["Static Groups"]
	}
	`

	testAccCreateGroupSuccess := testAccProvider + preReqs + `	
	resource "ome_static_group" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1"
		parent_id = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = data.ome_device.devs.devices[*].id
	}
	`

	testAccUpdateGroupSuccess := testAccProvider + preReqs + `
	resource "ome_static_group" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [data.ome_device.devs.devices[0].id]
	}
	`

	testAccDuplicateNameNeg := testAccProvider + preReqs + `
	resource "ome_static_group" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [data.ome_device.devs.devices[0].id]
	}
	resource "ome_static_group" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [data.ome_device.devs.devices[1].id]
	}
	`

	testAccInvalidDeviceNeg := testAccProvider + preReqs + `
	resource "ome_static_group" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [data.ome_device.devs.devices[0].id]
	}
	resource "ome_static_group" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [-1]
	}
	`

	testAccCreate2 := testAccProvider + preReqs + `
	resource "ome_static_group" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [data.ome_device.devs.devices[0].id]
	}
	resource "ome_static_group" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
	}
	`

	testAccUpdateMultipleDevices := testAccProvider + preReqs + `
	resource "ome_static_group" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [data.ome_device.devs.devices[0].id]
	}
	resource "ome_static_group" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = data.ome_device.devs.devices[*].id
	}
	`

	testAccNoDevices := testAccProvider + preReqs + `
	resource "ome_static_group" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = [data.ome_device.devs.devices[0].id]
	}
	resource "ome_static_group" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = []
	}
	`

	testGroupWoDesc := testAccProvider + preReqs + `
	resource "ome_static_group" "terraform-acceptance-test" {
		name       = "` + DeviceGroup1 + `"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = data.ome_device.devs.devices[*].id
	}
	`

	testUpdateParentID := testAccProvider + preReqs + `
	resource "ome_static_group" "parent" {
		name       = "` + DeviceGroup1 + `_parent"
		parent_id  = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
		device_ids = data.ome_device.devs.devices[*].id
	}
	resource "ome_static_group" "terraform-acceptance-test" {
		name       = "` + DeviceGroup1 + `"
		parent_id  = ome_static_group.parent.id
		device_ids = data.ome_device.devs.devices[*].id
	}
	`

	if FunctionMocker != nil {
		FunctionMocker.Release()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testGroupWoDesc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test", "name", DeviceGroup1),
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test", "description", ""),
				),
			},
			{
				Config: testUpdateParentID,
			},
			{
				Config: testAccCreateGroupSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test-1", "name", DeviceGroup1),
				),
			},
			{
				Config: testAccUpdateGroupSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test-1", "name", DeviceGroup1Update),
				),
			},
			{
				ResourceName:      "ome_static_group.terraform-acceptance-test-1",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       nil,
				ImportStateId:     DeviceGroup1Update,
			},
			{
				ResourceName:      "ome_static_group.terraform-acceptance-test-1",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("Error importing group"),
				ImportStateId:     "invalid",
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock((*clients.Client).CreateGroup).Return(0, fmt.Errorf("Duplicate name error")).Build()
				},
				// create group with existing group name
				Config:      testAccDuplicateNameNeg,
				ExpectError: regexp.MustCompile("Error while creation"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				// create group with invalid device ids
				Config:      testAccInvalidDeviceNeg,
				ExpectError: regexp.MustCompile("Error while adding group devices"),
			},
			{
				// this just sets up 2 groups for update test cases next
				Config: testAccCreate2,
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock((*clients.Client).UpdateGroup).Return(fmt.Errorf("Duplicate name error")).Build()
				},
				// update group name to existing group name
				Config:      testAccDuplicateNameNeg,
				ExpectError: regexp.MustCompile("Error while updation"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				// update group by adding invalid device id
				Config:      testAccInvalidDeviceNeg,
				ExpectError: regexp.MustCompile("Error while adding group devices"),
			},
			{
				// update Group to add 2 devices at once
				Config: testAccUpdateMultipleDevices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test-2", "device_ids.#", "2"),
				),
			},
			{
				// Recreate sub group with zero devices
				Taint:  []string{"ome_static_group.terraform-acceptance-test-2"},
				Config: testAccNoDevices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test-2", "device_ids.#", "0"),
				),
			},
			{
				// update Group to add 2 devices at once
				Config: testAccUpdateMultipleDevices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test-2", "device_ids.#", "2"),
				),
			},
			{
				// remove all devices from sub group
				Config: testAccNoDevices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_static_group.terraform-acceptance-test-2", "device_ids.#", "0"),
				),
			},
		},
	})
}
