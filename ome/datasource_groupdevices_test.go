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
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_ids.#", "2"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_servicetags.#", "2")),
			},
			{
				Config: testgroupDeviceDSInvalidGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_servicetags.#", "0")),
			},
			{
				Config: testgroupDeviceDSEmptyGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("data.ome_groupdevices_info.gd", "device_servicetags.#", "0")),
			},
		},
	})
}

var testgroupDeviceDS = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_groupdevices_info" "gd" {
		id = "0"
		device_group_names = ["test_device_group"]
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
		device_group_names = ["NO_GROUP"]
	}
`
