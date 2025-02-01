package stream_subscription

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Stream Subscription

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Stream-Subscriptions`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"stream_id": schema.StringAttribute{
				Description: "The uuid of the stream that is the target of the stream subscription",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the stream subscription request",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Customer-provided stream subscription name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Customer-provided stream subscription description",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Stream subscription enabled status",
				Required:    true,
			},
			"filters": schema.ListNestedAttribute{
				Description: "List of filters to apply to the stream subscription selectors. Maximum of 8. All will be AND'd together with 1 of the 8 being a possible OR group of 3",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property": schema.StringAttribute{
							Description: "Property to apply the filter to",
							Required:    true,
						},
						"operator": schema.StringAttribute{
							Description: "Operation applied to the values of the filter",
							Required:    true,
						},
						"values": schema.ListAttribute{
							Description: "List of values to apply the operation to for the specified property",
							Required:    true,
							ElementType: types.StringType,
						},
						"or": schema.BoolAttribute{
							Description: "Boolean value to specify if this filter is a part of the OR group. Has a maximum of 3 and only counts for 1 of the 8 possible filters",
							Optional:    true,
						},
					},
				},
			},
			"metric_selector": schema.SingleNestedAttribute{
				Description: "Lists of metrics to be included/excluded on the stream subscription",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"include": schema.ListAttribute{
						Description: "List of metrics to include",
						ElementType: types.StringType,
						Required:    true,
					},
					"except": schema.ListAttribute{
						Description: "List of metrics to exclude",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			"event_selector": schema.SingleNestedAttribute{
				Description: "Lists of events to be included/excluded on the stream subscription",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"include": schema.ListAttribute{
						Description: "List of events to include",
						ElementType: types.StringType,
						Required:    true,
					},
					"except": schema.ListAttribute{
						Description: "List of events to exclude",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			"sink": schema.SingleNestedAttribute{
				Description: "The details of the subscriber to the Equinix Stream",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"uri": schema.StringAttribute{
						Description: "Publicly reachable http endpoint destination for data stream",
						Optional:    true,
					},
					"type": schema.StringAttribute{
						Description: "Type of the subscriber",
						Required:    true,
					},
					"batch_enabled": schema.BoolAttribute{
						Description: "Boolean switch enabling batch delivery of data",
						Optional:    true,
					},
					"batch_size_max": schema.Int32Attribute{
						Description: "Maximum size of the batch delivery if enabled",
						Optional:    true,
					},
					"batch_wait_time_max": schema.Int32Attribute{
						Description: "Maximum time to wait for batch delivery if enabled",
						Optional:    true,
					},
					"host": schema.StringAttribute{
						Description: "Known hostname of certain data stream subscription products. Not to be confused with a variable URI",
						Optional:    true,
					},
					"credential": schema.SingleNestedAttribute{
						Description: "Access details for the specified sink type",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "Type of the credential being passed",
								Required:    true,
							},
							"access_token": schema.StringAttribute{
								Description: "Passed as Authorization header value",
								Optional:    true,
							},
							"integration_key": schema.StringAttribute{
								Description: "Passed as Authorization header value",
								Optional:    true,
							},
							"api_key": schema.StringAttribute{
								Description: "Passed as Authorization header value",
								Optional:    true,
							},
							"username": schema.StringAttribute{
								Description: "Passed as Authorization header value",
								Optional:    true,
							},
							"password": schema.StringAttribute{
								Description: "Passed as Authorization header value",
								Optional:    true,
							},
						},
					},
					"settings": schema.SingleNestedAttribute{
						Description: "Stream subscription sink settings",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"event_index": schema.StringAttribute{
								Optional: true,
							},
							"metric_index": schema.StringAttribute{
								Optional: true,
							},
							"source": schema.StringAttribute{
								Optional: true,
							},
							"application_key": schema.StringAttribute{
								Optional: true,
							},
							"event_uri": schema.StringAttribute{
								Optional: true,
							},
							"metric_uri": schema.StringAttribute{
								Optional: true,
							},
							"transform_alerts": schema.BoolAttribute{
								Optional: true,
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
				CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
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
		},
	}
}
