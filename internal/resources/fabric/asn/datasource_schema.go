package asn

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAllAsnsSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows users to search for Equinix Autonomous System Numbers (ASNs) with pagination details.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/ASNs`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"data": schema.ListNestedAttribute{
				Description: "Returned list of ASN objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BaseAsnModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getAsnSchema(ctx),
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"property": schema.StringAttribute{
						Description: "Field name to use on filters, e.g. \"/project/projectId\", \"/uuid\", \"/state\", \"/type\"",
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
				Description: "Pagination details for the returned ASNs list",
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
						Description: "The total number of ASNs available to the user making the request",
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
						Description: "The property name to use in sorting, e.g. \"/state\". Allowed values: /changeLog/createdDateTime, /uuid, /state, /project/projectId, /type",
						Optional:    true,
					},
				},
			},
		},
	}
}

func dataSourceSingleAsnSchema(ctx context.Context) schema.Schema {
	baseSchema := getAsnSchema(ctx)
	baseSchema["id"] = framework.IDAttributeDefaultDescription()
	baseSchema["asn_id"] = schema.StringAttribute{
		Description: "The UUID of the ASN this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allows users to fetch an Autonomous System Number (ASN) by UUID.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/ASNs`,
		Attributes: baseSchema,
	}
}

func getAsnSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique identifier for the ASN record",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the ASN",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "State of the ASN. One of: REQUESTED, PROVISIONING, PENDING, PROVISIONED, ATTACHED, FAILED, PROVISIONING_FAILED, DEPROVISIONING, DEPROVISIONING_FAILED, DEPROVISIONED",
			Computed:    true,
		},
		"number": schema.Int64Attribute{
			Description: "Autonomous System Number",
			Computed:    true,
		},
		"rir_region": schema.StringAttribute{
			Description: "Regional Internet Registry region of the ASN. One of: AFRINIC, APNIC, ARIN, LACNIC, RIPE_NCC",
			Computed:    true,
		},
		"range": schema.StringAttribute{
			Description: "Range of the ASN. One of: 16_BIT, 32_BIT",
			Computed:    true,
		},
		"project": schema.SingleNestedAttribute{
			Description: "Project the ASN belongs to",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ProjectModel](ctx),
			Attributes: map[string]schema.Attribute{
				"project_id": schema.StringAttribute{
					Description: "Project identifier",
					Computed:    true,
				},
			},
		},
		"location": schema.SingleNestedAttribute{
			Description: "Location of the ASN",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[LocationModel](ctx),
			Attributes: map[string]schema.Attribute{
				"metro_href": schema.StringAttribute{Description: "URI to the metro", Computed: true},
				"metro_code": schema.StringAttribute{Description: "Metro code", Computed: true},
				"metro_name": schema.StringAttribute{Description: "Metro name", Computed: true},
				"region":     schema.StringAttribute{Description: "Region of the metro", Computed: true},
				"ibx":        schema.StringAttribute{Description: "IBX site code (deprecated)", Computed: true},
			},
		},
		"ownership_challenge_details": schema.SingleNestedAttribute{
			Description: "Details of the ASN ownership challenge",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[OwnershipChallengeModel](ctx),
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "Ownership challenge type. One of: LOA, RPKI",
					Computed:    true,
				},
			},
		},
		"change": schema.SingleNestedAttribute{
			Description: "Current state of the latest change on the ASN",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
			Attributes: map[string]schema.Attribute{
				"href": schema.StringAttribute{
					Description: "URI to the change resource",
					Computed:    true,
				},
			},
		},
		"assets": schema.ListNestedAttribute{
			Description: "Assets (services) the ASN is attached to",
			Computed:    true,
			CustomType:  fwtypes.NewListNestedObjectTypeOf[AssetModel](ctx),
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"uuid":             schema.StringAttribute{Description: "Asset identifier", Computed: true},
					"announcing_metro": schema.StringAttribute{Description: "Metro announcing the ASN", Computed: true},
					"service": schema.SingleNestedAttribute{
						Description: "Service the ASN is attached to",
						Computed:    true,
						CustomType:  fwtypes.NewObjectTypeOf[AssetServiceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"href": schema.StringAttribute{Description: "URI to the service", Computed: true},
							"type": schema.StringAttribute{Description: "Service type. One of: IA", Computed: true},
							"uuid": schema.StringAttribute{Description: "Service identifier", Computed: true},
						},
					},
					"ip_blocks": schema.ListNestedAttribute{
						Description: "IP blocks announced with the ASN",
						Computed:    true,
						CustomType:  fwtypes.NewListNestedObjectTypeOf[AssetIpBlockModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"href": schema.StringAttribute{Description: "URI to the IP block", Computed: true},
								"type": schema.StringAttribute{Description: "IP block reference type. One of: IPV4_IP_BLOCK, IPV6_IP_BLOCK", Computed: true},
								"uuid": schema.StringAttribute{Description: "IP block identifier", Computed: true},
							},
						},
					},
				},
			},
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Details of the last change on the ASN",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
			Attributes: map[string]schema.Attribute{
				"created_date_time": schema.StringAttribute{Description: "Creation date and time", Computed: true},
				"updated_date_time": schema.StringAttribute{Description: "Last update date and time", Computed: true},
				"deleted_date_time": schema.StringAttribute{Description: "Deletion date and time", Computed: true},
			},
		},
	}
}
