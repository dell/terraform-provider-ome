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

# Get application certificate information of ome
data "ome_application_certificate" "cert_info" {
}

output "fetched_cert" {
  value = {
    "valid_to" : data.ome_application_certificate.cert.valid_to,
    "valid_from" : data.ome_application_certificate.cert.valid_from,
    "issued_to_distinguished_name" : data.ome_application_certificate.cert.issued_to.distinguished_name,
    "issued_by_distinguished_name" : data.ome_application_certificate.cert.issued_by.distinguished_name
  }
}
