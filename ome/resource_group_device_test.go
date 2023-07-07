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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	DeviceGroup1       = "test_acc_group_device_1"
	DeviceGroup1Update = "test_acc_group_device_1_updated"
)

func TestDeviceGroup(t *testing.T) {

	testAccProvider := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}
	`
	testAccCreateGroupSuccess := testAccProvider + `	
	resource "ome_group_device" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1"
		parent_id = 1021
		device_ids = [` + DeviceID1 + `, ` + DeviceID2 + `]
	}
	`

	testAccUpdateGroupSuccess := testAccProvider + `	
	resource "ome_group_device" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = [` + DeviceID1 + `]
	}
	`

	testAccDuplicateNameNeg := testAccProvider + `	
	resource "ome_group_device" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = [` + DeviceID1 + `]
	}
	resource "ome_group_device" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = [` + DeviceID2 + `]
	}
	`

	testAccInvalidDeviceNeg := testAccProvider + `	
	resource "ome_group_device" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = [` + DeviceID1 + `]
	}
	resource "ome_group_device" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = [-1]
	}
	`

	testAccCreate2 := testAccProvider + `	
	resource "ome_group_device" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = [` + DeviceID1 + `]
	}
	resource "ome_group_device" "terraform-acceptance-test-2" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = []
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateGroupSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_group_device.terraform-acceptance-test-1", "name", DeviceGroup1),
					// resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "refdevice_servicetag", DeviceSvcTag1),
					// resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "fqdds", "iDRAC,niC"),

					// resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "name", TemplateName2),
					// resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "refdevice_id", DeviceID2),
					// resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "fqdds", "All"),
					// resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "description", "This is sample description"),
				),
			},
			{
				Config: testAccUpdateGroupSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_group_device.terraform-acceptance-test-1", "name", DeviceGroup1Update),
				),
			},
			{
				ResourceName:      "ome_group_device.terraform-acceptance-test-1",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       nil,
				ImportStateId:     DeviceGroup1Update,
			},
			{
				ResourceName:      "ome_group_device.terraform-acceptance-test-1",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("Error importing group"),
				ImportStateId:     "invalid",
			},
			{
				// create group with existing group name
				Config:      testAccDuplicateNameNeg,
				ExpectError: regexp.MustCompile("Error while creation"),
			},
			{
				// create group with invalid device ids
				Config:      testAccInvalidDeviceNeg,
				ExpectError: regexp.MustCompile("Error while adding group devices"),
			},
			{
				// this just sets up 2 groups for update test cases next
				Config: testAccCreate2,
			},
			{
				// update group name to existing group name
				Config:      testAccDuplicateNameNeg,
				ExpectError: regexp.MustCompile("Error while updation"),
			},
			{
				// update group by adding invalid device id
				Config:      testAccInvalidDeviceNeg,
				ExpectError: regexp.MustCompile("Error while adding group devices"),
			},
		},
	})
}
