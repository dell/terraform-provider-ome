package ome

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func OmeNetworkSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"adapter_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Adapter Setting",
			Description:         "Ome Adapter Setting",
			Optional:            true,
			Attributes:          OmeAdapterSettingSchema(),
		},

		"session_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Session Setting",
			Description:         "Ome Session Setting",
			Optional:            true,
			Attributes:          OmeSessionSettingSchema(),
		},

		"time_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Time Setting",
			Description:         "Ome Time Setting",
			Optional:            true,
			Attributes:          OmeTimeSettingSchema(),
		},

		"proxy_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Proxy Setting",
			Description:         "Ome Proxy Setting",
			Optional:            true,
			Attributes:          OmeProxySettingSchema(),
		},
	}
}

func OmeAdapterSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_nic": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable Network Interface Card (NIC) configuration",
			Description:         "Enable or disable Network Interface Card (NIC) configuration",
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(true)),
			},
		},

		"interface_name": schema.StringAttribute{
			MarkdownDescription: "If there are multiple interfaces, network configuration changes can be applied to a single interface using the `interface name` of the NIC.If this option is not specified, Primary interface is chosen by default.",
			Description:         "If there are multiple interfaces, network configuration changes can be applied to a single interface using the `interface name` of the NIC.If this option is not specified, Primary interface is chosen by default.Interface Name",
			Optional:            true,
			Computed:            true,
		},

		"ipv4_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "IPv4 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv4 address",
			Description:         "IPv4 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv4 address",
			Optional:            true,
			Attributes:          OmeIPv4ConfigSchema(),
		},

		"ipv6_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "IPv6 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv6 address",
			Description:         "IPv6 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv6 address",
			Optional:            true,
			Attributes:          OmeIPv6ConfigSchema(),
		},

		"management_vlan": schema.SingleNestedAttribute{
			MarkdownDescription: "vLAN configuration. settings are applicable for OpenManage Enterprise Modular",
			Description:         "vLAN configuration. settings are applicable for OpenManage Enterprise Modular",
			Optional:            true,
			Attributes:          OmeManagementVLANSchema(),
		},

		"dns_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "Domain Name System(DNS) settings",
			Description:         "Domain Name System(DNS) settings",
			Optional:            true,
			Attributes:          OmeDNSConfigSchema(),
		},

		"reboot_delay": schema.Int64Attribute{
			MarkdownDescription: "The time in seconds, after which settings are applied",
			Description:         "The time in seconds, after which settings are applied",
			Optional:            true,
			Computed:            true,
		},

		"job_id": schema.Int64Attribute{
			MarkdownDescription: "Job ID ",
			Description:         "Job ID",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeIPv4ConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ipv4": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable access to the network using IPv4.",
			Description:         "Enable or disable access to the network using IPv4.",
			Required: true,
		},

		"enable_dhcp": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable the automatic request to get an IPv4 address from the IPv4 Dynamic Host Configuration Protocol (DHCP) server. If enable_dhcp option is true, OpenManage Enterprise retrieves the IP configuration—IPv4 address, subnet mask, and gateway from a DHCP server on the existing network.",
			Description:         "Enable or disable the automatic request to get an IPv4 address from the IPv4 Dynamic Host Configuration Protocol (DHCP) server. If enable_dhcp option is true, OpenManage Enterprise retrieves the IP configuration—IPv4 address, subnet mask, and gateway from a DHCP server on the existing network.",
			Optional:            true,
			Computed:            true,
		},

		"static_ip_address": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 address. This option is applicable when \"enable_dhcp\" is false.",
			Description:         "Static IPv4 address. This option is applicable when \"enable_dhcp\" is false.",
			Optional:            true,
			Computed:            true,
		},

		"static_subnet_mask": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 subnet mask address. This option is applicable when \"enable_dhcp\" is false.",
			Description:         "Static IPv4 subnet mask address.",
			Optional:            true,
			Computed:            true,
		},

		"static_gateway": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 gateway address. This option is applicable when \"enable_dhcp\" is false.",
			Description:         "Static IPv4 gateway address. This option is applicable when \"enable_dhcp\" is false.",
			Optional:            true,
			Computed:            true,
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "This option allows to automatically request and obtain a DNS server IPv4 address from the DHCP server. This option is applicable when \"enable_dhcp\" is true.",
			Description:         "This option allows to automatically request and obtain a DNS server IPv4 address from the DHCP server. This option is applicable when \"enable_dhcp\" is true.",
			Optional:            true,
			Computed:            true,
		},

		"static_preferred_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv4 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
		},

		"static_alternate_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv4 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeIPv6ConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ipv6": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable access to the network using the IPv6.",
			Description:         "Enable or disable access to the network using the IPv6.",
			Required: true,
		},

		"enable_auto_configuration": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable the automatic request to get an IPv6 address from the IPv6 DHCP server or router advertisements(RA). If \"enable_auto_configuration\" is true, OME retrieves IP configuration-IPv6 address, prefix, and gateway, from a DHCPv6 server on the existing network",
			Description:        "Enable or disable the automatic request to get an IPv6 address from the IPv6 DHCP server or router advertisements(RA). If \"enable_auto_configuration\" is true, OME retrieves IP configuration-IPv6 address, prefix, and gateway, from a DHCPv6 server on the existing network",
			Optional:            true,
			Computed:            true,
		},

		"static_ip_address": schema.StringAttribute{
			MarkdownDescription: "Static IPv6 address. This option is applicable when \"enable_auto_configuration\" is false.",
			Description: "Static IPv6 address. This option is applicable when \"enable_auto_configuration\" is false.",
			Optional: true,
			Computed: true,
		},

		"static_prefix_length": schema.Int64Attribute{
			MarkdownDescription: "Static IPv6 prefix length. This option is applicable when \"enable_auto_configuration\" is false.",
			Description:         "Static IPv6 prefix length. This option is applicable when \"enable_auto_configuration\" is false.",
			Optional:            true,
			Computed:            true,
		},

		"static_gateway": schema.StringAttribute{
			MarkdownDescription: "Static IPv6 gateway address. This option is applicable when \"enable_auto_configuration\" is false.",
			Description:         "Static IPv6 gateway address. This option is applicable when \"enable_auto_configuration\" is false.",
			Optional:            true,
			Computed:            true,
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "This option allows to automatically request and obtain a DNS server IPv6 address from the DHCP server. This option is applicable when \"enable_auto_configuration\" is true.",
			Description:         "This option allows to automatically request and obtain a DNS server IPv6 address from the DHCP server. This option is applicable when \"enable_auto_configuration\" is true.",
			Optional:            true,
			Computed:            true,
		},

		"static_preferred_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv6 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv6 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
		},

		"static_alternate_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv6 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv6 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeManagementVLANSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_vlan": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable vLAN for management.The vLAN configuration cannot be updated if the \"register_with_dns\" field under \"dns_configuration\" is true. WARNING: Ensure that the network cable is plugged to the correct port after the vLAN configurationchanges have been made. If not, the configuration change may not be effective.",
			Description:         "Enable or disable vLAN for management.The vLAN configuration cannot be updated if the \"register_with_dns\" field under \"dns_configuration\" is true. WARNING: Ensure that the network cable is plugged to the correct port after the vLAN configurationchanges have been made. If not, the configuration change may not be effective.",
			Optional:            true,
			Computed:            true,
		},

		"id": schema.Int64Attribute{
			MarkdownDescription: "vLAN ID. This option is applicable when \"enable_vlan\" is true.",
			Description:         "vLAN ID. This option is applicable when \"enable_vlan\" is true.",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeDNSConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"register_with_dns": schema.BoolAttribute{
			MarkdownDescription: "Register With DNS",
			Description:         "Register With DNS",
			Optional:            true,
			Computed:            true,
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "Use DHCPfor DNSServer Names",
			Description:         "Use DHCPfor DNSServer Names",
			Optional:            true,
			Computed:            true,
		},

		"dns_name": schema.StringAttribute{
			MarkdownDescription: "DNSName",
			Description:         "DNSName",
			Optional:            true,
			Computed:            true,
		},

		"dns_domain_name": schema.StringAttribute{
			MarkdownDescription: "DNSDomain Name",
			Description:         "DNSDomain Name",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeSessionSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_universal_timeout": schema.BoolAttribute{
			MarkdownDescription: "Enable Universal Timeout",
			Description:         "Enable Universal Timeout",
			Optional:            true,
			Computed:            true,
		},

		"universal_timeout": schema.Float64Attribute{
			MarkdownDescription: "Unversal Timeout",
			Description:         "Unversal Timeout",
			Optional:            true,
			Computed:            true,
		},

		"api_timeout": schema.Float64Attribute{
			MarkdownDescription: "APITimeout",
			Description:         "APITimeout",
			Optional:            true,
			Computed:            true,
		},

		"api_session": schema.Int64Attribute{
			MarkdownDescription: "APISession",
			Description:         "APISession",
			Optional:            true,
			Computed:            true,
		},

		"gui_timeout": schema.Float64Attribute{
			MarkdownDescription: "GUITimeout",
			Description:         "GUITimeout",
			Optional:            true,
			Computed:            true,
		},

		"gui_session": schema.Int64Attribute{
			MarkdownDescription: "GUISession",
			Description:         "GUISession",
			Optional:            true,
			Computed:            true,
		},

		"ssh_timeout": schema.Float64Attribute{
			MarkdownDescription: "SSHTimeout",
			Description:         "SSHTimeout",
			Optional:            true,
			Computed:            true,
		},

		"ssh_session": schema.Int64Attribute{
			MarkdownDescription: "SSHSession",
			Description:         "SSHSession",
			Optional:            true,
			Computed:            true,
		},

		"serial_timeout": schema.Float64Attribute{
			MarkdownDescription: "Serial Timeout",
			Description:         "Serial Timeout",
			Optional:            true,
			Computed:            true,
		},

		"serial_session": schema.Int64Attribute{
			MarkdownDescription: "Serial Session",
			Description:         "Serial Session",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeTimeSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ntp": schema.BoolAttribute{
			MarkdownDescription: "Enable NTP",
			Description:         "Enable NTP",
			Optional:            true,
			Computed:            true,
		},

		"system_time": schema.StringAttribute{
			MarkdownDescription: "System Time",
			Description:         "System Time",
			Optional:            true,
			Computed:            true,
		},

		"time_zone": schema.StringAttribute{
			MarkdownDescription: "Time Zone",
			Description:         "Time Zone",
			Optional:            true,
			Computed:            true,
		},

		"primary_ntp_address": schema.StringAttribute{
			MarkdownDescription: "Primary NTPAddress",
			Description:         "Primary NTPAddress",
			Optional:            true,
			Computed:            true,
		},

		"secondary_ntp_address1": schema.StringAttribute{
			MarkdownDescription: "Secondary NTPAddress1",
			Description:         "Secondary NTPAddress1",
			Optional:            true,
			Computed:            true,
		},

		"secondary_ntp_address2": schema.StringAttribute{
			MarkdownDescription: "Secondary NTPAddress2",
			Description:         "Secondary NTPAddress2",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeProxySettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_proxy": schema.BoolAttribute{
			MarkdownDescription: "Enable Proxy",
			Description:         "Enable Proxy",
			Optional:            true,
			Computed:            true,
		},

		"ip_address": schema.StringAttribute{
			MarkdownDescription: "IPAddress",
			Description:         "IPAddress",
			Optional:            true,
			Computed:            true,
		},

		"proxy_port": schema.Int64Attribute{
			MarkdownDescription: "Proxy Port",
			Description:         "Proxy Port",
			Optional:            true,
			Computed:            true,
		},

		"enable_authentication": schema.BoolAttribute{
			MarkdownDescription: "Enable Authentication",
			Description:         "Enable Authentication",
			Optional:            true,
			Computed:            true,
		},

		"Username": schema.StringAttribute{
			MarkdownDescription: "Username",
			Description:         "Username",
			Optional:            true,
			Computed:            true,
		},

		"Password": schema.StringAttribute{
			MarkdownDescription: "Password",
			Description:         "Password",
			Optional:            true,
			Computed:            true,
		},
	}
}
