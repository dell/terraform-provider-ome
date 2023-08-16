package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// NetworkAdapterSetting for network adapter setting
type NetworkAdapterSetting struct {
	InterfaceName     string            `json:"InterfaceName"`
	ProfileName       string            `json:"ProfileName"`
	EnableNIC         bool              `json:"EnableNIC"`
	Ipv4Configuration Ipv4Configuration `json:"Ipv4Configuration,omitempty"`
	Ipv6Configuration Ipv6Configuration `json:"Ipv6Configuration,omitempty"`
	ManagementVLAN    ManagementVLAN    `json:"ManagementVLAN,omitempty"`
	DNSConfiguration  DNSConfiguration  `json:"DnsConfiguration,omitempty"`
	CurrentSettings   CurrentSettings   `json:"CurrentSettings,omitempty"`
	Delay             int               `json:"Delay"`
	PrimaryInterface  bool              `json:"PrimaryInterface"`
}

// Ipv4Configuration for IPv4 Configuration
type Ipv4Configuration struct {
	Enable                   bool   `json:"Enable"`
	EnableDHCP               bool   `json:"EnableDHCP"`
	StaticIPAddress          string `json:"StaticIPAddress"`
	StaticSubnetMask         string `json:"StaticSubnetMask"`
	StaticGateway            string `json:"StaticGateway"`
	UseDHCPForDNSServerNames bool   `json:"UseDHCPForDNSServerNames"`
	StaticPreferredDNSServer string `json:"StaticPreferredDNSServer"`
	StaticAlternateDNSServer string `json:"StaticAlternateDNSServer"`
}

// Ipv6Configuration for IPv6 Configuration
type Ipv6Configuration struct {
	Enable                   bool   `json:"Enable"`
	EnableAutoConfiguration  bool   `json:"EnableAutoConfiguration"`
	StaticIPAddress          string `json:"StaticIPAddress"`
	StaticPrefixLength       int    `json:"StaticPrefixLength"`
	StaticGateway            string `json:"StaticGateway"`
	UseDHCPForDNSServerNames bool   `json:"UseDHCPForDNSServerNames"`
	StaticPreferredDNSServer string `json:"StaticPreferredDNSServer"`
	StaticAlternateDNSServer string `json:"StaticAlternateDNSServer"`
}

// ManagementVLAN for management vlan configuration
type ManagementVLAN struct {
	EnableVLAN bool `json:"EnableVLAN"`
	ID         int  `json:"Id"`
}

// DNSConfiguration for DNS Configuration
type DNSConfiguration struct {
	RegisterWithDNS               bool   `json:"RegisterWithDNS"`
	DNSName                       string `json:"DnsName"`
	UseDHCPForDNSDomainName       bool   `json:"UseDHCPForDNSDomainName"`
	DNSDomainName                 string `json:"DnsDomainName"`
	FqdndomainName                string `json:"FqdndomainName"`
	Ipv4CurrentPreferredDNSServer string `json:"Ipv4CurrentPreferredDNSServer"`
	Ipv4CurrentAlternateDNSServer string `json:"Ipv4CurrentAlternateDNSServer"`
	Ipv6CurrentPreferredDNSServer string `json:"Ipv6CurrentPreferredDNSServer"`
	Ipv6CurrentAlternateDNSServer string `json:"Ipv6CurrentAlternateDNSServer"`
}

// Ipv4Settings for IPv4 Setting
type Ipv4Settings struct {
	Enable                   bool     `json:"Enable"`
	EnableDhcp               bool     `json:"EnableDhcp"`
	CurrentIPAddress         []string `json:"CurrentIPAddress"`
	CurrentSubnetMask        string   `json:"CurrentSubnetMask"`
	CurrentGateway           string   `json:"CurrentGateway"`
	UseDHCPForDNSServerNames bool     `json:"UseDHCPForDNSServerNames"`
	Ipv4DNS                  []string `json:"Ipv4Dns"`
}

// Ipv6Settings for IPv6 Setting
type Ipv6Settings struct {
	Enable                   bool   `json:"Enable"`
	EnableAutoConfiguration  bool   `json:"EnableAutoConfiguration"`
	CurrentIPAddress         []any  `json:"CurrentIPAddress"`
	CurrentGateway           string `json:"CurrentGateway"`
	CurrentLinkLocalAddress  string `json:"CurrentLinkLocalAddress"`
	UseDHCPForDNSServerNames bool   `json:"UseDHCPForDNSServerNames"`
	Ipv6DNS                  []any  `json:"Ipv6Dns"`
}

// DNSSetting for DNS Setting
type DNSSetting struct {
	DNSFQDName    string `json:"DnsFQDName"`
	DNSDomainName string `json:"DnsDomainName"`
}

// CurrentSettings for getting current setting of network adapter
type CurrentSettings struct {
	Ipv4Settings Ipv4Settings `json:"Ipv4Settings"`
	Ipv6Settings Ipv6Settings `json:"Ipv6Settings"`
	DNSSetting   DNSSetting   `json:"DnsSetting"`
}

