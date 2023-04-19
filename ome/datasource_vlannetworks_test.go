package ome

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSource_ReadVlanNetworks(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testVlanNetworksWrongCreds,
				ExpectError: regexp.MustCompile(".*invalid credentials.*"),
			},
			{
				Config: testVlanNetworks,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_vlannetworks_info.vlans", "vlan_networks.#", "3")),
			},
		},
	})
}

var testVlanNetworksWrongCreds = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "invalid"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_vlannetworks_info" "vlans" {
	}
`
var testVlanNetworks = `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}

	data "ome_vlannetworks_info" "vlans" {
	}
`
