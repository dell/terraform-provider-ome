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
	"fmt"
	"strings"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"terraform-provider-ome/utils"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllCatalogFirmware get all catalog firmware
func GetAllCatalogFirmware(client *clients.Client) (*models.Catalogs, error) {
	return client.GetAllCatalogFirmware()
}

// GetSpecificCatalogFirmware get specific catalog firmware
func GetSpecificCatalogFirmware(client *clients.Client, id int64) (models.CatalogsModel, error) {
	return client.GetSpecificCatalogFirmware(id)
}

// UpdateCatalogFirmware update catalog firmware
func UpdateCatalogFirmware(client *clients.Client, id int64, payload models.CatalogsModel) (models.CatalogsModel, error) {
	return client.UpdateCatalogFirmware(id, payload)
}

// CreateCatalogFirmware create catalog firmware
func CreateCatalogFirmware(client *clients.Client, payload models.CatalogsModel) (models.CatalogsModel, error) {
	return client.CreateCatalogFirmware(payload)
}

// DeleteCatalogFirmware delete catalog firmware
func DeleteCatalogFirmware(client *clients.Client, id int64) error {
	ids := []int64{id}
	return client.DeleteCatalogFirmware(ids)
}

// SetStateCatalogFirmware set state catalog firmware
func SetStateCatalogFirmware(ctx context.Context, cat models.CatalogsModel, plan models.OmeSingleCatalogResource) (models.OmeSingleCatalogResource, error) {
	var state models.OmeSingleCatalogResource
	// Copy a majority of the state fields
	errCopy := utils.CopyFields(ctx, cat, &state)
	if errCopy != nil {
		return state, errCopy
	}
	// Map the associated baselines
	mappedBaselines, mapDiags := MapAssociatedBaselines(cat.AssociatedBaselines)
	if mapDiags.HasError() {
		return state, fmt.Errorf("failed to map associated baselines")
	}
	state.AssociatedBaselines = mappedBaselines

	// Map Repository Object
	mappedRepository, mapRepoDiags := MapFirmwareCatalogRepository(cat.Repository)
	if mapRepoDiags.HasError() {
		return state, fmt.Errorf("failed to map repository")
	}
	state.Repository = mappedRepository
	// Set the user input values to the state
	state.Name = plan.Name
	state.CatalogUpdateType = plan.CatalogUpdateType
	state.ShareType = plan.ShareType
	state.CatalogFilePath = plan.CatalogFilePath
	state.CatalogRefreshSchedule = plan.CatalogRefreshSchedule
	state.ShareUser = plan.ShareUser
	state.SharePassword = plan.SharePassword
	state.Domain = plan.Domain
	state.ShareAddress = plan.ShareAddress

	return state, nil
}

// GetIDFromNameFirmwareCatalog - Get the ID after create, for whatever reason the create api does not return the actual ID Instead it returns 0. The only way to get the true id is to get all of the catalogs and find the one that matches by name (names are required to be unique for catalogs)
func GetIDFromNameFirmwareCatalog(client *clients.Client, name string) (int64, error) {
	allCats, allErr := GetAllCatalogFirmware(client)
	if allErr != nil {
		return 0, fmt.Errorf(`Unable to get the Id of the catalog for catalog: `+name+`.`, allErr.Error())
	}
	for _, cm := range allCats.Value {
		if cm.Repository.Name == name {
			return cm.ID, nil
		}
	}
	return 0, fmt.Errorf(`Unable to get the Id of the catalog for catalog: ` + name + `. catalog was unable to be found.`)
}

// MapAssociatedBaselines map the associated baselines to the terraform list attribute
func MapAssociatedBaselines(baselines []models.AssociatedBaselineModel) (types.List, diag.Diagnostics) {
	var genObjects []attr.Value
	typeKey := map[string]attr.Type{
		"baseline_id":   types.Int64Type,
		"baseline_name": types.StringType,
	}
	for _, gen := range baselines {
		genMap := make(map[string]attr.Value)
		genMap["baseline_id"] = types.Int64Value(gen.BaselineID)
		genMap["baseline_name"] = types.StringValue(gen.BaselineName)

		genObject, _ := types.ObjectValue(typeKey, genMap)
		genObjects = append(genObjects, genObject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: typeKey}, genObjects)
}

