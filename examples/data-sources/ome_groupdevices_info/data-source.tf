# Get Deviceid's or servicetags from a specified list of groups
data "ome_groupdevices_info" "gd" {
  device_group_names = ["WINDOWS"]
}