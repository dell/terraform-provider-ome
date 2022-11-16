terraform {
  required_providers {
    ome = {
      version = "1.0.0"
      source  = "dell/ome"
    }
  }
}

provider "ome" {
  username = var.username
  password = var.password
  host = var.host
  skipssl = var.skipssl
}

# Create Baseline using device ids
resource "ome_configuration_baseline" "baseline1" {
	baseline_name = "baseline1"
	ref_template_id = 745
	device_ids = [10001, 10002]
	description = "baseline description"
}


# Create Baseline using device servicetag with daily notification scheduled 
resource "ome_configuration_baseline" "baseline2" {
	baseline_name = "baseline2"
	ref_template_id = 745
	device_servicetags = ["MXL1234", "MXL1235"]
	description = "baseline description"
	schedule_notification = true
	notification_on_schedule = true
	email_addresses = ["test@testmail.com"]
	cron = "0 30 11 * * ? *"
	output_format = "csv"
}


# Create Baseline using device ids with daily notification on status changing to non-compliant 
resource "ome_configuration_baseline" "baseline3" {
	baseline_name = "baseline3"
	ref_template_id = 745
	device_ids = [10001, 10002]
	description = "baseline description"
	schedule_notification = true
	email_addresses = ["test@testmail.com"]
  output_format = "pdf"
}
