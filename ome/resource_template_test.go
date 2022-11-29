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
	TemplateName1                           = "test_acc_template-1"
	TemplateName2                           = "test_acc_template-2"
	TemplateNameUpdate2                     = "test_acc_template_update-2"
	ResourceName1                           = "ome_template.terraform-acceptance-test-6"
	ReferenceDeploymentTemplateNameForClone = "test-clone-reference-template"
	ReferenceComplianceTemplateNameForClone = "test-compliance-template"
)

func TestTemplateCreation_CreateAndUpdateTemplateSuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
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
				Config:      testAccUpdateTemplateWithInvalidAttributeID,
				ExpectError: regexp.MustCompile("Unable to update the template"),
			},
			{
				Config:      testAccUpdateTemplateWithExistingName,
				ExpectError: regexp.MustCompile("Unable to update the template"),
			},
			{
				Config:      testAccUpdateTemplateWithInvalidVlanNetworkID,
				ExpectError: regexp.MustCompile("Unable to update the template"),
			},
			{
				Config: testAccUpdateTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "name", TemplateName1),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "refdevice_servicetag", DeviceSvcTag1),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "fqdds", "iDRAC,niC"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "description", "This is a test template"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "vlan.bonding_technology", "NoTeaming"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "vlan.vlan_attributes.0.untagged_network", "10133"),

					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "name", TemplateName2),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "refdevice_id", DeviceID2),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "description", "This is sample description"),
				),
			},
		},
	})
}

// The identity pool and Vlans does not get cloned into the new template in OME.
func TestTemplateCreation_CreateTemplateByCloningSuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloneTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_template.clone-template-test", "name", "clone-template-test"),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "view_type", "Deployment"),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "view_type_id", "2"),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "reftemplate_name", ReferenceDeploymentTemplateNameForClone),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "refdevice_id", DeviceID1),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "description", "This is a template for testing deployments in acceptance testcases. Please do not delete this template."),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "content", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-test", "attributes.0.display_name", "EventFilters,EventFilters.Audit.1,Event Filters,Action FSD 4 2"),

					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "name", "clone-template-deployment-compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "view_type", "compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "view_type_id", "1"),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "reftemplate_name", ReferenceDeploymentTemplateNameForClone),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "refdevice_id", DeviceID1),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "description", "This is a template for testing deployments in acceptance testcases. Please do not delete this template."),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "content", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-deployment-compliance", "attributes.0.display_name", "EventFilters,EventFilters.Audit.1,Event Filters,Action FSD 4 2"),

					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "name", "clone-template-compliance-compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "view_type", "Compliance"),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "view_type_id", "1"),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "reftemplate_name", ReferenceComplianceTemplateNameForClone),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "refdevice_id", DeviceID1),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "description", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "content", ""),
					resource.TestCheckResourceAttr("ome_template.clone-template-compliance-compliance", "attributes.0.display_name", "BIOS,BIOS Boot Settings,Boot Sequence"),
				),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplatesInvalidScenarios(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccCreateTemplateInvalidSvcTag,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
			{
				Config:      testAccCreateTemplateMutuallyExclusive1,
				ExpectError: regexp.MustCompile("error creating/cloning the template"),
			},
			{
				Config:      testAccCreateTemplateMutuallyExclusive2,
				ExpectError: regexp.MustCompile("error creating/cloning the template"),
			},
			{
				Config:      testAccCreateTemplateMutuallyExclusive3,
				ExpectError: regexp.MustCompile("error creating/cloning the template"),
			},
			{
				Config:      testAccCreateTemplateWithIOAndVlan,
				ExpectError: regexp.MustCompile("error creating the template"),
			},
			{
				Config:      testAccCreateTemplateInvaliddevID,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
			{
				Config:      testAccCreateTemplateEmptyDevice,
				ExpectError: regexp.MustCompile("either reftemplate_name or refdevice_id or refdevice_servicetag"),
			},
			{
				Config:      testAccCreateTemplateInvalidFqdds,
				ExpectError: regexp.MustCompile(clients.ErrInvalidFqdds),
			},
			{
				Config:      testAccCreateTemplateInvalidViewType,
				ExpectError: regexp.MustCompile(clients.ErrInvalidTemplateViewType),
			},
			{
				Config:      testAccCloneTemplateFailure,
				ExpectError: regexp.MustCompile("error cloning the template with given reference template name"),
			},
			{
				Config:      testAccCloneTemplateFailureForComplainceToDeployment,
				ExpectError: regexp.MustCompile("cannot clone compliance template as deployment template."),
			},
			{
				Config:      testAccCloneTemplateFailureForDescription,
				ExpectError: regexp.MustCompile("description will be copied from the reference template."),
			},
		},
	})
}

