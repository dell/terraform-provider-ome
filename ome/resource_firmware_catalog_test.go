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
	"terraform-provider-ome/helper"
	"terraform-provider-ome/models"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var localFunctionalMocker *Mocker

func TestFirmwareCatalogResourceCreate(t *testing.T) {
	if os.Getenv("TF_ACC") == "0" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	var catalogTfName = "ome_firmware_catalog.cat_1"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFirmwareCatalogResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(catalogTfName, "repository.backup_existing_catalog", "false"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.check_certificate", "false"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.description", "tfacc_firmware_catalog_resource terraform catalog"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.domain_name", "example"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.editable", "true"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.name", CatalogResource),
					resource.TestCheckResourceAttr(catalogTfName, "repository.repository_type", "HTTPS"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.source", "https://1.2.2.1"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.username", "example-user"),
				),
			},
			{
				Config: updateFirmwareCatalogResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(catalogTfName, "repository.backup_existing_catalog", "false"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.check_certificate", "false"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.description", "tfacc_firmware_catalog_resource_update terraform catalog"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.domain_name", "example_update"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.editable", "true"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.name", CatalogResource+"_update"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.repository_type", "HTTPS"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.source", "https://2.1.1.2"),
					resource.TestCheckResourceAttr(catalogTfName, "repository.username", "example-user_update"),
				),
			},
			// Import testing
			{
				ResourceName: catalogTfName,
				ImportState:  true,
			},
		},
	})
}

func TestFirmwareCatalogResourceValidationCreateError(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      createFirmwareCatalogResourceValidateErrorAutomatic,
				ExpectError: regexp.MustCompile(".*invalid automatic update configuration.*"),
			},
			{
				Config:      createFirmwareCatalogResourceValidateErrorNFS,
				ExpectError: regexp.MustCompile(".*invalid NFS share configuration.*"),
			},
			{
				Config:      createFirmwareCatalogResourceValidateErrorCIFS,
				ExpectError: regexp.MustCompile(".*invalid CIFS share configuration.*"),
			},
			{
				Config:      createFirmwareCatalogResourceValidateErrorHTTP,
				ExpectError: regexp.MustCompile(".*invalid HTTP share configuration.*"),
			},
			{
				Config:      createFirmwareCatalogResourceValidateErrorHTTPS,
				ExpectError: regexp.MustCompile(".*invalid HTTPS share configuration.*"),
			},
		},
	})
}

