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

# Use the device datasource to get the ids of the required devices
# Their ids shall be used to run actions on them
data "ome_device" "devs" {
  filters = {
    device_service_tags = ["CZMC1T2", "4111H63"]
  }
}

# refresh inventory of devices immediately on apply
# The resource creation will fail if the inventory refresh job fails or doesnt complete within `timeout` in minutes (here 8 minutes).
resource "ome_device_action" "code_1" {
  device_ids      = data.ome_device.devs.devices[*].id
  action          = "inventory_refresh"
  job_name        = "inventory-refresh-job"
  job_description = "Job to refresh inventory of CZMC1T2 and 4111H63 devices"
  timeout         = 8
  lifecycle {
    ignore_changes = [
      timeout,
    ]
  }
}

# refresh inventory of devices sometime in the future
# The resource creation succeeds when the job is created on OME
resource "ome_device_action" "code_2" {
  device_ids      = data.ome_device.devs.devices[*].id
  action          = "inventory_refresh"
  job_name        = "inventory-refresh-job-cron"
  job_description = "Job to refresh inventory of CZMC1T2 and 4111H63 devices in future"
  cron            = "0 * */10 * * ? *"
}

# Rerunning the same action is done by forcing recreation of the resource
# Option 1: Taint the resource 
#     https://developer.hashicorp.com/terraform/cli/commands/taint
# Option 2: From Terraform 1.5, one can use the -replace directive 
#     https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address
# Option 3 (shown below): From Terraform 1.2, one can use the replace-triggered-by lifecycle method 
#     https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#replace_triggered_by


# In below example, changing the string firmware_v1 to firmware_v2 will rerun the action
resource "terraform_data" "devices_firmware" {
  input = "firmware_v1"
}

resource "ome_device_action" "code_3" {
  device_ids      = data.ome_device.devs.devices[*].id
  action          = "inventory_refresh"
  job_name        = "inventory-refresh-job"
  job_description = "Job to refresh inventory of CZMC1T2 and 4111H63 devices when any of their firwares is upgraded"
  timeout         = 8
  lifecycle {
    ignore_changes = [
      timeout,
    ]
    # From Terraform 1.2, one can use the replace-triggered-by lifecycle method 
    # https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#replace_triggered_by
    replace_triggered_by = [
      terraform_data.devices_firmware
    ]
  }
}

