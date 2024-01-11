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

package ome

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func omeFirmwareCatalogDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "Dummy ID of the datasource.",
			Description:         "Dummy ID of the datasource.",
			Computed:            true,
		},
		"firmware_catalogs": schema.ListNestedAttribute{
			MarkdownDescription: "Devices fetched.",
			Description:         "Devices fetched.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeSingleCatalogFirmwareDataSchema()},
		},
		"names": schema.ListAttribute{
			MarkdownDescription: "A list of catalog names which can filter the datasource, if blank will return all catalogs.",
			Description:         "A list of catalog names which can filter the datasource, if blank will return all catalogs.",
			Optional:            true,
			ElementType:         types.StringType,
		},
	}
}

func omeSingleCatalogFirmwareDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"associated_baselines": schema.ListNestedAttribute{
			MarkdownDescription: "Associated Baselines.",
			Description:         "Associated Baselines.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"baseline_id": schema.Int64Attribute{
						MarkdownDescription: "Baseline ID.",
						Description:         "Baseline ID.",
						Computed:            true,
					},
					"baseline_name": schema.StringAttribute{
						MarkdownDescription: "Baseline Name.",
						Description:         "Baseline Name.",
						Computed:            true,
					},
				},
			},
		},
		"baseline_location": schema.StringAttribute{
			MarkdownDescription: "Baseline Location.",
			Description:         "Baseline Location.",
			Computed:            true,
		},
		"bundles_count": schema.Int64Attribute{
			MarkdownDescription: "Bundles Count.",
			Description:         "Bundles Count.",
			Computed:            true,
		},
		"source_path": schema.StringAttribute{
			MarkdownDescription: "Source path.",
			Description:         "Source path.",
			Computed:            true,
		},
		"create_date": schema.StringAttribute{
			MarkdownDescription: "Create Date.",
			Description:         "Create Date.",
			Computed:            true,
		},
		"filename": schema.StringAttribute{
			MarkdownDescription: "Filename.",
			Description:         "Filename.",
			Computed:            true,
		},
		"last_update": schema.StringAttribute{
			MarkdownDescription: "Last Update.",
			Description:         "Last Update.",
			Computed:            true,
		},
		"id": schema.Int64Attribute{
			MarkdownDescription: "id.",
			Description:         "id.",
			Computed:            true,
		},
		"manifest_identifier": schema.StringAttribute{
			MarkdownDescription: "Manifest Identifier.",
			Description:         "Manifest Identifier.",
			Computed:            true,
		},
		"manifest_version": schema.StringAttribute{
			MarkdownDescription: "Manifest Version.",
			Description:         "Manifest Version.",
			Computed:            true,
		},
		"next_update": schema.StringAttribute{
			MarkdownDescription: "Next Update.",
			Description:         "Next Update.",
			Computed:            true,
		},
		"owner_id": schema.Int64Attribute{
			MarkdownDescription: "Owner ID.",
			Description:         "Owner ID.",
			Computed:            true,
		},
		"predcessor_identifier": schema.StringAttribute{
			MarkdownDescription: "Predcessor Identifier.",
			Description:         "Predcessor Identifier.",
			Computed:            true,
		},
		"release_identifier": schema.StringAttribute{
			MarkdownDescription: "Release Identifier.",
			Description:         "Release Identifier.",
			Computed:            true,
		},
		"repository": schema.SingleNestedAttribute{
			MarkdownDescription: "Repository.",
			Description:         "Repository.",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"backup_existing_catalog": schema.BoolAttribute{
					MarkdownDescription: "Backup Existing Catalog.",
					Description:         "Backup Existing Catalog.",
					Computed:            true,
				},
				"check_certificate": schema.BoolAttribute{
					MarkdownDescription: "Check Certificate.",
					Description:         "Check Certificate.",
					Computed:            true,
				},
				"description": schema.StringAttribute{
					MarkdownDescription: "Description.",
					Description:         "Description.",
					Computed:            true,
				},
				"domain_name": schema.StringAttribute{
					MarkdownDescription: "Domain Name.",
					Description:         "Domain Name.",
					Computed:            true,
				},
				"editable": schema.BoolAttribute{
					MarkdownDescription: "Editable.",
					Description:         "Editable.",
					Computed:            true,
				},
				"id": schema.Int64Attribute{
					MarkdownDescription: "ID.",
					Description:         "ID.",
					Computed:            true,
				},
				"name": schema.StringAttribute{
					MarkdownDescription: "Name.",
					Description:         "Name.",
					Computed:            true,
				},
				"repository_type": schema.StringAttribute{
					MarkdownDescription: "Repository Type.",
					Description:         "Repository Type.",
					Computed:            true,
				},
				"source": schema.StringAttribute{
					MarkdownDescription: "Source.",
					Description:         "Source.",
					Computed:            true,
				},
				"username": schema.StringAttribute{
					MarkdownDescription: "Username.",
					Description:         "Username.",
					Computed:            true,
				},
			},
		},
		"schedule": schema.SingleNestedAttribute{
			MarkdownDescription: "Schedule.",
			Description:         "Schedule.",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"cron": schema.StringAttribute{
					MarkdownDescription: "Cron.",
					Description:         "Cron.",
					Computed:            true,
				},
				"end_time": schema.StringAttribute{
					MarkdownDescription: "End Time.",
					Description:         "End Time.",
					Computed:            true,
				},
				"start_time": schema.StringAttribute{
					MarkdownDescription: "Start Time.",
					Description:         "Start Time.",
					Computed:            true,
				},
				"run_later": schema.BoolAttribute{
					MarkdownDescription: "Run later.",
					Description:         "Run later.",
					Computed:            true,
				},
				"run_now": schema.BoolAttribute{
					MarkdownDescription: "Run Now.",
					Description:         "Run Now.",
					Computed:            true,
				},
			},
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "Status.",
			Description:         "Status.",
			Computed:            true,
		},
		"task_id": schema.Int64Attribute{
			MarkdownDescription: "Task ID.",
			Description:         "Task ID.",
			Computed:            true,
		},
	}
}
