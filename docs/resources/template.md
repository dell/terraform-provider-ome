---
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://mozilla.org/MPL/2.0/
# 
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "ome_template resource"
linkTitle: "ome_template"
page_title: "ome_template Resource - terraform-provider-ome"
subcategory: ""
description: |-
  Resource for managing template on OpenManage Enterprise.
---

# ome_template (Resource)

Resource for managing template on OpenManage Enterprise.

~> **Note:** Exactly one of `reftemplate_name`, `refdevice_id`, `refdevice_servicetag` and `content` are required.

## Example Usage

```terraform
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
  name         = "template_3"
  refdevice_id = 10001
  fqdds        = "NIC"
}

# used to fetch vlan network data
data "ome_vlannetworks_info" "vlans" {
}


data "ome_template_info" "template_data" {
  name = "template_4"
}

#use locals to fetch vlan network ID from vlan name for updating vlan template attributes.
locals {
  vlan_network_map = { for vlan_network in data.ome_vlannetworks_info.vlans.vlan_networks : vlan_network.name => vlan_network.vlan_id }
}

#use locals to modify the attributes required for updating a template for assigning identity pool.
locals {
  attributes_value = tomap({
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy" : "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy" : "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered" : "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered" : "WarmReset, ColdReset, ACPowerLoss"
    "iDRAC,IO Identity Optimization,IOIDOpt 1 IOIDOpt Enable" : "Enabled"

  })
  attributes_is_ignored = tomap({
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Initiator Persistence Policy" : false
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Storage Target Persistence Policy" : false
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Auxiliary Powered" : false
    "iDRAC,IO Identity Optimization,IOIDOpt 1 Virtual Address Persistence Policy Non Auxiliary Powered" : false
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
    2740260 : "One Way"
    2743100 : "Disabled"
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
    "iDRAC,Time Zone Configuration Information,Time 1 Time Zone String" : "IST"
    "System,Server Topology,ServerTopology 1 Aisle Name" : "Aisle-123"
    "iDRAC,User Domain,UserDomain 1 User Domain Name" : "TestDomain1"
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
  name    = "template_7"
  content = file("../testdata/test_acc_template.xml")
}

# Create a compliance template from a XML.
resource "ome_template" "template_8" {
  name      = "template_8"
  content   = file("../testdata/test_acc_template.xml")
  view_type = "Compliance"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the template resource.

### Optional

- `attributes` (List of Object) List of attributes associated with a template. This field is ignored while creating a template. (see [below for nested schema](#nestedatt--attributes))
- `content` (String) The XML content of template. Cannot be updated.
- `description` (String) Description of the template
- `device_type` (String) OME template device type, supported types are Server, Chassis. Cannot be updated and is applicable only for importing xml. Valid values are `Server` and `Chassis`. Default value is `Server`.
- `fqdds` (String) Comma seperated values of components from a specified server. Valid values are `iDRAC`, `System`, `BIOS`, `NIC`, `LifeCycleController`, `RAID`, `EventFilters` and `All`. Default value is `All`. Cannot be updated.
- `identity_pool_name` (String) Identity Pool name to be attached with template.
- `job_retry_count` (Number) Number of times the job has to be polled to get the final status of the resource. Default value is `5`.
- `refdevice_id` (Number) Target device id from which the template needs to be created. Cannot be updated.
- `refdevice_servicetag` (String) Target device servicetag from which the template needs to be created. Cannot be updated.
- `reftemplate_name` (String) Reference Template name from which the template needs to be cloned. Cannot be updated.
- `sleep_interval` (Number) Sleep time interval for job polling in seconds. Default value is `30`.
- `view_type` (String) OME template view type. Valid values are `Deployment` and `Compliance`. Default value is `Deployment`. Cannot be updated.
- `vlan` (Object) VLAN details to be attached with template. (see [below for nested schema](#nestedatt--vlan))

### Read-Only

- `id` (String) ID of the template resource.
- `identity_pool_id` (Number) ID of the Identity Pool attached with template.
- `view_type_id` (Number) OME template view type id.

<a id="nestedatt--attributes"></a>
### Nested Schema for `attributes`

Optional:

- `attribute_id` (Number)
- `display_name` (String)
- `is_ignored` (Boolean)
- `value` (String)


<a id="nestedatt--vlan"></a>
### Nested Schema for `vlan`

Optional:

- `bonding_technology` (String)
- `propogate_vlan` (Boolean)
- `vlan_attributes` (List of Object) (see [below for nested schema](#nestedobjatt--vlan--vlan_attributes))

<a id="nestedobjatt--vlan--vlan_attributes"></a>
### Nested Schema for `vlan.vlan_attributes`

Optional:

- `is_nic_bonded` (Boolean)
- `nic_identifier` (String)
- `port` (Number)
- `tagged_networks` (Set of Number)
- `untagged_network` (Number)

## Import

Import is supported using the following syntax:

```shell
terraform import ome_template.citctest "<existing_template_name>"
```