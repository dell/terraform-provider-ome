---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ome Provider"
subcategory: ""
description: |-
  The Terraform Provider for OpenManage Enterprise (OME) is a plugin for Terraform that allows the resource management of PowerEdge servers using OME
---

# ome Provider

The Terraform Provider for OpenManage Enterprise (OME) is a plugin for Terraform that allows the resource management of PowerEdge servers using OME

## Example Usage

```terraform
provider "ome" {
  username = "username"
  password = "password"
  host = "yourhost.host.com"
  skipssl = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host` (String) OpenManage Enterprise IP address or hostname.
- `password` (String, Sensitive) OpenManage Enterprise password.
- `username` (String) OpenManage Enterprise username.

### Optional

- `port` (Number) OpenManage Enterprise HTTPS port.
- `skipssl` (Boolean) Skips SSL certificate validation on OpenManage Enterprise
- `timeout` (Number) HTTPS timeout for OpenManage Enterprise client
