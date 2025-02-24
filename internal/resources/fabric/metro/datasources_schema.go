package metro

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAllMetroSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"presence": schema.StringAttribute{
				Description: "User On Boarded Metros based on Fabric resource availability",
				Optional:    true,
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned metro list",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[PaginationModel](ctx),
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
						Description: "The total number of metro returned",
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
				Description: "Returned list of metro objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[Model](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getMetroSchema(ctx),
				},
			},
		},
	}
}

func dataSourceSingleMetroSchema(ctx context.Context) schema.Schema {
	baseMetroSchema := getMetroSchema(ctx)
	baseMetroSchema["id"] = framework.IDAttributeDefaultDescription()
	baseMetroSchema["metro_code"] = schema.StringAttribute{
		Description: "The metro code this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Attributes: baseMetroSchema,
	}
}

func getMetroSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"href": schema.StringAttribute{
			Description: "The canonical URL at which the resource resides",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Indicator of a fabric metro",
			Computed:    true,
		},
		"code": schema.StringAttribute{
			Description: "Code assigned to an Equinix IBX data center in a specified metropolitan area",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "Board geographical area in which the data center is located",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Name of the region in which the data center is located",
			Computed:    true,
		},
		"equinix_asn": schema.Int64Attribute{
			Description: "Autonomous system number (ASN) for a specified Fabric metro. The ASN is a unique identifier that carries the network routing protocol and exchanges that data with other internal systems via border gateway protocol.",
			Computed:    true,
		},
		"local_vc_bandwidth_max": schema.Int64Attribute{
			Description: "This field holds Max Connection speed within the metro.",
			Computed:    true,
		},
		"geo_coordinates": schema.SingleNestedAttribute{
			Description: "Geographic location data of Fabric Metro",
			CustomType:  fwtypes.NewObjectTypeOf[GeoCoordinatesModel](ctx),
			Attributes: map[string]schema.Attribute{
				"latitude": schema.Float64Attribute{
					Description: "Latitude of the Metro",
					Computed:    true,
				},
				"longitude": schema.Float64Attribute{
					Description: "Longitude of the Metro",
					Computed:    true,
				},
			},
			Computed: true,
		},
		"connected_metros": schema.ListAttribute{
			Description: "Arrays of objects containing latency data for the specified metro",
			CustomType:  fwtypes.NewListNestedObjectTypeOf[ConnectedMetroModel](ctx),
			ElementType: fwtypes.NewObjectTypeOf[ConnectedMetroModel](ctx),
			Computed:    true,
		},
		"geo_scopes": schema.ListAttribute{
			Description: "List of supported geographic boundaries of a Fabric Metro. Example values: CANADA, CONUS.",
			CustomType:  fwtypes.ListOfStringType,
			ElementType: types.StringType,
			Computed:    true,
		},
	}
}
