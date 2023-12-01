package internal

import (
	"context"
	"fmt"
	"regexp"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/metal_bgp_session"
	"github.com/equinix/terraform-provider-equinix/internal/metal_connection"
	"github.com/equinix/terraform-provider-equinix/internal/metal_gateway"
	"github.com/equinix/terraform-provider-equinix/internal/metal_ip_attachment"
	"github.com/equinix/terraform-provider-equinix/internal/metal_organization"
	"github.com/equinix/terraform-provider-equinix/internal/metal_organization_member"
	"github.com/equinix/terraform-provider-equinix/internal/metal_port"
	"github.com/equinix/terraform-provider-equinix/internal/metal_project"
	"github.com/equinix/terraform-provider-equinix/internal/metal_reserved_ip_block"
	"github.com/equinix/terraform-provider-equinix/internal/metal_ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/metal_vlan"
	"github.com/equinix/terraform-provider-equinix/internal/metal_vrf"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var urlRE = regexp.MustCompile(`^https?://(?:www\.)?[a-zA-Z0-9./]+$`)

type FrameworkProvider struct {
	ProviderVersion string
	Meta            *config.Config
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
				Description: "The Equinix API base URL to point out desired environment. Defaults to " + config.DefaultBaseURL,
				Validators: []validator.String{
					stringvalidator.RegexMatches(urlRE, "must be a valid URL with http or https schema"),
				},
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
				Description:  fmt.Sprintf("The duration of time, in seconds, that the Equinix Platform API Client should wait before canceling an API request.  Defaults to %d", config.DefaultTimeout),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"response_max_page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum number of records in a single response for REST queries that produce paginated responses",
				Validators: []validator.Int64{
					int64validator.AtLeast(100),
				},
			},
			"max_retries": schema.Int64Attribute{
				Optional:    true,
				// Description: "Maximum number of retries.",
			},
			"max_retry_wait_seconds": schema.Int64Attribute{
				Optional:    true,
				// Description: "Maximum number of seconds to wait before retrying a request.",
			},
		},
	}
}

func (p *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		metal_bgp_session.NewResource,
		metal_ssh_key.NewResource,
		metal_project.NewResource,
		metal_organization.NewResource,
		metal_organization_member.NewResource,
		metal_port.NewResource,
		metal_vlan.NewResource,
		metal_connection.NewResource,
		func() resource.Resource {
            return metal_gateway.NewResource(ctx)
        },
		metal_ip_attachment.NewResource,
		func() resource.Resource {
            return metal_reserved_ip_block.NewResource(ctx)
        },
		metal_vrf.NewResource,
	}
}

func (p *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// return nil
	return []func() datasource.DataSource{
		// metal_ssh_key.NewDataSource,
	}
}
