package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

var (
	l2Types = []string{"layer2-individual", "layer2-bonded"}
	l3Types = []string{"layer3", "hybrid", "hybrid-bonded"}
)

func resourceMetalPort() *schema.Resource {
	return &schema.Resource{
		Read:   resourceMetalPortRead,
		Create: resourceMetalPortCreate,

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

func getPortByResourceData(d *schema.ResourceData, client *packngo.Client) (*packngo.Port, error) {
	portId, portIdOk := d.GetOk("id")
	deviceId, deviceIdOk := d.GetOk("device_id")
	portName, portNameOk := d.GetOk("name")

	var port *packngo.Port

	if portIdOk && (deviceIdOk || portNameOk) {
		return nil, fmt.Errorf("You must specify either id or (device_id and name)")
	}
	getOpts := &packngo.GetOptions{Includes: []string{
		"native_virtual_network",
		"virtual_networks",
	}}
	if portIdOk {
		var err error
		port, _, err = client.Ports.Get(portId.(string), getOpts)
		if err != nil {
			return nil, err
		}
	} else {
		if !(deviceIdOk && portNameOk) {
			return nil, fmt.Errorf("If you don't use port_id, you must supply both device_id and name")
		}
		device, _, err := client.Devices.Get(deviceId.(string), getOpts)
		if err != nil {
			return nil, err
		}
		port, err = device.GetPortByName(portName.(string))
		if err != nil {
			return nil, err
		}
	}

	return port, nil
}

func getSpecifiedNative(d *schema.ResourceData) string {
	nativeRaw, nativeOk := d.GetOk("native_vlan_id")
	specifiedNative := ""
	if nativeOk {
		specifiedNative = nativeRaw.(string)
	}
	return specifiedNative
}

func getCurrentNative(p *packngo.Port) string {
	currentNative := ""
	if p.NativeVirtualNetwork != nil {
		currentNative = p.NativeVirtualNetwork.ID
	}
	return currentNative
}

func attachedVlanIds(p *packngo.Port) []string {
	attached := []string{}
	for _, v := range p.AttachedVirtualNetworks {
		attached = append(attached, v.ID)
	}
	return attached
}

func specifiedVlanIds(d *schema.ResourceData) []string {
	vlanIdsRaw, vlanIdsOk := d.GetOk("vlan_ids")
	specified := []string{}
	if vlanIdsOk {
		specified = convertStringArr(vlanIdsRaw.([]interface{}))
	}
	return specified
}

type portVlanAction func(*packngo.PortAssignRequest) (*packngo.Port, *packngo.Response, error)

func processVlansOnPort(portId string, vlanIds []string, f portVlanAction) error {
	par := packngo.PortAssignRequest{PortID: portId}
	for _, vId := range vlanIds {
		par.VirtualNetworkID = vId
		_, _, err := f(&par)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Create and Update will be probably identical in case of metal_port

I figured we should process port action in order:
 - native VLAN detachment
 - VLANs detachment
 - bond break
 - l2/l3 conversion
 - bond creation
 - VLANs assignment
 - native VLAN assignment
*/

func resourceMetalPortCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	port, err := getPortByResourceData(d, client)
	if err != nil {
		return err
	}

	// Constraint: Only bond ports have layer2 mode
	l2, l2Ok := d.GetOk("layer2")
	if port.Type != "NetworkBondPort" && l2Ok {
		return fmt.Errorf("layer2 flag can be specified only for bond ports")
	}

	// remove native vlan
	currentNative := getCurrentNative(port)
	specifiedNative := getSpecifiedNative(d)
	if (currentNative != specifiedNative) && currentNative != "" {
		_, _, err = client.DevicePorts.UnassignNative(port.ID)
		if err != nil {
			return err
		}
	}

	// remove vlans
	vlansToRemove := difference(
		attachedVlanIds(port),
		specifiedVlanIds(d),
	)
	err = processVlansOnPort(port.ID, vlansToRemove, client.DevicePorts.Unassign)
	if err != nil {
		return err
	}

	// disbond
	bonded, bondedOk := d.GetOk("bonded")
	if bondedOk && !bonded.(bool) && port.Data.Bonded {
		port, _, err = client.DevicePorts.Disbond(port, false)
		if err != nil {
			return err
		}
	}

	// convert to layer2
	isLayer2 := contains(l2Types, port.NetworkType)

	if l2Ok && l2.(bool) && !isLayer2 {
		port, _, err = client.Ports.ConvertToLayerTwo(port.ID)
		if err != nil {
			return err
		}
	}

	// convert to layer3
	if l2Ok && !l2.(bool) && isLayer2 {
		ips := []packngo.AddressRequest{
			{AddressFamily: 4, Public: true},
			{AddressFamily: 4, Public: false},
			{AddressFamily: 6, Public: true},
		}
		port, _, err = client.Ports.ConvertToLayerThree(port.ID, ips)
		if err != nil {
			return err
		}
	}

	// constructive phase - bond creation and vlan assignment
	// ======================================================

	port, err = getPortByResourceData(d, client)
	if err != nil {
		return err
	}

	// bond
	if bondedOk && bonded.(bool) && !port.Data.Bonded {
		port, _, err = client.DevicePorts.Bond(port, false)
		if err != nil {
			return err
		}
	}

	// assign VLANs
	vlansToAssign := difference(
		specifiedVlanIds(d),
		attachedVlanIds(port),
	)
	err = processVlansOnPort(port.ID, vlansToAssign, client.DevicePorts.Assign)
	if err != nil {
		return err
	}

	// assign Native VLAN
	currentNative = getCurrentNative(port)
	specifiedNative = getSpecifiedNative(d)
	if (currentNative != specifiedNative) && currentNative != "" {
		par := packngo.PortAssignRequest{
			PortID:           port.ID,
			VirtualNetworkID: specifiedNative,
		}
		_, _, err = client.DevicePorts.AssignNative(&par)
		if err != nil {
			return err
		}
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
