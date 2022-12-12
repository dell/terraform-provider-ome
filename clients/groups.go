package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

// GetGroupByName - method to get a group object by name.
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
	var err error
	for _, groupName := range groupNames {
		groupDevices, err := c.GetDevicesByGroupName(groupName)
		if err != nil && len(groupDevices.Value) == 0 {
			return []models.Device{}, err
		}
		devices = append(devices, groupDevices.Value...)
	}
	return devices, err
}
