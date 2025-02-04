package route_aggregation

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAllRouteAggregationsSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Streams with pagination details
Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Streams`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"data": schema.ListNestedAttribute{
				Description: "Returned list of stream objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BaseRouteAggregationModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getRouteAggregationSchema(ctx),
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Required:    true,
				//CustomType:  fwtypes.NewObjectTypeOf[FilterModel](ctx),
				Attributes: map[string]schema.Attribute{
					"property": schema.StringAttribute{
						Description: fmt.Sprintf("possible field names to use on filters. One of %v", fabricv4.AllowedRouteFiltersSearchFilterItemPropertyEnumValues),
						Required:    true,
					},
					"operator": schema.StringAttribute{
						Description: "Operators to use on your filtered field with the values given. One of [ =, !=, >, >=, <, <=, BETWEEN, NOT BETWEEN, LIKE, NOT LIKE, IN, NOT IN, IS NOT NULL, IS NULL]",
						Required:    true,
					},
					"values": schema.ListAttribute{
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
			//"filter": schema.ObjectAttribute{
			//	Description: "Filters for the Data Source Search Request",
			//	Required:    true,
			//	AttributeTypes: map[string]attr.Type{
			//		"property": schema.String{
			//			Description: "The property to be used in the filter condition (e.g., 'status', 'type')",
			//			Required:    true,
			//		},
			//		"operator": schema.StringAttribute{
			//			Description: "The operator to be used in the filter condition (e.g., '=', '>', 'IN')",
			//			Required:    true,
			//		},
			//		"values": schema.ListAttribute{
			//			Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
			//			ElementType: types.StringType,
			//			Required:    true,
			//		},
			//	},
			//}
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned route aggregations list",
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
						Description: "The total number of streams available to the user making the request",
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
			"sort": schema.SingleNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Optional:    true,
				CustomType:  fwtypes.NewObjectTypeOf[SortModel](ctx),
				Attributes: map[string]schema.Attribute{
					"direction": schema.StringAttribute{
						Description: "The sorting direction. Can be one of: [DESC, ASC], Defaults to DESC",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("DESC", "ASC"),
						},
					},
					"property": schema.StringAttribute{
						Description: fmt.Sprintf("The property name to use in sorting. One of %v Defaults to /name", fabricv4.AllowedRouteFiltersSearchFilterItemPropertyEnumValues),
						Optional:    true,
					},
				},
			},
		},
	}
}
func dataSourceSingleRouteAggregationSchema(ctx context.Context) schema.Schema {
	baseRouteAggregationSchema := getRouteAggregationSchema(ctx)
	baseRouteAggregationSchema["id"] = framework.IDAttributeDefaultDescription()
	baseRouteAggregationSchema["route_aggregation_id"] = schema.StringAttribute{
		Description: "The uuid of the route aggregation this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Stream by UUID
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Route-Aggregations`,
		Attributes: baseRouteAggregationSchema,
	}

}

func getRouteAggregationSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "Equinix defined Route Aggregation Type; BGP_IPv4_PREFIX_AGGREGATION, BGP_IPv6_PREFIX_AGGREGATION",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Customer provided name of the route aggregation",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Customer-provided route aggregation description",
			Computed:    true,
		},
		"project": schema.SingleNestedAttribute{
			Description: "Equinix Project attribute object",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ProjectModel](ctx),
			Attributes: map[string]schema.Attribute{
				"project_id": schema.StringAttribute{
					Description: "Equinix Subscriber-assigned project ID",
					Computed:    true,
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
			CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
			Attributes: map[string]schema.Attribute{
				"uuid": schema.StringAttribute{
					Description: "Equinix-assigned unique id for a change",
					Computed:    true,
				},
				"type": schema.StringAttribute{
					Description: "Equinix defined Route Aggregation Change Type",
					Computed:    true,
				},
				"href": schema.StringAttribute{
					Description: "Equinix auto generated URI to the route aggregation change",
					Computed:    true,
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
