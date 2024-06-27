package gateway

import (
	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (m *ResourceModel) parse(gw *metalv1.FindMetalGatewayById200Response) diag.Diagnostics {
	// Convert Metal Gateway data to the Terraform state
	if gw.MetalGateway != nil {
		m.ID = types.StringValue(gw.MetalGateway.GetId())
		m.ProjectID = types.StringValue(gw.MetalGateway.Project.GetId())
		m.VlanID = types.StringValue(gw.MetalGateway.VirtualNetwork.GetId())
		m.VrfID = types.StringNull()

		if gw.MetalGateway.IpReservation != nil {
			m.IPReservationID = types.StringValue(gw.MetalGateway.IpReservation.GetId())
		} else {
			m.IPReservationID = types.StringNull()
		}

		m.PrivateIPv4SubnetSize = calculateSubnetSize(gw.MetalGateway.IpReservation)
		m.State = types.StringValue(string(gw.MetalGateway.GetState()))
	} else {
		m.ID = types.StringValue(gw.VrfMetalGateway.GetId())
		m.ProjectID = types.StringValue(gw.VrfMetalGateway.Project.GetId())
		m.VlanID = types.StringValue(gw.VrfMetalGateway.VirtualNetwork.GetId())
		m.VrfID = types.StringValue(gw.VrfMetalGateway.Vrf.GetId())

		if gw.VrfMetalGateway.IpReservation != nil {
			m.IPReservationID = types.StringValue(gw.VrfMetalGateway.IpReservation.GetId())
		} else {
			m.IPReservationID = types.StringNull()
		}

		m.PrivateIPv4SubnetSize = calculateSubnetSize(gw.VrfMetalGateway.IpReservation)
		m.State = types.StringValue(string(gw.VrfMetalGateway.GetState()))
	}
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

func (m *DataSourceModel) parse(gw *metalv1.FindMetalGatewayById200Response) diag.Diagnostics {
	if gw.MetalGateway != nil {
		// Convert Metal Gateway data to the Terraform state
		m.ID = types.StringValue(gw.MetalGateway.GetId())
		m.ProjectID = types.StringValue(gw.MetalGateway.Project.GetId())
		m.VlanID = types.StringValue(gw.MetalGateway.VirtualNetwork.GetId())
		m.VrfID = types.StringNull()

		if gw.MetalGateway.IpReservation != nil {
			m.IPReservationID = types.StringValue(gw.MetalGateway.IpReservation.GetId())
		} else {
			m.IPReservationID = types.StringNull()
		}

		m.PrivateIPv4SubnetSize = calculateSubnetSize(gw.MetalGateway.IpReservation)
		m.State = types.StringValue(string(gw.MetalGateway.GetState()))
	} else {
		// Convert Metal Gateway data to the Terraform state
		m.ID = types.StringValue(gw.VrfMetalGateway.GetId())
		m.ProjectID = types.StringValue(gw.VrfMetalGateway.Project.GetId())
		m.VlanID = types.StringValue(gw.VrfMetalGateway.VirtualNetwork.GetId())
		m.VrfID = types.StringValue(gw.VrfMetalGateway.Vrf.GetId())

		if gw.VrfMetalGateway.IpReservation != nil {
			m.IPReservationID = types.StringValue(gw.VrfMetalGateway.IpReservation.GetId())
		} else {
			m.IPReservationID = types.StringNull()
		}

		m.PrivateIPv4SubnetSize = calculateSubnetSize(gw.VrfMetalGateway.IpReservation)
		m.State = types.StringValue(string(gw.VrfMetalGateway.GetState()))
	}
	return nil
}

type ipReservationCommon interface {
	GetCidr() int32
	GetPublic() bool
	GetAddressFamily() int32
}

func calculateSubnetSize(ip ipReservationCommon) basetypes.Int64Value {
	privateIPv4SubnetSize := uint64(0)
	if !ip.GetPublic() && ip.GetAddressFamily() == 4 {
		privateIPv4SubnetSize = 1 << (32 - ip.GetCidr())
		return types.Int64Value(int64(privateIPv4SubnetSize))
	}
	return types.Int64Null()
}
