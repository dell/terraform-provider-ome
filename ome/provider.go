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
	"terraform-provider-ome/clients"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultPort             int64         = 443
	defaultTimeoutInSeconds int           = 30
	defaultTimeout          time.Duration = time.Second * time.Duration(defaultTimeoutInSeconds)
)

var (
	_ provider.Provider = &omeProvider{}
)

// New - returns new provider struct definition.
func New() provider.Provider {
	return &omeProvider{}
}

type omeProvider struct {
	// client options can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	clientOpt *clients.ClientOptions

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	Username types.String `tfsdk:"username"`
	Host     types.String `tfsdk:"host"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	SkipSSL  types.Bool   `tfsdk:"skipssl"`
	Timeout  types.Int64  `tfsdk:"timeout"`
}

// Metadata - provider metadata AKA name.
func (p *omeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ome_"
}

// Configure - provider pre-initiate calle function.
func (p *omeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	tflog.Trace(ctx, "Started configuring the provider")
	data := providerData{}
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if data.Username.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return
	}

	if data.Username.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Unable to find username",
			"Username cannot be an empty string",
		)
		return
	}

	if data.Password.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return
	}

	if data.Password.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Unable to find ome password",
			"password cannot be an empty string",
		)
		return
	}

	if data.Host.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return
	}

	if data.Host.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Unable to find ome host.",
			"host cannot be an empty string",
		)
		return
	}

	//Default port to 443
	port := defaultPort
	if data.Port.ValueInt64() != 0 {
		port = data.Port.ValueInt64()
	}
	//Default timeout to 30 sec
	timeout := defaultTimeout
	if data.Timeout.ValueInt64() != 0 {
		timeout = time.Second * time.Duration(data.Timeout.ValueInt64())
	}

	if data.SkipSSL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as SkipSSL",
		)
		return
	}

	url := clients.GetURL(data.Host.ValueString(), port)

	tflog.Info(ctx, "Collected all data creating client options")

	clientOptions := clients.ClientOptions{
		Username:       data.Username.ValueString(),
		Password:       data.Password.ValueString(),
		URL:            url,
		SkipSSL:        data.SkipSSL.ValueBool(),
		Timeout:        timeout,
		Retry:          clients.Retries,
		PreRequestHook: clients.ClientPreReqHook,
	}
	p.clientOpt = &clientOptions

	p.configured = true
	resp.DataSourceData = p
	resp.ResourceData = p

	tflog.Trace(ctx, p.clientOpt.Username)
	tflog.Trace(ctx, "Finished configuring the provider")
}

func (p *omeProvider) createOMESession(ctx context.Context, caller string) (*clients.Client, diag.Diagnostics) {
	//Create Session and defer the remove session
	var d diag.Diagnostics
	omeClient, err := clients.NewClient(*p.clientOpt)
	if err != nil {
		d.AddError(
			clients.ErrCreateClient,
			err.Error(),
		)
		return nil, d
	}

	tflog.Trace(ctx, fmt.Sprintf("resource_configuration_baseline %s Creating Session", caller))
	_, err = omeClient.CreateSession()
	if err != nil {
		d.AddError(
			clients.ErrCreateSession,
			err.Error(),
		)
		return nil, d
	}
	return omeClient, d
}

func (p *omeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTemplateResource,
		NewDeploymentResource,
		NewConfigurationBaselineResource,
		NewConfigurationComplianceResource,
		NewStaticGroupResource,
		NewUserResource,
	}
}

func (p *omeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTemplateDataSource,
		NewGroupDevicesDatasource,
		NewVlanNetworkDataSource,
		NewConfigurationReportDataSource,
	}
}

func (p *omeProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Terraform Provider for OpenManage Enterprise (OME) is a plugin for Terraform that allows the resource management of PowerEdge servers using OME",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "OpenManage Enterprise IP address or hostname.",
				Description:         "OpenManage Enterprise IP address or hostname.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "OpenManage Enterprise username.",
				Description:         "OpenManage Enterprise username.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "OpenManage Enterprise password.",
				Description:         "OpenManage Enterprise password.",
				Required:            true,
				Sensitive:           true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "OpenManage Enterprise HTTPS port." +
					fmt.Sprintf(" Default value is `%d`.", defaultPort),
				Description: "OpenManage Enterprise HTTPS port." +
					fmt.Sprintf(" Default value is '%d'.", defaultPort),
				Optional: true,
			},
			"skipssl": schema.BoolAttribute{
				MarkdownDescription: "Skips SSL certificate validation on OpenManage Enterprise." +
					" Default value is `false`.",
				Description: "Skips SSL certificate validation on OpenManage Enterprise." +
					" Default value is 'false'.",
				Optional: true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "HTTPS timeout in seconds for OpenManage Enterprise client." +
					fmt.Sprintf(" Default value is `%d`.", defaultTimeoutInSeconds),
				Description: "HTTPS timeout in seconds for OpenManage Enterprise client." +
					fmt.Sprintf(" Default value is '%d'.", defaultTimeoutInSeconds),
				Optional: true,
			},
		},
	}
}

// New method is used to create a new provider via a RPC call or from main
/* func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
} */

/* // convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
//
//lint:ignore U1000 used by the internal provider, to be checked
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
*/
