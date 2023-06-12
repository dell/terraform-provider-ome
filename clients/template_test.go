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
	"reflect"

	"terraform-provider-ome/models"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

const (
	TestDeploymentTemplateID       = 123
	TestClonedDeploymentTemplateID = 124
	TestComplianceTemplateID       = 125
	TestClonedComplianceTemplateID = 126
	DeploymentViewTypeID           = 2
	ComplainceViewTypeID           = 1
)

// TestGetTemplateAttributes
func TestGetTemplateAttributes(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		templateID      int64
		stateAttributes []models.Attribute
	}
	tests := []struct {
		name          string
		args          args
		omeAttributes []models.OmeAttribute
	}{
		{"Get Template Attribute - Attribute Found", args{31, []models.Attribute{
			{
				AttributeID: types.Int64Value(110),
				DisplayName: types.StringValue("BIOS,BIOS Boot Settings,Boot Sequence"),
				Value:       types.StringValue("HardDisk.List.1-1"),
				IsIgnored:   types.BoolValue(false),
			},
			{
				AttributeID: types.Int64Value(120),
				DisplayName: types.StringValue("NIC,NIC.Integrated.1-1-1,iSCSI General Parameters,Boot to Target"),
				Value:       types.StringValue("Enabled"),
				IsIgnored:   types.BoolValue(false),
			},
		}}, []models.OmeAttribute{
			{
				AttributeID: 110,
				DisplayName: "BIOS,BIOS Boot Settings,Boot Sequence",
				Value:       "HardDisk.List.1-1",
				IsIgnored:   false,
			},
			{
				AttributeID: 120,
				DisplayName: "NIC,NIC.Integrated.1-1-1,iSCSI General Parameters,Boot to Target",
				Value:       "Enabled",
				IsIgnored:   false,
			},
		}},
		{"Get Template Attribute - Attribute Not found", args{32, []models.Attribute{
			{
				AttributeID: types.Int64Value(110),
				DisplayName: types.StringValue("BIOS,BIOS Boot Settings,Boot Test Sequence"),
				Value:       types.StringValue("HardDisk.List.1-1"),
				IsIgnored:   types.BoolValue(false),
			},
		}}, nil},
		{"Get Template Attribute - Attribute Group Not found", args{33, []models.Attribute{
			{
				AttributeID: types.Int64Value(110),
				DisplayName: types.StringValue("BIOS TEST,BIOS Boot Settings,Boot Test Sequence"),
				Value:       types.StringValue("HardDisk.List.1-1"),
				IsIgnored:   types.BoolValue(false),
			},
		}}, nil},
		{"Get Template Attribute - SubAttribute Group Not found", args{34, []models.Attribute{
			{
				AttributeID: types.Int64Value(110),
				DisplayName: types.StringValue("BIOS,BIOS TEST Boot Settings,Boot Test Sequence"),
				Value:       types.StringValue("HardDisk.List.1-1"),
				IsIgnored:   types.BoolValue(false),
			},
		}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetTemplateAttributes(tt.args.templateID, tt.args.stateAttributes, false)
			if tt.args.templateID == 31 {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.omeAttributes[0].AttributeID, response[0].AttributeID)
				assert.Equal(t, tt.omeAttributes[0].DisplayName, response[0].DisplayName)
				assert.Equal(t, tt.omeAttributes[0].Value, response[0].Value)
				assert.Equal(t, tt.omeAttributes[0].IsIgnored, response[0].IsIgnored)
			}
			if tt.args.templateID == 32 {
				assert.NotNil(t, err)
				assert.Equal(t, 0, len(response))
				assert.ErrorContains(t, err, "attribute could not be found")
			}
			if tt.args.templateID == 33 {
				assert.NotNil(t, err)
				assert.Equal(t, 0, len(response))
				assert.ErrorContains(t, err, "AttributeGroup could not be found")
			}
			if tt.args.templateID == 34 {
				assert.NotNil(t, err)
				assert.Equal(t, 0, len(response))
				assert.ErrorContains(t, err, "SubAttributeGroup could not be found")
			}
		})
	}
}

