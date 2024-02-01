package helper

import (
	"context"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"terraform-provider-ome/utils"
)

// GetAllDeviceComplianceReport get all device compliance report
func GetAllDeviceComplianceReport(clients *clients.Client, baseLineID int64) (models.DeviceComplianceReportModel, error) {
	return clients.GetComplianceReportDetails(baseLineID)
}

// SetStateDeviceComplianceReport set state device compliance report
func SetStateDeviceComplianceReport(ctx context.Context, device models.DeviceComplianceReportModel) ([]models.DeviceComplianceData, error) {
	vals := make([]models.DeviceComplianceData, 0)
	for _, v := range device.Value {
		val := models.DeviceComplianceData{}
		err := utils.CopyFields(ctx, v, &val)
		if err != nil {
			return nil, err
		}
		complianceReports, reportsErr := setComponentComplianceReports(ctx, v.ComponentComplianceReports)
		if reportsErr != nil {
			return nil, reportsErr
		}
		val.ComponentComplianceReport = complianceReports
		vals = append(vals, val)
	}
	return vals, nil
}

// setComponentComplianceReports set component compliance reports
func setComponentComplianceReports(ctx context.Context, componentComplianceReports []models.ComponentComplianceReportModel) ([]models.ComponentComplianceReportsData, error) {
	vals := make([]models.ComponentComplianceReportsData, 0)
	for _, v := range componentComplianceReports {
		val := models.ComponentComplianceReportsData{}
		err := utils.CopyFields(ctx, v, &val)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}

	return vals, nil
}
