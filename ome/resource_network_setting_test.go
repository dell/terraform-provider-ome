package ome

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestNetworkSettingProxy(t *testing.T) {
	testAccCreateNetworkProxySuccess := testProvider + `
	resource "ome_network_setting" "code_1" {
		proxy_setting = {
		  enable_proxy = true
		  ip_address = "` + DeviceIP1 + `"
		  proxy_port = 443
		  enable_authentication = false
		}
	  }
	`

	testAccUpdateNetworkProxySuccess := testProvider + `
	resource "ome_network_setting" "code_1" {
		proxy_setting = {
		  enable_proxy = true
		  ip_address = "` + DeviceIP1 + `"
		  proxy_port = 443
		  enable_authentication = true
		  username = "root" 
		  password = "root"
		}
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateNetworkProxySuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_1", "proxy_setting.enable_proxy", "true"),
				),
			},
			{
				Config: testAccUpdateNetworkProxySuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_1", "proxy_setting.username", "root"),
				),
			},
		},
	})
}

func TestNetworkSettingProxyIsInfraChangeDetected(t *testing.T) {
	testAccCreateNetworkProxy := testProvider + `
	resource "ome_network_setting" "code_2" {
		proxy_setting = {
		  enable_proxy = true
		  ip_address = "` + DeviceIP2 + `"
		  proxy_port = 446
		  enable_authentication = false
		}
	  }
	`
	testAccUpdateNetworkProxy := testProvider + `
	resource "ome_network_setting" "code_2" {
		proxy_setting = {
		  enable_proxy = false
		}
	  }
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateNetworkProxy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_2", "proxy_setting.enable_proxy", "true"),
				),
			},
			{
				Config: testAccUpdateNetworkProxy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_2", "proxy_setting.enable_proxy", "false"),
				),
			},
		},
	})
}

func TestNetworkSettingProxyInValidConfig(t *testing.T) {
	testAccNetworkProxyInvalid := testProvider + `
	resource "ome_network_setting" "code_1" {
		proxy_setting = {
			enable_proxy = false
			ip_address = "` + DeviceIP1 + `"
		}
	}
	`
	testAccNetworkProxyInvalid1 := testProvider + `
	resource "ome_network_setting" "code_1" {
		proxy_setting = {
			enable_proxy = true
		}
	}
	`
	testAccNetworkProxyInvalid2 := testProvider + `
	resource "ome_network_setting" "code_1" {
		proxy_setting = {
			enable_proxy = true
			ip_address = "` + DeviceIP1 + `"
		  	proxy_port = 443
		  	enable_authentication = true
		}
	}
	`

	testAccNetworkProxyInvalid3 := testProvider + `
	resource "ome_network_setting" "code_1" {
		proxy_setting = {
			enable_proxy = true
			ip_address = "` + DeviceIP1 + `"
		  	proxy_port = 443
		  	enable_authentication = false
			username = "root"
		}
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkProxyInvalid,
				ExpectError: regexp.MustCompile(`.*please ensure enable proxy should be set to true*.`),
			},
			{
				Config:      testAccNetworkProxyInvalid1,
				ExpectError: regexp.MustCompile(`.*please ensure that you set both the IP address and port*.`),
			},
			{
				Config:      testAccNetworkProxyInvalid2,
				ExpectError: regexp.MustCompile(`.*please ensure that you set both the username and password*.`),
			},
			{
				Config:      testAccNetworkProxyInvalid3,
				ExpectError: regexp.MustCompile(`.*please ensure enable authentication should be set to true*.`),
			},
		},
	})
}
