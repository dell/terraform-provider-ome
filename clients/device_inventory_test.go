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

func TestClient_GetDeviceInventory(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	var v models.DeviceInventory

	_, err := c.GetDeviceInventory(123000)
	assert.NotNil(t, err)

	v, err = c.GetDeviceInventory(123456)
	assert.Nil(t, err)
	assert.NotEmpty(t, v.DeviceManagement)
	assert.NotEmpty(t, v.DeviceCapabilities)

	_, err = c.GetDeviceInventoryByType(123456, "unknown")
	assert.NotNil(t, err)

	v, err = c.GetDeviceInventoryByType(123456, "serverDeviceCards")
	assert.Nil(t, err)
	assert.NotEmpty(t, v.ServerDeviceCards)
	assert.Empty(t, v.DeviceManagement)
}
