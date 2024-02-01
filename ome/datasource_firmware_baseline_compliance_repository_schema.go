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

package ome

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func omeFirmwareBaselineComplianceRepositoryDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"fbc_repositories": schema.ListNestedAttribute{
			MarkdownDescription: "Repositories fetched.",
			Description:         "Repositories fetched.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeSingleRepositoryDataSchema()},
		},
		"id": schema.Int64Attribute{
			MarkdownDescription: "Dummy ID of the datasource.",
			Description:         "Dummy ID of the datasource.",
			Computed:            true,
		},
		"names": schema.SetAttribute{
			MarkdownDescription: "A list of repository names which can filter the datasource. Length should be at least 1.",
			Description:         "A list of repository names which can filter the datasource. Length should be at least 1.",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
				),
			},
		},
	}
}

func omeSingleRepositoryDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of the repository.",
			Description:         "ID of the repository.",
			Computed:            true,
		},
		"backup_existing_catalog": schema.BoolAttribute{
			MarkdownDescription: "Catalog will take backup automatically if it is true.",
			Description:         "Catalog will take backup automatically if it is true.",
			Computed:            true,
		},
		"check_certificate": schema.BoolAttribute{
			MarkdownDescription: "If certificate check must be done for HTTPS repository.",
			Description:         "If certificate check must be done for HTTPS repository.",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Description of the repository.",
			Description:         "Description of the repository.",
			Computed:            true,
		},
		"domain_name": schema.StringAttribute{
			MarkdownDescription: "Domain Name for user credentials.",
			Description:         "Domain Name for user credentials.",
			Computed:            true,
		},
		"editable": schema.BoolAttribute{
			MarkdownDescription: "True, if the catalog can be editable.",
			Description:         "True, if the catalog can be editable.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Unique name of the repository.",
			Description:         "Unique name of the repository.",
			Computed:            true,
		},
		"repository_type": schema.StringAttribute{
			MarkdownDescription: "Source of the repository.",
			Description:         "Source of the repository.",
			Computed:            true,
		},
		"source": schema.StringAttribute{
			MarkdownDescription: "URL or IP or FQDN of the repository host.",
			Description:         "URL or IP or FQDN of the repository host.",
			Computed:            true,
		},
		"username": schema.StringAttribute{
			MarkdownDescription: "Username to access the share containing the catalog (CIFS/HTTPS).",
			Description:         "Username to access the share containing the catalog (CIFS/HTTPS).",
			Computed:            true,
		},
	}
}
