/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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

# Get Deviceid's and servicetags of all devices that belong to a specified list of groups
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
