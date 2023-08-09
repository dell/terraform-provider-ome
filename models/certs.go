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

// CSRConfig - CSR generation form
type CSRConfig struct {
	DistinguishedName string `json:"DistinguishedName"`
	DepartmentName    string `json:"DepartmentName"`
	BusinessName      string `json:"BusinessName"`
	Locality          string `json:"Locality"`
	State             string `json:"State"`
	Country           string `json:"Country"`
	Email             string `json:"Email"`
	Sans              string `json:"San,omitempty"`
}

// CertInfo - Certificate Information received from OME
type CertInfo struct {
	IssuedTo  CSRConfig `json:"IssuedTo"`
	IssuedBy  CSRConfig `json:"IssuedBy"`
	ValidTo   string    `json:"ValidTo"`
	ValidFrom string    `json:"ValidFrom"`
}

// tfsdk structs

// CSRConfigModel - CSR generation tfsdk form
type CSRConfigModel struct {
	DistinguishedName types.String `tfsdk:"distinguished_name"`
	DepartmentName    types.String `tfsdk:"department_name"`
	BusinessName      types.String `tfsdk:"business_name"`
	Locality          types.String `tfsdk:"locality"`
	State             types.String `tfsdk:"state"`
	Country           types.String `tfsdk:"country"`
	Email             types.String `tfsdk:"email"`
	Sans              types.String `tfsdk:"subject_alternate_names"`
}

// CertInfoModel - Certificate Information tfsdk received from OME
type CertInfoModel struct {
	ID        types.String   `tfsdk:"id"`
	IssuedTo  CSRConfigModel `tfsdk:"issued_to"`
	IssuedBy  CSRConfigModel `tfsdk:"issued_by"`
	ValidTo   types.String   `tfsdk:"valid_to"`
	ValidFrom types.String   `tfsdk:"valid_from"`
}

// NewCSRConfigModel - Converts CSRConfig to CSRConfigModel
func NewCSRConfigModel(input CSRConfig) CSRConfigModel {
	return CSRConfigModel{
		DistinguishedName: types.StringValue(input.DistinguishedName),
		DepartmentName:    types.StringValue(input.DepartmentName),
		BusinessName:      types.StringValue(input.BusinessName),
		Locality:          types.StringValue(input.Locality),
		State:             types.StringValue(input.State),
		Country:           types.StringValue(input.Country),
		Email:             types.StringValue(input.Email),
		Sans:              types.StringValue(input.Sans),
	}
}

// NewCertInfoModel - Converts CertInfo to CertInfoModel
func NewCertInfoModel(info CertInfo) CertInfoModel {
	return CertInfoModel{
		ID:        types.StringValue("dummy"),
		IssuedTo:  NewCSRConfigModel(info.IssuedBy),
		IssuedBy:  NewCSRConfigModel(info.IssuedTo),
		ValidTo:   types.StringValue(info.ValidTo),
		ValidFrom: types.StringValue(info.ValidFrom),
	}
}
