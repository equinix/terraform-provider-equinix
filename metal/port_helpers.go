package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

type portVlanAction func(*packngo.PortAssignRequest) (*packngo.Port, *packngo.Response, error)

type ClientPortResource struct {
	Client   *packngo.Client
	Port     *packngo.Port
	Resource *schema.ResourceData
}

func getPortByResourceData(d *schema.ResourceData, client *packngo.Client) (*packngo.Port, error) {
	portId, portIdOk := d.GetOk("port_id")
	deviceId, deviceIdOk := d.GetOk("device_id")
	portName, portNameOk := d.GetOk("name")

	// check parameter sanity only for a new (not-yet-created) resource
	if d.Id() == "" {
		if portIdOk && (deviceIdOk || portNameOk) {
			return nil, fmt.Errorf("You must specify either id or (device_id and name)")
		}
	}

	var port *packngo.Port

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
		return device.GetPortByName(portName.(string))
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

func processVlansOnPort(port *packngo.Port, vlanIds []string, f portVlanAction) (*packngo.Port, error) {
	par := packngo.PortAssignRequest{PortID: port.ID}
	for _, vId := range vlanIds {
		par.VirtualNetworkID = vId
		var err error
		port, _, err = f(&par)
		if err != nil {
			return nil, err
		}
	}
	return port, nil
}

func removeVlans(cpr *ClientPortResource) error {
	vlansToRemove := difference(
		attachedVlanIds(cpr.Port),
		specifiedVlanIds(cpr.Resource),
	)
	port, err := processVlansOnPort(cpr.Port, vlansToRemove, cpr.Client.DevicePorts.Unassign)
	if err != nil {
		return err
	}
	*(cpr.Port) = *port
	return nil
}

func assignVlans(cpr *ClientPortResource) error {
	// assign VLANs
	vlansToAssign := difference(
		specifiedVlanIds(cpr.Resource),
		attachedVlanIds(cpr.Port),
	)
	port, err := processVlansOnPort(cpr.Port, vlansToAssign, cpr.Client.DevicePorts.Assign)
	if err != nil {
		return err
	}
	*(cpr.Port) = *port
	return nil
}

func removeNativeVlan(cpr *ClientPortResource) error {
	currentNative := getCurrentNative(cpr.Port)
	specifiedNative := getSpecifiedNative(cpr.Resource)
	if (currentNative != specifiedNative) && currentNative != "" {
		port, _, err := cpr.Client.DevicePorts.UnassignNative(cpr.Port.ID)
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}

func assignNativeVlan(cpr *ClientPortResource) error {
	// assign Native VLAN
	currentNative := getCurrentNative(cpr.Port)
	specifiedNative := getSpecifiedNative(cpr.Resource)
	if (currentNative != specifiedNative) && currentNative != "" {
		par := packngo.PortAssignRequest{
			PortID:           cpr.Port.ID,
			VirtualNetworkID: specifiedNative,
		}
		port, _, err := cpr.Client.DevicePorts.AssignNative(&par)
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}

func processBondAction(cpr *ClientPortResource, actionIsBond bool) error {
	wantsBondedRaw, bondedSpecified := cpr.Resource.GetOk("bonded")
	wantsBonded := wantsBondedRaw.(bool)

	// only act if the necessary action is the one specified in doBond
	if bondedSpecified && (wantsBonded == actionIsBond) {
		// act if the current Bond state of the port is different than the spcified
		if wantsBonded != cpr.Port.Data.Bonded {
			action := cpr.Client.DevicePorts.Disbond
			if wantsBonded {
				action = cpr.Client.DevicePorts.Bond
			}
			port, _, err := action(cpr.Port, false)
			if err != nil {
				return err
			}
			*(cpr.Port) = *port
		}
	}
	return nil
}

func makeBond(cpr *ClientPortResource) error {
	return processBondAction(cpr, true)
}

func makeDisbond(cpr *ClientPortResource) error {
	return processBondAction(cpr, false)
}

func convertToL2(cpr *ClientPortResource) error {
	l2, l2Ok := cpr.Resource.GetOk("layer2")
	isLayer2 := contains(l2Types, cpr.Port.NetworkType)

	if l2Ok && l2.(bool) && !isLayer2 {
		port, _, err := cpr.Client.Ports.ConvertToLayerTwo(cpr.Port.ID)
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}

func convertToL3(cpr *ClientPortResource) error {
	l2, l2Ok := cpr.Resource.GetOk("layer2")
	isLayer2 := contains(l2Types, cpr.Port.NetworkType)
	if l2Ok && !l2.(bool) && isLayer2 {
		ips := []packngo.AddressRequest{
			{AddressFamily: 4, Public: true},
			{AddressFamily: 4, Public: false},
			{AddressFamily: 6, Public: true},
		}
		port, _, err := cpr.Client.Ports.ConvertToLayerThree(cpr.Port.ID, ips)
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}
