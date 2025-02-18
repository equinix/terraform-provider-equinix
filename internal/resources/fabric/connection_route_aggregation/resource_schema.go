package connection_route_aggregation

import (
	"context"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"route_aggregation_id": schema.StringAttribute{
				Description: "UUID of the Route Aggregation to apply this Rule to",
				Required:    true,
			},
			"connection_id": schema.StringAttribute{
				Description: "Equinix Assigned UUID of the Equinix Connection to attach the Route Aggregation Policy to",
				Required:    true,
			},
			"href": schema.StringAttribute{
				Description: "URI to the attached Route Aggregation Policy on the Connection",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Route Aggregation Type. One of [\"BGP_IPv4_PREFIX_AGGREGATION\"]",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix Assigned ID for Route Aggregation Policy",
				Computed:    true,
			},
			"attachment_status": schema.StringAttribute{
				Description: "Status of the Route Aggregation Policy attachment lifecycle",
				Computed:    true,
			},
		},
	}
}
