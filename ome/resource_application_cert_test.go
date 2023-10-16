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
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCert(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	var testCreateCertInvalid = testProvider + `
	resource "ome_application_certificate" "cert1" {
		certificate_base64 = "invalid"
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testCreateCertInvalid,
				ExpectError: regexp.MustCompile(".*failed to decode base64.*"),
			},
			{
				SkipFunc: func() (bool, error) {
					t.Log("Skipping as INV_CERT is not set")
					return (InvCert == ""), nil
				},
				Config:      testCertBad,
				ExpectError: regexp.MustCompile(".*certificate[[:space:]]file[[:space:]]provided[[:space:]]is[[:space:]]invalid.*"),
			},
		},
	})

}

func TestCertPos(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	var testCreateCert = testProvider + `
	provider "random" {
	}
	resource "random_string" "csr_seed" {
	  length = 16
	  special = false
	  keepers = {
		"key" = "%d"
	  }
	}
	resource "ome_application_csr" "csr1" {
		specs = {
			distinguished_name = "terraform${random_string.csr_seed.result}.ome.com"
			department_name = "Terraform Server Solutions ${random_string.csr_seed.result}"
			business_name = "Dell Terraform"
			locality = "RedRock"
			state = "Texas"
			country = "US"
			email = "abc@gmail.com"
		}
	}
	
	locals {
	  formatted_private_key = replace(
		  replace(ome_application_csr.csr1.csr, "-----BEGIN CERTIFICATE REQUEST-----", "-----BEGIN CERTIFICATE REQUEST-----\n"),
		  "-----END CERTIFICATE REQUEST-----", "\n-----END CERTIFICATE REQUEST-----\n"
		)
	  cert_file = "${path.cwd}/certificateTf.pem.crt"
	  cert_file_cmd = "%s"
	}
	
	resource "local_file" "csr_file" {
	  filename = "${path.cwd}/reqTf.pem"
	  content = local.formatted_private_key
	}
	
	resource "terraform_data" "certificate" {
	  depends_on = [ local_file.csr_file ]
	  triggers_replace = [
		local_file.csr_file.content
	  ]
	
	  provisioner "local-exec" {
		command = "${local.cert_file_cmd} ${local_file.csr_file.filename} ${local.cert_file}"
	  }
	}
	
	data "local_file" "cert" {
	  depends_on = [ terraform_data.certificate ]
	  filename = local.cert_file
	}

	resource "ome_application_certificate" "ome_cert" {
		certificate_base64 = data.local_file.cert.content_base64
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if InvCert == "" || CertScript == "" {
				t.Skip("Skipping because INV_CERT=", InvCert, " and CERT_SCRIPT=", CertScript)
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
			"local": {
				Source: "hashicorp/local",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testCreateCert, 0, CertScript),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 20)
				},
				Config: fmt.Sprintf(testCreateCert, 1, CertScript),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 20)
				},
				Config:      testCertBad,
				ExpectError: regexp.MustCompile(".*certificate[[:space:]]file[[:space:]]provided[[:space:]]is[[:space:]]invalid.*"),
			},
		},
	})

}

var testCertBad = testProvider + `
resource "ome_application_certificate" "ome_cert" {
	certificate_base64 = filebase64("` + InvCert + `")
}
`
