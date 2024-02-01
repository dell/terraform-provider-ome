package helper

import (
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
)

// GetAllDeviceComplianceReport get all device compliance report
func GetAllDeviceComplianceReport(clients *clients.Client, baseLineID int64) (models.DeviceComplianceReportModel, error) {
	return clients.GetComplianceReportDetails(baseLineID)
}
