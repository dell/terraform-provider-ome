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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	DeviceGroup1       = "test_acc_group_device_1"
	DeviceGroup1Update = "test_acc_group_device_1_updated"
)

func TestDeviceGroupCreation(t *testing.T) {

	testAccCreateGroupSuccess := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}
	
	resource "ome_group_device" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1 + `"
		description = "Device Group for Acceptance Test 1"
		parent_id = 1021
		device_ids = [` + DeviceID1 + `, ` + DeviceID2 + `]
	}
	`

	testAccUpdateGroupSuccess := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}
	
	resource "ome_group_device" "terraform-acceptance-test-1" {
		name = "` + DeviceGroup1Update + `"
		description = "Device Group for Acceptance Test 1 Updated"
		parent_id = 1021
		device_ids = [` + DeviceID1 + `]
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
				// ImportStateCheck: assertTFImportState,
				ExpectError:   nil,
				ImportStateId: DeviceGroup1Update,
			},
			// {
			// 	Config:      testAccUpdateTemplateWithExistingName,
			// 	ExpectError: regexp.MustCompile(clients.ErrUpdateTemplate),
			// },
			// {
			// 	Config:      testAccUpdateTemplateWithInvalidVlanNetworkID,
			// 	ExpectError: regexp.MustCompile(clients.ErrUpdateTemplate),
			// },
		},
	})
}
