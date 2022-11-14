terraform {
  required_providers {
    ome = {
      version = "1.0.0"
      source  = "dell/ome"
    }
  }
}

provider "ome" {
  username = var.username
  password = var.password
  host     = var.host
  skipssl  = var.skipssl
}

# create a template with reference device id.
resource "ome_template" "template_1" {
  name         = "template_1"
  refdevice_id = 10001
}

# create a template with reference servicetag.
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


# create a template and attach the identity pool and vlan.
# identity pool and vlan attributes are applicable as part of the update operation, please uncomment it to update after the create is finished.
resource "ome_template" "template_4" {
  name                 = "template_4"
  refdevice_servicetag = "MXL1234"
#   identity_pool_name   = "IO1"
#   vlan = {
#     propogate_vlan     = true
#     bonding_technology = "NoTeaming"
#     vlan_attributes = [
#       {
#         untagged_network = 12001
#         tagged_networks  = [0]
#         is_nic_bonded    = false
#         port             = 1
#         nic_identifier   = "NIC in Mezzanine 1A"
#       },
#       {
#         untagged_network = 0
#         tagged_networks  = [12002]
#         is_nic_bonded    = false
#         port             = 2
#         nic_identifier   = "NIC in Mezzanine 1A"
#       }
#     ]
#   }
}


# get the template details 
data "ome_template_info" "template_dat_5" {
  name = "template_5"
}


#use locals to modify the attributes required for updating a template.
locals {
  template_attributes = data.ome_template_info.data-template-1.attributes != null ? [
    for attr in data.ome_template_info.data-template-1.attributes : tomap({
      attribute_id = attr.attribute_id
      is_ignored   = attr.is_ignored
      display_name = attr.display_name
      value        = attr.display_name == "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String" ? "IST" : attr.value
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
