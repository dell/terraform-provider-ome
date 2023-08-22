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
	"context"
	"fmt"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &resourceCert{}
	_ resource.ResourceWithConfigure = &resourceCert{}
)

// NewCertResource is new resource for application Cert
func NewCertResource() resource.Resource {
	return &resourceCert{}
}

type resourceCert struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceCert) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (r resourceCert) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "application_certificate"
}

// Cert Resource schema
func (r resourceCert) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for generating application Certificate Signing Request from OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID for application Cert resource.",
				Description:         "ID for application Cert resource.",
				Computed:            true,
			},
			"certificate_base64": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded certificate.",
				Description:         "Base64 encoded certificate.",
				Required:            true,
			},
		},
	}
}

// Create a new resource
func (r resourceCert) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_Cert create: started")
	var plan models.CertResModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, dgs := r.uploadCert(ctx, plan)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceCert) uploadCert(ctx context.Context, plan models.CertResModel) (models.CertResModel, diag.Diagnostics) {
	//Create Session and defer the remove session
	omeClient, dgs := r.p.createOMESession(ctx, "resource_cert Upload")
	if dgs.HasError() {
		return plan, dgs
	}
	defer omeClient.RemoveSession()

	devs, errd := omeClient.GetAllDevices(nil)
	fmt.Sprintln(devs, errd)

	tflog.Info(ctx, "resource_cert uploading Cert")

	_, err := omeClient.PostCert(plan.Cert.ValueString())
	if err != nil {
		dgs.AddError(
			"Error uploading Cert.",
			err.Error(),
		)
		return plan, dgs
	}

	state := models.CertResModel{
		ID:   types.StringValue("dummy"),
		Cert: plan.Cert,
	}
	return state, dgs
}

// Read resource information
func (r resourceCert) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// read refresh changes nothing
	resp.State = req.State
}

// Update resource
func (r resourceCert) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_Cert update: started")
	var plan models.CertResModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, dgs := r.uploadCert(ctx, plan)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r resourceCert) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Just remove State Data
	resp.State.RemoveResource(ctx)
}
