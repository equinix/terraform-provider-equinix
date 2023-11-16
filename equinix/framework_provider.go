package equinix

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/equinix/helper"
	"github.com/equinix/terraform-provider-equinix/equinix/metal_ssh_key"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type FrameworkProvider struct {
	ProviderVersion string
	Meta            *helper.FrameworkProviderMeta
}

func CreateFrameworkProvider(version string) provider.Provider {
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
				Description: "The Equinix API base URL to point out desired environment. Defaults to " + DefaultBaseURL,
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "API Consumer Key available under My Apps section in developer portal",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "API Consumer secret available under My Apps section in developer portal",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "API token from the developer sandbox",
			},
			"auth_token": schema.StringAttribute{
				Optional:    true,
				Description: "The Equinix Metal API auth key for API operations",
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "The duration of time, in seconds, that the Equinix Platform API Client should wait before canceling an API request. Defaults to " + fmt.Sprint(DefaultTimeout),
			},
			"response_max_page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum number of records in a single response for REST queries that produce paginated responses",
			},
			"max_retries": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of retries.",
			},
			"max_retry_wait_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of seconds to wait before retrying a request.",
			},
		},
	}
}

func (p *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		metal_ssh_key.NewResource,
	}
}

func (p *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// return nil
	return []func() datasource.DataSource{
		metal_ssh_key.NewDataSource,
	}
}
