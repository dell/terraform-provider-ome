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
	"net/http"
	"strconv"
	"terraform-provider-ome/models"
)

// CreateDiscoveryJob - create a discovery job in OME.
func (c *Client) CreateDiscoveryJob(discoveryJob models.DiscoveryJob) (models.DiscoveryJob, error) {
	data, _ := c.JSONMarshal(discoveryJob)
	response, err := c.Post(DiscoveryJobAPI, nil, data)
	if err != nil {
		return models.DiscoveryJob{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeDiscoveryJob := models.DiscoveryJob{}
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	return omeDiscoveryJob, err
}

// UpdateDiscoveryJob - update a discovery job in OME.
func (c *Client) UpdateDiscoveryJob(discoveryJob models.DiscoveryJob) (models.DiscoveryJob, error) {
	data, _ := c.JSONMarshal(discoveryJob)
	queryParams := map[string]string{
		"groupId": strconv.Itoa(discoveryJob.DiscoveryConfigGroupID),
	}
	response, err := c.Do(http.MethodPost, DiscoveryJobAPI, nil, queryParams, data)
	if err != nil {
		return models.DiscoveryJob{}, err
	}
	respData, _ := c.GetBodyData(response.Body)
	omeDiscoveryJob := models.DiscoveryJob{}
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	return omeDiscoveryJob, err
}

// DeleteDiscoveryJob - delete a discovery job in OME.
func (c *Client) DeleteDiscoveryJob(discoveryGroupIds models.DiscoveryJobDeletePayload) (string, error) {
	data, _ := c.JSONMarshal(discoveryGroupIds)
	resp, err := c.Post(DiscoveryJobRemoveAPI, nil, data)
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

// GetDiscoveryJobByGroupID - get a discovery job from discovery group id.
func (c *Client) GetDiscoveryJobByGroupID(groupID int64) (models.DiscoveryJob, error) {
	omeDiscoveryJob := models.DiscoveryJob{}
	endpoint := fmt.Sprintf(DiscoveryJobByGroupIDAPI, groupID)
	response, err := c.Get(endpoint, nil, nil)
	if err != nil {
		return omeDiscoveryJob, err
	}
	respData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(respData, &omeDiscoveryJob)
	return omeDiscoveryJob, err
}