// UpdateNetworkAdapterSetting to update network adapter
type UpdateNetworkAdapterSetting struct {
	InterfaceName     string            `json:"InterfaceName"`
	ProfileName       string            `json:"ProfileName"`
	EnableNIC         bool              `json:"EnableNIC"`
	Ipv4Configuration Ipv4Configuration `json:"Ipv4Configuration"`
	Ipv6Configuration Ipv6Configuration `json:"Ipv6Configuration"`
	DNSConfiguration  DNSConfiguration  `json:"DnsConfiguration"`
	Delay             int               `json:"Delay"`
	PrimaryInterface  bool              `json:"PrimaryInterface"`
}

// NetworkSessions to get network session
type NetworkSessions struct {
	SessionList []SessionInfo `json:"value"`
}

// SessionInfo to get session info
type SessionInfo struct {
	SessionType    string `json:"SessionType"`
	MaxSessions    int    `json:"MaxSessions"`
	SessionTimeout int    `json:"SessionTimeout"`
}

// TimeConfiguration to get time configuration
type TimeConfiguration struct {
	OdataContext         string `json:"@odata.context"`
	OdataType            string `json:"@odata.type"`
	OdataID              string `json:"@odata.id"`
	TimeZone             string `json:"TimeZone"`
	TimeZoneIDLinux      string `json:"TimeZoneIdLinux"`
	TimeZoneIDWindows    string `json:"TimeZoneIdWindows"`
	EnableNTP            bool   `json:"EnableNTP"`
	PrimaryNTPAddress    any    `json:"PrimaryNTPAddress"`
	SecondaryNTPAddress1 any    `json:"SecondaryNTPAddress1"`
	SecondaryNTPAddress2 any    `json:"SecondaryNTPAddress2"`
	SystemTime           string `json:"SystemTime"`
	TimeSource           string `json:"TimeSource"`
	UtcTime              string `json:"UtcTime"`
}

// TimeZones to get all time zones.
type TimeZones struct {
	OdataContext string     `json:"@odata.context"`
	OdataCount   int        `json:"@odata.count"`
	Value        []TimeZone `json:"value"`
}

// TimeZone for one time zone.
type TimeZone struct {
	OdataType        string `json:"@odata.type"`
	Utcoffsetminutes int    `json:"Utcoffsetminutes"`
	ID               string `json:"Id"`
	Name             string `json:"Name"`
}

// TimeConfigPayload to get time configuration payload.
type TimeConfigPayload struct {
	TimeZone             string `json:"TimeZone"`
	EnableNTP            bool   `json:"EnableNTP"`
	PrimaryNTPAddress    any    `json:"PrimaryNTPAddress"`
	SecondaryNTPAddress1 any    `json:"SecondaryNTPAddress1"`
	SecondaryNTPAddress2 any    `json:"SecondaryNTPAddress2"`
	SystemTime           string `json:"SystemTime"`
}

// TimeConfigResponse to get time config response.
type TimeConfigResponse struct {
	TimeZone             string `json:"TimeZone"`
	TimeZoneIDLinux      any    `json:"TimeZoneIdLinux"`
	TimeZoneIDWindows    any    `json:"TimeZoneIdWindows"`
	EnableNTP            bool   `json:"EnableNTP"`
	PrimaryNTPAddress    any    `json:"PrimaryNTPAddress"`
	SecondaryNTPAddress1 any    `json:"SecondaryNTPAddress1"`
	SecondaryNTPAddress2 any    `json:"SecondaryNTPAddress2"`
	SystemTime           any    `json:"SystemTime"`
	TimeSource           string `json:"TimeSource"`
	UtcTime              any    `json:"UtcTime"`
	JobID                any    `json:"JobId"`
}

// ProxyConfiguration to get proxy configuration
type ProxyConfiguration struct {
	IPAddress            string `json:"IpAddress"`
	PortNumber           int    `json:"PortNumber"`
	Username             string `json:"Username"`
	Password             string `json:"Password"`
	EnableAuthentication bool   `json:"EnableAuthentication"`
	EnableProxy          bool   `json:"EnableProxy"`
	SslCheckDisabled     bool   `json:"SslCheckDisabled"`
	ProxyExclusionList   any    `json:"ProxyExclusionList"`
}

// PayloadProxyConfiguration for api payload of proxy configuration
type PayloadProxyConfiguration struct {
	IPAddress            string `json:"IpAddress"`
	PortNumber           int    `json:"PortNumber"`
	EnableAuthentication bool   `json:"EnableAuthentication"`
	EnableProxy          bool   `json:"EnableProxy"`
	Username             string `json:"Username"`
	Password             string `json:"Password"`
}

// tfsdk struct definition

// OmeNetworkSetting for network terraform attribute
type OmeNetworkSetting struct {
	ID                types.String       `tfsdk:"id"`
	OmeAdapterSetting *OmeAdapterSetting `tfsdk:"adapter_setting"`
	OmeSessionSetting *OmeSessionSetting `tfsdk:"session_setting"`
	OmeTimeSetting    *OmeTimeSetting    `tfsdk:"time_setting"`
	OmeProxySetting   *OmeProxySetting   `tfsdk:"proxy_setting"`
}