func TestGetTemplateAttributes_refreshAll(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		templateID      int64
		stateAttributes []models.Attribute
	}
	tests := []struct {
		name          string
		args          args
		omeAttributes []models.OmeAttribute
	}{
		{"Get Template Attribute - Get all attributes", args{35, []models.Attribute{}}, []models.OmeAttribute{
			{
				AttributeID: 110,
				DisplayName: "BIOS,BIOS Boot Settings,Boot Sequence",
				Value:       "HardDisk.List.1-1",
				IsIgnored:   false,
			},
			{
				AttributeID: 120,
				DisplayName: "NIC,NIC.Integrated.1-1-1,iSCSI General Parameters,Boot to Target",
				Value:       "Enabled",
				IsIgnored:   false,
			},
		}},
		{"Get Template Attribute - Get all attributes - Invalid template", args{36, []models.Attribute{}}, []models.OmeAttribute{}},
		{"Get Template Attribute - Get all attributes - Invalid json", args{37, []models.Attribute{}}, []models.OmeAttribute{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetTemplateAttributes(tt.args.templateID, tt.args.stateAttributes, true)
			if tt.args.templateID == 35 {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, 2, len(response))
				assert.Equal(t, tt.omeAttributes[0].AttributeID, response[0].AttributeID)
				assert.Equal(t, tt.omeAttributes[0].DisplayName, response[0].DisplayName)
				assert.Equal(t, tt.omeAttributes[0].Value, response[0].Value)
				assert.Equal(t, tt.omeAttributes[0].IsIgnored, response[0].IsIgnored)
				assert.Equal(t, tt.omeAttributes[1].AttributeID, response[1].AttributeID)
				assert.Equal(t, tt.omeAttributes[1].DisplayName, response[1].DisplayName)
				assert.Equal(t, tt.omeAttributes[1].Value, response[1].Value)
				assert.Equal(t, tt.omeAttributes[1].IsIgnored, response[1].IsIgnored)
			}
			if tt.args.templateID == 36 {
				assert.NotNil(t, err)
				assert.Equal(t, 0, len(response))

			}
		})
	}
}

func TestClient_GetIdentityPoolByNameUnAuth(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8235, mockPortUnAuth)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	response, err := c.GetIdentityPoolByName("1234")
	assert.NotNil(t, err)
	assert.Equal(t, "", response.Name)

	response, err = c.GetIdentityPoolByID(1234)
	assert.NotNil(t, err)
	assert.Equal(t, "", response.Name)
}

func TestClient_GetIdentityPoolInvalidJson(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	response, err := c.GetIdentityPoolByName("2234")
	assert.NotNil(t, err)
	assert.Equal(t, "", response.Name)

	response, err = c.GetIdentityPoolByID(2234)
	assert.NotNil(t, err)
	assert.Equal(t, "", response.Name)
}

func TestClient_GetIdentityPoolByName(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		IdentityPoolName string
	}
	tests := []struct {
		name         string
		args         args
		identityPool models.IdentityPool
		errorMessage string
	}{
		{"Get IdentityPool By Name - Get IdentityPool ID for valid name", args{"IdPool1"}, models.IdentityPool{
			Name: "IdPool1",
			ID:   1,
		}, ""},
		{"Get IdentityPool By Name - Get IdentityPool ID for invalid name", args{"IdPool2"}, models.IdentityPool{}, "IdentityPool: 'IdPool2' is not available in the appliance"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetIdentityPoolByName(tt.args.IdentityPoolName)
			if tt.args.IdentityPoolName == "IdPool1" {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, int64(1), response.ID)
			}
			if tt.args.IdentityPoolName == "IdPool2" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.errorMessage)
			}
		})
	}
}

func TestClient_GetIdentityPoolByID(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		IdentityPoolID int64
	}
	tests := []struct {
		name         string
		args         args
		identityPool models.IdentityPool
		errorMessage string
	}{
		{"Get IdentityPool By Name - Get IdentityPool for valid id", args{123}, models.IdentityPool{
			Name: "IdPool1",
			ID:   123,
		}, ""},
		{"Get IdentityPool By Name - Get IdentityPool for invalid id", args{124}, models.IdentityPool{}, "Unable to process the request because the Identity Pool ID 124 provided is invalid."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetIdentityPoolByID(tt.args.IdentityPoolID)
			if tt.errorMessage != "" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.errorMessage)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, int64(123), response.ID)
			}
		})
	}
}

