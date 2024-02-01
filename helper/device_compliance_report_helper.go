/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://mozilla.org/MPL/2.0/
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"context"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"terraform-provider-ome/utils"
)

// GetAllDeviceComplianceReport get all device compliance report
func GetAllDeviceComplianceReport(clients *clients.Client, plan models.OMEDeviceComplianceData) ([]models.FirmwareBaselinesGetModel, error) {
	devices, err := clients.GetDevices(utils.ConvertListValueToStringSlice(plan.DeviceServiceTags), utils.ConvertListValueToIntSlice(plan.DeviceIDs), utils.ConvertListValueToStringSlice(plan.DeviceGroupNames))
	if err != nil {
		return []models.FirmwareBaselinesGetModel{}, err
	}
	_, uniqueDeviceIds, _ := clients.GetUniqueDevicesIdsAndServiceTags(devices)
	return clients.GetComplianceReportDetails(uniqueDeviceIds)
}

// SetStateDeviceComplianceReport set state device compliance report
func SetStateDeviceComplianceReport(ctx context.Context, device []models.FirmwareBaselinesGetModel) ([]models.DeviceComplianceData, error) {

	// Extract all the reports
	allDeviceReports := make([]models.DeviceComplianceModel, 0)
	for _, v := range device {
		allDeviceReports = append(allDeviceReports, v.DeviceComplianceReport...)
	}

	vals := make([]models.DeviceComplianceData, 0)
	for _, v := range allDeviceReports {
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
