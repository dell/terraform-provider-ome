package clients

import (
	"fmt"
	"terraform-provider-ome/models"
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

	if err == nil {
		devices := models.Devices{}
		bodyData, _ := c.GetBodyData(response.Body)
		err = c.JSONUnMarshal(bodyData, &devices)
		if err == nil {
			if len(devices.Value) > 0 {
				device = devices.Value[0]
				err = nil
			} else {
				err = fmt.Errorf(ErrInvalidDeviceIdentifiers+" %s", val)
			}
		}
	}
	return device, err
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
