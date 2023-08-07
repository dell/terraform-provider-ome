package ome

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

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
			MarkdownDescription: "Enable Nic",
			Description:         "Enable Nic",
			Optional:            true,
			Computed:            true,
		},

		"interface_name": schema.StringAttribute{
			MarkdownDescription: "Interface Name",
			Description:         "Interface Name",
			Optional:            true,
			Computed:            true,
		},

		"ipv4_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "IPV4Config",
			Description:         "IPV4Config",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeIPv4ConfigSchema(),
		},

		"ipv6_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "IPV6Config",
			Description:         "IPV6Config",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeIPv6ConfigSchema(),
		},

		"management_vlan": schema.SingleNestedAttribute{
			MarkdownDescription: "Management VLAN",
			Description:         "Management VLAN",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeManagementVLANSchema(),
		},

		"dns_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "DNSConfig",
			Description:         "DNSConfig",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeDNSConfigSchema(),
		},

		"reboot_delay": schema.Int64Attribute{
			MarkdownDescription: "Reboot Delay",
			Description:         "Reboot Delay",
			Optional:            true,
			Computed:            true,
		},

		"job_id": schema.Int64Attribute{
			MarkdownDescription: "Job ID",
			Description:         "Job ID",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeIPv4ConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ipv4": schema.BoolAttribute{
			MarkdownDescription: "Enable IPv4",
			Description:         "Enable IPv4",
			Optional:            true,
			Computed:            true,
		},

		"enable_dhcp": schema.BoolAttribute{
			MarkdownDescription: "Enable DHCP",
			Description:         "Enable DHCP",
			Optional:            true,
			Computed:            true,
		},

		"static_ip_address": schema.StringAttribute{
			MarkdownDescription: "Static IPAddress",
			Description:         "Static IPAddress",
			Optional:            true,
			Computed:            true,
		},

		"static_subnet_mask": schema.StringAttribute{
			MarkdownDescription: "Static Subnet Mask",
			Description:         "Static Subnet Mask",
			Optional:            true,
			Computed:            true,
		},

		"static_gateway": schema.StringAttribute{
			MarkdownDescription: "Static Gateway",
			Description:         "Static Gateway",
			Optional:            true,
			Computed:            true,
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "Use DHCPfor DNSServer Names",
			Description:         "Use DHCPfor DNSServer Names",
			Optional:            true,
			Computed:            true,
		},

		"static_preferred_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static Preferred DNSServer",
			Description:         "Static Preferred DNSServer",
			Optional:            true,
			Computed:            true,
		},

		"static_alternate_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static Alternate DNSServer",
			Description:         "Static Alternate DNSServer",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeIPv6ConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ipv6": schema.BoolAttribute{
			MarkdownDescription: "Enable IPv6",
			Description:         "Enable IPv6",
			Optional:            true,
			Computed:            true,
		},

		"enable_auto_configuration": schema.BoolAttribute{
			MarkdownDescription: "Enable Auto Configuration",
			Description:         "Enable Auto Configuration",
			Optional:            true,
			Computed:            true,
		},

		"static_ip_address": schema.StringAttribute{
			MarkdownDescription: "Static IPAddress",
			Description:         "Static IPAddress",
			Optional:            true,
			Computed:            true,
		},

		"static_prefix_length": schema.Int64Attribute{
			MarkdownDescription: "Static Prefix Length",
			Description:         "Static Prefix Length",
			Optional:            true,
			Computed:            true,
		},

		"static_gateway": schema.StringAttribute{
			MarkdownDescription: "Static Gateway",
			Description:         "Static Gateway",
			Optional:            true,
			Computed:            true,
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "Use DHCPfor DNSServer Names",
			Description:         "Use DHCPfor DNSServer Names",
			Optional:            true,
			Computed:            true,
		},

		"static_preferred_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static Preferred DNSServer",
			Description:         "Static Preferred DNSServer",
			Optional:            true,
			Computed:            true,
		},

		"static_alternate_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static Alternate DNSServer",
			Description:         "Static Alternate DNSServer",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeManagementVLANSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_vlan": schema.BoolAttribute{
			MarkdownDescription: "Enable VLAN",
			Description:         "Enable VLAN",
			Optional:            true,
			Computed:            true,
		},

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",
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
