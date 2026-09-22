package ipblock

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAllIpBlocksSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows users to search for Equinix Fabric IP Blocks with pagination details.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/IP-Blocks`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"data": schema.ListNestedAttribute{
				Description: "Returned list of IP block objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BaseIpBlockModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getIpBlockSchema(ctx),
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"property": schema.StringAttribute{
						Description: "Field name to use on filters",
						Required:    true,
					},
					"operator": schema.StringAttribute{
						Description: "Operator to use on your filtered field with the values given. One of: [=, !=, >, >=, <, <=, IN, NOT IN, LIKE, NOT LIKE, IS NOT NULL, IS NULL]",
						Required:    true,
					},
					"values": schema.ListAttribute{
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned IP blocks list",
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
						Description: "The total number of IP blocks available to the user making the request",
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
				Description: "Sort criteria for the Data Source Search Request",
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
						Description: "The property name to use in sorting. One of \"/uuid\" \"/state\" \"/prefixLength\" \"/changeLog/createdDateTime\" \"/changeLog/updatedDateTime\". Defaults to \"/changeLog/updatedDateTime\"",
						Optional:    true,
					},
				},
			},
		},
	}
}

func dataSourceSingleIpBlockSchema(ctx context.Context) schema.Schema {
	baseSchema := getIpBlockSchema(ctx)
	baseSchema["id"] = framework.IDAttributeDefaultDescription()
	baseSchema["ip_block_id"] = schema.StringAttribute{
		Description: "The UUID of the IP block this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows users to fetch an Equinix Fabric IP Block by UUID.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/IP-Blocks`,
		Attributes: baseSchema,
	}
}

func getIpBlockSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique identifier for the IP block",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the IP block resource",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Type of the IP block product",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "Provisioning state of the IP block",
			Computed:    true,
		},
		"ownership": schema.StringAttribute{
			Description: "Ownership of the IP block. One of: EQUINIX, CUSTOMER",
			Computed:    true,
		},
		"prefix_length": schema.Int32Attribute{
			Description: "Prefix length of the IP block (e.g. 24 for a /24 block)",
			Computed:    true,
		},
		"prefix": schema.StringAttribute{
			Description: "The assigned IP prefix (e.g. 192.0.2.0/24)",
			Computed:    true,
		},
		"location": schema.SingleNestedAttribute{
			Description: "Location of the IP block",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[LocationModel](ctx),
			Attributes: map[string]schema.Attribute{
				"metro_href": schema.StringAttribute{
					Description: "URI to the metro resource",
					Computed:    true,
				},
				"metro_code": schema.StringAttribute{
					Description: "Two-letter metro code",
					Computed:    true,
				},
			},
		},
		"order": schema.SingleNestedAttribute{
			Description: "Order details for the IP block",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[OrderModel](ctx),
			Attributes: map[string]schema.Attribute{
				"href": schema.StringAttribute{
					Description: "URI to the order resource",
					Computed:    true,
				},
				"order_number": schema.StringAttribute{
					Description: "Order number",
					Computed:    true,
				},
			},
		},
		"account": schema.SingleNestedAttribute{
			Description: "Account associated with the IP block",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[AccountModel](ctx),
			Attributes: map[string]schema.Attribute{
				"account_number": schema.StringAttribute{
					Description: "Account number",
					Computed:    true,
				},
			},
		},
		"project": schema.SingleNestedAttribute{
			Description: "Project associated with the IP block",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ProjectModel](ctx),
			Attributes: map[string]schema.Attribute{
				"href": schema.StringAttribute{
					Description: "URI to the project resource",
					Computed:    true,
				},
				"project_id": schema.StringAttribute{
					Description: "Equinix subscriber-assigned project ID",
					Computed:    true,
				},
			},
		},
		"assets": schema.ListNestedAttribute{
			Description: "List of assets (connections or services) using this IP block",
			Computed:    true,
			CustomType:  fwtypes.NewListNestedObjectTypeOf[AssetModel](ctx),
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Type of the asset",
						Computed:    true,
					},
					"uuid": schema.StringAttribute{
						Description: "Unique identifier of the asset",
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: "URI to the asset resource",
						Computed:    true,
					},
				},
			},
		},
		"change": schema.SingleNestedAttribute{
			Description: "Current state of the latest change on the IP block",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
			Attributes: map[string]schema.Attribute{
				"href": schema.StringAttribute{
					Description: "URI to the change resource",
					Computed:    true,
				},
			},
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Audit log of changes to the IP block",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
			Attributes: map[string]schema.Attribute{
				"created_date_time": schema.StringAttribute{
					Description: "Creation time of the IP block",
					Computed:    true,
				},
				"updated_date_time": schema.StringAttribute{
					Description: "Last update time of the IP block",
					Computed:    true,
				},
				"deleted_date_time": schema.StringAttribute{
					Description: "Deletion time of the IP block",
					Computed:    true,
				},
			},
		},
	}
}
