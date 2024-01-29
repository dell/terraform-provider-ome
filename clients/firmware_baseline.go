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

package clients

import (
	"fmt"
	"strconv"
	"terraform-provider-ome/models"
)

// CreateFirmwareBaseline - Creates a new baseline in the catalog
func (c *Client) CreateFirmwareBaseline(payload models.CreateUpdateFirmwareBaseline) (int64, error) {
	data, _ := c.JSONMarshal(payload)
	response, err := c.Post(FirmwareBaselineAPI, nil, data)
	if err != nil {
		return -1, err
	}
	respData, _ := c.GetBodyData(response.Body)
	val, _ := strconv.ParseInt(string(respData), 10, 64)
	return val, nil
}

// GetFirmwareBaselineWithID - Gets the baseline details by baseline ID
func (c *Client) GetFirmwareBaselineWithID(id int64) (models.FirmwareBaselinesModel, error) {
	omeBaseline := models.FirmwareBaselinesModel{}
	response, err := c.Get(fmt.Sprintf(FirmwareBaselineAPI+"(%d)", id), nil, nil)
	if err != nil {
		return omeBaseline, err
	}

	respData, _ := c.GetBodyData(response.Body)

	err = c.JSONUnMarshal(respData, &omeBaseline)
	if err != nil {
		return omeBaseline, err
	}

	return omeBaseline, nil
}

// GetFirmwareBaselineWithName - Gets the baseline details by baseline name
func (c *Client) GetFirmwareBaselineWithName(name string) (models.FirmwareBaselinesModel, error) {
	omeBaseline := []models.FirmwareBaselinesModel{}
	err := c.GetPaginatedDataWithQueryParam(FirmwareBaselineAPI, map[string]string{"$expand": "DeviceComplianceReports"}, &omeBaseline)
	if err != nil {
		return models.FirmwareBaselinesModel{}, err
	}
	if len(omeBaseline) == 0 {
		return models.FirmwareBaselinesModel{}, nil
	}

	for _, baseline := range omeBaseline {
		if baseline.Name == name {
			return baseline, nil
		}
	}
	return models.FirmwareBaselinesModel{}, nil
}

// DeleteFirmwareBaseline - Deletes the specified baseline
func (c *Client) DeleteFirmwareBaseline(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	delBody := map[string]any{"BaselineIds": ids}
	body, err := c.JSONMarshal(delBody)
	if err != nil {
		return err
	}
	_, delErr := c.Post(RemoveFirmwareBaseline, nil, body)
	return delErr
}

// UpdateFirmwareBaseline - Updates the specified baseline
func (c *Client) UpdateFirmwareBaseline(baseline models.CreateUpdateFirmwareBaseline) (int64, error) {
	data, _ := c.JSONMarshal(baseline)
	response, err := c.Put(fmt.Sprintf(FirmwareBaselineAPI+"(%d)", baseline.ID), nil, data)
	if err != nil {
		return -1, err
	}
	respData, _ := c.GetBodyData(response.Body)
	val, _ := strconv.ParseInt(string(respData), 10, 64)
	return val, nil
}
