/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	InvalidBaselineID = "1000910"
)

func TestConfigurationRemediationErrors(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	temps := initTemplates(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselinewithDeviceTag + temps.templateSvcTag1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName)),
			},
			{ // Check if the device part of baseline
				Config:      testConfigureBaselineRemediationDevicePartOfBaseline + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(clients.ErrGnrBaseLineCreateRemediation),
			},
		},
	})

}
func TestConfigurationRemediationInvScenarios(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Check for the valid baseline with name
				Config:      testConfigureBaselineRemediationInvalidBaselineName,
				ExpectError: regexp.MustCompile(".*baseline not found.*"),
			},
			{ //  Check for the valid baseline with id
				Config:      testConfigureBaselineRemediationInvalidBaselineID,
				ExpectError: regexp.MustCompile(".*baseline not found.*"),
			},
			{ //  any one baseline name or id required
				Config:      testConfigureBaselineRemediationBaselineMutually,
				ExpectError: regexp.MustCompile(".*either baseline name or id is required.*"),
			},
			{ //  target device size
				Config:      testConfigureBaselineRemediationBaselneDevicesRequired,
				ExpectError: regexp.MustCompile(".*Attribute target_devices set must contain at least 1 elements.*"),
			},
			{ //  invalid compliance status
				Config:      testConfigureBaselineRemediationComplianceStatus,
				ExpectError: regexp.MustCompile(".*Error: Invalid Attribute Value Match.*"),
			},
			{ //  if both baseline name or id not specfied
				Config:      testConfigureBaselineRemediationBaselineInfo,
				ExpectError: regexp.MustCompile(clients.ErrGnrBaseLineCreateRemediation),
			},
		},
	})

}

var testConfigureBaselineRemediationDevicePartOfBaseline = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		description = "baseline description"
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_name = "` + BaselineName + `"
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag2 + `"
				compliance_status = "Compliant"
			},
		  ]
	}
`

var testConfigureBaselineRemediationInvalidBaselineName = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_name = "` + "InValidBaselineName" + `"
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag1 + `"
				compliance_status = "Compliant"
			},
		  ]
	}
`
var testConfigureBaselineRemediationInvalidBaselineID = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_id = ` + InvalidBaselineID + `
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag1 + `"
				compliance_status = "Compliant"
			},
		  ]
	}
`
var testConfigureBaselineRemediationBaselineMutually = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		description = "baseline description"
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_id = ` + InvalidBaselineID + `
		baseline_name = "` + BaselineName + `"
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag1 + `"
				compliance_status = "Compliant"
			},
		  ]
	}
`

var testConfigureBaselineRemediationBaselineInfo = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		description = "baseline description"
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag1 + `"
				compliance_status = "Compliant"
			},
		  ]
	}
`

var testConfigureBaselineRemediationBaselneDevicesRequired = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_name = "` + BaselineName + `"
		target_devices = [
		  ]
	}
`

var testConfigureBaselineRemediationComplianceStatus = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `", "` + DeviceSvcTag2 + `"]
		description = "baseline description"
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_name = "` + BaselineName + `"
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag2 + `"
				compliance_status = "NonCompliant"
			},
		  ]
	}
`

func TestConfigurationRemediation(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselineRemediation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "target_devices.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "target_devices.0.device_service_tag", DeviceSvcTag1),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testConfigureBaselineRemediationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "target_devices.#", "2"),
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "target_devices.0.device_service_tag", DeviceSvcTag2),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

var testConfigureBaselineRemediation = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "avenger-temp"
		device_servicetags = ["` + DeviceSvcTag1 + `","` + DeviceSvcTag2 + `"]
		description = "baseline description"
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_name = ome_configuration_baseline.create_baseline.baseline_name
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag1 + `"
				compliance_status = "Compliant"
			},
		  ]
	}
`

var testConfigureBaselineRemediationUpdate = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "avenger-temp"
		device_servicetags = ["` + DeviceSvcTag1 + `","` + DeviceSvcTag2 + `"]
		description = "baseline description"
	}

	resource "ome_configuration_compliance" "baseline_remediation" {
		baseline_name = ome_configuration_baseline.create_baseline.baseline_name
		target_devices = [
			{
				device_service_tag = "` + DeviceSvcTag1 + `"
				compliance_status = "Compliant"
			},
			{
				device_service_tag = "` + DeviceSvcTag2 + `"
				compliance_status = "Compliant"
			},
		  ]
	}
`
