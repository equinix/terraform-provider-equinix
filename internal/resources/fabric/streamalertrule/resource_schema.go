package streamalertrule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Stream Alert Rules'
}


Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Stream-Alert-Rules`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"stream_id": schema.StringAttribute{
				Description: "The stream UUID that contains this alert rule",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the stream alert rule",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Customer-provided stream alert rule name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Customer-provided stream alert rule description",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Stream alert rule enabled status",
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				Computed:    true,
			},
			"resource_selector": schema.SingleNestedAttribute{
				Description: "Resource selector for the stream alert rule",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[selectorModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"include": schema.ListAttribute{
						Description: "List of metrics to include",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
			"metric_selector": schema.SingleNestedAttribute{
				Description: "Metric selector for the stream alert rule",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[selectorModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"include": schema.ListAttribute{
						Description: "List of metrics to include",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
			"detection_method": schema.SingleNestedAttribute{
				Description: "Detection method for stream alert rule",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[selectorModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Stream Alert Rule detection method type",
						Required:    true,
					},
					"window_size": schema.StringAttribute{
						Description: "Stream alert rule metric window size",
						Optional:    true,
						Computed:    true,
					},
					"operand": schema.StringAttribute{
						Description: "Stream alert rule metric operand",
						Optional:    true,
						Computed:    true,
					},
					"warning_threshold": schema.StringAttribute{
						Description: "Stream alert rule metric warning threshold",
						Optional:    true,
						Computed:    true,
					},
					"critical_threshold": schema.StringAttribute{
						Description: "Stream alert rule metric critical threshold",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"href": schema.StringAttribute{
				Description: "Equinix assigned URI of the stream alert rule",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix assigned unique identifier for the stream alert rule",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "Value representing provisioning status for the stream alert rule",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"change_log": schema.SingleNestedAttribute{
				Description: "Details of the last change on the stream resource",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[changeLogModel](ctx),
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
