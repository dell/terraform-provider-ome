package ome

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDeviceActionRes(t *testing.T) {
	testAccCreateDevicesResSuccess := testProvider + `
	data "ome_device" "devs" {
		filters = {
			device_service_tags = ["` + DeviceSvcTag1 + `", "` + DeviceSvcTag2 + `"]
		}
	}
	resource "ome_device_action" "code_1" {
		device_ids = data.ome_device.devs.devices[*].id
	}
	`
	// testAccCreateDevicesResIDSuccess := testProvider + `
	// resource "ome_device_action" "code_2" {
	// 	devices = [
	// 		{
	// 			id = ` + DeviceID1 + `
	// 		},
	// 		{
	// 			id = ` + DeviceID2 + `
	// 		}
	// 	]
	// }
	// `
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_device_action.code_1", "id"),
				),
			},
			// {
			// 	Config:        testAccCreateDevicesResSuccess,
			// 	ImportState:   true,
			// 	ResourceName:  "ome_device_action.code_1",
			// 	ImportStateId: fmt.Sprintf("id:%s,%s", DeviceID1, "not-integer"),
			// 	ExpectError:   regexp.MustCompile(".*ID could not be converted to an Int64.*"),
			// },
			// {
			// 	Config:        testAccCreateDevicesResSuccess,
			// 	ImportState:   true,
			// 	ResourceName:  "ome_device_action.code_1",
			// 	ImportStateId: fmt.Sprintf("ip:%s,%s", "ip1", "ip2"),
			// 	ExpectError:   regexp.MustCompile(".*Identifier of type ip is not recognised.*"),
			// },
			// {
			// 	Config:            testAccCreateDevicesResSuccess,
			// 	ImportState:       true,
			// 	ResourceName:      "ome_device_action.code_1",
			// 	ImportStateId:     fmt.Sprintf("svc_tag:%s,%s", DeviceSvcTag1, DeviceSvcTag2),
			// 	ImportStateVerify: true,
			// },
			// {
			// 	Config: testAccCreateDevicesResIDSuccess,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("ome_device_action.code_2", "devices.#", "2"),
			// 	),
			// },
			// {
			// 	Config:            testAccCreateDevicesResIDSuccess,
			// 	ImportState:       true,
			// 	ResourceName:      "ome_device_action.code_2",
			// 	ImportStateId:     fmt.Sprintf("id:%s,%s", DeviceID1, DeviceID2),
			// 	ImportStateVerify: true,
			// },
			// {
			// 	Config:            testAccCreateDevicesResIDSuccess,
			// 	ImportState:       true,
			// 	ResourceName:      "ome_device_action.code_2",
			// 	ImportStateId:     fmt.Sprintf("%s,%s", DeviceID1, DeviceID2),
			// 	ImportStateVerify: true,
			// },
		},
	})

}
