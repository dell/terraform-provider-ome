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
	_ "embed"
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed json_data/payloadCreateDiscovery.json
	payloadCreateDiscovery []byte
	//go:embed json_data/payloadUpdateDiscovery.json
	payloadUpdateDiscovery []byte
	//go:embed json_data/payloadDeleteDiscovery.json
	payloadDeleteDiscovery []byte
)

func TestClient_DiscoveryCreateJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	var createDiscoveryJobPayloadSuccess models.DiscoveryJob

	err := c.JSONUnMarshal(payloadCreateDiscovery, &createDiscoveryJobPayloadSuccess)
	if err != nil {
		t.Error(err)
	}
	t.Log(createDiscoveryJobPayloadSuccess.CreateGroup)
	tests := []struct {
		name string
		args models.DiscoveryJob
	}{
		{"Create Discovery Job Successfully", createDiscoveryJobPayloadSuccess},
		{"Create Discovery Job Failed", models.DiscoveryJob{DiscoveryConfigGroupName: "invalid-create"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discoveryJob, err := c.CreateDiscoveryJob(tt.args)
			t.Log(discoveryJob, err)
			if err == nil {
				assert.Equal(t, createDiscoveryJobPayloadSuccess.DiscoveryConfigGroupName, discoveryJob.DiscoveryConfigGroupName)
			}
		})
	}
}

func TestClient_DiscoveryUpdateJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	var updateDiscoveryJobSuccess models.DiscoveryJob
	t.Log(string(payloadUpdateDiscovery))
	err := c.JSONUnMarshal(payloadUpdateDiscovery, &updateDiscoveryJobSuccess)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args models.DiscoveryJob
	}{
		{"Update Discovery Job Successfully", updateDiscoveryJobSuccess},
		{"Update Discovery Job Failed", models.DiscoveryJob{DiscoveryConfigGroupName: "invalid-update"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discoveryJob, err := c.UpdateDiscoveryJob(tt.args)
			t.Log(discoveryJob, err)
			if err == nil {
				assert.Equal(t, updateDiscoveryJobSuccess.DiscoveryConfigGroupName, discoveryJob.DiscoveryConfigGroupName)
			}
		})
	}
}

func TestClient_DiscoveryDeleteJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	var deleteDiscoveryJobPayloadSuccess models.DiscoveryJobDeletePayload

	err := c.JSONUnMarshal(payloadDeleteDiscovery, &deleteDiscoveryJobPayloadSuccess)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		args models.DiscoveryJobDeletePayload
	}{
		{"Delete Discovery Job Successfully", deleteDiscoveryJobPayloadSuccess},
		{"Delete Discovery Job Failed", models.DiscoveryJobDeletePayload{
			DiscoveryGroupIds: []int{-1},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.DeleteDiscoveryJob(tt.args)
			t.Log(resp)
			if err == nil {
				assert.Equal(t, resp, "204 No Content")
			}
		})
	}
}

func TestClient_DiscoveryGetJobByGroupID(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	tests := []struct {
		name string
		args int64
	}{
		{"Get Discovery Job By Group ID Successfully", 51},
		{"Get Discovery Job By Group ID Failed", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.GetDiscoveryJobByGroupID(tt.args)
			t.Log(resp)
			if err == nil {
				assert.Equal(t, resp.DiscoveryConfigGroupID, tt.args)
			}
		})
	}
}
