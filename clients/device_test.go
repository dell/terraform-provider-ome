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
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetDevice(t *testing.T) {
	type args struct {
		id         int64
		serviceTag string
	}
	tests := []struct {
		name      string
		assertVal int64
		args      args
	}{
		{"test-service-tag-success", 123456, args{0, "SVT123"}},
		{"test-service-tag-failure", -1, args{0, "SV6789"}},
		{"test-service-tag-invalidjson", 5, args{0, "INVJSON"}},
		{"test-service-tag-402-no-auth", 5, args{0, "NOAUTH"}},
		{"test-service-tag-id-empty", 5, args{0, ""}},
		{"test-id-success", 123456, args{123456, ""}},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetDevice(tt.args.serviceTag, tt.args.id)
			if tt.assertVal >= 123456 {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, response.ID, tt.assertVal)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_GetDeviceByIP(t *testing.T) {
	type args struct {
		ips     []string
		ids     []int64
		isError bool
	}
	tests := []struct {
		name string
		args
	}{
		{"test-valid", args{
			ips: []string{
				"192.35.0.1",
				"10.36.0.0-192.36.0.255",
				"fe80::ffff:ffff:ffff:ffff",
				"fe80::ffff:192.0.2.0/125",
				"fe80::ffff:ffff:ffff:1111-fe80::ffff:ffff:ffff:ffff",
				"192.37.0.0/24",
			},
			ids: []int64{123450, 123458},
		}},
		{"test-invalid", args{
			ips:     []string{"192.35.0.344"},
			ids:     []int64{},
			isError: true,
		}},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetDeviceByIps(tt.args.ips)
			if !tt.args.isError {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				for i, id := range tt.args.ids {
					assert.Equal(t, id, response[i].ID)
				}
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClientValidateDevice(t *testing.T) {
	type args struct {
		serviceTag string
		devID      int64
	}
	tests := []struct {
		name      string
		assertVal int64
		args      args
	}{
		{"test-service-tag-success", 123456, args{"SVT123", 0}},
		{"test-service-tag-failure", -1, args{"SV6789", 0}},
		{"test-service-tag-invalidjson", 5, args{"INVJSON", 0}},
		{"test-service-tag-402-no-auth", 5, args{"NOAUTH", 0}},
		{"test-device-id-valid", 123456, args{"", 123456}},
		{"test-device-id-invalid", -1, args{"", 123457}},
		{"test-device-id-invalid", 123456, args{"SV6789", 123456}},
		{"test-invalid-service-tag-device-id", -1, args{"", 0}},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.ValidateDevice(tt.args.serviceTag, tt.args.devID)
			if tt.assertVal >= 123456 {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.assertVal, response)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClientGetDeviceIds(t *testing.T) {
	type args struct {
		serviceTags []string
		devIDs      []int64
		groupNames  []string
	}
	tests := []struct {
		name              string
		expectedDeviceIds []int64
		args              args
	}{
		{"GetDeviceIds - fetch by device IDs", []int64{123456, 223456}, args{[]string{}, []int64{123456, 223456}, []string{}}},
		{"GetDeviceIds - fetch by service tags", []int64{123456}, args{[]string{"SVT123", "SVT223"}, []int64{}, []string{}}},
		{"GetDeviceIds - fetch by group names", []int64{10337, 10338}, args{[]string{}, []int64{}, []string{"valid_group1", "valid_group2"}}},
		{"GetDeviceIds - invalid device IDs", []int64{10337, 10338}, args{[]string{}, []int64{123457}, []string{}}},
		{"GetDeviceIds - invalid service tags", []int64{10337, 10338}, args{[]string{"SV6789"}, []int64{}, []string{}}},
		{"GetDeviceIds - invalid group name", []int64{10337, 10338}, args{[]string{}, []int64{}, []string{"invalid_group1"}}},
		{"GetDeviceIds - empty params", []int64{}, args{[]string{}, []int64{}, []string{}}},
		{"GetDeviceIds - fetch by dev ID, service tag and group name", []int64{123456, 10337}, args{[]string{"SVT123"}, []int64{123456}, []string{"valid_group1"}}},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devices, err := c.GetDevices(tt.args.serviceTags, tt.args.devIDs, tt.args.groupNames)
			if len(tt.args.groupNames) == 0 && len(tt.args.devIDs) == 0 && len(tt.args.serviceTags) == 0 {
				assert.NotNil(t, err)
				assert.Empty(t, devices)
			} else if (len(tt.args.devIDs) > 0 && tt.args.devIDs[0] == 123457) ||
				(len(tt.args.serviceTags) > 0 && tt.args.serviceTags[0] == "SV6789") ||
				(len(tt.args.groupNames) > 0 && tt.args.groupNames[0] == "invalid_group1") {
				assert.NotNil(t, err)
				assert.Empty(t, devices)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedDeviceIds[0], devices[0].ID)
			}
		})
	}
}

func TestClientGetUniqueDevices(t *testing.T) {
	type args struct {
		devices []models.Device
	}

	tests := []struct {
		name            string
		expectedDevices []models.Device
		args            args
	}{
		{"GetUniqueDevices - filterDuplicate devices", []models.Device{
			{
				ID:               1,
				DeviceServiceTag: "TestDevice",
			},
		}, args{[]models.Device{
			{
				ID:               1,
				DeviceServiceTag: "TestDevice",
			},
			{
				ID:               1,
				DeviceServiceTag: "TestDevice",
			},
		}},
		},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devices := c.GetUniqueDevices(tt.args.devices)
			assert.Equal(t, devices, tt.expectedDevices)
		})
	}
}

func TestClientGetUniqueDevicesIdsAndServiceTags(t *testing.T) {
	type args struct {
		devices []models.Device
	}

	tests := []struct {
		name               string
		expectedDevices    []models.Device
		expectedDevicesIds []int64
		expectedDevicesSTs []string
		args               args
	}{
		{"GetUniqueDevices - filterDuplicate devices", []models.Device{
			{
				ID:               1,
				DeviceServiceTag: "TestDevice",
			},
		}, []int64{1}, []string{"TestDevice"}, args{[]models.Device{
			{
				ID:               1,
				DeviceServiceTag: "TestDevice",
			},
			{
				ID:               1,
				DeviceServiceTag: "TestDevice",
			},
		}},
		},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devices, deviceIDs, deviceServiceTags := c.GetUniqueDevicesIdsAndServiceTags(tt.args.devices)
			assert.Equal(t, devices, tt.expectedDevices)
			assert.Equal(t, deviceIDs, tt.expectedDevicesIds)
			assert.Equal(t, deviceServiceTags, tt.expectedDevicesSTs)
		})
	}
}
