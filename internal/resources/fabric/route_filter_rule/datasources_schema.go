package route_filter_rule

import (
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix Assigned ID for Route Filter Rule to retrieve data for",
		},
		"prefix": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP Address Prefix to Filter on",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the Route Filter",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Optional description to add to the Route Filter you will be creating",
		},
		"prefix_match": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Prefix matching operator. One of [ orlonger, exact ] Default: \"orlonger\"",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Route Filter Type. One of [ BGP_IPv4_PREFIX_FILTER_RULE, BGP_IPv6_PREFIX_FILTER_RULE ] ",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Route filter rules URI",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State of the Route Filter Rule in its lifecycle",
		},
		"action": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Action that will be taken on IP Addresses matching the rule",
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "An object with the details of the previous change applied on the Route Filter",
			Elem:        changeSch(),
		},
		"change_log": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
	}
}

func routeFilterSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "UUID of the Route Filter Policy the rule is attached to",
	}
}

func paginationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The page offset for the pagination request. Index of the first element. Default is 0.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of elements to be requested per page. Number must be between 1 and 100. Default is 20",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of elements returned.",
			},
			"next": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL relative to the last item in the response.",
			},
			"previous": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL relative to the first item in the response.",
			},
		},
	}
}

func dataSourceByUUIDSchema() map[string]*schema.Schema {
	baseSchema := dataSourceBaseSchema()
	baseSchema["route_filter_id"] = routeFilterSchema()

	return baseSchema
}

func dataSourceAllRulesForRouteFilterSchema() map[string]*schema.Schema {
	baseSchema := dataSourceBaseSchema()
	baseSchema["uuid"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Equinix Assigned ID for Route Filter Rule to retrieve data for",
	}
	routeFilterRulesSchema := map[string]*schema.Schema{
		"route_filter_id": routeFilterSchema(),
		"offset": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "The page offset for the pagination request. Index of the first element. Default is 0.",
		},
		"limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     20,
			Description: "Number of elements to be requested per page. Number must be between 1 and 100. Default is 20",
		},
		"pagination": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Pagination details for the Data Source Search Request",
			Elem:        paginationSchema(),
		},
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The list of Rules attached to the given Route Filter Policy UUID",
			Elem: &schema.Resource{
				Schema: baseSchema,
			},
		},
	}

	return routeFilterRulesSchema
}
