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
			MarkdownDescription: "ID of the OME user.",
			Description:         "ID of the OME user.",
			Computed:            true,
		},

		"user_type_id": schema.Int64Attribute{
			MarkdownDescription: "User Type ID of the OME user.",
			Description:         "User Type ID of the OME user.",
			Optional:            true,
			Computed:            true,
		},

		"directory_service_id": schema.Int64Attribute{
			MarkdownDescription: "Directory Service ID of the OME user.",
			Description:         "Directory Service ID of the OME user.",
			Optional:            true,
			Computed:            true,
		},

		"description": schema.StringAttribute{
			MarkdownDescription: "Description of the OME user.",
			Description:         "Description of the OME user.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"password": schema.StringAttribute{
			MarkdownDescription: "Password of the OME user.",
			Description:         "Password of the OME user.",
			Required:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"username": schema.StringAttribute{
			MarkdownDescription: "Username of the OME user.",
			Description:         "Username of the OME user.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"role_id": schema.StringAttribute{
			MarkdownDescription: "Role ID of the OME user.",
			Description:         "Role ID of the OME user.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"locked": schema.BoolAttribute{
			MarkdownDescription: "Lock OME user.",
			Description:         "Lock OME user.",
			Optional:            true,
			Computed:            true,
		},

		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable OME user.",
			Description:         "Enable OME user.",
			Optional:            true,
			Computed:            true,
		},
	}
}
