package metal_virtual_circuit

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

func metalVirtualCircuitResourceSchema(ctx context.Context) *schema.Schema {
    return &schema.Schema{
        Attributes: map[string]schema.Attribute{
            "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
                Delete: true,
            }),
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "Unique identifier of the Virtual Circuit",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "connection_id": schema.StringAttribute{
                Required:    true,
                Description: "UUID of Connection where the VC is scoped to",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "project_id": schema.StringAttribute{
                Required:    true,
                Description: "UUID of the Project where the VC is scoped to",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "port_id": schema.StringAttribute{
                Required:    true,
                Description: "UUID of the Connection Port where the VC is scoped to",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "name": schema.StringAttribute{
                Optional:    true,
                Description: "Name of the Virtual Circuit resource",
            },
            "description": schema.StringAttribute{
                Optional:    true,
                Description: "Description of the Virtual Circuit resource",
            },
            "speed": schema.StringAttribute{
                Optional:    true,
                Computed:    true,
                Description: "Description of the Virtual Circuit speed",
                // TODO: implement logic similar to SuppressDiffFunc for input with units to bps without units
            },
            "tags": schema.ListAttribute{
                Optional:    true,
                Description: "Tags attached to the reserved block",
                ElementType: types.StringType,
            },
            "nni_vlan": schema.Int64Attribute{
                Optional:    true,
                Description: "Equinix Metal network-to-network VLAN ID (optional when the connection has mode=tunnel)",
                PlanModifiers: []planmodifier.Int64{
                    int64planmodifier.RequiresReplace(),
                },
            },
            "vlan_id": schema.StringAttribute{
                Optional:    true,
                Description:  "UUID of the VLAN to associate",
                Validators: []validator.String{
                    stringvalidator.ExactlyOneOf(path.Expressions{
                        path.MatchRoot("vrf_id"),
                    }...),
                },
            },
            "vrf_id": schema.StringAttribute{
                Optional:    true,
                Description:  "UUID of the VRF to associate",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.AlsoRequires(path.Expressions{
                        path.MatchRoot("peer_asn"),
                        path.MatchRoot("subnet"),
                        path.MatchRoot("metal_ip"),
                        path.MatchRoot("customer_ip"),
                    }...),
                },
            },
            "peer_asn": schema.Int64Attribute{
                Optional:    true,
                Description:  "The BGP ASN of the peer. The same ASN may be the used across several VCs, but it cannot be the same as the local_asn of the VRF.",
                PlanModifiers: []planmodifier.Int64{
                    int64planmodifier.RequiresReplace(),
                },
            },
            "subnet": schema.StringAttribute{
                Optional:    true,
                Description: `A subnet from one of the IP blocks associated with the VRF that we will help create an IP reservation for. Can only be either a /30 or /31.
                    * For a /31 block, it will only have two IP addresses, which will be used for the metal_ip and customer_ip.
                    * For a /30 block, it will have four IP addresses, but the first and last IP addresses are not usable. We will default to the first usable IP address for the metal_ip.`,
            },
            "metal_ip": schema.StringAttribute{
                Optional:    true,
                Description: "The Metal IP address for the SVI (Switch Virtual Interface) of the VirtualCircuit. Will default to the first usable IP in the subnet.",
            },
            "customer_ip": schema.StringAttribute{
                Optional:    true,
                Description: "The Customer IP address which the CSR switch will peer with. Will default to the other usable IP in the subnet.",
            },
            "md5": schema.StringAttribute{
                Optional:    true,
                Sensitive:   true,
                Description: "The password that can be set for the VRF BGP peer",
            },
            "vnid": schema.Int64Attribute{
                Computed:    true,
                Description: "VNID VLAN parameter, see https://metal.equinix.com/developers/docs/networking/fabric/",
            },
            "nni_vnid": schema.Int64Attribute{
                Computed:    true,
                Description: "Nni VLAN ID parameter, see https://metal.equinix.com/developers/docs/networking/fabric/",
            },
            "status": schema.StringAttribute{
                Computed:    true,
                Description: "Status of the virtual circuit resource",
            },
        },
    }
}
