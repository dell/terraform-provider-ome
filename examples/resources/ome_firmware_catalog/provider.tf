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
      source  = "registry.terraform.io/dell/ome"
    }
  }
}

provider "ome" {
  username = ""
  password = ""
  host     = ""
  skipssl  = true

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
