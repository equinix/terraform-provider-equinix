package connectionrouteaggregation

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func dataSourceAllConnectionRouteAggregationSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Connection Route Aggregations with pagination details
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Route-Aggregations`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"connection_id": schema.StringAttribute{
				Description: "The uuid of the connection this data source should retrieve",
				Required:    true,
			},
			"data": schema.ListNestedAttribute{
				Description: "Returned list of connection route aggregation objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[baseConnectionRouteAggregationModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getConnectionRouteAggregationSchema(ctx),
				},
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned connection route aggregations list",
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
						Description: "The total number of connection route aggregations available to the user making the request",
						Optional:    true,
						Computed:    true,
					},
					"next": schema.StringAttribute{
						Description: "The URL relative to the next item in the response",
						Optional:    true,
						Computed:    true,
					},
					"previous": schema.StringAttribute{
						Description: "The URL relative to the previous item in the response",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

func dataSourceSingleConnectionRouteAggregationSchema(ctx context.Context) schema.Schema {
	baseConnectionRouteAggregationSchema := getConnectionRouteAggregationSchema(ctx)
	baseConnectionRouteAggregationSchema["id"] = framework.IDAttributeDefaultDescription()
	baseConnectionRouteAggregationSchema["route_aggregation_id"] = schema.StringAttribute{
		Description: "The uuid of the route aggregation this data source should retrieve",
		Required:    true,
	}
	baseConnectionRouteAggregationSchema["connection_id"] = schema.StringAttribute{
		Description: "The uuid of the connection this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Connection Route Aggregation by UUID
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Route-Aggregations`,
		Attributes: baseConnectionRouteAggregationSchema,
	}
}

func getConnectionRouteAggregationSchema(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"route_aggregation_id": schema.StringAttribute{
			Description: "UUID of the Route Aggregation to attach this Connection to",
			Required:    true,
		},
		"connection_id": schema.StringAttribute{
			Description: "UUID of the Connection to attach this Route Aggregation to",
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
	}
}
