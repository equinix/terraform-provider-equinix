package gateway

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

func (m *ResourceModel) parse(gw *packngo.MetalGateway) diag.Diagnostics {
    // Convert Metal Gateway data to the Terraform state
    m.ID = types.StringValue(gw.ID)
    m.ProjectID = types.StringValue(gw.Project.ID)
    m.VlanID = types.StringValue(gw.VirtualNetwork.ID)

    if gw.VRF != nil {
        m.VrfID = types.StringValue(gw.VRF.ID)
    } else {
        m.VrfID = types.StringNull()
    }

    if gw.IPReservation != nil {
        m.IPReservationID = types.StringValue(gw.IPReservation.ID)
    } else {
        m.IPReservationID = types.StringNull()
    }

    m.PrivateIPv4SubnetSize = calculateSubnetSize(gw.IPReservation)
    m.State = types.StringValue(string(gw.State))
    return nil
}

type DataSourceModel struct {
    ID                    types.String `tfsdk:"id"`
	GatewayID             types.String `tfsdk:"gateway_id"`
	ProjectID             types.String `tfsdk:"project_id"`
    VlanID                types.String `tfsdk:"vlan_id"`
    VrfID                 types.String `tfsdk:"vrf_id"`
    IPReservationID       types.String `tfsdk:"ip_reservation_id"`
    PrivateIPv4SubnetSize types.Int64  `tfsdk:"private_ipv4_subnet_size"`
    State                 types.String `tfsdk:"state"`
}

func (m *DataSourceModel) parse(gw *packngo.MetalGateway) diag.Diagnostics {

    // Convert Metal Gateway data to the Terraform state
    m.ID = types.StringValue(gw.ID)
    m.ProjectID = types.StringValue(gw.Project.ID)
    m.VlanID = types.StringValue(gw.VirtualNetwork.ID)

	if gw.VRF != nil {
        m.VrfID = types.StringValue(gw.VRF.ID)
    } else {
        m.VrfID = types.StringNull()
    }

    if gw.IPReservation != nil {
        m.IPReservationID = types.StringValue(gw.IPReservation.ID)
    } else {
        m.IPReservationID = types.StringNull()
    }

    m.PrivateIPv4SubnetSize = calculateSubnetSize(gw.IPReservation)
    m.State = types.StringValue(string(gw.State))
	return nil
}

func calculateSubnetSize(ip  *packngo.IPAddressReservation) basetypes.Int64Value {
    privateIPv4SubnetSize := uint64(0)
    if !ip.Public {
        privateIPv4SubnetSize = 1 << (32 - ip.CIDR)
        return types.Int64Value(int64(privateIPv4SubnetSize))
    }
    return types.Int64Null()
}