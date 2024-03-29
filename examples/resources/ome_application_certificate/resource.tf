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

data "local_file" "cert" {
  # this is the path to the certificate that we want to upload.
  filename = "./signed-certificate.pem.crt"
}

resource "ome_application_certificate" "ome_cert" {
  certificate_base64 = data.local_file.cert.content_base64
}
