package gateway

import (
	"github.com/equinix/terraform-provider-equinix/internal/deprecations"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var dataSourceSchema = schema.Schema{
	DeprecationMessage: deprecations.MetalDeprecationMessage,
	Attributes: map[string]schema.Attribute{
		"id": framework.IDAttributeDefaultDescription(),
		"gateway_id": schema.StringAttribute{
			Description: "UUID of the Metal Gateway to fetch",
			Required:    true,
		},
		"project_id": schema.StringAttribute{
			Description: "UUID of the Project where the Gateway is scoped to",
			Computed:    true,
		},
		"vlan_id": schema.StringAttribute{
			Description: "UUID of the associated VLAN",
			Computed:    true,
		},
		"vrf_id": schema.StringAttribute{
			Description: "UUID of the VRF associated with the IP Reservation",
			Computed:    true,
		},
		"ip_reservation_id": schema.StringAttribute{
			Description: "UUID of the associated IP Reservation",
			Computed:    true,
		},
		"private_ipv4_subnet_size": schema.Int64Attribute{
			Description: "Size of the private IPv4 subnet to create for this gateway",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "Status of the gateway resource",
			Computed:    true,
		},
	},
}
