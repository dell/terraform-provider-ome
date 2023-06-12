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
	"testing"
)

type testTemplates struct {
	templateSvcTag1       string
	templateSvcTag2       string
	templateSvcTag1Full   string
	templateDeploySvcTag1 string
}

func initTemplates(t *testing.T) testTemplates {
	omeTestdataDir, _ := os.LookupEnv("OME_TESTDATA_DIR")
	if omeTestdataDir == "" {
		t.Error("The environment variable OME_TESTDATA_DIR must be set for this test.")
	}
	templateSvcTag1FullFileName := "test_acc_template_full_svc_tag_1.xml"
	templateSvcTag1FileName := "test_acc_template_compliance_svc_tag_1.xml"
	templateSvcTag2FileName := "test_acc_template_compliance_svc_tag_2.xml"
	templateDeploySvcTag1FileName := "test_acc_template_deploy_svc_tag_1.xml"
	if _, err := os.Stat(omeTestdataDir + "/" + templateSvcTag1FileName); err != nil {
		t.Error(err.Error())
	}
	if _, err := os.Stat(omeTestdataDir + "/" + templateSvcTag2FileName); err != nil {
		t.Error(err.Error())
	}
	if _, err := os.Stat(omeTestdataDir + "/" + templateSvcTag1FullFileName); err != nil {
		t.Error(err.Error())
	}
	if _, err := os.Stat(omeTestdataDir + "/" + templateDeploySvcTag1FileName); err != nil {
		t.Error(err.Error())
	}
	templateSvcTag1 := `
	resource "ome_template" "terraform-acceptance-test-1" {
		view_type = "Compliance"
		name = "%s"
		content = file("%s/%s")
	}
	`
	templateSvcTag2 := `
	resource "ome_template" "terraform-acceptance-test-2" {
		view_type = "Compliance"
		name = "%s"
		content = file("%s/%s")
	}
	`

	templateSvcTag1Full := `
	resource "ome_template" "terraform-acceptance-test-1" {
		view_type = "Compliance"
		name = "%s"
		content = file("%s/%s")
	}
	`
	templateDeploySvcTag1 := `
	resource "ome_template" "terraform-acceptance-test-1" {
		name = "%s"
		content = file("%s/%s")
	}
	`
	var ret testTemplates
	ret.templateSvcTag1 = fmt.Sprintf(templateSvcTag1, TestRefTemplateName, omeTestdataDir, templateSvcTag1FileName)
	ret.templateSvcTag1Full = fmt.Sprintf(templateSvcTag1Full, TestRefTemplateName, omeTestdataDir, templateSvcTag1FullFileName)
	ret.templateSvcTag2 = fmt.Sprintf(templateSvcTag2, TestRefTemplateNameUpdate, omeTestdataDir, templateSvcTag2FileName)
	ret.templateDeploySvcTag1 = fmt.Sprintf(templateDeploySvcTag1, TestAccTemplateName, omeTestdataDir, templateDeploySvcTag1FileName)
	return ret
}
