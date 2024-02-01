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

resource "ome_application_csr" "csr1" {
  specs = {
    distinguished_name      = "localhost"
    department_name         = "Terraform Server Solutions"
    business_name           = "Dell Terraform"
    locality                = "RedRock"
    state                   = "Texas"
    country                 = "US"
    email                   = "noreply@gmail.com"
    subject_alternate_names = ["dell.com", "amer.dell.com", "10.36.0.124", "2607:f2b1:f006:127::10"]
  }
}

resource "local_file" "csr_file" {
  content  = ome_application_csr.csr
  filename = "foo.pem"
}
