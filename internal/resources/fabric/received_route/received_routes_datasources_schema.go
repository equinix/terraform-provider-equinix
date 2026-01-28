package advertised_route

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceReceivedRoutesSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned received routes list",
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
				Description: "Returned list of received routes objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[receivedRoutesBaseModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getReceivedRoutesSchema(ctx),
				},
			},
		},
	}
}

func getReceivedRoutesSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "Indicator of a received route",
			Computed:    true,
		},
		"protocol_type": schema.StringAttribute{
			Description: "Received Route protocol type",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "State of the Received Route",
			Computed:    true,
		},
		"prefix": schema.StringAttribute{
			Description: "Prefix of the Received Route",
			Computed:    true,
		},
		"next_hop": schema.StringAttribute{
			Description: "Next Hop of the Received Route",
			Computed:    true,
		},
		"med": schema.Int32Attribute{
			Description: "Multi-Exit Discriminator for the Received Route",
			Computed:    true,
		},
		"local_preference": schema.Int32Attribute{
			Description: "This field holds local preference of the Received route.",
			Computed:    true,
		},
		"as_path": schema.ListAttribute{
			Description: "List of supported AS Paths for the Received Routes.",
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
