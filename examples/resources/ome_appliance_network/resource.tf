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
    # configuring number of parallel session 
    api_session = 10
    gui_session = 11
  }
  // to change the ome appliance time setting
  proxy_setting = {
    enable_proxy = true
    # ip address and port is require when enable proxy is true 
    ip_address            = "1.0.0.1"
    proxy_port            = 443
    enable_authentication = true
    # username and password is required when enable authentication is true 
    username = "root"
    password = "root"
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

# session setting without universal_timeout
resource "ome_appliance_network" "net2" {
  // to change the ome appliance session setting
  session_setting = {
    # configuring timeout attributes
    enable_universal_timeout = false
    api_timeout              = 20
    gui_timeout              = 20
    # configuring number of parallel session 
    api_session = 10
    gui_session = 11
  }
}

# static IP adapter setting 
resource "ome_appliance_network" "net3" {
  adapter_setting = {
    enable_nic     = true
    interface_name = "ens160"
    reboot_delay   = 0
    ipv6_configuration = {
      enable_ipv6               = true
      enable_auto_configuration = false
      # static ipv6 configuration
      static_ip_address    = "<static-ipv6>"
      static_gateway       = "<static-ipv6-gateway>"
      static_prefix_length = 64
      # static dns configuration
      use_dhcp_for_dns_server_names = false
      static_preferred_dns_server   = "<dns-server-1>"
      static_alternate_dns_server   = "<dns-server-2>"
    }

    ipv4_configuration = {
      enable_ipv4 = true
      enable_dhcp = false
      # static ipv4 configuration
      static_ip_address  = "10.10.10.10"
      static_subnet_mask = "255.255.255.0"
      static_gateway     = "10.10.10.1"
      # static dns configuration
      use_dhcp_for_dns_server_names = false
      static_alternate_dns_server   = "10.10.10.2"
      static_preferred_dns_server   = "10.10.10.3"
    }
  }
}
