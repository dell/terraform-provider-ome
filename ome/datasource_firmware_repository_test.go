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
	"os"
	"regexp"
	"terraform-provider-ome/helper"
	"terraform-provider-ome/utils"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSource_Fbc_Repository_Read(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// No filter so should return all values
			{
				Config: allRepos,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("fbc-repository", "true"),
				),
			},
			// Empty Filter should show error
			{
				Config:      filterReposEmptyFilter,
				ExpectError: regexp.MustCompile(`.*Attribute names set must contain at least 1 elements.*`),
			},
			// Using the name filter
			{
				Config: filterRepos,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_firmware_repository.fbc-repository-name-filter", "fbc_repositories.#", "1"),
				),
			},
			// Error getting FBC Repositories
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetAllCatalogFirmware).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      allRepos,
				ExpectError: regexp.MustCompile(`.*Error Reading Repositories*`),
			},
			// Processing error test
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(utils.CopyFields).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      allRepos,
				ExpectError: regexp.MustCompile(`.*Error Copying values for repositories*`),
			},
		},
	})
}

var allRepos = testProvider + `
	data "ome_firmware_repository" "fbc-repository-all" {
	}
	output "fbc-repository" {
		value = length(data.ome_firmware_repository.fbc-repository-all.fbc_repositories) != 0
	}
`
var filterRepos = testProvider + `
	data "ome_firmware_repository" "fbc-repository-name-filter" {
			names = ["` + Repository + `"]
	}
	output "fbc-repository-name"{
		value = data.ome_firmware_repository.fbc-repository-name-filter
	}
`
var filterReposEmptyFilter = testProvider + `
	data "ome_firmware_repository" "fbc-repository-empty-filter" {
		names = []
	}
`