// MapFirmwareCatalogRepository map the repository to the terraform list attribute
func MapFirmwareCatalogRepository(repo models.RepositoryModel) (types.Object, diag.Diagnostics) {
	typeKey := map[string]attr.Type{
		"backup_existing_catalog": types.BoolType,
		"check_certificate":       types.BoolType,
		"description":             types.StringType,
		"domain_name":             types.StringType,
		"editable":                types.BoolType,
		"id":                      types.Int64Type,
		"name":                    types.StringType,
		"repository_type":         types.StringType,
		"source":                  types.StringType,
		"username":                types.StringType,
	}

	genMap := make(map[string]attr.Value)
	genMap["backup_existing_catalog"] = types.BoolValue(repo.BackupExistingCatalog)
	genMap["check_certificate"] = types.BoolValue(repo.CheckCertificate)
	genMap["description"] = types.StringValue(repo.Description)
	genMap["domain_name"] = types.StringValue(repo.DomainName)
	genMap["editable"] = types.BoolValue(repo.Editable)
	genMap["id"] = types.Int64Value(repo.ID)
	genMap["name"] = types.StringValue(repo.Name)
	genMap["repository_type"] = types.StringValue(repo.RepositoryType)
	genMap["source"] = types.StringValue(repo.Source)
	genMap["username"] = types.StringValue(repo.Username)

	return types.ObjectValue(typeKey, genMap)
}

// MakeCatalogJSONModel  create catalog json model for create and update requests
func MakeCatalogJSONModel(id int64, repoID int64, plan models.OmeSingleCatalogResource) models.CatalogsModel {
	sourcePath, fileName := extractSourcePathAndFilename(plan.CatalogFilePath.ValueString())
	// For create requests
	if id == 0 || repoID == 0 {
		return models.CatalogsModel{
			Filename:   fileName,
			SourcePath: sourcePath,
			Schedule: models.ScheduleModel{
				Cron: CreateCronValue(plan.CatalogRefreshSchedule, plan.CatalogUpdateType.ValueString()),
			},
			Repository: models.RepositoryModel{
				Name:                  plan.Name.ValueString(),
				Description:           plan.Name.ValueString() + " terraform catalog",
				Source:                plan.ShareAddress.ValueString(),
				DomainName:            plan.Domain.ValueString(),
				Username:              plan.ShareUser.ValueString(),
				Password:              plan.SharePassword.ValueString(),
				RepositoryType:        plan.ShareType.ValueString(),
				BackupExistingCatalog: false,
				Editable:              true,
			},
		}
		// For Update Requests
	}
	// else for update requests
	return models.CatalogsModel{
		Filename:   fileName,
		SourcePath: sourcePath,
		ID:         id,
		Schedule: models.ScheduleModel{
			Cron: CreateCronValue(plan.CatalogRefreshSchedule, plan.CatalogUpdateType.ValueString()),
		},
		Repository: models.RepositoryModel{
			Name:                  plan.Name.ValueString(),
			Description:           plan.Name.ValueString() + " terraform catalog",
			Source:                plan.ShareAddress.ValueString(),
			DomainName:            plan.Domain.ValueString(),
			Username:              plan.ShareUser.ValueString(),
			Password:              plan.SharePassword.ValueString(),
			RepositoryType:        plan.ShareType.ValueString(),
			BackupExistingCatalog: false,
			Editable:              true,
			ID:                    repoID,
		},
	}
}

func extractSourcePathAndFilename(path string) (string, string) {
	parts := strings.Split(path, "/")
	filename := parts[len(parts)-1]
	return strings.TrimSuffix(path, "/"+filename), filename
}

// CreateCronValue create the cron value based on the users inputs for the refresh schedule
func CreateCronValue(schedule models.CatalogRefreshSchedule, catalogType string) string {
	// Cron value if manual schedule is set
	cronValue := "startnow"
	if catalogType == "Automatic" {
		// Need to subtract 1 since the cron job expects values 0-23
		hourlyValue := schedule.TimeOfDay.ValueInt64() - 1
		// If PM then add 12 to the value to make it in the 24 hour format
		if schedule.AmPm.ValueString() == "PM" {
			hourlyValue = hourlyValue + 12
		}
		if schedule.Cadence.ValueString() == "Daily" {
			cronValue = fmt.Sprintf("0 0 %v * * ? *", hourlyValue)
			// Else it is weekly
		} else {
			dayOfWeek := "mon"
			switch schedule.DayOfTheWeek.ValueString() {
			case "Monday":
				dayOfWeek = "mon"
			case "Tuesday":
				dayOfWeek = "tue"
			case "Wednesday":
				dayOfWeek = "wed"
			case "Thursday":
				dayOfWeek = "thu"
			case "Friday":
				dayOfWeek = "fri"
			case "Saturday":
				dayOfWeek = "sat"
			case "Sunday":
				dayOfWeek = "sun"
			}
			cronValue = fmt.Sprintf("0 0 %v ? * %s *", hourlyValue, dayOfWeek)
		}
	}
	return cronValue
}

