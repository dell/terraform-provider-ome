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

// GetComplianceReportDetails - get compliance report details
func (c *Client) GetComplianceReportDetails(ids []int64) ([]models.FirmwareBaselinesGetModel, error) {
	response := []models.FirmwareBaselinesGetModel{}
	if len(ids) == 0 {
		return response, nil
	}
	var strSlice []string
	for _, num := range ids {
		str := strconv.FormatInt(num, 10)
		strSlice = append(strSlice, str)
	}

	getBody := map[string]any{"Ids": strSlice}
	body, err := c.JSONMarshal(getBody)
	if err != nil {
		return response, err
	}
	resp, err := c.Post(DeviceComplianceReportAPI, nil, body)
	if err != nil {
		return response, err
	}
	bodyData, _ := c.GetBodyData(resp.Body)
	err = c.JSONUnMarshal(bodyData, &response)
	if err != nil {
		err = fmt.Errorf(ErrBaselineDeviceReportsID+" %w", err)
	}
	return response, err
}
