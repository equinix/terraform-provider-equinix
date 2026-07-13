package routeaggregationrule

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/equinix/terraform-provider-equinix/internal/slice"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceSearchRouteAggregationRulesSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Route Aggregation Rules with pagination details
Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Route-Aggregations`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"route_aggregation_id": schema.StringAttribute{
				Description: "The UUID of the route aggregation from which this data source retrieves its rules.",
				Required:    true,
			},
			"data": schema.ListNestedAttribute{
				Description: "Returned list of Route Aggregation Rule objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[baseRouteAggregationRuleModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getRouteAggregationRuleSchema(ctx),
				},
			},
			"filter": schema.ListNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Optional:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[filterModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property": schema.StringAttribute{
							Required:    true,
							Description: "Possible field names to use on filters. One of [ /type, /name, /uuid, /state, /prefix]",
						},
						"operator": schema.StringAttribute{
							Required:    true,
							Description: "Operators to use on your filtered field with the values given. One of [ =, !=, LIKE, NOT LIKE, IN, NOT IN, ILIKE]",
						},
						"values": schema.ListAttribute{
							Required:    true,
							Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
							CustomType:  fwtypes.ListOfStringType,
							ElementType: types.StringType,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(8),
				},
			},
			"outer_operator": schema.StringAttribute{
				Description: "Determines if the filter list will be grouped by AND or by OR. One of [AND, OR]",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AND", "OR")},
			},
			"sort": schema.ListNestedAttribute{
				Description: "Sort criteria for the Data Source Search Request",
				Optional:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[sortModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							Description: fmt.Sprintf("The sorting direction. Can be one of: %v, Defaults to DESC", fabricv4.AllowedRouteAggregationRuleSortDirectionEnumValues),
							Optional:    true,
							Validators: []validator.String{stringvalidator.OneOf(
								slice.Map(fabricv4.AllowedRouteAggregationRuleSortDirectionEnumValues, func(r fabricv4.RouteAggregationRuleSortDirection) string { return string(r) })...,
							)},
						},
						"property": schema.StringAttribute{
							Description: fmt.Sprintf("The property name to use in sorting. One of %v. Defaults to /changeLog/updatedDateTime", fabricv4.AllowedRouteFilterRuleSortByEnumValues),
							Optional:    true,
							Validators: []validator.String{stringvalidator.OneOf(
								slice.Map(fabricv4.AllowedRouteAggregationRuleSortByEnumValues, func(r fabricv4.RouteAggregationRuleSortBy) string {
									return string(r)
								})...,
							)},
						},
					},
				},
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned Route Aggregation Rules list",
				Optional:    true,
				CustomType:  fwtypes.NewObjectTypeOf[paginationModel](ctx),
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
		Description: "The UUID of the Route Aggregation Rule this data source should retrieve",
		Required:    true,
	}
	baseRouteAggregationRuleSchema["route_aggregation_id"] = schema.StringAttribute{
		Description: "The UUID of the route aggregation this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Route Aggregation Rule by UUID
Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Route-Aggregations`,
		Attributes: baseRouteAggregationRuleSchema,
	}

}

func getRouteAggregationRuleSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"route_aggregation_id": schema.StringAttribute{
			Description: "UUID of the Route Aggregation that the rule is applied to",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the Route Aggregation Rule resource",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Equinix defined Route Aggregation Type; BGP_IPv4_PREFIX_AGGREGATION, BGP_IPv6_PREFIX_AGGREGATION",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique id for the Route Aggregation Rule resource",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Customer provided name of the Route Aggregation Rule",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Customer-provided Route Aggregation Rule description",
			Optional:    true,
		},
		"state": schema.StringAttribute{
			Description: "Value representing provisioning status for the Route Aggregation Rule resource",
			Computed:    true,
		},
		"prefix": schema.StringAttribute{
			Description: "Customer-provided Route Aggregation Rule prefix",
			Computed:    true,
		},
		"change": schema.SingleNestedAttribute{
			Description: "Current state of latest Route Aggregation Rule change",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[changeModel](ctx),
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
	}
}