// FilterCatalogFirmware filter catalog firmware
func FilterCatalogFirmware(ctx context.Context, filterElements []string, cat *models.Catalogs) ([]models.OmeSigleCatalogData, error) {
	vals := make([]models.OmeSigleCatalogData, 0)
	for _, v := range cat.Value {
		if len(filterElements) == 0 || utils.ContainsString(ctx, filterElements, v.Repository.Name) {
			val := models.OmeSigleCatalogData{}
			err := utils.CopyFields(ctx, v, &val)
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
		}
	}

	if len(filterElements) != 0 && len(filterElements) != len(vals) {
		return nil, fmt.Errorf("one of the filtered names (%s) does not exist in the list of catalogs", filterElements)
	}

	return vals, nil
}

// ValidateCatalogUpdate validates catalog update for the different share_type cases
func ValidateCatalogUpdate(plan models.OmeSingleCatalogResource, state models.OmeSingleCatalogResource) error {
	if plan.ShareType != state.ShareType {
		return fmt.Errorf("catalog share type is not allowed to be updated after create")
	}
	return nil
}

// ValidateCatalogCreate validates catalog create for the different share_type cases
func ValidateCatalogCreate(plan models.OmeSingleCatalogResource) error {

	// Validate Automatic Update type
	if plan.CatalogUpdateType.ValueString() == "Automatic" {
		if plan.CatalogRefreshSchedule.Cadence.ValueString() == "" || plan.CatalogRefreshSchedule.TimeOfDay.ValueInt64() == 0 || plan.CatalogRefreshSchedule.AmPm.ValueString() == "" {
			return fmt.Errorf("invalid automatic update configuration, please provide refresh schedule values 'cadence', 'timeOfDay' and 'amPm'")
		}
		if plan.CatalogRefreshSchedule.Cadence.ValueString() == "Weekly" && plan.CatalogRefreshSchedule.DayOfTheWeek.ValueString() == "" {
			return fmt.Errorf("invalid automatic update configuration, please provide 'dayOfTheWeek' if using the 'Weekly' cadence")
		}
	}

	// Validate Share Type Values
	switch plan.ShareType.ValueString() {
	case "NFS":
		if plan.ShareAddress.ValueString() == "" ||
			plan.CatalogFilePath.ValueString() == "" {
			return fmt.Errorf("invalid NFS share configuration, please provide 'share_address' and 'catalog_file_path'")
		}
	case "CIFS":
		if plan.ShareAddress.ValueString() == "" ||
			plan.CatalogFilePath.ValueString() == "" ||
			plan.ShareUser.ValueString() == "" ||
			plan.SharePassword.ValueString() == "" {
			return fmt.Errorf("invalid CIFS share configuration, please provide 'share_address', 'catalog_file_path', 'share_user', and 'share_password'")
		}
	case "HTTP":
		if plan.ShareAddress.ValueString() == "" ||
			plan.CatalogFilePath.ValueString() == "" {
			return fmt.Errorf("invalid HTTP share configuration, please provide 'share_address' and 'catalog_file_path'")
		}
	case "HTTPS":
		if plan.ShareAddress.ValueString() == "" ||
			plan.CatalogFilePath.ValueString() == "" {
			return fmt.Errorf("invalid HTTPS share configuration, please provide 'share_address' and 'catalog_file_path'")
		}
	}
	return nil
}

// GetCatalogFirmwareByName filter catalog firmware
func GetCatalogFirmwareByName(client *clients.Client, name string) (*models.CatalogsModel, error) {
	// Get all catalog firmware
	catalogFirmware, err := GetAllCatalogFirmware(client)
	if err != nil {
		return nil, err
	}

	// Filter catalog firmware based on name
	var filteredCatalogFirmware *models.CatalogsModel
	for _, catalog := range catalogFirmware.Value {
		if catalog.Repository.Name == name {
			// to resolve implicit memory aliasing
			cat := catalog
			filteredCatalogFirmware = &cat
			break
		}
	}

	if filteredCatalogFirmware == nil {
		return nil, fmt.Errorf("catalog firmware %s not found", name)
	}

	return filteredCatalogFirmware, nil
}
