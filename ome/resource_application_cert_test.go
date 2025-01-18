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
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCert(t *testing.T) {
	if os.Getenv("TF_ACC") == "0" {
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
					if InvCert == "" {
						t.Log("Skipping as INV_CERT is not set")
						return true, nil
					}
					return false, nil
				},
				Config:      testCertBad,
				ExpectError: regexp.MustCompile(".*certificate[[:space:]]file[[:space:]]provided[[:space:]]is[[:space:]]invalid.*"),
			},
		},
	})

}

func TestCertPos(t *testing.T) {
	if os.Getenv("TF_ACC") == "0" {
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

	locals {
		ca_key = "${path.cwd}/ca_private_key.pem"
		ca_cert = "${path.cwd}/ca_cert.pem"
		ca_subj = "/C=US/ST=California/L=The Cloud/O=Dell/OU=ISG/CN=test_acc_user"
	}
	
	resource "terraform_data" "ca" {
		provisioner "local-exec" {
			command = "openssl req -x509 -days 365 -newkey rsa:4096 -keyout ${local.ca_key} -out ${local.ca_cert} -nodes -subj \"${local.ca_subj}\""
		}
		provisioner "local-exec" {
			when = destroy
			command = "rm ca_* "
		}
	}
	resource "ome_application_csr" "csr1" {
		depends_on = [ terraform_data.ca ]
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
		cert_file = "${path.cwd}/my_signed_cert.pem"
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
			command = "openssl x509 -req -in ${local_file.csr_file.filename} -days 365 -CA ${local.ca_cert} -CAkey ${local.ca_key} -CAcreateserial -out my_signed_cert.pem"
		}
		provisioner "local-exec" {
			when = destroy
			command = "rm my_signed_cert.pem"
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
			if InvCert == "" {
				t.Skip("Skipping because INV_CERT=", InvCert)
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
				Config: fmt.Sprintf(testCreateCert, 0),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 20)
				},
				Config: fmt.Sprintf(testCreateCert, 1),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 20)
				},
				SkipFunc: func() (bool, error) {
					if InvCert == "" {
						t.Log("Skipping as INV_CERT is not set")
						return true, nil
					}
					return false, nil
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
