package port

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
	PortID           types.String `tfsdk:"port_id"`
	Bonded           types.Bool   `tfsdk:"bonded"`
	Layer2           types.Bool   `tfsdk:"layer2"`
	NativeVlanID     types.String `tfsdk:"native_vlan_id"`
	VXLANIDs         types.Set    `tfsdk:"vxlan_ids"`
	VLANIDs          types.Set    `tfsdk:"vlan_ids"`
	ResetOnDelete    types.Bool   `tfsdk:"reset_on_delete"`
	Name             types.String `tfsdk:"name"`
	NetworkType      types.String `tfsdk:"network_type"`
	DisbondSupported types.Bool   `tfsdk:"disbond_supported"`
	BondName         types.String `tfsdk:"bond_name"`
	BondID           types.String `tfsdk:"bond_id"`
	Type             types.String `tfsdk:"type"`
	MAC              types.String `tfsdk:"mac"`
}

 
