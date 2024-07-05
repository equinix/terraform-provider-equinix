package equinix

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/converters"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
	"github.com/pkg/errors"
)

type ClientPortResource struct {
	Client   *packngo.Client
	Port     *packngo.Port
	Resource *schema.ResourceData
}

func getClientPortResource(d *schema.ResourceData, meta interface{}) (*ClientPortResource, *packngo.Response, error) {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	port_id := d.Get("port_id").(string)

	getOpts := &packngo.GetOptions{Includes: []string{
		"native_virtual_network",
		"virtual_networks",
	}}
	port, resp, err := client.Ports.Get(port_id, getOpts)
	if err != nil {
		return nil, resp, err
	}

	cpr := &ClientPortResource{
		Client:   client,
		Port:     port,
		Resource: d,
	}
	return cpr, resp, nil
}

func getPortByResourceData(d *schema.ResourceData, client *packngo.Client) (*packngo.Port, error) {
	portId, portIdOk := d.GetOk("port_id")
	resourceId := d.Id()

	// rely on d.Id in imported resources
	if !portIdOk {
		if resourceId != "" {
			portId = resourceId
			portIdOk = true
		}
	}
	deviceId, deviceIdOk := d.GetOk("device_id")
	portName, portNameOk := d.GetOk("name")

	// check parameter sanity only for a new (not-yet-created) resource
	if resourceId == "" {
		if portIdOk && (deviceIdOk || portNameOk) {
			return nil, fmt.Errorf("you must specify either id or (device_id and name)")
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
	// either vlan_ids or vxlan_ids should be set, TF should ensure that
	vlanIdsRaw, vlanIdsOk := d.GetOk("vlan_ids")
	if vlanIdsOk {
		return converters.IfArrToStringArr(vlanIdsRaw.(*schema.Set).List())
	}

	vxlanIdsRaw, vxlanIdsOk := d.GetOk("vxlan_ids")
	if vxlanIdsOk {
		return converters.IfArrToIntStringArr(vxlanIdsRaw.(*schema.Set).List())
	}
	return []string{}
}

func batchVlans(ctx context.Context, start time.Time, removeOnly bool) func(*ClientPortResource) error {
	return func(cpr *ClientPortResource) error {
		var vlansToAssign []string
		var currentNative string
		vlansToRemove := converters.Difference(
			attachedVlanIds(cpr.Port),
			specifiedVlanIds(cpr.Resource),
		)
		if !removeOnly {
			currentNative = getCurrentNative(cpr.Port)

			vlansToAssign = converters.Difference(
				specifiedVlanIds(cpr.Resource),
				attachedVlanIds(cpr.Port),
			)
		}
		vacr := &packngo.VLANAssignmentBatchCreateRequest{}
		for _, v := range vlansToRemove {
			vacr.VLANAssignments = append(vacr.VLANAssignments, packngo.VLANAssignmentCreateRequest{
				VLAN:  v,
				State: packngo.VLANAssignmentUnassigned,
			})
		}

		for _, v := range vlansToAssign {
			native := currentNative == v
			vacr.VLANAssignments = append(vacr.VLANAssignments, packngo.VLANAssignmentCreateRequest{
				VLAN:   v,
				State:  packngo.VLANAssignmentAssigned,
				Native: &native,
			})
		}
		return createAndWaitForBatch(ctx, start, cpr, vacr)
	}
}

func createAndWaitForBatch(ctx context.Context, start time.Time, cpr *ClientPortResource, vacr *packngo.VLANAssignmentBatchCreateRequest) error {
	if len(vacr.VLANAssignments) == 0 {
		return nil
	}

	portID := cpr.Port.ID
	c := cpr.Client

	b, _, err := c.VLANAssignments.CreateBatch(portID, vacr, nil)
	if err != nil {
		return fmt.Errorf("vlan assignment batch could not be created: %w", err)
	}

	deadline, _ := ctx.Deadline()
	// originally set timeout in ctx by TF
	ctxTimeout := deadline.Sub(start)

	stateChangeConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{string(packngo.VLANAssignmentBatchQueued), string(packngo.VLANAssignmentBatchInProgress)},
		Target:     []string{string(packngo.VLANAssignmentBatchCompleted)},
		MinTimeout: 5 * time.Second,
		Timeout:    ctxTimeout - time.Since(start) - 30*time.Second,
		Refresh: func() (result interface{}, state string, err error) {
			b, _, err := c.VLANAssignments.GetBatch(portID, b.ID, nil)
			switch b.State {
			case packngo.VLANAssignmentBatchFailed:
				return b, string(packngo.VLANAssignmentBatchFailed),
					fmt.Errorf("vlan assignment batch %s provisioning failed: %s", b.ID, strings.Join(b.ErrorMessages, "; "))
			case packngo.VLANAssignmentBatchCompleted:
				return b, string(packngo.VLANAssignmentBatchCompleted), nil
			default:
				if err != nil {
					return b, "", fmt.Errorf("vlan assignment batch %s could not be polled: %w", b.ID, err)
				}
				return b, string(b.State), err
			}
		},
	}
	if _, err = stateChangeConf.WaitForStateContext(ctx); err != nil {
		return errors.Wrapf(err, "vlan assignment batch %s is not complete after timeout", b.ID)
	}
	return nil
}

func updateNativeVlan(cpr *ClientPortResource) error {
	currentNative := getCurrentNative(cpr.Port)
	specifiedNative := getSpecifiedNative(cpr.Resource)

	if currentNative != specifiedNative {
		var port *packngo.Port
		var err error
		if specifiedNative == "" && currentNative != "" {
			port, _, err = cpr.Client.Ports.UnassignNative(cpr.Port.ID)
		} else {
			port, _, err = cpr.Client.Ports.AssignNative(cpr.Port.ID, specifiedNative)
		}
		if err != nil {
			return err
		}
		*(cpr.Port) = *port
	}
	return nil
}

func processBondAction(cpr *ClientPortResource, actionIsBond bool) error {
	wantsBondedRaw, wantsBondedOk := cpr.Resource.GetOkExists("bonded")
	wantsBonded := wantsBondedRaw.(bool)
	// only act if the necessary action is the one specified in doBond
	if wantsBondedOk && (wantsBonded == actionIsBond) {
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
			getOpts := &packngo.GetOptions{Includes: []string{
				"native_virtual_network",
				"virtual_networks",
			}}
			port, _, err = cpr.Client.Ports.Get(port.ID, getOpts)
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
	l2, l2Ok := cpr.Resource.GetOkExists("layer2")
	isLayer2 := slices.Contains(l2Types, cpr.Port.NetworkType)

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
	l2, l2Ok := cpr.Resource.GetOkExists("layer2")
	isLayer2 := slices.Contains(l2Types, cpr.Port.NetworkType)

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

func portSanityChecks(cpr *ClientPortResource) error {
	isBondPort := cpr.Port.Type == "NetworkBondPort"

	// Constraint: Only bond ports have layer2 mode
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
		vlans := specifiedVlanIds(cpr.Resource)
		if !slices.Contains(vlans, nativeVlan) {
			return fmt.Errorf("the native VLAN to be set is not (being) assigned to the port")
		}
		if len(vlans) < 2 {
			return fmt.Errorf("native VLAN can only be set if more than one VLAN are assigned to the port ")
		}
	}

	return nil
}

func portProperlyDestroyed(port *packngo.Port) error {
	var errs []string
	if !port.Data.Bonded {
		errs = append(errs, fmt.Sprintf("port %s wasn't bonded after equinix_metal_port destroy;", port.ID))
	}
	if port.Type == "NetworkBondPort" && port.NetworkType != "layer3" {
		errs = append(errs, "bond port should be in layer3 type after destroy;")
	}
	if port.NativeVirtualNetwork != nil {
		errs = append(errs, "port should not have native VLAN assigned after destroy;")
	}
	if len(port.AttachedVirtualNetworks) != 0 {
		errs = append(errs, "port should not have VLANs attached after destroy")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", errs)
	}

	return nil
}
