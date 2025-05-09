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
	"encoding/json"
	"fmt"
	"terraform-provider-ome/models"
)

// GetDeviceInventory returns the inventory of a device
func (c *Client) GetDeviceInventory(deviceID int64) (models.DeviceInventory, error) {
	inv := models.NewDeviceInventory()
	path := fmt.Sprintf(DeviceInventoryAPI, deviceID)
	response, err := c.Get(path, nil, nil)
	if err != nil {
		return inv, fmt.Errorf("error querying device inventory: %w", err)
	}
	bodyData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return inv, getBodyError
	}
	err = c.JSONUnMarshalValue(bodyData, &inv)
	return inv, err
}

// GetDeviceInventoryByType returns the inventory of a device of a particular type
func (c *Client) GetDeviceInventoryByType(deviceID int64, inventoryType string) (models.DeviceInventory, error) {
	inv := models.NewDeviceInventory()
	path := fmt.Sprintf(DeviceInventorySingleAPI, deviceID, inventoryType)
	response, err := c.Get(path, nil, nil)
	if err != nil {
		return inv, fmt.Errorf("error querying device inventory with type %s: %w", inventoryType, err)
	}
	bodyData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return inv, getBodyError
	}
	temp := models.DeviceInventoryInfo{}
	if err := json.Unmarshal(bodyData, &temp); err != nil {
		return inv, fmt.Errorf("error unmarshalling device type: %w", err)
	}
	err = inv.AddInfo(temp)
	return inv, err
}

// RefreshDeviceInventory - creates a job to refresh inventory of devices
func (c *Client) RefreshDeviceInventory(deviceIDs []int64, opts JobOpts) (JobResp, error) {
	targets := make([]models.JobTargetType, 0)
	for _, id := range deviceIDs {
		targets = append(targets, models.JobTargetType{
			ID:         id,
			TargetType: models.DeviceTargetType,
		})
	}
	payload := models.JobPayload{
		Enabled:        true,
		JobName:        opts.Name,
		JobDescription: opts.Description,
		Schedule:       opts.getSchedule(),
		JobType:        models.InventoryRefreshJobType,
		Params: map[string]string{
			"action":                   "CONFIG_INVENTORY",
			"isCollectDriverInventory": "true",
		},
		Targets: targets,
	}
	response, err := c.CreateJob(payload)
	if err != nil {
		return JobResp{}, fmt.Errorf("error creating device inventory refresh job: %w", err)
	}
	return response, nil
}
