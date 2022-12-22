package clients

import (
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateBaseline(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name string
		args models.ConfigurationBaselinePayload
	}{
		{"Create Baseline Successfully", models.ConfigurationBaselinePayload{
			Name:        "TestAccCreateBaseline",
			Description: "Test Acc Description for create baseline",
			TemplateID:  326,
			BaselineTargets: []models.BaselineTarget{
				{
					ID: 10093,
					Type: models.BaselineTargetType{
						ID:   1,
						Name: "DEVICE",
					},
				},
				{
					ID: 10104,
					Type: models.BaselineTargetType{
						ID:   1,
						Name: "DEVICE",
					},
				},
			},
			NotificationSettings: &models.NotificationSettings{
				NotificationType: "NOTIFY_ON_SCHEDULE",
				EmailAddresses:   []string{"test@testdell.com"},
				Schedule: models.BaselineNotificationSchedule{
					Cron: "0 00 00 * * ? *",
				},
				OutputFormat: "HTML",
			},
		}},
		{"Create Baseline Failure - Invalid template ID", models.ConfigurationBaselinePayload{
			Name:        "TestAccCreateBaselineFailure",
			Description: "Test Acc Description for create baseline failure",
			TemplateID:  -1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseline, err := c.CreateBaseline(tt.args)
			if err != nil {
				assert.NotNil(t, err)
				assert.Empty(t, baseline.ID)
				assert.ErrorContains(t, err, "Unable to process the request because the template ID -1 provided is invalid.")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.args.Name, baseline.Name)
				assert.NotNil(t, baseline.ID)
				assert.Equal(t, 2, len(baseline.BaselineTargets))
				assert.Equal(t, "NOTIFY_ON_SCHEDULE", baseline.NotificationSettings.NotificationType)
				assert.Equal(t, 1, len(baseline.NotificationSettings.EmailAddresses))
				assert.Equal(t, "0 00 00 * * ? *", baseline.NotificationSettings.Schedule.Cron)
				assert.Equal(t, "HTML", baseline.NotificationSettings.OutputFormat)
			}
		})
	}
}

func TestClient_UpdateBaseline(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name string
		args models.ConfigurationBaselinePayload
	}{
		{"Update Baseline Successfully", models.ConfigurationBaselinePayload{
			ID:          100,
			Name:        "TestAccCreateBaseline",
			Description: "Test Acc Description for create baseline",
			TemplateID:  326,
			BaselineTargets: []models.BaselineTarget{
				{
					ID: 10093,
					Type: models.BaselineTargetType{
						ID:   1,
						Name: "DEVICE",
					},
				},
				{
					ID: 10104,
					Type: models.BaselineTargetType{
						ID:   1,
						Name: "DEVICE",
					},
				},
			},
			NotificationSettings: &models.NotificationSettings{
				NotificationType: "NOTIFY_ON_SCHEDULE",
				EmailAddresses:   []string{"test@testdell.com"},
				Schedule: models.BaselineNotificationSchedule{
					Cron: "0 00 00 * * ? *",
				},
				OutputFormat: "HTML",
			},
		}},
		{"Update Baseline Failure - Invalid template ID", models.ConfigurationBaselinePayload{
			ID:          101,
			Name:        "TestAccCreateBaselineFailure",
			Description: "Test Acc Description for create baseline failure",
			TemplateID:  -1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseline, err := c.UpdateBaseline(tt.args)
			if err != nil {
				assert.NotNil(t, err)
				assert.Empty(t, baseline.ID)
				assert.ErrorContains(t, err, "Unable to process the request because the template ID -1 provided is invalid.")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.args.Name, baseline.Name)
				assert.NotNil(t, baseline.ID)
				assert.Equal(t, 2, len(baseline.BaselineTargets))
				assert.Equal(t, "NOTIFY_ON_SCHEDULE", baseline.NotificationSettings.NotificationType)
				assert.Equal(t, 1, len(baseline.NotificationSettings.EmailAddresses))
				assert.Equal(t, "0 00 00 * * ? *", baseline.NotificationSettings.Schedule.Cron)
				assert.Equal(t, "HTML", baseline.NotificationSettings.OutputFormat)
			}
		})
	}
}

func TestClient_DeleteBaseline(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name  string
		args  []int64
		isErr bool
	}{
		{"Delete Baseline Successfully", []int64{10001}, false},
		{"Delete Baseline Failure - Invalid template ID", []int64{10002}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.DeleteBaseline(tt.args)
			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestClient_GetBaselineByID(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name       string
		baselineID int64
	}{
		{"Get Baseline By ID Successfully", 1},
		{"Get Baseline By ID Failure", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseline, err := c.GetBaselineByID(tt.baselineID)
			if tt.baselineID == -1 {
				assert.NotNil(t, err)
				assert.Empty(t, baseline.ID)
				assert.ErrorContains(t, err, "Unable to process the request because an error occurred.")
			} else {
				assert.Nil(t, err)
				assert.Equal(t, "Baseline Name", baseline.Name)
				assert.NotNil(t, baseline.ID)
				assert.Equal(t, 2, len(baseline.BaselineTargets))
				assert.Equal(t, "NOTIFY_ON_NON_COMPLIANCE", baseline.NotificationSettings.NotificationType)
				assert.Equal(t, 1, len(baseline.NotificationSettings.EmailAddresses))
				assert.Equal(t, "0 00 00 * * ? *", baseline.NotificationSettings.Schedule.Cron)
				assert.Equal(t, "html", baseline.NotificationSettings.OutputFormat)
			}
		})
	}
}

