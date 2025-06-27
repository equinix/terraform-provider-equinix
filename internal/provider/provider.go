// Package provider implements the Terraform provider for Equinix, including provider configuration,
// resource and data source registration, and integration with the Terraform Plugin Framework.
package provider

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/provider/services"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// FrameworkProvider implements the Terraform provider, holding version and configuration metadata.
type FrameworkProvider struct {
	ProviderVersion string
	Meta            *config.Config
}

// CreateFrameworkProvider initializes a new FrameworkProvider with the specified version.
func CreateFrameworkProvider(version string) provider.ProviderWithMetaSchema {
	return &FrameworkProvider{
		ProviderVersion: version,
	}
}

// Metadata returns the provider's metadata, such as its type name, to the Terraform framework.
func (p *FrameworkProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "equinixcloud"
}

// Schema returns the provider's schema, which defines the configuration options available to users.
func (p *FrameworkProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: fmt.Sprintf("The Equinix API base URL to point out desired environment. This argument can also be specified with the `EQUINIX_API_ENDPOINT` shell environment variable. (Defaults to `%s`)", config.DefaultBaseURL),
				Validators: []validator.String{
					equinix_validation.URLWithScheme("http", "https"),
				},
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "API Consumer Key available under \"My Apps\" in developer portal. This argument can also be specified with the `EQUINIX_API_CLIENTID` shell environment variable.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "API Consumer secret available under \"My Apps\" in developer portal. This argument can also be specified with the `EQUINIX_API_CLIENTSECRET` shell environment variable.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "API tokens are generated from API Consumer clients using the [OAuth2 API](https://developer.equinix.com/dev-docs/fabric/getting-started/getting-access-token#request-access-and-refresh-tokens). This argument can also be specified with the `EQUINIX_API_TOKEN` shell environment variable.",
			},
			"auth_token": schema.StringAttribute{
				Optional:    true,
				Description: "The Equinix Metal API auth key for API operations",
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: fmt.Sprintf("The duration of time, in seconds, that the Equinix Platform API Client should wait before canceling an API request. Canceled requests may still result in provisioned resources. (Defaults to `%d`)", config.DefaultTimeout),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"response_max_page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum number of records in a single response for REST queries that produce paginated responses. (Default is client specific)",
				Validators: []validator.Int64{
					int64validator.AtLeast(100),
				},
			},
			"sts_auth_scope": schema.StringAttribute{
				Optional:    true,
				Description: "The scope of the authentication token. Must be an access policy ERN or a string of the form roleassignments:<org_id> This argument can also be specified with the `EQUINIX_STS_AUTH_SCOPE` shell environment variable.",
			},
			"sts_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: fmt.Sprintf("The STS API base URL to point out desired environment. This argument can also be specified with the `EQUINIX_STS_ENDPOINT` shell environment variable. (Defaults to `%s`)", config.DefaultStsBaseURL),
				Validators: []validator.String{
					equinix_validation.URLWithScheme("http", "https"),
				},
			},
			"sts_source_token": schema.StringAttribute{
				Optional:    true,
				Description: "The source token to use for STS authentication. Must be an OIDC ID token issued by an OIDC provider trusted by Equinix STS. This argument can also be specified with the `EQUINIX_STS_SOURCE_TOKEN` shell environment variable.",
			},
			"max_retries": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of retries in case of network failure.",
			},
			"max_retry_wait_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of seconds to wait before retrying a request.",
			},
		},
	}
}

// MetaSchema returns the provider's metadata schema, which defines additional metadata attributes.
func (p *FrameworkProvider) MetaSchema(
	_ context.Context,
	_ provider.MetaSchemaRequest,
	resp *provider.MetaSchemaResponse,
) {
	resp.Schema = metaschema.Schema{
		Attributes: map[string]metaschema.Attribute{
			"module_name": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Resources returns a list of resource constructors that the provider supports.
func (p *FrameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	resources := []func() resource.Resource{}
	resources = append(resources, services.FabricResources()...)
	resources = append(resources, services.MetalResources()...)
	resources = append(resources, services.NetworkEdgeResources()...)

	return resources
}

// DataSources returns a list of data source constructors that the provider supports.
func (p *FrameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	datasources := []func() datasource.DataSource{}
	datasources = append(datasources, services.FabricDatasources()...)
	datasources = append(datasources, services.MetalDatasources()...)
	datasources = append(datasources, services.NetworkEdgeDatasources()...)

	return datasources
}
