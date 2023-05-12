# create a template with reference device id.
resource "ome_template" "template_1" {
  name         = "template_1"
  refdevice_id = 10001
}

# create a template with reference device servicetag.
resource "ome_template" "template_2" {
  name                 = "template_2"
  refdevice_servicetag = "MXL1234"
}

# create a template with fqdds as NIC.
resource "ome_template" "template_3" {
  name            = "template_3"
  refdevice_id    = 10001
  fqdds           = "NIC"
}

# used to fetch vlan network data
data "ome_vlannetworks_info" "vlans" {
}


data "ome_template_info" "template_data" {
  name = "template_4"
}

#use locals to fetch vlan network ID from vlan name for updating vlan template attributes.
locals {
  vlan_network_map = {for vlan_network in  data.ome_vlannetworks_info.vlans.vlan_networks : vlan_network.name => vlan_network.vlan_id}
}

#use locals to modify the attributes required for updating a template for assigning identity pool.
locals {
  attributes_value = tomap({
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy": "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy": "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered": "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered": "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 IOIDOpt Enable" : "Enabled"

  })
  attributes_is_ignored = tomap({
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy": false
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy": false
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered": false
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered": false
    "iDRAC,IO Identity Optimization,IOIDOpt 1 IOIDOpt Enable" : false

  })
  
  template_attributes = data.ome_template_info.template_data.attributes != null ? [
    for attr in data.ome_template_info.template_data.attributes : tomap({
      attribute_id = attr.attribute_id
      is_ignored   = lookup(local.attributes_is_ignored, attr.display_name, attr.is_ignored)
      display_name = attr.display_name
      value        = lookup(local.attributes_value, attr.display_name, attr.value)
  })] : null
}

# create a template and uncomment the below in a resouce to update the template to attach the identity pool and vlan.
# identity pool and vlan attributes are applicable as part of the update operation, please uncomment it to update after the create is finished.
resource "ome_template" "template_4" {
  name                 = "template_4"
  refdevice_servicetag = "MXL1234"
  # attributes = local.template_attributes
  # identity_pool_name   = "IO1"
  # vlan = {
  #     propogate_vlan     = true
  #     bonding_technology = "NoTeaming"
  #     vlan_attributes = [
  #       {
  #         untagged_network = lookup(local.vlan_network_map, "VLAN1", 0)
  #         tagged_networks  = [0]
  #         is_nic_bonded    = false
  #         port             = 1
  #         nic_identifier   = "NIC in Mezzanine 1A"
  #       },
  #       {
  #         untagged_network = 0
  #         tagged_networks  = [lookup(local.vlan_network_map, "VLAN1", 0), lookup(local.vlan_network_map, "VLAN2", 0), lookup(local.vlan_network_map, "VLAN3", 0)]
  #         is_nic_bonded    = false
  #         port             = 1
  #         nic_identifier   = "NIC in Mezzanine 1B"
  #       },
  #     ]
  #   }
}


# get the template details 
data "ome_template_info" "template_data1" {
  name = "template_5"
}


#use locals to modify the attributes required for updating a template using attribute ids.
locals {
  attributes_map = tomap({
    2740260: "One Way"
    2743100: "Disabled"
  })
  
  template_attributes = data.ome_template_info.template_data1.attributes != null ? [
    for attr in data.ome_template_info.template_data1.attributes : tomap({
      attribute_id = attr.attribute_id
      is_ignored   = attr.is_ignored
      display_name = attr.display_name
      value        = lookup(local.attributes_map, attr.attribute_id, attr.value)
  })] : null
}

#use locals to modify the attributes required for updating a template using display name.
locals {
  attributes_map = tomap({
    "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String": "IST"
    "System,Server Topology,ServerTopology 1 Aisle Name": "Aisle-123"
    "iDRAC,User Domain,UserDomain 1 User Domain Name": "TestDomain1"
  })
  
  template_attributes = data.ome_template_info.template_data1.attributes != null ? [
    for attr in data.ome_template_info.template_data1.attributes : tomap({
      attribute_id = attr.attribute_id
      is_ignored   = attr.is_ignored
      display_name = attr.display_name
      value        = lookup(local.attributes_map, attr.display_name, attr.value)
  })] : null
}

# create a template and update the attributes of the template
# attributes are only updatable and is not applicable during create operation.
# attributes existing list can be fetched from a template with a datasource - ome_template_info as defined above.
# modified attributes list should be passed to update the attributes for a template
resource "ome_template" "template_5" {
  name                 = "template_5"
  refdevice_servicetag = "MXL1234"
  attributes           = local.template_attributes
}


# create multiple templates with template names and reference devices.
resource "ome_template" "templates" {
  count                = length(var.ome_template_names)
  name                 = var.ome_template_names[count.index]
  refdevice_servicetag = var.ome_template_servicetags[count.index]
}

# Clone a deploy template to create compliance template.
resource "ome_template" "template_6" {
  name             = "template_6"
  reftemplate_name = "template_5"
  view_type        = "Compliance"
}

# Create a deployment template from a XML.
resource "ome_template" "template_7" {
	name = "template_7"
	content = file("../testdata/test_acc_template.xml")
}

# Create a compliance template from a XML.
resource "ome_template" "template_8" {
	name = "template_8"
	content = file("../testdata/test_acc_template.xml")
	view_type = "Compliance"
}