func TestClient_GetNetworkAttributes(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name         string
		templateID   int64
		errorMessage string
	}{
		{"GetNetworkAttributes - For a valid templateID and ViewID", 50, ""},
		{"GetNetworkAttributes - Invalid template", 51,
			"Unable to complete the operation because the value provided for TemplateId is invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetVlanNetworkModel(tt.templateID)
			if tt.templateID == 50 {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, response.NetworkAttributeGroups[0].NetworkAttributes[0].DisplayName, "Nic Bonding Technology")
				assert.Equal(t, response.NetworkAttributeGroups[0].NetworkAttributes[0].Value, "NoTeaming")
				assert.Equal(t, response.NetworkAttributeGroups[1].SubAttributeGroups[0].SubAttributeGroups[0].SubAttributeGroups[0].NetworkAttributes[0].DisplayName, "NIC Bonding Enabled")
				assert.Equal(t, response.NetworkAttributeGroups[1].SubAttributeGroups[0].SubAttributeGroups[0].SubAttributeGroups[0].NetworkAttributes[0].Value, "false")
				assert.Equal(t, response.NetworkAttributeGroups[1].SubAttributeGroups[0].SubAttributeGroups[0].SubAttributeGroups[0].NetworkAttributes[1].DisplayName, "Vlan Tagged")
				assert.Equal(t, response.NetworkAttributeGroups[1].SubAttributeGroups[0].SubAttributeGroups[0].SubAttributeGroups[0].NetworkAttributes[1].Value, "10133")
				assert.Equal(t, response.NetworkAttributeGroups[1].SubAttributeGroups[0].SubAttributeGroups[0].SubAttributeGroups[0].NetworkAttributes[2].DisplayName, "Vlan UnTagged")
				assert.Equal(t, response.NetworkAttributeGroups[1].SubAttributeGroups[0].SubAttributeGroups[0].SubAttributeGroups[0].NetworkAttributes[2].Value, "0")
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, len(response.NetworkAttributeGroups), 0)
				assert.ErrorContains(t, err, "Unable to complete the operation because the value provided for TemplateId is invalid")
			}
		})
	}

}

func TestClient_GetNetworkAttributesUnAuth(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8235, mockPortUnAuth)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	response, err := c.GetVlanNetworkModel(1234)
	assert.NotNil(t, err)
	assert.Equal(t, len(response.NetworkAttributeGroups), 0)
}

func TestClient_GetNetworkAttributesInvalidJson(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	response, err := c.GetVlanNetworkModel(2234)
	assert.NotNil(t, err)
	assert.Equal(t, len(response.NetworkAttributeGroups), 0)
}
func TestClient_createTemplate(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name         string
		templateid   int64
		errorMessage string
		args         models.CreateTemplate
	}{
		{"Create Template Successfully", 1, "", models.CreateTemplate{Fqdds: "All",
			ViewTypeID:     2,
			SourceDeviceID: 12345,
			Name:           "TestTemplate"}},
		{"Create Template Existing Fail", -1, "Unable to create the template because the template name ExtTemplate already exists.", models.CreateTemplate{Fqdds: "All",
			ViewTypeID:     2,
			SourceDeviceID: 12345,
			Name:           "ExtTemplate"}},
		{"Create Template Fail", -1, "Unable to create or deploy the template because the device ID 23456 is invalid.", models.CreateTemplate{Fqdds: "All",
			ViewTypeID:     2,
			SourceDeviceID: 23456,
			Name:           "TemplateInvalidDevice"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateID, err := c.CreateTemplate(tt.args)
			if tt.errorMessage == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.errorMessage)
			}
			assert.Equal(t, tt.templateid, templateID)
		})
	}
}

