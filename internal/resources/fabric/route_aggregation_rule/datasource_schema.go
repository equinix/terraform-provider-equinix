package route_aggregation_rule

import (
	"context"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func dataSourceAllRouteAggregationRulesSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Route Aggregation Rules with pagination details
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Route-Aggregations`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"route_aggregation_id": schema.StringAttribute{
				Description: "The uuid of the route aggregation rule this data source should retrieve",
				Required:    true,
			},
			"data": schema.ListNestedAttribute{
				Description: "Returned list of route aggregation rule objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BaseRouteAggregationRuleModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getRouteAggregationRuleSchema(ctx),
				},
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned route aggregation rules list",
				Optional:    true,
				CustomType:  fwtypes.NewObjectTypeOf[PaginationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"offset": schema.Int32Attribute{
						Description: "Index of the first item returned in the response. The default is 0",
						Optional:    true,
						Computed:    true,
					},
					"limit": schema.Int32Attribute{
						Description: "Maximum number of search results returned per page. Number must be between 1 and 100, and the default is 20",
						Optional:    true,
						Computed:    true,
					},
					"total": schema.Int32Attribute{
						Description: "The total number of route agrgegation rules available to the user making the request",
						Computed:    true,
					},
					"next": schema.StringAttribute{
						Description: "The URL relative to the next item in the response",
						Computed:    true,
					},
					"previous": schema.StringAttribute{
						Description: "The URL relative to the previous item in the response",
						Computed:    true,
					},
				},
			},
		},
	}
}

func dataSourceSingleRouteAggregationRuleSchema(ctx context.Context) schema.Schema {
	baseRouteAggregationRuleSchema := getRouteAggregationRuleSchema(ctx)
	baseRouteAggregationRuleSchema["id"] = framework.IDAttributeDefaultDescription()
	baseRouteAggregationRuleSchema["route_aggregation_rule_id"] = schema.StringAttribute{
		Description: "The uuid of the route aggregation rule this data source should retrieve",
		Required:    true,
	}
	baseRouteAggregationRuleSchema["route_aggregation_id"] = schema.StringAttribute{
		Description: "The uuid of the route aggregation this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Route Aggregation Rule by UUID
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Route-Aggregations`,
		Attributes: baseRouteAggregationRuleSchema,
	}

}

func getRouteAggregationRuleSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"route_aggregation_id": schema.StringAttribute{
			Description: "UUID of the Route Aggregation to apply this Rule to",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the route aggregation rule resource",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Equinix defined Route Aggregation Type; BGP_IPv4_PREFIX_AGGREGATION, BGP_IPv6_PREFIX_AGGREGATION",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique id for the route aggregation rule resource",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Customer provided name of the route aggregation rule",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Customer-provided route aggregation rule description",
			Optional:    true,
		},
		"state": schema.StringAttribute{
			Description: "Value representing provisioning status for the route aggregation rule resource",
			Computed:    true,
		},
		"prefix": schema.StringAttribute{
			Description: "Customer-provided route aggregation rule prefix",
			Computed:    true,
		},
		"change": schema.SingleNestedAttribute{
			Description: "Current state of latest route aggregation rule change",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
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
					Computed:    true,
				},
			},
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
	}
}