// OmeAdapterSetting for adapter_setting terraform attribute.
type OmeAdapterSetting struct {
	EnableNic      types.Bool         `tfsdk:"enable_nic"`
	InterfaceName  types.String       `tfsdk:"interface_name"`
	IPV4Config     *OmeIPv4Config     `tfsdk:"ipv4_configuration"`
	IPV6Config     *OmeIPv6Config     `tfsdk:"ipv6_configuration"`
	ManagementVLAN *OmeManagementVLAN `tfsdk:"management_vlan"`
	DNSConfig      *OmeDNSConfig      `tfsdk:"dns_configuration"`
	RebootDelay    types.Int64        `tfsdk:"reboot_delay"`
	JobID          types.Int64        `tfsdk:"job_id"`
}

// OmeIPv4Config for ipv4_configuration terraform attribute
type OmeIPv4Config struct {
	EnableIPv4               types.Bool   `tfsdk:"enable_ipv4"`
	EnableDHCP               types.Bool   `tfsdk:"enable_dhcp"`
	StaticIPAddress          types.String `tfsdk:"static_ip_address"`
	StaticSubnetMask         types.String `tfsdk:"static_subnet_mask"`
	StaticGateway            types.String `tfsdk:"static_gateway"`
	UseDHCPforDNSServerNames types.Bool   `tfsdk:"use_dhcp_for_dns_server_names"`
	StaticPreferredDNSServer types.String `tfsdk:"static_preferred_dns_server"`
	StaticAlternateDNSServer types.String `tfsdk:"static_alternate_dns_server"`
}

// OmeIPv6Config for ipv6_configuration terraform attribute
type OmeIPv6Config struct {
	EnableIPv6               types.Bool   `tfsdk:"enable_ipv6"`
	EnableAutoConfiguration  types.Bool   `tfsdk:"enable_auto_configuration"`
	StaticIPAddress          types.String `tfsdk:"static_ip_address"`
	StaticPrefixLength       types.Int64  `tfsdk:"static_prefix_length"`
	StaticGateway            types.String `tfsdk:"static_gateway"`
	UseDHCPforDNSServerNames types.Bool   `tfsdk:"use_dhcp_for_dns_server_names"`
	StaticPreferredDNSServer types.String `tfsdk:"static_preferred_dns_server"`
	StaticAlternateDNSServer types.String `tfsdk:"static_alternate_dns_server"`
}

// OmeManagementVLAN for management_vlan terraform attribute.
type OmeManagementVLAN struct {
	EnableVLAN types.Bool  `tfsdk:"enable_vlan"`
	ID         types.Int64 `tfsdk:"id"`
}

// OmeDNSConfig for dns_configuration terraform attribute.
type OmeDNSConfig struct {
	RegisterWithDNS          types.Bool   `tfsdk:"register_with_dns"`
	UseDHCPforDNSServerNames types.Bool   `tfsdk:"use_dhcp_for_dns_server_names"`
	DNSName                  types.String `tfsdk:"dns_name"`
	DNSDomainName            types.String `tfsdk:"dns_domain_name"`
}

// OmeSessionSetting for session_setting terraform attribute.
type OmeSessionSetting struct {
	EnableUniversalTimeout types.Bool    `tfsdk:"enable_universal_timeout"`
	UniversalTimeout       types.Float64 `tfsdk:"universal_timeout"`
	APITimeout             types.Float64 `tfsdk:"api_timeout"`
	APISession             types.Int64   `tfsdk:"api_session"`
	GUITimeout             types.Float64 `tfsdk:"gui_timeout"`
	GUISession             types.Int64   `tfsdk:"gui_session"`
	SSHTimeout             types.Float64 `tfsdk:"ssh_timeout"`
	SSHSession             types.Int64   `tfsdk:"ssh_session"`
	SerialTimeout          types.Float64 `tfsdk:"serial_timeout"`
	SerialSession          types.Int64   `tfsdk:"serial_session"`
}

// OmeTimeSetting for time_setting terraform attribute.
type OmeTimeSetting struct {
	EnableNTP            types.Bool   `tfsdk:"enable_ntp"`
	SystemTime           types.String `tfsdk:"system_time"`
	TimeZone             types.String `tfsdk:"time_zone"`
	PrimaryNTPAddress    types.String `tfsdk:"primary_ntp_address"`
	SecondaryNTPAddress1 types.String `tfsdk:"secondary_ntp_address1"`
	SecondaryNTPAddress2 types.String `tfsdk:"secondary_ntp_address2"`
}

// OmeProxySetting for proxy_setting terraform attribute.
type OmeProxySetting struct {
	EnableProxy          types.Bool   `tfsdk:"enable_proxy"`
	IPAddress            types.String `tfsdk:"ip_address"`
	ProxyPort            types.Int64  `tfsdk:"proxy_port"`
	EnableAuthentication types.Bool   `tfsdk:"enable_authentication"`
	Username             types.String `tfsdk:"username"`
	Password             types.String `tfsdk:"password"`
}