func TestClient_GetViewTypeID(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name       string
		viewType   string
		viewTypeID int64
	}{
		{"Get View Type ID Successfully - Deployment", "Deployment", 2},
		{"Get View Type ID Successfully - Complaince", "Compliance", 1},
		{"Get View Type ID Fail", "test-view-type", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viewTypeID, err := c.GetViewTypeID(tt.viewType)
			if tt.viewTypeID != -1 {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.viewTypeID, viewTypeID)
		})
	}
}

func TestClient_GetViewTypeIDUnAuth(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8235, mockPortUnAuth)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	viewTypeID, err := c.GetViewTypeID("Deployment")
	assert.NotNil(t, err)
	assert.Equal(t, int64(-1), viewTypeID)
}

func TestClient_GetDeviceTypeIDUnAuth(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8235, mockPortUnAuth)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	viewTypeID, err := c.GetDeviceTypeID("Server")
	assert.NotNil(t, err)
	assert.Equal(t, int64(-1), viewTypeID)
}

func TestClient_GetViewTypeIDInvalidJson(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	viewTypeID, err := c.GetViewTypeID("Deployment")
	assert.NotNil(t, err)
	assert.Equal(t, int64(-1), viewTypeID)
}

func TestClient_GetDeviceTypeIDInvalidJson(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	viewTypeID, err := c.GetDeviceTypeID("Server")
	assert.NotNil(t, err)
	assert.Equal(t, int64(-1), viewTypeID)
}

func TestClient_GetDeviceTypeID(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name         string
		deviceType   string
		deviceTypeID int64
	}{
		{"Get Device Type ID Successfully - Server", "Server", 2},
		{"Get Device Type ID Successfully - Chassis", "Chassis", 4},
		{"Get Device Type ID Fail", "test-device-type", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deviceTypeID, err := c.GetDeviceTypeID(tt.deviceType)
			if tt.deviceTypeID != -1 {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.deviceTypeID, deviceTypeID)
		})
	}
}

func TestClient_GetTemplateByID(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name       string
		templateID int64
		template   models.OMETemplate
	}{
		{"Get Template by ID Successfully", 23, models.OMETemplate{
			ID:          23,
			Name:        "Test Template",
			Description: "This is a test template",
		}},
		{"Get Template by ID Fail", 24, models.OMETemplate{}},
		{"Get Template by ID Fail", 25, models.OMETemplate{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := c.GetTemplateByID(tt.templateID)
			if tt.templateID <= 23 {
				assert.Nil(t, err)
				assert.Equal(t, tt.template.ID, template.ID)
				assert.Equal(t, tt.template.Name, template.Name)
				assert.Equal(t, tt.template.Description, template.Description)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, template)
			}
		})
	}
}

