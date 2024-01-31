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
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestFwBaselineCompReportDatasource_Read is a test function for the Read method of fwBaselineCompReportDatasource.
//
// It creates a new instance of fwBaselineCompReportDatasource, creates a context, creates a mock datasource.ReadRequest,
// creates a mock datasource.ReadResponse, calls the Read method of fwBaselineCompReportDatasource, and performs
// assertions to check the expected behavior.
func TestDataSource_ReadFwCompReport(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testFwBaselineCompReportDS,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("all", "true")),
			},
		},
	})
}

// TestDataSource_ReadFwCompReportBaselineErr is a Go function that tests the behavior of the ReadFwCompReportBaselineErr function.
//
// This function takes a testing.T object as a parameter and performs a series of test steps using the resource.Test function.
// The test steps include setting up a pre-check function, providing ProtoV6ProviderFactories, and defining the configuration for the test.
// The test steps include a single step that configures the testFwBaselineCompReportDSInvalidBaseline configuration and expects an error message that matches the regular expression ".*Error fetching baseline*".
//
// The function does not return any values.
func TestDataSource_ReadFwCompReportBaselineErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testFwBaselineCompReportDSInvalidBaseline,
				ExpectError: regexp.MustCompile(".*Error fetching baseline*"),
			},
		},
	})
}

func TestDataSource_ReadFwCompReportInvAttrib(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testFwBaselineCompReportDSFilterKeyErr,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Match*"),
			},
		},
	})
}

// TestDataSource_ReadFwCompReportErrRet tests the read function for firmware compliance report error return.
//
// This function tests the read operation for the firewall compliance report in case of an error return. It mocks the
// GetFwBaselineComplianceReport function to return an error and checks if the error message matches the expected error.
//
// Parameters:
// - t: *testing.T - The testing object used for assertions and reporting test results.
//
// Return type: None.
func TestDataSource_ReadFwCompReportErrRet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error getting firmware reports
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFwBaselineComplianceReport, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      testFwBaselineCompReportDSErr,
				ExpectError: regexp.MustCompile(`.*Mock error*.`),
			},
		},
	})
}

var testFwBaselineCompReportDSInvalidBaseline = testProvider + `
	
	data "ome_firmware_baseline_compliance_report" "cr" {
		baseline_name = "tfacc_baseline_dell_invalid"
	}
`
var testFwBaselineCompReportDSErr = testProvider + `

	data "ome_firmware_baseline_compliance_report" "report3" {
		baseline_name = "tfacc_baseline_dell_1"

	}

	output "all" {
		value = length(data.ome_firmware_baseline_compliance_report.report3) > 0
	}
`
var testFwBaselineCompReportDSFilterKeyErr = testProvider + `

	data "ome_firmware_baseline_compliance_report" "report2" {
		
		filter {
			key = "DeviceModel-Error"
			value = "PowerEdge R640"
		}
		baseline_name = "tfacc_baseline_dell_1"
	}

	output "all" {
		value = length(data.ome_firmware_baseline_compliance_report.report2) > 0
	}
`

var testFwBaselineCompReportDS = testProvider + `

	data "ome_firmware_baseline_compliance_report" "report" {
		
		filter {
			key = "DeviceModel"
			value = "PowerEdge R640"
		}
		baseline_name = "tfacc_baseline_dell_1"
	}

	output "all" {
		value = length(data.ome_firmware_baseline_compliance_report.report) > 0
	}
`
