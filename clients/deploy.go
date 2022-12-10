package clients

import (
	"fmt"
	"strconv"
	"terraform-provider-ome/models"
)

// CreateDeployment creates a deployment for a specific template
func (c *Client) CreateDeployment(deploymentRequest models.OMETemplateDeployRequest) (int64, error) {
	data, _ := c.JSONMarshal(deploymentRequest)

	response, err := c.Post(DeployAPI, nil, data)
	if err != nil {
		return -1, err
	}

	respData, _ := c.GetBodyData(response.Body)
	val, _ := strconv.ParseInt(string(respData), 10, 64)
	return val, nil
}

// GetServerProfileInfoByTemplateName returns the profile information for a templateName
func (c *Client) GetServerProfileInfoByTemplateName(name string) (models.OMEServerProfiles, error) {
	response, err := c.Get(ProfileAPI, nil, map[string]string{"$filter": fmt.Sprintf("%s eq '%s'", "TemplateName", name)})
	if err != nil {
		return models.OMEServerProfiles{}, err
	}
	b, _ := c.GetBodyData(response.Body)

	omeServerProfiles := models.OMEServerProfiles{}
	err = c.JSONUnMarshal(b, &omeServerProfiles)
	if err != nil {
		return models.OMEServerProfiles{}, err
	}
	if len(omeServerProfiles.Value) == 0 {
		return models.OMEServerProfiles{}, nil
	}
	return omeServerProfiles, nil
}

// DeleteDeployment unassigns and deletes the profile corresponding to the deployment
func (c *Client) DeleteDeployment(deleteDeploymentReq models.ProfileDeleteRequest) error {
	data, _ := c.JSONMarshal(&deleteDeploymentReq)
	response, err := c.Post(UnAssignProfileAPI, nil, data)
	if err != nil {
		return err
	}

	respData, _ := c.GetBodyData(response.Body)
	jobID, _ := strconv.ParseInt(string(respData), 10, 64)

	if jobID != 0 {
		jobStatus, statusMessage := c.TrackJob(jobID, 10, 10)
		if !jobStatus {
			return fmt.Errorf(statusMessage)
		}
	}
	_, err = c.Post(DeleteProfileAPI, nil, data)
	if err != nil {
		return err
	}
	return nil
}