func TestClient_GetTemplateByName(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name         string
		templateName string
	}{
		{"Empty Template response", "ValidEmptyTemplate"},
		{"Sinle Template response", "ValidSingleTemplate"},
		{"Multiple Template response", "ValidMultipleTemplate"},
		{"Multiple Template response with pagination", "ValidTemplatePagination"},
		{"Unauthorised Template response", "UnauthorisedTemplate"},
		{"Unmarshal error Template response", "UnmarshalErrTemplate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetTemplateByName(tt.templateName)
			if tt.templateName == "ValidEmptyTemplate" {
				assert.Nil(t, err)
				assert.Equal(t, response.Name, "")
			}
			if tt.templateName == "ValidSingleTemplate" {
				assert.Nil(t, err)
				assert.Equal(t, tt.templateName, response.Name)
			}
			if tt.templateName == "ValidMultipleTemplate" || tt.templateName == "ValidTemplatePagination" {
				assert.Nil(t, err)
				assert.Equal(t, tt.templateName, response.Name)
			}
			if tt.templateName == "UnauthorisedTemplate" || tt.templateName == "UnmarshalErrTemplate" {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_UpdateTemplate(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name       string
		templateID int64
		template   models.UpdateTemplate
	}{
		{"Update Template Successfully", 123, models.UpdateTemplate{
			ID:          123,
			Name:        "Test Template update",
			Description: "This is a test template update",
		}},
		{"Update Template Fail", 124, models.UpdateTemplate{
			ID:          124,
			Name:        "Test Template update",
			Description: "This is a test template update",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.UpdateTemplate(tt.template)
			if tt.templateID != 124 {
				assert.Nil(t, err)
			} else {
				assert.ErrorContains(t, err, "Unable to complete the operation because the requested URI is invalid.")
			}
		})
	}
}

func TestClient_GetPayloadVlanAttribute(t *testing.T) {

	networkViewData := models.NetworkSpecificView{
		ViewID: 4,
		NetworkAttributeGroups: []models.NetworkAttributeGroup{
			{
				GroupNameID: 1005,
				DisplayName: "NicBondingTechnology",
				NetworkAttributes: []models.NetworkAttribute{
					{
						ComponentID: 0,
						DisplayName: "Nic Bonding Technology",
						Value:       "NoTeaming",
					},
				},
			},
			{
				GroupNameID: 1006,
				DisplayName: "NICModel",
				SubAttributeGroups: []models.NetworkAttributeGroup{
					{
						GroupNameID: 1,
						DisplayName: "Integrated NIC 1",
						SubAttributeGroups: []models.NetworkAttributeGroup{
							{
								GroupNameID: 1,
								DisplayName: "Port ",
								SubAttributeGroups: []models.NetworkAttributeGroup{
									{
										GroupNameID: 1,
										DisplayName: "Partition ",
										NetworkAttributes: []models.NetworkAttribute{
											{
												ComponentID: 1049,
												DisplayName: "NIC Bonding Enabled",
												Value:       "false",
											},
											{
												ComponentID: 1049,
												DisplayName: "Vlan Tagged",
												Value:       "10133, 10594",
											},
											{
												ComponentID: 1049,
												DisplayName: "Vlan UnTagged",
												Value:       "0",
											},
										},
									},
								},
							},
						},
					},
				},
				NetworkAttributes: []models.NetworkAttribute{
					{
						ComponentID: 0,
						DisplayName: "Nic Bonding Technology",
						Value:       "NoTeaming",
					},
				},
			},
		},
	}

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		networkView   models.NetworkSpecificView
		nicIdentifier string
		port          int64
	}
	tests := []struct {
		name         string
		args         args
		want         models.PayloadVlanAttribute
		errorMessage string
	}{
		{"Get Vlan Attributes - Successful", args{networkViewData, "Integrated NIC 1", 1}, models.PayloadVlanAttribute{
			ComponentID: 1049,
			Untagged:    0,
			Tagged:      []int64{10133, 10594},
			IsNICBonded: false,
		}, ""},
		{"Get Vlan Attributes - Invalid Nic", args{networkViewData, "NIC Invalid 1", 1}, models.PayloadVlanAttribute{}, ErrInvalidNetworkDetails},
		{"Get Vlan Attributes - Invalid port", args{networkViewData, "Integrated NIC 1", 0}, models.PayloadVlanAttribute{}, ErrInvalidNetworkDetails},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.GetPayloadVlanAttribute(tt.args.networkView, tt.args.nicIdentifier, tt.args.port)
			if tt.errorMessage == "" {
				assert.Nil(t, err)

			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.errorMessage)
			}
			assert.True(t, reflect.DeepEqual(resp, tt.want))
		})
	}
}

func TestGetSchemaVlanData(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	tests := []struct {
		name       string
		templateID int64
		want       models.OMEVlan
	}{
		{"Test GetVlanAttributes successfully", 50, models.OMEVlan{
			PropagateVLAN:     false,
			BondingTechnology: "NoTeaming",
			OMEVlanAttributes: []models.OMEVlanAttribute{
				{
					NicIdentifier: "Integrated NIC 1",
					Port:          1,
					ComponentID:   1049,
					Untagged:      0,
					Tagged:        []int64{10133},
					IsNICBonded:   false,
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.GetSchemaVlanData(tt.templateID)
			assert.NotNil(t, resp)
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(resp, tt.want))
		})
	}
}

func TestClient_GetSchemaVlanDataUnAuth(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8235, mockPortUnAuth)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	resp, err := c.GetSchemaVlanData(50)
	assert.NotNil(t, err)
	assert.Equal(t, len(resp.OMEVlanAttributes), 0)
}

