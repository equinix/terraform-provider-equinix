package cloud_router

import (
	"context"

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
			// add ID etc like in metros?
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
		},
	}
}

func getAdvertisedRoutesSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "Indicator of a advertised route",
			Computed:    true,
		},
		"protocolType": schema.StringAttribute{
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
		"nextHop": schema.StringAttribute{
			Description: "Next Hop of the Advertised Route",
			Computed:    true,
		},
		"MED": schema.Int32Attribute{
			Description: "Multi-Exit Discriminator for the Advertised Route",
			Computed:    true,
		},
		"localPreference": schema.Int32Attribute{
			Description: "This field holds local preference of the advertised route.",
			Computed:    true,
		},
		"asPath": schema.ListAttribute{
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
		"changeLog": schema.SingleNestedAttribute{
			Description: "Change Log of the route table entry",
			CustomType:  fwtypes.NewObjectTypeOf[connectionModel](ctx),
			Attributes: map[string]schema.Attribute{
				"createdBy": schema.StringAttribute{
					Description: "Created by User Key",
					Computed:    true,
				},
				"createdByFullName": schema.StringAttribute{
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
