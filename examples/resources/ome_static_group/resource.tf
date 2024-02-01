# /*
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# get all devices in the CIDR "10.10.10.10/26" with model PowerEdge MX840c
data "ome_device" "devs" {
  filters = {
    ip_expressions    = ["10.10.10.10/26"]
    filter_expression = "Model eq 'PowerEdge MX840c'"
  }
}

# get the root group of all static groups in OME
# we are mainly concerned with the ID of this group which we shall use to create a child group
data "ome_groupdevices_info" "ome_root" {
  device_group_names = ["Static Groups"]
}

# Create a group of all the devices fetched by the device datasource
# Its parent group will be "Static Groups"
resource "ome_static_group" "pE-slash-26" {
  name        = "Group46"
  description = "Group of all devices in the CIDR '10.10.10.10/26' with model PowerEdge MX840c"
  parent_id   = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
  device_ids  = data.ome_device.devs.devices[*].id
}
