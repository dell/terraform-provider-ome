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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	User       = "test_acc_user_1"
	UserUpdate = "test_acc_user_2"
)

func TestUser(t *testing.T) {

	var userID string

	testAccProvider := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
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
					func(s *terraform.State) error {
						userID = s.RootModule().Resources["ome_user.code_1"].Primary.Attributes["id"]
						return nil
					},
				),
			},
			{
				ResourceName:      "ome_user.code_1",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       nil,
				ImportStateId:     userID,
			},
			// {
			// 	// create group with existing group name
			// 	Config:      testAccDuplicateNameNeg,
			// 	ExpectError: regexp.MustCompile("Error while creation"),
			// },
		},
	})
}
