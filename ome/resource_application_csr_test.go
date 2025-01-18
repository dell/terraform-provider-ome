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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCsr(t *testing.T) {
	if os.Getenv("TF_ACC") == "0" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testCreateCSRWoDN,
				ExpectError: regexp.MustCompile(".*\"distinguished_name\"[[:space:]]is[[:space:]]required.*"),
			},
			{
				Config:      testCreateCSRWoDeptName,
				ExpectError: regexp.MustCompile(".*\"department_name\"[[:space:]]is[[:space:]]required.*"),
			},
			{
				Config:      testCreateCSRWoBN,
				ExpectError: regexp.MustCompile(".*\"business_name\"[[:space:]]is[[:space:]]required.*"),
			},
			{
				Config:      testCreateCSRWoLocality,
				ExpectError: regexp.MustCompile(".*\"locality\"[[:space:]]is[[:space:]]required.*"),
			},
			{
				Config:      testCreateCSRWoState,
				ExpectError: regexp.MustCompile(".*\"state\"[[:space:]]is[[:space:]]required.*"),
			},
			{
				Config:      testCreateCSRWoCountry,
				ExpectError: regexp.MustCompile(".*\"country\"[[:space:]]is[[:space:]]required.*"),
			},
			{
				Config:      testCreateCSRWoEmail,
				ExpectError: regexp.MustCompile(".*\"email\"[[:space:]]is[[:space:]]required.*"),
			},
			{
				Config:      testCreateCSRWithMoreThan4San,
				ExpectError: regexp.MustCompile(".*at[[:space:]]most[[:space:]]4[[:space:]]elements.*"),
			},
			{
				Config:      testCreateCSRWithZeroSan,
				ExpectError: regexp.MustCompile(".*at[[:space:]]least[[:space:]]1[[:space:]]elements.*"),
			},
			{
				Config:      testCreateCSRWithEmptyDN,
				ExpectError: regexp.MustCompile(".*DistinguishedName[[:space:]]is[[:space:]]missing.*"),
			},
			{
				Config: testCreateCSR,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_application_csr.csr1", "csr")),
			},
			{
				Config: testUpdateCSRWithSan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ome_application_csr.csr1", "csr")),
			},
			{
				// this is actually an update test
				Config:      testCreateCSRWithEmptyDN,
				ExpectError: regexp.MustCompile(".*DistinguishedName[[:space:]]is[[:space:]]missing.*"),
			},
		},
	})

}

var testCreateCSR = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
	}
}
`

var testUpdateCSRWithSan = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
		subject_alternate_names = ["` + omeHost + `"]
	}
}
`

var testCreateCSRWithEmptyDN = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = ""
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
	}
}
`

var testCreateCSRWithMoreThan4San = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
		subject_alternate_names = ["aa","bb","cc","dd","ee"]
	}
}
`

var testCreateCSRWithZeroSan = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
		subject_alternate_names = []
	}
}
`

var testCreateCSRWoDN = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
	}
}
`
var testCreateCSRWoDeptName = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
	}
}
`

var testCreateCSRWoBN = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		locality = "RedRock"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
	}
}
`

var testCreateCSRWoLocality = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		state = "Texas"
		country = "US"
		email = "abc@gmail.com"
	}
}
`
var testCreateCSRWoState = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		country = "US"
		email = "abc@gmail.com"
	}
}
`
var testCreateCSRWoCountry = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		email = "abc@gmail.com"
	}
}
`
var testCreateCSRWoEmail = testProvider + `
resource "ome_application_csr" "csr1" {
	specs = {
		distinguished_name = "terraform.ome.com"
		department_name = "Terraform Server Solutions"
		business_name = "Dell Terraform"
		locality = "RedRock"
		state = "Texas"
		country = "US"
	}
}
`