func TestFirmwareCatalogResourceReadCreateUpdateErrors(t *testing.T) {
	if os.Getenv("TF_ACC") == "0" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	createMock := models.CatalogsModel{
		ID: 1,
		Repository: models.RepositoryModel{
			BackupExistingCatalog: false,
			CheckCertificate:      false,
			Description:           "tfacc_firmware_catalog_resource terraform catalog",
			DomainName:            "example",
			Editable:              true,
			Name:                  CatalogResource,
			RepositoryType:        "HTTPS",
			Source:                "https://1.2.2.1",
			Username:              "example-user",
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.CreateCatalogFirmware, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      createFirmwareCatalogResource,
				ExpectError: regexp.MustCompile(`.*Unable to create catalog*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateCatalogFirmware, OptGeneric).Return(createMock, nil).Build()
					localFunctionalMocker = Mock(helper.SetStateCatalogFirmware).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      createFirmwareCatalogResource,
				ExpectError: regexp.MustCompile(`.*Unable to process catalog after create*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localFunctionalMocker != nil {
						localFunctionalMocker.UnPatch()
					}
					localFunctionalMocker = Mock(helper.GetSpecificCatalogFirmware, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      createFirmwareCatalogResource,
				ExpectError: regexp.MustCompile(`.*Unable to read specific firmware catalog*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localFunctionalMocker != nil {
						localFunctionalMocker.UnPatch()
					}
				},
				Config:      updateFirmwareVaildationError,
				ExpectError: regexp.MustCompile(`.*Unable to update catalog, validation error:*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localFunctionalMocker != nil {
						localFunctionalMocker.UnPatch()
					}
					localFunctionalMocker = Mock(helper.UpdateCatalogFirmware, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      updateFirmwareMockError,
				ExpectError: regexp.MustCompile(`.*Unable to update catalog:*.`),
			},
		},
	})
}

var updateFirmwareMockError = testProvider + `
	resource "ome_firmware_catalog" "cat_1" {
		name = "` + CatalogResource + `_update"
		catalog_update_type = "Manual"
		share_type = "HTTPS"
		catalog_file_path = "catalogs/example_catalog_1.xml"
        share_address = "https://2.1.1.2"
        catalog_refresh_schedule = {
          cadence = "Weekly"
          day_of_the_week = "Wednesday"
          time_of_day = "6"
          am_pm = "PM"
        }
        domain = "example"
        share_user = "example-user"
        share_password = "example-pass"
	}
`

var updateFirmwareVaildationError = testProvider + `
	resource "ome_firmware_catalog" "cat_1" {
		name = "` + CatalogResource + `"
		catalog_update_type = "Manual"
		share_type = "CIFS"
		catalog_file_path = "catalogs/example_catalog_1.xml"
        share_address = "https://1.2.2.1"
        catalog_refresh_schedule = {
          cadence = "Weekly"
          day_of_the_week = "Wednesday"
          time_of_day = "6"
          am_pm = "PM"
        }
        domain = "example"
        share_user = "example-user"
        share_password = "example-pass"
	}
`

var createFirmwareCatalogResourceValidateErrorAutomatic = testProvider + `
	resource "ome_firmware_catalog" "cat_1" {
		name = "` + CatalogResource + `_validate"
		catalog_update_type = "Automatic"		
	}
`

var createFirmwareCatalogResourceValidateErrorNFS = testProvider + `
	resource "ome_firmware_catalog" "cat_1" {
		name = "` + CatalogResource + `_validate"
		catalog_update_type = "Manual"
		share_type = "NFS"
	}
`

var createFirmwareCatalogResourceValidateErrorCIFS = testProvider + `
	resource "ome_firmware_catalog" "cat_1" {
		name = "` + CatalogResource + `_validate"
		catalog_update_type = "Manual"
		share_type = "CIFS"
	}
`

var createFirmwareCatalogResourceValidateErrorHTTP = testProvider + `
	resource "ome_firmware_catalog" "cat_1" {
		name = "` + CatalogResource + `_validate"
		catalog_update_type = "Manual"
		share_type = "HTTP"
	}
`

var createFirmwareCatalogResourceValidateErrorHTTPS = testProvider + `
	resource "ome_firmware_catalog" "cat_1" {
		name = "` + CatalogResource + `_validate"
		catalog_update_type = "Manual"
		share_type = "HTTPS"
	}
`

var createFirmwareCatalogResource = testProvider + `
    resource "ome_firmware_catalog" "cat_1" {
        name = "` + CatalogResource + `"
        catalog_update_type = "Manual"
        share_type = "HTTPS"
        catalog_file_path = "catalogs/example_catalog_1.xml"
        share_address = "https://1.2.2.1"
        catalog_refresh_schedule = {
          cadence = "Weekly"
          day_of_the_week = "Wednesday"
          time_of_day = "6"
          am_pm = "PM"
        }
        domain = "example"
        share_user = "example-user"
        share_password = "example-pass"
    }
`

var updateFirmwareCatalogResource = testProvider + `
resource "ome_firmware_catalog" "cat_1" {
    name = "` + CatalogResource + "_update" + `"
    catalog_update_type = "Automatic"
    share_type = "HTTPS"
    catalog_file_path = "catalogs/example_catalog_1.xml"
    share_address = "https://2.1.1.2"
    catalog_refresh_schedule = {
      cadence = "Weekly"
      day_of_the_week = "Tuesday"
      time_of_day = "8"
      am_pm = "AM"
    }
    domain = "example_update"
    share_user = "example-user_update"
    share_password = "example-pass_update"
}
`
