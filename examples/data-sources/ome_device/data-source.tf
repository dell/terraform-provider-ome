/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://mozilla.org/MPL/2.0/
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

# get device by ids
data "ome_device" "devi" {
  filters = {
    ids = [1001, 1002, 1003]
  }
}

# get device by their network
data "ome_device" "devn" {
  filters = {
    ip_expressions = [
      "10.10.10.10",
      "10.36.0.0-192.36.0.255",
      "fe80::ffff:ffff:ffff:ffff",
      "fe80::ffff:192.0.2.0/125",
      "fe80::ffff:ffff:ffff:1111-fe80::ffff:ffff:ffff:ffff"
    ]
  }
}

# get device by service tags
data "ome_device" "devs" {
  filters = {
    device_service_tags = ["CZNF1T2", "CZMC1T2"]
  }
}

# get device by valid OME OData filter query
data "ome_device" "devf" {
  filters = {
    filter_expression = "Model eq 'PowerEdge MX840c'"
  }
}

# get all devices in the CIDR "10.10.10.10/26" with model PowerEdge MX840c

data "ome_device" "devs" {
  filters = {
    ip_expressions    = ["10.10.10.10/26"]
    filter_expression = "Model eq 'PowerEdge MX840c'"
  }
}

# get device with inventory
# to get inventory of a device, query only one device per datasource

data "ome_device" "dev_invent_full" {
  filters = {
    ip_expressions = ["10.10.10.10"]
  }
}

data "ome_device" "dev_invent" {
  filters = {
    ip_expressions = ["10.10.10.10"]
  }
  inventory_types = ["serverNetworkInterfaces", "serverArrayDisks"]
}

output "dev_inv_out" {
  value = {
    full                    = data.ome_device.dev_invent_full.devices.detailed_inventory
    serverNetworkInterfaces = data.ome_device.dev_invent.devices.detailed_inventory.nics
    serverArrayDisks        = data.ome_device.dev_invent.devices.detailed_inventory.disks
  }
}
