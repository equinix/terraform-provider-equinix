package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMetalGateway() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceMetalGatewayRead,

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "UUID of the Metal Gateway to fetch",
			},
			"project_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "UUID of the Project where the Gateway is scoped to",
			},
			"vlan_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "UUID of the VLAN to associate",
			},
			"ip_reservation_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "UUID of the IP Reservation to associate, must be in the same metro as the VLAN",
			},
			"private_ipv4_subnet_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: fmt.Sprintf("Size of the private IPv4 subnet to create for this gateway, one of %v", subnetSizes),
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the virtual circuit resource",
			},
		},
	}
}

func dataSourceMetalGatewayRead(d *schema.ResourceData, meta interface{}) error {
	gatewayId, _ := d.Get("gateway_id").(string)

	d.SetId(gatewayId)
	return resourceMetalGatewayRead(d, meta)
}
