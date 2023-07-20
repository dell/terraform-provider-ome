resource "ome_user" "code_1" {
  user_type_id         = 1
  directory_service_id = 0
  description          = "Avengers alpha"
  password             = "Dell123!"
  username             = "Dell"
  role_id              = "10"
  locked               = false
  enabled              = false
}