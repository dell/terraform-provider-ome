# Manage baseline using Device Servicetags
resource "ome_discovery" "discovery_job" {
  discovery_config_group_name = "Discovery-Job-Resource"
  discovery_config_models = [
    {
      discovery_config_id          = 331105536
      discovery_config_description = ""
      discovery_config_status      = ""
      discovery_config_targets = [
        {
          discovery_config_target_id = 0
          network_address_detail     = "10.226.197.136"
          address_type               = 30
          disabled                   = false
          exclude                    = false
        }
      ]
      connection_profile_id             = 0
      connection_profile                = "{\"profileName\":\"\",\"profileDescription\":\"\",\"type\":\"DISCOVERY\",\"credentials\":[{\"id\":0,\"type\":\"WSMAN\",\"authType\":\"Basic\",\"modified\":false,\"credentials\":{\"username\":\"root\",\"password\":\"calvin\",\"caCheck\":false,\"cnCheck\":false,\"port\":443,\"retries\":3,\"timeout\":60,\"isHttp\":false,\"keepAlive\":false}}]}"
      device_type                       = [1000]
      discovery_config_vendor_platforms = []
    }
  ]
  schedule = {
    run_now    = false
    run_later  = true
    cron       = "0 55 19 16 9 ? 2023"
    start_time = ""
    end_time   = ""
    recurring = {
      hourly = {

      }
      daily = {
        time = {

        }
      }
      weekley = {
        time = {

        }
      }
    }
  }
  discovery_config_task_param = []
  discovery_config_tasks      = []
  create_group                = true
  trap_destination            = false
}