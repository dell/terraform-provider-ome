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

# Gets the firmware baseline compliance report
data "ome_firmware_baseline_compliance_report" "report1" {
  baseline_name = "tfacc_baseline_dell_1"

  # Supported filter keys are: DeviceName, DeviceModel, ServiceTag
  # Only one filter key/value can be used at a time
  # filter {
			# key = "DeviceModel"
			# value = "Valid Name"
	# }
}

output "all" {
  value = data.ome_firmware_baseline_compliance_report.report1
}