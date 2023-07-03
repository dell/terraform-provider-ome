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

// GetGroupByName - method to get a group object by name.
func (c *Client) GetGroupById(id int64) (models.Group, error) {
	path := fmt.Sprintf(GroupServiceAPI, id)
	response, err := c.Get(path, nil, nil)
	if err != nil {
		return models.Group{}, err
	}

	bodyData, _ := c.GetBodyData(response.Body)

	group := models.Group{}
	err = c.JSONUnMarshal(bodyData, &group)
	if err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func (c *Client) DeleteGroupDevice(id int64) error {
	path := fmt.Sprintf(GroupServiceAPI, id)
	fmt.Println(path)
	_, err := c.Delete(path, nil, nil)
	return err
}

// GetSingleGroupByName - method to get a single group object by name.
func (c *Client) GetSingleGroupByName(groupName string) (models.Group, error) {
	groups, err := c.GetGroupByName(groupName)
	if err != nil {
		return models.Group{}, nil
	}
	if num := len(groups.Value); num != 1 {
		return models.Group{},
			fmt.Errorf("received %d groups by name %s, while expecting only 1", num, groupName)
	}
	return groups.Value[0], nil
}

// GetGroupByName - method to get a groups object by name.
func (c *Client) GetGroupByName(groupName string) (models.Groups, error) {
	response, err := c.Get(GroupAPI, nil, map[string]string{"Name": groupName})
	if err != nil {
		return models.Groups{}, err
	}
	groups := models.Groups{}
	bodyData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(bodyData, &groups)
	if err != nil {
		return models.Groups{}, err
	}
	return groups, nil
}

// GetDevicesByGroupID - method to get device objects by group id.
func (c *Client) GetDevicesByGroupID(groupID int64) (models.Devices, error) {
	response, err := c.Get(fmt.Sprintf(GroupServiceDevicesAPI, groupID), nil, nil)
	if err != nil {
		return models.Devices{}, err
	}
	allDevices := models.Devices{}
	devices := models.Devices{}
	bodyData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(bodyData, &devices)
	if err != nil {
		return models.Devices{}, err
	}
	allDevices.Value = append(allDevices.Value, devices.Value...)
	for devices.NextLink != "" {
		response, err := c.Get(devices.NextLink, nil, nil)
		if err != nil {
			return allDevices, err
		}
		devices = models.Devices{}
		bodyData, _ := c.GetBodyData(response.Body)
		err = c.JSONUnMarshal(bodyData, &devices)
		if err != nil {
			return allDevices, err
		}
		allDevices.Value = append(allDevices.Value, devices.Value...)
	}
	return allDevices, nil
}

// GetDevicesByGroupName - method to get device objects by group name.
func (c *Client) GetDevicesByGroupName(groupName string) (models.Devices, error) {
	groups, err := c.GetGroupByName(groupName)
	if err != nil {
		return models.Devices{}, err
	}

	if len(groups.Value) == 0 {
		return models.Devices{}, nil
	}
	devices, err := c.GetDevicesByGroupID(groups.Value[0].ID)
	return devices, err
}

// GetDevicesByGroups - returns the list of device by group names
func (c *Client) GetDevicesByGroups(groupNames []string) ([]models.Device, error) {
	devices := []models.Device{}
	for _, groupName := range groupNames {
		groupDevices, err := c.GetDevicesByGroupName(groupName)
		if err != nil && len(groupDevices.Value) == 0 {
			return []models.Device{}, err
		}
		devices = append(devices, groupDevices.Value...)
	}
	return devices, nil
}

// CreateGroupDevice - Creates a new static device group and returns its id
func (c *Client) CreateGroupDevice(group models.Group) (int64, error) {
	payload := map[string]any{
		"GroupModel": group,
	}
	payloadb, err := c.JSONMarshal(payload)
	if err != nil {
		return 0, err
	}
	path := fmt.Sprintf(GroupServiceActionsAPI, "Create")
	response, err2 := c.Post(path, nil, payloadb)
	if err2 != nil {
		return 0, err2
	}
	respData, _ := c.GetBodyData(response.Body)
	val, _ := strconv.ParseInt(string(respData), 10, 64)
	return val, nil
}

// UpdateGroupDevice - Updates a static device group
func (c *Client) UpdateGroupDevice(group models.Group) error {
	payload := map[string]any{
		"GroupModel": group,
	}
	payloadb, err := c.JSONMarshal(payload)
	if err != nil {
		return err
	}
	path := fmt.Sprintf(GroupServiceActionsAPI, "Update")
	_, err2 := c.Post(path, nil, payloadb)
	if err2 != nil {
		return err2
	}
	return nil
}
