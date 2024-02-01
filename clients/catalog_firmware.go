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

package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

// GetAllCatalogFirmware - Get All catalog firmware
func (c *Client) GetAllCatalogFirmware() (*models.Catalogs, error) {
	response := models.Catalogs{}
	err := c.GetValueWithPagination(RequestOptions{
		URL: CatalogFirmwareAPI,
	}, &response.Value)
	return &response, err
}

// GetSpecificCatalogFirmware - Get specific catalog firmware
func (c *Client) GetSpecificCatalogFirmware(id int64) (models.CatalogsModel, error) {
	catalog := models.CatalogsModel{}
	resp, err := c.Get(fmt.Sprintf(CatalogFirmwareSpecificAPI, id), nil, nil)
	if err != nil {
		return catalog, err
	}
	bodyData, _ := c.GetBodyData(resp.Body)
	err = c.JSONUnMarshal(bodyData, &catalog)
	if err != nil {
		err = fmt.Errorf(ErrInvalidFirmwareCatalogIdentifiers+" %w", err)
	}
	return catalog, err
}

// CreateCatalogFirmware - Create catalog firmware
func (c *Client) CreateCatalogFirmware(payload models.CatalogsModel) (models.CatalogsModel, error) {
	data, _ := c.JSONMarshal(payload)
	response, err := c.Post(CatalogFirmwareAPI, nil, data)
	var returnVal = models.CatalogsModel{}
	if err != nil {
		return returnVal, err
	}
	bodyData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(bodyData, &returnVal)
	return returnVal, err
}

// DeleteCatalogFirmware - Deletes firmware catalogs
func (c *Client) DeleteCatalogFirmware(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	bodyv := map[string]any{"CatalogIds": ids}
	body, errb := c.JSONMarshal(bodyv)
	if errb != nil {
		return errb
	}
	_, err := c.Post(DeleteFirmwareCatalogAPI, nil, body)
	return err
}

// UpdateCatalogFirmware - Update firmware catalog details
func (c *Client) UpdateCatalogFirmware(id int64, payload models.CatalogsModel) (models.CatalogsModel, error) {
	data, _ := c.JSONMarshal(payload)
	response, err := c.Put(fmt.Sprintf(CatalogFirmwareSpecificAPI, id), nil, data)
	var returnVal = models.CatalogsModel{}
	if err != nil {
		return returnVal, err
	}

	bodyData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(bodyData, &returnVal)
	return returnVal, err
}
