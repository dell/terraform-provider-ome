# Terraform provider for OpenManage Enterprise

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.1%20adopted-ff69b4.svg)](https://github.com/dell/terraform-provider-ome/blob/main/about/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/github/license/dell/terraform-provider-ome)](https://github.com/dell/terraform-provider-ome/blob/main/LICENSE)
[![Go version](https://img.shields.io/badge/go-1.19+-blue.svg)](https://go.dev/dl/)
[![Terraform version](https://img.shields.io/badge/terraform-1.0+-blue.svg)](https://www.terraform.io/downloads)
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

## Supported Platforms
  * Dell OpenManage Enterprise versions 3.9.0 (Build 55) and above.

## Prerequisites
  * [Terraform >= 1.3.2](https://www.terraform.io)
  * Go >= 1.19

## List of DataSources in Terraform Provider for Dell OME
  * Configuration Report
  * Device Groups
  * Template
  * VLAN Networks

## List of Resources in Terraform Provider for Dell OME
  * Configuration Baseline
  * Configuration Compliance
  * Deployment
  * Template

## Installation
Install Terraform Provider for OpenManage Enterprise from terraform registry by adding the following block
```terraform
terraform {
  required_providers {
    ome = {
      version = "1.0.0"
      source  = "dell/ome"
    }
  }
}
````
For adding resources, please refer [examples](https://github.com/dell/terraform-provider-ome/blob/main/docs)

## About
Terraform Provider for OpenManage Enterprise is 100% open source and community-driven. All components are available under [MPL-2.0 license](https://www.mozilla.org/en-US/MPL/2.0/) on GitHub.

## Documentation
For more detailed information on the provider, please refer to [Dell Terraform Providers Documentation](https://dell.github.io/terraform-docs/).
