package metal_vrf

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var metalVrfResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            Computed:    true,
            Description: "The unique identifier of the VRF",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.UseStateForUnknown(),
            },
        },
        "name": schema.StringAttribute{
            Required:    true,
            Description: "User-supplied name of the VRF, unique to the project",
        },
        "metro": schema.StringAttribute{
            Required:    true,
            Description: "Metro Code",
            PlanModifiers: []planmodifier.String{ // NOTE (ocobles) it wasn't mark as required in legacy sdk schema but it cannot be updated
                stringplanmodifier.RequiresReplace(),
            },
        },
        "project_id": schema.StringAttribute{
            Required:    true,
            Description: "ID of the project where the connection is scoped to. Required with type \"shared\"",
            PlanModifiers: []planmodifier.String{ // NOTE (ocobles) it wasn't mark as required in legacy sdk schema but it cannot be updated
                stringplanmodifier.RequiresReplace(),
            },
        },
        "description": schema.StringAttribute{
            Required:    true,
            Description: "Description of the VRF",
        },
        "local_asn": schema.Int64Attribute{
            Optional:    true,
            Computed:    true,
            Description: "The 4-byte ASN set on the VRF",
        },
        "ip_ranges": schema.ListAttribute{
            Optional:    true,
            Description: "All IPv4 and IPv6 Ranges that will be available to BGP Peers.",
        },
        // TODO: created_by, created_at, updated_at, href
    },
}
