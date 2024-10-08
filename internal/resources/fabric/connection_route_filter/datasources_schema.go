package connection_route_filter

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceByUUIDSchema() map[string]*schema.Schema {
	dsSchema := baseSchema()
	dsSchema["connection_id"] = connectionIdSchema()
	dsSchema["route_filter_id"] = routeFilterIdSchema()
	return dsSchema
}

func dataSourceAllFiltersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_id": connectionIdSchema(),
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
				Schema: baseSchema(),
			},
		},
	}
}

func baseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"direction": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Direction of the filtering of the attached Route Filter Policy",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Route Filter Type. One of [ \"BGP_IPv4_PREFIX_FILTER\", \"BGP_IPv6_PREFIX_FILTER\" ] ",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "URI to the attached Route Filter Policy on the Connection",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix Assigned ID for Route Filter Policy",
		},
		"attachment_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the Route Filter Policy attachment lifecycle",
		},
	}
}

func routeFilterIdSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Equinix Assigned UUID of the Route Filter Policy to attach to the Equinix Connection",
	}
}

func connectionIdSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Equinix Assigned UUID of the Equinix Connection to attach the Route Filter Policy to",
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
