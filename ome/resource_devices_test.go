package ome

import (
	"fmt"
	"regexp"
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDevicesRes(t *testing.T) {
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
	testAccCreateDevicesResMixedSuccess := testProvider + `
	resource "ome_devices" "code_3" {
		devices = [
			{
				service_tag = "` + DeviceSvcTagRmv + `"
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
