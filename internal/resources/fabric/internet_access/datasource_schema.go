package internetaccess

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceSingleInternetAccessServiceSchema(ctx context.Context) schema.Schema {
	attrs := getInternetAccessServiceSchema(ctx)
	attrs["internet_access_service_id"] = schema.StringAttribute{
		Description: "UUID of the Internet Access service to retrieve",
		Required:    true,
	}
	attrs["id"] = framework.IDAttributeDefaultDescription()
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows retrieval of an Equinix Internet Access (EIA) service by UUID.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Internet-Access-Services`,
		Attributes: attrs,
	}
}

func dataSourceAllInternetAccessServicesSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows users to search for Equinix Internet Access (EIA) services with pagination details.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Internet-Access-Services`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"data": schema.ListNestedAttribute{
				Description: "Returned list of Internet Access service objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BaseInternetAccessServiceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getInternetAccessServiceSchema(ctx),
				},
			},
			"filter": schema.ListNestedAttribute{
				Description: "List of filter conditions (AND logic) for the Data Source Search Request",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property": schema.StringAttribute{
							Description: "Field name to use on filters",
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
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned Internet Access services list",
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
						Description: "The total number of Internet Access services available to the user making the request",
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
						Description: "The property name to use in sorting. One of \"/name\" \"/uuid\" \"/state\" \"/changeLog/createdDateTime\" \"/changeLog/updatedDateTime\". Defaults to \"/changeLog/updatedDateTime\"",
						Optional:    true,
					},
				},
			},
		},
	}
}

func getInternetAccessServiceSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the Internet Access service resource",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Type of the Internet Access service",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique identifier for the Internet Access service",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Name of the Internet Access service",
			Computed:    true,
		},
		"bandwidth": schema.Int32Attribute{
			Description: "Bandwidth of the service in Mbps",
			Computed:    true,
		},
		"bandwidth_commit": schema.Int32Attribute{
			Description: "Minimum bandwidth commit for burst billing variant of the service",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "Provisioning state of the Internet Access service",
			Computed:    true,
		},
		"use_case": schema.StringAttribute{
			Description: "Use case of the Internet Access service. One of: MAIN, MANAGEMENT_ACCESS",
			Computed:    true,
		},
		"change": schema.SingleNestedAttribute{
			Description: "Current state of the latest change on the Internet Access service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
			Attributes: map[string]schema.Attribute{
				"href": schema.StringAttribute{
					Description: "Equinix auto generated URI to the change resource",
					Computed:    true,
				},
			},
		},
		"locations": schema.ListNestedAttribute{
			Description: "List of locations associated with the service",
			Computed:    true,
			CustomType:  fwtypes.NewListNestedObjectTypeOf[LocationModel](ctx),
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"metro_href": schema.StringAttribute{
						Description: "URI to the metro resource",
						Computed:    true,
					},
					"metro_code": schema.StringAttribute{
						Description: "Two-letter metro code",
						Computed:    true,
					},
					"region": schema.StringAttribute{
						Description: "Geographic region of the metro",
						Computed:    true,
					},
					"ibx": schema.StringAttribute{
						Description: "IBX data center code",
						Computed:    true,
					},
				},
			},
		},
		"routing_protocol": schema.SingleNestedAttribute{
			Description: "Routing protocol configuration for the Internet Access service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[RoutingProtocolModel](ctx),
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "Routing protocol type. One of: BGP, DIRECT, STATIC",
					Computed:    true,
				},
				"customer_routes": schema.ListNestedAttribute{
					Description: "List of customer routes (IP block allocations) attached to the service",
					Computed:    true,
					CustomType:  fwtypes.NewListNestedObjectTypeOf[CustomerRouteModel](ctx),
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"ip_block": schema.SingleNestedAttribute{
								Description: "IP block reference",
								Computed:    true,
								CustomType:  fwtypes.NewObjectTypeOf[IpBlockRefModel](ctx),
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{
										Description: "URI to the IP block resource",
										Computed:    true,
									},
									"uuid": schema.StringAttribute{
										Description: "Unique identifier of the IP block",
										Computed:    true,
									},
								},
							},
						},
					},
				},
				"connections": schema.ListNestedAttribute{
					Description: "List of Fabric connections associated with the service",
					Computed:    true,
					CustomType:  fwtypes.NewListNestedObjectTypeOf[ConnectionRefModel](ctx),
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"href": schema.StringAttribute{
								Description: "URI to the connection resource",
								Computed:    true,
							},
							"uuid": schema.StringAttribute{
								Description: "Unique identifier of the connection",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		"billing": schema.SingleNestedAttribute{
			Description: "Billing details for the Internet Access service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[BillingModel](ctx),
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "Billing type of the service",
					Computed:    true,
				},
				"enabled": schema.BoolAttribute{
					Description: "Whether billing is enabled for the service",
					Computed:    true,
				},
				"start_date": schema.StringAttribute{
					Description: "Billing start date in RFC3339 format",
					Computed:    true,
				},
			},
		},
		"account": schema.SingleNestedAttribute{
			Description: "Account details for the Internet Access service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[AccountModel](ctx),
			Attributes: map[string]schema.Attribute{
				"account_number": schema.StringAttribute{
					Description: "Account number",
					Computed:    true,
				},
				"href": schema.StringAttribute{
					Description: "URI to the account resource",
					Computed:    true,
				},
			},
		},
		"project": schema.SingleNestedAttribute{
			Description: "Project associated with the Internet Access service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ProjectModel](ctx),
			Attributes: map[string]schema.Attribute{
				"project_id": schema.StringAttribute{
					Description: "Equinix subscriber-assigned project ID",
					Computed:    true,
				},
			},
		},
		"order": schema.SingleNestedAttribute{
			Description: "Order details for the Internet Access service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[OrderModel](ctx),
			Attributes: map[string]schema.Attribute{
				"purchase_order_number": schema.StringAttribute{
					Description: "Purchase order number",
					Computed:    true,
				},
				"customer_reference_number": schema.StringAttribute{
					Description: "Customer reference number",
					Computed:    true,
				},
				"billing_tier": schema.StringAttribute{
					Description: "Billing tier",
					Computed:    true,
				},
				"order_id": schema.StringAttribute{
					Description: "Order ID",
					Computed:    true,
				},
				"order_number": schema.StringAttribute{
					Description: "Order number",
					Computed:    true,
				},
				"term_length": schema.Int32Attribute{
					Description: "Term length in months",
					Computed:    true,
				},
				"contracted_bandwidth": schema.Int32Attribute{
					Description: "Contracted bandwidth in Mbps",
					Computed:    true,
				},
				"href": schema.StringAttribute{
					Description: "URI to the order resource",
					Computed:    true,
				},
			},
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Details of the last change on the Internet Access service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
			Attributes: map[string]schema.Attribute{
				"created_by": schema.StringAttribute{
					Description: "User name of the creator",
					Computed:    true,
				},
				"created_by_full_name": schema.StringAttribute{
					Description: "Full name of the creator",
					Computed:    true,
				},
				"created_by_email": schema.StringAttribute{
					Description: "Email of the creator",
					Computed:    true,
				},
				"created_date_time": schema.StringAttribute{
					Description: "Creation time",
					Computed:    true,
				},
				"updated_by": schema.StringAttribute{
					Description: "User name of the last updater",
					Computed:    true,
				},
				"updated_by_full_name": schema.StringAttribute{
					Description: "Full name of the last updater",
					Computed:    true,
				},
				"updated_by_email": schema.StringAttribute{
					Description: "Email of the last updater",
					Computed:    true,
				},
				"updated_date_time": schema.StringAttribute{
					Description: "Last update time",
					Computed:    true,
				},
				"deleted_by": schema.StringAttribute{
					Description: "User name of the deleter",
					Computed:    true,
				},
				"deleted_by_full_name": schema.StringAttribute{
					Description: "Full name of the deleter",
					Computed:    true,
				},
				"deleted_by_email": schema.StringAttribute{
					Description: "Email of the deleter",
					Computed:    true,
				},
				"deleted_date_time": schema.StringAttribute{
					Description: "Deletion time",
					Computed:    true,
				},
			},
		},
	}
}
