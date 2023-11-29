package metal_ip_attachment

import (
	"fmt"
    "path"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
)

type MetalIPAttachmentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	DeviceID       types.String `tfsdk:"device_id"`
	CIDRNotation   types.String `tfsdk:"cidr_notation"`
	Address        types.String `tfsdk:"address"`
	Gateway        types.String `tfsdk:"gateway"`
	Network        types.String `tfsdk:"network"`
	Netmask        types.String `tfsdk:"netmask"`
	AddressFamily  types.Int64  `tfsdk:"address_family"`
	CIDR           types.Int64  `tfsdk:"cidr"`
	Public         types.Bool   `tfsdk:"public"`
	Global         types.Bool   `tfsdk:"global"`
	Manageable     types.Bool   `tfsdk:"manageable"`
	Management     types.Bool   `tfsdk:"management"`
	VrfID          types.String `tfsdk:"vrf_id"`
}

func (rm *MetalIPAttachmentResourceModel) parse(assignment *packngo.IPAddressAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	rm.ID = types.StringValue(assignment.ID)
	rm.DeviceID = types.StringValue(path.Base(assignment.AssignedTo.Href))
	rm.CIDRNotation = types.StringValue(fmt.Sprintf("%s/%d", assignment.Network, assignment.CIDR))
	rm.Address = types.StringValue(assignment.Address)
	rm.Gateway = types.StringValue(assignment.Gateway)
	rm.Network = types.StringValue(assignment.Network)
	rm.Netmask = types.StringValue(assignment.Netmask)
	rm.AddressFamily = types.Int64Value(int64(assignment.AddressFamily))
	rm.CIDR = types.Int64Value(int64(assignment.CIDR))
	rm.Public = types.BoolValue(assignment.Public)
	rm.Global = types.BoolValue(assignment.Global)
	rm.Manageable = types.BoolValue(assignment.Manageable)
	rm.Management = types.BoolValue(assignment.Management)

    if assignment.VRF != nil {
        rm.VrfID = types.StringValue(assignment.VRF.ID)
    }

	return diags
}