func TestClient_GetBaselineDeviceComplianceReportByID(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name       string
		baselineID int64
		deviceID   int64
	}{
		{"Get Baseline Device Compliance Reports By ID Successfully", 14, 11803},
		{"Get Baseline Device Compliance Reports By ID Failure", -1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baselineDevComplianceReport, err := c.GetBaselineDevComplianceReportsByID(tt.baselineID)
			if tt.baselineID == 14 {
				assert.Nil(t, err)
				assert.NotEmpty(t, baselineDevComplianceReport)
				assert.Equal(t, baselineDevComplianceReport[0].ID, tt.deviceID)
			} else if tt.baselineID == -1 {
				assert.NotNil(t, err)
				assert.Empty(t, baselineDevComplianceReport)
			}

		})
	}
}

func TestClient_GetBaselineDeviceAttrComplianceReportByID(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name       string
		baselineID int64
		deviceID   int64
	}{
		{"Get Baseline Device Compliance Attribute Reports By ID Successfully", 14, 11803},
		{"Get Baseline Device Compliance Attribute Reports By ID Failure", -1, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baselineDevAttrComplianceReportStr, err := c.GetBaselineDevAttrComplianceReportsByID(tt.baselineID, tt.deviceID)
			if tt.baselineID == 14 && tt.deviceID == 11803 {
				assert.Nil(t, err)
				assert.NotEmpty(t, baselineDevAttrComplianceReportStr)
				assert.Contains(t, baselineDevAttrComplianceReportStr, "\"DeviceId\": 11803")
				assert.Contains(t, baselineDevAttrComplianceReportStr, "\"BaselineId\": 14")
			} else if tt.baselineID == -1 && tt.deviceID == -1 {
				assert.NotNil(t, err)
				assert.Empty(t, baselineDevAttrComplianceReportStr)
			}

		})
	}
}

func TestClient_GetBaselineByName(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name         string
		baselineName string
	}{
		{"Get Baseline By Name from first page in response successfully", "TestAccBaseline1"},
		{"Get Baseline By Name from second page in response successfully", "TestAccBaseline2"},
		{"Get Baseline By Name from third page in response successfully", "TestAccBaseline3"},
		{"Get Baseline By Invalid Name failure", "test_acc_invalid_baseline"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseline, err := c.GetBaselineByName(tt.baselineName)
			if tt.baselineName == "test_acc_invalid_baseline" {
				assert.NotNil(t, err)
				assert.Equal(t, models.OmeBaseline{}, baseline)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.baselineName, baseline.Name)
			}
		})
	}
}

func TestBaselineUnAuth(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8235, mockPortUnAuth)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	response, err := c.GetBaselineByName("unauth_baseline1")
	assert.NotNil(t, err)
	assert.Equal(t, "", response.Name)
}

func TestBaselineInvalidJson(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	response, err := c.GetBaselineByName("invalid_json_baseline1")
	assert.NotNil(t, err)
	assert.Equal(t, "", response.Name)
}

func TestClient_RemediateBaseLineDevices(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name    string
		cr      models.ConfigurationRemediationPayload
		want    int64
		wantErr bool
	}{
		{"RemediateBaseline Success", models.ConfigurationRemediationPayload{
			ID:        100,
			DeviceIDS: []int64{12345},
			Schedule: models.OMESchedule{
				RunNow:   true,
				RunLater: false,
			},
		}, 12345, false},
		{"RemediateBaseline Failure", models.ConfigurationRemediationPayload{
			ID:        101,
			DeviceIDS: []int64{12345},
			Schedule: models.OMESchedule{
				RunNow:   true,
				RunLater: false,
			},
		}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.RemediateBaseLineDevices(tt.cr)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_ConfiBaselineDeviceReport(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name       string
		baseLineID int64
		want       []models.OMEDeviceComplianceReport
		wantErr    bool
	}{
		{"Get Config Device Comp Report", 185, []models.OMEDeviceComplianceReport{
			{ID: 12328},
			{ID: 12329},
			{ID: 12330},
		}, false},
		{"Get Config Device Comp Report error", 186, []models.OMEDeviceComplianceReport{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetAllConfiBaselineDeviceReport(tt.baseLineID)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, len(tt.want), len((got)))
			}
		})
	}
}
