package ome

import (
	"fmt"
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselinewithDeviceTag,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName)),
			},
			{ // Check if the device part of baseline
				Config:      testConfigureBaselineRemediationDevicePartOfBaseline,
				ExpectError: regexp.MustCompile(clients.ErrGnrBaseLineCreateRemediation),
			},
			{ // Check for the valid baseline with name
				Config:      testConfigureBaselineRemediationInvalidBaselineName,
				ExpectError: regexp.MustCompile(clients.ErrGnrBaseLineCreateRemediation),
			},
			{ //  Check for the valid baseline with id
				Config:      testConfigureBaselineRemediationInvalidBaselineID,
				ExpectError: regexp.MustCompile(clients.ErrGnrBaseLineCreateRemediation),
			},
			{ //  any one baseline name or id required
				Config:      testConfigureBaselineRemediationBaselineMutaully,
				ExpectError: regexp.MustCompile(clients.ErrGnrBaseLineCreateRemediation),
			},
			{ //  target device size
				Config:      testConfigureBaselineRemediationBaselneDevicesRequired,
				ExpectError: regexp.MustCompile(fmt.Sprintf(clients.ErrBaseLineTargetsSize, 1)),
			},
			{ //  invalid compliance status
				Config:      testConfigureBaselineRemediationCompianceStatus,
				ExpectError: regexp.MustCompile(fmt.Sprintf(clients.ErrBaseLineComplianceStatus, clients.ValidComplainceStatus)),
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

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		description = "baseline description"
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

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		description = "baseline description"
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
var testConfigureBaselineRemediationBaselineMutaully = `
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

var testConfigureBaselineRemediationCompianceStatus = `
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
				// ExpectNonEmptyPlan: true,
			},
			{
				Config: testConfigureBaselineRemediationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "target_devices.#", "2"),
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "target_devices.0.device_service_tag", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_compliance.baseline_remediation", "target_devices.1.device_service_tag", DeviceSvcTag2),
				),
				// ExpectNonEmptyPlan: true,
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
