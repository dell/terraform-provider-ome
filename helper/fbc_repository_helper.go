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
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"
	"terraform-provider-ome/utils"
)

// GetAllRepositories get all FBC Repositories
func GetAllRepositories(client *clients.Client) ([]models.RepositoryModel, error) {
	catalogs, err := GetAllCatalogFirmware(client)
	if err != nil {
		return nil, err
	}
	// array of catalogs
	catalogsArray := catalogs.Value
	if len(catalogsArray) == 0 {
		return nil, fmt.Errorf("no data found for repositories")
	}
	repositoryArray := make([]models.RepositoryModel, 0)
	for _, catalog := range catalogsArray {
		repositoryArray = append(repositoryArray, catalog.Repository)
	}
	return repositoryArray, nil
}

// GetFilteredRepositoriesByName get all FBC Repositories by name
func GetFilteredRepositoriesByName(context context.Context, repos []models.RepositoryModel, plan models.OMERepositoryData) ([]models.RepositoryModel, error) {
	var filteredArray []models.RepositoryModel
	var repoNames []string

	for _, name := range plan.Names {
		repoNames = append(repoNames, name.ValueString())
	}
	for _, repo := range repos {
		if utils.ContainsString(context, repoNames, repo.Name) {
			filteredArray = append(filteredArray, repo)
		}
	}

	if len(repoNames) != 0 && len(filteredArray) != len(repoNames) {
		return nil, fmt.Errorf("one of the filtered names (%s) does not exist in the list of repositories", repoNames)
	}

	return filteredArray, nil
}
