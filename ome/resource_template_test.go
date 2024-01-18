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
	"regexp"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	TemplateName1                           = "test_acc_template-1"
	TemplateName2                           = "test_acc_template-2"
	TemplateNameUpdate2                     = "test_acc_template_update-2"
	ResourceName1                           = "ome_template.terraform-acceptance-test-6"
	ReferenceDeploymentTemplateNameForClone = "test_acc_clone_deployment_template"
	ReferenceComplianceTemplateNameForClone = "test_acc_clone_compliance_template"
	// ContentFilePath                         = "../testdata/test_acc_template.xml"
)

var ContentFilePath = getTestData("test_acc_template.xml")

func init() {
	resource.AddTestSweepers("ome_template", &resource.Sweeper{
		Name:         "ome_template",
		Dependencies: []string{"ome_deployment"},
		F: func(region string) error {
			omeClient, err := getSweeperClient(region)
			if err != nil {
				log.Println("Error getting sweeper client")
				return nil
			}

			_, err = omeClient.CreateSession()
			if err != nil {
				log.Println("Error creating client session for sweeper ")
				return nil
			}
			defer omeClient.RemoveSession()

			templateURL := fmt.Sprintf(clients.TemplateNameContainsAPI, SweepTestsTemplateIdentifier)
			templateResp, templateErr := omeClient.Get(templateURL, nil, nil)
			if templateErr != nil {
				log.Println("failed to fetch templates containing " + SweepTestsTemplateIdentifier)
				return nil
			}

			templateBody, _ := omeClient.GetBodyData(templateResp.Body)
			omeTemplates := models.OMETemplates{}
			//nolint: errcheck
			omeClient.JSONUnMarshal(templateBody, &omeTemplates)

			for _, omeTemplateValue := range omeTemplates.Value {
				_, err = omeClient.Delete(fmt.Sprintf(clients.TemplateAPI+"(%d)", omeTemplateValue.ID), nil, nil)
				if err != nil {
					log.Println("failed to sweep dangling templates.")
					return nil
				}
			}
			return nil
		},
	})
}

func TestTemplateCreation_CreateAndUpdateTemplateSuccess(t *testing.T) {

	t.Log("Create")
	t.Log(testAccCreateTemplateSuccess)
	t.Log("Update")
	t.Log(testAccUpdateTemplateSuccess)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "name", TemplateName1),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "refdevice_servicetag", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "fqdds", "iDRAC,niC"),

					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "name", TemplateName2),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "refdevice_id", DeviceID2),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "description", "This is sample description"),
				),
			},
			{
				Config:      testAccUpdateTemplateWithExistingName,
				ExpectError: regexp.MustCompile(clients.ErrUpdateTemplate),
			},
			{
				Config:      testAccUpdateTemplateWithInvalidVlanNetworkID,
				ExpectError: regexp.MustCompile(clients.ErrUpdateTemplate),
			},

			{
				Config: testAccUpdateTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "name", TemplateName1),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "refdevice_servicetag", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "fqdds", "iDRAC,niC"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "description", "This is a test template"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "vlan.bonding_technology", "NoTeaming"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "vlan.vlan_attributes.0.untagged_network", "10172"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "identity_pool_name", "IO1"),

					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "name", TemplateName2),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "refdevice_id", DeviceID2),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "description", "This is sample description"),
				),
			},
		},
	})
}

