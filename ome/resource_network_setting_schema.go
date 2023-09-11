package ome

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkSettingSchema for network setting schema
func NetworkSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID of the ome network setting.",
			Description:         "ID of the ome network setting.",
			Computed:            true,
		},

		"adapter_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Adapter Setting",
			Description:         "Ome Adapter Setting",
			Optional:            true,
			Attributes:          AdapterSettingSchema(),
		},

		"session_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Session Setting",
			Description:         "Ome Session Setting",
			Optional:            true,
			Attributes:          SessionSettingSchema(),
		},

		"time_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Time Setting",
			Description:         "Ome Time Setting",
			Optional:            true,
			Attributes:          TimeSettingSchema(),
		},

		"proxy_setting": schema.SingleNestedAttribute{
			MarkdownDescription: "Ome Proxy Setting",
			Description:         "Ome Proxy Setting",
			Optional:            true,
			Attributes:          ProxySettingSchema(),
		},
	}
}

// AdapterSettingSchema for adapter setting schema
func AdapterSettingSchema() map[string]schema.Attribute {
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
			MarkdownDescription: "If there are multiple interfaces, network configuration changes can be applied to a single interface using the `interface name` of the NIC.",
			Description:         "If there are multiple interfaces, network configuration changes can be applied to a single interface using the `interface name` of the NIC.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"ipv4_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "IPv4 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv4 address",
			Description:         "IPv4 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv4 address",
			Optional:            true,
			Attributes:          IPv4ConfigSchema(),
		},

		"ipv6_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "IPv6 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv6 address",
			Description:         "IPv6 network configuration. (Warning) Ensure that you have an alternate interface to access OpenManage Enterprise as these options can change the current IPv6 address",
			Optional:            true,
			Attributes:          IPv6ConfigSchema(),
		},

		"management_vlan": schema.SingleNestedAttribute{
			MarkdownDescription: "vLAN configuration. settings are applicable for OpenManage Enterprise Modular",
			Description:         "vLAN configuration. settings are applicable for OpenManage Enterprise Modular",
			Optional:            true,
			Attributes:          ManagementVLANSchema(),
		},

		"dns_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "Domain Name System(DNS) settings",
			Description:         "Domain Name System(DNS) settings",
			Optional:            true,
			Attributes:          DNSConfigSchema(),
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

// IPv4ConfigSchema for IPv4 Configuration Schema
func IPv4ConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ipv4": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable access to the network using IPv4.",
			Description:         "Enable or disable access to the network using IPv4.",
			Required:            true,
		},

		"enable_dhcp": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable the automatic request to get an IPv4 address from the IPv4 Dynamic Host Configuration Protocol (DHCP) server. If enable_dhcp option is true, OpenManage Enterprise retrieves the IP configuration—IPv4 address, subnet mask, and gateway from a DHCP server on the existing network.",
			Description:         "Enable or disable the automatic request to get an IPv4 address from the IPv4 Dynamic Host Configuration Protocol (DHCP) server. If enable_dhcp option is true, OpenManage Enterprise retrieves the IP configuration—IPv4 address, subnet mask, and gateway from a DHCP server on the existing network.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"static_ip_address": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 address. This option is applicable when \"enable_dhcp\" is false.",
			Description:         "Static IPv4 address. This option is applicable when \"enable_dhcp\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"static_subnet_mask": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 subnet mask address. This option is applicable when \"enable_dhcp\" is false.",
			Description:         "Static IPv4 subnet mask address.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"static_gateway": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 gateway address. This option is applicable when \"enable_dhcp\" is false.",
			Description:         "Static IPv4 gateway address. This option is applicable when \"enable_dhcp\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "This option allows to automatically request and obtain a DNS server IPv4 address from the DHCP server. This option is applicable when \"enable_dhcp\" is true.",
			Description:         "This option allows to automatically request and obtain a DNS server IPv4 address from the DHCP server. This option is applicable when \"enable_dhcp\" is true.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"static_preferred_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv4 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"static_alternate_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv4 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv4 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}
}

