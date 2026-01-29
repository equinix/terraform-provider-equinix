package advertised_route

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAdvertisedRoutesSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"connection_id": schema.StringAttribute{
				Description: "The uuid of the routes this data source should retrieve",
				Required:    true,
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned advertised routes list",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[paginationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"offset": schema.Int32Attribute{
						Description: "Index of the first item returned in the response.",
						Optional:    true,
						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},
					"limit": schema.Int32Attribute{
						Description: "Maximum number of search results returned per page.",
						Optional:    true,
						Validators: []validator.Int32{
							int32validator.Between(1, 100),
						},
					},
					"total": schema.Int32Attribute{
						Description: "The total number of elements returned",
						Computed:    true,
						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},
					"next": schema.StringAttribute{
						Description: "URL relative to the next item in the response.",
						Computed:    true,
					},
					"previous": schema.StringAttribute{
						Description: "URL relative to the previous item in the response.",
						Computed:    true,
					},
				},
			},
			"data": schema.ListNestedAttribute{
				Description: "Returned list of advertised routes objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[advertisedRoutesBaseModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getAdvertisedRoutesSchema(ctx),
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Required:    true,
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
			"sort": schema.SingleNestedAttribute{
				Description: "Sort details for the returned advertised routes list",
				Optional:    true,
				CustomType:  fwtypes.NewObjectTypeOf[sortModel](ctx),
				Attributes: map[string]schema.Attribute{
					"direction": schema.StringAttribute{
						Description: "Sort direction, one of [ASC, DESC]",
						Optional:    true,
					},
					"property": schema.StringAttribute{
						Description: "Property name to sort by",
						Optional:    true,
					},
				},
			},
		},
	}
}

func getAdvertisedRoutesSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "Indicator of a advertised route",
			Computed:    true,
		},
		"connection_id": schema.StringAttribute{
			Description: "The uuid of the routes this data source should retrieve",
			Required:    true,
		},
		"protocol_type": schema.StringAttribute{
			Description: "Advertised Route protocol type",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "State of the advertised Route",
			Computed:    true,
		},
		"prefix": schema.StringAttribute{
			Description: "Prefix of the Advertised Route",
			Computed:    true,
		},
		"next_hop": schema.StringAttribute{
			Description: "Next Hop of the Advertised Route",
			Computed:    true,
		},
		"med": schema.Int32Attribute{
			Description: "Multi-Exit Discriminator for the Advertised Route",
			Computed:    true,
		},
		"local_preference": schema.Int32Attribute{
			Description: "This field holds local preference of the advertised route.",
			Computed:    true,
		},
		"as_path": schema.ListAttribute{
			Description: "List of supported AS Paths for the Advertised Routes.",
			CustomType:  fwtypes.ListOfStringType,
			ElementType: types.StringType,
			Computed:    true,
		},
		"connection": schema.SingleNestedAttribute{
			Description: "connection of the route table entry",
			CustomType:  fwtypes.NewObjectTypeOf[connectionModel](ctx),
			Attributes: map[string]schema.Attribute{
				"uuid": schema.StringAttribute{
					Description: "UUID of the Connection",
					Computed:    true,
				},
				"name": schema.StringAttribute{
					Description: "Name of the Connection",
					Computed:    true,
				},
				"href": schema.StringAttribute{
					Description: "HREF of the Connection",
					Computed:    true,
				},
			},
			Computed: true,
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Change Log of the route table entry",
			CustomType:  fwtypes.NewObjectTypeOf[connectionModel](ctx),
			Attributes: map[string]schema.Attribute{
				"created_by": schema.StringAttribute{
					Description: "Created by User Key",
					Computed:    true,
				},
				"created_by_full_name": schema.StringAttribute{
					Description: "Created by User Full Name",
					Computed:    true,
				},
				"href": schema.StringAttribute{
					Description: "HREF of the Connection",
					Computed:    true,
				},
			},
			Computed: true,
		},
	}
}
