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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDiscoveryOne(t *testing.T) {
	testAccCreateDiscoverySuccess := testProvider + `
	resource "ome_discovery" "code_1" {
		name = "test-create"
		schedule = "RunLater"
		cron = "0 * */10 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = [ "` + DeviceIP1 + `", "` + DeviceIP2 + `"]
		  redfish = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		}]
	  }
	`
	testAccUpdateDiscoverySuccess := testProvider + `
	resource "ome_discovery" "code_1" {
		name = "test-update"
		schedule = "RunLater"
		cron = "0 * */12 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = [ "` + DeviceIP1 + `"]
		  redfish = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		}]
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDiscoverySuccess,
			},
			{
				Config: testAccUpdateDiscoverySuccess,
			},
		},
	})

}

func TestAccDiscoveryTwo(t *testing.T) {
	testAccCreateDiscoveryDebug := testProvider + `
	resource "ome_discovery" "code_1" {
		name = "shiva-ganga"
		schedule = "RunLater"
		cron = "0 * */10 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = [ "` + DeviceIP1 + `","` + DeviceIP2 + `","` + DeviceIPExt + `"]
		  # redfish = {
		  #  username = "root"
		  #  password = "calvin" 
		  # }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		}]
	  }
	`
	testAccUpdateDiscoveryDebug := testProvider + `
	resource "ome_discovery" "code_1" {
		name = "kashi-ganga"
		schedule = "RunLater"
		cron = "0 * */12 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = ["` + DeviceIPExt + `"]
		  redfish = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		}]
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDiscoveryDebug,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_discovery.code_1", "name", "shiva-ganga"),
					resource.TestCheckResourceAttr("ome_discovery.code_1", "cron", "0 * */10 * * ? *"),
				),
			},
			{
				Config: testAccUpdateDiscoveryDebug,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_discovery.code_1", "name", "kashi-ganga"),
					resource.TestCheckResourceAttr("ome_discovery.code_1", "cron", "0 * */12 * * ? *"),
				),
			},
		},
	})

}

func TestAccDiscoveryThree(t *testing.T) {

	invalidDiscoveryConfigOne := testProvider + `
	resource "ome_discovery" "code_3" {
		name = "invalid-config"
		schedule = "RunNow"
		cron = "abc"
		discovery_config_targets = [ {
		  device_type = ["abc"]
		  network_address_detail = ["x.x.x.x"]
		  } 
		]
	  }
	`

	invalidDiscoveryConfigtwo := testProvider + `
	resource "ome_discovery" "code_4" {
		name = "invalid-config"
		schedule = "RunNow"
		timeout = 10
		ignore_partial_failure = true
		discovery_config_targets = []
	  }
	`

	invalidDiscoveryConfigThree := testProvider + `
	resource "ome_discovery" "code_4" {
		name = "invalid-config"
		schedule = "RunNow"
		timeout = 10
		ignore_partial_failure = true
		discovery_config_targets = [{
		  device_type = []
		  network_address_detail = []
		}]
	  }
	`

	invalidDiscoveryConfigFour := testProvider + `
	resource "ome_discovery" "code_4" {
		name     = "invalid-config"
		schedule = "RunLater"
		timeout = 10
		ignore_partial_failure = true
		discovery_config_targets = [{
		  device_type            = ["SERVER"]
		  network_address_detail = ["9.0.0.1"]
		  wsman = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		}]
	  }
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      invalidDiscoveryConfigOne,
				ExpectError: regexp.MustCompile(`.*value must be one of:*.`),
			},
			{
				Config:      invalidDiscoveryConfigOne,
				ExpectError: regexp.MustCompile(`.*With Schedule as RunNow, CRON can't be set*.`),
			},
			{
				Config:      invalidDiscoveryConfigOne,
				ExpectError: regexp.MustCompile(`.*Atleast one of protocol should be configured for the discovery targets*.`),
			},
			{
				Config:      invalidDiscoveryConfigtwo,
				ExpectError: regexp.MustCompile(`.*Attribute discovery_config_targets set must contain at least 1 elements*.`),
			},
			{
				Config:      invalidDiscoveryConfigThree,
				ExpectError: regexp.MustCompile(`.*list must contain at least 1 elements*.`),
			},
			{
				Config:      invalidDiscoveryConfigFour,
				ExpectError: regexp.MustCompile(`.*With Schedule as RunLater, Partial Failure can't be set*.`),
			},
			{
				Config:      invalidDiscoveryConfigFour,
				ExpectError: regexp.MustCompile(`.*With Schedule as RunLater, cron must be set*.`),
			},
			{
				Config:      invalidDiscoveryConfigFour,
				ExpectError: regexp.MustCompile(`.*With Schedule as RunLater, Timeout can't be set*.`),
			},
		},
	})
}

func TestAccDiscoveryFour(t *testing.T) {
	TrackDiscoveryJob := testProvider + `
	resource "ome_discovery" "discover1" {
		name = "discover-lab"
		schedule = "RunNow"
		timeout = 10
		ignore_partial_failure = true
		discovery_config_targets = [
		  {
		  network_address_detail = ["` + DeviceIP1 + `","` + DeviceIP2 + `","` + DeviceIPExt + `", "127.0.0.1","0.42.42.42","1.1.1.1","8.8.8.8","192.168.1.1"]
		  device_type = ["SERVER"]
		  wsman = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		}]
	  }
	`
	TrackDiscoveryJobUpdate := testProvider + `
	resource "ome_discovery" "discover1" {
		name = "discover-up-lab"
		schedule = "RunNow"
		timeout = 5
		ignore_partial_failure = true
		discovery_config_targets = [
		  {
		  network_address_detail = ["` + DeviceIP1 + `"]
		  device_type = ["SERVER"]
		  wsman = {
			username = "` + IdracUsername + `"
			password = "` + IdracPassword + `"
		  }
		}]
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TrackDiscoveryJob,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_discovery.discover1", "name", "discover-lab"),
				),
			},
			{
				Config: TrackDiscoveryJobUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_discovery.discover1", "name", "discover-up-lab"),
				),
			},
		},
	},
	)
}
