---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://mozilla.org/MPL/2.0/
# 
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "ome_configuration_report_info data source"
linkTitle: "ome_configuration_report_info"
page_title: "ome_configuration_report_info Data Source - terraform-provider-ome"
subcategory: ""
description: |-
  This Terraform DataSource is used to query compliance configuration report of a compliance template baseline data from OME. The information fetched from this data source can be used for getting the details / for further processing in resource block.
---

# ome_configuration_report_info (Data Source)

This Terraform DataSource is used to query compliance configuration report of a compliance template baseline data from OME. The information fetched from this data source can be used for getting the details / for further processing in resource block.

## Example Usage

```terraform
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

# Get configuration compliance report for a baseline
data "ome_configuration_report_info" "cr" {
  baseline_name = "BaselineName"
}
```

After the successful execution of above said block, We can see the output value by executing `terraform output` command.
Also, we can use the fetched information by the variable `data.ome_configuration_report_info.cr`

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `baseline_name` (String) Name of the Baseline.

### Optional

- `fetch_attributes` (Boolean) Fetch  device compliance attribute report.
- `id` (String) ID for baseline compliance data source.

### Read-Only

- `compliance_report_device` (Attributes List) Device compliance report. (see [below for nested schema](#nestedatt--compliance_report_device))

<a id="nestedatt--compliance_report_device"></a>
### Nested Schema for `compliance_report_device`

Read-Only:

- `compliance_status` (String) Device compliance status.
- `device_compliance_details` (String) Device compliance details.
- `device_id` (Number) Device ID
- `device_name` (String) Device Name.
- `device_servicetag` (String) Device servicetag.
- `device_type` (Number) Device type
- `inventory_time` (String) Inventory Time.
- `model` (String) Device model.
