package ome

import (
	"os"
	"regexp"
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSource_ReadConfigurationReport(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfiguratiponReportDSInvalid,
				ExpectError: regexp.MustCompile(clients.ErrGnrConfigurationReport),
			},
			{
				Config: testConfiguratiponReportDS,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_configuration_report_info.cr", "compliance_report_device.#", "2")),
			},
		},
	})
}

var testConfiguratiponReportDSInvalid = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_configuration_report_info" "cr" {
		id = "0"
		baseline_name = "` + "InvalidBaseline" + `"
		fetch_attributes = true
	}
`

var testConfiguratiponReportDS = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_template" "terraform-acceptance-test-1" {
		name = "` + TestRefTemplateName + `"
		refdevice_servicetag = "` + DeviceSvcTag1 + `"
		fqdds = "EventFilters"
		view_type = "Compliance"
		job_retry_count = 20
		sleep_interval = 30
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `","` + DeviceSvcTag2 + `"]
		description = "baseline description"
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}

	data "ome_configuration_report_info" "cr" {
		id = "0"
		baseline_name = "` + BaselineName + `"
		fetch_attributes = true
		depends_on =[
			"ome_configuration_baseline.create_baseline"
		]
	}
`
