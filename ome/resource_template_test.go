package ome

import (
	"fmt"
	"regexp"
	"strconv"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func TestTemplateCreation_CreateTemplatesSuccess(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(checkResourceCreation(t, testProvider),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "reftemplate_name", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "reftemplate_name", "")),
			},
		},
	})
}

// The identity pool and Vlans does not get cloned into the new template in OME.
func TestTemplateCreation_CreateTemplateByCloningSuccess(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloneTemplateSuccess,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "name", "clone-template-test"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "view_type", "Deployment"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "view_type_id", "2"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "reftemplate_name", ReferenceDeploymentTemplateNameForClone),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "refdevice_id", DeviceID1),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "description", "This is a template for testing deployments in acceptance testcases. Please do not delete this template."),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "content", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "attributes.0.display_name", "EventFilters,EventFilters.Audit.1,Event Filters,Action FSD 4 2")),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateByCloningSuccessForDeploymentToCompliance(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloneTemplateSuccessForDeploymentToCompliance,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "name", "clone-template-test"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "view_type", "Compliance"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "view_type_id", "1"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "reftemplate_name", ReferenceDeploymentTemplateNameForClone),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "refdevice_id", DeviceID1),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "description", "This is a template for testing deployments in acceptance testcases. Please do not delete this template."),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "content", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "attributes.0.display_name", "EventFilters,EventFilters.Audit.1,Event Filters,Action FSD 4 2")),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateByCloningSuccessForComplianceToCompliance(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloneTemplateSuccessForComplianceToCompliance,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "name", "clone-template-test"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "view_type", "Compliance"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "view_type_id", "1"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "reftemplate_name", ReferenceComplianceTemplateNameForClone),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "refdevice_id", DeviceID1),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "fqdds", "All"),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "refdevice_servicetag", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "description", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "content", ""),
					resource.TestCheckResourceAttr("ome_template.terraform-clone-template-test", "attributes.0.display_name", "BIOS,BIOS Boot Settings,Boot Sequence")),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateByCloningFailureForRefTemplateAndDevice(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloneTemplateFailureForRefTemplateAndDeviceID,
				ExpectError: regexp.MustCompile("please provide either reftemplate_name or refdevice_id/refdevice_servicetag"),
			},
			{
				Config:      testAccCloneTemplateFailureForRefTemplateAndDeviceServiceTag,
				ExpectError: regexp.MustCompile("please provide either reftemplate_name or refdevice_id/refdevice_servicetag"),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateByCloningFailure(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloneTemplateFailure,
				ExpectError: regexp.MustCompile("error cloning the template with given reference template name"),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateByCloningFailureForComplainceToDeployment(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloneTemplateFailureForComplainceToDeployment,
				ExpectError: regexp.MustCompile("cannot clone compliance template as deployment template."),
			},
		},
	})
}

