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

# Manage baseline using Device Servicetags
resource "ome_configuration_baseline" "baseline_name" {
  baseline_name      = "Baseline Name"
  device_servicetags = ["MXL1234", "MXL1235"]
}

# Create Baseline using device ids
resource "ome_configuration_baseline" "baseline1" {
  baseline_name   = "baseline1"
  ref_template_id = 745
  device_ids      = [10001, 10002]
  description     = "baseline description"
}


# Create Baseline using device servicetag with daily notification scheduled 
resource "ome_configuration_baseline" "baseline2" {
  baseline_name      = "baseline2"
  ref_template_id    = 745
  device_servicetags = ["MXL1234", "MXL1235"]
  description        = "baseline description"
  schedule           = true
  notify_on_schedule = true
  email_addresses    = ["test@testmail.com"]
  cron               = "0 30 11 * * ? *"
  output_format      = "csv"
}


# Create Baseline using device ids with daily notification on status changing to non-compliant 
resource "ome_configuration_baseline" "baseline3" {
  baseline_name   = "baseline3"
  ref_template_id = 745
  device_ids      = [10001, 10002]
  description     = "baseline description"
  schedule        = true
  email_addresses = ["test@testmail.com"]
  output_format   = "pdf"
}
