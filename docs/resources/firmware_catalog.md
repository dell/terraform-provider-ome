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

title: "ome_firmware_catalog resource"
linkTitle: "ome_firmware_catalog"
page_title: "ome_firmware_catalog Resource - terraform-provider-ome"
subcategory: ""
description: |-
  This terraform resource is used to manage firmware catalogs entity on OME.We can Create, Update and Delete OME firmware catalogs using this resource. We can also do an 'Import' an existing 'firmware catalog' from OME .
---

# ome_firmware_catalog (Resource)

This terraform resource is used to manage firmware catalogs entity on OME.We can Create, Update and Delete OME firmware catalogs using this resource. We can also do an 'Import' an existing 'firmware catalog' from OME .

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

# Resource to manage a new firmware catalog
resource "ome_firmware_catalog" "firmware_catalog_example" {
  # Name of the catalog required
  name = "example_catalog_1"
  
  # Catalog Update Type required.
  # Sets to Manual or Automatic on schedule catalog updates of the catalog. 
  # Defaults to manual.
  catalog_update_type = "Automatic"
  
  # Share type required.
  # Sets the different types of shares (DELL_ONLINE, NFS, CIFS, HTTP, HTTPS)
  # Defaults to DELL_ONLINE
  share_type = "HTTPS"

  # Catalog file path, required for share types (NFS, CIFS, HTTP, HTTPS)
  # Start directory path without leading '/' and use alphanumeric characters. 
  catalog_file_path = "catalogs/example_catalog_1.xml"

  # Share Address required for share types (NFS, CIFS, HTTP, HTTPS)
  # Must be a valid ipv4 (x.x.x.x), ipv6(xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx), or fqdn(example.com)
  # And include the protocol prefix ie (https://)
  share_address = "https://1.2.2.1"
 
  # Catalog refresh schedule, Required for catalog_update_type Automatic.
  # Sets the frequency of the catalog refresh.
  # Will be ignored if catalog_update_type is set to manual.
  catalog_refresh_schedule = {
    # Sets to (Weekly or Daily)
    cadence = "Weekly"
    # Sets the day of the week (Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)
    day_of_the_week = "Wednesday"
    # Sets the hour of the day (1-12)
    time_of_day = "6"
    # Sets (AM or PM)
    am_pm = "PM"
  }
  
  # Domain optional value for the share (CIFS), for other share types this will be ignored
  domain = "example"

  # Share user required value for the share (CIFS), optional value for the share (HTTPS)
  share_user = "example-user"

  # Share password required value for the share (CIFS), optional value for the share (HTTPS)
  share_password = "example-pass"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the new catalog.

### Optional

- `catalog_file_path` (String) Catalog File Path. Path on the share to gather catalog data. This field is required for share_types (NFS, CIFS, HTTP, HTTPS)
- `catalog_refresh_schedule` (Attributes) Catalog Refresh Schedule, when using automatic catalog update the schedule is required for cadence of the update. If catalog_update_type is set to manual, this field is ignored. (see [below for nested schema](#nestedatt--catalog_refresh_schedule))
- `catalog_update_type` (String) Catalog Update Type. Sets the frequency of catalog updates. Defaults to Manual. If set to automatic, the catalog_refresh_schedule field will need to be set. Options are (Manual, Automatic).
- `domain` (String) Domain. The domain for the catalog. This field is optional and only used for share_types (CIFS).
- `share_address` (String) Share Address. Gives the Ipv4, Ipv6, or FQDN of the share. This field is required for share_types (NFS, CIFS, HTTP, HTTPS)
- `share_password` (String, Sensitive) Share Password. The password related to the share address. This field is required for share_types (CIFS, HTTPS)
- `share_type` (String) Share Type, the type of share the catalog will pull from, Defaults to Dell. The different options will have different required fields to work properly. Options are (DELL, NFS, CIFS, HTTP, HTTPS).
- `share_user` (String) Share User. The username related to the share address. This field is required for share_types (CIFS, HTTPS).

### Read-Only

- `associated_baselines` (Attributes List) Associated Baselines. (see [below for nested schema](#nestedatt--associated_baselines))
- `baseline_location` (String) Baseline Location.
- `bundles_count` (Number) Bundles Count.
- `create_date` (String) Create Date.
- `filename` (String) Filename.
- `id` (Number) id.
- `last_update` (String) Last Update.
- `manifest_identifier` (String) Manifest Identifier.
- `manifest_version` (String) Manifest Version.
- `next_update` (String) Next Update.
- `owner_id` (Number) Owner ID.
- `predcessor_identifier` (String) Predcessor Identifier.
- `release_identifier` (String) Release Identifier.
- `repository` (Object) Repository. (see [below for nested schema](#nestedatt--repository))
- `source_path` (String) Source path.
- `status` (String) Status.
- `task_id` (Number) Task ID.

<a id="nestedatt--catalog_refresh_schedule"></a>
### Nested Schema for `catalog_refresh_schedule`

Optional:

- `am_pm` (String) AM/PM for the schedule. Options are (AM, PM).
- `cadence` (String) Cadence. Options are(Weekly, Daily).
- `day_of_the_week` (String) Day of the Week, only useful for weekly schedules. Options are(Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday).
- `time_of_day` (Number) Time of Day for the schedule in hour increments. Options are (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12).


<a id="nestedatt--associated_baselines"></a>
### Nested Schema for `associated_baselines`

Read-Only:

- `baseline_id` (Number) Baseline ID.
- `baseline_name` (String) Baseline Name.


<a id="nestedatt--repository"></a>
### Nested Schema for `repository`

Read-Only:

- `backup_existing_catalog` (Boolean)
- `check_certificate` (Boolean)
- `description` (String)
- `domain_name` (String)
- `editable` (Boolean)
- `id` (Number)
- `name` (String)
- `repository_type` (String)
- `source` (String)
- `username` (String)

## Import

Import is supported using the following syntax:

```shell
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://mozilla.org/MPL/2.0/


# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The command is
# terraform import ome_firmware_catalog.cat_1 <id>
# Example:
terraform import ome_firmware_catalog.cat_1 1
# after running this command, populate the name field in the config file to start managing this resource
```