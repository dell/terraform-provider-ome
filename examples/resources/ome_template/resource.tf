# create a template with reference device.
resource "ome_template" "template_2" {
  name = "template_2"
  refdevice_servicetag = "MXL1234"
}