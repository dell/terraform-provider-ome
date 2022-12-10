package ome

import (
	"os"
	"regexp"
	"terraform-provider-ome/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	TestAccTemplateName       = "test_acc_create_deployment"
	TestAccUpdateTemplateName = "test_acc_update_deployment"
)

func TestTemplateDeploy_InvalidTemplate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testTemplateDeploymentIDOrNameRequired,
				ExpectError: regexp.MustCompile(clients.ErrInvalidTemplate),
			},
			{
				Config:      testTemplateDeploymentIDSTGNMutuallyExclusive1,
				ExpectError: regexp.MustCompile(clients.ErrTemplateDeploymentCreate),
			},
			{
				Config:      testTemplateDeploymentDeviceInfoRequired,
				ExpectError: regexp.MustCompile(clients.ErrTemplateDeploymentCreate),
			},
		},
	})
}

func TestTemplateDeploy_CreateAndUpdateDeploySuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testTemplateDeploymentSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "template_name", TestAccTemplateName),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.0", DeviceSvcTag1),
				),
			},
			{
				Config:      testUpdateTemplateDeployWithInvalidTemplate,
				ExpectError: regexp.MustCompile(clients.ErrTemplateDeploymentUpdate),
			},
			{
				Config: testTemplateUpdateDeploymentSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "template_name", TestAccTemplateName),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.0", DeviceSvcTag2),
				),
			},
		},
	})
}

func TestTemplateDeploy_CreateUpdateDeployWithScheduleSuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testTemplateDeploymentSuccessWithSchedule,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "template_name", TestAccTemplateName),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "run_later", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "cron", "0 00 11 14 02 ? 2032")),
			},
			{
				Config: testTemplateUpdateDeploymentWithScheduleSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "template_name", TestAccTemplateName),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.0", DeviceSvcTag2),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "run_later", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "cron", "0 00 11 14 02 ? 2032")),
			},
		},
	})
}

func TestTemplateDeploy_ImportDeploymentError(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.Equal(t, 1, len(s))
		assert.Equal(t, TestAccTemplateName, s[0].Attributes["template_name"])
		assert.Equal(t, DeviceSvcTag1, s[0].Attributes["device_servicetags.0"])
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:            testAccImportDeploymentError,
				ResourceName:      "ome_deployment.import-deployment-error",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "invalid_state_id",
				ExpectError:       regexp.MustCompile(clients.ErrImportDeployment),
			},
			{
				Config:            testAccImportDeploymentError,
				ResourceName:      "ome_deployment.import-deployment-error",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     TestAccUpdateTemplateName,
				ExpectError:       regexp.MustCompile(clients.ErrImportDeployment),
			},
			{
				Config: testTemplateDeploymentSuccess,
			},
			{
				Config:           testAccImportDeploymentSuccess,
				ResourceName:     "ome_deployment.import-deployment-success",
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TestAccTemplateName,
			},
		},
	})
}

func TestTemplateDeploy_CreateDeployBootNetworkISOSuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testTemplateDeploymentbootToNetworkISOSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "template_name", TestAccTemplateName),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.0", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.boot_to_network", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_type", "CIFS"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.iso_timeout", "240"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.iso_path", "/cifsshare/unattended/unattended_rocky8.6.iso"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.ip_address", ShareIP),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.share_name", ""),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.work_group", ""),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.user", ShareUser),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.password", SharePassword),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "power_state_off", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "forced_shutdown", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "options_continue_on_warning", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "options_strict_checking_vlan", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "options_precheck_only", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_attributes.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_attributes.0.attributes.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_attributes.0.attributes.0.display_name", "ServerTopology 1 Aisle Name")),
			},
			{
				Config: testTemplateUpdateDeployWithParamsSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "template_name", TestAccTemplateName),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_servicetags.0", DeviceSvcTag2),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.boot_to_network", "true"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_type", "CIFS"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.iso_timeout", "240"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.iso_path", "/cifsshare/unattended/unattended_rocky8.6.iso"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.ip_address", ShareIP),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.share_name", ""),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.work_group", ""),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.user", ShareUser),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "boot_to_network_iso.share_detail.password", SharePassword),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_attributes.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_attributes.0.attributes.#", "1"),
					resource.TestCheckResourceAttr("ome_deployment.deploy-template-3", "device_attributes.0.attributes.0.display_name", "ServerTopology 1 Aisle Name"),
				),
			},
		},
	})
}

// Add resource as applicable
var testTemplateDeploymentIDSTGNMutuallyExclusive1 = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-1" {
		template_name = "demo_template_1"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		device_ids = [` + DeviceID1 + `]
	}
`

var testTemplateDeploymentIDOrNameRequired = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-1" {

	}
`

var testTemplateDeploymentDeviceInfoRequired = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-2" {
		template_name = "demo_template_1"
	}
`
var testUpdateTemplateDeployWithInvalidTemplate = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-3" {
		template_name = "invalid_template_name"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
	}
`

var testTemplateDeploymentSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-3" {
		template_name = "` + TestAccTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
	}
`

var testTemplateUpdateDeploymentSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-3" {
		template_name = "` + TestAccTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag2 + `"]
	}
`

var testTemplateDeploymentbootToNetworkISOSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-3" {
		template_name = "` + TestAccTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		power_state_off = true 
		forced_shutdown = true
		options_continue_on_warning = true
		options_strict_checking_vlan = true
		options_precheck_only = true
		boot_to_network_iso = {
			boot_to_network = true
			share_type = "CIFS"
			iso_timeout = 240
			iso_path = "/cifsshare/unattended/unattended_rocky8.6.iso"
			share_detail = {
				ip_address = "` + ShareIP + `"
				share_name = ""
				work_group = ""
				user = "` + ShareUser + `"
				password = "` + SharePassword + `"
			}
		}
		device_attributes = [
			{
				device_servicetags = ["` + DeviceSvcTag1 + `"]
				attributes = [
					{
						attribute_id = 1197404
						display_name = "ServerTopology 1 Aisle Name"
						value = "IST"
						is_ignored = false
					}
				]
			}
		]
		job_retry_count = 30
	}
`

var testTemplateDeploymentSuccessWithSchedule = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-3" {
		template_name = "` + TestAccTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag1 + `"]
		run_later = true
		cron = "0 00 11 14 02 ? 2032"
	}
`

var testTemplateUpdateDeployWithParamsSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-3" {
		template_name = "` + TestAccTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag2 + `"]
		boot_to_network_iso = {
			boot_to_network = true
			share_type = "CIFS"
			iso_timeout = 240
			iso_path = "/cifsshare/unattended/unattended_rocky8.6.iso"
			share_detail = {
				ip_address = "` + ShareIP + `"
				share_name = ""
				work_group = ""
				user = "` + ShareUser + `"
				password = "` + SharePassword + `"
			}
		}
		job_retry_count = 30
		device_attributes = [
			{
				device_servicetags = ["` + DeviceSvcTag2 + `"]
				attributes = [
					{
						attribute_id = 1197404
						display_name = "ServerTopology 1 Aisle Name"
						value = "IST"
						is_ignored = false
					}
				]
			}
		]
	}
`

var testTemplateUpdateDeploymentWithScheduleSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "deploy-template-3" {
		template_name = "` + TestAccTemplateName + `"
		device_servicetags = ["` + DeviceSvcTag2 + `"]
		run_later = true
		cron = "0 00 11 14 02 ? 2032"
	}
`

var testAccImportDeploymentError = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "import-deployment-error" {
	}
`

var testAccImportDeploymentSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_deployment" "import-deployment-success" {
	}
`
