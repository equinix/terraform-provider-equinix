package metal_gateway

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

var subnetSizes = []int64{8, 16, 32, 64, 128}

func metalGatewayResourceSchema(ctx context.Context) *schema.Schema {
    return &schema.Schema{
        Attributes: map[string]schema.Attribute{
            "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
                Delete: true,
            }),
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
                // NOTE (ocobles)
                //DiffSuppressFunc does not exist in fw, but I think it would not be necessary anyway and with computed in conflict it should work fine
                //DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
                // Suppress diff of IP reservation ID if private_ipv4_subnet_size has been set.
                // When the subnet size is set, the API will create a private subnet and return its ID
                // in this field, which generates a diff (ip_reservation_id is unset in HCL,
                // but the refreshed state shows there's an UUID of the new IPv4 block).
                    // 	if d.Get("private_ipv4_subnet_size").(int) != 0 {
                    // 		return true
                    // 	}
                    // 	return false
                    // },
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
