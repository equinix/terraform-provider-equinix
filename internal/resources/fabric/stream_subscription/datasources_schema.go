package streamsubscription

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAllStreamSubscriptionsSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data source that allows user to fetch Equinix Fabric Stream Subscriptions with pagination

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Stream-Subscriptions`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"stream_id": schema.StringAttribute{
				Description: "The uuid of the stream that is the target of the stream subscription",
				Required:    true,
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned streams list",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[paginationModel](ctx),
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
				CustomType:  fwtypes.NewListNestedObjectTypeOf[baseStreamSubscriptionModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getStreamSubscriptionSchema(ctx),
				},
			},
		},
	}
}

func dataSourceStreamSubscriptionByID(ctx context.Context) schema.Schema {
	baseStreamSchema := getStreamSubscriptionSchema(ctx)
	baseStreamSchema["id"] = framework.IDAttributeDefaultDescription()
	baseStreamSchema["stream_id"] = schema.StringAttribute{
		Description: "The uuid of the stream that is the target of the stream subscription",
		Required:    true,
	}
	baseStreamSchema["subscription_id"] = schema.StringAttribute{
		Description: "The uuid of the stream subscription",
		Required:    true,
	}

	return schema.Schema{
		Description: `Fabric V4 API compatible data source that allows user to fetch Equinix Fabric Stream Subscription by Stream Id and Subscription Id

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Stream-Subscriptions`,
		Attributes: baseStreamSchema,
	}
}

func getStreamSubscriptionSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "Type of the stream subscription request",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Customer-provided stream subscription name",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Customer-provided stream subscription description",
			Computed:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Stream subscription enabled status",
			Computed:    true,
		},
		"filters": schema.ListNestedAttribute{
			Description: "List of filters to apply to the stream subscription selectors. Maximum of 8. All will be AND'd together with 1 of the 8 being a possible OR group of 3",
			Computed:    true,
			CustomType:  fwtypes.NewListNestedObjectTypeOf[filterModel](ctx),
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"property": schema.StringAttribute{
						Description: "Property to apply the filter to",
						Computed:    true,
					},
					"operator": schema.StringAttribute{
						Description: "Operation applied to the values of the filter",
						Computed:    true,
					},
					"values": schema.ListAttribute{
						Description: "List of values to apply the operation to for the specified property",
						Computed:    true,
						ElementType: types.StringType,
					},
					"or": schema.BoolAttribute{
						Description: "Boolean value to specify if this filter is a part of the OR group. Has a maximum of 3 and only counts for 1 of the 8 possible filters",
						Computed:    true,
					},
				},
			},
		},
		"metric_selector": schema.SingleNestedAttribute{
			Description: "Lists of metrics to be included/excluded on the stream subscription",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[selectorModel](ctx),
			Attributes: map[string]schema.Attribute{
				"include": schema.ListAttribute{
					Description: "List of metrics to include",
					ElementType: types.StringType,
					Computed:    true,
				},
				"except": schema.ListAttribute{
					Description: "List of metrics to exclude",
					ElementType: types.StringType,
					Computed:    true,
				},
			},
		},
		"event_selector": schema.SingleNestedAttribute{
			Description: "Lists of events to be included/excluded on the stream subscription",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[selectorModel](ctx),
			Attributes: map[string]schema.Attribute{
				"include": schema.ListAttribute{
					Description: "List of events to include",
					ElementType: types.StringType,
					Computed:    true,
				},
				"except": schema.ListAttribute{
					Description: "List of events to exclude",
					ElementType: types.StringType,
					Computed:    true,
				},
			},
		},
		"sink": schema.SingleNestedAttribute{
			Description: "The details of the subscriber to the Equinix Stream",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[sinkModel](ctx),
			Attributes: map[string]schema.Attribute{
				"uri": schema.StringAttribute{
					Description: "Publicly reachable http endpoint destination for data stream",
					Computed:    true,
				},
				"type": schema.StringAttribute{
					Description: "Type of the subscriber",
					Computed:    true,
				},
				"batch_enabled": schema.BoolAttribute{
					Description: "Boolean switch enabling batch delivery of data",
					Computed:    true,
				},
				"batch_size_max": schema.Int32Attribute{
					Description: "Maximum size of the batch delivery if enabled",
					Computed:    true,
				},
				"batch_wait_time_max": schema.Int32Attribute{
					Description: "Maximum time to wait for batch delivery if enabled",
					Computed:    true,
				},
				"host": schema.StringAttribute{
					Description: "Known hostname of certain data stream subscription products. Not to be confused with a variable URI",
					Computed:    true,
				},
				"credential": schema.SingleNestedAttribute{
					Description: "Access details for the specified sink type",
					Computed:    true,
					CustomType:  fwtypes.NewObjectTypeOf[sinkCredentialModel](ctx),
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Type of the credential being passed",
							Computed:    true,
						},
						"access_token": schema.StringAttribute{
							Description: "Passed as Authorization header value",
							Computed:    true,
						},
						"integration_key": schema.StringAttribute{
							Description: "Passed as Authorization header value",
							Computed:    true,
						},
						"api_key": schema.StringAttribute{
							Description: "Passed as Authorization header value",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "Passed as Authorization header value",
							Computed:    true,
						},
						"password": schema.StringAttribute{
							Description: "Passed as Authorization header value",
							Computed:    true,
						},
					},
				},
				"settings": schema.SingleNestedAttribute{
					Description: "Stream subscription sink settings",
					Computed:    true,
					CustomType:  fwtypes.NewObjectTypeOf[sinkCredentialModel](ctx),
					Attributes: map[string]schema.Attribute{
						"event_index": schema.StringAttribute{
							Computed: true,
						},
						"metric_index": schema.StringAttribute{
							Computed: true,
						},
						"source": schema.StringAttribute{
							Computed: true,
						},
						"application_key": schema.StringAttribute{
							Computed: true,
						},
						"event_uri": schema.StringAttribute{
							Computed: true,
						},
						"metric_uri": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
		"href": schema.StringAttribute{
			Description: "Equinix assigned URI of the stream subscription resource",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix assigned unique identifier of the stream subscription resource",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "Value representing provisioning status for the stream resource",
			Computed:    true,
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Details of the last change on the stream resource",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[changeLogModel](ctx),
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
