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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

const (
	//SkipTestMsg
	SkipTestMsg                  = "Skipping the test because eith TF_ACC or ACC_DETAIL is not set to 1"
	SweepTestsTemplateIdentifier = "test_acc"
)

var testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
var omeUserName = os.Getenv("OME_USERNAME")
var omeHost = os.Getenv("OME_HOST")
var omePassword = os.Getenv("OME_PASSWORD")
var DeviceSvcTag1 = os.Getenv("DEVICESVCTAG1")
var DeviceSvcTag2 = os.Getenv("DEVICESVCTAG2")
var GroupID1 = os.Getenv("GROUPID1")
var DeviceID1 = os.Getenv("DEVICEID1")
var DeviceID2 = os.Getenv("DEVICEID2")
var DeviceID3 = os.Getenv("DEVICEID3")             // Not capable for deployment
var DeviceSvcTagRmv = os.Getenv("DEVICESVCTAGRMV") // Will be removed
var ShareUser = os.Getenv("SHAREUSERNAME")
var SharePassword = os.Getenv("SHAREPASSWORD")
var ShareIP = os.Getenv("SHAREIP")
var DeviceIP1 = os.Getenv("DEVICEIP1")
var DeviceIP2 = os.Getenv("DEVICEIP2")
var DeviceIP3 = os.Getenv("DEVICEIP3")

// the following must be absolute paths
// script that accepts csr and output file as inputs respectively and generates cert
var CertScript = os.Getenv("CERT_SCRIPT")

// an invalid cert
var InvCert = os.Getenv("INV_CERT")

// idrac username 
var IdracUsername = os.Getenv("IDRAC_USERNAME")
// idrac password 
var IdracPassword = os.Getenv("IDRAC_PASSWORD")

var testProvider = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	skipssl = true
}
`

func init() {
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	err := godotenv.Load("ome_test.env")
	if err != nil {
		panic(err)
	}
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		// newProvider is an example function that returns a tfsdk.Provider
		"ome": providerserver.NewProtocol6WithError(New()),
	}

}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OME_USERNAME"); v == "" {
		t.Fatal("OME_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("OME_PASSWORD"); v == "" {
		t.Fatal("OME_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("OME_HOST"); v == "" {
		t.Fatal("OME_HOST must be set for acceptance tests")
	}

	// testProvider.Configure(context.Background(), tfsdk.ConfigureProviderRequest{}, &tfsdk.ConfigureProviderResponse{})

}

func skipTest() bool {
	return os.Getenv("TF_ACC") == "" || os.Getenv("ACC_DETAIL") == ""
}

func getTestData(fileName string) string {
	wd, _ := os.Getwd()
	parent := filepath.Dir(wd)
	fileP := filepath.Join(parent, "testdata", fileName)
	return strings.ReplaceAll(fileP, "\\", "/")
}
