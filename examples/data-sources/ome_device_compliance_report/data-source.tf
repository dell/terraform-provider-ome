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

// Get the details of the device compliance report from device id, service tag or group name
data "ome_device_compliance_report" "device_compliance_report_data" {
  // This will get the device compliance report for devices based on the filters used. 
  // Only one of the below filters can be used at a time.
  # device_ids = [10102]
  # device_service_tags = ["HRPB0M3"]
  # device_group_names = ["Servers"]
}

output "device_compliance_report_data" {
  value = data.ome_device_compliance_report.device_compliance_report_data
}