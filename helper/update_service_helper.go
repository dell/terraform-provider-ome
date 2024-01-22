/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetFwBaselineComplianceReport gets compliance report for a firmware baseline.
func GetFwBaselineComplianceReport(ctx context.Context, client *clients.Client, baselineName int64, filterkey string, filterval string) (*models.ComplianceReport, error) {
	return client.GetFwBaselineComplianceReport(ctx, baselineName, filterkey, filterval)
}

// NewOmeFwComplianceReportList converts list of client.OmeFwComplianceReport to list of models.OmeFwComplianceReport
func NewOmeFwComplianceReportList(inputs []models.DeviceComplianceReport) []models.OmeFwComplianceReport {
	out := make([]models.OmeFwComplianceReport, 0)
	for _, input := range inputs {
		out = append(out, newOmeFwComplianceReport(input))
	}
	return out
}

// newOmeFwComplianceReport converts client.OmeFwComplianceReport to models.OmeFwComplianceReport
func newOmeFwComplianceReport(input models.DeviceComplianceReport) models.OmeFwComplianceReport {
	return models.OmeFwComplianceReport{
		ComplianceStatus:                types.StringValue(string(*input.ComplianceStatus)),
		ComponentComplianceReports:      newOmeComponentComplianceReportList(input.ComponentComplianceReports),
		DeviceFirmwareUpdateCapable:     types.BoolValue(input.DeviceFirmwareUpdateCapable),
		DeviceID:                        types.Int64Value(int64(input.DeviceID)),
		DeviceModel:                     types.StringValue(input.DeviceModel),
		DeviceName:                      types.StringValue(input.DeviceName),
		DeviceTypeID:                    types.Int64Value(int64(input.DeviceTypeID)),
		DeviceTypeName:                  types.StringValue(input.DeviceTypeName),
		DeviceUserFirmwareUpdateCapable: types.BoolValue(input.DeviceUserFirmwareUpdateCapable),
		FirmwareStatus:                  types.StringValue(input.FirmwareStatus),
		ID:                              types.Int64Value(int64(input.ID)),
		RebootRequired:                  types.BoolValue(input.RebootRequired),
		ServiceTag:                      types.StringValue(input.ServiceTag),
	}
}

// newOmeComponentComplianceReportList converts list of client.OmeComponentComplianceReport to list of models.OmeComponentComplianceReport
func newOmeComponentComplianceReportList(inputs []models.ComponentComplianceReport) []models.OmeComponentComplianceReport {
	out := make([]models.OmeComponentComplianceReport, 0)
	for _, input := range inputs {
		out = append(out, newOmeComponentComplianceReport(input))
	}
	return out
}

// newOmeComponentComplianceReport converts client.OmeComponentComplianceReport to models.OmeComponentComplianceReport
func newOmeComponentComplianceReport(input models.ComponentComplianceReport) models.OmeComponentComplianceReport {
	return models.OmeComponentComplianceReport{
		ComplianceDependencies:    newOmeComplianceDependencyList(input.ComplianceDependencies),
		ComplianceStatus:          types.StringValue(input.ComplianceStatus),
		ComponentType:             types.StringValue(input.ComponentType),
		Criticality:               types.StringValue(input.Criticality),
		CurrentVersion:            types.StringValue(input.CurrentVersion),
		DependencyUpgradeRequired: types.BoolValue(input.DependencyUpgradeRequired),
		ID:                        types.Int64Value(int64(input.ID)),
		ImpactAssessment:          types.StringValue(input.ImpactAssessment),
		Name:                      types.StringValue(input.Name),
		Path:                      types.StringValue(input.Path),
		PrerequisiteInfo:          types.StringValue(input.PrerequisiteInfo),
		RebootRequired:            types.BoolValue(input.RebootRequired),
		SourceName:                types.StringValue(input.SourceName),
		TargetIdentifier:          types.StringValue(input.TargetIdentifier),
		UniqueIdentifier:          types.StringValue(input.UniqueIdentifier),
		UpdateAction:              types.StringValue(input.UpdateAction),
		URI:                       types.StringValue(input.URI),
		Version:                   types.StringValue(input.Version),
	}
}

// newOmeComplianceDependencyList converts list of client.OmeComplianceDependency to list of models.OmeComplianceDependency
func newOmeComplianceDependencyList(inputs []models.ComplianceDependency) []models.OmeComplianceDependency {
	out := make([]models.OmeComplianceDependency, 0)
	for _, input := range inputs {
		out = append(out, newOmeComplianceDependency(input))
	}
	return out
}

// newOmeComplianceDependency converts client.OmeComplianceDependency to models.OmeComplianceDependency
func newOmeComplianceDependency(input models.ComplianceDependency) models.OmeComplianceDependency {
	return models.OmeComplianceDependency{
		ComplianceDependencyID: types.Int64Value(int64(input.ComplianceDependencyID)),
		IsHardDependency:       types.BoolValue(input.IsHardDependency),
		Name:                   types.StringValue(input.Name),
		Path:                   types.StringValue(input.Path),
		RebootRequired:         types.BoolValue(input.RebootRequired),
		SourceName:             types.StringValue(input.SourceName),
		UniqueIdentifier:       types.StringValue(input.UniqueIdentifier),
		UpdateAction:           types.StringValue(input.UpdateAction),
		URI:                    types.StringValue(input.URI),
		Version:                types.StringValue(input.Version),
	}
}
