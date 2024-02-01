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
	"terraform-provider-ome/utils"
)

// GetDevice is used to get device using serviceTag or devID in OME
func (c *Client) GetDevice(serviceTag string, devID int64) (models.Device, error) {

	device := models.Device{}
	var err error
	val := fmt.Sprintf("'%s'", serviceTag)
	key := "Identifier"
	if devID != 0 {
		val = fmt.Sprintf("%d", devID)
		key = "Id"
	}

	if val == "''" {
		return device, fmt.Errorf(ErrEmptyDeviceDetails)
	}

	response, err := c.Get(DeviceAPI, nil, map[string]string{"$filter": fmt.Sprintf("%s eq %s", key, val)})
	if err != nil {
		return device, err
	}

	bodyData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshalSingleValue(bodyData, &device)
	if err != nil {
		err = fmt.Errorf(ErrInvalidDeviceIdentifiers+" %s: %w", val, err)
	}
	return device, err
}

// RemoveDevices - function to remove specified list of devices by id
func (c *Client) RemoveDevices(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	bodyv := map[string]any{"DeviceIds": ids}
	body, errb := c.JSONMarshal(bodyv)
	if errb != nil {
		return errb
	}
	_, err := c.Post(DeviceRemovalAPI, nil, body)
	return err
}

// ValidateDevice is used to get deviceID using serviceTag in OME
func (c *Client) ValidateDevice(serviceTag string, devID int64) (int64, error) {

	var deviceID int64 = -1
	var err error
	val := fmt.Sprintf("'%s'", serviceTag)
	key := "Identifier"
	if devID != 0 {
		val = fmt.Sprintf("%d", devID)
		key = "Id"
	}

	if val == "''" {
		return deviceID, fmt.Errorf(ErrEmptyDeviceDetails)
	}

	response, err := c.Get(DeviceAPI, nil, map[string]string{"$filter": fmt.Sprintf("%s eq %s", key, val)})

	if err == nil {
		devices := models.Devices{}
		bodyData, _ := c.GetBodyData(response.Body)
		err = c.JSONUnMarshal(bodyData, &devices)
		if err == nil {
			if len(devices.Value) > 0 {
				deviceID = devices.Value[0].ID
				err = nil
			} else {
				err = fmt.Errorf(ErrInvalidDeviceIdentifiers+" %s", val)
			}
		}
	}
	return deviceID, err
}

// GetDevices - method to get all the devices associated with serviceTags, devIDs and groupNames.
func (c *Client) GetDevices(serviceTags []string, devIDs []int64, groupNames []string) ([]models.Device, error) {
	validDevices := []models.Device{}
	inValidDevices := []models.Device{}
	var invalidDevIDs []int64
	if len(devIDs) > 0 {
		for _, devID := range devIDs {
			device, err := c.GetDevice("", devID)
			if err != nil {
				inValidDevices = append(inValidDevices, device)
				invalidDevIDs = append(invalidDevIDs, devID)
			} else {
				validDevices = append(validDevices, device)
			}
		}
		if len(inValidDevices) > 0 {
			return nil, fmt.Errorf("invalid device ids: %v", invalidDevIDs)
		}
	}

	var invalidServiceTags []string
	if len(serviceTags) > 0 {
		for _, serviceTag := range serviceTags {
			device, err := c.GetDevice(serviceTag, 0)
			if err != nil {
				inValidDevices = append(inValidDevices, device)
				invalidServiceTags = append(invalidServiceTags, serviceTag)
			} else {
				validDevices = append(validDevices, device)
			}
		}
		if len(invalidServiceTags) > 0 {
			return nil, fmt.Errorf("invalid service tags: %v", invalidServiceTags)
		}
	}

	var err error
	var devices models.Devices
	if len(groupNames) > 0 {
		for _, groupName := range groupNames {
			devices, err = c.GetDevicesByGroupName(groupName)
			if err != nil && len(devices.Value) == 0 {
				return []models.Device{}, err
			}
			validDevices = append(validDevices, devices.Value...)
		}
	}

	if len(validDevices) > 0 {
		uniqueDevices := c.GetUniqueDevices(validDevices)
		return uniqueDevices, err
	}

	return []models.Device{}, fmt.Errorf("unable to fetch valid device ids")
}

// GetUniqueDevices return the unique device from a list of a devices
func (c *Client) GetUniqueDevices(devices []models.Device) []models.Device {
	keys := make(map[int64]bool)
	uniqueDevices := []models.Device{}
	for _, device := range devices {
		if _, value := keys[device.ID]; !value {
			keys[device.ID] = true
			uniqueDevices = append(uniqueDevices, device)

		}
	}
	return uniqueDevices
}

// GetUniqueDevicesIdsAndServiceTags return the unique device from a list of a devices
func (c *Client) GetUniqueDevicesIdsAndServiceTags(devices []models.Device) ([]models.Device, []int64, []string) {
	keys := make(map[int64]bool)
	uniqueDevices := []models.Device{}
	uniqueDevicesIDs := []int64{}
	uniqueDevicesSTs := []string{}
	for _, device := range devices {
		if _, value := keys[device.ID]; !value {
			keys[device.ID] = true
			uniqueDevices = append(uniqueDevices, device)
			uniqueDevicesIDs = append(uniqueDevicesIDs, device.ID)
			uniqueDevicesSTs = append(uniqueDevicesSTs, device.DeviceServiceTag)
		}
	}
	return uniqueDevices, uniqueDevicesIDs, uniqueDevicesSTs
}

// GetDeviceByIps - method to get device using ips in OME
func (c *Client) GetDeviceByIps(networks []string) ([]models.Device, error) {
	devices, err := c.GetAllDevices(nil)
	if err != nil {
		return make([]models.Device, 0), err
	}
	return FilterDeviceByIps(devices.Value, networks)
}

// FilterDeviceByIps - method to filter device using ips
func FilterDeviceByIps(devices []models.Device, networks []string) ([]models.Device, error) {
	ret := make([]models.Device, 0)
	pool, err := utils.ParseNetworks(networks)
	if err != nil {
		return ret, err
	}
	for _, v := range devices {
		if v.BelongsToPool(pool) {
			ret = append(ret, v)
		}
	}
	return ret, err
}

// GetAllDevices - method to fetch all devices filtered by input queries
func (c *Client) GetAllDevices(queries map[string]string) (models.Devices, error) {
	devices := models.Devices{}
	err := c.GetValueWithPagination(RequestOptions{
		URL:         DeviceAPI,
		QueryParams: queries,
	}, &devices.Value)
	return devices, err
}

// GetValidDevicesByNames retrieves devices based on their names.
func (c *Client) GetValidDevicesByNames(names []string) ([]models.Device, error) {
	// Retrieve all devices
	allDevices, err := c.GetAllDevices(nil)
	if err != nil {
		return nil, err
	}

	// Filter devices based on names
	var filteredDevices []models.Device
	for _, device := range allDevices.Value {
		for _, name := range names {
			if device.DeviceName == name {
				filteredDevices = append(filteredDevices, device)
				break
			}
		}
	}
	if len(filteredDevices) == 0 {
		return nil, fmt.Errorf("no devices found")
	}
	return filteredDevices, nil
}
