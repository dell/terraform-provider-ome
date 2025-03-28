<!--
Copyright (c) 2024-2025 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Terraform provider for OpenManage Enterprise

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.1%20adopted-ff69b4.svg)](https://github.com/dell/terraform-provider-ome/blob/main/about/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/github/license/dell/terraform-provider-ome)](https://github.com/dell/terraform-provider-ome/blob/main/LICENSE)
[![Go version](https://img.shields.io/badge/go-1.20+-blue.svg)](https://go.dev/dl/)
[![Terraform version](https://img.shields.io/badge/terraform-1.4+-blue.svg)](https://www.terraform.io/downloads)
[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/dell/terraform-provider-ome?include_prereleases&label=latest&style=flat-square)](https://github.com/dell/terraform-provider-ome/releases)


The Terraform Provider for OpenManage Enterprise is a provider for terraform that allows the resource management of PowerEdge servers using OpenManage Enterprise (OME). For more details on OME, please refer to OME official webpage [here][ome-website].

For general information about Terraform, visit the [official website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io
[tf-github]: https://github.com/hashicorp/terraform
[ome-website]:  https://www.dell.com/support/kbdoc/en-in/000175879/support-for-openmanage-enterprise?lang=en



## Table of Contents

  * [Code of Conduct](https://github.com/dell/terraform-provider-ome/blob/main/about/CODE_OF_CONDUCT.md)
  * [Committer Guide](https://github.com/dell/terraform-provider-ome/blob/main/about/COMMITTER_GUIDE.md)
  * [Contributing Guide](https://github.com/dell/terraform-provider-ome/blob/main/about/CONTRIBUTING.md)
  * [Maintainers](https://github.com/dell/terraform-provider-ome/blob/main/about/MAINTAINERS.md)
  * [Support](https://github.com/dell/terraform-provider-ome/blob/main/about/SUPPORT.md)
  * [Security](https://github.com/dell/terraform-provider-ome/blob/main/about/SECURITY.md)
  * [Additional Information](https://github.com/dell/terraform-provider-ome/blob/main/about/ADDITIONAL_INFORMATION.md)
  * [Attribution](https://github.com/dell/terraform-provider-ome/blob/main/about/ATTRIBUTION.md)
  * [New to Terraform?](#new-to-terraform)

## Prerequisites
 | **Terraform Provider** | **OME Version** | **OS** | **Terraform** | **Golang** |
 |------------------------|-----------------|--------|---------------|------------|
 | v1.2.2 |  3.10.x <br> 4.0.1 <br> 4.1.0 | Ubuntu22.04 <br> RHEL9.x | 1.8.x <br> 1.9.x <br> | 1.22


## List of DataSources in Terraform Provider for Dell OME
  * Configuration Report
  * Device Groups
  * Template
  * VLAN Networks
  * Device Datasource
  * Device Compliance Report
  * Application Certificate Datasource
  * Firmware Repository
  * Firmware Baseline Compliance Report
  * Firmware Catalog
  

## List of Resources in Terraform Provider for Dell OME
  * Configuration Baseline
  * Configuration Compliance
  * Deployment
  * Template
  * User Resource
  * Static Group Resource
  * Discovery Resource
  * Devices Resource
  * Device Action Resource
  * Application CSR Resource
  * Application Certificate Resource
  * Appliance Network Resource
  * Firmware Catalog
  * Firmware Baselines

## Installation
Install Terraform Provider for OpenManage Enterprise from terraform registry by adding the following block
```terraform
terraform {
  required_providers {
    ome = {
      version = "1.2.2"
      source  = "dell/ome"
    }
  }
}

# Provider Details
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

````
For adding resources, please refer [examples](https://github.com/dell/terraform-provider-ome/blob/main/docs)

## About
Terraform Provider for OpenManage Enterprise is 100% open source and community-driven. All components are available under [MPL-2.0 license](https://www.mozilla.org/en-US/MPL/2.0/) on GitHub.

## Releasing, Maintenance and Deprecation

Terraform Provider for Dell Technnologies PowerMax follows [Semantic Versioning](https://semver.org/).

New versions will be released regularly if significant changes (bug fix or new feature) are made in the provider.

Released code versions are located on tags in the form of "vx.y.z" where x.y.z corresponds to the version number.

## Documentation
For more detailed information on the provider, please refer to [Dell Terraform Providers Documentation](https://dell.github.io/terraform-docs/).

## New to Terraform?
**Here are some helpful links to get you started if you are new to terraform before using our provider:**

- Intro to Terraform: https://developer.hashicorp.com/terraform/intro 
- Providers: https://developer.hashicorp.com/terraform/language/providers 
- Resources: https://developer.hashicorp.com/terraform/language/resources
- Datasources: https://developer.hashicorp.com/terraform/language/data-sources
- Import: https://developer.hashicorp.com/terraform/language/import
- Variables: https://developer.hashicorp.com/terraform/language/values/variables
- Modules: https://developer.hashicorp.com/terraform/language/modules
- State: https://developer.hashicorp.com/terraform/language/state
- Environment Variables: https://developer.hashicorp.com/terraform/cli/config/environment-variables 