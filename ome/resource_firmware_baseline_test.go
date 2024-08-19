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

var localMocker *Mocker
var localMocker2 *Mocker

const (
	FirmwareBaselineName          = "test_acc_fw_baseline"
	FirmwareBaseline2Name         = "test_acc_fw_baseline_2"
	FirmwareBaseline3Name         = "test_acc_fw_baseline_3"
	FirmwareBaselineErrName       = "test_acc_fw_baseline_err"
	FirmwareBaselineErr2Name      = "test_acc_fw_baseline_err_2"
	FirmwareBaselineNameUpdate    = "test_acc_update_fw_baseline"
	FirmwareBaselineNameUpdateErr = "test_acc_update_fw_baseline_err"
)

func TestFirmwareBaselineResourceUpdateFail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:       createFirmwareBaselineDeviceResourceupd,
				ResourceName: "ome_firmware_baseline.firmware_baseline_upd",
			},
			{
				Config:       updateFirmwareBaselineResourceErr,
				ResourceName: "ome_firmware_baseline.firmware_baseline_upd",
				ExpectError:  regexp.MustCompile(".*Unable to Update Baseline:*"),
			},
		},
	})
}

func TestFirmwareBaselineResourceCreate(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	var fimwareBaselineCreate = "ome_firmware_baseline.firmware_baseline"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: createFirmwareBaselineDeviceResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "catalog_name", "tfacc_catalog_dell_online_1"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "name", FirmwareBaselineName),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "filter_no_reboot_required", "false"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "targets.#", "1"),
				),
			},
			{
				Config: updateFirmwareBaselineResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "catalog_name", "tfacc_catalog_dell_online_1"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "name", FirmwareBaselineNameUpdate),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "filter_no_reboot_required", "false"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "targets.#", "1"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "description", "test baseline updated"),
				),
			},
			// Import testing
			{
				ResourceName: fimwareBaselineCreate,
				ImportState:  true,
			},
		},
	})
}

func TestFirmwareBaselineResourceCreateGroup(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	var fimwareBaselineCreate = "ome_firmware_baseline.firmware_baseline2"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFirmwareBaselineGroupResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "catalog_name", "tfacc_catalog_dell_online_1"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "name", FirmwareBaseline2Name),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "filter_no_reboot_required", "false"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "targets.#", "5"),
				),
			},
		},
	})
}

func TestFirmwareBaselineResourceCreateServiceTag(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	var fimwareBaselineCreate = "ome_firmware_baseline.firmware_baseline3"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFirmwareBaselineServiceTagResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "catalog_name", "tfacc_catalog_dell_online_1"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "name", FirmwareBaseline3Name),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "filter_no_reboot_required", "false"),
					resource.TestCheckResourceAttr(fimwareBaselineCreate, "targets.#", "1"),
				),
			},
		},
	})
}

func TestFirmwareBaselineResourceCreateError(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      createFirmwareBaselineResourceError,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Combination.*"),
			},
			{
				Config:      createFirmwareBaselineResourceError2,
				ExpectError: regexp.MustCompile(".*Unable to create target model for*"),
			},
			{
				Config:      createFirmwareBaselineResourceError3,
				ExpectError: regexp.MustCompile(".*Catalog details not found*"),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = Mock(helper.CreateTargetModel).Return([]models.TargetModel{}, nil).Build()

				},
				Config:      createFirmwareBaselineResourceError4,
				ExpectError: regexp.MustCompile(".*Unable to create target model for*"),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = MockGeneric(helper.CreateFirmwareBaseline).Return(int64(0), fmt.Errorf("mock Error")).Build()

				},
				Config:      createFirmwareBaselineDeviceResource,
				ExpectError: regexp.MustCompile(".*Unable to create Baseline*"),
			},
		},
	})
}

