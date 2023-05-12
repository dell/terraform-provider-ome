# Get Deviceid's and servicetags of all devices that belong to a specified list of groups
data "ome_groupdevices_info" "gd" {
  device_group_names = ["WINDOWS"]
}