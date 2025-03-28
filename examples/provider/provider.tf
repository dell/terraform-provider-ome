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

terraform {
  required_providers {
    ome = {
      version = "1.2.2"
      source  = "registry.terraform.io/dell/ome"
    }
  }
}

provider "ome" {
  username = "username"
  password = "password"
  host     = "yourhost.host.com"
  timeout  = 30
  port     = 443
  protocol = "https"
  skipssl  = false

  ## Can also be set using environment variables
  ## If environment variables are set it will override this configuration
  ## Example environment variables
  # OME_USERNAME="username"
  # OME_PASSWORD="password"
  # OME_HOST="yourhost.host.com"
  # OME_PORT="443"
  # OME_SKIP_SSL="true"
  # OME_TIMEOUT="30"
  # OME_PROTOCOL="https"
}

# creating baseline from a reference device and making other devices complaint with that baseline. 
# ------------------------------------------------------------------------------------------------
resource "ome_template" "template-1" {
  name            = "template-compliance-1"
  view_type       = "Compliance"
  refdevice_id    = var.DeviceIDRef
  fqdds           = "EventFilters"
  description     = "This is server template"
  job_retry_count = 10
  sleep_interval  = 60
}

resource "ome_configuration_baseline" "baseline-1" {
  depends_on         = [ome_template.template-1]
  baseline_name      = "baseline-1"
  ref_template_name  = ome_template.template-1.name
  device_servicetags = [var.DeviceSvcTag1, var.DeviceSvcTag2]
  description        = "baseline description"
}

resource "ome_configuration_compliance" "baseline_remediation" {
  depends_on    = [ome_configuration_baseline.baseline-1]
  baseline_name = ome_configuration_baseline.baseline-1.baseline_name
  target_devices = [
    {
      device_service_tag = var.DeviceSvcTag1
      # when device is non-compliant, terraform will show a configuration drift at this field.
      compliance_status = "Compliant"
    }
  ]
}


# discovering devices and refreshing the inventory 
# ------------------------------------------------
resource "ome_discovery" "discovery_1" {
  name                   = local.disc_name
  schedule               = "RunNow"
  timeout                = 10
  ignore_partial_failure = false
  discovery_config_targets = [
    {
      device_type            = ["SERVER"]
      network_address_detail = [var.DeviceIP1, var.DeviceIP2]
      redfish = {
        username = var.discovery_redfish_username
        password = var.discovery_redfish_password
      }
    }
  ]
}

data "ome_device" "discovered_devices" {
  depends_on = [ome_discovery.discovery_1]
  filters = {
    ip_expressions = ome_discovery.discovery_1.discovery_config_targets[*].network_address_detail[*]
  }
}

data "ome_groupdevices_info" "root_group" {
  device_group_names = ["Static Groups"]
}

resource "ome_static_group" "discovered_group" {
  name        = "Discovered-Group"
  description = "Group of discovered devices"
  parent_id   = data.ome_groupdevices_info.ome_root.device_groups["Static Groups"].id
  device_ids  = data.ome_device.discovered_devices.devices[*].id
}

resource "ome_device_action" "inventory_refresh_action" {
  device_ids      = data.ome_device.discovered_devices.devices[*].id
  job_name        = "refresh-job"
  job_description = "r-job-desc"
  timeout         = 5
}
