package metal_reserved_ip_block

import (
	"fmt"
    "path"
	"encoding/json"
    "context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type MetalReservedIPBlockResourceModel struct {
    ID                types.String `tfsdk:"id"`
    ProjectID         types.String `tfsdk:"project_id"`
    Facility          types.String `tfsdk:"facility"`
    Metro             types.String `tfsdk:"metro"`
    Description       types.String `tfsdk:"description"`
    Quantity          types.Int64  `tfsdk:"quantity"`
    Type              types.String `tfsdk:"type"`
    CidrNotation      types.String `tfsdk:"cidr_notation"`
    Tags              types.List   `tfsdk:"tags"`
    CustomData        types.String `tfsdk:"custom_data"`
    WaitForState      types.String `tfsdk:"wait_for_state"`
    VrfID             types.String `tfsdk:"vrf_id"`
    Network           types.String `tfsdk:"network"`
    Cidr              types.Int64  `tfsdk:"cidr"`
    Address           types.String `tfsdk:"address"`
    AddressFamily     types.Int64  `tfsdk:"address_family"`
    Gateway           types.String `tfsdk:"gateway"`
    Netmask           types.String `tfsdk:"netmask"`
    Manageable        types.Bool   `tfsdk:"manageable"`
    Management        types.Bool   `tfsdk:"management"`
    Global            types.Bool   `tfsdk:"global"`
    Public            types.Bool   `tfsdk:"public"`
    Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func (rm *MetalReservedIPBlockResourceModel) parse(ctx context.Context, reservedBlock *packngo.IPAddressReservation) diag.Diagnostics {
    var diags diag.Diagnostics

    rm.ID            = types.StringValue(reservedBlock.ID)
    rm.ProjectID     = types.StringValue(path.Base(reservedBlock.Project.Href))
    rm.Address       = types.StringValue(reservedBlock.Address)
    rm.AddressFamily = types.Int64Value(int64(reservedBlock.AddressFamily))
    rm.Cidr          = types.Int64Value(int64(reservedBlock.CIDR))
    rm.Gateway       = types.StringValue(reservedBlock.Gateway)
    rm.Network       = types.StringValue(reservedBlock.Network)
    rm.Netmask       = types.StringValue(reservedBlock.Netmask)
    rm.Public        = types.BoolValue(reservedBlock.Public)
    rm.Management    = types.BoolValue(reservedBlock.Management)
    rm.Manageable    = types.BoolValue(reservedBlock.Manageable)
    rm.Type          = types.StringValue(string(reservedBlock.Type))

    // Optional fields
    if reservedBlock.Facility != nil {
        rm.Facility = types.StringValue(reservedBlock.Facility.Code)
    }
    if reservedBlock.Metro != nil {
        rm.Metro = types.StringValue(reservedBlock.Metro.Code)
    }
    if reservedBlock.VRF != nil {
        rm.VrfID = types.StringValue(reservedBlock.VRF.ID)
    }
    if reservedBlock.Description != nil && (*(reservedBlock.Description) != "") {
        rm.Description = types.StringPointerValue(reservedBlock.Description)
    }

    // Handling tags as a list
    tags, diags := types.ListValueFrom(ctx, types.StringType, reservedBlock.Tags)
    if diags.HasError() {
        return diags
    }
    rm.Tags = tags

    // Custom data (assuming it's a JSON string)
    if reservedBlock.CustomData != nil {
        customDataJSON, err := json.Marshal(reservedBlock.CustomData)
        if err != nil {
            diags.AddError(
                "Error parsing Reserved IP Block",
                fmt.Sprintf("Error marshaling custom data to JSON: %s", err.Error()),
            )
            return diags
        } else {
            rm.CustomData = types.StringValue(string(customDataJSON))
        }
    }

    // Description 
    if reservedBlock.Description != nil {
        rm.Description = types.StringPointerValue(reservedBlock.Description)
    }

    // CIDR notation
    rm.CidrNotation = types.StringValue(fmt.Sprintf("%s/%d", reservedBlock.Network, reservedBlock.CIDR))

    quantity := 0
	if reservedBlock.AddressFamily == 4 {
		quantity = 1 << (32 - reservedBlock.CIDR)
	} else {
		// In Equinix Metal, a reserved IPv6 block is allocated when a device is
		// run in a project. It's always /56, and it can't be created with
		// Terraform, only imported. The longest assignable prefix is /64,
		// making it max 256 subnets per block. The following logic will hold as
		// long as /64 is the smallest assignable subnet size.
		bits := 64 - reservedBlock.CIDR
		if bits > 30 {
            diags.AddError(
                "Error parsing Reserved IP Block",
                fmt.Sprintf("strange (too small) CIDR prefix: %d", reservedBlock.CIDR),
            )
            return diags
		}
		quantity = 1 << uint(bits)
	}
    rm.Quantity = types.Int64Value(int64(quantity))

    return diags
}
