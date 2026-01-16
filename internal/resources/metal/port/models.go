// Package port represents network ports for instances.
package port

import (
	"context"
	"slices"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceModel struct {
	// Basic Port Attributes
	PortID           types.String `tfsdk:"port_id"`
	Type             types.String `tfsdk:"type"`
	Name             types.String `tfsdk:"name"`
	NetworkType      types.String `tfsdk:"network_type"`
	MAC              types.String `tfsdk:"mac"`
	Bonded           types.Bool   `tfsdk:"bonded"`
	DisbondSupported types.Bool   `tfsdk:"disbond_supported"`

	// Reduction of Network Types to boolean L2/L3
	Layer2 types.Bool `tfsdk:"layer2"`

	// Only one VLAN ID can be native at a time, this is the VLAN for untagged traffic from the port.
	NativeVlanID types.String `tfsdk:"native_vlan_id"`

	// VXLAN IDs are the metro local VXLAN tags e.g. (2-4096)
	VXLANIDs types.Set `tfsdk:"vxlan_ids"`

	// VLAN IDs are the UUID of the Virtual Network resource
	VLANIDs types.Set `tfsdk:"vlan_ids"`

	// BondName is the name of bond if this port is a member of a bond
	BondName types.String `tfsdk:"bond_name"`

	// BondID is the UUID for the Bond port resource
	BondID types.String `tfsdk:"bond_id"`

	// TODO: define this
	ResetOnDelete types.Bool `tfsdk:"reset_on_delete"`
}

func (m *resourceModel) ToExecutionPlan(_ context.Context, port *metalv1.Port) ([]string, diag.Diagnostics) {
	plannedOperations := []string{}
	if port.Data.GetBonded() != m.Bonded.ValueBool() {
		if m.Bonded.ValueBool() {
			plannedOperations = append(plannedOperations, "bond")
		} else {
			plannedOperations = append(plannedOperations, "disbond")
		}
	}

	if !m.Layer2.IsNull() {
		alreadyL2 := slices.Contains(l2Types, port.GetNetworkType())
		wantsL2 := m.Layer2.ValueBool()
		wantsL3 := !m.Layer2.ValueBool()

		if alreadyL2 && wantsL3 {
			plannedOperations = append(plannedOperations, "toLayer3")
		}

		if wantsL2 && !alreadyL2 {
			plannedOperations = append(plannedOperations, "toLayer2")
		}

	}

	currentNativeVlan := port.GetNativeVirtualNetwork()
	desiredNativeVlan := m.NativeVlanID

	if desiredNativeVlan.IsNull() {
		if currentNativeVlan.GetId() != "" {
			plannedOperations = append(plannedOperations, "removeNativeVlan")
		}
	} else {
		if currentNativeVlan.GetId() != desiredNativeVlan.ValueString() {
			plannedOperations = append(plannedOperations, "assignNativeVlan")
		}
	}

	return plannedOperations, nil
}

func (m *resourceModel) parse(ctx context.Context, port *metalv1.Port) diag.Diagnostics {
	var diags diag.Diagnostics

	m.PortID = types.StringValue(port.GetId())
	m.Type = types.StringValue(string(port.GetType()))
	m.Name = types.StringValue(port.GetName())
	m.NetworkType = types.StringValue(string(port.GetNetworkType()))
	m.MAC = types.StringValue(port.Data.GetMac())
	m.Bonded = types.BoolValue(port.Data.GetBonded())
	m.DisbondSupported = types.BoolValue(port.GetDisbondOperationSupported())

	l2 := slices.Contains(l2Types, port.GetNetworkType())
	l3 := slices.Contains(l3Types, port.GetNetworkType())

	if l2 {
		m.Layer2 = types.BoolValue(true)
	}
	if l3 {
		m.Layer2 = types.BoolValue(false)
	}

	if port.NativeVirtualNetwork != nil {
		m.NativeVlanID = types.StringValue(port.NativeVirtualNetwork.GetId())
	}

	vlans := []string{}
	vxlans := []int{}

	for _, n := range port.VirtualNetworks {
		vlans = append(vlans, n.GetId())
		vxlans = append(vxlans, int(n.GetVxlan()))
	}

	vlanIDs, diags := types.SetValueFrom(ctx, types.StringType, vlans)
	if diags != nil {
		return diags
	}

	m.VLANIDs = vlanIDs

	vxlanIDs, diags := types.SetValueFrom(ctx, types.Int32Type, vxlans)
	if diags != nil {
		return diags
	}

	m.VXLANIDs = vxlanIDs

	if port.Bond != nil {
		m.BondID = types.StringValue(port.Bond.GetId())
		m.BondName = types.StringValue(port.Bond.GetName())
	}

	return diags
}

type datasourceModel struct {
	PortID           types.String `tfsdk:"port_id"`
	DeviceID         types.String `tfsdk:"device_id"`
	Name             types.String `tfsdk:"name"`
	NetworkType      types.String `tfsdk:"network_type"`
	Type             types.String `tfsdk:"type"`
	MAC              types.String `tfsdk:"mac"`
	BondID           types.String `tfsdk:"bond_id"`
	BondName         types.String `tfsdk:"bond_name"`
	Bonded           types.Bool   `tfsdk:"bonded"`
	DisbondSupported types.Bool   `tfsdk:"disbond_supported"`
	NativeVlanID     types.String `tfsdk:"native_vlan_id"`
	VLANIDs          types.Set    `tfsdk:"vlan_ids"`
	VXLANIDs         types.Set    `tfsdk:"vxlan_ids"`
	Layer2           types.Bool   `tfsdk:"layer2"`
}

func (m *datasourceModel) parse(ctx context.Context, port *metalv1.Port) diag.Diagnostics {
	var diags diag.Diagnostics

	m.PortID = types.StringValue(port.GetId())
	m.Type = types.StringValue(string(port.GetType()))
	m.Name = types.StringValue(port.GetName())
	m.NetworkType = types.StringValue(string(port.GetNetworkType()))
	m.MAC = types.StringValue(port.Data.GetMac())
	m.Bonded = types.BoolValue(port.Data.GetBonded())
	m.DisbondSupported = types.BoolValue(port.GetDisbondOperationSupported())

	l2 := slices.Contains(l2Types, port.GetNetworkType())
	l3 := slices.Contains(l3Types, port.GetNetworkType())

	if l2 {
		m.Layer2 = types.BoolValue(true)
	}
	if l3 {
		m.Layer2 = types.BoolValue(false)
	}

	if port.NativeVirtualNetwork != nil {
		m.NativeVlanID = types.StringValue(port.NativeVirtualNetwork.GetId())
	}

	vlans := []string{}
	vxlans := []int{}

	for _, n := range port.VirtualNetworks {
		vlans = append(vlans, n.GetId())
		vxlans = append(vxlans, int(n.GetVxlan()))
	}

	vlanIDs, diags := types.SetValueFrom(ctx, types.StringType, vlans)
	if diags != nil {
		return diags
	}

	m.VLANIDs = vlanIDs

	vxlanIDs, diags := types.SetValueFrom(ctx, types.Int32Type, vxlans)
	if diags != nil {
		return diags
	}

	m.VXLANIDs = vxlanIDs

	if port.Bond != nil {
		m.BondID = types.StringValue(port.Bond.GetId())
		m.BondName = types.StringValue(port.Bond.GetName())
	}

	return diags
}
