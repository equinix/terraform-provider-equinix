package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalPort() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalPortRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "UUID of the port to lookup",
				ConflictsWith: []string{"device_id", "name"},
			},
			"device_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Device UUID where to lookup the port",
				ConflictsWith: []string{"id"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Name of the port to look up, e.g. bond0, eth1",
				ConflictsWith: []string{"id"},
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "One of " + NetworkTypeListHB,
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
		},
	}
}

func dataSourceMetalPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	portId, portIdOk := d.GetOk("id")
	deviceId, deviceIdOk := d.GetOk("device_id")
	portName, portNameOk := d.GetOk("name")

	var port *packngo.Port

	if portIdOk && (deviceIdOk || portNameOk) {
		return fmt.Errorf("You must specify either id or (device_id and name)")
	}
	if portIdOk {
		var err error
		port, _, err = client.Ports.Get(
			portId.(string),
			&packngo.GetOptions{Includes: []string{
				"native_virtual_network",
				"virtual_networks",
			}},
		)
		if err != nil {
			return err
		}
	} else {
		if !(deviceIdOk && portNameOk) {
			return fmt.Errorf("If you don't use port_id, you must supply both device_id and name")
		}
		device, _, err := client.Devices.Get(deviceId.(string), nil)
		if err != nil {
			return err
		}
		port, err = device.GetPortByName(portName.(string))
		if err != nil {
			return err
		}
	}
	m := map[string]interface{}{
		"type":              port.Type,
		"name":              port.Name,
		"network_type":      port.NetworkType,
		"mac":               port.Data.MAC,
		"bonded":            port.Data.Bonded,
		"disbond_supported": port.DisbondOperationSupported,
	}

	if port.NativeVirtualNetwork != nil {
		m["native_vlan_id"] = port.NativeVirtualNetwork.ID
	}

	if len(port.AttachedVirtualNetworks) > 0 {
		vlans := []string{}
		for _, n := range port.AttachedVirtualNetworks {
			vlans = append(vlans, n.ID)
		}
		m["vlan_ids"] = vlans
	}

	if port.Bond != nil {
		m["bond_id"] = port.Bond.ID
		m["bond_name"] = port.Bond.Name
	}

	d.SetId(port.ID)
	return setMap(d, m)
}