/*
root@lglap049:~/terraform-provider-ome# go test ./ome -v -run TestTemplateCreation_CreateAndUpdateTemplateSuccess
=== RUN   TestTemplateCreation_CreateAndUpdateTemplateSuccess
    resource_template_test.go:85: Create
    resource_template_test.go:86:
                provider "ome" {
                        username = "admin"
                        password = "Password123!"
                        host = "10.225.105.1"
                        skipssl = true
                }

                resource "ome_template" "terraform-acceptance-test-1" {
                        name = "test_acc_template-1"
                        refdevice_servicetag = "HRPB0M3"
                        fqdds = "iDRAC,niC"
                }

                resource "ome_template" "terraform-acceptance-test-2" {
                        name = "test_acc_template-2"
                        refdevice_id = 24341
                        description = "This is sample description"
                        job_retry_count  = 10
                        sleep_interval = 60
                }

    resource_template_test.go:87: Update
    resource_template_test.go:88:
                provider "ome" {
                        username = "admin"
                        password = "Password123!"
                        host = "10.225.105.1"
                        skipssl = true
                }

                data "ome_template_info" "template_data" {
                        name = "test_acc_template-1"
                        id = 0
                }

                data "ome_vlannetworks_info" "vlans" {
                }


                resource "ome_template" "terraform-acceptance-test-1" {
                        name = "test_acc_template-1"
                        refdevice_servicetag = "HRPB0M3"
                        attributes = local.template_attributes
                        description = "This is a test template"
                        fqdds = "iDRAC,niC"
                        identity_pool_name   = "IO1"
                        vlan = {
                                propogate_vlan = true
                                bonding_technology = "NoTeaming"
                                vlan_attributes = local.vlan_attributes
                        }
                }

                resource "ome_template" "terraform-acceptance-test-2" {
                        name = "test_acc_template-2"
                        refdevice_id = 24341
                        job_retry_count  = 10
                        sleep_interval = 60
                }


                locals {
                        attributes_value = tomap({
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy": "WarmReset, ColdReset, ACPowerLoss"
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy": "WarmReset, ColdReset, ACPowerLoss"
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered": "WarmReset, ColdReset, ACPowerLoss"
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered": "WarmReset, ColdReset, ACPowerLoss"
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 IOIDOpt Enable" : "Enabled"
                          "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String": "IST"
                        })
                        attributes_is_ignored = tomap({
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy": false
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy": false
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered": false
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered": false
                          "iDRAC,IO Identity Optimization,IOIDOpt 1 IOIDOpt Enable" : false

                        })

                        template_attributes = data.ome_template_info.template_data.attributes != null ? [
                          for attr in data.ome_template_info.template_data.attributes : tomap({
                                attribute_id = attr.attribute_id
                                is_ignored   = lookup(local.attributes_is_ignored, attr.display_name, attr.is_ignored)
                                display_name = attr.display_name
                                value        = lookup(local.attributes_value, attr.display_name, attr.value)
                        })] : null

                        vlan_network_map = {for vlan_network in  data.ome_vlannetworks_info.vlans.vlan_networks : vlan_network.name => vlan_network.vlan_id}

                        vlan_attributes_to_change = tomap({
                                "Integrated NIC 1 - 1": lookup(local.vlan_network_map, "VLAN1", 0)
                                "Integrated NIC 1 - 3": lookup(local.vlan_network_map, "VLAN1", 0)
                        })
                        vlan_attributes = [ for attr in data.ome_template_info.template_data.vlan.vlan_attributes : {
                                is_nic_bonded = attr.is_nic_bonded
                                nic_identifier   = attr.nic_identifier
                                port = attr.port
                                tagged_networks = attr.tagged_networks
                                untagged_network        = lookup(local.vlan_attributes_to_change, "${attr.nic_identifier} - ${attr.port}", attr.untagged_network)
                        }]
                  }

--- PASS: TestTemplateCreation_CreateAndUpdateTemplateSuccess (0.00s)
PASS
ok      terraform-provider-ome/ome      0.033s
*/