func TestClient_GetSchemaVlanDataInvalidJson(t *testing.T) {

	ts := createNewTLSServerWithPort(t, 8236, mockPortInValidJSON)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	resp, err := c.GetSchemaVlanData(50)
	assert.NotNil(t, err)
	assert.Equal(t, len(resp.OMEVlanAttributes), 0)
}

func TestUpdateNetworkConfig(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	tests := []struct {
		name         string
		nwConfig     *models.UpdateNetworkConfig
		errormessage string
	}{
		{"Test Update network config successfully", &models.UpdateNetworkConfig{
			TemplateID:        50,
			BondingTechnology: "NoTeaming",
			PropagateVLAN:     true,
			VLANAttributes: []models.PayloadVlanAttribute{
				{
					ComponentID: 1055,
					Untagged:    0,
					Tagged:      []int64{10122},
					IsNICBonded: false,
				},
			},
		}, ""},
		{"Test Update network config: Error", &models.UpdateNetworkConfig{
			TemplateID:        51,
			BondingTechnology: "None",
			PropagateVLAN:     true,
			VLANAttributes: []models.PayloadVlanAttribute{
				{
					ComponentID: 1055,
					Untagged:    0,
					Tagged:      []int64{0},
					IsNICBonded: false,
				},
			},
		}, "invalid value is entered"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.UpdateNetworkConfig(tt.nwConfig)
			if tt.nwConfig.TemplateID == 50 {
				assert.Nil(t, err)
			} else if tt.nwConfig.TemplateID == 51 {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.errormessage)
			}
		})
	}
}

func TestClient_GetTemplateByIDOrName(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	invalidTemplateID := int64(24)

	tests := []struct {
		name         string
		templateID   int64
		templateName string
	}{
		{"Get Template by ID Successfully", 23, "invalid_template_name"},
		{"Get empty template by name response", invalidTemplateID, "ValidEmptyTemplate"},
		{"Get single Template by Name", invalidTemplateID, "ValidSingleTemplate"},
		{"Get multiple Template by Name", invalidTemplateID, "ValidMultipleTemplate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := c.GetTemplateByIDOrName(tt.templateID, tt.templateName)
			if tt.templateID == 23 && (tt.templateID == invalidTemplateID && tt.templateName != "ValidEmptyTemplate") {
				assert.Nil(t, err)
				assert.Equal(t, tt.templateID, template.ID)
				assert.Equal(t, tt.templateName, template.Name)
			}
			if tt.templateID == invalidTemplateID && tt.templateName == "ValidEmptyTemplate" {
				assert.Empty(t, template.ID)
				assert.Empty(t, template.Name)
			}
		})
	}
}

func Test_getStringArrayToArrayInt(t *testing.T) {
	type args struct {
		strSlice []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"String Empty list", args{[]string{}}},
		{"String first element is 0", args{[]string{"0"}}},
		{"String to int conv", args{[]string{"10", "20"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringArrayToArrayInt(tt.args.strSlice)
			if len(tt.args.strSlice) >= 2 {
				assert.Contains(t, got, int64(10))
				assert.Contains(t, got, int64(20))
			} else {
				assert.Empty(t, got)
			}
		})
	}
}

