package ome

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ============================================= Adapter Setting Test ===========================================
func TestNetworkSettingAdapter(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	if os.Getenv("TF_ACC_NWAD") == "" {
		t.Skip("Skipping as TF_ACC_NWAD is not set")
	}

	testAccCreateAdapterSetting := testProvider + `
	resource "ome_appliance_network" "net1" {
		adapter_setting = {
		  enable_nic     = true
		  interface_name = "ens160"
		  reboot_delay   = 0
		  management_vlan = {
			  enable_vlan = false
		  }
		  ipv6_configuration = {
			enable_ipv6                   = false
			enable_auto_configuration = false
		  }
		dns_configuration = {
		  register_with_dns = false
		}
		}
	  }
	`
	testAccUpdateAdapterSetting := testProvider + `
	resource "ome_appliance_network" "net1" {
		adapter_setting = {
		  enable_nic     = true
		  interface_name = "ens160"
		  reboot_delay   = 0
		  management_vlan = {
			  enable_vlan = false
		  }
		  ipv6_configuration = {
			enable_ipv6                   = true
			enable_auto_configuration = true
		  }
		dns_configuration = {
		  register_with_dns = false
		}
		}
	  }
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateAdapterSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_appliance_network.net1", "adapter_setting.ipv6_configuration.enable_ipv6", "false"),
				),
			},
			{
				Config: testAccUpdateAdapterSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_appliance_network.net1", "adapter_setting.ipv6_configuration.enable_ipv6", "true"),
				),
			},
		},
	})
}

func TestNetworkSettingAdapterInvalidConfig(t *testing.T) {
	testAccNetworkAdapterInvalid1 := testProvider + `
	resource "ome_appliance_network" "invalid1" {
		adapter_setting = {
		  interface_name = "invalid"
		  enable_nic = true
		  ipv4_configuration = {
			enable_ipv4        = true
			enable_dhcp        = true
			static_ip_address  = "0.0.0.0"
			static_subnet_mask = "1.1.1.1"
			static_gateway     = "2.2.2.2"
		  }
		}
	  }
	`

	testAccNetworkAdapterInvalid2 := testProvider + `
	resource "ome_appliance_network" "invalid2" {
		adapter_setting = {
		  enable_nic = true
		  interface_name = "invalid"
		  ipv4_configuration = {
			enable_ipv4 = true
			use_dhcp_for_dns_server_names = true
			static_preferred_dns_server    = "3.3.3.3"
			static_alternate_dns_server    = "4.4.4.4"
		  }
		}
	  }
	`
	testAccNetworkAdapterInvalid3 := testProvider + `
	resource "ome_appliance_network" "invalid3" {
		adapter_setting = {
		  enable_nic = true
		  interface_name = "invalid"
		  ipv6_configuration = {
			enable_ipv6 = true
			enable_auto_configuration = true
			static_ip_address = "0.0.0.0"
			static_prefix_length = 0
			static_gateway = "1.1.1.1"
		  }
		}
	  }
	`

	testAccNetworkAdapterInvalid4 := testProvider + `
	resource "ome_appliance_network" "invalid4" {
		adapter_setting = {
		  interface_name = "invalid"
		  enable_nic = true
		  ipv6_configuration = {
			enable_ipv6 = true
			use_dhcp_for_dns_server_names = true
			static_preferred_dns_server    = "3.3.3.3"
			static_alternate_dns_server    = "4.4.4.4"
		  }
		}
	  }
	`

	testAccNetworkAdapterInvalid5 := testProvider + `
	resource "ome_appliance_network" "invalid5" {
		adapter_setting = {
		  interface_name = "invalid"
		  enable_nic = true
		  management_vlan = {
			enable_vlan = false
			id = 1
		  }
		}
	  }
	`

	testAccNetworkAdapterInvalid6 := testProvider + `
	resource "ome_appliance_network" "invalid6" {
		adapter_setting = {
		  interface_name = "invalid"
		  enable_nic = true
		  dns_configuration = {
			register_with_dns              = false
			dns_name                       = "err"
			use_dhcp_for_dns_server_names = true
			dns_domain_name                = "err"
		  }
		}
	  }
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkAdapterInvalid1,
				ExpectError: regexp.MustCompile(`.*static_ip_address / static_subnet_mask / static_gateway should not be set*.`),
			},
			{
				Config:      testAccNetworkAdapterInvalid2,
				ExpectError: regexp.MustCompile(`.*static_ip_address / static_subnet_mask / static_gateway are required*.`),
			},
			{
				Config:      testAccNetworkAdapterInvalid3,
				ExpectError: regexp.MustCompile(`.*static_ip_address / static_prefix_length / static_gateway should not be set*.`),
			},
			{
				Config:      testAccNetworkAdapterInvalid4,
				ExpectError: regexp.MustCompile(`.*static_ip_address / static_prefix_length / static_gateway are required*.`),
			},
			{
				Config:      testAccNetworkAdapterInvalid5,
				ExpectError: regexp.MustCompile(`.*please validate enable_vlan is true*.`),
			},
			{
				Config:      testAccNetworkAdapterInvalid6,
				ExpectError: regexp.MustCompile(`.*please validate register_with_dn is true*.`),
			},
		},
	})
}

