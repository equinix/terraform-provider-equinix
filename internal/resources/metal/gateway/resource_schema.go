package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var subnetSizes = []int64{8, 16, 32, 64, 128}

func resourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this Metal Gateway",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "UUID of the Project where the Gateway is scoped to",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vlan_id": schema.StringAttribute{
				Description: "UUID of the VLAN to associate",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vrf_id": schema.StringAttribute{
				Description: "UUID of the VRF associated with the IP Reservation",
				Computed:    true,
			},
			"ip_reservation_id": schema.StringAttribute{
				Description: "UUID of the Public or VRF IP Reservation to associate",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("private_ipv4_subnet_size"),
					}...),
				},
			},
			"private_ipv4_subnet_size": schema.Int64Attribute{
				Description: "Size of the private IPv4 subnet to create for this gateway",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.OneOf(subnetSizes...),
					int64validator.ConflictsWith(path.Expressions{
						path.MatchRoot("ip_reservation_id"),
						path.MatchRoot("vrf_id"),
					}...),
				},
			},
			"state": schema.StringAttribute{
				Description: "Status of the gateway resource",
				Computed:    true,
			},
		},
	}
}
