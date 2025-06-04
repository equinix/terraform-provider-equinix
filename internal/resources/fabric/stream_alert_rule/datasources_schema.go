package stream_alert_rule

import (
	"context"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAllStreamAlertRulesSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data source that allows user to fetch Equinix Fabric Stream Alert Rules with pagination

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Stream-Alert-Rules`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"stream_id": schema.StringAttribute{
				Description: "The uuid of the stream that is the target of the stream alert rule",
				Required:    true,
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned stream alert rules list",
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
						Description: "The total number of alert rules available to the user making the request",
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
				CustomType:  fwtypes.NewListNestedObjectTypeOf[baseStreamAlertRulesModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getStreamAlertRuleSchema(ctx),
				},
			},
		},
	}
}

func dataSourceStreamAlertRuleByID(ctx context.Context) schema.Schema {
	baseStreamAlertRuleSchema := getStreamAlertRuleSchema(ctx)
	baseStreamAlertRuleSchema["id"] = framework.IDAttributeDefaultDescription()
	baseStreamAlertRuleSchema["stream_id"] = schema.StringAttribute{
		Description: "The uuid of the stream that is the target of the stream alert rule",
		Required:    true,
	}
	baseStreamAlertRuleSchema["alert_rule_id"] = schema.StringAttribute{
		Description: "The uuid of the stream alert rule",
		Required:    true,
	}

	return schema.Schema{
		Description: `Fabric V4 API compatible data source that allows user to fetch Equinix Fabric Stream Alert Rule by Stream Id and Alert Rule Id

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Stream-Alert-Rules`,
		Attributes: baseStreamAlertRuleSchema,
	}
}

func getStreamAlertRuleSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "Type of the stream alert rule",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "Customer-provided stream alert rule name",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Customer-provided stream alert rule description",
			Computed:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Stream subscription enabled status",
			Computed:    true,
		},
		"resource_selector": schema.SingleNestedAttribute{
			Description: "Lists of metrics to be included/excluded on the stream alert rule",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[selectorModel](ctx),
			Attributes: map[string]schema.Attribute{
				"include": schema.ListAttribute{
					Description: "List of metrics to include",
					ElementType: types.StringType,
					Computed:    true,
				},
			},
		},
		"metric_name": schema.StringAttribute{
			Description: "Stream alert rule metric name",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix assigned unique identifier of the stream subscription resource",
			Computed:    true,
		},

		"window_size": schema.StringAttribute{
			Description: "Stream alert rule metric window size",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix assigned URI of the stream alert rule resource",
			Computed:    true,
		},
		"operand": schema.StringAttribute{
			Description: "Stream alert rule metric operand",
			Computed:    true,
		},
		"warning_threshold": schema.StringAttribute{
			Description: "Stream alert rule metric warning threshold",
			Computed:    true,
		},
		"critical_threshold": schema.StringAttribute{
			Description: "Stream alert rule metric critical threshold",
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