func TestTemplateImport_ImportTemplateError(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:            testAccImportTemplateError,
				ResourceName:      ResourceName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "invalid_state_id",
				ExpectError:       regexp.MustCompile("unable to get template"),
			},
		},
	})
}

func TestTemplateImport_ImportTemplateSuccess(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	assertTFImportState := func(s []*terraform.InstanceState) error {
		assert.NotEmpty(t, s[0].Attributes["attributes.12.display_name"])
		assert.NotEmpty(t, s[0].Attributes["vlan.bonding_technology"])
		assert.Equal(t, TemplateName1, s[0].Attributes["name"])
		assert.Equal(t, "", s[0].Attributes["fqdds"])
		assert.Equal(t, 1, len(s))
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
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

// func checkResourceUpdate_attachIdentityPool(t *testing.T, p tfsdk.Provider) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		templateAPI := "/api/TemplateService/Templates?$filter=Name eq '%s'"
// 		template1URL := fmt.Sprintf(templateAPI, TemplateName1)
// 		provider := p.(*provider)
// 		omeClient, err := clients.NewClient(*provider.clientOpt)
// 		if err != nil {
// 			return fmt.Errorf("Unable to create client %s", err.Error())
// 		}

// 		_, err = omeClient.CreateSession()
// 		if err != nil {
// 			return fmt.Errorf("Error creating client session %s", err.Error())
// 		}

// 		response, err := omeClient.Get(template1URL, nil, nil)
// 		assert.Nil(t, err)
// 		assert.NotNil(t, response)
// 		b, _ := omeClient.GetBodyData(response.Body)

// 		omeTemplates1 := models.OMETemplates{}
// 		err = omeClient.JSONUnMarshal(b, &omeTemplates1)
// 		if err != nil {
// 			fmt.Printf("Unable to create client %s", err.Error())
// 		}

// 		omeTemplate1 := omeTemplates1.Value[0]
// 		assert.Equal(t, int32(1), omeTemplate1.IdentityPoolID)
// 		return nil

// 	}
// }

// func checkResourceUpdate_detachIdentityPool(t *testing.T, p tfsdk.Provider) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		templateAPI := "/api/TemplateService/Templates?$filter=Name eq '%s'"
// 		template2URL := fmt.Sprintf(templateAPI, TemplateName2)
// 		provider := p.(*provider)
// 		omeClient, err := clients.NewClient(*provider.clientOpt)
// 		if err != nil {
// 			return fmt.Errorf("Unable to create client %s", err.Error())
// 		}

// 		_, err = omeClient.CreateSession()
// 		if err != nil {
// 			return fmt.Errorf("Error creating client session %s", err.Error())
// 		}

// 		response, err := omeClient.Get(template2URL, nil, nil)
// 		assert.Nil(t, err)
// 		assert.NotNil(t, response)
// 		b, _ := omeClient.GetBodyData(response.Body)

// 		omeTemplates2 := models.OMETemplates{}
// 		err = omeClient.JSONUnMarshal(b, &omeTemplates2)
// 		if err != nil {
// 			fmt.Printf("Unable to create client %s", err.Error())
// 		}

// 		omeTemplate2 := omeTemplates2.Value[0]
// 		assert.Equal(t, int32(0), omeTemplate2.IdentityPoolID)
// 		return nil

// 	}
// }

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

	data "ome_template_info" "terraform-template-data-1" {
		name = "` + TemplateName1 + `"
		id = 0
	}

	resource "ome_template" "terraform-acceptance-test-1" {
		name = "` + TemplateName1 + `"
		refdevice_servicetag = "` + DeviceSvcTag1 + `"
		attributes = local.template_attributes
		description = "This is a test template"
		fqdds = "iDRAC,niC"
		vlan = {
			propogate_vlan = true
			bonding_technology = "NoTeaming"
			vlan_attributes = [
				{
					untagged_network = 10133
					tagged_networks = [0]
					is_nic_bonded = false
					port = 1
					nic_identifier = "NIC in Mezzanine 1A"
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
		job_retry_count  = 10
		sleep_interval = 60
	}

	locals {
		template_attributes = data.ome_template_info.terraform-template-data-1.attributes != null ? [
		  for attr in  data.ome_template_info.terraform-template-data-1.attributes: tomap({
			  attribute_id = attr.attribute_id
		is_ignored = attr.is_ignored
		display_name = attr.display_name
		value = attr.display_name == "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String" ?  "IST" : attr.value
		  })] : null
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
	refdevice_id = 12328
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
	refdevice_servicetag = "MX1404"
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
	refdevice_id = 12328
	refdevice_servicetag = "MX1404"
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
	refdevice_id = 12328
	identity_pool_name = "IO1"
	vlan = {
		propogate_vlan = true
		bonding_technology = "NoTeaming"
		vlan_attributes = [
			{
				untagged_network = 10133
				tagged_networks = [0]
				is_nic_bonded = false
				port = 1
				nic_identifier = "NIC in Mezzanine 1A"
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

var testAccUpdateTemplateWithInvalidAttributeID = `

provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

data "ome_template_info" "terraform-template-data-1" {
	name = "` + TemplateName1 + `"
	id = 0
}

resource "ome_template" "terraform-acceptance-test-1" {
	name = "` + TemplateName1 + `"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
	attributes = local.template_attributes
	fqdds = "iDRAC,niC"
}

resource "ome_template" "terraform-acceptance-test-2" {
	name = "` + TemplateName2 + `"
	refdevice_id = ` + DeviceID2 + `
	description = "This is sample description"
	job_retry_count  = 10
	sleep_interval = 60
}

locals {
	template_attributes = data.ome_template_info.terraform-template-data-1.attributes != null ? [
	  for attr in  data.ome_template_info.terraform-template-data-1.attributes: tomap({
	is_ignored = attr.is_ignored
	display_name = attr.display_name
	attribute_id = attr.display_name == "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String" ?  123 : attr.attribute_id
	value = attr.display_name == "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String" ?  "IST" : attr.value
	  })] : null
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

var testAccCloneTemplateSuccess = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "` + ReferenceDeploymentTemplateNameForClone + `"
}

resource "ome_template" "clone-template-deployment-compliance" {
	name = "clone-template-deployment-compliance"
	reftemplate_name = "` + ReferenceDeploymentTemplateNameForClone + `"
	view_type = "compliance"
}

resource "ome_template" "clone-template-compliance-compliance" {
	name = "clone-template-compliance-compliance"
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

// var testAccUpdateTemplateSuccess_attachIdentityPool = `

// 	provider "ome" {
// 		username = "` + omeUserName + `"
// 		password = "` + omePassword + `"
// 		host = "` + omeHost + `"
// 		skipssl = true
// 	}

// 	data "ome_template_info" "terraform-template-data-1" {
// 		name = "` + TemplateName1 + `"
// 		id = 0
// 	}

// 	resource "ome_template" "terraform-acceptance-test-1" {
// 		name = "` + TemplateName1 + `"
// 		refdevice_servicetag = "` + DeviceSvcTag1 + `"
// 		attributes = data.ome_template_info.terraform-template-data-1.attributes
// 		identity_pool_name = "IO1"
// 		fqdds = "iDRAC,NIC"
// 	}
// `

// var testAccUpdateSuccess_detachIdentityPool = `
// 	provider "ome" {
// 		username = "` + omeUserName + `"
// 		password = "` + omePassword + `"
// 		host = "` + omeHost + `"
// 		skipssl = true
// 	}

// 	data "ome_template_info" "terraform-template-data-1" {
// 		name = "` + TemplateName1 + `"
// 		id = 0
// 	}

// 	resource "ome_template" "terraform-acceptance-test-1" {
// 		name = "` + TemplateName1 + `"
// 		refdevice_servicetag = "` + DeviceSvcTag1 + `"
// 		attributes = data.ome_template_info.terraform-template-data-1.attributes
// 		fqdds = "iDRAC"
// 		identity_pool_name=""
// 	}
// 	`
