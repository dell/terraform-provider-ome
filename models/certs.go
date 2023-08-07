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

type CertInfo struct {
	IssuedTo  CSRConfig `json:"IssuedTo"`
	IssuedBy  CSRConfig `json:"IssuedBy"`
	ValidTo   string    `json:"ValidTo"`
	ValidFrom string    `json:"ValidFrom"`
}
