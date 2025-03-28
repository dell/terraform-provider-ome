/*
Copyright (c) 2024-2025 Dell Inc., or its subsidiaries. All Rights Reserved.
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
	"os"
	"strconv"
	"strings"
	"terraform-provider-ome/clients"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultPort             int64         = 443
	defaultTimeoutInSeconds int           = 30
	defaultTimeout          time.Duration = time.Second * time.Duration(defaultTimeoutInSeconds)
	defaultProtocol         string        = "https"
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
	Protocol types.String `tfsdk:"protocol"`
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

	// Attempt to read from Env Variables, otherwise use defaults
	usernameEnv := os.Getenv("OME_USERNAME")
	if usernameEnv != "" {
		data.Username = types.StringValue(usernameEnv)
	}

	if data.Username.IsUnknown() && usernameEnv == "" {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return
	}

	// Attempt to read from Env Variables, otherwise use defaults
	hostEnv := os.Getenv("OME_HOST")
	if hostEnv != "" {
		data.Host = types.StringValue(hostEnv)
	}
	passEnv := os.Getenv("OME_PASSWORD")
	if passEnv != "" {
		data.Password = types.StringValue(passEnv)
	}
	protocolEnv := os.Getenv("OME_PROTOCOL")
	if protocolEnv != "" {
		if protocolEnv != "http" && protocolEnv != "https" {
			resp.Diagnostics.AddError(
				"Invalid protocol",
				"Protocol must be http or https",
			)
			return
		}
		data.Protocol = types.StringValue(protocolEnv)
	}

	portEnv, errPort := strconv.ParseInt(os.Getenv("OME_PORT"), 10, 64)
	if errPort == nil {
		data.Port = types.Int64Value(portEnv)
	}
	skipSslEnv, errSkipSsl := strconv.ParseBool(os.Getenv("OME_SKIP_SSL"))
	if errSkipSsl == nil {
		data.SkipSSL = types.BoolValue(skipSslEnv)
	}
	timeoutEnv, errTimeout := strconv.ParseInt(os.Getenv("OME_TIMEOUT"), 10, 64)
	if errTimeout == nil {
		data.Timeout = types.Int64Value(timeoutEnv)
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
	//Default https to https
	https := defaultProtocol
	if !data.Protocol.IsNull() {
		https = data.Protocol.ValueString()
	}

	if data.SkipSSL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as SkipSSL",
		)
		return
	}

	url := clients.GetURL(https, data.Host.ValueString(), port)

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
		NewDiscoveryResource,
		NewNetworkSettingResource,
		NewCsrResource,
		NewDevicesResource,
		NewCertResource,
		NewDeviceActionResource,
		NewFirmwareCatalogResource,
		NewFirmwareBaselineResource,
	}
}

func (p *omeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTemplateDataSource,
		NewGroupDevicesDatasource,
		NewVlanNetworkDataSource,
		NewConfigurationReportDataSource,
		NewDeviceDatasource,
		NewAppCertDataSource,
		NewFirmwareCatalogDataSource,
		NewFirmwareBaselineComplianceRepositoryDatasource,
		NewfwBaselineCompReportDatasource,
		NewDeviceComplianceReportDataSource,
	}
}

func (p *omeProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Terraform Provider for OpenManage Enterprise (OME) is a plugin for Terraform that allows the resource management of PowerEdge servers using OME",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "OpenManage Enterprise IP address or hostname. This can also be set using the environment variable OME_HOST",
				Description:         "OpenManage Enterprise IP address or hostname. This can also be set using the environment variable OME_HOST",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "OpenManage Enterprise username. This can also be set using the environment variable OME_USERNAME",
				Description:         "OpenManage Enterprise username. This can also be set using the environment variable OME_USERNAME",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "OpenManage Enterprise password. This can also be set using the environment variable OME_PASSWORD",
				Description:         "OpenManage Enterprise password. This can also be set using the environment variable OME_PASSWORD",
				// This should remain optional so user can use environment variables if they choose.
				Optional:  true,
				Sensitive: true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "OpenManage Enterprise HTTPS port. This can also be set using the environment variable OME_PORT" +
					fmt.Sprintf(" Default value is `%d`.", defaultPort),
				Description: "OpenManage Enterprise HTTPS port. This can also be set using the environment variable OME_PORT" +
					fmt.Sprintf(" Default value is '%d'.", defaultPort),
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"skipssl": schema.BoolAttribute{
				MarkdownDescription: "Skips SSL certificate validation on OpenManage Enterprise. This can also be set using the environment variable OME_SKIP_SSL" +
					" Default value is `false`.",
				Description: "Skips SSL certificate validation on OpenManage Enterprise. This can also be set using the environment variable OME_SKIP_SSL" +
					" Default value is 'false'.",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "HTTPS timeout in seconds for OpenManage Enterprise client. This can also be set using the environment variable OME_TIMEOUT" +
					fmt.Sprintf(" Default value is `%d`.", defaultTimeoutInSeconds),
				Description: "HTTPS timeout in seconds for OpenManage Enterprise client. This can also be set using the environment variable OME_TIMEOUT" +
					fmt.Sprintf(" Default value is '%d'.", defaultTimeoutInSeconds),
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Set the Http protocol for OpenManage Enterprise client. This can also be set using the environment variable OME_PROTOCOL" +
					fmt.Sprintf(" Default value is `%s`.", defaultProtocol),
				Description: "Set the Http protocol for OpenManage Enterprise client. This can also be set using the environment variable OME_PROTOCOL" +
					fmt.Sprintf(" Default value is '%s'.", defaultProtocol),
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"http",
						"https",
					}...),
				},
			},
		},
	}
}

// helper for schema generation of accepted values
func makeSchemaAcceptedValues(inputs []string, quote string) string {
	inQuote := make([]string, 0)
	for _, in := range inputs {
		inQuote = append(inQuote, fmt.Sprintf("%s%s%s", quote, in, quote))
	}
	return fmt.Sprintf(" Accepted values are %s.", strings.Join(inQuote, ", "))
}
