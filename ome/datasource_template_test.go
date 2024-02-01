/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
