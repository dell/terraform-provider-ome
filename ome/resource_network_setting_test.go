package ome

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ============================================= Time Setting Test ==============================================
func TestNetworkSettingTime(t *testing.T) {
	testAccCreateNetworkTimeSuccess := testProvider + `
	resource "ome_network_setting" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  system_time = "2023-08-18 07:50:08.387"
		}
	  }
	`
	testAccUpdateNetworkTimeSuccess := testProvider + `
	resource "ome_network_setting" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  enable_ntp = true 
		  primary_ntp_address = "` + DeviceIP1 + `"
		  secondary_ntp_address1 = "` + DeviceIP2 + `"
		  secondary_ntp_address2 = "` + DeviceIP3 + `"
		}
	  }
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateNetworkTimeSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.its_ome_time", "time_setting.time_zone", "TZ_ID_65"),
					resource.TestCheckResourceAttr("ome_network_setting.its_ome_time", "time_setting.system_time", "2023-08-18 07:50:08.387"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccUpdateNetworkTimeSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.its_ome_time", "time_setting.time_zone", "TZ_ID_65"),
					resource.TestCheckResourceAttr("ome_network_setting.its_ome_time", "time_setting.enable_ntp", "true"),
					resource.TestCheckResourceAttr("ome_network_setting.its_ome_time", "time_setting.primary_ntp_address", DeviceIP1),
					resource.TestCheckResourceAttr("ome_network_setting.its_ome_time", "time_setting.secondary_ntp_address1", DeviceIP2),
					resource.TestCheckResourceAttr("ome_network_setting.its_ome_time", "time_setting.secondary_ntp_address2", DeviceIP3),
				),
			},
		},
	})
}

func TestNetworkSettingTimeInvalidConfig(t *testing.T) {
	testAccNetworkTimeInvalid := testProvider + `
	resource "ome_network_setting" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  enable_ntp = true
		  system_time = "2023-08-19 07:50:08.387"
		}
	}
	`
	testAccNetworkTimeInvalid1 := testProvider + `
	resource "ome_network_setting" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  enable_ntp = true
		}
	}
	`

	testAccNetworkTimeInvalid2 := testProvider + `
	resource "ome_network_setting" "its_ome_time" {
		time_setting = {
		  enable_ntp = true
		}
	}
	`
	testAccNetworkTimeInvalid3 := testProvider + `
	resource "ome_network_setting" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		primary_ntp_address = "10.x.x.1"
		}
	}
	`

	testAccNetworkTimeInvalid4 := testProvider + `
	resource "ome_network_setting" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		}
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkTimeInvalid,
				ExpectError: regexp.MustCompile(`.*please validate that the system_time is unset when enable_ntp is active*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid1,
				ExpectError: regexp.MustCompile(`.*please validate that the primary_ntp_address is set when enable_ntp is active*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid2,
				ExpectError: regexp.MustCompile(`.*please validate that the time_zone is set*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid3,
				ExpectError: regexp.MustCompile(`.*please validate that primary_ntp_address, secondary_ntp_address1*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid4,
				ExpectError: regexp.MustCompile(`.*please validate that the system_time is set*.`),
			},
		},
	})
}

// ============================================= Session Setting Test ===========================================

func TestNetworkSettingSession(t *testing.T) {
	testAccCreateNetworkSessionSuccess := testProvider + `
	resource "ome_network_setting" "code_ome" {
		session_setting = {
		  enable_universal_timeout = true 
		  universal_timeout = 20
		  api_session = 10 
		  gui_session = 11
		}
	  }
	`
	testAccUpdateNetworkSessionSuccess := testProvider + `
	resource "ome_network_setting" "code_ome" {
		session_setting = {
		  enable_universal_timeout = false
		  api_session = 15
		  api_timeout = 40
		  gui_session = 20
		  gui_timeout = 40
		}
	  }
	`

	testAccUpdateNetworkSessionSuccess1 := testProvider + `
	resource "ome_network_setting" "code_ome" {
		session_setting = {
		  enable_universal_timeout = true
		  universal_timeout = 30
		  api_session = 10
		  gui_session = 10
		}
	  }
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateNetworkSessionSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.enable_universal_timeout", "true"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.universal_timeout", "20"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.api_session", "10"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.gui_session", "11"),
				),
			},
			{
				Config: testAccUpdateNetworkSessionSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.enable_universal_timeout", "false"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.api_session", "15"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.api_timeout", "40"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.gui_session", "20"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.gui_timeout", "40"),
				),
			},
			{
				Config: testAccUpdateNetworkSessionSuccess1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.enable_universal_timeout", "true"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.universal_timeout", "30"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.api_session", "10"),
					resource.TestCheckResourceAttr("ome_network_setting.code_ome", "session_setting.gui_session", "10"),
				),
			},
		},
	})
}

func TestNetworkSettingSessionInValidConfig(t *testing.T) {
	testAccNetworkSessionInvalid := testProvider + `
	resource "ome_network_setting" "code_1" {
		session_setting = {
		  enable_universal_timeout = true
		}
	  }
	`

	testAccNetworkSessionInvalid1 := testProvider + `
	resource "ome_network_setting" "code_1" {
		session_setting = {
		  enable_universal_timeout = true
		  universal_timeout = 10
		  api_timeout = 20
		}
	  }
	`

	testAccNetworkSessionInvalid2 := testProvider + `
	resource "ome_network_setting" "code_1" {
		session_setting = {
			universal_timeout = 10
		}
	  }
	`

	testAccNetworkSessionInvalid3 := testProvider + `
	resource "ome_network_setting" "code_1" {
		session_setting = {
			ssh_timeout = 10
		}
	}
	`

	testAccNetworkSessionInvalid4 := testProvider + `
	resource "ome_network_setting" "code_1" {
		session_setting = {
			serial_timeout = 10
		}
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkSessionInvalid,
				ExpectError: regexp.MustCompile(`.*please ensure universal_timeout is set*.`),
			},
			{
				Config:      testAccNetworkSessionInvalid1,
				ExpectError: regexp.MustCompile(`.*please validate that the configuration for api_timeout*.`),
			},
			{
				Config:      testAccNetworkSessionInvalid2,
				ExpectError: regexp.MustCompile(`.*please ensure universal_timeout is unset*.`),
			},
			{
				Config:      testAccNetworkSessionInvalid3,
				ExpectError: regexp.MustCompile(`.*please verify that the SSH Session is unset*.`),
			},
			{
				Config:      testAccNetworkSessionInvalid4,
				ExpectError: regexp.MustCompile(`.*please verify that the Serial Session is unset*.`),
			},
		},
	})
}

// ============================================== Proxy Setting Test ============================================

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

func TestNetworkSettingProxyNil(t *testing.T) {
	testAccNetworkProxyCreateNil := testProvider + `
	resource "ome_network_setting" "code_3" {
	}
	`
	testAccNetworkProxyCreateUpdateNil1 := testProvider + `
	resource "ome_network_setting" "code_4" {
		proxy_setting = {
			enable_proxy = false
		}
	}
	`
	testAccNetworkProxyCreateUpdateNil2 := testProvider + `
	resource "ome_network_setting" "code_4" {
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkProxyCreateNil,
			},
			{
				Config: testAccNetworkProxyCreateUpdateNil1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_network_setting.code_4", "proxy_setting.enable_proxy", "false"),
				),
			},
			{
				Config: testAccNetworkProxyCreateUpdateNil2,
			},
		},
	})
}
