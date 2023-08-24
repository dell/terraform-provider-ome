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
	"encoding/json"
	"fmt"
	"terraform-provider-ome/models"
)

func (c *Client) CreateJob(payload models.JobPayload) (JobResp, error) {
	payloadb, errb := json.Marshal(payload)
	if errb != nil {
		return JobResp{}, fmt.Errorf("error marshalling payload: %w", errb)
	}
	response, err := c.Post(JobAPI, nil, payloadb)
	if err != nil {
		return JobResp{}, err
	}
	bodyData, _ := c.GetBodyData(response.Body)
	temp := JobResp{}
	if err := json.Unmarshal(bodyData, &temp); err != nil {
		return temp, fmt.Errorf("error unmarshalling job creation response: %w", err)
	}
	return temp, nil
}

func (c *Client) DeleteJob(id int64) error {
	path := fmt.Sprintf(GetJobAPI, id)
	_, err := c.Delete(path, nil, nil)
	return err
}
