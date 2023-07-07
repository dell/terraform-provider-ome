# /*
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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

# TODO: Get parent id from ome_groupdevice_info datasource

# TODO: Get device ids from ome_device datasource by ips

resource "ome_static_group" "linux-group" {
  name        = "Linux"
  description = "All linux servers"
  parent_id   = 1011
  device_ids  = [10056, 10057]
}