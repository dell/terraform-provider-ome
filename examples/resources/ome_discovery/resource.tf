resource "ome_discovery" "discover1" {
  name     = "discover-lab"
  schedule = "RunNow"
  # timeout to track the discovery job till x number of minutes.
  timeout = 10
  # ignore_partial_failure is used to control the terraform error in case of undiscovered ips after the discovery.
  ignore_partial_failure = true

  # discovery_config_targets is used to provide the server details
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
