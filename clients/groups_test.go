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
	"log"
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
				assert.Equal(t, "valid_group1", response.Value[0].Name)
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
				assert.Nil(t, err)
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
		{"GetDeviceByGroups -  get devices for valid/invalid groups", args{[]string{"valid_group1", "invalid_group1"}}, []models.Device{{
			ID:               10337,
			DeviceServiceTag: "SvcTag-1",
		}}, false},
		{"GetDeviceByGroups -  get devices for only invalid groups", args{[]string{"invalid_group1"}}, []models.Device{}, false},
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

func TestClient_CreateGroup(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		Name          string
		parentGroupID int64
		isValid       bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Create group with existing name", args{"Extroup", 1015, false}},
		{"Create group with non-existing parent id", args{"TestGroup", 1015, false}},
		{"Create group success", args{"TestGroup", 1011, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println("Hola " + tt.name)
			response, err := c.CreateGroup(models.Group{
				Name:        tt.args.Name,
				Description: "dummy",
				ParentID:    tt.args.parentGroupID,
			})
			if tt.args.isValid {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.EqualValues(t, 1012, response)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_ModifyGroup(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		Name    string
		GroupID int64
		isValid bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Update group to existing name", args{"ExtGroup", 1011, false}},
		{"Update non-existing group", args{"TestGroup", 1055, false}},
		{"Update group success", args{"TestGroup", 1011, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.UpdateGroup(models.Group{
				Name:        tt.args.Name,
				Description: "dummy",
				ParentID:    tt.args.GroupID,
			})
			if tt.args.isValid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_UpdateGroupMembers(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		DeviceID int64
		GroupID  int64
		Add      bool
		isValid  bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Add pre existing device to group", args{10056, 1011, true, false}},
		{"Add device to non existent group", args{10057, 1055, true, false}},
		{"Add device to group success", args{10057, 1011, true, true}},
		{"Remove non existent device from group", args{10057, 1011, false, false}},
		{"Remove device from non existent group", args{10056, 1055, false, false}},
		{"Remove device from group success", args{10056, 1011, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.args.Add {
				err = c.AddGroupMembers(models.GroupMemberPayload{
					GroupID:   tt.args.GroupID,
					DeviceIds: []int64{tt.args.DeviceID},
				})
			} else {
				err = c.RemoveGroupMembers(models.GroupMemberPayload{
					GroupID:   tt.args.GroupID,
					DeviceIds: []int64{tt.args.DeviceID},
				})
			}
			if tt.args.isValid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_ReadGroup(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		GroupName string
		GroupID   int64
		toID      bool
		isValid   bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Get group by existing name", args{"valid_group1", 1011, false, true}},
		{"Get group by non existing name", args{"invalid_group1", 0, false, false}},
		{"Get group by existing id", args{"group 1011", 1011, true, true}},
		{"Get group by non existing id", args{"", -1, true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var group models.Group
			if tt.args.toID {
				group, err = c.GetGroupByID(tt.args.GroupID)
			} else {
				group, err = c.GetSingleGroupByName(tt.args.GroupName)
			}
			if tt.args.isValid {
				assert.Nil(t, err)
				assert.EqualValues(t, tt.args.GroupID, group.ID)
				assert.EqualValues(t, tt.args.GroupName, group.Name)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_ReadExpandedGroupByName(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	t.Run("Invalid expansion", func(t *testing.T) {
		groupName := "Dummy"
		_, err := c.GetExpandedGroupByName(groupName, "InvalidExpansion")
		assert.NotNil(t, err)
	})

	t.Run("Invalid name", func(t *testing.T) {
		groupName := "invalid_group1"
		_, err := c.GetExpandedGroupByName(groupName, "")
		t.Logf(err.Error())
		assert.NotNil(t, err)
	})

	t.Run("Valid name and expansion", func(t *testing.T) {
		// t.Log(string(getExpandedGroupResponse))
		groupName := "Dummy"
		group, err := c.GetExpandedGroupByName(groupName, "")
		assert.Nil(t, err)
		assert.EqualValues(t, 1011, group.ID)
		assert.EqualValues(t, "The Dummy group", group.Description)
		assert.EqualValues(t, 2, len(group.SubGroups))
		assert.EqualValues(t, "dummy-test2", group.SubGroups[0].Name)
		assert.EqualValues(t, "dummy-test", group.SubGroups[1].Name)
	})
}

func TestClient_DeleteGroup(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		GroupID int64
		isValid bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Delete existing group", args{1011, true}},
		{"Delete non existing group", args{1055, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.DeleteGroup(tt.args.GroupID)
			if tt.args.isValid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
