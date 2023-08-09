package ome

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
