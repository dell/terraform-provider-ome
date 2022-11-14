package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// TemplateDataSource Schema object for data source
type TemplateDataSource struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ViewTypeID     types.Int64  `tfsdk:"view_type_id"`
	DeviceTypeID   types.Int64  `tfsdk:"device_type_id"`
	RefdeviceID    types.Int64  `tfsdk:"refdevice_id"`
	Content        types.String `tfsdk:"content"`
	Description    types.String `tfsdk:"description"`
	Attributes     types.List   `tfsdk:"attributes"`
	IdentityPoolID types.Int64  `tfsdk:"identity_pool_id"`
	Vlan           types.Object `tfsdk:"vlan"`
}

// Template Schema object
type Template struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	FQDDS               types.String `tfsdk:"fqdds"`
	ViewType            types.String `tfsdk:"view_type"`
	ViewTypeID          types.Int64  `tfsdk:"view_type_id"`
	RefdeviceServicetag types.String `tfsdk:"refdevice_servicetag"`
	RefdeviceID         types.Int64  `tfsdk:"refdevice_id"`
	ReftemplateName     types.String `tfsdk:"reftemplate_name"`
	Description         types.String `tfsdk:"description"`
	Attributes          types.List   `tfsdk:"attributes"`
	JobRetryCount       types.Int64  `tfsdk:"job_retry_count"`
	SleepInterval       types.Int64  `tfsdk:"sleep_interval"`
	IdentityPoolName    types.String `tfsdk:"identity_pool_name"`
	IdentityPoolID      types.Int64  `tfsdk:"identity_pool_id"`
	Vlan                types.Object `tfsdk:"vlan"`
}

// Attribute template attributes
type Attribute struct {
	AttributeID types.Int64  `tfsdk:"attribute_id"`
	DisplayName types.String `tfsdk:"display_name"`
	Value       types.String `tfsdk:"value"`
	IsIgnored   types.Bool   `tfsdk:"is_ignored"`
}

// Vlan for templates
type Vlan struct {
	PropogateVlan     types.Bool   `tfsdk:"propogate_vlan"`
	BondingTechnology types.String `tfsdk:"bonding_technology"`
	VlanAttributes    types.List   `tfsdk:"vlan_attributes"`
}

// VlanAttributes list of vlan attributes
type VlanAttributes struct {
	UntaggedNetwork types.Int64  `tfsdk:"untagged_network"`
	TaggedNetworks  types.List   `tfsdk:"tagged_networks"`
	IsNicBonded     types.Bool   `tfsdk:"is_nic_bonded"`
	Port            types.Int64  `tfsdk:"port"`
	NicIdentifier   types.String `tfsdk:"nic_identifier"`
}

// AttributeGroup resembles the AttributeGroup field in the response of GET AttributeDetails call
type AttributeGroup struct {
	GroupNameID        int64               `json:"GroupNameId"`
	DisplayName        string              `json:"DisplayName"`
	SubAttributeGroups []SubAttributeGroup `json:"SubAttributeGroups"`
	Attributes         []OmeAttribute      `json:"Attributes"`
}

// SubAttributeGroup resembles the SubAttributeGroup in the response of GET AttributeDetails call
type SubAttributeGroup struct {
	GroupNameID        int64               `json:"GroupNameId"`
	DisplayName        string              `json:"DisplayName"`
	SubAttributeGroups []SubAttributeGroup `json:"SubAttributeGroups"`
	Attributes         []OmeAttribute      `json:"Attributes"`
}

// OmeAttribute resembles the Attribute in the response of Get AttributeDetails call
type OmeAttribute struct {
	AttributeID int64  `json:"AttributeId"`
	DisplayName string `json:"DisplayName"`
	Value       string `json:"Value"`
	IsIgnored   bool   `json:"IsIgnored"`
}

// CreateTemplate - payload to create a template
type CreateTemplate struct {
	Fqdds          string `json:"Fqdds"`
	ViewTypeID     int64  `json:"ViewTypeId"`
	SourceDeviceID int64  `json:"SourceDeviceId"`
	Name           string `json:"Name"`
	Description    string `json:"Description"`
}

// UpdateTemplate - payload to update a template
type UpdateTemplate struct {
	ID          int64             `json:"Id"`
	Name        string            `json:"Name"`
	Attributes  []UpdateAttribute `json:"Attributes"`
	Description string            `json:"Description"`
}

// UpdateAttribute resembles the attribute in the update template payload to OME
type UpdateAttribute struct {
	ID        int64  `json:"Id"`
	Value     string `json:"Value"`
	IsIgnored bool   `json:"IsIgnored"`
}

