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

# remediate baseline for the specified target devices 
resource "ome_configuration_compliance" "remeditation0" {
  baseline_name = "baseline_name"
  target_devices = [
    {
      device_service_tag = "MX12345"
      compliance_status  = "Compliant"
    }
  ]
}

# remediate baseline for the specified target devices with scheduling
resource "ome_configuration_compliance" "remeditation1" {
  baseline_name = "baseline_name"
  target_devices = [
    {
      device_service_tag = "MX12345"
      compliance_status  = "Compliant"
    }
  ]
  run_later = true
  cron      = "0 00 11 14 02 ? 2032"
}

# Manage a baseline and also remediate it
# create baseline 
resource "ome_configuration_baseline" "baseline" {
  baseline_name      = var.baselinename
  ref_template_name  = "Mytemplate"
  device_servicetags = ["MX12345"]
  description        = "baseline description"
}

# create a compliance resource from above baseline
resource "ome_configuration_compliance" "remeditation" {
  baseline_name = var.baselinename
  target_devices = [
    {
      device_service_tag = "MX12345"
      compliance_status  = "Compliant"
    }
  ]
  depends_on = [
    ome_configuration_baseline.baseline
  ]
}