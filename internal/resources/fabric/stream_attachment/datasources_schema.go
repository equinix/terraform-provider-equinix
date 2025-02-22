package streamattachment

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func dataSourceAllStreamAttachmentsSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Stream Attached Assets with filters and pagination details

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Streams`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"filters": schema.ListNestedAttribute{
				Description: "List of filters to apply to the stream attachment get request. Maximum of 8. All will be AND'd together with 1 of the 8 being a possible OR group of 3",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[FilterModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property": schema.StringAttribute{
							Description: "Property to apply the filter to",
							Required:    true,
						},
						"operator": schema.StringAttribute{
							Description: "Operation applied to the values of the filter",
							Required:    true,
						},
						"values": schema.ListAttribute{
							Description: "List of values to apply the operation to for the specified property",
							Required:    true,
							ElementType: types.StringType,
						},
						"or": schema.BoolAttribute{
							Description: "Boolean value to specify if this filter is a part of the OR group. Has a maximum of 3 and only counts for 1 of the 8 possible filters",
							Optional:    true,
						},
					},
				},
			},
			"sort": schema.ListNestedAttribute{
				Description: "The list of sort criteria for the stream assets search request",
				Optional:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[SortModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							Description: "The sorting direction of the property chosen. ASC or DESC",
							Required:    true,
						},
						"property": schema.StringAttribute{
							Description: "The field name the sorting is performed on",
							Required:    true,
						},
					},
				},
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned streams list",
				Required:    true,
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
			"data": schema.ListNestedAttribute{
				Description: "Returned list of stream objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BaseAssetModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getStreamSchema(ctx),
				},
			},
		},
	}
}

func dataSourceByIDsSchema(ctx context.Context) schema.Schema {
	baseStreamSchema := getStreamSchema(ctx)
	baseStreamSchema["id"] = framework.IDAttributeDefaultDescription()
	baseStreamSchema["stream_id"] = schema.StringAttribute{
		Description: "The uuid of the stream this data source should retrieve",
		Required:    true,
	}
	baseStreamSchema["asset"] = schema.StringAttribute{
		Description: "Equinix defined asset category. Matches the product name the asset is a part of",
		Required:    true,
	}
	baseStreamSchema["asset_id"] = schema.StringAttribute{
		Description: "Equinix defined UUID of the asset being attached to the stream",
		Required:    true,
	}

	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Stream Asset Attachment by IDs

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Streams`,
		Attributes: baseStreamSchema,
	}
}

func getStreamSchema(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"metrics_enabled": schema.BoolAttribute{
			Description: "Boolean value indicating enablement of metrics for this asset stream attachment",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Equinix defined type for the asset stream attachment",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the stream attachment in Equinix Portal",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique id for the stream attachment",
			Computed:    true,
		},
		"attachment_status": schema.StringAttribute{
			Description: "Value representing status for the stream attachment",
			Computed:    true,
		},
	}
}
