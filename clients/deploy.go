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
	"fmt"
	"strconv"
	"terraform-provider-ome/models"
)

// CreateDeployment creates a deployment for a specific template
func (c *Client) CreateDeployment(deploymentRequest models.OMETemplateDeployRequest) (int64, error) {
	data, errMarshal := c.JSONMarshal(deploymentRequest)
	if errMarshal != nil {
		return -1, errMarshal
	}
	response, err := c.Post(DeployAPI, nil, data)
	if err != nil {
		return -1, err
	}

	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return -1, getBodyError
	}
	val, parseErr := strconv.ParseInt(string(respData), 10, 64)
	if parseErr != nil {
		return -1, parseErr
	}
	return val, nil
}

// GetServerProfileInfoByTemplateName returns the profile information for a templateName
func (c *Client) GetServerProfileInfoByTemplateName(name string) (models.OMEServerProfiles, error) {
	omeServerProfileResp := []models.OMEServerProfile{}
	err := c.GetPaginatedDataWithQueryParam(ProfileAPI, map[string]string{"$filter": fmt.Sprintf("%s eq '%s'", "TemplateName", name)}, &omeServerProfileResp)
	if err != nil {
		return models.OMEServerProfiles{}, err
	}
	if len(omeServerProfileResp) == 0 {
		return models.OMEServerProfiles{}, nil
	}
	omeServerProfileFilteredResp := []models.OMEServerProfile{}
	for _, filteredServerProfile := range omeServerProfileResp {
		if filteredServerProfile.TemplateName == name {
			omeServerProfileFilteredResp = append(omeServerProfileFilteredResp, filteredServerProfile)
		}
	}
	return models.OMEServerProfiles{Value: omeServerProfileFilteredResp}, nil
}

// DeleteDeployment unassigns and deletes the profile corresponding to the deployment
func (c *Client) DeleteDeployment(deleteDeploymentReq models.ProfileDeleteRequest) error {
	data, errMarshal := c.JSONMarshal(&deleteDeploymentReq)
	if errMarshal != nil {
		return errMarshal
	}
	response, err := c.Post(UnAssignProfileAPI, nil, data)
	if err != nil {
		return err
	}

	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return getBodyError
	}
	jobID, parseErr := strconv.ParseInt(string(respData), 10, 64)
	if parseErr != nil {
		return parseErr
	}

	if jobID != 0 {
		jobStatus, statusMessage := c.TrackJob(jobID, 10, 10)
		if !jobStatus {
			return fmt.Errorf("%s", statusMessage)
		}
	}
	_, err = c.Post(DeleteProfileAPI, nil, data)
	if err != nil {
		return err
	}
	return nil
}
