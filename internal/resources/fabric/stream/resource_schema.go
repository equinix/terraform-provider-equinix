package stream

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Stream

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Streams`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"type": schema.StringAttribute{
				Description: "Equinix defined Streaming Type",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Customer-provided name of the stream resource",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Customer-provided description of the stream resource",
				Required:    true,
			},
			"project": schema.SingleNestedAttribute{
				Description: "Equinix Project attribute object",
				Optional:    true,
				Computed:    true,
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
				Description: "Equinix auto generated URI to the stream resource in Equinix Portal",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix-assigned unique id for the stream resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "Value representing provisioning status for the stream resource",
				Computed:    true,
			},
			"assets_count": schema.Int32Attribute{
				Description: "Count of the streaming assets attached to the stream resource",
				Computed:    true,
			},
			"stream_subscriptions_count": schema.Int32Attribute{
				Description: "Count of the client subscriptions on the stream resource",
				Computed:    true,
			},
			"change_log": schema.SingleNestedAttribute{
				Description: "Details of the last change on the stream resource",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
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
