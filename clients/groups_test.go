package clients

import (
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientGetGroupIdByGroupName(t *testing.T) {
	tests := []struct {
		name      string
		groupName string
	}{
		{"GetGroupIdByGroupName - valid names", "valid_group1"},
		{"GetGroupIdByGroupName - invalid group names", "invalid_group1"},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetGroupByName(tt.groupName)
			if tt.groupName == "valid_group1" {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, 1, len(response.Value))
				assert.Equal(t, int64(1011), response.Value[0].ID)
				assert.Equal(t, "Linux Servers", response.Value[0].Name)
			}
			if tt.groupName == "invalid_group1" {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Zero(t, len(response.Value))
			}
		})
	}
}

func TestClientGetDevicesByGroupID(t *testing.T) {
	tests := []struct {
		name    string
		groupID int64
		devices models.Devices
	}{
		{"GetDevicesByGroupID - valid groupID", 1011, models.Devices{Value: []models.Device{{ID: 10337}}}},
		{"GetDevicesByGroupID - valid groupID with pagination", 1013, models.Devices{Value: []models.Device{{ID: 10337}, {ID: 10338}}}},
		{"GetDevicesByGroupID - valid groupID with failure in pagination", 1014, models.Devices{Value: []models.Device{{ID: 10337}}}},
		{"GetDevicesByGroupID - valid groupID with failure in pagination", 1015, models.Devices{Value: []models.Device{}}},
		{"GetDevicesByGroupID - valid groupID with failure in pagination", 1016, models.Devices{Value: []models.Device{{ID: 10337}}}},
		{"GetDevicesByGroupID - invalid groupID", -1, models.Devices{}},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetDevicesByGroupID(tt.groupID)
			if tt.groupID == 1011 {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, 1, len(response.Value))
				assert.Equal(t, tt.devices.Value[0].ID, response.Value[0].ID)
			} else if tt.groupID == 1013 {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, 2, len(response.Value))
				assert.Equal(t, tt.devices.Value[0].ID, response.Value[0].ID)
				assert.Equal(t, tt.devices.Value[1].ID, response.Value[1].ID)
			} else if tt.groupID == 1014 || tt.groupID == 1016 {
				assert.NotNil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, 1, len(response.Value))
				assert.Equal(t, tt.devices.Value[0].ID, response.Value[0].ID)
			} else if tt.groupID == -1 || tt.groupID == 1015 {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_GetDevicesByGroupName(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name      string
		groupName string
	}{
		{"Get devices by group name Successfully", "valid_group1"},
		{"Get empty devices for invalid group name", "invalid_group1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devices, err := c.GetDevicesByGroupName(tt.groupName)
			if tt.groupName == "valid_group1" {
				assert.Nil(t, err)
				assert.Equal(t, int64(10337), devices.Value[0].ID)
				assert.Equal(t, 1, len(devices.Value))
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, devices.Value)
			}
		})
	}
}

func TestClient_GetDevicesByGroupIDAndNameUnAuth(t *testing.T) {
	ts := createNewTLSServerWithPort(t, 8235, mockPortUnAuth)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	_, err := c.GetDevicesByGroupID(123456)
	assert.NotNil(t, err)

	_, err = c.GetDevicesByGroupName("123456")
	assert.NotNil(t, err)

	response, err := c.GetGroupByName("invalid_group_id")
	assert.NotNil(t, err)
	assert.Empty(t, response.Value)
}

func TestClient_GetDevicesByGroupIDAndNameInvalidJson(t *testing.T) {
	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	_, err := c.GetDevicesByGroupID(123456)
	assert.NotNil(t, err)

	_, err = c.GetDevicesByGroupName("123456")
	assert.NotNil(t, err)

	_, err = c.GetGroupByName("invalid_group_id")
	assert.NotNil(t, err)
}

func TestClient_GetDevicesByGroups(t *testing.T) {
	type args struct {
		groupNames []string
	}
	tests := []struct {
		name        string
		args        args
		expected    []models.Device
		expectError bool
	}{
		{"GetDevicesByGroups - get devices for multiple groups", args{[]string{"valid_group1", "valid_group2"}}, []models.Device{
			{
				ID:               10337,
				DeviceServiceTag: "SvcTag-1",
			},
			{
				ID:               10338,
				DeviceServiceTag: "SvcTag-2",
			},
		}, false},
		{"GetDeviceByGroups -  get devices for valid/invalid groups", args{[]string{"valid_group1", "invalid_group1"}}, []models.Device{}, true},
		{"GetDeviceByGroups -  get devices for only invalid groups", args{[]string{"invalid_group1"}}, []models.Device{}, true},
	}
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetDevicesByGroups(tt.args.groupNames)
			if tt.expectError {
				assert.NotNil(t, err)
				assert.Empty(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
