package gateway

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type ResourceModel struct {
    ID                    types.String   `tfsdk:"id"`
    ProjectID             types.String   `tfsdk:"project_id"`
    VlanID                types.String   `tfsdk:"vlan_id"`
    VrfID                 types.String   `tfsdk:"vrf_id"`
    IPReservationID       types.String   `tfsdk:"ip_reservation_id"`
    PrivateIPv4SubnetSize types.Int64    `tfsdk:"private_ipv4_subnet_size"`
    State                 types.String   `tfsdk:"state"`
    Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

func (rm *ResourceModel) parse(mg *packngo.MetalGateway) diag.Diagnostics {
    var diags diag.Diagnostics

    // Convert Metal Gateway data to the Terraform state
    rm.ID = types.StringValue(mg.ID)
    rm.ProjectID = types.StringValue(mg.Project.ID)
    rm.VlanID = types.StringValue(mg.VirtualNetwork.ID)

    if mg.VRF != nil {
        rm.VrfID = types.StringValue(mg.VRF.ID)
    } else {
        rm.VrfID = types.StringNull()
    }

    if mg.IPReservation != nil {
        rm.IPReservationID = types.StringValue(mg.IPReservation.ID)
    } else {
        rm.IPReservationID = types.StringNull()
    }

    // Calculate subnet size if it's a private IPv4 subnet
    privateIPv4SubnetSize := uint64(0)
    if !mg.IPReservation.Public {
        privateIPv4SubnetSize = 1 << (32 - mg.IPReservation.CIDR)
        rm.PrivateIPv4SubnetSize = types.Int64Value(int64(privateIPv4SubnetSize))
    } else {
        rm.PrivateIPv4SubnetSize = types.Int64Null()
    }

    rm.State = types.StringValue(string(mg.State))
    return diags
}