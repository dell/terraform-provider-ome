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
