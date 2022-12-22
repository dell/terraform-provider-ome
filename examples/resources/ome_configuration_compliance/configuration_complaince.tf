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

# remediate baseline for the specified target devices 
resource "ome_configuration_compliance" "remeditation" {
	baseline_name = "baseline_name"
	target_devices = [
      {
        device_service_tag = "MX12345"
        compliance_status = "Compliant"
      }
    ]
}

# remediate baseline for the specified target devices with scheduling
resource "ome_configuration_compliance" "remeditation1" {
	baseline_name = "baseline_name"
	target_devices = [
      {
        device_service_tag = "MX12345"
        compliance_status = "Compliant"
      }
    ]
    run_later = true
    cron = "0 00 11 14 02 ? 2032"
}

# create baseline 
resource "ome_configuration_baseline" "baseline" {
	baseline_name = var.baselinename
	ref_template_name = "Mytemplate"
	device_servicetags = ["MX12345"]
	description = "baseline description"
}

resource "ome_configuration_compliance" "remeditation" {
		baseline_name = var.baselinename
		target_devices = [
      {
        device_service_tag = "MX12345"
        compliance_status = "Compliant"
      }
    ]
    depends_on = [
      ome_configuration_baseline.baseline
    ]
}