// The identity pool and Vlans does not get cloned into the new template in OME.
func TestTemplateCreation_CreateTemplateByCloningSuccess(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateTemplateForClone,
			},
			{
				Config: testAccCloneTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_template.clone-template-test", "name", "test_acc_clone_template_test"),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "view_type", "Deployment"),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "view_type_id", "2"),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "reftemplate_name", ReferenceDeploymentTemplateNameForClone),
					// resource.TestCheckResourceAttr("ome_template.clone-template-test", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "description", "This is a template for testing deployments in acceptance testcases. Please do not delete this template"),

					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "name", "test_acc_clone_template_deployment_compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "view_type", "Compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "view_type_id", "1"),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "reftemplate_name", ReferenceDeploymentTemplateNameForClone),
					// resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "description", "This is a template for testing deployments in acceptance testcases. Please do not delete this template"),

					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "name", "test_acc_clone_template_compliance_compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "view_type", "Compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "view_type_id", "1"),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "reftemplate_name", ReferenceComplianceTemplateNameForClone),
					// resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "description", ""),
				),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplatesInvalidScenarios(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// step 1
				Config:      testAccCreateTemplateInvalidSvcTag,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
			{
				// step 2
				Config:      testAccCreateTemplateMutuallyExclusive1,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 3
				Config:      testAccCreateTemplateMutuallyExclusive2,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 4
				Config:      testAccCreateTemplateMutuallyExclusive3,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 5
				Config:      testAccCreateTemplateMutuallyExclusive4,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 6
				Config:      testAccCreateTemplateWithIOAndVlan,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 7
				Config:      testAccCreateTemplateInvaliddevID,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
			{
				// step 8
				Config:      testAccCreateTemplateEmptyDevice,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 9
				Config:      testAccCreateTemplateInvalidFqdds,
				ExpectError: regexp.MustCompile(clients.ErrInvalidFqdds),
			},
			{
				// step 10
				Config:      testAccCreateTemplateInvalidViewType,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Match.*"),
			},
			{
				// step 11
				Config:      testAccCloneTemplateFailure,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 12
				Config:      testAccCloneTemplateFailureForComplainceToDeployment,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
			{
				// step 13
				Config:      testAccCloneTemplateFailureForDescription,
				ExpectError: regexp.MustCompile(clients.ErrCreateTemplate),
			},
		},
	})
}

func TestTemplateImport_ImportTemplates(t *testing.T) {
	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.NotEmpty(t, s[0].Attributes["attributes.12.display_name"])
		assert.NotEmpty(t, s[0].Attributes["vlan.bonding_technology"])
		assert.Equal(t, TemplateName1, s[0].Attributes["name"])
		assert.Equal(t, "All", s[0].Attributes["fqdds"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            testAccImportTemplateError,
				ResourceName:      ResourceName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "invalid_state_id",
				ExpectError:       regexp.MustCompile(clients.ErrImportTemplate),
			},
			{
				Config: testAccImportTemplateSuccess,
			},
			{
				ResourceName:     ResourceName1,
				ImportState:      true,
				ImportStateCheck: assertTFImportState,
				ExpectError:      nil,
				ImportStateId:    TemplateName1,
			},
		},
	})
}

func TestTemplateCreation_CreateImportTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateImportTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_template.citdtest", "name", "test_acc_import_content_d"),
					resource.TestCheckResourceAttr("ome_template.citdtest", "view_type", "Deployment"),
					resource.TestCheckResourceAttr("ome_template.citdtest", "view_type_id", "2"),
					resource.TestCheckResourceAttr("ome_template.citdtest", "device_type", "Server"),

					resource.TestCheckResourceAttr("ome_template.citctest", "name", "test_acc_import_content_c"),
					resource.TestCheckResourceAttr("ome_template.citctest", "view_type", "Compliance"),
					resource.TestCheckResourceAttr("ome_template.citctest", "view_type_id", "1"),
					resource.TestCheckResourceAttr("ome_template.citctest", "device_type", "Server"),
				),
			},
		},
	})
}

var testAccCreateTemplateForClone = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_template" "terraform-acceptance-test-1" {
		name = "` + ReferenceDeploymentTemplateNameForClone + `"
		refdevice_servicetag = "` + DeviceSvcTag1 + `"
		description = "This is a template for testing deployments in acceptance testcases. Please do not delete this template"
	}

	resource "ome_template" "terraform-acceptance-test-2" {
		name = "` + ReferenceComplianceTemplateNameForClone + `"
		refdevice_servicetag = "` + DeviceSvcTag2 + `"
		view_type = "Compliance"
	}
`

var testAccCreateTemplateSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_template" "terraform-acceptance-test-1" {
		name = "` + TemplateName1 + `"
		refdevice_servicetag = "` + DeviceSvcTag1 + `"
		fqdds = "iDRAC,niC"
	}

	resource "ome_template" "terraform-acceptance-test-2" {
		name = "` + TemplateName2 + `"
		refdevice_id = ` + DeviceID2 + `
		description = "This is sample description"
		job_retry_count  = 10
		sleep_interval = 60
	}
`

var testAccUpdateTemplateSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_template_info" "template_data" {
		name = "` + TemplateName1 + `"
		id = 0
	}

	data "ome_vlannetworks_info" "vlans" {
	}


	resource "ome_template" "terraform-acceptance-test-1" {
		name = "` + TemplateName1 + `"
		refdevice_servicetag = "` + DeviceSvcTag1 + `"
		attributes = local.template_attributes
		description = "This is a test template"
		fqdds = "iDRAC,niC"
		identity_pool_name   = "IO1"
		vlan = {
			propogate_vlan = true
			bonding_technology = "NoTeaming"
			vlan_attributes = local.vlan_attributes
		}
	}

	resource "ome_template" "terraform-acceptance-test-2" {
		name = "` + TemplateName2 + `"
		refdevice_id = ` + DeviceID2 + `
		job_retry_count  = 10
		sleep_interval = 60
	}


	locals {
		attributes_value = tomap({
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy": "WarmReset, ColdReset, ACPowerLoss"
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy": "WarmReset, ColdReset, ACPowerLoss"
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered": "WarmReset, ColdReset, ACPowerLoss"
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered": "WarmReset, ColdReset, ACPowerLoss"
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 IOIDOpt Enable" : "Enabled"
		  "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String": "IST"
		})
		attributes_is_ignored = tomap({
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy": false
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy": false
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered": false
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered": false
		  "iDRAC,IO Identity Optimization,IOIDOpt 1 IOIDOpt Enable" : false
	  
		})
		
		template_attributes = data.ome_template_info.template_data.attributes != null ? [
		  for attr in data.ome_template_info.template_data.attributes : tomap({
			attribute_id = attr.attribute_id
			is_ignored   = lookup(local.attributes_is_ignored, attr.display_name, attr.is_ignored)
			display_name = attr.display_name
			value        = lookup(local.attributes_value, attr.display_name, attr.value)
		})] : null

		vlan_network_map = {for vlan_network in  data.ome_vlannetworks_info.vlans.vlan_networks : vlan_network.name => vlan_network.vlan_id}
		
		vlan_attributes_to_change = tomap({
			"Integrated NIC 1 - 1": lookup(local.vlan_network_map, "VLAN1", 0)
			"Integrated NIC 1 - 3": lookup(local.vlan_network_map, "VLAN1", 0)
		})
		vlan_attributes = [ for attr in data.ome_template_info.template_data.vlan.vlan_attributes : {
			is_nic_bonded = attr.is_nic_bonded
			nic_identifier   = attr.nic_identifier
			port = attr.port
			tagged_networks = attr.tagged_networks
			untagged_network        = lookup(local.vlan_attributes_to_change, "${attr.nic_identifier} - ${attr.port}", attr.untagged_network)
		}]
	  }
`

var testAccCreateTemplateInvalidSvcTag = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-3" {
	name = "test_acc_template-3"
	refdevice_servicetag = "TEST"
}
`

var testAccCreateTemplateMutuallyExclusive1 = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-3" {
	name = "test_acc_template-3"
	reftemplate_name = "template_6"
	refdevice_id = 10112
}
`

var testAccCreateTemplateMutuallyExclusive2 = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-3" {
	name = "test_acc_template-3"
	reftemplate_name = "template_6"
	refdevice_servicetag = "CZMC1T2"
}
`

var testAccCreateTemplateMutuallyExclusive3 = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-3" {
	name = "test_acc_template-3"
	refdevice_id = 10112
	refdevice_servicetag = "CZMC1T2"
}
`

var testAccCreateTemplateMutuallyExclusive4 = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-3" {
	name = "test_acc_template-3"
	refdevice_id = 10112
	content = "CZMC1T2"
}
`

var testAccCreateTemplateWithIOAndVlan = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-3" {
	name = "test_acc_template-3"
	refdevice_id = 10112
	identity_pool_name = "IO1"
	vlan = {
		propogate_vlan = true
		bonding_technology = "NoTeaming"
		vlan_attributes = [
			{
				is_nic_bonded =  false,
				nic_identifier = "Integrated NIC 1"
				port = 1
				tagged_networks = [0]
				untagged_network = 0
			},
			{
				is_nic_bonded =  false
				nic_identifier = "Integrated NIC 1"
				port = 2
				tagged_networks = [0]
				untagged_network = 0
			},
			{
				is_nic_bonded =  false,
				nic_identifier = "Integrated NIC 1"
				port = 3
				tagged_networks = [0]
				untagged_network = 0
			},
			{
				is_nic_bonded =  false
				nic_identifier = "Integrated NIC 1"
				port = 4
				tagged_networks = [0]
				untagged_network = 0 
			}
		]
	}
}`

var testAccCreateTemplateEmptyDevice = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-4" {
	name = "test_acc_template-4"
}
`

