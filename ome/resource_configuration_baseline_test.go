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
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	BaselineName              = "test_acc_create_baseline"
	BaselineNameUpdate        = "test_acc_update_baseline"
	TestRefTemplateName       = "test_acc_compliance_template"
	TestRefTemplateNameUpdate = "test_acc_compliance_template_update"
)

func init() {
	resource.AddTestSweepers("ome_configuration_baseline", &resource.Sweeper{
		Name: "ome_configuration_baseline",
		F: func(region string) error {
			fmt.Println("Sweepers for baseline invoked")
			omeClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client")
				return nil
			}

			_, err = omeClient.CreateSession()
			if err != nil {
				log.Println("Error creating client session for sweeper")
				return nil
			}
			defer omeClient.RemoveSession()

			omeBaselines := []models.OmeBaseline{}
			err = omeClient.GetPaginatedData(clients.BaselineAPI, &omeBaselines)
			if err != nil {
				log.Println("failed to fetch baseline details for the name " + SweepTestsTemplateIdentifier)
				return nil
			}

			var baselineIDs []int64
			for _, omeBaseline := range omeBaselines {
				if strings.Contains(omeBaseline.Name, SweepTestsTemplateIdentifier) {
					baselineIDs = append(baselineIDs, omeBaseline.ID)
				}
			}

			err = omeClient.DeleteBaseline(baselineIDs)
			if err != nil {
				log.Println("failed to sweep dangling baselines.")
				return nil
			}
			return nil
		},
	})
}

func TestCreateBaseline_TestValidations(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Step 1
				Config:      testCreateBaselineFailureWithBothTemplateIDName,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 2
				Config:      testCreateBaselineValidationFailureEmptyTemplateIDName,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 3
				Config:      testCreateBaselineValidationFailureInvalidTemplateName,
				ExpectError: regexp.MustCompile("reference template id or name should be of type compliance"),
			},
			{ // Step 4
				Config:      testCreateBaselineValidationFailureNonComplianceTemplateID,
				ExpectError: regexp.MustCompile("reference template id or name should be of type compliance"),
			},
		},
	})
}

func TestCreateBaseline_TestNotificationValidations(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: justProvider + `
				resource "ome_configuration_baseline" "create_baseline_validation_failure_unexpected_cron" {
					ref_template_name = "not-needed"
					baseline_name = "test_acc_create_baseline"
					device_servicetags = ["10129"]
					schedule = true
					notify_on_schedule = false
					cron = "abc"
				}
				`,
				ExpectError: regexp.MustCompile(".*attributes `cron` is not accepted*"),
				// ExpectError: regexp.MustCompile(clients.ErrBaseLineNotifyValid),
			},
			{
				Config: justProvider + `
				resource "ome_configuration_baseline" "create_baseline_validation_failure_unexpected_cron" {
					ref_template_name = "not-needed"
					baseline_name = "test_acc_create_baseline"
					device_servicetags = ["10129"]
					cron = "abc"
				}
				`,
				ExpectError: regexp.MustCompile(".*attributes `cron` and `email_addresses` are accepted only when `schedule` is.*"),
				// ExpectError: regexp.MustCompile(".*" + clients.ErrCronRequired + ".*"),
			},
			{
				Config: justProvider + `
				resource "ome_configuration_baseline" "create_baseline_validation_failure_unexpected_cron" {
					ref_template_name = "not-needed"
					baseline_name = "test_acc_create_baseline"
					device_servicetags = ["10129"]
					email_addresses=["abc@gmail.com"]
				}
				`,
				ExpectError: regexp.MustCompile("attributes `cron` and `email_addresses` are accepted only when `schedule` is.*"),
				// ExpectError: regexp.MustCompile(".*" + clients.ErrCronRequired + ".*"),
			},
		},
	})
}