// IPv6ConfigSchema for IPv6 Configuration Schema
func IPv6ConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ipv6": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable access to the network using the IPv6.",
			Description:         "Enable or disable access to the network using the IPv6.",
			Required:            true,
		},

		"enable_auto_configuration": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable the automatic request to get an IPv6 address from the IPv6 DHCP server or router advertisements(RA). If \"enable_auto_configuration\" is true, OME retrieves IP configuration-IPv6 address, prefix, and gateway, from a DHCPv6 server on the existing network",
			Description:         "Enable or disable the automatic request to get an IPv6 address from the IPv6 DHCP server or router advertisements(RA). If \"enable_auto_configuration\" is true, OME retrieves IP configuration-IPv6 address, prefix, and gateway, from a DHCPv6 server on the existing network",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"static_ip_address": schema.StringAttribute{
			MarkdownDescription: "Static IPv6 address. This option is applicable when \"enable_auto_configuration\" is false.",
			Description:         "Static IPv6 address. This option is applicable when \"enable_auto_configuration\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
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
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "This option allows to automatically request and obtain a DNS server IPv6 address from the DHCP server. This option is applicable when \"enable_auto_configuration\" is true.",
			Description:         "This option allows to automatically request and obtain a DNS server IPv6 address from the DHCP server. This option is applicable when \"enable_auto_configuration\" is true.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"static_preferred_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv6 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv6 DNS preferred server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"static_alternate_dns_server": schema.StringAttribute{
			MarkdownDescription: "Static IPv6 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Description:         "Static IPv6 DNS alternate server. This option is applicable when \"use_dhcp_for_dns_server_names\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}
}

// ManagementVLANSchema for management vlan schema
func ManagementVLANSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_vlan": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable vLAN for management.The vLAN configuration cannot be updated if the \"register_with_dns\" field under \"dns_configuration\" is true. WARNING: Ensure that the network cable is plugged to the correct port after the vLAN configurationchanges have been made. If not, the configuration change may not be effective.",
			Description:         "Enable or disable vLAN for management.The vLAN configuration cannot be updated if the \"register_with_dns\" field under \"dns_configuration\" is true. WARNING: Ensure that the network cable is plugged to the correct port after the vLAN configurationchanges have been made. If not, the configuration change may not be effective.",
			Required:            true,
		},

		"id": schema.Int64Attribute{
			MarkdownDescription: "vLAN ID. This option is applicable when \"enable_vlan\" is true.",
			Description:         "vLAN ID. This option is applicable when \"enable_vlan\" is true.",
			Optional:            true,
			Computed:            true,
		},
	}
}

// DNSConfigSchema for DNS Configuration Schema
func DNSConfigSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"register_with_dns": schema.BoolAttribute{
			MarkdownDescription: "Register/Unregister I(dns_name) on the DNS Server.This option cannot be updated if vLAN configuration changes.",
			Description:         "Register/Unregister I(dns_name) on the DNS Server.This option cannot be updated if vLAN configuration changes.",
			Optional:            true,
			Computed:            true,
		},

		"use_dhcp_for_dns_server_names": schema.BoolAttribute{
			MarkdownDescription: "Get the \"dns_domain_name\" using a DHCP server.",
			Description:         "Get the \"dns_domain_name\" using a DHCP server.",
			Optional:            true,
			Computed:            true,
		},

		"dns_name": schema.StringAttribute{
			MarkdownDescription: "DNS name for \"hostname\". This is applicable when \"register_with_dns\" is true.",
			Description:         "DNS name for \"hostname\". This is applicable when \"register_with_dns\" is true.DNSName",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"dns_domain_name": schema.StringAttribute{
			MarkdownDescription: "Static DNS domain name. This is applicable when \"use_dhcp_for_dns_domain_name\" is false.",
			Description:         "Static DNS domain name. This is applicable when \"use_dhcp_for_dns_domain_name\" is false.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}
}

// SessionSettingSchema for the session setting schema.
func SessionSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_universal_timeout": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable the universal inactivity timeout.",
			Description:         "Enable or disable the universal inactivity timeout.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"universal_timeout": schema.Float64Attribute{
			MarkdownDescription: "Duration of inactivity in minutes after which all sessions end. This is applicable when \"enable_universal_timeout\" is true. This is mutually exclusive with \"api_timeout\", \"gui_timeout\", \"ssh_timeout\" and \"serial_timeout\".",
			Description:         "Duration of inactivity in minutes after which all sessions end. This is applicable when \"enable_universal_timeout\" is true. This is mutually exclusive with \"api_timeout\", \"gui_timeout\", \"ssh_timeout\" and \"serial_timeout\".",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Float64{
				float64validator.Between(1, 1400),
			},
		},

		"api_timeout": schema.Float64Attribute{
			MarkdownDescription: "Duration of inactivity in minutes after which the API session ends. This is mutually exclusive with \"universal_timeout\".",
			Description:         "Duration of inactivity in minutes after which the API session ends. This is mutually exclusive with \"universal_timeout\".",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Float64{
				float64validator.Between(1, 1400),
			},
		},

		"api_session": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of API sessions to be allowed.",
			Description:         "The maximum number of API sessions to be allowed.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 100),
			},
		},

		"gui_timeout": schema.Float64Attribute{
			MarkdownDescription: "Duration of inactivity in minutes after which the web interface of Graphical User Interface (GUI) session ends. This is mutually exclusive with \"universal_timeout\".",
			Description:         "Duration of inactivity in minutes after which the web interface of Graphical User Interface (GUI) session ends. This is mutually exclusive with \"universal_timeout\".",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Float64{
				float64validator.Between(1, 1400),
			},
		},

		"gui_session": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of GUI sessions to be allowed.",
			Description:         "The maximum number of GUI sessions to be allowed.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 100),
			},
		},

		"ssh_timeout": schema.Float64Attribute{
			MarkdownDescription: "Duration of inactivity in minutes after which the SSH session ends. This is applicable only for OpenManage Enterprise Modular. This is mutually exclusive with \"universal_timeout\".",
			Description:         "Duration of inactivity in minutes after which the SSH session ends. This is applicable only for OpenManage Enterprise Modular. This is mutually exclusive with \"universal_timeout\".",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Float64{
				float64validator.Between(1, 1400),
			},
		},

		"ssh_session": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of SSH sessions to be allowed. This is applicable to OME-M only.",
			Description:         "The maximum number of SSH sessions to be allowed. This is applicable to OME-M only.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 100),
			},
		},

		"serial_timeout": schema.Float64Attribute{
			MarkdownDescription: "Duration of inactivity in minutes after which the serial console session ends.This is applicable only for OpenManage Enterprise Modular. This is mutually exclusive with \"universal_timeout\".",
			Description:         "Duration of inactivity in minutes after which the serial console session ends.This is applicable only for OpenManage Enterprise Modular. This is mutually exclusive with \"universal_timeout\"",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Float64{
				float64validator.Between(1, 1400),
			},
		},

		"serial_session": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of serial console sessions to be allowed. This is applicable only for OpenManage Enterprise Modular.",
			Description:         "The maximum number of serial console sessions to be allowed. This is applicable only for OpenManage Enterprise Modular.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 100),
			},
		},
	}
}

