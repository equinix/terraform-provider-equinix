package metal_ip_attachment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var metalIPAttachmentResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            Computed:    true,
            Description: "The unique identifier of the IP attachment",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.UseStateForUnknown(),
            },
        },
        "device_id": schema.StringAttribute{
            Required:    true,
            Description: "UUID of the device to which the IP is attached",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
        },
        "cidr_notation": schema.StringAttribute{
            Required:    true,
            Description: "CIDR notation of the IP address",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
        },
        //TODO (ocobles) is not described in the legacy sdk documentation
        "address": schema.StringAttribute{
            Computed:    true,
            Description: "The IP address",
        },
        "gateway": schema.StringAttribute{
            Computed:    true,
            Description: "The gateway IP address",
        },
        "network": schema.StringAttribute{
            Computed:    true,
            Description: "The network IP address portion of the block specification",
        },
        "netmask": schema.StringAttribute{
            Computed:    true,
            Description: "The mask in decimal notation",
        },
        "address_family": schema.Int64Attribute{
            Computed:    true,
            Description: "Address family as integer (4 or 6)",
        },
        "cidr": schema.Int64Attribute{
            Computed:    true,
            Description: "Length of CIDR prefix of the block as integer",
        },
        "public": schema.BoolAttribute{
            Computed:    true,
            Description: "Flag indicating whether IP block is addressable from the Internet",
        },
        //TODO (ocobles) is not described in the legacy sdk documentation
        "global": schema.BoolAttribute{
            Computed:    true,
            Description: "Flag indicating whether IP block is global (i.e., assignable in any location)",
        },
        //TODO (ocobles) is not described in the legacy sdk documentation
        "manageable": schema.BoolAttribute{
            Computed:    true,
            Description: "Flag indicating whether the IP block is manageable",
        },
        //TODO (ocobles) is not described in the legacy sdk documentation
        "management": schema.BoolAttribute{
            Computed:    true,
            Description: "Flag indicating whether the IP block is for management",
        },
        //TODO (ocobles) it wasn't returned in the legacy sdk resource
        "vrf_id": schema.StringAttribute{
            Computed:    true,
            Description: "UUID of the VRF associated with the IP Reservation",
        },
    },
}
