package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type UserPayload struct {
	UserTypeID         int    `json:"UserTypeId,omitempty"`
	DirectoryServiceID int    `json:"DirectoryServiceId,omitempty"`
	Description        string `json:"Description,omitempty"`
	Name               string `json:"Name,omitempty"`
	Password           string `json:"Password,omitempty"`
	UserName           string `json:"UserName,omitempty"`
	RoleID             string `json:"RoleId,omitempty"`
	Locked             bool   `json:"Locked,omitempty"`
	Enabled            bool   `json:"Enabled,omitempty"`
}

type User struct {
	ID                 string `json:"Id,omitempty"`
	UserTypeID         int    `json:"UserTypeId,omitempty"`
	DirectoryServiceID int    `json:"DirectoryServiceId,omitempty"`
	Description        string `json:"Description,omitempty"`
	Name               string `json:"Name,omitempty"`
	Password           string `json:"Password,omitempty"`
	UserName           string `json:"UserName,omitempty"`
	RoleID             string `json:"RoleId,omitempty"`
	Locked             bool   `json:"Locked,omitempty"`
	Enabled            bool   `json:"Enabled,omitempty"`
}

type OmeUser struct {
	ID                 types.String `tfsdk:"id"`
	UserTypeID         types.Int64  `tfsdk:"user_type_id"`
	DirectoryServiceID types.Int64  `tfsdk:"directory_service_id"`
	Description        types.String `tfsdk:"description"`
	Name               types.String `tfsdk:"name"`
	Password           types.String `tfsdk:"password"`
	UserName           types.String `tfsdk:"username"`
	RoleID             types.String `tfsdk:"role_id"`
	Locked             types.Bool   `tfsdk:"locked"`
	Enabled            types.Bool   `tfsdk:"enabled"`
}