// ============================================= Time Setting Test ==============================================
func TestNetworkSettingTime(t *testing.T) {
	testAccCreateNetworkTimeSuccess := testProvider + `
	resource "ome_appliance_network" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  system_time = "2023-08-18 07:50:08.387"
		}
	  }
	`
	testAccUpdateNetworkTimeSuccess := testProvider + `
	resource "ome_appliance_network" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  enable_ntp = true
		  primary_ntp_address = "` + DeviceIP1 + `"
		  secondary_ntp_address1 = "` + DeviceIP2 + `"
		  secondary_ntp_address2 = "` + DeviceIPExt + `"
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
					resource.TestCheckResourceAttr("ome_appliance_network.its_ome_time", "time_setting.time_zone", "TZ_ID_65"),
					resource.TestCheckResourceAttr("ome_appliance_network.its_ome_time", "time_setting.system_time", "2023-08-18 07:50:08.387"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccUpdateNetworkTimeSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_appliance_network.its_ome_time", "time_setting.time_zone", "TZ_ID_65"),
					resource.TestCheckResourceAttr("ome_appliance_network.its_ome_time", "time_setting.enable_ntp", "true"),
					resource.TestCheckResourceAttr("ome_appliance_network.its_ome_time", "time_setting.primary_ntp_address", DeviceIP1),
					resource.TestCheckResourceAttr("ome_appliance_network.its_ome_time", "time_setting.secondary_ntp_address1", DeviceIP2),
					resource.TestCheckResourceAttr("ome_appliance_network.its_ome_time", "time_setting.secondary_ntp_address2", DeviceIPExt),
				),
			},
		},
	})
}

func TestNetworkSettingTimeInvalidConfig(t *testing.T) {
	testAccNetworkTimeInvalid := testProvider + `
	resource "ome_appliance_network" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  enable_ntp = true
		  system_time = "2023-08-19 07:50:08.387"
		}
	}
	`
	testAccNetworkTimeInvalid1 := testProvider + `
	resource "ome_appliance_network" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		  enable_ntp = true
		}
	}
	`

	testAccNetworkTimeInvalid2 := testProvider + `
	resource "ome_appliance_network" "its_ome_time" {
		time_setting = {
		  enable_ntp = true
		}
	}
	`
	testAccNetworkTimeInvalid3 := testProvider + `
	resource "ome_appliance_network" "its_ome_time" {
		time_setting = {
		  time_zone = "TZ_ID_65"
		primary_ntp_address = "10.x.x.1"
		}
	}
	`

	testAccNetworkTimeInvalid4 := testProvider + `
	resource "ome_appliance_network" "its_ome_time" {
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
				ExpectError: regexp.MustCompile(`.*system_time should not be set when enable_ntp is active*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid1,
				ExpectError: regexp.MustCompile(`.*primary_ntp_address should be set when enable_ntp is active*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid2,
				ExpectError: regexp.MustCompile(`.*Inappropriate value for attribute "time_setting": attribute "time_zone"*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid3,
				ExpectError: regexp.MustCompile(`.*primary_ntp_address, secondary_ntp_address1 and secondary_ntp_address2*.`),
			},
			{
				Config:      testAccNetworkTimeInvalid4,
				ExpectError: regexp.MustCompile(`.*system_time should be set when enable_ntp is disable*.`),
			},
		},
	})
}

// ============================================= Session Setting Test ===========================================

func TestNetworkSettingSession(t *testing.T) {
	testAccCreateNetworkSessionSuccess := testProvider + `
	resource "ome_appliance_network" "code_ome" {
		session_setting = {
		  enable_universal_timeout = true
		  universal_timeout = 20
		  api_session = 10
		  gui_session = 11
		}
	  }
	`
	testAccUpdateNetworkSessionSuccess := testProvider + `
	resource "ome_appliance_network" "code_ome" {
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
	resource "ome_appliance_network" "code_ome" {
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
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.enable_universal_timeout", "true"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.universal_timeout", "20"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.api_session", "10"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.gui_session", "11"),
				),
			},
			{
				Config: testAccUpdateNetworkSessionSuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.enable_universal_timeout", "false"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.api_session", "15"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.api_timeout", "40"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.gui_session", "20"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.gui_timeout", "40"),
				),
			},
			{
				Config: testAccUpdateNetworkSessionSuccess1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.enable_universal_timeout", "true"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.universal_timeout", "30"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.api_session", "10"),
					resource.TestCheckResourceAttr("ome_appliance_network.code_ome", "session_setting.gui_session", "10"),
				),
			},
		},
	})
}

