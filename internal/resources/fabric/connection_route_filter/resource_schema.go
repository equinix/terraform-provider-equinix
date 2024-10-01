package connection_route_filter

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Equinix Assigned UUID of the Equinix Connection to attach the Route Filter Policy to",
		},
		"route_filter_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Equinix Assigned UUID of the Route Filter Policy to attach to the Equinix Connection",
		},
		"direction": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"INBOUND", "OUTBOUND"}, false),
			Description:  "Direction of the filtering of the attached Route Filter Policy",
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
