package models


type NetworkAdapterSetting struct {
	OdataContext      string            `json:"@odata.context"`
	OdataType         string            `json:"@odata.type"`
	OdataID           string            `json:"@odata.id"`
	InterfaceName     string            `json:"InterfaceName"`
	ProfileName       string            `json:"ProfileName"`
	EnableNIC         bool              `json:"EnableNIC"`
	Ipv4Configuration Ipv4Configuration `json:"Ipv4Configuration"`
	Ipv6Configuration Ipv6Configuration `json:"Ipv6Configuration"`
	ManagementVLAN    ManagementVLAN    `json:"ManagementVLAN"`
	DNSConfiguration  DNSConfiguration  `json:"DnsConfiguration"`
	CurrentSettings   CurrentSettings   `json:"CurrentSettings"`
	Delay             int               `json:"Delay"`
	PrimaryInterface  bool              `json:"PrimaryInterface"`
}
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
type ManagementVLAN struct {
	EnableVLAN bool `json:"EnableVLAN"`
	ID         int  `json:"Id"`
}
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
type Ipv4Settings struct {
	Enable                   bool     `json:"Enable"`
	EnableDhcp               bool     `json:"EnableDhcp"`
	CurrentIPAddress         []string `json:"CurrentIPAddress"`
	CurrentSubnetMask        string   `json:"CurrentSubnetMask"`
	CurrentGateway           string   `json:"CurrentGateway"`
	UseDHCPForDNSServerNames bool     `json:"UseDHCPForDNSServerNames"`
	Ipv4DNS                  []string `json:"Ipv4Dns"`
}
type Ipv6Settings struct {
	Enable                   bool   `json:"Enable"`
	EnableAutoConfiguration  bool   `json:"EnableAutoConfiguration"`
	CurrentIPAddress         []any  `json:"CurrentIPAddress"`
	CurrentGateway           string `json:"CurrentGateway"`
	CurrentLinkLocalAddress  string `json:"CurrentLinkLocalAddress"`
	UseDHCPForDNSServerNames bool   `json:"UseDHCPForDNSServerNames"`
	Ipv6DNS                  []any  `json:"Ipv6Dns"`
}
type DNSSetting struct {
	DNSFQDName    string `json:"DnsFQDName"`
	DNSDomainName string `json:"DnsDomainName"`
}
type CurrentSettings struct {
	Ipv4Settings Ipv4Settings `json:"Ipv4Settings"`
	Ipv6Settings Ipv6Settings `json:"Ipv6Settings"`
	DNSSetting   DNSSetting   `json:"DnsSetting"`
}


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

type NetworkSessions struct {
	OdataContext string  `json:"@odata.context"`
	OdataCount   int     `json:"@odata.count"`
	Value        []SessionInfo `json:"value"`
}
type SessionInfo struct {
	OdataType                  string `json:"@odata.type"`
	SessionType                string `json:"SessionType"`
	MaxSessions                int    `json:"MaxSessions"`
	SessionTimeout             int    `json:"SessionTimeout"`
	MinSessionTimeout          int    `json:"MinSessionTimeout"`
	MaxSessionTimeout          int    `json:"MaxSessionTimeout"`
	MinSessionsAllowed         int    `json:"MinSessionsAllowed"`
	MaxSessionsAllowed         int    `json:"MaxSessionsAllowed"`
	MaxSessionsConfigurable    bool   `json:"MaxSessionsConfigurable"`
	SessionTimeoutConfigurable bool   `json:"SessionTimeoutConfigurable"`
}

type UpdateNetworkSessions []struct {
	SessionType    string `json:"SessionType"`
	MaxSessions    int    `json:"MaxSessions"`
	SessionTimeout int    `json:"SessionTimeout"`
}

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

type TimeZones struct {
	OdataContext string  `json:"@odata.context"`
	OdataCount   int     `json:"@odata.count"`
	Value        []TimeZone `json:"value"`
}
type TimeZone struct {
	OdataType        string `json:"@odata.type"`
	Utcoffsetminutes int    `json:"Utcoffsetminutes"`
	ID               string `json:"Id"`
	Name             string `json:"Name"`
}

type TimeConfigPayload struct {
	TimeZone             string `json:"TimeZone"`
	EnableNTP            bool   `json:"EnableNTP"`
	PrimaryNTPAddress    any    `json:"PrimaryNTPAddress"`
	SecondaryNTPAddress1 any    `json:"SecondaryNTPAddress1"`
	SecondaryNTPAddress2 any    `json:"SecondaryNTPAddress2"`
	SystemTime           string `json:"SystemTime"`
}

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

type ProxyConfiguration struct {
	IPAddress            string `json:"IpAddress"`
	PortNumber           int    `json:"PortNumber"`
	Username             string `json:"Username"`
	Password             any    `json:"Password"`
	EnableAuthentication bool   `json:"EnableAuthentication"`
	EnableProxy          bool   `json:"EnableProxy"`
	SslCheckDisabled     bool   `json:"SslCheckDisabled"`
	ProxyExclusionList   any    `json:"ProxyExclusionList"`
}

type PayloadProxyConfiguration struct {
	IPAddress            string `json:"IpAddress"`
	PortNumber           int    `json:"PortNumber"`
	EnableAuthentication bool   `json:"EnableAuthentication"`
	EnableProxy          bool   `json:"EnableProxy"`
	Username             string `json:"Username"`
	Password             string `json:"Password"`
}