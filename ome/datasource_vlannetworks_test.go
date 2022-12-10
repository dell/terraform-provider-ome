package ome

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSource_ReadVlanNetworks(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testVlanNetworks,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ome_vlannetworks_info.vlans", "vlan_networks.#", "3")),
			},
		},
	})
}

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