// ViewTypes resembles the OME response of GET TemplateViewTypes
type ViewTypes struct {
	Value []ViewType `json:"value"`
}

// ViewType holds the details of template View type
type ViewType struct {
	ID          int64  `json:"Id"`
	Description string `json:"Description"`
}

// DeviceTypes resembles the OME response of GET Devices
type DeviceTypes struct {
	Value []DeviceType `json:"value"`
}

// DeviceType holds the details of the device
type DeviceType struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name"`
}

// OMETemplate resembles the GET template response from OME
type OMETemplate struct {
	ID                   int64             `json:"Id"`
	Name                 string            `json:"Name"`
	Description          string            `json:"Description"`
	Content              string            `json:"Content"`
	SourceDeviceID       int64             `json:"SourceDeviceId"`
	TypeID               int64             `json:"TypeId"`
	ViewTypeID           int64             `json:"ViewTypeId"`
	TaskID               int64             `json:"TaskId"`
	Status               int64             `json:"Status"`
	IdentityPoolID       int64             `json:"IdentityPoolId"`
	ViewsNavigationLink  string            `json:"Views@odata.navigationLink"`
	AttributeDetailsLink map[string]string `json:"AttributeDetails"`
}

// OMETemplateAttrGroups is used to parse the output of Get Templates(<templateId>)/AttributeDetails call
type OMETemplateAttrGroups struct {
	AttributeGroups []AttributeGroup `json:"AttributeGroups"`
}

// OMETemplates is used to parse the output of Get templates call with filter
type OMETemplates struct {
	Value []OMETemplate `json:"value"`
}

// IdentityPool holds the details of the IdentityPool
type IdentityPool struct {
	Name string `json:"Name"`
	ID   int64  `json:"Id"`
}

// OMEIdentityPools is used to parse the output of Get IdentityPools
type OMEIdentityPools struct {
	Value []IdentityPool `json:"value"`
}

// UpdateNetworkConfig - payload to updateNetworkConfig API (To update identityPool and Vlan attributes)
type UpdateNetworkConfig struct {
	TemplateID        int64                  `json:"TemplateId"`
	IdentityPoolID    int64                  `json:"IdentityPoolId"`
	BondingTechnology string                 `json:"BondingTechnology"`
	PropagateVLAN     bool                   `json:"PropagateVlan"`
	VLANAttributes    []PayloadVlanAttribute `json:"VlanAttributes"`
}

// OMEVlan - model that represents vlan in the resource schema in Go types.
type OMEVlan struct {
	BondingTechnology string
	PropagateVLAN     bool
	OMEVlanAttributes []OMEVlanAttribute
}

// OMEVlanAttribute represents each attribute in the list of OMEVlanAttributes in vlan in the resource schema in Go types
type OMEVlanAttribute struct {
	ComponentID   int64
	Untagged      int64
	Tagged        []int64
	IsNICBonded   bool
	Port          int64
	NicIdentifier string
}

// PayloadVlanAttribute - Vlan attribute model representing the vlan attributes in the updateNetworkConfig API payload
type PayloadVlanAttribute struct {
	ComponentID int64   `json:"ComponentId"`
	Untagged    int64   `json:"Untagged"`
	Tagged      []int64 `json:"Tagged"`
	IsNICBonded bool    `json:"IsNicBonded"`
}

// NetworkSpecificView - model that represents the output of Template(<templateId>)/Views(4)/AttributeViewDetails API call
type NetworkSpecificView struct {
	ViewID                 int64                   `json:"Id"`
	NetworkAttributeGroups []NetworkAttributeGroup `json:"AttributeGroups"`
}

// NetworkAttributeGroup - represents each AttributeGroup in the response of Template(<templateId>)/Views(4)/AttributeViewDetails API call
type NetworkAttributeGroup struct {
	GroupNameID        int64                   `json:"GroupNameId"`
	DisplayName        string                  `json:"DisplayName"`
	SubAttributeGroups []NetworkAttributeGroup `json:"SubAttributeGroups"`
	NetworkAttributes  []NetworkAttribute      `json:"Attributes"`
}

// NetworkAttribute represents each attribute in the list of Attributes in the response of Template(<templateId>)/Views(4)/AttributeViewDetails API call
type NetworkAttribute struct {
	ComponentID int64  `json:"CustomId"`
	DisplayName string `json:"DisplayName"`
	Value       string `json:"Value"`
}

// OMECloneTemplate - model used to clone template from a reference template id.
type OMECloneTemplate struct {
	SourceTemplateID int64  `json:"SourceTemplateId"`
	NewTemplateName  string `json:"NewTemplateName"`
	ViewTypeID       int64  `json:"ViewTypeId"`
}
