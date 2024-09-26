package route_filter

import (
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Route Filter Type. One of [ \"BGP_IPv4_PREFIX_FILTER\", \"BGP_IPv6_PREFIX_FILTER\" ] ",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the Route Filter",
		},
		"project": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "The Project object that contains project_id and href that is related to the Fabric Project containing connections the Route Filter can be attached to",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"project_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Project id associated with Fabric Project",
					},
					"href": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "URI of the Fabric Project",
					},
				},
			},
		},
		"description": {
			Type:        schema.TypeString,
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
			Required:    true,
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
			Elem: &schema.Resource{
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
			},
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

func paginationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The page offset for the pagination request. Index of the first element. Default is 0.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of elements to be requested per page. Number must be between 1 and 100. Default is 20",
			},
			"total": {
				Type:        schema.TypeInt,
				Optional:    true,
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

func dataSourceSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Route Filters",
			Elem: &schema.Resource{
				Schema: dataSourceBaseSchema(),
			},
		},
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Filters for the Data Source Search Request. Maximum of 8 total filters.",
			MaxItems:    10,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The API response property which you want to filter your request on. Can be one of the following: \"/type\", \"/name\", \"/project/projectId\", \"/uuid\", \"/state\"",
						ValidateFunc: validation.StringInSlice([]string{"/type", "/name", "/project/projectId", "/uuid", "/state"}, true),
					},
					"operator": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Possible operators to use on the filter property. Can be one of the following: [ \"=\", \"!=\", \"[NOT] LIKE\", \"[NOT] IN\", \"ILIKE\" ]",
						ValidateFunc: validation.StringInSlice([]string{"=", "!=", "[NOT] LIKE", "[NOT] IN", "ILIKE"}, true),
					},
					"values": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"pagination": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Pagination details for the Data Source Search Request",
			MaxItems:    1,
			Elem:        paginationSchema(),
		},
		"sort": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Filters for the Data Source Search Request",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"direction": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "DESC",
						Description:  "The sorting direction. Can be one of: [DESC, ASC], Defaults to DESC",
						ValidateFunc: validation.StringInSlice([]string{"DESC", "ASC"}, true),
					},
					"property": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "/changeLog/updatedDateTime",
						Description:  "The property name to use in sorting. Can be one of the following: [/type, /uuid, /name, /project/projectId, /state, /notMatchedRuleAction, /connectionsCount, /changeLog/createdDateTime, /changeLog/updatedDateTime], Defaults to /changeLog/updatedDateTime",
						ValidateFunc: validation.StringInSlice([]string{"/type", "/uuid", "/name", "/project/projectId", "/state", "/notMatchedRuleAction", "/connectionsCount", "/changeLog/createdDateTime", "/changeLog/updatedDateTime"}, true),
					},
				},
			},
		},
	}
}
