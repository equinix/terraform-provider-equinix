package port

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/batch"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// Deprecated: empty port assignment input that is required
	// for some endpoints; probably indicates a bug in the API spec
	dummy = metalv1.PortAssignInput{}
)

type clientPortResource struct {
	Client   *metalv1.APIClient
	Port     *metalv1.Port
	Resource *schema.ResourceData
}

func getClientPortResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (*clientPortResource, *http.Response, error) {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	portID := d.Get("port_id").(string)

	getOpts := []string{
		"native_virtual_network",
		"virtual_networks",
	}
	port, resp, err := client.PortsApi.FindPortById(ctx, portID).Include(getOpts).Execute()
	if err != nil {
		return nil, resp, err
	}

	cpr := &clientPortResource{
		Client:   client,
		Port:     port,
		Resource: d,
	}
	return cpr, resp, nil
}

func getPortByResourceData(ctx context.Context, d *schema.ResourceData, client *metalv1.APIClient) (*metalv1.Port, *http.Response, error) {
	portID, portIDOk := d.GetOk("port_id")
	resourceID := d.Id()

	// rely on d.Id in imported resources
	if !portIDOk {
		if resourceID != "" {
			portID = resourceID
			portIDOk = true
		}
	}
	deviceID, deviceIDOk := d.GetOk("device_id")
	portName, portNameOk := d.GetOk("name")

	// check parameter sanity only for a new (not-yet-created) resource
	if resourceID == "" {
		if portIDOk && (deviceIDOk || portNameOk) {
			return nil, nil, fmt.Errorf("you must specify either id or (device_id and name)")
		}
	}

	var port *metalv1.Port
	var resp *http.Response
	var err error

	getOpts := []string{
		"native_virtual_network",
		"virtual_networks",
	}
	if portIDOk {
		port, resp, err = client.PortsApi.FindPortById(ctx, portID.(string)).Include(getOpts).Execute()
		if err != nil {
			return nil, resp, err
		}
	} else {
		if !deviceIDOk || !portNameOk {
			return nil, nil, fmt.Errorf("if you don't use port_id, you must supply both device_id and name")
		}
		var device *metalv1.Device
		device, resp, err = client.DevicesApi.FindDeviceById(ctx, deviceID.(string)).Include(getOpts).Execute()
		if err != nil {
			return nil, resp, err
		}
		port, err = getPortByName(device, portName.(string))
		return port, nil, err
	}

	return port, resp, nil
}

func getSpecifiedNative(d *schema.ResourceData) string {
	nativeRaw, nativeOk := d.GetOk("native_vlan_id")
	specifiedNative := ""
	if nativeOk {
		specifiedNative = nativeRaw.(string)
	}
	return specifiedNative
}

func getCurrentNative(p *metalv1.Port) string {
	currentNative := ""
	if p.NativeVirtualNetwork != nil {
		currentNative = p.NativeVirtualNetwork.GetId()
	}
	return currentNative
}

func attachedVlanIDs(p *metalv1.Port) []string {
	attached := []string{}
	for _, v := range p.VirtualNetworks {
		attached = append(attached, v.GetId())
	}
	return attached
}

func specifiedVlanIDs(d *schema.ResourceData) []string {
	// either vlan_ids or vxlan_ids should be set, TF should ensure that
	vlanIDsRaw, vlanIDsOk := d.GetOk("vlan_ids")
	if vlanIDsOk {
		return converters.IfArrToStringArr(vlanIDsRaw.(*schema.Set).List())
	}

	vxlanIDsRaw, vxlanIDsOk := d.GetOk("vxlan_ids")
	if vxlanIDsOk {
		return converters.IfArrToIntStringArr(vxlanIDsRaw.(*schema.Set).List())
	}
	return []string{}
}

func batchVlans(removeOnly bool) func(context.Context, *clientPortResource) error {
	return func(ctx context.Context, cpr *clientPortResource) error {
		var vlansToAssign []string
		var currentNative string
		vlansToRemove := converters.Difference(
			attachedVlanIDs(cpr.Port),
			specifiedVlanIDs(cpr.Resource),
		)
		if !removeOnly {
			currentNative = getCurrentNative(cpr.Port)

			vlansToAssign = converters.Difference(
				specifiedVlanIDs(cpr.Resource),
				attachedVlanIDs(cpr.Port),
			)
		}

		vlanBatch := batch.NewVlanBatch(cpr.Port.GetId())

		for _, v := range vlansToRemove {
			vlanBatch.RemoveAssignment(v)
		}

		for _, v := range vlansToAssign {
			if currentNative == v {
				vlanBatch.AddNativeAssignment(v)
			} else {
				vlanBatch.AddAssignment(v)
			}
		}

		_, _, err := vlanBatch.Execute(ctx, cpr.Client)
		return err
	}
}

