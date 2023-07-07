package clients

import (
	"io/ioutil"
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateDiscoveryJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	var createDiscoveryJobPayloadSuccess models.DiscoveryJobPayload
	jsonData, err := ioutil.ReadFile("json_data/payloadCreateDiscovery.json")
	if err != nil {
		t.Error(err)
	}
	err = c.JSONUnMarshal(jsonData, &createDiscoveryJobPayloadSuccess)
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

func TestClient_UpdateDiscoveryJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	var updateDiscoveryJobPayloadSuccess models.DiscoveryJobPayload
	jsonData, err := ioutil.ReadFile("json_data/payloadUpdateDiscovery.json")
	if err != nil {
		t.Error(err)
	}
	err = c.JSONUnMarshal(jsonData, &updateDiscoveryJobPayloadSuccess)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args models.DiscoveryJobPayload
	}{
		{"Update Discovery Job Successfully", updateDiscoveryJobPayloadSuccess},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discoveryJob, err := c.UpdateDiscoveryJob(tt.args)
			t.Log(discoveryJob, err)
			if err == nil {
				assert.Equal(t, updateDiscoveryJobPayloadSuccess.DiscoveryConfigGroupName, discoveryJob.DiscoveryConfigGroupName)
			}
		})
	}
}

func TestClient_DeleteDiscoveryJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	var deleteDiscoveryJobPayloadSuccess models.DiscoveryJobDeletePayload
	jsonData, err := ioutil.ReadFile("json_data/payloadDeleteDiscovery.json")
	if err != nil {
		t.Error(err)
	}
	err = c.JSONUnMarshal(jsonData, &deleteDiscoveryJobPayloadSuccess)
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

func TestClient_GetDiscoveryJobByGroupID(t *testing.T) {
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

// func IntiateClient(t *testing.T, notMock bool) *Client {
// 	var c *Client
// 	ts := createNewTLSServer(t)
// 	defer ts.Close()

// 	opts := initOptions(ts)

// 	optsPlatform := ClientOptions{
// 		URL:        "",
// 		SkipSSL:    true,
// 		RootCaPath: "",
// 		Timeout:    time.Second * 30,
// 		Retry:      1,
// 		Username:   "admin",
// 		Password:   "Password123!",
// 	}
// 	if notMock {
// 		c, _ = NewClient(optsPlatform)
// 		c.CreateSession()
// 	} else {
// 		c, _ = NewClient(opts)
// 	}
// 	return c
// }
