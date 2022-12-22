# Manage baseline using Device Servicetags
resource "ome_configuration_baseline" "baseline_name" {
	baseline_name = "Baseline Name"
	device_servicetags = ["MXL1234","MXL1235"]
}
