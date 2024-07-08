package port

import (
	"github.com/equinix/terraform-provider-equinix/internal/network"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: resourceMetalPortRead,

		Schema: map[string]*schema.Schema{
			"port_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "UUID of the port to lookup",
				ConflictsWith: []string{"device_id", "name"},
			},
			"device_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Device UUID where to lookup the port",
				ConflictsWith: []string{"port_id"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Name of the port to look up, e.g. bond0, eth1",
				ConflictsWith: []string{"port_id"},
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "One of " + network.NetworkTypeListHB,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port type",
			},
			"mac": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MAC address of the port",
			},
			"bond_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the bond port",
			},
			"bond_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the bond port",
			},
			"bonded": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the port is bonded",
			},
			"disbond_supported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the port can be removed from a bond",
			},
			"native_vlan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of native VLAN of the port",
			},
			"vlan_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "UUIDs of attached VLANs",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vxlan_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "UUIDs of attached VLANs",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"layer2": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the port is in layer2 (or layer3) mode",
			},
		},
	}
}