func TestNetworkSettingSessionInValidConfig(t *testing.T) {
	testAccNetworkSessionInvalid := testProvider + `
	resource "ome_appliance_network" "code_1" {
		session_setting = {
		  enable_universal_timeout = true
		}
	  }
	`

	testAccNetworkSessionInvalid1 := testProvider + `
	resource "ome_appliance_network" "code_1" {
		session_setting = {
		  enable_universal_timeout = true
		  universal_timeout = 10
		  api_timeout = 20
		}
	  }
	`

	testAccNetworkSessionInvalid2 := testProvider + `
	resource "ome_appliance_network" "code_1" {
		session_setting = {
			universal_timeout = 10
		}
	  }
	`

	testAccNetworkSessionInvalid3 := testProvider + `
	resource "ome_appliance_network" "code_1" {
		session_setting = {
			ssh_timeout = 10
		}
	}
	`

	testAccNetworkSessionInvalid4 := testProvider + `
	resource "ome_appliance_network" "code_1" {
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
				ExpectError: regexp.MustCompile(`.*universal_timeout should be set*.`),
			},
			{
				Config:      testAccNetworkSessionInvalid1,
				ExpectError: regexp.MustCompile(`.*api_timeout, gui_timeout, ssh_timeout and serial_timeout*.`),
			},
			{
				Config:      testAccNetworkSessionInvalid2,
				ExpectError: regexp.MustCompile(`.*universal_timeout should not be set*.`),
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
	resource "ome_appliance_network" "code_1" {
		proxy_setting = {
		  enable_proxy = true
		  ip_address = "` + DeviceIP1 + `"
		  proxy_port = 443
		  enable_authentication = false
		}
	  }
	`

	testAccUpdateNetworkProxySuccess := testProvider + `
	resource "ome_appliance_network" "code_1" {
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
					resource.TestCheckResourceAttr("ome_appliance_network.code_1", "proxy_setting.enable_proxy", "true"),
				),
			},
			{
				Config: testAccUpdateNetworkProxySuccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_appliance_network.code_1", "proxy_setting.username", "root"),
				),
			},
		},
	})
}

func TestNetworkSettingProxyIsInfraChangeDetected(t *testing.T) {
	testAccCreateNetworkProxy := testProvider + `
	resource "ome_appliance_network" "code_2" {
		proxy_setting = {
		  enable_proxy = true
		  ip_address = "` + DeviceIP2 + `"
		  proxy_port = 446
		  enable_authentication = false
		}
	  }
	`
	testAccUpdateNetworkProxy := testProvider + `
	resource "ome_appliance_network" "code_2" {
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
					resource.TestCheckResourceAttr("ome_appliance_network.code_2", "proxy_setting.enable_proxy", "true"),
				),
			},
			{
				Config: testAccUpdateNetworkProxy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ome_appliance_network.code_2", "proxy_setting.enable_proxy", "false"),
				),
			},
		},
	})
}

func TestNetworkSettingProxyInValidConfig(t *testing.T) {
	testAccNetworkProxyInvalid := testProvider + `
	resource "ome_appliance_network" "code_1" {
		proxy_setting = {
			enable_proxy = false
			ip_address = "` + DeviceIP1 + `"
		}
	}
	`
	testAccNetworkProxyInvalid1 := testProvider + `
	resource "ome_appliance_network" "code_1" {
		proxy_setting = {
			enable_proxy = true
		}
	}
	`
	testAccNetworkProxyInvalid2 := testProvider + `
	resource "ome_appliance_network" "code_1" {
		proxy_setting = {
			enable_proxy = true
			ip_address = "` + DeviceIP1 + `"
		  	proxy_port = 443
		  	enable_authentication = true
		}
	}
	`

	testAccNetworkProxyInvalid3 := testProvider + `
	resource "ome_appliance_network" "code_1" {
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
				ExpectError: regexp.MustCompile(`.*enable proxy should be set to true*.`),
			},
			{
				Config:      testAccNetworkProxyInvalid1,
				ExpectError: regexp.MustCompile(`.*both IP address and port are required*.`),
			},
			{
				Config:      testAccNetworkProxyInvalid2,
				ExpectError: regexp.MustCompile(`.*both username and password are required*.`),
			},
			{
				Config:      testAccNetworkProxyInvalid3,
				ExpectError: regexp.MustCompile(`.*enable authentication should be set to true*.`),
			},
		},
	})
}

func TestNetworkSettingProxyNil(t *testing.T) {
	testAccNetworkProxyCreateNil := testProvider + `
	resource "ome_appliance_network" "code_3" {
	}
	`
	testAccNetworkProxyCreateUpdateNil1 := testProvider + `
	resource "ome_appliance_network" "code_4" {
		proxy_setting = {
			enable_proxy = false
		}
	}
	`
	testAccNetworkProxyCreateUpdateNil2 := testProvider + `
	resource "ome_appliance_network" "code_4" {
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
					resource.TestCheckResourceAttr("ome_appliance_network.code_4", "proxy_setting.enable_proxy", "false"),
				),
			},
			{
				Config: testAccNetworkProxyCreateUpdateNil2,
			},
		},
	})
}
