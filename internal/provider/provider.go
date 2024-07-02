package provider

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	metalconnection "github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	metalgateway "github.com/equinix/terraform-provider-equinix/internal/resources/metal/gateway"
	metalorganization "github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	metalorganizationmember "github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization_member"
	metalproject "github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	metalprojectsshkey "github.com/equinix/terraform-provider-equinix/internal/resources/metal/project_ssh_key"
	metalsshkey "github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type FrameworkProvider struct {
	ProviderVersion string
	Meta            *config.Config
}

func CreateFrameworkProvider(version string) provider.ProviderWithMetaSchema {
	return &FrameworkProvider{
		ProviderVersion: version,
	}
}

func (p *FrameworkProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "equinixcloud"
}

func (p *FrameworkProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
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

func (p *FrameworkProvider) MetaSchema(
	ctx context.Context,
	req provider.MetaSchemaRequest,
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

func (p *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		metalgateway.NewResource,
		metalproject.NewResource,
		metalprojectsshkey.NewResource,
		metalsshkey.NewResource,
		metalconnection.NewResource,
		metalorganization.NewResource,
		metalorganizationmember.NewResource,
		vlan.NewResource,
	}
}

func (p *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		metalgateway.NewDataSource,
		metalproject.NewDataSource,
		metalprojectsshkey.NewDataSource,
		metalconnection.NewDataSource,
		metalorganization.NewDataSource,
		vlan.NewDataSource,
	}
}
