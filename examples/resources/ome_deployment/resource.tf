# Deploy template using Device Service tags
resource "ome_deployment" "deploy-template-1" {
	template_name = "deploy-template-1"
	device_servicetags = ["MXL1234","MXL1235"]
}
