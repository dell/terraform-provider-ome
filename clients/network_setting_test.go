package clients

import (
	_ "embed"
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed json_data/payloadUpdateNetworkSuccess.json
	payloadUpdateNetworkSuccess []byte
	//go:embed json_data/payloadUpdateNetworkSession.json
	payloadUpdateNetworkSession []byte
	//go:embed json_data/payloadUpdateNetworkTime.json
	payloadUpdateNetworkTime []byte
	//go:embed json_data/payloadUpdateNetworkProxy.json
	payloadUpdateNetworkProxy []byte
)

func TestNetwork_GetNetworkAdapterConfigByInterface(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	tests := []struct {
		name string
		args string
	}{
		{"Get Network Adapter Config Successfully", "ens160"},
		{"Get Network Adapter Config Failed", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNetAdp, err := c.GetNetworkAdapterConfigByInterface(tt.args)
			t.Log(getNetAdp, err)
			if err == nil {
				assert.Equal(t, tt.args, getNetAdp.InterfaceName)
			}
		})
	}
}

func TestNetwork_UpdateNetworkAdapterConfig(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	var updateNetworkAdapterSuccess models.UpdateNetworkAdapterSetting
	t.Logf(string(payloadUpdateNetworkSuccess))
	err := c.JSONUnMarshal(payloadUpdateNetworkSuccess, &updateNetworkAdapterSuccess)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args models.UpdateNetworkAdapterSetting
	}{
		{"Update Network Adapter Config Successfully", updateNetworkAdapterSuccess},
		{"Update Network Adapter Config Failed", models.UpdateNetworkAdapterSetting{InterfaceName: "invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			networkJob, err := c.UpdateNetworkAdapterConfig(tt.args)
			t.Log(networkJob, err)
			if err == nil {
				assert.Equal(t, "OMERealtime_Task", networkJob.JobName)
			}
		})
	}
}

func TestNetwork_GetNetworkSessions(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	tests := []struct {
		name string
	}{
		{"Get Network Sessions Successfully"},
		{"Get Network Sessions Failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNetSession, err := c.GetNetworkSessions()
			t.Log(getNetSession, err)
			if err == nil {
				assert.Equal(t, getNetSession.SessionList[0].SessionType, "GUI")
			}
		})
	}
}

func TestNetwork_UpdateNetworkSessions(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	var sessionPayload []models.SessionInfo
	t.Logf(string(payloadUpdateNetworkSession))
	err := c.JSONUnMarshal(payloadUpdateNetworkSession, &sessionPayload)
	if err != nil {
		t.Error(err)
	}
	var invalidUpdateNetworkSessions []models.SessionInfo
	invalidSession := models.SessionInfo{
		SessionType:    "invalid",
		MaxSessions:    1,
		SessionTimeout: 1000,
	}
	invalidUpdateNetworkSessions = append(invalidUpdateNetworkSessions, invalidSession)
	tests := []struct {
		name string
		args []models.SessionInfo
	}{
		{"Update Network Session Successfully", sessionPayload},
		{"Update Network Session Failed", invalidUpdateNetworkSessions},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			networkSession, err := c.UpdateNetworkSessions(tt.args)
			t.Log(networkSession, err)
			if err == nil {
				assert.Equal(t, networkSession[0].SessionType, "GUI")
			}
		})
	}
}

func TestNetwork_GetTimeConfiguration(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	tests := []struct {
		name string
	}{
		{"Get Network Time Config Successfully"},
		{"Get Network Time Config Failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNetTime, err := c.GetTimeConfiguration()
			t.Log(getNetTime, err)
			if err == nil {
				assert.Equal(t, getNetTime.TimeZone, "TZ_ID_33")
			}
		})
	}
}

func TestNetwork_UpdateTimeConfiguration(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	var payloadTC models.TimeConfig
	t.Logf(string(payloadUpdateNetworkTime))
	err := c.JSONUnMarshal(payloadUpdateNetworkTime, &payloadTC)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args models.TimeConfig
	}{
		{"Update Network Time Successfully", payloadTC},
		{"Update Network Time Failed", models.TimeConfig{TimeZone: "invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			networkTime, err := c.UpdateTimeConfiguration(tt.args)
			t.Log(networkTime, err)
			if err == nil {
				assert.Equal(t, networkTime.TimeZone, "TZ_ID_65")
			}
		})
	}
}

func TestNetwork_GetTimeZone(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	tests := []struct {
		name string
	}{
		{"Get Network Time Zone Successfully"},
		{"Get Network Time Zone Failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNetTimeZone, err := c.GetTimeZone()
			t.Log(getNetTimeZone, err)
			if err == nil {
				assert.Equal(t, getNetTimeZone.TimeZoneList[0].Name, "TZ_ID_38")
			}
		})
	}
}

func TestNetwork_GetProxyConfig(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	tests := []struct {
		name string
	}{
		{"Get Network Proxy Successfully"},
		{"Get Network Proxy Failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNetProxy, err := c.GetProxyConfig()
			t.Log(getNetProxy, err)
			if err == nil {
				assert.Equal(t, getNetProxy.Username, "admin")
			}
		})
	}
}

func TestNetwork_UpdateProxyConfig(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	var payloadProxy models.PayloadProxyConfiguration
	t.Logf(string(payloadUpdateNetworkProxy))
	err := c.JSONUnMarshal(payloadUpdateNetworkProxy, &payloadProxy)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args models.PayloadProxyConfiguration
	}{
		{"Update Network Time Successfully", payloadProxy},
		{"Update Network Time Failed", models.PayloadProxyConfiguration{Username: "invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			networkProxy, err := c.UpdateProxyConfig(tt.args)
			t.Log(networkProxy, err)
			if err == nil {
				assert.Equal(t, networkProxy.Username, "admin")
			}
		})
	}
}
