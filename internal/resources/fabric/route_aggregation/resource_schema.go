package route_aggregation

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
			"type": schema.StringAttribute{
				Description: "Equinix defined Route Aggregation Type; BGP_IPv4_PREFIX_AGGREGATION, BGP_IPv6_PREFIX_AGGREGATION",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Customer provided name of the route aggregation",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Customer-provided route aggregation description",
				Optional:    true,
			},
			"project": schema.SingleNestedAttribute{
				Description: "Equinix Project attribute object",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				CustomType: fwtypes.NewObjectTypeOf[ProjectModel](ctx),
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Description: "Equinix Subscriber-assigned project ID",
						Required:    true,
					},
				},
			},
			"href": schema.StringAttribute{
				Description: "Equinix auto generated URI to the route aggregation resource",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix-assigned unique id for the route aggregation resource",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "Value representing provisioning status for the route aggregation resource",
				Computed:    true,
			},
			"change": schema.SingleNestedAttribute{
				Description: "Current state of latest Route Aggregation change",
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				CustomType: fwtypes.NewObjectTypeOf[ChangeModel](ctx),
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Description: "Equinix-assigned unique id for a change",
						Required:    true,
					},
					"type": schema.StringAttribute{
						Description: "Equinix defined Route Aggregation Change Type",
						Required:    true,
					},
					"href": schema.StringAttribute{
						Description: "Equinix auto generated URI to the route aggregation change",
						Required:    true,
					},
				},
			},
			"connections_count": schema.Int32Attribute{
				Description: "Number of Connections attached to route aggregation",
				Computed:    true,
			},
			"rules_count": schema.Int32Attribute{
				Description: "Number of Rules attached to route aggregation",
				Computed:    true,
			},
			"change_log": schema.SingleNestedAttribute{
				Description: "Details of the last change on the stream resource",
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				CustomType: fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
				Attributes: map[string]schema.Attribute{
					"created_by": schema.StringAttribute{
						Description: "User name of creator of the stream resource",
						Computed:    true,
					},
					"created_by_full_name": schema.StringAttribute{
						Description: "Legal name of creator of the stream resource",
						Computed:    true,
					},
					"created_by_email": schema.StringAttribute{
						Description: "Email of creator of the stream resource",
						Computed:    true,
					},
					"created_date_time": schema.StringAttribute{
						Description: "Creation time of the stream resource",
						Computed:    true,
					},
					"updated_by": schema.StringAttribute{
						Description: "User name of last updater of the stream resource",
						Computed:    true,
					},
					"updated_by_full_name": schema.StringAttribute{
						Description: "Legal name of last updater of the stream resource",
						Computed:    true,
					},
					"updated_by_email": schema.StringAttribute{
						Description: "Email of last updater of the stream resource",
						Computed:    true,
					},
					"updated_date_time": schema.StringAttribute{
						Description: "Last update time of the stream resource",
						Computed:    true,
					},
					"deleted_by": schema.StringAttribute{
						Description: "User name of deleter of the stream resource",
						Computed:    true,
					},
					"deleted_by_full_name": schema.StringAttribute{
						Description: "Legal name of deleter of the stream resource",
						Computed:    true,
					},
					"deleted_by_email": schema.StringAttribute{
						Description: "Email of deleter of the stream resource",
						Computed:    true,
					},
					"deleted_date_time": schema.StringAttribute{
						Description: "Deletion time of the stream resource",
						Computed:    true,
					},
				},
			},
		},
	}
}
