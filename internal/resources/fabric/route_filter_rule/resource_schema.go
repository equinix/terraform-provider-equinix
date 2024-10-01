package route_filter_rule

import (
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"route_filter_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "UUID of the Route Filter Policy to apply this Rule to",
		},
		"prefix": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "IP Address Prefix to Filter on",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Name of the Route Filter",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional description to add to the Route Filter you will be creating",
		},
		"prefix_match": {
			Type:        schema.TypeString,
			Optional:    true,
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
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix Assigned ID for Route Filter Rule",
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

func changeSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI of the previous Route Filter Rule Change",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of change. One of [ \"BGP_IPv4_PREFIX_FILTER_RULE_UPDATE\",\"BGP_IPv4_PREFIX_FILTER_RULE_CREATION\",\"BGP_IPv4_PREFIX_FILTER_RULE_DELETION\",\"BGP_IPv6_PREFIX_FILTER_RULE_UPDATE\",\"BGP_IPv6_PREFIX_FILTER_RULE_CREATION\",\"BGP_IPv6_PREFIX_FILTER_RULE_DELETION\" ]",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the previous change",
			},
		},
	}
}
