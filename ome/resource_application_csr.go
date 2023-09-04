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
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &resourceCsr{}
	_ resource.ResourceWithConfigure = &resourceCsr{}
)

// NewCsrResource is new resource for application csr
func NewCsrResource() resource.Resource {
	return &resourceCsr{}
}

type resourceCsr struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *resourceCsr) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata implements resource.Resource
func (r resourceCsr) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "application_csr"
}

// CSR Resource schema
func (r resourceCsr) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for generating application Certificate Signing Request from OpenManage Enterprise.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID for application csr resource.",
				Description:         "ID for application csr resource.",
				Computed:            true,
			},
			"specs": schema.SingleNestedAttribute{
				MarkdownDescription: "CSR specifications." +
					" Terraform will replace (delete and recreate) this resource if this attribute is modified.",
				Description: "CSR specifications." +
					" Terraform will replace (delete and recreate) this resource if this attribute is modified.",
				Attributes: r.specSchema(),
				Required:   true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"csr": schema.StringAttribute{
				MarkdownDescription: "CSR in single line PEM format returned from OME.",
				Description:         "CSR in single line PEM format returned from OME.",
				Computed:            true,
			},
		},
	}
}

func (*resourceCsr) specSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"distinguished_name": schema.StringAttribute{
			MarkdownDescription: "Distinguished Name.",
			Description:         "Distinguished Name.",
			Required:            true,
		},
		"department_name": schema.StringAttribute{
			MarkdownDescription: "Department Name.",
			Description:         "Department Name.",
			Required:            true,
		},
		"business_name": schema.StringAttribute{
			MarkdownDescription: "Business Name.",
			Description:         "Business Name.",
			Required:            true,
		},
		"locality": schema.StringAttribute{
			MarkdownDescription: "Locality of the business.",
			Description:         "Locality of the business.",
			Required:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "State of the business.",
			Description:         "State of the business.",
			Required:            true,
		},
		"country": schema.StringAttribute{
			MarkdownDescription: "Country of the business.",
			Description:         "Country of the business.",
			Required:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "Email address.",
			Description:         "Email address.",
			Required:            true,
		},
		"subject_alternate_names": schema.ListAttribute{
			MarkdownDescription: "Subject Alternate names. Maximum 4.",
			Description:         "Subject Alternate names. Maximum 4.",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeBetween(1, 4),
			},
		},
	}
}

// Create a new resource
func (r resourceCsr) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Get Plan Data
	tflog.Trace(ctx, "resource_csr create: started")
	var plan models.CsrResModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_csr Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	tflog.Info(ctx, "resource_csr generating csr")

	state, err := r.genCSR(ctx, plan, omeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating CSR.",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceCsr) genCSR(ctx context.Context, plan models.CsrResModel, c *clients.Client) (models.CsrResModel, error) {
	state := models.CsrResModel{
		ID:    types.StringValue("dummy"),
		Specs: plan.Specs,
	}
	csr, err := c.GetCSR(plan.Specs.GetCsrConfig(ctx))
	state.Csr = types.StringValue(csr)
	return state, err
}

// Read resource information
func (r resourceCsr) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// read refresh changes nothing
	resp.State = req.State
}

// Update resource
func (r resourceCsr) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update should never happen
	resp.Diagnostics.AddError(
		"Error updating CSR.",
		"An update plan of CSR should never be invoked. This resource is supposed to be replaced on update.",
	)
}

// Delete resource
func (r resourceCsr) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Just remove State Data
	resp.State.RemoveResource(ctx)
}
