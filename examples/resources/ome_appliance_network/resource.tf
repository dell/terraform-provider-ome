resource "ome_appliance_network" "net1" {
  // to change the ome appliance time setting
  time_setting = {
    time_zone              = "TZ_ID_65"
    enable_ntp             = true
    primary_ntp_address    = "1.0.0.0"
    secondary_ntp_address1 = "2.0.0.0"
    secondary_ntp_address2 = "3.0.0.0"
  }
  // to change the ome appliance session setting
  session_setting = {
    enable_universal_timeout = true
    universal_timeout        = 20
    api_session              = 10
    gui_session              = 11
  }
  // to change the ome appliance time setting
  proxy_setting = {
    enable_proxy          = true
    ip_address            = "1.0.0.1"
    proxy_port            = 443
    enable_authentication = true
    username              = "root"
    password              = "root"
  }
  // to change the ome appliance adapter setting
  adapter_setting = {
    enable_nic     = true
    interface_name = "ens160"
    reboot_delay   = 0
    management_vlan = {
      enable_vlan = false
    }
    ipv6_configuration = {
      enable_ipv6               = true
      enable_auto_configuration = true
    }
    dns_configuration = {
      register_with_dns = false
    }
  }
}
