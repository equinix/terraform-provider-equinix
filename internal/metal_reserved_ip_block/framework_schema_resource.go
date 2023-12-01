package metal_reserved_ip_block

import (
    "context"

	customstringvalidator "github.com/equinix/terraform-provider-equinix/internal/schema_validators/strings"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

func metalReservedIpBlockResourceSchema(ctx context.Context) *schema.Schema {
    return &schema.Schema{
        Attributes: map[string]schema.Attribute{
            "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
                Delete: true,
            }),
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "The unique identifier of the reserved IP block",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                Required:    true,
                Description: "The metal project ID where to allocate the address block",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "facility": schema.StringAttribute{
                Optional:    true,
                Description: "Facility where to allocate the public IP address block",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplaceIfConfigured(),
                },
                Validators: []validator.String{
                    stringvalidator.ConflictsWith(path.Expressions{
                        path.MatchRoot("metro"),
                    }...),
                },
                // NOTE (ocobles)
                //DiffSuppressFunc does not exist in fw
                // Let's try with RequiresReplaceIfConfigured ans see if it works as expected
                // otherwise replace it with appropriate logic in the Update function
                //
                // DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
                //     // suppress diff when unsetting facility
                //     if len(old) > 0 && new == "" {
                //         return true
                //     }
                //     return old == new
                // },
            },
            "metro": schema.StringAttribute{
                Optional:    true,
                Description: "Metro where to allocate the public IP address block",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                // TODO (ocobles)
                //
                // DiffSuppressFunc: func(k, fromState, fromHCL string, d *schema.ResourceData) bool {
                //     _, facOk := d.GetOk("facility")
                //     // if facility is not in state, treat the diff normally, otherwise do following messy checks:
                //     if facOk {
                //         // If metro from HCL is specified, but not present in state, suppress the diff.
                //         // This is legacy, and I think it's here because of migration, so that old
                //         // facility reservations are not recreated when metro is specified ???)
                //         if fromHCL != "" && fromState == "" {
                //             return true
                //         }
                //         // If metro is present in state but not present in HCL, suppress the diff.
                //         // This is for "facility-specified" reservation blocks created after ~July 2021.
                //         // These blocks will have metro "computed" to the TF state, and we don't want to
                //         // emit a diff if the metro field is empty in HCL.
                //         if fromHCL == "" && fromState != "" {
                //             return true
                //         }
                //     }
                //     return fromState == fromHCL
                // },
                //
                // TODO (ocobles)
                //
                // StateFunc doesn't exist in terraform, it requires implementation of bespoke logic before storing state, for instance in resource Create method
                // StateFunc: toLower
            },
            "description": schema.StringAttribute{
                Optional:    true,
                Description: "Arbitrary description for the reserved IP block",
            },
            "quantity": schema.Int64Attribute{
                Optional:    true,
                Computed:    true,
                Description: "The number of allocated /32 addresses, a power of 2",
                Validators: []validator.Int64{
                    int64validator.ExactlyOneOf(path.Expressions{
                        path.MatchRoot("vrf_id"),
                    }...),
                },
            },
            "type": schema.StringAttribute{
                Optional:    true,
                Description: "Either global_ipv4, public_ipv4, or vrf. Defaults to public_ipv4.",
                Default:     stringdefault.StaticString("public_ipv4"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf(
                        "public_ipv4",
                        "global_ipv4",
                        "vrf",
                    ),
                },
            },
            "tags": schema.ListAttribute{
                Optional:    true,
                Description: "Tags attached to the reserved block",
                ElementType: types.StringType,
            },
            "custom_data": schema.StringAttribute{
                Optional:    true,
                Description: "Custom Data in JSON format assigned to the IP Reservation",
                Default:     stringdefault.StaticString("{}"),
                Validators: []validator.String{
                    // NOTE (ocobles) StringIsJSON doesn't exist in framework,
                    // This is a custom implementation I made and we need to ensure it is working as expected
                    customstringvalidator.StringIsJSON(),
                },
                // TODO (ocobles) https://discuss.hashicorp.com/t/diffsuppressfunc-alternative-in-terraform-framework/52578/4
                // DiffSuppressFunc: structure.SuppressJsonDiff,
            },
            "wait_for_state": schema.StringAttribute{
                Optional:    true,
                Description: "Wait for the IP reservation block to reach a desired state on resource creation",
                Default:     stringdefault.StaticString(string(packngo.IPReservationStateCreated)),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf(
                        string(packngo.IPReservationStateCreated),
                        string(packngo.IPReservationStatePending),
                    ),
                },
            },
            "vrf_id": schema.StringAttribute{
                Optional:    true,
                Description: "VRF ID for type=vrf reservations",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.AlsoRequires(path.Expressions{
                        path.MatchRoot("network"),
                        path.MatchRoot("cidr"),
                    }...),
                },
            },
            "network": schema.StringAttribute{
                Optional:    true,
                Computed:    true,
                Description: "An unreserved network address from an existing vrf ip_range",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "cidr": schema.Int64Attribute{
                Optional:    true,
                Computed:    true,
                Description: "The size of the network to reserve from an existing vrf ip_range",
                PlanModifiers: []planmodifier.Int64{
                    int64planmodifier.RequiresReplace(),
                },
            },
            "cidr_notation": schema.StringAttribute{
                Computed:    true,
                Description: "CIDR notation of the IP address",
            },
            "address": schema.StringAttribute{
                Computed:    true,
                Description: "The IP address",
            },
            "address_family": schema.Int64Attribute{
                Computed:    true,
                Description: "Address family as integer (4 or 6)",
            },
            "gateway": schema.StringAttribute{
                Computed:    true,
                Description: "The gateway IP address",
            },
            "netmask": schema.StringAttribute{
                Computed:    true,
                Description: "The mask in decimal notation",
            },
            "manageable": schema.BoolAttribute{
                Computed:    true,
                Description: "Flag indicating whether the IP block is manageable",
            },
            "management": schema.BoolAttribute{
                Computed:    true,
                Description: "Flag indicating whether the IP block is for management",
            },
            "public": schema.BoolAttribute{
                Computed:    true,
                Description: "Flag indicating whether IP block is addressable from the Internet",
            },
            "global": schema.BoolAttribute{
                Computed:    true,
                Description: "Flag indicating whether IP block is global (i.e., assignable in any location)",
            },
        },
    }
}
