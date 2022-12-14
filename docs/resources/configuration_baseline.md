---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ome_configuration_baseline Resource - terraform-provider-ome"
subcategory: ""
description: |-
  Resource for managing configuration baselines on OpenManage Enterprise. Updates are supported for the following parameters: baseline_name, description, device_ids, device_servicetags, schedule_notification, notification_on_schedule, email_addresses, output_format, cron, job_retry_count, sleep_interval.
---

# ome_configuration_baseline (Resource)

Resource for managing configuration baselines on OpenManage Enterprise. Updates are supported for the following parameters: `baseline_name`, `description`, `device_ids`, `device_servicetags`, `schedule_notification`, `notification_on_schedule`, `email_addresses`, `output_format`, `cron`, `job_retry_count`, `sleep_interval`.

## Example Usage

```terraform
# Manage baseline using Device Servicetags
resource "ome_configuration_baseline" "baseline_name" {
	baseline_name = "Baseline Name"
	device_servicetags = ["MXL1234","MXL1235"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `baseline_name` (String) Name of the Baseline.

### Optional

- `cron` (String) Cron expression for notification schedule.
- `description` (String) Description of the baseline.
- `device_ids` (List of Number) List of the device id on which the baseline compliance needs to be run.
- `device_servicetags` (List of String) List of the device servicetag on which the baseline compliance needs to be run.
- `email_addresses` (List of String) Email addresses for notification.
- `job_retry_count` (Number) Number of times the job has to be polled to get the final status of the resource.
- `notification_on_schedule` (Boolean) Flag to set scheduled notification via cron.
- `output_format` (String) Output format type, the input is case senitive.
- `ref_template_id` (Number) Reference template ID.
- `ref_template_name` (String) Reference template name.
- `schedule_notification` (Boolean) Flag to schedule notification via email.
- `sleep_interval` (Number) Sleep time interval for job polling in seconds.

### Read-Only

- `id` (Number) ID of the resource.
- `task_id` (Number) Task id associated with baseline.

