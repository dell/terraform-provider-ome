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

// BaseLineModel - BaseLineModel responses from api.
type BaseLineModel struct {
	OdataContext string                  `json:"@odata.context"`
	OdataCount   int64                   `json:"@odata.count"`
	Value        []UpdateServiceBaseline `json:"value"`
}

// ComplianceSummaryData - ComplianceSummaryData responses from api.
type ComplianceSummaryData struct {
	ComplianceStatus  string `json:"ComplianceStatus"`
	NumberOfCritical  int64  `json:"NumberOfCritical"`
	NumberOfDowngrade int64  `json:"NumberOfDowngrade"`
	NumberOfNormal    int64  `json:"NumberOfNormal"`
	NumberOfUnknown   int64  `json:"NumberOfUnknown"`
	NumberOfWarning   int64  `json:"NumberOfWarning"`
}

// TypeData - TypeData structure
type TypeData struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name"`
}

// TargetData - TargetData structure
type TargetData struct {
	ID   int64    `json:"Id"`
	Type TypeData `json:"Type"`
}

// UpdateServiceBaseline - UpdateService's Baseline structure
type UpdateServiceBaseline struct {
	OdataID                                    string                `json:"@odata.id"`
	OdataType                                  string                `json:"@odata.type"`
	CatalogID                                  int64                 `json:"CatalogId"`
	ComplianceSummary                          ComplianceSummaryData `json:"ComplianceSummary"`
	Description                                string                `json:"Description"`
	DeviceComplianceReportsOdataNavigationLink string                `json:"DeviceComplianceReports@odata.navigationLink"`
	DowngradeEnabled                           bool                  `json:"DowngradeEnabled"`
	FilterNoRebootRequired                     bool                  `json:"FilterNoRebootRequired"`
	ID                                         int64                 `json:"Id"`
	Is64Bit                                    bool                  `json:"Is64Bit"`
	LastRun                                    string                `json:"LastRun"`
	Name                                       string                `json:"Name"`
	RepositoryID                               int64                 `json:"RepositoryId"`
	RepositoryName                             string                `json:"RepositoryName"`
	RepositoryType                             string                `json:"RepositoryType"`
	Targets                                    []TargetData          `json:"Targets"`
	TaskID                                     int64                 `json:"TaskId"`
	TaskStatusID                               int64                 `json:"TaskStatusId"`
}
