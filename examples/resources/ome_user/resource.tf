resource "ome_user" "code_1" {
  # required params
  username = "Dell"
  password = "Dell123!"
  role_id  = "10"
  # optional params
  user_type_id         = 1
  directory_service_id = 0
  description          = "Avengers alpha"
  locked               = false
  enabled              = false
}