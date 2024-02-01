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

import "github.com/hashicorp/terraform-plugin-framework/types"

// UserPayload - to store the create user payload
type UserPayload struct {
	UserTypeID         int    `json:"UserTypeId,omitempty"`
	DirectoryServiceID int    `json:"DirectoryServiceId,omitempty"`
	Description        string `json:"Description,omitempty"`
	Password           string `json:"Password,omitempty"`
	UserName           string `json:"UserName,omitempty"`
	RoleID             string `json:"RoleId,omitempty"`
	Locked             bool   `json:"Locked"`
	Enabled            bool   `json:"Enabled"`
}

// User - to store the ome user info in json tag struct
type User struct {
	ID                 string `json:"Id,omitempty"`
	UserTypeID         int    `json:"UserTypeId,omitempty"`
	DirectoryServiceID int    `json:"DirectoryServiceId,omitempty"`
	Description        string `json:"Description,omitempty"`
	Password           string `json:"Password,omitempty"`
	UserName           string `json:"UserName,omitempty"`
	RoleID             string `json:"RoleId,omitempty"`
	Locked             bool   `json:"Locked"`
	Enabled            bool   `json:"Enabled"`
}

// OmeUser - to store the ome user info in tfsdk tag struct
type OmeUser struct {
	ID                 types.String `tfsdk:"id"`
	UserTypeID         types.Int64  `tfsdk:"user_type_id"`
	DirectoryServiceID types.Int64  `tfsdk:"directory_service_id"`
	Description        types.String `tfsdk:"description"`
	Password           types.String `tfsdk:"password"`
	UserName           types.String `tfsdk:"username"`
	RoleID             types.String `tfsdk:"role_id"`
	Locked             types.Bool   `tfsdk:"locked"`
	Enabled            types.Bool   `tfsdk:"enabled"`
}
