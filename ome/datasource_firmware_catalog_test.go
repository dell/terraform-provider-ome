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
	"terraform-provider-ome/helper"
	"terraform-provider-ome/utils"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSource_FirmwareCatalogRead(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// No filter so should return all values
			{
				Config: firmwareCatAll + outputs,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_many", "true"),
					resource.TestCheckOutput("fetched_one", "false"),
				),
			},
			// Empty Filter should return error
			{
				Config:      firmwareCatAllEmptyNameFilter + outputs,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value*.`),
			},
			// Using the name filter
			{
				Config: firmwareCatFilter + outputs,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched_any", "true"),
					resource.TestCheckOutput("fetched_many", "false"),
					resource.TestCheckOutput("fetched_one", "true"),
				),
			},
			// Error getting firmware catalogs
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetAllCatalogFirmware).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      firmwareCatAll,
				ExpectError: regexp.MustCompile(`.*Error fetching firmware catalogs*.`),
			},
			// Processing error test
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(utils.CopyFields).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      firmwareCatAll,
				ExpectError: regexp.MustCompile(`.*Error processing firmware catalogs*.`),
			},
		},
	})
}

var outputs = `
output "fetched_many" {
	value = length(data.ome_firmware_catalog.data-catalog.firmware_catalogs) > 1
}
  
output "fetched_any" {
	value = length(data.ome_firmware_catalog.data-catalog.firmware_catalogs) != 0
}

output "fetched_one" {
	value = length(data.ome_firmware_catalog.data-catalog.firmware_catalogs) == 1
}
`

var firmwareCatAll = testProvider + `
data "ome_firmware_catalog" "data-catalog" {
}
`

var firmwareCatAllEmptyNameFilter = testProvider + `
data "ome_firmware_catalog" "data-catalog" {
	names = []
}  
`

var firmwareCatFilter = testProvider + `

data "ome_firmware_catalog" "data-catalog" {
	names = ["` + Catalog1 + `"]
}
`
