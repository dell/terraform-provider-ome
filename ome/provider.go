package ome

import (
	"context"
	"fmt"
	"terraform-provider-ome/clients"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultPort    int64         = 443
	defaultTimeout time.Duration = time.Second * 30
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.Provider = &provider{}

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client options can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	clientOpt *clients.ClientOptions

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
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

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	tflog.Trace(ctx, "Started configuring the provider")
	data := providerData{}
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if data.Username.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return
	}

	if data.Username.Value == "" {
		resp.Diagnostics.AddError(
			"Unable to find username",
			"Username cannot be an empty string",
		)
		return
	}

	if data.Password.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return
	}

	if data.Password.Value == "" {
		resp.Diagnostics.AddError(
			"Unable to find ome password",
			"password cannot be an empty string",
		)
		return
	}

	if data.Host.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return
	}

	if data.Host.Value == "" {
		resp.Diagnostics.AddError(
			"Unable to find ome host.",
			"host cannot be an empty string",
		)
		return
	}

	//Default port to 443
	port := defaultPort
	if data.Port.Value != 0 {
		port = data.Port.Value
	}
	//Default timeout to 30 sec
	timeout := defaultTimeout
	if data.Timeout.Value != 0 {
		timeout = time.Second * time.Duration(data.Timeout.Value)
	}

	if data.SkipSSL.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as SkipSSL",
		)
		return
	}

	url := clients.GetURL(data.Host.Value, port)

	tflog.Info(ctx, "Collected all data creating client options")

	clientOptions := clients.ClientOptions{
		Username:       data.Username.Value,
		Password:       data.Password.Value,
		URL:            url,
		SkipSSL:        data.SkipSSL.Value,
		Timeout:        timeout,
		Retry:          clients.Retries,
		PreRequestHook: clients.ClientPreReqHook,
	}
	p.clientOpt = &clientOptions

	p.configured = true

	tflog.Trace(ctx, "Finished configuring the provider")
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"ome_template":                  resourceTemplateType{},
		"ome_deployment":                resourceDeploymentType{},
		"ome_configuration_baseline":    resourceConfigurationBaselineType{},
		"ome_configuration_compliance": resourceConfigurationComplianceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"ome_template_info":     templateDataSourceType{},
		"ome_groupdevices_info": groupDevicesDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "The Terraform Provider for OpenManage Enterprise (OME) is a plugin for Terraform that allows the resource management of PowerEdge servers using OME",
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				MarkdownDescription: "OpenManage Enterprise IP address or hostname.",
				Description:         "OpenManage Enterprise IP address or hostname.",
				Type:                types.StringType,
				Required:            true,
			},
			"username": {
				MarkdownDescription: "OpenManage Enterprise username.",
				Description:         "OpenManage Enterprise username.",
				Type:                types.StringType,
				Required:            true,
			},
			"password": {
				MarkdownDescription: "OpenManage Enterprise password.",
				Description:         "OpenManage Enterprise password.",
				Type:                types.StringType,
				Required:            true,
				Sensitive:           true,
			},
			"port": {
				MarkdownDescription: "OpenManage Enterprise HTTPS port.",
				Description:         "OpenManage Enterprise HTTPS port.",
				Type:                types.Int64Type,
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 443}),
				},
			},
			"skipssl": {
				MarkdownDescription: "Skips SSL certificate validation on OpenManage Enterprise",
				Description:         "Skips SSL certificate validation on OpenManage Enterprise",
				Type:                types.BoolType,
				Optional:            true,
			},
			"timeout": {
				MarkdownDescription: "HTTPS timeout for OpenManage Enterprise client",
				Description:         "HTTPS timeout for OpenManage Enterprise client",
				Type:                types.Int64Type,
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					DefaultAttribute(types.Int64{Value: 30}),
				},
			},
		},
	}, nil
}

// New method is used to create a new provider via a RPC call or from main
func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
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
