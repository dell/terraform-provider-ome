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
	"os"
	"regexp"
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSource_ReadConfigurationReport(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	temps := initTemplates(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfiguratiponReportDSInvalidCreds,
				ExpectError: regexp.MustCompile(".*invalid credentials.*"),
			},
			{
				Config:      testConfiguratiponReportDSInvalid,
				ExpectError: regexp.MustCompile(clients.ErrGnrConfigurationReport),
			},
			{
				Config: testConfiguratiponReportDS + temps.templateSvcTag1Full,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_configuration_report_info.cr", "compliance_report_device.#", "2")),
			},
		},
	})
}

var testConfiguratiponReportDSInvalidCreds = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "invalid"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_configuration_report_info" "cr" {
		id = "0"
		baseline_name = "` + "InvalidBaseline" + `"
		fetch_attributes = true
	}
`

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
