package ome

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceActionResInvalid(t *testing.T) {
	test := testProvider + `
	resource "ome_device_action" "code_1" {
		job_name = "refresh-job"
		device_ids = [100, 200]
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      test,
				ExpectError: regexp.MustCompile(".*Error creating job.*"),
			},
		},
	})
}

func TestAccDeviceActionRes(t *testing.T) {
	getDeviceIds := `
	data "ome_device" "devs" {
		filters = {
			device_service_tags = ["` + DeviceSvcTag1 + `", "` + DeviceSvcTag2 + `"]
		}
	}
	`
	testAccCreateDevicesResSuccess := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		device_ids = data.ome_device.devs.devices[*].id
		job_name = "refresh-job"
		job_description = "r-job-desc"
		timeout = 5
	}
	`

	testAccUpdateDevicesRes := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		device_ids = data.ome_device.devs.devices[*].id
		job_name = "refresh-job"
		job_description = "r-job-desc"
		timeout = 10
	}
	`

	testAccCreateDevicesResCron := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		device_ids = data.ome_device.devs.devices[*].id
		job_name = "refresh-job"
		job_description = "r-job-desc"
		cron = "0 * */10 * * ? *"
	}
	`
	testAccNoDesc := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		device_ids = data.ome_device.devs.devices[*].id
		job_name = "refresh-job"
		cron = "0 * */10 * * ? *"
	}
	`
	testAccNoCronOrTimeout := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		device_ids = data.ome_device.devs.devices[*].id
		job_name = "refresh-job"
	}
	`
	testAccNoJobNameNeg := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		device_ids = data.ome_device.devs.devices[*].id
	}
	`
	testAccNoDevicesNeg := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		job_name = "refresh-job"
	}
	`
	testAccZeroDevicesNeg := testProvider + getDeviceIds + `
	resource "ome_device_action" "code_1" {
		job_name = "refresh-job"
		device_ids = []
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccNoJobNameNeg,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(".*The argument \"job_name\" is required.*"),
			},
			{
				Config:      testAccNoDevicesNeg,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(".*The argument \"device_ids\" is required.*"),
			},
			{
				Config:      testAccZeroDevicesNeg,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(".*must contain at least 1 elements.*"),
			},
			{
				Config: testAccCreateDevicesResCron,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_device_action.code_1", "id"),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "action", "inventory_refresh"),
					resource.TestCheckNoResourceAttr("ome_device_action.code_1", "timeout"),
				),
			},
			{
				Taint:  []string{"ome_device_action.code_1"},
				Config: testAccNoDesc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_device_action.code_1", "id"),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "job_description", ""),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "action", "inventory_refresh"),
					resource.TestCheckNoResourceAttr("ome_device_action.code_1", "timeout"),
				),
			},
			{
				Config: testAccNoCronOrTimeout,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_device_action.code_1", "id"),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "action", "inventory_refresh"),
					resource.TestCheckNoResourceAttr("ome_device_action.code_1", "timeout"),
					resource.TestCheckNoResourceAttr("ome_device_action.code_1", "cron"),
				),
			},
			{
				Taint:  []string{"ome_device_action.code_1"},
				Config: testAccCreateDevicesResSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_device_action.code_1", "id"),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "action", "inventory_refresh"),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "timeout", "5"),
					resource.TestCheckNoResourceAttr("ome_device_action.code_1", "cron"),
				),
			},
			{
				Config: testAccUpdateDevicesRes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_device_action.code_1", "id"),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "action", "inventory_refresh"),
					resource.TestCheckResourceAttr("ome_device_action.code_1", "timeout", "10"),
					resource.TestCheckNoResourceAttr("ome_device_action.code_1", "cron"),
				),
			},
		},
	})

}
