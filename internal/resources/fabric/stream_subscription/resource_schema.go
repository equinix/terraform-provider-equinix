package stream_subscription

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Stream Subscription

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Stream-Subscriptions`,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "Type of the stream subscription request",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Customer-provided stream subscription name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Customer-provided stream subscription description",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Stream subscription enabled status",
				Required:    true,
			},
			"filters": schema.ListNestedAttribute{
				Description: "List of filters to apply to the stream subscription selectors. Maximum of 8. All will be AND'd together with 1 of the 8 being a possible OR group of 3",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property": schema.StringAttribute{
							Description: "Property to apply the filter to",
							Required: true,
						},
						"operator": schema.StringAttribute{
							Description: "Operation applied to the values of the filter",
							Required: true,
						}
						"values": schema.ListAttribute{
							Description: "List of values to apply the operation to for the specified property",
							Required: true,
							ElementType: types.StringType,
						},
						"or": schema.BoolAttribute{
							Description: "Boolean value to specify if this filter is a part of the OR group. Has a maximum of 3 and only counts for 1 of the 8 possible filters",
							Optional: true,
						},
					},
				},
			},
			"metric_selector": schema.SingleNestedAttribute{
				Description: "Lists of metrics to be included/excluded on the stream subscription",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"include": schema.ListAttribute{
						Description: "List of metrics to include",
						ElementType: types.StringType,
						Required:    true,
					},
					"except": schema.ListAttribute{
						Description: "List of metrics to exclude",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			"event_selector": schema.SingleNestedAttribute{
				Description: "Lists of events to be included/excluded on the stream subscription",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"include": schema.ListAttribute{
						Description: "List of events to include",
						ElementType: types.StringType,
						Required:    true,
					},
					"except": schema.ListAttribute{
						Description: "List of events to exclude",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			"sink": schema.SingleNestedAttribute{
				
			},
		},
	}
}
