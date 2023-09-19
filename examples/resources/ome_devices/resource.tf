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

# Resource to manage all devices
resource "ome_devices" "dev_list_1" {
}

# Resource to manage specific devices
resource "ome_devices" "dev_list_2" {
  devices = [
    {
      service_tag = "ABCD34"
    },
    {
      id = 2000
    },
    # removing this block will remove device with service tag `QWX321` from the OME
    {
      service_tag = "QWX321"
    }
  ]
}

