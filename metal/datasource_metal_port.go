package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				ConflictsWith: []string{"device_id", "port_name"},
			},
			"device_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Device UUID where to lookup the port",
				ConflictsWith: []string{"id"},
			},
			"port_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Name of the port to look up, e.g. bond0, eth1",
				ConflictsWith: []string{"id"},
			},
			"bond_port": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether port is a bond port or physical port",
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
				// is the listing correct?
				Description: "One of layer2-bonded, layer2-individual, layer3, hybrid, hybrid-bonded",
			},
		},
	}
}

func dataSourceMetalPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	portId, portIdOk := d.GetOk("id")
	deviceId, deviceIdOk := d.GetOk("device_id")
	portName, portNameOk := d.GetOk("port_name")

	var port *packngo.Port

	if portIdOk && (deviceIdOk || portNameOk) {
		return fmt.Errorf("You must specify either id or (device_id and port_name)")
	}
	if portIdOk {
		var err error
		port, _, err = client.Ports.Get(portId.(string), nil)
		if err != nil {
			return err
		}
	} else {
		if !(deviceIdOk && portNameOk) {
			return fmt.Errorf("If you don't use port_id, you must supply both device_id and port_name")
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

	d.Set("bond_port", port.Type == "NetworkBondPort")
	d.Set("port_name", port.Name)
	d.Set("device_id", port.Name)
	d.Set("network_type", port.NetworkType)

	d.SetId(port.ID)

	return nil
}