func TestTemplateCreation_CreateUpdateSuccess(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{

			{
				Config: testAccCreateTemplateSuccess,
				Check:  resource.ComposeTestCheckFunc(checkResourceCreation(t, testProvider)),
			},

			{
				Config: testAccUpdateTemplateSuccess,
				Check:  resource.ComposeTestCheckFunc(checkResourceUpdate(t, testProvider)),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplatesInvalidSvcTag(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{

			{
				Config:      testAccCreateTemplateInvalidSvcTag,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplatesInvalidDevID(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{

			{
				Config:      testAccCreateTemplateInvaliddevID,
				ExpectError: regexp.MustCompile(clients.ErrInvalidDeviceIdentifiers),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateEmptyDeviceDetails(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{

			{
				Config:      testAccCreateTemplateEmptyDevice,
				ExpectError: regexp.MustCompile(clients.ErrEmptyDeviceDetails),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateInvalidFqdds(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{

			{
				Config:      testAccCreateTemplateInvalidFqdds,
				ExpectError: regexp.MustCompile(clients.ErrInvalidFqdds),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateInvalidTemplateViewType(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{

			{
				Config:      testAccCreateTemplateInvalidViewType,
				ExpectError: regexp.MustCompile(clients.ErrInvalidTemplateViewType),
			},
		},
	})
}

func TestTemplateUpdation_UpdateTemplateWithInvalidAttributeId(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateTemplateSuccess,
			},
			{
				Config:      testAccUpdateTemplateWithInvalidAttributeID,
				ExpectError: regexp.MustCompile("Unable to update the template"),
			},
		},
	})
}

func TestTemplateCreation_CreateTemplateWithExistingName(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateTemplateWithExistingName,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-1", "name", TemplateName1),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "name", TemplateName2),
					resource.TestCheckResourceAttr("ome_template.terraform-acceptance-test-2", "description", "")),
			},
			{
				Config:      testAccUpdateTemplateWithExistingName,
				ExpectError: regexp.MustCompile("Unable to update the template"),
			},
		},
	})
}

func TestTemplateImport_ImportTemplateError(t *testing.T) {
	if skipTest() {
		t.Skip(SkipTestMsg)
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
	if skipTest() {
		t.Skip(SkipTestMsg)
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

func checkResourceCreation(t *testing.T, p tfsdk.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		templateAPI := "/api/TemplateService/Templates?$filter=Name eq '%s'"
		template1URL := fmt.Sprintf(templateAPI, TemplateName1)
		fmt.Println("Ome template 1 url", template1URL)
		template2URL := fmt.Sprintf(templateAPI, TemplateName2)
		provider := p.(*provider)
		omeClient, err := clients.NewClient(*provider.clientOpt)
		if err != nil {
			return fmt.Errorf("Unable to create client %s", err.Error())
		}

		_, err = omeClient.CreateSession()
		if err != nil {
			return fmt.Errorf("Error creating client session %s", err.Error())
		}
		response, err := omeClient.Get(template1URL, nil, nil)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		b, _ := omeClient.GetBodyData(response.Body)

		omeTemplates1 := models.OMETemplates{}
		err = omeClient.JSONUnMarshal(b, &omeTemplates1)
		if err != nil {
			fmt.Printf("Unable to create client %s", err.Error())
		}

		omeTemplate1 := omeTemplates1.Value[0]
		fmt.Println("Ome template 1", omeTemplate1)
		assert.Equal(t, TemplateName1, omeTemplate1.Name)
		assert.Equal(t, int64(2), omeTemplate1.ViewTypeID)

		devID, _ := strconv.ParseInt(DeviceID1, 10, 64)
		assert.Equal(t, devID, omeTemplate1.SourceDeviceID)
		attrURL := fmt.Sprintf("%s(%d)/AttributeDetails", clients.TemplateAPI, omeTemplate1.ID)
		response, _ = omeClient.Get(attrURL, nil, nil)
		b, _ = omeClient.GetBodyData(response.Body)
		omeTemplateAttrGroups := models.OMETemplateAttrGroups{}
		omeClient.JSONUnMarshal(b, &omeTemplateAttrGroups)

		assert.Equal(t, 2, len(omeTemplateAttrGroups.AttributeGroups))
		assert.Equal(t, "iDRAC", omeTemplateAttrGroups.AttributeGroups[0].DisplayName)

		response, err = omeClient.Get(template2URL, nil, nil)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		b, _ = omeClient.GetBodyData(response.Body)
		omeTemplates2 := models.OMETemplates{}
		omeClient.JSONUnMarshal(b, &omeTemplates2)

		omeTemplate2 := omeTemplates2.Value[0]
		assert.Equal(t, TemplateName2, omeTemplate2.Name)
		// assert.Equal(t, int64(1), omeTemplate2.ViewTypeID)
		devID, _ = strconv.ParseInt(DeviceID2, 10, 64)
		assert.Equal(t, devID, omeTemplate2.SourceDeviceID)

		_, err = omeClient.RemoveSession()
		if err != nil {
			fmt.Println("Error on removing sessions ", err.Error())
		}
		return nil
	}
}

func checkResourceUpdate(t *testing.T, p tfsdk.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		templateAPI := "/api/TemplateService/Templates?$filter=Name eq '%s'"
		template1URL := fmt.Sprintf(templateAPI, TemplateName1)
		template2URL := fmt.Sprintf(templateAPI, TemplateNameUpdate2)
		provider := p.(*provider)
		omeClient, err := clients.NewClient(*provider.clientOpt)
		if err != nil {
			return fmt.Errorf("Unable to create client %s", err.Error())
		}

		_, err = omeClient.CreateSession()
		if err != nil {
			return fmt.Errorf("Error creating client session %s", err.Error())
		}
		response, err := omeClient.Get(template1URL, nil, nil)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		b, _ := omeClient.GetBodyData(response.Body)

		omeTemplates1 := models.OMETemplates{}
		err = omeClient.JSONUnMarshal(b, &omeTemplates1)
		if err != nil {
			fmt.Printf("Unable to create client %s", err.Error())
		}

		omeTemplate1 := omeTemplates1.Value[0]
		assert.Equal(t, TemplateName1, omeTemplate1.Name)
		assert.Equal(t, "This is a test template", omeTemplate1.Description)

		attrURL := fmt.Sprintf("%s(%d)/AttributeDetails", clients.TemplateAPI, omeTemplate1.ID)
		response, _ = omeClient.Get(attrURL, nil, nil)
		b, _ = omeClient.GetBodyData(response.Body)
		omeTemplateAttrGroups := models.OMETemplateAttrGroups{}
		omeClient.JSONUnMarshal(b, &omeTemplateAttrGroups)

		assert.Equal(t, 2, len(omeTemplateAttrGroups.AttributeGroups))
		assert.Equal(t, "iDRAC", omeTemplateAttrGroups.AttributeGroups[0].DisplayName)
		assert.Contains(t, string(b), "\"Value\":\"IST\"")
		assert.Nil(t, err)

		vlanAttrs, err := omeClient.GetSchemaVlanData(omeTemplate1.ID)
		if err != nil {
			fmt.Printf("Unable to get vlan attributes %s", err.Error())
		}
		assert.Equal(t, vlanAttrs.BondingTechnology, "NoTeaming")
		assert.True(t, len(vlanAttrs.OMEVlanAttributes) > 0)

		response, err = omeClient.Get(template2URL, nil, nil)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		b, _ = omeClient.GetBodyData(response.Body)

		omeTemplates2 := models.OMETemplates{}
		err = omeClient.JSONUnMarshal(b, &omeTemplates2)
		if err != nil {
			fmt.Printf("Unable to create client %s", err.Error())
		}

		omeTemplate2 := omeTemplates2.Value[0]
		assert.Equal(t, omeTemplate2.Name, TemplateNameUpdate2)
		assert.Equal(t, "This is sample description", omeTemplate2.Description)
		return nil
	}
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
		fqdds = "iDRAC,NIC"
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
		fqdds = "iDRAC,NIC"
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
		refdevice_servicetag = "` + DeviceSvcTag2 + `"
		identity_pool_name = "IO1"
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
	fqdds = "iDRAC,NIC"
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

var testAccCreateTemplateWithExistingName = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-acceptance-test-1" {
	name = "` + TemplateName1 + `"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
	fqdds = "iDRAC"
}

resource "ome_template" "terraform-acceptance-test-2" {
	name = "` + TemplateName2 + `"
	refdevice_id = ` + DeviceID2 + `
	job_retry_count  = 15
	sleep_interval = 60
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
	fqdds = "iDRAC"
}

resource "ome_template" "terraform-acceptance-test-2" {
	name = "` + TemplateName1 + `"
	refdevice_servicetag = "` + DeviceSvcTag2 + `"
	job_retry_count  = 15
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

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "` + ReferenceDeploymentTemplateNameForClone + `"
}
`

var testAccCloneTemplateSuccessForDeploymentToCompliance = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "` + ReferenceDeploymentTemplateNameForClone + `"
	view_type = "Compliance"
}
`

var testAccCloneTemplateSuccessForComplianceToCompliance = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
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

var testAccCloneTemplateFailureForRefTemplateAndDeviceID = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "invalid-template-name"
	refdevice_servicetag = "` + DeviceSvcTag2 + `"
}
`

var testAccCloneTemplateFailureForRefTemplateAndDeviceServiceTag = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}

resource "ome_template" "terraform-clone-template-test" {
	name = "clone-template-test"
	reftemplate_name = "invalid-template-name"
	refdevice_servicetag = "` + DeviceSvcTag1 + `"
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
	reftemplate_name = "test-compliance-template"
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
