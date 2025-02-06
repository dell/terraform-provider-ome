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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-ome/clients"
)

func TestAccDevicesRes(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	testAccCreateDevicesResSuccess := testProvider + `
	resource "ome_devices" "code_1" {
		devices = [
			{
				service_tag = "` + DeviceSvcTag1 + `"
			},
			{
				service_tag = "` + DeviceSvcTag2 + `"
			}
		]
	}
	`
	testAccCreateDevicesResIDSuccess := testProvider + `
	resource "ome_devices" "code_2" {
		devices = [
			{
				id = ` + DeviceID1 + `
			},
			{
				id = ` + DeviceID2 + `
			}
		]
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_devices.code_1", "devices.#", "2"),
				),
			},
			{
				Config:        testAccCreateDevicesResSuccess,
				ImportState:   true,
				ResourceName:  "ome_devices.code_1",
				ImportStateId: fmt.Sprintf("id:%s,%s", DeviceID1, "not-integer"),
				ExpectError:   regexp.MustCompile(".*ID could not be converted to an Int64.*"),
			},
			{
				Config:        testAccCreateDevicesResSuccess,
				ImportState:   true,
				ResourceName:  "ome_devices.code_1",
				ImportStateId: fmt.Sprintf("ip:%s,%s", "ip1", "ip2"),
				ExpectError:   regexp.MustCompile(".*Identifier of type ip is not recognised.*"),
			},
			{
				Config:            testAccCreateDevicesResSuccess,
				ImportState:       true,
				ResourceName:      "ome_devices.code_1",
				ImportStateId:     fmt.Sprintf("svc_tag:%s,%s", DeviceSvcTag1, DeviceSvcTag2),
				ImportStateVerify: true,
			},
			{
				Config: testAccCreateDevicesResIDSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_devices.code_2", "devices.#", "2"),
				),
			},
			{
				Config:            testAccCreateDevicesResIDSuccess,
				ImportState:       true,
				ResourceName:      "ome_devices.code_2",
				ImportStateId:     fmt.Sprintf("id:%s,%s", DeviceID1, DeviceID2),
				ImportStateVerify: true,
			},
			{
				Config:            testAccCreateDevicesResIDSuccess,
				ImportState:       true,
				ResourceName:      "ome_devices.code_2",
				ImportStateId:     fmt.Sprintf("%s,%s", DeviceID1, DeviceID2),
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccDevicesResUpdate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	if DeviceIPExt == "" {
		t.Skip("Skipping this test as there is no Device to remove.")
	}

	DiscoveryJob := `
	resource "ome_discovery" "discover1" {
		name = "discover-up-lab"
		schedule = "RunNow"
		ignore_partial_failure = true
		timeout = 10
		discovery_config_targets = [
			{
				network_address_detail = ["` + DeviceIPExt + `"]
				device_type = ["SERVER"]
				wsman = {
					username = "` + IdracUsername + `"
					password = "` + IdracPassword + `"
				}
			}
		]
	}
	`

	testAccCreateDevicesResMixedSuccess := testProvider + DiscoveryJob + `
	data "ome_device" "dev" {
		depends_on = [ome_discovery.discover1]
		filters = {
			ip_expressions = ["` + DeviceIPExt + `"]
		}
	}
	resource "ome_devices" "code_3" {
		devices = [
			{
				service_tag = data.ome_device.dev.devices[0].device_service_tag
			},
			{
				id = ` + DeviceID1 + `
			},
			{
				service_tag = "` + DeviceSvcTag2 + `"
			}
		]
	}
	`
	testAccUpdateDevicesResSuccess := testProvider + `
	resource "ome_devices" "code_3" {
		devices = [
			{
				id = ` + DeviceID1 + `
			},
			{
				service_tag = "` + DeviceSvcTag2 + `"
			}
		]
	}
	`
	testAccAddInvDevice := testProvider + `
	resource "ome_devices" "code_3" {
		devices = [
			{
				id = ` + DeviceID1 + `
			},
			{
				service_tag = "` + DeviceSvcTag2 + `"
			},
			{
				service_tag = "invalid"
			}
		]
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDevicesResMixedSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_devices.code_3", "devices.#", "3"),
				),
			},
			{
				// check that device can be removed from the list
				// device Svc Tag Rmv will be removed
				Config: testAccUpdateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_devices.code_3", "devices.#", "2"),
				),
			},
			{
				// check that invalid devices cannot be added to the list
				Config:      testAccAddInvDevice,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
		},
	})

}

func TestAccDevicesResUnk(t *testing.T) {
	testAccCreateDevicesResSuccess := testProvider + `
	resource "ome_devices" "code_1" {
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_devices.code_1", "devices.#"),
				),
			},
			{
				Config:       testAccCreateDevicesResSuccess,
				ImportState:  true,
				ResourceName: "ome_devices.code_1",
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return "", nil
				},
				ImportStateVerify: true,
			},
		},
	})

}
