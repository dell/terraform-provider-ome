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
	"fmt"
	"terraform-provider-ome/models"
)

// GetComplianceReportDetails - get compliance report details
func (c *Client) GetComplianceReportDetails(baselineID int64) (models.DeviceComplianceReportModel, error) {
	response := models.DeviceComplianceReportModel{}
	resp, err := c.Get(fmt.Sprintf(DeviceComplianceReportAPI, baselineID), nil, nil)
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