func updateNativeVlan(ctx context.Context, cpr *clientPortResource) error {
	currentNative := getCurrentNative(cpr.Port)
	specifiedNative := getSpecifiedNative(cpr.Resource)

	if currentNative != specifiedNative {
		var port *metalv1.Port
		var err error
		if specifiedNative == "" && currentNative != "" {
			port, _, err = cpr.Client.PortsApi.DeleteNativeVlan(ctx, cpr.Port.GetId()).Execute()
		} else {
			port, _, err = cpr.Client.PortsApi.AssignNativeVlan(ctx, cpr.Port.GetId()).Vnid(specifiedNative).Execute()
		}
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}

func processBondAction(ctx context.Context, cpr *clientPortResource, actionIsBond bool) error {
	// There's no good alternative to GetOkExists until metal_port
	// is converted to terraform-plugin-framework
	// nolint:staticcheck
	wantsBondedRaw, wantsBondedOk := cpr.Resource.GetOkExists("bonded")
	wantsBonded := wantsBondedRaw.(bool)
	// only act if the necessary action is the one specified in doBond
	if wantsBondedOk && (wantsBonded == actionIsBond) {
		// act if the current Bond state of the port is different than the spcified
		if wantsBonded != cpr.Port.Data.GetBonded() {
			var port *metalv1.Port
			var err error
			if wantsBonded {
				port, _, err = cpr.Client.PortsApi.BondPort(ctx, cpr.Port.GetId()).Execute()
			} else {

				port, _, err = cpr.Client.PortsApi.DisbondPort(ctx, cpr.Port.GetId()).Execute()
			}

			if err != nil {
				return err
			}
			getOpts := []string{
				"native_virtual_network",
				"virtual_networks",
			}
			port, _, err = cpr.Client.PortsApi.FindPortById(ctx, port.GetId()).Include(getOpts).Execute()
			if err != nil {
				return err
			}

			*(cpr.Port) = *port
		}
	}
	return nil
}

func makeBond(ctx context.Context, cpr *clientPortResource) error {
	return processBondAction(ctx, cpr, true)
}

func makeDisbond(ctx context.Context, cpr *clientPortResource) error {
	return processBondAction(ctx, cpr, false)
}

func convertToL2(ctx context.Context, cpr *clientPortResource) error {
	// There's no good alternative to GetOkExists until metal_port
	// is converted to terraform-plugin-framework
	// nolint:staticcheck
	l2, l2Ok := cpr.Resource.GetOkExists("layer2")
	isLayer2 := slices.Contains(l2Types, cpr.Port.GetNetworkType())

	if l2Ok && l2.(bool) && !isLayer2 {
		port, _, err := cpr.Client.PortsApi.ConvertLayer2(ctx, cpr.Port.GetId()).PortAssignInput(dummy).Execute()
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}

func convertToL3(ctx context.Context, cpr *clientPortResource) error {
	// There's no good alternative to GetOkExists until metal_port
	// is converted to terraform-plugin-framework
	// nolint:staticcheck
	l2, l2Ok := cpr.Resource.GetOkExists("layer2")
	isLayer2 := slices.Contains(l2Types, cpr.Port.GetNetworkType())

	if l2Ok && !l2.(bool) && isLayer2 {
		ips := metalv1.PortConvertLayer3Input{
			RequestIps: []metalv1.PortConvertLayer3InputRequestIpsInner{
				{AddressFamily: metalv1.PtrInt32(4), Public: metalv1.PtrBool(true)},
				{AddressFamily: metalv1.PtrInt32(4), Public: metalv1.PtrBool(false)},
				{AddressFamily: metalv1.PtrInt32(6), Public: metalv1.PtrBool(true)},
			},
		}

		port, _, err := cpr.Client.PortsApi.ConvertLayer3(ctx, cpr.Port.GetId()).PortConvertLayer3Input(ips).Execute()
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}

func portSanityChecks(_ context.Context, cpr *clientPortResource) error {
	isBondPort := cpr.Port.GetType() == "NetworkBondPort"

	// Constraint: Only bond ports have layer2 mode
	// There's no good alternative to GetOkExists until metal_port
	// is converted to terraform-plugin-framework
	// nolint:staticcheck
	l2Raw, l2Ok := cpr.Resource.GetOkExists("layer2")
	if !isBondPort && l2Ok {
		return fmt.Errorf("layer2 flag can be set only for bond ports")
	}

	l2 := l2Raw.(bool)

	bonded := cpr.Resource.Get("bonded").(bool)

	// Constraint: L3 unbonded is not really allowed for Bond port
	if isBondPort && !l2 && !bonded {
		return fmt.Errorf("bond port in Layer3 can't be unbonded")
	}

	// Constraint: native vlan ..
	// - must be one of assigned vlans
	// - there must be more than one vlan assigned to the port
	nativeVlanRaw, nativeVlanOk := cpr.Resource.GetOk("native_vlan_id")
	if nativeVlanOk {
		nativeVlan := nativeVlanRaw.(string)
		vlans := specifiedVlanIDs(cpr.Resource)
		if !slices.Contains(vlans, nativeVlan) {
			return fmt.Errorf("the native VLAN to be set is not (being) assigned to the port")
		}
		if len(vlans) < 2 {
			return fmt.Errorf("native VLAN can only be set if more than one VLAN are assigned to the port ")
		}
	}

	return nil
}

func getPortByName(d *metalv1.Device, name string) (*metalv1.Port, error) {
	for _, port := range d.NetworkPorts {
		if port.GetName() == name {
			return &port, nil
		}
	}
	return nil, fmt.Errorf("port %s not found in device %s", name, d.GetId())
}
