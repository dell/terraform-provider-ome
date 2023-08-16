package ome

import (
	// "regexp"
	"regexp"
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	testAccUpdateDevicesResSuccess := testProvider + `
	resource "ome_devices" "code_1" {
		devices = [
			{
				service_tag = "` + DeviceSvcTag2 + `"
			}
		]
	}
	`
	testAccAddInvDevice := testProvider + `
	resource "ome_devices" "code_1" {
		devices = [
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
				Config: testAccCreateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_devices.code_1", "devices.#", "2"),
				),
			},
			{
				// check if device can be removed from the list
				Config: testAccUpdateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_devices.code_1", "devices.#", "1"),
				),
			},
			{
				// check if device can be removed from the list
				Config:      testAccAddInvDevice,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
			{
				// check if device can be added to the list
				Config: testAccCreateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_devices.code_1", "devices.#", "2"),
				),
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
		},
	})

}
