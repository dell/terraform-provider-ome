resource "ome_discovery" "discover1" {
  name                   = "discover-lab"
  schedule               = "RunNow"
  timeout                = 10
  ignore_partial_failure = true
  discovery_config_targets = [
    {
      network_address_detail = ["10.0.0.0"]
      device_type            = ["SERVER"]
      wsman = {
        username = "user"
        password = "password"
      }
  }]
}
