package ome

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSource_ReadTemplate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	temps := initTemplates(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: justProvider + temps.templateSvcTag1Full,
			},
			{
				Config: testReadTemplate + temps.templateSvcTag1Full,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_template_info.template", "name", TestRefTemplateName),
					resource.TestCheckResourceAttr("data.ome_template_info.template", "view_type_id", "1"),
					resource.TestCheckResourceAttrPair("data.ome_template_info.template", "attributes", "ome_template.terraform-acceptance-test-1", "attributes"),
					resource.TestCheckResourceAttr("data.ome_template_info.template", "description", "Imported from a file."),
				),
			},
			{
				Config: testReadInvalidTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_template_info.template", "name", "InvalidTemplate"),
					resource.TestCheckNoResourceAttr("data.ome_template_info.template", "view_type_id"),
					resource.TestCheckNoResourceAttr("data.ome_template_info.template", "refdevice_id"),
					resource.TestCheckNoResourceAttr("data.ome_template_info.template", "description"),
				),
			},
		},
	})
}

var testReadTemplate = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_template_info" "template" {
		id = "0"
		name = "` + TestRefTemplateName + `"
		depends_on = ["ome_template.terraform-acceptance-test-1"]
	}
`

var testReadInvalidTemplate = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_template_info" "template" {
		id = "0"
		name = "` + "InvalidTemplate" + `"
	}
`
