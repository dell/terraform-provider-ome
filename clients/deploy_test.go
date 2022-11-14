package clients

import (
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostDeployTemplate(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name         string
		ProfileJobID int64
		errorMessage string
		request      models.OMETemplateDeployRequest
	}{
		{"Create Deployment Successfully", 1234, "", models.OMETemplateDeployRequest{
			ID:        1,
			TargetIDS: []int64{1, 2, 3},
		}},
		{"Create Deployment Failure - Deployment exist for the device id and template id", 2, "Unable to deploy the template test_deployment because 100.96.24.28 has a profile assigned.",
			models.OMETemplateDeployRequest{
				ID:        2,
				TargetIDS: []int64{1},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.CreateDeployment(tt.request)
			if tt.errorMessage == "" {
				assert.Nil(t, err)
				assert.Equal(t, tt.ProfileJobID, response)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.errorMessage)
			}
		})
	}
}

func TestClient_GetServerProfileInfoByTemplateName(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name         string
		templateName string
	}{
		{"Empty Profile response", "ValidEmptyProfileTemplateName"},
		{"Single Profile response", "ValidSingleProfileTemplateName"},
		{"Multiple Profile response", "ValidMultipleProfileTemplateName"},
		{"Unauthorised Profile response", "UnauthorisedProfileTemplateName"},
		{"Unmarshal error Profile response", "UnmarshalErrProfileTemplateName"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := c.GetServerProfileInfoByTemplateName(tt.templateName)
			if tt.templateName == "ValidEmptyTemplate" {
				assert.Nil(t, err)
				assert.Equal(t, response.Value[0].ID, "")
			}
			if tt.templateName == "ValidSingleTemplate" {
				assert.Nil(t, err)
				assert.Equal(t, 10848, response.Value[0].ID)
			}
			if tt.templateName == "ValidMultipleTemplate" {
				assert.Nil(t, err)
				assert.Equal(t, 10849, response.Value[1].ID)
			}
			if tt.templateName == "UnauthorisedTemplate" || tt.templateName == "UnmarshalErrTemplate" {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_DeleteDeployment(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		deleteDeploymentReq models.ProfileDeleteRequest
	}
	tests := []struct {
		name      string
		args      args
		expectErr bool
	}{
		{"DeleteDeployment - deletes the deployment successfully", args{models.ProfileDeleteRequest{
			ProfileIds: []int64{10850},
		}}, false},
		{"DeleteDeployment - deployment deletion failure (fails unassigning profile)", args{models.ProfileDeleteRequest{
			ProfileIds: []int64{10851},
		}}, true},
		{"DeleteDeployment - deployment deletion failure (fails in deleteProfile API call)", args{models.ProfileDeleteRequest{
			ProfileIds: []int64{10852},
		}}, true},
		{"DeleteDeployment - deployment deletion failure (fails in track job)", args{models.ProfileDeleteRequest{
			ProfileIds: []int64{10853},
		}}, true},
		{"DeleteDeployment - deployment deletion failure (fails in json marshalling)", args{models.ProfileDeleteRequest{
			ProfileIds: nil,
		}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.DeleteDeployment(tt.args.deleteDeploymentReq)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
