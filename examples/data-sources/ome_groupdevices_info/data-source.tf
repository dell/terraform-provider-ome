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

# Get Group Info, Deviceids and servicetags of all devices that belong to a specified list of groups
data "ome_groupdevices_info" "gd" {
  device_group_names = ["WINDOWS", "WINDOWS-10"]
}

output "out" {
  value = {
    "common_device_ids"            = data.ome_groupdevices_info.gd.device_ids,
    "common_device_scvtags"        = data.ome_groupdevices_info.gd.device_servicetags,
    "group_windows"                = data.ome_groupdevices_info.gd.device_groups["WINDOWS"],
    "group_windows_subgroup_names" = data.ome_groupdevices_info.gd.device_groups["WINDOWS"].sub_groups[*].name,
    "group_windows_device_ids"     = data.ome_groupdevices_info.gd.device_groups["WINDOWS"].devices[*].id,
  }
}

# Get Sub Group Info of all specified groups
locals {
  # gd_names_with_non_zero_children is the subset of the specified list of group names where the group
  # has one or more child groups.
  # In the below code, we iterate over the specified list of groups and check if the sub_groups
  # field for that group in the ome_groupdevices_info data source has any entries. If it does, then that name is added
  # to a list. Then we use the to_set function to convert the list of names to a set to ensure that there is no 
  # duplication of names. Then that set is assigned to gd_names_with_non_zero_children.
  gd_names_with_non_zero_children = toset([for i in data.ome_groupdevices_info.gd.device_group_names : i if
  length(data.ome_groupdevices_info.gd.device_groups[i].sub_groups) > 0])
}

data "ome_groupdevices_info" "gd_children" {
  id                 = "1"
  device_group_names = data.ome_groupdevices_info.gd.device_groups[each.key].sub_groups[*].name
  for_each           = local.gd_names_with_non_zero_children
}