func TestCreateBaseline_TestValidationsWithValidTemplate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	temps := initTemplates(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: justProvider + temps.templateSvcTag1,
			},
			{ // Step 2
				Config:      testCreateBaselineValidationFailureEmptyDevice + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(clients.ErrDeviceRequired),
			},
			{ // Step 3
				Config:      testCreateBaselineValidationFailureInvalidDevice + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile("invalid service tags:"),
			},
			{ // Step 4
				Config: testCreateBaselineValidationFailureScheduleNotificationEmptyEmail + temps.templateSvcTag1,
				// ExpectError: regexp.MustCompile(clients.ErrScheduleNotification),
				ExpectError: regexp.MustCompile(".*please provide a valid email address.*"),
			},
			{ // Step 5
				Config:      testCreateBaselineValidationFailureScheduleNotificationInvalidEmail + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(fmt.Sprintf(clients.ErrInvalidEmailAddress, "abc")),
			},
			{ // Step 6
				Config:      testCreateBaselineValidationFailureNotificationOnScheduleEmptyCron + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(clients.ErrInvalidCronExpression),
			},
			{ // Step 7
				Config:      testCreateBaselineValidationInvalidOutputFormatCase + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			{ // Step 8
				Config:      testCreateBaselineValidationInvalidOutputFormat + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			{ // Step 9
				Config:      testCreateBaselineValidationDeviceCapable + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
		},
	})
}

var justProvider = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}
`

var testCreateBaselineFailureWithBothTemplateIDName = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_failure_template_id_name" {
		ref_template_name = "` + TestRefTemplateName + `"
		ref_template_id = ` + "123" + `
		baseline_name = "` + BaselineName + `"
	}
`

var testCreateBaselineValidationFailureEmptyTemplateIDName = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_failure_empty_template_id_name" {
		baseline_name = "` + BaselineName + `"
	}
`

var testCreateBaselineValidationFailureInvalidTemplateName = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_failure_empty_template_id_name" {
		ref_template_name = "invalid"
		baseline_name = "` + BaselineName + `"
	}
`

var testCreateBaselineValidationFailureNonComplianceTemplateID = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_failure_non_compliance_template_id" {
		ref_template_name = "` + TestRefTemplateName + `"
		baseline_name = "` + BaselineName + `"
	}
`

var testCreateBaselineValidationFailureEmptyDevice = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_failure_empty_device" {
		ref_template_name = "` + TestRefTemplateName + `"
		baseline_name = "` + BaselineName + `"
	}
`

var testCreateBaselineValidationFailureInvalidDevice = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_invalid_device" {
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = [` + "10001" + `]
		baseline_name = "` + BaselineName + `"
	}
`

var testCreateBaselineValidationFailureNotificationOnScheduleEmptyCron = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_empty_cron" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		schedule=true
		email_addresses=["abc@gmail.com"]
		notify_on_schedule=true
	}
`
var testCreateBaselineValidationFailureScheduleNotificationInvalidEmail = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_empty_email" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		schedule=true
		email_addresses= ["abc"]
	}

`

var testCreateBaselineValidationFailureScheduleNotificationEmptyEmail = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_empty_email" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		schedule=true
	}
`

var testCreateBaselineValidationInvalidOutputFormatCase = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_empty_email" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		output_format = "CSV" // cannot be uppercase
	}
`

var testCreateBaselineValidationInvalidOutputFormat = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_empty_email" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		output_format = "abc" // invalid value be uppercase
	}
`

var testCreateBaselineValidationDeviceCapable = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline_validation_empty_email" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = [` + DeviceID2 + `]
	}
`

func TestCreateBaseline_BaselineWithDeviceIDAndTags(t *testing.T) {
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
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0")),
			},
			{
				Config: testConfigureBaselinewithDeviceTagUpdate + temps.templateSvcTag1 + temps.templateSvcTag2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineNameUpdate),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description updated"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateNameUpdate),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag2),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0")),
			},
			{
				Config: testConfigureBaselinewithDeviceID + temps.templateSvcTag1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName+"-1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.0", DeviceID1)),
			},
			{
				Config: testConfigureBaselinewithDeviceIDUpdate + temps.templateSvcTag1 + temps.templateSvcTag2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineNameUpdate+"-1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description updated"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateNameUpdate),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.0", DeviceID2)),
			},
			{
				Config:        testImportConfigurationBaseline,
				ResourceName:  "ome_configuration_baseline.import_baseline",
				ImportState:   true,
				ExpectError:   nil,
				ImportStateId: BaselineNameUpdate + "-1",
			},
		},
	})
}

var testConfigureBaselinewithDeviceTag = `
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
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

var testConfigureBaselinewithDeviceTagUpdate = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineNameUpdate + `"
		ref_template_name = "` + TestRefTemplateNameUpdate + `"
		device_servicetags = ["` + DeviceSvcTag2 + `"]
		description = "baseline description updated"
		depends_on = ["ome_template.terraform-acceptance-test-2"]
	}
`

