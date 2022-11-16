package ome

import (
	"fmt"
	"os"
	"regexp"
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	BaselineName              = "test_acc_create_baseline"
	BaselineNameUpdate        = "test_acc_update_baseline"
	TestRefTemplateName       = "test_acc_compliance_template"
	TestRefTemplateNameUpdate = "test_acc_compliance_template_update"
)

func TestCreateBaseline_TestValidations(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
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
				Config:      testCreateBaselineValidationFailureNonComplianceTemplateID,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 4
				Config:      testCreateBaselineValidationFailureEmptyDevice,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 5
				Config:      testCreateBaselineValidationFailureInvalidDevice,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 6
				Config:      testCreateBaselineValidationFailureScheduleNotificationEmptyEmail,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 7
				Config:      testCreateBaselineValidationFailureScheduleNotificationInvalidEmail,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 8
				Config:      testCreateBaselineValidationFailureNotificationOnScheduleEmptyCron,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
			{ // Step 9
				Config:      testCreateBaselineValidationInvalidOutputFormatCase,
				ExpectError: regexp.MustCompile(fmt.Sprintf("Allowed values are one of  :  %s", clients.ValidOutputFormat)),
			},
			{ // Step 10
				Config:      testCreateBaselineValidationInvalidOutputFormat,
				ExpectError: regexp.MustCompile(fmt.Sprintf("Allowed values are one of  :  %s", clients.ValidOutputFormat)),
			},
			{ // Step 11
				Config:      testCreateBaselineValidationDeviceCapable,
				ExpectError: regexp.MustCompile(clients.ErrGnrCreateBaseline),
			},
		},
	})
}

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
		schedule_notification=true
		email_addresses=["abc@gmail.com"]
		notification_on_schedule=true
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
		schedule_notification=true
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
		schedule_notification=true
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
		device_servicetags = [` + DeviceID3 + `]
	}
`

func TestCreateBaseline_BaselineWithDeviceID(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, BaselineNameUpdate, s[0].Attributes["baseline_name"])
		assert.Equal(t, DeviceSvcTag2, s[0].Attributes["device_servicetags.0"])
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselinewithDeviceID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0")),
			},
			{
				Config: testConfigureBaselinewithDeviceIDUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineNameUpdate),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description updated"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateNameUpdate),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag2),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0")),
			},
			{
				Config:           testImportConfigurationBaseline,
				ResourceName:     "ome_configuration_baseline.import_baseline",
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    BaselineNameUpdate,
			},
		},
	})
}

var testConfigureBaselinewithDeviceID = `
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
`

var testConfigureBaselinewithDeviceIDUpdate = `
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
	}
`

func TestCreateBaseline_CreateBaselineWithSchedule(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, BaselineName, s[0].Attributes["baseline_name"])
		assert.Equal(t, DeviceSvcTag1, s[0].Attributes["device_servicetags.0"])
		assert.Equal(t, "true", s[0].Attributes["schedule_notification"])
		assert.Equal(t, "test@mail.com", s[0].Attributes["email_addresses.0"])
		assert.Equal(t, "0 50 8 * * ? *", s[0].Attributes["cron"])
		assert.Equal(t, "true", s[0].Attributes["notification_on_schedule"])
		assert.Equal(t, "html", s[0].Attributes["output_format"])
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselineWithSchedule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule_notification", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.0", "test@testmail.com"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "cron", "0 49 8 * * ? *"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "notification_on_schedule", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "output_format", "csv"),
				),
			},
			{
				Config: testConfigureBaselineWithScheduleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule_notification", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.0", "test@mail.com"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "cron", "0 50 8 * * ? *"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "notification_on_schedule", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "output_format", "html"),
				),
			},
			{
				Config:           testImportConfigurationBaseline,
				ResourceName:     "ome_configuration_baseline.import_baseline",
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    BaselineName,
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
		schedule_notification = true
		notification_on_schedule = true
		email_addresses = ["test@testmail.com"]
		cron = "0 49 8 * * ? *"
		output_format = "csv"
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
		schedule_notification = true
		notification_on_schedule = true
		email_addresses = ["test@mail.com"]
		cron = "0 50 8 * * ? *"
		output_format = "html"
	}
`

func TestCreateBaseline_CreateBaselineWithScheduleNonCompliant(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testConfigureBaselineScheduleNonCompliant,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule_notification", "true"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "email_addresses.0", "test@testmail.com"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "output_format", "html"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "notification_on_schedule", "false")),
			},
			{
				Config: testConfigureBaselineScheduleNonCompliantUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "baseline_name", BaselineName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "description", "baseline description"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "ref_template_name", TestRefTemplateName),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "device_ids.#", "0"),
					resource.TestCheckResourceAttr("ome_configuration_baseline.create_baseline", "schedule_notification", "false"),
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
		schedule_notification = true
		email_addresses = ["test@testmail.com"]
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
