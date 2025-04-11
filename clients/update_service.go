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
	"context"
	"fmt"
	"net/http"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetFwBaselineComplianceReport retrieves the compliance report for a firmware baseline.
//
// ctx: The context.Context object for the request.
// baseLineID: The ID of the firewall baseline.
// Returns a pointer to a models.ComplianceReport object and an error.
func (c *Client) GetFwBaselineComplianceReport(ctx context.Context, baseLineID int64, filterkey string, filterval string) (*models.ComplianceReport, error) {
	compReports := models.ComplianceReport{}
	var resp *http.Response
	var err error
	if filterkey == "" || filterval == "" {
		resp, err = c.Get(fmt.Sprintf(FwBaselineComplianceReportsAPI, baseLineID), nil, nil)
	} else {
		tflog.Info(ctx, fmt.Sprintf("Filtering on %s: %s", filterkey, filterval))
		resp, err = c.Get(fmt.Sprintf(FwBaselineComplianceReportsAPI, baseLineID), nil, map[string]string{"$filter": fmt.Sprintf("%s eq '%s'", filterkey, filterval)})
	}
	if err != nil {
		return &models.ComplianceReport{}, err
	}
	respData, getBodyError := c.GetBodyData(resp.Body)
	if getBodyError != nil {
		return &models.ComplianceReport{}, getBodyError
	}
	err = c.JSONUnMarshal(respData, &compReports)

	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("Comp Reports %v", compReports))
		return &models.ComplianceReport{}, err
	}
	if len(compReports.Value) != 0 {
		return &compReports, nil
	}
	return &models.ComplianceReport{}, fmt.Errorf(ErrFwBaselineReport, baseLineID)
}

// GetUpdateServiceBaselineIDByName retrieves the Update Service baseline ID
// by its name.
//
// name: the name of the baseline to retrieve.
// int64: the ID of the baseline if found.
// error: an error if the baseline is not found or if there is an error during
// the retrieval process.
func (c *Client) GetUpdateServiceBaselineIDByName(name string) (int64, error) {
	baseLines := models.BaseLineModel{}
	resp, err := c.Get(FirmwareBaselineAPI, nil, nil)
	if err != nil {
		return -1, err
	}
	respData, getBodyError := c.GetBodyData(resp.Body)
	if getBodyError != nil {
		return -1, getBodyError
	}
	err = c.JSONUnMarshal(respData, &baseLines)

	if err != nil {
		return -1, err
	}
	for _, usBaseline := range baseLines.Value {
		if usBaseline.Name == name {
			return usBaseline.ID, nil
		}
	}
	return -1, fmt.Errorf(ErrBaselineNameNotFound, name)
}
