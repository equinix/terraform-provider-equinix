package route_filter

import (
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"BGP_IPv4_PREFIX_FILTER", "BGP_IPv6_PREFIX_FILTER"}, false),
			Description:  "Route Filter Type. One of [ \"BGP_IPv4_PREFIX_FILTER\", \"BGP_IPv6_PREFIX_FILTER\" ] ",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the Route Filter",
		},
		"project": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "The Project object that contains project_id and href that is related to the Fabric Project containing connections the Route Filter can be attached to",
			Elem:        projectSch(),
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional description to add to the Route Filter you will be creating",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Route filter URI",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix Assigned ID for Route Filter",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State of the Route Filter in its lifecycle",
		},
		"not_matched_rule_action": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The action that will be taken on ip ranges that don't match the rules present within the Route Filter",
		},
		"connections_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The number of Fabric Connections that this Route Filter is attached to",
		},
		"rules_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The number of Route Filter Rules attached to this Route Filter",
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

func projectSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project id associated with Fabric Project",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI of the Fabric Project",
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
				Description: "The URI of the previous Route Filter Change",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of change. One of [ \"BGP_IPv4_PREFIX_FILTER_UPDATE\",\"BGP_IPv4_PREFIX_FILTER_CREATION\",\"BGP_IPv4_PREFIX_FILTER_DELETION\",\"BGP_IPv6_PREFIX_FILTER_UPDATE\",\"BGP_IPv6_PREFIX_FILTER_CREATION\",\"BGP_IPv6_PREFIX_FILTER_DELETION\" ]",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the previous change",
			},
		},
	}
}
