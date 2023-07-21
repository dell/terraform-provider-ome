package ome

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// UserSchema - schema for terraform config of ome user
func UserSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},

		"user_type_id": schema.Int64Attribute{
			MarkdownDescription: "User Type ID",
			Description:         "User Type ID",
			Optional:            true,
			Computed:            true,
		},

		"directory_service_id": schema.Int64Attribute{
			MarkdownDescription: "Directory Service ID",
			Description:         "Directory Service ID",
			Optional:            true,
			Computed:            true,
		},

		"description": schema.StringAttribute{
			MarkdownDescription: "Description",
			Description:         "Description",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"password": schema.StringAttribute{
			MarkdownDescription: "Password",
			Description:         "Password",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"username": schema.StringAttribute{
			MarkdownDescription: "User Name",
			Description:         "User Name",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"role_id": schema.StringAttribute{
			MarkdownDescription: "Role ID",
			Description:         "Role ID",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"locked": schema.BoolAttribute{
			MarkdownDescription: "Locked",
			Description:         "Locked",
			Optional:            true,
			Computed:            true,
		},

		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Enabled",
			Description:         "Enabled",
			Optional:            true,
			Computed:            true,
		},
	}
}
