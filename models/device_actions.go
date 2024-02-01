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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceActionModel - Tfsdk model for device action resource
type DeviceActionModel struct {
	ID             types.Int64  `tfsdk:"id"`
	DeviceIDs      []int64      `tfsdk:"device_ids"`
	Action         types.String `tfsdk:"action"`
	Cron           types.String `tfsdk:"cron"`
	Timeout        types.Int64  `tfsdk:"timeout"`
	JobName        types.String `tfsdk:"job_name"`
	JobDescription types.String `tfsdk:"job_description"`
	NextRunTime    types.String `tfsdk:"next_run_time"`
	LastRunTime    types.String `tfsdk:"last_run_time"`
	JobStatus      types.String `tfsdk:"current_status"`
	LastRunStatus  types.String `tfsdk:"last_run_status"`
	StartTime      types.String `tfsdk:"start_time"`
	EndTime        types.String `tfsdk:"end_time"`
}
