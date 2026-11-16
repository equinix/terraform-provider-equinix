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

func dataSourceAllInternetAccessSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows users to search for Equinix Internet Access (EIA) services with pagination details.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Internet-Access`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"data": schema.ListNestedAttribute{
				Description: "Returned list of Internet Access service objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BaseInternetAccessModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getInternetAccessSchema(ctx),
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"property": schema.StringAttribute{
						Description: "Field name to use on filters, e.g. \"/project/projectId\", \"/name\", \"/state\", \"/type\"",
						Required:    true,
					},
					"operator": schema.StringAttribute{
						Description: "Operator to use on your filtered field with the values given. One of: [=, !=, IN, NOT IN, LIKE, ILIKE, BETWEEN, >, >=, <, <=, ~*]",
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
						Description: "The sorting direction. Can be one of: [DESC, ASC], Defaults to ASC",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("DESC", "ASC"),
						},
					},
					"property": schema.StringAttribute{
						Description: "The property name to use in sorting, e.g. \"/name\". Defaults to \"/name\". Note: sorting by \"/changeLog/createdDateTime\" currently makes the API return a server error",
						Optional:    true,
					},
				},
			},
		},
	}
}

func dataSourceSingleInternetAccessSchema(ctx context.Context) schema.Schema {
	baseSchema := getInternetAccessSchema(ctx)
	baseSchema["id"] = framework.IDAttributeDefaultDescription()
	baseSchema["internet_access_service_id"] = schema.StringAttribute{
		Description: "The UUID of the Internet Access service this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows users to fetch an Equinix Internet Access (EIA) service by UUID.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Internet-Access`,
		Attributes: baseSchema,
	}
}

func peeringSchema(ctx context.Context, ipVersion string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: ipVersion + " peering details of the connection",
		Computed:    true,
		CustomType:  fwtypes.NewObjectTypeOf[PeeringModel](ctx),
		Attributes: map[string]schema.Attribute{
			"prefix": schema.StringAttribute{
				Description: ipVersion + " prefix for the peering connection",
				Computed:    true,
			},
			"prefix_length": schema.Int32Attribute{
				Description: "Determines the size of the subnet",
				Computed:    true,
			},
			"equinix_peer_ip": schema.StringAttribute{
				Description: "Peering IP address on the Equinix side",
				Computed:    true,
			},
			"customer_peer_ip": schema.StringAttribute{
				Description: "Peering IP address on the customer side",
				Computed:    true,
			},
			"equinix_vrrp_ip": schema.StringAttribute{
				Description: "VRRP IP address on the Equinix side",
				Computed:    true,
			},
			"customer_vrrp_ip": schema.StringAttribute{
				Description: "VRRP IP address on the customer side",
				Computed:    true,
			},
		},
	}
}

func getInternetAccessSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique identifier for the Internet Access service",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the Internet Access service",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Type of the Internet Access service. One of: SINGLE_IA, DUAL_IA",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Name of the Internet Access service",
			Computed:    true,
		},
		"bandwidth": schema.Int32Attribute{
			Description: "Bandwidth of the service",
			Computed:    true,
		},
		"bandwidth_commit": schema.Int32Attribute{
			Description: "Minimum bandwidth commit for the burst billing variant of the service",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "Provisioning state of the service. One of: PROVISIONING, PROVISIONED, FAILED, DEPROVISIONING, DEPROVISIONED",
			Computed:    true,
		},
		"use_case": schema.StringAttribute{
			Description: "Use case of the service. One of: MAIN, MANAGEMENT_ACCESS",
			Computed:    true,
		},
		"change": schema.SingleNestedAttribute{
			Description: "Current state of the latest change on the service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
			Attributes: map[string]schema.Attribute{
				"href": schema.StringAttribute{
					Description: "URI to the change resource",
					Computed:    true,
				},
			},
		},
		"locations": schema.ListNestedAttribute{
			Description: "Locations associated with the service",
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
						Description: "Region of the location. One of: AMER, EMEA, APAC",
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
			Description: "Routing protocol configuration of the service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[RoutingProtocolModel](ctx),
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "Routing protocol type. One of: BGP, DIRECT, STATIC",
					Computed:    true,
				},
				"customer_routes": schema.ListNestedAttribute{
					Description: "Customer routes announced through the service",
					Computed:    true,
					CustomType:  fwtypes.NewListNestedObjectTypeOf[CustomerRouteModel](ctx),
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"ip_block": schema.SingleNestedAttribute{
								Description: "IP block of the route",
								Computed:    true,
								CustomType:  fwtypes.NewObjectTypeOf[IpBlockModel](ctx),
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{
										Description: "URI to the IP block",
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
					Description: "Fabric connections attached to the service",
					Computed:    true,
					CustomType:  fwtypes.NewListNestedObjectTypeOf[ConnectionModel](ctx),
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"href": schema.StringAttribute{
								Description: "URI to the connection",
								Computed:    true,
							},
							"uuid": schema.StringAttribute{
								Description: "Unique identifier of the connection",
								Computed:    true,
							},
							"peering_ipv4": peeringSchema(ctx, "IPv4"),
							"peering_ipv6": peeringSchema(ctx, "IPv6"),
						},
					},
				},
			},
		},
		"billing": schema.SingleNestedAttribute{
			Description: "Billing details of the service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[BillingModel](ctx),
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "Billing type. One of: FIXED, USAGE_BASED, BURST_BASED",
					Computed:    true,
				},
				"enabled": schema.BoolAttribute{
					Description: "Whether billing is enabled",
					Computed:    true,
				},
				"start_date": schema.StringAttribute{
					Description: "Start date of the billing period",
					Computed:    true,
				},
			},
		},
		"account": schema.SingleNestedAttribute{
			Description: "Account associated with the service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[AccountModel](ctx),
			Attributes: map[string]schema.Attribute{
				"account_number": schema.StringAttribute{
					Description: "Account number",
					Computed:    true,
				},
				"href": schema.StringAttribute{
					Description: "URI to the account",
					Computed:    true,
				},
			},
		},
		"project": schema.SingleNestedAttribute{
			Description: "Project associated with the service",
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
			Description: "Order details of the service",
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
					Description: "Billing tier for the connection bandwidth",
					Computed:    true,
				},
				"order_id": schema.StringAttribute{
					Description: "Order identification",
					Computed:    true,
				},
				"order_number": schema.StringAttribute{
					Description: "Order reference number",
					Computed:    true,
				},
				"term_length": schema.Int32Attribute{
					Description: "Term length in months",
					Computed:    true,
				},
				"contracted_bandwidth": schema.Int32Attribute{
					Description: "Contracted bandwidth",
					Computed:    true,
				},
				"href": schema.StringAttribute{
					Description: "URI to the order",
					Computed:    true,
				},
			},
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Audit log of changes to the service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
			Attributes: map[string]schema.Attribute{
				"created_date_time": schema.StringAttribute{
					Description: "Creation time of the service",
					Computed:    true,
				},
				"updated_date_time": schema.StringAttribute{
					Description: "Last update time of the service",
					Computed:    true,
				},
				"deleted_date_time": schema.StringAttribute{
					Description: "Deletion time of the service",
					Computed:    true,
				},
			},
		},
	}
}