// TimeSettingSchema for time setting schema
func TimeSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_ntp": schema.BoolAttribute{
			MarkdownDescription: "Enables or disables Network Time Protocol(NTP).If \"enable_ntp\" is false, then the NTP addresses reset to their default values.",
			Description:         "Enables or disables Network Time Protocol(NTP).If \"enable_ntp\" is false, then the NTP addresses reset to their default values.",
			Optional:            true,
			Computed:            true,
		},

		"system_time": schema.StringAttribute{
			MarkdownDescription: "Time in the current system. This option is only applicable when \"enable_ntp\" is false. This option must be provided in following format 'yyyy-mm-dd hh:mm:ss'.",
			Description:         "Time in the current system. This option is only applicable when \"enable_ntp\" is false. This option must be provided in following format 'yyyy-mm-dd hh:mm:ss'.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"time_zone": schema.StringAttribute{
			MarkdownDescription: "The valid timezone ID to be used. This option is applicable for both system time and NTP time synchronization.",
			Description:         "The valid timezone ID to be used. This option is applicable for both system time and NTP time synchronization.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"primary_ntp_address": schema.StringAttribute{
			MarkdownDescription: "The primary NTP address. This option is applicable when \"enable_ntp\" is true.",
			Description:         "The primary NTP address. This option is applicable when \"enable_ntp\" is true.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"secondary_ntp_address1": schema.StringAttribute{
			MarkdownDescription: "The first secondary NTP address. This option is applicable when \"enable_ntp\" is true.",
			Description:         "The first secondary NTP address. This option is applicable when \"enable_ntp\" is true.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"secondary_ntp_address2": schema.StringAttribute{
			MarkdownDescription: "The second secondary NTP address. This option is applicable when \"enable_ntp\" is true.",
			Description:         "The second secondary NTP address. This option is applicable when \"enable_ntp\" is true.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}
}

// ProxySettingSchema for proxy setting schema
func ProxySettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"enable_proxy": schema.BoolAttribute{
			MarkdownDescription: "Enables or disables the HTTP proxy configuration. If \"enable proxy\" is false, then the HTTP proxy configuration is set to its default value.",
			Description:         "Enables or disables the HTTP proxy configuration. If \"enable proxy\" is false, then the HTTP proxy configuration is set to its default value.",
			Required:            true,
		},

		"ip_address": schema.StringAttribute{
			MarkdownDescription: "Proxy server address. This option is mandatory when \"enable_proxy\" is true.",
			Description:         "Proxy server address. This option is mandatory when \"enable_proxy\" is true.",
			Optional:            true,
			Computed:            true,
		},

		"proxy_port": schema.Int64Attribute{
			MarkdownDescription: "Proxy server's port number. This option is mandatory when \"enable_proxy\" is true.",
			Description:         "Proxy server's port number. This option is mandatory when \"enable_proxy\" is true.",
			Optional:            true,
			Computed:            true,
		},

		"enable_authentication": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable proxy authentication. If \"enable_authentication\" is true, \"proxy_username\" and \"proxy_password\" must be provided. If \"enable_authentication\" is false, the proxy username and password are set to its default values.",
			Description:         "Enable or disable proxy authentication. If \"enable_authentication\" is true, \"proxy_username\" and \"proxy_password\" must be provided. If \"enable_authentication\" is false, the proxy username and password are set to its default values.",
			Optional:            true,
			Computed:            true,
		},

		"username": schema.StringAttribute{
			MarkdownDescription: "Proxy server username. This option is mandatory when \"enable_authentication\" is true.",
			Description:         "Proxy server username. This option is mandatory when \"enable_authentication\" is true.",
			Optional:            true,
			Computed:            true,
		},

		"password": schema.StringAttribute{
			MarkdownDescription: "Proxy server password. This option is mandatory when \"enable_authentication\" is true.",
			Description:         "Proxy server password. This option is mandatory when \"enable_authentication\" is true.",
			Optional:            true,
			Computed:            true,
		},
	}
}