func TestFirmwareBaselineResourceGetFailure(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					if localMocker != nil {
						localMocker.Release()
					}
					if localMocker2 != nil {
						localMocker2.Release()
					}
					FunctionMocker = Mock(helper.CreateFirmwareBaseline, OptGeneric).Return(int64(0), nil).Build()
					localMocker = Mock(helper.GetFirmwareBaselineWithName).Return(models.FirmwareBaselinesModel{}, nil).Build()
					localMocker2 = Mock(helper.SetStateBaseline).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      createFirmwareBaselineDeviceResource,
				ExpectError: regexp.MustCompile(".*Could not copy Baseline data*"),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					if localMocker != nil {
						localMocker.Release()
					}
					if localMocker2 != nil {
						localMocker2.Release()
					}

					localMocker = Mock(helper.GetAllCatalogFirmware).Return(nil, fmt.Errorf("Mock error")).Build()
					FunctionMocker = Mock(helper.GetCatalogFirmwareByName).Return(nil, nil).Build()

				},

				Config:      createFirmwareBaselineDeviceResource,
				ExpectError: regexp.MustCompile(".*Catalog details not found*"),
			},
		},
	})
}

func TestFirmwareBaselineResourceImportFail(t *testing.T) {
	var fimwareBaselineCreate = "ome_firmware_baseline.firmware_baseline"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            importFailureBaselineResource,
				ResourceName:      fimwareBaselineCreate,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "invalid_state_id",
				ExpectError:       regexp.MustCompile(".*Unable to import firmware baseline*"),
			},
		},
	})
}

var importFailureBaselineResource = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline" {
	catalog_name = "` + Catalog1 + `"
   	device_names = ["10.226.197.28"]
	name = "` + FirmwareBaselineName + `"
	description = "test baseline"
  }
`
var createFirmwareBaselineDeviceResource = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline" {
	catalog_name = "` + Catalog1 + `"
   	device_names = ["10.226.197.28"]
	name = "` + FirmwareBaselineName + `"
	description = "test baseline"
  }
`
var createFirmwareBaselineGroupResource = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline2" {
	catalog_name = "` + Catalog1 + `"
	group_names = ["HCI Appliances","Hyper-V Servers"]
	name =  "` + FirmwareBaseline2Name + `"
	description = "test baseline"
  }
`
var createFirmwareBaselineServiceTagResource = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline3" {
	catalog_name = "` + Catalog1 + `"
	device_service_tags = ["HRPB0M3"]
	name =  "` + FirmwareBaseline3Name + `"
	description = "test baseline"
  }
`
var createFirmwareBaselineResourceError = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline" {
	catalog_name = "` + Catalog1 + `"
	device_service_tags = ["HRPB0M3"]
	group_names = ["HCI Appliances","Hyper-V Servers"]
	name =  "` + FirmwareBaseline2Name + `"
	#is_64_bit = false
	#filter_no_reboot_required = true
	description = "test baseline"
  }
`
var createFirmwareBaselineResourceError2 = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline" {
	catalog_name = "` + Catalog1 + `"
	name =  "` + FirmwareBaseline2Name + `"
	#is_64_bit = false
	#filter_no_reboot_required = true
	description = "test baseline"
  }
`
var createFirmwareBaselineResourceError3 = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline" {
	catalog_name = "Invalid Catalog"
	device_service_tags = ["HRPB0M3"]
	name = "` + FirmwareBaselineErrName + `"
	#is_64_bit = false
	#filter_no_reboot_required = true
	description = "test baseline"
  }
`
var createFirmwareBaselineResourceError4 = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline" {
	catalog_name =  "` + Catalog1 + `"
	device_service_tags = ["empty-target"]
	name = "` + FirmwareBaselineErrName + `"
	#is_64_bit = false
	#filter_no_reboot_required = true
	description = "test baseline"
  }
`

var updateFirmwareBaselineResource = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline" {
	catalog_name = "` + Catalog1 + `"
   	device_names = ["10.226.197.28"]
	name = "` + FirmwareBaselineNameUpdate + `"
	#is_64_bit = false
	#filter_no_reboot_required = true
	description = "test baseline updated"
  }
`
var createFirmwareBaselineDeviceResourceupd = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline_upd" {
	catalog_name = "` + Catalog1 + `"
   	device_names = ["10.226.197.28"]
	name =  "` + FirmwareBaselineName + `"
	description = "test baseline"
  }
`
var updateFirmwareBaselineResourceErr = testProvider + `
resource "ome_firmware_baseline" "firmware_baseline_upd" {
	catalog_name = "InvalidCat"
	device_names = ["10.226.197.28"]
	name = "` + FirmwareBaselineNameUpdateErr + `"
	#is_64_bit = false
	#filter_no_reboot_required = true
	description = "test baseline updated"
  }
`
