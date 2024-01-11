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
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// FirmwareCatalogSchema returns the schema for the FirmwareCatalog resource
func FirmwareCatalogSchema() map[string]schema.Attribute {
	objd, _ := basetypes.NewObjectValue(
		map[string]attr.Type{
			"cadence":         types.StringType,
			"day_of_the_week": types.StringType,
			"time_of_day":     types.Int64Type,
			"am_pm":           types.StringType,
		},
		map[string]attr.Value{
			"cadence":         types.StringValue("Daily"),
			"day_of_the_week": types.StringValue("Monday"),
			"time_of_day":     types.Int64Value(1),
			"am_pm":           types.StringValue("AM"),
		},
	)

	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "Name of the new catalog.",
			Description:         "Name of the new catalog.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
					"must contain only alphanumeric characters and _-",
				),
			},
		},
		"catalog_update_type": schema.StringAttribute{
			MarkdownDescription: "Catalog Update Type. Sets the frequency of catalog updates. Defaults to Manual. If set to automatic, the catalog_refresh_schedule field will need to be set. Options are (Manual, Automatic).",
			Description:         "Catalog Update Type. Sets the frequency of catalog updates. Defaults to Manual. If set to automatic, the catalog_refresh_schedule field will need to be set. Options are (Manual, Automatic).",
			Default:             stringdefault.StaticString("Manual"),
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("Manual", "Automatic"),
			},
		},
		"share_type": schema.StringAttribute{
			MarkdownDescription: "Share Type, the type of share the catalog will pull from, Defaults to Dell. The different options will have different required fields to work properly. Options are (DELL, NFS, CIFS, HTTP, HTTPS).",
			Description:         "Share Type, the type of share the catalog will pull from, Defaults to Dell. The different options will have different required fields to work properly. Options are (DELL, NFS, CIFS, HTTP, HTTPS).",
			Default:             stringdefault.StaticString("DELL_ONLINE"),
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("DELL_ONLINE", "NFS", "CIFS", "HTTP", "HTTPS"),
			},
		},
		"catalog_file_path": schema.StringAttribute{
			MarkdownDescription: "Catalog File Path. Path on the share to gather catalog data. This field is required for share_types (NFS, CIFS, HTTP, HTTPS)",
			Description:         "Catalog File Path. Path on the share to gather catalog data. This field is required for share_types (NFS, CIFS, HTTP, HTTPS)",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[a-zA-Z0-9\/_-]*.*\/.*$`),
					"Start directory path without leading '/' and contain only alphanumeric characters and /_/-",
				),
			},
		},
		"catalog_refresh_schedule": schema.SingleNestedAttribute{
			MarkdownDescription: "Catalog Refresh Schedule, when using automatic catalog update the schedule is required for cadence of the update. If catalog_update_type is set to manual, this field is ignored.",
			Description:         "Catalog Refresh Schedule, when using automatic catalog update the schedule is required for cadence of the update. If catalog_update_type is set to manual, this field is ignored.",
			Optional:            true,
			Default:             objectdefault.StaticValue(objd),
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"cadence": schema.StringAttribute{
					MarkdownDescription: "Cadence. Options are(Weekly, Daily).",
					Description:         "Cadence. Options are(Weekly, Daily).",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("Weekly", "Daily"),
					},
					Computed: true,
				},
				"day_of_the_week": schema.StringAttribute{
					MarkdownDescription: "Day of the Week, only useful for weekly schedules. Options are(Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday).",
					Description:         "Day of the Week, only useful for weekly schedules. Options are(Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday).",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"),
					},
					Computed: true,
				},
				"time_of_day": schema.Int64Attribute{
					MarkdownDescription: "Time of Day for the schedule in hour increments. Options are (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12).",
					Description:         "Time of Day for the schedule in hour increments. Options are (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12).",
					Optional:            true,
					Validators: []validator.Int64{
						int64validator.OneOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12),
					},
					Computed: true,
				},
				"am_pm": schema.StringAttribute{
					MarkdownDescription: "AM/PM for the schedule. Options are (AM, PM).",
					Description:         "AM/PM for the schedule. Options are (AM, PM).",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("AM", "PM"),
					},
					Computed: true,
				},
			},
		},
		"share_address": schema.StringAttribute{
			MarkdownDescription: "Share Address. Gives the Ipv4, Ipv6, or FQDN of the share. This field is required for share_types (NFS, CIFS, HTTP, HTTPS)",
			Description:         "Share Address. Gives the Ipv4, Ipv6, or FQDN of the share. This field is required for share_types (NFS, CIFS, HTTP, HTTPS)",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[a-zA-Z0-9_\-\:\/\.]*$`),
					"must be ivp4, ipv6, or fqdn and start with desired protocol (e.g. https://)",
				),
			},
		},
		"domain": schema.StringAttribute{
			MarkdownDescription: "Domain. The domain for the catalog. This field is optional and only used for share_types (CIFS).",
			Description:         "Domain. The domain for the catalog. This field is optional and only used for share_types (CIFS).",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
					"must contain only alphanumeric characters and _-",
				),
			},
		},
		"catalog_user": schema.StringAttribute{
			MarkdownDescription: "Catalog User. The username related to the catalog share. This field is required for share_types (CIFS, HTTPS).",
			Description:         "Catalog User. The username related to the catalog share. This field is required for share_types (CIFS, HTTPS)",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
					"must contain only alphanumeric characters and _-",
				),
			},
		},
		"catalog_password": schema.StringAttribute{
			MarkdownDescription: "Catalog Password. The password related to the catalog share. This field is required for share_types (CIFS, HTTPS)",
			Description:         "Catalog Password. The password related to the catalog share. This field is required for share_types (CIFS, HTTPS)",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[a-zA-Z0-9_-]*$`),
					"must contain only alphanumeric characters and _-",
				),
			},
		},

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
		"repository": schema.ObjectAttribute{
			MarkdownDescription: "Repository.",
			Description:         "Repository.",
			Computed:            true,
			AttributeTypes: map[string]attr.Type{
				"backup_existing_catalog": types.BoolType,
				"check_certificate":       types.BoolType,
				"description":             types.StringType,
				"domain_name":             types.StringType,
				"editable":                types.BoolType,
				"id":                      types.Int64Type,
				"name":                    types.StringType,
				"repository_type":         types.StringType,
				"source":                  types.StringType,
				"username":                types.StringType,
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
