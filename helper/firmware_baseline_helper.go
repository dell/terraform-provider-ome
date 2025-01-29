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
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"terraform-provider-ome/utils"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CreateTargetModel create the target model based on the input plan
func CreateTargetModel(client *clients.Client, plan models.FirmwareBaselineResource) ([]models.TargetModel, error) {
	var targets []models.TargetModel
	var filterList []string

	if len(plan.DeviceNames.Elements()) != 0 {
		filterList = utils.ConvertListValueToStringSlice(plan.DeviceNames)
		// Get all devices based on the names in the list
		devices, err := client.GetValidDevicesByNames(filterList)
		if err != nil {
			return targets, err
		}

		for _, device := range devices {
			targetModel := models.TargetModel{}
			targetModel.ID = device.ID
			targetModeltype := models.TargetTypeModel{}
			targetModeltype.ID = device.Type
			targetModeltype.Name = "DEVICE"
			targetModel.Type = targetModeltype
			targets = append(targets, targetModel)
		}
	} else if len(plan.GroupNames.Elements()) != 0 {
		filterList = utils.ConvertListValueToStringSlice(plan.GroupNames)
		// Get all groups and subgroups based on the names in the list
		groups, err := client.GetValidGroupsByNames(filterList)
		if err != nil {
			return targets, err
		}
		for _, group := range groups {
			targetModel := models.TargetModel{}
			targetModel.ID = group.ID
			targetModeltype := models.TargetTypeModel{}
			targetModeltype.ID = group.TypeID
			targetModeltype.Name = "GROUP"
			targetModel.Type = targetModeltype
			targets = append(targets, targetModel)
		}
	} else if len(plan.DeviceServiceTags.Elements()) != 0 {
		filterList = utils.ConvertListValueToStringSlice(plan.DeviceServiceTags)
		// Get all devices based on the service tags in the list
		devices, err := client.GetDevices(filterList, nil, nil)
		if err != nil {
			return targets, err
		}

		for _, device := range devices {
			targetModel := models.TargetModel{}
			targetModel.ID = device.ID
			targetModeltype := models.TargetTypeModel{}
			targetModeltype.ID = device.Type
			targetModeltype.Name = "DEVICE"
			targetModel.Type = targetModeltype
			targets = append(targets, targetModel)
		}
	} else {
		return targets, fmt.Errorf("either device_names or group_names or device_service_tags is required")
	}

	return targets, nil
}

// SetStateBaseline set state baseline
func SetStateBaseline(ctx context.Context, baseline models.FirmwareBaselinesModel, plan models.FirmwareBaselineResource) (models.FirmwareBaselineResource, error) {
	var state models.FirmwareBaselineResource
	// Copy a majority of the state fields
	errCopy := utils.CopyFields(ctx, baseline, &state)
	if errCopy != nil {
		return state, errCopy
	}

	// Map compliance summary object
	mappedComplianceSummary, errMap := MapBaselineComplianceSummary(baseline.ComplianceSummary)

	if errMap.HasError() {
		return state, fmt.Errorf("failed to map compliance summary")
	}
	state.ComplianceSummary = mappedComplianceSummary

	// Map Targets
	mappedTargets, mapTargetDiags := MapBaselineTargets(baseline.Targets)
	if mapTargetDiags.HasError() {
		return state, fmt.Errorf("failed to map targets")
	}
	state.Targets = mappedTargets
	if baseline.ID != nil {
		state.ID = types.Int64Value(int64(*baseline.ID))
	}
	if baseline.TaskStatusID != nil {
		state.TaskStatus = types.StringValue(GetJobStatus(int64(*baseline.TaskStatusID)))
	}
	if baseline.TaskID != nil {
		state.TaskID = types.Int64Value(int64(*baseline.TaskID))
	}
	//state.LastRun = types.StringPointerValue(baseline.LastRun)

	// state.Description = types.StringPointerValue(baseline.Description)
	// state.Is64Bit = types.BoolValue(baseline.Is64Bit)
	// state.FilterNoRebootRequired = types.BoolPointerValue(baseline.FilterNoRebootRequired)

	// Set the user input values to the state
	state.Name = plan.Name
	state.CatalogName = plan.CatalogName

	return state, nil
}

// MapBaselineComplianceSummary maps the compliance summary model to a types.Object.
func MapBaselineComplianceSummary(complianceSummary models.ComplianceSummaryModel) (types.Object, diag.Diagnostics) {

	typeKey := map[string]attr.Type{
		"compliance_status":   types.StringType,
		"number_of_critical":  types.Int64Type,
		"number_of_downgrade": types.Int64Type,
		"number_of_normal":    types.Int64Type,
		"number_of_warning":   types.Int64Type,
		"number_of_unknown":   types.Int64Type,
	}

	complianceMap := make(map[string]attr.Value)
	complianceMap["compliance_status"] = types.StringPointerValue(complianceSummary.ComplianceStatus)
	complianceMap["number_of_critical"] = types.Int64PointerValue(complianceSummary.NumberOfCritical)
	complianceMap["number_of_downgrade"] = types.Int64PointerValue(complianceSummary.NumberOfDowngrade)
	complianceMap["number_of_normal"] = types.Int64PointerValue(complianceSummary.NumberOfNormal)
	complianceMap["number_of_warning"] = types.Int64PointerValue(complianceSummary.NumberOfWarning)
	complianceMap["number_of_unknown"] = types.Int64PointerValue(complianceSummary.NumberOfUnknown)
	dcrObject, diags := types.ObjectValue(typeKey, complianceMap)

	return dcrObject, diags

}

