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
	"terraform-provider-ome/utils"
)

// GetAllCatalogFirmware get all catalog firmware
func GetAllCatalogFirmware(client *clients.Client) (*models.Catalogs, error) {
	return client.GetAllCatalogFirmware()
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

	return vals, nil
}
