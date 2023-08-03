package ome

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDiscoveryOne(t *testing.T) {
	testAccProvider := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}
	`

	testAccCreateDiscoverySuccess := testAccProvider + `
	resource "ome_discovery" "code_1" {
		name = "test-create"
		schedule = "RunLater"
		cron = "0 * */10 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = [ "` + DeviceIP1 + `", "` + DeviceIP2 + `"]
		  redfish = {
		   username = "root"
		   password = "calvin" 
		  }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "root"
			password = "calvin"
		  }
		}]
	  }
	`
	testAccUpdateDiscoverySuccess := testAccProvider + `
	resource "ome_discovery" "code_1" {
		name = "test-update"
		schedule = "RunLater"
		cron = "0 * */12 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = [ "` + DeviceIP1 + `"]
		  redfish = {
		   username = "root"
		   password = "calvin" 
		  }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "root"
			password = "calvin"
		  }
		}]
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDiscoverySuccess,
			},
			{
				Config: testAccUpdateDiscoverySuccess,
			},
		},
	})

}

func TestDiscoveryTwo(t *testing.T) {
	testAccProvider := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}
	`

	testAccCreateDiscoveryDebug := testAccProvider + `
	resource "ome_discovery" "code_1" {
		name = "shiva-ganga"
		schedule = "RunLater"
		cron = "0 * */10 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = [ "`+ DeviceIP1 + `","` + DeviceIP2+ `","` + DeviceIP3 + `"]
		  # redfish = {
		  #  username = "root"
		  #  password = "calvin" 
		  # }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "root"
			password = "calvin"
		  }
		}]
	  }
	`
	testAccUpdateDiscoveryDebug := testAccProvider + `
	resource "ome_discovery" "code_1" {
		name = "kashi-ganga"
		schedule = "RunLater"
		cron = "0 * */10 * * ? *"
		discovery_config_targets = [{
		  device_type = [ "SERVER" ]
		  network_address_detail = ["` + DeviceIP3 + `"]
		  redfish = {
		   username = "root"
		   password = "calvin" 
		  }
		  snmp = {
			community = "public"
		  }
		  ssh = {
			username = "root"
			password = "calvin"
		  }
		}]
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateDiscoveryDebug,
			},
			{
				Config: testAccUpdateDiscoveryDebug,
			},
		},
	})

}

