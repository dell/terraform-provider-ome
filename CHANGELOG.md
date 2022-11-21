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