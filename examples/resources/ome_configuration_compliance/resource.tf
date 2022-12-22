# Manage remediation using baseline and Device Servicetags
resource "ome_configuration_compliance" "remeditation" {
	baseline_name = "baseline_name"
	target_devices = [
      {
        device_service_tag = "MX12345"
        compliance_status = "Compliant"
      }
    ]
}