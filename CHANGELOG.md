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
