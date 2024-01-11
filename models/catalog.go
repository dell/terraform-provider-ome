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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// Catalogs - catalog responses from catalog api.
type Catalogs struct {
	Context string          `json:"@odata.context"`
	Count   int             `json:"@odata.count"`
	Value   []CatalogsModel `json:"value"`
}

// AssociatedBaselineModel - model for associated baselines response
type AssociatedBaselineModel struct {
	// ID of baselines associated with the catalog.
	BaselineID int64 `json:"BaselineId,omitempty"`
	// Name of a baseline associated with the catalog.
	BaselineName string `json:"BaselineName,omitempty"`
}

// CatalogsModel - model for catalog response
type CatalogsModel struct {
	// The baselines that are associated with the catalog. This is an array.
	AssociatedBaselines []AssociatedBaselineModel `json:"AssociatedBaselines,omitempty"`
	// The repository location for the catalog (for example for online catalogs, the location is downloads.dell.com)
	BaseLocation string `json:"BaseLocation"`
	// Number of bundles referenced in the catalog
	BundlesCount int64 `json:"BundlesCount,omitempty"`
	// Catalog release date.
	CreatedDate string `json:"CreatedDate,omitempty"`
	// File name of the catalog
	Filename string `json:"Filename,omitempty"`
	ID       int64  `json:"Id,omitempty"`
	// Timestamp of catalog instance being updated in the appliance.
	LastUpdated string `json:"LastUpdated,omitempty"`
	// Catalog manifest ID.
	ManifestIdentifier string `json:"ManifestIdentifier,omitempty"`
	// Catalog manifest version.
	ManifestVersion string `json:"ManifestVersion,omitempty"`
	// Next update time
	NextUpdate string `json:"NextUpdate,omitempty"`
	// ID of the owner of the catalog creator
	OwnerID int64 `json:"OwnerId,omitempty"`
	// ID for the catalog version that was published before the current catalog.
	PredecessorIdentifier string `json:"PredecessorIdentifier,omitempty"`
	// Catalog release ID.
	ReleaseIdentifier string          `json:"ReleaseIdentifier,omitempty"`
	Repository        RepositoryModel `json:"Repository"`
	// Schedule information
	Schedule ScheduleModel `json:"Schedule,omitempty"`
	// Repository source-Full path including     subfolders to where the catalog file is located. For example-downloads.dell.com/catalog/catalog.gz
	SourcePath string `json:"SourcePath,omitempty"`
	// Status of creating the catalog.
	Status string `json:"Status,omitempty"`
	// The ID of the task or job that is created to download the catalog.
	TaskID int64 `json:"TaskId,omitempty"`
}

// ScheduleModel - model for schedule response
type ScheduleModel struct {
	Cron      string `json:"Cron"`
	EndTime   string `json:"EndTime"`
	RunLater  bool   `json:"RunLater"`
	RunNow    bool   `json:"RunNow"`
	StartTime string `json:"StartTime"`
}

// RepositoryModel - model for repository response
type RepositoryModel struct {
	// Catalog will take backup automatically if it is true
	BackupExistingCatalog bool `json:"BackupExistingCatalog,omitempty"`
	// If certificate check must be done for HTTPS repository
	CheckCertificate bool `json:"CheckCertificate,omitempty"`
	// Description of the repository
	Description string `json:"Description,omitempty"`
	// Domain Name for user credentials
	DomainName string `json:"DomainName,omitempty"`
	// True, if the catalog can be editable
	Editable bool  `json:"Editable,omitempty"`
	ID       int64 `json:"Id,omitempty"`
	// Name of the repository. The name must be unique for each repository that is created
	Name string `json:"Name"`
	// Password to access the share containing the catalog (CIFS/HTTPS)
	Password string `json:"Password,omitempty"`
	// Source of the repository. For potential values, check out the Enumerator Guide
	RepositoryType string `json:"RepositoryType"`
	// URL or IP or FQDN of the repository host
	Source string `json:"Source,omitempty"`
	// Username to access the share containing the catalog (CIFS/HTTPS)
	Username string `json:"Username,omitempty"`
}

// OMECatalogData represents the OME Firmware Catalog
type OMECatalogData struct {
	ID      types.Int64           `tfsdk:"id"`
	Catalog []OmeSigleCatalogData `tfsdk:"firmware_catalogs"`
	Names   types.List            `tfsdk:"names"`
}

// OmeSigleCatalogData represents the OME Firmware Catalog
type OmeSigleCatalogData struct {
	AssociatedBaselines   []AssociatedBaselines `tfsdk:"associated_baselines"`
	BaseLocation          types.String          `tfsdk:"baseline_location"`
	BundlesCount          types.Int64           `tfsdk:"bundles_count"`
	CreatedDate           types.String          `tfsdk:"create_date"`
	Filename              types.String          `tfsdk:"filename"`
	ID                    types.Int64           `tfsdk:"id"`
	LastUpdated           types.String          `tfsdk:"last_update"`
	ManifestIdentifier    types.String          `tfsdk:"manifest_identifier"`
	ManifestVersion       types.String          `tfsdk:"manifest_version"`
	NextUpdate            types.String          `tfsdk:"next_update"`
	OwnerID               types.Int64           `tfsdk:"owner_id"`
	PredecessorIdentifier types.String          `tfsdk:"predcessor_identifier"`
	ReleaseIdentifier     types.String          `tfsdk:"release_identifier"`
	Repository            CatalogRepository     `tfsdk:"repository"`
	Schedule              ScheduleCatalog       `tfsdk:"schedule"`
	SourcePath            types.String          `tfsdk:"source_path"`
	Status                types.String          `tfsdk:"status"`
	TaskID                types.Int64           `tfsdk:"task_id"`
}

// ScheduleCatalog - model for schedule tf
type ScheduleCatalog struct {
	Cron      types.String `tfsdk:"cron"`
	EndTime   types.String `tfsdk:"end_time"`
	RunLater  types.Bool   `tfsdk:"run_later"`
	RunNow    types.Bool   `tfsdk:"run_now"`
	StartTime types.String `tfsdk:"start_time"`
}

// CatalogRepository - model for repository tf
type CatalogRepository struct {
	BackupExistingCatalog types.Bool   `tfsdk:"backup_existing_catalog"`
	CheckCertificate      types.Bool   `tfsdk:"check_certificate"`
	Description           types.String `tfsdk:"description"`
	DomainName            types.String `tfsdk:"domain_name"`
	Editable              types.Bool   `tfsdk:"editable"`
	ID                    types.Int64  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	RepositoryType        types.String `tfsdk:"repository_type"`
	Source                types.String `tfsdk:"source"`
	Username              types.String `tfsdk:"username"`
}

// AssociatedBaselines - model for associated baselines tf
type AssociatedBaselines struct {
	BaselineID   types.Int64  `tfsdk:"baseline_id"`
	BaselineName types.String `tfsdk:"baseline_name"`
}
