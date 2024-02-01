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

func TestDataSource_ReadCert(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testGetCert + certDataOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("fetched", "true"),
				),
			},
		},
	})
}

var certDataOut = `
output "fetched" {
	value = alltrue([
		length(data.ome_application_certificate.cert.valid_to) > 0,
		length(data.ome_application_certificate.cert.valid_from) > 0,
		length(data.ome_application_certificate.cert.issued_to.distinguished_name) > 0,
		length(data.ome_application_certificate.cert.issued_by.distinguished_name) > 0,
	])
}
`

var testGetCert = testProvider + `
data "ome_application_certificate" "cert" {
}
`
