# v1.2.3

- Addresses Github Issues: #152, #126, #69, #68
- Addresses security vulnerablilites

# v1.2.2

# Release Summary

The release addresses security vulnerabilities.

All existing resources and datasources are qualified against OME v4.3.0.

## Bug Fixes

- ome_discovery in perpetual drift #135
- discovery_config_targets.network_address_detail can't handle unknown in OME Discovery Resource #136

# v1.2.1

# Release Summary

The release addresses security vulnerabilities.

All existing resources and datasources are qualified against OME v4.1.0.

# v1.2.0

# Release Summary

The release supports resources and data sources mentioned in the Features section for Dell OME.

## Features

- List of new Resources and supported operations in Terraform Provider for Dell OME.

  * Firmware Catalog
  * Firmware Baselines
  
- List of new DataSources and supported operations in Terraform Provider for Dell OME.

  * Firmware Repository
  * Firmware Baseline Compliance Report
  * Firmware Catalog
  * Device Compliance Report

## Enhancements

N/A

## Bug Fixes

N/A


# v1.1.1

# Release Summary

The release address security vulnerability reported in gRPC library.

Implemented functionality derived from the netdata library.

# v1.1.0

## Release Summary

The release supports resources and data sources mentioned in the Features section for Dell OME.

## Features

- List of new Resources and supported operations in Terraform Provider for Dell OME.

  * User Resource - Create, Read, Update, Delete, Import
  * Static Group Resource -  Create, Read, Update, Delete, Import
  * Discovery Resource -  Create, Read, Update, Delete, Import
  * Devices Resource -  Create(Implicit Import), Read, Update, Delete, Import
  * Device Action Resource -  Create, Read, Update, Delete
  * Application CSR Resource -  Create, Read, Update, Delete
  * Application Certificate Resource -  Create, Read, Update, Delete
  * Appliance Network Resource -  Create (Implicit Import), Read, Update, Delete
- List of new DataSources and supported operations in Terraform Provider for Dell OME.

  * Device Datasource - Read
  * Application Certificate Datasource - Read

## Enhancements

Group Devices Info Datasource was enhanced to add field `device_groups` which contains information about groups.

## Bug Fixes

N/A

# v1.0.0

## Release Summary

First major release for terraform provider for Dell OpenManage Enterprise (OME).

## Features

Migrated to terraform-provider-framework v1.1 but no new feature added.

## Enhancements

N/A

## Bug Fixes

N/A

# v1.0.0-beta

## Release Summary

The release supports the resources and data sources to manage baseline and remediation of baseline devices in OpenManage Enterprise (OME).

## Features

Configuring the OME Baselines for configuration compliance.

### Data Sources

* `ome_configration_report_info` to list the compliance configuration report of a baseline from OME.

### Resources

* `ome_configuration_baseline` for managing configuration baselines on OME.
* `ome_configuration_compliance` for managing configuration baseline remediation on OME.

### Others

N/A

## Enhancements

* `ome_template` enhanced to support template creation from the XML file.

## Bug Fixes

* `ome_template` fixed template import and clone, where template name matches with multiple templates when using `eq` condition as a filter in an API.

# v1.0.0-alpha

## Release Summary

The release supports a terraform client to query OpenManage Enterprise (OME) and the resources and data sources to manage templates and deployment of template in PowerEdge Servers.

## Features

### Data Sources

* `ome_template_info` to list the template details from the OME.
* `ome_groupdevices_info` to list all the devices in the group from OME.
* `ome_vlannetworks_info` to list the vlan networks from OME.

### Resources

* `ome_template` for managing template(deployment and compliance) on OME.
* `ome_deployment` for managing template deployment on OME.

### Others

N/A

## Enhancements

N/A

## Bug Fixes

N/A
