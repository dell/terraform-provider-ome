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

func TestClient_GetAllVlanNetworks(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name    string
		want    []models.VLanNetworks
		wantErr bool
	}{
		// {"Successfully return page 2 data", []models.VLanNetworks{{ID: 1234, Name: "VLAN1"}, {ID: 1235, Name: "VLAN2"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetAllVlanNetworks()
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_GetAllVlanNetworksAnAuth(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	resp, err := c.GetAllVlanNetworks()
	assert.NotNil(t, err)
	assert.Equal(t, []models.VLanNetworks{}, resp)
}
