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

# Deploy template using Device Service tags
resource "ome_deployment" "deploy-template-1" {
	template_name = "deploy-template-1"
	device_servicetags = ["MXL1234","MXL1235"]
  job_retry_count = 30
  sleep_interval = 10
}

# Deploy template using Device Id's
resource "ome_deployment" "deploy-template-2" {
	template_name = "deploy-template-2"
	device_ids = [10001, 10002]
 }

# Get Deviceid's or servicetags from a specified list of groups
data "ome_groupdevices_info" "gd" {
  device_group_names = ["WINDOWS"]
}

# Deploy template for group by fetching the device ids using data sources
resource "ome_deployment" "deploy-template-3" {
	template_name = "deploy-template-3"
	device_ids = data.ome_groupdevices_info.gd.device_ids
}

# Deploy template using Device Service tags with Schedule
resource "ome_deployment" "deploy-template-4" {
	template_name = "deploy-template-4"
	device_servicetags = ["MXL1234"]
	run_later = true
	cron = "0 45 12 19 10 ? 2022"
}

# Deploy template using Device ids and deploy device attributes
resource "ome_deployment" "deploy-template-5" {
	template_name = "deploy-template-5"
	device_ids = [10001, 10002]
	device_attributes = [
		{
			device_ids = [10001, 10002]
			attributes = [
				{
					attribute_id = 1197967
					display_name = "ServerTopology 1 Aisle Name"
					value = "aisle updated value"
					is_ignored = false
				}
			]
		}
	]
}

# Deploy template using Device ids and boot to network iso
resource "ome_deployment" "deploy-template-6" {
  template_name = "deploy-template-6"
	device_ids = [10001, 10002]
	boot_to_network_iso = {
		boot_to_network = true
		share_type = "CIFS"
		iso_timeout = 240
		iso_path = "/cifsshare/unattended/unattended_rocky8.6.iso"
		share_detail = {
			ip_address = "192.168.0.2"
			share_name = ""
			work_group = ""
			user = "username"
			password = "password"
		}
	}
	job_retry_count = 30
}

# Deploy template using Device ids by changing the job_retry_count and sleep_interval and ignore the same during updates
resource "ome_deployment" "deploy-template-7" {
  device_servicetags = ["MXL1234"]
  job_retry_count = 30
  sleep_interval = 10

  lifecycle {
    ignore_changes = [
      job_retry_count,
      sleep_interval
    ]
  }
}

# Deploy template using Device service tags and groupnames
resource "ome_deployment" "deploy-template-8" {
	template_id = 614
	device_servicetags = concat(data.ome_groupdevices_info.gd.device_servicetags, ["MXL1235"])
}