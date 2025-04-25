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
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	//SkipTestMsg
	SkipTestMsg                  = "Skipping the test because eith TF_ACC or ACC_DETAIL is not set to 1"
	SweepTestsTemplateIdentifier = "test_acc"
)

// Used for Mocking responses from functions
var FunctionMocker *Mocker

var globalEnvMap = getEnvMap()
var testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
var omeUserName = globalEnvMap["OME_USERNAME"]
var omeHost = globalEnvMap["OME_HOST"]
var omePassword = globalEnvMap["OME_PASSWORD"]
var port = setDefault(globalEnvMap["OME_PORT"], "443")
var protocol = setDefault(globalEnvMap["OME_PROTOCOL"], "https")
var DeviceSvcTag1 = globalEnvMap["DEVICESVCTAG1"]
var DeviceSvcTag2 = globalEnvMap["DEVICESVCTAG2"]
var DeviceID1 = globalEnvMap["DEVICEID1"]
var DeviceID2 = globalEnvMap["DEVICEID2"]     // Not capable for deployment
var DeviceIPExt = globalEnvMap["DEVICEIPEXT"] // Must be external to OME environment but discoverable
var ShareUser = globalEnvMap["SHAREUSERNAME"]
var SharePassword = globalEnvMap["SHAREPASSWORD"]
var ShareIP = globalEnvMap["SHAREIP"]
var DeviceIP1 = globalEnvMap["DEVICEIP1"]
var DeviceIP2 = globalEnvMap["DEVICEIP2"]
var Catalog1 = setDefault(globalEnvMap["CATALOG1"], "tfacc_catalog_dell_online_1")
var Repository = setDefault(globalEnvMap["REPOSITORY"], "tfacc_catalog_dell_online_1")
var CatalogResource = setDefault(globalEnvMap["CATALOG_RESOURCE"], "tfacc_firmware_catalog_resource")

// Device Model to be used in DS test
// Must have multiple devices of this model
var DeviceModel = globalEnvMap["DEVICE_MODEL"]

// an invalid cert - can be set to any text file
var InvCert = globalEnvMap["INV_CERT"]

// idrac username
var IdracUsername = globalEnvMap["IDRAC_USERNAME"]

// idrac password
var IdracPassword = globalEnvMap["IDRAC_PASSWORD"]

var testProvider = `
provider "ome" {
	username = "` + omeUserName + `"
	password = "` + omePassword + `"
	host = "` + omeHost + `"
	port = "` + port + `"
	protocol = "` + protocol + `"
	skipssl = true
}
`

func init() {
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	os.Setenv("TF_ACC", globalEnvMap["TF_ACC"])
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		// newProvider is an example function that returns a tfsdk.Provider
		"ome": providerserver.NewProtocol6WithError(New()),
	}

}

func getEnvMap() map[string]string {
	envMap, err := loadEnvFile("ome_test.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return envMap
	}
	return envMap
}

func loadEnvFile(path string) (map[string]string, error) {
	envMap := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}

func testAccPreCheck(t *testing.T) {
	if v := omeUserName; v == "" {
		t.Fatal("OME_USERNAME must be set for acceptance tests")
	}

	if v := omePassword; v == "" {
		t.Fatal("OME_PASSWORD must be set for acceptance tests")
	}

	if v := omeHost; v == "" {
		t.Fatal("OME_HOST must be set for acceptance tests")
	}

	// Make sure to unpatch before each new test is run
	if FunctionMocker != nil {
		FunctionMocker.UnPatch()
	}
	if localFunctionalMocker != nil {
		localFunctionalMocker.UnPatch()
	}
	if localMocker != nil {
		localMocker.UnPatch()
	}
	if localMocker2 != nil {
		localMocker2.UnPatch()
	}

	// testProvider.Configure(context.Background(), tfsdk.ConfigureProviderRequest{}, &tfsdk.ConfigureProviderResponse{})

}

func skipTest() bool {
	return globalEnvMap["TF_ACC"] == "" || globalEnvMap["ACC_DETAIL"] == ""
}

func getTestData(fileName string) string {
	wd, _ := os.Getwd()
	parent := filepath.Dir(wd)
	fileP := filepath.Join(parent, "testdata", fileName)
	return strings.ReplaceAll(fileP, "\\", "/")
}

// if there is no os setting set, then use the default value
func setDefault(osInput string, defaultStr string) string {
	if osInput == "" {
		return defaultStr
	}
	return osInput
}
