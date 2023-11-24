package metal_port

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/packethost/packngo"
	"github.com/pkg/errors"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
)

type ClientPortData struct {
	Client	*packngo.Client
	Port	*packngo.Port
	Data	MetalPortResourceModel
}

func updatePort(ctx context.Context, client *packngo.Client, plan MetalPortResourceModel) error {
	start := time.Now()
	cpd, _, err := getPortData(client, plan)
	if err != nil {
		return helper.FriendlyError(err)
	}

	for _, f := range [](func(*ClientPortData) error){
		portSanityChecks(ctx),
		batchVlans(ctx, start, true),
		makeDisbond,
		convertToL2,
		makeBond,
		convertToL3,
		batchVlans(ctx, start, false),
		updateNativeVlan,
	} {
		if err := f(cpd); err != nil {
			return helper.FriendlyError(err)
		}
	}

	return nil
}

func getPortData(client *packngo.Client, data MetalPortResourceModel) (*ClientPortData, *packngo.Response, error) {
	getOpts := &packngo.GetOptions{Includes: []string{
		"native_virtual_network",
		"virtual_networks",
	}}
	port, resp, err := client.Ports.Get(data.PortID.ValueString(), getOpts)
	if err != nil {
		return nil, resp, helper.FriendlyError(err)
	}

	cpd := &ClientPortData{
		Client: client,
		Port:   port,
		Data:   data,
	}
	return cpd, resp, nil
}

func getPortByResourceData(d MetalPortResourceModel, client *packngo.Client) (*packngo.Port, error) {
	portId := d.PortID
	resourceId := d.ID

	// rely on d.Id in imported resources
	if portId.IsNull() {
		if !resourceId.IsNull() {
			portId = resourceId
		}
	}
	deviceId := d.DeviceID
	portName := d.Name

	// check parameter sanity only for a new (not-yet-created) resource
	if resourceId.IsNull() {
		if !portId.IsNull() && (!deviceId.IsNull() || !portName.IsNull()) {
			return nil, fmt.Errorf("you must specify either port_id or (device_id and name)")
		}
	}

	var port *packngo.Port

	getOpts := &packngo.GetOptions{Includes: []string{
		"native_virtual_network",
		"virtual_networks",
	}}
	if !portId.IsNull() {
		var err error
		port, _, err = client.Ports.Get(portId.ValueString(), getOpts)
		if err != nil {
			return nil, err
		}
	} else {
		if deviceId.IsNull() && portName.IsNull() {
			return nil, fmt.Errorf("if you don't use port_id, you must supply both device_id and name")
		}
		device, _, err := client.Devices.Get(deviceId.ValueString(), getOpts)
		if err != nil {
			return nil, err
		}
		return device.GetPortByName(portName.ValueString())
	}

	return port, nil
}

