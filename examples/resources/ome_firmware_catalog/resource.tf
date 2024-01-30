/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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