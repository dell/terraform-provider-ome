# Motivation
Template creation from a reference device is a long process. To speed up our tests, we shall create templates from xml content wherever possible.

> The `testdata` directory contains sample template files. Do not use them directly as the templates need to be in sync with `DEVICESVCTAG1` and `DEVICESVCTAG2` environment variables.

Below are a list of expected templates and how to build them prior to starting the tests.

# How to create these testdata

Set the aceptance test environment variables.

Create the following templates using terraform provider ome:

## main.tf

```terraform

variable "DEVICESVCTAG1" {
    type = string
}
variable "DEVICESVCTAG2" {
    type = string
}

resource "ome_template" "terraform-acceptance-test-1" {
    name = "BuildComplianceSvcTag1"
    refdevice_servicetag = var.DEVICESVCTAG1
    fqdds = "EventFilters"
    view_type = "Compliance"
    job_retry_count = 20
    sleep_interval = 30
}

resource "ome_template" "terraform-acceptance-test-1-full" {
    name = "BuildComplianceSvcTag1Full"
    refdevice_servicetag = var.DEVICESVCTAG1
    view_type = "Compliance"
    job_retry_count = 20
    sleep_interval = 30
}

resource "ome_template" "terraform-acceptance-test-2" {
    name = "BuildComplianceSvcTag2"
    refdevice_servicetag = var.DEVICESVCTAG2
    fqdds = "EventFilters"
    view_type = "Compliance"
    job_retry_count = 20
    sleep_interval = 30
}

resource "ome_template" "terraform-acceptance-test-deploy-1" {
    name = "BuildDeploySvcTag1"
    refdevice_servicetag = var.DEVICESVCTAG1
    fqdds = "System"
    job_retry_count = 20
    sleep_interval = 30
}

```
Run as

```sh
TF_VAR_DEVICESVCTAG1=${DEVICESVCTAG1} TF_VAR_DEVICESVCTAG2=${DEVICESVCTAG2} terraform apply
```

Then export the templates to files in one folder:

| **Template Name**          | **Test File Name**                         | **Device Service Tag** | **FQDDs**      |
|----------------------------|--------------------------------------------|------------------------|----------------|
| BuildComplianceSvcTag1     | test_acc_template_compliance_svc_tag_1.xml | DEVICESVCTAG1          | `EventFilters` |
| BuildComplianceSvcTag2     | test_acc_template_compliance_svc_tag_2.xml | DEVICESVCTAG2          | `EventFilters` |
| BuildComplianceSvcTag1Full | test_acc_template_full_svc_tag_1.xml       | DEVICESVCTAG1          | `All`          |
| BuildDeployeSvcTag1        | test_acc_template_deploy_svc_tag_1.xml     | DEVICESVCTAG1          | `System`       |

After exporting these files, destroy the templates as:

```sh
TF_VAR_DEVICESVCTAG1=${DEVICESVCTAG1} TF_VAR_DEVICESVCTAG2=${DEVICESVCTAG2} terraform destroy --auto-approve
```

# How to use the testdata in tests

The tests expect the environment variable `OME_TESTDATA_DIR` to be set. And its expected that its value will be a directory containing these testdata files.
So, set that environment variable with the folder where the exported test files are kept.
Example
```sh
export OME_TESTDATA_DIR="~/terraform-provider-ome/testdata"
```
> Relative paths will not work for acceptance tests as the tests are run from a temporary working directory.