var testAccCreateTemplateInvaliddevID = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-4" {
	name = "test_acc_template-5"
	refdevice_id = 123
}
`

var testAccCreateTemplateInvalidFqdds = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-5" {
	name = "test_acc_template-6"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
	fqdds = "Test fqdds"
}
`

var testAccCreateTemplateInvalidViewType = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-5" {
	name = "test_acc_template-7"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
	view_type = "Test View type"
}
`

var testAccImportTemplateError = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_template" "terraform-acceptance-test-6" {
	}
`

var testAccImportTemplateSuccess = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	resource "ome_template" "terraform-acceptance-test-6" {
		name = "` + TemplateName1 + `"
		refdevice_servicetag = "` + DeviceSvcTag1 + `"
	}
`

var testAccUpdateTemplateWithExistingName = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-1" {
	name = "` + TemplateName1 + `"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
	fqdds = "iDRAC,niC"
}

resource "ome_template" "terraform-acceptance-test-2" {
	name = "` + TemplateName1 + `"
	refdevice_id = ` + DeviceID2 + `
	description = "This is sample description"
	job_retry_count  = 10
	sleep_interval = 60
}
`

var testAccUpdateTemplateWithInvalidVlanNetworkID = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-1" {
	name = "` + TemplateName1 + `"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
	fqdds = "iDRAC,niC"
	vlan = {
		propogate_vlan     = true
		bonding_technology = "NoTeaming"
		vlan_attributes = [
		  {
			untagged_network = 0
			tagged_networks  = [123]
			is_nic_bonded    = false
			port             = 1
			nic_identifier   = "NIC in Mezzanine 1A"
		  },
		  {
			untagged_network = 0
			tagged_networks = [0]
			is_nic_bonded = false
			port = 2
			nic_identifier = "NIC in Mezzanine 1A"
		}
		]
	}
}

resource "ome_template" "terraform-acceptance-test-2" {
	name = "` + TemplateName2 + `"
	refdevice_id = ` + DeviceID2 + `
	description = "This is sample description"
	job_retry_count  = 10
	sleep_interval = 60
}
`
var testAccCreateImportTemplateSuccess = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "citdtest" {
	name = "test_acc_import_content_d"
	content = file("` + ContentFilePath + `")
}

resource "ome_template" "citctest" {
	name = "test_acc_import_content_c"
	content = file("` + ContentFilePath + `")
	view_type = "Compliance"
}
`

var testAccCloneTemplateSuccess = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-1" {
	name = "` + ReferenceDeploymentTemplateNameForClone + `"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
}

resource "ome_template" "terraform-acceptance-test-2" {
	name = "` + ReferenceComplianceTemplateNameForClone + `"
	refdevice_servicetag = "` + DeviceSvcTag2 + `"
	view_type = "Compliance"
}

resource "ome_template" "clone-template-test" {
	name = "test_acc_clone_template_test"
	reftemplate_name = "` + ReferenceDeploymentTemplateNameForClone + `"
}

resource "ome_template" "clone-template-deployment-compliance" {
	name = "test_acc_clone_template_deployment_compliance"
	reftemplate_name = "` + ReferenceDeploymentTemplateNameForClone + `"
	view_type = "Compliance"
}

resource "ome_template" "clone-template-compliance-compliance" {
	name = "test_acc_clone_template_compliance_compliance"
	reftemplate_name = "` + ReferenceComplianceTemplateNameForClone + `"
	view_type = "Compliance"
}
`

var testAccCloneTemplateFailure = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "invalid-template-name"
}
`

var testAccCloneTemplateFailureForComplainceToDeployment = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "test_acc_compliance_template"
}
`

var testAccCloneTemplateFailureForDescription = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "ReferenceDeploymentTemplateNameForClone"
	description = "This is invalid desc."
}
`
