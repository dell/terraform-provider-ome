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

# Resource to manage a new firmware baseline
resource "ome_firmware_baseline" "firmware_baseline" {
  // Required Fields
  # Name of the catalog
  catalog_name = "tfacc_catalog_dell_online_1"
  # Name of the Baseline
  name = "baselinetest"
  
  // Only one of the following fields (device_names, group_names , device_service_tags) is required
  # List of the Device names to associate with the firmware baseline.
  device_names = ["10.2.2.1"]
  
  # List of the Group names to associate with the firmware baseline.
  # group_names = ["HCI Appliances","Hyper-V Servers"]
  
  # List of the Device service tags to associate with the firmware baseline.
  # device_service_tags = ["HRPB0M3"]
  
  // Optional Fields
  // This must always be set to true. The size of the DUP files used is 64 bits."
  #is_64_bit = true

  // Filters applicable updates where no reboot is required during create baseline for firmware updates. This field is set to false by default.
  #filter_no_reboot_required = true
  
  # Description of the firmware baseline 
  description = "test baseline"
}