// MapBaselineTargets maps a list of TargetModel objects into a List of attr.Value objects.
func MapBaselineTargets(targets []models.TargetModel) (types.List, diag.Diagnostics) {
	var targetsMap []attr.Value
	targetKey := map[string]attr.Type{
		"id": types.Int64Type,
		"type": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.Int64Type,
				"name": types.StringType,
			},
		},
	}

	for _, target := range targets {
		targetMap := make(map[string]attr.Value)
		targetMap["id"] = types.Int64Value(int64(target.ID))

		typeMap := make(map[string]attr.Value)
		typeMap["id"] = types.Int64Value(int64(target.Type.ID))
		typeMap["name"] = types.StringValue(target.Type.Name)
		typeObject, diags := types.ObjectValue(
			map[string]attr.Type{
				"id":   types.Int64Type,
				"name": types.StringType,
			}, typeMap)
		if diags.HasError() {
			return types.List{}, diags
		}

		targetMap["type"] = typeObject
		targetObject, diags := types.ObjectValue(targetKey, targetMap)
		if diags.HasError() {
			return types.List{}, diags
		}
		targetsMap = append(targetsMap, targetObject)
	}

	return types.ListValue(types.ObjectType{AttrTypes: targetKey}, targetsMap)
}

// GetFirmwareBaselineWithName gets a baseline by name
func GetFirmwareBaselineWithName(client clients.Client, name string) (models.FirmwareBaselinesModel, error) {
	return client.GetFirmwareBaselineWithName(name)
}

// GetFirmwareBaselineWithID gets a baseline by id
func GetFirmwareBaselineWithID(client clients.Client, id int64) (models.FirmwareBaselinesModel, error) {
	return client.GetFirmwareBaselineWithID(id)
}

// CreateFirmwareBaseline - Creates a new Firmware baseline
func CreateFirmwareBaseline(client *clients.Client, payload models.CreateUpdateFirmwareBaseline) (int64, error) {
	return client.CreateFirmwareBaseline(payload)
}

// DeleteFirmwareBaseline deletes the given firmware baseline
func DeleteFirmwareBaseline(client clients.Client, id int64) error {
	baselineIds := []int64{id}
	return client.DeleteFirmwareBaseline(baselineIds)
}

// UpdateFirmwareBaseline updates the given firmware baseline
func UpdateFirmwareBaseline(client clients.Client, state models.FirmwareBaselineResource, plan models.FirmwareBaselineResource) (int64, error) {
	payload := models.CreateUpdateFirmwareBaseline{}
	payload.ID = state.ID.ValueInt64()

	if plan.CatalogName.ValueString() != "" && plan.CatalogName.ValueString() != state.CatalogName.ValueString() {
		catalog, err := GetCatalogFirmwareByName(&client, plan.CatalogName.ValueString())
		if err != nil {
			return -1, err
		}
		if catalog == nil {
			return -1, fmt.Errorf("catalog %s not found", plan.CatalogName.ValueString())
		}
		payload.CatalogID = catalog.ID
		payload.RepositoryID = catalog.Repository.ID
	} else {
		payload.CatalogID = state.CatalogID.ValueInt64()
		payload.RepositoryID = state.RepositoryID.ValueInt64()
	}
	if plan.Is64Bit.ValueBool() != state.Is64Bit.ValueBool() {
		payload.Is64Bit = plan.Is64Bit.ValueBool()
	} else {
		payload.Is64Bit = state.Is64Bit.ValueBool()
	}
	if plan.FilterNoRebootRequired.ValueBool() != state.FilterNoRebootRequired.ValueBool() {
		payload.FilterNoRebootRequired = plan.FilterNoRebootRequired.ValueBool()
	} else {
		payload.FilterNoRebootRequired = state.FilterNoRebootRequired.ValueBool()
	}
	if plan.Description.ValueString() != "" && plan.Description.ValueString() != state.Description.ValueString() {
		payload.Description = plan.Description.ValueString()
	} else {
		payload.Description = state.Description.ValueString()
	}
	if plan.Name.ValueString() != "" && plan.Name.ValueString() != state.Name.ValueString() {
		payload.Name = plan.Name.ValueString()
	} else {
		payload.Name = state.Name.ValueString()
	}

	targets, err := CreateTargetModel(&client, plan)

	if err != nil {
		return -1, fmt.Errorf("unable to create target model for: %s. details: %s", plan.Name.ValueString(), err.Error())
	}
	payload.Targets = targets

	id, err := client.UpdateFirmwareBaseline(payload)
	if err != nil {
		return -1, err
	}
	return id, nil
}
