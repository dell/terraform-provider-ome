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
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	User        = "test_create_user"
	UserUpdate  = "test_update_user"
	User1       = "test_create"
	UserUpdate1 = "test_update"
)

func TestUser(t *testing.T) {

	testAccProvider := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		port = "` + port + `"
 		protocol = "` + protocol + `"
		skipssl = true
	}
	`
	testAccCreateUserSuccess := testAccProvider + `	
	resource "ome_user" "code_1" {
		user_type_id =  1
		directory_service_id = 0
		description = "Avengers alpha"
		password = "Avenger1232$"
		username = "` + User + `"
		role_id = "10"
		locked = true
		enabled = true
	}
	`

	testAccUpdateGroupSuccess := testAccProvider + `	
	resource "ome_user" "code_1" {
		user_type_id =  1
		directory_service_id = 0
		description = "Avengers alpha"
		password = "Avenger1232$"
		username = "` + UserUpdate + `"
		role_id = "10"
		locked = true
		enabled = false
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateUserSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_user.code_1", "username", User),
				),
			},
			{
				Config: testAccUpdateGroupSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_user.code_1", "username", UserUpdate),
				),
			},
			{
				ResourceName:      "ome_user.code_1",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       nil,
				ImportStateIdFunc: testAccImportStateIDFunc("ome_user.code_1"),
			},
			{
				ResourceName:      "ome_user.code_1",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(clients.ErrGnrImportUser),
				ImportStateId:     "invalid",
			},
		},
	})
}

func TestUserNegative(t *testing.T) {

	if os.Getenv("TF_ACC") == "0" {
		t.Skip("Dont run with units tests because negative cases we are not running with mock server")
	}

	testAccProvider := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		port = "` + port + `"
 		protocol = "` + protocol + `"
		skipssl = true
	}
	`

	testAccCreateFailure := testAccProvider + `
	resource "ome_user" "code_2" {
		username = "123456789123456789"
		password = "Abcde123!"
		role_id = "101"
	}
	`

	testAccCreateUpdate := testAccProvider + `
	resource "ome_user" "code_3" {
		username = "` + User1 + `"
		password = "Abcde123!"
		role_id = "101"
	}
	`

	testAccUpdateFailure := testAccProvider + `
	resource "ome_user" "code_3" {
		username = "` + UserUpdate1 + `"
		password = "Abcde123!"
		role_id = "invalid"
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCreateFailure,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateUser),
			},
			{
				Config: testAccCreateUpdate,
			},
			{
				Config:      testAccUpdateFailure,
				ExpectError: regexp.MustCompile(clients.ErrGnrUpdateUser),
			},
		},
	})

}

func testAccImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["id"], rs.Primary.Attributes["password"]), nil
	}
}
