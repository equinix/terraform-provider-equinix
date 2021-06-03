package metal

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

var (
	l2Types = []string{"layer2-individual", "layer2-bonded"}
	l3Types = []string{"layer3", "hybrid", "hybrid-bonded"}
)

func resourceMetalPort() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Read: resourceMetalPortRead,
		// Create and Update are the same func
		Create: resourceMetalPortUpdate,
		Update: resourceMetalPortUpdate,
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"bonded": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Flag indicating whether the port should be bonded",
			},
			"layer2": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Flag indicating whether the port is in layer2 (or layer3) mode",
			},
			"native_vlan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "UUID of native VLAN of the port",
			},
			"vlan_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "UUIDs VLANs to attach",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "One of layer2-bonded, layer2-individual, layer3, hybrid and hybrid-bonded",
			},
			"disbond_supported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the port can be removed from a bond",
			},
			"bond_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the bond port",
			},
			"bond_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the bond port",
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
		},
	}
}

func resourceMetalPortUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	port, err := getPortByResourceData(d, client)
	if err != nil {
		return err
	}
	cpr := &ClientPortResource{
		Client:   client,
		Port:     port,
		Resource: d,
	}

	// Constraint: Only bond ports have layer2 mode
	_, l2Ok := d.GetOk("layer2")
	if port.Type != "NetworkBondPort" && l2Ok {
		return fmt.Errorf("layer2 flag can be specified only for bond ports")
	}

	err = removeNativeVlan(cpr)
	if err != nil {
		return err
	}

	err = removeVlans(cpr)
	if err != nil {
		return err
	}

	err = makeDisbond(cpr)
	if err != nil {
		return err
	}

	err = convertToL2(cpr)
	if err != nil {
		return err
	}

	err = convertToL3(cpr)
	if err != nil {
		return err
	}

	err = makeBond(cpr)
	if err != nil {
		return err
	}

	err = assignVlans(cpr)
	if err != nil {
		return err
	}

	err = assignNativeVlan(cpr)
	if err != nil {
		return err
	}
	return resourceMetalPortRead(d, meta)
}

func resourceMetalPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	port, err := getPortByResourceData(d, client)
	if err != nil {
		return err
	}
	m := map[string]interface{}{
		"layer2":            contains(l2Types, port.NetworkType),
		"port_id":           port.ID,
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