func Test_getAttributeValue(t *testing.T) {
	type args struct {
		nas         []models.NetworkAttribute
		displayName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Attribute empty", args{[]models.NetworkAttribute{
			{ComponentID: 1234, DisplayName: "TestDisplay", Value: "12345"}}, "TestDisplay1"}, ""},
		{"Attribute not empty", args{[]models.NetworkAttribute{
			{ComponentID: 1234, DisplayName: "TestDisplay", Value: "12345"}}, "TestDisplay"}, "12345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAttributeValue(tt.args.nas, tt.args.displayName)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestClient_CloneTemplateByRefTemplateID(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name                 string
		cloneTemplateRequest models.OMECloneTemplate
		errorMessage         string
		newTemplateID        int64
	}{
		{"CloneTemplateByRefTemplateID - From deployment template to deployment", models.OMECloneTemplate{
			SourceTemplateID: TestDeploymentTemplateID,
			NewTemplateName:  "dep-dep-template",
			ViewTypeID:       DeploymentViewTypeID,
		}, "", TestClonedDeploymentTemplateID},
		{"CloneTemplateByRefTemplateID - From deployment template to compliance", models.OMECloneTemplate{
			SourceTemplateID: TestDeploymentTemplateID,
			NewTemplateName:  "dep-comp-template",
			ViewTypeID:       ComplainceViewTypeID,
		}, "", TestClonedComplianceTemplateID},
		{"CloneTemplateByRefTemplateID - From compliance template to compliance", models.OMECloneTemplate{
			SourceTemplateID: TestComplianceTemplateID,
			NewTemplateName:  "comp-comp-template",
			ViewTypeID:       ComplainceViewTypeID,
		}, "", TestClonedComplianceTemplateID},
		{"CloneTemplateByRefTemplateID - Invalid source template", models.OMECloneTemplate{
			SourceTemplateID: -1,
			NewTemplateName:  "test-invalid-template-id",
			ViewTypeID:       ComplainceViewTypeID,
		}, "Unable to clone the template clone example because Source template does not exist.", -1},
		{"CloneTemplateByRefTemplateID - Invalid view type id", models.OMECloneTemplate{
			SourceTemplateID: TestComplianceTemplateID,
			NewTemplateName:  "test-invalid-viewtype-id",
			ViewTypeID:       -1,
		}, "", TestClonedComplianceTemplateID},
		{"CloneTemplateByRefTemplateID - template name already exist", models.OMECloneTemplate{
			SourceTemplateID: TestComplianceTemplateID,
			NewTemplateName:  "test-existing-template-name",
			ViewTypeID:       DeploymentViewTypeID,
		}, "Unable to create the template because the template name test-existing-template-name already exists.", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTemplateID, err := c.CloneTemplateByRefTemplateID(tt.cloneTemplateRequest)
			if tt.cloneTemplateRequest.NewTemplateName == "test-invalid-template-id" || tt.cloneTemplateRequest.NewTemplateName == "test-existing-template-name" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.errorMessage)
				assert.Equal(t, int64(-1), newTemplateID)
			} else {
				assert.Equal(t, tt.newTemplateID, newTemplateID)
				assert.Nil(t, err)
			}
		})
	}
}

func TestClient_ImportTemplate(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name          string
		importRequest models.OMEImportTemplate
		isError       bool
		newTemplateID int64
	}{
		{"OME Import server deployment Template", models.OMEImportTemplate{
			Type:       1,
			Name:       "server-dep-template",
			Content:    "<a>server-dep-template</a>",
			ViewTypeID: DeploymentViewTypeID,
		}, false, 123},
		{"OME Import chassis compliance Template", models.OMEImportTemplate{
			Type:       2,
			Name:       "chassis-dep-template",
			Content:    "<a>chassis-dep-template</a>",
			ViewTypeID: DeploymentViewTypeID,
		}, false, 124},
		{"OME Import server compliance Template", models.OMEImportTemplate{
			Type:       1,
			Name:       "server-comp-template",
			Content:    "<a>server-comp-template</a>",
			ViewTypeID: ComplainceViewTypeID,
		}, false, 125},
		{"OME Import chassis compliance Template", models.OMEImportTemplate{
			Type:       2,
			Name:       "chassis-comp-template",
			Content:    "<a>chassis-comp-template</a>",
			ViewTypeID: ComplainceViewTypeID,
		}, false, 126},
		{"OME Import invalid compliance Template", models.OMEImportTemplate{
			Type:       2,
			Name:       "invalid-template-content",
			Content:    "<a>invalid-template-content<a>",
			ViewTypeID: ComplainceViewTypeID,
		}, true, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTemplateID, err := c.ImportTemplate(tt.importRequest)
			if tt.isError {
				assert.NotNil(t, err)
				assert.Equal(t, tt.newTemplateID, newTemplateID)
			} else {
				assert.Equal(t, tt.newTemplateID, newTemplateID)
				assert.Nil(t, err)
			}
		})
	}
}