func getSpecifiedNative(d MetalPortResourceModel) string {
	specifiedNative := ""
	if !d.NativeVLANID.IsNull() {
		specifiedNative = d.NativeVLANID.ValueString()
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

func specifiedVlanIds(ctx context.Context, d MetalPortResourceModel) ([]string, error) {
	if !d.VlanIDs.IsNull() {
		var ids []string
		diags := d.VlanIDs.ElementsAs(ctx, &ids, false)
		if diags.HasError(){
			return nil, fmt.Errorf("failed to validate vlan IDs: %s", diags.Errors())
			
		}
	}

	if !d.VxlanIDs.IsNull() {
		var ids []string
		diags := d.VxlanIDs.ElementsAs(ctx, &ids, false)
		if diags.HasError(){
			return nil, fmt.Errorf("failed to validate vxlan IDs: %s", diags.Errors())
		}
	}

	return []string{}, nil
}

func batchVlans(ctx context.Context, start time.Time, removeOnly bool) func(*ClientPortData) error {
	return func(cpd *ClientPortData) error {
		var vlansToAssign []string
		var currentNative string
		specifiedVlanIds, err := specifiedVlanIds(ctx, cpd.Data)
		if err != nil {
			return err
		}
		vlansToRemove := helper.Difference(
			attachedVlanIds(cpd.Port),
			specifiedVlanIds,
		)
		if !removeOnly {
			currentNative = getCurrentNative(cpd.Port)

			vlansToAssign = helper.Difference(
				specifiedVlanIds,
				attachedVlanIds(cpd.Port),
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
		return createAndWaitForBatch(ctx, start, cpd, vacr)
	}
}

func createAndWaitForBatch(ctx context.Context, start time.Time, cpd *ClientPortData, vacr *packngo.VLANAssignmentBatchCreateRequest) error {
	if len(vacr.VLANAssignments) == 0 {
		return nil
	}

	portID := cpd.Port.ID
	c := cpd.Client

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

func updateNativeVlan(cpd *ClientPortData) error {
	currentNative := getCurrentNative(cpd.Port)
	specifiedNative := getSpecifiedNative(cpd.Data)

	if currentNative != specifiedNative {
		var port *packngo.Port
		var err error
		if specifiedNative == "" && currentNative != "" {
			port, _, err = cpd.Client.Ports.UnassignNative(cpd.Port.ID)
		} else {
			port, _, err = cpd.Client.Ports.AssignNative(cpd.Port.ID, specifiedNative)
		}
		if err != nil {
			return err
		}
		*(cpd.Port) = *port
	}
	return nil
}

func processBondAction(cpd *ClientPortData, actionIsBond bool) error {
	wantsBonded := cpd.Data.Bonded.ValueBool()
	// only act if the necessary action is the one specified in doBond
	if wantsBonded == actionIsBond {
		// act if the current Bond state of the port is different than the spcified
		if wantsBonded != cpd.Port.Data.Bonded {
			action := cpd.Client.Ports.Disbond
			if wantsBonded {
				action = cpd.Client.Ports.Bond
			}

			port, _, err := action(cpd.Port.ID, false)
			if err != nil {
				return err
			}
			getOpts := &packngo.GetOptions{Includes: []string{
				"native_virtual_network",
				"virtual_networks",
			}}
			port, _, err = cpd.Client.Ports.Get(port.ID, getOpts)
			if err != nil {
				return err
			}

			*(cpd.Port) = *port
		}
	}
	return nil
}

func makeBond(cpd *ClientPortData) error {
	return processBondAction(cpd, true)
}

func makeDisbond(cpd *ClientPortData) error {
	return processBondAction(cpd, false)
}

func convertToL2(cpd *ClientPortData) error {
	l2 := cpd.Data.Layer2
	isLayer2 := helper.Contains(l2Types, cpd.Port.NetworkType)

	if l2.ValueBool() && !isLayer2 {
		port, _, err := cpd.Client.Ports.ConvertToLayerTwo(cpd.Port.ID)
		if err != nil {
			return err
		}
		*(cpd.Port) = *port
	}
	return nil
}

func convertToL3(cpd *ClientPortData) error {
	l2 := cpd.Data.Layer2
	isLayer2 := helper.Contains(l2Types, cpd.Port.NetworkType)

	if !l2.ValueBool() && isLayer2 {
		ips := []packngo.AddressRequest{
			{AddressFamily: 4, Public: true},
			{AddressFamily: 4, Public: false},
			{AddressFamily: 6, Public: true},
		}
		port, _, err := cpd.Client.Ports.ConvertToLayerThree(cpd.Port.ID, ips)
		if err != nil {
			return err
		}
		*(cpd.Port) = *port
	}
	return nil
}

func portSanityChecks(ctx context.Context) func(*ClientPortData) error {
	return func(cpd *ClientPortData) error {
		isBondPort := cpd.Port.Type == "NetworkBondPort"

		// Constraint: Only bond ports have layer2 mode
		l2 := cpd.Data.Layer2.ValueBool()
		if !isBondPort && l2 {
			return fmt.Errorf("layer2 flag can be set only for bond ports")
		}

		bonded := cpd.Data.Bonded.ValueBool()

		// Constraint: L3 unbonded is not really allowed for Bond port
		if isBondPort && !l2 && !bonded {
			return fmt.Errorf("bond port in Layer3 can't be unbonded")
		}

		// Constraint: native vlan ..
		// - must be one of assigned vlans
		// - there must be more than one vlan assigned to the port
		nativeVlan := cpd.Data.NativeVLANID
		if !nativeVlan.IsNull() {
			vlans, err := specifiedVlanIds(ctx, cpd.Data)
			if err != nil {
				return err
			}
			if !helper.Contains(vlans, nativeVlan.ValueString()) {
				return fmt.Errorf("the native VLAN to be set is not (being) assigned to the port")
			}
			if len(vlans) < 2 {
				return fmt.Errorf("native VLAN can only be set if more than one VLAN are assigned to the port ")
			}
		}

		return nil
	}
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
