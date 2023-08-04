package ome

import (
	"regexp"
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
		  network_address_detail = [ "` + DeviceIP1 + `","` + DeviceIP2 + `","` + DeviceIP3 + `"]
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

func TestDiscoveryThree(t *testing.T) {
	testAccProvider := `
	provider "ome" {
		username = "` + omeUserName + `"
		password = "` + omePassword + `"
		host = "` + omeHost + `"
		skipssl = true
	}
	`

	invalidDiscoveryConfigOne := testAccProvider + `
	resource "ome_discovery" "code_3" {
		name = "invalid-config"
		schedule = "RunNow"
		cron = "abc"
		discovery_config_targets = [ {
		  device_type = ["abc"]
		  network_address_detail = ["x.x.x.x"]
		  } 
		]
	  }
	`

	invalidDiscoveryConfigtwo := testAccProvider + `
	resource "ome_discovery" "code_4" {
		name = "invalid-config"
		discovery_config_targets = []
	  }
	`

	invalidDiscoveryConfigThree := testAccProvider + `
	resource "ome_discovery" "code_4" {
		name = "invalid-config"
		discovery_config_targets = [{
		  device_type = []
		  network_address_detail = []
		}]
	  }
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      invalidDiscoveryConfigOne,
				ExpectError: regexp.MustCompile(`.*The device type list should contain the following values*.`),
			},
			{
				Config:      invalidDiscoveryConfigOne,
				ExpectError: regexp.MustCompile(`.*With Schedule as RunNow, CRON can't be set*.`),
			},
			{
				Config:      invalidDiscoveryConfigOne,
				ExpectError: regexp.MustCompile(`.*Atleast one of protocol should be configured for the discovery targets*.`),
			},
			{
				Config: invalidDiscoveryConfigtwo,
				ExpectError: regexp.MustCompile(`.*Define at least one discovery configuration target in the list.*.`),
			},
			{
				Config: invalidDiscoveryConfigThree,
				ExpectError: regexp.MustCompile(`.*Atleast one of device type should be configured*.`),
			},
			{
				Config: invalidDiscoveryConfigThree,
				ExpectError: regexp.MustCompile(`.*Atleast one of network address detail should be configured*.`),
			},
		},
	})
}