var testConfigureBaselinewithDeviceID = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `-1"
		ref_template_name = "` + TestRefTemplateName + `"
		device_ids = ["` + DeviceID1 + `"]
		description = "baseline description"
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

var testConfigureBaselinewithDeviceIDUpdate = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineNameUpdate + `-1"
		ref_template_name = "` + TestRefTemplateNameUpdate + `"
		device_ids = ["` + DeviceID2 + `"]
		description = "baseline description updated"
		depends_on = ["ome_template.terraform-acceptance-test-2"]
	}
`

func TestCreateBaseline_CreateBaselineWithSchedule(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	temps := initTemplates(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselineWithSchedule + temps.templateSvcTag1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.0", "test@testmail.com"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "cron", "0 49 8 * * ? *"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "notify_on_schedule", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "output_format", "csv"),
				),
			},
			{
				Config: testConfigureBaselineWithScheduleUpdate + temps.templateSvcTag1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.0", "test@mail.com"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "cron", "0 50 8 * * ? *"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "notify_on_schedule", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "output_format", "html"),
				),
			},
			{
				Config:        testImportConfigurationBaseline,
				ResourceName:  "ome_configuration_baseline.import_baseline",
				ImportState:   true,
				ExpectError:   nil,
				ImportStateId: BaselineName,
			},
		},
	})
}

var testConfigureBaselineWithSchedule = `
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
		schedule = true
		notify_on_schedule = true
		email_addresses = ["test@testmail.com"]
		cron = "0 49 8 * * ? *"
		output_format = "csv"
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`
var testConfigureBaselineWithScheduleUpdate = `
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
		schedule = true
		notify_on_schedule = true
		email_addresses = ["test@mail.com"]
		cron = "0 50 8 * * ? *"
		output_format = "html"
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

func TestCreateBaseline_CreateBaselineWithScheduleNonCompliant(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	temps := initTemplates(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselineScheduleNonCompliant + temps.templateSvcTag1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.0", "test@testmail.com"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "output_format", "html"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "notify_on_schedule", "false")),
			},
			{
				Config: testConfigureBaselineScheduleNonCompliantUpdate + temps.templateSvcTag1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule", "false"),
				),
			},
		},
	})
}

var testConfigureBaselineScheduleNonCompliant = `
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
		schedule = true
		email_addresses = ["test@testmail.com"]
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

var testConfigureBaselineScheduleNonCompliantUpdate = `
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
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

var testImportConfigurationBaseline = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "import_baseline" {
	}
`

func TestCreateBaseline_Update(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	temps := initTemplates(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testCreateBaselineWrongCreds,
				ExpectError: regexp.MustCompile(".*invalid credentials.*"),
			},
			{
				Config: testConfigureBaselinewithDeviceTag + temps.templateSvcTag1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0")),
			},
			{
				Config:      testUpdateBaselinewithInvalidDeviceSvcTag + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(".*invalid service tags.*"),
			},
			{
				Config:      testUpdateBaselinewithInvalidTemplate + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(".*reference template id or name should be of type compliance.*"),
			},
			{
				Config:      testUpdateBaselinewithUnexpectedCron + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(".*attributes `cron` and `email_addresses` are accepted only when `schedule` is.*"),
			},
			{
				Config:      testUpdateBaselinewithInvalidSchedule + temps.templateSvcTag1,
				ExpectError: regexp.MustCompile(".*please provide a valid email address.*"),
			},
		},
	})
}

var testCreateBaselineWrongCreds = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "invalid"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["invalid"]
		description = "baseline description"
	}
`

var testUpdateBaselinewithInvalidDeviceSvcTag = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "` + TestRefTemplateName + `"
		device_servicetags = ["invalid"]
		description = "baseline description"
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

var testUpdateBaselinewithInvalidTemplate = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_configuration_baseline" "create_baseline" {
		baseline_name = "` + BaselineName + `"
		ref_template_name = "invalid"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		description = "baseline description"
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

var testUpdateBaselinewithUnexpectedCron = `
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
		depends_on = ["ome_template.terraform-acceptance-test-1"]
		cron = "abc"
	}
`

var testUpdateBaselinewithInvalidSchedule = `
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
		depends_on = ["ome_template.terraform-acceptance-test-1"]
		schedule=true
	}
`
