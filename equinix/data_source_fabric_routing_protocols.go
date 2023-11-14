package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FabricRoutingProtocolsDataResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Connection URI associated with Routing Protocol",
		},
		"direct_routing_protocol_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Direct Routing Protcol UUID",
		},
		"bgp_routing_protocol_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BGP Routing Protcol UUID",
		},
		"direct_routing_protocol": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Routing Protocol Direct Details",
			Elem: &schema.Resource{
				Schema: DirectRoutingProtocolSch(),
			},
		},
		"bgp_routing_protocol": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "BGP Routing Protocol Details",
			Elem: &schema.Resource{
				Schema: BGPRoutingProtocolSch(),
			},
		},
	}
}

func dataSourceRoutingProtocols() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoutingProtocolsRead,
		Schema:      FabricRoutingProtocolsDataResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch routing protocols for the given routing protocol ids\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func dataSourceRoutingProtocolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	directUUID, _ := d.Get("direct_routing_protocol_uuid").(string)
	bgpUUID, _ := d.Get("bgp_routing_protocol_uuid").(string)
	id := directUUID
	if bgpUUID != "" {
		id = directUUID + "/" + bgpUUID
	}

	d.SetId(id)
	return resourceFabricRoutingProtocolsRead(ctx, d, meta)
}
