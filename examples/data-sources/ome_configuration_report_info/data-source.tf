# Get configuration compliance report for a baseline
data "ome_configuration_report_info" "cr" {
  baseline_name = "BaselineName"
}