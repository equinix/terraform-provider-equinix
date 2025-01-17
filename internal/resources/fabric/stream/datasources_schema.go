package stream

import (
	"context"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func dataSourceAllStreamsSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned streams list",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"offset": schema.Int32Attribute{
						Description: "The page offset for the returned streams list",
						Optional:    true,
						Computed:    true,
					},
					"limit": schema.Int32Attribute{
						Description: "The page size for the returned streams list",
						Optional:    true,
						Computed:    true,
					},
					"total": schema.Int32Attribute{
						Description: "The total number of streams available to the user making the request",
						Computed:    true,
					},
				},
			},
			"data": schema.ListNestedAttribute{
				Description: "Returned list of stream objects",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: getStreamSchema(),
				},
			},
		},
	}
}

func dataSourceSingleStreamSchema(ctx context.Context) schema.Schema {
	baseStreamSchema := getStreamSchema()
	baseStreamSchema["id"] = framework.IDAttributeDefaultDescription()
	baseStreamSchema["stream_id"] = schema.StringAttribute{
		Description: "The uuid of the stream this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Attributes: baseStreamSchema,
	}
}

func getStreamSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "Equinix defined Streaming Type",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Customer-provided name of the stream resource",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Customer-provided description of the stream resource",
			Computed:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Boolean switch enabling streaming data for the stream resource",
			Computed:    true,
		},
		"project": schema.SingleNestedAttribute{
			Description: "Equinix Project attribute object",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"project_id": schema.StringAttribute{
					Description: "Equinix Subscriber-assigned project ID",
					Computed:    true,
				},
			},
		},
		"href": schema.StringAttribute{
			Description: "Equinix auto generated URI to the stream resource in Equinix Portal",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix-assigned unique id for the stream resource",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "Value representing provisioning status for the stream resource",
			Computed:    true,
		},
		"assets_count": schema.Int32Attribute{
			Description: "Count of the streaming assets attached to the stream resource",
			Computed:    true,
		},
		"stream_subscriptions_count": schema.Int32Attribute{
			Description: "Count of the client subscriptions on the stream resource",
			Computed:    true,
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Details of the last change on the stream resource",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"created_by": schema.StringAttribute{
					Description: "User name of creator of the stream resource",
					Computed:    true,
				},
				"created_by_full_name": schema.StringAttribute{
					Description: "Legal name of creator of the stream resource",
					Computed:    true,
				},
				"created_by_email": schema.StringAttribute{
					Description: "Email of creator of the stream resource",
					Computed:    true,
				},
				"created_date_time": schema.StringAttribute{
					Description: "Creation time of the stream resource",
					Computed:    true,
				},
				"updated_by": schema.StringAttribute{
					Description: "User name of last updater of the stream resource",
					Computed:    true,
				},
				"updated_by_full_name": schema.StringAttribute{
					Description: "Legal name of last updater of the stream resource",
					Computed:    true,
				},
				"updated_by_email": schema.StringAttribute{
					Description: "Email of last updater of the stream resource",
					Computed:    true,
				},
				"updated_date_time": schema.StringAttribute{
					Description: "Last update time of the stream resource",
					Computed:    true,
				},
				"deleted_by": schema.StringAttribute{
					Description: "User name of deleter of the stream resource",
					Computed:    true,
				},
				"deleted_by_full_name": schema.StringAttribute{
					Description: "Legal name of deleter of the stream resource",
					Computed:    true,
				},
				"deleted_by_email": schema.StringAttribute{
					Description: "Email of deleter of the stream resource",
					Computed:    true,
				},
				"deleted_date_time": schema.StringAttribute{
					Description: "Deletion time of the stream resource",
					Computed:    true,
				},
			},
		},
	}
}
