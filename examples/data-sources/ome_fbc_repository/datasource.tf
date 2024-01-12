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

# Get details of all the firmware repositories 
data "ome_fbc_repository" "fbc-repository-all" {

}
output "fbc-repository" {
  value = data.ome_fbc_repository.fbc-repository-all
}

# Get filtered Repositories
# data "ome_fbc_repository" "fbc-repository-name-filter" {
#        # If names filter is added, it should atleast have one name.Otherwise you will see an error.
#        # If you do not want to filter and get all repositories then do not pass the name filter(Above configuration)
#         names = [
#           "tfacc_catalog_1",
#           "tfacc_catalog_dell_online_1",
#         ]
# }
# output "fbc-repository-name" {
#   value = data.ome_fbc_repository.fbc-repository-name-filter
# }