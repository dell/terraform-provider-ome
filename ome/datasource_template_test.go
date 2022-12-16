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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testReadTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_template_info.template", "name", TestAccTemplateName),
					resource.TestCheckResourceAttr("data.ome_template_info.template", "view_type_id", "2"),
					resource.TestCheckResourceAttr("data.ome_template_info.template", "refdevice_id", "12152"),
					resource.TestCheckResourceAttr("data.ome_template_info.template", "description", ""),
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

	resource "ome_template" "terraform-acceptance-test-1" {
		name = "` + TestRefTemplateName + `"
		refdevice_servicetag = "` + DeviceSvcTag1 + `"
		fqdds = "EventFilters"
		view_type = "Compliance"
		job_retry_count = 20
		sleep_interval = 30
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
