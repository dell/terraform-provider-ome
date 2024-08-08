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
	"terraform-provider-ome/helper"
	"terraform-provider-ome/utils"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSource_DeviceComplianceReport(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	//Test to retrive the device compliance report
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Should fetch at least 1 device compliance report using device id
			{
				Config: deviceComplianceReportDeviceID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetch", "true"),
				),
			},
			// Should fetch at least 1 device compliance report using service tag
			{
				Config: deviceComplianceReportServiceTag,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetch", "true"),
				),
			},
			// Should fetch at least 1 device compliance report using group name
			{
				Config: deviceComplianceReportGroupName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetch", "true"),
				),
			},
			// Error getting device reports
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetAllDeviceComplianceReport).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      deviceComplianceReportDeviceID,
				ExpectError: regexp.MustCompile(`.*Error reading device compliance report*.`),
			},
			// Error getting device reports
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(utils.CopyFields).Return(fmt.Errorf("Mock error")).Build()
				},
				Config:      deviceComplianceReportDeviceID,
				ExpectError: regexp.MustCompile(`.*Error processing device compliance report*.`),
			},
		},
	})
}

var deviceComplianceReportDeviceID = testProvider + `
data "ome_device_compliance_report" "device_compliance_report_data" {
	device_ids = [` + DeviceID1 + `]
}
output "fetch" {
	value = length(data.ome_device_compliance_report.device_compliance_report_data.device_compliance_reports) > 0
}
`

var deviceComplianceReportServiceTag = testProvider + `
data "ome_device_compliance_report" "device_compliance_report_data" {
	device_service_tags = ["` + DeviceSvcTag1 + `"]
}
output "fetch" {
	value = length(data.ome_device_compliance_report.device_compliance_report_data.device_compliance_reports) > 0
}
`

var deviceComplianceReportGroupName = testProvider + `
data "ome_device_compliance_report" "device_compliance_report_data" {
	device_group_names = ["Servers"]
}
output "fetch" {
	value = length(data.ome_device_compliance_report.device_compliance_report_data.device_compliance_reports) > 0
}
`
