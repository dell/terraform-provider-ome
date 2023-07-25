package clients

import (
	_ "embed"
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed json_data/payloadCreateDiscovery.json
	jsonData1 []byte
	// go:embed json_data/payloadUpdateDiscovery.json
	jsonData2 []byte
	//go:embed json_data/payloadDeleteDiscovery.json
	jsonData3 []byte
)

func TestClient_DiscoveryCreateJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	var createDiscoveryJobPayloadSuccess models.DiscoveryJobPayload

	err := c.JSONUnMarshal(jsonData1, &createDiscoveryJobPayloadSuccess)
	if err != nil {
		t.Error(err)
	}
	t.Log(createDiscoveryJobPayloadSuccess.CreateGroup)
	tests := []struct {
		name string
		args models.DiscoveryJobPayload
	}{
		{"Create Discovery Job Successfully", createDiscoveryJobPayloadSuccess},
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

	err := c.JSONUnMarshal(jsonData2, &updateDiscoveryJobSuccess)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args models.DiscoveryJob
	}{
		{"Update Discovery Job Successfully", updateDiscoveryJobSuccess},
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

	err := c.JSONUnMarshal(jsonData3, &deleteDiscoveryJobPayloadSuccess)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		args models.DiscoveryJobDeletePayload
	}{
		{"Delete Discovery Job Successfully", deleteDiscoveryJobPayloadSuccess},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.DeleteDiscoveryJob(deleteDiscoveryJobPayloadSuccess)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.GetDiscoveryJobByGroupID(51)
			t.Log(resp)
			if err == nil {
				assert.Equal(t, resp.DiscoveryConfigGroupID, 51)
			}
		})
	}
}
